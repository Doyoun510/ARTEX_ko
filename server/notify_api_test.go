package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/notify"
)

// 이 파일은 알림 전송 기능의 엔드투엔드 동작을 다룹니다: 취약점 등록 → 이벤트 → 분배 → 실제 HTTP 전송.
//
// 안전 관련 주의 사항: 이 사례들은 **전역 Notifier.step()을 호출하지 않으며**, 직접 생성한
// 채널에만 stepRealtime/stepDigest를 호출합니다. step()이 DB의 모든 활성 채널을 순회하므로,
// 실제 DingTalk/WeCom(기업용 위챗) 봇이 설정된 개발 DB에서 테스트하면 전역 step이 테스트 중
// 생성한 취약점을 해당 그룹에 실제로 전송합니다. 채널별 호출로 영향 범위를 테스트가 만든 가짜 수신 측에만 엄격히 제한합니다.
//
// 정리: 사례가 끝나면 이 사례에서 생성한 이벤트(전송 항목 연쇄 삭제)와 채널을 삭제해 실제 채널에 전송 적체를 남기지 않습니다.
//
// 검증 기준: stepRealtime/stepDigest는 반환값 없이 내부 로그를 기록하므로, 여기서는 함수의
// 반환값 대신 **관측 가능한 외부 동작**(가짜 수신 측이 받은 내용, 전송 행의 최종 상태)을
// 검증합니다. 반환값을 stub 처리하는 것보다 실제 호출 경로에 더 가깝습니다.

// notifyFixture는 이 파일의 사례에서 공유하는 fixture입니다.
type notifyFixture struct {
	s       *Server
	pg      *db.DB
	request func(method, path, body string) *httptest.ResponseRecorder
	n       *Notifier
	// 직접 생성한 task/exploration: 사례의 취약점을 여기에 기록해 다른 사례의 데이터와 격리합니다.
	taskID int64
	expID  int64
	// cleanupMark 이후에 생성한 이벤트는 정리 시 함께 삭제합니다.
	cleanupMark int64
}

func newNotifyFixture(t *testing.T) *notifyFixture {
	t.Helper()
	// 이 파일의 가짜 수신 측은 모두 127.0.0.1에서 실행되며, 전송은 기본적으로 루프백 주소를 거부합니다.
	// (SSRF로 같은 시스템의 서비스와 클라우드 메타데이터에 접근하는 것을 방지합니다.) 테스트에서는 이 스위치를 명시적으로 켭니다.
	// 연결 보호 검사의 '기본 거부' 동작은 notify 패키지의 ssrf_test.go가 다룹니다.
	t.Setenv(notify.AllowLocalTargetsEnv, "1")
	s, _, request := trafficEvidenceServer(t)
	pg := s.m.pg

	// task를 직접 생성합니다: 공유 fixture인 trafficEvidenceServer가 만든 task에서는 exploration id를 얻을 수 없지만,
	// 취약점 기록에는 이 값이 반드시 필요합니다.
	task, err := s.m.CreateTask("通知推送测试", "推送行为验证", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pg.Exec(`DELETE FROM tasks WHERE id=$1`, taskID) })

	var mark int64
	if err := pg.QueryRow(`SELECT COALESCE(max(id),0) FROM notification_events`).Scan(&mark); err != nil {
		t.Fatal(err)
	}
	// fixture가 자체적으로 완결되도록, fixture 생성 전에 존재하던 이벤트를 한 번에 분배 완료로 표시합니다.
	//
	// 필요한 이유: FanOutPendingEvents는 **전역**으로 동작하며 DB의 미분배 이벤트를 모두
	// 매칭되는 모든 채널에 분배합니다. 공유 fixture인 trafficEvidenceServer도 취약점 하나를 기록하고
	// (바로 반환하는 최초 finding), 다른 사례의 잔여 항목도 있을 수 있습니다. 격리하지 않으면
	// 이 불필요한 이벤트가 현재 사례의 채널에 분배되어 '전송 항목 N개가 있어야 함' 같은 검사가
	// 간헐적으로 실패합니다. 실패 여부가 사례 실행 순서에 따라 달라져 즉시 실패하는 경우보다 찾기 어렵습니다.
	if _, err := pg.Exec(`UPDATE notification_events SET fanned_out = true WHERE id <= $1 AND NOT fanned_out`, mark); err != nil {
		t.Fatal(err)
	}

	f := &notifyFixture{s: s, pg: pg, request: request, n: newNotifier(s), taskID: taskID, expID: task.ExpID, cleanupMark: mark}
	t.Cleanup(func() {
		if _, err := pg.Exec(`DELETE FROM notification_events WHERE id > $1`, f.cleanupMark); err != nil {
			t.Logf("알림 이벤트 정리 실패: %v", err)
		}
	})
	// 전체 스위치는 반드시 켜져 있어야 합니다(다른 사례가 껐을 수 있습니다).
	if err := pg.SetBool(settingNotifyEnabled, true); err != nil {
		t.Fatal(err)
	}
	return f
}

// record는 실제 증거 쓰기 경로로 취약점 하나를 등록하고 finding id를 반환합니다.
// 이 경로는 **같은 트랜잭션**에서 알림 전송 이벤트를 등록하며, 이 기능이 연결되는 지점입니다.
func (f *notifyFixture) record(t *testing.T, vulnclass, severity string) int64 {
	t.Helper()
	out, err := f.s.evidenceStore().Record(context.Background(), db.RecordFindingInput{
		TaskID:        f.taskID,
		ExplorationID: f.expID,
		Worker:        "test",
		VulnClass:     vulnclass,
		Name:          vulnclass,
		Severity:      severity,
		Summary:       vulnclass + " 的摘要",
		Evidence:      "poc",
	}, nil)
	if err != nil {
		t.Fatalf("취약점 기록 실패: %v", err)
	}
	return out.FindingID
}

// channel은 채널 설정을 읽어 옵니다(채널별 stepX 호출에 사용).
func (f *notifyFixture) channel(t *testing.T, id int64) *db.NotificationChannel {
	t.Helper()
	ch, err := f.pg.NotificationChannelByID(context.Background(), id)
	if err != nil {
		t.Fatalf("채널 읽기 실패: %v", err)
	}
	return ch
}

