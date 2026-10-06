package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// weComMarkdownLimit는 WeCom(기업용 위챗) 그룹 봇 markdown content의 엄격한 상한입니다(바이트이며, 문자가 아님).
// 여섯 채널 중 가장 엄격한 제한이며, TruncateBytes가 존재하는 주요 이유입니다.
const weComMarkdownLimit = 4096

// weComChannel은 WeCom 그룹 봇을 구현합니다.
//
// 플랫폼 특성:
//   - URL의 key로만 인증·인가 확인을 하며 서명 추가를 지원하지 않습니다. 따라서 webhook 주소 자체가 자격 증명 전부입니다.
//   - markdown content 상한은 4096 **바이트**이며, 초과하면 메시지 전체를 거부합니다(잘리는 것이 아님). 중국어는 글자당 3바이트이므로,
//     본문에는 천여 글자만 쓸 수 있어 반드시 클라이언트에서 잘라야 합니다.
//   - 전송 속도 제한은 분당 20건이며, 마찬가지로 클라이언트에서 제한해야 합니다.
type weComChannel struct{}

func (weComChannel) Kind() string { return KindWeCom }

func (weComChannel) DefaultRatePerMin() int { return 20 }

// WeCom에는 Webhook 한 곳의 자격 증명만 있고(URL의 key), 서명 추가를 지원하지 않으므로,
// 주소 전체가 자격 증명 전부이며 마스킹할 다른 필드는 없습니다.
func (weComChannel) SecretKeys() []string { return []string{"webhook"} }

// WeCom은 Webhook 필드 하나뿐이며 대상이자 자격 증명이므로, "주소 변경 후 남아 있는 자격 증명"은 없습니다.
func (weComChannel) DestinationKeys() []string { return []string{"webhook"} }

func (weComChannel) Validate(cfg map[string]any) error {
	hook := cfgString(cfg, "webhook")
	if hook == "" {
		return errors.New("Webhook 주소가 없습니다")
	}
	if err := validateHTTPURL(hook); err != nil {
		return fmt.Errorf("Webhook 주소가 유효하지 않습니다: %w", err)
	}
	return nil
}

func (c weComChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	// 모아 보내기 배치는 길어질 수 있어(50건 × 건당 한 줄 + 접두사), 4096바이트를 쉽게 넘습니다.
	// 플랫폼 오류에 기대지 않고 여기서 자릅니다. 거부되면 배치 전체를 잃지만, 자르면 최소한 처음 일부 항목은 전달됩니다.
	content, kept := markdownBody(m, weComMarkdownLimit)
	payload := map[string]any{
		"msgtype":  "markdown",
		"markdown": map[string]any{"content": content},
	}
	raw, err := doJSON(ctx, "POST", cfgString(cfg, "webhook"), nil, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("WeCom 응답 파싱 실패: %w (%s)", err, snippet(raw))
	}
	if res.ErrCode != 0 {
		// 45009는 인터페이스 호출 제한 초과입니다. 플랫폼의 전송 속도 제한 시간 구간은 이동하므로 백오프 후 재시도가 효과적이며,
		// 따라서 명시적으로 재시도 가능한 것으로 분류합니다. 여기에 도달했다면 클라이언트 rate_per_min 설정이 지나치게 높다는 뜻이며,
		// 재시도는 보완 수단일 뿐, 실제 해결 방법은 해당 채널의 전송 속도 제한 값을 낮추는 것입니다.
		if res.ErrCode == 45009 {
			return 0, fmt.Errorf("WeCom 전송 속도 제한 %d: %s", res.ErrCode, res.ErrMsg)
		}
		// 93000은 webhook key가 유효하지 않다는 뜻으로, 영구 실패이며 재시도로 해결되지 않습니다.
		return 0, Permanent(fmt.Errorf("WeCom 오류 반환 %d: %s", res.ErrCode, res.ErrMsg))
	}
	return kept, nil
}
