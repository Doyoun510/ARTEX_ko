package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/intercept"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// Planner is the event-driven LLM planner (docs §4.3): each time the asset or
// exploration graph changes (debounced), it reads the exploration route, queries
// assets, judges whether the task goal is met, and emits 0..N exploration intents
// into the frontier. It is the sole intent generator.
type Planner struct {
	findingRecorder   FindingRecorder
	prov              llm.Provider
	model             string
	tx                *transcript.Store                      // raw LLM conversation persistence (nil = off)
	window            int                                    // context window in tokens (for compaction)
	windowFn          func() int                             // optional dynamic task-chain minimum
	maxTurns          int                                    // max agent turns per run (0 = unlimited)
	killWork          func(intentID int64) error             // engine callback to terminate a running work (nil = off)
	steerWork         func(intentID int64, msg string) error // engine callback to steer a running work mid-run (nil = off)
	proxyAddr         string                                 // recording proxy for WebFetch (empty = direct)
	proxyCACert       string                                 // recording proxy's CA cert path (HTTPS verify)
	webSearch         WebSearchOpts                          // web_search tool backend selection (off by default)
	workDir           string                                 // shared work dir (surfaced in prompt as artifact-output target)
	injectConstraints func() bool                            // resolver: inject task operation constraints into system prompt? (nil = yes)
	nonStreamingFn    func() bool                            // resolver: use non-streaming (Complete) path? (nil = streaming)
	noaEnabledFn      func() bool                            // resolver: use experimental noa compaction? (nil = off)
	maxTokensFn       func() int                             // resolver: per-reply output cap (nil/0 = send no cap)
	compactor         *Compactor                             // cold-node compaction (§7); nil = disabled

	// todos keeps ONE plan-scratchpad per task (keyed by exploration id) so the
	// planner's multi-step plan survives across wake-ups — each Plan() is a fresh
	// session, but the shared store lets it record a serial exploit chain once and
	// dispatch it step-by-step over rounds instead of front-loading it in parallel.
	todoMu sync.Mutex
	todos  map[int64]*actool.TodoStore
}

func NewPlanner(prov llm.Provider, model, workDir string, tx *transcript.Store, window, maxTurns int) *Planner {
	return &Planner{prov: prov, model: model, workDir: workDir, tx: tx, window: window, maxTurns: maxTurns, todos: map[int64]*actool.TodoStore{}}
}

func (p *Planner) SetCompactionWindowResolver(fn func() int) { p.windowFn = fn }

// SetCompactor wires the cold-node compactor (cold-digest §7). Called each
// planner wake-up to advance the round counter, maintain cold stamps, and
// (off the hot path) fold cold nodes into digests. nil = feature disabled.
func (p *Planner) SetCompactor(c *Compactor) { p.compactor = c }

// SetNonStreaming wires a resolver deciding whether runs use the non-streaming
// model path (true = non-streaming). nil/unset = streaming (default).
func (p *Planner) SetNonStreaming(fn func() bool) { p.nonStreamingFn = fn }

func (p *Planner) nonStreaming() bool { return p.nonStreamingFn != nil && p.nonStreamingFn() }

// SetNoaEnabled wires a resolver deciding whether runs use the experimental noa
// context-compression mechanism. nil/unset = off (built-in compaction). Read per
// run so the settings toggle takes effect without rebuilding the agent.
func (p *Planner) SetNoaEnabled(fn func() bool) { p.noaEnabledFn = fn }

// SetMaxTokens wires a resolver for the per-reply output cap. nil/unset or 0 =
// send no cap and let the endpoint decide. Read per run, like nonStreaming.
func (p *Planner) SetMaxTokens(fn func() int) { p.maxTokensFn = fn }

func (p *Planner) maxTokens() int {
	if p.maxTokensFn == nil {
		return 0
	}
	return p.maxTokensFn()
}

func (p *Planner) compactionWindow() int {
	if p.windowFn != nil {
		return p.windowFn()
	}
	return p.window
}

// SetProxy points the planner's WebFetch at the recording proxy plus the CA cert
// it trusts to verify HTTPS through it (empty addr = direct).
func (p *Planner) SetProxy(addr, caCert string) { p.proxyAddr, p.proxyCACert = addr, caCert }

// SetWebSearch selects the web_search backend for the planner (off by default).
func (p *Planner) SetWebSearch(o WebSearchOpts) { p.webSearch = o }

// SetConstraintInject wires a resolver deciding whether this task's operation
// constraints get injected into the planner system prompt. Read per round so the
// settings toggle takes effect without rebuilding the agent. nil = inject (default).
func (p *Planner) SetConstraintInject(fn func() bool) { p.injectConstraints = fn }

// wantConstraints reports whether constraint injection is enabled (default yes).
func (p *Planner) wantConstraints() bool { return p.injectConstraints == nil || p.injectConstraints() }

// todoFor returns the task's persistent planning todo store, creating it on first
// use. Shared across all of this task's planner wake-ups.
func (p *Planner) todoFor(expID int64) *actool.TodoStore {
	p.todoMu.Lock()
	defer p.todoMu.Unlock()
	s := p.todos[expID]
	if s == nil {
		s = actool.NewTodoStore()
		p.todos[expID] = s
	}
	return s
}

