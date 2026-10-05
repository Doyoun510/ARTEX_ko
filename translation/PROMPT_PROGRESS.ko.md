# 프롬프트 한국어화 진행 (U1~U9)

[`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) 기준. 단위별 작업 기록.

---

## U1 · 공용 기반  [done]
- **대상 파일**: `agent/promptcatalog.go` · `agent/prompt.go` · `agent/assembly.go` · `agent/chat.go`
- **내용**: 내장 agent 공용 프롬프트(Auto 운영 도우미 / 독립 침투 agent / 보고서 작성 agent / DefaultAssistantPrompt) + 렌더·조립·chat 경로 주석·문자열.
- **표준 출력 문구 확정**: `DefaultAssistantPrompt` = `"...간결하고 정확한 한국어로 답하라..."` → 이후 단위는 이 어투를 따른다. (`全程中文` → `전 과정 한국어`)
- **검증**: 중국어 잔여 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱 40=40·따옴표 84=84·중괄호·`%s/%v`) ✅ / `()`는 전각→반각 정규화로 균형 증가(정상) ✅
  - **go build / go test / 템플릿 렌더 / 1 run 스모크 = 미검증**(이 환경에 Go 툴체인 없음 — 로컬에서 `go build ./agent/...` 확인 필요)
- **남긴 중국어 + 보존 이유(DNT)**:
  - 도구명(`report_finding`·`update_finding_report`·`get_finding_traffic`·`traffic_search`·`insert_assets`·`list_tasks`·`spawn_task` 등) / 파라미터·키(`task_id`·`node_id`·`finding_id`·`vulnclass`·`severity`·`evidence_version`·`kind`·`exec`·`schema`·`props` 등) / enum·리터럴(`inferred`·`flag`) / 반환 리터럴(`finding recorded: <id>`) / 템플릿 치환자(`{{.Goal}}`) / 경로(`/tmp`·`<workDir>/noa/<SessionID>`) — 전부 원문 유지.
- **조정 포인트(다음 단위에 알림)**:
  - `ReporterDefaultPrompt`가 참조하는 컨텍스트 라벨 `"작업: #<id>"` — 이 라벨을 **실제로 생성하는 Go**(컨텍스트 조립부, 추정 worker/sidequestion)도 같은 "작업:"으로 맞춰야 완전 일치. 지금은 한국어로 선반영, 해당 단위에서 동기화 필요.
  - 에이전트 표시명: `渗透测试`→`침투 테스트`, `报告撰写`→`보고서 작성` (U4·U8 등에서 재등장 시 동일 표기).

---

## U2 · 플래너 계열  [done]
- **대상 파일**: `agent/planner.go` · `goals.go` · `constraints.go` (동반 `*_test.go` 없음)
- **내용**: 계획자 시스템 프롬프트(plannerDefaultTmpl) + 트리거/태세 문자열 + 목표 분해기(goalsDefaultTmpl·goalsScopeTail) + 동작 제약 주입 블록.
- **검증**: 중국어 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱·따옴표·중괄호) ✅ / `%s·%d` 개수·순서 ✅ / `{{.Goal}}`·`{{ }}` 템플릿 보존 ✅
  - **go build / test / 1 run 스모크 = 미검증**(Go 툴체인 없음 — 로컬 `go build ./agent/...`)
- **남긴 중국어/DNT**: 도구명(set_constraints·set_goals·add_task_scope·add_intent·prove_goal·TodoWrite·graph_overview·node_detail·list_facts/findings/assets·get_worker_trace/output·steer_work·kill_work 등)·JSON 키(summary·asset_ids·parent_ids·goal_id·evidence_id·reason·frontier_open·running_intents·recent_done/facts·coverage.pct 등)·state enum(open·running·done·exhausted·blocked·met·observed·inferred)·`{{.Goal}}`·`out-of-scope`·세션 id(`exp%d-goals` 등) 전부 원문.
- **조정 포인트**: userMsg 라벨 `작업 목표:`/`작업 설명:`(L141/143)과 goalsDefaultTmpl의 `'작업 목표 / 작업 설명'` 참조를 **함께** 번역해 일치시킴.

## U3~U15  [todo]
(각 단위 완료 시 위와 같은 형식으로 추가. 트랙 A=U1~U9·U15 / 트랙 B=U10~U13 / 트랙 C=U14. [`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) §1 커버리지 맵 참고)
