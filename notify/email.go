package notify

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// emailDialTimeout / emailSessionTimeout은 각각 연결과 전체 SMTP 세션을 제한합니다.
// net/smtp 자체에는 타임아웃 기능이 없으므로 두 제한을 설정하지 않으면, 상대 측이 멈췄을 때
// 전송 goroutine이 계속 멈춘 상태로 남습니다. dispatcher는 단일 goroutine에서 순차 처리하므로,
// 알림 시스템 전체가 멈추는 셈입니다.
const (
	emailDialTimeout    = 10 * time.Second
	emailSessionTimeout = 45 * time.Second
)

// emailChannel은 SMTP 메일 전송을 구현합니다.
type emailChannel struct{}

func (emailChannel) Kind() string { return KindEmail }

// 메일에는 플랫폼 전송 속도 제한이 없지만 메시지를 과도하게 보내서는 안 되므로, 넉넉한 기본값을 제공합니다.
func (emailChannel) DefaultRatePerMin() int { return 60 }

// 비밀번호만 마스킹합니다. SMTP 호스트·계정·수신자는 비밀 정보가 아니며, 마스킹하면 편집만 어려워집니다.
func (emailChannel) SecretKeys() []string { return []string{"password"} }

// host/port는 어느 서버에 비밀번호를 전달할지 결정하고, tls는 암호화 전송 여부를 결정합니다. 셋 중 하나라도 바뀌면
// 비밀번호를 다시 명시해야 합니다. "TLS 끄기"도 단순 변경이 아니라 자격 증명을 명시적으로 포함해야 하도록 합니다.
func (emailChannel) DestinationKeys() []string { return []string{"host", "port", "tls"} }

func (emailChannel) Validate(cfg map[string]any) error {
	if cfgString(cfg, "host") == "" {
		return errors.New("SMTP 서버 주소가 없습니다")
	}
	port := cfgInt(cfg, "port")
	if port <= 0 || port > 65535 {
		return errors.New("SMTP 포트가 유효하지 않습니다(1-65535여야 합니다)")
	}
	if cfgString(cfg, "from") == "" {
		return errors.New("발신자 주소가 없습니다")
	}
	if len(cfgStrings(cfg, "to")) == 0 {
		return errors.New("수신자 주소가 하나 이상 필요합니다")
	}
	return nil
}

