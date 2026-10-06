package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 이 파일은 전달 작업의 획득과 상태 전이다.
//
// 획득은 긴 트랜잭션 대신 '리스'를 쓴다: 행을 sending으로 두고 next_attempt_at을 미래로 밀어
// 리스 만료 시간으로 삼은 뒤, 트랜잭션을 커밋하고 나서 네트워크 전달을 한다. 이렇게 하면 전달 중 DB 락을 쥐지 않는다 ——
// 네트워크 요청은 수 초 걸릴 수 있어(클라이언트 타임아웃 15초), 행 락을 계속 쥐면 같은 DB의 다른 쓰기를 마비시킨다.
//
// 대가는 프로세스가 전달 도중 크래시하면 행이 sending에 멈추는 것이다. 이건 **자가 치유 가능**하다: 리스 만료 후
// next_attempt_at이 과거로 떨어지면, 다음 라운드 획득이 같은 행을 다시 집어 올린다(획득 조건의
// state IN ('pending','sending') 참조). 재시도 카운트는 획득 시 이미 +1 되므로, 크래시가
// 무한 재시도를 유발하지 않는다 —— MaxNotifyAttempts회 기회를 다 쓰면 failed로 떨어져 수동 처리를 기다린다.

// MaxNotifyAttempts는 한 전달의 최대 시도 횟수다(첫 시도 포함).
// 전달 엔진이 아니라 여기에 정의한다: 상태 기계 자체의 정책이고, 엔진은 실행자일 뿐이다.
const MaxNotifyAttempts = 3

// MaxDigestBatchSize는 단일 요약 배치가 한 번에 병합할 수 있는 최대 전달 건수다.
//
// 존재 이유는 자원이다: 한 요약 주기에 취약점 수만 개가 나오면(충분히 가능 —— 전량 스캔
// 한 번이면 된다), 상한이 없으면 획득이 전체 행을 메모리에 읽어 초장문 메시지로 렌더링하고,
// 그다음 채널 길이 상한에 절반 이상 잘린다 —— 메모리도 낭비하고 잘린 취약점을 **조용히 손실**한다.
// 상한을 두면 초과분은 DB에 남아 다음 배치가 되고, 다음 주기에 자연히 보내져 손실되지 않는다.
//
// 500을 택한 근거: 메시지로 렌더링했을 때 WeCom 4096바이트 상한 안에서도 '읽을 내용이 있는' 수준이다;
// 더 키워도 절단 위치만 더 뒤로 밀 뿐이다.
const MaxDigestBatchSize = 500

