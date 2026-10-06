package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// 이 파일은 doJSON의 HTTP 계층 오류 분류를 검사합니다.
//
// 따로 검사하는 이유: 각 채널 어댑터는 플랫폼 자체의 업무 오류 코드(DingTalk errcode,
// Feishu code, Telegram ok 필드)만 담당하며, **HTTP 계층** 분류는 doJSON이 통일하여 수행합니다.
// 두 분류는 독립된 방어선입니다. 이것이 없으면 503을 반환하는 중계 게이트웨이를 영구 실패로 보고,
// 재시도를 바로 포기합니다. 반대로 403을 재시도 가능으로 보고 불필요하게 세 번 백오프하게 됩니다.

func replyServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestDoJSONClassifiesHTTPStatus(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		permanent bool
	}{
		{"200 成功不算错误", 200, false},
		{"429 限流可重试", 429, false},
		{"408 请求超时可重试", 408, false},
		{"500 服务端错误可重试", 500, false},
		{"502 网关错误可重试", 502, false},
		{"503 服务不可用可重试", 503, false},
		{"400 参数错误永久失败", 400, true},
		{"401 鉴权失败永久失败", 401, true},
		{"403 禁止访问永久失败", 403, true},
		{"404 地址不存在永久失败", 404, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := replyServer(t, tc.status, `{"detail":"upstream says no"}`)
			_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
			if tc.status < 300 {
				if err != nil {
					t.Fatalf("2xx에는 오류가 나서는 안 됩니다: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("2xx가 아니면 오류가 나야 합니다")
			}
			if got := IsPermanent(err); got != tc.permanent {
				t.Fatalf("HTTP %d permanent 판정 오류: 기대 %v 실제 %v (%v)",
					tc.status, tc.permanent, got, err)
			}
			// 상태 코드는 반드시 오류에 나타나야 합니다. 그렇지 않으면 사용자는 설정 오류인지 상대 측 장애인지 판단할 수 없습니다.
			// Go의 영어 StatusText 대신 숫자를 검사합니다. 이 패키지의 문구는 한국어이며,
			// (프로젝트의 다른 부분과 일치), 숫자가 언어에 의존하지 않고 안정적으로 검사할 수 있는 부분입니다.
			if !strings.Contains(err.Error(), strconv.Itoa(tc.status)) {
				t.Errorf("오류 메시지에 HTTP 상태 코드 %d 포함 필요. 실제: %v", tc.status, err)
			}
		})
	}
}

// TestDoJSONIncludesResponseSnippet은 snippet을 검사합니다. 상대 측이 반환한 오류 설명을 포함해야 하며,
// 그렇지 않으면 사용자는 "실패했다"는 것만 알고 상대 측이 왜 거부했는지는 모릅니다.
func TestDoJSONIncludesResponseSnippet(t *testing.T) {
	srv := replyServer(t, 400, `{"error":"invalid webhook token"}`)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("오류가 나야 합니다")
	}
	if !strings.Contains(err.Error(), "invalid webhook token") {
		t.Errorf("오류 메시지에 상대 측 설명을 포함해야 합니다. 실제: %v", err)
	}
}

// TestDoJSONSnippetIsSingleLineAndBounded는 snippet의 형태를 제한합니다.
// 상대 측 응답은 last_error 열과 프런트 표에 그대로 들어가므로, 여러 줄/긴 내용은 레이아웃과 데이터 크기를 망가뜨립니다.
func TestDoJSONSnippetIsSingleLineAndBounded(t *testing.T) {
	// 줄바꿈·탭과 5000문자의 긴 내용을 포함한 응답입니다.
	long := strings.Repeat("x", 5000)
	srv := replyServer(t, 500, "line1\nline2\r\n\tline3 "+long)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("오류가 나야 합니다")
	}
	msg := err.Error()
	if strings.ContainsAny(msg, "\r\n\t") {
		t.Errorf("오류 메시지는 한 줄로 압축해야 합니다. 실제: %q", msg)
	}
	// snippet 상한 200문자 + 고정 접두사이며, 전체는 원래 응답보다 훨씬 작아야 합니다.
	if len(msg) > 400 {
		t.Errorf("오류 메시지가 너무 깁니다(%d바이트). snippet으로 잘라야 합니다: %q", len(msg), msg)
	}
}

// TestDoJSONRejectsOversizedResponse는 읽기 상한을 확인합니다. 상대 측이 매우 큰 내용을 반환할 때,
// 응답 전체를 메모리에 읽어서는 안 됩니다(전송 이력의 각 항목마다 last_error를 저장함).
func TestDoJSONRejectsOversizedResponse(t *testing.T) {
	huge := strings.Repeat("A", 1<<20) // 1 MiB
	srv := replyServer(t, 400, huge)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("오류가 나야 합니다")
	}
	if len(err.Error()) > 400 {
		t.Errorf("큰 응답은 읽기 길이를 제한하고 잘라야 합니다. 오류 메시지 길이: %d", len(err.Error()))
	}
}