// SetKillWork wires the engine's per-work terminate callback so the planner's
// kill_work tool can stop a single running worker.
func (p *Planner) SetKillWork(fn func(intentID int64) error) { p.killWork = fn }

// SetSteerWork wires the engine's per-work steering callback so the planner's
// steer_work tool can inject a mid-run course-correction into a running worker.
func (p *Planner) SetSteerWork(fn func(intentID int64, msg string) error) { p.steerWork = fn }

// renderPlannerTodos formats the persistent planning todo for injection into the
// wake-up prompt (empty when there are no todos yet — first wake-up).
func renderPlannerTodos(items []actool.Todo) string {
	if len(items) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\n[당신의 계획 할 일(깨어남 간 유지, 지난 라운드에 당신이 쓴 것)]:\n")
	for _, it := range items {
		mark := map[actool.TodoStatus]string{actool.TodoPending: "☐", actool.TodoInProgress: "▶", actool.TodoCompleted: "✔"}[it.Status]
		if mark == "" {
			mark = "☐"
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", mark, it.Content))
	}
	b.WriteString("이에 따라 진행: [선행 단계가 완료됨 / 그것이 의존하는 fact가 이미 존재함]인 다음 단계에만 의도를 파견; TodoWrite로 목록을 갱신(fact로 충족된 단계를 completed로 표시). 목록에 이미 pending/in_progress인 단계를 중복 파견하지 마라.")
	return b.String()
}

// TriggerEvent describes what concretely caused this planning round to fire, so
// the planner looks first at the actual change instead of re-scanning the whole
// overview. Kind:
//
//	"done"    — a worker finished intent IntentID (its output conclusion is fetched).
//	"finding" — a worker reported a finding on intent IntentID (Detail = 요약).
//	"goal"    — the human (via 메인 agent의 set_goals) added one OR MORE goals in a
//	            single call (Goals = 이번에 추가한 목표 텍스트, 1+개; set_goals 일괄 지원).
//	"goal_deleted" — the human deleted a goal from 개요의 목표 관리 (Detail = 삭제된 목표 텍스트).
//	"goal_edited"  — the human edited a goal from 개요의 목표 관리 (OldGoal→NewGoal 텍스트).
//	"cancelled" — the human deleted intent IntentID (Detail = 삭제 사유). The intent is
//	            stopped (not deleted) and the reason is attached to it as a fact.
type TriggerEvent struct {
	Kind     string
	IntentID int64
	Detail   string
	Summary  string   // Kind=="cancelled" 전용: 삭제 전 캡처한 의도 요약(실제 삭제 후 노드가 없어져 다시 조회 불가)
	Goals    []string // Kind=="goal" 전용: 이번 set_goals로 추가한 목표 텍스트(1개 또는 여러 개)
	OldGoal  string   // Kind=="goal_edited" 전용: 수정 전 목표 텍스트
	NewGoal  string   // Kind=="goal_edited" 전용: 수정 후 목표 텍스트
	Hints    []string // Kind=="hint" 전용: 이번 add_hint로 추가한 힌트 텍스트(1개 또는 여러 개)
}

// renderTriggers spells out the change(s) that fired this round: for a finished
// worker — which intent + its output conclusion; for a finding — which intent +
// what was found. Empty for time/heartbeat wakes. Reads the store (best-effort;
// a blank field never blocks the round).
func renderTriggers(ts *db.ExplorationStore, evs []TriggerEvent) string {
	if len(evs) == 0 || ts == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\n[이번에 이 라운드를 트리거한 실제 변동(먼저 여기를 보고 방향 보충 여부를 결정)]:")
	for _, ev := range evs {
		switch ev.Kind {
		case "goal":
			if len(ev.Goals) == 1 {
				b.WriteString(fmt.Sprintf("\n- 사람(메인 agent)이 목표 하나를 추가: %s —— 새로 달성해야 할 목표, 이에 따라 탐색 방향을 보충하라(대응 의도가 아직 없으면).", ev.Goals[0]))
			} else {
				b.WriteString(fmt.Sprintf("\n- 사람(메인 agent)이 목표 %d개를 추가: %s —— 모두 새로 달성해야 할 목표, 대응 의도가 아직 없는 목표마다 탐색 방향을 보충하라.", len(ev.Goals), strings.Join(ev.Goals, ", ")))
			}
		case "hint":
			if len(ev.Hints) == 1 {
				b.WriteString(fmt.Sprintf("\n- 사람(메인 agent)이 전략 힌트 하나를 추가: %s —— 탐색 그래프에 연결됨, 이에 따라 탐색 방향을 조정/보충하라(대응 의도가 아직 없으면).", ev.Hints[0]))
			} else {
				b.WriteString(fmt.Sprintf("\n- 사람(메인 agent)이 전략 힌트 %d개를 추가: %s —— 모두 탐색 그래프에 연결됨, 하나씩 이에 따라 탐색 방향을 조정/보충하라.", len(ev.Hints), strings.Join(ev.Hints, ", ")))
			}
		case "goal_deleted":
			b.WriteString(fmt.Sprintf("\n- 사람이 이 목표를 삭제: %s —— 해당 목표가 제거됨, 이에 따라 남은 목표/방향을 재판단하라(더는 그것을 위해 의도를 파견할 필요 없음).", ev.Detail))
		case "goal_edited":
			b.WriteString(fmt.Sprintf("\n- 사람이 목표를 수정, '%s'에서 '%s'로 —— 새 목표에 따라 탐색 방향을 조정하라(기존 방향이 더는 맞지 않으면 파견 중단).", ev.OldGoal, ev.NewGoal))
		case "finding":
			b.WriteString(fmt.Sprintf("\n- 의도 #%d(%s)의 worker가 finding 하나를 보고: %s", ev.IntentID, intentSummary(ts, ev.IntentID), ev.Detail))
		case "cancelled":
			// 의도 내용은 삭제 시 캡처한 Summary를 우선 사용(실제 삭제 후 노드가 없어 intentSummary로 조회 불가).
			sm := ev.Summary
			if sm == "" {
				sm = intentSummary(ts, ev.IntentID)
			}
			b.WriteString(fmt.Sprintf("\n- 사용자가 의도 #%d 삭제, 의도 내용: %s, 삭제 사유: %s. 이 의도는 삭제됨(더는 실행 안 함); 이에 따라 다시 계획하라.", ev.IntentID, sm, ev.Detail))
		default: // "done"
			b.WriteString(fmt.Sprintf("\n- 의도 #%d(%s)의 worker 종료, 출력 결론: %s", ev.IntentID, intentSummary(ts, ev.IntentID), workerOutput(ts, ev.IntentID)))
			if fids := factIDsYielded(ts, ev.IntentID); fids != "" {
				b.WriteString(fmt.Sprintf("; 이 의도가 새로 낳은 사실 id: %s ", fids))
			}
		}
	}
	b.WriteString("\n(완전한 세부는 node_detail / get_worker_output / list_findings로 다시 조회 가능.)")
	return b.String()
}

// factIDsYielded lists the fact ids an intent produced this run as "#12, #15", so the
// planner can jump straight to the round's incremental facts. Empty (best-effort) when
// the intent yielded no facts or the lookup fails.
func factIDsYielded(ts *db.ExplorationStore, id int64) string {
	ids, err := ts.FactsYielded(id)
	if err != nil || len(ids) == 0 {
		return ""
	}
	parts := make([]string, len(ids))
	for i, fid := range ids {
		parts[i] = fmt.Sprintf("#%d", fid)
	}
	return strings.Join(parts, ", ")
}

// intentSummary reads an intent node's one-line summary (best-effort, "?" on miss).
func intentSummary(ts *db.ExplorationStore, id int64) string {
	n, err := ts.GetNode(id)
	if err != nil || n == nil {
		return "?"
	}
	var p map[string]any
	if json.Unmarshal(n.Payload, &p) == nil {
		if s, ok := p["summary"].(string); ok && s != "" {
			return s
		}
	}
	return "?"
}

// workerOutput returns the finished worker's conclusion for an intent — the last
// 'result' (else 'text') activity's full detail, truncated. Same source get_worker_output uses.
func workerOutput(ts *db.ExplorationStore, id int64) string {
	acts, _, err := ts.ActivityList(&id, 0, 1000)
	if err != nil {
		return "(출력 가져오기 실패)"
	}
	var pick *db.Activity
	for i := range acts {
		if acts[i].Kind == "result" {
			pick = &acts[i]
		} else if acts[i].Kind == "text" && pick == nil {
			pick = &acts[i]
		}
	}
	if pick == nil {
		return "(이 work에 아직 출력 기록 없음)"
	}
	out, _ := ts.ActivityDetail(pick.ID)
	if out == "" {
		out = pick.Summary
	}
	return truncOutput(out, 800)
}

// truncOutput caps a worker-output blob so the trigger context doesn't bloat the
// system prompt every round; full text is one get_worker_output call away.
func truncOutput(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + " …(잘림, 전체는 get_worker_output 참고)"
}

// renderGraphOverview folds the pre-computed graph_overview snapshot into the
// wake-up prompt so the planner starts each round with the full situation in
// hand — saving the round-trip it would otherwise spend calling the tool. It is
// the exact same JSON graph_overview would return; deeper detail is still one
// tool call away (node_detail / list_facts / …).
func renderGraphOverview(data map[string]any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "" // fall back to the model calling graph_overview itself
	}
	return "\n\n[이번 라운드 상황(graph_overview 미리 가져옴, 그 도구를 호출한 반환과 동일; 세부가 필요하면 node_detail/list_facts 등을 호출)]:\n" + string(b)
}

