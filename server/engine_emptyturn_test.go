package server

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

// 무진행 턴(생각만 있고 본문과 도구가 없음)의 식별과 실행 재개는 steerHooks.Stop 참조.

func assistantThinking(text string) llm.Message {
	return llm.Message{Role: llm.RoleAssistant, Content: []llm.ContentBlock{
		{Type: llm.BlockThinking, Thinking: text, Signature: "sig"},
	}}
}

func TestIsThinkingOnlyTurn(t *testing.T) {
	toolUse := llm.Message{Role: llm.RoleAssistant, Content: []llm.ContentBlock{
		{Type: llm.BlockThinking, Thinking: "先扫端口"},
		{Type: llm.BlockToolUse, ID: "t1", Name: "run_nuclei"},
	}}
	cases := []struct {
		name string
		msgs []llm.Message
		want bool
	}{
		{"생각만 있음", []llm.Message{llm.UserText("开始"), assistantThinking("想想")}, true},
		{"생각+도구", []llm.Message{llm.UserText("开始"), toolUse}, false},
		{"생각+본문", []llm.Message{assistantThinking("想想"), {
			Role:    llm.RoleAssistant,
			Content: []llm.ContentBlock{{Type: llm.BlockThinking, Thinking: "x"}, llm.TextBlock("结论")},
		}}, false},
		{"본문에 공백 문자만 있음", []llm.Message{{
			Role:    llm.RoleAssistant,
			Content: []llm.ContentBlock{{Type: llm.BlockThinking, Thinking: "x"}, llm.TextBlock("  \n ")},
		}}, true},
		{"완전히 빈 assistant 턴", []llm.Message{{Role: llm.RoleAssistant}}, true},
		// 도구 결과의 역할은 user이므로 반드시 앞에 있는 assistant까지 거슬러 올라가 판정하고, 가장 가까운 메시지로 잘못 판정해서는 안 됩니다.
		{"마지막 메시지가 도구 결과", []llm.Message{toolUse, {
			Role:    llm.RoleUser,
			Content: []llm.ContentBlock{{Type: llm.BlockToolResult, ToolUseID: "t1"}},
		}}, false},
		{"assistant 메시지 없음", []llm.Message{llm.UserText("开始")}, false},
		{"빈 이력", nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isThinkingOnlyTurn(c.msgs); got != c.want {
				t.Fatalf("isThinkingOnlyTurn = %v, want %v", got, c.want)
			}
		})
	}
}

// fakeHooks는 프로그래밍 가능한 inner HookRunner로, steerHooks가 inner의 결정을 존중하는지 검증합니다.
type fakeHooks struct {
	prevent  bool
	blocking []string
	msg      string
}

func (f fakeHooks) PreToolUse(context.Context, string, []byte) (bool, string, []byte) {
	return false, "", nil
}
func (f fakeHooks) PostToolUse(context.Context, string, []byte, []byte, bool) {}
func (f fakeHooks) Stop(context.Context, []llm.Message) (bool, []string, string) {
	return f.prevent, f.blocking, f.msg
}

func TestSteerHooksStopNudgesEmptyTurn(t *testing.T) {
	empty := []llm.Message{assistantThinking("我应该先枚举子域名")}

	t.Run("무진행 턴에 실행 재개 지시 추가", func(t *testing.T) {
		h := steerHooks{nudges: &atomic.Int64{}, limit: defaultEmptyTurnNudges, label: "worker-1 · #1"}
		prevent, blocking, _ := h.Stop(context.Background(), empty)
		if prevent {
			t.Fatal("무진행 턴에서는 강제로 중지하면 안 됩니다")
		}
		if len(blocking) != 1 || blocking[0] != emptyTurnNudge {
			t.Fatalf("blocking = %v, want [emptyTurnNudge]", blocking)
		}
	})

	t.Run("본문이나 도구가 있으면 개입하지 않음", func(t *testing.T) {
		h := steerHooks{nudges: &atomic.Int64{}, limit: defaultEmptyTurnNudges}
		normal := []llm.Message{{
			Role:    llm.RoleAssistant,
			Content: []llm.ContentBlock{llm.TextBlock("已完成扫描，未发现开放端口")},
		}}
		if _, blocking, _ := h.Stop(context.Background(), normal); blocking != nil {
			t.Fatalf("정상적인 마무리를 진행 없음으로 잘못 판정: %v", blocking)
		}
		if n := h.nudges.Load(); n != 0 {
			t.Fatalf("개입하지 않았을 때는 횟수를 세면 안 됩니다, got %d", n)
		}
	})

	t.Run("상한에 도달하면 마무리를 허용", func(t *testing.T) {
		const limit = 5 // 사용자가 '빈 응답 재시도 횟수'를 5로 설정합니다.
		h := steerHooks{nudges: &atomic.Int64{}, limit: limit}
		for i := 1; i <= limit; i++ {
			if _, blocking, _ := h.Stop(context.Background(), empty); len(blocking) != 1 {
				t.Fatalf("%d번째는 아직 할당량 이내여야 합니다, blocking = %v", i, blocking)
			}
		}
		if _, blocking, _ := h.Stop(context.Background(), empty); blocking != nil {
			t.Fatalf("상한을 넘었는데도 추가하고 있습니다: %v", blocking)
		}
	})

	// '빈 응답 재시도 횟수'를 -1로 설정하면 이 계층을 끄며, emptyTurnNudgeLimit는 0으로 해석됩니다.
	t.Run("설정이 비활성화되면 개입하지 않음", func(t *testing.T) {
		h := steerHooks{nudges: &atomic.Int64{}, limit: 0}
		if _, blocking, _ := h.Stop(context.Background(), empty); blocking != nil {
			t.Fatalf("비활성화했는데도 추가하고 있습니다: %v", blocking)
		}
	})

	t.Run("inner가 강제 중지를 결정하면 추가하지 않음", func(t *testing.T) {
		h := steerHooks{inner: fakeHooks{prevent: true, msg: "guard 拒绝收场"}, nudges: &atomic.Int64{}, limit: defaultEmptyTurnNudges}
		prevent, blocking, msg := h.Stop(context.Background(), empty)
		if !prevent || msg != "guard 拒绝收场" || blocking != nil {
			t.Fatalf("inner의 강제 중지가 바뀌었습니다: prevent=%v blocking=%v msg=%q", prevent, blocking, msg)
		}
		if n := h.nudges.Load(); n != 0 {
			t.Fatalf("inner에 맡길 때는 할당량을 소비하면 안 됩니다, got %d", n)
		}
	})

	t.Run("inner가 이미 계속 실행을 요구하면 추가하지 않음", func(t *testing.T) {
		h := steerHooks{inner: fakeHooks{blocking: []string{"guard 的续跑理由"}}, nudges: &atomic.Int64{}, limit: defaultEmptyTurnNudges}
		_, blocking, _ := h.Stop(context.Background(), empty)
		if len(blocking) != 1 || blocking[0] != "guard 的续跑理由" {
			t.Fatalf("inner의 계속 실행 메시지가 바뀌었습니다: %v", blocking)
		}
	})

	t.Run("카운터가 없으면 동작을 바꾸지 않음", func(t *testing.T) {
		h := steerHooks{limit: defaultEmptyTurnNudges} // 예를 들어 이후 다른 호출 지점에서 nudges 전달을 빠뜨리는 경우입니다.
		if _, blocking, _ := h.Stop(context.Background(), empty); blocking != nil {
			t.Fatalf("카운터가 없을 때는 추가하면 안 됩니다: %v", blocking)
		}
	})
}
