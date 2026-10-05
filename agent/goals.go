package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
	acperm "github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// goalsDefaultTmpl is the built-in EDITABLE body (섹션 [A]) of the goals-decomposer
// prompt, seeded into agent_prompts. No template vars are used today.
const goalsDefaultTmpl = `당신은 침투 테스트 목표 분해기다. 당신의 역할은 사용자 입력에서 **최종적으로 달성할 결과**를 식별하는 것이지, 공격 단계를 계획하는 것이 아니다.

**1단계(목표 분해 전에 먼저): 동작 제약 추출**
'작업 목표 / 작업 설명'에서 운영자가 [무엇을 할 수 있고, 어떤 동작을 할 수 없는지]에 대해 명시한 규정을 식별해, set_constraints로 하나씩 등록한다(설명·목표에 동작 제약이 없으면 동작 제약 추출을 하지 않아도 된다):
- type=deny: 금지된 동작(예: '포트 스캔 금지' '운영 환경에 쓰기/삭제 동작 금지' '무차별 대입 금지' '특정 서브도메인 접촉 금지').
- type=allow: 명시적으로 허용/한정된 동작 범위(예: '패시브 정찰만 허용' '특정 도메인만 대상').
- 제약 ≠ 목표, 또한 ≠ 공격 단계: 동작 행위의 경계에 대한 규정이다.
- **제약은 반드시 [자기완결적이고 구체적 대상을 명시]해야 한다**: '현재 목표/현재 포트/현재 IP/현재 도메인/본 사이트' 같은 **지시 대명사**를 작업 목표/설명의 **구체적 값**으로 치환한다. 제약은 실행 단계 프롬프트에 단독 주입되므로, 맥락을 벗어나면 지시 대명사가 누구를 가리키는지 판단할 수 없다.
  예: 목표가 https://abc.example.net → '현재 목표만 테스트 허용'이 아니라 'abc.example.net만 테스트 허용'으로; '현재 포트만 테스트'가 아니라 '목표 포트 443만 테스트, 다른 포트는 스캔 금지'로. 원문이 '현재 목표'만 말해도 목표 주소가 명확하면 주소를 채워 넣는다.
- **목표/설명에 [명시적으로 쓰였거나 강조된] 제약만 등록하고, 지어내기 엄금**; 유형이 불확실하면 deny를 쓴다(더 보수적).
- 목표/설명에 동작 제약이 전혀 없으면 set_constraints를 **호출하지 마라**.
제약 등록(있으면)을 마친 뒤, 아래 목표 분해를 진행한다.

**목표 = 최종적으로 전달 가능/검증 가능한 결과**

**목표가 아닌 것(하위 목표로 열거 금지)**:
- 정보 수집, 정찰, 엔드포인트 스캔
- 취약점 분석과 검증 과정
- 공격 단계, 익스플로잇 수단
- 결과 검증 단계

**분해 원칙**:
- 사용자가 기술한 최종 목표가 하나뿐 → 하나 출력
- **서로 독립적인** 최종 전달물이 여럿 → 각각 열거
- 명확한 취약점 유형에 대응되면 vulnclass 표기; 정보 수집/비즈니스 로직류 목표는 비움
- 사용자가 언급하지 않은 목표를 지어내기 엄금

set_goals를 호출해 결과를 제출한다.`