// deliver는 이벤트를 분배하고 지정한 채널에만 한 번 전송을 수행합니다.
func (f *notifyFixture) deliver(t *testing.T, chID int64, baseURL string) {
	t.Helper()
	ctx := context.Background()
	if _, _, err := f.pg.FanOutPendingEvents(ctx, 500); err != nil {
		t.Fatalf("분배 실패: %v", err)
	}
	f.n.stepRealtime(ctx, f.channel(t, chID), 50, baseURL)
}

// createChannel은 HTTP 엔드포인트로 채널을 생성하며 엔드포인트 자체의 검증 경로도 다룹니다.
func (f *notifyFixture) createChannel(t *testing.T, payload map[string]any) int64 {
	t.Helper()
	raw, _ := json.Marshal(payload)
	r := f.request("POST", "/api/notify/channels", string(raw))
	if r.Code != 200 {
		t.Fatalf("채널 생성 실패 %d: %s", r.Code, r.Body)
	}
	var res struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &res); err != nil || res.ID == 0 {
		t.Fatalf("채널 생성 반환값이 비정상입니다: %s (%v)", r.Body, err)
	}
	t.Cleanup(func() { f.pg.Exec(`DELETE FROM notification_channels WHERE id=$1`, res.ID) })
	return res.ID
}

// fakeWebhook은 받은 요청 본문을 기록하는 가짜 수신 측입니다.
type fakeWebhook struct {
	*httptest.Server
	mu     sync.Mutex
	bodies []map[string]any
}

func newFakeWebhook(t *testing.T) *fakeWebhook {
	t.Helper()
	f := &fakeWebhook{}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		f.mu.Lock()
		f.bodies = append(f.bodies, body)
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"errcode":0,"errmsg":"ok"}`)
	}))
	t.Cleanup(f.Close)
	return f
}

func (f *fakeWebhook) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.bodies)
}

func (f *fakeWebhook) body(t *testing.T, i int) map[string]any {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if i >= len(f.bodies) {
		t.Fatalf("가짜 수신 측이 받은 요청은 %d개뿐이므로 %d번째 요청을 가져올 수 없습니다", len(f.bodies), i)
	}
	return f.bodies[i]
}

func (f *fakeWebhook) last(t *testing.T) map[string]any {
	t.Helper()
	if f.count() == 0 {
		t.Fatal("가짜 수신 측이 요청을 하나도 받지 못했습니다")
	}
	return f.body(t, f.count()-1)
}

// markdownText는 요청 본문에서 본문을 추출하며, 각 제품의 필드명 차이를 처리합니다:
// DingTalk markdown은 `text`, ActionCard는 `text`, WeCom(기업용 위챗) markdown은 `content`를 사용합니다.
func markdownText(t *testing.T, body map[string]any) string {
	t.Helper()
	for _, key := range []string{"markdown", "actionCard"} {
		section, ok := body[key].(map[string]any)
		if !ok {
			continue
		}
		for _, field := range []string{"text", "content"} {
			if s, ok := section[field].(string); ok && s != "" {
				return s
			}
		}
	}
	t.Fatalf("요청 본문에 인식할 수 있는 본문이 없습니다: %v", body)
	return ""
}

// agePendingBatch는 해당 채널의 전송 대기 항목 시각을 과거로 바꿔 모아 보내기 배치의 만료를 테스트합니다.
func (f *notifyFixture) agePendingBatch(t *testing.T, chID int64) {
	t.Helper()
	if _, err := f.pg.Exec(`UPDATE notification_deliveries SET created_at = now() - interval '2 hours'
WHERE channel_id=$1 AND state=$2`, chID, db.NotifyStatePending); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyEndToEndRealtimeDelivery(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "实时推送",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
	})
	f.record(t, "SQL注入", "high")
	f.deliver(t, chID, "")

	if hook.count() != 1 {
		t.Fatalf("메시지 1개를 전송해야 하지만 실제로는 %d개입니다", hook.count())
	}
	text := markdownText(t, hook.last(t))
	for _, want := range []string{"SQL注入", "높음", "**요약**:"} {
		if !strings.Contains(text, want) {
			t.Fatalf("메시지 본문에 다음 항목이 없습니다: %q\n%s", want, text)
		}
	}
	// 전송 상태는 sent로 바뀌어야 합니다.
	var pending int
	if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1 AND state <> $2`,
		chID, db.NotifyStateSent).Scan(&pending); err != nil {
		t.Fatal(err)
	}
	if pending != 0 {
		t.Fatalf("전송 후에도 sent로 표시되지 않은 항목이 %d개 남았습니다", pending)
	}
}

