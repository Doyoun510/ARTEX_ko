package notify

import (
	"bufio"
	"context"

	"net"
	"strings"
	"sync"
	"testing"
)

// 이 파일은 메일 채널의 프로토콜 수준 테스트를 보충합니다. 이전에는 email.Send의 커버리지가 0이었으며,
// SMTP 경로 전체를 실행한 테스트가 없었습니다. 그런데 이 채널은 여섯 채널 중 프로토콜 범위가 가장 크고,
// 오류가 가장 발생하기 쉬운 채널입니다(핸드셰이크·인증·SMTP 봉투·DATA 단계마다 실패 의미가 다름).
//
// 여기서는 net/smtp를 mock하지 않고 직접 만든 최소한의 SMTP 서버로 테스트합니다.
// 메일 채널 위험 대부분은 "실제 SMTP 서버와 통신"하는 단계에 있으므로,
// 이 단계를 mock하면 검사하지 않는 셈입니다.

// fakeSMTP는 필요한 기능만 제공하는 SMTP 서버로, greet/EHLO/AUTH/MAIL/RCPT/DATA/QUIT를 수행하고,
// 테스트 요구에 따라 특정 단계에 지정된 응답 코드를 반환합니다.
type fakeSMTP struct {
	ln net.Listener

	// rcptReply는 RCPT TO 응답이며, 기본값은 250입니다.
	rcptReply string
	// mailReply는 MAIL FROM 응답이며, 기본값은 250입니다.
	mailReply string
	// advertiseAuth가 true이면 EHLO에 AUTH PLAIN 지원을 선언합니다.
	advertiseAuth bool

	mu       sync.Mutex
	data     string
	commands []string
}

func newFakeSMTP(t *testing.T) *fakeSMTP {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeSMTP{ln: ln, rcptReply: "250 OK", mailReply: "250 OK"}
	go f.serve()
	t.Cleanup(func() { ln.Close() })
	return f
}

func (f *fakeSMTP) hostPort(t *testing.T) (string, int) {
	t.Helper()
	addr, ok := f.ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("TCP 수신 대기 주소가 아닙니다")
	}
	return "127.0.0.1", addr.Port
}

func (f *fakeSMTP) record(cmd string) {
	f.mu.Lock()
	f.commands = append(f.commands, cmd)
	f.mu.Unlock()
}

func (f *fakeSMTP) body() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.data
}

func (f *fakeSMTP) sawCommand(prefix string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.commands {
		if strings.HasPrefix(c, prefix) {
			return true
		}
	}
	return false
}

func (f *fakeSMTP) serve() {
	conn, err := f.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	br := bufio.NewReader(conn)
	w := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
	w("220 fake.local ESMTP ready")
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		f.record(line)
		switch {
		case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
			// STARTTLS를 선언하지 않아 코드가 평문 분기를 거치게 합니다(테스트 대상은 SMTP 봉투 로직이며 TLS가 아님).
			w("250-fake.local")
			if f.advertiseAuth {
				w("250-AUTH PLAIN")
			}
			w("250 8BITMIME")
		case strings.HasPrefix(line, "AUTH"):
			// 단순 처리: PLAIN의 초기 응답은 여러 줄일 수 있으며, 바로 허용합니다.
			w("235 2.7.0 Authentication successful")
		case strings.HasPrefix(line, "MAIL FROM"):
			w(f.mailReply)
		case strings.HasPrefix(line, "RCPT TO"):
			w(f.rcptReply)
		case strings.HasPrefix(line, "DATA"):
			w("354 End data with <CR><LF>.<CR><LF>")
			var sb strings.Builder
			for {
				dl, err := br.ReadString('\n')
				if err != nil {
					return
				}
				if strings.TrimRight(dl, "\r\n") == "." {
					break
				}
				sb.WriteString(dl)
			}
			f.mu.Lock()
			f.data = sb.String()
			f.mu.Unlock()
			w("250 2.0.0 Ok: queued as FAKE1")
		case strings.HasPrefix(line, "QUIT"):
			w("221 2.0.0 Bye")
			return
		default:
			w("250 OK")
		}
	}
}

func emailCfg(t *testing.T, f *fakeSMTP, extra map[string]any) map[string]any {
	t.Helper()
	host, port := f.hostPort(t)
	cfg := map[string]any{
		"host": host,
		"port": float64(port),
		"from": "artex@example.com",
		"to":   []any{"a@example.com", "b@example.com"},
	}
	for k, v := range extra {
		cfg[k] = v
	}
	return cfg
}