// goalsScopeTail is the code-owned tail appended after the editable goals body
// WHEN an asset store + task context are available. It teaches the decomposer to
// also lift the explicit asset scope out of the goal/description and register it
// via add_task_scope. Kept in code (not the DB-editable body) so it always applies
// on released DBs and can't be edited away — same pattern as the trafficTool tail.
const goalsScopeTail = `

**추가 역할: 테스트 자산 범위 등록**
목표 분해 외에도, '작업 목표 / 작업 설명'에서 **명시적으로 주어진 테스트 자산 범위**를 식별해 add_task_scope로 등록한다(이 작업의 권한 경계이자 자산 테스트 커버리지의 분모). **최소 범위 원칙: 사용자가 명확히 지목한 그 하나의 대상만 등록하고, 절대 임의로 확대하지 마라.**
- 목표가 URL 또는 호스트명이 있는 주소(예: https://xxx.example.com/path, app.example.com) → 그 **완전한 호스트명**을 취해 kind=subdomain, value=완전한 호스트명.
  예: 목표 https://a1b2c3.lab.example.net/path → kind=subdomain, value=a1b2c3.lab.example.net(**example.net가 아님**).
  서브도메인이 있는 호스트명을 루트 도메인으로 줄이는 것을 **엄금**한다——xxx.example.com을 보고 example.com 전체를 등록하면 범위가 사용자 목표 밖으로 확장돼 최소 범위 원칙에 위배된다.
- 사용자가 준 것이 **서브도메인이 없는 루트 도메인일 때**(예: example.com을 그대로 씀), 또는 '전체 사이트 / 모든 서브도메인 / 전체 도메인'이라고 명시할 때만 → kind=root_domain, value=example.com을 쓴다.
- 순수 IP 또는 네트워크 대역 → kind=ip / cidr, value=IP 또는 CIDR.
- 회사 범위(company)를 등록하지 **마라**——작업이 막 생성돼 자산 시스템에 보통 아직 이 회사가 없어 등록되지 않으며, 회사 수준 범위는 이후 plan 단계가 처리한다.
기타 규칙:
- 목표/설명에 **명시적으로 쓰인** 범위만 등록; 언급되지 않은 도메인/IP를 지어내거나 추론하기 엄금.
- reason에 어느 문장에 근거했는지 간단히 적어 감사에 용이하게 한다.
- 목표/설명에 명확한 자산 범위가 전혀 없으면 add_task_scope를 **호출하지 마라**.
먼저 add_task_scope로 범위를 등록(있으면)한 뒤, set_goals를 호출해 목표를 제출한다.`

// GoalSpec is one decomposed objective.
type GoalSpec struct {
	Text      string `json:"text"`
	VulnClass string `json:"vulnclass,omitempty"`
}

// DecomposeGoals asks the LLM to break a pentest task goal into discrete,
// independently-verifiable objectives (each becomes a goal node). Returns nil if
// no provider is configured or the call yields nothing — the caller then falls
// back to a rule-based split so goal nodes always exist.
//
// prov is supplied by the caller (rather than built here from a Config) so goal
// decomposition rides the SAME provider instance as the rest of the engine — it
// shares the rate limiter, gets recorded by llmrec, and participates in LLM
// failover instead of quietly bypassing all three.
//
// desc is the task's free-text description (배경: 타깃 범위/flag 수/교전 설명 등).
// It is fed alongside the goal so the decomposer no longer splits blind — the
// prompt still forbids inventing anything the two texts don't state.
//
// emit, when non-nil, receives every LLM step (thinking/tool_use/result) with
// Worker="planner" so the round-0 goal-decomposition activity is visible in the UI.
//
// as + taskID, when non-nil/positive, wire the add_task_scope tool so the
// decomposer can register the explicit asset scope it extracts from the goal.
//
// ts is the task's exploration store: set_goals writes the decomposed goal nodes
// straight into it (the same managed tool the main agent uses to add goals at
// runtime). The returned specs are read back from the store so callers can emit
// per-goal activity and detect the "LLM produced nothing" case for their fallback.
func DecomposeGoals(ctx context.Context, prov llm.Provider, dataDir, goalText, desc string, as *db.AssetStore, ts *db.ExplorationStore, taskID int64, emit func(db.Activity)) []GoalSpec {
	if prov == nil {
		return nil
	}
	return DecomposeGoalsWithProvider(ctx, prov, dataDir, goalText, desc, as, ts, taskID, false, 0, emit)
}

