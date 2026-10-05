package agent

import (
	"context"
	"encoding/json"

	actool "github.com/Autumn-27/norma/tool"
)

// 이 파일은 '내장 도구'를 순수 코드에서 열거 및 DB 재정의가 가능한 카탈로그로 전환합니다:
//   - BuiltinToolSeeds(): 세 실행 agent의 내장 도구 모음을 key +
//     설명 + 파라미터 schema + 기본 바인딩 agent로 구성된 seed 레코드로 펼쳐, 서버 시작 시 tools 테이블에 멱등으로 시드합니다.
//   - ToolResolve 훅: 런타임에 DB의 tools 행을 기준으로 조립된 도구를 'agent별 필터링 +
//     설명/schema 재정의 + 파라미터 기본값 주입' 처리합니다. key/handler는 계속 코드 계층에 있으며, DB는 '산문과 기본값'만 변경합니다.
// handler(Call 동작)는 항상 코드에서 가져옵니다. DB에서는 이를 바꿀 수 없으며 모델에 보이는 설명과 기본 입력 파라미터만 바꿀 수 있습니다.

// ToolSeed는 내장 도구의 시드 가능한 스냅샷입니다. key는 CoreTool.Name()이며(handler와 고정 바인딩,
// UI 읽기 전용), Desc/Schema는 코드의 도구 정의에서 가져오고 Agents는 코드에서 기본으로 할당한 agent 목록입니다.
type ToolSeed struct {
	Key    string         // = CoreTool.Name(), 기본 키, 변경 불가
	Desc   string         // 최상위 설명(UI에서 재정의 가능)
	Schema map[string]any // 파라미터 JSON-Schema(구조는 읽기 전용, description/default는 UI에서 변경 가능)
	Agents []string       // 기본으로 바인딩된 agent key(worker/planner/mainagent)
}

// builtinToolsByAgent는 '읽기 전용 빈 틀' ToolSet(nil stores)을 사용하여 각 실행 agent의
// 도메인 도구 모음을 구성합니다. 도구 생성자는 클로저를 Spec에 넣기만 하고 생성 중 store를 역참조하지 않으므로 nil이어도 안전합니다.
// 이 도구는 여기서 Name()/Description()/InputSchema()를 읽는 데만 사용하며 절대 Call하지 않습니다.
//
// SDK 공통 도구 actool.DefaultTools()(Read/Write/Edit/MultiEdit/LS/Glob/
// Grep/Bash)는 의도적으로 [포함하지 않습니다]. 모든 agent가 항상 보유하여 '누구에게 바인딩할지' 선택할 필요가 없고, 설명도 대부분 Prompt()
// 안에 있기 때문입니다(이 테이블은 Description()만 재정의하므로 일부만 덮어써서 오해를 일으킬 수 있음). seed 없음 → DB 행 없음 → ToolResolve가
// 변경 없이 통과시키며, 이전과 동일하게 동작합니다. artex 자체 도메인 도구만 테이블에 등록하여 관리합니다.
func builtinToolsByAgent() map[string][]actool.CoreTool {
	ts := NewToolSet(nil, "")
	return map[string][]actool.CoreTool{
		"mainagent": ts.MainAgentTools(),
		"planner":   ts.PlannerTools(),
		"worker":    ts.WorkerTools(),
		// goals(목표 분해기)는 기본적으로 set_goals + set_constraints에 바인딩됩니다. 이를 통해 분해한 목표와
		// 추출한 동작 제약 조건을 저장소에 기록합니다. mainagent와 동일한 관리 대상 도구를 공유하며, web에서 설명/schema를 변경하고 agent별로 선택할 수 있습니다.
		"goals": {ts.setGoals(), ts.setConstraints()},
		// auto는 기본적으로 취약점 보고 + 자산 관리 도구에 바인딩되며, 다른 도메인 도구는 UI에서 필요에 따라 선택할 수 있습니다.
		// 새 DB에는 이 seed로 기록하고, 기존 DB는 seedAutoDefaultBindings로 마이그레이션합니다.
		"auto": {ts.addFinding(), ts.insertAssets(), ts.addCompanyScope(), ts.listAssets(), ts.listCompanies()},
		// pentest(독립 침투 테스트 agent)의 기본 바인딩: 자산 조회 / 자산 등록 / 취약점 보고 / 취약점 조회 / 회사 조회.
		// 새 DB에는 이 seed로 기록하고, 기존 DB는 seedPentestDefaultBindings로 마이그레이션합니다.
		"pentest": {ts.listAssets(), ts.insertAssets(), ts.addFinding(), ts.listFindings(), ts.listCompanies()},
	}
}

// defaultUnbound: 이 system 도구는 평소처럼 카탈로그에 등록되지만(web에서 표시되고 agent별로 직접 선택 가능),
// 기본적으로 [어떤 agent에도 바인딩하지 않습니다]. ToolResolve는 바인딩이 비어 있는 도구를 모든 agent에서 제거하므로 명시적인 opt-in이 필요합니다.
// 그래도 특정 agent의 base 도구 모음에 두는 이유(예: PlannerTools의 goal_met)는 첫째, seed가
// 도구를 구성하여 desc/schema를 가져올 수 있고, 둘째, 사용자가 직접 다시 바인딩하면 런타임 base에 도구가 있어 ToolResolve가 유지할 수 있기 때문입니다.
//
// goal_met: 각 prove_goal을 우회하고 전역에서 [전체 작업 완료]를 바로 선언합니다. 영향이 크고 오판 위험이 있으며,
// 'prove_goal이 마지막 목표를 표시 → 자동 완료'와도 중복되므로 기본적으로 어떤 agent에도 할당하지 않고 필요할 때 직접 바인딩합니다.
var defaultUnbound = map[string]bool{"goal_met": true}

