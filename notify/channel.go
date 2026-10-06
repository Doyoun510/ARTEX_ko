package notify

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Channel은 알림 채널 어댑터입니다. 구현은 반드시 **상태를 가지지 않아야** 합니다. 동일한 인스턴스를 여러 채널
// 설정에서 동시에 재사용하며, 자격 증명은 모두 cfg 파라미터로 전달합니다.
type Channel interface {
	// Kind는 채널 유형 식별자를 반환하며, 등록 목록의 키와 일치해야 합니다.
	Kind() string
	// Validate는 설정 저장 시 호출하며 필수 필드와 형식을 검증합니다. 반환 오류는 설정하는 사람에게 직접
	// 표시되므로, 문구는 막연한 "설정이 유효하지 않음" 대신 "어떤 필드가 없는지" 설명해야 합니다.
	Validate(cfg map[string]any) error
	// Send는 메시지를 한 번 전송하며 **실제로 전달한 항목 수**와 오류를 반환합니다.
	//
	// 건수를 반환하는 이유: 각 플랫폼에는 메시지 길이 상한이 있어, 모아 보내기 메시지에 배치 전체를 담지 못하면 내용이 잘립니다.
	// 호출자가 무조건 배치 전체를 전달됨으로 표시하면 잘린 항목은 사라집니다. 메시지에도 없고,
	// 전송 이력에는 성공으로 표시되어 취약점이 전송되지 않았다는 사실을 어디에서도 확인할 수 없습니다. kept를 반환하면
	// 호출자는 처음 kept건만 표시하고, 나머지는 다음 배치로 남깁니다.
	//
	// 오류 반환은 전송 실패를 뜻하며, *PermanentError는 재시도해서는 안 됨을 나타냅니다.
	// 실패 시 kept는 의미가 없으므로 호출자는 무시해야 합니다.
	Send(ctx context.Context, cfg map[string]any, m Message) (int, error)
	// DefaultRatePerMin은 채널의 공식 권장 분당 전송 상한을 반환하며, 새 채널 인스턴스 생성
	// 시 전송 속도 제한 기본값으로 사용합니다. 0 반환은 알려진 제한이 없음을 뜻합니다.
	DefaultRatePerMin() int
	// SecretKeys는 채널 설정에서 자격 증명에 해당하는 키 이름을 반환합니다. API 응답 시 이 키의 값은 마스킹되며,
	// 업데이트 시 마스킹된 값을 받으면 DB의 기존 값을 유지합니다. 어떤 필드가 자격 증명인지는 구현 자체만 알 수
	// 있으므로(WeCom(기업용 위챗)은 Webhook 주소 전체가 자격 증명이고, DingTalk은 그중 secret만 해당),
	// 이 정보는 채널이 제공해야 하며 상위 계층에서 추측해서는 안 됩니다.
	SecretKeys() []string
	// DestinationKeys는 채널 설정에서 "메시지를 어디로 보낼지" 결정하는 키 이름을 반환합니다.
	//
	// SecretKeys와 마찬가지로 보안 관련 정보입니다. 대상 주소와 자격 증명은 서로 독립된 필드이며,
	// "주소만 변경하고 자격 증명은 그대로 유지"하도록 허용하면 채널 설정을 변경할 수 있는 누구나 DB의 실제 자격 증명을
	// 자신이 제어하는 서버로 보낼 수 있어, 채널 설정 마스킹의 의미가 완전히 사라집니다.
	// 자세한 내용은 PrepareConfigUpdate를 참조하세요.
	DestinationKeys() []string
}

// registry는 채널 등록 목록입니다. init() 자체 등록 대신 명시적 리터럴을 의도적으로 사용하여 "어떤 채널이 있는지"
// 한곳에서 전부 확인할 수 있고, 새 채널 추가 시 런타임 부수 효과 대신 컴파일 시점에 누락이 드러나게 합니다.
var registry = map[string]Channel{
	KindDingTalk: dingTalkChannel{},
	KindFeishu:   feishuChannel{},
	KindWeCom:    weComChannel{},
	KindWebhook:  webhookChannel{},
	KindTelegram: telegramChannel{},
	KindEmail:    emailChannel{},
}

