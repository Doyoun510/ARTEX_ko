package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Autumn-27/artex/notify"
)

// 이 파일의 케이스는 모두 실제로 PostgreSQL에 연결한다(DB 없으면 건너뜀). 이 SQL들은
// FOR UPDATE SKIP LOCKED·make_interval·JSONB·여러 행 IN(...) 자리표시자 조립을 쓰는데,
// 모두 '컴파일은 되지만 런타임에 오류 날 수 있는' 작성법이라, 실제로 돌려야 검증된 것이다.

func notifyTestDB(t *testing.T) *DB {
	t.Helper()
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v) — skipping", err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

// newTestChannel은 채널을 하나 만들고, 테스트 종료 시 자동 삭제한다.
func newTestChannel(t *testing.T, d *DB, kind, mode string, filter string) *NotificationChannel {
	t.Helper()
	if filter == "" {
		filter = `{}`
	}
	ch := &NotificationChannel{
		Name:       "测试渠道-" + t.Name(),
		Kind:       kind,
		Mode:       mode,
		Config:     json.RawMessage(`{"webhook":"https://example.com/hook"}`),
		Filter:     json.RawMessage(filter),
		RatePerMin: 100,
	}
	id, err := d.SaveNotificationChannel(context.Background(), ch)
	if err != nil {
		t.Fatalf("채널 생성 실패: %v", err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_channels WHERE id=$1`, id) })
	ch.ID = id
	return ch
}

// addTestEvent는 이벤트를 하나 직접 쓴다(finding 경유 안 함), 분배·전달 테스트용.
func addTestEvent(t *testing.T, d *DB, kind string, findingID int64, snap notify.Snapshot) int64 {
	t.Helper()
	snap.Kind = kind
	snap.FindingID = findingID
	id, err := d.AddNotificationEvent(context.Background(), kind, findingID, snap)
	if err != nil {
		t.Fatalf("이벤트 쓰기 실패: %v", err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_events WHERE id=$1`, id) })
	return id
}

