package agent

import (
	"context"
	"fmt"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/intercept"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// MainAgent is the thin human-interface orchestrator (docs §4.2 / §7). The human
// chats with it; it observes (read tools), and steers by injecting hints
// (→planner) or direct high-priority intents (→frontier). It does NOT run the
// autonomous intent-generation loop (that is the planner's job).
type MainAgent struct {
	findingRecorder FindingRecorder
	prov            llm.Provider
	model           string
	tx              *transcript.Store                      // raw LLM conversation persistence (nil = off)
	window          int                                    // context window in tokens (for compaction)
	windowFn        func() int                             // optional dynamic task-chain minimum
	maxTurns        int                                    // max agent turns per run (0 = unlimited)
	proxyAddr       string                                 // recording proxy for WebFetch (empty = direct)
	proxyCACert     string                                 // recording proxy's CA cert path (HTTPS verify)
	webSearch       WebSearchOpts                          // web_search tool backend selection (off by default)
	workDir         string                                 // shared work dir (surfaced in prompt as artifact-output target)
	steerWork       func(intentID int64, msg string) error // engine callback: steer a running work (nil = off)
	nonStreamingFn  func() bool                            // resolver: use non-streaming (Complete) path? (nil = streaming)
	noaEnabledFn    func() bool                            // resolver: use experimental noa compaction? (nil = off)
	maxTokensFn     func() int                             // resolver: per-reply output cap (nil/0 = send no cap)
}

// SetNoaEnabled wires a resolver deciding whether runs use the experimental noa
// context-compression mechanism. nil/unset = off (built-in compaction). Read per
// run so the settings toggle takes effect without rebuilding the agent.
func (m *MainAgent) SetNoaEnabled(fn func() bool) { m.noaEnabledFn = fn }

// SetNonStreaming wires a resolver deciding whether runs use the non-streaming
// model path (true = non-streaming). nil/unset = streaming (default).
func (m *MainAgent) SetNonStreaming(fn func() bool) { m.nonStreamingFn = fn }

func (m *MainAgent) nonStreaming() bool { return m.nonStreamingFn != nil && m.nonStreamingFn() }

// SetMaxTokens wires a resolver for the per-reply output cap. nil/unset or 0 =
// send no cap and let the endpoint decide. Read per run, like nonStreaming.
func (m *MainAgent) SetMaxTokens(fn func() int) { m.maxTokensFn = fn }

func (m *MainAgent) maxTokens() int {
	if m.maxTokensFn == nil {
		return 0
	}
	return m.maxTokensFn()
}

func NewMainAgent(prov llm.Provider, model, workDir string, tx *transcript.Store, window, maxTurns int) *MainAgent {
	return &MainAgent{prov: prov, model: model, workDir: workDir, tx: tx, window: window, maxTurns: maxTurns}
}

func (m *MainAgent) SetCompactionWindowResolver(fn func() int) { m.windowFn = fn }

func (m *MainAgent) compactionWindow() int {
	if m.windowFn != nil {
		return m.windowFn()
	}
	return m.window
}

// SetProxy points the main agent's WebFetch at the recording proxy plus the CA
// cert it trusts to verify HTTPS through it (empty addr = direct).
func (m *MainAgent) SetProxy(addr, caCert string) { m.proxyAddr, m.proxyCACert = addr, caCert }

// SetWebSearch selects the web_search backend for the main agent (off by default).
func (m *MainAgent) SetWebSearch(o WebSearchOpts) { m.webSearch = o }

// SetSteerWork wires the engine callback that lets the main agent's steer_work
// tool inject a mid-run course-correction into a running work (nil = tool off).
func (m *MainAgent) SetSteerWork(fn func(intentID int64, msg string) error) { m.steerWork = fn }

// mainAgentDefaultTmpl is the built-in EDITABLE body (섹션 [A]) of the main agent
// prompt, seeded into agent_prompts. Goal is a {{.Goal}} template var; the 중간
// 산출물 출력 규약 tail is code-owned (artifactSpec), appended after rendering.
const mainAgentDefaultTmpl = `당신은 승인된 침투 테스트 시스템의 "메인 agent"로, 인간 운영자의 인터페이스다. 당신은 직접 탐색하지 않으며, 스스로 의도를 연속 생성하지도 않는다(그것은 planner의 일이다). 당신의 역할:

1. 관찰: graph_overview / list_findings / list_facts / list_assets / get_worker_output 로 현재 진행 상황에 대한 사람의 질문에 답한다.
2. 조타(사람의 의도를 시스템에 반영):
   - 사람이 "방향 변경 / 특정 유형 취약점 강조 / 특정 영역 집중"을 원하면 → add_hint 로 힌트를 쓴다(planner가 다음번에 읽는다).
   - 사람이 "특정 목표를 즉시 테스트"하기를 원하면 → add_intent 로 고우선순위 의도 하나를 직접 주입한다(priority 8-10). 시스템이 자동으로 완료된 작업을 실행 상태로 되돌리고, worker가 이 의도를 받아 실행하며, 끝나면 다시 완료 상태로 돌아간다.
     **작업 목표가 이미 전부 달성된 경우**(graph_overview 에서 goals 가 모두 met): 내리기 전에 이 의도 뒤에 "새로, 달성해야 할 결과"가 암시되어 있는지 먼저 판단한다. 암시되어 있으면, 당신이 추측한 목표를 한 문장으로 사람에게 되풀이해 말하고, **정식 목표로 등록할지 되묻는다** —— 사람이 원하면 → set_goals 로 등록한다(작업은 이후 일반 계획에 들어가고, planner가 자율적으로 이어서 추진한다); 사람이 원치 않거나 / 그저 임시로 한번 탐색해 보려는 것이면 → add_intent 로 이 한 건만 내리고, worker가 작업을 다 실행하면 완료 상태로 돌아간다(자율적으로 계속하지 않는다). 이 의도가 명백히 일회성 확인일 뿐 새 목표를 암시하지 않으면, 바로 add_intent 하면 되며 매번 물을 필요는 없다.
   - 사람이 "실행 중인 어떤 의도(work)를 실시간 교정(더 이상 X로 가지 말고, Y에 집중)"하기를 원하면 → steer_work 를 쓴다(중단하지 않고, 기존 진행을 잃지 않으며, worker의 다음 동작 전에 적용된다); 먼저 get_worker_output 으로 무엇을 하고 있는지 본다. 방향이 통째로 틀렸으면 add_intent 로 새 의도를 따로 내린다.
   - 사람이 "달성해야 할 최종 목표를 새로 추가"하기를 원하면 → set_goals 로 목표를 보충한다. 시스템이 그 목표를 작업 그래프에 쓰고 **완료/일시정지된 작업을 자동으로 실행 상태로 되돌려 계속 실행한다**(planner가 이후 이를 근거로 달성 여부를 다시 판정한다). 사람이 다시 복구를 누를 필요가 없다.
   - 사람이 "테스트 제약을 추가/변경(특정 동작을 허용/금지, 예: '현재 포트만 테스트' '무차별 대입 금지' '패시브 정찰만')"하기를 원하면 → set_constraints 로 등록한다(type=allow 허용 / type=deny 금지). 제약은 다음 라운드 계획 시 planner/worker 의 프롬프트에 주입되어 탐색 경계를 설정한다; 개요의 '제약 관리'에서 추가·삭제·수정할 수도 있다.
3. 쉬운 말로 간결하게 답하고, 당신이 무엇을 했는지 설명한다.

현재 작업 목표: {{.Goal}}

발견을 지어내지 마라; 도구가 반환한 실제 데이터에만 근거해 답하라.`

func mainAgentSystem(goal, dataDir, workDir string) string {
	body := renderSystem("mainagent", mainAgentDefaultTmpl, MainVars{Goal: goal, DataDir: dataDir, Now: nowStr()})
	return body + artifactSpec(workDir)
}

// Chat handles one human message and returns the assistant reply. emit, if
// non-nil, receives each execution step (thinking / tool_use / tool_result /
// text / result) so the main-agent session shows its work — exactly like the
// worker/planner sessions — not just the final answer.
func (m *MainAgent) Chat(ctx context.Context, taskID int64, mainSeg int, as *db.AssetStore, ts *db.ExplorationStore, goal, message string, emit func(db.Activity), notify, resume func(), notifyGoal, notifyHint func([]string)) (string, error) {
	tsx := NewToolSet(ts, "human")
	tsx.SetFindingRecorder(m.findingRecorder)
	if as != nil {
		tsx.SetAssetStore(as, as.Companies())
	}
	tsx.SetTaskID(taskID)
	tsx.SetCoverageEnabled(as == nil || as.CoverageEnabled(taskID))
	tsx.SetNotify(notify)         // 범용 깨우기(전용 콜백이 없는 쓰기 작업은 이것을 사용, debounced)
	tsx.SetResumeTask(resume)     // set_goals 목표 추가 → 완료/일시정지된 작업을 running 으로 되돌린다
	tsx.SetNotifyGoal(notifyGoal) // set_goals 목표 추가 → planner 에 '사람이 목표를 추가함: …' 트리거 하나를 기록
	tsx.SetNotifyHint(notifyHint) // add_hint 힌트 추가 → planner 에 '사람이 전략 힌트 N개를 추가함: …' 트리거 하나를 기록
	tsx.steerWork = m.steerWork   // enable steer_work tool (nil = unavailable)
	// 도메인 도구 + 기본 디폴트 도구 집합(Read/Write/Edit/MultiEdit/LS/Glob/Grep/Bash)
	// 자산 커버리지 기능이 꺼지면 add_task_scope/list_untested_assets 를 제거한다(prompt에 들어가지 않음).
	base := append(tsx.DropCoverageTools(tsx.MainAgentTools()), actool.DefaultTools()...)
	ctx = WithRunInfo(ctx, RunInfo{TaskID: taskID, ExplorationID: explorationID(ts)})
	tools, def, cleanup := AugmentTools(ctx, "mainagent", base)
	defer cleanup()
	// 본 작업의 작업 디렉터리 <workDir>/tasks/<taskID>, 먼저 만들어 둔다.
	mainDir := ensureRunDir(m.workDir, taskID, 0)
	ctx = intercept.WithReviewWorkingDirectory(ctx, mainDir)
	system, boundary := deferredSystem(mainAgentSystem(goal, m.workDir, mainDir), def)
	opts := agentcore.Options{
		Provider:        m.prov,
		SystemPrompt:    system,
		DynamicBoundary: boundary,
		Tools:           tools,
		DeferredTools:   def.Deferred,
		UnlockSet:       def.Unlock,
		PermissionMode:  permission.ModeBypass,
		EnableWebFetch:  true, // 기록 프록시를 거쳐 흔적을 남긴다; 프록시 CA를 로드해 MITM이 재서명한 HTTPS 인증서를 검증
		WebFetchProxy:   m.proxyAddr,
		WebFetchCACert:  m.proxyCACert,
		// 인터넷 검색(선택). ddgs 는 key 불필요; brave-free 는 BraveKey 필요; tavily 는 TavilyKey 필요.
		// WebSearchProxy 는 독립 출구 프록시(http/https/socks5)로, 트래픽을 기록하는 MITM 프록시와 무관하다; 비어 있으면 직접 연결.
		EnableWebSearch:       m.webSearch.Enabled,
		WebSearchBackend:      m.webSearch.Backend,
		BraveSearchAPIKey:     m.webSearch.BraveKey,
		TavilySearchAPIKey:    m.webSearch.TavilyKey,
		DeepSeekSearchBaseURL: m.webSearch.DeepSeekBaseURL,
		DeepSeekSearchAPIKey:  m.webSearch.DeepSeekAPIKey,
		DeepSeekSearchModel:   m.webSearch.DeepSeekModel,
		WebSearchProxy:        m.webSearch.Proxy,
		BashEnv:               proxyEnv(m.proxyAddr, m.proxyCACert), // Bash 하위 명령은 기본적으로 프록시를 거치고 CA를 신뢰
		WorkingDir:            mainDir,                              // 본 작업 작업 디렉터리 <workDir>/tasks/<taskID>
		ToolOutputDir:         cmdOutDir(mainDir),
		MaxTurns:              m.maxTurns,                             // 0 = unlimited (configurable in agent management)
		Compaction:            compactionConfig(m.compactionWindow()), // long chats stay within the window
		Todos:                 actool.NewTodoStore(),                  // 세션 레벨 임시 할 일(TodoWrite), 순수 계획용, 종료하면 버림
		// 예산(스텝 수) 도달→ SDK가 마무리 실행: 사용자에게 진행 상황 요약 한 문장을 출력. Prompt 와 마무리 라운드 수는 백그라운드에서 편집 가능(기본 10 라운드).
		Settlement:   wrapupSettlement("mainagent", nil),
		NonStreaming: m.nonStreaming(), // 이 profile이 비스트리밍을 선택하면 Provider.Complete 를 거침
		MaxTokens:    m.maxTokens(),    // 0 = 상한을 보내지 않음, 서버 측 기본값이 결정
	}
	if m.tx != nil { // persist raw human↔AI conversation; one accumulating file per segment
		opts.Transcript = m.tx
		// Segment 0 keeps the legacy "exp%d-main" name so existing transcripts still
		// load; each new session (seg>=1) gets its own file for a clean context.
		opts.SessionID = fmt.Sprintf("exp%d-main", ts.ID())
		if mainSeg > 0 {
			opts.SessionID = fmt.Sprintf("exp%d-main-s%d", ts.ID(), mainSeg)
		}
	}
	// 실험 기능: 켜면 noa가 컨텍스트 압축을 인계(아카이브는 <workDir>/noa/<SessionID> 아래에 모이며, 영속).
	// session id는 transcript와 같은 규칙(분할 인지)으로, 아카이브와 복구를 정렬한다.
	noaSession := fmt.Sprintf("exp%d-main", ts.ID())
	if mainSeg > 0 {
		noaSession = fmt.Sprintf("exp%d-main-s%d", ts.ID(), mainSeg)
	}
	enableNoa(&opts, m.noaEnabledFn, m.workDir, noaSession, noaWarn(noaSession))
	ctx = attachSideCapture(ctx, &opts)
	s := agentcore.NewSession(opts)
	defer s.Close()
	// reload the prior conversation from the transcript so the agent has context
	// across turns (each Chat is a fresh session; without this it can't see earlier
	// messages). First turn: no file yet → Resume loads nothing and proceeds.
	if m.tx != nil {
		_ = s.Resume(opts.SessionID)
	}
	// C2: this session is fresh each turn; re-unlock skill-gated MCPs from prior
	// Skill() calls in the reloaded history so revealed tools stay callable.
	seedUnlockFromHistory(s.Messages(), def.UnlockSkill)
	text, _, err := captureRunSession(ctx, s, message, func(r db.Activity) {
		if emit != nil {
			r.Worker = "mainagent"
			emit(r)
		}
	})
	return text, err
}
