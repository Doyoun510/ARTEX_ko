package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/notify"
)

// 전역 설정 키(settings 키값 테이블에 저장, 테이블 생성 불필요).
const (
	// settingNotifyEnabled은 푸시 마스터 스위치. 기본 켜짐: 유지보수 기간에 원클릭으로 지혈하는 용도이고,
	// 기능 활성화 조건이 아니다 —— 진짜 활성화 조건은 '채널을 설정했는지'다.
	settingNotifyEnabled = "notify_enabled"
	// settingNotifyPublicBaseURL은 취약점 상세 회신 링크를 생성하는 외부 접근 주소
	// (예: https://artex.example.com). 비우면 메시지에 회신 링크 버튼이 없다.
	// 프로젝트에 재사용할 외부 주소 설정이 없어 여기에 하나 추가.
	settingNotifyPublicBaseURL = "notify_public_base_url"
	// settingNotifyDigestMinutes는 요약 모드의 주기(분).
	settingNotifyDigestMinutes = "notify_digest_interval_min"
)

const (
	// notifyTick은 전달 엔진의 폴링 간격. 3초는 이 엔진 실시간성의 상한이고,
	// '취약점 DB 기록'부터 '메시지가 IM에 도달'까지의 주요 지연 원인이다.
	notifyTick = 3 * time.Second
	// notifyLease는 전달을 수령할 때의 lease 기간. 단일 전달의 최악 소요 시간
	// (notify 패키지 HTTP 클라이언트 타임아웃 15초)보다 뚜렷이 커야 하며, 아니면 같은 행을 두
	// dispatcher가 동시에 전달하는 일이 생긴다.
	notifyLease = 3 * time.Minute
	// notifyFanOutPerTick은 라운드마다 디스패치하는 이벤트 수를 제한해, 채널을 처음 활성화할 때
	// 과거 적체를 한 번에 전부 전달 작업으로 펼치는 것을 방지.
	notifyFanOutPerTick = 200
	// notifyDefaultDigestMinutes는 요약 주기의 기본값.
	notifyDefaultDigestMinutes = 30
	// notifyUnlimitedBurstPerTick은 채널이 레이트 리밋 미설정일 때 라운드당 전달 상한.
	// '채널 하나를 무제한으로 설정 + 한 번에 수천 건 취약점 스캔'이
	// 단일 라운드 루프를 장시간 블록으로 끌고 가는 것을 방지하는 데 의의가 있다.
	notifyUnlimitedBurstPerTick = 50
	// notifyMaxSendsPerChannelPerTick은 단일 채널이 라운드당 최대 몇 건 전달하는지.
	//
	// 이 상한은 **lease 기간**에서 역산: 수령 시 행에 찍는 것이 lease(notifyLease = 3분)이고,
	// 한 라운드에서 직렬 전달하는 건수가 많아 최악 소요가 lease를 넘으면, 뒤 몇 건은 다 보내기 전에 lease가 만료된다.
	// 단일 프로세스 내에서는 무관(Run은 단일 goroutine 직렬 실행, tick은 재진입 안 함)하지만, **두
	// 프로세스가 같은 DB에 연결**되면, 상대가 lease 만료된 행을 다시 수령해 중복 전송하고,
	// attempts를 이중으로 증가시키며, 원 프로세스가 아직 전달 중일 때 실패로 판정한다.
	//
	// 값: 3분 lease / 30초 단일 타임아웃 = 6은 **lease를 딱 꽉 채움**, 여유 0이라,
	// 쓸 수 없다; 5를 쓰면 최악 소요 150초로 30초 여유를 남긴다. 이 관계는
	// TestNotifyTickBudgetFitsWithinLease가 고정한다 —— notifyLease,
	// notifySendTimeout 또는 이 값 중 어느 하나라도 바꾸면 그 단언이 실패한다.
	notifyMaxSendsPerChannelPerTick = 5
	// notifySendTimeout은 단일 전달의 타임아웃. 위 상수의 값도 함께 결정하며,
	// 둘을 곱한 값이 notifyLease를 넘으면 안 된다, TestNotifyTickBudgetFitsWithinLease 참고.
	notifySendTimeout = 30 * time.Second
)

