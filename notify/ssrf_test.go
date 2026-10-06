package notify

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// 이 파일은 두 가지 관련 보안 강화를 검사합니다.
//   ① 전송 주소로 서버를 경유 서버로 삼아 내부망 / 클라우드 메타데이터에 접근해서는 안 됩니다(SSRF).
//   ② 주소 검증 오류 메시지에 주소의 자격 증명이 노출되어서는 안 됨
//
// 테스트 환경: 이 패키지의 많은 테스트가 127.0.0.1의 httptest 가짜 수신처를 사용하며, 연결 보호 검사는 기본적으로 이를
// 차단합니다. 따라서 TestMain에서 AllowLocalTargetsEnv를 일괄 켜고, 아래의 각 SSRF 테스트는
// 명시적으로 이를 해제하여 **기본적으로 거부**하는 동작을 검증합니다.

func TestMain(m *testing.M) {
	// 일반 테스트가 로컬 가상 수신 측에 연결할 수 있게 합니다. SSRF 테스트는 직접 임시로 비웁니다.
	_ = os.Setenv(AllowLocalTargetsEnv, "1")
	os.Exit(m.Run())
}

// TestDialGuardRejectsLoopbackByDefault는 SSRF 보호의 핵심 검사입니다.
// 기본 설정에서는 루프백 주소로의 전송을 반드시 **연결 계층**에서 거부해야 합니다.
func TestDialGuardRejectsLoopbackByDefault(t *testing.T) {
	var hit bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit = true
		_, _ = io.WriteString(w, `{"errcode":0}`)
	}))
	defer srv.Close()

	t.Setenv(AllowLocalTargetsEnv, "") // 명시적 허용 끄기 = 기본 동작
	_, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": srv.URL + "/robot/send"}, Message{Items: []Item{{Severity: "high"}}})
	if err == nil {
		t.Fatal("기본적으로 루프백 주소로의 전송을 허용해서는 안 됩니다")
	}
	if hit {
		t.Fatal("요청이 이미 로컬 서비스에 도달했습니다. 연결 보호 검사가 작동하지 않았습니다")
	}
	// 오류 메시지는 사용자에게 허용 방법을 안내해야 합니다(로컬 SMTP 중계는 정상 설정).
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("거부 메시지에는 명시적으로 허용하는 방법을 설명해야 합니다: %v", err)
	}
}