func TestNotifyChannelAPIMasksSecretsAndPreservesOnUpdate(t *testing.T) {
	f := newNotifyFixture(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "掩码用例",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": "https://oapi.dingtalk.com/robot/send?access_token=abc123456", "secret": "SECabcdef123456"},
	})

	r := f.request("GET", "/api/notify/channels", "")
	if r.Code != 200 {
		t.Fatalf("채널 목록 조회 실패 %d: %s", r.Code, r.Body)
	}
	if strings.Contains(r.Body.String(), "abc123456") || strings.Contains(r.Body.String(), "SECabcdef123456") {
		t.Fatalf("엔드포인트의 반환값이 자격 증명을 노출했습니다: %s", r.Body)
	}
	var listed struct {
		Channels []struct {
			ID         int64          `json:"id"`
			Config     map[string]any `json:"config"`
			SecretKeys []string       `json:"secret_keys"`
		} `json:"channels"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	var mine *struct {
		ID         int64          `json:"id"`
		Config     map[string]any `json:"config"`
		SecretKeys []string       `json:"secret_keys"`
	}
	for i := range listed.Channels {
		if listed.Channels[i].ID == chID {
			mine = &listed.Channels[i]
		}
	}
	if mine == nil {
		t.Fatal("생성한 채널이 목록에 없습니다")
	}
	if !notify.IsMasked(fmt.Sprint(mine.Config["webhook"])) || !notify.IsMasked(fmt.Sprint(mine.Config["secret"])) {
		t.Fatalf("자격 증명 필드는 마스킹된 값이어야 합니다: %v", mine.Config)
	}
	if len(mine.SecretKeys) == 0 {
		t.Fatal("엔드포인트는 어떤 필드가 자격 증명인지 프런트엔드에 알려야 합니다")
	}

	// PATCH로 이름만 변경하고 마스킹된 자격 증명을 반환할 때 실제 자격 증명은 그대로 보존해야 합니다.
	body, _ := json.Marshal(map[string]any{
		"name":   "改名后",
		"config": map[string]any{"webhook": fmt.Sprint(mine.Config["webhook"]), "secret": fmt.Sprint(mine.Config["secret"])},
	})
	if r := f.request("PATCH", fmt.Sprintf("/api/notify/channels/%d", chID), string(body)); r.Code != 200 {
		t.Fatalf("업데이트 실패 %d: %s", r.Code, r.Body)
	}
	cfg := f.channelConfig(t, chID)
	if cfg["webhook"] != "https://oapi.dingtalk.com/robot/send?access_token=abc123456" {
		t.Fatalf("반환한 마스킹 값이 실제 자격 증명을 덮어썼습니다: %v", cfg["webhook"])
	}
	if cfg["secret"] != "SECabcdef123456" {
		t.Fatalf("반환한 마스킹 값이 secret을 덮어썼습니다: %v", cfg["secret"])
	}
	if f.channel(t, chID).Name != "改名后" {
		t.Fatal("이름이 업데이트되지 않았습니다")
	}

	// secret을 명시적으로 비우면 반영되어야 합니다('마스킹된 값 반환=변경 없이 유지'와 구분).
	body, _ = json.Marshal(map[string]any{"config": map[string]any{"secret": ""}})
	if r := f.request("PATCH", fmt.Sprintf("/api/notify/channels/%d", chID), string(body)); r.Code != 200 {
		t.Fatalf("secret 비우기 실패 %d: %s", r.Code, r.Body)
	}
	if _, still := f.channelConfig(t, chID)["secret"]; still {
		t.Fatal("빈 문자열이면 secret을 비워야 합니다")
	}
}

func (f *notifyFixture) channelConfig(t *testing.T, id int64) map[string]any {
	t.Helper()
	var cfg map[string]any
	if err := json.Unmarshal(f.channel(t, id).Config, &cfg); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestNotifyChannelAPICreateValidation(t *testing.T) {
	f := newNotifyFixture(t)
	cases := []struct {
		name    string
		payload map[string]any
		wantSub string
	}{
		{"유효하지 않은 유형", map[string]any{"name": "x", "kind": "nope", "config": map[string]any{}}, "채널 유형이 유효하지 않음"},
		{"이름 없음", map[string]any{"kind": notify.KindDingTalk, "config": map[string]any{"webhook": "https://e.com/h"}}, "채널 이름이 없습니다"},
		{"webhook 없음", map[string]any{"name": "x", "kind": notify.KindDingTalk, "config": map[string]any{}}, "Webhook"},
		{"유효하지 않은 webhook 프로토콜", map[string]any{"name": "x", "kind": notify.KindDingTalk, "config": map[string]any{"webhook": "file:///etc/passwd"}}, "Webhook 주소가 유효하지 않습니다"},
		{"유효하지 않은 모드", map[string]any{"name": "x", "kind": notify.KindDingTalk, "mode": "sometimes", "config": map[string]any{"webhook": "https://e.com/h"}}, "푸시 모드가 유효하지 않음"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, _ := json.Marshal(tc.payload)
			r := f.request("POST", "/api/notify/channels", string(raw))
			if r.Code != 400 {
				t.Fatalf("400을 반환해야 하지만 실제로는 %d입니다: %s", r.Code, r.Body)
			}
			if !strings.Contains(r.Body.String(), tc.wantSub) {
				t.Fatalf("오류 메시지에 %q가 있어야 하지만 실제로는 %s입니다", tc.wantSub, r.Body)
			}
		})
	}
	if r := f.request("DELETE", "/api/notify/channels/99999999", ""); r.Code != 404 {
		t.Fatalf("존재하지 않는 채널 삭제는 404여야 하지만 실제로는 %d입니다", r.Code)
	}
}

func TestNotifyFilterBlocksBelowThreshold(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "仅严重",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
		"filter": map[string]any{"min_severity": "critical"},
	})
	f.record(t, "低危问题", "low")
	if _, _, err := f.pg.FanOutPendingEvents(context.Background(), 500); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1`, chID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("임계값보다 낮은 취약점은 전송 항목을 생성하면 안 되지만 실제로는 %d개입니다", n)
	}
	f.n.stepRealtime(context.Background(), f.channel(t, chID), 50, "")
	if hook.count() != 0 {
		t.Fatal("필터링된 취약점은 메시지를 전송하면 안 됩니다")
	}
}

