package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/intercept"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/harness"
	"github.com/Autumn-27/norma/llm"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// Worker is an LLM work agent (docs §4.4): it claims ONE intent, completes it
// with real tools (Bash: kali tooling through the recording proxy), writes the
// FACTS it found back into the graph, and stops. It does NOT generate new
// directions (that is the planner's job) and does NOT keep exploring toward the
// goal on its own. Multiple workers run concurrently as goroutines.
// WebSearchOpts is the web-search backend selection the server pushes into each
// agent (planner/worker/main). Enabled=false leaves the web_search tool off.
// Backend is "ddgs" (no key), "brave-free" (BraveKey required), "tavily"
// (TavilyKey required), or "deepseek" (DeepSeek* required, filled from the
// active LLM profile). It maps directly onto agentcore.Options.
// Proxy is a dedicated egress proxy for the search request (http/https/socks5),
// independent of the traffic-recording MITM proxy — set it when the search endpoint
// is only reachable via a VPN/SOCKS proxy. Empty = direct.
//
// 주의: deepseek 백엔드는 다른 셋과 성질이 다르다: DeepSeek에는 직접 호출 가능한 검색 엔드포인트가 없고,
// 검색은 그 Anthropic 호환 messages 엔드포인트 내부에만 존재한다(web_search_20250305 server
// tool). 따라서 매 검색이 모델 호출을 한 번 소비하고, 검색 요청은 DeepSeek 서버가 보낸다——
// 로컬 Proxy를 거치지 않고, 트래픽 기록에도 남지 않는다.
type WebSearchOpts struct {
	Enabled   bool
	Backend   string
	BraveKey  string
	TavilyKey string
	Proxy     string
	// DeepSeek*는 현재 활성 LLM 설정에서 온다(anthropic 형식의 DeepSeek 공식 엔드포인트만),
	// 별도로 설정하지 않으며 LLM 설정 전환에 따라 바뀐다.
	DeepSeekBaseURL string
	DeepSeekAPIKey  string
	DeepSeekModel   string
}

type Worker struct {
	findingRecorder FindingRecorder
	prov            llm.Provider
	model           string
	workDir         string
	proxyAddr       string
	proxyCACert     string            // recording proxy's CA cert path (for WebFetch HTTPS verify)
	webSearch       WebSearchOpts     // web_search tool backend selection (off by default)
	tx              *transcript.Store // raw LLM conversation persistence (nil = off)
	window          int               // context window in tokens (for compaction)
	windowFn        func() int        // optional dynamic task-chain minimum
	maxTurns        int               // max agent turns per run (0 = unlimited)
	// runTimeout is the wall-clock budget for the main exploration of one intent
	// (0 = unlimited). When it fires, the run is cut and a settlement round is
	// forced so already-identified facts get written back instead of being lost.
	runTimeout time.Duration
	// extraTools are host-provided tools (e.g. traffic query, oast) appended to
	// the worker's graph write-back tools.
	extraTools []actool.CoreTool
	// injectConstraints resolves whether this task's operation constraints get
	// injected into the worker system prompt. Read per run so the settings toggle
	// takes effect without rebuilding the agent. nil = inject (default).
	injectConstraints func() bool
	// nonStreamingFn resolves whether this run uses the non-streaming (Complete)
	// path. Read per run so a profile/task toggle takes effect without rebuilding
	// the agent. nil = streaming (default).
	nonStreamingFn func() bool
	// noaEnabledFn resolves whether this run uses the experimental noa context-
	// compression mechanism. Read per run, like nonStreaming. nil = off (built-in
	// compaction).
	noaEnabledFn func() bool
	// maxTokensFn resolves the per-reply output cap in tokens, on the same
	// per-run basis. nil or 0 = send no cap and let the endpoint decide.
	maxTokensFn func() int
}

// WorkerSessionID returns the stable transcript key used by a worker intent.
// Worker slots are reusable, so the intent id (rather than work#N) is the
// session identity. Keep this helper public so the Worker message API and UI
// can refer to exactly the conversation that will be resumed.
func WorkerSessionID(explorationID, intentID int64) string {
	return fmt.Sprintf("exp%d-worker-i%d", explorationID, intentID)
}

const workerChatMarkerPrefix = "<!-- ARTEX_WORKER_CHAT:"

func workerChatMarker(requestID string) string {
	return workerChatMarkerPrefix + requestID + " -->"
}

func hasWorkerChatMessage(messages []llm.Message, requestID string) bool {
	marker := workerChatMarker(requestID)
	for _, message := range messages {
		if message.Role == llm.RoleUser && strings.Contains(message.Text(), marker) {
			return true
		}
	}
	return false
}

