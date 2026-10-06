package notify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

// dingTalkChannel은 DingTalk 사용자 지정 봇을 구현합니다.
//
// 플랫폼 특성(여기서의 구현 선택을 결정):
//   - 봇 하나의 전송 속도 제한은 분당 20건이며, 초과분은 조용히 버려집니다(HTTP는 여전히 200일 수 있음).
//     따라서 클라이언트에서 반드시 전송 속도를 제한해야 합니다. DefaultRatePerMin을 참조하세요.
//   - 보안 설정은 서명 추가 / 사용자 지정 키워드 / IP 허용 목록 중 하나를 선택합니다. 서명 추가는 메시지 내용에 의존하지 않는
//     유일한 방식이므로 서명 추가만 지원합니다(세 가지를 모두 켜지 않은 webhook도 지원).
//   - 성공/실패 모두 HTTP 200을 반환하며, body의 errcode로 구분합니다. errcode를 확인하지 않으면
//     전송 실패를 성공으로 기록하게 됩니다.
type dingTalkChannel struct{}

func (dingTalkChannel) Kind() string { return KindDingTalk }

func (dingTalkChannel) DefaultRatePerMin() int { return 20 }

// DingTalk의 Webhook 주소에는 access_token이 있어 주소 자체가 자격 증명이므로 전체를 마스킹합니다.
func (dingTalkChannel) SecretKeys() []string { return []string{"webhook", "secret"} }

// 대상은 DingTalk의 Webhook 주소 자체입니다. 주소를 바꾸면 새 주소의 서명 키도 다시 명시해야 합니다.
func (dingTalkChannel) DestinationKeys() []string { return []string{"webhook"} }

func (dingTalkChannel) Validate(cfg map[string]any) error {
	hook := cfgString(cfg, "webhook")
	if hook == "" {
		return errors.New("Webhook 주소가 없습니다")
	}
	if err := validateHTTPURL(hook); err != nil {
		return fmt.Errorf("Webhook 주소가 유효하지 않습니다: %w", err)
	}
	return nil
}

// Send는 메시지를 한 번 전송합니다. 상세 링크가 있고 단일 메시지이면 ActionCard(버튼 포함), 그렇지 않으면 markdown을 사용합니다.
func (c dingTalkChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	hook := cfgString(cfg, "webhook")
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	endpoint, err := dingTalkSignedURL(hook, cfgString(cfg, "secret"), time.Now())
	if err != nil {
		return 0, Permanent(err)
	}

	title := markdownTitle(m)
	// DingTalk markdown 본문에는 명확한 바이트 상한이 없지만, 증거 필드가 비정상적으로 커지는 것을 막기 위해 상한 보호를 적용합니다.
	text, kept := markdownBody(m, 20000)

	var payload any
	if !m.Batch && len(m.Items) == 1 && m.Items[0].DetailURL != "" {
		payload = map[string]any{
			"msgtype": "actionCard",
			"actionCard": map[string]any{
				"title":          title,
				"text":           text,
				"btnOrientation": "0",
				"singleTitle":    "상세 보기",
				"singleURL":      m.Items[0].DetailURL,
			},
		}
	} else {
		payload = map[string]any{
			"msgtype":  "markdown",
			"markdown": map[string]any{"title": title, "text": text},
		}
	}

	raw, err := doJSON(ctx, "POST", endpoint, nil, payload)
	if err != nil {
		return 0, err
	}
	// DingTalk은 업무 오류를 200 응답에 넣습니다.
	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("DingTalk 응답 파싱 실패: %w (%s)", err, snippet(raw))
	}
	if res.ErrCode != 0 {
		// 301000은 서명 검증 실패, 310000은 키워드 불일치이며, 둘 다 설정 오류이므로,
		// 재시도로 해결되지 않습니다.
		return 0, Permanent(fmt.Errorf("DingTalk 오류 반환 %d: %s", res.ErrCode, res.ErrMsg))
	}
	return kept, nil
}

// dingTalkSignedURL은 공식 서명 추가 규칙에 따라 webhook에 timestamp와 sign 파라미터를 덧붙입니다.
//
// 규칙: 서명할 문자열 = timestamp + "\n" + secret. HMAC-SHA256의 **키도 secret**이며,
// 결과를 base64로 바꾼 뒤 URL 인코딩합니다. timestamp는 밀리초이며, secret이 비어 있으면 그대로 반환하여,
// 서명 추가를 켜지 않은 봇을 지원합니다.
func dingTalkSignedURL(hook, secret string, now time.Time) (string, error) {
	if secret == "" {
		return hook, nil
	}
	ts := strconv.FormatInt(now.UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "\n" + secret))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	u, err := url.Parse(hook)
	if err != nil {
		// err를 그대로 전달하지 않습니다. url.Parse 오류 텍스트에는 전체 주소(access_token 포함)가 있습니다.
		return "", fmt.Errorf("Webhook 주소 파싱 실패: %s", redactRequestTarget(hook))
	}
	q := u.Query()
	q.Set("timestamp", ts)
	q.Set("sign", sign)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// validateHTTPURL은 주소 사용 가능 여부와 지원 프로토콜을 검증하며, 리터럴 IP 대상이 내부 네트워크인지 판정합니다.
//
// 주의할 두 가지 사항:
//
//  1. **오류 정보는 반드시 민감정보를 제거해야 합니다**. url.Parse가 반환하는 *url.Error의 Error()에는
//     **원본 주소 전체**가 포함되며, 이 기능에서 사용하는 주소에는 자격 증명이 들어 있습니다(DingTalk access_token,
//     WeCom key, Telegram의 bot token, Feishu hook id). 이전에는 여기서 바로 `return err`를 했으므로,
//     '주소 형식이 유효하지 않음' 오류에 자격 증명이 포함되어 테스트 인터페이스의 400 응답으로 전달됐고,
//     매번 전송 시 DB에 저장되는 last_error, 서버 로그, 전송 이력 인터페이스로도 전달됐습니다.
//
//  2. **IP 리터럴은 바로 내부망 여부를 판정**하고, 도메인은 연결 단계에서 판정합니다(blockInternalDial이 최종
//     적용 지점이며 DNS 리바인딩도 처리). 여기서 한 번 검사하는 이유는 첫 전송 실패를 기다리지 않고,
//     설정 저장 시 안내를 받을 수 있게 하기 위해서입니다.
//
// 프로토콜 제한은 방어적 조치입니다. file:///gopher:// 등은 http.Client에 예상하지 못한
// 동작을 유발할 수 있습니다(scheme 검사가 이미 차단하지만, 허용할 이유가 없음).
func validateHTTPURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("주소를 파싱할 수 없습니다(%s)", redactRequestTarget(raw))
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("http/https만 지원합니다. 받은 값: %q", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("호스트 이름이 없습니다")
	}
	if ip := net.ParseIP(u.Hostname()); ip != nil && isBlockedDialIP(ip) && !allowLocalTargets() {
		return fmt.Errorf("로컬/link-local 주소로의 전송 거부: %s(로컬 서비스에 반드시 전송해야 한다면 %s=1 설정)", ip, AllowLocalTargetsEnv)
	}
	return nil
}
