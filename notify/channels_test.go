package notify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"
)

// singleMsg는 따옴표와 줄바꿈이 있는 단일 메시지를 구성합니다. 의도적으로 `"`와 `\n`이 있는 제목/요약을 사용합니다.
// 이는 템플릿 값 삽입에서 유효하지 않은 JSON을 가장 쉽게 만드는 입력입니다.
func singleMsg() Message {
	return Message{
		Items: []Item{{
			FindingID: 42,
			Name:      `登录处 "SQL注入" 风险`,
			VulnClass: "SQL注入",
			Severity:  "high",
			Summary:   "参数 id\n未过滤 导致注入",
			Assets:    []string{"a.example.com", "b.example.com"},
			DetailURL: "https://artex.local/function/findings/detail?id=42",
		}},
	}
}

// batchMsg는 모아 보내기 메시지 배치를 구성합니다.
func batchMsg(n int) Message {
	m := Message{Batch: true, WindowMinutes: 30, HomeURL: "https://artex.local/function/findings"}
	for i := 0; i < n; i++ {
		m.Items = append(m.Items, Item{
			FindingID: int64(i + 1),
			Name:      "漏洞" + itoa(i+1),
			VulnClass: "XSS",
			Severity:  "medium",
			Summary:   "反射型跨站脚本",
			Assets:    []string{"target.example.com"},
		})
	}
	return m
}

// capturePost는 가상 수신 측을 시작하여 받은 요청 본문과 헤더를 검사 함수에 전달합니다.
func capturePost(t *testing.T, respBody string, assert func(t *testing.T, body map[string]any, r *http.Request)) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Errorf("요청 본문이 유효한 JSON이 아닙니다: %v\n원문: %s", err, raw)
			}
		}
		if assert != nil {
			assert(t, body, r)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, respBody)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestDingTalkSendsActionCardWhenLinkPresent(t *testing.T) {
	srv := capturePost(t, `{"errcode":0,"errmsg":"ok"}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		if body["msgtype"] != "actionCard" {
			t.Fatalf("상세 링크가 있으면 actionCard를 보내야 합니다. 실제: %v", body["msgtype"])
		}
		card, _ := body["actionCard"].(map[string]any)
		if card["singleURL"] != "https://artex.local/function/findings/detail?id=42" {
			t.Errorf("상세 링크 누락: %v", card["singleURL"])
		}
	})
	if _, err := (dingTalkChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

func TestDingTalkFallsBackToMarkdownForBatch(t *testing.T) {
	srv := capturePost(t, `{"errcode":0,"errmsg":"ok"}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		if body["msgtype"] != "markdown" {
			t.Fatalf("모아 보내기 메시지는 markdown으로 보내야 합니다. 실제: %v", body["msgtype"])
		}
		md, _ := body["markdown"].(map[string]any)
		if !strings.Contains(md["text"].(string), "최근 30분") {
			t.Errorf("모아 보내기 본문에 시간 구간 누락: %v", md["text"])
		}
	})
	if _, err := (dingTalkChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, batchMsg(3)); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

// TestDingTalkBusinessErrorIsPermanent는 "HTTP 200이지만 errcode가 0이 아님" 판정을 고정합니다.
// errcode를 확인하지 않으면 전송 실패를 성공으로 기록하게 됩니다. 이는 중국 IM 플랫폼들의 공통적인 주의점입니다.
func TestDingTalkBusinessErrorIsPermanent(t *testing.T) {
	srv := capturePost(t, `{"errcode":310000,"errmsg":"keywords not in content"}`, nil)
	_, err := (dingTalkChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, singleMsg())
	if err == nil {
		t.Fatal("errcode가 0이 아니면 오류가 나야 합니다")
	}
	if !IsPermanent(err) {
		t.Fatalf("키워드 불일치는 설정 오류이므로 영구 실패로 표시해야 합니다. 실제: %v", err)
	}
	if !strings.Contains(err.Error(), "310000") {
		t.Errorf("오류 메시지에 플랫폼 오류 코드를 포함해야 합니다. 실제: %v", err)
	}
}

