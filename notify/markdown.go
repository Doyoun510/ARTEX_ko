package notify

import (
	"fmt"
	"strings"
)

// 이 파일은 "Markdown 계열" 채널(DingTalk, WeCom(기업용 위챗))이 공통으로 사용하는 메시지 렌더링입니다.
// Feishu는 카드 JSON, Telegram은 HTML, 메일은 HTML을 사용하며 각 어댑터에서 렌더링합니다.

// maxAssetsShown은 메시지에 나열할 최대 자산 수입니다. 취약점 하나가 자산 수십 개에 연결될 수 있지만,
// 전부 나열하면 메시지가 너무 커지고 정보 가치도 없습니다. IM에서 4번째 이후의 도메인을 읽는 사람은 없습니다.
const maxAssetsShown = 3

// maxSummaryRunes는 요약을 압축할 문자 수입니다. IM 메시지는 "상세 내용을 보도록 안내"하는 것으로,
// 보고서 자체가 아니며 전체 내용은 플랫폼에 있습니다.
const maxSummaryRunes = 120

// markdownReservedBytes는 메시지 머리(모아 보내기 줄 + 심각도 분포 + 발생할 수 있는 잘림 안내)에 미리 할당한 크기이며,
// 꼬리(플랫폼 링크)도 포함합니다. 항목 전체 단위로 묶을 때 이 부분을 예산에서 빼 머리와 꼬리가 잘리지 않게 합니다.
// 머리와 꼬리가 잘리면 독자는 "어떤 배치인지, 몇 건이 표시되지 않았는지"조차 알 수 없습니다.
const markdownReservedBytes = 320

// markdownEscape는 markdown 메타문자를 이스케이프합니다.
//
// 필수인 이유: 취약점 제목·요약·유형·자산 표시 이름은 모두 **신뢰할 수 없는 출처**에서 옵니다.
// 제목과 요약은 모델 출력(모델이 읽는 것은 테스트 대상의 응답)에서, 자산의 url은 스캔하여
// 얻은 전체 URL(대상이 제어할 수 있는 쿼리 문자열 포함)에서 옵니다. 이스케이프하지 않으면 제목이 아래와 같은
//
//	登录口 SQL 注入\n[紧急：点此验证账号](http://attacker.tld)
//
// 취약점은 보안 엔지니어의 DingTalk/Feishu에서 **클릭 가능한 외부 링크**로 렌더링됩니다.
// `![](http://attacker.tld/beacon)`은 렌더링 시 클라이언트가 가져오므로, "이 취약점을 확인했다"는 사실을
// 알리고 독자의 IP를 노출하는 셈입니다. 악의 없는 내용이라도 삽입된 굵은 글씨나
// 인용 블록이 아래의 심각한 취약점을 접힌 줄 뒤로 밀어낼 수 있습니다.
//
// 이스케이프 집합은 제목/링크/강조/목록/인용/취소선처럼 구조를 바꾸거나 클릭 가능한
// 요소를 만드는 문자를 포함합니다. `\`를 반드시 먼저 처리해야 뒤에 추가한 역슬래시가 다시 이스케이프되지 않습니다.
func markdownEscape(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		"`", "\\`",
		"*", `\*`,
		"_", `\_`,
		"[", `\[`,
		"]", `\]`,
		"(", `\(`,
		")", `\)`,
		"!", `\!`,
		"#", `\#`,
		">", `\>`,
		"|", `\|`,
		"~", `\~`,
	)
	return replacer.Replace(s)
}

// markdownText는 신뢰할 수 없는 텍스트를 한 줄로 압축하고 이스케이프하여 markdown 본문에 사용합니다.
// 한 줄로 만드는 것도 이스케이프만큼 중요합니다. 줄바꿈 자체로 새 목록 항목이나 인용 블록을 만들 수 있으며,
// 문자 이스케이프로는 막을 수 없습니다.
func markdownText(s string, maxRunes int) string {
	return markdownEscape(OneLine(s, maxRunes))
}

