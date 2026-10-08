package selfupdate

import (
	"archive/zip"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// testPaths는 격리된 업그레이드 디렉터리를 만든다. ResolvePaths()를 직접 쓰면 안 된다 —— 그건 테스트
// 바이너리 자체를 가리켜, 실행하면 go test의 실행 파일 이름을 바꿔버린다.
func testPaths(t *testing.T) Paths {
	t.Helper()
	dir := t.TempDir()
	return Paths{
		Dir:     dir,
		Current: filepath.Join(dir, "artex"),
		New:     filepath.Join(dir, "artex.new"),
		Sum:     filepath.Join(dir, "artex.new.sha256"),
		Old:     filepath.Join(dir, "artex.old"),
		Marker:  filepath.Join(dir, "artex.upgrade.json"),
	}
}

// fakeBin은 artex를 가장하는 실행 가능한 셸 스크립트를 작성한다. smokeTest는 단지 -h로 기동해 종료 코드만 보므로,
// 스크립트로 충분하고, 진짜 바이너리를 컴파일하는 것보다 훨씬 빠르다.
func fakeBin(t *testing.T, path, marker string, exitCode int) {
	t.Helper()
	script := "#!/bin/sh\necho " + marker + "\nexit " + itoa(exitCode) + "\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("가짜 바이너리 쓰기 %s: %v", path, err)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	return string(rune('0' + n))
}

// stage는 bin을 "스테이징되어 교체 대기" 상태로 배치한다: artex.new와 그 체크섬을 작성한다.
func stage(t *testing.T, p Paths, marker string, exitCode int) {
	t.Helper()
	fakeBin(t, p.New, marker, exitCode)
	sum, err := fileSHA256(p.New)
	if err != nil {
		t.Fatalf("체크섬 계산: %v", err)
	}
	if err := os.WriteFile(p.Sum, []byte(sum), 0o644); err != nil {
		t.Fatalf("체크섬 쓰기: %v", err)
	}
}

func readAll(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s 읽기: %v", path, err)
	}
	return string(b)
}

func requireUnix(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("가짜 바이너리는 sh 스크립트라 Windows에서 실행 안 됨")
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b       string
		want       int
		comparable bool
	}{
		{"0.3.7", "0.3.8", -1, true},
		{"0.3.8", "0.3.7", 1, true},
		{"0.3.7", "0.3.7", 0, true},
		{"v0.3.7", "0.3.8", -1, true}, // build.sh는 v 제거, tag는 v 포함, 양쪽 다 인정해야 함
		{"0.3.7", "v0.3.7", 0, true},
		{"0.9.0", "0.10.0", -1, true}, // 사전순이 아니라 숫자로 비교
		{"1.0.0", "0.99.99", 1, true},
		// 개발 빌드는 비교 불가로 판정되어야 한다. 그렇지 않으면 정식 버전이 커밋하지 않은 변경을 덮어쓴다.
		{"dev", "0.3.8", 0, false},
		{"0.3.7-2-gabc1234", "0.3.8", 0, false},
		{"0.3.7-dirty", "0.3.8", 0, false},
		{"0.3", "0.3.8", 0, false},
		{"", "0.3.8", 0, false},
	}
	for _, c := range cases {
		got, ok := CompareVersions(c.a, c.b)
		if ok != c.comparable {
			t.Errorf("CompareVersions(%q,%q) comparable=%v, 기대 %v", c.a, c.b, ok, c.comparable)
			continue
		}
		if ok && got != c.want {
			t.Errorf("CompareVersions(%q,%q)=%d, 기대 %d", c.a, c.b, got, c.want)
		}
	}
}

func TestResolvePathsNaming(t *testing.T) {
	p, err := ResolvePaths()
	if err != nil {
		t.Fatalf("ResolvePaths: %v", err)
	}
	// 핵심 불변식: 모든 업그레이드 파일은 실행 파일과 같은 디렉터리에 있다. CWD에 떨어지면 서비스로 실행
	// (작업 디렉터리가 /일 수 있음)할 때 교체가 완전히 무효화된다.
	for name, path := range map[string]string{"New": p.New, "Sum": p.Sum, "Old": p.Old, "Marker": p.Marker} {
		if filepath.Dir(path) != p.Dir {
			t.Errorf("%s 경로가 실행 파일 디렉터리 아래에 없음: %s (기대 %s)", name, path, p.Dir)
		}
	}
	// Windows에서 .new/.old는 반드시 .exe를 유지해야 하며, 그렇지 않으면 스모크 테스트와 교체 후 실행이 모두 실패한다.
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(p.New, ".exe") || !strings.HasSuffix(p.Old, ".exe") {
			t.Errorf("Windows에서 .new/.old는 반드시 .exe로 끝나야 함: new=%s old=%s", p.New, p.Old)
		}
	}
}

