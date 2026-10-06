package notify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// feishuChannel은 Feishu(Lark 포함) 사용자 지정 봇을 구현하며, 대화형 카드를 사용합니다.
//
// 플랫폼 특성:
//   - 서명 추가 알고리즘이 DingTalk과 **다르고**, 잘못 작성하기 쉬우므로 feishuSign 주석을 참조하세요.
//   - DingTalk과 마찬가지로 업무 오류를 HTTP 200의 body에 넣습니다(code != 0).
//   - 카드 header는 색상 템플릿을 지원합니다. 심각도를 색상에 매핑하여 메시지 목록에서 심각도를 한눈에 알아볼 수 있게 합니다.
type feishuChannel struct{}

func (feishuChannel) Kind() string { return KindFeishu }

// Feishu 사용자 지정 봇은 약 초당 5회이며, 분당 100회로 환산합니다.
func (feishuChannel) DefaultRatePerMin() int { return 100 }

// Webhook 주소 마지막 부분은 봇 고유 식별자이므로 자격 증명에 해당합니다.
func (feishuChannel) SecretKeys() []string { return []string{"webhook", "secret"} }

// 마찬가지로 Webhook 주소를 바꾸면 새 주소의 서명 키도 다시 명시해야 합니다.
func (feishuChannel) DestinationKeys() []string { return []string{"webhook"} }

func (feishuChannel) Validate(cfg map[string]any) error {
	hook := cfgString(cfg, "webhook")
	if hook == "" {
		return errors.New("Webhook 주소가 없습니다")
	}
	if err := validateHTTPURL(hook); err != nil {
		return fmt.Errorf("Webhook 주소가 유효하지 않습니다: %w", err)
	}
	return nil
}

