package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// 발견 페이지의 '내보내기' 렌더링: findings 테이블 행 묶음을 종합 Markdown·단건 Markdown·
// CSV로 렌더링합니다. JSON은 server 계층에서 DTO로 직접 직렬화하며 여기서는 처리하지 않습니다.

// sortFindingsForExport 심각도 내림차순, 그다음 시간 내림차순으로 정렬하며 종합 보고서의 그룹과 일치합니다.
func sortFindingsForExport(fs []*db.DBFinding) {
	sort.SliceStable(fs, func(i, j int) bool {
		ri, rj := sevRank[fs[i].Severity], sevRank[fs[j].Severity]
		if ri != rj {
			return ri < rj // sevRank 값이 작을수록 더 심각합니다
		}
		return fs[i].CreatedAt.After(fs[j].CreatedAt)
	})
}

// findingTitle 취약점의 읽기 쉬운 제목을 선택합니다: 이름 → 분류 → '미분류'.
func findingTitle(f *db.DBFinding) string {
	return nz(f.Name, nz(f.VulnClass, "미분류"))
}

// FindingsMarkdown findings 묶음을 하나의 종합 보고서로 통합합니다(요약 + 심각도별 그룹,
// 각 항목에 분류/상태/소속 작업/증거/상세 보고서 포함).
func FindingsMarkdown(fs []*db.DBFinding, generatedAt time.Time) string {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var b strings.Builder
	b.WriteString("# 취약점 발견 종합 보고서\n\n")
	fmt.Fprintf(&b, "- **생성 시간**: %s\n", generatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **전체 발견 수**: %d개\n\n", len(items))

	// 요약: 심각도별 개수.
	counts := map[string]int{}
	for _, f := range items {
		counts[f.Severity]++
	}
	b.WriteString("## 요약\n\n")
	b.WriteString("| 심각도 | 개수 |\n| --- | --- |\n")
	for _, s := range []struct{ key, label string }{
		{"critical", "심각"}, {"high", "높음"}, {"medium", "중간"}, {"low", "낮음"},
	} {
		fmt.Fprintf(&b, "| %s | %d |\n", s.label, counts[s.key])
	}
	b.WriteString("\n")

	if len(items) == 0 {
		b.WriteString("_일치하는 취약점이 없습니다._\n")
		return b.String()
	}

	b.WriteString("## 취약점 상세\n\n")
	for i, f := range items {
		fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
		if f.VulnClass != "" {
			fmt.Fprintf(&b, "- **분류**: %s\n", f.VulnClass)
		}
		fmt.Fprintf(&b, "- **상태**: %s\n", nz(f.Status, "pending"))
		if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
			fmt.Fprintf(&b, "- **소속 작업**: %s\n", desc)
		}
		fmt.Fprintf(&b, "- **발견 시간**: %s\n\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		if s := strings.TrimSpace(f.Summary); s != "" {
			fmt.Fprintf(&b, "%s\n\n", s)
		}
		if e := strings.TrimSpace(f.Evidence); e != "" {
			fmt.Fprintf(&b, "**증거:**\n\n```\n%s\n```\n\n", e)
		}
		if rep := strings.TrimSpace(f.Report); rep != "" {
			b.WriteString("**상세 보고서:**\n\n")
			b.WriteString(rep)
			b.WriteString("\n\n")
		}
		b.WriteString(findingTrafficMarkdown(f, false))
		b.WriteString("---\n\n")
	}
	return b.String()
}

// SingleFindingMarkdown 단일 취약점을 독립 Markdown으로 렌더링합니다('취약점 하나당 파일 하나'로 압축할 때 사용).
func SingleFindingMarkdown(f *db.DBFinding, generatedAt time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# [%s] %s\n\n", strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
	if f.VulnClass != "" {
		fmt.Fprintf(&b, "- **분류**: %s\n", f.VulnClass)
	}
	fmt.Fprintf(&b, "- **심각도**: %s\n", nz(f.Severity, "info"))
	fmt.Fprintf(&b, "- **상태**: %s\n", nz(f.Status, "pending"))
	if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
		fmt.Fprintf(&b, "- **소속 작업**: %s\n", desc)
	}
	fmt.Fprintf(&b, "- **발견 시간**: %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **생성 시간**: %s\n\n", generatedAt.Format("2006-01-02 15:04:05"))
	if s := strings.TrimSpace(f.Summary); s != "" {
		fmt.Fprintf(&b, "## 개요\n\n%s\n\n", s)
	}
	if e := strings.TrimSpace(f.Evidence); e != "" {
		fmt.Fprintf(&b, "## 증거\n\n```\n%s\n```\n\n", e)
	}
	if rep := strings.TrimSpace(f.Report); rep != "" {
		b.WriteString("## 상세 보고서\n\n")
		b.WriteString(rep)
		b.WriteString("\n")
	}
	b.WriteString(findingTrafficMarkdown(f, true))
	return b.String()
}

