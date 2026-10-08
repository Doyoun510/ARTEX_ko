package agent

import (
	"strings"

	"github.com/Autumn-27/norma/harness"
)

// 마무리 프롬프트(wrap-up / settlement prompt): agent가 [스텝 소진(MaxTurns)] 또는
// [타임아웃(run_seconds/MaxDuration)]으로 종료될 때, SDK의 settlement 단계가 이 프롬프트를 주입해
// agent가 먼저 식별했지만 쓰지 않은 내용을 DB에 저장하고 한 문장 요약을 출력하게 해 미완 종료를 막는다.
//
// 각 agent의 마무리 프롬프트는 백엔드에서 필요 시 오버라이드 가능(agents.wrapup_prompt에 저장), 비우면 여기
// 내장 기본값 사용. [프롬프트 본문]만 편집 가능; 어떤 도구를 비활성화할지·마무리 자체에 몇 라운드 예산을 줄지는 코드 고정 정책.

// WrapupOverride, if set, returns the stored wrap-up prompt for an agent key and
// whether a non-empty one exists. Wired by the server to the agents table (like
// PromptOverride for system prompts). nil / empty → the built-in default is used.
var WrapupOverride func(agentKey string) (string, bool)

// WrapupMaxTurnsOverride, if set, returns the admin-configured turn budget for the
// wrap-up phase of an agent and whether a positive one exists. Wired to the agents
// table. nil / ≤0 → the built-in per-agent default (wrapupTurnDefaults) is used.
var WrapupMaxTurnsOverride func(agentKey string) (int, bool)

// 내장 기본 마무리 프롬프트, agent key로 인덱싱. worker는 과거부터 하드코딩된 settleWrapUpPrompt
// (worker.go에 정의)를 재사용, planner/mainagent는 각각 한 벌; 미매칭(커스텀 agent)은 공통 폴백.
var wrapupDefaults = map[string]string{
	"worker":    settleWrapUpPrompt,
	"planner":   plannerWrapUpDefault,
	"mainagent": mainAgentWrapUpDefault,
}

// wrapupTurnDefaults: 각 agent 마무리 단계 [자체]의 라운드 예산 내장 기본값(백엔드에서 >0으로 오버라이드 가능).
// 모두 10라운드를 줘 마무리 단계에 DB 저장 스텝이 충분하도록. 미매칭은 genericWrapupTurns.
var wrapupTurnDefaults = map[string]int{
	"worker":    10,
	"planner":   10,
	"mainagent": 10,
}

const genericWrapupTurns = 10

const plannerWrapUpDefault = "이번 라운드 계획의 스텝이 곧 소진된다——[이 라운드]만 끝나는 것이지, 시스템은 이후에도 상황 변화에 따라 당신을 다시 깨워 계획을 이어가게 하며 작업 종료가 아니므로, 여기서 전체 계획을 마무리할 필요는 없다. 이번 라운드에 이미 정리된 결론을 기록해 이 라운드를 헛되게 하지 말되, [마무리를 위해 억지로 의도를 채우지도 마라](이번 라운드 0 의도도 완전히 정상적인 결과다): (1) [지금 바로 파견해야 할] 탐색 방향을 판단했으면 add_intent로 한 번에 일괄 제출(정한 건 묵혀 두지 마라); (2) 어떤 발견/사실로 달성이 증명된 목표는 prove_goal로 met 표시(누락 금지); (3) 단계를 나눠야 하는 직렬 익스플로잇 체인을 식별했으면 TodoWrite로 기록해 다음 깨어남에 이어 파견하기 쉽게 하라. 끝나면 바로 이번 라운드를 종료하고, 요약 텍스트를 출력할 필요 없다."

const mainAgentWrapUpDefault = "당신의 스텝이 곧 소진되어 이번 상호작용이 끝난다. 더는 새 탐색/동작을 시작하지 마라. **한 문장 순수 텍스트로 단독으로** 현재 진행·핵심 결론·권장 다음 단계를 사용자에게 요약하라."

