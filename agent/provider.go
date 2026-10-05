// Package agent wires real LLM-driven planner and work agents (on top of the
// agent-core SDK) to the dual SQLite graph. See docs/ARTEX-架构设计.md
// §4.3 (planner) and §4.4 (work agent).
//
// Provider configuration is read from the environment so the system runs with
// any Anthropic- or OpenAI-format endpoint. If no key is configured, FromEnv
// returns ok=false and the exploration engine stays idle (an LLM is required).
package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Autumn-27/artex/llmrec"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/compaction"
	"github.com/Autumn-27/norma/llm"
	acperm "github.com/Autumn-27/norma/permission"
	"github.com/Autumn-27/norma/transcript"
)

// Config describes the LLM backend resolved from the environment.
type Config struct {
	Format  llm.Format
	BaseURL string
	APIKey  string
	Model   string
	// Proxy routes all LLM requests through the given proxy URL (http/https/socks5,
	// optionally with user:pass@ credentials). Empty means direct — it does NOT
	// fall back to the standard *_PROXY environment variables.
	Proxy string
	// RatePerSecond / RatePerMinute cap the shared request rate across ALL agents
	// using the provider (0 = that window unlimited).
	RatePerSecond float64
	RatePerMinute float64
	// ContextWindowK is the model's context window in K tokens (user-configured),
	// used to size compaction thresholds. 0 = default; see CompactionWindow.
	ContextWindowK int
	// ThinkingType은 사고 '스위치' 필드(thinking.type)를 독립적으로 제어합니다:
	//   "" = 전송하지 않음(기본값, 이 필드를 지원하지 않는 모델과 호환); "disabled" = 명시적으로 끔;
	//   "enabled" = 켬. ReasoningEffort와 완전히 분리됩니다. 일부 인터페이스에는 thinking 필드가 없고
	//   강도 파라미터만으로 사고를 활성화할 수 있으므로, 두 값을 각각 독립적으로 설정할 수 있습니다.
	ThinkingType string
	// ReasoningEffort는 사고 '강도' 필드를 독립적으로 제어합니다:
	//   "" = 전송하지 않음(기본값); "low"/"medium"/"high"/"xhigh"/"max" = 해당 강도.
	//   OpenAI에서는 최상위 reasoning_effort로, Anthropic에서는 output_config.effort로 매핑됩니다.
	ReasoningEffort string
	// Stream은 이 profile이 스트리밍(SSE) 인터페이스를 사용할지 제어합니다. true(기본값) = 스트리밍;
	// false = 완전한 비스트리밍(stream:false를 보내고 전체 JSON을 한 번에 받아 Provider.Complete 경유).
	// 비스트리밍은 일부 게이트웨이의 불완전한 SSE 구현(빈 프레임, 사고 필드 프레임 누락)을 우회할 수 있지만,
	// 실행 중 실시간 진행 상황과 실시간 token 집계를 잃습니다. agentcore.Options.NonStreaming = !Stream으로 매핑됩니다.
	Stream bool
	// MaxTokens는 단일 응답의 출력 상한(token)입니다. 0 = 이 필드를 전송하지 않고 서버 기본값을 사용
	// (기존 동작). ContextWindowK와 달리, 후자는 모델의 전체 용량이며 로컬에서 압축 임계값을 계산할 때만
	// 사용되어 요청에는 포함되지 않습니다. 이 값은 매 요청에 포함됩니다. agentcore.Options.MaxTokens로 매핑됩니다.
	MaxTokens int
	// MaxTokensField는 MaxTokens에 사용할 요청 필드명을 선택하며, format=openai에만 적용됩니다:
	//   "" = max_tokens(기본값); "max_completion_tokens" = 새 필드.
	// OpenAI 추론 모델(o 시리즈/GPT-5)은 후자만 인식하며 max_tokens를 받으면 즉시 unsupported_parameter를 반환합니다.
	// 반면 대부분의 호환 게이트웨이는 전자만 인식하므로 자동으로 추론하지 않고 사용자가 엔드포인트에 맞게 선택하도록 합니다.
	MaxTokensField string
	// SessionHeaderKey가 비어 있지 않으면 모든 LLM 요청에 사용자 지정 HTTP 헤더를 추가합니다. 헤더 이름은
	// 이 값이고 헤더 값은 [현재 세션의 session id]입니다(chat 세션=conv-<id>, worker=exp<x>-worker-i<intent>
	// 등, WorkerSessionID 참조). session-id 헤더로 프롬프트 캐시나 고정 라우팅을 수행하는 일부 게이트웨이에 사용합니다.
	// 빈 값 = 전송하지 않음. transcript.WithSessionID가 요청 context에 값을 연결하고 RoundTripper가 읽어서
	// 채우므로, 동일한 공유 provider도 세션별로 다른 헤더 값을 전송할 수 있습니다.
	SessionHeaderKey string
	// Retry는 이 설정을 해석한 재시도 파라미터입니다(profile 재정의 → 전역 정책 → 내장 기본값 순서로
	// server 측에서 해석). 세 계층의 의미는 RetryConfig를 참조하며, 0 값은 내장 기본값을 그대로 사용합니다.
	Retry RetryConfig
}

