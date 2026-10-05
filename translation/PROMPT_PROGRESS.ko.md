# 프롬프트 한국어화 진행 (U1~U15)

[`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) 기준. 단위별 작업 기록.
※ 상태 표기는 **정적 검토**와 **실행 검증**을 분리한다(§3). Go 툴체인이 없는 환경에서 한 작업은 `[정적검토 done · 실행 미검증]`이며, 실행 검증 전까지 "완료(done)"로 단정하지 않는다.

---

## U1 · 공용 기반  [정적검토 done · 실행 미검증]
- **대상 파일**: `agent/promptcatalog.go` · `agent/prompt.go` · `agent/assembly.go` · `agent/chat.go`
- **내용**: 내장 agent 공용 프롬프트(Auto 운영 도우미 / 독립 침투 agent / 보고서 작성 agent / DefaultAssistantPrompt) + 렌더·조립·chat 경로 주석·문자열.
- **출력-언어 지시 어휘 기준**: `DefaultAssistantPrompt`의 원문 요구(간결·정확)를 유지한 채 언어만 한국어로(=`"...간결하고 정확한 한국어로 답하라..."`). 이는 **강제 통일 표준이 아니라** 이후 단위가 참고할 번역 어휘 기준이다 — 각 프롬프트의 원문 조건·범위는 보존(§2.5). (`全程中文` → `전 과정 한국어`)
- **검증**: 중국어 잔여 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱 40=40·따옴표 84=84·중괄호·`%s/%v`) ✅ / `()`는 전각→반각 정규화로 균형 증가(정상) ✅
  - **go build / go test / 템플릿 렌더 / 1 run 스모크 = 미검증**(이 환경에 Go 툴체인 없음 — 로컬에서 `go build ./agent/...` 확인 필요)
- **남긴 중국어 + 보존 이유(DNT)**:
  - 도구명(`report_finding`·`update_finding_report`·`get_finding_traffic`·`traffic_search`·`insert_assets`·`list_tasks`·`spawn_task` 등) / 파라미터·키(`task_id`·`node_id`·`finding_id`·`vulnclass`·`severity`·`evidence_version`·`kind`·`exec`·`schema`·`props` 등) / enum·리터럴(`inferred`·`flag`) / 반환 리터럴(`finding recorded: <id>`) / 템플릿 치환자(`{{.Goal}}`) / 경로(`/tmp`·`<workDir>/noa/<SessionID>`) — 전부 원문 유지.
- **조정 포인트(다음 단위에 알림)**:
  - **(미확인)** `ReporterDefaultPrompt`가 참조하는 컨텍스트 라벨 `"작업: #<id>"` — 이 라벨을 **실제로 생성하는 Go 위치가 어디인지, 번역/동기화가 필요한지 아직 미확인**(추정 worker/sidequestion). 지금은 프롬프트 내 참조만 한국어로 선반영했고, **U9 착수 시 생성부를 찾아 확인**해 일치 여부를 판단한다. 확인 전까지 계약으로 간주.
  - 에이전트 표시명: `渗透测试`→`침투 테스트`, `报告撰写`→`보고서 작성` (U4·U8 등에서 재등장 시 동일 표기).

---

## U2 · 플래너 계열  [정적검토 done · 실행 미검증]
- **대상 파일**: `agent/planner.go` · `goals.go` · `constraints.go` (동반 `*_test.go` 없음)
- **내용**: 계획자 시스템 프롬프트(plannerDefaultTmpl) + 트리거/태세 문자열 + 목표 분해기(goalsDefaultTmpl·goalsScopeTail) + 동작 제약 주입 블록.
- **검증**: 중국어 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱·따옴표·중괄호) ✅ / `%s·%d` 개수·순서 ✅ / `{{.Goal}}`·`{{ }}` 템플릿 보존 ✅
  - **go build / test / 1 run 스모크 = 미검증**(Go 툴체인 없음 — 로컬 `go build ./agent/...`)
- **남긴 중국어/DNT**: 도구명(set_constraints·set_goals·add_task_scope·add_intent·prove_goal·TodoWrite·graph_overview·node_detail·list_facts/findings/assets·get_worker_trace/output·steer_work·kill_work 등)·JSON 키(summary·asset_ids·parent_ids·goal_id·evidence_id·reason·frontier_open·running_intents·recent_done/facts·coverage.pct 등)·state enum(open·running·done·exhausted·blocked·met·observed·inferred)·`{{.Goal}}`·`out-of-scope`·세션 id(`exp%d-goals` 등) 전부 원문.
- **조정 포인트**: userMsg 라벨 `작업 목표:`/`작업 설명:`(L141/143)과 goalsDefaultTmpl의 `'작업 목표 / 작업 설명'` 참조를 **함께** 번역해 일치시킴.
- **검토 반영(2026-10-05, 정적)**:
  1. `规划者/执行者` → 용어집 §(L87-89) 확정대로 **`planner`/`worker`**로 수정(계획자/실행자 오역 교정; planner.go·promptcatalog.go pentest 프롬프트). `审计者`=감사자(에이전트 아님, 메타포 유지).
  2. renderTriggers의 `strings.Join(ev.Goals/Hints, "；")`는 §7.2 허용 목록에 없어 **원문 `；`로 복원**(factIDsYielded의 `、`만 §7.2 허용 대상이라 그건 `, ` 유지).
  3. `의도 #%d가 사용자에 의해 삭제됨` → **`사용자가 의도 #%d 삭제`**(변수 뒤 조사 회피).
  4. goals.go 오타/직역: `닮`→`단다`, `벌거벗은 루트 도메인`→`서브도메인이 없는 루트 도메인`.
  5. **관련 테스트(이름 다름) `agent/prompt_test.go` 정적 확인**: `plannerSystem` 출력에 `中间产物输出规约` 포함을 기대하는데, 이 문자열은 **`worker.go:261/268`의 `artifactSpec`(U3, 미번역)**에서 나온다. U2는 주석만 건드려 **테스트 계약 무결(정적)**. ⚠️ **U3 조정 포인트**: U3에서 artifactSpec의 `中间产物输出规约`를 번역하면 `prompt_test.go`의 기대값 3곳(L41·L42·L63)도 같은 커밋에서 갱신해야 한다.
- **상태**: 번역·정적 검토 완료 / 관련 테스트는 정적 대조만(실행 안 함) / **go build·test·1 run = 미실행**.

## U3~U15  [todo]
(각 단위 완료 시 위와 같은 형식으로 추가. 트랙 A=U1~U9·U15 / 트랙 B=U10~U13 / 트랙 C=U14. [`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) §1 커버리지 맵 참고)