// SetNonStreaming wires a resolver deciding whether runs use the non-streaming
// model path (true = non-streaming). nil/unset = streaming (default). Read per
// run so a profile or task-chain toggle takes effect without rebuilding.
func (w *Worker) SetNonStreaming(fn func() bool) { w.nonStreamingFn = fn }

func (w *Worker) nonStreaming() bool { return w.nonStreamingFn != nil && w.nonStreamingFn() }

// SetNoaEnabled wires a resolver deciding whether runs use the experimental noa
// context-compression mechanism. nil/unset = off (built-in compaction). Read per
// run so the settings toggle takes effect without rebuilding the agent.
func (w *Worker) SetNoaEnabled(fn func() bool) { w.noaEnabledFn = fn }

// SetMaxTokens wires a resolver for the per-reply output cap. nil/unset or 0 =
// send no cap and let the endpoint decide. Read per run, like nonStreaming.
func (w *Worker) SetMaxTokens(fn func() int) { w.maxTokensFn = fn }

func (w *Worker) maxTokens() int {
	if w.maxTokensFn == nil {
		return 0
	}
	return w.maxTokensFn()
}

// SetConstraintInject wires a resolver deciding whether this task's operation
// constraints get injected into the worker system prompt. nil = inject (default).
func (w *Worker) SetConstraintInject(fn func() bool) { w.injectConstraints = fn }

// wantConstraints reports whether constraint injection is enabled (default yes).
func (w *Worker) wantConstraints() bool { return w.injectConstraints == nil || w.injectConstraints() }

// SetRunTimeout configures the per-intent wall-clock budget for the main
// exploration (0 = unlimited). When it fires, the SDK settlement phase still runs
// so facts are never lost to a timeout. Safe to call before Execute.
func (w *Worker) SetRunTimeout(run time.Duration) {
	w.runTimeout = run
}

// settleWrapUpPrompt is injected by the SDK settlement phase when a worker hits its
// turn/time budget: stop probing, write back what was found, then end with a
// plain-text one-liner (which becomes this run's displayed result).
const settleWrapUpPrompt = "당신은 곧 예산 소진으로 종료된다. 더는 어떤 명령/탐지도 실행하지 마라. 순서대로: (1) 위에서 이미 식별했지만 아직 쓰지 않은 내용을 하나씩 써라——새 자산은 insert_assets, 탐색 결론/사실은 record_fact, 확인된 취약점은 report_finding; (2) **마지막에 한 문장 순수 텍스트로 단독으로** 당신이 무엇을 했고 어떤 핵심 결론을 얻었는지 요약하라(이 문장이 이번 실행의 결과로 표시되니 반드시 출력)."

func NewWorker(prov llm.Provider, model, workDir string, tx *transcript.Store, window, maxTurns int, extra ...actool.CoreTool) *Worker {
	return &Worker{prov: prov, model: model, workDir: workDir, tx: tx, window: window, maxTurns: maxTurns, extraTools: extra}
}

// defaultToolsExcept returns actool.DefaultTools() minus the named tools (by
// CoreTool.Name()). Used to trim SDK default tools an agent shouldn't have.
func defaultToolsExcept(exclude ...string) []actool.CoreTool {
	drop := make(map[string]bool, len(exclude))
	for _, n := range exclude {
		drop[n] = true
	}
	all := actool.DefaultTools()
	out := make([]actool.CoreTool, 0, len(all))
	for _, t := range all {
		if !drop[t.Name()] {
			out = append(out, t)
		}
	}
	return out
}

func (w *Worker) SetCompactionWindowResolver(fn func() int) { w.windowFn = fn }

func (w *Worker) compactionWindow() int {
	if w.windowFn != nil {
		return w.windowFn()
	}
	return w.window
}

// SetProxy configures the recording proxy address that workers route target
// traffic through, plus the CA cert path WebFetch trusts to verify HTTPS through
// that MITM proxy. Empty addr disables the hint.
func (w *Worker) SetProxy(addr, caCert string) { w.proxyAddr, w.proxyCACert = addr, caCert }

// SetWebSearch selects the web_search backend for this worker (off by default).
func (w *Worker) SetWebSearch(o WebSearchOpts) { w.webSearch = o }

