package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// webhookChannel은 일반 Webhook 어댑터입니다. 사용자가 URL, 메서드, 요청 헤더와 JSON 템플릿을 지정합니다.
// 이 어댑터로 Slack / Mattermost / Discord / 자체 시스템마다 구현을 따로 작성할 필요 없이,
// 설정 가능한 템플릿 하나로 해당 플랫폼을 모두 지원할 수 있습니다.
type webhookChannel struct{}

func (webhookChannel) Kind() string { return KindWebhook }

// 일반 Webhook에는 공식 제한이 없습니다. 0 반환은 기본적으로 전송 속도를 제한하지 않음을 뜻하며, 사용자가 상대 측 성능에 맞춰 지정합니다.
func (webhookChannel) DefaultRatePerMin() int { return 0 }

// url과 headers를 마스킹합니다. 대상 주소 자체에 token이 포함되는 경우가 많고 사용자 지정 헤더에는 보통 인증·인가 확인 자격 증명을 넣으며,
// 둘 다 API 응답에 나타나므로 모두 막아야 합니다.
// 그 대신 편집 시 헤더 하나를 바꾸려면 헤더 전체를 다시 입력해야 합니다(마스킹된 값은 "기존 값 유지"로 해석).
// 이 선택은 의도적입니다. 한 번 더 입력하더라도 자격 증명을 브라우저에 반환하지 않습니다.
func (webhookChannel) SecretKeys() []string { return []string{"url", "headers"} }

// 대상은 url입니다. url을 바꾸면 headers를 다시 명시해야 합니다. 그렇지 않으면 기존 Authorization 헤더를
// 새 주소로 그대로 보내게 되며, 이는 마스킹 우회의 주요 경로입니다.
func (webhookChannel) DestinationKeys() []string { return []string{"url"} }

// webhookDefaultTemplate은 템플릿 미지정 시 사용하는 기본 요청 본문입니다. 단순한 JSON 구조로,
// 'JSON 한 건을 받아 DB에 저장'하는 대부분의 자체 구축 수신처를 처리합니다.
const webhookDefaultTemplate = `{
  "title": {{json .Title}},
  "batch": {{.Batch}},
  "count": {{.Count}},
  "items": [
{{- range $i, $it := .Items}}
{{- if $i}},{{end}}
    {
      "finding_id": {{$it.FindingID}},
      "name": {{json $it.Name}},
      "vulnclass": {{json $it.VulnClass}},
      "severity": {{json $it.Severity}},
      "summary": {{json $it.Summary}},
      "assets": {{json $it.Assets}},
      "detail_url": {{json $it.DetailURL}}
    }
{{- end}}
  ]
}`

// webhookTemplateData는 사용자 템플릿에 노출하는 컨텍스트입니다.
type webhookTemplateData struct {
	Title   string
	Batch   bool
	Count   int
	Items   []webhookItem
	HomeURL string
	// SentAt은 이번 전송 시간(RFC3339)이며, 수신 측에서 기록할 수 있습니다.
	SentAt string
}

type webhookItem struct {
	FindingID     int64
	Name          string
	VulnClass     string
	Severity      string
	SeverityLabel string
	Summary       string
	Assets        []string
	DetailURL     string
	FromStatus    string
	ToStatus      string
	// StatusLabel은 "처리 대기 → 수정 완료"와 같은 읽기 쉬운 상태 변경 설명입니다. 상태 변경이 아니면 비어 있습니다.
	StatusLabel string
}

func (webhookChannel) Validate(cfg map[string]any) error {
	raw := cfgString(cfg, "url")
	if raw == "" {
		return errors.New("대상 URL이 없습니다")
	}
	if err := validateHTTPURL(raw); err != nil {
		return fmt.Errorf("대상 URL이 유효하지 않습니다: %w", err)
	}
	if m := strings.ToUpper(cfgString(cfg, "method")); m != "" && m != http.MethodGet && m != http.MethodPost && m != http.MethodPut && m != http.MethodPatch {
		return fmt.Errorf("지원하지 않는 메서드 %s(사용 가능: GET/POST/PUT/PATCH)", m)
	}
	if tpl := cfgString(cfg, "body_template"); tpl != "" {
		if _, err := parseWebhookTemplate(tpl); err != nil {
			return fmt.Errorf("요청 본문 템플릿 구문 오류: %w", err)
		}
	}
	return nil
}