// Get은 유형에 따라 채널 구현을 가져옵니다.
func Get(kind string) (Channel, bool) {
	c, ok := registry[kind]
	return c, ok
}

// ValidKind는 kind가 지원하는 채널 유형인지 반환합니다.
func ValidKind(kind string) bool {
	_, ok := registry[kind]
	return ok
}

// Kinds는 지원하는 모든 채널 유형을 사전순으로 반환합니다(UI 드롭다운에 일정하게 표시).
func Kinds() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// PermanentError는 재시도해서는 안 되는 전송 실패를 표시합니다. 자격 증명 오류, 대상의 거부, 유효하지 않은 요청 본문 등이 해당합니다.
// 재시도는 일시적 오류(네트워크 불안정, 전송 속도 제한, 상대 측 5xx)에만 의미가 있습니다. 영구 실패에 반복해서 백오프 후 재시도하면
// 성공할 수도 없고, 실제 오류가 재시도 로그에 묻히게 됩니다.
type PermanentError struct{ Err error }

func (e *PermanentError) Error() string { return e.Err.Error() }
func (e *PermanentError) Unwrap() error { return e.Err }

// Permanent는 err를 영구 실패로 표시합니다. err가 nil이면 nil을 반환하여,
// `return Permanent(someCheck())` 형태로 쉽게 작성할 수 있습니다.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &PermanentError{Err: err}
}

// IsPermanent는 err 체인에 영구 실패 표시가 있는지 반환합니다.
func IsPermanent(err error) bool {
	var pe *PermanentError
	return errors.As(err, &pe)
}

// ---- 설정 읽기 helper ----
//
// 채널 설정은 DB의 JSONB 열에서 오며, encoding/json으로 역직렬화하면 map[string]any이고,
// 숫자는 모두 float64, 배열은 []any입니다. 아래 helper는 이 변환을 통일하고 사용자가
// UI에서 비워 두어 생기는 타입 차이(포트를 문자열로 입력하는 경우 등)를 허용합니다.

// cfgString은 문자열 설정 항목을 가져오며, 웹 폼에서 복사·붙여넣기할 때 쉽게 포함되는 앞뒤 공백을 모두 제거합니다.
func cfgString(cfg map[string]any, key string) string {
	v, ok := cfg[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

// cfgInt는 정수 설정 항목을 가져오며, float64(JSON 기본값)와 문자열 두 가지 출처를 지원합니다.
func cfgInt(cfg map[string]any, key string) int {
	switch v := cfg[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

// cfgBool은 boolean 설정 항목을 가져오며, 문자열 "true"/"1"을 지원합니다.
func cfgBool(cfg map[string]any, key string) bool {
	switch v := cfg[key].(type) {
	case bool:
		return v
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		return s == "true" || s == "1" || s == "yes"
	default:
		return false
	}
}

// cfgStrings는 문자열 배열 설정 항목을 가져오며, 공백을 자동으로 제거하고 빈 문자열은 버립니다.
func cfgStrings(cfg map[string]any, key string) []string {
	raw, ok := cfg[key].([]any)
	if !ok {
		// 값이 하나인 폼 제출을 쉽게 하도록 단일 문자열도 허용합니다.
		if s := cfgString(cfg, key); s != "" {
			return []string{s}
		}
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// cfgMap은 문자열 매핑 설정 항목(사용자 지정 HTTP 헤더 등)을 가져오며, 키와 값의 공백을 제거하고 빈 키는 버립니다.
func cfgMap(cfg map[string]any, key string) map[string]string {
	raw, ok := cfg[key].(map[string]any)
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[k] = s
	}
	return out
}