func TestWeComTruncatesCJKWithinByteLimit(t *testing.T) {
	var contentLen int
	srv := capturePost(t, `{"errcode":0,"errmsg":"ok"}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		md, _ := body["markdown"].(map[string]any)
		content, _ := md["content"].(string)
		contentLen = len(content)
		if !utf8.ValidString(content) {
			t.Fatal("자른 후 유효한 UTF-8이 아니므로 WeCom이 메시지 전체를 거부합니다")
		}
	})
	// 충분히 긴 중국어 모아 보내기 배치를 만들어 4096바이트를 반드시 넘게 합니다.
	m := batchMsg(200)
	if _, err := (weComChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, m); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
	if contentLen > weComMarkdownLimit {
		t.Fatalf("본문 %d바이트가 WeCom 상한 %d 초과", contentLen, weComMarkdownLimit)
	}
	if contentLen == 0 {
		t.Fatal("본문이 비어 있습니다")
	}
}

func TestWeComRateLimitIsRetryableButKeyErrorIsPermanent(t *testing.T) {
	limited := capturePost(t, `{"errcode":45009,"errmsg":"api freq out of limit"}`, nil)
	_, err := (weComChannel{}).Send(context.Background(), map[string]any{"webhook": limited.URL}, singleMsg())
	if err == nil || IsPermanent(err) {
		t.Fatalf("45009는 이동하는 시간 구간의 전송 속도 제한이므로 재시도할 수 있어야 합니다. 실제: %v", err)
	}

	badKey := capturePost(t, `{"errcode":93000,"errmsg":"invalid webhook url"}`, nil)
	_, err = (weComChannel{}).Send(context.Background(), map[string]any{"webhook": badKey.URL}, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("93000은 key가 유효하지 않음으로, 재시도로 해결되지 않아 영구 실패여야 합니다. 실제: %v", err)
	}
}

func TestFeishuCardStructureAndSign(t *testing.T) {
	const secret = "SECtest123"
	srv := capturePost(t, `{"code":0,"msg":"success"}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		if body["msg_type"] != "interactive" {
			t.Fatalf("대화형 카드를 보내야 합니다. 실제: %v", body["msg_type"])
		}
		card, _ := body["card"].(map[string]any)
		header, _ := card["header"].(map[string]any)
		if header["template"] != "orange" {
			t.Errorf("high 심각도는 orange 색상이어야 합니다. 실제: %v", header["template"])
		}
		// secret을 설정하면 서명 추가 파라미터를 반드시 포함해야 합니다. 그렇지 않으면 Feishu가 19021로 거부합니다.
		if body["sign"] == nil || body["timestamp"] == nil {
			t.Fatalf("서명 추가 파라미터 누락: %v", body)
		}
		// 카드 요소에는 취약점 상세 내용을 가리키는 url의 버튼이 하나 있어야 합니다.
		elements, _ := card["elements"].([]any)
		foundButton := false
		for _, e := range elements {
			em, _ := e.(map[string]any)
			if em["tag"] != "action" {
				continue
			}
			actions, _ := em["actions"].([]any)
			for _, a := range actions {
				am, _ := a.(map[string]any)
				if am["url"] == "https://artex.local/function/findings/detail?id=42" {
					foundButton = true
				}
			}
		}
		if !foundButton {
			t.Fatal("카드에 상세 페이지로 연결되는 버튼이 없습니다")
		}
	})
	cfg := map[string]any{"webhook": srv.URL, "secret": secret}
	if _, err := (feishuChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

func TestFeishuWithoutSecretOmitsSign(t *testing.T) {
	srv := capturePost(t, `{"code":0,"msg":"success"}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		if body["sign"] != nil || body["timestamp"] != nil {
			t.Fatalf("secret을 설정하지 않으면 서명 추가 파라미터를 포함해서는 안 됩니다: %v", body)
		}
	})
	if _, err := (feishuChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

func TestTelegramEscapesHTMLInUntrustedContent(t *testing.T) {
	var text string
	srv := capturePost(t, `{"ok":true}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		text, _ = body["text"].(string)
		if body["parse_mode"] != "HTML" {
			t.Fatalf("HTML 파싱 모드를 사용해야 합니다. 실제: %v", body["parse_mode"])
		}
	})
	m := Message{Items: []Item{{
		Severity: "high",
		// 제목과 요약은 테스트 대상/모델 출력에서 온 신뢰할 수 없는 내용입니다.
		Name:    `<script>alert(1)</script>`,
		Summary: "a & b < c",
	}}}
	if _, err := (telegramChannel{}).Send(context.Background(),
		map[string]any{"bot_token": "tok", "chat_id": "1", "base_url": srv.URL}, m); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
	if strings.Contains(text, "<script>") {
		t.Fatalf("HTML을 이스케이프하지 않아 인젝션 발생: %q", text)
	}
	if !strings.Contains(text, "&lt;script&gt;") {
		t.Fatalf("이스케이프한 엔티티가 있어야 합니다. 실제: %q", text)
	}
	if !strings.Contains(text, "a &amp; b") {
		t.Fatalf("&가 이스케이프되지 않았습니다. 실제: %q", text)
	}
}