// notifyBackoff은 실패 재시도의 백오프 시퀀스, 인덱스는 이미 시도한 횟수.
// 3번 기회(첫 시도 포함)는 db.MaxNotifyAttempts와 대응, 둘은 반드시 함께 바꾼다.
var notifyBackoff = []time.Duration{
	time.Second,
	5 * time.Second,
	30 * time.Second,
}

// Notifier는 취약점 푸시의 전달 엔진.
//
// Scheduler와 병렬로, 독립 goroutine으로 실행(server.New 참고). 의도적으로 재사용하지 않는
// Scheduler의 tick: 푸시의 실시간성 요구(3초)가 트리거의 업무 리듬과 다르고,
// 둘의 실패는 서로 영향 없음 —— 푸시가 막혀도 agent 트리거에 영향을 줘선 안 된다.
type Notifier struct {
	s  *Server
	pg *db.DB

	// mu는 buckets를 보호. 채널 수가 적고 경합이 낮아 뮤텍스 하나면 충분하며,
	// 더 세분화된 구조를 도입할 가치가 없다.
	mu      sync.Mutex
	buckets map[int64]*notifyBucket
}

// notifyBucket은 단일 채널의 토큰 버킷.
//
// '분당 카운트 후 리셋'하는 슬라이딩 윈도우 대신 토큰 버킷을 쓰는 건 후자의 경계 효과가 나쁘기 때문:
// 윈도우 끝에서 20건을 꽉 채워 보내고 다음 순간 또 20건을 보내면, 플랫폼 입장에서는 1초에 40건이라,
// 레이트 리밋에 걸린다; 토큰 버킷은 일정 속도로 보충해 이런 버스트를 자연히 피한다.
type notifyBucket struct {
	tokens   float64
	lastFill time.Time
}

func newNotifier(s *Server) *Notifier {
	return &Notifier{s: s, pg: s.m.pg, buckets: map[int64]*notifyBucket{}}
}

// Run은 ctx가 끝날 때까지 루프. server.New가 한 번 시작.
func (n *Notifier) Run(ctx context.Context) {
	if n.pg == nil {
		return
	}
	t := time.NewTicker(notifyTick)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n.step(ctx)
		}
	}
}

// step은 한 라운드를 돈다: 먼저 새 이벤트를 디스패치하고, 다음에 만기된 작업을 전달.
//
// 어느 단계가 실패해도 로그만 남기고 루프를 중단하지 않음 —— 알림 시스템의 장애가 프로세스 레벨 문제로 번져선 절대 안 된다.
// 각 tick은 독립적이라, 다음 라운드가 자연히 재시도한다.
func (n *Notifier) step(ctx context.Context) {
	if !n.enabled() {
		return
	}
	if _, _, err := n.pg.FanOutPendingEvents(ctx, notifyFanOutPerTick); err != nil {
		log.Printf("[notify] 이벤트 디스패치 실패: %v", err)
		return
	}
	channels, err := n.pg.ListNotificationChannels(ctx)
	if err != nil {
		log.Printf("[notify] 채널 읽기 실패: %v", err)
		return
	}
	baseURL := n.publicBaseURL()
	for _, ch := range channels {
		if !ch.IsEnabled() {
			continue
		}
		// 토큰 버킷의 계량 단위는 **메시지 건수**(HTTP 요청 수와 동등)이지 취약점 건수가 아니다.
		// 실시간 모드에서는 둘이 같고(취약점 하나당 메시지 하나); 요약 모드에서는 한 배치의 취약점을 합쳐
		// 메시지 하나로 만들어 토큰 하나만 소비한다.
		//
		// 두 모드 모두 먼저 토큰 버킷에 묻고, 다음에 할당량만큼 수령 —— 순서를 뒤집으면 안 되며, 아니면 레이트 리밋에 막힌
		// 전달이 이미 재시도 횟수를 소비한 뒤가 된다.
		now := time.Now()
		if ch.Mode == db.NotifyModeDigest {
			tokens, claimLimit := digestTickPlan()
			if n.takeTokens(ch.ID, ch.RatePerMin, tokens, now) <= 0 {
				continue
			}
			n.stepDigest(ctx, ch, claimLimit, baseURL)
			continue
		}
		allow := n.takeTokens(ch.ID, ch.RatePerMin, notifyMaxSendsPerChannelPerTick, now)
		if allow <= 0 {
			continue
		}
		n.stepRealtime(ctx, ch, allow, baseURL)
	}
}

