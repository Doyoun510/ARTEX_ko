// Package selfupdate implements ARTEX의 페이지 원클릭 업데이트: GitHub Release에서 새 버전
// 바이너리를 가져와 검증·스테이징하고, 다음 시작 시 원자적으로 교체한다.
//
// 전체 역할 분담(start.sh / start.bat 참조):
//
//	시작 스크립트 = 단순 감시 루프, "프로세스 종료 후 종료 코드에 따라 재기동 여부 결정"만 담당
//	이 패키지   = 오류 나기 쉬운 모든 로직(다운로드 / SHA256 검증 / 스모크 / 교체 / 실패 롤백)
//
// 교체를 스크립트가 아니라 Go에 둔 이유는, sha256 검증과 스모크 테스트를 sh와 bat에서
// 두 벌(sha256sum / shasum / certutil)로 작성해야 하는데, 이것이 바로 가장 틀려선 안 되는 부분이기 때문이다 —— 실행되지
// 않는 바이너리로 교체하면 감시 프로세스가 충실히 반복해서 기동하고, 사용자는 기계에 직접 올라가 수동으로 복구할 수밖에 없다.
//
// 한 번의 완전한 업그레이드는 세 번의 프로세스 시작을 거친다:
//
//	① 구버전 server가 /api/update/apply 수신 → 다운로드·검증 → artex.new 스테이징 → exit 75
//	② 스크립트가 구버전 재기동 → Bootstrap이 artex.new 발견 → 검증+스모크 → 교체 → exit 75
//	③ 스크립트 재기동, 이때는 이미 신버전 → Bootstrap이 시도 1회 기록 → 시작 성공 후 마커 제거
//
// 어느 단계든 실패하면 구버전으로 되돌린다: ② 검증 실패 시 스테이징 파일을 삭제하고 구버전을 계속 실행; ③ 마커 제거까지
// 세 번 연속 살아남지 못하면(기동 실패로 크래시) artex.old를 자동으로 되돌린다.
package selfupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ExitRestart는 "감시 프로세스가 나를 재기동"하라는 종료 코드(EX_TEMPFAIL)다. 시작 스크립트가 이를 보면
// 즉시 재실행하고 크래시 백오프에 포함하지 않는다. 0은 사용자 정상 종료(스크립트 루프 탈출), 나머지는 전부 크래시로 간주한다.
const ExitRestart = 75

// maxAttempts는 교체 후 허용되는 시작 시도 횟수다. 신버전은 시작할 때마다 카운트를 +1 하고,
// settleDelay를 넘겨 살아남으면 마커를 제거한다; maxAttempts회 연속 크래시면 신버전이 아예 기동 못 하는 것이므로 자동 롤백한다.
const maxAttempts = 3

// Paths는 한 번의 업그레이드에 관여하는 모든 파일로, 전부 **실행 파일이 있는 디렉터리** 아래에 둔다.
// CWD를 일부러 쓰지 않는다: 서비스로 실행할 때 작업 디렉터리가 / 이거나 임의 경로일 수 있어, CWD를 쓰면 스테이징 파일이
// 다른 곳에 떨어져 교체 로직이 바로 무효화된다.
type Paths struct {
	Dir     string // 실행 파일이 있는 디렉터리
	Current string // 현재 실행 중인 바이너리        artex      / artex.exe
	New     string // 스테이징된 신버전            artex.new  / artex.new.exe
	Sum     string // 신버전의 sha256(hex)  artex.new.sha256 / artex.new.exe.sha256
	Old     string // 교체 전 백업한 구버전      artex.old  / artex.old.exe
	Marker  string // 업그레이드 상태 마커            artex.upgrade.json
}

