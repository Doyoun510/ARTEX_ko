package selfupdate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// smokeEnv는 스모크 테스트가 기동한 자식 프로세스가 Bootstrap을 바로 건너뛰게 한다.
//
// 엄밀히는 없어도 문제없다: 자식 프로세스의 os.Executable()은 artex.new라, 도출된
// 모든 경로가 .new 접두를 가져 진짜 업그레이드 파일에 닿지 않는다. 하지만 이런 우연에 의존하는 건 너무 취약하고,
// 명시적 단락(short-circuit)이 한눈에 보이며, 자식 프로세스의 불필요한 디스크 탐색도 아낀다.
const smokeEnv = "ARTEX_SELFUPDATE_SMOKE"

// Action은 Bootstrap이 main에게 주는 지시다.
type Action int

const (
	// Continue: 평소대로 server를 시작한다.
	Continue Action = iota
	// Restart: 즉시 ExitRestart로 종료해 감시 스크립트가 재기동하게 한다.
	Restart
)

// State는 이번 시작 시의 업그레이드 상태를 기술해, /api/update/check가 프런트에
// "지난 업그레이드가 성공했는지 롤백됐는지"를 사실대로 알리도록 한다.
type State struct {
	Pending     bool   // 교체 후 아직 안정 확인 안 됨
	RolledBack  bool   // 이번 시작에서 방금 자동 롤백을 수행함
	FailedStage bool   // 스테이징 파일 검증/스모크 미통과, 폐기됨
	Detail      string // 사용자 대상 한 줄 설명
}

// Bootstrap은 main 맨 앞에서 실행되며, 어떤 포트 리스닝이나 데이터베이스 열기보다 먼저 호출해야 한다.
//
// 세 가지 상황:
//
//	① 스테이징 파일 artex.new 존재  → 검증 + 스모크, 통과하면 교체하고 재시작 요청; 미통과면 폐기하고 구버전 계속 실행
//	② 마커 파일만 남음          → 방금 교체 완료를 의미, 시도 1회 누적; 연속 실패가 충분히 많으면 롤백
//	③ 아무것도 없음            → 정상 시작
func Bootstrap() (Action, State) {
	if os.Getenv(smokeEnv) != "" {
		return Continue, State{}
	}
	p, err := ResolvePaths()
	if err != nil {
		log.Printf("[update] 부트스트랩 건너뜀: %v", err)
		return Continue, State{}
	}

	if _, err := os.Stat(p.New); err == nil {
		return applyStaged(p)
	}

	m, ok := readMarker(p.Marker)
	if !ok {
		return Continue, State{}
	}
	return confirmOrRollback(p, m)
}

// applyStaged는 "스테이징 파일 존재" 상황을 처리한다: 검증 통과하면 교체하고, 실패하면 폐기한다.
//
// 여기가 전체 업그레이드 경로에서 실행 파일을 덮어쓰는 유일한 곳이자 마지막 관문이다 —— 스모크 테스트가
// 다운로드 손상·아키텍처 오선택·동적 링크 누락 같은 문제를 막는다. 실행되지 않는 바이너리를 한 번 통과시키면,
// 감시 스크립트가 지치지 않고 반복해서 기동하는데 Go 코드는 아예 실행될 기회가 없어, 자동 롤백도 불가능해진다.
func applyStaged(p Paths) (Action, State) {
	m, _ := readMarker(p.Marker)

	if err := verifyStaged(p); err != nil {
		log.Printf("[update] 스테이징된 신버전이 검증을 통과하지 못해 폐기하고 현재 버전을 계속 실행: %v", err)
		cleanStaged(p)
		_ = os.Remove(p.Marker)
		return Continue, State{FailedStage: true, Detail: "신버전 검증 실패로 폐기함: " + err.Error()}
	}

	if err := swap(p); err != nil {
		log.Printf("[update] 교체 실패로 현재 버전을 계속 실행: %v", err)
		cleanStaged(p)
		_ = os.Remove(p.Marker)
		return Continue, State{FailedStage: true, Detail: "교체 실패: " + err.Error()}
	}

	// 교체 성공. 마커를 유지해, 다음 시작(바로 신버전 실행)에 안정 여부 확인을 맡긴다.
	m.Attempts = 0
	if m.StagedAt == 0 {
		m.StagedAt = time.Now().Unix()
	}
	if err := writeMarker(p.Marker, m); err != nil {
		log.Printf("[update] 업그레이드 마커 쓰기 실패(자동 롤백 능력 상실): %v", err)
	}
	log.Printf("[update] %s로 교체됨, 재시작을 위해 종료(exit %d)", orUnknown(m.To), ExitRestart)
	return Restart, State{Pending: true}
}

