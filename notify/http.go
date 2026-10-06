package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"
)

// allowLocalTargets는 루프백 / link-local 주소로 메시지 전송을 허용할지 결정합니다.
//
// 기본적으로 거부합니다. 이 주소 범위에는 IM 봇이나 공개 메일 서버가 없지만, 접근할 수 있는 대상은
// 민감합니다. 동일 기기의 다른 서비스 관리 포트와 클라우드 환경의 메타데이터 엔드포인트
// (169.254.169.254, 인스턴스 자격 증명 조회 가능)입니다. 전송 주소는 관리자가 설정하지만, XSS/CSRF로
// 악용된 관리 세션이나 동일한 JWT를 공유하는 다른 사람도 설정을 바꾸어 응답 내용을 읽을 수 있습니다.
// doJSON이 4xx/5xx 응답 본문의 처음 200바이트를 last_error에 쓰고, 전송 이력 인터페이스가
// 이를 다시 표시하므로 오류 응답을 통한 부분 정보 읽기가 가능합니다.
//
// 하지만 "로컬 SMTP 중계"(127.0.0.1:25의 postfix)는 자체 메일 서버의 흔한 설정이므로,
// 일괄 차단하면 사용이 막힙니다. 따라서 하드코딩으로 허용하는 대신 명시적인 허용 방법을 제공합니다.
// ARTEX_NOTIFY_ALLOW_LOCAL=1로 설정하면 허용합니다.
//
// AllowLocalTargetsEnv로 공개한 이유는 테스트에서 명시적으로 켤 수 있게 하기 위해서입니다. 이 패키지와 server 패키지는
// 127.0.0.1의 httptest 가짜 수신처를 많이 사용하며, 켜지 않으면 모두 연결 보호 검사에서 차단됩니다.
const AllowLocalTargetsEnv = "ARTEX_NOTIFY_ALLOW_LOCAL"

func allowLocalTargets() bool {
	v := strings.TrimSpace(os.Getenv(AllowLocalTargetsEnv))
	return v == "1" || strings.EqualFold(v, "true")
}

// isBlockedDialIP는 대상 IP가 "기본적으로 전송을 허용하지 않는" 주소 대역에 속하는지 반환합니다.
//
// 루프백·link-local 주소(클라우드 메타데이터 169.254.169.254 포함)·미지정 주소·멀티캐스트만 거부합니다.
// RFC1918 사설 네트워크는 **거부하지 않습니다**. 내부에서 운영하는 Mattermost / SMTP 중계는 흔한 정상 사용 방식이며,
// 함께 차단하면 실제 환경에서 기능을 사용할 수 없게 됩니다. 이 선택은 의도적입니다.
// 실제로 민감한 대상은 보호하면서 정상적인 배포까지 막아서는 안 됩니다.
func isBlockedDialIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	// IPv4-mapped IPv6(::ffff:127.0.0.1)는 IPv4로 복원한 뒤 판정해야 합니다. 그렇지 않으면 검사를 우회합니다.
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
}

// blockInternalDial은 http.Transport 연결 함수의 Control 훅으로, **연결을 설정할 때**
// 대상 주소를 검사합니다.
//
// 설정 저장 시 검증에만 의존하지 않고 연결 단계에 두는 이유: 이곳이 최종적으로 적용되는 위치입니다.
// 설정 검증을 우회하는 두 경우를 함께 처리합니다. DNS 리바인딩(검증 시에는 공개 IP로 조회되지만,
// 실제 연결 시에는 내부망으로 조회됨)과 리디렉션(다른 호스트로의 이동은 이미 거부하지만, 같은 호스트로의 이동도
// 경로를 다른 곳으로 향하게 할 수 있음)입니다.
func blockInternalDial(_, address string, _ syscall.RawConn) error {
	if allowLocalTargets() {
		return nil
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("대상 주소를 파싱할 수 없습니다: %q", host)
	}
	if isBlockedDialIP(ip) {
		return fmt.Errorf("로컬/link-local 주소로의 전송 거부: %s(로컬 서비스에 반드시 전송해야 한다면 %s=1 설정)", ip, AllowLocalTargetsEnv)
	}
	return nil
}

// notifyTransport는 기본 Transport에 연결 보호 검사만 추가합니다.
// Clone으로 기본 조정 사항 전체(연결 풀, HTTP/2, 타임아웃, proxy 등)를 유지하여,
// 검사 하나를 추가하기 위해 다른 동작이 바뀌는 것을 방지합니다.
var notifyTransport = func() *http.Transport {
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Transport{}
	}
	clone := t.Clone()
	clone.DialContext = (&net.Dialer{Timeout: 10 * time.Second, Control: blockInternalDial}).DialContext
	return clone
}()