func TestTelegramErrorClassification(t *testing.T) {
	rateLimited := capturePost(t, `{"ok":false,"error_code":429,"description":"Too Many Requests"}`, nil)
	_, err := (telegramChannel{}).Send(context.Background(),
		map[string]any{"bot_token": "tok", "chat_id": "1", "base_url": rateLimited.URL}, singleMsg())
	if err == nil || IsPermanent(err) {
		t.Fatalf("429는 재시도할 수 있어야 합니다. 실제: %v", err)
	}

	forbidden := capturePost(t, `{"ok":false,"error_code":403,"description":"bot was blocked by the user"}`, nil)
	_, err = (telegramChannel{}).Send(context.Background(),
		map[string]any{"bot_token": "tok", "chat_id": "1", "base_url": forbidden.URL}, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("403은 설정 문제이므로 영구 실패여야 합니다. 실제: %v", err)
	}
}

func TestWebhookDefaultTemplateProducesValidJSON(t *testing.T) {
	// 이것이 기본 템플릿이 존재하는 이유입니다. 제목에 따옴표와 줄바꿈이 있으면 단순한
	// `"title": "{{.Title}}"` 표현은 유효하지 않은 JSON을 만듭니다. {{json .}}는 그렇지 않습니다.
	srv := capturePost(t, `{"ok":true}`, func(t *testing.T, body map[string]any, _ *http.Request) {
		if body["title"] != `[🟠 높음] 登录处 "SQL注入" 风险` {
			t.Errorf("제목이 올바르게 복원되지 않았습니다: %v", body["title"])
		}
		items, _ := body["items"].([]any)
		if len(items) != 1 {
			t.Fatalf("items 수는 1이어야 합니다. 실제: %d", len(items))
		}
		it, _ := items[0].(map[string]any)
		if it["summary"] != "参数 id\n未过滤 导致注入" {
			t.Errorf("요약이 올바르게 복원되지 않았습니다: %v", it["summary"])
		}
		// 숫자는 문자열이 아니라 JSON 숫자여야 합니다(json:"...,string" 등의 표현을 사용하면 이 문제가 생김).
		if _, ok := it["finding_id"].(float64); !ok {
			t.Errorf("finding_id는 숫자여야 합니다. 실제: %T", it["finding_id"])
		}
	})
	if _, err := (webhookChannel{}).Send(context.Background(), map[string]any{"url": srv.URL}, singleMsg()); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

func TestWebhookCustomTemplateAndHeaders(t *testing.T) {
	srv := capturePost(t, `{"ok":true}`, func(t *testing.T, body map[string]any, r *http.Request) {
		if r.Header.Get("X-Token") != "s3cret" {
			t.Errorf("사용자 지정 헤더 누락: %v", r.Header)
		}
		if body["msg"] != "3 条" {
			t.Errorf("사용자 지정 템플릿 렌더링 오류: %v", body["msg"])
		}
		if body["first"] != "漏洞1" {
			t.Errorf("range 추출 오류: %v", body["first"])
		}
	})
	cfg := map[string]any{
		"url":           srv.URL,
		"headers":       map[string]any{"X-Token": "s3cret"},
		"body_template": `{"msg": {{json (printf "%d 条" .Count)}}, "first": {{json (index .Items 0).Name}}}`,
	}
	if _, err := (webhookChannel{}).Send(context.Background(), cfg, batchMsg(3)); err != nil {
		t.Fatalf("전송 실패: %v", err)
	}
}

func TestWebhookRejectsNonJSONRenderResult(t *testing.T) {
	cfg := map[string]any{"url": "https://example.com/hook", "body_template": `not json at all`}
	_, err := (webhookChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("JSON이 아닌 렌더링 결과는 영구 실패여야 합니다(잘못된 템플릿은 재시도가 무의미). 실제: %v", err)
	}
}

func TestWebhookValidateCatchesBadConfigEarly(t *testing.T) {
	bad := []map[string]any{
		{},
		{"url": "file:///etc/passwd"},
		{"url": "https://example.com", "method": "DELETE"},
		{"url": "https://example.com", "body_template": `{{.Items.`},
	}
	for i, cfg := range bad {
		if err := (webhookChannel{}).Validate(cfg); err == nil {
			t.Errorf("%d번째 설정 묶음은 거부되어야 합니다: %v", i, cfg)
		}
	}
}

func TestEmailMessageIsWellFormed(t *testing.T) {
	msg, err := buildEmailMessage("artex@example.com", []string{"a@example.com", "b@example.com"}, singleMsg())
	if err != nil {
		t.Fatalf("메일 구성 실패: %v", err)
	}
	if !strings.HasPrefix(msg, "From: artex@example.com\r\n") {
		t.Fatalf("From 헤더 오류:\n%s", msg)
	}
	if !strings.Contains(msg, "To: a@example.com, b@example.com\r\n") {
		t.Fatalf("To 헤더 오류:\n%s", msg)
	}
	// 중국어 제목은 RFC 2047로 인코딩해야 합니다. 그렇지 않으면 클라이언트에 깨진 문자로 표시됩니다.
	if !strings.Contains(msg, "Subject: =?utf-8?") {
		t.Fatalf("제목이 RFC 2047로 인코딩되지 않았습니다:\n%s", msg)
	}
	if dec, err := new(mime.WordDecoder).DecodeHeader(mustExtractHeader(t, msg, "Subject")); err != nil {
		t.Fatalf("제목을 디코딩할 수 없습니다: %v", err)
	} else if !strings.Contains(dec, "SQL注入") {
		t.Fatalf("제목 디코딩 후 내용 오류: %q", dec)
	}

	// 본문은 base64이며, 디코딩하면 유효한 HTML이어야 합니다.
	parts := strings.SplitN(msg, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatal("메일에 헤더/본문 구분이 없습니다")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(parts[1]), "\r\n", ""))
	if err != nil {
		t.Fatalf("본문 base64 디코딩 실패: %v", err)
	}
	html := string(decoded)
	if !strings.HasPrefix(html, "<div") {
		t.Fatalf("본문이 HTML이 아닙니다: %.80s", html)
	}
	// 제목은 텍스트 위치에 그대로 나타납니다. HTML 텍스트 내용의 큰따옴표는 유효한 문자이며 이스케이프가 필요하지 않습니다.
	// 여기서 "그대로 유지"를 검사하는 이유는 따옴표 이스케이프를 잘못 추가하여 중국어의 따옴표가
	// &quot;로 표시되는 일을 막기 위함입니다.
	if !strings.Contains(html, `"SQL注入"`) {
		t.Fatalf("제목의 따옴표는 텍스트 위치에서 그대로 유지해야 합니다: %.200s", html)
	}
}