func (c feishuChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	card, kept := feishuCard(m)
	payload := map[string]any{
		"msg_type": "interactive",
		"card":     card,
	}
	// 서명 추가 파라미터는 메시지와 같은 계층에 있으며, secret을 설정한 경우에만 포함합니다.
	if secret := cfgString(cfg, "secret"); secret != "" {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		payload["timestamp"] = ts
		payload["sign"] = feishuSign(ts, secret)
	}
	raw, err := doJSON(ctx, "POST", cfgString(cfg, "webhook"), nil, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		// 일부 Feishu hook 버전은 이 필드 이름을 사용하므로 함께 지원합니다.
		StatusCode    int    `json:"StatusCode"`
		StatusMessage string `json:"StatusMessage"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("Feishu 응답 파싱 실패: %w (%s)", err, snippet(raw))
	}
	if res.Code != 0 {
		return 0, Permanent(fmt.Errorf("Feishu 오류 반환 %d: %s", res.Code, res.Msg))
	}
	if res.StatusCode != 0 {
		return 0, Permanent(fmt.Errorf("Feishu 오류 반환 %d: %s", res.StatusCode, res.StatusMessage))
	}
	return kept, nil
}

// feishuSign은 Feishu 공식 규칙에 따라 서명을 계산합니다.
//
// 실수하기 쉬운 부분입니다. 공식 예시는 다음과 같습니다.
//
//	hmac.new(string_to_sign.encode(), digestmod=sha256)
//
// 즉, **key = timestamp + "\n" + secret이며, message는 비어 있습니다**. 직관적인
// "key=secret, message=stringToSign" 방식은 DingTalk의 알고리즘입니다. 두 알고리즘은 정확히 반대이므로,
// 다른 쪽 구현을 그대로 따라 쓰면 반드시 서명 검증에 실패합니다(19021 반환).
func feishuSign(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret
	mac := hmac.New(sha256.New, []byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// feishuSeverityTemplate은 취약점 심각도를 카드 header 색상 템플릿에 매핑합니다.
// 알 수 없는 심각도는 grey를 사용합니다. low와 혼동하지 않도록 blue는 사용하지 않습니다.
func feishuSeverityTemplate(severity string) string {
	switch severity {
	case "critical":
		return "red"
	case "high":
		return "orange"
	case "medium":
		return "yellow"
	case "low":
		return "blue"
	default:
		return "grey"
	}
}

// feishuMaxCardBytes는 카드 내용의 보수적인 상한입니다. Feishu는 카드 크기를 제한하며, 초과하면 메시지 전체를 거부합니다.
// 공식 상한보다 확실히 작은 값을 사용하여 JSON으로 감싸는 데 드는 크기도 포함합니다.
const feishuMaxCardBytes = 24000

// feishuCard는 대화형 카드를 구성하며, 카드와 **실제로 작성한 항목 수**를 반환합니다.
// kept의 용도는 markdownBody와 같습니다. 실제로 카드에 들어간 항목만 전달됨으로 표시해야 합니다.
func feishuCard(m Message) (map[string]any, int) {
	elements := []any{}
	kept := 0
	if m.Batch {
		// 먼저 항목 전체 단위로 묶은 뒤 머리를 구성합니다. 머리에는 "나머지 N건은 다음 메시지에서 이어집니다"를 써야 하며,
		// N은 실제로 담은 건수에서 계산해야 합니다.
		kept = packItemCount(m.Items, feishuMaxCardBytes, markdownReservedBytes, "", byteSize, func(it Item, idx int) string {
			return feishuBatchLine(it, idx+1)
		})
		items := m.Items[:kept]
		elements = append(elements, feishuMarkdownDiv(markdownBatchIntro(m, items, len(m.Items))))
		for i, it := range items {
			elements = append(elements, feishuMarkdownDiv(feishuBatchLine(it, i+1)))
		}
		if m.HomeURL != "" {
			elements = append(elements, feishuButton("플랫폼에서 전체 보기", m.HomeURL))
		}
	} else if len(m.Items) > 0 {
		kept = 1
		it := m.Items[0]
		elements = append(elements, feishuMarkdownDiv(feishuItemLines(it)))
		if it.DetailURL != "" {
			elements = append(elements, feishuButton("상세 보기", it.DetailURL))
		}
	}

	card := map[string]any{
		"config":   map[string]any{"wide_screen_mode": true},
		"header":   map[string]any{"title": map[string]any{"tag": "plain_text", "content": markdownTitle(m)}},
		"elements": elements,
	}
	if len(m.Items) > 0 {
		card["header"].(map[string]any)["template"] = feishuSeverityTemplate(m.Items[0].Severity)
	}
	return card, kept
}

func feishuMarkdownDiv(content string) map[string]any {
	return map[string]any{"tag": "div", "text": map[string]any{"tag": "lark_md", "content": content}}
}

func feishuButton(label, url string) map[string]any {
	return map[string]any{
		"tag": "action",
		"actions": []any{map[string]any{
			"tag":  "button",
			"text": map[string]any{"tag": "lark_md", "content": label},
			"url":  url,
			"type": "primary",
		}},
	}
}

// feishuItemLines는 취약점 하나의 lark_md 본문을 렌더링합니다.
//
// lark_md와 markdown은 같은 계열의 텍스트 형식이며, 링크와 강조도 파싱합니다. 따라서 외부에서 온
// 필드는 모두 markdownText(한 줄로 압축 + 이스케이프)를 거쳐야 합니다. 그렇지 않으면 취약점 제목 하나가
// Feishu에서 클릭 가능한 외부 링크로 바뀔 수 있습니다.
func feishuItemLines(it Item) string {
	out := fmt.Sprintf("**%s · %s**", SeverityLabel(it.Severity), markdownText(it.Title(), 0))
	if it.IsStatusChange() {
		out += fmt.Sprintf("\n**상태 변경**: %s → %s",
			markdownText(StatusLabel(it.FromStatus), 0), markdownText(StatusLabel(it.ToStatus), 0))
	}
	if it.VulnClass != "" && it.VulnClass != it.Title() {
		out += fmt.Sprintf("\n**유형**: %s", markdownText(it.VulnClass, 0))
	}
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		out += fmt.Sprintf("\n**자산**: %s", markdownText(a, 0))
	}
	if it.Summary != "" {
		if s := markdownText(it.Summary, maxSummaryRunes); s != "" {
			out += fmt.Sprintf("\n**요약**: %s", s)
		}
	}
	return out
}

// feishuBatchLine은 모아 보내기 카드의 한 항목을 렌더링합니다.
func feishuBatchLine(it Item, index int) string {
	line := fmt.Sprintf("**%d. %s · %s**", index, SeverityLabel(it.Severity), markdownText(it.Title(), 0))
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		line += " — " + markdownText(a, 0)
	}
	return line
}