// httpClient는 모든 채널 전송에 공통으로 사용하는 클라이언트입니다.
//
// 프로젝트의 전역 아웃바운드 프록시(server 측 GlobalProxy)를 의도적으로 재사용하지 **않습니다**. 해당 프록시는 침투 테스트
// 대상 트래픽용이며 불안정한 터널인 경우가 많고, 알림 가용성이 대상 네트워크의 불안정에 좌우되어서는 안 됩니다.
// IM 알림 전송은 직접 연결하면 됩니다. 타임아웃은 15초이며, 이보다 느린 상대 측은 사실상 이미 장애 상태입니다.
//
// 다른 호스트로의 리디렉션을 거부합니다. 이 기능의 전송 주소는 모두 "고정 endpoint 하나" 형태이며, 일반적으로
// 다른 호스트로 리디렉션하지 않습니다. 이 플랫폼들의 자격 증명(DingTalk의 access_token, WeCom(기업용 위챗)의 key, Telegram의
// bot token)은 **URL 안에** 있으므로, 다른 호스트로 이동하면 리디렉션 대상에게 자격 증명을 전달하게 됩니다. 같은 호스트로의
// 이동(끝에 슬래시를 붙이는 경우 등)은 여전히 허용합니다.
var httpClient = &http.Client{
	Timeout:   15 * time.Second,
	Transport: notifyTransport,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("리디렉션 횟수가 너무 많습니다")
		}
		if len(via) > 0 && req.URL.Host != via[0].URL.Host {
			return fmt.Errorf("다른 호스트로의 리디렉션 거부(%s → %s)", via[0].URL.Host, req.URL.Host)
		}
		return nil
	},
}

// respBodyLimit는 읽을 응답 본문의 크기를 제한합니다. 상대 측 오류로 매우 큰 내용이 반환될 수 있지만, 필요한 것은
// 전송 이력에 표시할 오류 코드와 짧은 오류 설명뿐입니다.
const respBodyLimit = 8 << 10

// doJSON은 요청을 한 번 보내고 응답 본문을 반환합니다(길이 제한 적용).
//
// payload가 nil이면 빈 body를 보냅니다(GET 또는 플랫폼이 body를 요구하지 않는 경우).
// headers의 키와 값은 그대로 붙이며, 일반 Webhook의 사용자 지정 헤더에 사용합니다.
//
// 오류 분류는 이 함수의 핵심 역할입니다. 네트워크 계층 실패와 5xx/408/429는 "재시도 가능"으로,
// 나머지 4xx는 "영구 실패"로 분류합니다. 403을 재시도해도 같은 오류를 로그에 3번 반복할 뿐입니다.
func doJSON(ctx context.Context, method, url string, headers map[string]string, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			// 직렬화 실패는 로컬 bug(설정 필드 타입이 잘못됨)이므로, 재시도해도 나아지지 않습니다.
			return nil, Permanent(fmt.Errorf("요청 본문 구성 실패: %w", err))
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		// URL이 유효하지 않으면 대부분 사용자가 주소를 잘못 입력한 경우이며, 영구 실패입니다.
		// 여기서도 err를 그대로 전달해서는 안 됩니다. url.Parse 오류 텍스트에는 전체 주소가 포함됩니다.
		return nil, Permanent(fmt.Errorf("요청 주소가 유효하지 않습니다: %s", redactRequestTarget(url)))
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		// 연결 거부·DNS 실패·타임아웃은 대부분 일시적 오류이므로 백오프 후 재시도에 맡깁니다.
		//
		// 오류 텍스트는 외부로 전달하기 전에 반드시 민감정보를 제거해야 합니다. http.Client.Do는 *url.Error를 반환하며,
		// Error()는 `Op "전체URL": 하위 오류`이고 이 기능에서 사용하는 자격 증명은 **URL 안에** 있습니다
		// (DingTalk access_token, WeCom key, Feishu hook id, Telegram /bot<token>/).
		// 민감정보를 제거하지 않으면 자격 증명이 이 오류 문자열을 통해 네 곳으로 전달됩니다. notification_deliveries의
		// last_error(평문으로 DB에 저장), 전송 이력 인터페이스 응답(**채널 설정 마스킹 우회**),
		// 서버 로그, 테스트 전송 인터페이스가 프런트엔드에 반환하는 502 텍스트입니다.
		return nil, fmt.Errorf("요청 실패: %s", redactTransportError(err))
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, respBodyLimit))
	if readErr != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", readErr)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return raw, nil
	}
	// 429(전송 속도 제한)와 408(타임아웃)은 재시도할 가치가 있습니다. 나머지 4xx는 설정 또는 권한 문제이므로 재시도는 의미가 없습니다.
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout {
		return nil, fmt.Errorf("상대 측 전송 속도 제한 또는 타임아웃 (HTTP %d): %s", resp.StatusCode, snippet(raw))
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("상대 측 서비스 오류 (HTTP %d): %s", resp.StatusCode, snippet(raw))
	}
	return nil, Permanent(fmt.Errorf("상대 측 요청 거부 (HTTP %d): %s", resp.StatusCode, snippet(raw)))
}