// digestTickPlan은 요약 채널의 이번 라운드 토큰 소비와 배치 크기 상한을 반환.
//
// 두 반환값은 **서로 다른 단위**이며, 이것이 함수로 분리한 이유다:
//
//   - tokens는 메시지 건수. 한 배치의 취약점을 메시지 하나로 합쳐 HTTP 요청 한 번을 보내므로 항상 1.
//     따라서 rate_per_min은 digest에도 여전히 적용된다(분당 최대 이만큼의 요약 메시지).
//   - claimLimit은 이 배치가 최대 몇 건의 취약점을 담는지. 메모리 상한에만 제약되고 요청 예산과는 무관.
//
// 예전에 rate_per_min을 digest에 적용하려고, 라운드당 요청 예산
// (notifyMaxSendsPerChannelPerTick, lease에서 역산)을 그대로 배치 크기로 넘겼다.
// 그 결과 rate_per_min=20 채널이 3초 tick에서 토큰을 1개만 보충받아, 요약
// 메시지가 취약점 1건만 담아 —— digest가 '요약 문구가 붙은 실시간 푸시'로 퇴화하고, 독자는
// '최근 30분 취약점 1건 추가'의 연속을 받으며, db.MaxDigestBatchSize는 영원히 도달 불가.
//
// 이 증상은 엔드투엔드 테스트에서 발견하기 어렵다(기존 케이스는 모두 충분히 큰 limit을 수동으로
// stepDigest에 넘겨 step의 할당량 계산을 우회했다), 그래서 결정을 여기로 모아
// TestDigestTickPlanDecouplesBatchSizeFromSendBudget가 직접 고정한다.
func digestTickPlan() (tokens, claimLimit int) {
	return 1, db.MaxDigestBatchSize
}

// stepRealtime은 어떤 채널의 실시간 작업을 수령해 전달, 취약점 하나당 메시지 하나.
func (n *Notifier) stepRealtime(ctx context.Context, ch *db.NotificationChannel, allow int, baseURL string) {
	deliveries, err := n.pg.ClaimRealtimeDeliveries(ctx, ch.ID, allow, notifyLease)
	if err != nil {
		log.Printf("[notify] 실시간 전달 수령 실패 channel=%d: %v", ch.ID, err)
		return
	}
	if len(deliveries) == 0 {
		return
	}
	channel, cfg, ok := n.adapt(ch)
	if !ok {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), fmt.Sprintf("채널 유형 %q 미등록", ch.Kind))
		return
	}
	for _, dl := range deliveries {
		msg, err := n.renderSingle(ctx, dl, baseURL)
		if err != nil {
			// 렌더링 실패는 로컬 데이터 문제라, 재시도해도 나아지지 않는다.
			_ = n.pg.FailDeliveries(ctx, []int64{dl.ID}, err.Error())
			continue
		}
		n.send(ctx, channel, cfg, msg, []*db.NotificationDelivery{dl})
	}
}