func TestNotifyDigestBatchesMultipleFindingsIntoOneMessage(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "汇总推送",
		"kind":   notify.KindDingTalk,
		"mode":   db.NotifyModeDigest,
		"config": map[string]any{"webhook": hook.URL},
	})
	for i := 0; i < 3; i++ {
		f.record(t, fmt.Sprintf("汇总漏洞%d", i+1), "high")
	}
	ctx := context.Background()
	if _, _, err := f.pg.FanOutPendingEvents(ctx, 500); err != nil {
		t.Fatal(err)
	}
	ch := f.channel(t, chID)

	// 아직 전송 시점이 되지 않았으면 전송하지 않습니다.
	f.n.stepDigest(ctx, ch, 50, "")
	if hook.count() != 0 {
		t.Fatal("모아 보내기 배치의 전송 시점 전에 전송되었습니다")
	}

	// 배치 시각을 과거로 바꾼 뒤에는 세 항목을 메시지 하나로 모읍니다.
	f.agePendingBatch(t, chID)
	f.n.stepDigest(ctx, ch, 50, "")
	if got := hook.count(); got != 1 {
		t.Fatalf("세 항목을 메시지 하나로 모아야 하지만 실제로는 %d개를 전송했습니다", got)
	}
	text := markdownText(t, hook.last(t))
	if !strings.Contains(text, "최근") || !strings.Contains(text, "새 취약점 3개") {
		t.Fatalf("모아 보내기 메시지에 개수/시간 범위 안내가 없습니다:\n%s", text)
	}
	for i := 1; i <= 3; i++ {
		if !strings.Contains(text, fmt.Sprintf("汇总漏洞%d", i)) {
			t.Fatalf("모아 보내기 메시지에 %d번째 항목이 없습니다:\n%s", i, text)
		}
	}
	// 같은 배치는 batch_id를 공유해야 합니다.
	var distinct, total int
	if err := f.pg.QueryRow(`SELECT count(DISTINCT batch_id), count(*) FROM notification_deliveries WHERE channel_id=$1`, chID).Scan(&distinct, &total); err != nil {
		t.Fatal(err)
	}
	if total != 3 || distinct != 1 {
		t.Fatalf("세 전송 항목이 하나의 batch_id를 공유해야 합니다, got distinct=%d total=%d", distinct, total)
	}
}

func TestNotifyDisabledChannelDoesNotSend(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":    "停用渠道",
		"kind":    notify.KindDingTalk,
		"enabled": false,
		"config":  map[string]any{"webhook": hook.URL},
	})
	f.record(t, "停用期间的漏洞", "critical")
	if _, _, err := f.pg.FanOutPendingEvents(context.Background(), 500); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1`, chID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("비활성화된 채널은 전송 항목을 생성하면 안 되지만 실제로는 %d개입니다", n)
	}
}

func TestNotifyStatusChangeDelivery(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "状态变更订阅",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
		"filter": map[string]any{"on_status_change": true},
	})
	finding := f.record(t, "状态变更用例", "high")
	r := f.request("PATCH", fmt.Sprintf("/api/exploration/findings/%d", finding), `{"status":"fixed"}`)
	if r.Code != 200 {
		t.Fatalf("상태 변경 실패 %d: %s", r.Code, r.Body)
	}
	f.deliver(t, chID, "")

	// 두 항목이 있어야 합니다: fixed 항목은 상태 변경이며 finding_created 항목도 같은 회차에 전송될 수 있습니다.
	// 상태 변경 항목이 실제로는 더 늦게 생성되지만, 순서에 의존하지 않고 모두 검색합니다.
	found := false
	for i := 0; i < hook.count(); i++ {
		text := markdownText(t, hook.body(t, i))
		if strings.Contains(text, "상태 변경") && strings.Contains(text, "수정 완료") {
			found = true
		}
	}
	if !found {
		t.Fatalf("'状态变更 → 已修复'가 포함된 메시지를 받지 못했습니다(총 %d개)", hook.count())
	}
}

func TestNotifyStatusChangeSuppressedByDefault(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "不订阅状态变更",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
	})
	finding := f.record(t, "不订阅变更", "high")
	if r := f.request("PATCH", fmt.Sprintf("/api/exploration/findings/%d", finding), `{"status":"false_positive"}`); r.Code != 200 {
		t.Fatalf("상태 변경 실패 %d: %s", r.Code, r.Body)
	}
	if _, _, err := f.pg.FanOutPendingEvents(context.Background(), 500); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries d
JOIN notification_events e ON e.id = d.event_id
WHERE d.channel_id=$1 AND e.kind=$2`, chID, notify.EventFindingStatusChanged).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("상태 변경을 구독하지 않은 채널은 상태 변경 전송을 받으면 안 되지만 실제로는 %d개입니다", n)
	}
}

func TestNotifyTestMessageEndpoint(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "测试发送",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
	})
	if r := f.request("POST", fmt.Sprintf("/api/notify/channels/%d/test", chID), ""); r.Code != 200 {
		t.Fatalf("테스트 전송 실패 %d: %s", r.Code, r.Body)
	}
	if hook.count() != 1 {
		t.Fatalf("가짜 수신 측이 테스트 메시지 1개를 받아야 하지만 실제로는 %d개입니다", hook.count())
	}
	// 테스트 메시지는 테스트임을 즉시 알아볼 수 있어야 하며 실제 취약점으로 오인해서는 안 됩니다.
	if text := markdownText(t, hook.last(t)); !strings.Contains(text, "테스트") {
		t.Fatalf("테스트 메시지는 테스트임을 표시해야 합니다: %s", text)
	}
	// 설정이 잘못되면 채널의 원본 오류를 사용자에게 그대로 반환해야 합니다.
	badID := f.createChannel(t, map[string]any{
		"name":   "坏地址",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": "http://127.0.0.1:1/hook"},
	})
	if r := f.request("POST", fmt.Sprintf("/api/notify/channels/%d/test", badID), ""); r.Code != 502 {
		t.Fatalf("전송 실패는 502를 반환해야 하지만 실제로는 %d입니다: %s", r.Code, r.Body)
	}
}

