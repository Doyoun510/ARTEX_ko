package notify

import "testing"

func TestParseFilterMalformedFallsBackToMatchAll(t *testing.T) {
	// 잘못된 JSON·빈 입력·타입이 다른 필드는 모두 영값 Filter가 되어야 하며,
	// 즉 "필터 없음"을 뜻합니다. 이 불변 조건은 "더 보내더라도 누락하지 않는다"는 원칙을 구현합니다.
	// 여기서 오류를 내거나 일부만 파싱하도록 바꾸면 사용자가 문자 하나를 잘못 입력해 높음 심각도 알림 전체를 조용히 잃게 됩니다.
	cases := []struct {
		name string
		raw  string
	}{
		{"빈 입력", ""},
		{"유효하지 않은 JSON", `{not json`},
		{"잘린 JSON", `{"min_severity":`},
		{"타입 불일치", `{"min_severity": 123, "task_ids": "abc"}`},
		{"최상위가 배열", `[1,2,3]`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := ParseFilter([]byte(tc.raw))
			if f.MinSeverity != "" || len(f.TaskIDs) != 0 || len(f.AssetIDs) != 0 {
				t.Fatalf("잘못된 설정은 영값 Filter가 되어야 합니다. 실제: %+v", f)
			}
			// 영값 Filter는 모든 이벤트에 매칭되어야 합니다.
			ev := Snapshot{Kind: EventFindingCreated, Severity: "low", VulnClass: "XSS"}
			if !Match(f, ev) {
				t.Fatal("영값 Filter는 모든 이벤트에 매칭되어야 합니다")
			}
		})
	}
}

func TestMatchSeverityThreshold(t *testing.T) {
	ev := func(sev string) Snapshot {
		return Snapshot{Kind: EventFindingCreated, Severity: sev}
	}
	cases := []struct {
		min    string
		sev    string
		expect bool
	}{
		{"", "low", true},
		{"", "critical", true},
		{"high", "critical", true},
		{"high", "high", true},
		{"high", "medium", false},
		{"high", "low", false},
		{"critical", "high", false},
		{"critical", "critical", true},
		// 알 수 없는 심각도의 순위는 0이므로 비어 있지 않은 모든 기준에서 차단해야 합니다(의심스러우면 보내지 않음).
		{"low", "", false},
		{"low", "unknown", false},
		{"", "", true},
	}
	for _, tc := range cases {
		got := Match(Filter{MinSeverity: tc.min}, ev(tc.sev))
		if got != tc.expect {
			t.Errorf("min=%q sev=%q: 기대 %v 실제 %v", tc.min, tc.sev, tc.expect, got)
		}
	}
}

func TestMatchScopeRestrictions(t *testing.T) {
	ev := Snapshot{
		Kind:      EventFindingCreated,
		Severity:  "high",
		TaskID:    7,
		AssetIDs:  []int64{10, 20},
		VulnClass: "SQL注入",
	}
	cases := []struct {
		name   string
		filter Filter
		expect bool
	}{
		{"빈 범위=제한 없음", Filter{}, true},
		{"작업 매칭", Filter{TaskIDs: []int64{7}}, true},
		{"작업 미매칭", Filter{TaskIDs: []int64{8}}, false},
		{"여러 작업 선택에 매칭 포함", Filter{TaskIDs: []int64{8, 7}}, true},
		{"자산 교집합 있음", Filter{AssetIDs: []int64{20, 99}}, true},
		{"자산 교집합 없음", Filter{AssetIDs: []int64{99}}, false},
		{"작업과 자산 모두 매칭", Filter{TaskIDs: []int64{7}, AssetIDs: []int64{10}}, true},
		{"작업은 매칭되지만 자산은 미매칭", Filter{TaskIDs: []int64{7}, AssetIDs: []int64{99}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Match(tc.filter, ev); got != tc.expect {
				t.Errorf("기대 %v 실제 %v", tc.expect, got)
			}
		})
	}
}

func TestMatchVulnClassKeywords(t *testing.T) {
	ev := func(class string) Snapshot {
		return Snapshot{Kind: EventFindingCreated, Severity: "high", VulnClass: class}
	}
	cases := []struct {
		name   string
		filter Filter
		class  string
		expect bool
	}{
		{"include 비어 있음=모두 수신", Filter{}, "任意类型", true},
		{"include 매칭", Filter{VulnClassInclude: []string{"SQL"}}, "SQL注入", true},
		{"include 미매칭", Filter{VulnClassInclude: []string{"命令执行"}}, "SQL注入", false},
		{"include 여러 단어 중 하나라도 매칭", Filter{VulnClassInclude: []string{"命令执行", "SQL"}}, "SQL注入", true},
		{"대소문자 구분 안 함", Filter{VulnClassInclude: []string{"sql"}}, "SQL注入", true},
		{"exclude 매칭 시 제외", Filter{VulnClassExclude: []string{"信息泄露"}}, "信息泄露", false},
		{"exclude 미매칭 시 허용", Filter{VulnClassExclude: []string{"信息泄露"}}, "SQL注入", true},
		// 제외가 포함보다 우선하며, 동시에 매칭되면 제외해야 합니다.
		{"제외가 포함보다 우선", Filter{
			VulnClassInclude: []string{"SQL"},
			VulnClassExclude: []string{"注入"},
		}, "SQL注入", false},
		// 공백뿐인 키워드는 무시해야 합니다. 그렇지 않으면 "공백을 포함한 모든 문자열에 매칭"하게 됩니다.
		{"공백 키워드 무시", Filter{VulnClassInclude: []string{"", "  "}}, "SQL注入", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Match(tc.filter, ev(tc.class)); got != tc.expect {
				t.Errorf("기대 %v 실제 %v", tc.expect, got)
			}
		})
	}
}

func TestMatchStatusChangeRequiresOptIn(t *testing.T) {
	ev := Snapshot{Kind: EventFindingStatusChanged, Severity: "critical", FromStatus: "pending", ToStatus: "fixed"}
	// 기본적으로 꺼져 있습니다. 대부분의 사람에게 "취약점 알림 전송"은 상태 기록이 아니라 새 취약점 발견을 뜻합니다.
	if Match(Filter{MinSeverity: "low"}, ev) {
		t.Fatal("활성화하지 않은 상태 변경 이벤트는 건너뛰어야 합니다")
	}
	if !Match(Filter{OnStatusChange: true}, ev) {
		t.Fatal("on_status_change 활성화 후 상태 변경 이벤트에 매칭되어야 합니다")
	}
	// 생성 이벤트는 on_status_change의 영향을 받지 않습니다.
	created := Snapshot{Kind: EventFindingCreated, Severity: "critical"}
	if !Match(Filter{MinSeverity: "low"}, created) {
		t.Fatal("생성 이벤트는 on_status_change에 의존해서는 안 됩니다")
	}
}