// snippet은 응답 본문을 짧은 한 줄 텍스트로 압축하여 오류 메시지에 사용합니다. 응답에 줄바꿈과 많은 공백이 포함될 수 있어,
// last_error에 바로 넣으면 전송 이력 페이지의 레이아웃이 깨집니다.
func snippet(raw []byte) string {
	return OneLine(string(raw), 200)
}

// redactRequestTarget은 전송 주소를 "scheme://host/…"로 압축하여 오류 메시지에 사용합니다.
//
// 이것이 이 패키지의 유일한 주소 민감정보 제거 기준이며, 의도적으로 **충분히 과감하게** 처리합니다. scheme과 host를 제외한
// 모든 부분을 버립니다. URL의 어느 부분이 자격 증명인지 판단할 '일반적이고 안전한' 방법이 없기 때문입니다.
//
//	DingTalk   자격 증명은 query에 있음      /robot/send?access_token=xxx
//	WeCom   자격 증명은 query에 있음      /cgi-bin/webhook/send?key=xxx
//	Feishu   자격 증명은 **경로 끝부분**에 있음 /open-apis/bot/v2/hook/<hook_id>
//	Telegram 자격 증명은 **경로 중간**에 있음 /bot<token>/sendMessage
//
// "유용한 부분만 유지"하려면 채널별 처리가 필요하지만, 하나라도 누락하면 자격 증명이 유출됩니다.
// host만 유지해도 원인 파악에 충분합니다(DNS 조회 실패·연결 불가·잘못된 인증서 등을 확인할 수 있음).
// 어떤 봇인지는 채널 설정의 마스킹된 마지막 자리 힌트로 구분합니다.
//
// 파싱 실패 시 고정 자리표시자를 반환하며, 원래 문자열은 절대로 반환하지 않습니다.
func redactRequestTarget(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "(주소를 파싱할 수 없음)"
	}
	return u.Scheme + "://" + u.Host + "/…"
}

// redactTransportError는 전송 계층 오류에서 주소를 제거하고 원래 원인만 유지합니다.
//
// *url.Error 구조는 {Op, URL, Err}이며, Error()는 URL도 함께 출력합니다.
// 여기서는 Err 필드를 명시적으로 가져와 Error()를 거치지 않습니다. 사후 문자열 치환보다 더 확실하며,
// 치환으로는 URL 인코딩/이스케이프의 다양한 형태를 올바르게 처리해야 하므로 누락하기 쉽습니다.
func redactTransportError(err error) string {
	var uerr *url.Error
	if errors.As(err, &uerr) {
		host := ""
		if u, parseErr := url.Parse(uerr.URL); parseErr == nil {
			host = u.Host
		}
		if uerr.Err != nil {
			return fmt.Sprintf("%s %s: %s", uerr.Op, host, uerr.Err)
		}
		return fmt.Sprintf("%s %s: 알 수 없는 오류", uerr.Op, host)
	}
	// *url.Error가 아닌 오류(리디렉션 정책이 반환하는 오류 등)에도 주소가 포함될 수 있으므로 모두 민감정보를 제거합니다.
	return redactURLsInText(err.Error())
}

// redactURLsInText는 텍스트에 나타나는 http(s) 주소를 민감정보가 제거된 형태로 바꿉니다.
//
// 구조화된 필드를 얻을 수 없는 오류(리디렉션 정책 오류, 외부 라이브러리의 사용자 지정 오류)에 대체 처리를 적용합니다.
// http/https 접두사만 인식하며 공백과 따옴표로 나눕니다. 주소에는 이 두 종류의 문자가 포함되지 않습니다.
func redactURLsInText(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		rest := s[i:]
		if strings.HasPrefix(rest, "http://") || strings.HasPrefix(rest, "https://") {
			end := len(rest)
			if j := strings.IndexAny(rest, " \t\n\"'"); j >= 0 {
				end = j
			}
			b.WriteString(redactRequestTarget(rest[:end]))
			i += end
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