func TestNotificationAssetNamesResolvesAndPreservesOrder(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	// 세 가지 자산은 각자 표시 방식이 다르다: 도메인, IP, URL.
	insertAsset := func(query, value string) int64 {
		t.Helper()
		var id int64
		if err := d.QueryRow(query, value).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	domID := insertAsset(`INSERT INTO assets(type, domain) VALUES('subdomain',$1) RETURNING id`, "a.example.com")
	ipID := insertAsset(`INSERT INTO assets(type, ip) VALUES('ip',$1) RETURNING id`, "10.1.2.3")
	svcID := insertAsset(`INSERT INTO assets(type, url) VALUES('service',$1) RETURNING id`, "https://a.example.com/admin")
	t.Cleanup(func() {
		d.Exec(`DELETE FROM assets WHERE id IN ($1,$2,$3)`, domID, ipID, svcID)
	})

	// 입력 순서를 일부러 섞고, 존재하지 않는 id 하나를 포함한다.
	got, err := d.NotificationAssetNames(ctx, []int64{svcID, 999999999, domID, ipID, svcID})
	if err != nil {
		t.Fatalf("자산명 해석 실패: %v", err)
	}
	want := []string{"https://a.example.com/admin", "a.example.com", "10.1.2.3"}
	if len(got) != len(want) {
		t.Fatalf("자산명 개수 불일치, 기대 %v 얻음 %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("순서/값 불일치, 기대 %v 얻음 %v", want, got)
		}
	}
}

// TestRecordNotificationEventTxUnwindsOnFailure는 세이브포인트 메커니즘의 핵심 케이스다:
// 트랜잭션에서 먼저 notification_events 쓰기를 반드시 실패하게 만든 뒤(항상 false인 제약을 임시 추가),
// ① 이 함수가 false를 반환하고 ② 트랜잭션이 aborted 상태가 아니어서 이후 문장이 실행되는지 단언한다.
//
// 세이브포인트가 없으면 PostgreSQL은 전체 트랜잭션을 무효화하고, 이후 어떤 문장도
// "current transaction is aborted"로 실패한다 —— 그것이 바로 '알림 테이블 하나의 문제로
// 취약점을 저장하지 못하는' 장애 경로다.
//
// 여기서는 일부러 **COMMIT이 아니라 ROLLBACK으로 마무리**한다: ALTER TABLE은 PG에서 트랜잭션성이라,
// 한번 커밋하면 그 임시 제약이 schema에 영구히 남아 이후 모든 케이스를 함께 망가뜨린다.
// 롤백은 DDL을 자동으로 취소해 수동 정리가 필요 없다. 단언은 '트랜잭션이 아직 살아 있음'만 필요하고,
// 실제로 커밋할 필요는 없다.
func TestRecordNotificationEventTxUnwindsOnFailure(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	// 방어적 정리: 과거 실행이 이 제약을 남겼으면 먼저 제거한다.
	if _, err := d.Exec(`ALTER TABLE notification_events DROP CONSTRAINT IF EXISTS notify_test_never`); err != nil {
		t.Fatal(err)
	}

	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback() //nolint:errcheck // 임시 제약 취소, 함수 주석 참조

	// NOT VALID: 이후 쓰이는 행만 제약하고, DB에 이미 있는 과거 이벤트는 검증하지 않는다
	// (안 그러면 기존 행 위반으로 제약을 추가할 수 없다).
	if _, err := tx.ExecContext(ctx, `ALTER TABLE notification_events ADD CONSTRAINT notify_test_never CHECK (false) NOT VALID`); err != nil {
		t.Fatalf("임시 제약 추가 실패: %v", err)
	}
	if RecordNotificationEventTx(ctx, tx, notify.EventFindingCreated, 1, notify.Snapshot{Severity: "high"}) {
		t.Fatal("반드시 실패하는 제약 아래에서도 쓰기 성공을 보고함")
	}
	// 핵심 단언: 트랜잭션이 아직 쓸 수 있다.
	var one int
	if err := tx.QueryRowContext(ctx, `SELECT 1`).Scan(&one); err != nil {
		t.Fatalf("트랜잭션이 오염됨(세이브포인트 미작동): %v", err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatalf("롤백 실패: %v", err)
	}
	// DDL이 롤백과 함께 취소됐는지 확인해, 이후 케이스에 지뢰를 남기지 않는다.
	var exists bool
	if err := d.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_constraint WHERE conname='notify_test_never')`).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("임시 제약이 롤백으로 취소되지 않아 이후 케이스를 오염시킴")
	}
}

func TestFanOutRoutesEventsByFilter(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	all := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	onlyCritical := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{"min_severity":"critical"}`)
	sqlOnly := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{"vulnclass_include":["SQL"]}`)

	highSQL := addTestEvent(t, d, notify.EventFindingCreated, 1001, notify.Snapshot{Severity: "high", VulnClass: "SQL注入"})
	lowXSS := addTestEvent(t, d, notify.EventFindingCreated, 1002, notify.Snapshot{Severity: "low", VulnClass: "XSS"})
	criticalXSS := addTestEvent(t, d, notify.EventFindingCreated, 1003, notify.Snapshot{Severity: "critical", VulnClass: "XSS"})

	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatalf("분배 실패: %v", err)
	}

	cases := []struct {
		name    string
		eventID int64
		channel int64
		want    bool
	}{
		{"전수신 채널 high 수신", highSQL, all.ID, true},
		{"전수신 채널 low 수신", lowXSS, all.ID, true},
		{"심각만 채널 high 건너뜀", highSQL, onlyCritical.ID, false},
		{"심각만 채널 critical 수신", criticalXSS, onlyCritical.ID, true},
		{"SQL만 채널 SQL 수신", highSQL, sqlOnly.ID, true},
		{"SQL만 채널 XSS 건너뜀", lowXSS, sqlOnly.ID, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var exists bool
			if err := d.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM notification_deliveries WHERE event_id=$1 AND channel_id=$2)`,
				tc.eventID, tc.channel).Scan(&exists); err != nil {
				t.Fatal(err)
			}
			if exists != tc.want {
				t.Fatalf("전달 존재 여부: 기대 %v 얻음 %v", tc.want, exists)
			}
		})
	}

	// 다시 분배해도 중복 전달이 생기면 안 된다(fanned_out 멱등).
	events, deliveries, err := d.FanOutPendingEvents(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}
	if events != 0 || deliveries != 0 {
		t.Fatalf("이미 분배된 이벤트는 다시 처리되면 안 됨, 얻음 events=%d deliveries=%d", events, deliveries)
	}
}