var unsafeFilenameChars = regexp.MustCompile(`[^\p{Han}\p{L}\p{N}._-]+`)

// FindingFilename '취약점 하나당 파일 하나'에 사용할 안전한 .md 파일명을 생성합니다. 예:
// `critical_SQL注入_#123.md`. 경로 구분자와 제어 문자를 제거해 zip 내부의 잘못된 경로를 방지합니다.
func FindingFilename(f *db.DBFinding) string {
	sev := nz(f.Severity, "info")
	title := findingTitle(f)
	name := fmt.Sprintf("%s_%s_#%d", sev, title, f.ID)
	name = unsafeFilenameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = fmt.Sprintf("finding_%d", f.ID)
	}
	// 방어 처리: 경로를 한 번 더 제거해 zip slip을 방지합니다.
	name = path.Base(name)
	if len(name) > 120 {
		name = name[:120]
	}
	return name + ".md"
}

// FindingsCSV findings 묶음을 CSV로 렌더링합니다(UTF-8 BOM 포함, Excel에서 중국어를 올바르게 인식하도록 지원).
// 긴 report/evidence 전문은 포함하지 않고 요약 필드만 넣습니다. 전문이 필요하면 Markdown/JSON으로 내보냅니다.
func FindingsCSV(fs []*db.DBFinding) []byte {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "이름", "분류", "심각도", "상태", "소속 작업", "발견 시간", "개요", "트래픽 증거 수", "트래픽 증거 ID"})
	for _, f := range items {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			findingTitle(f),
			f.VulnClass,
			nz(f.Severity, "info"),
			nz(f.Status, "pending"),
			f.TaskDescription,
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			strings.TrimSpace(f.Summary),
			fmt.Sprint(len(f.TrafficBindings)), findingTrafficIDs(f),
		})
	}
	w.Flush()
	return buf.Bytes()
}

func findingTrafficIDs(f *db.DBFinding) string {
	ids := make([]string, 0, len(f.TrafficBindings))
	for _, b := range f.TrafficBindings {
		ids = append(ids, fmt.Sprint(b.ID))
	}
	return strings.Join(ids, ",")
}

func findingTrafficMarkdown(f *db.DBFinding, attachments bool) string {
	stale := f.Report != "" && f.EvidenceVersion != f.ReportEvidenceVersion
	if len(f.TrafficBindings) == 0 && !stale {
		return ""
	}
	var out strings.Builder
	out.WriteString("\n## 연관 트래픽 증거\n\n")
	fmt.Fprintf(&out, "증거 버전: %d, 바인딩 수: %d.\n\n", f.EvidenceVersion, len(f.TrafficBindings))
	if stale {
		out.WriteString("증거가 변경되어 상세 보고서 업데이트가 필요합니다.\n\n")
	}
	for i, b := range f.TrafficBindings {
		fmt.Fprintf(&out, "%d. **증거 #%d · %s** — `%s %s`, 상태 코드 %d\n", i+1, b.ID, b.Role, b.Snapshot.Method, strings.ReplaceAll(b.Snapshot.URL, "`", "%60"), b.Snapshot.Status)
		if b.Note != "" {
			fmt.Fprintf(&out, "   %s\n", strings.ReplaceAll(b.Note, "\n", "\n   "))
		}
		if attachments {
			fmt.Fprintf(&out, "   [요청 메시지](evidence/%d/%d/request.http) · [응답 메시지](evidence/%d/%d/response.http)\n", f.ID, b.ID, f.ID, b.ID)
		}
	}
	out.WriteString("\n")
	return out.String()
}
