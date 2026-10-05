package agent

import (
	"context"
	"errors"
	"fmt"
)

// AbortCause names why an agent run's context was cancelled. Every cancellation
// site should attach one so the activity trace can report the real initiator.
type AbortCause struct {
	Code  string
	Short string
	Text  string
}

func (c *AbortCause) Error() string { return c.Text }

func cause(code, short, text string) *AbortCause {
	return &AbortCause{Code: code, Short: short, Text: text}
}

// Causef builds a cause that includes runtime-specific detail.
func Causef(code, short, format string, args ...any) *AbortCause {
	return &AbortCause{Code: code, Short: short, Text: fmt.Sprintf(format, args...)}
}

var (
	// Task-level execution context.
	AbortPausedByUser = cause("paused_by_user", "사용자가 작업을 일시 중지했습니다",
		"사용자가 작업 제어 인터페이스(POST /api/tasks/{id}/control, action=pause)를 통해 작업을 일시 중지했습니다. 이번 Planner/Worker 실행은 의도적으로 취소됩니다. 실행 중인 의도는 frontier(open)로 돌아가며, 작업 재개 후 다시 가져와 처음부터 실행합니다")
	AbortPausedByOrchestrator = cause("paused_by_orchestrator", "오케스트레이션 Agent가 작업을 일시 중지했습니다",
		"오케스트레이션 Agent가 pause_task 도구를 호출해 이 작업을 일시 중지했습니다. 이번 Planner/Worker 실행은 의도적으로 취소됩니다. 실행 중인 의도는 frontier(open)로 돌아가며, 재개 후 다시 실행합니다")
	AbortTaskDeleted = cause("task_deleted", "작업이 삭제되었습니다",
		"작업을 삭제하는 중입니다(DELETE /api/tasks/{id}). 삭제 보호 장치가 이 작업에서 실행 중인 Planner, Worker, 메인 agent를 취소했으며, 이번 실행 결과는 더 이상 사용되지 않습니다")
	AbortPausedOnReload = cause("paused_on_reload", "백엔드가 작업의 일시 중지 상태를 복원했습니다",
		"백엔드가 시작할 때 데이터베이스에 저장된 상태에 따라 작업의 일시 중지 상태를 복원했습니다. 이번 실행은 취소됩니다. 정상적인 경우 복원 단계에는 실행 중인 Agent가 없습니다")
	AbortGoalMet = cause("goal_met", "planner가 작업 목표 달성을 판정했습니다",
		"planner가 작업 목표 달성을 판정하고 작업을 done 상태로 설정한 뒤, 아직 실행 중인 Worker를 취소합니다. 해당 의도는 실패가 아니라 stopped로 표시됩니다")
	AbortSettleDrainTimeout = cause("settle_drain_timeout", "작업 시간 초과 후 마무리를 기다리는 시간이 끝났습니다",
		"작업이 timeout에 도달한 후 실행 중인 Worker의 정상적인 마무리를 기다렸지만 90초의 drain 유예 시간도 부족해 강제 취소합니다. 의도는 exhausted로 표시되며, 마무리 단계에서 이미 기록한 사실과 자산은 보존됩니다")

	// Per-work context.
	AbortKilledByPlanner = cause("killed_by_planner", "planner가 이 의도를 종료했습니다",
		"planner가 kill_work를 호출해 이 의도를 직접 종료했습니다. 보통 진행 방향이 어긋났거나 계속할 가치가 없다는 뜻입니다. 의도는 stopped로 표시되며 자동으로 다시 가져오지 않습니다")
	AbortWorkPausedByUser = cause("work_paused_by_user", "사용자가 이 Worker 의도를 일시 중지했습니다",
		"사용자가 실행 중인 Worker를 일시 중지했습니다. 이번 호출은 취소되고 의도는 paused 상태로 전환됩니다. 이미 등록된 의도, 사실, 취약점, 활동 기록은 모두 보존되며, 재개 후 처음부터 다시 실행합니다")
	AbortWorkCancelledByUser = cause("work_cancelled_by_user", "사용자가 이 Worker 의도를 삭제했습니다",
		"사용자가 실행 중인 Worker를 삭제했습니다. 이번 호출은 취소됩니다. Worker가 쓰기 구간을 벗어난 뒤 서버는 사용자가 선택한 삭제 방식에 따라 해당 의도를 처리합니다. 논리 삭제는 삭제됨으로만 표시하고 모든 산출물을 보존하며, 물리 삭제는 해당 의도와 해당 의도로만 뒷받침되는 하위 노드를 연쇄 삭제합니다")
	AbortWorkFinished = cause("work_finished", "Worker가 정상적으로 종료되고 context를 해제했습니다",
		"Worker가 정상적으로 종료되어 엔진이 detachWork에서 context 리소스를 해제합니다. 이는 실행 중단이 아닙니다. 중단 메시지에 이 사유가 나타난다면 취소와 마무리 이벤트 사이에 Race Condition이 발생한 것입니다")
	AbortPausedRaceGuard = cause("paused_race_guard", "작업 일시 중지 중 새 실행의 시작을 거부했습니다",
		"작업이 일시 중지된 동안 엔진은 새 실행 context 생성을 거부합니다. 이는 claim과 일시 중지 사이의 Race Condition으로 Worker가 계속 시작되는 것을 방지하기 위한 것입니다. 이미 가져온 의도는 frontier로 돌아갑니다")

	// Main Agent and standalone conversation contexts.
	AbortChatStoppedByUser = cause("chat_stopped_by_user", "사용자가 이번 대화를 중지했습니다",
		"사용자가 중지를 눌러 이번 메인 agent 또는 대화 Agent 실행을 직접 중단했습니다. 이미 생성된 활동 기록은 보존되며, 다음 메시지를 계속 보낼 수 있습니다")
	AbortChatPausedWithTask = cause("chat_paused_with_task", "작업 일시 중지로 메인 agent 대화가 중단되었습니다",
		"사용자가 작업을 일시 중지할 때 실행 중이던 메인 agent 대화도 함께 취소됩니다. 이미 생성된 활동 기록은 보존됩니다. 작업 재개 후 이번 메시지가 자동으로 다시 재생되지는 않습니다")
	AbortChatTurnFinished = cause("chat_turn_finished", "이번 대화가 정상적으로 종료되고 context를 해제했습니다",
		"이번 대화는 정상적으로 종료되었고, 서버가 해당 턴의 context 리소스를 해제하고 있습니다. 이는 실행 중단이 아닙니다. 중단 메시지에 이 사유가 나타난다면 취소와 마무리 이벤트 사이에 Race Condition이 발생한 것입니다")

	// Process-level and per-run hard backstop.
	AbortShutdown = cause("shutdown", "백엔드 프로세스가 종료 중입니다",
		"백엔드 프로세스가 SIGINT 또는 SIGTERM을 받아 재시작, 업데이트 또는 종료 중입니다. 실행 중인 모든 Agent는 취소됩니다. 재시작 후 남아 있는 running 의도는 open으로 재설정되어 다시 실행됩니다")
	AbortRunHardTimeout = cause("run_hard_timeout", "단일 실행의 하드 타임아웃 최종 보호 조치가 작동했습니다",
		"단일 실행이 소프트 실제 경과 시간 예산과 추가 유예 시간을 초과했습니다. 이는 모델 요청이나 특정 도구가 오랫동안 반환되지 않아 정상적인 턴 경계 마무리를 실행할 수 없음을 의미합니다. 중단 직전에 반환되지 않은 마지막 도구 호출을 중점적으로 확인하세요")
)

// AbortReason resolves the named cause attached to a cancelled run context.
func AbortReason(ctx context.Context) (code, short, text string, ok bool) {
	c := context.Cause(ctx)
	if c == nil {
		return "", "", "", false
	}
	var ac *AbortCause
	if errors.As(c, &ac) {
		return ac.Code, ac.Short, ac.Text, true
	}
	switch {
	case errors.Is(c, context.DeadlineExceeded):
		return "deadline_exceeded", "상위 context가 deadline에 도달했습니다",
			"상위 context가 deadline에 도달했지만, 설정한 쪽에서 WithTimeoutCause를 통해 코드가 지정된 취소 사유를 첨부하지 않았습니다: " + c.Error(), true
	case errors.Is(c, context.Canceled):
		return "canceled_no_cause", "취소한 쪽에서 코드가 지정된 취소 사유를 첨부하지 않았습니다",
			"상위 context가 취소되었지만, 취소한 쪽에서 context.WithCancelCause를 통해 코드가 지정된 취소 사유를 첨부하지 않았습니다. agent/cancelcause.go에 사유를 등록하고 해당 취소 지점에 연결하세요", true
	default:
		return "other", firstLine(c.Error(), 80), c.Error(), true
	}
}
