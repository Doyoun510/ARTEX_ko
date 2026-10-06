package selfupdate

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// sumsAsset은 release.yml이 생성하는 체크섬 목록으로, Release의 모든 zip을 포함한다.
const sumsAsset = "SHA256SUMS"

// maxBinarySize는 압축 해제된 바이너리 크기를 제한해, 기형 zip이 디스크를 가득 채우는 것을 방지한다.
const maxBinarySize = 512 << 20 // 512 MiB

// Phase는 업그레이드 과정의 단계로, SSE 이벤트의 phase 필드로 바로 사용한다.
type Phase string

const (
	PhaseIdle     Phase = "idle"
	PhaseDownload Phase = "downloading"
	PhaseVerify   Phase = "verifying"
	PhaseExtract  Phase = "extracting"
	PhaseStaged   Phase = "staged"
	PhaseFailed   Phase = "failed"
)

// Progress는 호출자가 제공하며, 진행 상황을 프런트에 전달하는 데 쓴다. pct는 다운로드 단계에서만 의미가 있고(0-100),
// 나머지 단계는 -1을 전달한다.
type Progress func(ph Phase, pct int, msg string)

// Stage는 지정한 Release의 현재 플랫폼 패키지를 다운로드하고, 검증 후 새 바이너리를 artex.new로 스테이징한다.
//
// 순수 바이너리가 아니라 완전한 zip을 쓰는데, 이유는 두 가지다: 기존 Release의 SHA256SUMS가 원래
// zip만 포함하므로, zip을 쓰면 CI를 바꿀 필요가 없고 이미 배포된 과거 버전과도 호환된다; zip에는
// skills/도 들어 있어 향후 내장 skill 동기화를 위한 여지를 남긴다. 대가는 skills 몇백 KB를 더 받는 것뿐이다.
//
// 함수가 반환되면 스테이징 완료를 의미하며, 호출자는 이어서 정상 종료하고 ExitRestart로 나간다.
func Stage(ctx context.Context, c *http.Client, rel *Release, currentVersion string, prog Progress) error {
	if prog == nil {
		prog = func(Phase, int, string) {}
	}
	p, err := ResolvePaths()
	if err != nil {
		return err
	}
	if err := checkWritable(p.Dir); err != nil {
		return err
	}

	name := AssetName(rel.TagName, runtime.GOOS, runtime.GOARCH)
	asset, ok := rel.FindAsset(name)
	if !ok {
		return fmt.Errorf("이 버전은 %s/%s 패키지를 제공하지 않습니다(%s 없음)", runtime.GOOS, runtime.GOARCH, name)
	}

	prog(PhaseDownload, 0, "체크섬 목록 가져오는 중…")
	sums, err := fetchSums(ctx, c, rel)
	if err != nil {
		return err
	}
	want, ok := sums[name]
	if !ok {
		return fmt.Errorf("%s에 %s가 수록되지 않아, 검증되지 않은 바이너리 설치를 거부합니다", sumsAsset, name)
	}

	// 임시 파일을 모두 대상 디렉터리에 두어, 마지막 rename이 같은 파일 시스템 내의 원자적 연산이 되도록 보장한다
	// (장치 간 rename은 실패하고, /tmp은 흔히 독립 마운트 지점이다).
	zipPath := p.New + ".zip.part"
	binPath := p.New + ".part"
	defer func() {
		_ = os.Remove(zipPath)
		_ = os.Remove(binPath)
	}()

	prog(PhaseDownload, 0, fmt.Sprintf("%s 다운로드 중(%s)…", name, humanSize(asset.Size)))
	got, err := download(ctx, c, asset, zipPath, prog)
	if err != nil {
		return err
	}

	prog(PhaseVerify, -1, "SHA256 검증 중…")
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("SHA256 불일치: 기대 %s, 실제 %s(다운로드 손상 또는 변조)", short(want), short(got))
	}

	prog(PhaseExtract, -1, "압축 해제 및 스모크 테스트 중…")
	if err := extractBinary(zipPath, binPath); err != nil {
		return err
	}
	if err := smokeTest(binPath); err != nil {
		return fmt.Errorf("신버전이 현재 시스템에서 실행되지 않습니다: %w", err)
	}

	// 스테이징 파일 자체의 sha256을 별도로 저장한다: 다음 시작 때 교체 전에 한 번 더 검증해,
	// 스테이징 후 재시작 전 사이에 파일이 변경되거나 손상되는 것을 방지한다.
	binSum, err := fileSHA256(binPath)
	if err != nil {
		return fmt.Errorf("새 바이너리 체크섬 계산: %w", err)
	}
	if err := os.WriteFile(p.Sum, []byte(binSum), 0o644); err != nil {
		return fmt.Errorf("체크섬 쓰기: %w", err)
	}
	if err := os.Rename(binPath, p.New); err != nil {
		_ = os.Remove(p.Sum)
		return fmt.Errorf("신버전 스테이징: %w", err)
	}

	if err := writeMarker(p.Marker, marker{
		From:     currentVersion,
		To:       strings.TrimPrefix(rel.TagName, "v"),
		StagedAt: time.Now().Unix(),
	}); err != nil {
		// 마커는 자동 롤백 능력에만 영향을 주고, 스테이징 파일 자체는 이미 자리 잡았으므로 이 때문에 업그레이드를 중단하지 않는다.
		prog(PhaseStaged, -1, "경고: 업그레이드 마커 쓰기 실패, 이번 업그레이드는 자동 롤백 보호가 없습니다")
	}

	prog(PhaseStaged, 100, "신버전 준비 완료, 재시작 중…")
	return nil
}

