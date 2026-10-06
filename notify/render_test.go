package notify

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncateBytesKeepsValidUTF8(t *testing.T) {
	// 이 패키지의 가장 중요한 불변 조건입니다. WeCom(기업용 위챗)은 **바이트**로 길이를 제한하며, 중국어는 글자당 3바이트이므로,
	// 바이트 단위로 바로 자르는 구현은 한 글자를 반으로 잘라 유효하지 않은 UTF-8이 되어 플랫폼에서 거부됩니다.
	// 길이가 서로소인 다양한 중국어·영어 혼합 입력으로 가능한 모든 잘림 위치를 검사합니다.
	inputs := []string{
		"中文测试内容",
		"混合 mixed 内容 content",
		"a中b文c测d试e",
		"🔴🟠🟡🔵", // 4바이트 emoji이므로 잘못 자르면 더 뚜렷하게 드러남
		strings.Repeat("漏洞", 100),
	}
	for _, in := range inputs {
		for max := 1; max <= len(in)+2; max++ {
			got := TruncateBytes(in, max)
			if !utf8.ValidString(got) {
				t.Fatalf("입력 %q max=%d: 유효하지 않은 UTF-8 생성 %q", in, max, got)
			}
			if len(got) > max {
				t.Fatalf("입력 %q max=%d: 결과 %d바이트가 상한 초과", in, max, len(got))
			}
			// 자르지 않았으면 내용을 변경해서는 안 됩니다.
			if len(in) <= max && got != in {
				t.Fatalf("입력 %q max=%d: 상한을 넘지 않았는데 내용 변경 -> %q", in, max, got)
			}
		}
	}
}

func TestTruncateBytesZeroMeansUnlimited(t *testing.T) {
	long := strings.Repeat("x", 10000)
	if got := TruncateBytes(long, 0); got != long {
		t.Fatal("max=0은 제한 없음을 뜻해야 합니다")
	}
	if got := TruncateBytes(long, -5); got != long {
		t.Fatal("max<0은 제한 없음을 뜻해야 합니다")
	}
}

func TestTruncateBytesEllipsisBudget(t *testing.T) {
	// max가 말줄임표 자체보다 작을 때, 말줄임표를 덧붙여 상한을 초과해서는 안 됩니다.
	got := TruncateBytes("abcdefgh", 1)
	if len(got) > 1 {
		t.Fatalf("max=1일 때 결과 %q 길이 %d가 상한 초과", got, len(got))
	}
	// 정상적인 경우에는 말줄임표가 있어야 합니다.
	if got := TruncateBytes("abcdefgh", 5); !strings.HasSuffix(got, ellipsis) {
		t.Fatalf("말줄임표가 있어야 합니다. 실제: %q", got)
	}
}

func TestTruncateRunesCountsCharactersNotBytes(t *testing.T) {
	// TruncateBytes와의 기준 차이를 유지해야 합니다. Telegram은 문자 수로 길이를 제한하며,
	// 바이트 기준을 사용하면 중국어 메시지가 3분의 1만 남게 됩니다.
	s := "一二三四五六七八九十"
	got := TruncateRunes(s, 5)
	if n := utf8.RuneCountInString(got); n != 5 {
		t.Fatalf("기대 5문자, 실제 %d문자 (%q)", n, got)
	}
	// 같은 문자열을 바이트 기준으로 자르면 확실히 더 짧아야 합니다.
	if utf8.RuneCountInString(TruncateBytes(s, 5)) >= 5 {
		t.Fatal("바이트 기준 결과의 문자 수가 문자 기준과 같아서는 안 됩니다")
	}
}

func TestOneLineCollapsesWhitespace(t *testing.T) {
	got := OneLine("第一行\n\n第二行\t带制表   多空格", 0)
	if strings.ContainsAny(got, "\n\t") {
		t.Fatalf("모든 공백을 합쳐야 합니다. 실제: %q", got)
	}
	if strings.Contains(got, "  ") {
		t.Fatalf("연속 공백을 유지해서는 안 됩니다. 실제: %q", got)
	}
	// 자른 후에도 읽을 수 있고 유효해야 합니다.
	got = OneLine("一二三四五六七八九十", 4)
	if n := utf8.RuneCountInString(got); n != 4 {
		t.Fatalf("기대 4문자, 실제 %d (%q)", n, got)
	}
}