// RetryConfig는 하나의 LLM 설정에 적용되는 재시도 파라미터입니다. 각 계층의 '횟수'는 같은 의미입니다:
// 0 = 내장 기본 횟수 사용; 음수 = 해당 계층 재시도 끄기; >0 = 해당 값 사용. 각 계층의 '간격':
// 0 = 해당 계층의 기존 지수 백오프 사용; >0 = 이 고정 간격 사용.
type RetryConfig struct {
	// ConnectAttempts/ConnectInterval: SDK 연결 재시도(연결 재설정/타임아웃/429/5xx, 스트림 시작 전),
	// llm.Config.MaxRetries / RetryInterval로 직접 매핑됩니다. 기본 3회, 0.5s부터 지수적으로 증가(최대 8s).
	ConnectAttempts int
	ConnectInterval time.Duration
	// EmptyAttempts/EmptyInterval: SDK 빈 응답 재시도(완료되었지만 content block 없음, openai
	// 형식만 해당), llm.Config.EmptyResponseRetries / EmptyResponseInterval로 매핑됩니다.
	// 기본 2회, 동일한 지수 증가 단계를 사용합니다.
	EmptyAttempts int
	EmptyInterval time.Duration
	// StreamAttempts/StreamInterval: 동일 provider의 안전 구간 재시도. 이 프로젝트가 SDK 위에 추가한
	// 계층으로, '호출자에게 아직 어떤 출력도 전달하지 않은' 경우에만 스트림 중단/과부하/스트림 내 429를 재시도합니다.
	// SDK에서는 이를 볼 수 없으며 server/task_llm.go가 소비합니다. 기본 2회, 0.5s부터 지수적으로 증가(최대 4s).
	StreamAttempts int
	StreamInterval time.Duration
}

// compaction window resolution bounds (in K tokens). Below the floor the
// threshold math (window − summary reserve − buffer) would go non-positive and
// compaction would fire every turn; above the cap it would never fire.
const (
	defaultWindowK = 200  // unset → assume a 200K window (Claude default)
	minWindowK     = 32   // floor so effectiveWindow stays comfortably positive
	maxWindowK     = 1000 // cap at 1M tokens (user request)
)

// CompactionWindow returns the model context window in TOKENS for compaction
// thresholds, resolved from the user-configured size (ContextWindowK). 0/unset →
// a 200K default; otherwise clamped to [32K, 1M] so compaction stays effective.
func (c Config) CompactionWindow() int {
	k := c.ContextWindowK
	if k <= 0 {
		k = defaultWindowK
	}
	if k < minWindowK {
		k = minWindowK
	}
	if k > maxWindowK {
		k = maxWindowK
	}
	return k * 1000
}

// compactionConfig builds the agent-core compaction config for a context window
// in tokens. agentcore.NewSession wires the summarizer (same provider) when this
// is set on Options.Compaction.
func compactionConfig(windowTokens int) *compaction.Config {
	if windowTokens <= 0 {
		windowTokens = defaultWindowK * 1000
	}
	return &compaction.Config{ContextWindow: windowTokens}
}