// markdownTitle은 메시지 제목(IM 플랫폼의 제목 표시줄/카드 제목)을 반환하며, 내용은 **이스케이프하지 않은 원문**입니다.
//
// 여기서는 의도적으로 이스케이프하지 않습니다. 제목을 네 가지 컨텍스트의 렌더러가 공유합니다. markdown 본문, Telegram의
// HTML, Feishu 카드의 plain_text, 일반 Webhook의 JSON 및 메일 제목입니다. 각 컨텍스트의
// 이스케이프 규칙은 모두 다르며(markdown 이스케이프를 HTML에 넣으면 눈에 보이는 역슬래시가 남고, JSON에 넣으면
// 데이터가 변형됨), 각 출력 측에서 이스케이프를 담당해야 합니다. writeItem / feishuItemLines /
// telegramEscape를 참조하세요. 공유 함수에 markdown 이스케이프를 추가한 적이 있으며, 그 결과 Telegram 메시지에
// `\(1\)`과 같이 눈에 보이는 역슬래시가 나타났습니다.
func markdownTitle(m Message) string {
	if m.Batch {
		return fmt.Sprintf("취약점 모아 보내기 · 총 %d건", len(m.Items))
	}
	if len(m.Items) == 0 {
		return "취약점 알림"
	}
	it := m.Items[0]
	return fmt.Sprintf("[%s] %s", SeverityLabel(it.Severity), OneLine(it.Title(), 0))
}

// markdownBody는 메시지 본문을 렌더링하며, 본문과 **실제로 작성한 항목 수**를 반환합니다.
//
// 반환값 kept는 이번 전송에서 실제로 전달한 항목 수입니다. 호출자는 이 값을 기준으로 처음 kept건만
// 전달됨으로 표시합니다. 채널 길이 상한으로 담지 못한 항목은 다음 배치로 남겨야 하며, 함께
// 성공으로 표시해서는 안 됩니다. 이것이 "조용한 손실"의 원인입니다. 메시지는 잘렸는데 전송 기록에는 전부 전달됨으로 표시되어,
// 뒤쪽 내용이 전송되지 않았다는 사실을 어디에서도 확인할 수 없습니다.
//
// maxBytes<=0이면 제한하지 않습니다.
func markdownBody(m Message, maxBytes int) (string, int) {
	if !m.Batch {
		if len(m.Items) == 0 {
			return "", 0
		}
		var b strings.Builder
		writeItem(&b, m.Items[0], "", true)
		// 단일 메시지가 너무 길어도 전송합니다(최종 잘림 처리로 보호). 취약점 한 건의 일부 정보라도
		// 전송하는 편이 한 건도 전송하지 않는 것보다 낫습니다.
		return TruncateBytes(b.String(), maxBytes), 1
	}

	footer := ""
	if m.HomeURL != "" {
		footer = fmt.Sprintf("\n[플랫폼에서 전체 보기](%s)\n", m.HomeURL)
	}
	kept := packItemCount(m.Items, maxBytes, markdownReservedBytes, footer, byteSize, func(it Item, idx int) string {
		var b strings.Builder
		writeItem(&b, it, fmt.Sprintf("%d. ", idx+1), false)
		return b.String()
	})

	items := m.Items[:kept]
	var b strings.Builder
	b.WriteString(markdownBatchIntro(m, items, len(m.Items)))
	for i, it := range items {
		writeItem(&b, it, fmt.Sprintf("%d. ", i+1), false)
	}
	b.WriteString(footer)
	return TruncateBytes(b.String(), maxBytes), kept
}

