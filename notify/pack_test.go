package notify

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// 이 파일은 "항목 전체 단위로 묶기" 수정을 검사합니다. 모아 보내기 메시지가 채널 길이 상한을 넘으면 반드시 **항목 전체 단위로**
// 자르고, 담지 못한 항목 수를 정확하게 보고하여 호출자가 실제로 전달한 항목만 표시하도록 해야 합니다.
//
// 이전 방식은 전체를 렌더링한 뒤 자르고 배치 전체를 전달됨으로 표시했습니다. 메시지 뒤쪽이 사라졌지만,
// 전송 이력에는 전부 성공으로 나타나 취약점이 사라졌으며, 어디에서도 발견할 수 없었습니다.

func TestMarkdownBodyPacksWholeItemsWithinByteLimit(t *testing.T) {
	// 중국어 모아 보내기 200건은 WeCom(기업용 위챗)의 4096바이트를 확실히 넘습니다.
	m := batchMsg(200)
	body, kept := markdownBody(m, weComMarkdownLimit)

	if len(body) > weComMarkdownLimit {
		t.Fatalf("본문 %d바이트가 상한 %d 초과", len(body), weComMarkdownLimit)
	}
	if !utf8.ValidString(body) {
		t.Fatal("본문이 유효한 UTF-8이 아닙니다")
	}
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("일부만 담아야 합니다(0 < kept < %d). 실제: %d", len(m.Items), kept)
	}
	// 머리에는 이 메시지에 몇 건만 포함되고 나머지는 몇 건인지 반드시 정확하게 설명해야 하며, 그렇지 않으면 독자는 머리의
	// 숫자를 전체로 생각하게 됩니다.
	if !strings.Contains(body, "나머지") || !strings.Contains(body, "다음 메시지에서 이어집니다") {
		t.Fatalf("머리에 이 메시지에 포함되지 않은 건수를 설명해야 합니다:\n%s", body[:minInt(400, len(body))])
	}
	// 처음 kept건만 포함해야 합니다.
	for i := 0; i < kept; i++ {
		if !strings.Contains(body, "漏洞"+itoa(i+1)) {
			t.Fatalf("%d번째 항목은 이 메시지에 있어야 합니다:\n%s", i+1, body)
		}
	}
	if strings.Contains(body, "漏洞"+itoa(kept+1)) {
		t.Fatalf("%d번째 항목은 나타나서는 안 됩니다(다음 배치에 속함)", kept+1)
	}
}

func TestMarkdownBodyKeepsEverythingWhenUnderLimit(t *testing.T) {
	m := batchMsg(3)
	body, kept := markdownBody(m, 0) // 0 = 제한 없음
	if kept != len(m.Items) {
		t.Fatalf("길이 제한이 없으면 전부 유지해야 합니다. 실제: kept=%d", kept)
	}
	if strings.Contains(body, "나머지") {
		t.Fatalf("자르지 않았으면 잘림 안내가 나타나서는 안 됩니다:\n%s", body)
	}
}

func TestMarkdownBodyAlwaysKeepsAtLeastOneItem(t *testing.T) {
	// 예산이 한 건도 담을 수 없을 만큼 작아도 한 건은 전송해야 합니다(최종 잘림 처리로 보호).
	// 그렇지 않으면 긴 취약점 하나가 배치 전체를 영원히 막습니다. 매번 가져와도 담지 못하고 전송하지 않게 됩니다.
	m := batchMsg(5)
	_, kept := markdownBody(m, 50)
	if kept != 1 {
		t.Fatalf("1건 이상 유지해야 합니다. 실제: %d", kept)
	}
}

func TestMarkdownBodySingleReturnsOne(t *testing.T) {
	_, kept := markdownBody(singleMsg(), 4096)
	if kept != 1 {
		t.Fatalf("단일 메시지는 전달 1건을 보고해야 합니다. 실제: %d", kept)
	}
	// 빈 메시지에는 전달할 수 있는 항목이 없습니다.
	if _, k := markdownBody(Message{}, 4096); k != 0 {
		t.Fatalf("빈 메시지는 0건을 보고해야 합니다. 실제: %d", k)
	}
}

func TestTelegramPackingUsesRuneBudget(t *testing.T) {
	m := batchMsg(200)
	text, kept := telegramHTML(m)
	// Telegram은 **문자 수**로 길이를 제한합니다. 바이트 기준을 사용하면 중국어 메시지가 3분의 1로 압축됩니다.
	if n := utf8.RuneCountInString(text); n > telegramTextLimit {
		t.Fatalf("본문 %d문자가 상한 %d 초과", n, telegramTextLimit)
	}
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("일부만 담아야 합니다. 실제: %d", kept)
	}
	if !strings.Contains(text, "다음 메시지에서 이어집니다") {
		t.Fatalf("아직 포함하지 않은 나머지가 있음을 설명해야 합니다:\n%.300s", text)
	}
}

