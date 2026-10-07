package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/traffic"
	actool "github.com/Autumn-27/norma/tool"
)

func workflowCall(t *testing.T, ctx context.Context, tool actool.CoreTool, input any, wantError bool) string {
	t.Helper()
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	r, err := tool.Call(ctx, raw, nil)
	if err != nil || r.IsError != wantError {
		t.Fatalf("%s: error=%v result=%s", tool.Name(), err, r.Flatten())
	}
	return r.Flatten()
}

func workflowTool(t *testing.T, tools []actool.CoreTool, name string) actool.CoreTool {
	t.Helper()
	for _, tool := range tools {
		if tool.Name() == name {
			return tool
		}
	}
	t.Fatalf("missing tool: %s", name)
	return nil
}

func TestFindingWorkflowAutoHintToPlannerAndSetting(t *testing.T) {
	s, initial, request := trafficEvidenceServer(t)
	pg := s.m.pg
	ctx := context.Background()
	if _, err := pg.Exec(`DELETE FROM settings WHERE key=$1`, settingAgentTrafficBinding); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pg.SetBool(settingAgentTrafficBinding, false) })
	if s.settingsPayload()[settingAgentTrafficBinding] != false {
		t.Fatal("missing setting must default off")
	}
	seedServerEvidenceFlow(t, s, "handoff-proof", []byte("verified local proof"))
	seedServerEvidenceFlow(t, s, "handoff-baseline", []byte("local baseline"))
	seedServerEvidenceFlow(t, s, "handoff-verification", []byte{0, 1, 255})
	// Manual binding is independent of the Agent switch.
	if r := request("POST", fmt.Sprintf("/api/exploration/findings/%d/traffic", initial.FindingID), `{"traffic_refs":[{"traffic_id":"handoff-proof"}]}`); r.Code != 200 {
		t.Fatal(r.Code, r.Body)
	}
	bind := s.toolBindFindingTraffic()
	input := map[string]any{"finding_id": initial.FindingID, "traffic_refs": []db.TrafficRef{{TrafficID: "handoff-baseline"}}}
	workflowCall(t, ctx, bind, input, true)
	if r := request("PUT", "/api/settings", `{"agent_traffic_binding":true}`); r.Code != 200 {
		t.Fatal(r.Code, r.Body)
	}
	if !pg.GetBool(settingAgentTrafficBinding, false) {
		t.Fatal("switch was not persisted")
	}
	f, _ := pg.GetFinding(initial.FindingID)
	task := s.m.ResolveTask(fmt.Sprint(*f.TaskID))
	refs := []db.TrafficRef{{TrafficID: "handoff-baseline", Role: "baseline", Note: "normal response"}, {TrafficID: "handoff-proof", Role: "proof", Note: "proves the finding"}, {TrafficID: "handoff-verification", Role: "verification", Note: "binary verification"}}
	// Real cross-task Auto hint handler -> persisted graph -> Planner tool.
	hintResult := workflowCall(t, ctx, s.toolAddHint(), map[string]any{"task_id": task.ID, "hints": []any{map[string]any{"text": "Report the confirmed local finding with its verified evidence", "traffic_refs": refs}}}, false)
	var hints struct {
		IDs []int64 `json:"ids"`
	}
	if err := json.Unmarshal([]byte(hintResult), &hints); err != nil || len(hints.IDs) != 1 || hints.IDs[0] == 0 {
		t.Fatal(hintResult, err)
	}
	graph := workflowCall(t, ctx, s.toolGetTaskGraph(), map[string]any{"task_id": task.ID}, false)
	for _, ref := range refs {
		if !strings.Contains(graph, ref.TrafficID) {
			t.Fatal("handoff lost reference", ref.TrafficID)
		}
	}
	for range 6 {
		if _, err := task.Store.AddNode(db.KindFact, map[string]any{"summary": "ID separation"}, 0, "confirmed", "test", nil); err != nil {
			t.Fatal(err)
		}
	}
	ts := agent.NewToolSet(task.Store, "planner")
	ts.SetTaskID(*f.TaskID)
	ts.SetFindingRecorder(s.evidenceStore())
	notices := 0
	ts.SetNotifyFinding(func(int64, string) { notices++ })
	tools, def, cleanup := agent.AugmentTools(ctx, "planner", ts.PlannerTools())
	defer cleanup()
	if !strings.Contains(def.FindingGuidance, "evidence_hint_id") || !strings.Contains(def.FindingGuidance, "Worker 를 취소") {
		t.Fatal("Planner missed runtime guidance")
	}
	report := workflowTool(t, tools, "report_finding")
	result := workflowCall(t, ctx, report, map[string]any{"vulnclass": "TEST", "severity": "low", "summary": "Planner handoff fixture", "evidence_hint_id": hints.IDs[0]}, false)
	var recorded db.RecordedFinding
	if err := json.Unmarshal([]byte(strings.SplitN(result, "\n", 2)[1]), &recorded); err != nil {
		t.Fatal(err)
	}
	if recorded.FindingID == recorded.NodeID || len(recorded.Traffic.Bindings) != 3 || notices != 1 {
		t.Fatal("bad finding/traffic result", result)
	}
	for i, b := range recorded.Traffic.Bindings {
		if b.Snapshot.SourceTrafficID != refs[i].TrafficID || b.Role != refs[i].Role || b.Note != refs[i].Note {
			t.Fatalf("handoff mismatch: %+v", b)
		}
	}
	list := workflowCall(t, ctx, s.toolListTaskFindings(), map[string]any{"task_id": task.ID}, false)
	var nodes []struct {
		ID            int64 `json:"id"`
		FindingID     int64 `json:"finding_id"`
		FindingNodeID int64 `json:"finding_node_id"`
		Count         int   `json:"traffic_count"`
	}
	if err := json.Unmarshal([]byte(list), &nodes); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, n := range nodes {
		if n.ID == recorded.NodeID {
			found = n.FindingID == recorded.FindingID && n.FindingNodeID == recorded.NodeID && n.Count == 3
		}
	}
	if !found {
		t.Fatal("canonical IDs/count missing", list)
	}
	workflowCall(t, ctx, s.toolGetFindingTraffic(), map[string]any{"finding_id": recorded.FindingID}, false)
	wrong := workflowCall(t, ctx, s.toolGetFindingTraffic(), map[string]any{"finding_id": int64(900000000000000000)}, true)
	if !strings.Contains(wrong, "독립 취약점 기록 ID") {
		t.Fatal("ambiguous ID error", wrong)
	}
	input["finding_id"] = recorded.FindingID
	// Duplicate append preserves metadata and version.
	workflowCall(t, ctx, bind, input, false)
	after, err := pg.GetFindingTraffic(ctx, recorded.FindingID)
	if err != nil || len(after.Bindings) != 3 || after.Version != recorded.Traffic.Version || after.Bindings[0].Note != refs[0].Note {
		t.Fatal("duplicate changed evidence", after, err)
	}
	// A malformed handoff must not create a partial finding or notify Planner.
	badHint := workflowCall(t, ctx, s.toolAddHint(), map[string]any{"task_id": task.ID, "text": "missing packet", "traffic_refs": []db.TrafficRef{{TrafficID: "missing"}}}, false)
	var badID int64
	fmt.Sscanf(badHint, "hint added: %d", &badID)
	workflowCall(t, ctx, report, map[string]any{"vulnclass": "TEST", "severity": "low", "summary": "must fail", "evidence_hint_id": badID}, true)
	if notices != 1 {
		t.Fatal("failed handoff notified before commit")
	}
	child, err := s.m.CreateTaskWithOptions("inherited handoff", "fixture", db.TaskCreateOptions{SourceTaskIDs: []int64{*f.TaskID}})
	if err != nil {
		t.Fatal(err)
	}
	childID, _ := strconv.ParseInt(child.ID, 10, 64)
	childRow, _ := pg.GetTask(childID)
	t.Cleanup(func() { pg.DeleteTask(childRow.ID) })
	childCtx := agent.WithRunInfo(ctx, agent.RunInfo{TaskID: childRow.ID})
	workflowCall(t, childCtx, bind, input, true)
	workflowCall(t, childCtx, s.toolGetFindingTraffic(), map[string]any{"finding_id": recorded.FindingID}, false)
	// Existing, already assembled tools must honor a later switch-off.
	if r := request("PUT", "/api/settings", `{"agent_traffic_binding":false}`); r.Code != 200 {
		t.Fatal(r.Code, r.Body)
	}
	workflowCall(t, ctx, bind, input, true)
	// Turning binding off must still save a confirmed finding from an already
	// assembled tool, even if the model sends the old optional evidence fields.
	offResult := workflowCall(t, ctx, report, map[string]any{"vulnclass": "TEST", "severity": "low", "summary": "off", "evidence_hint_id": hints.IDs[0]}, false)
	var offRecord struct {
		db.RecordedFinding
		EvidenceStatus string `json:"evidence_status"`
		EvidenceNote   string `json:"evidence_note"`
	}
	if err := json.Unmarshal([]byte(strings.SplitN(offResult, "\n", 2)[1]), &offRecord); err != nil {
		t.Fatal(err)
	}
	if offRecord.FindingID <= 0 || len(offRecord.Traffic.Bindings) != 0 || offRecord.EvidenceStatus != "not_bound" || !strings.Contains(offRecord.EvidenceNote, "비활성화") {
		t.Fatal("disabled binding discarded finding or bound evidence", offResult)
	}
	workflowCall(t, ctx, report, map[string]any{"vulnclass": "TCP", "severity": "low", "summary": "no packet needed"}, false)
	if notices != 3 {
		t.Fatal("optional no-packet report failed")
	}
	workflowCall(t, ctx, s.toolGetFindingTraffic(), map[string]any{"finding_id": recorded.FindingID}, false)
	tools, def, closeTools := agent.AugmentTools(ctx, "planner", ts.PlannerTools())
	defer closeTools()
	if def.FindingGuidance != "" {
		t.Fatal("disabled feature still injects guidance")
	}
	for _, tool := range tools {
		if tool.Name() == "bind_finding_traffic" {
			t.Fatal("disabled binding tool exposed")
		}
	}
}