// FromEnv resolves the LLM provider config:
//
//	ARTEX_LLM_PROVIDER = anthropic|openai (default: inferred from keys)
//	ARTEX_LLM_MODEL    = model id        (default: per provider)
//	ARTEX_LLM_BASE_URL = endpoint        (optional)
//	ARTEX_LLM_PROXY    = proxy URL        (optional; http/https/socks5)
//	ANTHROPIC_API_KEY / OPENAI_API_KEY         = credentials
func FromEnv() (Config, bool) {
	prov := os.Getenv("ARTEX_LLM_PROVIDER")
	anthKey := os.Getenv("ANTHROPIC_API_KEY")
	oaiKey := os.Getenv("OPENAI_API_KEY")

	if prov == "" {
		switch {
		case anthKey != "":
			prov = "anthropic"
		case oaiKey != "":
			prov = "openai"
		default:
			return Config{}, false
		}
	}

	c := Config{
		BaseURL: os.Getenv("ARTEX_LLM_BASE_URL"),
		Model:   os.Getenv("ARTEX_LLM_MODEL"),
		Proxy:   strings.TrimSpace(os.Getenv("ARTEX_LLM_PROXY")),
		// 기본은 스트리밍이며, ARTEX_LLM_STREAM=false/0/off이면 명시적으로 끄고 비스트리밍을 사용합니다.
		Stream: !isFalsy(os.Getenv("ARTEX_LLM_STREAM")),
	}
	switch prov {
	case "openai":
		c.Format = llm.FormatOpenAI
		c.APIKey = oaiKey
		if c.Model == "" {
			c.Model = "gpt-4o"
		}
	case "openai-responses":
		c.Format = llm.FormatOpenAIResponses
		c.APIKey = oaiKey
		if c.Model == "" {
			c.Model = "gpt-5"
		}
	default:
		c.Format = llm.FormatAnthropic
		c.APIKey = anthKey
		if c.Model == "" {
			c.Model = "claude-opus-4-8"
		}
	}
	if c.APIKey == "" {
		return Config{}, false
	}
	return c, true
}

// ConfigFrom builds a Config from UI-provided strings (provider defaults to
// anthropic; model defaults per provider). Inputs are trimmed and the base URL
// is normalized to the API base the provider expects (the provider appends the
// endpoint path itself), so a full endpoint URL is tolerated.
func ConfigFrom(provider, model, baseURL, apiKey, proxy string) Config {
	c := Config{
		Model:   strings.TrimSpace(model),
		BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		APIKey:  strings.TrimSpace(apiKey),
		Proxy:   strings.TrimSpace(proxy),
		Stream:  true, // 기본은 스트리밍이며, 호출자가 profile에 따라 재정의합니다.
	}
	switch strings.TrimSpace(provider) {
	case "openai":
		c.Format = llm.FormatOpenAI
		// provider appends "/chat/completions"; tolerate a full endpoint URL.
		c.BaseURL = strings.TrimRight(strings.TrimSuffix(c.BaseURL, "/chat/completions"), "/")
		if c.Model == "" {
			c.Model = "gpt-4o"
		}
	case "openai-responses":
		c.Format = llm.FormatOpenAIResponses
		// provider appends "/responses"; tolerate a full endpoint URL.
		c.BaseURL = strings.TrimRight(strings.TrimSuffix(c.BaseURL, "/responses"), "/")
		if c.Model == "" {
			c.Model = "gpt-5"
		}
	default:
		c.Format = llm.FormatAnthropic
		// provider appends "/v1/messages".
		c.BaseURL = strings.TrimRight(strings.TrimSuffix(c.BaseURL, "/v1/messages"), "/")
		if c.Model == "" {
			c.Model = "claude-opus-4-8"
		}
	}
	return c
}

// isFalsy reports whether an env-var string explicitly requests "off". Empty or
// unrecognized → false (so an unset var keeps the streaming default).
func isFalsy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "0", "false", "off", "no":
		return true
	}
	return false
}

// Provider returns the short provider name ("anthropic"/"openai").
func (c Config) Provider() string {
	switch c.Format {
	case llm.FormatOpenAI:
		return "openai"
	case llm.FormatOpenAIResponses:
		return "openai-responses"
	}
	return "anthropic"
}