// plannerDefaultTmpl is the built-in EDITABLE body (섹션 [A]) of the planner prompt,
// seeded into agent_prompts. Goal is a {{.Goal}} template var; the 중간 산출물 출력 규약
// tail is code-owned (artifactSpec) and appended by plannerSystem after rendering.
const plannerDefaultTmpl = `당신은 사이버보안 플랫폼의 승인된 침투 테스트 시스템의 "planner"로, 자주 깨어난다(그래프가 바뀔 때마다 깨어남). 역할: 상황 읽기 → 목표 판정 → **커버되지 않은 새 방향이 확실히 있을 때만** 탐색 의도를 보충한다. 당신은 planner이지 worker가 아니다: 이번 라운드의 모든 산출물은 [의도 생성/명확화] 또는 [목표 판정]뿐이며, plan 안에서 직접 일을 해서는 절대 안 된다.

작업 목표: {{.Goal}}

**이번 라운드에 의도를 몇 개 산출할지(먼저 이것부터 분명히)**:
- **하드 레드라인(최우선)**: [목표 미달성]이고 [현재 open 또는 running 의도가 하나도 없음](frontier_open=0 이고 running_intents 비어 있음)이면, 이번 라운드에 [반드시] 목표로 전진하는 의도를 최소 하나 산출한다——실행 중인 work도 없고 대기 중인 방향도 없을 때 0 의도 산출=작업 정지; 알려진 방향이 recent_done에만 있더라도 아래 done/exhausted/blocked 판단에 따라 새로 하나 열거나 이어서 하나 파견한다.
- 하드 레드라인 외에는, **0개 의도 산출도 정상적인 결과이나 정당한 이유가 있어야 한다**("적게 파견하는 게 안전"이라는 기본값이 아니다): ①**이미 커버됨**——떠올린 방향이 모두 아직 open/running인 의도로 처리 중(표현만 바꿔 이미 있는 의도를 중복 생성하는 것은 심각한 오류); ②**의존 대기**——다음 단계가 현재 실행 중인 work의 산출에 의존하는데 아직 안 나옴(이때 억지로 파견하면 하위가 전제를 못 받아 헛돈다, 다음 깨어남에 그래프가 갱신된 뒤 파견한다).
- 반대로: [커버되지 않았고 실행 중 work에 의존하지 않는] 새 방향이 확실히 있거나, 목표 미달성이고 범위 내에 아직 미테스트 면이 있으면 파견해야 한다——0 의도를 게으른 기본값으로 삼지 마라.

**매 깨어남의 의사결정 흐름**:

1. **완전한 상황이 이 프롬프트 아래에 첨부돼 있다**(graph_overview의 반환 그대로, 다시 호출할 필요 없음): task(원래 제목+목표/루트 노드), 자산 개수, goals+상태, open/running/recent_done 의도, sites_without_endpoints(엔드포인트 없는 사이트, 탐색 대기 방향 암시), facts(탐색 사실 수, 취약점과는 별개 범주), recent_facts({id,summary,confidence?}).
   - **범위**: 탐색 노드(goals/의도/facts/findings)는 본 작업만; **자산 그래프는 전역 공유**(여러 작업이 같은 하나, 자산 개수는 범위 내 전역이며 본 작업 전용 아님)——본 작업과 무관한 자산이 나오면 무시.
   - **계보**: 각 의도는 parents(상위: 어떤 사실/의도에서 파생됐는지)와 yields(하위: 어떤 사실/발견을 낳았는지)를 가지며, recent_facts 각 항목은 from_intent를 가진다; 이를 통해 "어떤 사실이 어느 방향에서 왔는지, 종합해 새 방향을 낼 수 있는지"를 이해한다.
   - **부정/의심 관찰**(recent_facts의 "포트 닫힘/주입 불가" 등)은 worker의 관찰이지 확정된 결론이 아니다: 믿기 전에 먼저 node_detail(id)로 evidence를 본다——evidence가 탄탄하고 confidence=observed이며 수단을 다 써봤을 때만 그 방향이 잠시 막힌 것으로 본다; evidence가 없거나 "그래 보임/한 번만 탐색"이거나 confidence=inferred인 것은 [아직 미규명]으로 처리하고, 범위 내이며 다른 의도가 커버하지 않으면 기본적으로 재확인 의도 하나를 파견해 입증 또는 반증한다(**같은 부정 방향은 최대 한 번만 재확인**; 재확인 후에도 부정이고 증거가 합리적이면 그 결론을 존중하고 더 파견하지 않는다).
   - **더 깊은 세부는 필요할 때만 호출**: list_facts(페이지네이션, 최신순, 기본 20, q 필터·before 페이지 넘김 가능, total/has_more 포함), list_findings(전체 취약점), node_detail(id)(완전한 증거/상세; 목록/recent_facts는 요약만), list_assets(pull: q 검색, type/company_id/task_id 필터, 페이지네이션, 또는 id/ids 직접 조회), asset_neighbors. 자산은 전역 공유이니 기본으로 전량 가져오지 마라.

2. **목표 판정(핵심 역할)**: goals 필드에 목표와 상태가 이미 들어 있다; 어떤 발견/사실로 증명된 미달성 목표는 prove_goal(goal_id, evidence_id, reason)을 호출해 met로 표시한다. **당신이 표시한 것이 마침 마지막 미완료 목표이면, 시스템이 작업 전체 완료를 자동 판정한다**——마무리는 오직 prove_goal을 하나씩 하는 것으로만 이뤄지며, 다른 "원클릭 완료" 수단은 없다.
   - ⚠️ **목표 달성의 정량적 확인(조기 도장 엄금)**: 목표에 정량 조건(커버리지 X% 달성, flag N개 획득, 특정 권한 획득)이 있으면, prove_goal 전에 [반드시] 위 graph_overview의 실측값(coverage.pct, findings_total 개수 등)을 대조한다: 미달이면 prove_goal [금지]하고, 의도를 파견해 차이를 메운다; "대체로 달성/핵심은 확보"를 이유로 조기에 met로 표시하지 마라. 예: 커버리지 100% 요구인데 실측 coverage.pct=40% → 미달성, 계속 보완 의도를 파견.

3. **(선택, 개시 때만, 극히 경량) 탐지로 이해**: 그래프에 fact가 거의 없고(recent_facts가 거의 비었고 작업이 막 시작됨) 상황만으로 초기 의도를 구체화할 수 없을 때만, Bash 등으로 대상에 극소량·읽기 전용 탐지를 한다(예: 1–2회 curl로 첫 페이지/지문 확인). **유일한 합법 산출물은 더 정확한 의도 설명 한 문장**이다——취약점의 발견/검증/익스플로잇이 절대 아니고, 엔드포인트/디렉터리/파라미터 열거 결과도 아니다(그것은 worker의 일, 의도로 써서 파견). 세 가지 하드 경계:
   - 그래프에 이미 worker가 낸 fact가 있으면(facts>0 / recent_facts 비어 있지 않음) → 스스로 탐지하는 것 [금지], 모든 판단을 기존 fact에 근거하며, 이번 라운드 산출물은 "새 의도 파견" 또는 "종료"뿐; 어떤 단서를 깊이 파고 싶으면 → 의도를 파견해 worker가 조사하게 하라, 직접 curl하지 마라.
   - 개시 때라도 최대 ≤3회만 탐지하고 손을 뗀다, 오직 초기 의도를 분명히 하기 위해서; "방향을 빨리 정하기"가 아니라 "깊이 파고드는 검증"을 하고 있음을 발견하면(엔드포인트/디렉터리 하나씩 열거, id 하나씩 시도, 디코딩 체인, 같은 엔드포인트 반복 탐색, 모든 인젝션/권한 우회/취약점 테스트·검증——전부 worker의 중노동) 즉시 멈추고 의도로 쓴다.
   - 기존 사실/상황으로 판단할 수 있으면 탐지할 필요가 전혀 없다.

4. **어떤 새 방향을 보충할지 결정**: **여기서 "절제"는 [이미 있는 의도를 중복하지 않음]만을 뜻하지, "적게 파견할수록 좋다"가 아니다**——목표 미달성 시 기본 물음은 "목표에 다가가기 위해 더 깊고, 더 강하고, 아직 커버되지 않은 수법이 뭐가 있나"이지 "마무리해도 되나"가 아니다. 의도는 [열린 탐색 방향]이며(고정 유형/메뉴가 아님), 기존 사실·자산·목표를 결합해 스스로 방향을 판단하고, open + running + recent_done와 하나씩 대조한다:
   - 이미 open/running이 커버 → 생성하지 않음(처리 중).
   - recent_done에 나온 적 있음 → **먼저 그 의도의 state(각 항목에 있음)로 어떻게 멈췄는지 분별한 뒤 결정**:
     · **done(정상 완료)**: 커버됨 → 원형 그대로 재파견 안 함; 막다른 길인지는 state가 아니라 그것이 yields한 fact 결론으로 본다; [실질적인 새 메커니즘](새 사실/자산/파라미터/명백히 다른 수법)이 있을 때만 재파견하고, summary에 지난번과의 차이를 분명히 쓴다; 표현만 바꾸거나 "한 번 더 하면 될지도"는 인정 안 됨, 재시도 금지.
     · **exhausted(예산 소진, 절반 탐색하다 끊김, 일부만 기록) / blocked(모델 또는 네트워크 실패, 거의 탐색 못 함)**: 둘 다 중도에 마무리 못 하고 정보 불완전——먼저 get_worker_trace / get_worker_output으로 실제 어디까지 했고 어디서 막혔는지 본 뒤 아래에서 고른다: 돌파 직전 예산에 끊김 → "지난 진행에 이어 계속" 파견; 순수 외부 장애로 못 돎(blocked가 흔히 그러함) → 같은 방향 바로 재파견; 매번 같은 곳에서 막힘 → 수법/방향 변경. 근거는 언제나 trace의 실제 진행이지 state 자체가 아니다.
   - 어떤 의도도 전혀 커버하지 않은 완전히 새 방향 → 생성.
   - 알려진 모든 방향이 아직 open/running인 의도로 커버됨 → 생성하지 않고 바로 종료(실행 중/대기 중 work가 있으니 그것들이 진행되길 기다림); 단 recent_done 커버만 남고 open/running이 없으며 목표 미달성이면 → 맨 위 하드 레드라인에 따라 반드시 새로 열거나 이어서 파견.
   - **깊이가 커버리지보다 우선**: coverage는 하한/달성 검증 항목이지 탐색 목표 자체가 아니다; 고가치 진입점(RCE/권한 상승/데이터 유출로 이어질 수 있는)을 발견한 뒤에는, 커버리지를 맞추려 넓게 펼쳐 자산을 하나씩 얕게 테스트하지 말고, 그 경로를 [끝까지 뚫는] 의도를 우선 파견한다.
   - **경로 다양성 유지, 조기 수렴 금지**: 목표 미달성 시, 기존 의도가 모두 같은 경로/진입점에 몰려 있고 [본질적으로 다른] 미커버 방향(다른 진입면/다른 종류 자산/다른 익스플로잇 체인)이 있으면, 같은 선상에 동의어 의도를 더하지 말고 그 분기 방향을 우선 보충한다(표현이 아니라 실질 차이를 본다); 그 분기 방향이 이미 기존 의도로 커버되면 역시 생성하지 않는다. 이상적인 건 메커니즘이 다른 2–3개 경로 병존(예: "업로드 체인으로 공격"과 "인증 우회로 공격")이며, 어느 하나가 [목표 근접] 증거를 내놓은 뒤에 자원을 거기로 집중한다. **단 다양성은 언제나 맨 위 [동작 제약]에 종속된다**: 제약으로 배제된 진입면/포트/호스트/동작은 본질적으로 다르더라도 절대 의도를 생성하지 않는다.

   **직렬 익스플로잇 체인: 단계별로 파견하고, 병렬로 쪼개지 마라.** 강한 의존의 직렬 체인(①→②→③, 뒤 단계가 앞 단계의 실제 산출에 의존): 한 번에 병렬로 내려보내지 마라(하위가 아직 없는 전제를 못 받아 중복/헛돌기만 함); TodoWrite로 체인 전체를 할 일로 기록하고(각 단계를 한 항목), 이번 라운드엔 "전제가 충족된" 단계(보통 첫 단계)만 파견하며, 그것이 fact를 낸 뒤 다음 깨어남(프롬프트에 할 일 목록이 딸려 옴)에 다음 단계를 파견하고 충족된 것을 completed로 표시한다. "같은 일"을 두 개로 쪼개지 마라("트리거 지점 확인"과 "트리거 지점 발동"은 같은 단계); [병렬이며 서로 의존하지 않는] 차원(예: 무관한 여러 엔드포인트 열거)만 다중 의도로 병렬 처리한다.

5. **제출**: [한 번의] add_intent로 선별한 새 방향을 일괄 제출한다(intents 배열, 최고가치 최대 4개, 하나씩 여러 번 호출하지 마라):
   - **summary**: 그 방향을 자연어 한 문장으로 기술(테스트 대상 전체 주소 + 무엇을 + 왜), 고정 분류에 끼워 맞추지 말 것; 중복 제거는 주로 이것을 기존 의도와 비교해 이뤄진다.
   - **asset_ids**: 이 방향에서 테스트/공격할 대상 자산 id(가능하면 전달, 0/1/여러 개, list_assets에서)——방향이 구체적 자산(사이트/엔드포인트/파라미터/호스트)을 중심으로 하면 반드시 전달하며, 커버리지 중복 제거·자산 체인 연결에 쓰이고, 여러 자산에 걸치면 모두 전달; 순수 전역 정찰이고 구체적 자산이 없을 때만 비운다.
   - **parent_ids**: 이 방향이 어떤 상위 노드를 종합해 나왔는지(선택, 0/1/여러 개)——여러 사실이 결합해 한 의도를 낳으면 모두 전달하고, 어떤 상위 의도/발견에서 파생됐으면 그 id도 전달하며, 최상위의 완전히 새 방향은 비운다.

중복하지 말고 억지로 채우지 마라; 단 목표 미달성이고 미커버이며 더 깊은 수법이 있을 때는 파견해야 한다. 간결하게, 집중해서, 효율적으로.`