func TestFindingWorkflowMigrationPreservesUserConfiguration(t *testing.T) {
	s, _, _ := trafficEvidenceServer(t)
	pg := s.m.pg
	for _, key := range []string{"report_finding", "add_hint", "add_task_hint", "traffic_search", "traffic_get", "traffic_blob", "bind_finding_traffic", "get_finding_traffic"} {
		old, err := pg.GetTool(key)
		if err != nil || old == nil {
			t.Fatal("missing tool", key, err)
		}
		t.Cleanup(func() {
			bindings, _ := json.Marshal(old.Agents)
			pg.UpdateTool(key, old.Description, old.Schema, bindings, old.Enabled)
		})
	}
	custom := json.RawMessage(`{"type":"object","properties":{"text":{"type":"string","description":"USER TEXT"},"hints":{"type":"array","items":{"type":"object","properties":{"text":{"type":"string","description":"USER ITEM"}}}}},"required":["text"]}`)
	for _, key := range []string{"report_finding", "add_hint", "add_task_hint"} {
		if err := pg.UpdateTool(key, "USER DESCRIPTION", custom, json.RawMessage(`["custom-agent"]`), false); err != nil {
			t.Fatal(err)
		}
	}
	for _, key := range []string{"traffic_search", "traffic_get"} {
		row, _ := pg.GetTool(key)
		bindings := json.RawMessage(`["custom-agent"]`)
		if key == "traffic_search" {
			bindings = json.RawMessage(`["worker"]`)
		}
		if err := pg.UpdateTool(key, row.Description, row.Schema, bindings, false); err != nil {
			t.Fatal(err)
		}
	}
	if err := pg.SetSetting("finding_workflow_tools_v2_reporter", "false"); err != nil {
		t.Fatal(err)
	}
	s.seedFindingWorkflowTools()
	for _, key := range []string{"report_finding", "add_hint", "add_task_hint"} {
		row, _ := pg.GetTool(key)
		if row.Enabled || row.Description != "USER DESCRIPTION" || len(row.Agents) != 1 || row.Agents[0] != "custom-agent" {
			t.Fatal("changed user tool configuration", key)
		}
		if !strings.Contains(string(row.Schema), "USER TEXT") || !strings.Contains(string(row.Schema), "USER ITEM") {
			t.Fatal("changed user schema", key)
		}
		field := "traffic_refs"
		if key == "report_finding" {
			field = "evidence_hint_id"
		}
		if !strings.Contains(string(row.Schema), field) {
			t.Fatal("missing optional field", key)
		}
	}
	search, _ := pg.GetTool("traffic_search")
	get, _ := pg.GetTool("traffic_get")
	if !strings.Contains(search.Description, "호스트 단독, 호스트:포트 또는 완전한 URL을 지원") {
		t.Fatal("traffic_search description migration missing host/port guidance")
	}
	if search.Enabled || !contains(search.Agents, "reporter") || get.Enabled || len(get.Agents) != 1 || get.Agents[0] != "custom-agent" {
		t.Fatal("default/custom reader binding migration incorrect")
	}
	if err := pg.RemoveAgentFromTool("reporter", "traffic_search"); err != nil {
		t.Fatal(err)
	}
	s.seedFindingWorkflowTools()
	search, _ = pg.GetTool("traffic_search")
	if contains(search.Agents, "reporter") {
		t.Fatal("one-time migration undid later unbinding")
	}
}