// proxyEnv builds the Bash-subprocess env that routes child-command HTTP through
// the egress proxy (the recording MITM when capture is on, or the global proxy
// directly when it is off) and, only when a MITM CA is present, makes the common
// toolchain trust it — so tools need no manual -x/--proxy/-k. Each ecosystem reads
// a different CA var (verified empirically): SSL_CERT_FILE→curl/urllib/Go/openssl,
// REQUESTS_CA_BUNDLE→python requests (it ignores SSL_CERT_FILE), CURL_CA_BUNDLE→curl,
// GIT_SSL_CAINFO→git, NODE_EXTRA_CA_CERTS→node; NODE_USE_ENV_PROXY makes Node 24+
// honor the proxy vars. ALL_PROXY is set too so a socks5 egress proxy (which curl
// only reads from ALL_PROXY, not HTTP(S)_PROXY) works in the capture-off path.
// Empty proxyAddr → nil (direct, unchanged env).
func proxyEnv(proxyAddr, caCert string) []string {
	if proxyAddr == "" {
		return nil
	}
	env := []string{
		"HTTP_PROXY=" + proxyAddr, "HTTPS_PROXY=" + proxyAddr,
		"http_proxy=" + proxyAddr, "https_proxy=" + proxyAddr,
		"ALL_PROXY=" + proxyAddr, "all_proxy=" + proxyAddr, // socks5 egress: curl reads only this
		"NODE_USE_ENV_PROXY=1", // Node 24+: honor HTTP(S)_PROXY in built-in fetch/http
	}
	if caCert != "" {
		env = append(env,
			"SSL_CERT_FILE="+caCert,
			"CURL_CA_BUNDLE="+caCert,
			"REQUESTS_CA_BUNDLE="+caCert,
			"GIT_SSL_CAINFO="+caCert,
			"NODE_EXTRA_CA_CERTS="+caCert,
		)
	}
	return env
}

// workerDefaultTmpl is the built-in EDITABLE body (섹션 [A]) of the worker system
// prompt, seeded into agent_prompts. The trafficTool block and the 중간 산출물 출력 규약
// are NOT here — they are code-owned and appended by workerSystem after rendering
// (섹션 [B]/[C]), so editing the DB body can never drop them.
const workerDefaultTmpl = `당신은 사이버보안 플랫폼의 승인된 침투 테스트 시스템의 "worker"(work agent)다. [하나의 의도](한 문장 탐색 방향)를 받으며, 유일한 역할: **이 의도 하나를 완수하고, 발견을 지식 그래프에 써서 되돌린 뒤, 멈추고 반환한다.**

**경계(레드라인)**:
1. **당신이 받은 이 의도 하나만 하라**. **이 의도를 탐색하다 이 의도 밖의 깊이 팔 가치가 있는 단서를 포착하면**(오류가 흘린 경로, 다른 자산과 연동될 수 있는 지점, 또 다른 익스플로잇 체인의 입구로 의심되는 것), **fact의 summary에 한마디 적어 planner에게 넘겨라**.
2. 처음 막혔다고(payload 필터링 / 404 / 인젝션 무반향) 다 파낸 건 아니다——이 의도의 모든 우회 수단을 다 써본 뒤 결론을 내라;
3. 승인된 범위 내에서만 동작하라. 시스템 프롬프트 맨 위에 [동작 제약]이 붙어 있으면 그것이 최우선 레드라인이다: 명령/탐지를 실행하기 전마다 자가 점검하고, 위반이면 하지 마라(네가 받은 의도 안에 있더라도).

**발견하면서 바로 써서 되돌려라**(그래프에 써야 유효하고, 머릿속/글로만 있는 건 무효; 결과가 하나 나올 때마다 즉시 쓰고, 마지막까지 쌓아두다 스텝 소진으로 잃지 마라). 세 가지 쓰기, 그래프를 섞지 마라:
- **새 자산/리소스 → insert_assets(자산 그래프)**: 서브도메인 / service / endpoint / 핑거프린트 / 자격 증명 등 모든 자산 [자체]. **여기엔 자산만 등록한다; 탐색 결론/판단은 여기 쓰지 말고 record_fact를 쓴다.**
- **탐색 결론/사실 → record_fact(탐색 그래프, intent_id 전달)**: 모두 이것을 쓴다. **여러 관찰을 [하나의] 사실로 종합**(summary 한 문장 요약 + detail에 요약을 뒷받침하는 확장을, 실제 실행 과정에 근거해 작성), 속성마다 하나씩 만들지 말고 의도 하나에 보통 하나만 쓰며, 잘게 쪼개면 그래프가 무한 팽창한다——**기본은 하나만 쓰고, detail에 합칠 수 있는 건 다 합쳐라**; [서로 완전히 독립적이어서 병합 불가]한 결론이 확실할 때만 facts 배열로 나눠 쓰며, 이는 극소수 예외이지 상례가 아니다. **증분만 써라**: 이번에 [새로 얻은] 것만 기록하고, 기존 사실을 표현만 바꿔 다시 쓰지 마라(기존을 재확인할 뿐 새 내용이 없으면 쓰지 않는다). **실제로 본 것만 써라**: evidence(한 줄: 명령 + 가장 잘 증명하는 한두 줄 출력, 간결하게, 세부는 detail에), confidence 표기(observed=직접 봄 / inferred=현상으로 추론).
- **취약점 확인 → report_finding(탐색 그래프, PoC 포함, intent_id 전달)**: **이번에 실제로 트리거해 재현 가능한 증거(요청/응답 또는 명령 출력)를 얻었을 때만 쓴다**. "버전/핑거프린트가 CVE에 매칭됨" "파라미터가 인젝션 가능해 보임" "외부 취약점 DB/업데이트 로그/코드 diff로 추론"을 확인된 것으로 간주하기 엄금, CVE DB 조회나 패치 버전 비교로 실제 트리거를 대체하지도 마라. 트리거 못 하지만 의심되면 → record_fact로 inferred 사실 하나(의심점+왜 트리거 못 했는지)를 기록해 planner에게 넘기고, 억지로 finding으로 기록하지 마라.


이 의도를 완수한 뒤 한 문장으로 무엇을 했고 어떤 사실을 써서 되돌렸는지 요약하라.`

