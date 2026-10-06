package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
)

// jsonResult marshals v to a JSON tool result.
func jsonResult(v any) (actool.Result, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return actool.Errorf(err.Error()), nil
	}
	return actool.Text(string(b)), nil
}

// 이 파일은 P2 '크로스 작업 오케스트레이션 도구 집합'(docs/跑分编排 §2 P2)을 구현한다. 이들은 host 도구다 —— Manager
// (임의 작업의 Store)·Engine(일시정지)·작업 생성 흐름에 접근해야 하므로, server 레이어에 둔다.
// 읽기류 도구는 '기존 per-task 도구'를 대상 작업의 store 위에서 돌도록 리다이렉트한다(임시 ToolSet를
// 만들어 해당 도구를 Call), 이로써 완전히 동일한 로직을 재사용한다; 제어류(spawn/pause)는 Manager/Engine을 직접 호출.
// 이들은 트래픽 도구처럼 tools 테이블에 seed되고, agent별로 바인딩된다(오케스트레이션 agent에만 바인딩해 보임).

// hostTools is the runtime host-tool provider fed to ToolAugment: traffic tools
// (gated by capture) + cross-task orchestration tools + user-defined custom tools.
// The second return is the names of custom tools flagged `deferred` (schema
// withheld, routed via SearchExtraTools/ExecuteExtraTool). Per-agent binding still
// decides who actually sees any of them.
//
//nolint:unused // used as the hostTools provider in wireAgentAugment
func (s *Server) hostTools() ([]actool.CoreTool, map[string][]string) {
	tools := append(s.m.HostTools(), s.orchestrationTools()...)
	tools = append(tools, s.findingRetestTools()...)
	tools = append(tools, s.platformTools()...) // 플랫폼 동작 도구(skill/도구/MCP 생성·수정, Auto용)
	custom, err := s.customTools()
	if err != nil {
		log.Printf("[custom-tool] 로드 실패: %v", err)
		return tools, nil
	}
	tools = append(tools, custom...)
	// deferred custom tools → name -> its bound agent keys. ToolAugment turns a
	// name into a deferred entry only for agents it's actually bound to (so we don't
	// advertise a tool the per-agent binding will drop from the callable set).
	deferred := map[string][]string{}
	rows, _ := s.m.pg.ListCustomTools()
	for _, t := range rows {
		if t.Deferred && t.Enabled {
			deferred[t.Key] = t.Agents
		}
	}
	return tools, deferred
}

// orchestrationTools returns the cross-task tool set. Bound per-agent via the
// tools table (default: no binding — opt-in for orchestration agents).
func (s *Server) orchestrationTools() []actool.CoreTool {
	return []actool.CoreTool{
		s.toolListTasks(),
		s.toolListLLMProfiles(),
		s.toolSpawnTask(),
		s.toolPauseTask(),
		s.toolGetTaskGraph(),
		s.toolListTaskFindings(),
		s.toolAddHint(),
		s.toolGetWorkerTrace(),
		s.toolListWorkerTraces(),
		s.toolSearchWorkerTraces(),
		s.toolGetTaskNodeDetail(),
		s.toolUpdateFindingReport(),
		s.toolGetFindingTraffic(),
		s.toolBindFindingTraffic(),
	}
}

// --- schema helpers ---

func strParam(desc string) map[string]any {
	return map[string]any{"type": "string", "description": desc}
}

// parseProfileID reads an LLM profile id from a tool arg that may arrive as a JSON
// number (5) or a numeric string ("5"); returns 0 when absent/unparseable.
func parseProfileID(raw json.RawMessage) int64 {
	if len(raw) == 0 {
		return 0
	}
	var n int64
	if json.Unmarshal(raw, &n) == nil {
		return n
	}
	var str string
	if json.Unmarshal(raw, &str) == nil {
		v, _ := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
		return v
	}
	return 0
}

func objSchema(props map[string]any, required ...string) map[string]any {
	m := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		req := make([]any, len(required))
		for i, r := range required {
			req[i] = r
		}
		m["required"] = req
	}
	return m
}

func roTool(name, desc string, schema map[string]any, run func(context.Context, json.RawMessage) (actool.Result, error)) actool.CoreTool {
	return actool.Build(actool.Spec{
		Name: name, Description: desc, Schema: schema,
		ReadOnly:   func(json.RawMessage) bool { return true },
		Concurrent: func(json.RawMessage) bool { return true },
		Permissions: func(context.Context, json.RawMessage, permission.Context) permission.Decision {
			return permission.Allowed()
		},
		Run: func(ctx context.Context, in json.RawMessage, _ *actool.ToolContext) (actool.Result, error) {
			return run(ctx, in)
		},
	})
}

func wrTool(name, desc string, schema map[string]any, run func(context.Context, json.RawMessage) (actool.Result, error)) actool.CoreTool {
	return actool.Build(actool.Spec{
		Name: name, Description: desc, Schema: schema,
		Permissions: func(context.Context, json.RawMessage, permission.Context) permission.Decision {
			return permission.Allowed()
		},
		Run: func(ctx context.Context, in json.RawMessage, _ *actool.ToolContext) (actool.Result, error) {
			return run(ctx, in)
		},
	})
}

// delegateToTask resolves the `task_id` in the input, builds a ToolSet bound to
// that task's store, strips task_id, and calls the chosen per-task tool — so the
// cross-task read reuses the exact in-task logic against another task.
func (s *Server) delegateToTask(ctx context.Context, in json.RawMessage, pick func(*agent.ToolSet) actool.CoreTool) (actool.Result, error) {
	var head struct {
		TaskID string `json:"task_id"`
	}
	_ = json.Unmarshal(in, &head)
	if strings.TrimSpace(head.TaskID) == "" {
		return actool.Errorf("task_id는 필수입니다"), nil
	}
	t, ok := s.m.Task(head.TaskID)
	if !ok {
		return actool.Errorf("task가 존재하지 않습니다: " + head.TaskID), nil
	}
	var m map[string]json.RawMessage
	_ = json.Unmarshal(in, &m)
	delete(m, "task_id")
	inner, _ := json.Marshal(m)
	tsx := agent.NewToolSet(t.Store, "orchestrator")
	if s.m.Assets() != nil {
		tsx.SetAssetStore(s.m.Assets(), s.m.Assets().Companies())
	}
	tsx.SetNotify(t.Notify)         // 범용 깨우기(전용 콜백이 없는 쓰기 작업은 이것을 사용; 읽기 도구는 no-op)
	tsx.SetNotifyHint(t.NotifyHint) // add_hint → '사람이 전략 힌트 N개를 추가함: …' 트리거 하나를 기록하고 planner를 깨움
	return pick(tsx).Call(ctx, inner, nil)
}

