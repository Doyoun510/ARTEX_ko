package notify

import (
	"strings"
	"unicode/utf8"
)

const ellipsis = "…"

// TruncateBytes는 s를 max바이트 이하로 자르며, 결과가 유효한 UTF-8이고 문자가 잘리지 않도록 보장합니다.
//
// 문자 경계에서 잘라야 하는 이유: WeCom(기업용 위챗) 그룹 봇의 markdown에는 4096 **바이트**의 엄격한 상한이 있으며(문자 수가
// 아님), 중국어 한 글자는 3바이트입니다. 바이트 단위로 바로 자르면 한 글자를 반으로 잘라 유효하지 않은
// UTF-8이 됩니다. 플랫폼에서 메시지 전체를 거부하거나 깨진 네모로 표시합니다. 여기서는 예산 위치부터
// 가장 가까운 rune 시작 바이트까지 뒤로 이동합니다(utf8.RuneStart로 연속 바이트 0b10xxxxxx 판정).
//
// max<=0이면 제한하지 않습니다. max가 말줄임표조차 담지 못할 만큼 작지 않으면 자른 뒤 말줄임표를 덧붙입니다.
func TruncateBytes(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	budget := max - len(ellipsis)
	suffix := ellipsis
	if budget < 0 {
		// max가 말줄임표보다 짧으면 말줄임표 없이 자르기만 하여, 결과가 오히려 max를 넘는 것을 막습니다.
		budget = max
		suffix = ""
	}
	cut := budget
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + suffix
}

// OneLine은 여러 줄 텍스트를 한 줄로 압축합니다. 모든 공백을 합친 뒤 문자 수에 따라 자릅니다.
// IM 메시지 제목 줄에 사용합니다. 요약에는 줄바꿈이 많아 표/제목에 바로 넣으면 레이아웃이 깨질 수 있습니다.
// max<=0이면 길이를 제한하지 않습니다.
func OneLine(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	return TruncateRunes(s, max)
}

// TruncateRunes는 s를 max문자 이하로 자릅니다(바이트가 아님). 초과하면 말줄임표를 덧붙입니다.
// max<=0이면 제한하지 않습니다.
//
// TruncateBytes와의 차이는 플랫폼 기준입니다. WeCom은 바이트 수로, Telegram은 문자 수로 길이를 제한합니다.
// 잘못된 기준을 사용해도 오류가 나지 않고 예상보다 메시지가 훨씬 짧아집니다(중국어 1글자 = 3바이트이며,
// 4096바이트로 자르면 약 1365글자만 남음). 따라서 두 함수를 모두 유지하며 채널별로 선택해야 합니다.
func TruncateRunes(s string, max int) string {
	if max <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 1 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + ellipsis
}

// TruncateHTML은 HTML 조각을 문자 수에 따라 자르며, 태그가 중간에서 잘리지 않도록 보장합니다.
//
// HTML을 바로 문자 단위로 자르면 `<a href="htt`와 같은 불완전한 태그가 생겨, 플랫폼 파서가
// 오류로 메시지 전체를 거부하거나 뒤의 본문을 속성값으로 해석할 수 있습니다. 여기서는 먼저 문자 수에 따라 자른 뒤,
// 끝부분에 닫히지 않은 `<`가 있는지 확인하여 있으면 그 앞까지 되돌립니다.
//
// 태그 짝 맞추기(</b> 등을 보충)는 하지 않습니다. Telegram의 HTML 파서는 닫히지 않은 태그를 자동으로 닫으며,
// 직접 구현하면 속성의 따옴표·주석·스스로 닫는 태그까지 처리해야 하므로 복잡도에 비해 이점이 적습니다.
func TruncateHTML(s string, max int) string {
	if max <= 0 || len([]rune(s)) <= max {
		return s
	}
	cut := TruncateRunes(s, max)
	// 끝부분이 `<`로 시작하는 조각이면(마지막 `<` 뒤에 `>`가 없음), `<` 앞까지 되돌립니다.
	if lt := strings.LastIndex(cut, "<"); lt >= 0 && !strings.Contains(cut[lt:], ">") {
		cut = cut[:lt]
	}
	// 끝부분에 잘린 HTML 엔티티가 있으면(`&amp;`가 `&amp`로 잘린 경우 등), 마찬가지로 되돌립니다.
	// 엔티티 조각으로 인해 엔티티만 인식하는 파서가 **메시지 전체**를 거부할 수 있습니다. 길이
	// 상한을 넘는 모아 보내기 메시지는 흔하며, 그 때문에 알림 전체를 잃을 이유는 없습니다.
	if amp := strings.LastIndex(cut, "&"); amp >= 0 && !strings.Contains(cut[amp:], ";") {
		cut = cut[:amp]
	}
	return cut
}

