package notify

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// 이 파일은 불변 조건 테스트입니다. 채널 구현에서 나오는 **어떤** 오류 텍스트에도 자격 증명이 포함되어서는 안 됩니다.
//
// 별도 파일로 분리하는 이유: 처음 채널 테스트는 성공 경로와 플랫폼 업무 오류만 검사했으며,
// 전송 계층 실패는 전혀 확인하지 않았습니다. 그러나 전송 계층 오류(연결 거부/DNS 실패/타임아웃)가 가장 위험합니다.
// http.Client.Do가 반환하는 *url.Error는 **전체 URL**을 오류 텍스트에 넣으며, 이 기능에서 사용하는
// 플랫폼들의 자격 증명은 URL 안에 있습니다. 자격 증명은 이 문자열을 통해 네 곳으로 전달됩니다.
//
//	notification_deliveries.last_error  → 평문으로 DB 저장
//	GET /api/notify/deliveries 응답     → 채널 설정 마스킹을 우회하여 브라우저에 반환
//	서버 로그                          → 외부로 보내 보관하는 경우가 많음
//	테스트 전송 API의 502 응답             → 프런트에 바로 표시
//
// 따라서 여기서는 함수 하나만 검사하지 않고, 각 채널에서 반드시 실패하는 요청을 한 번씩 실제로 보내 오류 텍스트에
// 해당 자격 증명이 없는지 검사합니다.

// credentialCases는 "URL에 자격 증명이 있음"인 모든 채널 형태를 검사합니다.
// DingTalk/WeCom(기업용 위챗)은 query, Feishu는 경로 끝부분, Telegram은 경로 중간에 있습니다.
var credentialCases = []struct {
	name   string
	ch     Channel
	cfg    map[string]any
	secret string
}{
	{
		name:   "钉钉 access_token 在 query",
		ch:     dingTalkChannel{},
		cfg:    map[string]any{"webhook": "http://127.0.0.1:1/robot/send?access_token=" + leakProbeToken},
		secret: leakProbeToken,
	},
	{
		name:   "企业微信 key 在 query",
		ch:     weComChannel{},
		cfg:    map[string]any{"webhook": "http://127.0.0.1:1/cgi-bin/webhook/send?key=" + leakProbeToken},
		secret: leakProbeToken,
	},
	{
		name:   "飞书 hook id 在路径末段",
		ch:     feishuChannel{},
		cfg:    map[string]any{"webhook": "http://127.0.0.1:1/open-apis/bot/v2/hook/" + leakProbeToken},
		secret: leakProbeToken,
	},
	{
		name:   "Telegram bot token 在路径中段",
		ch:     telegramChannel{},
		cfg:    map[string]any{"bot_token": leakProbeToken, "chat_id": "1", "base_url": "http://127.0.0.1:1"},
		secret: leakProbeToken,
	},
	{
		name:   "钉钉加签密钥",
		ch:     dingTalkChannel{},
		cfg:    map[string]any{"webhook": "http://127.0.0.1:1/robot/send", "secret": leakProbeToken},
		secret: leakProbeToken,
	},
}

// leakProbeToken은 실제 자격 증명일 가능성이 전혀 없는 센티널 값이며, 오류 텍스트에서 검색하는 데 사용합니다.
const leakProbeToken = "LEAKPROBE0123456789abcdef"

// TestChannelErrorsNeverLeakCredentials는 핵심 불변 조건입니다.
func TestChannelErrorsNeverLeakCredentials(t *testing.T) {
	for _, tc := range credentialCases {
		t.Run(tc.name, func(t *testing.T) {
			// 반드시 실패하는 상대 측: 127.0.0.1:1에는 수신 대기하는 서비스가 없어 연결 거부 경로를 거칩니다.
			_, err := tc.ch.Send(context.Background(), tc.cfg, Message{
				Items: []Item{{FindingID: 1, Severity: "high", Name: "泄露探针"}},
			})
			if err == nil {
				t.Fatal("접근할 수 없는 주소에는 오류가 나야 합니다")
			}
			assertNoSecret(t, err.Error(), tc.secret)
		})
	}
}

