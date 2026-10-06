// Package notify는 취약점 발견의 IM / 메일 알림 전송 채널 어댑터 계층을 구현합니다.
//
// 계층: 이 패키지는 **최하위 패키지**로 표준 라이브러리에만 의존하며, DB나 server를 알지 못합니다. 채널 설정은
// map[string]any로 전달합니다(notification_channels.config JSONB 열에 대응).
// 전송할 내용은 Message로 전달합니다. 이렇게 분리하면 서명 계산·UTF-8 잘림·필터 매칭처럼
// 오류가 발생하기 쉬운 부분을 PostgreSQL 없이 단위 테스트할 수 있고, 호스트는 server 측에서 오케스트레이션만 하면 됩니다.
//
// 동시성 규약: Channel 구현은 반드시 **상태를 가지지 않아야** 합니다. 동일한 Channel 인스턴스를 여러 채널 설정에서
// (동일 채널의 여러 봇 인스턴스에서도) 동시에 재사용하며, 모든 자격 증명은 cfg 파라미터로 전달합니다.
// webhook URL 등을 구현 자체의 필드에 캐시해서는 안 됩니다.
package notify

// 채널 유형 식별자입니다. 값은 notification_channels.kind의 유효한 집합이기도 하며, server 측에서
// 허용 목록으로 검증합니다(findings.status와 마찬가지로, 추후 채널 추가를 쉽게 하도록 DB CHECK는 사용하지 않음).
const (
	KindDingTalk = "dingtalk" // DingTalk 사용자 지정 봇
	KindFeishu   = "feishu"   // Feishu(Lark 포함) 사용자 지정 봇
	KindWeCom    = "wecom"    // WeCom(기업용 위챗) 그룹 봇
	KindWebhook  = "webhook"  // 일반 Webhook: 사용자 지정 메서드/헤더/JSON 템플릿
	KindTelegram = "telegram" // Telegram Bot API
	KindEmail    = "email"    // SMTP 메일
)

// 이벤트 유형이며 notification_events.kind에 대응합니다.
const (
	EventFindingCreated       = "finding_created"
	EventFindingStatusChanged = "finding_status_changed"
)

// InitKind는 config의 kind가 비어 있을 때 사용하는 기본값입니다.
const InitKind = KindDingTalk

// severityRank는 취약점 심각도를 비교 가능한 순위로 매핑합니다. 알 수 없는 심각도는 0을 반환하므로, 모든
// min_severity 설정이 알 수 없는 심각도를 차단합니다. 의심스러우면 보내지 않아 잘못된 알림이 화면을 채우는 것을 막습니다.
var severityRank = map[string]int{
	"low":      1,
	"medium":   2,
	"high":     3,
	"critical": 4,
}

// SeverityRank는 심각도의 순위를 반환합니다. 알 수 없는 심각도는 0을 반환합니다.
func SeverityRank(severity string) int { return severityRank[severity] }

// SeverityLabel은 메시지 제목과 카드 색상에 사용할 emoji가 포함된 한국어 심각도 이름을 반환합니다.
// 알 수 없는 심각도는 그대로 반환하며, 임의로 만들어 내지 않습니다.
func SeverityLabel(severity string) string {
	switch severity {
	case "critical":
		return "🔴 심각"
	case "high":
		return "🟠 높음"
	case "medium":
		return "🟡 중간"
	case "low":
		return "🔵 낮음"
	default:
		return severity
	}
}

// StatusLabel은 상태 변경 메시지에 사용할 처리 상태를 한국어로 번역합니다.
func StatusLabel(status string) string {
	switch status {
	case "pending":
		return "처리 대기"
	case "in_progress":
		return "처리 중"
	case "confirmed":
		return "확인됨"
	case "resolved":
		return "처리됨"
	case "fixed":
		return "수정 완료"
	case "false_positive":
		return "오탐"
	case "ignored":
		return "무시"
	case "duplicate":
		return "중복"
	case "risk_accepted":
		return "위험 수용"
	default:
		return status
	}
}

// AtLeast는 severity가 min 기준 이상인지 판정합니다. min이 비어 있으면 기준을 두지 않고 모두 허용합니다.
// 알 수 없는 severity의 순위는 0이므로 비어 있지 않은 모든 min에서 거부됩니다(severityRank 주석 참조).
func AtLeast(severity, min string) bool {
	if min == "" {
		return true
	}
	return SeverityRank(severity) >= SeverityRank(min)
}