// NewProvider builds an llm.Provider from the config. When a rate is set, the
// limiter lives on the single provider instance — so planner + all workers +
// main agent (which share this provider) are bounded by one shared rate limit.
func (c Config) NewProvider() (llm.Provider, error) {
	client, err := quotaAwareHTTPClient(c.Proxy, c.SessionHeaderKey)
	if err != nil {
		return nil, err
	}
	lc := llm.Config{
		Format:     c.Format,
		BaseURL:    c.BaseURL,
		APIKey:     c.APIKey,
		Model:      c.Model,
		HTTPClient: client,
	}
	// 사고 스위치와 강도 필드를 각각 그대로 전달합니다(빈 값 = 해당 필드를 전송하지 않음). 두 필드는 분리됩니다:
	// thinking.type만, effort만, 둘 다 전송하거나 둘 다 전송하지 않을 수 있습니다.
	lc.ThinkingType = c.ThinkingType
	lc.ReasoningEffort = c.ReasoningEffort
	// 출력 상한의 필드명을 선택합니다(빈 값 = max_tokens 사용). 상한의 '값'은 여기에서 정하지 않습니다.
	// 매 턴 agentcore.Options.MaxTokens를 따르며, provider는 값을 넣을 키만 결정합니다.
	lc.MaxTokensField = c.MaxTokensField
	// 재시도 파라미터는 SDK와 같은 의미이며(횟수 0=기본값/음수=끄기, 간격 0=지수 백오프/>0=고정), 그대로 전달합니다.
	lc.MaxRetries = c.Retry.ConnectAttempts
	lc.RetryInterval = c.Retry.ConnectInterval
	lc.EmptyResponseRetries = c.Retry.EmptyAttempts
	lc.EmptyResponseInterval = c.Retry.EmptyInterval
	if c.RatePerSecond > 0 || c.RatePerMinute > 0 {
		lc.RateLimit = &llm.RateLimit{PerSecond: c.RatePerSecond, PerMinute: c.RatePerMinute}
	}
	return llm.NewProvider(lc)
}

// IsQuotaExhaustedMessage deliberately recognizes only explicit balance,
// billing, credit, or quota-exhaustion signals. Generic 429/rate-limit text,
// authentication failures, network errors, and server failures are excluded.
var nonFailoverHTTPStatus = regexp.MustCompile(`(?:status(?:\s+code)?|http(?:\s+status)?)\s*[=:]?\s*(?:401|403|5\d\d)\b`)
var transientQuotaLimit = regexp.MustCompile(`(?i)(?:\b(?:rpm|tpm|rpd|qps)\b|quota[_\s-]*metric|rate[_\s-]*limit|too many requests|(?:requests?|tokens?)\s+(?:per|/)\s*(?:second|minute)|(?:per|/)\s*(?:second|minute)\s+(?:requests?|tokens?)|generate[_\s-]*requests[_\s-]*per[_\s-]*(?:minute|second)|tokens?[_\s-]*per[_\s-]*(?:minute|second))`)

func IsQuotaExhaustedMessage(message string) bool {
	message = strings.ToLower(message)
	// Authentication/authorization and provider-side 5xx failures never rotate,
	// even when a gateway happens to echo a quota-looking phrase in the body.
	if nonFailoverHTTPStatus.MatchString(message) {
		return false
	}
	// Provider APIs frequently describe an ordinary rate limit as "quota
	// exceeded", especially Google-style responses containing a quota metric.
	// These limits recover with time and must stay on the current provider.
	if transientQuotaLimit.MatchString(message) {
		return false
	}
	markers := []string{
		"insufficient_quota", "quota_exceeded", "quota exceeded", "quota exhausted",
		"exceeded your current quota", "billing_hard_limit_reached",
		"billing hard limit", "billing_not_active", "credit balance", "insufficient credit",
		"insufficient balance", "balance is too low", "payment required", "status 402",
		"余额不足", "额度不足", "额度已用尽", "欠费",
	}
	for _, marker := range markers {
		if strings.Contains(message, marker) {
			return true
		}
	}
	// gRPC RESOURCE_EXHAUSTED is overloaded for both account quota and ordinary
	// request-rate limiting. Preserve it as an explicit exhaustion signal only
	// when the same error does not identify a transient rate limit.
	return strings.Contains(message, "resource_exhausted") &&
		!strings.Contains(message, "rate limit") &&
		!strings.Contains(message, "too many requests")
}

