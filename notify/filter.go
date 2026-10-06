package notify

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// Filter는 notification_channels.filter JSONB 열의 계약으로, 채널 인스턴스의 필터 조건입니다.
// 모든 필드는 선택 사항이며 생략하면 '필터링하지 않음'입니다. 이는 잘못된 설정의 기본 처리 의미이므로 ParseFilter를 참조하세요.
type Filter struct {
	// MinSeverity는 최소 심각도 기준(low/medium/high/critical)이며, 비어 있으면 기준을 두지 않습니다.
	MinSeverity string `json:"min_severity"`
	// TaskIDs / AssetIDs가 빈 배열이면 제한하지 않으며, 비어 있지 않으면 이벤트와 공통 항목이 있어야 합니다.
	TaskIDs  []int64 `json:"task_ids"`
	AssetIDs []int64 `json:"asset_ids"`
	// VulnClassInclude가 비어 있으면 모두 받으며, 비어 있지 않으면 vulnclass가 키워드 중 하나 이상과 매칭되어야 합니다.
	// VulnClassExclude의 키워드 중 하나라도 매칭되면 제외합니다(제외가 포함보다 우선).
	// 매칭은 대소문자를 구분하지 않는 부분 문자열 매칭입니다. 정규 표현식보다 안전하며, 잘못된 정규 표현식 설정으로 채널이 조용히 동작하지 않게 되는 일을 막습니다.
	VulnClassInclude []string `json:"vulnclass_include"`
	VulnClassExclude []string `json:"vulnclass_exclude"`
	// OnStatusChange는 채널이 취약점 상태 변경 이벤트를 받을지 결정합니다(realtime 모드에서만 의미 있음).
	OnStatusChange bool `json:"on_status_change"`
}

// ParseFilter는 채널 필터 설정을 파싱합니다.
//
// **error를 절대로 반환하지 않습니다.** 의도적인 설계 선택으로, 필터 조건 설정이 잘못되면 항상 영값
// Filter(= 필터 없음 = 전부 매칭)를 사용합니다. 취약점 알림 시스템에서는 **알림 하나를 더 보내는 것이
// 높음 심각도 알림 하나를 조용히 누락하는 것보다 훨씬 낫기** 때문입니다. 파싱 실패가 "알림 전송 없음"이 되면 설정이 된 것처럼 보이지만,
// 실제로는 아무것도 보내지 않는 채널이 됩니다. 이는 가장 나쁜 실패 방식입니다.
func ParseFilter(raw []byte) Filter {
	var f Filter
	if len(raw) == 0 {
		return f
	}
	// 파싱에 실패하면 f는 영값을 유지하며, 필터링하지 않습니다.
	_ = json.Unmarshal(raw, &f)
	return f
}

// ValidMinSeverity는 s가 유효한 심각도 기준인지 반환합니다(빈 문자열이면 기준 없음).
func ValidMinSeverity(s string) bool {
	if s == "" {
		return true
	}
	_, ok := severityRank[s]
	return ok
}

// Validate는 필터 설정에서 **값이 제한된** 필드를 검증하며, 채널 저장 시 호출합니다.
//
// 쓰기 시 반드시 차단해야 하는 이유: Match는 알 수 없는 기준을 `rank >= 0`으로 판정하며 항상 참입니다.
// 따라서 min_severity의 철자 하나를 틀리면("hgih") 필터가 **오류를 알리지 않고 작동하지 않게** 되어
// '전부 전송'합니다. 이는 이 패키지의 '빠뜨리느니 더 많이 전송'한다는 선택 방향과 일치하지만(누락 없음),
// 사용자는 심각도별로 전송한다고 생각하는 동안 실제로는 모든 취약점이 그룹에 쏟아지며,
// 설정 오류를 알리는 징후도 없습니다. 이러한 '오류를 알리지 않고 필터링이 해제되는 동작'은 진입점에서 차단해야 합니다.
//
// Validate는 **쓰기** 경로에서만 사용합니다. 읽기 경로는 여전히 ParseFilter의 관대한 처리를 사용하여,
// 과거 데이터에 이미 존재하는 잘못된 값으로 인해 채널 전체를 읽지 못하는 일이 없도록 합니다.
func (f Filter) Validate() error {
	if !ValidMinSeverity(f.MinSeverity) {
		return fmt.Errorf("최소 심각도 %q: 유효하지 않습니다. 선택 가능: low / medium / high / critical. 비워 두면 제한하지 않습니다", f.MinSeverity)
	}
	return nil
}

// Match는 이벤트를 해당 필터 조건의 채널로 전송해야 하는지 판정합니다.
//
// **error를 절대로 반환하지 않으며**, 이유는 ParseFilter와 같습니다. 모든 내부 오류는 "매칭"으로 처리합니다.
// 판정 순서: 이벤트 유형 → 심각도 기준 → 작업/자산 범위 → 취약점 유형 키워드.
func Match(f Filter, s Snapshot) bool {
	// 상태 변경 이벤트는 명시적으로 활성화한 채널만 받습니다. 대부분의 사용자는
	// "알림 전송"을 모든 상태 전이를 따라가는 기록이 아니라 "새 취약점 발견"으로 기대하므로 기본적으로 꺼져 있습니다.
	if s.Kind == EventFindingStatusChanged && !f.OnStatusChange {
		return false
	}
	if !AtLeast(s.Severity, f.MinSeverity) {
		return false
	}
	if len(f.TaskIDs) > 0 && !slices.Contains(f.TaskIDs, s.TaskID) {
		return false
	}
	if len(f.AssetIDs) > 0 && !intersectsInt(f.AssetIDs, s.AssetIDs) {
		return false
	}
	// 제외 우선: 제외 키워드 중 하나라도 매칭되면 포함 목록에도 동시에 매칭되더라도 제외합니다.
	if len(f.VulnClassExclude) > 0 && containsAnyFold(s.VulnClass, f.VulnClassExclude) {
		return false
	}
	if len(f.VulnClassInclude) > 0 && !containsAnyFold(s.VulnClass, f.VulnClassInclude) {
		return false
	}
	return true
}

func intersectsInt(a, b []int64) bool {
	// 작은 집합은 선형 탐색이면 충분합니다. 양쪽 모두 "사람이 직접 선택한 수십 개" 수준이므로,
	// map 생성 비용이 이점보다 큽니다.
	for _, v := range b {
		if slices.Contains(a, v) {
			return true
		}
	}
	return false
}

// containsAnyFold는 s에 keywords의 키워드 중 하나라도 포함되는지 반환합니다(대소문자 구분 없음).
func containsAnyFold(s string, keywords []string) bool {
	lower := strings.ToLower(s)
	for _, kw := range keywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw != "" && strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