// TestFanOutMarksEventsWithNoMatchingChannel은 '이벤트가 어떤 채널에도 적중하지 않은' 경우를 다룬다.
// 이런 이벤트도 똑같이 분배 완료로 표시되어야 한다, 안 그러면 영원히 분배 대기 집합에 남아 매 tick 재스캔된다.
func TestFanOutMarksEventsWithNoMatchingChannel(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	pick := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{"vulnclass_include":["绝不匹配的类型"]}`)
	_ = pick

	ev := addTestEvent(t, d, notify.EventFindingCreated, 2001, notify.Snapshot{Severity: "high", VulnClass: "XSS"})
	_, deliveries, err := d.FanOutPendingEvents(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deliveries != 0 {
		t.Fatalf("전달이 생기면 안 됨, 얻음 %d", deliveries)
	}
	var fanned bool
	if err := d.QueryRowContext(ctx, `SELECT fanned_out FROM notification_events WHERE id=$1`, ev).Scan(&fanned); err != nil {
		t.Fatal(err)
	}
	if !fanned {
		t.Fatal("채널에 적중하지 않은 이벤트도 분배 완료로 표시되어야 함, 안 그러면 무한 재스캔됨")
	}
}

func TestClaimRealtimeDeliveriesHonorsLeaseAndMode(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	realtime := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	digest := newTestChannel(t, d, notify.KindDingTalk, NotifyModeDigest, `{}`)

	addTestEvent(t, d, notify.EventFindingCreated, 3001, notify.Snapshot{Severity: "high", VulnClass: "XSS"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}

	// 실시간 획득은 realtime 채널의 것만 가져와야 하고, digest 채널 것은 건드리면 안 된다.
	got, err := d.ClaimRealtimeDeliveries(ctx, realtime.ID, 10, time.Minute)
	if err != nil {
		t.Fatalf("획득 실패: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("1건 획득해야 함, 얻음 %d", len(got))
	}
	if got[0].State != NotifyStateSending || got[0].Attempts != 1 {
		t.Fatalf("획득 후 sending이고 attempts=1이어야 함, 얻음 state=%s attempts=%d", got[0].State, got[0].Attempts)
	}
	// 조인으로 로드한 렌더링 컨텍스트가 완전해야 한다(채널 설정 + 이벤트 스냅샷 + finding id).
	if got[0].Channel == nil || len(got[0].Channel.Config) == 0 {
		t.Fatal("획득 결과에 채널 설정이 없어 렌더링이 실패함")
	}
	if got[0].FindingID != 3001 {
		t.Fatalf("finding id가 이벤트에서 전달되지 않음, 얻음 %d", got[0].FindingID)
	}

	// lease 미만료라 두 번째 획득은 비어야 한다 —— 이것이 '같은 행을 두 dispatcher가 동시에 전달하지 않음'
	// 의 보장이다.
	again, err := d.ClaimRealtimeDeliveries(ctx, realtime.ID, 10, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 0 {
		t.Fatalf("lease 기간 내 중복 획득은 안 됨, 얻음 %d건", len(again))
	}

	// digest 채널의 전달은 실시간 획득에 걸리면 안 된다.
	left, err := d.ClaimRealtimeDeliveries(ctx, digest.ID, 10, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Fatalf("실시간 획득은 digest 채널 전달을 가져오면 안 됨, 얻음 %d건", len(left))
	}
}

// TestClaimExpiredLeaseRecovers는 크래시 자가 치유를 다룬다: 프로세스가 전달 도중 죽으면 sending
// 행이 남는데, lease 만료 후 다시 획득될 수 있어야 한다, 안 그러면 이 전달은 영원히 막힌다.
func TestClaimExpiredLeaseRecovers(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	addTestEvent(t, d, notify.EventFindingCreated, 4001, notify.Snapshot{Severity: "high"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	first, err := d.ClaimRealtimeDeliveries(ctx, ch.ID, 10, time.Minute)
	if err != nil || len(first) != 1 {
		t.Fatalf("첫 획득 실패: %v (%d건)", err, len(first))
	}
	// lease를 수동으로 과거로 밀어 'lease 만료'를 시뮬레이션한다.
	if _, err := d.Exec(`UPDATE notification_deliveries SET next_attempt_at = now() - interval '1 minute' WHERE id=$1`, first[0].ID); err != nil {
		t.Fatal(err)
	}
	second, err := d.ClaimRealtimeDeliveries(ctx, ch.ID, 10, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(second) != 1 {
		t.Fatalf("lease 만료된 sending 행은 다시 획득될 수 있어야 함, 얻음 %d건", len(second))
	}
	if second[0].Attempts != 2 {
		t.Fatalf("재획득은 시도 횟수를 누적해야 함, 얻음 %d", second[0].Attempts)
	}
}

func TestClaimSkipsDisabledChannel(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	addTestEvent(t, d, notify.EventFindingCreated, 5001, notify.Snapshot{Severity: "high"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	// 비활성화는 기존 발송 대기 전달을 함께 skipped로 표시한다.
	if err := d.SetNotificationChannelEnabled(ctx, ch.ID, false); err != nil {
		t.Fatal(err)
	}
	var state string
	if err := d.QueryRow(`SELECT state FROM notification_deliveries WHERE channel_id=$1`, ch.ID).Scan(&state); err != nil {
		t.Fatal(err)
	}
	if state != NotifyStateSkipped {
		t.Fatalf("비활성화 채널의 기존 발송 대기 전달은 skipped로 표시되어야 함, 얻음 %s", state)
	}
	got, err := d.ClaimRealtimeDeliveries(ctx, ch.ID, 10, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("비활성화 채널은 획득되면 안 됨, 얻음 %d건", len(got))
	}
}

func TestDigestBatchDueAndStableBatchID(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeDigest, `{}`)
	for i := 0; i < 3; i++ {
		addTestEvent(t, d, notify.EventFindingCreated, int64(6000+i), notify.Snapshot{Severity: "high"})
	}
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}

	// 배치가 막 생겨 나이 0이라, 30분 주기에서는 만기되면 안 된다.
	due, err := d.DigestBatchDue(ctx, ch.ID, 30*time.Minute)
	if err != nil {
		t.Fatalf("배치 만기 판정 실패: %v", err)
	}
	if due {
		t.Fatal("막 생긴 배치는 즉시 만기되면 안 됨")
	}

	// 세 전달의 생성 시간을 함께 과거로 밀어, 주기를 충분히 채운 배치를 시뮬레이션한다.
	if _, err := d.Exec(`UPDATE notification_deliveries SET created_at = now() - interval '40 minutes' WHERE channel_id=$1`, ch.ID); err != nil {
		t.Fatal(err)
	}
	due, err = d.DigestBatchDue(ctx, ch.ID, 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !due {
		t.Fatal("주기를 넘긴 배치는 만기로 판정되어야 함")
	}

	batch, err := d.ClaimDigestBatch(ctx, ch.ID, MaxDigestBatchSize, time.Minute)
	if err != nil {
		t.Fatalf("요약 배치 획득 실패: %v", err)
	}
	if len(batch) != 3 {
		t.Fatalf("요약은 한 번에 3건 전부 가져가야 함, 얻음 %d건", len(batch))
	}
	if batch[0].BatchID == nil {
		t.Fatal("요약 배치는 batch_id를 써야 함, 안 그러면 이력에서 함께 보낸 것을 알 수 없음")
	}
	firstBatchID := *batch[0].BatchID
	for _, dl := range batch {
		if dl.BatchID == nil || *dl.BatchID != firstBatchID {
			t.Fatalf("같은 배치는 batch_id를 공유해야 함, 얻음 %v vs %d", dl.BatchID, firstBatchID)
		}
	}

	// 이 배치를 **전체** 실패 재배치 후 다시 획득해도 batch_id는 원래 값을 유지해야 한다(COALESCE의 역할):
	// 안 그러면 한 번의 재시도로 '이 배치는 함께 보냈다'는 사실이 지워진다.
	//
	// 한 건만이 아니라 배치 전체를 재배치해야 한다 —— 전달 엔진이 요약 메시지를 보낼 때 그렇게 처리한다
	// (메시지 하나가 배치 전체를 대표하므로 성패를 함께한다). 한 건만 재배치하면 나머지는 아직 lease 기간 내라,
	// 재획득이 당연히 그 한 건만 가져온다.
	allIDs := make([]int64, 0, len(batch))
	for _, dl := range batch {
		allIDs = append(allIDs, dl.ID)
	}
	if err := d.RescheduleDeliveries(ctx, allIDs, time.Second, "模拟失败"); err != nil {
		t.Fatal(err)
	}
	// lease를 과거로 밀어 백오프 시간이 도달한 것을 시뮬레이션한다.
	if _, err := d.Exec(`UPDATE notification_deliveries SET next_attempt_at = now() - interval '1 minute' WHERE channel_id=$1`, ch.ID); err != nil {
		t.Fatal(err)
	}
	reclaimed, err := d.ClaimDigestBatch(ctx, ch.ID, MaxDigestBatchSize, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(reclaimed) != 3 {
		t.Fatalf("재획득은 3건 전부 가져와야 함, 얻음 %d", len(reclaimed))
	}
	if reclaimed[0].BatchID == nil || *reclaimed[0].BatchID != firstBatchID {
		t.Fatalf("재시도 후 batch_id는 원래 값 %d를 유지해야 함, 얻음 %v", firstBatchID, reclaimed[0].BatchID)
	}
}

func TestDeliveryStateTransitions(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	addTestEvent(t, d, notify.EventFindingCreated, 7001, notify.Snapshot{Severity: "high"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	got, err := d.ClaimRealtimeDeliveries(ctx, ch.ID, 10, time.Minute)
	if err != nil || len(got) != 1 {
		t.Fatalf("획득 실패: %v (%d)", err, len(got))
	}
	id := got[0].ID

	if err := d.RescheduleDeliveries(ctx, []int64{id}, time.Second, "网络抖动"); err != nil {
		t.Fatal(err)
	}
	var state, lastErr string
	if err := d.QueryRow(`SELECT state, last_error FROM notification_deliveries WHERE id=$1`, id).Scan(&state, &lastErr); err != nil {
		t.Fatal(err)
	}
	if state != NotifyStatePending || lastErr != "网络抖动" {
		t.Fatalf("재배치 후 pending이고 원인이 기록되어야 함, 얻음 state=%s err=%q", state, lastErr)
	}

	if err := d.FailDeliveries(ctx, []int64{id}, "重试耗尽"); err != nil {
		t.Fatal(err)
	}
	if err := d.QueryRow(`SELECT state FROM notification_deliveries WHERE id=$1`, id).Scan(&state); err != nil {
		t.Fatal(err)
	}
	if state != NotifyStateFailed {
		t.Fatalf("failed여야 함, 얻음 %s", state)
	}

	// 수동 재전송은 재시도 카운트를 0으로 하고 즉시 만기시켜야 한다, 안 그러면 구 실패 예산을 상속한다.
	if err := d.RetryNotificationDelivery(ctx, id); err != nil {
		t.Fatalf("재전송 실패: %v", err)
	}
	var attempts int
	var next time.Time
	if err := d.QueryRow(`SELECT state, attempts, next_attempt_at FROM notification_deliveries WHERE id=$1`, id).Scan(&state, &attempts, &next); err != nil {
		t.Fatal(err)
	}
	if state != NotifyStatePending || attempts != 0 {
		t.Fatalf("재전송 후 pending이고 attempts=0이어야 함, 얻음 state=%s attempts=%d", state, attempts)
	}
	if next.After(time.Now().Add(time.Second)) {
		t.Fatal("재전송은 즉시 획득 가능해야 함")
	}

	// 전달됨 상태의 전달은 재전송되면 안 된다.
	if err := d.MarkDeliveriesSent(ctx, []int64{id}); err != nil {
		t.Fatal(err)
	}
	if err := d.RetryNotificationDelivery(ctx, id); err == nil {
		t.Fatal("전달됨 상태의 전달은 재전송을 허용하면 안 됨")
	}
}

func TestListNotificationDeliveriesPagingAndFilter(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	for i := 0; i < 5; i++ {
		addTestEvent(t, d, notify.EventFindingCreated, int64(8000+i), notify.Snapshot{Severity: "high", Name: "分页测试"})
	}
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	if _, err := d.ClaimRealtimeDeliveries(ctx, ch.ID, 10, time.Minute); err != nil {
		t.Fatal(err)
	}

	page1, total, err := d.ListNotificationDeliveries(ctx, NotificationDeliveryFilter{ChannelID: ch.ID, State: NotifyStateSending}, 1, 2)
	if err != nil {
		t.Fatalf("조회 실패: %v", err)
	}
	if total != 5 {
		t.Fatalf("총수는 5여야 함, 얻음 %d", total)
	}
	if len(page1) != 2 {
		t.Fatalf("페이지당 2건, 얻음 %d", len(page1))
	}
	// 새 것이 앞: 첫 페이지 첫 건의 id가 두 번째 페이지 첫 건보다 커야 한다.
	page2, _, err := d.ListNotificationDeliveries(ctx, NotificationDeliveryFilter{ChannelID: ch.ID, State: NotifyStateSending}, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 || page2[0].ID >= page1[0].ID {
		t.Fatalf("페이지 순서는 새 것이 앞이어야 함, 얻음 page1[0]=%d page2[0]=%d", page1[0].ID, page2[0].ID)
	}
	// 렌더링 컨텍스트가 이력과 함께 반환되어야 한다, 안 그러면 목록이 '무엇을 푸시했는지' 표시할 수 없다.
	if page1[0].ChannelName == "" || page1[0].FindingID == 0 {
		t.Fatalf("이력 항목에 표시 필드가 없음: %+v", page1[0])
	}

	// 상태로 필터: pending이 없는 것.
	pending, totalPending, err := d.ListNotificationDeliveries(ctx, NotificationDeliveryFilter{ChannelID: ch.ID, State: NotifyStatePending}, 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	if totalPending != 0 || len(pending) != 0 {
		t.Fatalf("pending 전달이 없어야 함, 얻음 %d건 (total=%d)", len(pending), totalPending)
	}
}

func TestSetFindingStatusWithNotifyOnlyEmitsOnRealChange(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	tk, err := d.CreateTask("通知状态变更测试", "目标", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer d.DeleteTask(tk.ID)
	es := d.Exploration(tk.ExplorationID)
	f, err := es.RecordFinding(ctx, RecordFindingInput{
		TaskID: tk.ID, Worker: "test", VulnClass: "SQL注入", Name: "状态变更用例",
		Severity: "high", Summary: "摘要",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_events WHERE finding_id=$1`, f.FindingID) })

	// 저장 시 finding_created 이벤트가 하나 등록되므로, 먼저 그것을 세어 기준선으로 삼는다.
	var base int
	if err := d.QueryRow(`SELECT count(*) FROM notification_events WHERE finding_id=$1`, f.FindingID).Scan(&base); err != nil {
		t.Fatal(err)
	}
	if base < 1 {
		t.Fatal("취약점 저장은 같은 트랜잭션에서 푸시 이벤트 하나를 등록해야 함")
	}

	// 같은 상태로 변경: 이벤트가 생기면 안 됨(반복 제출로 푸시 노이즈를 내는 것 방지).
	from, found, notified, err := d.SetFindingStatusWithNotify(ctx, f.FindingID, "pending")
	if err != nil || !found {
		t.Fatalf("상태 설정 실패: found=%v err=%v", found, err)
	}
	if notified {
		t.Fatal("상태가 안 바뀌면 푸시 이벤트를 등록하면 안 됨")
	}
	if from != "pending" {
		t.Fatalf("변경 전 상태 pending을 반환해야 함, 얻음 %q", from)
	}

	// 실제 변경: 이벤트를 등록하고 from/to를 기록해야 한다.
	from, found, notified, err = d.SetFindingStatusWithNotify(ctx, f.FindingID, "fixed")
	if err != nil || !found {
		t.Fatalf("상태 설정 실패: found=%v err=%v", found, err)
	}
	if !notified {
		t.Fatal("상태가 실제로 바뀌면 푸시 이벤트를 등록해야 함")
	}
	if from != "pending" {
		t.Fatalf("from은 pending이어야 함, 얻음 %q", from)
	}
	var snapshot []byte
	if err := d.QueryRow(`SELECT snapshot FROM notification_events WHERE finding_id=$1 AND kind=$2`,
		f.FindingID, notify.EventFindingStatusChanged).Scan(&snapshot); err != nil {
		t.Fatalf("상태 변경 이벤트를 찾지 못함: %v", err)
	}
	var snap notify.Snapshot
	if err := json.Unmarshal(snapshot, &snap); err != nil {
		t.Fatal(err)
	}
	if snap.FromStatus != "pending" || snap.ToStatus != "fixed" {
		t.Fatalf("스냅샷의 상태 전이가 틀림: %s → %s", snap.FromStatus, snap.ToStatus)
	}
	// 스냅샷은 렌더링에 필요한 필드를 가져야 한다, 안 그러면 상태 변경 메시지가 빈 껍데기가 된다.
	if snap.VulnClass != "SQL注入" || snap.Severity != "high" || snap.Name != "状态变更用例" {
		t.Fatalf("스냅샷에 렌더링 필드가 없음: %+v", snap)
	}
	var status string
	if err := d.QueryRow(`SELECT status FROM findings WHERE id=$1`, f.FindingID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "fixed" {
		t.Fatalf("상태가 fixed로 갱신되어야 함, 얻음 %s", status)
	}

	// 존재하지 않는 취약점: found=false, 오류 없음.
	if _, found, _, err := d.SetFindingStatusWithNotify(ctx, 999999999, "fixed"); err != nil || found {
		t.Fatalf("존재하지 않는 취약점은 found=false이고 오류 없어야 함, 얻음 found=%v err=%v", found, err)
	}
}