// TestDialGuardAllowsLoopbackWhenOptedIn은 반대 테스트입니다. 명시적으로 켜면 반드시 사용할 수 있어야 하며,
// 그렇지 않으면 로컬 postfix / 내부 중계와 같은 정상적인 배포를 일괄적으로 막게 됩니다.
func TestDialGuardAllowsLoopbackWhenOptedIn(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"errcode":0,"errmsg":"ok"}`)
	}))
	defer srv.Close()

	t.Setenv(AllowLocalTargetsEnv, "1")
	if _, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": srv.URL + "/robot/send"}, Message{Items: []Item{{Severity: "high"}}}); err != nil {
		t.Fatalf("명시적 허용 후 전송할 수 있어야 합니다: %v", err)
	}
}

func TestIsBlockedDialIP(t *testing.T) {
	blocked := []string{
		"127.0.0.1", "127.1.2.3", "::1",
		"169.254.169.254", // 클라우드 메타데이터 엔드포인트이며, 이 함수가 존재하는 주요 이유
		"169.254.1.1", "fe80::1",
		"0.0.0.0", "::",
		"224.0.0.1", "ff02::1",
		"::ffff:127.0.0.1", // IPv4-mapped 형태는 복원 후 판정해야 하며, 그렇지 않으면 우회 경로가 됨
		"",
	}
	for _, s := range blocked {
		if !isBlockedDialIP(net.ParseIP(s)) {
			t.Errorf("%s: 거부되어야 합니다", s)
		}
	}
	// RFC1918 사설 네트워크는 **의도적으로 허용**합니다. 내부에서 운영하는 Mattermost / SMTP 중계는 흔한 정상 사용 방식입니다.
	// 이 검사는 해당 선택을 고정합니다. 나중에 사설 네트워크 판정을 추가하면 여기서 실패하므로,
	// 의식적인 결정을 요구합니다(배포 일부를 조용히 무효화하지 않음).
	allowed := []string{"10.0.0.5", "172.16.3.4", "192.168.1.10", "8.8.8.8", "2606:4700::1111"}
	for _, s := range allowed {
		if isBlockedDialIP(net.ParseIP(s)) {
			t.Errorf("%s: 허용되어야 합니다(사설 네트워크는 흔한 정상 전송 대상)", s)
		}
	}
}

// TestValidateHTTPURLRejectsLiteralPrivateTargets는 설정 단계의 사전 안내를 검사합니다.
// 리터럴 IP는 첫 전송 실패까지 기다리지 않고 저장 시 거부되어야 합니다.
func TestValidateHTTPURLRejectsLiteralPrivateTargets(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "")
	for _, raw := range []string{
		"http://127.0.0.1:8080/hook",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]:8080/hook",
	} {
		if err := validateHTTPURL(raw); err == nil {
			t.Errorf("%s: 설정 단계에서 거부되어야 합니다", raw)
		}
	}
	// 공인 주소와 사설 주소는 정상적으로 통과합니다(사설 네트워크는 연결 단계에 맡기며, 그곳에서는 차단하지 않음).
	for _, raw := range []string{"https://oapi.dingtalk.com/robot/send", "http://10.0.0.9/hook"} {
		if err := validateHTTPURL(raw); err != nil {
			t.Errorf("%s: 검증에 통과해야 합니다: %v", raw, err)
		}
	}
}

// TestValidateHTTPURLErrorNeverLeaksCredentials는 감사에서 발견한 이전 수정의 누락 분기입니다.
//
// url.Parse가 **실패**하면 *url.Error를 반환하며, Error()에는 원본 주소 전체가 포함됩니다. 이전 수정에서는
// http.Client.Do의 반환 오류에서만 민감정보를 제거하여 이 부분을 놓쳤습니다. 당시 추가한 '영구 실패 경로' 테스트의
// file://·gopher://·ftp://는 모두 url.Parse가 성공하여 scheme 분기를 거쳤으므로,
// 전부 통과해도 이 경로의 안전성을 입증하지 못합니다. 잘못된 보장이었습니다.
func TestValidateHTTPURLErrorNeverLeaksCredentials(t *testing.T) {
	cases := []string{
		"http://127.0.0.1/%zz?access_token=" + leakProbeToken,         // 유효하지 않은 퍼센트 이스케이프
		"https://a.example.com:port/x?access_token=" + leakProbeToken, // 포트가 숫자가 아님
		"http://[::1?access_token=" + leakProbeToken,                  // 괄호 짝이 맞지 않음
	}
	for _, raw := range cases {
		// 먼저 입력이 **실제로** url.Parse를 실패하게 하는지 확인합니다. 하지 않으면 테스트가
		// 아무런 인지 없이 다른 분기를 거칠 수 있습니다(이전의 잘못된 보장이 생긴 원인).
		if _, err := url.Parse(raw); err == nil {
			t.Errorf("%q: 파싱에 실패해야 하며, 그렇지 않으면 이 테스트가 대상 분기를 검사하지 못합니다", raw)
			continue
		}
		err := validateHTTPURL(raw)
		if err == nil {
			t.Errorf("%q: 검증에 실패해야 합니다", raw)
			continue
		}
		assertNoSecret(t, err.Error(), leakProbeToken)
	}
	// 채널 계층에서 감싸더라도 주소가 노출되지 않는지 확인합니다.
	t.Setenv(AllowLocalTargetsEnv, "")
	err := (dingTalkChannel{}).Validate(map[string]any{"webhook": cases[0]})
	if err == nil {
		t.Fatal("유효하지 않은 주소는 검증에 실패해야 합니다")
	}
	assertNoSecret(t, err.Error(), leakProbeToken)
}

// TestEmailDialGuardRejectsLoopbackByDefault는 SMTP 채널의 연결 보호 검사를 검사합니다.
//
// 메일 채널은 이전에 보호 검사가 없는 net.Dialer를 사용하여 전체 SSRF 보호의 유일한 빈틈이었습니다. host에
// 169.254.169.254 또는 127.0.0.1을 입력하면 바로 연결되고, smtp.NewClient 핸드셰이크 실패 시
// 상대 측 반환 한 줄이 오류에 포함되어 last_error를 거쳐 전송 이력 인터페이스에 표시됩니다. 이는 다른 채널에서
// 이미 차단한 오류 응답을 통한 부분 정보 읽기이며, 연결 거부 vs 타임아웃의 소요 시간 차이로 포트도 탐지할 수 있습니다.
//
// 이 패키지의 TestMain은 AllowLocalTargetsEnv를 전역으로 켭니다(많은 테스트가 127.0.0.1의
// 가짜 수신처를 사용). 따라서 이 테스트는 직접 이를 해제해야 합니다. 그렇지 않으면 연결 보호 검사 유무와 관계없이 통과하며,
// 이것이 해당 빈틈을 처음에 어떤 테스트도 발견하지 못한 이유입니다.
func TestEmailDialGuardRejectsLoopbackByDefault(t *testing.T) {
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)

	t.Setenv(AllowLocalTargetsEnv, "") // 명시적 허용 끄기 = 기본 동작
	_, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil {
		t.Fatal("기본적으로 루프백 주소로 메일 전송을 허용해서는 안 됩니다")
	}
	// 연결 자체가 생성되어서는 안 됩니다. Control 훅에서 연결 보호 검사가 차단하므로 EHLO는 전송될 수 없습니다.
	if f.sawCommand("EHLO") || f.sawCommand("HELO") {
		t.Fatal("SMTP 세션이 이미 생성되었습니다. 연결 보호 검사가 작동하지 않았습니다")
	}
	// 오류 메시지는 사용자에게 허용 방법을 안내해야 합니다(로컬 postfix 중계는 정상 설정).
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("거부 메시지에는 명시적으로 허용하는 방법을 설명해야 합니다: %v", err)
	}
}

// TestEmailDialGuardAllowsLoopbackWhenOptedIn은 짝을 이루는 반대 테스트입니다. 명시적으로 켜면
// 정상 전송이 가능해야 합니다. 내부망 자체 SMTP / 로컬 중계는 매우 흔한 배포 형태이므로 연결 보호 검사가 무조건 차단해서는 안 됩니다.
func TestEmailDialGuardAllowsLoopbackWhenOptedIn(t *testing.T) {
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)

	t.Setenv(AllowLocalTargetsEnv, "1")
	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("명시적 허용 후 로컬 SMTP로 전송할 수 있어야 합니다: %v", err)
	}
	if !f.sawCommand("EHLO") {
		t.Fatal("EHLO가 없어 세션이 실제로 설정되지 않았습니다")
	}
}