func TestFindingWorkflowHostSearchDescriptionMigration(t *testing.T) {
	s, _, _ := trafficEvidenceServer(t)
	pg := s.m.pg
	const hostSearchDescriptionFlag = "finding_workflow_tools_v3_host_search_description"
	const reporterFlag = "finding_workflow_tools_v2_reporter"
	// bf425ef의 직접 부모 42e1deb에서 확보한 독립 fixture. SQL 비교 계약이므로 원문을 보존합니다.
	const historicalDescription = "查询记录代理已抓取的目标流量（必须指定 host，可再按 URL 子串或正文关键词过滤）。body_contains 会在已抓取的请求/响应头与正文中做全文搜索，支持任意子串和中文（至少 3 个字符），可用来找响应里的密码、密钥、报错、内网地址等。仅返回极轻量索引(id/method/url/status/resp_len)，不含任何响应内容。默认只返回 3 条、每页最多 10 条；结果多时用 page 翻页（page=0 起）；要看某条的请求/响应原文用 traffic_get(id)。回看已访问资源、找端点先用它，避免重复 curl 同一 URL。"
	const customDescription = "USER CUSTOM traffic_search DESCRIPTION\n사용자가 편집한 호스트·응답 본문 안내"
	customSchema := json.RawMessage(`{"type":"object","properties":{"host":{"type":"string","description":"USER HOST DESCRIPTION"}},"required":["host"]}`)
	customBindings := json.RawMessage(`["custom-agent"]`)

	// 두 사례는 순차 실행하며 각 사례가 변경한 행과 플래그를 다음 사례 전에 복원합니다.
	for _, tc := range []struct {
		name        string
		description string
		want        string
	}{
		{"역사 기본 설명 교체", historicalDescription, traffic.TrafficSearchDescription},
		{"사용자 편집 설명 보존", customDescription, customDescription},
	} {
		t.Run(tc.name, func(t *testing.T) {
			old, err := pg.GetTool("traffic_search")
			if err != nil || old == nil || !old.System {
				t.Fatalf("system traffic_search 행이 없습니다: %v", err)
			}
			t.Cleanup(func() {
				bindings, err := json.Marshal(old.Agents)
				if err != nil {
					t.Errorf("기존 바인딩 직렬화: %v", err)
					return
				}
				if err := pg.UpdateTool(old.Key, old.Description, old.Schema, bindings, old.Enabled); err != nil {
					t.Errorf("기존 traffic_search 행 복원: %v", err)
				}
			})
			for _, flag := range []string{hostSearchDescriptionFlag, reporterFlag} {
				value, exists, err := pg.GetSetting(flag)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if exists {
						if err := pg.SetSetting(flag, value); err != nil {
							t.Errorf("기존 플래그 %s 복원: %v", flag, err)
						}
					} else if _, err := pg.Exec(`DELETE FROM settings WHERE key=$1`, flag); err != nil {
						t.Errorf("추가한 플래그 %s 삭제: %v", flag, err)
					}
				})
			}
			if err := pg.UpdateTool("traffic_search", tc.description, customSchema, customBindings, false); err != nil {
				t.Fatal(err)
			}
			prepared, err := pg.GetTool("traffic_search")
			if err != nil || prepared == nil || !prepared.System || prepared.Description != tc.description || prepared.Enabled || len(prepared.Agents) != 1 || prepared.Agents[0] != "custom-agent" {
				t.Fatalf("traffic_search fixture 준비 실패: row=%+v error=%v", prepared, err)
			}
			// v3 설명 교체만 검사하며, 별도의 v2 스키마/바인딩 마이그레이션은 격리합니다.
			if err := pg.SetSetting(reporterFlag, "true"); err != nil {
				t.Fatal(err)
			}
			if err := pg.SetSetting(hostSearchDescriptionFlag, "false"); err != nil {
				t.Fatal(err)
			}
			if value, exists, err := pg.GetSetting(hostSearchDescriptionFlag); err != nil || !exists || value != "false" {
				t.Fatalf("설명 마이그레이션 미완료 조건 준비 실패: value=%q exists=%v error=%v", value, exists, err)
			}

			s.seedFindingWorkflowTools()

			after, err := pg.GetTool("traffic_search")
			if err != nil || after == nil {
				t.Fatal("seed 이후 traffic_search 조회 실패", err)
			}
			if after.Description != tc.want {
				t.Fatalf("설명 불일치: got=%q want=%q", after.Description, tc.want)
			}
			if !after.System || string(after.Schema) != string(prepared.Schema) || after.Enabled != prepared.Enabled || len(after.Agents) != 1 || after.Agents[0] != "custom-agent" {
				t.Fatalf("설명 마이그레이션이 사용자 스키마/바인딩/활성 설정을 변경했습니다: %+v", after)
			}
			// 사용자 설명이 매칭되지 않아도 SQL 자체가 성공하면 현재 구현은 완료로 기록합니다.
			if value, exists, err := pg.GetSetting(hostSearchDescriptionFlag); err != nil || !exists || value != "true" {
				t.Fatalf("설명 마이그레이션 완료 플래그 불일치: value=%q exists=%v error=%v", value, exists, err)
			}
			if value, exists, err := pg.GetSetting(reporterFlag); err != nil || !exists || value != "true" {
				t.Fatalf("격리한 reporter 플래그가 변경됐습니다: value=%q exists=%v error=%v", value, exists, err)
			}
		})
	}
}