func plannerSystem(goal, dataDir, workDir string) string {
	body := renderSystem("planner", plannerDefaultTmpl, PlannerVars{Goal: goal, DataDir: dataDir, Now: nowStr()})
	return body + artifactSpec(workDir)
}

// Plan runs one planning round. emit, if non-nil, receives the planner's execution
// steps (so users can see how it reads the situation and judges goals — the
// planner is the intent generator and was previously a black box). Returns whether
// the planner judged the goal met.
// triggers carries the concrete change(s) that fired this round — worker(s) done
// and/or finding(s) reported (may be several — the engine debounces a burst; empty
// for time/heartbeat wakes). They are spelled out at the top of the prompt so the
// planner looks first at the actual change (which intent, its output/finding).
func (p *Planner) Plan(ctx context.Context, taskID int64, as *db.AssetStore, ts *db.ExplorationStore, goal string, triggers []TriggerEvent, emit func(db.Activity)) (met bool, reason string, err error) {
	// cold-digest §2.3/§7: advance this task's planner-round counter, maintain the
	// cold_since_round stamps, and (if a threshold is hit) kick off background
	// compaction. Synchronous part is cheap (a few queries); the LLM compaction
	// runs in a detached goroutine so it never adds latency to this round.
	p.compactor.OnPlannerRound(ctx, ts)
	tsx := NewToolSet(ts, "planner")
	tsx.SetFindingRecorder(p.findingRecorder)
	if as != nil {
		tsx.SetAssetStore(as, as.Companies())
	}
	tsx.SetTaskID(taskID)
	tsx.SetCoverageEnabled(as == nil || as.CoverageEnabled(taskID))
	tsx.killWork = p.killWork   // enable kill_work tool (nil = unavailable)
	tsx.steerWork = p.steerWork // enable steer_work tool (nil = unavailable)
	if origin, _ := ts.OriginFactID(); origin > 0 {
		tsx.SetOwnerNode(origin) // planner-side anchors default to the task root (origin fact)
	}
	// 도메인 도구 + 기본 디폴트 도구 집합(Read/Write/Edit/MultiEdit/LS/Glob/Grep/Bash)
	// 자산 커버리지 기능이 꺼지면 list_untested_assets만 제거하며, add_task_scope와 insertAssets의 task_scope 누적은 유지한다.
	base := append(tsx.DropCoverageTools(tsx.PlannerTools()), actool.DefaultTools()...)
	ctx = WithRunInfo(ctx, RunInfo{TaskID: taskID, ExplorationID: explorationID(ts)})
	tools, def, cleanup := AugmentTools(ctx, "planner", base)
	defer cleanup()
	// 핵심 상황(방금 완료된 의도 + 미리 가져온 전체 그래프)을 [이번 라운드 user 입력](아래 input 참고)으로 옮김, system은
	// 정적 계획 본문만 남긴다. move-out으로 system이 매 라운드 안정돼 캐시에 유리; 대가는 단일 라운드가 길어지면 상황이
	// compaction에 압축될 수 있음(planner 단일 라운드는 보통 짧아 위험 낮음). situational은 아래 input에 결합된다.
	situational := renderTriggers(ts, triggers) + renderGraphOverview(tsx.graphOverviewData())
	// 작업 수준 deadline / 종국 모드(ctx로 주입, taskclock.go 참고). 종국 라운드는 작업 타임아웃
	// planner 마무리 문구를 [이번 라운드 동작 지시]로 이번 라운드 user 입력(situational과 함께)에 결합해, 마지막
	// 목표 판정만 하고 새 의도를 내지 않게 한다.
	tc := taskClockFrom(ctx)
	if tc.Final {
		situational += "\n\n[작업 종국 마무리(이번 라운드 특수 지시, 위 일반 계획 흐름을 덮어씀)]:" + resolveTaskTimeoutWrapup("planner")
	}
	// 이 작업의 작업 디렉터리 <workDir>/tasks/<taskID>, 먼저 생성.
	taskDir := ensureRunDir(p.workDir, taskID, 0)
	ctx = intercept.WithReviewContext(ctx, taskDir, intercept.ReviewBackground{})
	sysBody := plannerSystem(goal, p.workDir, taskDir)
	if p.wantConstraints() {
		sysBody += constraintBlock(ts) // 동작 제약(있으면)을 시스템 프롬프트에 주입, 탐색 경계 설정
	}
	system, boundary := deferredSystem(sysBody, def)
	// planner는 자체 벽시계 예산 없음; deadline이 있으면 MaxDuration을 남은 시간으로 좁혀, 실행 중인 계획 라운드가
	// 작업 만료 시 마무리에 들어가게 한다(타임아웃 때문→작업 타임아웃 문구, 스텝 때문→per-run 문구).
	maxDur, clamped := clampMaxDuration(tc.DeadlineUnix, 0)
	settle := wrapupSettlement("planner", nil)
	if tc.DeadlineUnix > 0 {
		settle = wrapupSettlementForTask("planner", nil, clamped)
	}
	opts := agentcore.Options{
		Provider:        p.prov,
		SystemPrompt:    system,
		DynamicBoundary: boundary,
		Tools:           tools,
		DeferredTools:   def.Deferred,
		UnlockSet:       def.Unlock,
		PermissionMode:  permission.ModeBypass,
		EnableWebFetch:  true, // 레코딩 프록시를 거쳐 흔적을 남김; 프록시 CA를 로드해 MITM이 재서명한 HTTPS 인증서 검증
		WebFetchProxy:   p.proxyAddr,
		WebFetchCACert:  p.proxyCACert,
		// 인터넷 검색(선택). ddgs는 key 불필요; brave-free는 BraveKey 필요; tavily는 TavilyKey 필요.
		// WebSearchProxy는 독립 아웃바운드 프록시(http/https/socks5), 트래픽을 기록하는 MITM 프록시와 무관; 비면 직접 연결.
		EnableWebSearch:       p.webSearch.Enabled,
		WebSearchBackend:      p.webSearch.Backend,
		BraveSearchAPIKey:     p.webSearch.BraveKey,
		TavilySearchAPIKey:    p.webSearch.TavilyKey,
		DeepSeekSearchBaseURL: p.webSearch.DeepSeekBaseURL,
		DeepSeekSearchAPIKey:  p.webSearch.DeepSeekAPIKey,
		DeepSeekSearchModel:   p.webSearch.DeepSeekModel,
		WebSearchProxy:        p.webSearch.Proxy,
		BashEnv:               proxyEnv(p.proxyAddr, p.proxyCACert), // Bash 하위 명령은 기본적으로 프록시 경유 + CA 신뢰
		WorkingDir:            taskDir,                              // 이 작업 작업 디렉터리 <workDir>/tasks/<taskID>
		ToolOutputDir:         cmdOutDir(taskDir),
		MaxTurns:              p.maxTurns, // 0 = unlimited (configurable in agent management)
		MaxDuration:           maxDur,     // 0=무제한; deadline 있으면=deadline까지 남은 시간
		Compaction:            compactionConfig(p.compactionWindow()),
		// 깨어남 간 공유하는 계획 할 일: 직렬 체인이 여러 라운드에 걸쳐 유지되게 한다(session은 새것, store는 아님).
		Todos: p.todoFor(ts.ID()),
		// [이번 라운드] 스텝 예산 도달→ SDK가 마무리 실행: 이번 라운드에 정리된 결론을 기록(파견할 add_intent,
		// 증명 가능한 prove_goal, 직렬 체인은 TodoWrite로 기록)하며, 계획을 멈추지 않음——planner는 이후에도 반복 깨어남.
		// clamped(작업 deadline으로 좁혀짐) 시 PromptByReason 사용(wrapupSettlementForTask 참고).
		Settlement:   settle,
		NonStreaming: p.nonStreaming(), // 이 profile이 비스트리밍을 선택하면 Provider.Complete 경유
		MaxTokens:    p.maxTokens(),    // 0 = 상한 미전송, 서버 기본값으로 결정
	}
	if p.tx != nil { // persist raw LLM conversation; one accumulating file per task's planner
		opts.Transcript = p.tx
		opts.SessionID = fmt.Sprintf("exp%d-planner", ts.ID())
	}
	// 실험 기능: 켜면 noa가 컨텍스트 압축을 인계(아카이브는 <workDir>/noa/<SessionID> 아래에 모이며, 영속).
	noaSession := fmt.Sprintf("exp%d-planner", ts.ID())
	enableNoa(&opts, p.noaEnabledFn, p.workDir, noaSession, noaWarn(noaSession))
	// 상황(방금 완료된 의도 + 전체 그래프)은 이제 이번 라운드 user 입력에 결합된다(아래 input 참고). user에는 또
	// 지시 + 깨어남 간 할 일(todo는 모델 자신의 계획 메모라 재생성 가능, user에 두면 됨)이 있다.
	// 서두는 '이번 라운드에 구체적 변동이 있는지'에 따라 둘로 나뉨: 변동 있음 → 아래 [실제 변동] 블록을 가리킴; 변동 없음
	// (heartbeat 정기 점검 / hint / 복구 등) → "그래프가 바뀌었다"고 거짓말하지 말고, 실행 중 의도를 함께 재점검하도록 안내.
	lead := "방금 구체적 변동이 있음(아래 [이번에 이 라운드를 트리거한 실제 변동] 참고), 이에 따라 다음 단계를 계획하라: "
	if len(triggers) == 0 {
		lead = "이번 라운드는 **정기 점검(heartbeat 도달)/구체적 변동 신호 없음**의 깨어남이다——그래프에 새 변동이 꼭 있는 건 아니다. 실행 중 의도를 함께 재점검: 오래 진전 없거나 빗나간 것은 steer_work로 교정, 방향이 통째로 틀린 것은 kill_work로 손절; 그 뒤 목표를 판정하고 방향 보충 여부를 결정하라: "
		// heartbeat/무변동 깨어남 시, 전체 그래프에 open 또는 running 의도가 하나도 없으면 → 탐색 정지(실행 중 worker 없고,
		// 대기 방향도 없음). planner에 명확히 알리고 이번 라운드에 새 방향을 반드시 내게 하며, 실행 중 의도만 재점검하고 헛도는 라운드를 막는다.
		if active, err := ts.HasActiveIntent(); err == nil && !active {
			lead = "이번 라운드는 **정기 점검(heartbeat 도달)**의 깨어남이며, 현재 **open 또는 running 의도가 하나도 없음**——실행 중 worker도 없고 대기 중 방향도 없어 탐색이 정지됨. 당신은 **반드시** 이번 라운드에 목표로 전진하며 그래프의 기존 의도와 **서로 중복되지 않는** 새 의도를 하나 이상 산출해야 한다(0 의도 산출 금지); 먼저 아래 상황으로 목표 달성 여부를 판정하고, 미달성이면 즉시 방향을 보충하라: "
		}
	}
	input := lead + situational + "\n\n위 상황에 따라 목표를 판정하라. 목표가 [진짜 달성됨](목표 성과를 획득/목표 취약점을 확인)이면 prove_goal로 하나씩 표시한다. **하드 레드라인: 목표가 아직 미달성이고 현재 open 또는 running 의도가 하나도 없으면(frontier_open=0 이고 running_intents 비어 있음), 이번 라운드에 목표로 전진하는 의도를 반드시 최소 하나 산출한다——이때 실행 중인 work도 없고 대기 중인 방향도 없어 0 의도 산출=작업 정지다. 이미 open/running 의도가 진행 중이거나 목표가 달성된 경우에만 이번 라운드에 새 의도를 산출하지 않아도 된다.**" +
		renderPlannerTodos(opts.Todos.List())
	// MaxDuration은 이제 벽시계 도달 시 실행 중 도구를 중단하고 그 자리에서 마무리에 들어감(살아있는 ctx에서), 단일 라운드 멈춤이 더는
	// 마무리를 우회하지 않으므로 외부 하드 ctx 폴백 불필요. ctx는 pause / kill / shutdown만 전달.
	_, _, err = captureRun(ctx, opts, input,
		func(r db.Activity) {
			if emit != nil {
				r.Worker = "planner" // planner activity has no intent_id (it generates them)
				emit(r)
			}
		})
	return tsx.GoalMet, tsx.Reason, err
}
