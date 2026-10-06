package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// telegramTextLimit는 Telegram sendMessage의 text 필드 상한입니다(문자 수).
const telegramTextLimit = 4096

// telegramChannel은 Telegram Bot API를 구현합니다.
//
// 플랫폼 특성:
//   - 인증·인가 확인은 전부 URL path(/bot<token>/sendMessage)에 있으며, 서명 추가는 필요하지 않습니다.
//   - MarkdownV2 대신 HTML 파싱 모드를 사용합니다. MarkdownV2는 `_*[]()~`>#+-=|{}.!`의
//     총 18문자를 이스케이프해야 하며, 하나라도 빠지면 메시지 전체가 거부됩니다. HTML은 & < > 세 문자만 이스케이프하면 됩니다.
//   - 업무 오류는 마찬가지로 HTTP 200에 포함되며, ok 필드로 판정합니다.
type telegramChannel struct{}

func (telegramChannel) Kind() string { return KindTelegram }

// Telegram은 개인 대화 약 초당 1건, 그룹은 분당 20건이며, 보수적인 값을 사용합니다.
func (telegramChannel) DefaultRatePerMin() int { return 20 }

// Bot Token이 자격 증명 전부이며, chat_id는 수신자일 뿐 비밀 정보가 아닙니다(알아도 Token이 없으면 메시지를 보낼 수 없음).
func (telegramChannel) SecretKeys() []string { return []string{"bot_token"} }

// base_url은 Token을 어느 API 엔드포인트(자체 구축 리버스 프록시 등)로 전송할지 결정하므로, 변경 시 Token을 반드시 다시 지정해야 합니다.
func (telegramChannel) DestinationKeys() []string { return []string{"base_url"} }

func (telegramChannel) Validate(cfg map[string]any) error {
	if cfgString(cfg, "bot_token") == "" {
		return errors.New("Bot Token이 없습니다")
	}
	if cfgString(cfg, "chat_id") == "" {
		return errors.New("Chat ID가 없습니다")
	}
	if base := cfgString(cfg, "base_url"); base != "" {
		if err := validateHTTPURL(base); err != nil {
			return fmt.Errorf("API 주소가 유효하지 않습니다: %w", err)
		}
	}
	return nil
}

func (c telegramChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	endpoint, err := telegramEndpoint(cfg)
	if err != nil {
		return 0, Permanent(err)
	}
	text, kept := telegramHTML(m)
	payload := map[string]any{
		"chat_id":                  cfgString(cfg, "chat_id"),
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": false,
	}
	raw, err := doJSON(ctx, "POST", endpoint, nil, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("Telegram 응답 파싱 실패: %w (%s)", err, snippet(raw))
	}
	if res.OK {
		return kept, nil
	}
	// 429는 전송 속도 제한으로, 백오프 후 재시도가 유효합니다. 나머지(400 파라미터 오류, 401 token 오류, 403 차단,
	// 404 chat 없음)는 모두 설정 문제이며 재시도로 해결되지 않습니다.
	if res.ErrorCode == 429 {
		return 0, fmt.Errorf("Telegram 전송 속도 제한: %s", res.Description)
	}
	return 0, Permanent(fmt.Errorf("Telegram 오류 반환 %d: %s", res.ErrorCode, res.Description))
}

// telegramEndpoint는 sendMessage 주소를 조합합니다. base_url이 비어 있으면 공식 API를 사용하고,
// 비어 있지 않으면 자체 구축 Bot API 리버스 프록시에 사용합니다(중국 내 네트워크에서 흔한 요구 사항).
func telegramEndpoint(cfg map[string]any) (string, error) {
	base := cfgString(cfg, "base_url")
	if base == "" {
		base = "https://api.telegram.org"
	}
	base = strings.TrimSuffix(base, "/")
	token := cfgString(cfg, "bot_token")
	raw := base + "/bot" + token + "/sendMessage"
	u, err := url.Parse(raw)
	if err != nil {
		// err를 그대로 전달하지 않습니다. 주소에 Bot Token이 포함되며, 이 시점에는 addr조차 반환해서는 안 됩니다.
		return "", fmt.Errorf("API 주소 조합 실패(API 주소: %s)", redactRequestTarget(base))
	}
	return u.String(), nil
}