func TestVerifyStagedRejectsTamperedBinary(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	stage(t, p, "new", 0)

	// 체크섬을 쓴 뒤 파일을 변경해, 다운로드 손상 / 바꿔치기를 시뮬레이션한다.
	fakeBin(t, p.New, "tampered", 0)
	if err := verifyStaged(p); err == nil {
		t.Fatal("SHA256 불일치가 거부되길 기대했으나 통과됨")
	}
}

func TestVerifyStagedRejectsUnrunnableBinary(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	stage(t, p, "broken", 1) // 실행은 되지만 종료 코드가 0이 아님

	if err := verifyStaged(p); err == nil {
		t.Fatal("스모크 테스트 실패가 거부되길 기대했으나 통과됨")
	}
}

func TestApplyStagedHappyPath(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "old", 0)
	stage(t, p, "new", 0)
	if err := writeMarker(p.Marker, marker{From: "0.3.7", To: "0.3.8"}); err != nil {
		t.Fatalf("마커 쓰기: %v", err)
	}

	action, st := applyStaged(p)
	if action != Restart {
		t.Fatalf("Restart 기대, %v 얻음", action)
	}
	if !st.Pending {
		t.Error("교체 후 상태는 Pending이어야 함")
	}
	if !strings.Contains(readAll(t, p.Current), "new") {
		t.Error("artex는 신버전으로 교체되어 있어야 함")
	}
	if !strings.Contains(readAll(t, p.Old), "old") {
		t.Error("구버전은 artex.old로 백업되어야 함")
	}
	if _, err := os.Stat(p.New); !os.IsNotExist(err) {
		t.Error("교체 후 artex.new는 사라져 있어야 함")
	}
	if _, err := os.Stat(p.Sum); !os.IsNotExist(err) {
		t.Error("교체 후 체크섬 파일은 정리되어 있어야 함")
	}
	// 마커는 반드시 남아 있어야 한다. 다음 시작(신버전 실행)이 이것으로 카운트하고 필요 시 롤백한다.
	if _, ok := readMarker(p.Marker); !ok {
		t.Error("교체 후 업그레이드 마커는 보존되어야 함")
	}
}

func TestApplyStagedKeepsCurrentWhenVerifyFails(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "old", 0)
	stage(t, p, "new", 0)
	fakeBin(t, p.New, "tampered", 0) // 체크섬 손상

	action, st := applyStaged(p)
	if action != Continue {
		t.Fatalf("검증 실패 시 Continue 기대, %v 얻음", action)
	}
	if !st.FailedStage {
		t.Error("상태는 FailedStage로 표시되어야 함")
	}
	if !strings.Contains(readAll(t, p.Current), "old") {
		t.Fatal("검증 실패 시 현재 버전을 절대 건드리면 안 됨")
	}
	if _, err := os.Stat(p.New); !os.IsNotExist(err) {
		t.Error("검증 실패한 스테이징 파일은 정리되어야 함, 그렇지 않으면 다음 시작 때 다시 시도")
	}
}

func TestSwapOverwritesPreviousBackup(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "v2", 0)
	fakeBin(t, p.Old, "v1", 0) // 지난 업그레이드가 남긴 백업
	stage(t, p, "v3", 0)

	if err := swap(p); err != nil {
		t.Fatalf("swap: %v", err)
	}
	if !strings.Contains(readAll(t, p.Current), "v3") {
		t.Error("v3로 교체되어야 함")
	}
	if !strings.Contains(readAll(t, p.Old), "v2") {
		t.Error("백업은 방금 교체된 v2로 갱신되어야 함")
	}
}