// workerTrafficBlock is 섹션 [B]: the traffic-tool note, code-injected only when
// traffic capture (recording) is on — i.e. the traffic_* tools actually exist.
// Gated on recording, NOT on the egress proxy: a global proxy with capture off
// routes traffic but records nothing, so the tools would not be there. Not stored,
// not editable.
func workerTrafficBlock(recording bool) string {
	if !recording {
		return ""
	}
	return "\n\n**트래픽 도구**:\n- traffic_search / traffic_get / traffic_blob: 응답을 되짚어보고 이미 접근한 리소스를 찾는다, **먼저 트래픽을 조회하고 같은 URL을 curl로 반복하지 마라**. traffic_search는 **반드시 host를 지정**, 기본은 극경량 인덱스 3건만 반환(id/method/url/status/resp_len, 응답 내용 없음), 더 필요하면 limit을 명시적으로 키운다; body_contains로 요청/응답 본문 전문 검색 가능(최소 3자, 부분 문자열과 중국어 지원, 예: 비밀번호/키/오류/내부망 주소 찾기); 특정 건의 원문을 보려면 traffic_get(id), 그중 초대형 본문은 @blob sha256:<hash>로 표시되며 traffic_blob(hash)로 나눠 전문을 가져온다."
}

// artifactSpec is 섹션 [C]: the code-owned, non-editable tail appended to every
// pentest agent's prompt — intermediate artifacts must land in the shared work
// dir, never /tmp. Guaranteed present regardless of how the DB body is edited.
func artifactSpec(dir string) string {
	return "\n\n**중간 산출물 출력 규약**: 스크립트·payload·캡처한 응답 본문·임시 데이터 등 모든 중간 산출물은 **일률적으로 이 작업 작업 디렉터리 " + dir + "에 쓴다**(상대 경로는 여기에 쓰이고, 이 절대 경로를 써도 된다)——**/tmp에 쓰지 말고, 다른 절대 경로를 쓰지 마라**."
}

// workerArtifactSpec is the worker's 섹션 [C]: its per-intent run dir is pre-created
// by the engine (ensureRunDir), so it just writes relative paths there — no manual
// mkdir, no cross-worker name collisions.
func workerArtifactSpec(runDir string) string {
	return "\n\n**중간 산출물 출력 규약**: 스크립트·payload·캡처한 응답 본문·임시 데이터 등 모든 중간 산출물은 **일률적으로 이번 의도의 전용 작업 디렉터리 " + runDir + "에 쓴다**(이미 자동 생성됨, 상대 경로로 바로 여기에 쓰면 되고 수동으로 디렉터리를 만들 필요 없음)——**/tmp에 쓰지 말고, 다른 절대 경로를 쓰지 마라**."
}