// stepDigest는 배치 만기 시 어떤 채널의 발송 대기 전달을 메시지 하나로 모아 보낸다.
func (n *Notifier) stepDigest(ctx context.Context, ch *db.NotificationChannel, allow int, baseURL string) {
	window := n.digestInterval()
	due, err := n.pg.DigestBatchDue(ctx, ch.ID, window)
	if err != nil {
		log.Printf("[notify] 요약 배치 판단 실패 channel=%d: %v", ch.ID, err)
		return
	}
	if !due {
		return
	}
	deliveries, err := n.pg.ClaimDigestBatch(ctx, ch.ID, allow, notifyLease)
	if err != nil {
		log.Printf("[notify] 요약 배치 수령 실패 channel=%d: %v", ch.ID, err)
		return
	}
	if len(deliveries) == 0 {
		return
	}
	channel, cfg, ok := n.adapt(ch)
	if !ok {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), fmt.Sprintf("채널 유형 %q 미등록", ch.Kind))
		return
	}
	msg, included, err := n.renderBatch(ctx, deliveries, baseURL, int(window.Minutes()))
	if err != nil {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), err.Error())
		return
	}
	// 스냅샷이 깨져 메시지에 못 들어간 전달은 명시적으로 실패 판정해야 한다. 그러지 않으면 그것들은
	// included 밖에 남아 메시지에도 실패 목록에도 들어가지 않고 —— 전송 성공 시 그 상태가
	// 이후 일괄 표시에서 누락되어, lease가 만료되어 반복 수령될 때까지 영원히 sending에 머문다.
	if skipped := excludeDeliveries(deliveries, included); len(skipped) > 0 {
		reason := "이벤트 스냅샷을 파싱할 수 없어 이 취약점을 메시지로 렌더링할 수 없음"
		if fErr := n.pg.FailDeliveries(ctx, deliveryIDs(skipped), reason); fErr != nil {
			log.Printf("[notify] 불량 스냅샷 전달 실패 표시 channel=%s ids=%v: %v", ch.Kind, deliveryIDs(skipped), fErr)
		}
		log.Printf("[notify] 스냅샷 파싱 불가 전달 %d건 건너뜀 channel=%d", len(skipped), ch.ID)
	}
	// 메시지에 들어간 것만 send에 넘긴다: included[i]와 msg.Items[i]가 엄격히 대응하고,
	// send는 이 대응 관계로 '채널이 앞 K건을 담았다고 보고'한 것을 올바른 전달 행에 반영한다.
	n.send(ctx, channel, cfg, msg, included)
}