func TestNotificationStatsSnapshot(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	addTestEvent(t, d, notify.EventFindingCreated, 9001, notify.Snapshot{Severity: "high"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	stats, err := d.NotificationStatsSnapshot(ctx)
	if err != nil {
		t.Fatalf("통계 실패: %v", err)
	}
	if stats.Channels < 1 || stats.ChannelsOn < 1 {
		t.Fatalf("채널 카운트가 틀림: %+v", stats)
	}
	if stats.Pending < 1 {
		t.Fatalf("발송 대기 전달이 집계되어야 함: %+v", stats)
	}
	// 막 생긴 전달의 적체 나이는 0에 가까워야 하며, 음수나 거대한 값이면 안 된다.
	if stats.BacklogAgeMS < 0 || stats.BacklogAgeMS > int64(time.Hour/time.Millisecond) {
		t.Fatalf("적체 나이가 유효하지 않음: %d ms", stats.BacklogAgeMS)
	}
	_ = ch
}

func TestNotificationChannelCRUDRoundTrip(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	ch := &NotificationChannel{
		Name:       "CRUD 往返",
		Kind:       notify.KindEmail,
		Mode:       NotifyModeDigest,
		Config:     json.RawMessage(`{"host":"smtp.example.com","port":587,"from":"a@b.c","to":["x@y.z"]}`),
		Filter:     json.RawMessage(`{"min_severity":"medium","on_status_change":true}`),
		RatePerMin: 42,
	}
	id, err := d.SaveNotificationChannel(ctx, ch)
	if err != nil {
		t.Fatalf("생성 실패: %v", err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_channels WHERE id=$1`, id) })

	got, err := d.NotificationChannelByID(ctx, id)
	if err != nil {
		t.Fatalf("읽기 실패: %v", err)
	}
	if got.Mode != NotifyModeDigest || got.RatePerMin != 42 || got.Name != "CRUD 往返" {
		t.Fatalf("왕복 필드 불일치: %+v", got)
	}
	if !got.IsEnabled() {
		t.Fatal("기본값은 활성이어야 함")
	}
	var cfg map[string]any
	if err := json.Unmarshal(got.Config, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg["host"] != "smtp.example.com" {
		t.Fatalf("설정이 올바르게 저장되지 않음: %v", cfg)
	}
	var filter notify.Filter
	if err := json.Unmarshal(got.Filter, &filter); err != nil {
		t.Fatal(err)
	}
	if filter.MinSeverity != "medium" || !filter.OnStatusChange {
		t.Fatalf("필터 조건이 올바르게 저장되지 않음: %+v", filter)
	}

	// 업데이트 후 다시 읽는다.
	got.Name = "改名了"
	off := false
	got.Enabled = &off
	if _, err := d.SaveNotificationChannel(ctx, got); err != nil {
		t.Fatal(err)
	}
	after, err := d.NotificationChannelByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if after.Name != "改名了" || after.IsEnabled() {
		t.Fatalf("업데이트가 적용되지 않음: %+v", after)
	}

	// 삭제 후 '존재하지 않음'을 보고해야 하며 조용히 성공하면 안 된다.
	if err := d.DeleteNotificationChannel(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := d.NotificationChannelByID(ctx, id); err != ErrNotificationChannelNotFound {
		t.Fatalf("ErrNotificationChannelNotFound 기대, 얻음 %v", err)
	}
	if err := d.DeleteNotificationChannel(ctx, id); err != ErrNotificationChannelNotFound {
		t.Fatalf("중복 삭제는 존재하지 않음을 보고해야 함, 얻음 %v", err)
	}
}

// TestSaveNotificationChannelKeepsExplicitZeroRate는 예전에 잘못 쓴 곳을 고정한다:
// **0은 유효한 설정이며 '레이트 리밋 없음'을 의미하므로, db 레이어가 '미지정'으로 보고 기본값으로 덮어써선 안 된다**.
//
// 과거 버그: SaveNotificationChannel에 `if RatePerMin <= 0 { 기본값 사용 }`이 있었고,
// 그래서 문서·UI 안내·takeTokens는 모두 '0=레이트 리밋 없음'으로 해석하는데, 저장 레이어에서만 몰래
// 20(DingTalk/WeCom/Telegram) 또는 100(Feishu)으로 바꿨다 —— 조작자는 리밋을 풀었다고 여기지만 실제로는 막히고,
// 아무 안내도 없었다. '미지정'과 '명시적 0'의 구분은 요청 본문만 표현할 수 있으므로,
// 기본값은 server 레이어가 채우고(notifyCreateChannel 참조), db 레이어는 저장만 한다.
func TestSaveNotificationChannelKeepsExplicitZeroRate(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	// 명시적 0(레이트 리밋 없음): 반드시 원래대로 저장해야 한다.
	unlimited := &NotificationChannel{
		Name: "不限流", Kind: notify.KindDingTalk, RatePerMin: 0,
		Config: json.RawMessage(`{"webhook":"https://example.com/h"}`),
	}
	id, err := d.SaveNotificationChannel(ctx, unlimited)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_channels WHERE id=$1`, id) })
	got, err := d.NotificationChannelByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.RatePerMin != 0 {
		t.Fatalf("명시적 0은 레이트 리밋 없음을 뜻하므로 원래대로 저장해야 함, 얻음 %d", got.RatePerMin)
	}
	if got.Mode != NotifyModeRealtime {
		t.Fatalf("기본 모드는 realtime이어야 함, 얻음 %s", got.Mode)
	}

	// 음수는 잘못된 입력이라, 조용히 다른 값으로 바꾸지 말고 거부해야 한다.
	bad := &NotificationChannel{
		Name: "负限流", Kind: notify.KindDingTalk, RatePerMin: -1,
		Config: json.RawMessage(`{"webhook":"https://example.com/h"}`),
	}
	if _, err := d.SaveNotificationChannel(ctx, bad); err == nil {
		t.Fatal("음수 레이트 리밋은 거부되어야 함")
	}
}

// TestDeleteChannelCascadesDeliveries는 외래 키 동작을 고정한다: 채널 삭제 후 그 전달 이력도 함께 사라지지만
// (설정이 없어지면 이력을 해석할 수 없음), 이벤트 자체는 남아야 한다 —— 다른 채널이 여전히 참조할 수 있다.
func TestDeleteChannelCascadesDeliveries(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, `{}`)
	ev := addTestEvent(t, d, notify.EventFindingCreated, 9101, notify.Snapshot{Severity: "high"})
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	var before int
	if err := d.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1`, ch.ID).Scan(&before); err != nil {
		t.Fatal(err)
	}
	if before == 0 {
		t.Fatal("전제 조건 불성립: 전달이 생기지 않음")
	}
	if err := d.DeleteNotificationChannel(ctx, ch.ID); err != nil {
		t.Fatal(err)
	}
	var after int
	if err := d.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1`, ch.ID).Scan(&after); err != nil {
		t.Fatal(err)
	}
	if after != 0 {
		t.Fatalf("채널 삭제 후 그 전달은 캐스케이드 삭제되어야 함, 여전히 %d건", after)
	}
	var evExists bool
	if err := d.QueryRow(`SELECT EXISTS(SELECT 1 FROM notification_events WHERE id=$1)`, ev).Scan(&evExists); err != nil {
		t.Fatal(err)
	}
	if !evExists {
		t.Fatal("채널 삭제가 이벤트 자체를 함께 삭제하면 안 됨")
	}
}