// --- tools ---

func (s *Server) toolListTasks() actool.CoreTool {
	return roTool("list_tasks",
		"모든 작업을 나열한다(id/설명/목표/상태/실행 시간/상위 작업/LLM 설정), 오케스트레이션 agent가 이것으로 전체를 파악하고, 어떤 작업이 너무 오래 막혔는지·각자 어떤 LLM을 쓰는지 본다. 실행 시간: 실행 중=생성→현재, 종료 상태=생성→마지막 활동(초). llm_profile: 작업 planner/worker가 쓰는 설정 이름, (활성 설정)=전역 활성을 따름.",
		objSchema(map[string]any{}),
		func(context.Context, json.RawMessage) (actool.Result, error) {
			lastAct, _ := s.m.PG().LastActivityAll()
			// id -> name to resolve each task's pinned LLM profile.
			profName := map[int64]string{}
			if profs, err := s.m.pg.ListProfiles(); err == nil {
				for _, p := range profs {
					profName[p.ID] = p.Name
				}
			}
			out := make([]map[string]any, 0)
			for _, t := range s.m.List() {
				status := s.deriveTaskStatus(t)
				end := lastAct[t.ExpID]
				if live := s.engine.LastActivity(t.ID); live > end {
					end = live
				}
				dur := int64(0)
				if status == "running" {
					dur = time.Now().Unix() - t.CreatedAt
				} else if end > t.CreatedAt {
					dur = end - t.CreatedAt
				}
				row := map[string]any{"id": t.ID, "description": t.Description, "goal": t.Goal, "status": status, "run_seconds": dur}
				if t.ParentRef != "" {
					row["parent_ref"] = t.ParentRef
				}
				llmState := t.llmStateSnapshot()
				if llmState.ProfileID == nil {
					row["llm_profile"] = "(활성 설정)"
				} else if n, ok := profName[*llmState.ProfileID]; ok {
					row["llm_profile"] = n
				} else {
					row["llm_profile"] = fmt.Sprintf("#%d(삭제됨)", *llmState.ProfileID)
				}
				out = append(out, row)
			}
			return jsonResult(out)
		})
}

// toolListLLMProfiles lists the available LLM profiles (name/model/active) so an
// orchestration agent can pick one for spawn_task's llm_profile. Never leaks keys.
func (s *Server) toolListLLMProfiles() actool.CoreTool {
	return roTool("list_llm_profiles",
		"사용 가능한 LLM 설정(profile)을 나열한다: id, 이름, 모델, 형식, 현재 활성 설정 여부. id로 spawn_task의 llm_profile_id 파라미터에 하위 작업 전용 LLM을 지정한다(예: 정찰은 저렴한 모델, 익스플로잇은 강한 모델). API Key는 포함하지 않는다.",
		objSchema(map[string]any{}),
		func(context.Context, json.RawMessage) (actool.Result, error) {
			profs, err := s.m.pg.ListProfiles()
			if err != nil {
				return actool.Errorf(err.Error()), nil
			}
			out := make([]map[string]any, 0, len(profs))
			for _, p := range profs {
				out = append(out, map[string]any{
					"id": p.ID, "name": p.Name, "model": p.Model, "format": p.Format, "is_active": p.IsDefault,
				})
			}
			return jsonResult(map[string]any{"profiles": out})
		})
}

func (s *Server) toolSpawnTask() actool.CoreTool {
	return wrTool("spawn_task",
		"하위 작업 하나를 새로 만들고 탐색 엔진을 시작한다, task_id를 반환. 어떤 일(예: 한 문제/한 목표)을 독립 작업으로 보내는 데 쓴다. parent_ref 선택: 현재 오케스트레이션이 연결된 상위 작업 id를 넣어 상·하위 연결.",
		objSchema(map[string]any{
			"description":            strParam("작업 설명(간단한 제목)"),
			"goal":                   strParam("작업 목표(무엇을 달성할지)"),
			"parent_ref":             strParam("선택: 상위 작업 id(상·하위 연결)"),
			"source_task_ids":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": fmt.Sprintf("선택: 읽기 전용 상속 출처 작업 id 목록(최대 %d개). 하위 작업은 이 작업들이 이미 밝혀낸 자산/결론을 읽기 전용으로 참조해 출발점으로 삼을 수 있다; parent_ref의 순수 상·하위 포인터와 달리, 이것은 내용 상속이다.", db.MaxTaskSourceCount)},
			"llm_profile_id":         map[string]any{"type": "integer", "description": "선택: 이 하위 작업의 planner/worker가 쓸 LLM 설정 id 지정(list_llm_profiles 참고); 비우면 상위 작업을 상속하고, 다시 전역 활성 설정으로 폴백"},
			"timeout_seconds":        map[string]any{"type": "integer", "description": "선택: 작업 레벨 타임아웃(초). 시점이 되면 우아한 마무리를 트리거하고 timeout 종료 상태에 진입; 비우거나 0 = 무제한"},
			"plan_heartbeat_seconds": map[string]any{"type": "integer", "description": "선택: planner 하트비트 트리거 간격(초). 지난 계획 종료/작업 시작으로부터 이 값만큼 지나고 그 사이 트리거가 없으면 → 계획 한 라운드 트리거(데드락 폴백 + 진행 중인 worker 감독을 위한 깨움). 비우거나 0 = 기본 600(10min);"},
			"seed_first_intent":      map[string]any{"type": "boolean", "description": "선택: 간단한 작업에는 켤 수 있음, 생성 시 시드 의도 하나를 바로 내려(내용=설명+목표) worker가 첫 라운드 planner를 안 기다리고 바로 테스트를 시작; 기본 false(표준인 선계획 후실행)."},
		}, "description", "goal"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			var a struct {
				Description          string          `json:"description"`
				Goal                 string          `json:"goal"`
				ParentRef            string          `json:"parent_ref"`
				SourceTaskIDs        []string        `json:"source_task_ids"`
				LLMProfileID         json.RawMessage `json:"llm_profile_id"`
				TimeoutSeconds       int             `json:"timeout_seconds"`
				PlanHeartbeatSeconds int             `json:"plan_heartbeat_seconds"`
				SeedFirstIntent      bool            `json:"seed_first_intent"`
			}
			_ = json.Unmarshal(in, &a)
			if strings.TrimSpace(a.Description) == "" {
				a.Description = "이름 없는 작업"
			}
			if strings.TrimSpace(a.Goal) == "" {
				return actool.Errorf("goal은 필수입니다"), nil
			}
			if a.TimeoutSeconds < 0 {
				a.TimeoutSeconds = 0
			}
			// 읽기 전용 상속 출처 작업: 개수 상한 + 각 id 유효/중복 제거/존재, 검증 규칙은 HTTP 작업 생성과 동일.
			if len(a.SourceTaskIDs) > db.MaxTaskSourceCount {
				return actool.Errorf(fmt.Sprintf("연관 작업은 최대 %d개까지 선택할 수 있습니다", db.MaxTaskSourceCount)), nil
			}
			sourceIDs := make([]int64, 0, len(a.SourceTaskIDs))
			seenSources := map[int64]bool{}
			for _, raw := range a.SourceTaskIDs {
				id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
				if err != nil || id <= 0 || seenSources[id] {
					return actool.Errorf("연관 작업 id가 유효하지 않거나 중복됩니다"), nil
				}
				if _, ok := s.m.Task(strconv.FormatInt(id, 10)); !ok {
					return actool.Errorf(fmt.Sprintf("연관 작업 #%d가 존재하지 않습니다", id)), nil
				}
				seenSources[id] = true
				sourceIDs = append(sourceIDs, id)
			}
			// LLM profile resolution: explicit id > inherit parent's pin > active(nil).
			var pin *int64
			if id := parseProfileID(a.LLMProfileID); id > 0 {
				if _, ok := s.loadProfileConfig(id); !ok {
					return actool.Errorf(fmt.Sprintf("LLM 설정 #%d가 없거나 API Key가 설정되지 않았습니다", id)), nil
				}
				pin = &id
			} else if a.ParentRef != "" {
				if pt, ok := s.m.Task(a.ParentRef); ok {
					pin = pt.LLMProfileID
				}
			}
			var llmIDs []int64
			if pin != nil {
				llmIDs = []int64{*pin}
			}
			t, err := s.m.CreateTaskWithOptions(a.Description, a.Goal, db.TaskCreateOptions{
				SourceTaskIDs:        sourceIDs,
				LLMProfileIDs:        llmIDs,
				TimeoutSeconds:       a.TimeoutSeconds,
				PlanHeartbeatSeconds: a.PlanHeartbeatSeconds,
			})
			if err != nil {
				return actool.Errorf(err.Error()), nil
			}
			if a.ParentRef != "" {
				t.ParentRef = a.ParentRef
				if id, e := strconv.ParseInt(t.ID, 10, 64); e == nil {
					_ = s.m.PG().SetParentRef(id, a.ParentRef)
				}
			}
			// 공유하는 생성 후 흐름, HTTP 작업 생성(server.go createTask)과 같은 launchTask를 재사용:
			// seed + 백그라운드에서 보이게 목표 분해(제0라운드/LLM단계/goal 하나씩) + engine.Run.
			// seed_first_intent 기본 false(표준인 선계획 후실행); 간단한 작업은 켜서 work 하나를 바로 내려 테스트.
			s.launchTask(t, a.Description+" "+a.Goal, a.SeedFirstIntent)
			return actool.Text(fmt.Sprintf("task created: %s", t.ID)), nil
		})
}