// send는 전달하고 결과에 따라 상태를 전이한다.
//
// 같은 배치의 전달(요약 모드에서는 수십 건일 수 있음)은 하나의 전송 결과를 공유한다: 송달되거나, 배치 전체 재시도.
// 건별 재시도는 하지 않음 —— 요약 메시지는 하나라, 그 일부만 재전송하면 배치 의미가 뒤틀린다.
//
// 유일한 예외는 **채널 길이 상한으로 인한 분할**: 채널이 실제로 앞 K건만 담았다고 보고하면,
// K+1건부터는 다음 배치로 남겨야 하며, 함께 성공으로 표시해선 안 된다. 아니면 잘려 나간
// 그 취약점들은 메시지에도 실패 목록에도 없이 완전히 사라진다.
func (n *Notifier) send(ctx context.Context, channel notify.Channel, cfg map[string]any, msg notify.Message, deliveries []*db.NotificationDelivery) {
	// 단일 전달에 상한을 두어, 어떤 채널이 막혀 이 라운드의 나머지 채널 전부를 끌고 가는 것을 방지.
	sendCtx, cancel := context.WithTimeout(ctx, notifySendTimeout)
	defer cancel()
	delivered, err := channel.Send(sendCtx, cfg, msg)
	if err == nil && delivered > 0 {
		if delivered > len(deliveries) {
			// 채널이 보고한 건수가 전달 수를 넘을 수 없다; 실제로 발생하면 렌더링 레이어가 잘못 계산한 것이라,
			// 전부 송달로 처리하고 문제를 기록하는 편이, 레코드를 어지럽히는 것보다 낫다.
			log.Printf("[notify] 채널 보고 송달 건수 %d가 전달 수 %d 초과 channel=%s, 전부 송달로 처리",
				delivered, len(deliveries), channel.Kind())
			delivered = len(deliveries)
		}
		sent, rest := deliveries[:delivered], deliveries[delivered:]
		if err := n.pg.MarkDeliveriesSent(ctx, deliveryIDs(sent)); err != nil {
			log.Printf("[notify] 송달 표시 실패 channel=%s ids=%v: %v", channel.Kind(), deliveryIDs(sent), err)
		}
		if len(rest) > 0 {
			// 이 메시지가 채널 길이 상한에 도달: 나머지는 즉시 큐로 복귀, 다음 tick이 이어 발송.
			// RescheduleDeliveries가 아니라 DeferDeliveries 사용 —— 이것은 실패가 아니라,
			// 재시도 예산을 소비해선 안 된다(수령 시 이미 낙관적으로 +1 했고, 거기서 되돌린다).
			if err := n.pg.DeferDeliveries(ctx, deliveryIDs(rest),
				fmt.Sprintf("이 메시지가 채널 길이 상한에 도달, 앞 %d건만 송달, 나머지는 다음 배치로", delivered)); err != nil {
				log.Printf("[notify] 분할 이어 발송 큐 등록 실패 channel=%s ids=%v: %v", channel.Kind(), deliveryIDs(rest), err)
			}
		}
		return
	}
	if err == nil {
		// 채널이 오류도 안 내고 몇 건 송달했는지도 말하지 않음. 실패로 처리(백오프),
		// 이 전달이 반복 수령되면서도 영원히 표시되지 않는 것을 방지.
		err = fmt.Errorf("채널이 송달 건수를 보고하지 않음(delivered=%d)", delivered)
	}

	// 실패 처리는 **건별**로 결정하며, 배치 전체의 최대 시도 횟수로 판단하지 않는다.
	//
	// 예전에는 `if maxAttempts(deliveries) >= MaxNotifyAttempts`로 배치 전체를 판정했지만, 배치 안
	// 각 건의 시도 횟수는 같지 않다: 이미 두 번 재시도한 오래된 전달(attempts=2)이 같은 배치의
	// 완전히 새 전달(attempts=1)을 함께 failed로 끌고 가 —— 새 취약점이 재시도 한 번 못 쓰고 영구 손실되며,
	// '오래된 행이 새 행을 끌어들이지 않게' 한다는 본래 취지와 정반대가 된다.
	permanent := notify.IsPermanent(err)
	var failIDs, exhaustedIDs []int64
	byDelay := map[time.Duration][]int64{}
	for _, dl := range deliveries {
		switch {
		case permanent:
			failIDs = append(failIDs, dl.ID)
		case dl.Attempts >= db.MaxNotifyAttempts:
			exhaustedIDs = append(exhaustedIDs, dl.ID)
		default:
			delay := notifyBackoff[min(dl.Attempts, len(notifyBackoff)-1)]
			byDelay[delay] = append(byDelay[delay], dl.ID)
		}
	}

	if len(failIDs) > 0 {
		if fErr := n.pg.FailDeliveries(ctx, failIDs, err.Error()); fErr != nil {
			log.Printf("[notify] 실패 상태 표시 오류 channel=%s ids=%v: %v", channel.Kind(), failIDs, fErr)
		}
	}
	if len(exhaustedIDs) > 0 {
		reason := fmt.Sprintf("%d번 재시도 후에도 실패: %s", db.MaxNotifyAttempts, err)
		if fErr := n.pg.FailDeliveries(ctx, exhaustedIDs, reason); fErr != nil {
			log.Printf("[notify] 실패 상태 표시 오류 channel=%s ids=%v: %v", channel.Kind(), exhaustedIDs, fErr)
		}
	}
	// 지연별로 그룹화해 재정렬: 백오프가 3단계뿐이라 그룹 수가 자연히 작아, 건별로 한 번씩
	// UPDATE를 보낼 필요 없다(그러면 500건 배치가 500번 왕복을 낸다).
	for delay, group := range byDelay {
		if rErr := n.pg.RescheduleDeliveries(ctx, group, delay, err.Error()); rErr != nil {
			log.Printf("[notify] 전달 재정렬 실패 channel=%s ids=%v: %v", channel.Kind(), group, rErr)
		}
	}
	if len(failIDs)+len(exhaustedIDs) > 0 {
		log.Printf("[notify] 전달 실패 channel=%d kind=%s 영구 실패=%d 재시도 소진=%d 재시도 대기=%d: %s",
			deliveries[0].ChannelID, channel.Kind(), len(failIDs), len(exhaustedIDs), len(byDelay), err)
	}
}

