package server

import (
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
)

// 재시도 정책의 서버 측 해석, docs/LLM重试设计.md 참고. 다섯 겹 중:
//   - 연결 수립 / 빈 응답 / 같은 provider 안전 구간은 '엔드포인트를 따르는' 것이라, 각 LLM 설정이
//     전역 기본값을 재정의할 수 있다(profile의 어떤 항목을 비우면 전역을 상속, 전역도 미설정이면 내장 기본값);
//   - 회로 차단 / 의도 재실행은 프로세스 레벨이라, 전역 하나뿐.
//
// 전역 정책은 DB의 settings 한 행을 한 번 읽으며, 호출 지점이 모두 저빈도 경로(provider 구축, work 마무리,
// 설정 저장)라, 캐시를 한 겹 더 둘 가치가 없다; 회로 차단 파라미터는 예외 —— 실패 경로에서 매번 읽어야 하므로,
// applyRetryPolicy가 Registry에 푸시해 저장한다.

// retryPolicy reads the global policy; a nil DB yields the zero policy (all
// layers on their built-in defaults).
func (s *Server) retryPolicy() db.LLMRetryPolicy {
	if s.m == nil || s.m.pg == nil {
		return db.LLMRetryPolicy{}
	}
	return s.m.pg.LLMRetryPolicy()
}

// resolveRetry layers one profile's override on top of the global policy and
// converts the result into the form agent.Config carries. Rules combine field by
// field, so a profile that only pins an interval still inherits the global count.
func resolveRetry(o db.RetryOverride, pol db.LLMRetryPolicy) agent.RetryConfig {
	connect := o.Connect.Or(pol.Connect)
	empty := o.Empty.Or(pol.Empty)
	stream := o.Stream.Or(pol.Stream)
	return agent.RetryConfig{
		// 횟수는 여기서 '0=기본 / 음수=끔'의 원래 의미를 유지: SDK의 MaxRetries /
		// EmptyResponseRetries와 완전히 동형이라, 자체 해석에 맡기면 된다.
		ConnectAttempts: connect.Attempts, ConnectInterval: connect.Interval(),
		EmptyAttempts: empty.Attempts, EmptyInterval: empty.Interval(),
		StreamAttempts: stream.Attempts, StreamInterval: stream.Interval(),
	}
}

// applyProfileRetry fills cfg.Retry for a profile read from the DB.
func (s *Server) applyProfileRetry(cfg *agent.Config, p *db.LLMProfile) {
	if p == nil {
		return
	}
	cfg.Retry = resolveRetry(p.Retry, s.retryPolicy())
}

// 회로 차단(순환 전환 재시도 대기 시간)의 기본값, llmpool 내장과 동일 —— 여기서는 '사용자가 값을 설정'했을 때만 재정의.
// 의도 재실행의 기본값은 engine.go의 modelErrorRetries / modelErrorRetryBackoff 참고.

// applyRetryPolicy pushes the process-wide layers of the policy into the objects
// that consume them on a hot path: the circuit-breaker registry. Called at
// startup and whenever the policy is saved.
func (s *Server) applyRetryPolicy() {
	pol := s.retryPolicy()
	if s.llmHealth != nil {
		s.llmHealth.SetPolicy(pol.Breaker.Attempts, pol.Breaker.Interval())
	}
}

// modelErrorRetryPolicy resolves the intent-level replay knobs (layer ⑤): how
// many times a model_error work is re-run and how long to back off between runs.
func (e *Engine) modelErrorRetryPolicy() (retries int, backoff time.Duration) {
	retries, backoff = modelErrorRetries, modelErrorRetryBackoff
	if e == nil || e.m == nil || e.m.pg == nil {
		return retries, backoff
	}
	rule := e.m.pg.LLMRetryPolicy().Intent
	if rule.Attempts != 0 {
		retries = max(rule.Attempts, 0)
	}
	if d := rule.Interval(); d > 0 {
		backoff = d
	}
	return retries, backoff
}

// emptyTurnNudgeLimit resolves how many empty-turn continuations one work may
// inject (see steerHooks.Stop). It deliberately reuses layer ②'s knob —— '빈 응답
// 재시도 횟수': 둘은 같은 일의 두 수단이다. SDK 그 겹은 '내용 블록이 하나도 없음'을 다루고, 수단은
// 같은 요청을 그대로 재전송; 여기서는 '사고만 있고 본문도 도구도 없음'을 다루고, 수단은 지시 하나를 덧붙여
// 모델이 기존 사고를 가지고 이어가게 한다(그대로 재전송은 이렇게 컨텍스트 형태로 결정되는 공회전에 의미 없음). 빈 것 판단 기준이
// 다른 건 SDK가 '이벤트를 yield한 적 있는지'를 기준으로 삼기 때문이고, 사고 증분 자체가 이벤트다 —— 하지만 사용자가
// '빈 응답 재시도 몇 번'을 설정할 때 표현하려는 건 '모델이 실질 내용을 산출하지 않으면 한 번 더'이고, 두 겹이 같은 횟수를 공유해야
// 이 심상에 맞는다.
//
// 특정 profile의 재정의가 아니라 전역 정책을 읽음: 한 run이 도중에 failover로 profile을 바꿀 수 있지만, 이것은
// 의도 전체의 총량 게이트라, 엔드포인트가 바뀐다고 따라 바뀌면 안 된다. 의미는 SDK의 emptyRetries()와 동형:
// 0 = 기본 defaultEmptyTurnNudges; -1(음수) = 공회전 이어 실행 끔; >0 = 그 값 사용.
func (e *Engine) emptyTurnNudgeLimit() int {
	if e == nil || e.m == nil || e.m.pg == nil {
		return defaultEmptyTurnNudges
	}
	switch n := e.m.pg.LLMRetryPolicy().Empty.Attempts; {
	case n == 0:
		return defaultEmptyTurnNudges
	case n < 0:
		return 0
	default:
		return n
	}
}