// TestChannelErrorsNeverLeakCredentialsInPermanentPath는 영구 실패 분기를 검사합니다.
// URL 검증 실패나 플랫폼 업무 오류도 오류 텍스트를 외부로 보내므로, 마찬가지로 자격 증명이 포함되어서는 안 됩니다.
func TestChannelErrorsNeverLeakCredentialsInPermanentPath(t *testing.T) {
	cases := []struct {
		name string
		ch   Channel
		cfg  map[string]any
	}{
		// 주소에 자격 증명이 있지만 형식이 유효하지 않음 → validateHTTPURL / url.Parse 분기 실행.
		{"钉钉地址非法", dingTalkChannel{}, map[string]any{"webhook": "file:///" + leakProbeToken}},
		{"企微地址非法", weComChannel{}, map[string]any{"webhook": "gopher://" + leakProbeToken}},
		{"飞书地址非法", feishuChannel{}, map[string]any{"webhook": "ftp://" + leakProbeToken + "/hook"}},
		{"Telegram API 地址非法", telegramChannel{}, map[string]any{"bot_token": "tok", "chat_id": "1", "base_url": "file://" + leakProbeToken}},
		{"通用 Webhook 地址非法", webhookChannel{}, map[string]any{"url": "javascript:" + leakProbeToken}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.ch.Send(context.Background(), tc.cfg, Message{Items: []Item{{Severity: "high"}}})
			if err == nil {
				t.Fatal("유효하지 않은 설정에는 오류가 나야 합니다")
			}
			assertNoSecret(t, err.Error(), leakProbeToken)
		})
	}
}

func assertNoSecret(t *testing.T, text, secret string) {
	t.Helper()
	if strings.Contains(text, secret) {
		t.Fatalf("오류 텍스트에 자격 증명 %q 노출:\n    %s", secret, text)
	}
}

func TestRedactRequestTargetKeepsOnlySchemeAndHost(t *testing.T) {
	cases := map[string]string{
		"https://oapi.dingtalk.com/robot/send?access_token=S1":    "https://oapi.dingtalk.com/…",
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=S2": "https://qyapi.weixin.qq.com/…",
		"https://open.feishu.cn/open-apis/bot/v2/hook/S3":         "https://open.feishu.cn/…",
		"https://api.telegram.org/botS4/sendMessage":              "https://api.telegram.org/…",
		"http://10.0.0.5:8080/hook":                               "http://10.0.0.5:8080/…",
	}
	for in, want := range cases {
		got := redactRequestTarget(in)
		if got != want {
			t.Errorf("redactRequestTarget(%q) = %q, 기대 %q", in, got, want)
		}
		// 민감정보 제거 결과에도 원본 주소의 경로/쿼리 조각이 포함되어서는 안 됩니다.
		if parts := strings.SplitN(in, "://", 2); len(parts) == 2 {
			if hostAndRest := strings.SplitN(parts[1], "/", 2); len(hostAndRest) == 2 && hostAndRest[1] != "" {
				if strings.Contains(got, hostAndRest[1]) {
					t.Errorf("민감정보 제거 후에도 경로/쿼리 조각 %q 포함: %q", hostAndRest[1], got)
				}
			}
		}
	}
	// 파싱할 수 없는 입력의 원문은 절대로 반환하지 않습니다.
	for _, bad := range []string{"", "://", "not a url", "http://"} {
		if got := redactRequestTarget(bad); strings.Contains(got, bad) && bad != "" {
			t.Errorf("파싱할 수 없는 입력 %q 반환됨: %q", bad, got)
		}
	}
}