// quotaAwareTransport preserves Norma's normal retry behavior except for a 429
// whose body explicitly says the account quota/balance is exhausted. Norma's
// retry loop treats every 429 as transient; normalizing only that response to
// 402 lets a task router fail over immediately while retaining the original
// response body for provider-specific classification and audit logs.
type quotaAwareTransport struct {
	base http.RoundTripper
	// sessionHeaderKey, when non-empty, is the HTTP header name each request
	// carries; its value is the session id read from the request context. Empty
	// disables it. See Config.SessionHeaderKey.
	sessionHeaderKey string
}

func (t quotaAwareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Custom session-id header: name is user-configured, value is THIS run's
	// session id (norma stashes it on the context via transcript.WithSessionID).
	// Stable across a session's turns and distinct across sessions — exactly what
	// a session-keyed prompt cache wants. Skipped when no session id is present.
	if t.sessionHeaderKey != "" {
		if sid := transcript.SessionIDFrom(req.Context()); sid != "" {
			req.Header.Set(t.sessionHeaderKey, sid)
		}
	}
	// When LLM recording is on, the Recorder puts a Capture on the context so the
	// raw wire bodies can be persisted. This is the only layer that still sees
	// them: norma builds the request body internally and decodes the SSE response
	// before either reaches the recorder.
	capt := llmrec.CaptureFrom(req.Context())
	capt.SetRequest(requestBodySnapshot(req))

	resp, err := t.base.RoundTrip(req)
	if err != nil || resp == nil {
		return resp, err
	}
	// Tee rather than read: a 200 is an SSE stream that must keep streaming. The
	// 429 branch below reads through this wrapper, so its body lands in the
	// capture before being replaced.
	resp.Body = capt.TeeResponse(resp.StatusCode, resp.Body)

	if resp.StatusCode != http.StatusTooManyRequests {
		return resp, nil
	}
	body, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	if readErr != nil {
		return resp, nil
	}
	if IsQuotaExhaustedMessage(string(body)) {
		resp.StatusCode = http.StatusPaymentRequired
		resp.Status = "402 Payment Required"
	}
	return resp, nil
}

// requestBodySnapshot copies an outgoing request body without consuming it.
// norma builds every model request from a *bytes.Reader, so net/http populates
// GetBody and the copy has no effect on what gets sent.
func requestBodySnapshot(req *http.Request) string {
	if req.GetBody == nil {
		return ""
	}
	rc, err := req.GetBody()
	if err != nil {
		return ""
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return ""
	}
	return string(b)
}

func quotaAwareHTTPClient(proxy, sessionHeaderKey string) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	proxy = strings.TrimSpace(proxy)
	if proxy == "" {
		transport.Proxy = nil // 빈 값 = 직접 연결, HTTP_PROXY/HTTPS_PROXY 환경 변수로 대체하지 않습니다.
	} else {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("llm: invalid proxy %q: %w", proxy, err)
		}
		switch proxyURL.Scheme {
		case "http", "https", "socks5":
		case "":
			return nil, fmt.Errorf("llm: proxy %q missing scheme (use http://, https:// or socks5://)", proxy)
		default:
			return nil, fmt.Errorf("llm: unsupported proxy scheme %q (use http, https or socks5)", proxyURL.Scheme)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	return &http.Client{Transport: quotaAwareTransport{base: transport, sessionHeaderKey: strings.TrimSpace(sessionHeaderKey)}}, nil
}

// logTestConnection prints the raw HTTP status code(s) and response body of a
// connection test to the server log, so "연결 테스트" leaves a diagnosable trail of
// exactly what the gateway returned — 401 bodies, quota text, empty frames — not
// just the collapsed ok/err the UI shows. Bodies are clipped to keep a chatty
// SSE stream from flooding the log.
func logTestConnection(c Config, capt *llmrec.Capture) {
	attempts := capt.Attempts()
	if len(attempts) == 0 {
		log.Printf("[llm-test] %s / %s @ %s — HTTP 요청을 보내지 못함(설정 분석 또는 연결 단계에서 실패)",
			c.Provider(), c.Model, c.BaseURL)
		return
	}
	for i, a := range attempts {
		log.Printf("[llm-test] %s / %s @ %s — 시도 %d/%d HTTP %d\n응답 본문: %s",
			c.Provider(), c.Model, c.BaseURL, i+1, len(attempts), a.Status, clipBody(a.Body))
	}
}