func TestNotifyDeliveriesHistoryAndRetry(t *testing.T) {
	f := newNotifyFixture(t)
	// 반드시 실패하는 주소를 지정해 failed 전송을 만듭니다.
	chID := f.createChannel(t, map[string]any{
		"name":   "失败重试",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": "http://127.0.0.1:1/hook"},
	})
	f.record(t, "会失败的推送", "high")
	ctx := context.Background()
	if _, _, err := f.pg.FanOutPendingEvents(ctx, 500); err != nil {
		t.Fatal(err)
	}
	ch := f.channel(t, chID)
	// 재시도 예산을 소진할 때까지 계속 전송합니다.
	for i := 0; i < db.MaxNotifyAttempts; i++ {
		f.n.stepRealtime(ctx, ch, 50, "")
		if _, err := f.pg.Exec(`UPDATE notification_deliveries SET next_attempt_at = now() - interval '1 minute' WHERE channel_id=$1`, chID); err != nil {
			t.Fatal(err)
		}
	}
	var state string
	if err := f.pg.QueryRow(`SELECT state FROM notification_deliveries WHERE channel_id=$1`, chID).Scan(&state); err != nil {
		t.Fatal(err)
	}
	if state != db.NotifyStateFailed {
		t.Fatalf("재시도 소진 후에는 failed여야 합니다, got %s", state)
	}

	r := f.request("GET", fmt.Sprintf("/api/notify/deliveries?channel_id=%d&state=failed", chID), "")
	if r.Code != 200 {
		t.Fatalf("이력 조회 실패 %d: %s", r.Code, r.Body)
	}
	var hist struct {
		Deliveries []struct {
			ID        int64  `json:"id"`
			State     string `json:"state"`
			LastError string `json:"last_error"`
			Attempts  int    `json:"attempts"`
			Title     string `json:"title"`
		} `json:"deliveries"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &hist); err != nil {
		t.Fatal(err)
	}
	if hist.Total != 1 || len(hist.Deliveries) != 1 {
		t.Fatalf("실패한 전송 항목 1개가 조회되어야 합니다, got total=%d len=%d", hist.Total, len(hist.Deliveries))
	}
	if hist.Deliveries[0].LastError == "" {
		t.Fatal("이력에 실패 원인이 있어야 사용자가 원인을 찾을 수 있습니다")
	}
	if hist.Deliveries[0].Attempts < db.MaxNotifyAttempts {
		t.Fatalf("시도 횟수가 기록되어야 합니다, got %d", hist.Deliveries[0].Attempts)
	}
	if hist.Deliveries[0].Title != "会失败的推送" {
		t.Fatalf("이력에 취약점 제목이 포함되어야 합니다, got %q", hist.Deliveries[0].Title)
	}

	// 수동 재전송 시 pending으로 돌아가고 횟수를 0으로 초기화해야 합니다.
	if r := f.request("POST", fmt.Sprintf("/api/notify/deliveries/%d/retry", hist.Deliveries[0].ID), ""); r.Code != 200 {
		t.Fatalf("재전송 실패 %d: %s", r.Code, r.Body)
	}
	var attempts int
	if err := f.pg.QueryRow(`SELECT state, attempts FROM notification_deliveries WHERE id=$1`, hist.Deliveries[0].ID).Scan(&state, &attempts); err != nil {
		t.Fatal(err)
	}
	if state != db.NotifyStatePending || attempts != 0 {
		t.Fatalf("재전송 후 pending이며 attempts=0이어야 합니다, got %s/%d", state, attempts)
	}
}

func TestNotifyMetaAndSettingsRoundTrip(t *testing.T) {
	f := newNotifyFixture(t)
	r := f.request("GET", "/api/notify/meta", "")
	if r.Code != 200 {
		t.Fatalf("meta 실패: %s", r.Body)
	}
	var meta struct {
		Kinds []struct {
			Kind       string   `json:"kind"`
			SecretKeys []string `json:"secret_keys"`
		} `json:"kinds"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &meta); err != nil {
		t.Fatal(err)
	}
	if len(meta.Kinds) != len(notify.Kinds()) {
		t.Fatalf("meta에는 채널 %d개가 모두 있어야 합니다, got %d", len(notify.Kinds()), len(meta.Kinds))
	}
	for _, k := range meta.Kinds {
		if len(k.SecretKeys) == 0 {
			t.Errorf("채널 %s: 자격 증명 필드가 보고되지 않았습니다", k.Kind)
		}
	}

	// 전역 설정 세 항목의 왕복 확인. 끝의 슬래시는 정규화하여 제거해야 하며, 그렇지 않으면 상세 링크에 "//function/..."가 생깁니다.
	if r := f.request("PUT", "/api/settings", `{"notify_public_base_url":"https://artex.example.com/","notify_digest_interval_min":15,"notify_enabled":true}`); r.Code != 200 {
		t.Fatalf("설정 쓰기 실패 %d: %s", r.Code, r.Body)
	}
	t.Cleanup(func() {
		f.pg.Exec(`DELETE FROM settings WHERE key IN ($1,$2)`, settingNotifyPublicBaseURL, settingNotifyDigestMinutes)
	})
	payload := f.s.settingsPayload()
	if payload["notify_public_base_url"] != "https://artex.example.com" {
		t.Fatalf("상세 링크 주소가 정규화되지 않았습니다: %v", payload["notify_public_base_url"])
	}
	if payload["notify_digest_interval_min"] != 15 {
		t.Fatalf("모아 보내기 주기가 반영되지 않았습니다: %v", payload["notify_digest_interval_min"])
	}

	// 유효하지 않은 값은 거부해야 합니다.
	for _, body := range []string{
		`{"notify_public_base_url":"ftp://x"}`,
		`{"notify_digest_interval_min":0}`,
		`{"notify_digest_interval_min":99999}`,
	} {
		if r := f.request("PUT", "/api/settings", body); r.Code != 400 {
			t.Errorf("%s는 400을 반환해야 합니다, got %d", body, r.Code)
		}
	}
}

// TestNotifyDeepLinkUsesPublicBaseURL은 상세 링크 구성을 다룹니다: public_base_url이 설정되어 있으면
// 단일 메시지는 반드시 버튼이 있는 ActionCard를 사용하며 링크는 취약점 상세 페이지를 가리켜야 합니다.
func TestNotifyDeepLinkUsesPublicBaseURL(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "回链",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
	})
	finding := f.record(t, "带回链的漏洞", "high")
	f.deliver(t, chID, "https://artex.example.com")

	body := hook.last(t)
	card, _ := body["actionCard"].(map[string]any)
	if card == nil {
		t.Fatalf("상세 링크가 있으면 ActionCard를 사용해야 합니다, got msgtype=%v", body["msgtype"])
	}
	want := fmt.Sprintf("https://artex.example.com/function/findings/detail?id=%d", finding)
	if card["singleURL"] != want {
		t.Fatalf("상세 링크가 잘못되었습니다\n기대값 %s\n실제 %v", want, card["singleURL"])
	}
}

