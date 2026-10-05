package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Autumn-27/norma/harness"
)

// runTrace retains the latest tool call so an interrupted run can identify the
// operation that was still in flight.
type runTrace struct {
	startedAt time.Time
	id        string
	name      string
	input     string
	at        time.Time
	pending   bool
}

func (t *runTrace) start(id, name, input string) {
	t.id, t.name, t.input, t.at, t.pending = id, name, input, time.Now(), true
}

func (t *runTrace) done(id string) {
	if id == t.id {
		t.pending = false
	}
}

var reasonHint = map[harness.TerminalReason]string{
	harness.ReasonCompleted:         "모델이 이번 실행을 정상적으로 종료했지만 텍스트 요약을 남기지 않았습니다. 사실과 자산은 이번 실행의 도구 호출 기록을 기준으로 합니다",
	harness.ReasonMaxTurns:          "스텝 수 상한(MaxTurns)에 도달했습니다: SDK가 마무리를 실행하고 사실과 자산을 기록했습니다. 의도는 exhausted로 표시되어 planner가 방향을 바꿔 계속하도록 하며, 실패로 처리하지 않습니다",
	harness.ReasonTimeout:           "단일 실행의 실제 경과 시간 예산(MaxDuration)에 도달했습니다: 시간이 되면 실행 중인 도구를 중단하고 그 자리에서 마무리하여 식별한 사실과 자산을 기록하며, 의도는 exhausted로 표시됩니다",
	harness.ReasonModelError:        "모델 또는 API 호출에 실패했습니다(네트워크, 인증·인가 확인, 요청 제한, 공급자 5xx 등). 재시도를 소진하면 의도는 blocked로 표시됩니다. 전송 계층 장애로 이 의도는 사실상 제대로 탐색되지 않았습니다. 실행 과정(get_worker_trace)을 확인한 뒤 다시 맡기거나 방법을 바꿀지 결정하세요",
	harness.ReasonBlockingLimit:     "컨텍스트 길이가 하드 상한에 도달하여 요청을 보내기 전에 차단되었습니다. 의도 단위를 더 작게 조정하거나 도구 반환 내용을 압축해야 합니다",
	harness.ReasonPromptTooLong:     "프롬프트가 너무 길고 컨텍스트 압축 재시도도 소진되어 실행을 계속할 수 없습니다",
	harness.ReasonImageError:        "현재 모델은 이번 실행의 멀티모달 내용을 지원하지 않습니다. 시각 입력을 지원하는 모델로 전환하거나 도구가 이미지를 반환하지 않도록 하세요",
	harness.ReasonStopHookPrevented: "Stop 훅이 이번 실행의 종료를 막았으며 이후 계속하지 못했습니다. 작업의 Guard 규칙이 지나치게 엄격한지 확인하세요",
	harness.ReasonHookStopped:       "도구 또는 훅이 실행 지속을 직접 중지했습니다. 예를 들면 범위 밖 대상이나 비활성화된 명령입니다. 마지막 tool_result의 차단 설명을 확인하세요",
	harness.ReasonAbortedStreaming:  "모델 출력의 스트리밍 생성 단계에서 실행이 취소되었습니다",
	harness.ReasonAbortedTools:      "도구 실행 단계에서 실행이 취소되었습니다",
}