func (s *Server) toolPauseTask() actool.CoreTool {
	return wrTool("pause_task", "지정한 작업을 일시정지한다(그 planner/worker 루프를 멈춤).",
		objSchema(map[string]any{"task_id": strParam("일시정지할 작업 id")}, "task_id"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			var a struct {
				TaskID string `json:"task_id"`
			}
			_ = json.Unmarshal(in, &a)
			t, ok := s.m.Task(a.TaskID)
			if !ok {
				return actool.Errorf("task가 존재하지 않습니다: " + a.TaskID), nil
			}
			if _, err := s.applyTaskControlWithCause(t, "pause", agent.AbortPausedByOrchestrator); err != nil {
				return actool.Errorf(err.Error()), nil
			}
			return actool.Text("task paused: " + a.TaskID), nil
		})
}

func (s *Server) toolGetTaskGraph() actool.CoreTool {
	return roTool("get_task_graph", "지정한 작업의 탐색 그래프 개요를 읽는다(graph_overview와 동일: 자산 수/frontier/발견/커버리지 등), task_id로 작업 지정.",
		objSchema(map[string]any{"task_id": strParam("작업 id")}, "task_id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).GraphOverviewTool)
		})
}

func (s *Server) toolListTaskFindings() actool.CoreTool {
	return roTool("list_task_findings", "지정한 작업의 확인된 취약점을 읽는다(flag/PoC 포함; 각 항목에 id/task_id/intent_id/vulnclass/severity/요약/상태), task_id로 작업 지정.",
		objSchema(map[string]any{"task_id": strParam("작업 id")}, "task_id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).ListFindingsTool)
		})
}

func (s *Server) toolAddHint() actool.CoreTool {
	return wrTool("add_task_hint", "지정한 작업에 전략 힌트를 주입한다(그 작업의 planner가 다음 라운드에 의도를 생성할 때 읽는다).\n"+
		"★배치 우선: 여러 힌트를 hints 배열에 넣어 한 번에 제출(ids 배열 반환, hints와 같은 길이·같은 순서, 실패 항목 id=0); 단건이면 hints를 생략하고 최상위 text에 바로 준다.",
		objSchema(map[string]any{
			"task_id":      strParam("작업 id"),
			"hints":        map[string]any{"type": "array", "description": "[이것을 우선 사용] 힌트 배열, 각 요소 필드는 최상위와 동일(text/asset_ids/traffic_refs).", "items": objSchema(map[string]any{"text": strParam("힌트 내용"), "asset_ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer"}}, "traffic_refs": agent.HintTrafficSchema()})},
			"text":         strParam("[단건] 힌트 내용"),
			"traffic_refs": agent.HintTrafficSchema(),
			"asset_ids":    map[string]any{"type": "array", "items": map[string]any{"type": "integer"}, "description": "앵커로 삼는 자산 id(선택, 0/1/여러 개; 그 작업 내의 자산 id)"},
		}, "task_id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).AddHintTool)
		})
}

func (s *Server) toolGetWorkerTrace() actool.CoreTool {
	return roTool("get_task_worker_trace",
		"지정한 작업 안 어떤 work(의도)의 실행 과정을 본다: get_task_worker_trace(task_id, intent_id)로 단계 요약을 보고; step_ids=[...]를 더해 그 단계들의 전체 내용을 가져온다(한 번에 최대 5개, 더 넘겨도 앞 5개만 반환).",
		objSchema(map[string]any{
			"task_id":   strParam("작업 id"),
			"intent_id": map[string]any{"type": "integer", "description": "의도 id(그 작업 안의 work)"},
			"step_ids":  map[string]any{"type": "array", "items": map[string]any{"type": "integer"}, "description": "선택: 전체 내용을 가져올 단계 id(한 번에 최대 5개, 더 넘겨도 앞 5개만 반환, 나머지는 omitted_step_ids에 나열)"},
		}, "task_id", "intent_id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).GetWorkerTraceTool)
		})
}

