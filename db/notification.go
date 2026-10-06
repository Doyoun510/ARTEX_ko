package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Autumn-27/artex/notify"
)

// 이 파일은 IM 푸시의 채널 설정·이벤트 레이어다. 전송 작업의 획득과 상태 전이는
// db/notification_delivery.go 참조.
//
// 두 불변식, 이 파일을 수정할 때 반드시 유지할 것:
//
//  1. 취약점 쓰기 트랜잭션(RecordFindingTx)은 InsertNotificationEventTx만 호출해 한 번 블라인드 삽입하고,
//     알림 관련 테이블을 읽지 않고 필터 매칭도 하지 않는다. 여기에 도입하는 어떤 읽기든
//     사용자가 잘못 설정한 필터 조건 때문에 취약점 쓰기 트랜잭션을 오염시키거나 중단시킬 수 있다.
//  2. 필터 매칭은 절대 오류를 내지 않는다: 설정이 기형이면 일괄 '적중'으로 처리(notify.Match 참조). 차라리 더 보낼지언정
//     빠뜨려서는 안 된다.

// ErrNotificationChannelNotFound 채널이 존재하지 않음.
var ErrNotificationChannelNotFound = errors.New("알림 채널이 존재하지 않습니다")

// 전송 상태.
const (
	NotifyStatePending = "pending" // 전송 대기
	NotifyStateSending = "sending" // 어떤 dispatcher가 획득, lease 미만료
	NotifyStateSent    = "sent"    // 전송됨
	NotifyStateFailed  = "failed"  // 재시도 소진 또는 영구 실패, 수동 재전송 가능
	NotifyStateSkipped = "skipped" // 채널 비활성화, 더 이상 전송 안 함
)

// 푸시 모드.
const (
	NotifyModeRealtime = "realtime"
	NotifyModeDigest   = "digest"
)

// ValidNotifyMode는 푸시 모드를 화이트리스트로 검증한다(findings.status와 동일 원리: DB CHECK를 쓰지 않아
// 이후 확장이 쉽다).
func ValidNotifyMode(m string) bool {
	return m == NotifyModeRealtime || m == NotifyModeDigest
}

