package db

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Finding asset tree — the '자산별' view of the global findings list.
//
// 계층은 BuildCoverageGraph와 동일 출처다(company → root_domain/ip/app → subdomain →
// service → endpoint), 그쪽은 '어떤 작업의 범위 내 자산'의 포스-디렉티드 그래프이고, 이쪽은
// '전체 DB에서 발견이 있는 자산'의 트리다: 발견이 있는 자산과 그 조상 체인만 수집하고, 노드는 서브트리 집계 카운트를 가진다.
// 둘의 부모-자식 우선순위 규칙은 반드시 일치해야 하므로, 한쪽을 바꾸면 task_scope.go를 대조해 다른 쪽도 바꿔라.
// ---------------------------------------------------------------------------

// FindingUnassignedAsset은 '미연관 자산' 노드 key이자, 목록 API의 필터 센티넬이다:
// asset_ids가 비었거나 가리키는 자산이 삭제된 발견에 적중한다.
const FindingUnassignedAsset = "__none__"

// findingAssetTreeMaxNodes는 프런트에 반환하는 노드 상한이다. 초과 시 바닥에서 위로 한 층씩 통째로 버린다
// (endpoint 우선, 그다음 service): 그 카운트는 이미 부모 노드에 누적되어, 노드는 버려도 숫자는 버리지 않는다.
const findingAssetTreeMaxNodes = 3000

// FindingAssetNode는 자산 트리의 한 노드다. Key는 커버리지 그래프와 동형이다: 자산 행은 "a:<id>",
// 회사는 "c:<id>", 자산 행이 없는 루트 도메인은 합성된 "r:<domain>", 미연관 버킷은 "__none__"이다.
type FindingAssetNode struct {
	Key       string `json:"key"`
	Parent    string `json:"parent,omitempty"`
	Kind      string `json:"kind"` // company|root_domain|subdomain|ip|service|app|endpoint|none
	Label     string `json:"label"`
	AssetID   int64  `json:"asset_id,omitempty"`
	CompanyID int64  `json:"company_id,omitempty"`
	// Self는 해당 자산에 직접 달린 발견 수; Total은 모든 자손을 포함하며 finding 기준으로 중복 제거한다
	// (한 발견이 여러 자산에 달리면 공통 조상에서만 한 번 계산한다).
	Self        int       `json:"self"`
	Total       int       `json:"total"`
	Critical    int       `json:"critical"`
	High        int       `json:"high"`
	Medium      int       `json:"medium"`
	Low         int       `json:"low"`
	LastFoundAt time.Time `json:"last_found_at"`
}

// FindingAssetTree는 전체 트리의 일회성 스냅샷이다. Nodes는 정렬되어 있다: 같은 부모 노드 아래에서 발견 수
// 내림차순, 라벨 오름차순으로, '미연관 자산'은 항상 마지막이다.
type FindingAssetTree struct {
	Nodes        []FindingAssetNode `json:"nodes"`
	FindingTotal int                `json:"finding_total"`
	// Truncated=true는 크기를 제어하려고 DroppedKinds의 계층을 버렸음을 뜻한다.
	Truncated    bool     `json:"truncated"`
	DroppedKinds []string `json:"dropped_kinds,omitempty"`
}

// assetRow는 트리 구성에 필요한 자산 필드의 부분집합이다.
type assetRow struct {
	id          int64
	kind        string
	companyID   int64
	domain      string
	rootDomain  string
	ip          string
	url         string
	port        int
	serviceType string
	appName     string
}

func (a *assetRow) coverageNode() CoverageGraphNode {
	return CoverageGraphNode{
		Kind: a.kind, Domain: a.domain, RootDomain: a.rootDomain, IP: a.ip,
		URL: a.url, Port: a.port, ServiceType: a.serviceType, AppName: a.appName,
	}
}