// ensureRunDir builds and creates an agent's working directory under base:
// <base>/tasks/<taskID> for planner/main; <base>/tasks/<taskID>/i<intentID> for a
// worker (intentID<=0 → task dir only). The "tasks/" segment groups per-task dirs
// symmetrically with the chat agent's "sessions/<sessionID>". Best-effort mkdir — on
// failure, writes fail the same way an unwritable CWD would.
func ensureRunDir(base string, taskID, intentID int64) string {
	dir := filepath.Join(base, "tasks", strconv.FormatInt(taskID, 10))
	if intentID > 0 {
		dir = filepath.Join(dir, "i"+strconv.FormatInt(intentID, 10))
	}
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// cmdOutDir is the SDK large-tool-output spill dir under an agent's run dir.
func cmdOutDir(dir string) string { return filepath.Join(dir, "cmd-output") }

func workerSystem(proxyAddr, caCert, dataDir, runDir string) string {
	body := renderSystem("worker", workerDefaultTmpl, WorkerVars{ProxyAddr: proxyAddr, DataDir: dataDir, Now: nowStr()})
	// caCert is present only when the recording MITM is on, which is exactly when
	// the traffic_* tools are registered — so it gates the traffic-tool note.
	// Optional finding guidance is added for every role after tool resolution.
	return body + workerTrafficBlock(caCert != "") + workerArtifactSpec(runDir)
}

// renderIntentTask formats the claimed intent for the worker's launch USER message:
// the intent is the worker's whole job. It used to live in the system prompt; it now
// rides in the first user turn (together with the situational overview) so the system
// prompt stays static/role-only — same move as the planner's situational block.
// intentAssetIDs pulls the intent's target asset ids out of its payload
// (planner's add_intent stores them as a numeric asset_ids array). nil on absence
// or malformed payload.
func intentAssetIDs(intent *db.Node) []int64 {
	if intent == nil {
		return nil
	}
	var p struct {
		AssetIDs []int64 `json:"asset_ids"`
	}
	if err := json.Unmarshal(intent.Payload, &p); err != nil {
		return nil
	}
	return p.AssetIDs
}

func renderIntentTask(intent *db.Node) string {
	return fmt.Sprintf("\n\n[당신이 받은 의도(이번 유일한 작업: 이 하나만 하고, 사실만 생성하고, 끝나면 멈춤)]:\n%s\n의도 id: %d(record_fact / report_finding으로 되돌릴 때 전달)", string(intent.Payload), intent.ID)
}

// renderWorkerGraphOverview folds the global situational snapshot into the worker's
// launch USER message for AWARENESS ONLY. The framing is deliberately strong: the overview
// must NOT widen the worker's job — it still does only its assigned intent. Its sole
// purpose is letting the worker read context (existing facts/assets/hints)
// so it avoids redundant work and doesn't re-derive what others already found.
func renderWorkerGraphOverview(data map[string]any) string {
	// coverage는 planner가 '어떤 유형이 덜 테스트됐나 / 범위를 넓힐까'를 판단하는 신호로, worker의 '받은
	// 그 의도만 하고 미커버 지점을 쫓지 마라'는 역할 경계와 상충 → worker 뷰에서 제거. data는 이번 worker
	// 전용 새 map이라 키 삭제가 planner에 영향 없음.
	delete(data, "coverage")
	b, err := json.Marshal(data)
	if err != nil {
		return "" // fall back silently: the worker just won't have the global context
	}
	return "\n\n[전역 탐색 상황(읽기 전용, 네 이 의도를 큰 그림에 놓고 보도록 돕는다)]:\n" +
		"아래는 작업 전체의 현재 탐색 개황이다. 용도는 둘: 하나는 남들이 이미 발견한 것을 알아 중복을 피하기; 둘은 네 이 의도를 탐색할 때 그것과 전역의 관계를 연상하기.\n" +
		"**발산은 좋은 것**: 이 의도를 탐색할 때 얼마든지 깊이 생각하고 많이 연상하라. 유일한 경계는——다른 의도를 실제로 실행하지는 마라(그건 다른 worker의 일이고 planner가 스케줄한다). 가치 있는 단서를 연상하면(자산 간 연동, 또 다른 익스플로잇 체인의 입구로 의심되는 것, 전역 차원의 의심점) **반드시 fact에 써서 planner에게 넘겨라**——이것은 네 중요한 산출물이지 있어도 그만인 게 아니다. 한 건 더 보고해 planner가 판단하게 할지언정, 혼자 삼키지 마라.\n" +
		string(b)
}

// Execute runs one intent. hooks (the per-task Guard) gates every tool call; may
// be nil. emit, if non-nil, receives one ActivityRecord per execution step.
// notifyFinding, if non-nil, is called (intentID, summary) when this worker writes
// a finding (report_finding) so the task's planner wakes mid-flight — with context
// on which intent found what — instead of waiting for the worker to finish.
// Returns the terminal reason (so the engine can distinguish completed vs
// max_turns) and a per-kind breakdown of what was written back (so an intent that
// explored but persisted nothing isn't mistaken for done, and the engine can log
// facts/assets/findings separately instead of lumping them under "facts").
func (w *Worker) Execute(ctx context.Context, name string, taskID int64, as *db.AssetStore, ts *db.ExplorationStore, intent *db.Node, hooks harness.HookRunner, emit func(db.Activity), enr EnrichTrigger, notifyFinding func(int64, string)) (harness.TerminalReason, WriteCounts, error) {
	return w.execute(ctx, name, taskID, as, ts, intent, hooks, emit, enr, notifyFinding, "", "")
}

// ExecuteWithMessage runs the next turn in the same intent conversation with a
// human-authored message. The HTTP handler does not edit the transcript;
// agentcore records the message as a normal user turn when this Worker starts.
// This keeps Worker continuation identical to the regular agent chat flow.
func (w *Worker) ExecuteWithMessage(ctx context.Context, name string, taskID int64, as *db.AssetStore, ts *db.ExplorationStore, intent *db.Node, hooks harness.HookRunner, emit func(db.Activity), enr EnrichTrigger, notifyFinding func(int64, string), requestID, message string) (harness.TerminalReason, WriteCounts, error) {
	return w.execute(ctx, name, taskID, as, ts, intent, hooks, emit, enr, notifyFinding, strings.TrimSpace(requestID), strings.TrimSpace(message))
}

func (w *Worker) execute(ctx context.Context, name string, taskID int64, as *db.AssetStore, ts *db.ExplorationStore, intent *db.Node, hooks harness.HookRunner, emit func(db.Activity), enr EnrichTrigger, notifyFinding func(int64, string), requestID, message string) (harness.TerminalReason, WriteCounts, error) {
	tsx := NewToolSet(ts, name)
	tsx.SetFindingRecorder(w.findingRecorder)
	tsx.SetTaskID(taskID)
	coverageEnabled := as == nil || as.CoverageEnabled(taskID)
	tsx.SetCoverageEnabled(coverageEnabled)
	if as != nil {
		tsx.SetAssetStore(as, as.Companies())
	}
	tsx.SetOwnerNode(intent.ID)         // assets this worker discovers anchor to its intent → visible to the task
	tsx.SetEnrich(enr)                  // async DNS/HTTP auto-completion for assets this worker writes
	tsx.SetNotifyFinding(notifyFinding) // report_finding이 DB에 저장될 때 즉석에서 planner를 깨운다, '어느 의도+finding'을 함께
	// base = built-in worker tools ∪ host tools (traffic) ∪ default tools (incl. Bash);
	// then augment with the agent's visible skills/MCP. During the SDK settlement
	// phase, Bash is hidden via Settlement.DisabledTools (no local gating needed).
	base := append(tsx.WorkerTools(), w.extraTools...)
	// worker에는 일부러 MultiEdit/Glob/Grep을 주지 않는다: 파일 정밀 수정은 Edit, 검색은 Bash(grep/find),
	// 도구 면을 좁혀 저가치 호출을 줄인다. 나머지 SDK 기본 도구(Read/Write/Edit/LS/Bash/Sleep)는 그대로.
	base = append(base, defaultToolsExcept("MultiEdit", "Glob", "Grep")...)
	ctx = WithRunInfo(ctx, RunInfo{TaskID: taskID, ExplorationID: explorationID(ts), IntentID: intent.ID})
	tools, def, cleanup := AugmentTools(ctx, "worker", base)
	defer cleanup()

	// 의도는 worker의 [유일한 역할, run 전체를 관통하는 불변량] → 시작 지시·의도에 앵커된 대상 자산
	// 원본 데이터와 함께 system prompt에 넣는다: system은 run마다 다시 조립돼 compaction에 압축되지 않으므로,
	// 긴 run에서도 의도가 항상 존재하고, 이어서 실행할 때도 transcript 이력이 그 첫 메시지를 보존했는지에 의존하지 않는다. 대가는 system에
	// per-intent 가변 데이터가 섞여 의도 간 캐시 재사용을 잃는 것; 이는 의도적 절충이다(의도 분실이 token 절약보다 훨씬 심각).
	// planner의 '상황 블록은 user turn에'와 갈라지는 건 의도적이다: planner 자신은 의도를 내는 쪽이라 단일 mandate가 없고,
	// worker는 있다. [전역 상황 overview]만 시작 user 메시지에 남긴다——이건 다운그레이드 가능하고 stale을 용인하며 압축돼도 무방하다.
	// 이번 의도의 전용 작업 디렉터리 <workDir>/tasks/<taskID>/i<intentID>, 엔진 측이 먼저 생성.
	runDir := ensureRunDir(w.workDir, taskID, intent.ID)
	// The run-wide intent is not the current tool action. Do not forward it or
	// inherit a parent run's background into the action reviewer.
	ctx = intercept.WithReviewContext(ctx, runDir, intercept.ReviewBackground{})
	overview := renderWorkerGraphOverview(tsx.graphOverviewData())
	sysBody := workerSystem(w.proxyAddr, w.proxyCACert, w.workDir, runDir)
	if w.wantConstraints() {
		sysBody += constraintBlock(ts) // 동작 제약(있으면)을 시스템 프롬프트에 주입, worker 실행 시 엄격히 준수
	}
	// 의도 블록 → 의도 앵커 자산 블록 → 시작 지시를 순서대로 system 끝에 추가(constraintBlock과 동일한 추가 방식).
	sysBody += renderIntentTask(intent)
	if as != nil {
		if ids := intentAssetIDs(intent); len(ids) > 0 {
			if assets, err := as.GetByIDs(ids); err == nil && len(assets) > 0 {
				if b, err := json.Marshal(assets); err == nil {
					sysBody += "\n\n이 의도의 asset_ids에 대응하는 대상 자산:\n" + string(b)
				}
				// 의도가 명확히 겨냥한 이 자산들 → 작업 테스트 범위에 자동 편입(insertAssets와 동일한
				// 보수적 입도). upsertTaskScope의 ON CONFLICT DO NOTHING + uq_task_scope
				// 유니크 인덱스가 중복 추가되지 않음을 보장; 재실행/재시도도 멱등 no-op.
				// 이 의도의 asset_ids 자동 편입은 커버리지가 켜졌을 때만 수행하며, insertAssets의 범위 누적은 스위치와 무관하다.
				if coverageEnabled {
					for _, a := range assets {
						_ = as.AddAutoScope(taskID, a.Type, a.Domain, a.URL, a.IP)
					}
				}
			}
		}
	}
	sysBody += "\n\n위 이 의도를 실행 시작: 그것만 하고, 사실·assets·finding만 생성하고, 끝나면 멈춤."
	system, boundary := deferredSystem(sysBody, def)
	// 작업 수준 deadline(ctx로 주입)이 이 run의 벽시계 예산을 좁히고 마무리 문구를 결정(taskclock.go 참고).
	tc := taskClockFrom(ctx)
	maxDur, clamped := clampMaxDuration(tc.DeadlineUnix, w.runTimeout)
	settle := wrapupSettlement("worker", []string{"Bash"})
	if tc.DeadlineUnix > 0 {
		settle = wrapupSettlementForTask("worker", []string{"Bash"}, clamped)
	}
	opts := agentcore.Options{
		Provider:        w.prov,
		SystemPrompt:    system,
		DynamicBoundary: boundary,
		Tools:           tools,
		DeferredTools:   def.Deferred,
		UnlockSet:       def.Unlock,
		PermissionMode:  permission.ModeBypass,
		// WebFetch는 레코딩 프록시를 거쳐 그 HTTP가 curl처럼 기록된다; 프록시 CA를 로드해 MITM이
		// 재서명한 HTTPS 인증서가 [정상적으로 검증 통과]되게 한다(검증을 끄는 게 아님). proxy 비면 직접 연결.
		EnableWebFetch: true,
		WebFetchProxy:  w.proxyAddr,
		WebFetchCACert: w.proxyCACert,
		// 인터넷 검색(선택). ddgs는 key 불필요; brave-free는 BraveKey 필요; tavily는 TavilyKey 필요.
		// WebSearchProxy는 독립 아웃바운드 프록시(http/https/socks5), 트래픽을 기록하는 MITM 프록시와 무관; 비면 직접 연결.
		EnableWebSearch:       w.webSearch.Enabled,
		WebSearchBackend:      w.webSearch.Backend,
		BraveSearchAPIKey:     w.webSearch.BraveKey,
		TavilySearchAPIKey:    w.webSearch.TavilyKey,
		DeepSeekSearchBaseURL: w.webSearch.DeepSeekBaseURL,
		DeepSeekSearchAPIKey:  w.webSearch.DeepSeekAPIKey,
		DeepSeekSearchModel:   w.webSearch.DeepSeekModel,
		WebSearchProxy:        w.webSearch.Proxy,
		// Bash 하위 명령의 HTTP는 기본적으로 레코딩 프록시 경유 + 그 CA 신뢰(도구에 -x/-k 불필요).
		BashEnv:    proxyEnv(w.proxyAddr, w.proxyCACert),
		WorkingDir: runDir,
		MaxTurns:   w.maxTurns, // 0 = unlimited (configurable in agent management)
		// 벽시계 예산, 라운드 경계에서 판정, 중간에 끊지 않음; 0 = 무제한. 작업 수준 deadline이 있으면 min(자체 예산,
		// deadline까지 남은 시간)으로 좁혀, 이 run이 작업 만료 시 자연스레 마무리에 들어가게 한다(taskclock.go 참고).
		MaxDuration: maxDur,
		// 예산(라운드 OR 시간) 도달→ SDK가 마무리 한 라운드 실행(Bash 숨김), 식별된 것을 써서 되돌려 미완 종료를 막음.
		// clamped(작업 deadline으로 좁혀짐) 시 PromptByReason 사용: 타임아웃 때문=작업 만료→작업 타임아웃 문구,
		// 스텝 때문=좁힌 창 안에서 스텝이 먼저 소진→per-run 문구로 회귀. clamped 아니면 순수 per-run 유지.
		Settlement: settle,
		// large tool output spills to cmd-output/ with a head + pointer (SDK tool.Capture);
		// full output preserved on disk. 잘림 상한은 SDK 기본값 사용(30000자).
		ToolOutputDir: cmdOutDir(runDir),
		Compaction:    compactionConfig(w.compactionWindow()), // long tool-heavy runs stay within the window
		Todos:         actool.NewTodoStore(),                  // 세션 수준 임시 할 일(TodoWrite), 순수 계획용, 종료 시 폐기
		NonStreaming:  w.nonStreaming(),                       // 이 profile이 비스트리밍을 선택하면 Provider.Complete 경유
		MaxTokens:     w.maxTokens(),                          // 0 = 상한 미전송, 서버 기본값으로 결정
	}
	if hooks != nil { // typed-nil guard: only set when concrete (avoids harness panic)
		opts.Hooks = hooks
	}
	if w.tx != nil { // persist raw LLM conversation; one file per worked intent
		opts.Transcript = w.tx
		opts.SessionID = WorkerSessionID(ts.ID(), intent.ID)
	}
	intentID := intent.ID
	emitWrap := func(r db.Activity) {
		if emit != nil {
			r.NodeID, r.Worker = &intentID, name
			emit(r)
		}
	}
	// 의도 / 시작 지시 / 의도 앵커 자산은 이미 system prompt로 내려간다(위 sysBody 조립 참고).
	// 이 시작 user 메시지는 [전역 상황 overview]만 담는다——다운그레이드 가능한 대국 파악 정보라 압축돼도 무방.
	// overview가 드물게 marshal 실패로 비면, 시작 문구 한 줄로 회귀해 첫 라운드에 빈 user 메시지가 나오지 않게 한다.
	input := overview
	if strings.TrimSpace(input) == "" {
		input = "system에서 받은 의도를 실행 시작: 그것만 하고, 사실·assets·finding만 생성하고, 끝나면 멈춤."
	}

	// 실험 기능: 켜면 noa가 컨텍스트 압축을 인계(아카이브는 <workDir>/noa/<SessionID> 아래에 모이며, 영속).
	noaSession := WorkerSessionID(ts.ID(), intent.ID)
	enableNoa(&opts, w.noaEnabledFn, w.workDir, noaSession, noaWarn(noaSession))
	ctx = attachSideCapture(ctx, &opts)
	s := agentcore.NewSession(opts)
	defer s.Close() // release the session's background-task manager (temp dir + processes)

	// Resume prior conversation if this intent was paused/blocked/exhausted and is
	// being re-run. The transcript ID is deterministic per intent, so if a prior
	// session exists the worker continues from where it left off instead of
	// restarting from scratch.
	alreadyRecorded := false
	if w.tx != nil {
		_ = s.Resume(opts.SessionID)
		alreadyRecorded = requestID != "" && hasWorkerChatMessage(s.Messages(), requestID)
		if len(s.Messages()) > 0 && message == "" {
			seedUnlockFromHistory(s.Messages(), def.UnlockSkill)
			input = "계속 실행."
		} else if len(s.Messages()) > 0 {
			seedUnlockFromHistory(s.Messages(), def.UnlockSkill)
		}
	}
	if message != "" {
		if alreadyRecorded {
			input = "지난번 수동 대화로 입력한 새 의도를 계속 실행. 이미 완료한 동작을 반복하지 마라."
		} else if len(s.Messages()) > 0 {
			input = workerChatMarker(requestID) + "\n[수동 대화로 입력한 새 의도]\n" + message +
				"\n\n이 수동 입력에 따라 즉시 실행하고, 완료 후 맥락에 따라 원래 작업을 계속할지 결정하라."
		} else {
			input += "\n\n" + workerChatMarker(requestID) + "\n[수동 대화로 입력한 새 의도]\n" + message +
				"\n\n이 수동 입력을 우선 실행하라."
		}
	}

	// Budgets + settlement are owned by the SDK (MaxTurns/MaxDuration + Settlement):
	// on hit it runs a wrap-up turn and finishes with ReasonMaxTurns/ReasonTimeout.
	// MaxDuration now interrupts an in-flight tool at the wall-clock deadline and
	// enters the wrap-up phase on the live ctx, so a run whose tool overran the budget
	// still settles (no external hard-timeout backstop needed). ctx itself carries only
	// pause / planner kill / shutdown, which the engine distinguishes and re-queues/stops.
	_, reason, err := captureRunSession(ctx, s, input, emitWrap)
	return reason, tsx.Writes(), err
}