// confirmOrRollback은 "교체 후의 시작"을 처리한다: 시도 횟수를 누적하고, 한도를 넘으면 구버전으로 되돌린다.
//
// 카운트는 Go 코드가 실행된 뒤에만 증가하므로, "실행은 되지만 초기화 때 크래시"
// (설정 비호환·포트 점유·DB 마이그레이션 터짐)하는 류의 장애를 담당한다; "exec 자체가 불가"는 교체 전
// 스모크 테스트가 막으며, 둘을 합쳐야 완전해진다.
func confirmOrRollback(p Paths, m marker) (Action, State) {
	m.Attempts++
	if m.Attempts > maxAttempts {
		if err := rollback(p); err != nil {
			// 롤백까지 실패했으면 더는 재시작하지 않는다, 그렇지 않으면 무한 재시작에 빠진다. 마커를 제거해,
			// 프로세스를 현재 상태로 기동하게 한다 —— 기동하지 못하면 사용자가 최소한 로그에서 원인을 볼 수 있다.
			log.Printf("[update] 신버전이 %d회 연속 시작에 실패하고 롤백도 실패: %v", maxAttempts, err)
			_ = os.Remove(p.Marker)
			return Continue, State{Detail: "신버전 시작 실패 및 롤백 실패: " + err.Error()}
		}
		log.Printf("[update] 신버전이 %d회 연속 시작에 실패해 %s로 롤백, 재시작을 위해 종료(exit %d)",
			maxAttempts, orUnknown(m.From), ExitRestart)
		_ = os.Remove(p.Marker)
		return Restart, State{RolledBack: true, Detail: fmt.Sprintf("신버전 시작 실패로 %s로 롤백함", orUnknown(m.From))}
	}
	if err := writeMarker(p.Marker, m); err != nil {
		log.Printf("[update] 업그레이드 마커 갱신 실패: %v", err)
	}
	log.Printf("[update] 신버전 시작 중(%d/%d번째 시도), 안정적으로 실행되면 업그레이드를 확정",
		m.Attempts, maxAttempts)
	return Continue, State{Pending: true}
}

// Settle은 신버전이 안정적으로 실행됨을 확인하고 업그레이드 마커를 제거한다.
//
// main이 HTTP 리스닝이 뜬 뒤 지연 호출한다: 이 시간을 넘겨 살아남아야 유효하며, 그렇지 않으면 마커가 그대로 남아
// 다음 시작에서 시도 횟수를 계속 누적해 롤백이 트리거될 때까지 이어진다.
func Settle() {
	p, err := ResolvePaths()
	if err != nil {
		return
	}
	settle(p)
}

func settle(p Paths) {
	if _, ok := readMarker(p.Marker); !ok {
		return // 업그레이드 후 시작이 아니라 할 일 없음
	}
	if err := os.Remove(p.Marker); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("[update] 업그레이드 마커 제거 실패: %v", err)
		return
	}
	log.Printf("[update] 신버전이 안정적으로 실행, 업그레이드 완료(이전 버전은 %s로 보존)", p.Old)
}

// SettleDelay는 "신버전이 살아남았다"고 판정하는 데 필요한 실행 시간이다.
const SettleDelay = 30 * time.Second