// NotificationDelivery는 렌더링에 필요한 채널 설정과 이벤트 스냅샷을 포함한 전달 작업 하나다.
type NotificationDelivery struct {
	ID            int64           `json:"id"`
	EventID       int64           `json:"event_id"`
	ChannelID     int64           `json:"channel_id"`
	State         string          `json:"state"`
	Attempts      int             `json:"attempts"`
	NextAttemptAt time.Time       `json:"next_attempt_at"`
	LastError     string          `json:"last_error"`
	BatchID       *int64          `json:"batch_id,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	SentAt        *time.Time      `json:"sent_at,omitempty"`
	Snapshot      json.RawMessage `json:"snapshot,omitempty"`
	// 조인으로 로드한 렌더링 컨텍스트, JSON에 들어가지 않음(server 레이어가 DTO 조립).
	Channel *NotificationChannel `json:"-"`
	// FindingID/EventKind는 이벤트에서 가져오며, 이력 목록이 취약점 상세로 바로 이동하도록 한다.
	FindingID int64  `json:"finding_id,string"`
	EventKind string `json:"event_kind"`
	// ChannelName/ChannelKind는 목록 표시용 중복 필드로, 프런트의 2차 조회를 없앤다.
	ChannelName string `json:"channel_name"`
	ChannelKind string `json:"channel_kind"`
}

const notificationDeliveryCols = `d.id, d.event_id, d.channel_id, d.state, d.attempts, d.next_attempt_at,
       d.last_error, d.batch_id, d.created_at, d.sent_at`

// joinedDeliveryQuery는 전달 행의 통합 읽기 형태다: 전달 + 이벤트 스냅샷 + 채널 설정.
// 메시지 하나를 렌더링하려면 셋 다 필수라, 따로 조회하면 왕복이 세 번이 된다.
const joinedDeliveryQuery = `SELECT ` + notificationDeliveryCols + `,
       e.snapshot, e.kind, e.finding_id,
       c.id, c.name, c.kind, c.enabled, c.config, c.mode, c.filter, c.rate_per_min
FROM notification_deliveries d
JOIN notification_events e ON e.id = d.event_id
JOIN notification_channels c ON c.id = d.channel_id`

func scanNotificationDelivery(sc interface{ Scan(...any) error }) (*NotificationDelivery, error) {
	var (
		dl        NotificationDelivery
		lastErr   sql.NullString
		batchID   sql.NullInt64
		sentAt    sql.NullTime
		snapshot  []byte
		eventKind string
		channel   NotificationChannel
		chEnabled bool
	)
	if err := sc.Scan(&dl.ID, &dl.EventID, &dl.ChannelID, &dl.State, &dl.Attempts, &dl.NextAttemptAt,
		&lastErr, &batchID, &dl.CreatedAt, &sentAt,
		&snapshot, &eventKind, &dl.FindingID,
		&channel.ID, &channel.Name, &channel.Kind, &chEnabled, &channel.Config, &channel.Mode, &channel.Filter, &channel.RatePerMin); err != nil {
		return nil, err
	}
	dl.LastError = lastErr.String
	if batchID.Valid {
		dl.BatchID = &batchID.Int64
	}
	if sentAt.Valid {
		dl.SentAt = &sentAt.Time
	}
	dl.Snapshot = json.RawMessage(snapshot)
	dl.EventKind = eventKind
	dl.ChannelName = channel.Name
	dl.ChannelKind = channel.Kind
	channel.Enabled = &chEnabled
	dl.Channel = &channel
	return &dl, nil
}

// claimQuery는 한 번의 획득을 기술한다: 먼저 sel로 후보를 골라 잠그고, 그다음 sending으로 두고
// 리스를 연장한다. sel의 lease 위치는 호출자가 $n으로 자리를 두고 직접 인자를 전달한다.
type claimQuery struct {
	sql  string
	args []any
}

// ClaimRealtimeDeliveries는 어떤 채널의 만기된 실시간 전달 한 묶음을 획득한다(최대 limit건).
//
// 일부러 **단일 채널** 단위로 획득한다, '전역으로 한 묶음 획득 후 골라 보내기'가 아니다: 레이트 리밋 게이트는 전달 엔진에서 채널별로
// 유지되므로, 이 채널이 이번 라운드에 몇 건 더 보낼 수 있는지 먼저 알고 그만큼만 행을 획득해야 레이트 리밋이
// 재시도 횟수를 소모하지 않는다. 반대로 먼저 획득하고 버리면, 레이트 리밋에 막힌 행은 이미 attempts가 한 번 계산되어
// 3회 예산이 순전히 대기로 소진되어 결국 failed로 떨어진다.
//
// 조건에 '리스 만료된 sending'을 포함한다 —— 그게 크래시 자가 치유의 지점이다. lease는 단일
// 전달의 최악 소요 시간(채널 HTTP 클라이언트 타임아웃 15초)보다 훨씬 커야 한다, 안 그러면 같은 행을 두 dispatcher가
// 동시에 전달한다. 비활성 채널도 함께 차단한다: 비활성화 작업이 기존 전달을 skipped로 표시하지만,
// 여기서 한 번 더 막아 비활성화와 획득이 동시에 일어날 때의 누락을 방지한다.
func (d *DB) ClaimRealtimeDeliveries(ctx context.Context, channelID int64, limit int, lease time.Duration) ([]*NotificationDelivery, error) {
	if limit <= 0 {
		return nil, nil
	}
	return d.claimDeliveries(ctx, lease, claimQuery{
		sql: `SELECT dd.id FROM notification_deliveries dd
JOIN notification_channels c ON c.id = dd.channel_id
WHERE dd.channel_id = $1 AND dd.state IN ($2,$3) AND dd.next_attempt_at <= now()
  AND c.enabled AND c.mode = $4
ORDER BY dd.next_attempt_at, dd.id
FOR UPDATE OF dd SKIP LOCKED
LIMIT $5`,
		args: []any{channelID, NotifyStatePending, NotifyStateSending, NotifyModeRealtime, limit},
	}, nil)
}

// DigestBatchDue는 이 채널이 만기 배치를 충분히 모았는지 보고한다: 발송 대기 전달이 있고, **가장 오래된 것**의
// 나이가 요약 주기에 도달했는지.
//
// 판정 근거는 벽시계가 아니라 가장 오래된 전달의 나이다: 그래서 막 만든 채널이 정시 정렬 때문에
// 한 건짜리 '요약'을 즉시 뱉지 않고, 오래 적체된 배치도 한 라운드를 더 헛되이 기다리지 않는다.
//
// ClaimDigestBatch와 분리한 것은 의미가 다르기 때문이다: 이 함수는 '보낼지 말지'만 답하고,
// 획득은 이 채널의 **전체** 대기 행을 가져간다(아직 나이가 안 찬 것 포함) —— 안 그러면 한 주기가
// 여러 메시지로 쪼개져 요약의 의미가 사라진다.
func (d *DB) DigestBatchDue(ctx context.Context, channelID int64, minAge time.Duration) (bool, error) {
	var due bool
	err := d.QueryRowContext(ctx, `SELECT EXISTS (
  SELECT 1 FROM notification_deliveries d
  JOIN notification_channels c ON c.id = d.channel_id
  WHERE d.channel_id = $1 AND d.state IN ($2,$3) AND c.enabled
  GROUP BY d.channel_id
  HAVING min(d.created_at) <= now() - make_interval(secs => $4)
)`, channelID, NotifyStatePending, NotifyStateSending, int64(minAge.Seconds())).Scan(&due)
	return due, err
}

// ClaimDigestBatch는 어떤 채널의 현재 만기된 발송 대기 전달을 하나의 요약 배치로 획득한다,
// 한 배치 최대 MaxDigestBatchSize건.
//
// 같은 배치의 모든 전달이 batch_id를 공유한다, 집합의 최소 id를 배치 번호로 쓴다(안정·가독·
// 추가 시퀀스 불필요). 재시도 시 COALESCE로 원래 배치 번호를 보존해, '이 배치 N건은 함께 보냈다'가
// 여러 번 재시도 후에도 성립한다.
//
// id 오름차순으로 앞 N건을 취한다(무작위 아님): 가장 먼저 생긴 전달이 먼저 나가, 적체 시
// '새 취약점 먼저, 오래된 취약점은 영원히 뒤'라는 기아가 생기지 않는다.
func (d *DB) ClaimDigestBatch(ctx context.Context, channelID int64, limit int, lease time.Duration) ([]*NotificationDelivery, error) {
	if limit <= 0 {
		return nil, nil
	}
	// limit는 **메모리 상한**이고, 호출자가 MaxDigestBatchSize를 전달한다; 여기서 한 번 더 조여,
	// 호출자가 더 큰 값을 넘기는 것을 방지한다.
	//
	// 일부러 '레이트 리밋 할당량'을 배치 크기로 받지 않는다: 레이트 리밋 단위는 메시지 건수다 —— 한 배치는
	// 메시지 한 건만 보내고 토큰 하나를 소비하며, server 레이어의 takeTokens가 차감한다 —— '한 배치에 취약점
	// 몇 건'과는 다른 차원이다. 예전엔 rate_per_min을 digest에 적용하려고 매 라운드
	// 요청 예산을 배치 크기로 넘겼는데, 결과적으로 rate=20/min 채널은 배치당 취약점 1건만 담아,
	// digest가 요약 문구가 붙은 실시간 푸시로 퇴화했다. 레이트 리밋을 바꾸려면 takeTokens의 want를 바꾸고,
	// 여기는 건들지 마라.
	if limit > MaxDigestBatchSize {
		limit = MaxDigestBatchSize
	}
	out, err := d.claimDeliveries(ctx, lease, claimQuery{
		sql: `SELECT dd.id FROM notification_deliveries dd
JOIN notification_channels c ON c.id = dd.channel_id
WHERE dd.channel_id = $1 AND dd.state IN ($2,$3) AND dd.next_attempt_at <= now() AND c.enabled
ORDER BY dd.id
FOR UPDATE OF dd SKIP LOCKED
LIMIT $4`,
		args: []any{channelID, NotifyStatePending, NotifyStateSending, limit},
	}, func(tx *sql.Tx, ids []int64) error {
		batchID := ids[0]
		for _, id := range ids {
			if id < batchID {
				batchID = id
			}
		}
		ph, idArgs := placeholders(2, ids)
		_, err := tx.ExecContext(ctx, `UPDATE notification_deliveries SET batch_id = COALESCE(batch_id, $1)
WHERE id IN (`+ph+`)`, append([]any{batchID}, idArgs...)...)
		return err
	})
	return out, err
}

// claimDeliveries는 '선택 + sending 설정·리스 연장 + 전체 행 읽기'를 전부 한 트랜잭션에서 수행한다.
// postClaim은 선택적 부가 단계다(요약 배치가 이를 써서 batch_id를 기록).
func (d *DB) claimDeliveries(ctx context.Context, lease time.Duration, cq claimQuery, postClaim func(*sql.Tx, []int64) error) ([]*NotificationDelivery, error) {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck // 커밋 성공 후에는 no-op

	ids, err := selectForClaim(ctx, tx, cq.sql, cq.args...)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, tx.Commit()
	}
	// sending으로 두고 next_attempt_at을 미래로 민다: 이 미래 시각이 곧 리스 만료 시간이라,
	// '리스 미만료'와 '재시도 시간 미도달'이 같은 조건식을 공유하게 되어, 새 열이 필요 없다.
	ph, idArgs := placeholders(3, ids)
	if _, err := tx.ExecContext(ctx, `UPDATE notification_deliveries
SET state=$1, attempts=attempts+1, next_attempt_at=now()+make_interval(secs => $2)
WHERE id IN (`+ph+`)`,
		append([]any{NotifyStateSending, lease.Seconds()}, idArgs...)...); err != nil {
		return nil, err
	}
	if postClaim != nil {
		if err := postClaim(tx, ids); err != nil {
			return nil, err
		}
	}
	out, err := loadDeliveriesTx(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	return out, tx.Commit()
}

func selectForClaim(ctx context.Context, tx *sql.Tx, query string, args ...any) ([]int64, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func loadDeliveriesTx(ctx context.Context, tx *sql.Tx, ids []int64) ([]*NotificationDelivery, error) {
	ph, args := placeholders(1, ids)
	rows, err := tx.QueryContext(ctx, joinedDeliveryQuery+` WHERE d.id IN (`+ph+`) ORDER BY d.id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*NotificationDelivery{}
	for rows.Next() {
		dl, err := scanNotificationDelivery(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, dl)
	}
	return out, rows.Err()
}