func TestFindingWorkflowReporterBindsBeforeWritingReport(t *testing.T) {
	s, f, request := trafficEvidenceServer(t)
	pg := s.m.pg
	ctx := context.Background()
	if err := pg.SetBool(settingAgentTrafficBinding, true); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pg.SetBool(settingAgentTrafficBinding, false); s.m.SetTrafficEnabled(false) })
	if err := s.m.SetTrafficEnabled(true); err != nil {
		t.Fatal(err)
	}
	seedServerEvidenceFlow(t, s, "reporter-proof", []byte("local proof payload"))
	seedServerEvidenceFlow(t, s, "reporter-baseline", []byte("local normal response"))
	tools, def, cleanup := agent.AugmentTools(ctx, "reporter", nil)
	defer cleanup()
	if !strings.Contains(def.FindingGuidance, "보고 전 트래픽 자동 연관") || !strings.Contains(def.FindingGuidance, "바인딩에 성공한 뒤") {
		t.Fatal("reporter did not receive binding workflow")
	}
	for _, name := range []string{"traffic_search", "traffic_get", "get_task_worker_trace", "get_task_node_detail", "bind_finding_traffic", "get_finding_traffic", "update_finding_report"} {
		workflowTool(t, tools, name)
	}
	finding, _ := pg.GetFinding(f.FindingID)
	detail := workflowCall(t, ctx, workflowTool(t, tools, "get_task_node_detail"), map[string]any{"task_id": fmt.Sprint(*finding.TaskID), "id": f.NodeID}, false)
	if !strings.Contains(detail, `"finding_id"`) || !strings.Contains(detail, `"finding_node_id"`) {
		t.Fatal("reporter lacks explicit ID mapping")
	}
	workflowCall(t, ctx, workflowTool(t, tools, "traffic_search"), map[string]any{"host": "evidence.local"}, false)
	for _, id := range []string{"reporter-baseline", "reporter-proof"} {
		packet := workflowCall(t, ctx, workflowTool(t, tools, "traffic_get"), map[string]any{"id": id}, false)
		if !strings.Contains(packet, "local") {
			t.Fatal("reporter cannot inspect packet", packet)
		}
	}
	refs := []db.TrafficRef{{TrafficID: "reporter-baseline", Role: "baseline"}, {TrafficID: "reporter-proof", Role: "proof", Note: "confirmed from recorded response"}}
	workflowCall(t, ctx, workflowTool(t, tools, "bind_finding_traffic"), map[string]any{"finding_id": f.FindingID, "traffic_refs": refs}, false)
	list := workflowCall(t, ctx, workflowTool(t, tools, "get_finding_traffic"), map[string]any{"finding_id": f.FindingID}, false)
	var summary struct {
		Version int64 `json:"version"`
	}
	if err := json.Unmarshal([]byte(list), &summary); err != nil || summary.Version != 1 {
		t.Fatal(list, err)
	}
	workflowCall(t, ctx, workflowTool(t, tools, "update_finding_report"), map[string]any{"finding_id": f.NodeID, "evidence_version": summary.Version, "report": "## Local report\n\nVerified baseline and proof using saved evidence."}, false)
	updated, err := pg.GetFinding(f.FindingID)
	if err != nil || updated.ReportEvidenceVersion != summary.Version || updated.EvidenceVersion != summary.Version {
		t.Fatal("report did not cover post-binding version", updated, err)
	}
	// Disabling automatic binding still lets Reporter read manually bound snapshots.
	if r := request("PUT", "/api/settings", `{"agent_traffic_binding":false}`); r.Code != 200 {
		t.Fatal(r.Code, r.Body)
	}
	off, offDef, closeOff := agent.AugmentTools(ctx, "reporter", nil)
	defer closeOff()
	if offDef.FindingGuidance != "" {
		t.Fatal("off reporter still receives auto-binding guidance")
	}
	for _, tool := range off {
		if tool.Name() == "traffic_search" || tool.Name() == "traffic_get" || tool.Name() == "traffic_blob" || tool.Name() == "bind_finding_traffic" {
			t.Fatal("auto-binding tool exposed while off", tool.Name())
		}
	}
	workflowCall(t, ctx, workflowTool(t, off, "get_finding_traffic"), map[string]any{"finding_id": f.FindingID}, false)
	workflowTool(t, off, "update_finding_report")
}