// label은 커버리지 그래프의 라벨 규칙(URL > domain > ip > app_name > root_domain)을 재사용하지만, URL이 없는
// 서비스(SMB·비 HTTP 포트 등)는 포트를 덧붙여야 한다: 안 그러면 그 라벨이 호스트 IP/도메인 행과
// 똑같아져, 트리에서 부모-자식 두 행이 완전히 중복돼 보인다.
func (a *assetRow) label() string {
	if a.kind == "service" && a.url == "" {
		if host, port := a.hostPort(); host != "" && port > 0 {
			return host + ":" + strconv.Itoa(port)
		}
	}
	n := a.coverageNode()
	n.Key = assetKey(a.id)
	return coverageNodeLabel(&n)
}

// hostPort는 커버리지 그래프와 일치한다: domain 우선, 그다음 URL의 host, 마지막 ip.
func (a *assetRow) hostPort() (string, int) {
	n := a.coverageNode()
	return hostPortOf(&n)
}

const findingAssetSelectCols = `a.id, a.type, COALESCE(a.company_id,0),
       COALESCE(a.domain,''), COALESCE(a.root_domain,''), COALESCE(a.ip,''),
       COALESCE(a.url,''), COALESCE(a.port,0), COALESCE(a.service_type,''),
       COALESCE(a.app_name,'')`

