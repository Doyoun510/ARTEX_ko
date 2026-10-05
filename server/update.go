package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Autumn-27/artex/selfupdate"
)

// 페이지 원클릭 업데이트의 HTTP 면. 실제 다운로드/검증/교체 로직은 모두 selfupdate 패키지에 있고,
// 여기서는 인증 경계·동시성 상호 배제·진행 브로드캐스트, 그리고 "이제 종료할 때"를 main에 알리는 것만 담당.
//
// 재시작은 이 프로세스가 하지 않는다: 새 버전을 임시 저장한 뒤 프로세스는 selfupdate.ExitRestart로 종료하고,
// 데몬 스크립트(start.sh / start.bat, Docker에서는 ENTRYPOINT)가 다시 띄운다.

// restartCh는 업그레이드 준비 완료 또는 롤백 완료 후 닫히며, main이 받으면 ExitRestart로 종료.
var (
	restartOnce sync.Once
	restartCh   = make(chan struct{})
)

// RestartRequested는 "종료하고 데몬이 나를 다시 띄우게 하라"일 때 닫히는 channel을 반환.
func RestartRequested() <-chan struct{} { return restartCh }

func requestRestart() { restartOnce.Do(func() { close(restartCh) }) }

// bootState는 이번 기동 시 selfupdate.Bootstrap의 결론(업그레이드 성공 / 방금 롤백 /
// 임시 저장 파일 폐기), main이 주입하며, /api/update/check가 지난 업그레이드의 결말을 프런트에 그대로 알리는 데 쓴다.
var (
	bootStateMu sync.Mutex
	bootState   selfupdate.State
)

// SetBootUpdateState는 main이 기동 시 한 번 호출.
func SetBootUpdateState(st selfupdate.State) {
	bootStateMu.Lock()
	defer bootStateMu.Unlock()
	bootState = st
}

func bootUpdateState() selfupdate.State {
	bootStateMu.Lock()
	defer bootStateMu.Unlock()
	return bootState
}

// releaseCache는 GitHub의 최신 버전 조회 결과를 캐시.
//
// 상단바의 "새 버전 있음" 안내는 전체 페이지 로드마다 한 번 조회하는데, 미인증 GitHub API는
// IP당 시간당 60회 —— 캐시하지 않으면 탭을 몇 개 더 열거나 페이지를 몇 번 새로고침하면 쿼터가 소진돼,
// 정작 업데이트하려 할 때 조회가 안 된다. 사용자가 "업데이트 확인"을 명시적으로 누르면 force로 캐시를 우회할 수 있다.
type releaseCache struct {
	mu  sync.Mutex
	rel *selfupdate.Release
	err error
	at  time.Time
	// fetch는 데이터 조회 함수로, 테스트용 주입 지점일 뿐; nil이면 실제 GitHub 조회를 탄다.
	fetch func(context.Context, *http.Client) (*selfupdate.Release, error)
}

const (
	releaseTTL = 30 * time.Minute
	// 실패 결과도 잠깐 캐시하며, 아니면 GitHub에 도달 못 할 때 페이지 로드마다 타임아웃을 한 번씩 헛되이 기다린다;
	// 단 TTL은 짧게, 네트워크 복구 후 금방 스스로 회복되도록.
	releaseErrTTL = 2 * time.Minute
	// 조회용 타임아웃. NewClient의 30분 타임아웃은 전체 패키지 다운로드용이라, 버전 조회는 그렇게 오래 기다릴 수 없다.
	releaseTimeout = 20 * time.Second
)

var relCache = &releaseCache{}

// get은 최신 Release를 반환하며, 캐시에 적중하면 네트워크에 접근하지 않는다.
//
// 조회 중 계속 락을 보유: 동시 요청은 같은 조회 결과를 줄 서서 기다리고, 각자 GitHub를 치지 않는다
// (페이지 로드 직후 여러 탭이 동시에 조회하는 것이 바로 레이트 리밋을 가장 쉽게 유발하는 순간).
func (c *releaseCache) get(ctx context.Context, client *http.Client, force bool) (*selfupdate.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !force {
		ttl := releaseTTL
		if c.err != nil {
			ttl = releaseErrTTL
		}
		if !c.at.IsZero() && time.Since(c.at) < ttl {
			return c.rel, c.err
		}
	}

	fetch := c.fetch
	if fetch == nil {
		fetch = selfupdate.FetchLatest
	}
	ctx, cancel := context.WithTimeout(ctx, releaseTimeout)
	defer cancel()
	rel, err := fetch(ctx, client)
	// 요청 취소(사용자가 탭을 닫음)는 GitHub 문제를 뜻하지 않으니 캐시에 쓰지 말 것,
	// 아니면 다음 방문자가 영문 모를 "취소됨" 오류를 받는다.
	if err != nil && ctx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
		return c.rel, err
	}
	c.rel, c.err, c.at = rel, err, time.Now()
	return rel, err
}