// clipBody trims a wire body for logging. 4K is plenty to show an error JSON or
// the head of an SSE stream while bounding a runaway response.
func clipBody(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "(비어 있음)"
	}
	const max = 4096
	if len(s) > max {
		return s[:max] + fmt.Sprintf("…(잘림, 총 %d바이트)", len(s))
	}
	return s
}

// TestConnection makes a minimal real completion to verify the provider/model/
// endpoint/key actually work. Returns the round-trip latency and the model's
// reply text.
func TestConnection(ctx context.Context, c Config) (time.Duration, string, error) {
	prov, err := c.NewProvider()
	if err != nil {
		return 0, "", err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// 원시 wire 메시지를 캡처합니다. 연결 테스트에서는 게이트웨이가 실제로 반환한 내용(상태 코드+응답 본문)을
	// 확인하는 것이 가장 중요하지만, norma가 응답을 StreamEvent로 디코딩하면 이 정보가 사라집니다.
	// quotaAwareTransport는 context에서 이 Capture를 찾아 각 HTTP 시도의 상태 코드와 body를 채웁니다.
	ctx, capt := llmrec.NewCapture(ctx)
	defer logTestConnection(c, capt)
	// 연결 테스트는 agentcore의 세션 루프를 거치지 않는 단일 요청 경로이므로 context에
	// session id를 연결하는 주체가 없습니다. SessionHeaderKey를 설정한 엔드포인트(예: opencode zen은
	// x-opencode-session 헤더를 필수로 요구하며, 없으면 즉시 400 MissingSessionID 반환)에서는 "대화는 정상인데
	// 연결 테스트는 400"인 차이가 발생합니다. 여기에서 일회성 무작위 session id를 연결하여 테스트와 실제 대화가
	// 동일한 헤더 전송 로직을 사용하게 합니다. SessionHeaderKey를 설정하지 않은 엔드포인트는 이를 읽지 않으므로 부수 효과가 없습니다.
	ctx = transcript.WithSessionID(ctx, "conntest-"+transcript.NewSessionID())
	start := time.Now()
	// MaxTokens는 충분해야 합니다. 추론 모델(예: deepseek-v4-pro)은 답변 전에 긴 사고 내용을 먼저 생성하며,
	// 실제로 "ping" 한 문장에도 ~2900 token을 소모할 수 있습니다. 32만 주면 모델이 계속 "사고 단계"에 머물다가
	// 출력 상한(finish=length)에 도달해 잘리고, 연결 테스트는 여전히 성공으로 처리되지만(err=nil)
	// "중단됨/length/resume"이 뒤섞여 표시됩니다. 예산을 충분히 주어 OK를 온전히 출력하게 합니다(finish=stop).
	// EscalateMaxTokens는 false로 유지합니다. 잘렸을 때 한도를 높여 재시도하지 않아 resume 루프에서 예산이 낭비되는 것을 방지합니다.
	reply, err := agentcore.Run(ctx, agentcore.Options{
		Provider:       prov,
		SystemPrompt:   []string{"연결 테스트입니다. 두 글자 OK만 바로 출력하세요. 사고하거나 설명하거나 다른 내용을 출력하지 마세요."},
		PermissionMode: acperm.ModeBypass,
		MaxTurns:       1,
		MaxTokens:      8192,
		NonStreaming:   !c.Stream, // 이 profile의 실제 송수신 방식으로 연결 테스트
	}, "ping")
	lat := time.Since(start)
	if err != nil {
		return lat, "", err
	}
	// err==nil만으로는 충분하지 않습니다. 요청은 성공했지만 모델이 한 글자도 출력하지 않는 경우가 실제로 있습니다
	// (사고에 예산을 모두 소모하거나, 안전 정책이 본문을 차단하거나, 호환 계층에서 content를 누락하는 경우).
	// 이런 설정은 세션에서 "응답 없음"인데도 테스트에서는 성공으로 보고됩니다. 이 차이를 없애기 위해 표시되는 본문이 없으면 모두 실패로 판정합니다.
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return lat, "", fmt.Errorf("모델 응답 내용 없음(요청은 성공했지만 텍스트가 반환되지 않음)")
	}
	return lat, reply, nil
}
