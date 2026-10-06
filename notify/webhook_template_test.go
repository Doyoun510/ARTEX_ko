package notify

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// 이 파일은 일반 Webhook 템플릿의 **기능 경계**를 고정합니다.
//
// 이 패키지에서 "사용자가 제공한 문자열을 코드로 평가"하는 유일한 곳이므로, 할 수 있는 일과
// 할 수 없는 일을 명확히 하고 테스트로 해당 특성을 고정해야 합니다. 그렇지 않으면 템플릿 컨텍스트에
// 메서드나 FuncMap의 readFile을 추가하면 기능 범위가 조용히 확장되지만, diff에는
// 무해한 작은 함수처럼 보일 수 있습니다.

// TestTemplateContextHasNoMethods는 가장 중요한 검사입니다.
//
// text/template은 노출된 메서드를 호출합니다({{.Foo}}로 필드 접근과 메서드 호출 모두 가능). 따라서 템플릿 컨텍스트에서
// 노출된 메서드가 있는 **어떤** 타입이라도 접근할 수 있으면, 해당 메서드를 템플릿 작성자에게 노출하는 셈입니다.
// 이 기능의 컨텍스트는 의도적으로 순수 데이터뿐입니다(노출된 필드만 있고 메서드는 없음).
//
// 이 검사가 실패하면 누군가 webhookTemplateData / webhookItem에 메서드를 추가한 것입니다.
// 허용하기 전에, 템플릿이 해당 메서드로 노출해서는 안 되는 내용을 읽을 수 있는지 확인하세요.
func TestTemplateContextHasNoMethods(t *testing.T) {
	for _, v := range []any{webhookTemplateData{}, webhookItem{}} {
		typ := reflect.TypeOf(v)
		if n := typ.NumMethod(); n != 0 {
			var names []string
			for i := 0; i < n; i++ {
				names = append(names, typ.Method(i).Name)
			}
			t.Fatalf("%s: 메서드 %d개 노출(%s). text/template이 호출할 수 있어,"+
				"해당 메서드의 기능을 템플릿 작성자에게 제공하는 셈입니다", typ.Name(), n, strings.Join(names, ", "))
		}
	}
}

// TestTemplateFuncsAreMinimal은 템플릿에 노출한 함수 집합을 고정합니다.
//
// FuncMap에 함수를 하나 더 넣으면 기능이 하나 더 생깁니다. 현재는 json / jsons뿐이며, 값을 직렬화하여
// JSON 조각으로 만드는 역할입니다. 파일 읽기·요청 전송·명령 실행은 할 수 없습니다.
func TestTemplateFuncsAreMinimal(t *testing.T) {
	var got []string
	for name := range webhookTemplateFuncs {
		got = append(got, name)
	}
	sort.Strings(got)
	want := []string{"json", "jsons"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("템플릿 함수 집합 변경: 실제 %v, 기대 %v. 함수 추가 전에 기능 범위가 확장되지 않는지 확인하세요"+
			"(파일 읽기·쓰기, 네트워크 요청, 명령 실행 금지)", got, want)
	}
}

// TestTemplateCannotReachUnknownData는 템플릿에서의 범위 밖 접근을 검사합니다.
// 존재하지 않는 내용에 접근하면 무언가 반환하는 대신 반드시 실패해야 하며, 실패 메시지에 내부 데이터가 노출되어서는 안 됩니다.
func TestTemplateCannotReachUnknownData(t *testing.T) {
	_, err := renderWebhookBody(`{"x": {{.Environment}}, "y": {{.Env}}}`, singleMsg())
	if err == nil {
		t.Fatal("없는 필드에 접근하면 오류가 나야 합니다")
	}
	// 오류에 템플릿 컨텍스트의 실제 내용(취약점 제목/요약)이 나타나서는 안 됩니다.
	for _, leak := range []string{"SQL注入", "参数 id"} {
		if strings.Contains(err.Error(), leak) {
			t.Errorf("템플릿 오류에 메시지 내용 %q 노출: %v", leak, err)
		}
	}
}

// TestTemplateRenderFailsPermanently는 잘못된 템플릿이 설정 오류이며 재시도로 해결되지 않음을 검사합니다.
// 재시도 가능으로 판정하면 잘못된 템플릿 하나로 매번 전송 시 불필요하게 세 번 백오프하게 됩니다.
func TestTemplateRenderFailsPermanently(t *testing.T) {
	cfg := map[string]any{
		"url":           "https://example.com/hook",
		"body_template": `{{.Items.`,
	}
	if err := (webhookChannel{}).Validate(cfg); err == nil {
		t.Fatal("템플릿 구문 오류는 저장 시 차단해야 합니다")
	}
	// 검증을 우회하여 바로 전송해도 반복 재시도 대신 반드시 영구 실패로 판정해야 합니다.
	_, err := (webhookChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("잘못된 템플릿은 영구 실패로 판정해야 합니다. 실제: %v", err)
	}
}

// TestTemplateCanOnlyProduceJSON은 "템플릿 렌더링 결과가 반드시 유효한 JSON이어야 함"을 검사합니다.
// 이는 "템플릿으로 일반 텍스트를 생성하여 다른 프로토콜을 실행"하는 사용도 차단합니다.
func TestTemplateCanOnlyProduceJSON(t *testing.T) {
	// 유효한 템플릿은 통과할 수 있습니다.
	ok := map[string]any{"url": "https://example.com/hook", "body_template": `{"t":{{json .Title}}}`}
	if err := (webhookChannel{}).Validate(ok); err != nil {
		t.Fatalf("유효한 템플릿은 검증에 통과해야 합니다: %v", err)
	}
	// JSON이 아닌 렌더링 결과는 반드시 거부해야 합니다(그대로 전송하지 않음).
	bad := map[string]any{"url": "http://127.0.0.1:1/hook", "body_template": `not json {{.Count}}`}
	_, err := (webhookChannel{}).Send(context.Background(), bad, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("JSON이 아닌 렌더링 결과는 영구 실패로 판정해야 합니다. 실제: %v", err)
	}
	if !strings.Contains(err.Error(), "유효한 JSON") {
		t.Errorf("오류 메시지에 JSON 문제임을 설명해야 합니다. 실제: %v", err)
	}
}