// markdownBatchIntro는 모아 보내기 메시지의 시작 부분을 렌더링합니다. 시간 구간·건수·심각도 분포를 표시합니다.
// 이 정보가 있으면 수신자는 플랫폼에 들어가지 않고도 배치를 즉시 처리해야 하는지 판단할 수 있습니다.
//
// items는 **실제로 담은** 항목이고, total은 이 배치의 전체 건수입니다. 둘이 다르면
// "다음 메시지에 몇 건이 남아 있는지" 반드시 명시해야 합니다. 그렇지 않으면 독자는 메시지 머리의 숫자가 전부라고 생각하며,
// 뒤의 전송되지 않은 항목은 화면에서도 전혀 보이지 않습니다.
func markdownBatchIntro(m Message, items []Item, total int) string {
	var b strings.Builder
	if m.WindowMinutes > 0 {
		fmt.Fprintf(&b, "**최근 %d분 새 취약점 %d개**", m.WindowMinutes, total)
	} else {
		fmt.Fprintf(&b, "**새 취약점 %d개**", total)
	}
	if extra := total - len(items); extra > 0 {
		fmt.Fprintf(&b, "(이 메시지에는 처음 %d건을 표시하며, 나머지 %d건은 다음 메시지에서 이어집니다)", len(items), extra)
	}
	// 심각도별 분포를 표시하여 심각한 항목이 있는지 한눈에 볼 수 있게 합니다. **이 메시지에 실제로 포함한**
	// 항목만 집계하여 "심각 3"과 아래에서 셀 수 있는 건수가 일치하도록 합니다.
	counts := map[string]int{}
	for _, it := range items {
		counts[it.Severity]++
	}
	var parts []string
	for _, sev := range []string{"critical", "high", "medium", "low"} {
		if n := counts[sev]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s %d", SeverityLabel(sev), n))
		}
	}
	if len(parts) > 0 {
		b.WriteString("\n" + strings.Join(parts, " · "))
	}
	b.WriteString("\n\n")
	return b.String()
}

// writeItem은 취약점 항목 하나를 렌더링합니다.
//
// prefix는 모아 보내기 목록의 순번에 사용합니다. single=true이면 전체 내용을 렌더링하며(요약과 상세 링크 포함),
// 모아 보내기 목록에서는 한 줄 요약만 렌더링합니다. 그렇지 않으면 50건의 모아 보내기가 긴 문서가 됩니다.
//
// 외부에서 온 모든 내용(제목/유형/자산/요약)은 markdownText를 거칩니다.
// 한 줄로 압축 + 이스케이프합니다. 상세 링크는 관리자가 설정한 public_base_url로 구성하며 신뢰할 수 없는 내용이 아니고,
// 클릭할 수 있어야 하므로 그대로 출력합니다.
func writeItem(b *strings.Builder, it Item, prefix string, single bool) {
	line := fmt.Sprintf("%s**%s · %s**", prefix, SeverityLabel(it.Severity), markdownText(it.Title(), 0))
	if !single {
		// 모아 보내기 모드: 한 줄로 표시하며, 자산과 요약을 압축하여 뒤에 붙입니다.
		var extras []string
		if a := assetLine(it.Assets, maxAssetsShown); a != "" {
			extras = append(extras, markdownText(a, 0))
		}
		if it.Summary != "" {
			extras = append(extras, markdownText(it.Summary, 60))
		}
		if len(extras) > 0 {
			line += " — " + strings.Join(extras, " · ")
		}
		b.WriteString(line + "\n")
		return
	}
	b.WriteString(line + "\n")
	if it.IsStatusChange() {
		fmt.Fprintf(b, "**상태 변경**: %s → %s\n",
			markdownText(StatusLabel(it.FromStatus), 0), markdownText(StatusLabel(it.ToStatus), 0))
	}
	if it.VulnClass != "" && it.VulnClass != it.Title() {
		fmt.Fprintf(b, "**유형**: %s\n", markdownText(it.VulnClass, 0))
	}
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		fmt.Fprintf(b, "**자산**: %s\n", markdownText(a, 0))
	}
	if it.Summary != "" {
		if s := markdownText(it.Summary, maxSummaryRunes); s != "" {
			fmt.Fprintf(b, "**요약**: %s\n", s)
		}
	}
	if it.DetailURL != "" {
		fmt.Fprintf(b, "[상세 보기](%s)\n", it.DetailURL)
	}
}