// fetchSums는 SHA256SUMS를 다운로드·파싱해, 파일명 → 16진수 다이제스트를 반환한다.
func fetchSums(ctx context.Context, c *http.Client, rel *Release) (map[string]string, error) {
	asset, ok := rel.FindAsset(sumsAsset)
	if !ok {
		return nil, fmt.Errorf("이 Release에 %s가 없어 무결성을 검증할 수 없으므로 업그레이드를 거부합니다", sumsAsset)
	}
	body, err := get(ctx, c, asset.URL)
	if err != nil {
		return nil, fmt.Errorf("%s 다운로드: %w", sumsAsset, err)
	}
	defer body.Close()

	raw, err := io.ReadAll(io.LimitReader(body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%s 읽기: %w", sumsAsset, err)
	}
	out := parseSums(string(raw))
	if len(out) == 0 {
		return nil, fmt.Errorf("%s 내용이 비었거나 형식을 인식할 수 없습니다", sumsAsset)
	}
	return out, nil
}

// parseSums는 sha256sum 스타일 목록을 파싱해, 파일명 → 16진수 다이제스트를 반환한다.
//
// 첫 번째 필드가 64자리 16진수여야만 수록한다. 단지 "정확히 두 필드"로만 판단하는 것은 부족하다 ——
// 두 단어짜리 설명 문구가 있는 아무 줄이나 유효 항목으로 간주되어 쓰레기 값을 다이제스트 표에 끼워 넣고,
// 진짜 자산이 오히려 잘못된 다이제스트에 매칭될 수 있다.
func parseSums(raw string) map[string]string {
	out := map[string]string{}
	for line := range strings.Lines(raw) {
		// 형식은 "<sha256>  <filename>"(sha256sum은 이중 공백 사용; shasum의 바이너리
		// 모드는 파일명에 * 접두를 붙인다).
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 || !isHexSHA256(fields[0]) {
			continue
		}
		name := strings.TrimPrefix(fields[1], "*")
		if name == "" {
			continue
		}
		out[name] = strings.ToLower(fields[0])
	}
	return out
}

func isHexSHA256(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9', c >= 'a' && c <= 'f', c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}

// download는 자산을 dst에 쓰면서 SHA256을 계산하고, Content-Length에 따라 진행 상황을 보고한다.
func download(ctx context.Context, c *http.Client, a Asset, dst string, prog Progress) (string, error) {
	body, err := get(ctx, c, a.URL)
	if err != nil {
		return "", fmt.Errorf("%s 다운로드: %w", a.Name, err)
	}
	defer body.Close()

	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("임시 파일 생성: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	pw := &progressWriter{total: a.Size, prog: prog, name: a.Name, last: time.Now()}
	if _, err := io.Copy(io.MultiWriter(f, h, pw), body); err != nil {
		return "", fmt.Errorf("다운로드 중단: %w", err)
	}
	if err := f.Sync(); err != nil {
		return "", fmt.Errorf("디스크 기록 실패: %w", err)
	}
	if a.Size > 0 && pw.written != a.Size {
		return "", fmt.Errorf("다운로드가 불완전합니다: 기대 %d 바이트, 실제 %d 바이트", a.Size, pw.written)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// get은 화이트리스트 제약을 받는 GET을 시작해, 응답 본문을 반환한다.
func get(ctx context.Context, c *http.Client, rawURL string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	if err := checkURL(req.URL); err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "artex-selfupdate")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// extractBinary는 릴리스 패키지에서 artex 실행 파일을 꺼낸다.
//
// 패키지 내부 구조는 artex-<버전>-<os>-<arch>/artex이지만, 여기서는 완전한 경로를 조합하지 않고 **기본 이름(base name)**으로
// 매칭한다: 버전 번호가 패키지 이름에 한 번 등장하는데 한 글자만 틀려도 업그레이드 전체가 실패하므로, 기본 이름으로 찾는 편이 변경에 더 강하다.
func extractBinary(zipPath, dst string) error {
	want := "artex"
	if runtime.GOOS == "windows" {
		want = "artex.exe"
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("릴리스 패키지 열기: %w", err)
	}
	defer zr.Close()

	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() || !strings.EqualFold(path.Base(entry.Name), want) {
			continue
		}
		rc, err := entry.Open()
		if err != nil {
			return fmt.Errorf("%s 읽기: %w", entry.Name, err)
		}
		defer rc.Close()

		f, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
		if err != nil {
			return fmt.Errorf("새 바이너리 쓰기: %w", err)
		}
		defer f.Close()

		n, err := io.Copy(f, io.LimitReader(rc, maxBinarySize+1))
		if err != nil {
			return fmt.Errorf("%s 압축 해제: %w", entry.Name, err)
		}
		if n > maxBinarySize {
			return fmt.Errorf("릴리스 패키지 내 실행 파일이 %s를 초과해 압축 해제를 거부합니다", humanSize(maxBinarySize))
		}
		if n == 0 {
			return fmt.Errorf("릴리스 패키지 내 %s가 빈 파일입니다", want)
		}
		return f.Sync()
	}
	return fmt.Errorf("릴리스 패키지에서 %s를 찾지 못했습니다", want)
}

// checkWritable은 디렉터리 쓰기 가능 여부를 미리 확인한다. 이 단계가 없으면 non-root로 실행하거나 바이너리가 시스템
// 디렉터리에 있을 때, 수십 MB를 다 받은 뒤 교체하는 순간에야 실패한다.
func checkWritable(dir string) error {
	probe, err := os.CreateTemp(dir, ".artex-update-probe-*")
	if err != nil {
		return fmt.Errorf("프로그램 디렉터리 %s에 쓸 수 없어 자동 업데이트할 수 없습니다(권한을 확인하거나 수동 업그레이드로 전환하세요): %w", dir, err)
	}
	name := probe.Name()
	_ = probe.Close()
	_ = os.Remove(name)
	return nil
}

// progressWriter는 쓴 바이트를 집계하고 보고 빈도를 제한해, 32KiB 청크마다 SSE를 하나씩 보내는 것을 방지한다.
type progressWriter struct {
	total   int64
	written int64
	name    string
	prog    Progress
	last    time.Time
}

func (w *progressWriter) Write(b []byte) (int, error) {
	w.written += int64(len(b))
	if time.Since(w.last) < 300*time.Millisecond {
		return len(b), nil
	}
	w.last = time.Now()
	pct := -1
	if w.total > 0 {
		pct = int(w.written * 100 / w.total)
	}
	w.prog(PhaseDownload, pct, fmt.Sprintf("다운로드 중 %s / %s", humanSize(w.written), humanSize(w.total)))
	return len(b), nil
}

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}

func short(sum string) string {
	if len(sum) > 12 {
		return sum[:12] + "…"
	}
	return sum
}
