package db

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// 자산 인터셉트 규칙의 매칭/실행 레이어. asset_intercept.go는 규칙 저장만 담당하고, 여기서는
// '대상 자산'의 도메인/IP/URL을 활성 규칙과 매칭한다. agent 도구(add_intent·
// insert_assets)가 의도 하달 / 자산 삽입 전에 호출하며, 매칭되면 거부한다.

// AssetInterceptKindLabel은 kind의 라벨을 반환하며, agent 설명 메시지에 쓴다.
func AssetInterceptKindLabel(kind string) string {
	switch kind {
	case "exact_domain":
		return "도메인(완전 일치)"
	case "exact_ip":
		return "IP(완전 일치)"
	case "exact_url":
		return "URL(완전 일치)"
	case "fuzzy_domain":
		return "도메인(퍼지)"
	case "fuzzy_ip":
		return "IP(퍼지)"
	case "fuzzy_url":
		return "URL(퍼지)"
	case "cidr":
		return "CIDR 대역"
	}
	return kind
}

// Reason은 읽을 수 있는 매칭 원인을 반환한다. 형식: 자산 인터셉트 규칙 매칭 [도메인(퍼지): .gov.cn](비고).
func (r AssetInterceptRule) Reason() string {
	s := fmt.Sprintf("자산 인터셉트 규칙 매칭 [%s: %s]", AssetInterceptKindLabel(r.Kind), r.Pattern)
	if note := strings.TrimSpace(r.Note); note != "" {
		s += "(" + note + ")"
	}
	return s
}

// matchOne은 단일 활성 규칙이 주어진 도메인/IP/URL 후보 문자열에 매칭되는지 판정하고, 매칭된 구체 값을 반환한다.
func matchOne(r AssetInterceptRule, domains, ips, urls []string) (string, bool) {
	p := strings.TrimSpace(r.Pattern)
	if p == "" {
		return "", false
	}
	switch r.Kind {
	case "exact_domain":
		for _, d := range domains {
			if strings.EqualFold(strings.TrimSpace(d), p) {
				return d, true
			}
		}
	case "exact_ip":
		for _, ip := range ips {
			if strings.TrimSpace(ip) == p {
				return ip, true
			}
		}
	case "exact_url":
		for _, u := range urls {
			if strings.TrimSpace(u) == p {
				return u, true
			}
		}
	case "fuzzy_domain":
		lp := strings.ToLower(p)
		for _, d := range domains {
			if d != "" && strings.Contains(strings.ToLower(d), lp) {
				return d, true
			}
		}
	case "fuzzy_ip":
		for _, ip := range ips {
			if ip != "" && strings.Contains(ip, p) {
				return ip, true
			}
		}
	case "fuzzy_url":
		lp := strings.ToLower(p)
		for _, u := range urls {
			if u != "" && strings.Contains(strings.ToLower(u), lp) {
				return u, true
			}
		}
	case "cidr":
		_, ipnet, err := net.ParseCIDR(p)
		if err != nil {
			return "", false
		}
		for _, ip := range ips {
			if pip := net.ParseIP(strings.TrimSpace(ip)); pip != nil && ipnet.Contains(pip) {
				return ip, true
			}
		}
	}
	return "", false
}

// MatchAssetInterceptRules는 주어진 도메인/IP/URL 후보 문자열에 매칭되는 첫 번째 활성 규칙과 매칭된 구체 값을 반환한다.
// insert_assets가 원본 입력(아직 저장되지 않은 assetInputItem)으로 매칭할 때 쓴다.
func MatchAssetInterceptRules(rules []AssetInterceptRule, domains, ips, urls []string) (AssetInterceptRule, string, bool) {
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		if v, ok := matchOne(r, domains, ips, urls); ok {
			return r, v, true
		}
	}
	return AssetInterceptRule{}, "", false
}

// interceptCandidates는 이미 저장된 자산에서 인터셉트 매칭용 도메인/IP/URL 후보 문자열을 추출한다.
// URL의 host가 분리·분류되어, 'URL만 있는' 서비스 자산도 도메인/IP 규칙에 매칭될 수 있다.
func (a *Asset) interceptCandidates() (domains, ips, urls []string) {
	add := func(dst *[]string, s string) {
		if s = strings.TrimSpace(s); s != "" {
			*dst = append(*dst, s)
		}
	}
	add(&domains, a.Domain)
	add(&domains, a.RootDomain)
	for _, d := range a.BoundDomains {
		add(&domains, d)
	}
	add(&ips, a.IP)
	add(&urls, a.URL)
	if a.URL != "" {
		if u, err := url.Parse(a.URL); err == nil {
			if h := u.Hostname(); h != "" {
				if net.ParseIP(h) != nil {
					add(&ips, h)
				} else {
					add(&domains, h)
				}
			}
		}
	}
	return domains, ips, urls
}

