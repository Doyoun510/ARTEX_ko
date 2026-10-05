package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/db"
)

// 개요의 '목표 관리' 수동 CRUD 인터페이스. agent 측 set_goals 도구와 같은 goal 노드 묶음에 쓰지만,
// 입구는 사람이 UI에서 직접 추가·삭제·수정하는 것이다; 추가/수정 후에는 '작업 되살리기' 로직을 재사용한다(admitTask resume:
// 종료 상태→running·일시정지 해제·필요 시 큐 대기), 삭제는 되살리지 않는다(제품 결정에 따라). 각 변경 handler는
// beginTaskOperation/decInflight를 거쳐, 작업 삭제와의 경쟁 상태를 피한다(의도 CRUD와 동일).

// listGoals는 본 작업의 모든 목표를 반환한다(text/vulnclass/state로 분해), 목표 관리 카드 렌더링용.
func (s *Server) listGoals(w http.ResponseWriter, r *http.Request) {
	t, ok := s.m.Task(r.PathValue("id"))
	if !ok {
		writeErr(w, 404, "task not found")
		return
	}
	goals, err := t.Store.ListByKind(db.KindGoal, 10000)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"goals": goalDTOs(goals)})
}

// addGoal 수동으로 목표 하나 추가: DB 저장(작업 루트 spawns 아래에 연결)→ '목표를 추가함' 트리거 하나를 기록해
// planner를 깨운다 → 작업 되살리기, planner가 새 목표를 근거로 달성 여부를 다시 판정하게 한다.
func (s *Server) addGoal(w http.ResponseWriter, r *http.Request) {
	t, ok := s.m.Task(r.PathValue("id"))
	if !ok {
		writeErr(w, 404, "task not found")
		return
	}
	if !s.engine.beginTaskOperation(t.ID) {
		writeErr(w, 409, "작업을 삭제하는 중이라 목표를 추가할 수 없습니다")
		return
	}
	defer s.engine.decInflight(t.ID)

	var body struct {
		Text      string `json:"text"`
		VulnClass string `json:"vulnclass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, "invalid JSON")
		return
	}
	text := strings.TrimSpace(body.Text)
	if text == "" {
		writeErr(w, 400, "목표 내용은 비울 수 없습니다")
		return
	}
	payload := map[string]any{"text": text}
	if vc := strings.TrimSpace(body.VulnClass); vc != "" {
		payload["vulnclass"] = vc
	}
	id, err := t.Store.AddGoal(payload, "human")
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	if of, _ := t.Store.OriginFactID(); of > 0 && id > 0 {
		_ = t.Store.Link(of, db.RelSpawns, id) // goal descends from the task root (origin fact)
	}
	t.NotifyGoal([]string{text}) // '사람이 목표를 추가함: …' 트리거를 기록하고 planner를 깨움
	s.reviveTask(t)              // 완료/일시정지된 작업을 실행 상태로 되돌려 계속 실행
	node, _ := t.Store.GetNode(id)
	if node == nil {
		writeErr(w, 500, "목표 저장 후 읽기 실패")
		return
	}
	writeJSON(w, 200, goalDTO(node))
}

// editGoal 수동으로 목표 텍스트(및 vulnclass) 수정: DB 수정 → '사용자가 목표를 old에서 new로 수정함' 기록 →
// 트리거로 planner를 깨운다 → 작업 되살리기, planner가 새 목표를 근거로 방향을 조정하게 한다.
func (s *Server) editGoal(w http.ResponseWriter, r *http.Request) {
	t, ok := s.m.Task(r.PathValue("id"))
	if !ok {
		writeErr(w, 404, "task not found")
		return
	}
	if !s.engine.beginTaskOperation(t.ID) {
		writeErr(w, 409, "작업을 삭제하는 중이라 목표를 수정할 수 없습니다")
		return
	}
	defer s.engine.decInflight(t.ID)

	gid, err := strconv.ParseInt(r.PathValue("gid"), 10, 64)
	if err != nil || gid <= 0 {
		writeErr(w, 400, "bad goal id")
		return
	}
	var body struct {
		Text      string `json:"text"`
		VulnClass string `json:"vulnclass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, 400, "invalid JSON")
		return
	}
	text := strings.TrimSpace(body.Text)
	if text == "" {
		writeErr(w, 400, "목표 내용은 비울 수 없습니다")
		return
	}
	node, err := t.Store.GetNode(gid)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	if node == nil || node.Kind != db.KindGoal {
		writeErr(w, 404, "목표가 존재하지 않습니다")
		return
	}
	oldText := goalDTO(node).Text
	if err := t.Store.UpdateGoalPayload(gid, text, strings.TrimSpace(body.VulnClass)); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	t.NotifyGoalEdited(oldText, text) // '사람이 목표를 old에서 new로 수정함' 트리거를 기록하고 planner를 깨움
	s.reviveTask(t)                   // 추가와 동일: 작업 되살리기, 새 목표를 근거로 다시 판정
	updated, _ := t.Store.GetNode(gid)
	if updated == nil {
		writeErr(w, 500, "목표 업데이트 후 읽기 실패")
		return
	}
	writeJSON(w, 200, goalDTO(updated))
}

// deleteGoal 수동으로 목표 하나 삭제(하드 삭제, 엣지/앵커 캐스케이드 삭제): DB 삭제 → '사용자가 해당 목표 X를 삭제함' 기록 →
// 트리거로 planner를 깨워 남은 목표를 다시 판정하게 한다. 제품 결정에 따라 삭제는 작업을 [되살리지 않는다].
func (s *Server) deleteGoal(w http.ResponseWriter, r *http.Request) {
	t, ok := s.m.Task(r.PathValue("id"))
	if !ok {
		writeErr(w, 404, "task not found")
		return
	}
	if !s.engine.beginTaskOperation(t.ID) {
		writeErr(w, 409, "작업을 삭제하는 중이라 목표를 삭제할 수 없습니다")
		return
	}
	defer s.engine.decInflight(t.ID)

	gid, err := strconv.ParseInt(r.PathValue("gid"), 10, 64)
	if err != nil || gid <= 0 {
		writeErr(w, 400, "bad goal id")
		return
	}
	node, err := t.Store.GetNode(gid)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	if node == nil || node.Kind != db.KindGoal {
		writeErr(w, 404, "목표가 존재하지 않습니다")
		return
	}
	text := goalDTO(node).Text
	if err := t.Store.DeleteGoal(gid); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	t.NotifyGoalDeleted(text) // '사람이 해당 목표를 삭제함: …' 트리거를 기록하고 planner를 깨움(작업 되살리지 않음)
	writeJSON(w, 200, map[string]bool{"ok": true})
}