// telegramHTML은 HTML 본문을 렌더링하며, 본문과 실제로 작성한 항목 수를 반환합니다(Channel.Send 참조).
func telegramHTML(m Message) (string, int) {
	var b strings.Builder
	b.WriteString("<b>" + telegramEscape(markdownTitle(m)) + "</b>\n")
	if m.Batch {
		// Telegram 상한은 **문자 수**이므로, 묶을 때도 문자 단위로 측정합니다(runeSize).
		footer := ""
		if m.HomeURL != "" {
			footer = fmt.Sprintf("\n\n<a href=\"%s\">플랫폼에서 전체 보기</a>", telegramEscapeAttr(m.HomeURL))
		}
		kept := packItemCount(m.Items, telegramTextLimit, telegramReservedRunes, footer, runeSize, func(it Item, idx int) string {
			return telegramBatchLine(it, idx+1)
		})
		items := m.Items[:kept]
		b.Reset()
		b.WriteString("<b>" + telegramEscape(telegramBatchTitle(m, items, len(m.Items))) + "</b>")
		for i, it := range items {
			b.WriteString("\n" + telegramEscape(telegramBatchLine(it, i+1)))
		}
		b.WriteString(footer)
		return TruncateHTML(b.String(), telegramTextLimit), kept
	}
	if len(m.Items) == 0 {
		return b.String(), 0
	}
	it := m.Items[0]
	if it.IsStatusChange() {
		b.WriteString(fmt.Sprintf("\n<b>상태 변경</b>: %s → %s",
			telegramEscape(StatusLabel(it.FromStatus)), telegramEscape(StatusLabel(it.ToStatus))))
	}
	if it.VulnClass != "" && it.VulnClass != it.Title() {
		b.WriteString("\n<b>유형</b>: " + telegramEscape(it.VulnClass))
	}
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		b.WriteString("\n<b>자산</b>: " + telegramEscape(a))
	}
	if s := OneLine(it.Summary, maxSummaryRunes); s != "" {
		b.WriteString("\n<b>요약</b>: " + telegramEscape(s))
	}
	if it.DetailURL != "" {
		b.WriteString(fmt.Sprintf("\n\n<a href=\"%s\">상세 보기</a>", telegramEscapeAttr(it.DetailURL)))
	}
	return TruncateHTML(b.String(), telegramTextLimit), 1
}

// telegramReservedRunes는 메시지 제목과 발생할 수 있는 잘림 안내에 미리 할당한 크기입니다(문자 수 기준).
const telegramReservedRunes = 160

// telegramBatchLine은 모아 보내기의 한 항목을 렌더링합니다(이스케이프하지 않으며, 호출자가 통일하여 처리).
func telegramBatchLine(it Item, idx int) string {
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		return fmt.Sprintf("%d. %s · %s — %s", idx, SeverityLabel(it.Severity), it.Title(), a)
	}
	return fmt.Sprintf("%d. %s · %s", idx, SeverityLabel(it.Severity), it.Title())
}

// telegramBatchTitle은 모아 보내기 메시지 제목 줄을 렌더링합니다. 건수는 **이 메시지에 실제로 포함한** 건수를 사용하며,
// 배치 전체 건수가 아닙니다. 그렇지 않으면 독자는 메시지 머리의 숫자가 전부라고 생각하게 됩니다.
func telegramBatchTitle(m Message, items []Item, total int) string {
	title := fmt.Sprintf("취약점 모아 보내기 · 총 %d건", total)
	if extra := total - len(items); extra > 0 {
		title += fmt.Sprintf("(처음 %d건을 표시하며, 나머지 %d건은 다음 메시지에서 이어집니다)", len(items), extra)
	}
	if m.WindowMinutes > 0 {
		title = fmt.Sprintf("최근 %d분 · %s", m.WindowMinutes, title)
	}
	return title
}

// telegramEscape는 HTML 텍스트 내용을 이스케이프합니다.
// Telegram은 이 세 엔티티만 인식합니다. 이스케이프하면 &amp;와 같은 기존 엔티티도 다시 이스케이프되며, 이는
// 의도한 동작입니다. 원래 문자를 표시하려는 것이지 사용자가 HTML을 인젝션하도록 허용하려는 것이 아닙니다.
func telegramEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// telegramEscapeAttr은 HTML 속성값을 이스케이프합니다. 텍스트 이스케이프에 더해 따옴표도 처리해야 하며,
// URL의 따옴표가 href 속성을 먼저 닫아 뒤의 내용을 인젝션 위치로 만들 수 있습니다.
func telegramEscapeAttr(s string) string {
	s = telegramEscape(s)
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