// updateProgress는 프런트에 푸시하는 진행 상황 하나.
type updateProgress struct {
	Phase   selfupdate.Phase `json:"phase"`
	Percent int              `json:"percent"` // 다운로드 단계에서만 의미 있음; 그 외는 -1
	Message string           `json:"message"`
	Version string           `json:"version,omitempty"`
	Error   string           `json:"error,omitempty"`
}

// updateHub는 한 업그레이드의 진행 상황을 보유하고 SSE 구독자에게 브로드캐스트.
//
// running은 상호 배제도 겸한다: 업그레이드 중 /api/update/apply를 다시 POST하면 바로 409,
// 두 goroutine이 같은 artex.new에 동시에 쓰는 것을 방지.
type updateHub struct {
	mu      sync.Mutex
	running bool
	cur     updateProgress
	subs    map[chan updateProgress]struct{}
}

var updHub = &updateHub{
	cur:  updateProgress{Phase: selfupdate.PhaseIdle, Percent: -1},
	subs: map[chan updateProgress]struct{}{},
}

// begin은 업그레이드 권한을 선점하며, 이미 진행 중이면 false 반환.
func (h *updateHub) begin(version string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.running {
		return false
	}
	h.running = true
	h.cur = updateProgress{Phase: selfupdate.PhaseDownload, Percent: 0, Message: "준비 중…", Version: version}
	h.fanout(h.cur)
	return true
}

// finish는 한 업그레이드를 끝낸다. err가 nil이면 임시 저장 성공, 재시작 대기.
func (h *updateHub) finish(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.running = false
	if err != nil {
		h.cur = updateProgress{Phase: selfupdate.PhaseFailed, Percent: -1, Message: "업데이트 실패", Error: err.Error(), Version: h.cur.Version}
	} else {
		h.cur = updateProgress{Phase: selfupdate.PhaseStaged, Percent: 100, Message: "새 버전 준비 완료, 재시작 중…", Version: h.cur.Version}
	}
	h.fanout(h.cur)
}

func (h *updateHub) publish(ph selfupdate.Phase, pct int, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cur = updateProgress{Phase: ph, Percent: pct, Message: msg, Version: h.cur.Version}
	h.fanout(h.cur)
}

// fanout은 h.mu를 보유한 채 호출해야 한다. 구독자 channel은 버퍼가 있고, 차면 버린다 ——
// 진행 상황은 버릴 수 있는 순간 정보라, 막힌 SSE 연결이 업그레이드 자체를 막게 해선 절대 안 된다.
func (h *updateHub) fanout(p updateProgress) {
	for ch := range h.subs {
		select {
		case ch <- p:
		default:
		}
	}
}

func (h *updateHub) snapshot() (updateProgress, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.cur, h.running
}

func (h *updateHub) subscribe() (<-chan updateProgress, func()) {
	ch := make(chan updateProgress, 64)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	var once sync.Once
	return ch, func() {
		once.Do(func() {
			h.mu.Lock()
			delete(h.subs, ch)
			h.mu.Unlock()
			close(ch)
		})
	}
}

// updateCheck는 GitHub의 최신 정식 버전을 조회해 현재 버전과 비교.
//
// 프런트도 api.github.com에 직접 연결하지만(GitHub의 CORS는 *), **이 인터페이스를 기준으로 한다**:
// 다운로드는 백엔드가 하므로, 백엔드가 GitHub에 접근할 수 있어야만 업데이트가 가능하다. 브라우저는 되고 서버는 안 되는
// 경우가 흔하며(서버가 내부망에 있거나, 프록시가 브라우저에만 설정됨), 그때 업데이트를 누르면 반드시 실패하니,
// 차라리 확인 단계에서 그대로 오류를 내는 편이 낫다.
func (s *Server) updateCheck(w http.ResponseWriter, r *http.Request) {
	current := BuildVersion
	mode := "binary"
	if selfupdate.InDocker() {
		mode = "docker"
	}
	boot := bootUpdateState()
	out := map[string]any{
		"current":     current,
		"mode":        mode,
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"has_backup":  selfupdate.HasBackup(),
		"repo":        selfupdate.Repo,
		"boot_notice": boot.Detail,
		"rolled_back": boot.RolledBack,
	}

	// 상단바 안내는 캐시를 탄다(기본); 사용자가 "업데이트 확인"을 누르면 force=1로 강제 재조회.
	force := r.URL.Query().Get("force") != ""
	client := selfupdate.NewClient(s.m.GlobalProxy())
	rel, err := relCache.get(r.Context(), client, force)
	if err != nil {
		out["error"] = err.Error()
		writeJSON(w, 200, out)
		return
	}

	latest := rel.TagName
	out["latest"] = latest
	out["notes"] = rel.Body
	out["html_url"] = rel.HTMLURL
	if !rel.PublishedAt.IsZero() {
		out["published_at"] = rel.PublishedAt.Format(time.RFC3339)
	}

	asset := selfupdate.AssetName(latest, runtime.GOOS, runtime.GOARCH)
	out["asset"] = asset
	if a, ok := rel.FindAsset(asset); ok {
		out["asset_available"] = true
		out["size"] = a.Size
	} else {
		out["asset_available"] = false
	}

	cmp, comparable := selfupdate.CompareVersions(current, latest)
	out["comparable"] = comparable
	out["has_update"] = comparable && cmp < 0
	if !comparable {
		// 개발 빌드(dev / git describe 접미사 포함)는 비교할 버전 번호가 없다. 허용하면
		// 정식 버전이 로컬에서 디버깅 중인 바이너리를 덮어쓸 뿐이라, 아예 업데이트를 주지 않는다.
		out["reason"] = fmt.Sprintf("현재 버전 %q는 정식 릴리스가 아니라 원클릭 업데이트가 비활성화됨", current)
	}
	writeJSON(w, 200, out)
}