func TestConfirmCountsAttemptsThenRollsBack(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "broken-new", 0)
	fakeBin(t, p.Old, "good-old", 0)
	m := marker{From: "0.3.7", To: "0.3.8"}

	// 처음 maxAttempts회 시작은 카운트만 누적해, 신버전이 스스로 자리 잡을 기회를 준다.
	for i := 1; i <= maxAttempts; i++ {
		action, st := confirmOrRollback(p, m)
		if action != Continue {
			t.Fatalf("%d번째 시도는 Continue 기대, %v 얻음", i, action)
		}
		if !st.Pending {
			t.Errorf("%d번째 시도 상태는 Pending이어야 함", i)
		}
		got, ok := readMarker(p.Marker)
		if !ok || got.Attempts != i {
			t.Fatalf("%d번째 시도 후 attempts=%d(ok=%v), 기대 %d", i, got.Attempts, ok, i)
		}
		m = got
	}

	// 한 번 더 크래시하면 한도를 초과해, 구버전을 자동으로 되돌린다.
	action, st := confirmOrRollback(p, m)
	if action != Restart {
		t.Fatalf("시도 한도 초과 시 Restart 기대, %v 얻음", action)
	}
	if !st.RolledBack {
		t.Error("상태는 RolledBack으로 표시되어야 함")
	}
	if !strings.Contains(readAll(t, p.Current), "good-old") {
		t.Fatal("구버전으로 롤백되어 있어야 함")
	}
	if _, err := os.Stat(p.Marker); !os.IsNotExist(err) {
		t.Error("롤백 후 마커는 제거되어야 함, 그렇지 않으면 무한 롤백")
	}
	// 기동하지 못한 그 버전은 조사용으로 남기고, 바로 삭제하지 않는다.
	if _, err := os.Stat(p.Current + ".failed"); err != nil {
		t.Error("실패한 버전은 조사용으로 .failed로 보존되어야 함")
	}
}

func TestManualRollbackIsReversible(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "v2", 0)
	fakeBin(t, p.Old, "v1", 0)

	// Rollback()은 ResolvePaths()를 거치므로, 여기서는 하위 교환 의미를 직접 테스트한다.
	tmp := p.Current + ".swap"
	if err := os.Rename(p.Current, tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(p.Old, p.Current); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp, p.Old); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readAll(t, p.Current), "v1") {
		t.Error("롤백 후 현재 버전은 v1이어야 함")
	}
	if !strings.Contains(readAll(t, p.Old), "v2") {
		t.Error("롤백 후 백업은 v2가 되어야 함, 이렇게 해야 다시 되돌릴 수 있음")
	}
}

func TestParseSums(t *testing.T) {
	const (
		linuxSum = "1111111111111111111111111111111111111111111111111111111111111111"
		winSum   = "ABCDEF0000000000000000000000000000000000000000000000000000000000"
	)
	// sha256sum 출력은 이중 공백 구분; shasum -a 256은 바이너리 모드에서 파일명에 *를 붙인다.
	raw := linuxSum + "  artex-0.3.8-linux-amd64.zip\n" +
		winSum + " *artex-0.3.8-windows-amd64.zip\n" +
		"\n" +
		"garbage line\n" + // 정확히 두 필드지만 첫 번째가 다이제스트가 아님
		"deadbeef  artex-0.3.8-darwin-arm64.zip\n" // 다이제스트 길이가 틀림

	out := parseSums(raw)
	if out["artex-0.3.8-linux-amd64.zip"] != linuxSum {
		t.Errorf("linux 항목 파싱 오류: %v", out)
	}
	// 다이제스트는 소문자로 통일해, 대조 시 대소문자 때문에 불일치로 오판하지 않도록 한다.
	if got := out["artex-0.3.8-windows-amd64.zip"]; got != strings.ToLower(winSum) {
		t.Errorf("windows 항목 오류(* 접두는 제거되고 다이제스트는 소문자로 변환되어야 함): %q", got)
	}
	if len(out) != 2 {
		t.Errorf("빈 줄·비 다이제스트 줄·길이가 틀린 줄은 무시되어야 함, %v 얻음", out)
	}
}