// terminalText renders a terminal event with no final text into a compact summary
// and a Markdown detail block.
func terminalText(ctx context.Context, term *harness.Terminal, tr *runTrace) (string, string) {
	reason := term.Reason
	aborted := reason == harness.ReasonAbortedStreaming || reason == harness.ReasonAbortedTools
	// Prompt may return ctx.Err directly without a terminal event. Preserve the
	// cancellation cause instead of falling back to an empty/unknown terminal reason.
	if reason == "" && ctx.Err() != nil {
		aborted = true
	}

	var sum string
	if aborted {
		_, short, _, ok := AbortReason(ctx)
		if !ok {
			short = "취소 사유를 가져오지 못했습니다"
		}
		stage := "실행 도중"
		switch reason {
		case harness.ReasonAbortedStreaming:
			stage = "모델 출력 단계"
		case harness.ReasonAbortedTools:
			stage = "도구 실행 단계"
		}
		sum = "(실행 중단: " + short + "; 중단 위치: " + stage + progressSuffix(term, tr) + ", 미완료)"
	} else if reason == harness.ReasonMaxTurns || reason == harness.ReasonTimeout {
		sum = "(실행 예산 상한 도달(" + string(reason) + "), 마무리 후 사실 기록 완료" + progressSuffix(term, tr) + "; 이번 실행에는 텍스트 요약이 없습니다)"
	} else {
		hint := terminalReasonHint(reason)
		sum = "(텍스트 요약 없음, 종료 상태 " + terminalReasonLabel(reason) + ": " + firstLine(hint, 80) + ")"
	}

	var b strings.Builder
	b.WriteString(sum)
	b.WriteString("\n\n")
	displayReason := terminalReasonLabel(reason)
	fmt.Fprintf(&b, "- **종료 상태**: `%s` - %s\n", displayReason, terminalReasonHint(reason))
	if aborted {
		code, _, why, ok := AbortReason(ctx)
		if ok {
			fmt.Fprintf(&b, "- **중단 사유** (`%s`): %s\n", code, why)
		} else {
			b.WriteString("- **중단 사유**: 가져올 수 없습니다. 취소한 쪽에서 context.WithCancelCause로 코드가 지정된 취소 사유를 첨부하지 않았을 수 있습니다\n")
		}
	}
	if term.Err != nil {
		fmt.Fprintf(&b, "- **하위 계층 오류**: `%v`\n", term.Err)
	}
	if aborted && strings.TrimSpace(term.Text) != "" {
		b.WriteString("- **취소 전에 생성된 일부 출력**:\n\n")
		b.WriteString(term.Text)
		b.WriteString("\n\n")
	}
	if term.Turns > 0 {
		fmt.Fprintf(&b, "- **실행한 모델 턴**: %d회\n", term.Turns)
	}
	if !tr.startedAt.IsZero() {
		fmt.Fprintf(&b, "- **이번 실행 소요 시간**: %s\n", roundDur(time.Since(tr.startedAt)))
	}
	if u := term.Usage; u.InputTokens+u.OutputTokens+u.CacheReadTokens+u.CacheWriteTokens > 0 {
		fmt.Fprintf(&b, "- **누적 token**: 입력 %d / 출력 %d / 캐시 읽기 %d / 캐시 쓰기 %d\n",
			u.InputTokens, u.OutputTokens, u.CacheReadTokens, u.CacheWriteTokens)
	}
	if tr.name == "" {
		b.WriteString("- **도구 호출**: 이번 실행은 도구 호출을 보내기 전에 종료되었습니다\n")
	} else if tr.pending {
		fmt.Fprintf(&b, "- **중단 시 실행 중이던 도구**: `%s`(실행 시간 %s, **결과 미반환**)\n\n  ```json\n  %s\n  ```\n",
			tr.name, roundDur(time.Since(tr.at)), firstLine(tr.input, 300))
	} else {
		fmt.Fprintf(&b, "- **중단 전 마지막 도구**: `%s`(정상 반환 완료)\n", tr.name)
	}
	return sum, b.String()
}

func terminalReasonLabel(reason harness.TerminalReason) string {
	if reason == "" {
		return "context_canceled"
	}
	return string(reason)
}

func terminalReasonHint(reason harness.TerminalReason) string {
	if hint := reasonHint[reason]; hint != "" {
		return hint
	}
	if reason == "" {
		return "실행의 context가 취소되었지만 하위 계층에서 Terminal 이벤트를 생성하지 않았습니다"
	}
	return "알 수 없는 종료 상태입니다. harness에 TerminalReason이 추가되었을 수 있으니 reasonHint를 보완하세요"
}

func progressSuffix(term *harness.Terminal, tr *runTrace) string {
	var parts []string
	if term.Turns > 0 {
		parts = append(parts, fmt.Sprintf("%d회", term.Turns))
	}
	if !tr.startedAt.IsZero() {
		parts = append(parts, roundDur(time.Since(tr.startedAt)))
	}
	if len(parts) == 0 {
		return ""
	}
	return ", 실행 경과: " + strings.Join(parts, " / ")
}

func roundDur(d time.Duration) string {
	switch {
	case d < time.Minute:
		return d.Round(100 * time.Millisecond).String()
	case d < time.Hour:
		return d.Round(time.Second).String()
	default:
		return d.Round(time.Minute).String()
	}
}