// DecomposeGoalsWithProvider is the task-runtime variant used when a task has an
// ordered provider chain. It preserves the same tools and write behavior while
// letting the caller own provider selection/failover. maxTokens is the profile's
// per-reply output cap (0 = send none).
func DecomposeGoalsWithProvider(ctx context.Context, prov llm.Provider, dataDir, goalText, desc string, as *db.AssetStore, ts *db.ExplorationStore, taskID int64, nonStreaming bool, maxTokens int, emit func(db.Activity)) []GoalSpec {
	if prov == nil {
		return nil
	}
	// 목표 분해는 일회성 호출: transcript store를 달지 않으므로 agentcore가 ctx에
	// session id를 달지 않음(writer가 있을 때만 단다, agentcore.Prompt 참고). 한편 session-id 헤더로
	// 프롬프트 캐시/고정 라우팅을 하는 게이트웨이(opencode zen은 x-opencode-session이 없으면 바로 400
	// MissingSessionID)가 읽는 게 바로 ctx의 이 값이다——안 채우면 '대화는 정상, 분해는 400'이 된다.
	// 안정적인 id를 명시적으로 단다: 같은 탐색의 분해 요청이 이를 공유(캐시 히트에 유리)하고, 명명이
	// planner/worker와 충돌하지 않아 llmrec.parseSession이 올바르게 귀속할 수 있다.
	if ts != nil {
		ctx = transcript.WithSessionID(ctx, fmt.Sprintf("exp%d-goals", ts.ID()))
	}
	// worker="goals" tags the goal nodes' provenance; ts/taskID let set_goals link
	// each goal under the task root. This is the catalog's real set_goals tool, so a
	// web-edited description/schema on it applies here too.
	tsx := &ToolSet{as: as, ts: ts, taskID: taskID, worker: "goals"}
	// Description rides in the user message (same channel as the goal), NOT via the
	// {{.EngagementDescription}} template var — else a prompt that references the var
	// would inject the description twice. System prompt stays pure static instructions.
	sys := renderSystem("goals", goalsDefaultTmpl, GoalsVars{DataDir: dataDir, Now: nowStr()})
	// set_constraints는 항상 사용 가능(asset store에 의존하지 않음): 본문에 이미 '먼저 동작 제약을 뽑고 목표를 분해' 단계가 포함돼
	// (agent 편집 페이지에서 문구 수정 가능), 여기서는 도구만 연결하면 된다.
	tools := []actool.CoreTool{tsx.setGoals(), tsx.setConstraints()}
	// Wire add_task_scope only when we have a real asset store + task to write to.
	// The scope-extraction tail is appended in lockstep so the prompt never asks for
	// a tool that isn't present.
	if as != nil && taskID > 0 {
		tools = append(tools, tsx.addTaskScope())
		sys += goalsScopeTail
	}
	userMsg := "작업 목표:\n" + goalText
	if d := strings.TrimSpace(desc); d != "" {
		userMsg += "\n\n작업 설명(배경 정보, 타깃 범위/flag 수/교전 설명을 포함할 수 있음; 참고용일 뿐, 언급되지 않은 내용을 지어내지 마라):\n" + d
	}
	// Use captureRun so every LLM step is emitted as an activity record (visible in
	// the plan tab under the round-0 marker). Falls back gracefully when emit is nil.
	captureEmit := func(r db.Activity) {
		if emit != nil {
			r.Worker = "planner"
			emit(r)
		}
	}
	captureRun(ctx, agentcore.Options{
		Provider:               prov,
		SystemPrompt:           []string{sys},
		Tools:                  tools,
		PermissionMode:         acperm.ModeBypass,
		DisableBackgroundTasks: true,
		// 3단계(제약 추출 → 범위 등록 → 목표 분해)마다 도구 호출 한 번씩 필요, 마무리 전 set_goals 누락을 막으려 회합을 충분히 준다.
		MaxTurns:     8,
		NonStreaming: nonStreaming, // 이 profile이 비스트리밍을 선택하면 Provider.Complete 경유
		MaxTokens:    maxTokens,    // 0 = 상한 미전송, 서버 기본값으로 결정
	}, userMsg, captureEmit)
	// set_goals persisted the goals directly; read them back so the caller sees what
	// was written (empty slice ⇒ the LLM produced nothing ⇒ caller falls back).
	if ts == nil {
		return nil
	}
	nodes, _ := ts.ListByKind(db.KindGoal, 10000)
	var out []GoalSpec
	for _, n := range nodes {
		var p struct {
			Text      string `json:"text"`
			VulnClass string `json:"vulnclass"`
		}
		_ = json.Unmarshal(n.Payload, &p)
		if strings.TrimSpace(p.Text) != "" {
			out = append(out, GoalSpec{Text: p.Text, VulnClass: p.VulnClass})
		}
	}
	return out
}