func TestExtractBinaryFindsNestedEntry(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("패키지 내 기본 이름은 Windows에서 artex.exe라, 이 케이스는 Unix 명명으로 구성")
	}
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "release.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	// 실제 릴리스 패키지 구조: artex-<버전>-<os>-<arch>/artex, 추가로 몇 개의 방해 파일.
	for name, body := range map[string]string{
		"artex-0.3.8-linux-amd64/README.md":           "readme",
		"artex-0.3.8-linux-amd64/skills/a.md":         "skill",
		"artex-0.3.8-linux-amd64/artex":               "#!/bin/sh\nexit 0\n",
		"artex-0.3.8-linux-amd64/config.example.json": "{}",
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()

	dst := filepath.Join(dir, "out")
	if err := extractBinary(zipPath, dst); err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	if got := readAll(t, dst); !strings.Contains(got, "exit 0") {
		t.Errorf("압축 해제된 것이 artex 실행 파일이 아님: %q", got)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Error("압축 해제된 바이너리는 실행 비트를 가져야 함")
	}
}

func TestExtractBinaryMissingEntry(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "release.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, _ := zw.Create("artex-0.3.8-linux-amd64/README.md")
	_, _ = w.Write([]byte("readme"))
	_ = zw.Close()
	f.Close()

	if err := extractBinary(zipPath, filepath.Join(dir, "out")); err == nil {
		t.Fatal("패키지에 실행 파일이 없으면 오류를 내야 함")
	}
}

func TestCheckURLRejectsNonGitHub(t *testing.T) {
	bad := []string{
		"http://github.com/x",           // 비 HTTPS
		"https://evil.com/artex.zip",    // 도메인이 화이트리스트에 없음
		"https://github.com.evil.com/x", // 접미사 위장
		"https://raw.githubusercontent.com.evil.com/x",
	}
	for _, raw := range bad {
		u := mustParse(t, raw)
		if err := checkURL(u); err == nil {
			t.Errorf("checkURL(%q)는 거부해야 함", raw)
		}
	}
	good := []string{
		"https://api.github.com/repos/x/releases/latest",
		"https://objects.githubusercontent.com/blah",
		"https://GitHub.com/x", // 도메인 대소문자 구분 안 함
	}
	for _, raw := range good {
		u := mustParse(t, raw)
		if err := checkURL(u); err != nil {
			t.Errorf("checkURL(%q)는 허용해야 하는데 오류: %v", raw, err)
		}
	}
}

func TestAssetNameMatchesBuildScript(t *testing.T) {
	// build.sh의 package_binary는 artex-<버전>-<os>-<arch>.zip을 사용하고, 버전 번호는
	// v 접두를 제거한다. 여기서 한 글자만 틀려도 모든 플랫폼의 원클릭 업데이트가 자산을 찾지 못한다.
	if got := AssetName("v0.3.8", "linux", "amd64"); got != "artex-0.3.8-linux-amd64.zip" {
		t.Errorf("AssetName = %q", got)
	}
	if got := AssetName("0.3.8", "windows", "amd64"); got != "artex-0.3.8-windows-amd64.zip" {
		t.Errorf("AssetName = %q", got)
	}
}

func mustParse(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%q 파싱: %v", raw, err)
	}
	return u
}

func TestSettleClearsMarkerAndStopsRollback(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "new", 0)
	fakeBin(t, p.Old, "old", 0)
	if err := writeMarker(p.Marker, marker{From: "0.3.7", To: "0.3.8", Attempts: 2}); err != nil {
		t.Fatal(err)
	}

	settle(p)

	if _, err := os.Stat(p.Marker); !os.IsNotExist(err) {
		t.Fatal("안정 확인 후 업그레이드 마커는 반드시 제거되어야 함")
	}
	// 마커가 사라지면, 이후 정상 재시작은 더 이상 횟수를 누적하지 않고 롤백을 잘못 트리거하지도 않는다.
	if _, ok := readMarker(p.Marker); ok {
		t.Error("마커 읽기는 실패해야 함")
	}
	// 백업은 남겨 두어야 사용자가 수동 롤백할 수 있다.
	if _, err := os.Stat(p.Old); err != nil {
		t.Error("안정 확인 후에도 이전 버전 백업은 보존되어야 함")
	}
}

func TestSettleIsNoopWithoutMarker(t *testing.T) {
	requireUnix(t)
	p := testPaths(t)
	fakeBin(t, p.Current, "cur", 0)
	settle(p) // 일반 시작 경로, panic도 안 되고 어떤 파일도 건드리면 안 됨
	if _, err := os.Stat(p.Current); err != nil {
		t.Error("마커가 없을 때 settle은 어떤 파일에도 영향을 주면 안 됨")
	}
}