func TestTruncateHTMLNeverCutsTagInHalf(t *testing.T) {
	// HTML을 바로 자르면 `<a href="htt`와 같은 조각이 생겨 플랫폼이 메시지 전체를 거부합니다.
	s := `<b>标题</b>正文正文正文<a href="https://example.com/very/long/path">查看详情</a>`
	for max := 1; max <= utf8.RuneCountInString(s)+2; max++ {
		got := TruncateHTML(s, max)
		if n := utf8.RuneCountInString(got); max > 0 && n > max {
			t.Fatalf("max=%d: 결과 %d문자가 상한 초과", max, n)
		}
		// 끝부분에 닫히지 않은 `<`가 있어서는 안 됩니다(즉, 마지막 부분에 `<`가 있지만 `>`는 없음).
		if lt := strings.LastIndex(got, "<"); lt >= 0 && !strings.Contains(got[lt:], ">") {
			t.Fatalf("max=%d: 끝부분 태그가 잘렸습니다 -> %q", max, got)
		}
	}
}

func TestAssetLineOmitsExcess(t *testing.T) {
	if got := assetLine(nil, 3); got != "" {
		t.Fatalf("자산이 없으면 빈 문자열을 반환해야 합니다. 실제: %q", got)
	}
	if got := assetLine([]string{"a", "b"}, 3); got != "a, b" {
		t.Fatalf("상한 이내이면 모두 나열해야 합니다. 실제: %q", got)
	}
	// 상한을 넘으면 반드시 총수를 표시해야 합니다. 그렇지 않으면 독자는 나열하지 않은 자산이 몇 개 있는지 알 수 없습니다.
	got := assetLine([]string{"a", "b", "c", "d", "e"}, 2)
	if !strings.Contains(got, "등 총 5개") {
		t.Fatalf("총수 5를 표시해야 합니다. 실제: %q", got)
	}
}

func TestSeverityAndStatusLabels(t *testing.T) {
	if AtLeast("", "low") {
		t.Fatal("빈 심각도의 순위는 0이므로 모든 기준에서 차단해야 합니다")
	}
	if !AtLeast("critical", "") {
		t.Fatal("빈 기준은 허용해야 합니다")
	}
	if got := StatusLabel("fixed"); got != "수정 완료" {
		t.Fatalf("알 수 없는 상태 매핑. 실제: %q", got)
	}
	// 알 수 없는 상태는 그대로 반환하며, 라벨을 임의로 만들지 않습니다.
	if got := StatusLabel("weird_status"); got != "weird_status" {
		t.Fatalf("알 수 없는 상태는 그대로 반환해야 합니다. 실제: %q", got)
	}
}

// TestTruncateHTMLNeverCutsEntity는 감사에서 발견한 누락을 검사합니다. 자를 때는 불완전한
// 태그뿐 아니라 잘린 HTML 엔티티도 피해야 합니다.
//
// `&amp;`가 `&amp`로 잘리면 엔티티만 인식하는 파서는 **메시지 전체**를 거부할 수 있으며,
// 긴 모아 보내기 메시지는 흔하므로 그 손실은 너무 큽니다.
func TestTruncateHTMLNeverCutsEntity(t *testing.T) {
	s := "aaaa&amp;bbbb&lt;cccc&quot;dddd"
	for max := 1; max <= utf8.RuneCountInString(s)+2; max++ {
		got := TruncateHTML(s, max)
		// 끝부분에 "&는 있지만 대응하는 ;는 없는" 엔티티 조각이 있어서는 안 됩니다.
		if amp := strings.LastIndex(got, "&"); amp >= 0 && !strings.Contains(got[amp:], ";") {
			t.Fatalf("max=%d: 끝부분에 엔티티 조각 %q 남음", max, got[amp:])
		}
		if strings.Contains(got, "&amp\x00") {
			t.Fatalf("max=%d: 잘못된 엔티티 발생", max)
		}
	}
}