// InterceptLabel은 자산의 짧은 식별자를 반환하며, agent 설명 메시지에 쓴다.
func (a *Asset) InterceptLabel() string {
	var target string
	switch {
	case a.Domain != "":
		target = a.Domain
	case a.URL != "":
		target = a.URL
	case a.IP != "":
		target = a.IP
	default:
		target = fmt.Sprintf("#%d", a.ID)
	}
	return fmt.Sprintf("자산#%d[%s] %s", a.ID, a.Type, target)
}

// hasEnabledRule은 규칙 집합에 활성 규칙이 하나라도 있는지 판정한다.
func hasEnabledRule(rules []AssetInterceptRule) bool {
	for _, r := range rules {
		if r.Enabled {
			return true
		}
	}
	return false
}

// AssetGateDecision은 '선 차단 후 허용' 게이트가 후보 문자열 집합에 대해 내린 판정 결과다.
type AssetGateDecision struct {
	Allowed bool
	Reason  string // 거부 사유(자산 식별자 미포함); Allowed=true면 빔
}

// EvaluateAssetGate는 작업 수준 게이트 판정을 수행한다:
//  1. 활성 blockRules 중 하나라도 매칭 → 거부(인터셉트 원인).
//  2. 그 외 allowRules에 활성 항목이 있고 모두 매칭되지 않으면 → 거부(허용 범위 밖).
//  3. 그 외에는 허용.
//
// allowRules가 비었거나 활성 항목이 없을 때는 허용 게이트가 작동하지 않는다(즉 화이트리스트 미적용, 전부 허용),
// '허용 규칙 미설정'이 모든 자산을 막아 버리는 것을 방지한다.
func EvaluateAssetGate(blockRules, allowRules []AssetInterceptRule, domains, ips, urls []string) AssetGateDecision {
	if rule, _, ok := MatchAssetInterceptRules(blockRules, domains, ips, urls); ok {
		return AssetGateDecision{Allowed: false, Reason: rule.Reason()}
	}
	if hasEnabledRule(allowRules) {
		if _, _, ok := MatchAssetInterceptRules(allowRules, domains, ips, urls); !ok {
			return AssetGateDecision{Allowed: false, Reason: "작업 허용(화이트리스트) 범위에 없어 테스트할 수 없습니다"}
		}
	}
	return AssetGateDecision{Allowed: true}
}

// AssetInterceptHit은 게이트에 거부된 자산을 기술한다(인터셉트 매칭 또는 허용 범위 밖).
type AssetInterceptHit struct {
	Asset  *Asset
	Reason string // 읽을 수 있는 원인
}

// Describe는 읽을 수 있는 설명을 반환한다: 자산 정보 + 원인.
func (h AssetInterceptHit) Describe() string {
	return fmt.Sprintf("%s → %s", h.Asset.InterceptLabel(), h.Reason)
}

// ListAssetInterceptRules는 *DB 동명 메서드의 패스스루로, AssetStore만 가진 호출자
// (예: agent 도구)도 규칙을 읽을 수 있게 한다.
func (s *AssetStore) ListAssetInterceptRules() ([]AssetInterceptRule, error) {
	return s.db.ListAssetInterceptRules()
}

// CheckAssetsIntercept는 id로 자산을 로드해 하나씩 '선 차단 후 허용' 게이트 판정을 수행하고, 거부된
// 자산을 모두 반환한다. 인터셉트 규칙 = 전역 ∪ 작업 수준 block; 허용 규칙 = 작업 수준 allow(이 작업만).
// id가 없으면 빠르게 반환한다. 전역 GetByIDs(작업 범위 필터를 받지 않음)를 써서 인터셉트가 scope에 의해 약화되지 않게 보장한다.
func (s *AssetStore) CheckAssetsIntercept(taskID int64, ids []int64) ([]AssetInterceptHit, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	blockRules, err := s.db.ListAssetInterceptRules()
	if err != nil {
		return nil, err
	}
	var allowRules []AssetInterceptRule
	if taskID > 0 {
		tb, ta, err := s.TaskInterceptRulesSplit(taskID)
		if err != nil {
			return nil, err
		}
		blockRules = append(blockRules, tb...)
		allowRules = ta
	}
	// 인터셉트 규칙도 활성 허용 규칙도 없으면 → 판정 불필요, 전부 허용.
	if len(blockRules) == 0 && !hasEnabledRule(allowRules) {
		return nil, nil
	}
	assets, err := s.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
	var hits []AssetInterceptHit
	for _, a := range assets {
		domains, ips, urls := a.interceptCandidates()
		if d := EvaluateAssetGate(blockRules, allowRules, domains, ips, urls); !d.Allowed {
			hits = append(hits, AssetInterceptHit{Asset: a, Reason: d.Reason})
		}
	}
	return hits, nil
}