// TestClaimDigestBatchHonorsCallerLimit은 감사에서 지적한 구멍을 다룬다:
// 요약 채널은 이전에 토큰 버킷을 완전히 우회했다 —— allow가 takeTokens로 차감됐지만 아무도 안 썼고,
// rate_per_min이 digest 모드에 전혀 작용하지 않았다. 이제 limit도 제약에 참여한다.
func TestClaimDigestBatchHonorsCallerLimit(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()
	ch := newTestChannel(t, d, notify.KindDingTalk, NotifyModeDigest, `{}`)
	for i := 0; i < 10; i++ {
		addTestEvent(t, d, notify.EventFindingCreated, int64(7000+i), notify.Snapshot{Severity: "high"})
	}
	if _, _, err := d.FanOutPendingEvents(ctx, 100); err != nil {
		t.Fatal(err)
	}
	// limit=3: 3건만 획득 가능하고 나머지는 DB에 남는다.
	got, err := d.ClaimDigestBatch(ctx, ch.ID, 3, time.Minute)
	if err != nil {
		t.Fatalf("획득 실패: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("호출자 레이트 리밋 할당량에 따라 3건만 획득해야 함, 얻음 %d", len(got))
	}
	// limit=0은 이번 라운드 할당량 소진을 뜻한다: 한 건도 획득하면 안 되고, 오류도 내면 안 된다.
	if got, err := d.ClaimDigestBatch(ctx, ch.ID, 0, time.Minute); err != nil || len(got) != 0 {
		t.Fatalf("할당량 0일 때 0건 획득이고 오류 없어야 함, 얻음 %d건 err=%v", len(got), err)
	}
}

// TestFinishFindingRetestEmitsStatusChange는 감사에서 지적한 완전성 공백을 다룬다:
// 재검증 결론이 '수정 완료'일 때 상태는 실제로 바뀌지만, 그 UPDATE는 직접 DB에 쓰여
// 알림 포함 버전을 우회했다 —— 그래서 on_status_change를 설정한 채널은 이런 상태 전이에
// 푸시를 전혀 못 받고, 화면에서 상태가 조용히 바뀌어, 운영자는 플랫폼을 열어야 안다.
//
// 이 케이스는 '상태를 바꾸는 모든 경로가 상태 변경 이벤트를 등록해야 한다'를 고정한다.
func TestFinishFindingRetestEmitsStatusChange(t *testing.T) {
	d := notifyTestDB(t)
	ctx := context.Background()

	tk, err := d.CreateTask("复测推送测试", "目标", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer d.DeleteTask(tk.ID)
	es := d.Exploration(tk.ExplorationID)
	f, err := es.RecordFinding(ctx, RecordFindingInput{
		TaskID: tk.ID, Worker: "test", VulnClass: "SQL注入", Name: "复测目标",
		Severity: "high", Summary: "摘要",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Exec(`DELETE FROM notification_events WHERE finding_id=$1`, f.FindingID) })

	// 재검증 기록을 하나 만들어 완료 상태로 바로 밀어 넣는다.
	rt, _, _, err := d.CreateFindingRetest(ctx, f.FindingID, "复核")
	if err != nil {
		t.Fatal(err)
	}
	if rt.ConversationID == nil {
		t.Fatal("재검증은 세션 하나를 연결해야 함")
	}
	// 재검증은 먼저 running에 들어가야 결론을 저장할 수 있다(실제 흐름과 일치).
	if ok, err := d.StartFindingRetest(ctx, rt.ID); err != nil || !ok {
		t.Fatalf("재검증 시작 실패: ok=%v err=%v", ok, err)
	}
	if err := d.RecordFindingRetestResult(ctx, *rt.ConversationID, "fixed", "已修复", "证据"); err != nil {
		t.Fatal(err)
	}
	if err := d.FinishFindingRetest(rt.ID, "completed", ""); err != nil {
		t.Fatalf("재검증 종료 실패: %v", err)
	}

	var status string
	if err := d.QueryRow(`SELECT status FROM findings WHERE id=$1`, f.FindingID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != FindingFixed {
		t.Fatalf("재검증이 수정 완료로 판정된 후 상태는 fixed여야 함, 얻음 %s", status)
	}

	// 핵심 단언: 상태 변경 이벤트가 하나 있어야 하고 from/to가 올바라야 한다.
	var snapshot []byte
	err = d.QueryRow(`SELECT snapshot FROM notification_events WHERE finding_id=$1 AND kind=$2 ORDER BY id DESC LIMIT 1`,
		f.FindingID, notify.EventFindingStatusChanged).Scan(&snapshot)
	if err != nil {
		t.Fatalf("재검증이 수정 완료로 판정되면 상태 변경 푸시 이벤트를 등록해야 함(안 그러면 on_status_change 설정 채널이 못 받음): %v", err)
	}
	var snap notify.Snapshot
	if err := json.Unmarshal(snapshot, &snap); err != nil {
		t.Fatal(err)
	}
	if snap.FromStatus != "pending" || snap.ToStatus != FindingFixed {
		t.Fatalf("스냅샷의 상태 전이가 틀림: %s → %s", snap.FromStatus, snap.ToStatus)
	}
	// 스냅샷은 렌더링에 필요한 필드를 가져야 한다, 안 그러면 푸시가 빈 껍데기로 나간다.
	if snap.Name != "复测目标" || snap.Severity != "high" {
		t.Fatalf("스냅샷에 렌더링 필드가 없음: %+v", snap)
	}
}