// MarkDeliveriesSent는 전달 한 묶음을 전달됨으로 표시한다.
func (d *DB) MarkDeliveriesSent(ctx context.Context, ids []int64) error {
	ph, args := placeholders(2, ids)
	if len(args) == 0 {
		return nil
	}
	_, err := d.ExecContext(ctx, `UPDATE notification_deliveries
SET state=$1, sent_at=now(), last_error='' WHERE id IN (`+ph+`)`, append([]any{NotifyStateSent}, args...)...)
	return err
}

// RescheduleDeliveries는 전달 한 묶음을 pending으로 되돌리고 재시도 시간을 미룬다.
//
// 새 중간 상태를 도입하지 않고 pending으로 되돌리는 것은, '기회가 몇 번 남았나'를 한 곳에서만
// 표현하기 위함이다(MaxNotifyAttempts). 상태 기계 분기가 재시도 정책에 따라 팽창하는 것을 막는다.
func (d *DB) RescheduleDeliveries(ctx context.Context, ids []int64, delay time.Duration, errMsg string) error {
	ph, args := placeholders(4, ids)
	if len(args) == 0 {
		return nil
	}
	_, err := d.ExecContext(ctx, `UPDATE notification_deliveries
SET state=$1, next_attempt_at=now()+make_interval(secs => $2), last_error=$3
WHERE id IN (`+ph+`)`,
		append([]any{NotifyStatePending, delay.Seconds(), truncateNotifyError(errMsg)}, args...)...)
	return err
}