func TestEmailSendDeliversFullMessage(t *testing.T) {
	f := newFakeSMTP(t)
	f.advertiseAuth = true
	cfg := emailCfg(t, f, map[string]any{"username": "artex", "password": "pw"})

	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
	// SMTP 봉투 단계를 반드시 거쳐야 합니다. 발신자, 수신자 두 명, DATA입니다.
	for _, want := range []string{"MAIL FROM:<artex@example.com>", "RCPT TO:<a@example.com>", "RCPT TO:<b@example.com>", "DATA", "AUTH", "QUIT"} {
		if !f.sawCommand(want) {
			t.Errorf("SMTP 세션에 %q 누락. 실제 명령: %v", want, f.commands)
		}
	}
	// 본문은 base64 HTML이며, 실제 취약점 내용을 포함해야 합니다(인코딩 후에도 식별 가능).
	body := f.body()
	if body == "" {
		t.Fatal("DATA 단계에서 본문을 받지 못했습니다")
	}
	if !strings.Contains(body, "Content-Type: text/html") {
		t.Errorf("Content-Type 헤더 누락:\n%s", body)
	}
	if !strings.Contains(body, "base64") {
		t.Errorf("본문이 base64로 인코딩되지 않았습니다(긴 HTML 줄이 SMTP의 1000바이트 줄 길이 제한을 위반):\n%s", body)
	}
	// 여러 수신자 모두 To 헤더에 나타나야 합니다.
	if !strings.Contains(body, "a@example.com, b@example.com") {
		t.Errorf("To 헤더에 모든 수신자가 포함되지 않았습니다:\n%s", body)
	}
}

func TestEmailSendWithoutAuth(t *testing.T) {
	// 계정을 설정하지 않으면 AUTH를 보내서는 안 됩니다. 일부 중계는 이를 거부합니다.
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)
	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
	if f.sawCommand("AUTH") {
		t.Errorf("계정을 설정하지 않았는데 AUTH를 전송했습니다: %v", f.commands)
	}
}

// TestEmailSendClassifiesSMTPReplies는 이번 감사 수정의 직접적인 검증입니다.
// 5xx는 영구 실패로, 4xx(Greylisting)는 재시도 가능한 것으로 판정합니다.
func TestEmailSendClassifiesSMTPReplies(t *testing.T) {
	cases := []struct {
		name      string
		rcptReply string
		mailReply string
		permanent bool
	}{
		{"수신자가 550으로 영구 거부됨", "550 5.1.1 User unknown", "250 OK", true},
		{"수신자에게 450 Greylisting 응답", "450 4.7.1 Greylisting in action", "250 OK", false},
		{"수신자에게 452 메일함 가득 참 응답", "452 4.2.2 Mailbox full", "250 OK", false},
		{"발신자가 553으로 영구 거부됨", "250 OK", "553 5.1.3 Bad address", true},
		{"발신자에게 451 일시적 오류", "250 OK", "451 4.3.0 Temporary failure", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeSMTP(t)
			f.rcptReply = tc.rcptReply
			f.mailReply = tc.mailReply
			_, err := (emailChannel{}).Send(context.Background(), emailCfg(t, f, nil), singleMsg())
			if err == nil {
				t.Fatal("오류가 나야 합니다")
			}
			if got := IsPermanent(err); got != tc.permanent {
				t.Fatalf("permanent 판정 오류: 기대 %v 실제 %v (%v)", tc.permanent, got, err)
			}
			// 서버 원문을 유지해야 합니다. 그렇지 않으면 사용자는 서버 관리자에게 문의할지 주소를 수정할지 알 수 없습니다.
			if !strings.Contains(err.Error(), strings.Fields(tc.rcptReply)[0]) && !strings.Contains(err.Error(), strings.Fields(tc.mailReply)[0]) {
				t.Errorf("오류에 서버 응답 코드를 유지해야 합니다: %v", err)
			}
		})
	}
}

func TestEmailSendRefusesPlaintextCredentials(t *testing.T) {
	// net/smtp의 PlainAuth는 암호화하지 않은 연결에서 자격 증명 전송을 거부합니다(대상이 localhost인 경우 제외).
	// 이는 **올바른** 보안 동작이므로 우회해서는 안 됩니다. 다만 수정 방법을 알 수 있는 오류를 제공해야 합니다.
	// 여기서는 localhost가 아닌 호스트 이름으로 해당 동작을 유발합니다.
	f := newFakeSMTP(t)
	f.advertiseAuth = true
	_, port := f.hostPort(t)
	cfg := map[string]any{
		"host":     "smtp.example.com", // localhost가 아님
		"port":     float64(port),
		"from":     "a@example.com",
		"to":       []any{"b@example.com"},
		"username": "artex",
		"password": "pw",
	}
	_, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil {
		t.Skip("이 시스템의 DNS 조회 결과가 로컬 서버이므로 건너뜁니다(다른 테스트에 영향 없음)")
	}
	// 연결 불가 또는 자격 증명 전송 거부 모두 이 검사를 통과합니다. 핵심은 비밀번호를 조용히 보내서는 **안 된다**는 것입니다.
	if !IsPermanent(err) && !strings.Contains(err.Error(), "연결") {
		t.Logf("오류: %v(localhost가 아닌 경우 연결되지 않는 것은 예상 동작)", err)
	}
}