// updateApply는 새 버전을 다운로드·임시 저장하고, 완료 후 프로세스를 종료해 데몬 스크립트가 재시작하게 한다.
//
// 즉시 202 반환, 실제 작업은 백그라운드 goroutine에서 돈다: 전체 패키지 다운로드는 몇 분 걸릴 수 있어,
// 요청에 매달면 리버스 프록시 타임아웃에 끊긴다. 진행 상황은 /api/update/stream으로.
func (s *Server) updateApply(w http.ResponseWriter, r *http.Request) {
	current := BuildVersion

	// 캐시를 탄다: 설치되는 것이 사용자가 화면에서 보고 확인한 그 버전임을 보장.
	client := selfupdate.NewClient(s.m.GlobalProxy())
	rel, err := relCache.get(r.Context(), client, false)
	if err != nil {
		writeErr(w, 502, err.Error())
		return
	}
	cmp, comparable := selfupdate.CompareVersions(current, rel.TagName)
	if !comparable {
		writeErr(w, 400, fmt.Sprintf("현재 버전 %q는 정식 릴리스가 아니라 원클릭 업데이트가 비활성화됨", current))
		return
	}
	if cmp >= 0 {
		writeErr(w, 400, fmt.Sprintf("이미 최신 버전입니다 %s", current))
		return
	}
	if !updHub.begin(rel.TagName) {
		writeErr(w, 409, "이미 업데이트가 진행 중입니다")
		return
	}

	go func() {
		// 의도적으로 요청 ctx가 아니라 s.ctx 사용: HTTP 응답이 반환되면 요청은 끝나,
		// 거기에 매달아 다운로드하면 즉시 취소된다.
		err := selfupdate.Stage(s.ctx, client, rel, current, func(ph selfupdate.Phase, pct int, msg string) {
			updHub.publish(ph, pct, msg)
		})
		updHub.finish(err)
		if err != nil {
			log.Printf("[update] 업데이트 실패: %v", err)
			return
		}
		log.Printf("[update] %s → %s 임시 저장됨, 교체 완료를 위해 곧 종료", current, rel.TagName)
		// 마지막 진행 상황을 프런트에 밀어줄 시간을 조금 남긴 뒤 종료를 트리거.
		time.Sleep(1500 * time.Millisecond)
		requestRestart()
	}()

	writeJSON(w, 202, map[string]any{"ok": true, "target": rel.TagName})
}

// updateRollback은 이전 버전으로 능동 롤백(교체 전 백업한 artex.old).
func (s *Server) updateRollback(w http.ResponseWriter, r *http.Request) {
	if _, running := updHub.snapshot(); running {
		writeErr(w, 409, "업데이트가 진행 중이라 롤백할 수 없습니다")
		return
	}
	if err := selfupdate.Rollback(); err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	log.Printf("[update] 이전 버전으로 수동 롤백됨, 전환 완료를 위해 곧 종료")
	writeJSON(w, 202, map[string]any{"ok": true})
	go func() {
		time.Sleep(500 * time.Millisecond)
		requestRestart()
	}()
}

// updateStream은 SSE로 업데이트 진행 상황을 푸시.
func (s *Server) updateStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, 500, "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch, unsub := updHub.subscribe()
	defer unsub()

	send := func(p updateProgress) {
		b, _ := json.Marshal(p)
		fmt.Fprintf(w, "data: %s\n\n", b)
		flusher.Flush()
	}
	// 먼저 현재 상태 하나를 보내, 페이지 새로고침 후 진행 중인 업그레이드를 즉시 볼 수 있게.
	cur, _ := updHub.snapshot()
	send(cur)

	ctx := r.Context()
	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case p, ok := <-ch:
			if !ok {
				return
			}
			send(p)
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		}
	}
}