// excludeDeliveries는 all 중 keep에 없는 것들을 반환(포인터 신원으로 비교).
// '메시지에 못 들어간' 전달을 찾는 데 사용 —— 그것들은 명시적으로 처리해야 하며, 회색 지대에 남겨선 안 된다.
func excludeDeliveries(all, keep []*db.NotificationDelivery) []*db.NotificationDelivery {
	inKeep := make(map[*db.NotificationDelivery]bool, len(keep))
	for _, dl := range keep {
		inKeep[dl] = true
	}
	var out []*db.NotificationDelivery
	for _, dl := range all {
		if !inKeep[dl] {
			out = append(out, dl)
		}
	}
	return out
}

// adapt은 채널 구현을 가져와 그 설정을 파싱.
// ok=false 반환은 유형 미등록을 뜻하며, 전달은 무한 재시도가 아니라 바로 실패 판정해야 한다.
func (n *Notifier) adapt(ch *db.NotificationChannel) (notify.Channel, map[string]any, bool) {
	channel, ok := notify.Get(ch.Kind)
	if !ok {
		return nil, nil, false
	}
	var cfg map[string]any
	if len(ch.Config) > 0 {
		// 설정 파싱 실패 시 빈 map 제공: 채널 자체의 Validate가 '어떤 필드가 빠졌는지' 알려주며,
		// 그 오류가 JSON 파싱 오류보다 사용자 수정을 더 잘 안내한다.
		_ = json.Unmarshal(ch.Config, &cfg)
	}
	if cfg == nil {
		cfg = map[string]any{}
	}
	return channel, cfg, true
}

// renderSingle은 단일 취약점 메시지를 렌더링.
func (n *Notifier) renderSingle(ctx context.Context, dl *db.NotificationDelivery, baseURL string) (notify.Message, error) {
	snap, err := parseSnapshot(dl)
	if err != nil {
		return notify.Message{}, err
	}
	item, err := n.itemFor(ctx, snap, baseURL)
	if err != nil {
		return notify.Message{}, err
	}
	return notify.Message{Items: []notify.Item{item}, HomeURL: baseURL}, nil
}