// TestRedactTransportErrorStripsURL은 구체적인 타입 *url.Error를 직접 검사합니다.
// 이는 http.Client.Do의 반환 타입이며, 유출이 처음 발생하는 위치입니다.
func TestRedactTransportErrorStripsURL(t *testing.T) {
	inner := errors.New("dial tcp 127.0.0.1:1: connect: connection refused")
	uerr := &url.Error{
		Op:  "Post",
		URL: "https://api.telegram.org/bot" + leakProbeToken + "/sendMessage",
		Err: inner,
	}
	got := redactTransportError(uerr)
	assertNoSecret(t, got, leakProbeToken)
	if !strings.Contains(got, "api.telegram.org") {
		t.Errorf("원인 파악을 위해 host를 유지해야 합니다. 실제: %q", got)
	}
	if !strings.Contains(got, "connection refused") {
		t.Errorf("원인 파악을 위해 원래 원인을 유지해야 합니다. 실제: %q", got)
	}
	// Op도 유지해야 합니다(POST인지 GET인지가 원인 파악에 의미 있음).
	if !strings.Contains(got, "Post") {
		t.Errorf("동작 이름을 유지해야 합니다. 실제: %q", got)
	}
}

// TestRedactURLsInTextHandlesFallback 대체 처리 경로: *url.Error가 아닌 사용자 지정 오류
// (리디렉션 정책이 반환하는 오류 등)의 주소도 제거해야 합니다.
func TestRedactURLsInTextHandlesFallback(t *testing.T) {
	in := fmt.Sprintf("拒绝跨主机重定向（a.example → http://b.example/bot%s/send）", leakProbeToken)
	got := redactURLsInText(in)
	assertNoSecret(t, got, leakProbeToken)
	if !strings.Contains(got, "http://b.example/…") {
		t.Errorf("주소를 민감정보가 제거된 형태로 바꿔야 합니다. 받은 값: %q", got)
	}
	// 주소가 없는 텍스트는 그대로 유지합니다.
	if plain := "dial tcp: connection refused"; redactURLsInText(plain) != plain {
		t.Error("주소가 없는 텍스트를 변경해서는 안 됩니다")
	}
}

// TestCrossHostRedirectRefused는 "URL의 자격 증명 + 다른 호스트로 이동 = 자격 증명 전달"을 검사합니다.
// httptest의 두 서비스는 127.0.0.1의 서로 다른 포트에서 수신 대기합니다. 포트가 다르면 Host가 다르므로,
// 다른 호스트로의 이동이 됩니다.
func TestCrossHostRedirectRefused(t *testing.T) {
	var hit bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit = true
		_, _ = io.WriteString(w, `{"errcode":0,"errmsg":"ok"}`)
	}))
	defer target.Close()

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/robot/send?access_token="+leakProbeToken, http.StatusTemporaryRedirect)
	}))
	defer redirector.Close()

	_, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": redirector.URL + "/robot/send?access_token=" + leakProbeToken},
		Message{Items: []Item{{Severity: "high"}}})
	if err == nil {
		t.Fatal("다른 호스트로의 리디렉션은 거부되어야 합니다")
	}
	if hit {
		t.Fatal("이동 대상에 접근했습니다. 리디렉션으로 자격 증명이 유출되었습니다")
	}
	assertNoSecret(t, err.Error(), leakProbeToken)
}

// TestSameHostRedirectAllowed는 반대 테스트입니다. 같은 호스트로의 이동(끝에 슬래시를 붙이는 경우 등)은 반드시 계속 사용할 수 있어야 하며,
// 그렇지 않으면 정상적인 작업 흐름까지 차단하게 됩니다.
func TestSameHostRedirectAllowed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robot/send" {
			// 같은 호스트·같은 포트로 이동합니다.
			http.Redirect(w, r, "/robot/send/", http.StatusTemporaryRedirect)
			return
		}
		_, _ = io.WriteString(w, `{"errcode":0,"errmsg":"ok"}`)
	}))
	defer srv.Close()

	if _, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": srv.URL + "/robot/send"},
		Message{Items: []Item{{Severity: "high"}}}); err != nil {
		t.Fatalf("같은 호스트로의 리디렉션은 거부되어서는 안 됩니다: %v", err)
	}
}