// BuiltinToolSeeds는 각 agent의 내장 도구 모음을 중복 제거하여 seed 목록으로 병합합니다. 같은 이름의 도구(예: 여러 agent가
// 보유한 list_assets)는 하나로 합치고 Agents는 합집합을 사용하며, defaultUnbound의 도구는 바인딩을 강제로 비웁니다.
func BuiltinToolSeeds() []ToolSeed {
	byAgent := builtinToolsByAgent()
	order := []string{"mainagent", "goals", "planner", "worker", "auto", "pentest"}

	type acc struct {
		tool   actool.CoreTool
		agents []string
	}
	m := map[string]*acc{}
	var keys []string
	for _, ak := range order {
		for _, t := range byAgent[ak] {
			a, ok := m[t.Name()]
			if !ok {
				a = &acc{tool: t}
				m[t.Name()] = a
				keys = append(keys, t.Name())
			}
			a.agents = append(a.agents, ak)
		}
	}

	out := make([]ToolSeed, 0, len(keys))
	for _, k := range keys {
		a := m[k]
		agents := a.agents
		if defaultUnbound[k] {
			agents = []string{} // 카탈로그에 등록하여 직접 바인딩할 수 있지만, 기본적으로 어떤 agent에도 할당하지 않음(null 대신 [] 저장, 다른 도구와 동일)
		}
		out = append(out, ToolSeed{
			Key:    k,
			Desc:   a.tool.Description(),
			Schema: a.tool.InputSchema(),
			Agents: agents,
		})
	}
	return out
}

// ToolResolve, if set, post-processes an agent's fully-assembled tool list against
// the DB tools table: it drops tools not bound to this agent (or globally disabled)
// and wraps the rest so the model sees the DB-overridden description/schema and
// 기본 입력 파라미터 get injected. Tools with no matching DB row (MCP/skill/host tools like
// traffic) pass through untouched. nil = tools unchanged. Wired in server/assembly.go.
var ToolResolve func(ctx context.Context, agentKey string, tools []actool.CoreTool) []actool.CoreTool

// DecorateTool wraps t so Description()/InputSchema() report the DB overrides and
// Call() injects scalar parameter defaults (from schema's "default" props) whenever
// the model omitted them. Name/Prompt/permission/scheduler flags delegate to t, so
// the tool's identity and handler are unchanged. Empty desc/schema fall back to t's.
func DecorateTool(t actool.CoreTool, desc string, schema map[string]any) actool.CoreTool {
	if desc == "" {
		desc = t.Description()
	}
	if len(schema) == 0 {
		schema = t.InputSchema()
	}
	return &overriddenTool{CoreTool: t, desc: desc, schema: schema}
}

// overriddenTool is a CoreTool decorator: it embeds the original (so all behavioral
// methods — Prompt/IsReadOnly/IsConcurrencySafe/CheckPermissions/Name — delegate)
// and overrides only the model-facing description/schema plus default injection.
type overriddenTool struct {
	actool.CoreTool
	desc   string
	schema map[string]any
}

func (o *overriddenTool) Description() string         { return o.desc }
func (o *overriddenTool) InputSchema() map[string]any { return o.schema }

func (o *overriddenTool) Call(ctx context.Context, in json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
	return o.CoreTool.Call(ctx, injectDefaults(in, o.schema), tc)
}

// injectDefaults fills scalar parameter defaults declared in the (possibly edited)
// schema into the input JSON whenever the model omitted the field or left it empty/
// null. Structure (names/types/required) is untouched — only 기본값 are merged in.
func injectDefaults(in json.RawMessage, schema map[string]any) json.RawMessage {
	defs := scalarDefaults(schema)
	if len(defs) == 0 {
		return in
	}
	m := map[string]json.RawMessage{}
	if len(in) > 0 {
		if err := json.Unmarshal(in, &m); err != nil {
			return in // non-object input: don't touch it
		}
	}
	changed := false
	for k, dv := range defs {
		if cur, ok := m[k]; !ok || isEmptyJSON(cur) {
			m[k] = dv
			changed = true
		}
	}
	if !changed {
		return in
	}
	b, err := json.Marshal(m)
	if err != nil {
		return in
	}
	return b
}

// scalarDefaults extracts properties[k]["default"] for scalar params (string/
// integer/number/boolean). Array/object defaults are skipped: merging them is
// ambiguous and not worth the surprise.
func scalarDefaults(schema map[string]any) map[string]json.RawMessage {
	props, _ := schema["properties"].(map[string]any)
	if len(props) == 0 {
		return nil
	}
	out := map[string]json.RawMessage{}
	for name, raw := range props {
		p, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		dv, ok := p["default"]
		if !ok || dv == nil {
			continue
		}
		switch p["type"] {
		case "string", "integer", "number", "boolean":
			if b, err := json.Marshal(dv); err == nil {
				out[name] = b
			}
		}
	}
	return out
}

func isEmptyJSON(raw json.RawMessage) bool {
	s := string(raw)
	return s == "null" || s == `""`
}