// TestNotifyNoDeepLinkWithoutBaseURL은 반대 조건을 다룹니다: 외부 주소가 설정되지 않으면 잘못된 링크를 만들면 안 됩니다.
// (예: localhost 또는 상대 경로를 가리키는 링크.) 순수 markdown으로 돌아가야 합니다.
func TestNotifyNoDeepLinkWithoutBaseURL(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	chID := f.createChannel(t, map[string]any{
		"name":   "无回链",
		"kind":   notify.KindDingTalk,
		"config": map[string]any{"webhook": hook.URL},
	})
	f.record(t, "无回链的漏洞", "high")
	f.deliver(t, chID, "")

	body := hook.last(t)
	if body["msgtype"] != "markdown" {
		t.Fatalf("외부 주소가 설정되지 않으면 markdown을 전송해야 합니다, got %v", body["msgtype"])
	}
	if text := markdownText(t, body); strings.Contains(text, "상세 보기") {
		t.Fatalf("외부 주소를 설정하지 않았을 때는 상세 링크가 나타나면 안 됩니다:\n%s", text)
	}
}

// TestNotifyDigestSegmentsAndDefersRemainder는 '알리지 않고 항목을 잃어버리는 문제' 수정의 엔드투엔드 증거입니다.
//
// 모아 보내기 메시지는 채널 길이 상한(WeCom 4096바이트)의 제약을 받으며, 배치 하나를 담지 못하면 반드시 **항목 단위**로 나눠야 합니다:
// 이 메시지에 담긴 항목은 전송 완료로 표시하고 나머지는 대기열로 돌려 다음 메시지를 기다립니다. 이전 구현은 전체 배치를
// 성공으로 표시했습니다. 잘린 항목은 메시지에도 실패 목록에도 없는데 전송 이력에는 성공으로 표시되어,
// 취약점이 그대로 사라졌습니다.
//
// 네 가지를 검증합니다: ① 실제로 담은 항목만 표시 ② 나머지는 전송 대기 유지 ③ 다음으로 미룬 항목은
// **재시도 횟수를 소비하지 않음** ④ 한 번 더 실행하면 나머지를 전송할 수 있음(진행이 막히지 않음).
func TestNotifyDigestSegmentsAndDefersRemainder(t *testing.T) {
	f := newNotifyFixture(t)
	hook := newFakeWebhook(t)
	// WeCom을 사용합니다: markdown 상한 4096바이트로 여섯 채널 중 가장 엄격합니다.
	chID := f.createChannel(t, map[string]any{
		"name":   "分段汇总",
		"kind":   notify.KindWeCom,
		"mode":   db.NotifyModeDigest,
		"config": map[string]any{"webhook": hook.URL},
	})
	const total = 60
	// 제목을 길게 만들어 60개 항목이 4096바이트를 훨씬 넘고 반드시 나뉘도록 합니다.
	longName := strings.Repeat("超长漏洞名称", 6)
	for i := 0; i < total; i++ {
		f.record(t, longName+strconv.Itoa(i+1), "high")
	}
	ctx := context.Background()
	if _, _, err := f.pg.FanOutPendingEvents(ctx, 500); err != nil {
		t.Fatal(err)
	}
	f.agePendingBatch(t, chID)
	ch := f.channel(t, chID)

	f.n.stepDigest(ctx, ch, 50, "")
	if hook.count() != 1 {
		t.Fatalf("메시지 하나만 전송해야 합니다, got %d", hook.count())
	}

	var sent, pending int
	if err := f.pg.QueryRow(`SELECT
    count(*) FILTER (WHERE state=$2),
    count(*) FILTER (WHERE state=$3)
  FROM notification_deliveries WHERE channel_id=$1`, chID, db.NotifyStateSent, db.NotifyStatePending).
		Scan(&sent, &pending); err != nil {
		t.Fatal(err)
	}
	if sent == 0 {
		t.Fatal("전송 완료로 표시된 항목이 있어야 합니다")
	}
	if pending == 0 {
		t.Fatalf("배치의 %d개 항목을 4096바이트에 모두 담을 수 없으므로 전송 대기 항목이 남아야 합니다; sent=%d", total, sent)
	}
	if sent+pending != total {
		t.Fatalf("항목 수가 일치하지 않습니다: sent=%d pending=%d total=%d(전송 완료도 대기도 아님=손실)", sent, pending, total)
	}
	// 메시지 본문은 이 메시지에 포함되지 않은 항목 수를 정확히 알려야 합니다.
	if text := markdownText(t, hook.last(t)); !strings.Contains(text, "나머지") {
		t.Fatalf("메시지는 이번 메시지에 포함하지 않은 항목이 남아 있음을 알려야 합니다:\n%.400s", text)
	}

	// 다음으로 미룬 항목은 재시도 예산을 소비하면 안 됩니다: 획득 시 attempts를 낙관적으로 +1 했으므로 미룰 때 다시 차감해야 합니다.
	var maxAttempts int
	if err := f.pg.QueryRow(`SELECT COALESCE(max(attempts),0) FROM notification_deliveries
WHERE channel_id=$1 AND state=$2`, chID, db.NotifyStatePending).Scan(&maxAttempts); err != nil {
		t.Fatal(err)
	}
	if maxAttempts > 0 {
		t.Fatalf("미룬 항목은 재시도 횟수를 소비하면 안 됩니다(그렇지 않으면 몇 번 만에 실패로 판정됨), got attempts=%d", maxAttempts)
	}

	// 모두 처리될 때까지 반복 실행합니다. **최종적으로 모두 전송**되며 중간에 실제로 여러 회차로 나뉘었는지를 검증합니다.
	// '두 번째 회차에 전송 완료'보다 강한 검사로, 분할 처리가 막히지 않고 나머지 항목도 잃지 않음을 증명합니다.
	rounds := 0
	for {
		var undelivered int
		if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries
WHERE channel_id=$1 AND state <> $2 AND state <> $3`, chID, db.NotifyStateSent, db.NotifyStateFailed).
			Scan(&undelivered); err != nil {
			t.Fatal(err)
		}
		if undelivered == 0 {
			break
		}
		rounds++
		if rounds > total+5 {
			t.Fatalf("분할 전송이 수렴하지 않습니다: %d회 실행한 뒤에도 미처리 항목 %d개가 남았습니다", rounds, undelivered)
		}
		before := hook.count()
		f.n.stepDigest(ctx, ch, 50, "")
		if hook.count() == before {
			t.Fatalf("%d번째 회차에 아무 진행이 없으므로 남은 항목 %d개가 영구적으로 처리되지 않습니다", rounds, undelivered)
		}
	}
	if rounds < 2 {
		t.Fatalf("4096바이트 메시지 하나에 긴 취약점 제목 %d개를 담을 수 없으므로 여러 회차로 전송해야 하지만 실제로는 %d회뿐입니다", total, rounds)
	}
	// 첫 회차 이후에는 모든 회차가 **나머지 항목의 후속 전송만** 수행하며, 채널이 거부한 항목은 없습니다.
	var failed int
	if err := f.pg.QueryRow(`SELECT count(*) FROM notification_deliveries WHERE channel_id=$1 AND state=$2`,
		chID, db.NotifyStateFailed).Scan(&failed); err != nil {
		t.Fatal(err)
	}
	if failed != 0 {
		t.Fatalf("가짜 수신 측은 항상 성공을 반환하므로 실패 항목이 없어야 합니다, got %d", failed)
	}
}

// TestNotifyBackoffTableMatchesAttemptBudget은 변경에 따른 불일치를 방지하는 검사입니다.
//
// 재시도 예산(db.MaxNotifyAttempts)과 백오프 간격 목록(notifyBackoff)은 두 패키지에 나뉘어 있습니다:
// 전자는 상태 머신의 정책이고 후자는 엔진의 실행 간격입니다. 한쪽만 바꾸면, 예를 들어 예산을 5회로
// 늘리면서 백오프 단계를 추가하지 않으면 코드 오류 없이 4번째와 5번째 재시도에 마지막 간격이 그대로 적용됩니다.
// '재시도가 갑자기 느려짐'으로 나타나므로, 원인을 찾을 때 이 부분을 떠올리기 어렵습니다.
// 두 길이가 같은지 검증하여 이런 불일치가 CI에서 드러나게 합니다.
func TestNotifyBackoffTableMatchesAttemptBudget(t *testing.T) {
	if len(notifyBackoff) != db.MaxNotifyAttempts {
		t.Fatalf("백오프 단계 수(%d)와 최대 시도 횟수(%d)가 다릅니다. 하나를 바꾸면 다른 하나도 함께 바꿔야 합니다",
			len(notifyBackoff), db.MaxNotifyAttempts)
	}
	// 백오프 간격은 줄어들면 안 됩니다. 그렇지 않으면 재시도가 점점 빨라져 전송 속도 제한을 더 악화시킵니다.
	for i := 1; i < len(notifyBackoff); i++ {
		if notifyBackoff[i] < notifyBackoff[i-1] {
			t.Fatalf("백오프 간격은 줄어들면 안 됩니다: %d번째 단계 %v < %d번째 단계 %v",
				i, notifyBackoff[i], i-1, notifyBackoff[i-1])
		}
	}
}

// TestNotifyRateLimitDoesNotConsumeRetryBudget은 '토큰을 먼저 얻고 전송 항목을 획득'하는 순서를 고정합니다.
// 순서가 반대라면(먼저 획득한 뒤 포기) 전송 속도 제한에 막힌 전송도 attempts를 한 번 늘린 상태여서,
// 단순히 기다리는 것만으로 예산을 소진하고 결국 failed가 됩니다.
func TestNotifyRateLimitDoesNotConsumeRetryBudget(t *testing.T) {
	// 토큰 버킷 자체만 검사하므로 Server는 필요하지 않습니다(이 검사를 위해 만들어서도 안 됩니다).
	n := &Notifier{buckets: map[int64]*notifyBucket{}}
	now := time.Now()
	// 분당 1개: 가득 찼을 때 최대 1개입니다.
	if got := n.takeTokens(1, 1, notifyMaxSendsPerChannelPerTick, now); got != 1 {
		t.Fatalf("가득 찼을 때 분당 1개이면 토큰 1개를 얻어야 합니다, got %d", got)
	}
	if got := n.takeTokens(1, 1, notifyMaxSendsPerChannelPerTick, now.Add(time.Millisecond)); got != 0 {
		t.Fatalf("토큰을 소진하면 즉시 0을 반환해야 합니다, got %d", got)
	}
	if got := n.takeTokens(1, 1, notifyMaxSendsPerChannelPerTick, now.Add(30*time.Second)); got != 0 {
		t.Fatalf("절반의 시간이 지났을 때 토큰 하나가 완전히 보충되면 안 됩니다, got %d", got)
	}
	if got := n.takeTokens(1, 1, notifyMaxSendsPerChannelPerTick, now.Add(time.Minute)); got != 1 {
		t.Fatalf("한 주기가 지나면 토큰 1개가 보충되어야 합니다, got %d", got)
	}
	// 전송 속도 제한이 없는 채널도 유한한 상한을 사용해 한 회차가 무한한 적체에 붙잡히지 않도록 합니다.
	if got := n.takeTokens(2, 0, notifyUnlimitedBurstPerTick+10, now); got != notifyUnlimitedBurstPerTick {
		t.Fatalf("전송 속도 제한이 없으면 회차별 상한 %d를 반환해야 합니다, got %d", notifyUnlimitedBurstPerTick, got)
	}
	// 채널별 토큰 저장량은 서로 독립적입니다.
	if got := n.takeTokens(1, 1, notifyMaxSendsPerChannelPerTick, now.Add(time.Millisecond)); got != 0 {
		t.Fatalf("채널 1의 토큰 저장량은 여전히 비어 있어야 합니다, got %d", got)
	}
}

// TestNotifyTakeTokensKeepsUnusedTokens는 'want개만 획득'하는 의미를 고정합니다.
//
// 이전 구현은 저장한 토큰을 모두 꺼낸 다음 호출자가 잘랐으므로, rate=100/min인 채널에 토큰이 가득 찼을 때도
// 한 회차에 5개만 쓰면 나머지 95개를 그대로 버렸습니다. 이 회차에 전송 대기 항목이 없을 때도 차감했습니다.
// 그 결과 주석의 '적체 시 한 번에 rate_per_min개 전송'을 어떤 경우에도 달성할 수 없었습니다.
func TestNotifyTakeTokensKeepsUnusedTokens(t *testing.T) {
	n := &Notifier{buckets: map[int64]*notifyBucket{}}
	now := time.Now()
	// 처음에 토큰이 가득 차 있으며(100), 이번 회차에는 5개만 필요합니다.
	if got := n.takeTokens(1, 100, 5, now); got != 5 {
		t.Fatalf("want=5이면 토큰을 정확히 5개 얻어야 합니다, got %d", got)
	}
	// 핵심 검사: 나머지 95개는 모두 꺼내 버리는 대신 반드시 그대로 남아 있어야 합니다.
	// 시간을 진행시키지 않아 획득하는 토큰이 새로 보충된 것이 아니라 기존 저장량에서 나온 것임을 보장합니다.
	if got := n.takeTokens(1, 100, 95, now); got != 95 {
		t.Fatalf("남은 토큰을 계속 사용할 수 있어야 합니다(기대값 95), got %d. 토큰을 전부 꺼내 버렸습니다", got)
	}
	if got := n.takeTokens(1, 100, 1, now); got != 0 {
		t.Fatalf("토큰을 모두 얻었으면 0을 반환해야 합니다, got %d", got)
	}
	// want<=0이면 토큰을 차감하지 않아야 합니다(빈 회차는 비용을 소비하지 않음).
	n2 := &Notifier{buckets: map[int64]*notifyBucket{}}
	if got := n2.takeTokens(1, 20, 0, now); got != 0 {
		t.Fatalf("want=0이면 0을 반환해야 합니다, got %d", got)
	}
	if got := n2.takeTokens(1, 20, 20, now); got != 20 {
		t.Fatalf("want=0 호출에서 토큰을 소비하면 안 되므로 20개를 모두 얻을 수 있어야 합니다, got %d", got)
	}
}

// TestDigestTickPlanDecouplesBatchSizeFromSendBudget은 요청 횟수와 배치당 취약점 개수의 측정 단위를 구분하도록 보장합니다.
//
// 모아 보내기 배치 크기를 회차별 요청 예산에 연결하면 rate_per_min=20인 채널은 매
// 3초 tick마다 토큰 1개만 보충하므로 모아 보내기 메시지 하나에 취약점 1개만 담게 되어, 사실상 모아 보내기가 없는 것과 같습니다.
// 그런데도 메시지 머리에는 「近 30 分钟新增 1 个漏洞」라고 적혀 있습니다. 이 기능 저하는 오류를 내지 않으며,
// 기존 엔드투엔드 사례에서도 드러나지 않습니다(stepDigest에 충분히 큰 limit를 직접 전달하여
// step의 할당량 계산을 건너뛰기 때문입니다). 따라서 여기서는 결정 자체를 직접 검증합니다.
func TestDigestTickPlanDecouplesBatchSizeFromSendBudget(t *testing.T) {
	tokens, claimLimit := digestTickPlan()
	// 배치 하나 = 메시지 하나 = 요청 한 번 = 토큰 하나. 토큰의 단위는 취약점이 아니라 메시지입니다.
	if tokens != 1 {
		t.Fatalf("모아 보내기 배치는 메시지 하나만 전송하므로 토큰을 정확히 1개 소비해야 합니다, got %d", tokens)
	}
	if claimLimit != db.MaxDigestBatchSize {
		t.Fatalf("모아 보내기 배치 크기는 메모리 상한 db.MaxDigestBatchSize=%d여야 합니다, got %d",
			db.MaxDigestBatchSize, claimLimit)
	}
	// 핵심 관계: 배치 크기는 회차별 요청 예산보다 훨씬 커야 합니다. 두 값이 같은 수준이면,
	// '메시지를 몇 개 전송하는가'와 '배치 하나에 취약점을 몇 개 담는가'를 하나의 수로 혼동한 것입니다.
	if claimLimit <= notifyMaxSendsPerChannelPerTick {
		t.Fatalf("모아 보내기 배치 크기 %d는 회차별 요청 예산 %d에 제한되어서는 안 됩니다. "+
			"요청 예산은 lease에서 역산한 요청 횟수이며, 요청 횟수와 배치당 취약점 개수는 측정 단위가 다릅니다",
			claimLimit, notifyMaxSendsPerChannelPerTick)
	}
}

// TestNotifyTickBudgetFitsWithinLease도 변경에 따른 불일치를 방지하는 검사입니다.
//
// 단일 채널의 회차별 전송 상한(notifyMaxSendsPerChannelPerTick)은 lease 기간에서 역산합니다:
// 한 회차의 순차 전송에 걸리는 최악의 시간은 반드시 lease보다 짧아야 합니다. 그렇지 않으면 뒤쪽 항목을 전송하기 전에 lease가 만료되어,
// 여러 인스턴스가 배포된 경우 다른 인스턴스가 다시 획득해 중복 전송합니다. 이 세 상수는 서로 다른 위치에 있으며,
// 어느 하나만 바꿔도 아무 오류 없이 이 관계가 깨질 수 있어 여기서 고정합니다.
func TestNotifyTickBudgetFitsWithinLease(t *testing.T) {
	worst := time.Duration(notifyMaxSendsPerChannelPerTick) * notifySendTimeout
	if worst >= notifyLease {
		t.Fatalf("단일 채널 한 회차의 최악의 소요 시간 %v는 lease %v 이상이 되면 안 됩니다"+
			"(notifyMaxSendsPerChannelPerTick=%d × notifySendTimeout=%v). "+
			"이 세 상수 중 하나를 바꾸면 다른 둘도 함께 확인해야 합니다",
			worst, notifyLease, notifyMaxSendsPerChannelPerTick, notifySendTimeout)
	}
}