func (c emailChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	host := cfgString(cfg, "host")
	port := cfgInt(cfg, "port")
	from := cfgString(cfg, "from")
	to := cfgStrings(cfg, "to")
	username := cfgString(cfg, "username")
	password := cfgString(cfg, "password")
	implicitTLS := cfgBool(cfg, "tls")

	msg, err := buildEmailMessage(from, to, m)
	if err != nil {
		return 0, Permanent(err)
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := emailDial(ctx, addr, host, implicitTLS)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// STARTTLS: 상대 측이 지원하면 전환합니다. 평문 세션에서는 자격 증명을 보낼 수 없습니다(아래 auth 설명 참조).
	if !implicitTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
				return 0, fmt.Errorf("STARTTLS 실패: %w", err)
			}
		}
	}
	if username != "" {
		if err := client.Auth(smtp.PlainAuth("", username, password, host)); err != nil {
			// smtp.PlainAuth는 암호화되지 않은 연결에서 자격 증명 전송을 거부합니다(대상이 localhost인 경우 제외).
			// 이는 **올바른** 보안 동작이므로 우회할 수 없지만, 원인을 명확히 번역해야 합니다.
			// 그렇지 않으면 사용자는 "unencrypted connection"만 보고 어떻게 해야 할지 알 수 없습니다.
			if strings.Contains(err.Error(), "unencrypted connection") {
				return 0, Permanent(fmt.Errorf("자격 증명 전송 거부: 연결이 암호화되지 않았습니다. TLS를 활성화하거나 465 포트(암시적 TLS)로 변경하거나 'TLS 활성화'를 선택하세요 (%w)", err))
			}
			return 0, Permanent(fmt.Errorf("SMTP 인증 실패: %w", err))
		}
	}
	if err := client.Mail(from); err != nil {
		return 0, smtpStageError(fmt.Sprintf("발신자 %s: 거부됨", from), err)
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return 0, smtpStageError(fmt.Sprintf("수신자 %s: 거부됨", rcpt), err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return 0, fmt.Errorf("SMTP DATA 실패: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return 0, fmt.Errorf("메일 본문 쓰기 실패: %w", err)
	}
	if err := w.Close(); err != nil {
		return 0, fmt.Errorf("메일 제출 실패: %w", err)
	}
	// Quit 실패는 "서버가 메일을 받았다"는 사실에 영향을 주지 않으므로 오류를 무시합니다.
	_ = client.Quit()
	// 메일은 길이에 따라 자르지 않으므로(HTML 본문 전체 전송), 배치 전체를 전달된 것으로 처리합니다.
	return len(m.Items), nil
}

// emailDial은 SMTP 연결을 설정합니다.
//
// implicitTLS=true이면 465와 같은 "연결 즉시 TLS" 방식이며, false이면 25/587에서 평문으로 연결한 뒤
// STARTTLS를 사용합니다. 두 방식을 섞으면 안 됩니다. 465 포트에 평문 greeting을 보내면 연결이 즉시 끊깁니다.
//
// 세션 제한 시간은 **연결 생성 위치**에서 설정합니다(나중에 추가하지 않음). net/smtp의 Client가 하위
// 연결을 외부로 공개하지 않은 필드에 숨겨 외부에서 가져올 수 없으며, 연결을 넘긴 뒤에는 미리 설정한 deadline을
// 최종 보호 조치로 사용할 수밖에 없기 때문입니다. 이는 핸드셰이크 단계에서 멈추는 경우도 처리합니다.
// Control에 blockInternalDial을 연결해 HTTP 계열 채널과 동일한 연결 보호 검사를 사용합니다. 연결하지 않으면 SMTP는
// 전체 SSRF 보호의 빈틈이 됩니다. host에 169.254.169.254 또는 127.0.0.1을 입력하면 바로 연결되고,
// smtp.NewClient 핸드셰이크 실패 시 상대 측에서 반환한 한 줄이 오류에 포함되어 last_error를 거쳐
// 전송 이력 인터페이스에 표시되므로 오류 응답을 통한 부분 정보 읽기가 가능합니다. 연결 거부 vs 타임아웃의 소요 시간 차이로
// 포트도 탐지할 수 있습니다. 연결 단계가 최종 적용 지점이며 DNS 리바인딩도 처리합니다.
func emailDial(ctx context.Context, addr, host string, implicitTLS bool) (*smtp.Client, error) {
	d := &net.Dialer{Timeout: emailDialTimeout, Control: blockInternalDial}
	var conn net.Conn
	var err error
	if implicitTLS {
		conn, err = tls.DialWithDialer(d, "tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("SMTP 서버 연결 실패: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(emailSessionTimeout))
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("SMTP 핸드셰이크 실패: %w", err)
	}
	return client, nil
}

// smtpStageError는 SMTP 응답 코드에 따라 특정 단계의 실패를 "재시도 가능"과 "영구 실패"로 구분합니다.
//
// 구분이 필수인 이유: SMTP의 4xx와 5xx는 의미가 완전히 다릅니다.
//   - 4xx(450 Greylisting, 451 로컬 오류, 452 저장 공간 부족)는 **일시적** 거부이며,
//     일반적인 처리는 나중에 재시도하는 것입니다. 특히 Greylisting은 거의 매번 첫 전송 시 발생합니다.
//   - 5xx(550 사용자 없음, 553 주소 유효하지 않음)는 영구 거부이며, 재시도는 의미가 없습니다.
//
// 모든 경우를 영구 실패로 판정하면 Greylisting을 사용하는 메일 서버에서는 **모든** 알림이 첫 번째
// 시도 후 failed로 넘어갑니다. 이러한 실패는 자동 재시도가 가장 효과를 발휘해야 하는 상황입니다.
// 응답 코드는 오류 텍스트의 처음 세 자리 숫자에서 가져옵니다. 코드를 가져올 수 없으면 재시도 가능으로 처리하며(한 번 더 시도하더라도,
// 파싱할 수 없다는 이유로 일시적 오류일 수 있는 실패를 영구 실패로 판정하지 않음).
func smtpStageError(what string, err error) error {
	code := smtpReplyCode(err.Error())
	if code >= 500 && code < 600 {
		return Permanent(fmt.Errorf("%s: %w", what, err))
	}
	return fmt.Errorf("%s: %w", what, err)
}

// smtpReplyCode는 SMTP 오류 텍스트의 첫 세 자리 응답 코드를 가져오며, 없으면 0을 반환합니다.
// net/smtp는 오류 코드 필드를 노출하지 않으므로 텍스트에서 가져와야 합니다. 형식은 "450 4.7.1 ..."입니다.
func smtpReplyCode(text string) int {
	if len(text) < 3 {
		return 0
	}
	n, err := strconv.Atoi(text[:3])
	if err != nil {
		return 0
	}
	return n
}

// buildEmailMessage는 전체 RFC 5322 메일을 구성합니다.
//
// 본문을 base64로 인코딩하는 데는 두 가지 이유가 있습니다. 첫째, SMTP는 한 줄을 1000바이트 이하로 제한하지만 HTML
// 본문(특히 모아 보내기 메일)에는 긴 줄이 쉽게 생깁니다. 둘째, base64에는 "."으로 시작하는
// 줄이 없어 SMTP 점 이스케이프를 처리할 필요가 없습니다.
func buildEmailMessage(from string, to []string, m Message) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(to, ", "))
	// 중국어 제목은 RFC 2047로 인코딩해야 합니다. 그렇지 않으면 클라이언트에서 깨진 문자로 표시됩니다.
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", htmlTitle(m)))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("Content-Transfer-Encoding: base64\r\n")
	// 메일에는 엄격한 길이 상한이 없으므로 본문을 자르지 않습니다.
	b.WriteString("\r\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(htmlBody(m, 0)))
	// base64는 76문자 단위로 줄을 나누어 RFC 2045를 준수합니다.
	for len(encoded) > 76 {
		b.WriteString(encoded[:76] + "\r\n")
		encoded = encoded[76:]
	}
	b.WriteString(encoded + "\r\n")
	return b.String(), nil
}