func scanAssetRows(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
	Close() error
}) ([]*assetRow, error) {
	defer rows.Close()
	var out []*assetRow
	for rows.Next() {
		a := &assetRow{}
		if err := rows.Scan(&a.id, &a.kind, &a.companyID, &a.domain, &a.rootDomain,
			&a.ip, &a.url, &a.port, &a.serviceType, &a.appName); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// findingAssetHit은 한 발견이 트리 구성 단계에서 필요한 최소 정보다.
type findingAssetHit struct {
	severity string
	ts       time.Time
	assetIDs []int64
}

// BuildFindingAssetTree는 현재 필터로 자산 트리를 구성한다. AssetScope 자체는 참여하지 않는다(안 그러면 트리가
// 선택 노드를 따라 한 체인으로 붕괴한다).
func (d *DB) BuildFindingAssetTree(f FindingFilter) (*FindingAssetTree, error) {
	return d.buildFindingAssetTree(f, findingAssetTreeMaxNodes)
}

// buildFindingAssetTree는 노드 상한이 있는 내부 구현이다. maxNodes<=0은 비절단을 의미한다 —— AssetScope를
// 파싱할 때는 반드시 이 모드를 써야 한다, 안 그러면 버려진 endpoint가 서브트리 id 집합을 불완전하게 만든다.
func (d *DB) buildFindingAssetTree(f FindingFilter, maxNodes int) (*FindingAssetTree, error) {
	f.AssetScope = ""
	f.assetIDs, f.assetNone, f.assetMiss = nil, false, false
	where, args := f.where()

	rows, err := d.Query(`SELECT COALESCE(f.severity,''), f.created_at,
       COALESCE(f.asset_ids::text,'[]')
FROM findings f LEFT JOIN tasks t ON f.task_id = t.id`+where, args...)
	if err != nil {
		return nil, err
	}
	hits := []findingAssetHit{}
	assetIDs := map[int64]bool{}
	for rows.Next() {
		var h findingAssetHit
		var aidsJSON string
		if err := rows.Scan(&h.severity, &h.ts, &aidsJSON); err != nil {
			rows.Close()
			return nil, err
		}
		_ = json.Unmarshal([]byte(aidsJSON), &h.assetIDs)
		for _, id := range h.assetIDs {
			if id > 0 {
				assetIDs[id] = true
			}
		}
		hits = append(hits, h)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	tree := &FindingAssetTree{Nodes: []FindingAssetNode{}, FindingTotal: len(hits)}
	byID, err := d.loadFindingAssetRows(assetIDs)
	if err != nil {
		return nil, err
	}

	nodes, parentOf := d.assembleFindingAssetNodes(byID)
	if err := d.attachCompanyNodes(nodes, parentOf); err != nil {
		return nil, err
	}

	// 카운트: 한 발견이 각 자산의 조상 체인을 따라 올라가며 중복 제거한 key 집합을 모아 하나씩 +1 하므로,
	// 부모 노드가 한 발견이 여러 하위 자산에 달렸다는 이유로 중복 계산되지 않는다.
	unassigned := &FindingAssetNode{Key: FindingUnassignedAsset, Kind: "none", Label: "미연관 자산"}
	touched := map[string]bool{}
	for _, h := range hits {
		clear(touched)
		var direct []*FindingAssetNode
		for _, id := range h.assetIDs {
			node := nodes[assetKey(id)]
			if node == nil {
				continue
			}
			direct = append(direct, node)
			for key := node.Key; key != ""; key = parentOf[key] {
				touched[key] = true
			}
		}
		if len(direct) == 0 {
			countFinding(unassigned, h)
			unassigned.Self++
			continue
		}
		for _, node := range direct {
			node.Self++
		}
		for key := range touched {
			countFinding(nodes[key], h)
		}
	}

	for _, node := range nodes {
		if node.Total > 0 {
			tree.Nodes = append(tree.Nodes, *node)
		}
	}
	if unassigned.Total > 0 {
		tree.Nodes = append(tree.Nodes, *unassigned)
	}
	sortFindingAssetNodes(tree.Nodes)
	truncateFindingAssetTree(tree, maxNodes)
	return tree, nil
}

// countFinding은 한 발견을 노드에 누적한다(총수 / 심각도 버킷 / 최근 발견 시간).
func countFinding(n *FindingAssetNode, h findingAssetHit) {
	if n == nil {
		return
	}
	n.Total++
	switch h.severity {
	case "critical":
		n.Critical++
	case "high":
		n.High++
	case "medium":
		n.Medium++
	case "low":
		n.Low++
	}
	if h.ts.After(n.LastFoundAt) {
		n.LastFoundAt = h.ts
	}
}

// loadFindingAssetRows는 적중한 자산 행을 읽고, 라운드마다 조상을 채운다(service의 호스트 도메인/IP,
// 서브도메인의 루트 도메인). 조상 자체는 발견이 없을 수 있지만, 트리가 형태를 갖추려면 필요하다.
func (d *DB) loadFindingAssetRows(ids map[int64]bool) (map[int64]*assetRow, error) {
	byID := map[int64]*assetRow{}
	if len(ids) == 0 {
		return byID, nil
	}
	idList := make([]int64, 0, len(ids))
	for id := range ids {
		idList = append(idList, id)
	}
	rows, err := d.Query(`SELECT `+findingAssetSelectCols+` FROM assets a WHERE a.id = ANY($1::bigint[])`, idList)
	if err != nil {
		return nil, err
	}
	found, err := scanAssetRows(rows)
	if err != nil {
		return nil, err
	}
	for _, a := range found {
		byID[a.id] = a
	}

	// 매 라운드 아직 부모 노드가 없는 호스트 식별자를 찾아 한 층을 일괄로 채운다; 층수는 고정이라(endpoint→service→
	// subdomain/ip→root_domain), 4라운드면 충분히 수렴한다.
	for range 4 {
		want := missingParents(byID)
		if want.empty() {
			break
		}
		added, err := d.loadAssetsByHost(want, byID)
		if err != nil {
			return nil, err
		}
		if added == 0 {
			break
		}
	}
	return byID, nil
}

// missingHosts는 한 라운드 채우기에서 DB에 찾아야 할 호스트 식별자로, 대상 자산 유형별로 분리한다.
type missingHosts struct {
	services []string // endpoint의 호스트(service 행 찾기)
	domains  []string // service/endpoint의 호스트 도메인(subdomain 행 찾기)
	ips      []string // service/endpoint의 호스트 IP(ip 행 찾기)
	roots    []string // 서브도메인의 루트 도메인(root_domain 행 찾기)
}

func (m missingHosts) empty() bool {
	return len(m.services) == 0 && len(m.domains) == 0 && len(m.ips) == 0 && len(m.roots) == 0
}

// missingParents는 아직 로드되지 않은 호스트를 집계한다: service(endpoint가 붙음), 서브도메인/IP(
// service와 endpoint가 붙음), 루트 도메인(서브도메인이 붙음).
func missingParents(byID map[int64]*assetRow) missingHosts {
	haveService := map[string]bool{}
	haveDomain := map[string]bool{}
	haveIP := map[string]bool{}
	haveRoot := map[string]bool{}
	for _, a := range byID {
		switch a.kind {
		case "service":
			if host, _ := a.hostPort(); host != "" {
				haveService[host] = true
			}
		case "subdomain":
			haveDomain[a.domain] = true
		case "ip":
			haveIP[a.ip] = true
		case "root_domain":
			haveRoot[a.domain] = true
		}
	}
	wantService := map[string]bool{}
	wantDomain := map[string]bool{}
	wantIP := map[string]bool{}
	wantRoot := map[string]bool{}
	for _, a := range byID {
		switch a.kind {
		case "service", "endpoint":
			host, _ := a.hostPort()
			// endpoint는 먼저 같은 호스트의 service를 찾는다; 포트가 안 맞는 service는 Total=0으로
			// 최종 필터링되어 트리를 오염시키지 않는다.
			if a.kind == "endpoint" && host != "" && !haveService[host] {
				wantService[host] = true
			}
			if host != "" && !haveDomain[host] && !haveIP[host] && !haveRoot[host] {
				if isIPLiteral(host) {
					wantIP[host] = true
				} else {
					wantDomain[host] = true
				}
			}
			if a.ip != "" && !haveIP[a.ip] {
				wantIP[a.ip] = true
			}
		case "subdomain":
			if a.rootDomain != "" && !haveRoot[a.rootDomain] {
				wantRoot[a.rootDomain] = true
			}
		}
	}
	return missingHosts{
		services: keysOf(wantService),
		domains:  keysOf(wantDomain),
		ips:      keysOf(wantIP),
		roots:    keysOf(wantRoot),
	}
}

func keysOf(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// isIPLiteral은 host가 IP 리터럴인지 대략 판정한다(ip 테이블에서 찾을지 subdomain 테이블에서 찾을지 결정에 사용).
func isIPLiteral(host string) bool {
	if strings.Contains(host, ":") {
		return true // IPv6
	}
	if host == "" {
		return false
	}
	for _, part := range strings.Split(host, ".") {
		if part == "" {
			return false
		}
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return strings.Count(host, ".") == 3
}

// loadAssetsByHost는 호스트 식별자로 자산 행을 일괄 채우고, 이번 라운드에 새로 추가된 행 수를 반환한다.
func (d *DB) loadAssetsByHost(want missingHosts, byID map[int64]*assetRow) (int, error) {
	added := 0
	load := func(q string, arg []string) error {
		if len(arg) == 0 {
			return nil
		}
		rows, err := d.Query(q, arg)
		if err != nil {
			return err
		}
		found, err := scanAssetRows(rows)
		if err != nil {
			return err
		}
		for _, a := range found {
			if _, ok := byID[a.id]; ok {
				continue
			}
			byID[a.id] = a
			added++
		}
		return nil
	}
	if err := load(`SELECT `+findingAssetSelectCols+` FROM assets a
WHERE a.type='service' AND (a.domain = ANY($1::text[]) OR a.ip = ANY($1::text[]))`, want.services); err != nil {
		return 0, err
	}
	if err := load(`SELECT `+findingAssetSelectCols+` FROM assets a
WHERE a.type='subdomain' AND a.domain = ANY($1::text[])`, want.domains); err != nil {
		return 0, err
	}
	if err := load(`SELECT `+findingAssetSelectCols+` FROM assets a
WHERE a.type='ip' AND a.ip = ANY($1::text[])`, want.ips); err != nil {
		return 0, err
	}
	if err := load(`SELECT `+findingAssetSelectCols+` FROM assets a
WHERE a.type='root_domain' AND a.domain = ANY($1::text[])`, want.roots); err != nil {
		return 0, err
	}
	return added, nil
}

// assembleFindingAssetNodes는 자산 행을 노드로 만들고 부모-자식 관계를 연결한다. 부모 노드가 없을 때(DB에
// 그 루트 도메인 자산이 아예 없음) "r:<domain>" 자리표시 노드를 합성하며, 커버리지 그래프 처리와 일치한다.
func (d *DB) assembleFindingAssetNodes(byID map[int64]*assetRow) (map[string]*FindingAssetNode, map[string]string) {
	nodes := map[string]*FindingAssetNode{}
	parentOf := map[string]string{}
	rootByDomain := map[string]string{}
	subByDomain := map[string]string{}
	ipByAddr := map[string]string{}
	svcByHost := map[string]string{}
	svcByHostPort := map[string]string{}

	for _, a := range byID {
		key := assetKey(a.id)
		nodes[key] = &FindingAssetNode{
			Key: key, Kind: a.kind, Label: a.label(),
			AssetID: a.id, CompanyID: a.companyID,
		}
		switch a.kind {
		case "root_domain":
			if a.domain != "" {
				rootByDomain[a.domain] = key
			}
		case "subdomain":
			if a.domain != "" {
				subByDomain[a.domain] = key
			}
		case "ip":
			if a.ip != "" {
				ipByAddr[a.ip] = key
			}
		case "service":
			if host, port := a.hostPort(); host != "" {
				svcByHost[host] = key
				svcByHostPort[host+"|"+strconv.Itoa(port)] = key
			}
		}
	}

	// 서브도메인의 루트 도메인이 DB에 자산 행이 없을 때 자리표시 루트를 합성해, 서브도메인이 최상위로 흩어지지 않게 한다.
	for _, a := range byID {
		if a.kind != "subdomain" || a.rootDomain == "" {
			continue
		}
		if _, ok := rootByDomain[a.rootDomain]; ok {
			continue
		}
		key := "r:" + a.rootDomain
		nodes[key] = &FindingAssetNode{Key: key, Kind: "root_domain", Label: a.rootDomain}
		rootByDomain[a.rootDomain] = key
	}

	firstOf := func(keys ...string) string {
		for _, k := range keys {
			if k != "" {
				if _, ok := nodes[k]; ok {
					return k
				}
			}
		}
		return ""
	}
	for _, a := range byID {
		key := assetKey(a.id)
		var parent string
		switch a.kind {
		case "subdomain":
			parent = firstOf(rootByDomain[a.rootDomain])
		case "service":
			host, _ := a.hostPort()
			parent = firstOf(subByDomain[a.domain], subByDomain[host],
				ipByAddr[a.ip], ipByAddr[host], rootByDomain[a.rootDomain], rootByDomain[host])
		case "endpoint":
			host, port := a.hostPort()
			parent = firstOf(svcByHostPort[host+"|"+strconv.Itoa(port)], svcByHost[host],
				subByDomain[host], subByDomain[a.domain], ipByAddr[host], ipByAddr[a.ip],
				rootByDomain[a.rootDomain], rootByDomain[host])
		}
		if parent != "" && parent != key {
			parentOf[key] = parent
			nodes[key].Parent = parent
		}
	}
	return nodes, parentOf
}

// attachCompanyNodes는 최상위 자산(루트 도메인 / IP / 앱)에 회사 부모 노드를 보충한다 —— 자산이 실제로
// 회사에 귀속된 경우에만 회사 계층이 나타나고, 귀속되지 않은 자산은 그대로 자기 자신이 최상위다.
func (d *DB) attachCompanyNodes(nodes map[string]*FindingAssetNode, parentOf map[string]string) error {
	want := map[int64]bool{}
	for _, n := range nodes {
		if n.Parent != "" || n.CompanyID <= 0 {
			continue
		}
		switch n.Kind {
		case "root_domain", "ip", "app":
			want[n.CompanyID] = true
		}
	}
	if len(want) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(want))
	for id := range want {
		ids = append(ids, id)
	}
	rows, err := d.Query(`SELECT id, COALESCE(name,'') FROM companies WHERE id = ANY($1::bigint[])`, ids)
	if err != nil {
		return err
	}
	defer rows.Close()
	names := map[int64]string{}
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		names[id] = name
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for id, name := range names {
		key := companyKey(id)
		if _, ok := nodes[key]; ok {
			continue
		}
		if name == "" {
			name = "회사 #" + strconv.FormatInt(id, 10)
		}
		nodes[key] = &FindingAssetNode{Key: key, Kind: "company", Label: name, CompanyID: id}
	}
	for _, n := range nodes {
		if n.Parent != "" || n.CompanyID <= 0 || n.Kind == "company" {
			continue
		}
		switch n.Kind {
		case "root_domain", "ip", "app":
			key := companyKey(n.CompanyID)
			if _, ok := nodes[key]; !ok {
				continue
			}
			n.Parent = key
			parentOf[n.Key] = key
		}
	}
	return nil
}

// sortFindingAssetNodes 정렬: 발견 많은 것이 앞, 같으면 라벨순; '미연관 자산'은 항상 마지막.
// 프런트는 배열 순서대로 자식 노드를 붙이므로, 같은 부모 노드 아래 상대 순서만 맞으면 된다.
func sortFindingAssetNodes(nodes []FindingAssetNode) {
	sort.SliceStable(nodes, func(i, j int) bool {
		a, b := nodes[i], nodes[j]
		if (a.Kind == "none") != (b.Kind == "none") {
			return b.Kind == "none"
		}
		if a.Total != b.Total {
			return a.Total > b.Total
		}
		return a.Label < b.Label
	})
}

// truncateFindingAssetTree는 노드가 너무 많으면 한 층을 통째로 버린다(endpoint 먼저, 그다음 service). 카운트는 이미
// 부모 노드에 누적되어, 버리는 것은 펼칠 수 있는 세부 계층뿐이다.
func truncateFindingAssetTree(tree *FindingAssetTree, maxNodes int) {
	if maxNodes <= 0 || len(tree.Nodes) <= maxNodes {
		return
	}
	for _, kind := range []string{"endpoint", "service"} {
		kept := tree.Nodes[:0]
		for _, n := range tree.Nodes {
			if n.Kind == kind {
				continue
			}
			kept = append(kept, n)
		}
		tree.Nodes = kept
		tree.Truncated = true
		tree.DroppedKinds = append(tree.DroppedKinds, kind)
		if len(tree.Nodes) <= maxNodes {
			return
		}
	}
}

// applyAssetScope는 AssetScope(노드 key)를 SQL에 쓸 수 있는 자산 id 집합으로 파싱한다. 노드 하나
// 선택은 그 전체 서브트리 선택과 같으므로, 먼저 트리를 만든 뒤 자손을 수집해야 한다.
func (d *DB) applyAssetScope(f FindingFilter) (FindingFilter, error) {
	scope := strings.TrimSpace(f.AssetScope)
	f.assetIDs, f.assetNone, f.assetMiss = nil, false, false
	if scope == "" {
		return f, nil
	}
	if scope == FindingUnassignedAsset {
		f.assetNone = true
		return f, nil
	}
	// 비절단: 버려진 endpoint도 id 수집에 참여해야 한다, 안 그러면 목록에 데이터가 빠진다.
	tree, err := d.buildFindingAssetTree(f, 0)
	if err != nil {
		return f, err
	}
	children := map[string][]FindingAssetNode{}
	byKey := map[string]FindingAssetNode{}
	for _, n := range tree.Nodes {
		byKey[n.Key] = n
		children[n.Parent] = append(children[n.Parent], n)
	}
	if _, ok := byKey[scope]; !ok {
		// 선택한 노드가 현재 필터에서 이미 존재하지 않으면, 결과는 비어야 하며 필터 없음으로 퇴화해서는 안 된다.
		f.assetMiss = true
		return f, nil
	}
	seen := map[string]bool{scope: true}
	queue := []string{scope}
	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]
		if id := byKey[key].AssetID; id > 0 {
			f.assetIDs = append(f.assetIDs, id)
		}
		for _, child := range children[key] {
			if seen[child.Key] {
				continue
			}
			seen[child.Key] = true
			queue = append(queue, child.Key)
		}
	}
	if len(f.assetIDs) == 0 {
		f.assetMiss = true
	}
	return f, nil
}

// assetIDContainments는 자산 id를 jsonb 포함 판정의 우변 피연산자 집합으로 변환하며,
// idx_findings_asset_ids(GIN jsonb_path_ops)와 함께 사용한다.
func assetIDContainments(ids []int64) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, "["+strconv.FormatInt(id, 10)+"]")
	}
	return out
}