func TestFeishuPackingReportsKept(t *testing.T) {
	m := batchMsg(2000)
	_, kept := feishuCard(m)
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("카드에는 일부만 담아야 합니다. 실제: %d", kept)
	}
}

func TestWebhookAndEmailReportAllItems(t *testing.T) {
	// 이 두 채널은 본문을 자르지 않아 배치 전체를 전달된 것으로 처리합니다.
	m := batchMsg(7)
	if n := len(m.Items); n != 7 {
		t.Fatal("사전 조건이 성립하지 않습니다")
	}
	// 렌더러 반환값으로 간접 확인합니다. markdownBody(0)으로 제한하지 않으면 전부 유지합니다.
	if _, k := markdownBody(m, 0); k != len(m.Items) {
		t.Fatalf("길이 제한이 없으면 전부 적용해야 합니다. 실제: %d", k)
	}
}

// TestMarkdownEscapesUntrustedContent는 "신뢰할 수 없는 내용이 메시지 구조를 바꾸어서는 안 됨"에 대한 회귀 테스트입니다.
// 제목과 요약은 모델 출력(모델은 테스트 대상 응답을 읽음)에서, 자산 이름은 테스트 대상 URL에서 옵니다.
func TestMarkdownEscapesUntrustedContent(t *testing.T) {
	cases := []struct {
		name  string
		item  Item
		must  []string // 결과에 반드시 나타나야 함(이스케이프된 형태)
		wrong []string // 결과에 나타나서는 안 됨(이스케이프하지 않은 형태)
	}{
		{
			name: "标题里的换行 + 外链",
			item: Item{
				Severity: "high",
				Name:     "登录口 SQL 注入\n[紧急：点此验证账号](http://attacker.tld)",
			},
			// 줄바꿈은 반드시 합쳐야 합니다(새 목록 항목/인용 블록을 만들 수 있기 때문).
			// 대괄호와 소괄호는 반드시 이스케이프해야 합니다(클릭 가능한 외부 링크가 되기 때문).
			must:  []string{`\[紧急：点此验证账号\]`, `\(http://attacker.tld\)`},
			wrong: []string{"\n[紧急", "\n\n[紧急"},
		},
		{
			name: "标题里的图片信标",
			item: Item{
				Severity: "high",
				Name:     "漏洞 ![](http://attacker.tld/beacon)",
			},
			must:  []string{`\!`, `\(http://attacker.tld/beacon\)`},
			wrong: []string{"![]("},
		},
		{
			name: "资产名里的强调与引用",
			item: Item{
				Severity: "high",
				Name:     "普通标题",
				Assets:   []string{"a.com/*注入*>引用"},
			},
			must:  []string{`\*注入\*`, `\>`},
			wrong: []string{"*注入*"},
		},
		{
			name: "摘要里的反引号与竖线",
			item: Item{
				Severity: "high",
				Name:     "标题",
				Summary:  "`code` | 表格",
			},
			must:  []string{"\\`code\\`", `\|`},
			wrong: []string{"`code`"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := Message{Items: []Item{tc.item}}
			// 단일 모드의 Item 작성은 세 markdown 채널이 공유하는 렌더링 경로입니다.
			var b strings.Builder
			writeItem(&b, tc.item, "", true)
			got := b.String()
			for _, want := range tc.must {
				if !strings.Contains(got, want) {
					t.Errorf("이스케이프된 형태 %q 누락:\n%s", want, got)
				}
			}
			for _, bad := range tc.wrong {
				if strings.Contains(got, bad) {
					t.Errorf("이스케이프하지 않은 형태 %q 발생(구조 또는 외부 링크를 인젝션하는 데 사용 가능):\n%s", bad, got)
				}
			}
			_ = m
		})
	}
}

// TestMarkdownEscapeBackslashFirst는 이스케이프 순서를 고정합니다. 역슬래시를 반드시 먼저 처리해야 하며,
// 그렇지 않으면 뒤에 덧붙인 역슬래시를 다시 이스케이프하여 출력에 이중 역슬래시가 나타납니다.
func TestMarkdownEscapeBackslashFirst(t *testing.T) {
	if got := markdownEscape(`a\b*c`); got != `a\\b\*c` {
		t.Fatalf("이스케이프 순서 오류. 실제: %q", got)
	}
}

// TestTelegramTitleHasNoMarkdownEscapes는 특정 회귀를 검사합니다.
// markdown 이스케이프가 Telegram의 HTML 출력에 나타나서는 안 됩니다(공유 제목 함수에
// 이스케이프를 추가한 적이 있으며, 그 결과 Telegram 메시지에 `\(1\)`처럼 눈에 보이는 역슬래시가 나타났음).
func TestTelegramTitleHasNoMarkdownEscapes(t *testing.T) {
	m := Message{Items: []Item{{Severity: "high", Name: "alert(1) *重点*"}}}
	text, _ := telegramHTML(m)
	if strings.Contains(text, `\(`) || strings.Contains(text, `\*`) {
		t.Fatalf("Telegram 본문에 markdown 역슬래시 이스케이프 발생:\n%s", text)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