func (s *Server) toolListWorkerTraces() actool.CoreTool {
	return roTool("list_task_worker_traces", "지정한 작업에서 어떤 work(의도)가 실행됐는지 + 각자 스텝 수를 나열한다, 어떤 work를 들여다볼 가치가 있는지 찾는 데 쓴다(이어서 get_task_worker_trace 사용).",
		objSchema(map[string]any{"task_id": strParam("작업 id")}, "task_id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).ListWorkerTracesTool)
		})
}

func (s *Server) toolSearchWorkerTraces() actool.CoreTool {
	return roTool("search_task_worker_traces", "지정한 작업에서 키워드로 모든 work의 실행 과정을 검색한다(매칭된 단계 요약 + intent_id 반환).",
		objSchema(map[string]any{"task_id": strParam("작업 id"), "q": strParam("검색 키워드")}, "task_id", "q"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).SearchWorkerTracesTool)
		})
}

func (s *Server) toolGetTaskNodeDetail() actool.CoreTool {
	return roTool("get_task_node_detail",
		"지정한 작업 안 어떤 탐색 그래프 노드의 전체 내용을 읽는다(발견/사실/의도/목표: 요약 + 상세/증거/PoC). id는 탐색 노드 id(예: report_finding 반환, 또는 list_task_findings 안의 id). 취약점 보고서를 쓰기 전에 이것으로 그 취약점의 전체 증거를 가져온다.",
		objSchema(map[string]any{
			"task_id": strParam("작업 id"),
			"id":      map[string]any{"type": "integer", "description": "탐색 그래프 노드 id(자산 id 아님)"},
		}, "task_id", "id"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			return s.delegateToTask(ctx, in, (*agent.ToolSet).NodeDetailTool)
		})
}

// toolUpdateFindingReport writes/overwrites a finding's detailed Markdown report.
// finding_id is the id report_finding returned ("finding recorded: <id>", the
// finding node id). The write (SetFindingReportByNodeID) is keyed by node_id and
// task-agnostic, so this host tool needs no task_id / exploration store.
func (s *Server) toolUpdateFindingReport() actool.CoreTool {
	return wrTool("update_finding_report",
		"이미 등록된 취약점에 [상세 보고서]를 쓰거나 갱신한다(Markdown 전문, 전체를 통째로 덮어씀). finding_id에는 report_finding이 반환한 그 id를 전달한다(\"finding recorded: <id>\" 안의 숫자). 보고서에는 취약점 개요, 영향과 위험, 재현 단계, 증거/PoC, 수정 권고를 포함할 것을 권장한다.",
		objSchema(map[string]any{
			"finding_id":       map[string]any{"type": "integer", "description": "대상 취약점 id(report_finding 반환 id)"},
			"report":           strParam("상세 보고서 전문, Markdown 형식"),
			"evidence_version": map[string]any{"type": "integer", "description": "get_finding_traffic이 반환한 증거 version; 보고서가 새 증거 변경을 덮어쓰는 것을 방지"},
		}, "finding_id", "report"),
		func(ctx context.Context, in json.RawMessage) (actool.Result, error) {
			var a struct {
				EvidenceVersion *int64          `json:"evidence_version"`
				FindingID       json.RawMessage `json:"finding_id"`
				Report          string          `json:"report"`
			}
			_ = json.Unmarshal(in, &a)
			nodeID := parseProfileID(a.FindingID) // '숫자 또는 숫자 문자열' 해석 재사용
			if nodeID <= 0 {
				return actool.Errorf("finding_id가 유효하지 않습니다"), nil
			}
			n, err := s.m.pg.SetFindingReportVersionByNodeID(ctx, nodeID, a.Report, a.EvidenceVersion)
			if err != nil {
				return actool.Errorf(err.Error()), nil
			}
			if n == 0 {
				return actool.Errorf(fmt.Sprintf("finding_id=%d에 해당하는 취약점 기록을 찾을 수 없습니다(먼저 report_finding으로 등록하세요)", nodeID)), nil
			}
			return actool.Text(fmt.Sprintf("finding %d report updated (%d chars)", nodeID, len(a.Report))), nil
		})
}

// deriveTaskStatus mirrors listTasks' status derivation for the list_tasks tool.
func (s *Server) deriveTaskStatus(t *Task) string {
	lifecycle := t.lifecycleSnapshot()
	switch {
	case isTerminalStatus(lifecycle.Status):
		return lifecycle.Status
	case lifecycle.Paused || s.engine.IsPaused(t.ID):
		return "paused"
	case s.engine.ReadyFor(t) && s.engine.Started(t.ID):
		return "running"
	}
	return "created"
}

// orchestrationToolSeeds seeds the cross-task tools into the tools table so they
// are bindable per-agent (default: bound to nobody — opt-in for orchestration
// agents). First-insert only, like the traffic seeds.
func (s *Server) seedOrchestrationTools() {
	// task-op + platform tools default-bind to the built-in Auto agent (플랫폼
	// 조작을 위해 태생적으로 쓰임). SeedTool 첫 삽입에 적용; 구 DB에 이미 seed된 행은 seedAutoDefaultBindings가 보강 바인딩.
	autoAgents, _ := json.Marshal([]string{"auto"})
	for _, t := range s.orchestrationTools() {
		schema, _ := json.Marshal(t.InputSchema())
		bindings := autoAgents
		if t.Name() == "bind_finding_traffic" {
			bindings = json.RawMessage(`["reporter"]`)
		}
		_ = s.m.PG().SeedTool(t.Name(), t.Description(), schema, bindings)
	}
	for _, t := range s.platformTools() {
		schema, _ := json.Marshal(t.InputSchema())
		_ = s.m.PG().SeedTool(t.Name(), t.Description(), schema, autoAgents)
	}
	s.refreshBuiltinToolSchemas()
	s.seedAutoDefaultBindings()
	s.seedPlannerDefaultBindings()
	s.seedPlannerListAssetsBinding()
	s.seedCompanyScopeRebind()
	s.seedWorkerReadToolsUnbind() // list_facts/list_companies/list_worker_traces를 worker 기본에서 바인딩 해제(일회성)
	s.seedWorkerReadbackRebind()  // 구 마이그레이션 오삭제 수정: search_all_worker_traces/get_worker_trace/node_detail을 worker에 다시 보강 바인딩(일회성)
	s.seedAutoReportFindingBinding()
	s.unbindGoalMetDefault()
	s.reseedGoalsPrompt()             // goals 프롬프트에 '동작 제약 추출' 단계 추가 → 구 DB에 새 기본값 한 버전 추가(일회성)
	s.reseedMainAgentPrompt()         // mainagent 프롬프트에 '목표 달성 후 add_intent 시 목표 등록 여부 되물음' 추가(일회성)
	s.reseedPlannerPrompt()           // planner 프롬프트: '0 의도' 정당 사유 재작성 + 정량 검수 점검 추가(일회성)
	s.reseedWorkerPrompt()            // worker 프롬프트: 부정 결론 증거 기준 추가(일회성)
	s.seedReporterAgent()             // '보고서 작성' agent + 도구 바인딩 + finding 트리거 사전 배치(일회성)
	s.upgradeReporterTriggerMessage() // 구 DB 보강 마이그레이션: reporter가 evidence_version을 되돌려 전달하게 함(일회성)
	s.seedFindingTrafficTools()       // 선택적 증거 파라미터 및 읽기 전용 증거 도구 추가, 사용자 설정 보존
	s.seedFindingWorkflowTools()
	// 주: pentest의 기본 도구 바인딩은 마이그레이션 불필요 —— BuiltinToolSeeds가 완전 새 초기화 때
	// list_assets/insert_assets/report_finding/list_findings/list_companies를
	// pentest와 함께 seed해 둔다(프로젝트에 아직 구 DB가 없어 마이그레이션하지 않음).
}

