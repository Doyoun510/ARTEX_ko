package notify

import (
	"net/url"
	"testing"
	"time"
)

// 서명 기준값은 이 패키지의 구현으로 생성하지 않고 OpenSSL로 독립적으로 계산했습니다.
// 그렇지 않으면 "코드가 바뀌지 않음"만 증명할 수 있고, "알고리즘이 맞음"은 증명할 수 없습니다.
//
//	TS=1700000000000, SECRET=SECtest123
//	DingTalk: printf '%s\n%s' "$TS" "$SECRET" | openssl dgst -sha256 -hmac "$SECRET" -binary | openssl base64 -A
//	      -> w3RMHXzixTMdzr8OHJUmVLS4IoPJVdu+Ut1LE48MePE=
//	Feishu: printf '' | openssl dgst -sha256 -hmac "$(printf '%s\n%s' "$TS" "$SECRET")" -binary | openssl base64 -A
//	      -> Hd4xFWQU6R6ad4nzy4ETIznzlqebqH7xcTFVmONTudo=
const (
	signTestTSMillis = int64(1700000000000)
	signTestSecret   = "SECtest123"
	dingTalkExpected = "w3RMHXzixTMdzr8OHJUmVLS4IoPJVdu+Ut1LE48MePE="
	feishuExpected   = "Hd4xFWQU6R6ad4nzy4ETIznzlqebqH7xcTFVmONTudo="
)

func TestDingTalkSignMatchesReference(t *testing.T) {
	got, err := dingTalkSignedURL("https://oapi.dingtalk.com/robot/send?access_token=tok", signTestSecret, time.UnixMilli(signTestTSMillis))
	if err != nil {
		t.Fatalf("서명 실패: %v", err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("생성된 주소를 파싱할 수 없습니다: %v", err)
	}
	q := u.Query()
	if q.Get("sign") != dingTalkExpected {
		t.Errorf("서명 불일치\n기대 %s\n실제 %s", dingTalkExpected, q.Get("sign"))
	}
	if q.Get("timestamp") != "1700000000000" {
		t.Errorf("타임스탬프는 밀리초이며 그대로 포함해야 합니다. 실제: %q", q.Get("timestamp"))
	}
	// 기존 query 파라미터(access_token)를 서명으로 덮어써서는 안 됩니다.
	if q.Get("access_token") != "tok" {
		t.Errorf("기존 query 파라미터 누락. 실제: %q", q.Get("access_token"))
	}
}

func TestFeishuSignMatchesReference(t *testing.T) {
	got := feishuSign("1700000000000", signTestSecret)
	if got != feishuExpected {
		t.Errorf("서명 불일치\n기대 %s\n실제 %s", feishuExpected, got)
	}
}

// TestSignAlgorithmsDiffer는 두 알고리즘의 차이를 고정합니다. 두 알고리즘의 파라미터 순서는 정확히 반대이며,
// (DingTalk key=secret, Feishu key=서명할 문자열), 다른 쪽을 그대로 따라 쓰면 반드시 검증에 실패하므로,
// 이 테스트는 향후 리팩터링에서 두 구현을 같은 함수로 합치지 않도록 보장합니다.
func TestSignAlgorithmsDiffer(t *testing.T) {
	ts := "1700000000000"
	dingURL, err := dingTalkSignedURL("https://example.com/hook", signTestSecret, time.UnixMilli(signTestTSMillis))
	if err != nil {
		t.Fatal(err)
	}
	dq, _ := url.Parse(dingURL)
	if dq.Query().Get("sign") == feishuSign(ts, signTestSecret) {
		t.Fatal("DingTalk과 Feishu 서명이 같으므로 한쪽의 알고리즘 구현이 잘못되었습니다")
	}
}

func TestDingTalkNoSecretLeavesURLUntouched(t *testing.T) {
	// 서명 추가를 켜지 않은 봇에는 timestamp/sign 파라미터를 임의로 추가해서는 안 됩니다.
	const hook = "https://oapi.dingtalk.com/robot/send?access_token=tok"
	got, err := dingTalkSignedURL(hook, "", time.UnixMilli(signTestTSMillis))
	if err != nil {
		t.Fatal(err)
	}
	if got != hook {
		t.Fatalf("secret을 설정하지 않으면 주소를 변경해서는 안 됩니다. 실제: %q", got)
	}
}

func TestValidateHTTPURL(t *testing.T) {
	ok := []string{"https://example.com/hook", "http://10.0.0.1:8080/x?y=1"}
	for _, s := range ok {
		if err := validateHTTPURL(s); err != nil {
			t.Errorf("%q: 허용되어야 합니다: %v", s, err)
		}
	}
	// file:// 등은 허용해서는 안 됩니다. http.Client의 처리가 예상 범위를 벗어납니다.
	bad := []string{"", "file:///etc/passwd", "ftp://example.com", "https://", "gopher://x"}
	for _, s := range bad {
		if err := validateHTTPURL(s); err == nil {
			t.Errorf("%q: 거부되어야 합니다", s)
		}
	}
}