func (c webhookChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	method := strings.ToUpper(cfgString(cfg, "method"))
	if method == "" {
		method = http.MethodPost
	}

	// GET에는 요청 본문을 넣지 않습니다. 내용을 query에 넣는 것은 템플릿 기능 범위를 벗어나며 GET의 의미에도 맞지 않으므로,
	// GET은 "매칭되면 바로 훅을 실행"하는 수신 측에만 적합합니다.
	var payload any
	if method != http.MethodGet {
		body, err := renderWebhookBody(cfgString(cfg, "body_template"), m)
		if err != nil {
			return 0, Permanent(err)
		}
		// 템플릿이 렌더링한 것은 문자열 형태의 JSON이므로, json.RawMessage로 바꿔 그대로 전송하여,
		// 이중 이스케이프로 사용자가 구성한 구조를 하나의 JSON 문자열 안에 넣는 일을 피합니다.
		if !json.Valid([]byte(body)) {
			return 0, Permanent(errors.New("요청 본문 템플릿 렌더링 결과가 유효한 JSON이 아닙니다"))
		}
		payload = json.RawMessage(body)
	}

	headers := cfgMap(cfg, "headers")
	if ct := cfgString(cfg, "content_type"); ct != "" {
		// 덮어쓰기를 허용하되 headers 다음에 적용하여 명시적 설정이 우선하도록 합니다.
		if headers == nil {
			headers = map[string]string{}
		}
		headers["Content-Type"] = ct
	}
	if _, err := doJSON(ctx, method, cfgString(cfg, "url"), headers, payload); err != nil {
		return 0, err
	}
	// 일반 Webhook은 본문을 자르지 않습니다(수신 측은 사용자의 서비스이며 크기는 body_template로 결정).
	// 따라서 배치 전체를 전달된 것으로 처리합니다.
	return len(m.Items), nil
}

// renderWebhookBody는 사용자 템플릿(또는 기본 템플릿)으로 요청 본문을 렌더링합니다.
func renderWebhookBody(tpl string, m Message) (string, error) {
	if strings.TrimSpace(tpl) == "" {
		tpl = webhookDefaultTemplate
	}
	t, err := parseWebhookTemplate(tpl)
	if err != nil {
		return "", fmt.Errorf("요청 본문 템플릿 구문 오류: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, newWebhookTemplateData(m)); err != nil {
		return "", fmt.Errorf("요청 본문 템플릿 렌더링 실패: %w", err)
	}
	return buf.String(), nil
}

// parseWebhookTemplate은 템플릿을 파싱합니다.
//
// missingkey=zero는 누락된 map 키를 오류 대신 영값으로 렌더링합니다. 다만 이 파일의 컨텍스트는 구조체이며,
// 주된 역할은 .Items가 비어 있어도 range에서 오류가 나지 않게 하는 것입니다. 실제로 방지해야 하는 것은 .Items가 nil인 경우입니다.
func parseWebhookTemplate(tpl string) (*template.Template, error) {
	return template.New("body").Funcs(webhookTemplateFuncs).Option("missingkey=zero").Parse(tpl)
}

// webhookTemplateFuncs는 템플릿에 노출하는 보조 함수입니다.
var webhookTemplateFuncs = template.FuncMap{
	// json은 임의의 값을 JSON으로 직렬화합니다.
	//
	// 이 함수는 선택적인 개선이 아니라 필수입니다. 없으면 사용자는 {{.Title}}로 직접 값을 삽입할 수밖에 없으며,
	// 취약점 제목에 따옴표나 줄바꿈이 있으면 요청 본문 전체가 유효한 JSON이 아니게 되어 수신 측에서
	// 거부합니다. 오류는 "JSON 파싱 실패"를 가리키므로 제목에 따옴표가 있다는 원인은 전혀 떠올리기 어렵습니다.
	"json": func(v any) (string, error) {
		raw, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(raw), nil
	},
	// jsons는 JSON 조각을 다른 JSON 문자열 값 안에 넣을 때 사용합니다(한 번의 문자열 이스케이프 수행).
	"jsons": func(v any) (string, error) {
		raw, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		quoted, err := json.Marshal(string(raw))
		if err != nil {
			return "", err
		}
		// 바깥 따옴표를 제거합니다. 따옴표를 붙일지는 호출자가 결정합니다.
		return string(quoted[1 : len(quoted)-1]), nil
	},
}

func newWebhookTemplateData(m Message) webhookTemplateData {
	d := webhookTemplateData{
		Title:   markdownTitle(m),
		Batch:   m.Batch,
		Count:   len(m.Items),
		HomeURL: m.HomeURL,
		SentAt:  time.Now().Format(time.RFC3339),
		Items:   make([]webhookItem, 0, len(m.Items)),
	}
	for _, it := range m.Items {
		wi := webhookItem{
			FindingID:     it.FindingID,
			Name:          it.Name,
			VulnClass:     it.VulnClass,
			Severity:      it.Severity,
			SeverityLabel: SeverityLabel(it.Severity),
			Summary:       it.Summary,
			Assets:        append([]string{}, it.Assets...),
			DetailURL:     it.DetailURL,
			FromStatus:    it.FromStatus,
			ToStatus:      it.ToStatus,
		}
		if it.IsStatusChange() {
			wi.StatusLabel = StatusLabel(it.FromStatus) + " → " + StatusLabel(it.ToStatus)
		}
		d.Items = append(d.Items, wi)
	}
	return d
}