// refreshBuiltinToolSchemas propagates code schema/description changes on the
// orchestration + platform tools into already-seeded rows ONCE per version flag —
// SeedTool is first-insert-only, so a new param (e.g. spawn_task의 llm_profile) never
// reaches an old DB otherwise. Preserves each tool's agent binding + enabled flag.
// Bump the flag whenever these tools' schemas/descriptions change in code.
func (s *Server) refreshBuiltinToolSchemas() {
	const flag = "tool_schema_refresh_v7_list_facts_paging"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	tools := append(s.orchestrationTools(), s.platformTools()...)
	for _, t := range tools {
		schema, _ := json.Marshal(t.InputSchema())
		if err := s.m.pg.RefreshToolDefaults(t.Name(), t.Description(), schema); err != nil {
			log.Printf("[tools] refresh %s schema failed: %v", t.Name(), err)
		}
	}
	// 동시에 일부 내장 agent 도구를 코드 기본값으로 갱신:
	//   - goal_met: 구 DB가 seed한 설명에 '이번 계획 라운드 종료'라는 오해 소지가 있어, planner가 그것을
	//     '빈 라운드 종료' 수단으로 여겨 막 시작하자마자 작업 전체를 완료로 오판하게 만든다.
	//   - insert_assets: related 입력 파라미터 추가(자산이 현재 작업과 관련 있는지 표시, 커버리지 포함 여부 결정),
	//     SeedTool 첫 삽입only이라, 구 DB에 이미 seed된 schema는 이 새 파라미터를 받지 못한다.
	//   - list_facts: 페이지네이션으로 변경, limit/before/q 입력 파라미터 추가; 구 DB에 이미 seed된 빈 schema는
	//     도구 관리 페이지에서 '파라미터 없음'으로 표시되고, 모델도 이 파라미터 설명들을 받지 못한다.
	refreshBuiltin := map[string]bool{"goal_met": true, "insert_assets": true, "list_facts": true}
	for _, sd := range agent.BuiltinToolSeeds() {
		if !refreshBuiltin[sd.Key] {
			continue
		}
		schema, _ := json.Marshal(sd.Schema)
		if err := s.m.pg.RefreshToolDefaults(sd.Key, sd.Desc, schema); err != nil {
			log.Printf("[tools] refresh %s desc failed: %v", sd.Key, err)
		}
	}
	_ = s.m.pg.SetSetting(flag, "true")
	log.Printf("[tools] orchestration/platform 도구 schema를 코드 기본값으로 갱신 완료(일회성)")
}