// DeferDeliveries는 전달 한 묶음을 pending으로 되돌려 즉시 재획득 가능하게 하고, **획득 시 계산한 그 시도 한 번을 취소**한다.
//
// 용도는 하나뿐이다: 요약 메시지를 채널 길이 상한에 맞춰 분할 전송할 때, 이 건에 못 담은 항목은 다음 배치로 남겨야 한다.
// 그건 실패가 아니라서 재시도 예산을 소모해선 안 된다 —— 획득 시 attempts가 낙관적으로 +1 되었으니,
// 여기서 반드시 되돌려야 한다. 안 그러면 500건 적체가 단락당 20건으로 25단락으로 쪼개져,
// 꼬리 항목이 3번째 단락에서 MaxNotifyAttempts로 failed 판정되는데, 그들은 아무 오류도 낸 적이 없다.
//
// GREATEST(...,0)은 '누가 수동 재전송으로 attempts를 0으로 만든 뒤 다시 여기 온' 경우를 막아,
// 카운트가 음수가 되지 않게 한다.
func (d *DB) DeferDeliveries(ctx context.Context, ids []int64, reason string) error {
	ph, args := placeholders(3, ids)
	if len(args) == 0 {
		return nil
	}
	_, err := d.ExecContext(ctx, `UPDATE notification_deliveries
SET state=$1, attempts=GREATEST(attempts-1, 0), next_attempt_at=now(), last_error=$2
WHERE id IN (`+ph+`)`,
		append([]any{NotifyStatePending, truncateNotifyError(reason)}, args...)...)
	return err
}