// renderBatch는 요약 메시지를 렌더링. 스냅샷을 건별로 파싱 —— 한 건이 깨지면 그 한 건만 건너뛰고,
// 그것이 배치 요약 전체를 날리지 않게 한다.
//
// 반환값 included는 msg.Items와 **엄격히 일대일 대응**(i번째 전달 ↔ i번째 항목).
// 이 대응 관계는 강한 요구사항: 호출자가 '채널이 앞 K건을 담았다고 보고'에 따라 앞 K개 전달을
// 송달로 표시한다. 여기서 불량 스냅샷을 건너뛰고도 건너뛴 전달을 included에서 빼지 않으면,
// 인덱스가 어긋나 —— 실패했어야 할 불량 항목이 송달로 표시되고, 정상 항목이 미송달로 오판된다.
// 깨진 것들은 호출자가 명시적으로 실패 표시, stepDigest 참고.
func (n *Notifier) renderBatch(ctx context.Context, deliveries []*db.NotificationDelivery, baseURL string, windowMinutes int) (notify.Message, []*db.NotificationDelivery, error) {
	items := make([]notify.Item, 0, len(deliveries))
	included := make([]*db.NotificationDelivery, 0, len(deliveries))
	for _, dl := range deliveries {
		snap, err := parseSnapshot(dl)
		if err != nil {
			// 불량 스냅샷은 메시지에도 included에도 들어가지 않음 —— 그 처리는 호출자가 담당
			// (명시적으로 실패 표시, '송달됨'에 섞여 얼버무리지 않음).
			log.Printf("[notify] 요약 배치에서 파싱 불가 스냅샷 건너뜀 delivery=%d: %v", dl.ID, err)
			continue
		}
		item, err := n.itemFor(ctx, snap, baseURL)
		if err != nil {
			return notify.Message{}, nil, err
		}
		items = append(items, item)
		included = append(included, dl)
	}
	if len(items) == 0 {
		return notify.Message{}, nil, fmt.Errorf("요약 배치 %d건 전달 전부 파싱 불가", len(deliveries))
	}
	return notify.Message{
		Items:         items,
		Batch:         true,
		WindowMinutes: windowMinutes,
		HomeURL:       baseURL,
	}, included, nil
}

// itemFor는 이벤트 스냅샷을 푸시 대기 항목으로 렌더링하며, 자산 이름과 상세 회신 링크도 파싱.
func (n *Notifier) itemFor(ctx context.Context, snap notify.Snapshot, baseURL string) (notify.Item, error) {
	assets, err := n.pg.NotificationAssetNames(ctx, snap.AssetIDs)
	if err != nil {
		// 자산 이름 파싱 실패가 푸시를 막아선 안 된다: 이름을 못 읽는 것이 알림을 못 받는 것보다 훨씬 가벼우며,
		// 메시지에서 자산 한 줄이 빠질 뿐이다.
		log.Printf("[notify] 자산 이름 파싱 실패 finding=%d: %v", snap.FindingID, err)
	}
	item := notify.Item{
		FindingID:  snap.FindingID,
		Name:       snap.Name,
		VulnClass:  snap.VulnClass,
		Severity:   snap.Severity,
		Summary:    snap.Summary,
		Assets:     assets,
		FromStatus: snap.FromStatus,
		ToStatus:   snap.ToStatus,
	}
	if baseURL != "" {
		// 상세 페이지 라우트는 web/src/app/(main)/function/findings/detail/page.tsx 참고,
		// query 파라미터 id에서 취약점 id를 읽는다.
		item.DetailURL = fmt.Sprintf("%s/function/findings/detail?id=%d", baseURL, snap.FindingID)
	}
	return item, nil
}