// verifyStaged는 스테이징 파일을 검증한다: 먼저 SHA256을 대조하고, 다음 실제로 기동해 한 번 실행해 본다.
func verifyStaged(p Paths) error {
	want, err := os.ReadFile(p.Sum)
	if err != nil {
		return fmt.Errorf("체크섬 읽기: %w", err)
	}
	got, err := fileSHA256(p.New)
	if err != nil {
		return fmt.Errorf("체크섬 계산: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(string(want)), got) {
		return errors.New("SHA256 불일치(다운로드 손상 또는 변조)")
	}
	return smokeTest(p.New)
}

// smokeTest는 -h로 새 바이너리를 기동해, 현재 시스템에서 실제로 실행되는지 확인한다.
// 이는 다운로드 절단·아키텍처 오선택(exec format error)·의존성 누락 같은 큰 부류의 문제를 막는다.
func smokeTest(bin string) error {
	if err := os.Chmod(bin, 0o755); err != nil {
		return fmt.Errorf("실행 권한 부여: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, "-h")
	cmd.Env = append(os.Environ(), smokeEnv+"=1")
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return errors.New("스모크 테스트 타임아웃(새 바이너리 무응답)")
	}
	if err != nil {
		snippet := strings.TrimSpace(string(out))
		if len(snippet) > 300 {
			snippet = snippet[:300] + "…"
		}
		return fmt.Errorf("스모크 테스트 실패: %v: %s", err, snippet)
	}
	return nil
}

// swap은 현재 바이너리를 스테이징된 신버전으로 교체한다.
//
// Unix와 Windows 모두 실행 중인 실행 파일의 rename을 허용한다(Windows가 금지하는 건 삭제와
// 덮어쓰기이고, rename은 거기 없다). 그래서 여기선 플랫폼을 나눌 필요도, 자기 자신을 먼저 멈출 필요도 없다.
func swap(p Paths) error {
	// Windows의 rename은 이미 존재하는 대상을 덮어쓰지 않으므로, 지난 업그레이드가 남긴 .old를 먼저 제거해야 한다.
	if err := os.Remove(p.Old); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("구 백업 정리 %s: %w", p.Old, err)
	}
	if err := os.Rename(p.Current, p.Old); err != nil {
		return fmt.Errorf("현재 버전 백업: %w", err)
	}
	if err := os.Rename(p.New, p.Current); err != nil {
		// 교체는 실패했지만 현재 버전은 이미 옮겨졌으므로, 원래대로 되돌려야 한다. 그렇지 않으면 다음 시작 때 실행 파일이 없다.
		if rerr := os.Rename(p.Old, p.Current); rerr != nil {
			return fmt.Errorf("신버전 설치 실패(%v), 게다가 현재 버전 복구도 실패: %w", err, rerr)
		}
		return fmt.Errorf("신버전 설치: %w", err)
	}
	_ = os.Remove(p.Sum)
	return nil
}

// rollback은 swap이 백업한 구버전을 되돌린다.
func rollback(p Paths) error {
	if _, err := os.Stat(p.Old); err != nil {
		return fmt.Errorf("롤백할 백업이 없음 %s: %w", p.Old, err)
	}
	// 기동하지 못한 신버전을 .failed로 옮겨 조사용으로 남기고, 바로 삭제하지 않는다.
	failed := p.Current + ".failed"
	_ = os.Remove(failed)
	if err := os.Rename(p.Current, failed); err != nil {
		return fmt.Errorf("실패한 버전 이동: %w", err)
	}
	if err := os.Rename(p.Old, p.Current); err != nil {
		return fmt.Errorf("구버전 복구: %w", err)
	}
	return nil
}

// Rollback은 /api/update/rollback의 구현이다: 능동적으로 이전 버전으로 되돌린다.
// 교체만 수행하고 재시작은 마찬가지로 감시 스크립트에 맡긴다(호출자는 이어서 ExitRestart로 종료).
func Rollback() error {
	p, err := ResolvePaths()
	if err != nil {
		return err
	}
	if _, err := os.Stat(p.Old); err != nil {
		return errors.New("롤백할 이전 버전이 없음(" + p.Old + " 없음)")
	}
	cleanStaged(p)
	if err := smokeTest(p.Old); err != nil {
		return fmt.Errorf("이전 버전을 실행할 수 없어 롤백을 거부합니다: %w", err)
	}
	// 현재와 백업을 교환한다: 롤백 후 다시 되돌릴 수 있다.
	tmp := p.Current + ".swap"
	_ = os.Remove(tmp)
	if err := os.Rename(p.Current, tmp); err != nil {
		return fmt.Errorf("현재 버전 이동: %w", err)
	}
	if err := os.Rename(p.Old, p.Current); err != nil {
		_ = os.Rename(tmp, p.Current)
		return fmt.Errorf("이전 버전 설치: %w", err)
	}
	if err := os.Rename(tmp, p.Old); err != nil {
		log.Printf("[update] 롤백 후 백업 정리 실패(실행에는 영향 없음): %v", err)
	}
	_ = os.Remove(p.Marker)
	return nil
}

// HasBackup은 롤백할 이전 버전이 있는지 보고해, 프런트가 롤백 버튼 표시 여부를 결정하도록 한다.
func HasBackup() bool {
	p, err := ResolvePaths()
	if err != nil {
		return false
	}
	_, err = os.Stat(p.Old)
	return err == nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func orUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "알 수 없는 버전"
	}
	return s
}