// TestEmailEscapesStructuralInjection은 메일 본문에서 실제로 방지해야 하는 인젝션을 검사합니다.
// 취약점 제목과 요약은 테스트 대상 및 모델 출력에서 온 신뢰할 수 없는 내용입니다. 텍스트 위치에서는 반드시
// & < >를 이스케이프해야 하며(태그 인젝션 방지), 속성 위치에서는 따옴표도 이스케이프해야 합니다(href를 닫는 것을 방지).
func TestEmailEscapesStructuralInjection(t *testing.T) {
	m := Message{
		Items: []Item{{
			Severity:  "high",
			Name:      `<script>alert(1)</script>`,
			Summary:   "a & b > c",
			DetailURL: `https://artex.local/x?a="onmouseover=alert(1)`,
		}},
	}
	html := htmlBody(m, 0)
	if strings.Contains(html, "<script>") {
		t.Fatalf("제목이 이스케이프되지 않아 태그 인젝션 가능: %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatalf("이스케이프한 엔티티가 있어야 합니다: %s", html)
	}
	if !strings.Contains(html, "a &amp; b &gt; c") {
		t.Fatalf("&와 >가 이스케이프되지 않았습니다: %s", html)
	}
	// 상세 링크는 관리자가 설정하는 public_base_url로, 신뢰도가 높지만 속성 위치에서는 여전히
	// 따옴표를 이스케이프해야 합니다. 그렇지 않으면 따옴표가 있는 주소가 href를 닫고 이벤트 처리기를 인젝션할 수 있습니다.
	if strings.Contains(html, `onmouseover=alert(1)">`) {
		t.Fatalf("href 속성이 올바르게 이스케이프되지 않았습니다: %s", html)
	}
	if !strings.Contains(html, "&quot;") {
		t.Fatalf("속성 위치의 따옴표는 이스케이프해야 합니다: %s", html)
	}
}

func mustExtractHeader(t *testing.T, msg, name string) string {
	t.Helper()
	for _, line := range strings.Split(msg, "\r\n") {
		if strings.HasPrefix(line, name+": ") {
			return strings.TrimPrefix(line, name+": ")
		}
	}
	t.Fatalf("%s 헤더를 찾지 못했습니다", name)
	return ""
}

func TestChannelValidateReportsMissingFields(t *testing.T) {
	// 검증 오류는 설정하는 사람에게 직접 표시되므로, 막연한 "설정이 유효하지 않음" 대신 무엇이 없는지 명확히 설명해야 합니다.
	cases := []struct {
		kind   string
		cfg    map[string]any
		substr string
	}{
		{KindDingTalk, map[string]any{}, "Webhook"},
		{KindFeishu, map[string]any{}, "Webhook"},
		{KindWeCom, map[string]any{}, "Webhook"},
		{KindTelegram, map[string]any{}, "Bot Token"},
		{KindTelegram, map[string]any{"bot_token": "t"}, "Chat ID"},
		{KindEmail, map[string]any{}, "SMTP"},
		{KindEmail, map[string]any{"host": "h"}, "포트"},
		{KindEmail, map[string]any{"host": "h", "port": 587, "from": "f"}, "수신자"},
	}
	for _, tc := range cases {
		ch, ok := Get(tc.kind)
		if !ok {
			t.Fatalf("채널 %s: 등록되지 않았습니다", tc.kind)
		}
		err := ch.Validate(tc.cfg)
		if err == nil {
			t.Errorf("%s 설정 %v: 검증에 실패해야 합니다", tc.kind, tc.cfg)
			continue
		}
		if !strings.Contains(err.Error(), tc.substr) {
			t.Errorf("%s 오류 메시지에 %q 포함 필요. 실제: %q", tc.kind, tc.substr, err.Error())
		}
	}
}

// TestEmailSMTPErrorClassification은 SMTP 4xx/5xx의 의미 차이를 고정합니다.
// 4xx도 영구 실패로 판정하면 Greylisting을 사용하는 메일 서버에서는 모든 알림이 첫 번째
// 시도 후 failed로 넘어갑니다. Greylisting은 자동 재시도가 가장 효과를 발휘해야 하는 상황입니다.
func TestEmailSMTPErrorClassification(t *testing.T) {
	cases := []struct {
		reply     string
		permanent bool
	}{
		{"450 4.7.1 Greylisting in action, please come back later", false},
		{"451 4.3.0 Temporary system failure", false},
		{"452 4.2.2 Mailbox full", false},
		{"550 5.1.1 User unknown", true},
		{"553 5.1.3 Bad address syntax", true},
		{"554 5.7.1 Relay access denied", true},
		// 응답 코드를 가져올 수 없으면 "재시도 가능"으로 처리합니다. 한 번 더 시도하더라도 일시적일 수 있는
		// 오류를 영구 실패로 판정하지 않아야 합니다.
		{"unexpected EOF", false},
		{"", false},
	}
	for _, tc := range cases {
		err := smtpStageError("收件人被拒", errors.New(tc.reply))
		if got := IsPermanent(err); got != tc.permanent {
			t.Errorf("응답 %q: 기대 permanent=%v 실제 %v", tc.reply, tc.permanent, got)
		}
		// 분류 방식과 관계없이 원문을 유지하여 사용자가 원인을 파악할 수 있어야 합니다.
		if tc.reply != "" && !strings.Contains(err.Error(), tc.reply) {
			t.Errorf("응답 %q 원문 누락: %v", tc.reply, err)
		}
	}
}

func TestRegistryCoversAllKinds(t *testing.T) {
	// 여섯 채널 모두 있어야 합니다. 하나라도 없으면 UI 드롭다운에서 조용히 사라집니다.
	want := []string{KindDingTalk, KindEmail, KindFeishu, KindTelegram, KindWebhook, KindWeCom}
	got := Kinds()
	if len(got) != len(want) {
		t.Fatalf("채널 수는 %d여야 합니다. 실제: %d: %v", len(want), len(got), got)
	}
	for _, k := range want {
		if !ValidKind(k) {
			t.Errorf("채널 %s: 등록되지 않았습니다", k)
		}
		if ch, ok := Get(k); !ok || ch.Kind() != k {
			t.Errorf("채널 %s: Kind()와 등록 키가 일치하지 않습니다", k)
		}
	}
	if ValidKind("nope") {
		t.Error("등록하지 않은 유형은 검증에 통과해서는 안 됩니다")
	}
}

func TestPermanentErrorUnwrap(t *testing.T) {
	base := &permanentSentinel{}
	err := Permanent(base)
	if !IsPermanent(err) {
		t.Fatal("영구 실패로 식별해야 합니다")
	}
	if !strings.Contains(err.Error(), "sentinel") {
		t.Fatalf("오류 메시지는 원래 오류를 그대로 전달해야 합니다: %v", err)
	}
	if Permanent(nil) != nil {
		t.Fatal("Permanent(nil)은 반드시 nil을 반환해야 합니다")
	}
	if IsPermanent(nil) {
		t.Fatal("nil은 영구 실패가 아닙니다")
	}
}

type permanentSentinel struct{}

func (*permanentSentinel) Error() string { return "sentinel" }