// takeTokens는 채널 토큰 버킷에서 **최대 want 개** 토큰을 가져오며, 실제로 가져온 수를 반환.
//
// 토큰 하나 = 메시지 하나(HTTP 요청 한 번). 실시간 모드에서는 호출자가 필요한 수만큼 넘기고;
// 요약 모드에서는 한 배치의 취약점을 메시지 하나로만 보내므로 1을 넘긴다.
//
// 버킷 용량은 그 채널의 분당 상한이고, 일정 속도로 보충. ratePerMin<=0은 레이트 리밋 없음을 뜻하며,
// 유한하지만 충분히 큰 값을 반환해, 단일 라운드 루프가 무한 적체에 끌려가는 것을 방지.
//
// want 이 상한은 필수: 없으면 버킷을 통째로 비울 수밖에 없는데, 호출자 자신도 라운드당 상한이 있어,
// 많이 가져온 토큰은 쓰이지도 않고 다음 보충 전에 허공으로 사라져 —— 쌓아 둔 버스트 용량은 영원히 도달 불가하고,
// '이번 라운드에 발송 대기 전달이 전혀 없음'조차 한 번 차감된다.
func (n *Notifier) takeTokens(channelID int64, ratePerMin, want int, now time.Time) int {
	if want <= 0 {
		return 0
	}
	if ratePerMin <= 0 {
		return min(want, notifyUnlimitedBurstPerTick)
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	b := n.buckets[channelID]
	if b == nil {
		b = &notifyBucket{tokens: float64(ratePerMin), lastFill: now}
		n.buckets[channelID] = b
	}
	// 경과한 실제 시간에 따라 보충, 속도는 ratePerMin/60 per second.
	if elapsed := now.Sub(b.lastFill).Seconds(); elapsed > 0 {
		b.tokens = minF(float64(ratePerMin), b.tokens+elapsed*float64(ratePerMin)/60)
		b.lastFill = now
	}
	// 아주 작은 epsilon을 더한 뒤 정수화: 토큰 수는 부동소수점 누적이라, 두 번에 나눠 채울 때
	// 0.5 + 0.5가 0.9999999999가 될 수 있고, 바로 int()하면 0으로 잘려 ——
	// 수학적으로는 가득 찬 버킷인데도 토큰을 못 꺼낸다. 1e-9은 토큰 하나보다 훨씬 작아, 진짜 부족분을 놓치지 않는다.
	take := min(int(b.tokens+1e-9), want)
	if take <= 0 {
		return 0
	}
	b.tokens -= float64(take)
	return take
}

// enabled는 마스터 스위치를 읽는다.
func (n *Notifier) enabled() bool {
	return n.pg.GetBool(settingNotifyEnabled, true)
}

// publicBaseURL은 회신 링크용 외부 주소를 반환하며, 끝 슬래시를 제거.
func (n *Notifier) publicBaseURL() string {
	v, ok, err := n.pg.GetSetting(settingNotifyPublicBaseURL)
	if err != nil || !ok {
		return ""
	}
	return trimTrailingSlash(v)
}

// digestInterval은 요약 주기를 반환하며, 잘못됐거나 미설정이면 기본값으로 폴백.
func (n *Notifier) digestInterval() time.Duration {
	v, ok, err := n.pg.GetSetting(settingNotifyDigestMinutes)
	if err != nil || !ok {
		return time.Duration(notifyDefaultDigestMinutes) * time.Minute
	}
	m := 0
	if _, err := fmt.Sscanf(v, "%d", &m); err != nil || m <= 0 {
		return time.Duration(notifyDefaultDigestMinutes) * time.Minute
	}
	return time.Duration(m) * time.Minute
}

// parseSnapshot은 전달에 대응하는 이벤트의 스냅샷을 파싱.
func parseSnapshot(dl *db.NotificationDelivery) (notify.Snapshot, error) {
	var snap notify.Snapshot
	if len(dl.Snapshot) == 0 {
		return snap, fmt.Errorf("전달 %d의 이벤트 스냅샷이 비어 있음", dl.ID)
	}
	if err := json.Unmarshal(dl.Snapshot, &snap); err != nil {
		return snap, fmt.Errorf("전달 %d의 이벤트 스냅샷 파싱 실패: %w", dl.ID, err)
	}
	if snap.Kind == "" {
		// 이벤트 유형은 이벤트 행을 기준으로 함, 스냅샷 안의 것은 구 버전이 썼을 수 있다.
		snap.Kind = dl.EventKind
	}
	return snap, nil
}

func deliveryIDs(deliveries []*db.NotificationDelivery) []int64 {
	out := make([]int64, 0, len(deliveries))
	for _, dl := range deliveries {
		out = append(out, dl.ID)
	}
	return out
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