// packItemCount는 예산 안에 **완전히** 담을 수 있는 건수를 계산하며, 모아 보내기 메시지를 항목 전체 단위로 묶는 데 사용합니다.
//
// 전체를 렌더링한 뒤 자르는 대신 항목 전체 단위로 묶는 이유: 자르면 뒤쪽 항목이 사라지지만,
// 전송 기록은 여전히 전달됨으로 표시됩니다. 메시지에도 없고 전송 이력에서도 확인할 수 없어,
// 취약점이 사라집니다. 항목 전체 단위로 묶으면 담지 못한 항목은 DB에 남아 다음 배치가 되며,
// 호출자가 받는 kept는 이번 메시지에 실제로 전달한 건수입니다.
//
// 파라미터: maxSize<=0이면 제한 없음. reserve는 메시지 머리/꼬리에 미리 할당한 양입니다.
// size는 크기를 측정합니다(플랫폼마다 기준이 다름: WeCom/DingTalk은 바이트, Telegram은 문자 수.
// 잘못된 기준을 사용해도 오류가 나지 않고 중국어 메시지가 상한보다 훨씬 작아질 뿐입니다).
// render는 idx번째 항목을 실제 텍스트로 렌더링합니다. 길이는 내용마다 달라 추정으로 처리할 수 없습니다.
//
// 항목이 남아 있으면 최소 1을 반환합니다. 단일 항목이 매우 길어도 전송하고 호출자가
// 최종 잘림 처리로 보호해야 합니다. 그렇지 않으면 긴 취약점 한 건이 배치 전체를 영구히 멈추게 합니다.
func packItemCount(items []Item, maxSize, reserve int, footer string, size func(string) int, render func(Item, int) string) int {
	if maxSize <= 0 {
		return len(items)
	}
	budget := maxSize - reserve - size(footer)
	if budget < 0 {
		budget = 0
	}
	used := 0
	for i, it := range items {
		used += size(render(it, i))
		if used > budget && i > 0 {
			return i
		}
	}
	return len(items)
}

// byteSize / runeSize는 packItemCount의 두 가지 크기 측정 기준입니다. 이름을 붙여 호출 위치에
// 이름 없는 func(s string) int 클로저가 나타나지 않도록 하며, 사용하는 기준을 쉽게 알아볼 수 있게 합니다.
func byteSize(s string) int { return len(s) }
func runeSize(s string) int { return utf8.RuneCountInString(s) }

// assetLine은 자산 목록을 한 줄의 표시 텍스트로 렌더링합니다. limit개를 넘으면 나머지를 생략하고 총수를 표시합니다.
// 취약점 하나가 자산 수십 개에 연결될 수 있어, 전부 나열하면 메시지가 너무 커집니다.
func assetLine(assets []string, limit int) string {
	if len(assets) == 0 {
		return ""
	}
	if limit <= 0 || len(assets) <= limit {
		return strings.Join(assets, ", ")
	}
	return strings.Join(assets[:limit], ", ") + " 등 총 " + itoa(len(assets)) + "개"
}

// itoa는 strconv.Itoa의 짧은 별칭이며, 표시 텍스트를 조합할 때만 사용하여 곳곳에서 import strconv를 하는 일을 피합니다.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