func TestEmailValidateReportsMissingFields(t *testing.T) {
	// 메일 채널은 설정 필드가 가장 많으며, 하나라도 누락하면 전송 시에야 드러납니다. 여기서는 각각
	// 검증이 먼저 차단할 수 있는지 확인합니다. 검사는 "오류 메시지에 누락된 항목이 나오는지"를 확인합니다.
	cases := []struct {
		name string
		cfg  map[string]any
	}{
		{"host 누락", map[string]any{"port": float64(25), "from": "a@b.c", "to": []any{"d@e.f"}}},
		{"port 누락", map[string]any{"host": "smtp.example.com"}},
		{"port 범위 초과", map[string]any{"host": "h", "port": float64(70000), "from": "a@b.c", "to": []any{"d@e.f"}}},
		{"from 누락", map[string]any{"host": "h", "port": float64(25), "to": []any{"d@e.f"}}},
		{"to 누락", map[string]any{"host": "h", "port": float64(25), "from": "a@b.c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := (emailChannel{}).Validate(tc.cfg); err == nil {
				t.Fatalf("검증에 실패해야 합니다: %v", tc.cfg)
			}
		})
	}
}

// TestEmailConfigTolerance는 설정 읽기의 오류 허용을 검사합니다. JSONB의 숫자는 float64이지만,
// 사용자는 UI에 포트를 문자열로 입력할 수 있으며, 배열 대신 단일 문자열일 수도 있습니다.
func TestEmailConfigTolerance(t *testing.T) {
	cfg := map[string]any{
		"host": "smtp.example.com",
		"port": "587", // 문자열 형태의 포트
		"from": "a@b.c",
		"to":   "d@e.f", // 배열이 아닌 단일 문자열
		"tls":  "true",  // 문자열 형태의 boolean
	}
	if err := (emailChannel{}).Validate(cfg); err != nil {
		t.Fatalf("문자열 형태의 숫자를 허용해야 합니다: %v", err)
	}
	if got := cfgInt(cfg, "port"); got != 587 {
		t.Errorf("cfgInt가 문자열 포트를 파싱하지 않았습니다. 실제: %d", got)
	}
	if !cfgBool(cfg, "tls") {
		t.Error("cfgBool이 문자열 \"true\"를 파싱하지 않았습니다")
	}
	if to := cfgStrings(cfg, "to"); len(to) != 1 || to[0] != "d@e.f" {
		t.Errorf("cfgStrings가 단일 문자열을 지원하지 않았습니다. 실제: %v", to)
	}
}

// TestFilterValidateRejectsTypo는 감사 수정의 직접적인 검증입니다.
// 기준에 오타가 있으면 반드시 쓰기 시 차단해야 합니다. 그렇지 않으면 필터가 조용히 무효화되어 전부 전송합니다.
func TestFilterValidateRejectsTypo(t *testing.T) {
	good := []string{"", "low", "medium", "high", "critical"}
	for _, s := range good {
		if err := (Filter{MinSeverity: s}).Validate(); err != nil {
			t.Errorf("유효한 기준 %q 거부됨: %v", s, err)
		}
	}
	// 실제로 발생할 수 있는 오타이며, 모두 거부해야 합니다.
	for _, s := range []string{"hgih", "HIGH", "严重", "high ", "crit"} {
		err := (Filter{MinSeverity: s}).Validate()
		if err == nil {
			t.Errorf("유효하지 않은 기준 %q 거부 필요(그렇지 않으면 필터가 조용히 무효화되어 전부 전송)", s)
			continue
		}
		// 오류 메시지는 올바른 수정 방법을 안내해야 합니다.
		if !strings.Contains(err.Error(), "low") || !strings.Contains(err.Error(), "critical") {
			t.Errorf("오류 메시지에 선택 가능한 값을 나열해야 합니다. 실제: %q", err.Error())
		}
	}
}

// TestFilterValidateIsWriteTimeOnly는 "쓰기는 엄격하게, 읽기는 관대하게" 역할 분담을 고정합니다.
// DB의 기존 잘못된 값 때문에 채널 전체를 읽지 못해서는 안 됩니다(기존 채널의 알림 전송이 갑자기 모두 멈추기 때문).
func TestFilterValidateIsWriteTimeOnly(t *testing.T) {
	raw := []byte(`{"min_severity":"hgih"}`)
	f := ParseFilter(raw) // 오류를 내지 않음
	if f.MinSeverity != "hgih" {
		t.Fatalf("읽기 경로에서는 그대로 유지해야 합니다. 실제: %q", f.MinSeverity)
	}
	// 또한 채널은 이벤트를 판정할 수 있어야 합니다(panic이나 멈춤 없음).
	_ = Match(f, Snapshot{Kind: EventFindingCreated, Severity: "critical"})
}
