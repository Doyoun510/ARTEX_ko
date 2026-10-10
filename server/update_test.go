package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Autumn-27/artex/selfupdate"
)

// releaseCache는 GitHub 할당량을 보호하는 계층입니다: 미인증 API는 60회/시간/IP만 허용되며,
// 상단 표시줄의 "새 버전 있음" 안내는 전체 페이지를 로드할 때마다 조회합니다. 캐시가 무효화되면 사용자가
// 탭을 몇 개 더 여는 것만으로 할당량을 소진해, 정작 업데이트하려 할 때 조회할 수 없게 됩니다.

func newTestCache(fetch func(context.Context, *http.Client) (*selfupdate.Release, error)) *releaseCache {
	return &releaseCache{fetch: fetch}
}

func TestReleaseCacheServesFromCache(t *testing.T) {
	calls := 0
	c := newTestCache(func(context.Context, *http.Client) (*selfupdate.Release, error) {
		calls++
		return &selfupdate.Release{TagName: "v0.3.8"}, nil
	})

	for range 5 {
		rel, err := c.get(t.Context(), nil, false)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if rel.TagName != "v0.3.8" {
			t.Fatalf("TagName = %q", rel.TagName)
		}
	}
	if calls != 1 {
		t.Errorf("5번 조회 시 원본 서버 조회는 1번만 해야 하며, 실제 %d번", calls)
	}
}

func TestReleaseCacheForceBypasses(t *testing.T) {
	calls := 0
	c := newTestCache(func(context.Context, *http.Client) (*selfupdate.Release, error) {
		calls++
		return &selfupdate.Release{TagName: "v0.3.8"}, nil
	})

	if _, err := c.get(t.Context(), nil, false); err != nil {
		t.Fatal(err)
	}
	// 사용자가 '업데이트 확인'을 누르면 반드시 실시간 결과를 받아야 합니다. 그렇지 않으면 방금 배포된 버전을 캐시가 만료될 때까지 볼 수 없습니다.
	if _, err := c.get(t.Context(), nil, true); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("force는 캐시를 건너뛰어야 하며, 원본 서버 조회는 2번이어야 하지만 실제 %d번", calls)
	}
}

func TestReleaseCacheExpiresAfterTTL(t *testing.T) {
	calls := 0
	c := newTestCache(func(context.Context, *http.Client) (*selfupdate.Release, error) {
		calls++
		return &selfupdate.Release{TagName: "v0.3.8"}, nil
	})

	if _, err := c.get(t.Context(), nil, false); err != nil {
		t.Fatal(err)
	}
	// 저장 시각을 방금 만료된 시점까지 앞당겨 TTL 만료를 재현합니다.
	c.at = time.Now().Add(-releaseTTL - time.Second)
	if _, err := c.get(t.Context(), nil, false); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("TTL 만료 후 원본 서버를 다시 조회해야 하며, 2번이어야 하지만 실제 %d번", calls)
	}
}

func TestReleaseCacheUsesShorterTTLForErrors(t *testing.T) {
	calls := 0
	c := newTestCache(func(context.Context, *http.Client) (*selfupdate.Release, error) {
		calls++
		return nil, errors.New("github 不可达")
	})

	if _, err := c.get(t.Context(), nil, false); err == nil {
		t.Fatal("오류 반환을 기대합니다")
	}
	// 실패 결과도 잠시 캐시해야 합니다. 그렇지 않으면 GitHub에 접근할 수 없을 때 페이지를 로드할 때마다 타임아웃을 헛되이 기다립니다.
	if _, err := c.get(t.Context(), nil, false); err == nil {
		t.Fatal("오류 반환을 기대합니다")
	}
	if calls != 1 {
		t.Errorf("오류는 짧게 캐시해야 하며, 원본 서버 조회는 1번이어야 하지만 실제 %d번", calls)
	}

	// 다만 오류 TTL은 반드시 성공 TTL보다 확실히 짧아야 네트워크가 복구된 뒤 빠르게 회복할 수 있습니다.
	if releaseErrTTL >= releaseTTL {
		t.Fatalf("오류 TTL(%v)은 반드시 성공 TTL(%v)보다 짧아야 합니다", releaseErrTTL, releaseTTL)
	}
	c.at = time.Now().Add(-releaseErrTTL - time.Second)
	if _, err := c.get(t.Context(), nil, false); err == nil {
		t.Fatal("오류 반환을 기대합니다")
	}
	if calls != 2 {
		t.Errorf("오류 TTL 만료 후 재시도해야 하며, 2번이어야 하지만 실제 %d번", calls)
	}
}

func TestReleaseCacheDoesNotPoisonOnCallerCancel(t *testing.T) {
	good := &selfupdate.Release{TagName: "v0.3.8"}
	c := newTestCache(func(ctx context.Context, _ *http.Client) (*selfupdate.Release, error) {
		return good, nil
	})
	if _, err := c.get(t.Context(), nil, false); err != nil {
		t.Fatal(err)
	}

	// 방문자가 탭을 닫으면 요청이 취소됩니다. 이는 GitHub에 문제가 있다는 뜻이 아니며, "취소됨"을 절대로
	// 캐시에 기록해서는 안 됩니다. 정상 결과 TTL은 30분, 일반 오류 TTL은 2분이며, 취소 오류는 캐시하지 않아야 합니다.
	c.fetch = func(ctx context.Context, _ *http.Client) (*selfupdate.Release, error) {
		return nil, ctx.Err()
	}
	c.at = time.Now().Add(-releaseTTL - time.Second) // 캐시를 만료시켜 원본 서버를 조회하도록 강제

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.get(ctx, nil, false); err == nil {
		t.Fatal("호출자가 이미 취소했으면 오류를 그대로 전달해야 합니다")
	}

	// 핵심 불변 조건: 취소된 호출은 어떤 흔적도 남기지 않습니다. 캐시에는 "취소됨" 오류가 없으며,
	// 직전의 정상 결과도 유지됩니다.
	if c.err != nil {
		t.Fatalf("취소 오류를 캐시에 기록하면 안 됩니다. 받은 값: %v", c.err)
	}
	if c.rel == nil || c.rel.TagName != "v0.3.8" {
		t.Fatalf("캐시는 직전의 정상 결과를 유지해야 합니다. 받은 값: %+v", c.rel)
	}

	// 그 취소로는 새 데이터를 얻지 못했으므로 다음 요청자는 원본 서버를 다시 조회해야 하며 정상적으로 결과를 얻어야 합니다.
	// 이전 취소의 영향을 받으면 안 됩니다.
	c.fetch = func(context.Context, *http.Client) (*selfupdate.Release, error) {
		return good, nil
	}
	rel, err := c.get(t.Context(), nil, false)
	if err != nil {
		t.Fatalf("취소 이후 정상 요청은 오류를 반환하면 안 됩니다: %v", err)
	}
	if rel == nil || rel.TagName != "v0.3.8" {
		t.Fatalf("정상 결과를 받아야 합니다. 받은 값: %+v", rel)
	}
}