const genericWrapUpDefault = "당신은 곧 예산 소진으로 종료된다. 먼저 완료했지만 DB에 저장하지 않은 결과를 써서 되돌린 뒤, **한 문장 순수 텍스트로 단독으로** 무엇을 했고 어떤 핵심 결론을 얻었는지 요약하라(이 문장이 이번 실행의 결과로 표시된다)."

// WrapupDefault returns the built-in default wrap-up prompt for an agent key —
// used by the admin UI as the "restore default" value and empty-field placeholder.
func WrapupDefault(agentKey string) string {
	if d, ok := wrapupDefaults[agentKey]; ok {
		return d
	}
	return genericWrapUpDefault
}

// WrapupTurnsDefault returns the built-in wrap-up turn budget for an agent key —
// used by the admin UI as the "0 = default N" hint.
func WrapupTurnsDefault(agentKey string) int {
	if n, ok := wrapupTurnDefaults[agentKey]; ok {
		return n
	}
	return genericWrapupTurns
}

// resolveWrapup returns the effective wrap-up prompt: the DB override (if set and
// non-empty) over the built-in default.
func resolveWrapup(agentKey string) string {
	if WrapupOverride != nil {
		if t, ok := WrapupOverride(agentKey); ok && strings.TrimSpace(t) != "" {
			return t
		}
	}
	return WrapupDefault(agentKey)
}

// resolveWrapupTurns returns the effective wrap-up turn budget: a positive DB
// override over the built-in per-agent default.
func resolveWrapupTurns(agentKey string) int {
	if WrapupMaxTurnsOverride != nil {
		if v, ok := WrapupMaxTurnsOverride(agentKey); ok && v > 0 {
			return v
		}
	}
	return WrapupTurnsDefault(agentKey)
}

// wrapupSettlement builds the settlement config for an agent's run. Prompt and the
// turn budget are admin-editable per agent; disabled tools are code-owned policy so
// a user can't edit away the "stop probing" guardrail. Resolved fresh each run
// (reads DB live), so edits apply on the next run without a restart.
func wrapupSettlement(agentKey string, disabledTools []string) *harness.Settlement {
	return &harness.Settlement{
		Prompt:        resolveWrapup(agentKey),
		DisabledTools: disabledTools,
		MaxTurns:      resolveWrapupTurns(agentKey),
	}
}

// ---------- 작업 수준 타임아웃 마무리 문구(docs/任务级超时与收尾设计.md 참고) ----------
//
// per-run 마무리 문구와는 [두 벌]이다: per-run은 "이번 run의 예산이 끝났다"; 작업 타임아웃은
// "작업 전체가 만료돼 곧 끝난다". 의미가 자주 반대다(특히 planner: per-run은 "멈추지 말고 계속 계획",
// 작업 타임아웃은 "만료됐으니 계획 중단, 마지막 판정"). worker/planner에만 설정.

// WrapupTaskTimeoutOverride / …TurnsOverride: 작업 타임아웃 마무리 문구와 라운드 수의 DB 오버라이드
// (agents.task_timeout_wrapup_prompt / _max_turns에 연결, worker/planner만).
var (
	WrapupTaskTimeoutOverride      func(agentKey string) (string, bool)
	WrapupTaskTimeoutTurnsOverride func(agentKey string) (int, bool)
)

var taskTimeoutWrapupDefaults = map[string]string{
	"worker":  workerTaskTimeoutDefault,
	"planner": plannerTaskTimeoutDefault,
}