// unbindGoalMetDefault removes goal_met's default "planner" binding ONCE (guarded by
// a settings flag), so existing DBs match the new default of NO agent. goal_met bypasses
// per-goal prove_goal to declare the whole task done — powerful/risky and redundant with
// the prove_goal→auto-complete path — so it ships unbound; users can re-bind it per agent
// in the UI. A user's own binding to another agent is untouched (we only strip planner).
func (s *Server) unbindGoalMetDefault() {
	const flag = "goal_met_unbind_default_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.RemoveAgentFromTool("planner", "goal_met"); err != nil {
		log.Printf("[tools] goal_met planner 바인딩 해제 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// reseedGoalsPrompt는 goals 목표 분해기의 프롬프트를 [현재 코드 기본값]으로 갱신한다 —— 기본 본문에
// '먼저 동작 제약 추출(set_constraints) 후 목표 분해'라는 단계가 추가됐는데, SeedPromptIfEmpty는 첫 삽입only라, 구 DB에
// 이미 있는 version 1은 이 단계를 받지 못한다. 여기서는 버전 관리로 [새 버전 하나 추가]하고 전환한다(ResetPromptToDefault),
// 구 버전은 이력에 남으므로, 사용자가 커스텀했다면 버전 기록에서 되찾을 수 있다. settings flag 가드 → 한 번만;
// 이후 기본값이 또 바뀌면 이 flag를 bump. 완전 새 DB는 처리 불필요(SeedPromptIfEmpty가 이미 최신 기본값 seed).
func (s *Server) reseedGoalsPrompt() {
	const flag = "goals_prompt_constraint_step_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 성공 여부와 무관하게 한 번만 시도
	a, err := s.m.pg.GetAgentByKey("goals")
	if err != nil || a == nil {
		return // 완전 새 DB에서 agent 행이 아직 없을 때는 seedPrompts가 최신 기본값을 바로 seed하므로 이 마이그레이션 불필요
	}
	tmpl := agent.BuiltinPromptSeeds()["goals"]
	if tmpl == "" {
		return
	}
	// 완전 새 DB는 seedPrompts가 이미 최신 기본값을 seed함 → 현재 버전이 코드 기본값과 같으므로 중복 버전을 추가할 필요 없음.
	if cur, err := s.m.pg.CurrentPrompt(a.ID); err == nil && cur == tmpl {
		return
	}
	if _, err := s.m.pg.ResetPromptToDefault(a.ID, tmpl); err != nil {
		log.Printf("[prompts] goals 프롬프트를 새 기본값으로 재갱신 실패: %v", err)
		return
	}
	log.Printf("[prompts] goals 프롬프트에 새 기본값 버전 추가 완료(동작 제약 추출 단계 추가, 일회성)")
}

// reseedMainAgentPrompt는 mainagent 프롬프트를 [현재 코드 기본값]으로 갱신한다 —— 기본 본문에 '목표가 전부
// 달성된 뒤 add_intent로 의도를 직접 투입할 때, 정식 목표로 등록할지 사람에게 되물음'이라는 안내가 추가됐는데, SeedPromptIfEmpty는 첫 삽입
// only라, 구 DB의 기존 버전은 받지 못한다. 버전 관리로 [새 버전 하나 추가]하고 전환한다(ResetPromptToDefault), 구 버전은
// 이력에 남으므로, 사용자가 커스텀했다면 버전 기록에서 되찾을 수 있다. settings flag 가드 → 한 번만. 완전 새 DB는 처리 불필요
// (SeedPromptIfEmpty가 이미 최신 기본값 seed). reseedGoalsPrompt와 완전히 동형.
func (s *Server) reseedMainAgentPrompt() {
	const flag = "mainagent_prompt_goalless_intent_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 성공 여부와 무관하게 한 번만 시도
	a, err := s.m.pg.GetAgentByKey("mainagent")
	if err != nil || a == nil {
		return // 완전 새 DB에서 agent 행이 아직 없을 때는 seedPrompts가 최신 기본값을 바로 seed하므로 이 마이그레이션 불필요
	}
	tmpl := agent.BuiltinPromptSeeds()["mainagent"]
	if tmpl == "" {
		return
	}
	// 완전 새 DB는 seedPrompts가 이미 최신 기본값을 seed함 → 현재 버전이 코드 기본값과 같으므로 중복 버전을 추가할 필요 없음.
	if cur, err := s.m.pg.CurrentPrompt(a.ID); err == nil && cur == tmpl {
		return
	}
	if _, err := s.m.pg.ResetPromptToDefault(a.ID, tmpl); err != nil {
		log.Printf("[prompts] mainagent 프롬프트를 새 기본값으로 재갱신 실패: %v", err)
		return
	}
	log.Printf("[prompts] mainagent 프롬프트에 새 기본값 버전 추가 완료(목표 달성 후 목표 등록 되물음 추가, 일회성)")
}

// reseedPlannerPrompt는 planner 프롬프트를 [현재 코드 기본값]으로 갱신한다 —— 기본 본문을 간소화 재구성했고, '절제'를
// 중복 제거만으로 강등, '깊이가 커버리지보다 우선'·'하드 하한선: 목표 미달성이고 실행 중 의도가 없으면 반드시 산출'을 추가, 부정 결론 재검토에 상한을 더했다.
// 기본값에 실질 변경이 있을 때마다 아래 flag를 bump(현재 v2)해 기존 구 DB를 다시 한 번 갱신한다. SeedPromptIfEmpty는 첫 삽입only라, 구 DB의 기존 버전은 받지 못하므로 버전 관리로
// [새 버전 하나 추가]하고 전환한다(ResetPromptToDefault), 구 버전은 이력에 남으므로, 사용자가 커스텀했다면 버전 기록에서
// 되찾을 수 있다. settings flag 가드 → 한 번만. 완전 새 DB는 처리 불필요(SeedPromptIfEmpty가 이미 최신 기본값 seed).
// reseedGoalsPrompt와 완전히 동형.
func (s *Server) reseedPlannerPrompt() {
	const flag = "planner_prompt_compact_realistic_v2"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 성공 여부와 무관하게 한 번만 시도
	a, err := s.m.pg.GetAgentByKey("planner")
	if err != nil || a == nil {
		return // 완전 새 DB에서 agent 행이 아직 없을 때는 seedPrompts가 최신 기본값을 바로 seed하므로 이 마이그레이션 불필요
	}
	tmpl := agent.BuiltinPromptSeeds()["planner"]
	if tmpl == "" {
		return
	}
	// 완전 새 DB는 seedPrompts가 이미 최신 기본값을 seed함 → 현재 버전이 코드 기본값과 같으므로 중복 버전을 추가할 필요 없음.
	if cur, err := s.m.pg.CurrentPrompt(a.ID); err == nil && cur == tmpl {
		return
	}
	if _, err := s.m.pg.ResetPromptToDefault(a.ID, tmpl); err != nil {
		log.Printf("[prompts] planner 프롬프트를 새 기본값으로 재갱신 실패: %v", err)
		return
	}
	log.Printf("[prompts] planner 프롬프트에 새 기본값 버전 추가 완료(간소화 재구성+절제 강등 중복 제거+깊이 우선+부정 재검토 상한, 일회성)")
}

// reseedWorkerPrompt는 worker 프롬프트를 [현재 코드 기본값]으로 갱신한다 —— 기본 본문 record_fact 단락에서 '부정류 결론은
// 관찰로 쓰고 탐색적 읽기법' 문장 전체를 삭제, confidence(observed/inferred)와 '이 의도의 수단을 소진했는지'를 디커플링했다(이들은 planner를 오도하기 쉬움),
// 동시에 facts 배열의 분조를 '서로 완전히 독립적이라 병합 불가'한 극소수 예외로 좁혔다. flag를 v3로 bump해 기존 구 DB를 다시 한 번 갱신한다.
// SeedPromptIfEmpty는 첫 삽입only라, 구 DB의 기존 버전은 받지 못하므로 버전 관리로 [새 버전 하나 추가]하고 전환한다, 구 버전은 이력에 남아 되찾을 수 있다.
// settings flag 가드 → 한 번만. 완전 새 DB는 처리 불필요. reseedGoalsPrompt와 완전히 동형.
func (s *Server) reseedWorkerPrompt() {
	const flag = "worker_prompt_compact_v4"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 성공 여부와 무관하게 한 번만 시도
	a, err := s.m.pg.GetAgentByKey("worker")
	if err != nil || a == nil {
		return // 완전 새 DB에서 agent 행이 아직 없을 때는 seedPrompts가 최신 기본값을 바로 seed하므로 이 마이그레이션 불필요
	}
	tmpl := agent.BuiltinPromptSeeds()["worker"]
	if tmpl == "" {
		return
	}
	// 완전 새 DB는 seedPrompts가 이미 최신 기본값을 seed함 → 현재 버전이 코드 기본값과 같으므로 중복 버전을 추가할 필요 없음.
	if cur, err := s.m.pg.CurrentPrompt(a.ID); err == nil && cur == tmpl {
		return
	}
	if _, err := s.m.pg.ResetPromptToDefault(a.ID, tmpl); err != nil {
		log.Printf("[prompts] worker 프롬프트를 새 기본값으로 재갱신 실패: %v", err)
		return
	}
	log.Printf("[prompts] worker 프롬프트에 새 기본값 버전 추가 완료(컨텍스트 조회 단락을 list_assets/list_findings로 수렴, list_facts/node_detail/asset_neighbors 제거, 일회성)")
}

// reporterToolCallMessage는 보고서를 쓰기 전에 무조건 get_finding_traffic을 한 번 읽도록 요구해야 한다.
// 이 도구는 읽기 전용이고 '캡처 스위치에 의존하지 않아', 자동 바인딩을 켜든 끄든 수동 바인딩한 증거를 읽을 수 있다. 여기를
// '자동 바인딩을 켜야 읽음'으로 쓰면, 기본 꺼짐 설정에서 reporter가 evidence_version을 전달하지 않고,
// SetFindingReportVersionByNodeID가 legacy 의미로 -1을 써서, 취약점 상세와 Markdown 내보내기가
// 그때부터 '증거가 변경됨, 보고서 갱신 대기'에 상주하는데, UI에는 그것을 지울 입구가 전혀 없다.
const reporterToolCallMessage = "위에서 방금 취약점 하나가 report_finding으로 등록됐다. 반환 JSON의 finding_id(독립 취약점 기록 ID)와 finding_node_id(탐색 노드 ID)를 읽어라, " +
	"먼저 get_finding_traffic(finding_id)으로 현재 증거 목록과 그 version을 읽는다(빈 목록은 정상이며, 평소대로 보고서를 쓴다); " +
	"실행 지침이 자동 바인딩을 켰다면, 읽기 전에 이번 취약점의 트래픽을 먼저 확인해 연관시켜라. 노드 상세는 finding_node_id를 쓴다. " +
	"마지막으로 update_finding_report(finding_id=finding_node_id, report, evidence_version=실제로 읽은 버전)을 호출해 저장하라, " +
	"evidence_version은 반드시 전달해야 하며, 아니면 보고서가 영구적으로 갱신 대기로 표시된다. 두 종류의 번호를 혼용하지 마라."

// 구 버전 트리거 메시지(0.3.8 및 이전). 여전히 이것과 바이트 단위로 동일한 기록만 마이그레이션으로 덮어쓰고, 사용자가 수정한 것은 원래대로 둔다.
const reporterToolCallMessageV1 = "上面刚有一个漏洞被 report_finding 登记。请从触发上下文里取出 finding_id" +
	"（工具返回 \"finding recorded: <id>\" 里的数字）与任务 id，按你的职责撰写该漏洞的详细报告，" +
	"最后调用 update_finding_report(finding_id, report) 保存。"

// upgradeReporterTriggerMessage는 구 DB에서 여전히 기본 문구인 reporter 트리거 메시지를 새 버전으로 갱신한다.
// seedReporterAgent는 reporter_agent_seed_v1 가드를 받고 새 agent 생성 시에만 트리거를 쓰므로,
// 업그레이드되어 온 DB는 새 문구를 받지 못한다 —— 도구 schema는 seedFindingTrafficTools가
// evidence_version을 보강했지만, reporter에게 그것을 쓰라고 알려주는 것이 아무것도 없다. 일회성, 그리고 수정되지 않은 문구만 덮어씀.
func (s *Server) upgradeReporterTriggerMessage() {
	const flag = "reporter_trigger_evidence_version_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 한 번만 시도
	triggers, err := s.m.pg.ListTriggersFor("reporter")
	if err != nil {
		log.Printf("[reporter] 트리거 읽기 실패: %v", err)
		return
	}
	for _, t := range triggers {
		if !t.OnToolCall || t.ToolCallMessage != reporterToolCallMessageV1 {
			continue // 사용자가 수정했거나 finding 트리거가 아님, 건드리지 않음.
		}
		t.ToolCallMessage = reporterToolCallMessage
		if err := s.m.pg.UpdateTrigger(t); err != nil {
			log.Printf("[reporter] 트리거 메시지 업그레이드 실패: %v", err)
			return
		}
		log.Printf("[reporter] 트리거 메시지를 evidence_version 읽기·되돌려 전달로 업그레이드 완료")
	}
}

// seedReporterAgent는 '보고서 작성' 커스텀 agent 하나를 사전 배치한다(builtin=false, UI에서 편집/삭제 가능):
// update_finding_report + 작업 조회 도구를 바인딩하고, 'report_finding 호출 즉시 트리거'되는
// 트리거를 하나 걸어 —— 취약점 하나를 등록할 때마다 그것을 깨워 상세 보고서를 쓰게 한다. 일회성(settings flag 가드): 사용자가 삭제하면 다시 만들지 않음.
// 의존: orchestration 도구가 이 함수 위쪽에서 이미 SeedTool로 입고돼 있어, 바인딩이 된다.
func (s *Server) seedReporterAgent() {
	const flag = "reporter_agent_seed_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	defer func() { _ = s.m.pg.SetSetting(flag, "true") }() // 성공 여부와 무관하게 한 번만 시도

	if exist, _ := s.m.pg.GetAgentByKey("reporter"); exist != nil {
		return // key가 이미 점유됨(사용자가 수동 생성)——덮어쓰지 않음
	}
	a, err := s.m.pg.CreateAgent("reporter", "보고서 작성",
		"취약점 상세 보고서 작성: 취약점 발견 시 자동 트리거, 증거와 실행 과정을 조회한 뒤 Markdown 보고서를 쓰고 되돌려 기록.")
	if err != nil {
		log.Printf("[reporter] agent 생성 실패: %v", err)
		return
	}
	if err := s.m.pg.SeedPromptIfEmpty(a.ID, agent.ReporterDefaultPrompt); err != nil {
		log.Printf("[reporter] seed prompt 실패: %v", err)
	}
	// 트리거 실행 정책: parallel + none —— 취약점 하나당 보고서 하나, 여러 finding을 동시에 각자 씀.
	// merge는 반드시 none: 아니면(기본 all) 한 무더기 finding이 한 번의 실행으로 병합돼 병렬이 무의미해짐.
	// maxParallel=5: 동시에 최대 5개 보고서 세션, 순간적으로 LLM 호출이 너무 많아지는 것을 방지.
	if err := s.m.pg.SetAgentTriggerBehavior("reporter", "parallel", "none", 5); err != nil {
		log.Printf("[reporter] 트리거 실행 정책 설정 실패: %v", err)
	}
	// 필요한 도구를 바인딩: 보고서 쓰기 + 증거/실행 과정/상황 읽기.
	if err := s.m.pg.AddAgentToToolBinding("reporter", []string{
		"update_finding_report", "get_task_node_detail", "list_task_findings",
		"get_task_worker_trace", "list_task_worker_traces", "search_task_worker_traces",
		"get_task_graph",
	}); err != nil {
		log.Printf("[reporter] 도구 바인딩 실패: %v", err)
	}
	// 트리거: report_finding 호출 즉시 트리거(도구가 "finding recorded: <id>"를 반환하며 finding_id를 실어 보내고,
	// 작업 id도 트리거 메시지 안에 있음).
	if _, err := s.m.pg.CreateTrigger(&db.AgentTrigger{
		AgentKey:        "reporter",
		Enabled:         true,
		OnToolCall:      true,
		ToolNames:       []string{"report_finding"},
		ToolCallMessage: reporterToolCallMessage,
	}); err != nil {
		log.Printf("[reporter] 트리거 생성 실패: %v", err)
	}
	log.Printf("[reporter] '보고서 작성' agent + finding 트리거 사전 배치 완료")
}

// seedAutoReportFindingBinding adds "auto" to report_finding's binding ONCE so
// conversation-context agents can call it without requiring an intent_id.
func (s *Server) seedAutoReportFindingBinding() {
	const flag = "auto_report_finding_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.AddAgentToToolBinding("auto", []string{"report_finding"}); err != nil {
		log.Printf("[auto] report_finding 기본 바인딩 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedPlannerDefaultBindings adds "planner" to report_finding's binding ONCE
// (guarded by a settings flag), so existing DBs — whose report_finding row was
// seeded as worker-only — also let the planner record findings. Fresh DBs already
// get it via PlannerTools(); this only backfills without overriding a user unbind.
func (s *Server) seedPlannerDefaultBindings() {
	const flag = "planner_report_finding_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.AddAgentToToolBinding("planner", []string{"report_finding"}); err != nil {
		log.Printf("[planner] report_finding 기본 바인딩 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedPlannerListAssetsBinding adds "planner" to list_assets's binding ONCE
// (guarded by a settings flag), so existing DBs — whose list_assets row was seeded
// as auto/pentest-only — also let the planner query the asset store by DSL. Fresh
// DBs already get it via PlannerTools(); this only backfills without overriding a
// user unbind.
func (s *Server) seedPlannerListAssetsBinding() {
	const flag = "planner_list_assets_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.AddAgentToToolBinding("planner", []string{"list_assets"}); err != nil {
		log.Printf("[planner] list_assets 기본 바인딩 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedCompanyScopeRebind changes add_company_scope's default binding ONCE on
// existing DBs (guarded by a settings flag): the tool moves off worker and onto
// planner — defining a company's asset scope is a planning/main/auto concern, not
// something a worker does mid-exploration. Fresh DBs already get planner via
// PlannerTools() and lack worker via WorkerTools(); this only backfills old rows.
// One-shot + flag-guarded so a user who later re-binds worker isn't overridden.
func (s *Server) seedCompanyScopeRebind() {
	const flag = "company_scope_rebind_v1" // worker→planner 기본 바인딩 전환
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.AddAgentToToolBinding("planner", []string{"add_company_scope"}); err != nil {
		log.Printf("[planner] add_company_scope 기본 바인딩 실패: %v", err)
		return
	}
	if err := s.m.pg.RemoveAgentFromTool("worker", "add_company_scope"); err != nil {
		log.Printf("[worker] add_company_scope 바인딩 해제 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedWorkerReadToolsUnbind strips the read-context tools off worker's default
// binding ONCE on existing DBs (guarded by a settings flag): a worker executes one
// intent and writes back — reading facts/companies and listing all workers' traces is
// a planning/main concern, not the executor's. Fresh DBs already lack these via
// WorkerTools(); this only backfills old rows without overriding a user who
// deliberately re-binds worker. Each RemoveAgentFromTool is per-tool +
// membership-guarded, so planner/mainagent bindings of the same tool are untouched.
//
// NOTE: search_all_worker_traces / get_worker_trace / node_detail are intentionally NOT
// unbound — worker owns them for cross-work look-back + node drill-down (see WorkerTools).
// They used to be in this list back when worker lacked them; seedWorkerReadbackRebind
// repairs DBs whose old run stripped them.
func (s *Server) seedWorkerReadToolsUnbind() {
	const flag = "worker_readtools_unbind_v1"
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	for _, k := range []string{
		"list_facts", "list_companies", "list_worker_traces",
	} {
		if err := s.m.pg.RemoveAgentFromTool("worker", k); err != nil {
			log.Printf("[worker] %s를 worker에서 바인딩 해제 실패: %v", k, err)
			return // 오류면 flag를 기록하지 않음, 다음 시작 때 재시도
		}
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedWorkerReadbackRebind re-binds the cross-work look-back / drill-down tools onto
// worker ONCE (guarded by a settings flag): an earlier seedWorkerReadToolsUnbind wrongly
// stripped search_all_worker_traces / get_worker_trace / node_detail from worker after
// they had been added to WorkerTools(), so any DB that ran that migration lost them.
// Fresh DBs already have them via WorkerTools() and this is a harmless no-op there.
// One-shot + flag-guarded so a user who later deliberately unbinds them isn't overridden.
func (s *Server) seedWorkerReadbackRebind() {
	const flag = "worker_readback_rebind_v2" // v2: node_detail 추가
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	if err := s.m.pg.AddAgentToToolBinding("worker", []string{
		"search_all_worker_traces", "get_worker_trace", "node_detail",
	}); err != nil {
		log.Printf("[worker] 회람/상세 도구 보강 바인딩 실패: %v", err)
		return // 오류면 flag를 기록하지 않음, 다음 시작 때 재시도
	}
	_ = s.m.pg.SetSetting(flag, "true")
}

// seedAutoDefaultBindings adds "auto" to the task-op + platform tools' bindings
// ONCE (guarded by a settings flag), so existing DBs whose tool rows were seeded
// before Auto existed still give Auto its default toolset — without re-adding it
// after a user deliberately unbinds.
func (s *Server) seedAutoDefaultBindings() {
	const flag = "auto_default_bindings_v3" // v3: 구 자산 도구명 교체, insert_assets/add_company_scope 추가
	if v, _, _ := s.m.pg.GetSetting(flag); v == "true" {
		return
	}
	keys := make([]string, 0, len(platformToolKeys)+12)
	for _, t := range s.orchestrationTools() {
		keys = append(keys, t.Name())
	}
	keys = append(keys, platformToolKeys...)
	// 자산 도구: Auto가 플랫폼을 조작하며 자산을 자주 조회/등록하고, 회사 범위를 관리.
	keys = append(keys, "insert_assets", "add_company_scope", "list_assets")
	if err := s.m.pg.AddAgentToToolBinding("auto", keys); err != nil {
		log.Printf("[auto] 기본 바인딩 실패: %v", err)
		return
	}
	_ = s.m.pg.SetSetting(flag, "true")
}
