package notify

// Snapshot은 notification_events.snapshot JSONB 열의 계약입니다. 쓰기는 db 계층의
// 취약점 저장 트랜잭션에서, 읽기는 server 계층의 전송 엔진과 필터 매칭에서 수행합니다. 이 패키지에 정의하는 이유는
// "알림 영역"의 데이터이기 때문입니다. db는 직렬화만 담당하며 필드의 의미를 해석하지 않습니다.
//
// 렌더링 시 다시 조회하지 않고 취약점 필드를 중복 저장하는 이유: 취약점 이름·심각도·상태는 나중에 변경될 수 있지만,
// 알림 전송 내용은 **이벤트 발생 당시**의 결론을 반영해야 합니다. 다시 조회하면 "나중에 low로 변경된" 결과를 가져와
// 위험하게 오해할 수 있습니다. 또한 fan-out과 렌더링에서 findings/tasks/assets 세 테이블을 JOIN할 필요도 없어집니다.
type Snapshot struct {
	// 이벤트 유형: finding_created / finding_status_changed
	Kind      string  `json:"kind"`
	FindingID int64   `json:"finding_id"`
	TaskID    int64   `json:"task_id"`
	VulnClass string  `json:"vulnclass"`
	Name      string  `json:"name"`
	Severity  string  `json:"severity"`
	Summary   string  `json:"summary"`
	AssetIDs  []int64 `json:"asset_ids"`
	// kind=finding_status_changed일 때만 비어 있지 않습니다.
	FromStatus string `json:"from_status,omitempty"`
	ToStatus   string `json:"to_status,omitempty"`
}

// Item은 채널 렌더링을 위한 전송 대기 취약점 하나입니다.
type Item struct {
	FindingID int64
	Name      string
	VulnClass string
	Severity  string
	Summary   string
	// Assets는 확인된 자산 표시 이름(도메인/IP 등)입니다. server 계층에서 채웁니다.
	// 이 패키지는 DB에 접근하지 않아 이름을 가져올 수 없습니다.
	Assets []string
	// DetailURL은 취약점 상세 링크입니다. 비어 있으면 public_base_url이 설정되지 않은 것으로 렌더링 시 생략합니다.
	DetailURL string
	// 상태 변경 이벤트 전용입니다. 두 항목이 모두 비어 있지 않으면 "처리 대기 → 수정 완료"로 렌더링합니다.
	FromStatus string
	ToStatus   string
}

// IsStatusChange는 해당 항목이 상태 변경 이벤트인지 반환합니다.
func (i Item) IsStatusChange() bool { return i.FromStatus != "" || i.ToStatus != "" }

// Title은 항목의 표시 제목을 반환합니다. 사람이 지정한 name을 우선하고, 없으면 취약점 유형 vulnclass를 사용하며,
// 둘 다 비어 있으면 자리표시자를 사용합니다. 빈 제목은 절대로 출력하지 않습니다.
func (i Item) Title() string {
	if i.Name != "" {
		return i.Name
	}
	if i.VulnClass != "" {
		return i.VulnClass
	}
	return "(이름 없는 취약점)"
}

// Message는 채널에서 한 번 전송할 전체 내용입니다.
type Message struct {
	// 단일 알림 전송 시 길이는 1이며, 모아 보내기(digest) 시에는 배치 전체입니다.
	// 빈 슬라이스는 유효하지 않으므로 호출자는 하나 이상의 항목을 보장해야 합니다.
	Items []Item
	// Batch=true이면 모아 보내기 메시지로 렌더링합니다(제목 변경, 시간 구간과 건수 포함).
	Batch bool
	// WindowMinutes는 모아 보내기 주기(분)이며, Batch=true일 때만 "최근 N분" 문구에 사용합니다.
	// 렌더링의 결정성을 유지하여 테스트하기 쉽도록, 렌더링 시 time.Since로 계산하지 않고 설정에서 명시적으로 전달합니다.
	WindowMinutes int
	// HomeURL은 플랫폼 패널 주소(전역 public_base_url)이며, 비어 있으면 패널로 이동하는 링크를 넣지 않습니다.
	HomeURL string
}