// FailDeliveries는 전달 한 묶음을 최종 실패로 표시하고, 전달 이력에서 수동 재전송을 기다린다.
func (d *DB) FailDeliveries(ctx context.Context, ids []int64, errMsg string) error {
	// 자리표시자는 $3부터 시작한다: $1은 state, $2는 last_error.
	ph, args := placeholders(3, ids)
	if len(args) == 0 {
		return nil
	}
	_, err := d.ExecContext(ctx, `UPDATE notification_deliveries SET state=$1, last_error=$2 WHERE id IN (`+ph+`)`,
		append([]any{NotifyStateFailed, truncateNotifyError(errMsg)}, args...)...)
	return err
}

// RetryNotificationDelivery는 전달 하나를 수동 재전송한다: pending으로 재설정·재시도 카운트 0으로·
// 즉시 만기. 카운트 초기화는 의도적이다 —— 사람이 '재전송'을 누른 것은 이전 실패 원인이 처리됐다는 뜻이라,
// 구 카운트로 다시 제한하는 것은 말이 안 된다.
func (d *DB) RetryNotificationDelivery(ctx context.Context, id int64) error {
	res, err := d.ExecContext(ctx, `UPDATE notification_deliveries
SET state=$2, attempts=0, next_attempt_at=now(), last_error=''
WHERE id=$1 AND state IN ($3,$4)`, id, NotifyStatePending, NotifyStateFailed, NotifyStateSkipped)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("전달 %d이(가) 존재하지 않거나 현재 상태에서 재전송할 수 없습니다", id)
	}
	return nil
}

// NotificationDeliveryFilter는 전달 이력의 조회 조건이다.
type NotificationDeliveryFilter struct {
	ChannelID int64
	State     string
	EventKind string
}

func (f NotificationDeliveryFilter) where() (string, []any) {
	var conds []string
	var args []any
	if f.ChannelID > 0 {
		args = append(args, f.ChannelID)
		conds = append(conds, fmt.Sprintf("d.channel_id=$%d", len(args)))
	}
	if f.State != "" {
		args = append(args, f.State)
		conds = append(conds, fmt.Sprintf("d.state=$%d", len(args)))
	}
	if f.EventKind != "" {
		args = append(args, f.EventKind)
		conds = append(conds, fmt.Sprintf("e.kind=$%d", len(args)))
	}
	if len(conds) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

// ListNotificationDeliveries는 전달 이력을 페이지네이션으로 반환한다, 새 것이 앞.
func (d *DB) ListNotificationDeliveries(ctx context.Context, f NotificationDeliveryFilter, page, pageSize int) ([]*NotificationDelivery, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 50
	}
	where, args := f.where()

	var total int
	if err := d.QueryRowContext(ctx, `SELECT count(*) FROM notification_deliveries d
JOIN notification_events e ON e.id = d.event_id`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := fmt.Sprintf("%s%s ORDER BY d.id DESC LIMIT $%d OFFSET $%d",
		joinedDeliveryQuery, where, len(args)+1, len(args)+2)
	rows, err := d.QueryContext(ctx, q, append(args, pageSize, (page-1)*pageSize)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*NotificationDelivery{}
	for rows.Next() {
		dl, err := scanNotificationDelivery(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, dl)
	}
	return out, total, rows.Err()
}

// truncateNotifyError는 에러 메시지를 열이 수용 가능한 길이로 자른다. 채널이 반환한 응답 본문은 매우 길 수 있어
// (범용 Webhook이 자체 구축 서비스를 칠 때 특히), 자르지 않으면 이력 목록 페이로드가 팽창한다.
func truncateNotifyError(msg string) string {
	const max = 500
	if len(msg) <= max {
		return msg
	}
	// 문자 경계로 되돌려, 반쪽 UTF-8 문자가 남아 프런트에 깨져 보이는 것을 방지한다.
	cut := max
	for cut > 0 && !isUTF8Start(msg[cut]) {
		cut--
	}
	return msg[:cut] + "…"
}

func isUTF8Start(b byte) bool { return b&0xC0 != 0x80 }

// placeholders는 start부터 시작하는 $n 자리표시자 문자열과 대응 인자를 생성한다, IN (...)에 사용.
// 예를 들어 start=3, ids=[7,8] → "$3,$4", [7,8].
func placeholders(start int, ids []int64) (string, []any) {
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		ph = append(ph, fmt.Sprintf("$%d", start+i))
		args = append(args, id)
	}
	return strings.Join(ph, ","), args
}