// NotificationChannel은 채널 인스턴스 설정 하나다. Config와 Filter는 원본 JSON을 유지하고,
// 파싱은 notify 패키지에 맡긴다 —— db 레이어는 그 필드 의미를 이해하지 않는다.
type NotificationChannel struct {
	ID     int64           `json:"id"`
	Name   string          `json:"name"`
	Kind   string          `json:"kind"`
	Mode   string          `json:"mode"`
	Config json.RawMessage `json:"config"`
	Filter json.RawMessage `json:"filter"`
	// Enabled가 포인터인 것은 '이 필드 미전달'과 '명시적 false 전달'을 구분하기 위함이다 ——
	// 프런트 토글 컨트롤은 변경된 필드만 제출한다.
	Enabled    *bool     `json:"enabled,omitempty"`
	RatePerMin int       `json:"rate_per_min"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IsEnabled는 채널 활성화 여부를 반환한다; Enabled가 nil(미로드)이면 활성으로 처리한다.
func (c *NotificationChannel) IsEnabled() bool { return c.Enabled == nil || *c.Enabled }

// NotificationEvent는 이벤트 사실 하나다.
type NotificationEvent struct {
	ID        int64           `json:"id"`
	Kind      string          `json:"kind"`
	FindingID int64           `json:"finding_id"`
	Snapshot  json.RawMessage `json:"snapshot"`
	CreatedAt time.Time       `json:"created_at"`
}

const notificationChannelCols = `id, name, kind, enabled, config, mode, filter, rate_per_min, created_at, updated_at`

func scanNotificationChannel(sc interface{ Scan(...any) error }) (*NotificationChannel, error) {
	var c NotificationChannel
	var enabled bool
	if err := sc.Scan(&c.ID, &c.Name, &c.Kind, &enabled, &c.Config, &c.Mode, &c.Filter, &c.RatePerMin, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	c.Enabled = &enabled
	return &c, nil
}

// ListNotificationChannels는 모든 채널 인스턴스를 반환한다. 활성이 앞, 동급은 id순.
// 정렬을 SQL에 둔 것은 UI와 dispatcher가 같은 안정 순서를 보게 하기 위함이다.
func (d *DB) ListNotificationChannels(ctx context.Context) ([]*NotificationChannel, error) {
	rows, err := d.QueryContext(ctx, `SELECT `+notificationChannelCols+` FROM notification_channels
ORDER BY enabled DESC, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*NotificationChannel{}
	for rows.Next() {
		c, err := scanNotificationChannel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// NotificationChannelByID는 단일 채널을 가져온다.
func (d *DB) NotificationChannelByID(ctx context.Context, id int64) (*NotificationChannel, error) {
	row := d.QueryRowContext(ctx, `SELECT `+notificationChannelCols+` FROM notification_channels WHERE id=$1`, id)
	c, err := scanNotificationChannel(row)
	if err == sql.ErrNoRows {
		return nil, ErrNotificationChannelNotFound
	}
	return c, err
}

// SaveNotificationChannel은 채널을 생성하거나 업데이트한다.
//
// 업데이트 시 호출자가 명시적으로 준 필드(non-nil / non-empty)만 덮어쓴다. 그래서 프런트가 일부만
// 수정한 드로어 폼을 제출할 수 있고, 표시하지 않은 config 필드를 되돌려 보낼 필요가 없다 —— 되돌려 보내면 오히려
// '마스킹 값이 실제 키를 덮어쓰는' 사고가 난다.
func (d *DB) SaveNotificationChannel(ctx context.Context, c *NotificationChannel) (int64, error) {
	if c.Mode == "" {
		c.Mode = NotifyModeRealtime
	}
	// 여기서는 일부러 0을 가공하지 **않는다**: 0은 유효한 설정이며 '전송 속도 제한 없음'을 의미한다.
	//
	// 예전엔 `if c.RatePerMin <= 0 { c.RatePerMin = 기본값 }`로 썼는데, 의도는 '미지정 시
	// 안전 기본값 제공'이었지만 그건 '명시적으로 0 설정'도 함께 삼켜버렸다 —— 문서·UI 안내·
	// takeTokens 모두 0을 전송 속도 제한 없음으로 해석하는데, 여기서만 몰래 20(DingTalk/WeCom/Telegram)
	// 또는 100(Feishu)으로 바꿔, 조작자는 전송 속도 제한을 풀었다고 여기지만 실제로는 20/분에 막히고 아무 안내도 없었다.
	//
	// '미지정'과 '명시적 0'의 구분은 호출자만 안다(요청 본문의 필드 누락 vs 명시적 0 전달),
	// 그래서 기본값은 server 레이어가 필드 누락 시 채운다. notifyCreateChannel 참조.
	if c.RatePerMin < 0 {
		return 0, errors.New("전송 속도 제한 값은 음수일 수 없습니다")
	}
	if c.Config == nil {
		c.Config = json.RawMessage(`{}`)
	}
	if c.Filter == nil {
		c.Filter = json.RawMessage(`{}`)
	}
	enabled := c.IsEnabled()

	if c.ID == 0 {
		var id int64
		err := d.QueryRowContext(ctx, `INSERT INTO notification_channels(name,kind,enabled,config,mode,filter,rate_per_min)
VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
			c.Name, c.Kind, enabled, string(c.Config), c.Mode, string(c.Filter), c.RatePerMin).Scan(&id)
		return id, err
	}
	res, err := d.ExecContext(ctx, `UPDATE notification_channels
SET name=$2, kind=$3, enabled=$4, config=$5, mode=$6, filter=$7, rate_per_min=$8
WHERE id=$1`,
		c.ID, c.Name, c.Kind, enabled, string(c.Config), c.Mode, string(c.Filter), c.RatePerMin)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, ErrNotificationChannelNotFound
	}
	return c.ID, nil
}

// SetNotificationChannelEnabled는 활성/비활성을 전환한다.
//
// 채널을 비활성화할 때 아직 보내지 않은 전송을 함께 skipped로 표시한다: 안 그러면 재활성화 후
// '비활성 기간에 쌓인' 구 취약점을 한꺼번에 받게 되어, 시의성이 사라지고 신규로 오판되기 쉽다.
func (d *DB) SetNotificationChannelEnabled(ctx context.Context, id int64, enabled bool) error {
	return d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `UPDATE notification_channels SET enabled=$2 WHERE id=$1`, id, enabled)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrNotificationChannelNotFound
		}
		if !enabled {
			if _, err := tx.ExecContext(ctx, `UPDATE notification_deliveries SET state=$2, last_error=$3
WHERE channel_id=$1 AND state IN ($4,$5)`,
				id, NotifyStateSkipped, "채널 비활성화됨", NotifyStatePending, NotifyStateSending); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteNotificationChannel은 채널을 삭제한다. 그 전송 이력은 외래 키 캐스케이드로 삭제된다
// (채널 설정이 없어지면 이력을 해석할 수 없다).
func (d *DB) DeleteNotificationChannel(ctx context.Context, id int64) error {
	res, err := d.ExecContext(ctx, `DELETE FROM notification_channels WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotificationChannelNotFound
	}
	return nil
}

// RecordNotificationEventTx는 호출자의 트랜잭션에서 푸시 이벤트 하나를 **최선을 다해** 기록한다.
//
// 취약점 쓰기 경로에서 유일한 알림 관련 변경이다: INSERT 한 번, 어떤 테이블도 읽지 않고 채널도 모르고
// 필터도 돌리지 않는다. 트랜잭션 커밋이 '취약점 저장'과 '푸시 작업 존재'의 원자적 일관성을 보장하며,
// 커밋은 성공했는데 큐에 안 들어가 메시지가 영구 손실되는 창은 없다.
//
// 두 핵심 설계, 둘 다 즉흥적으로 쓴 것이 아니다:
//
//  1. **왜 SAVEPOINT를 쓰나**: PostgreSQL에서 트랜잭션 내 한 문장이라도 오류가 나면 전체 트랜잭션이
//     aborted 상태가 되고, 이후 모든 문장(COMMIT 포함)이 전부 실패한다. 그래서 '이 INSERT
//     오류를 무시하고 호출자가 계속 커밋하게 하기'는 PG에서 불가능하다 —— 세이브포인트로 오류를 이
//     문장에 격리하지 않는 한. 세이브포인트가 없으면 '전체 롤백'이라는 선택지만 남는다.
//
//  2. **왜 전체 롤백이 틀렸나**: 푸시는 편의 기능이고, 취약점 기록이 제품 본체다. 알림
//     테이블 문제(구 DB 미마이그레이션·디스크 순간 장애)가 고위험 취약점을 저장 못 하게 해서는 안 된다. 그래서 여기서 오류를
//     격리하고 로그를 남기고 false를 반환해, 취약점 쓰기는 정상 커밋되게 한다 —— 대가는 이 푸시 하나를 잃는 것이다.
//     bool을 반환하고 error를 반환하지 않는 것은 의도적이다: 호출자는 이를 쓰기 성패에 영향을 주는 오류로 취급해선 안 된다.
func RecordNotificationEventTx(ctx context.Context, tx *sql.Tx, kind string, findingID int64, snap notify.Snapshot) bool {
	raw, err := json.Marshal(snap)
	if err != nil {
		log.Printf("[notify] 푸시 이벤트 직렬화 실패 finding=%d: %v", findingID, err)
		return false
	}
	if _, err := tx.ExecContext(ctx, `SAVEPOINT notify_event`); err != nil {
		log.Printf("[notify] 세이브포인트 생성 실패 finding=%d: %v", findingID, err)
		return false
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO notification_events(kind,finding_id,snapshot) VALUES($1,$2,$3)`,
		kind, findingID, string(raw)); err != nil {
		log.Printf("[notify] 푸시 이벤트 쓰기 실패 finding=%d(취약점 기록은 영향 없음): %v", findingID, err)
		// 세이브포인트로 롤백해 트랜잭션을 aborted 상태에서 되살린다.
		if _, rbErr := tx.ExecContext(ctx, `ROLLBACK TO SAVEPOINT notify_event`); rbErr != nil {
			log.Printf("[notify] 세이브포인트 롤백 실패 finding=%d: %v", findingID, rbErr)
		}
		return false
	}
	// 세이브포인트를 해제해 긴 트랜잭션에서 쓸모없는 세이브포인트가 쌓이는 것을 방지한다.
	_, _ = tx.ExecContext(ctx, `RELEASE SAVEPOINT notify_event`)
	return true
}

// AddNotificationEvent는 InsertNotificationEventTx의 독립 트랜잭션 버전으로, 기존
// 트랜잭션 밖 호출부가 사용한다(예: 채널의 '테스트 메시지 전송', 실제 finding이 없음).
func (d *DB) AddNotificationEvent(ctx context.Context, kind string, findingID int64, snap notify.Snapshot) (int64, error) {
	raw, err := json.Marshal(snap)
	if err != nil {
		return 0, fmt.Errorf("알림 이벤트 스냅샷 직렬화 실패: %w", err)
	}
	var id int64
	err = d.QueryRowContext(ctx, `INSERT INTO notification_events(kind,finding_id,snapshot) VALUES($1,$2,$3) RETURNING id`,
		kind, findingID, string(raw)).Scan(&id)
	return id, err
}

// FanOutPendingEvents는 아직 분배되지 않은 취약점 이벤트를 현재 활성 채널에 따라 전송 작업으로 펼치고,
// 이번 라운드에 처리한 이벤트 수와 새로 만든 전송 수를 반환한다.
//
// 한 라운드 전체가 한 트랜잭션 안에서: 이벤트는 FOR UPDATE SKIP LOCKED로 획득하므로, 여러 프로세스가 동시에 돌아도
// 각자 다른 행을 획득한다(프로젝트의 아카이브 큐 획득이 같은 기법을 쓴다,
// db/task_archives.go의 completeNextArchiveJob 참조).
//
// 필터 매칭을 일부러 SQL이 아니라 Go 쪽에 둔다: 채널 필터 조건은 선택 필드들의 JSONB라,
// SQL로 여섯 조합의 매칭을 표현하면 쿼리 유지보수가 어렵고, 채널 수는 '사람이 손으로 설정한 몇 개'라
// 전량 로드 후 메모리에서 하나씩 비교하는 편이 더 빠르고 테스트도 쉽다.
//
// 어떤 채널에도 적중하지 않은 이벤트도 fanned_out으로 표시된다 —— 안 그러면 영원히 분배 대기 집합에 남아
// 매 tick마다 다시 스캔된다.
func (d *DB) FanOutPendingEvents(ctx context.Context, limit int) (eventCount, deliveryCount int, err error) {
	if limit <= 0 {
		limit = 200
	}
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 커밋 성공 후에는 no-op

	channels, err := listEnabledNotificationChannelsTx(ctx, tx)
	if err != nil {
		return 0, 0, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT id, kind, finding_id, snapshot FROM notification_events
WHERE NOT fanned_out ORDER BY id FOR UPDATE SKIP LOCKED LIMIT $1`, limit)
	if err != nil {
		return 0, 0, err
	}
	var (
		events      []NotificationEvent
		parsedSnaps []notify.Snapshot
	)
	for rows.Next() {
		var ev NotificationEvent
		if err := rows.Scan(&ev.ID, &ev.Kind, &ev.FindingID, &ev.Snapshot); err != nil {
			rows.Close()
			return 0, 0, err
		}
		var snap notify.Snapshot
		// 스냅샷은 우리가 직접 쓴 것이라 이론상 반드시 파싱 가능하다; 파싱 실패는 전송 흐름을 막지 않지만,
		// 이 이벤트는 필드가 전부 비어 필터 조건이 있는 모든 채널이 건너뛴다 —— 차라리 하나 덜 보낼지언정
		// 나쁜 행 하나가 전체 큐를 막게 하지 않는다.
		_ = json.Unmarshal(ev.Snapshot, &snap)
		// kind는 행 내 값을 기준으로 한다: 스냅샷의 것은 렌더링용 사본이라 구 버전이 썼을 수 있다.
		snap.Kind = ev.Kind
		events = append(events, ev)
		parsedSnaps = append(parsedSnaps, snap)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, 0, err
	}
	if len(events) == 0 {
		return 0, 0, tx.Commit()
	}

	type pending struct {
		eventID   int64
		channelID int64
	}
	var toInsert []pending
	for i, snap := range parsedSnaps {
		for _, ch := range channels {
			if !notify.Match(notify.ParseFilter(ch.Filter), snap) {
				continue
			}
			toInsert = append(toInsert, pending{eventID: events[i].ID, channelID: ch.ID})
		}
	}
	if len(toInsert) > 0 {
		var (
			vals []string
			args []any
		)
		for _, p := range toInsert {
			vals = append(vals, fmt.Sprintf("($%d,$%d)", len(args)+1, len(args)+2))
			args = append(args, p.eventID, p.channelID)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO notification_deliveries(event_id,channel_id) VALUES `+strings.Join(vals, ","), args...); err != nil {
			return 0, 0, err
		}
	}

	// 이번 라운드 이벤트를 분배 완료로 표시한다. 어떤 채널에도 적중하지 않은 이벤트도 함께 표시한다(함수 주석 참조).
	ids := make([]string, 0, len(events))
	markArgs := make([]any, 0, len(events))
	for _, ev := range events {
		markArgs = append(markArgs, ev.ID)
		ids = append(ids, fmt.Sprintf("$%d", len(markArgs)))
	}
	if _, err := tx.ExecContext(ctx, `UPDATE notification_events SET fanned_out=true WHERE id IN (`+strings.Join(ids, ",")+`)`, markArgs...); err != nil {
		return 0, 0, err
	}
	return len(events), len(toInsert), tx.Commit()
}

// listEnabledNotificationChannelsTx는 트랜잭션에서 활성 채널을 가져온다. 수가 적어
// 페이지네이션도 캐시도 하지 않는다 —— 캐시는 '설정 변경이 언제 적용되나'라는 추가 타이밍 문제를 유발한다.
func listEnabledNotificationChannelsTx(ctx context.Context, tx *sql.Tx) ([]*NotificationChannel, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id, name, kind, config, mode, filter, rate_per_min
FROM notification_channels WHERE enabled ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*NotificationChannel{}
	for rows.Next() {
		var c NotificationChannel
		if err := rows.Scan(&c.ID, &c.Name, &c.Kind, &c.Config, &c.Mode, &c.Filter, &c.RatePerMin); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

// NotificationAssetNames는 자산 id를 짧은 표시명으로 해석하며, 푸시 메시지에 쓴다.
//
// 반환 순서는 입력 인자와 동일하고, 길이는 입력보다 작을 수 있다(존재하지 않는 id는 건너뜀). 입력 순서를 유지하는 것은
// 같은 취약점의 메시지가 여러 전송에서 자산 순서가 안정되게 하려는 것이다 —— 안 그러면 재시도 후 받은 메시지에서
// 자산 순서가 바뀌어 '자산이 바뀌었다'로 오독된다.
func (d *DB) NotificationAssetNames(ctx context.Context, ids []int64) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph, args := placeholders(1, ids)
	rows, err := d.QueryContext(ctx, `SELECT id, type, domain, ip, url, app_name, bundle_id FROM assets WHERE id IN (`+ph+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	labels := map[int64]string{}
	for rows.Next() {
		var (
			id                int64
			typ               string
			domain, ip, url   sql.NullString
			appName, bundleID sql.NullString
		)
		if err := rows.Scan(&id, &typ, &domain, &ip, &url, &appName, &bundleID); err != nil {
			return nil, err
		}
		labels[id] = assetDisplayName(typ, domain.String, ip.String, url.String, appName.String, bundleID.String)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(ids))
	seen := map[int64]bool{}
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if label, ok := labels[id]; ok && label != "" {
			out = append(out, label)
		}
	}
	return out, nil
}

// assetDisplayName은 자산 유형별로 가장 식별성 높은 식별자를 고른다.
// 최종적으로 빈 문자열을 반환해, 이름을 못 얻은 자산을 어떻게 표현할지는 호출자가 결정한다 —— 이 함수는 자리표시자를 지어내지 않는다,
// 안 그러면 '자산#42' 같은 노이즈가 푸시 메시지에 섞여 독자가 실제 도메인으로 오해한다.
func assetDisplayName(typ, domain, ip, url, appName, bundleID string) string {
	pick := func(vals ...string) string {
		for _, v := range vals {
			if strings.TrimSpace(v) != "" {
				return v
			}
		}
		return ""
	}
	switch typ {
	case "root_domain", "subdomain":
		return domain
	case "ip":
		return ip
	case "app":
		return pick(appName, bundleID)
	case "service", "endpoint":
		return pick(url, domain, ip)
	default:
		return pick(domain, ip, url, appName)
	}
}

// SetFindingStatusWithNotify는 취약점 처리 상태를 갱신하고, 같은 트랜잭션에서 상태 변경
// 푸시 이벤트 하나를 등록한다.
//
// 반환 from=변경 전 상태; found=취약점 존재 여부; notified=이벤트 등록 성공 여부.
//
// 의도적인 세 가지 동작:
//   - 상태가 실제로 변하지 않으면 이벤트를 등록하지 않는다. 프런트 드로어가 같은 값을 반복 제출하거나 자동화 스크립트가
//     멱등 재생하는 경우 모두 푸시 노이즈를 내서는 안 된다.
//   - 취약점이 없으면 found=false를 반환하고 아무 쓰기도 하지 않으며, 호출자가 404로 변환한다.
//   - 이벤트 등록 실패는 상태 갱신에 영향을 주지 않는다(RecordNotificationEventTx의 세이브포인트 설명 참조),
//     그래서 notified=false여도 상태는 이미 성공적으로 바뀐 것이니, 호출자는 이 때문에 오류를 내서는 안 된다.
func (d *DB) SetFindingStatusWithNotify(ctx context.Context, id int64, status string) (from string, found bool, notified bool, err error) {
	err = d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		var txErr error
		from, found, _, notified, txErr = SetFindingStatusTx(ctx, tx, id, status)
		return txErr
	})
	return from, found, notified, err
}

// SetFindingStatusTx는 **호출자의 트랜잭션** 안에서 취약점 상태를 갱신하고 상태 변경 푸시 이벤트를 등록한다.
//
// 트랜잭션 레벨 함수로 뽑은 것은 상태를 바꾸는 모든 경로가 같은 의미를 공유하게 하려는 것이다 —— 이전엔
// patchFinding만 알림 포함 버전을 썼고, **재검증 결론이 '수정 완료'일 때**(finding_retests의
// 그 `UPDATE findings SET status=...`)는 직접 DB에 썼기 때문에,
// `on_status_change`를 설정한 채널은 이런 상태 전이에 푸시를 전혀 못 받았다: 화면에서 상태가 조용히 바뀌고,
// 운영자는 플랫폼을 열어야 발견했다.
//
// 반환 from=변경 전 상태, found=취약점 존재 여부, changed=상태가 실제로 바뀌었는지,
// notified=이벤트 등록 성공 여부(등록 실패는 상태 갱신에 영향 없음, RecordNotificationEventTx 참조).
func SetFindingStatusTx(ctx context.Context, tx *sql.Tx, id int64, status string) (from string, found bool, changed bool, notified bool, err error) {
	var (
		vulnclass, name, severity, summary string
		taskID                             sql.NullInt64
		assetIDs                           []byte
	)
	scanErr := tx.QueryRowContext(ctx, `SELECT vulnclass, name, severity, summary, task_id, asset_ids, status
FROM findings WHERE id=$1 FOR UPDATE`, id).
		Scan(&vulnclass, &name, &severity, &summary, &taskID, &assetIDs, &from)
	if scanErr == sql.ErrNoRows {
		return "", false, false, false, nil
	}
	if scanErr != nil {
		return "", false, false, false, scanErr
	}
	found = true
	if from == status {
		// 상태가 실제로 변하지 않으면 이벤트를 등록하지 않는다: 같은 값 반복 제출·멱등 재생 모두
		// 푸시 노이즈를 내서는 안 된다.
		return from, true, false, false, nil
	}
	if _, err := tx.ExecContext(ctx, `UPDATE findings SET status=$2 WHERE id=$1`, id, status); err != nil {
		return from, true, false, false, err
	}
	var assets []int64
	_ = json.Unmarshal(assetIDs, &assets)
	notified = RecordNotificationEventTx(ctx, tx, notify.EventFindingStatusChanged, id, notify.Snapshot{
		Kind:       notify.EventFindingStatusChanged,
		FindingID:  id,
		TaskID:     taskID.Int64,
		VulnClass:  vulnclass,
		Name:       name,
		Severity:   severity,
		Summary:    summary,
		AssetIDs:   assets,
		FromStatus: from,
		ToStatus:   status,
	})
	return from, true, true, notified, nil
}

// NotificationStats는 알림 페이지 상단의 개요 카운트다.
type NotificationStats struct {
	Channels     int   `json:"channels"`
	ChannelsOn   int   `json:"channels_on"`
	Pending      int   `json:"pending"`
	Failed       int   `json:"failed"`
	SentToday    int   `json:"sent_today"`
	BacklogAgeMS int64 `json:"backlog_age_ms"` // 가장 오래된 전송 대기 항목의 현재까지 밀리초
}

// NotificationStatsSnapshot은 알림 시스템의 건강도를 집계한다.
// BacklogAgeMS는 '푸시가 막혔는지'의 가장 직접적인 지표다 —— pending 카운트보다 훨씬 유용한데,
// 적체 3건과 3건의 차이가 3초에서 3시간까지일 수 있기 때문이다.
func (d *DB) NotificationStatsSnapshot(ctx context.Context) (*NotificationStats, error) {
	var s NotificationStats
	if err := d.QueryRowContext(ctx, `SELECT
    (SELECT count(*) FROM notification_channels),
    (SELECT count(*) FROM notification_channels WHERE enabled),
    (SELECT count(*) FROM notification_deliveries WHERE state IN ($1,$2)),
    (SELECT count(*) FROM notification_deliveries WHERE state=$3),
    (SELECT count(*) FROM notification_deliveries WHERE state=$4 AND sent_at >= date_trunc('day', now())),
    COALESCE((SELECT EXTRACT(EPOCH FROM (now() - min(created_at))) * 1000 FROM notification_deliveries WHERE state=$1), 0)::bigint`,
		NotifyStatePending, NotifyStateSending, NotifyStateFailed, NotifyStateSent).
		Scan(&s.Channels, &s.ChannelsOn, &s.Pending, &s.Failed, &s.SentToday, &s.BacklogAgeMS); err != nil {
		return nil, err
	}
	return &s, nil
}