const workerTaskTimeoutDefault = "**작업 전체가 타임아웃 상한에 도달해 곧 끝난다**(네 이번 run의 예산이 아니라 탐색 전체가 만료됨). 마지막 기회다: (1) 식별했지만 아직 쓰지 않은 내용을 [전부] DB에 저장——새 자산 insert_assets, 탐색 결론/사실 record_fact, 확인 취약점 report_finding; (2) 더는 어떤 새 명령/탐지도 시작하지 마라; (3) **마지막에 한 문장 순수 텍스트로 단독으로** 이 의도에서의 핵심 결론을 요약하라."

const plannerTaskTimeoutDefault = "**작업 전체가 타임아웃 상한에 도달해 곧 끝난다**(이번 라운드가 아니라 작업 전체 종료). 현재 [전부]의 사실과 발견을 바탕으로 마지막 목표 판정을 하라: 증거로 달성이 증명된 목표는 prove_goal로 met 표시(누락 금지). **더는 어떤 새 의도도 생성하지 마라**(이때 의도를 파견해도 더는 실행되지 않는다). 판정 후 바로 마무리하고, 요약 텍스트를 출력할 필요 없다."

// TaskTimeoutWrapupDefault은 특정 agent의 작업 타임아웃 내장 기본 마무리 문구를 반환(백엔드 플레이스홀더/기본값 복원용).
func TaskTimeoutWrapupDefault(agentKey string) string {
	return taskTimeoutWrapupDefaults[agentKey] // 미설정(mainagent/chat)이면 빈 문자열 반환
}

// resolveTaskTimeoutWrapup: DB 오버라이드(비어 있지 않음) > 내장 기본값. 빈 문자열은 그 agent에 작업 타임아웃 문구가 없음을
// (worker/planner 아님) 뜻하며, 이때 호출자는 per-run 문구로 회귀해야 한다.
func resolveTaskTimeoutWrapup(agentKey string) string {
	if WrapupTaskTimeoutOverride != nil {
		if t, ok := WrapupTaskTimeoutOverride(agentKey); ok && strings.TrimSpace(t) != "" {
			return t
		}
	}
	return TaskTimeoutWrapupDefault(agentKey)
}

func resolveTaskTimeoutTurns(agentKey string) int {
	if WrapupTaskTimeoutTurnsOverride != nil {
		if v, ok := WrapupTaskTimeoutTurnsOverride(agentKey); ok && v > 0 {
			return v
		}
	}
	return resolveWrapupTurns(agentKey) // 기본은 per-run 라운드 수 그대로 사용
}

// wrapupSettlementForTask builds settlement for a worker/planner run that is aware
// of the task deadline. See §5 of the design doc:
//   - clamped=true  → 이번 run이 작업 deadline으로 좁혀짐: Timeout으로 마무리=작업 만료→작업 타임아웃 문구;
//     MaxTurns로 마무리=좁힌 창 안에서 스텝이 먼저 소진, 작업은 몇 분 남음→per-run 문구로 회귀.
//   - clamped=false → 작업 종료까지 아직 시간이 남음: 두 reason 모두 per-run 문구 사용(즉 wrapupSettlement로 퇴화).
//
// harness의 PromptByReason에 넘겨 마무리 시 [실제] reason에 따라 현장에서 고르므로, 빌드 시점 불일치가 없다.
func wrapupSettlementForTask(agentKey string, disabledTools []string, clamped bool) *harness.Settlement {
	perRun := resolveWrapup(agentKey)
	st := &harness.Settlement{
		Prompt:        perRun, // 폴백(clamped 아닐 때 두 reason의 값이기도 함)
		DisabledTools: disabledTools,
		MaxTurns:      resolveWrapupTurns(agentKey),
	}
	if clamped {
		if tt := resolveTaskTimeoutWrapup(agentKey); tt != "" {
			st.PromptByReason = map[harness.TerminalReason]string{
				harness.ReasonTimeout:  tt,     // 작업 만료
				harness.ReasonMaxTurns: perRun, // 스텝 먼저 소진, 작업은 시간 남음
			}
			st.MaxTurns = resolveTaskTimeoutTurns(agentKey)
		}
	}
	return st
}