// ResolvePaths는 현재 실행 파일을 기준으로 모든 업그레이드 경로를 도출한다.
//
// Windows에서 .new/.old도 반드시 .exe 접미사를 가져야 하며, 그렇지 않으면 스모크 테스트와 교체 후 실행이 모두 실패한다.
// 그래서 먼저 접미사를 떼고 다시 조합해야 두 플랫폼의 명명이 대칭이 된다.
func ResolvePaths() (Paths, error) {
	exe, err := os.Executable()
	if err != nil {
		return Paths{}, fmt.Errorf("실행 파일 위치 확인: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	dir := filepath.Dir(exe)
	name := filepath.Base(exe)
	ext := filepath.Ext(name) // Windows에서는 ".exe", Unix에서는 보통 비어 있음
	stem := strings.TrimSuffix(name, ext)

	join := func(suffix string) string { return filepath.Join(dir, stem+suffix+ext) }
	return Paths{
		Dir:     dir,
		Current: exe,
		New:     join(".new"),
		Sum:     join(".new") + ".sha256",
		Old:     join(".old"),
		Marker:  filepath.Join(dir, stem+".upgrade.json"),
	}, nil
}

// marker는 한 번의 교체 진행 상황을 기록하며, 신버전이 기동하지 못할 때 자동 롤백을 트리거한다.
type marker struct {
	From     string `json:"from"`     // 업그레이드 전 버전
	To       string `json:"to"`       // 목표 버전
	Attempts int    `json:"attempts"` // 교체 후 시작을 시도한 횟수
	StagedAt int64  `json:"staged_at"`
}

func readMarker(path string) (marker, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return marker{}, false
	}
	var m marker
	if json.Unmarshal(b, &m) != nil {
		return marker{}, false
	}
	return m, true
}

func writeMarker(path string, m marker) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// cleanStaged는 스테이징 파일을 제거한다. 교체 성공·검증 실패·사용자 취소 모두 이것을 거쳐, 잔존한
// artex.new가 다음 시작 때 다시 시도되는 것을 방지한다.
func cleanStaged(p Paths) {
	_ = os.Remove(p.New)
	_ = os.Remove(p.Sum)
}

// CompareVersions는 두 버전 번호를 비교해 -1/0/1(a<b / a==b / a>b)을 반환한다.
// ok=false는 적어도 한쪽이 비교 가능한 버전 번호가 아님을 뜻한다(예: 로컬 개발 빌드의 "dev" 또는
// git describe가 산출한 "0.3.7-2-gabc1234-dirty"). 이때 호출자는 원클릭 업데이트를 비활성화해야 하며,
// 그렇지 않으면 개발 중인 빌드를 정식 버전으로 "업그레이드"해 커밋하지 않은 변경을 덮어쓴다.
func CompareVersions(a, b string) (int, bool) {
	av, aok := parseVersion(a)
	bv, bok := parseVersion(b)
	if !aok || !bok {
		return 0, false
	}
	for i := range 3 {
		if av[i] != bv[i] {
			if av[i] < bv[i] {
				return -1, true
			}
			return 1, true
		}
	}
	return 0, true
}

// parseVersion은 "v0.3.7" / "0.3.7" 형식의 버전 번호를 [3]int로 파싱한다.
//
// 순수한 3단 형식만 허용한다: build.sh가 비-tag 빌드에서 git describe로 산출하는
// "0.3.7-2-gabc1234" 같은 접미사 붙은 버전은 반드시 비교 불가로 판정되어야 하고,
// 0.3.7 —— 그렇지 않으면 개발 빌드가 "이미 최신"으로 오판되거나 정식 버전으로 덮어쓰인다.
func parseVersion(s string) ([3]int, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	if s == "" {
		return [3]int{}, false
	}
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return [3]int{}, false
	}
	var out [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, false
		}
		out[i] = n
	}
	return out, true
}

// InDocker는 프로세스가 컨테이너 안에서 실행 중인지 보고한다. Docker에서는 교체가 컨테이너 쓰기 가능 레이어에 기록되어,
// `docker compose up -d`로 컨테이너를 재생성하면 이미지에 내장된 버전으로 되돌아간다 —— 이는 예상된 동작이며
// (그때는 사용자가 애초에 새 이미지를 받는 중이다), 프런트가 이에 맞춰 명확히 안내할 수 있어야 한다.
func InDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	b, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	s := string(b)
	return strings.Contains(s, "docker") || strings.Contains(s, "containerd")
}
