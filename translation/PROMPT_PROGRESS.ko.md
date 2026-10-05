# 프롬프트 한국어화 진행 (U1~U15)

[`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) 기준. 단위별 작업 기록.
※ 상태 표기는 **번역·정적 검토·보류·실행 검증**을 분리한다(GUIDE §3). 실행 환경이 없거나 사용자가 실행을 제한한 경우 각 실행 항목에 미검증·미실행 사유를 기록한다. 실행 검증 전까지 최종 완료(done)로 단정하지 않는다. 사람2의 현재 실행 검증은 **사용자 요청으로 미실행**이다.

※ 아래 U1·U2의 기존 번역·기호 비교 수치는 팀원의 작업 보고를 보존한 것이다. 이번 기준 문서 수정으로 소스·테스트를 새로 검증한 결과가 아니며, 관련 테스트 확인 보류도 별도로 기록한다.

---

## U1 · 공용 기반  [정적 검토 완료 보고 / 컨텍스트 라벨 계약 검토 보류 / 실행 미검증]
- **기존 보고 이력**: 팀원 GUIDE의 U1 "완료" 보고 및 PROGRESS의 `[정적검토 done · 실행 미검증]` 보고를 보존한다. 번역 내용·기호 비교·DNT 기록은 아래와 같으며, 컨텍스트 라벨 계약 확인·실행 검증의 새 근거는 이번 원격 기록에 없다. 현재 상태는 제목의 계약 검토 보류·실행 미검증을 유지한다.
- **대상 파일**: `agent/promptcatalog.go` · `agent/prompt.go` · `agent/assembly.go` · `agent/chat.go`
- **내용**: 내장 agent 공용 프롬프트(Auto 운영 도우미 / 독립 침투 agent / 보고서 작성 agent / DefaultAssistantPrompt) + 렌더·조립·chat 경로 주석·문자열.
- **출력-언어 지시 어휘 기준**: `DefaultAssistantPrompt`의 원문 요구(간결·정확)를 유지한 채 언어만 한국어로(=`"...간결하고 정확한 한국어로 답하라..."`). 이는 **강제 통일 표준이 아니라** 이후 단위가 참고할 번역 어휘 기준이다 — 각 프롬프트의 원문 조건·범위는 보존(§2.5). (`全程中文` → `전 과정 한국어`)
- **검증**: 중국어 잔여 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱 40=40·따옴표 84=84·중괄호·`%s/%v`) ✅ / `()`는 전각→반각 정규화로 균형 증가(정상) ✅
  - **go build / go test / 템플릿 렌더 / 1 run 스모크 = 미검증**(이 환경에 Go 툴체인 없음 — 로컬에서 `go build ./agent/...` 확인 필요)
- **남긴 중국어 + 보존 이유(DNT)**:
  - 도구명(`report_finding`·`update_finding_report`·`get_finding_traffic`·`traffic_search`·`insert_assets`·`list_tasks`·`spawn_task` 등) / 파라미터·키(`task_id`·`node_id`·`finding_id`·`vulnclass`·`severity`·`evidence_version`·`kind`·`exec`·`schema`·`props` 등) / enum·리터럴(`inferred`·`flag`) / 반환 리터럴(`finding recorded: <id>`) / 템플릿 치환자(`{{.Goal}}`) / 경로(`/tmp`·`<workDir>/noa/<SessionID>`) — 전부 원문 유지.
- **조정 포인트(다음 단위에 알림)**:
  - **(계약 검토 보류)** `ReporterDefaultPrompt`가 참조하는 컨텍스트 라벨 `"작업: #<id>"` — 이 라벨을 **실제로 생성하는 Go 위치가 어디인지, 번역/동기화가 필요한지 아직 미확인**(추정 worker/sidequestion). 현재 프롬프트 내 참조만 한국어로 선반영된 상태다. **U9 착수 시 생성부·참조부와 현재 표현을 찾아 불일치·수정 필요성을 보고**한다. 수정이 필요하면 실제 소유 단위와 변경 범위를 정한 뒤 처리하고, 위치 확인 전에는 임의로 동기화하지 않는다.
  - 에이전트 표시명: `渗透测试`→`침투 테스트`, `报告撰写`→`보고서 작성` (U4·U8 등에서 재등장 시 동일 표기).

---

## U2 · 플래너 계열  [번역·정적 검토 완료 보고 / 관련 테스트 확인 보류 / 실행 미검증]
- **원격 작업 보고 이력**: `aa0edc0`의 U2 번역과 팀원의 `[정적검토 done · 실행 미검증]` 보고를 보존한다. 기존 "동반 `*_test.go` 없음" 보고는 같은 이름의 테스트가 없다는 보고로 기록하며, 다른 이름의 관련 테스트까지 없다고 확정하지 않는다. 원격의 build·test·1 run 스모크 미검증 보고(Go 툴체인 없음, 로컬 `go build ./agent/...` 확인 필요)와 아래 번역·정적 검토·DNT·라벨 동시 변경 기록을 함께 보존한다.
- **대상 파일**: `agent/planner.go` · `goals.go` · `constraints.go` (같은 이름의 테스트 없음으로 보고됨; 관련 테스트 확인은 별도 필요)
- **내용**: 계획자 시스템 프롬프트(plannerDefaultTmpl) + 트리거/태세 문자열 + 목표 분해기(goalsDefaultTmpl·goalsScopeTail) + 동작 제약 주입 블록.
- **검증**: 중국어 0 ✅ / 전각부호 0 ✅ / 기호 parity(백틱·따옴표·중괄호) ✅ / `%s·%d` 개수·순서 ✅ / `{{.Goal}}`·`{{ }}` 템플릿 보존 ✅
  - **go build / go test / 1 run 스모크 = 미검증**(기존 보고 사유: Go 툴체인 없음)
  - **실제 템플릿 렌더 = 미검증**(실행 결과가 기록되지 않음; `{{ }}` 보존 정적 확인과 구분)
- **남긴 중국어/DNT**: 도구명(set_constraints·set_goals·add_task_scope·add_intent·prove_goal·TodoWrite·graph_overview·node_detail·list_facts/findings/assets·get_worker_trace/output·steer_work·kill_work 등)·JSON 키(summary·asset_ids·parent_ids·goal_id·evidence_id·reason·frontier_open·running_intents·recent_done/facts·coverage.pct 등)·state enum(open·running·done·exhausted·blocked·met·observed·inferred)·`{{.Goal}}`·`out-of-scope`·세션 id(`exp%d-goals` 등) 전부 원문.
- **조정 포인트**: userMsg 라벨 `작업 목표:`/`작업 설명:`(L141/143)과 goalsDefaultTmpl의 `'작업 목표 / 작업 설명'` 참조를 **함께** 번역해 일치시킴.
- **검토 반영(2026-10-05, 정적)**:
  1. `规划者/执行者` → 용어집 §(L87-89) 확정대로 **`planner`/`worker`**로 수정(계획자/실행자 오역 교정; planner.go·promptcatalog.go pentest 프롬프트). `审计者`=감사자(에이전트 아님, 메타포 유지).
  2. renderTriggers의 `strings.Join(ev.Goals/Hints, "；")`는 §7.2 허용 목록에 없어 **원문 `；`로 복원**(factIDsYielded의 `、`만 §7.2 허용 대상이라 그건 `, ` 유지).
  3. `의도 #%d가 사용자에 의해 삭제됨` → **`사용자가 의도 #%d 삭제`**(변수 뒤 조사 회피).
  4. goals.go 오타/직역: `닮`→`단다`, `벌거벗은 루트 도메인`→`서브도메인이 없는 루트 도메인`.
- **관련 테스트 확인(정적 대조, 실패 미확정)**: `agent/prompt_test.go`는 `plannerSystem` 결과에 `中间产物输出规约`가 포함되는지 검사한다. 최신 브랜치에서 정적 대조한 결과, 이 생성 문자열은 **`worker.go:261/268`의 `artifactSpec`(U3, 미번역)**에서 나오며 U2는 주석만 건드려 **생성 문자열·기대값이 모두 중국어 원문 그대로** = 계약 정적 무결. 이 기록은 **테스트 실패를 확정한 것이 아니며**, 테스트 실행은 **별도 허용 전에는 하지 않는다**. ⚠️ **U3 조정 포인트(편집 담당=U3 작업자)**: U3에서 `artifactSpec`의 `中间产物输出规约`를 번역하면 `prompt_test.go` 기대값 3곳(L41·L42·L63)도 같은 커밋에서 갱신해야 한다.
- **상태**: 번역·정적 검토 완료. **실행 검증 완료(2026-10-05)**: `go build ./agent/...` OK, `go test ./agent/` PASS(0.118s, `TestRenderSystemOverrideAndFallback` 포함 — plannerSystem+`中间产物输出规약` 계약 통과). Go 1.26.3. 1 run 스모크(실제 에이전트 구동)는 미실시.

## U3 · 워커 실행  [정적검토 done · 실행 검증 PASS(2026-10-05)]
- **대상 파일**: `agent/worker.go`·`wrapup.go`·`compaction.go`·`noa.go`·`capture.go`·`coldgraph.go` + **관련 테스트 `agent/prompt_test.go`**(이름 다름, artifactSpec 계약 동기화).
- **내용**: 워커 시스템 프롬프트(workerDefaultTmpl)·트래픽 블록·산출물 규약(artifactSpec/workerArtifactSpec)·의도/overview user 메시지·마무리 프롬프트(planner/mainagent/worker/generic + 작업 타임아웃)·압축 프롬프트(compressionSystemPrompt)·noa·캡처 오류 문구.
- **계약 동기화(핵심)**: `artifactSpec`의 `中间产物输出规约` → `중간 산출물 출력 규약`으로 번역하면서 **같은 커밋에서 `prompt_test.go` 기대값 3곳(`strings.Contains(... "中间产物输出规约")`)도 함께 변경**. (U2에서 예고한 U3 조정 포인트 처리 완료)
- **검증**: 중국어 0(의도 보존 제외) ✅ / 전각부호 0 ✅ / parity·`%verb`·`{{ }}` 보존 ✅ / **go build ./agent/... OK · go test ./agent/ PASS(0.062s, prompt_test 포함)**.
- **의도 보존(번역 안 함)**:
  - `wrapup.go:102` 주석의 `docs/任务级超时与收尾设计.md` — 실제 파일 경로 참조(DNT).
  - `prompt_test.go`의 픽스처(`目标:{{.Goal}} 范围:{{.Scope}}`·`走代理`·입력 `拿下X` 등) — renderSystem 치환 **메커니즘 검증용 테스트 데이터**라 GUIDE §1대로 일괄 번역 안 함(표시 문자열 아님).
- **DNT**: `执行者`→`worker`·`规划者`→`planner`(용어집), 도구명(insert_assets·record_fact·report_finding·traffic_search/get/blob·add_intent·prove_goal·TodoWrite 등)·키(intent_id·summary·detail·evidence·confidence·asset_ids 등)·enum(observed·inferred·met)·경로(`/tmp`·`<workDir>/...`·`@blob sha256`)·`workerChatMarker` 포맷·`{{.ProxyAddr}}` 등 전부 원문.

## U4 · 메인/리테스트/발견  [정적검토 done · 실행 검증 PASS(2026-10-05)]
- **대상 파일**: `agent/mainagent.go`·`retester.go`·`finding_workflow.go`·`finding_recorder.go`.
- **내용**: 메인 agent 시스템 프롬프트(mainAgentDefaultTmpl)·조립부 주석·문자열 / 재검증(retester) 프롬프트(RetesterDefaultPrompt) / 발견 워크플로의 도구 설명·note·가이던스(findingIDGuidance·report_finding/add_hint note·트래픽 증거 인계·보고 전 자동 연관 블록·HintTrafficSchema 설명·evidence_hint_id 오류 메시지) / 발견 기록 가이던스(findingTrafficGuidance).
- **관련 테스트 확인(이름 다른 것 포함)**:
  - `finding_workflow_test.go` — `traffic_refs`·`TCP`·`evidence_hint_id` 같은 **DNT 토큰만** 검사 → 보존하므로 변경 불필요.
  - `finding_recorder_test.go:102` — `strings.Contains(got, "未登记")` 는 **`agent/tools.go:1261`(U5 소관)** 의 문자열을 검사한다. U4 파일이 생성하지 않으므로 기대값·소스 모두 **건드리지 않음**(U5에서 tools.go 번역 시 동기화 대상).
  - `toolcatalog_test.go`·`side_questions_test.go` — U4 번역 문자열에 대한 기대값 없음(중국어 비교 없음). 변경 불필요.
- **검증(정적)**: 번역 대상 중국어·전각 문장부호 잔여 0(계약 앵커 없음) ✅ / `{{.Goal}}` 치환자 보존 ✅ / `%d`(evidence_hint_id 오류 메시지) 개수·순서 보존 ✅ / `\n`·`\n\n` 이스케이프·`[]`·JSON 구조 보존 ✅ / role enum `baseline / proof / verification / supporting` 원문 유지 ✅.
- **실행 검증(2026-10-05, Go 1.26.3)**: `gofmt -l` 무출력(정렬 OK) ✅ / `go build ./...` OK ✅ / `go test ./agent/` PASS(0.078s, finding_workflow_test·finding_recorder_test·prompt_test 포함) ✅. 1 run 스모크(실제 에이전트 구동)는 미실시.
- **용어·표기**:
  - `规划者`→`planner`·`执行者`→`worker`(용어집), `主 agent`→`"메인 agent"`, `报告 Agent`→`보고서 Agent`, `任务 Agent`→`작업 Agent`(소스 영문 `Auto`/`Planner`/`Worker` 대문자는 원문 유지).
  - `段 [A]`→`섹션 [A]`, `中间产物输出规约`→`중간 산출물 출력 규약`(U3 동기화), `收尾`→`마무리`, `接管`→`인계`(chat.go 선례).
  - 출력 언어 지시(§2.5): retester `使用简洁中文答复`→`간결한 한국어로 답하라`(간결 조건 보존). mainagent `用人话简洁回复`→`쉬운 말로 간결하게 답하고`.
  - verdict enum `reproduced`/`fixed`/`inconclusive` DNT, 괄호 주석만 번역. 발견 처리 상태 라벨 `「已修复」`→`'수정 완료'`(§확정#3, 내부 enum `fixed`는 별개 DNT).
  - `「」`→`'…'`(§6), 산문 전각부호→반각, 열거 `、`→`·`.
  - `evidence_hint_id=%d` 오류 메시지는 콜론 형식으로 변수 뒤 조사 회피(§변수 뒤 조사).
- **DNT**: 도구명 전부(report_finding·update_finding_report·get/bind_finding_traffic·list_findings·list_task_findings·node_detail·get_task_node_detail·get_finding_retest_context·record_finding_retest_result·traffic_search/get·add_hint·add_task_hint·set_goals·set_constraints·add_intent·steer_work·graph_overview·list_facts/assets·get_worker_output 등)·파라미터/키(finding_id·finding_node_id·traffic_refs·evidence_hint_id·intent_id·evidence_version·version·verdict·summary·evidence·traffic_id·role·note·hints·text·priority·type=allow/deny 등)·role enum(baseline/proof/verification/supporting)·약어(HTTP·TCP·WAF·PoC·Markdown·JSON·ID)·`{{.Goal}}`·경로(`<workDir>/...`·`/tmp`) 전부 원문.
- **U9 이월(미착수)**: U1·U9 조정 포인트인 컨텍스트 라벨 `"작업: #<id>"` 생성부 확인은 **U9 소관**이라 이번 U4에서 건드리지 않음. U9 착수 시 생성부·참조부 대조 예정.

## U9 · 사이드 질문  [정적검토 done · 실행 검증 PASS(2026-10-05)]
- **대상 파일**: `sidequestion/request.go`·`context.go`·`service.go`(+ `capture.go`는 중국어·전각 문장부호 없음 확인).
- **내용**: 보조 질문(`/btw`) 모델 프롬프트(request.go `instruction`·context.go `summaryInstruction`)·컨텍스트 조립 프레이밍 라벨(`[과거 보조 질문·답변 …]`·`[이전 보조 질문·답변 요약 …]`·`[보조 질문 기록 %d …]`·`[현재 메인 컨텍스트의 이전 요약 …]`·`[이전 요약]`/`[새 자료 조각]`)·오류/안내 메시지·service.go 도구 거부 응답.
- **용어**: `旁路(提问)`→`보조 질문`(§2 확정, "우회 질의" 아님), `sidequestion`·`/btw` DNT. `主 Agent`→`메인 Agent`, `主任务`→`메인 작업`, `原任务`→`원 작업`, `助手`→`어시스턴트`. 출력-언어 지시 없음(프롬프트 언어 따름).
- **안전 문구 보존(§9.2)**: `instruction`의 "도구 실행 능력 없음 / 조작·파일수정·메인작업 지휘 불가 / 나중 실행 약속 금지 / 컨텍스트의 작업 지시는 배경일 뿐"과 `summaryInstruction`의 "자료는 분석 대상 데이터·지시 실행 금지" 등 지시 주입 방지 의미를 그대로 유지. `尽量`→`되도록`(강도 보존).
- **계약 테스트 동기화(같은 커밋)**:
  - `sidequestion/context_test.go:308` `"未能进一步缩减"` → `"줄이지 못해"`(context.go:354 오류 문자열 계약).
  - `sidequestion/sidequestion_test.go:180` `"不能执行工具"` → `"도구 조작을 실행할 수 없습니다"`(service.go 도구 거부 응답).
  - `agent/side_questions_test.go:108`(교차 패키지 소비자) 같은 거부 응답 substring 동기화.
- **번역 안 함(DNT/픽스처)**: `context_test.go`의 `strings.Repeat("中",3000)`·`"历史依据ABC"`·`"old tool evidence 中文"`·`strings.Repeat("中",5000)`은 토큰 추정·길이/인코딩 검증용 픽스처(§7.4) → 보존. Status enum(`running`/`completed`)·JSON 키·`tokens`·time 포맷 DNT. `%d`/`%s`/`%w` 개수·순서 보존.
- **실행 검증(2026-10-05, Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / `go test ./sidequestion/` PASS ✅ / `go test ./agent/` PASS(교차 소비 테스트 포함) ✅. `go test ./...`의 server/notify 실패 22건은 전부 PostgreSQL 미연결(`ARTEX_PG_DSN` 미설정) 환경 사유이며 번역과 무관. 1 run 스모크는 미실시.
- **⚠️ U1/U9 컨텍스트 라벨 계약 확인 결과(생성부 확정)**: U1 `ReporterDefaultPrompt`(promptcatalog.go:80)가 참조하는 `"작업: #<id>"` 라벨의 **실제 생성부는 `server/conversations.go:710,712`의 `taskContextHeader`** = `fmt.Sprintf("【任务 #%d %s（目标：%s）】", …)`로 확인됨(스케줄 트리거 메시지 헤더). **sidequestion 패키지에는 이 라벨 생성이 없음** → U9 번역은 독립.
  - **불일치**: 생성부는 아직 중국어 `任务`(공백+`#`, 전각 `【】`), 참조부는 U1에서 한국어 `작업: #<id>`(콜론)로 선반영됨. **생성부 소유 = U11(server/, 현재 담당 미정)**.
  - **권고(U11 처리 대상)**: U11에서 `taskContextHeader` 번역 시 `任务`→`작업`으로 참조부 용어와 맞추고, 포맷(공백 vs 콜론·`【】`)을 U1 참조 문구와 통일할지 결정. 모델 힌트라 엄격 파서 계약은 아니나 **용어 일치 필요**.

## 라벨 계약 해소 (U1↔U11 교차, 사용자 지시로 처리 · 2026-10-05)  [정적검토 done · 실행 검증 PASS]
U9에서 확정한 `"작업 #<id>"` 라벨 계약을 사용자 지시로 처리했다. 생성부가 어떤 파서에도 소비되지 않음(백엔드·프런트 grep 0건)을 확인해 **§5.7 계약이 아니라 모델이 읽는 컨텍스트 헤더**로 판정, 번역 가능으로 처리.
- **대상(트리거 컨텍스트 조립 함수군, `server/conversations.go`)**: `taskContextHeader`(L710·712)·`mergeTriggeredRuns`(L736·751·755)·`mergeAllRuns`(L772·782·787) 8개 모델-읽기/표시 문자열. `任务 #%d`→`작업 #%d`, 전각 `【】`→`[]`, `（目标：…）`→`(목표: …)`, `触发`→`트리거`, `合并触发`→`병합 트리거`, `共 N 个任务`→`작업 N개`. `%d`/`%s` 개수·순서·`task#%d` 식별자 보존. 조사는 변수 뒤 받침 의존 조사 회피(`의`/`건`/`개` 사용).
- **참조부 통일(U1)**: `agent/promptcatalog.go:80` 리포터 프롬프트 참조를 `"작업: #<id>"`→`"작업 #<id>"`로 바꿔 생성부 `[작업 #%d …]`와 `작업 #` 마커 일치.
- **계약 테스트 동기화(같은 커밋)**: `server/trigger_merge_test.go` 4곳 — `── 触发 `→`── 트리거 `, `触发 39`→`트리거 39`, `共 2 个任务`→`작업 2개`, `【任务 #72`→`[작업 #72`. 픽스처 입력(`longGoal`·`f2-05 逆向`·`【本次由工具调用触发】…`·`定时触发正文`·`很`×5000)은 테스트 입력 데이터라 보존.
- **범위 한정**: `conversations.go`의 나머지 중국어(일반 U11)·`关联任务 #%d 不存在`류 에러 메시지는 **이번 범위 밖**(U11 본 작업에서 처리).
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / `go test ./server/ -run 'TestMerge|TestFinalTrigger|TestTaskContextHeader'` PASS(6건, DB 불필요 순수 함수) ✅ / `go test ./agent/ -count=1` PASS ✅.

## U11 · 서버 API (트랙 B)  [진행 중 — 묶음 A(실행 루프 코어)]
사용자 결정으로 사람1이 U11을 하위 묶음으로 나눠 착수. 묶음 A = 실행 루프 코어(트리거→런→목표 흐름). U11엔 모델-읽기 문자열이 섞여 있어(GUIDE 명시) 문자열별로 소비자(모델/사용자/로그/계약)를 구분해 규칙 적용.

### 묶음 A-1 (제어·목표·스케줄 클러스터)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/scheduler.go`·`task_control.go`·`goals.go`·`goals_api.go`·`engine_timeout.go`.
- **문자열 분류**: 트리거 컨텍스트 마커(`[이번은 …로 트리거됨]`)·도구 오류(의도/작업 제어)·Causef Short/Text(具名原因)=모델/사용자 읽기 → 번역. 런 타이틀·writeErr·로그·Summary=표시 → 번역. 코드 키 DNT.
- **DNT 보존**: state enum(running/paused/open/done/failed/timeout/completed/met/deleted)·Causef **code**(queued_for_admission·llm_unavailable_queued·first_run)·MaxDuration·SetTaskStatusGuarded·settling/drain/grace/deadline·도구/필드명·`task_ids`·`%s/%d/%v/%w`(개수·순서)·`docs/任务级超时与收尾设计.md`(실제 파일 경로, U3 선례와 동일 보존).
- **용어**: 触发→트리거, 终态→종료 상태, 收尾→마무리, 墙钟→실제 경과 시간(wall-clock), 并发→동시 실행, 已排队→대기열 등록됨, 复活(reviveTask)→작업 되살리기(恢复=재개와 구분), 规划者→planner. 전각 `【】`→`[]`, `「」`→`'…'`.
- **계약 확인**: 번역 문자열을 검사하는 테스트 없음(배치 A-1). `已截断`는 `chat_mentions.go`(별도 파일) 소유라 scheduler.go의 것과 무관. 프런트가 goal 에러(`目标不存在` 등)를 매칭하지 않음(§7.2 안전) 확인.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 순수 서버 테스트(trigger merge 6건) PASS ✅. DB 필요 서버 테스트는 PostgreSQL 미연결로 **미검증**.
- **주의(다음 배치/프런트)**: `conversations.go:112,443`의 `新对话`는 §5.7 프런트 매칭 계약 → **DNT 유지**(프런트 B레인 동시 변경 때 처리). 다른 파일의 동명 에러(`任务正在删除` 등)는 각자 독립 리터럴이라 해당 파일 번역 시 개별 처리.
### 묶음 A-2 (대화·작업 LLM 클러스터)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/conversations.go`(나머지, `新对话` 제외)·`task_llm.go`.
- **conversations.go**: 대화/메시지 HTTP writeErr·첨부/재검증 상태 문구·트리거 behavior 주석. `附件消息`→첨부 메시지, 复测→재검증.
- **task_llm.go**: 작업 LLM 설정 체인·같은 provider 안전 윈도우 재시도 로직 주석·로그·LLM 감사 라벨(`설정 #%d`·`기본 설정`·할당량/전환 요약). failover·committed·profile·provider·model_error·doStream·HTTP 코드 등 DNT.
- **DNT/보존**: `新对话`(§5.7 프런트 계약)·L561 `已手动停止`(norma SDK가 내보내는 문자열 참조, §5.5 DNT)는 **원문 유지**. `%s/%d/%v`·manual/failed/stopped enum 보존. 전각 `（）：`→반각.
- **계약 확인**: 번역 문자열을 검사하는 테스트 없음. 프런트가 이 에러/제목 문자열을 매칭하지 않음(§7.2) 확인.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅. DB 필요 서버 테스트는 미검증.

### 묶음 A-3 (엔진 루프)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/engine.go`.
- **모델-읽기 문자열(프롬프트 규칙 적용)**: `emptyTurnNudge`(공회전 라운드 이어 실행 지시)·planner 실시간 교정 주입 메시지·work 제어 오류·goalless 종료 Summary. 지시 강도·조건 보존.
- **로그·주석**: 엔진 루프(계획/워커 수명주기·타임아웃 마무리·model_error 재시도·공회전 nudge) 대량 번역.
- **DNT**: 상태 enum(open/running/paused/done/failed/timeout/exhausted/blocked/stopped/met)·model_error·first_run_at·deadline/drain/settling/frontier·cancelExec/SetTaskStatusGuarded·stop_reason/end_turn/tool_use·norma 경로·`%s/%d/%v/%w/%q`(개수·순서)·`docs/任务级超时与收尾设计.md`(파일 경로).
- **계약 테스트**: `engine_emptyturn_test.go`는 `emptyTurnNudge` **상수 참조**로 비교 → 상수 값 번역돼도 통과(기대값 변경 불필요). 테스트 내 중국어는 픽스처·케이스명(§7.4 DNT)이라 미변경.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 순수 테스트(empty-turn·trigger merge) PASS ✅. DB 필요 서버 테스트는 미검증.

### 묶음 A-4 (오케스트레이션)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/orchestration.go`. **묶음 A 완료**.
- **모델-읽기 도구 설명(U5급 DNT 주의)**: 크로스 작업 오케스트레이션 도구(list_tasks·list_llm_profiles·spawn_task·pause_task·get_task_graph·list_task_findings·add_task_hint·get_task_worker_trace·list_task_worker_traces·search_task_worker_traces·node_detail·update_finding_report)의 설명·파라미터 설명 번역. 도구명·파라미터 키·`finding recorded: <id>`·enum 전부 보존.
- **모델-읽기 프롬프트**: `reporterToolCallMessage`(리포터 트리거 메시지, 현재 시드값) 번역.
- **⚠️ 계약 DNT — `reporterToolCallMessageV1`**: 구 버전(0.3.8) 트리거 메시지로, `upgradeReporterTriggerMessage`가 **기존 DB 레코드와 바이트 단위 일치 비교**(L719 `!= reporterToolCallMessageV1`)해 마이그레이션 여부를 판정하는 **역사적 리터럴**. 번역하면 구 DB 매칭이 깨지므로 **원문 보존**(V1 상수만, 설명 주석은 번역).
- **표시/오류/로그/주석**: 작업 생성·바인딩 마이그레이션·reporter agent 시드 로그·주석 번역. 리포터 agent 표시명 `报告撰写`→`보고서 작성`(U1 동일).
- **DNT**: 도구/파라미터/상태 enum·settings flag(`*_v1/v2/v3`)·`finding recorded: <id>`·`docs/跑分编排`(설계 문서 참조)·포맷 지정자.
- **계약 확인**: 테스트/프런트 매칭 없음(`报告撰写`·`激活配置` 등 프런트 grep 0). 
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 순수 서버 테스트(trigger merge·empty-turn) + agent 회귀 PASS ✅. DB 필요 서버 테스트는 미검증.

### 묶음 A 요약
실행 루프 코어 9파일 완료(A-1 scheduler·task_control·goals·goals_api·engine_timeout / A-2 conversations·task_llm / A-3 engine / A-4 orchestration). DB 필요 서버 테스트는 전 배치 공통으로 PostgreSQL 미연결 미검증.
- **U11 잔여 묶음(todo)**: B(HTTP 핸들러/관리 API) / C(알림·동기화·업데이트). 그 외 server/ 미분류 파일 다수.

### 묶음 B-1 (인증·DTO·커스텀/플랫폼 도구)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/auth.go`·`dto.go`·`customtool.go`·`platform_tools.go`.
- **auth.go**: 로그인/비밀번호 HTTP 에러·JWT 로그. 프런트가 이 인증 에러를 substring 매칭하지 않음(표시용) 확인 → 번역.
- **dto.go**: DTO 구조체 필드/문서 주석(전부 주석). 식별자·enum(state=running/deleted·max_tokens 등) 보존.
- **customtool.go**: 커스텀 도구 실행기 주석·CRUD HTTP 에러(writeErr)·모델-읽기 도구 에러(actool.Errorf)·도구 설명. `该 key 已存在`는 프런트 매칭 페이지 없음(§5.7의 `已存在`는 skills 페이지 전용, customtool 아님) 확인 → 번역.
- **platform_tools.go**: Auto agent용 플랫폼 조작 도구(delete_assets_by_host·create_skill·update_skill_file·create_custom_tool·update_custom_tool·create_mcp·update_mcp) 설명·파라미터 설명·모델-읽기 에러. `skill 已存在`·`该 key 已存在`는 모두 **actool.Errorf(모델-읽기 도구 에러)**라 프런트 미노출 → 번역.
- **⚠️ §5.7 계약 위치 확인(묶음 B 핵심)**: 프런트가 substring 매칭하는 `已存在`는 **`web/src/.../system/skills/page.tsx:370`(`msg.includes("已存在")`)** 한 곳뿐이며, 이를 먹이는 백엔드는 **`server_mgmt.go`의 skill 생성 HTTP 에러**. 묶음 B-1 파일에는 해당 HTTP 계약 문자열이 없음. **B-2(server_mgmt.go) 번역 시 그 `已存在`는 DNT 유지 or 프런트(B레인) 동시 변경 필요.**
- **DNT**: 도구/파라미터/상태 enum·`docs/自定义工具设计.md`(파일 경로)·포맷 지정자 보존.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 테스트 커플링 없음. DB 필요 서버 테스트는 미검증.
### 묶음 B-2a (관리 API: server_mgmt.go)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/server_mgmt.go`(agent/도구/MCP/skill/LLM 설정·재시도 정책 관리 HTTP 핸들러).
- **⚠️ §5.7 계약 보존**: skill 업로드 중복 에러(L1361)는 프런트 `system/skills/page.tsx:370`의 `msg.includes("已存在")` 덮어쓰기 플로우가 매칭 → **`已存在`를 앵커로 보존**(한국어 번역문 안에 `(已存在)` 병기). agent 생성 에러(L266 `该 key 已存在`)는 프런트 매칭 없어 전문 번역.
- **번역**: HTTP writeErr·로그·JSON 필드 주석·프롬프트 템플릿 변수 설명(Now/DataDir)·재시도 정책 주석. `段 [A]`→`섹션 [A]`.
- **DNT**: `已存在` 앵커(L1361)·식별자(llm_profile_id·profile·SetAgentTriggerBehavior·Chat Completions·anthropic·openai·Registry·GBK/RLO/NBSP)·`{{.Var}}`·`SKILL.md`·포맷 지정자.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 테스트·프런트(기타 에러) 커플링 없음. DB 필요 서버 테스트 미검증.
### 묶음 B-2b (핵심 API: server.go)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/server.go`(라우팅·작업/의도/목표 CRUD·LLM 설정/연결 테스트·설정·채팅·seed). **묶음 B 완료**.
- **⚠️ §5.7 계약 DNT — fallbackChat 명령 접두**: `strings.HasPrefix(m, "意图")`·`strings.HasPrefix(m, "提示")`와 그 `TrimPrefix`(L3820-3825)는 **비교값이라 원문 보존**. 안내 문구는 §5.7대로 영문 `intent`/`hint` 사용을 안내하도록 번역(둘 다 HasPrefix 허용).
- **번역**: 라우트/셋업 주석·writeErr·로그·DTO 필드 주석·seed 의도 요약(모델-읽기)·LLM 연결 테스트/준비 상태 에러·설정 맵 주석·채팅 핸들러. `主 Agent`→`메인 Agent`, `规划者`→`planner`, 전각부호→반각.
- **DNT**: `意图`/`提示` 접두(L3820-3825)·상태 enum·설정 맵 키·`{{.Var}}`·식별자·포맷 지정자. seed asset source 값은 §7.1(새 DB) 번역.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / 순수 서버 테스트 PASS · 테스트/프런트 커플링 없음. DB 필요 서버 테스트는 미검증.

### 묶음 B 요약
HTTP 핸들러/관리 API 6파일 완료(B-1 auth·dto·customtool·platform_tools / B-2a server_mgmt / B-2b server). §5.7 계약 2건을 앵커 보존으로 처리: skill 업로드 `已存在`(server_mgmt.go, 프런트 skills 페이지 매칭)·fallbackChat `意图`/`提示` 접두(server.go). DB 필요 서버 테스트는 전 배치 공통 미검증.
### 묶음 C (알림·동기화·업데이트)  [정적검토 done · 실행 검증 일부 PASS]
- **대상(완료)**: `server/llmretry.go`·`sync_scopesentry.go`·`manager.go`·`update.go`·`notifier.go`·`notify_api.go`.
- **notifier.go/notify_api.go**: 취약점 IM 푸시 전달 엔진·채널 CRUD API 주석·로그·내부 에러·테스트 메시지. **푸시 메시지 템플릿·StatusLabel은 notify/ 패키지(U10) 소관이라 여기엔 없음**(전달/재시도/리스/토큰버킷 로직만).
- **계약 테스트 동기화**: `notify_api_test.go` 2곳 — `渠道类型无效`→`채널 유형이 유효하지 않음`, 테스트 메시지 매칭 `测试`→`테스트`(notify_api.go 테스트 메시지 Name/Summary 번역에 맞춤). e2e 테스트는 PostgreSQL 필요라 미실행.
- **sync_scopesentry.go**: ScopeSentry 자산 동기화 주석·에러·경고. ScopeSentry/MCP/list_projects_data 등 DNT.
- **manager.go/update.go/llmretry.go**: DTO 필드/설정/재시도 정책/자가 업데이트 주석·로그·SSE 진행 메시지. `docs/任务级超时与收尾设计.md`·`docs/LLM重试设计.md`(파일 경로) DNT.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / `go vet ./server/`(테스트 컴파일 포함) clean ✅ / 순수 서버 테스트(trigger merge·empty-turn·notify budget) PASS ✅. DB 필요 서버 테스트는 미검증.

### 묶음 D (server/ 긴 꼬리) — 진행 중
**D-1(계약·증거 클러스터)  [정적검토 done · 실행 검증 일부]**: `chat_mentions.go`·`task_archives.go`·`finding_retests.go`·`finding_traffic.go`·`finding_workflow.go`.
- **§5.7 계약 보존**: `chat_mentions.go` 멘션 토큰 정규식·맵 키(`漏洞/资产/…`)는 **DNT**(파서 계약), `已截断`은 번역(`잘림`)+테스트 동기화. `task_archives.go` `归档不存在`은 **앵커 보존**(한국어+병기, 프런트 tasks 페이지 매칭). `finding_retests.go` `内置默认`(프롬프트 버전 라벨)은 db/config.go(U12)와 일치 위해 **DNT 유지**(U12 번역 시 함께).
- 모델-읽기 도구 설명(get_finding_retest_context·record_finding_retest_result·get_finding_traffic·bind_finding_traffic·legacy traffic_search)·에러·안전 문구 번역. 도구/파라미터/enum·`finding recorded`·`report_finding` 등 DNT.
- **검증**: go build ./server/ OK, gofmt 무출력, 테스트 커플링 없음(chat_mentions_test 동기화 완료). DB 필요 테스트 미검증.

### U11 잔여(todo) — server/ 긴 꼬리 (D-2 이후)
A/B/C가 핵심 핸들러·오케스트레이션·알림을 덮었고, `server/`엔 소규모 파일이 다수 남음. **⚠️ §5.7 계약 주의 파일**:
- `chat_mentions.go`(12) — 멘션 토큰(`@[漏洞#id]` 종류명)·`已截断`(chat_mentions_test 매칭). 프런트·테스트 동시 고려.
- `task_archives.go`(19) — `归档不存在`(프런트 tasks 페이지 `includes("归档不存在")` 매칭).
- `finding_retests.go`·`finding_traffic.go`·`finding_workflow.go` — 재검증/증거 도구 설명(모델-읽기).
- `intercept.go`·`asset_intercept.go`·`task_intercept.go`·`side_questions.go`·`constraints_api.go`·`chatupload.go`·`skill_zip.go`·`workspace.go` 등.
- 그 외: assembly·assets·findings_groups·intent_intervention·llmpool·logsink·mcpdiscover·task_* ·triggers·webui_* 등.

## U5~U8·U10·U12~U15  [todo]
(각 단위 완료 시 위와 같은 형식으로 추가. 트랙 A=U1~U9·U15 / 트랙 B=U10~U13 / 트랙 C=U14. [`PROMPT_GUIDE.ko.md`](./PROMPT_GUIDE.ko.md) §1 커버리지 맵 참고)

### U6 · 첫 파일럿  [대상 번역·정적 검토 완료 / U6 잔여 미작업 / 실행 미검증]
- **완료 대상**: `agent/terminalreason.go`의 중국어 종료 설명·요약·상세 문구와 `agent/terminalreason_test.go`의 직접 연결된 기대 문자열. `cancelcause.go`·`provider.go`는 미작업이며 U6 전체 완료가 아니다.
- **진행 허용**: 사전 검토에서 U1·U2 보류와 직접적인 계약 의존을 발견하지 못한 이 파일럿에 한해 사용자가 진행을 허용했다. U1·U2 상태는 유지한다.
- **용어 보류 해소**: 终态·墙钟预算·钩子·多模态·具名原因·模型回合·意图粒度를 사용자 확정 용어로 등록하고 보류 10곳을 번역했다. 이번 대상에 추가 용어·계약 보류는 발견하지 못했다.
- **정적 검토**: 문자열 밖 코드·영어 주석·식별자·조건 분기·사유 코드·들여쓰기·개행·포맷 지정자와 순서·이스케이프·마크다운·JSON 구조를 보존했다. 원문 대비 텍스트 비교이며 구문·렌더·동작 검증을 대신하지 않는다.
- **기대값 대조**: `7회`는 `progressSuffix`의 요약과 `- **실행한 모델 턴**: %d회\n` 상세 항목 모두에 대응한다. 요약이 상세 첫 부분에도 포함됨을 확인했다. `결과 미반환`·`정상 반환 완료`·`실행 예산 상한`·`중단`의 대응 소스도 정적으로 대조했다. 테스트 로직은 변경하지 않았다.
- **잔여·보존 이유**: 소스의 번역 대상 중국어·전각 문장부호 잔여는 없다. 테스트의 `规划者`·`用户停止`·`用户暂停`은 외부 취소 사유 보존 검증, `已经生成的半段回答`는 실제 출력 보존 입력/기대값, `恢复`·`恢复恢复恢…`·`头\n尾`·`头`는 Unicode 길이·개행 검증 데이터로 원문 유지했다. 사유 코드·enum·실제 출력·오류·도구 입력·외부 취소 설명은 그대로 둔다.
- **범위 밖 잔여**: `AbortReason`이 반환하는 `cancelcause.go`의 취소 설명은 이번 대상이 아니므로 미번역 상태를 유지한다. 기존 DB 활동 기록의 문구도 변경하지 않는다.
- **실행 검증**: 설치·빌드·타입 검사·lint·테스트·실제 템플릿 렌더·앱·화면·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**. 실행 검증 완료나 U6 최종 완료로 표시하지 않는다.

### U6 · 실행 취소 사유 설명  [부분 번역·정적 검토 / 용어 검토 보류 / 실행 미검증]
- **대상**: `agent/cancelcause.go`의 `AbortCause` 자연어 `Short`·`Text`와 `AbortReason`의 자연어 대체 설명, `agent/terminalreason_test.go`의 직접 연결된 요약 기대값 3곳. `provider.go`는 미작업이며 U6 전체 완료가 아니다. 위 첫 파일럿 기록의 미번역·테스트 원문 표기는 당시 상태이고, 이번 변경은 아래에 별도로 기록한다. U1·U2 상태는 유지한다.
- **번역·정적 검토**: 사유 코드, 상태값, 경로, 도구명, 영어, 구조, 조건 분기, 들여쓰기와 개행을 보존했다. 16개 사유의 `Short`는 2곳을 보류하고 14곳을 번역했으며, `Text`는 8곳을 보류하고 8곳을 번역했다. `AbortReason`의 두 대체 설명은 코드 키를 유지하고 자연어만 번역했다. 연결된 테스트 기대값 3곳을 번역된 `Short`에 맞췄다. 모든 `Short`의 Unicode 문자 수는 40 이하이고 각 `Text`가 해당 `Short`보다 긴지 정적으로 대조했다. 테스트 실행 결과는 아니다.
- **용어 검토 보류**: `编排 Agent`(오케스트레이션 Agent), `删除屏障`(삭제 장벽), `优雅收尾`·`硬取消`(정상 마무리·강제 취소), `写入区`(쓰기 영역), `竞态`(경합 상태: 기존 Race Condition 확정은 결제·주문 공격에 한정), 단일 실행의 `兜底`·`软墙钟预算`(하드 타임아웃 대체 처리·소프트 실제 경과 시간 예산)은 이번 문맥의 확정 용어가 없어 해당 `Short`/`Text` 문장 전체를 원문으로 유지했다. 코드·계약 보류는 추가로 발견하지 못했다.
- **DNT·원문 입력**: 사유 코드와 `paused`·`stopped`·`exhausted`·`open` 등 상태값은 보존했다. 테스트의 실제 출력 `已经生成的半段回答`과 Unicode 길이·개행 데이터 `恢复`·`恢复恢复恢…`·`头\n尾`·`头`는 계속 원문으로 유지했다.
- **실행 검증**: 설치·빌드·타입 검사·lint·테스트·실제 템플릿 렌더·앱·화면·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**. 보류와 실행 미검증을 해소하기 전까지 U6를 최종 완료로 표시하지 않는다.

### U6 · 실행 취소 사유 보류 해소  [대상 번역·정적 검토 완료 / U6 잔여 미작업 / 실행 미검증]
- **용어 보류 해소**: 编排 Agent·删除屏障·优雅收尾·硬取消·写入区·竞态·兜底·软墙钟预算을 사용자 확정 용어로 등록하고, `agent/cancelcause.go`에 남아 있던 `Short` 2곳과 `Text` 8곳을 번역했다. 이번 취소 사유 설명에 추가 용어·계약 보류는 발견하지 못했다.
- **의미 보존**: 취소 주체·원인·발생 시점, 정상적인 마무리 대기, 90초 `drain` 유예, 재시작 조건, 논리 삭제·물리 삭제와 연쇄 삭제, `paused`·`stopped`·`exhausted`·`open` 상태 전이를 보존했다. 하드 타임아웃의 `兜底`는 이 문맥에 한정한 "최종 보호 조치"로 번역했으며 기존 다른 문맥의 번역은 유지했다.
- **소스·테스트 대응**: 이번에 번역한 두 `Short`는 `terminalreason_test.go`의 직접 기대값에 사용되지 않아 테스트 파일을 추가 수정하지 않았다. 기존에 맞춘 요약 기대값 3곳은 대응 소스와 일치한다. 사유 코드, 테스트 로직, 실제 출력 보존 입력, Unicode 길이·개행 검증 데이터는 보존했다.
- **최종 정적 검토**: 16개 사유의 모든 `Short`가 40자 이하이고 각 `Text`보다 짧음을 텍스트 분석으로 확인했다. 문자열 밖 코드·조건 분기·들여쓰기·개행과 사유 키·enum·영어·식별자·경로·숫자를 보존했다. `cancelcause.go`의 번역 대상 중국어·전각 문장부호 잔여는 없다.
- **U6 잔여·실행 검증**: `provider.go`는 미작업이므로 U6 전체 완료가 아니다. 설치·빌드·테스트·lint·실제 템플릿 렌더·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**했으며 실행 검증 완료로 표시하지 않는다.

### U6 · provider 설정·연결 경로  [번역·정적 검토 완료 / 실행 미검증]
- **대상·소비 경로**: `agent/provider.go`의 설정·재시도·스트리밍·세션 헤더·연결 테스트 설명과 주석, 사람에게 표시되는 서버 로그·오류, 모델이 읽는 연결 테스트 시스템 프롬프트를 번역했다. 설정 키·값과 provider 선택·호출·재시도·할당량 장애 조치·세션·타임아웃·컨텍스트 처리 로직은 보존했다.
- **계약·DNT**: 외부 provider의 할당량 오류를 판별하는 `余额不足`·`额度不足`·`额度已用尽`·`欠费`는 `IsQuotaExhaustedMessage`가 비교하는 계약 문자열이므로 원문을 유지했다. `docs/ARTEX-架构设计.md`는 경로이며, `OK`·`ping`·모델명·키·enum·URL·헤더·오류 판별 문자열도 보존했다.
- **관련 소비부·테스트 대조**: 연결 테스트 결과는 `server/server.go`에서 API 응답으로 전달되고, 할당량 판별은 `agent/provider_quota_test.go`와 `server/task_llm_test.go`, 원시 응답 캡처와 세션 헤더 동작은 `agent/provider_capture_test.go`·`agent/session_header_test.go`가 검증함을 읽기 전용으로 확인했다. 이번 번역 문자열을 직접 비교하는 테스트 기대값은 발견하지 못해 테스트 수정은 필요하지 않다.
- **정적 검토**: 문자열 밖 코드와 식별자·설정값·숫자·포맷 지정자 및 순서·정규식·개행·이스케이프·들여쓰기를 보존했다. `terminalreason.go`·`cancelcause.go`와 관련 테스트 기대값을 포함한 U6의 세 명시 소스는 번역·정적 검토를 마쳤으며 추가 용어·계약 보류는 발견하지 못했다.
- **실행 검증**: 설치·빌드·테스트·lint·실제 템플릿 렌더·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**했다. 따라서 U6 상태는 **번역·정적 검토 완료 / 실행 미검증**이며 최종 동작 검증 완료로 표시하지 않는다.

### U5 · 도구 카탈로그 첫 파일  [대상 번역·정적 검토 완료 / U5 잔여 미작업 / 실행 미검증]
- **대상·문자열 분류**: `agent/toolcatalog.go` 전체를 검토하고 내장 도구 카탈로그·DB 시드·agent 바인딩·런타임 재정의·기본 파라미터 주입을 설명하는 중국어 주석만 번역했다. 이 파일에는 번역 대상인 모델용 도구 설명·안내문이나 중국어 문자열 리터럴이 없다.
- **계약·코드 보존**: 도구명·카탈로그 key·agent key·파라미터명·JSON-Schema 키·enum·영어·식별자와 도구 등록·조회·선택·권한·노출 범위·로드 로직을 보존했다. 번역 대상 중국어와 전각 문장부호 잔여, 추가 용어·계약 보류는 발견하지 못했다.
- **생성부·소비부·관련 테스트 대조**: `server/assembly.go`·`server/server_mgmt.go`·`server/orchestration.go`와 `agent/assembly.go`의 시드·재정의 소비 경로를 확인했다. `agent/toolcatalog_test.go`·`server/tools_wire_test.go`·`agent/finding_workflow_test.go`는 key·바인딩·설명 전달·기본값 주입 동작을 검사하며, 번역한 주석을 비교하지 않으므로 기대값 수정은 필요하지 않다.
- **U5 잔여**: `agent/tools.go`·`agent/tools_insert.go`·`agent/tools_digest.go`는 이번 범위 밖 미작업이다. 따라서 U5 전체를 완료로 표시하지 않으며 U1·U2와 기존 U6 상태를 유지한다.
- **실행 검증**: 설치·빌드·테스트·lint·실제 템플릿 렌더·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**했다.

### U5 · cold digest 조회 도구  [대상 번역·정적 검토 완료 / U5 잔여 미작업 / 실행 미검증]
- **대상·문자열 분류**: `agent/tools_digest.go` 전체를 검토하고 `expand_digest`의 모델용 도구 설명·`id` 스키마 설명·도구 오류 결과와 중국어 개발자 주석을 번역했다. 지시 범위와 도구 사용 순서를 유지했다.
- **계약·코드 보존**: 도구명 `expand_digest`, 파라미터·JSON 키 `id`·`summary`·`state`·`confidence`·`error`, 조회 필드 `cold_digests`·`recent_facts`·`recent_done_intents`, enum과 출력 구조를 보존했다. 데이터 조회·집계·정렬·필터 로직, 숫자·포맷 지정자와 순서·개행·이스케이프·들여쓰기·원본 줄 수도 유지했다. 번역 대상 중국어·전각 문장부호 잔여와 추가 용어·계약 보류는 없다.
- **생성부·소비부·관련 테스트 대조**: `agent/tools.go`의 `graph_overview`가 cold digest 목록을 만들고 planner 도구 모음에 `expand_digest`를 등록하며, `agent/tools_insert.go`가 메인 agent 도구 모음에 등록함을 확인했다. `agent/tools_overview_test.go`와 `agent/coldgraph_test.go`를 포함한 관련 테스트를 검색했으나 이번 번역 문자열을 직접 비교하는 기대값이나 `expand_digest` 전용 테스트는 발견하지 못해 테스트 수정은 필요하지 않다.
- **U5 잔여**: `agent/tools.go`·`agent/tools_insert.go`는 이번 범위 밖 미작업이다. 따라서 U5 전체를 완료로 표시하지 않으며 U1·U2와 기존 U6 상태를 유지한다.
- **실행 검증**: 설치·빌드·테스트·lint·실제 템플릿 렌더·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**했다.

### U5 · 자산 등록·범위·조회 도구  [대상 번역·정적 검토 완료 / U5 잔여 미작업 / 실행 미검증]
- **대상·문자열 분류**: `agent/tools_insert.go` 전체를 검토하고 `insert_assets`·`add_company_scope`·`add_task_scope`·`list_untested_assets`·`list_assets`·`list_companies`의 모델용 도구 설명과 스키마 설명, 자체 오류·출처 설명 및 중국어 개발자 주석을 번역했다. 필수·금지·선택 조건, 자산 범위와 권한 경계, 입력 제약 및 도구 사용 목적을 유지했다.
- **계약·코드 보존**: 도구명·agent key·파라미터명·JSON/JSON-Schema 키·enum·실제 입력값·DSL 연산자와 필드·출력 구조를 보존했다. 도구 등록·권한·DB 쓰기·검증·조회·필터·조건 분기·반환 로직과 숫자·포맷 지정자 및 순서·개행·이스케이프·들여쓰기·원본 699줄을 유지했다. 자동 분류 계약인 `备案`은 표시 표기 `ICP 등록(备案)` 안에 보존했다.
- **생성부·소비부·관련 테스트 대조**: 모델 도구 설명과 스키마가 `writeTool`/`readTool`을 거쳐 카탈로그와 agent 도구 모음에 제공되고, 자체 오류는 도구 결과로 반환되며, 자산 출처 설명은 `SetTaskAssetSource`를 통해 저장·표시됨을 확인했다. `agent/insert_assets_test.go`·`agent/blackboard_inheritance_test.go`·`agent/tools_nil_store_test.go`·`agent/toolcatalog_test.go` 및 `db/task_assets_test.go`를 포함해 관련 테스트를 검색했으나 이번 번역 문자열의 정확한 전체 값을 비교하는 기대값은 발견하지 못해 테스트 수정은 필요하지 않다.
- **잔여·범위 밖**: 대상 파일의 한자 잔여는 자동 분류 계약 표기인 `备案`뿐이며 번역 누락·추가 용어 보류·계약 보류는 없다. 자산 차단 오류에 결합되는 `d.Reason`은 범위 밖 `db/asset_intercept_match.go`에서 생성되므로 그 안의 중국어 표시 문구는 이번 파일에서 수정하지 않았다. U1·U2의 기존 보류 상태도 유지한다.
- **U5 잔여**: `agent/tools.go`가 남아 있으므로 U5 전체를 완료로 표시하지 않는다. 완료된 `toolcatalog.go`·`tools_digest.go`와 U6 기록은 유지한다.
- **실행 검증**: 설치·빌드·타입 검사·테스트·lint·실제 템플릿 렌더·스모크·`gofmt`·Git 훅은 **사용자 요청으로 미실행**했다.

### U5 · tools.go 시작~graphOverviewData  [부분 번역·정적 검토 / 용어·계약 검토 보류 / 테스트 기대값 조정 필요 / 실행 미검증]
- **실제 범위**: `agent/tools.go` 1–552행, 파일 시작부터 `graphOverviewData()`의 닫는 중괄호까지 검토했다. `inheritedMap()`은 554행에서 시작하며 553행부터 파일 끝까지는 변경하지 않았다. 원본/현재 모두 2054줄이다.
- **번역·소비 경로**: `WriteCounts.String()`의 로그용 종류별 요약, `needExploration()`의 자체 오류, 콜백·도구 래퍼·개요·커버리지 관련 중국어 주석을 번역했다. `graphOverviewData()`는 `graph_overview` 도구 응답과 `planner.go`의 `renderGraphOverview()`를 거쳐 모델 컨텍스트에 들어간다. `worker.go`의 `renderWorkerGraphOverview()`에도 전달되며 그 경로는 `coverage`를 제외한다. DB에서 받은 목표·힌트·사실·취약점·작업 설명과 외부 오류는 변환하지 않았다.
- **용어·계약 보류**: 41·44행의 그래프 관계 `上游/下游`, 343행의 모델용 설명(`探索链路图`·`蒸馏`), 424·467–469·477·494–495·501–502·507행의 접기·콜드 영역·펼치기 관련 주석, 532행의 커버리지 안내(`容器型资产` 포함), 535행의 상태 문자열 `范围未锚定`는 문장 전체를 원문으로 유지했다. 마지막 항목은 `status` 출력값의 계약 여부도 확인 보류한다. 상세 원문·문맥·후보·근거는 `/tmp/artex-u5-tools-overview-holds.md`에 기록했다. 보류를 DNT나 번역 완료로 분류하지 않는다.
- **관련 테스트 확인**: `agent/tools_nil_store_test.go:47`의 `TestExplorationToolRefusesWithoutTask`는 `needExploration()` 오류에 `任务上下文`이 포함되는지 검사한다. 대응 소스 325행을 번역했으므로 이 기대값을 `작업 컨텍스트`로 맞춰야 한다. 테스트는 이번 허용 범위 밖이므로 수정하지 않았다. `tools_overview_test.go`·`blackboard_inheritance_test.go`·`coldgraph_test.go`와 관련 검증을 읽기 전용으로 대조했으며, 이번 변경에 직접 연결된 다른 문자열 기대값 변경은 발견하지 못했다. 실행 실패를 관측한 결과가 아니다.
- **정적 보존 검토**: 주석·허용된 자연어 문자열 밖 코드, 도구명·키·enum·ID·실제 값·조건 분기·조회·집계·필터·상속·커버리지 로직, 포맷 지정자의 종류·순서·개수와 인수, 영어·개행·이스케이프·들여쓰기·줄 수를 보존했다. 코드 구분자는 변경하지 않았다. 허용 범위 밖 내용은 작업 전과 동일하다.
- **U5·U1·U2 상태**: `tools.go` 나머지 함수와 이번 용어·계약 보류 및 테스트 기대값 조정이 남으므로 U5 전체 완료가 아니다. 이번 범위에서 U1의 작업 컨텍스트 라벨·U2의 중간 산출물 규약 테스트 보류와 직접 연결된 새 계약 의존은 발견하지 못했으며, 기존 U1·U2 상태와 U6 기록은 유지한다.
- **실행 검증**: 설치·빌드·타입 검사·lint·테스트·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 개요 범위 용어 반영  [부분 번역·정적 검토 / 두 항목 원문 유지 / 실행 미검증]
- **확정 용어 반영**: 上游·下游·探索链路图·蒸馏·折叠·冷·冷区·展开를 이번 계보·탐색 현황·cold digest 문맥으로 용어집에 등록했다. 기존 `探索图` 항목에 별칭 문맥을, `上游` 항목에 계보의 "상위" 문맥을 보충했으며 기존 레코딩 프록시의 "업스트림"은 보존했다. 직전 기록의 보류 중 이에 해당하는 9개 위치/블록을 번역했고, 의도 생성 시 힌트의 필수 반영과 계획 수립 시 먼저 호출하는 지시를 유지했다.
- **테스트 대응 해소**: `agent/tools_nil_store_test.go:47`의 `TestExplorationToolRefusesWithoutTask`에서 `任务上下文` 기대 문자열만 `작업 컨텍스트`로 변경해 `needExploration()`의 번역된 오류와 일치시켰다. 다른 테스트·입력·검증 조건·로직은 보존했다. 이전 기록의 기대값 조정 필요 상태는 이번에 정적으로 해소했으며 테스트를 실행한 결과는 아니다.
- **容器型资产 확인·보류**: 저장소의 해당 표현은 `tools.go:532`의 안내문에서만 확인했고 직접 정의는 발견하지 못했다. 자산 enum(`db/schema.sql:42–44`)에는 Docker/container 종류가 없으며, 계층은 endpoint→service→subdomain/ip→root_domain→company로 설명된다(`db/task_scope.go:386–387`). 커버리지는 범위·앵커의 자산을 유형별로 세고 fact의 직접 앵커만 테스트된 것으로 집계한다(`db/task_assets_context.go:24–53,123–152`). 따라서 하위 자산을 묶는 상위 자산을 뜻할 가능성이 높다는 추론과, 정확한 포함 종류는 미확정이라는 사실을 구분한다. 해당 용어를 확정하지 않고 532행 안내문 전체를 원문 유지했다.
- **范围未锚定 확인·원문 유지**: `tools.go:533–535`는 `cov.Denominator==0`일 때 생성한다. `db/task_scope.go:281`의 `ScopeRows==0` 설명 주석과 실제 조건을 혼동하지 않는다. `graph_overview` 도구 응답과 planner의 JSON 컨텍스트에서 자연어로 소비하며 worker의 자동 주입에서는 coverage 전체를 제외한다. `get_task_graph`도 같은 함수를 재사용하지만 현재 `delegateToTask`는 `SetTaskID`를 호출하지 않아 그 경로는 coverage 생성 조건을 충족하지 않는다. 저장소 소스·문서·관련 테스트의 정적 검색에서 이 값의 비교·분기·파싱·정규식·기대값·정확한 모델 출력 요구는 발견하지 못했다. 현재 근거로는 계약 enum이 아닌 자연어 설명값으로 판단하며 DNT로 단정하지 않는다. 사용자의 유지 지시에 따라 535행 값은 그대로 두었다. 외부 모델·DB 재정의·실행 동작은 미확인이다.
- **정적 보존·잔여**: `tools.go` 원본/현재 2054줄과 553행 이후 내용, 주석·허용 문자열 밖 코드, 영어·식별자·키·enum·포맷 지정자·인수 순서·개행·이스케이프·들여쓰기를 보존했다. 이번 범위의 중국어 잔여는 532·535행 두 원문 유지 항목이며 미분류 번역 누락은 발견하지 못했다. 상세 근거와 확인 범위는 `/tmp/artex-u5-tools-overview-contract-review.md`에 기록했다. U1·U2의 기존 보류와 팀원·U6 이력은 유지하며, `tools.go` 나머지와 두 항목이 남아 U5 전체 완료가 아니다.
- **실행 검증**: 설치·빌드·타입 검사·lint·테스트·템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅은 **사용자 요청으로 미실행**했다. 번역·정적 대조와 실행 검증을 구분한다.

### U5 · tools.go 첫 묶음 최종 정적 검토  [번역·정적 검토 완료 / 나머지 함수 미작업 / 실행 미검증]
- **남은 두 항목 해소**: 容器型资产는 이번 커버리지 안내의 "하위 자산을 묶는 상위 자산" 해석으로, 范围未锚定는 `coverage.status`의 자연어 설명값 "범위 앵커 미설정"으로 확정해 용어집에 등록했다. 직접 정의가 없는 상위 자산의 포함 종류를 임의로 열거하거나 일반 용어로 확대하지 않았다. 이전 원문 유지·확인 보류 기록은 이력으로 보존한다.
- **안내·조건 보존**: `tools.go:532`의 `coverage.note`는 엔드포인트 등 관련 자산을 포함한 대략적인 참고값, 현재 작업과 직접 관련된 작업의 scope·사실 앵커, 관련 scope의 읽기 전용, 상위 자산/대량 열거에 따른 낮은 비율, 테스트 완료로 판단하지 말라는 금지, `add_task_scope`·`list_untested_assets`의 용도 및 보통 후자를 호출하지 않고 작업을 진행한다는 지시를 유지했다. 535행은 자연어 값만 번역했으며 `Denominator==0` 조건·`pct=null`·키·출력 구조는 변경하지 않았다.
- **최종 범위·정적 보존**: 첫 묶음은 파일 시작부터 `graphOverviewData()` 끝인 1–552행이며 원본/현재 모두 2054줄이다. 553행부터 끝까지는 원본과 동일하다. 이번 추가 소스 변경은 532·535행의 두 문자열뿐이며 주석·허용 문자열 밖 코드, 영어·식별자·포맷 지정자·순서·들여쓰기·개행·이스케이프와 직전 번역을 보존했다.
- **테스트·잔여**: 직전 `tools_nil_store_test.go:47`의 "작업 컨텍스트" 기대값을 그대로 유지했으며 대응 소스의 오류와 정적으로 일치한다. 첫 묶음에서 중국어·전각 문장부호 잔여, 미분류 번역 누락, 추가 용어·계약 보류는 발견하지 못했다. DB 입력·외부 오류는 변환하지 않았으며 범위 밖 나머지 함수는 이번 완료 대상이 아니다.
- **단위 상태·실행 검증**: U5 전체는 나머지 `tools.go`가 남아 완료가 아니다. 기존 U1·U2 보류 및 팀원·U6 이력은 유지한다. 설치·빌드·타입 검사·lint·테스트·템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 검토를 최종 동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 관련 작업·읽기 도구  [부분 번역·정적 검토 / 용어 보류 / 실행 미검증]
- **실제 범위**: `inheritedMap()` 시작부터 `nodeDetail()` 끝인 `agent/tools.go` 554–1016행. 다음 `// --- planner write tools ---` 구분선은 1018행이다. 파일 원본/현재 모두 2054줄이며 1–553행과 1017행 이후는 작업 전과 동일하다.
- **번역·소비 경로**: 관련 작업 cold digest 안내·출처 주석, `list_findings`·`list_facts`·`node_detail`의 번역 가능한 도구/스키마 설명과 자체 오류를 번역했다. 도구 등록·planner/main/worker 도구 모음과 서버의 `ListFindingsTool`·`NodeDetailTool` 위임 경로, planner/worker JSON 컨텍스트 소비를 읽기 전용으로 확인했다. DB의 요약·상세·증거 및 외부 오류는 변환하지 않았다.
- **용어 보류**: 940행 `list_facts` 설명과 943행 `before` 설명은 미등록 `分页`·`游标` 때문에 문자열 전체를 원문 유지했다. 후보는 페이지네이션(문장에서는 페이지 단위 조회)·커서이며 원문 전체·문맥·근거는 `/tmp/artex-u5-tools-related-read-holds.md`에 기록했다. 이번 범위에 추가 계약 보류·후속 프롬프트 보류는 발견하지 못했다. 보류를 DNT나 완료로 분류하지 않는다.
- **관련 테스트·계약 검토**: `tools_overview_test.go`의 예산·UTF-8·문자 수, `blackboard_inheritance_test.go`의 직접 상속·출처·facts/findings/node detail 구조, `tools_nil_store_test.go`의 nil store 거부와 기존 "작업 컨텍스트" 기대값 및 DB 페이지 테스트를 확인했다. 전체/부분 문자열·접두사·정규식 의존 검색에서 이번 번역 문자열에 직접 연결된 테스트 기대값 변경은 발견하지 못했다. 테스트 파일은 변경하지 않았다.
- **정적 보존·잔여**: 주석·허용 자연어 문자열 밖 코드, 직접 관련 작업만 읽는 상속 깊이·읽기 전용·출처·예산·UTF-8·정렬·필터·페이지네이션·조회 권한·반환 구조, 키·enum·도구명·영어·숫자·포맷 지정자 및 순서·개행·이스케이프·들여쓰기를 보존했다. 중국어·전각 문장부호 잔여는 위 두 용어 보류 문자열이며 미분류 번역 누락과 중국어 DNT 잔여는 발견하지 못했다. 범위 밖 서버 설명·DB 입력·외부 오류는 이번 완료 대상이 아니다.
- **단위 상태·실행 검증**: `tools.go` 나머지와 두 용어 보류가 남아 U5 전체 완료가 아니다. 첫 묶음·팀원 기록·U1·U2 보류·U6 상태를 보존한다. 설치·빌드·타입 검사·lint·테스트·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅은 **사용자 요청으로 미실행**했다.

### U5 · tools.go 두 번째 묶음 보류 해소  [번역·정적 검토 완료 / 나머지 함수 미작업 / 실행 미검증]
- **확정 용어·보류 해소**: `分页`=페이지네이션, `游标`=커서를 중복·충돌 없이 용어집에 등록했다. 커서는 이번 `before`/`next_before` 페이지 조회 문맥으로 제한하며 파라미터명·JSON 키·실제 값은 DNT로 보존한다. 직전 940·943행의 보류 설명 두 문자열을 번역했다. 이전 부분 번역·보류 기록은 이력으로 보존한다.
- **설명 조건 대조**: 현재 작업과 직접 관련된 작업의 사실/결론 최신순 조회, 긴 요약의 잘림과 `node_detail(id)` 전체 조회, 모든 파라미터의 선택 사항, `limit` 기본 20·최대 100, 이전 응답의 `next_before` 커서로 더 오래된 페이지 조회, `before` 생략/0의 최신 페이지 및 그 값보다 작은 id만 반환, `q` 요약 키워드 필터, 필터 적용 후 전체 개수인 `total`, `has_more=true`의 계속 조회, `source_task_id/inherited=true`와 읽기 전용, `list_findings` 취약점 조회를 유지했다.
- **최종 정적 검토**: 두 번째 묶음은 554–1016행이며 원본/현재 모두 2054줄이다. 이번 추가 소스 변경은 940·943행의 두 문자열뿐이고 나머지 직전 번역과 테스트를 보존했다. 누적 변경에서도 1–553행과 1017행 이후는 원본과 동일하다. 주석·허용 문자열 밖 코드·키·영어·숫자·포맷·들여쓰기·개행·이스케이프는 보존했다. 두 번째 묶음의 번역 누락·중국어 DNT 잔여·추가 용어/계약/후속 단계 보류는 발견하지 못했다.
- **단위 상태·실행 검증**: 이번 두 문자열에 연결된 테스트 기대값 변경 필요는 발견하지 못했다. U5는 `tools.go` 나머지가 남아 전체 완료가 아니다. 기존 팀원 기록·U1·U2 보류와 U6 상태를 유지한다. 설치·빌드·타입 검사·lint·테스트·템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다.

### U5 · tools.go 세 번째 묶음(planner write tools)  [번역·정적 검토 완료 / 실행 미검증]
- **실제 범위**: `agent/tools.go`의 `// --- planner write tools ---`(1018행)부터 `goalMet()` 닫는 중괄호(1215행)까지. 줄 번호보다 함수 경계를 우선해, `intentItem`(1021) · `addOneIntent()`(1032) · `addIntent()`(1091) · `listGoals()`(1148) · `proveGoal()`(1157) · `goalMet()`(1205)와 관련 주석을 대상으로 했다. 다음 `// --- worker write tools ---`(1217행)부터는 수정하지 않았다. 파일 원본/현재 모두 2054줄이며, 1–1017행과 1217행 이후(완료된 첫·두 번째 묶음 포함)는 작업 전과 동일하다.
- **번역·소비 경로**: 모델용 도구 설명(`add_intent`·`list_goals`·`prove_goal`·`goal_met`), 스키마 설명(`intents`·`summary`·`asset_ids`·`parent_ids`·`priority`·`goal_id`·`evidence_id`·`reason`), 자체 오류 메시지, 중국어 개발자 주석을 DNT 판단 후 번역했다. `add_intent`→`writeExpTool`→frontier·탐색 체인 연결, `prove_goal`→목표 met 표시 및 전체 목표 met 시 자동 완료 판정, `goal_met`→`t.GoalMet`/`t.Reason` 설정을 거쳐 `planner.Plan`(agent/planner.go:355,472)이 `(met, reason, err)`로 반환하고 `server/engine.go:880`에서 소비함을 읽기 전용으로 확인했다. DB에서 받은 자산·노드 데이터와 외부 오류(`h.Describe()`, `%w`)는 변환하지 않았다.
- **의미 보존**: 의도 생성 근거(이미 확인된 fact/finding에만 앵커)와 허용 부모 종류, 금지되는 부모 관계(의도/목표/힌트 금지)와 origin fact 연결 조건(최상위 방향은 parent_ids 비우면 자동 연결), `parent_ids`·`asset_ids`·일괄(`intents`)/단건 입력과 기본값(priority 0→5, 반환 ids 길이·순서) 처리, prove_goal의 증거·참조 조건(본 작업/직접 관련 작업의 사실·취약점만, 관련 작업 목표는 읽기 전용), 개별 목표 달성(prove_goal)과 전체 작업 종료(goal_met)의 구분, 두 도구의 역할·제약·지시 강도(`반드시`/`~만`/`~하지 마세요`/`~해서는 안 됩니다`)를 보존했다.
- **계약·DNT 보존**: 도구명·파라미터명·JSON/스키마 키·enum·관계 종류·상태값(`open`·`running`·`done`·`met`·`exhausted`·`goalless`)·기존 영어·실제 값·숫자·경로·`frontier`·`origin fact`·`GoalMet`·`resumeTask`·`ToolSet`·`nil`·`no-op`를 보존했다. 영어 성공 응답(`intent created: %d`·`goal %d marked met`·`acknowledged: goal marked met`)은 번역하지 않았다. 검증·DB 쓰기(`AddIntent`·`Link`·`SetNodeState`)·계보 연결·상태 전이·콜백(`resumeTask`)·완료 판정·반환 구조·조건 분기는 변경하지 않았다.
- **형식 보존**: 포맷 지정자 집합·순서·인수(범위 내 %d 7·%q 1·%s 3·%w 1, HEAD와 동일), `\n` 2개, 마크다운 `**` 2쌍, Go 구조체 태그 백틱 16개를 HEAD와 동일하게 유지했다. 숫자(0-10·5·20·100)와 변수 뒤 조사를 확인했다. `parent_id` 오류 두 곳(1041·1044)은 숫자 `%d` 뒤에 조사를 붙이지 않도록 `parent_id %d: …` 콜론 형식으로 정리했고(%d·%q의 종류·개수·순서와 오류 의미 보존), 변수 `%s` 뒤에는 고정 명사 `의도`를 둔다. 산문 전각 문장부호(`，：；（）、「」【】""`)만 반각·한글·작은따옴표·`[]`로 정리했고, 코드 구분자는 변경하지 않았다.
- **용어 보류 해소(兜底)**: 직전 보류였던 `agent/tools.go:1030` 주석의 `兜底连 origin fact`를 사용자 확정에 따라 "기본 연결" 문맥으로 등록했다. 용어집 1장(그래프 노드/개념)에 `兜底(连 origin fact)=기본 연결` 항목을 추가하되, `parent_ids`를 비워 둔 최상위 의도를 origin fact에 기본 연결하는 이번 문맥에 한정하고 기존 다른 문맥의 兜底 번역(mock=기본 처리·인터셉트=보완 판정·LLM=대체 처리·U6=최종 보호 조치)은 보존했다. 1030행 문장은 "최상위의 완전히 새로운 방향은 parent_ids를 비워 두며, 기본적으로 origin fact에 연결한다."로 번역했고 기존 주석 줄 배치는 유지했다. 이번 묶음에 남은 용어·계약·후속 단계 보류는 없다.
- **관련 테스트 확인**: `agent/blackboard_inheritance_test.go`(add_intent 등 사용)와 `server/engine_cancelcause_test.go`(`AbortGoalMet` 코드 `goal_met` 검사, U6 영역)를 읽기 전용으로 확인했다. 전체/부분 문자열·접두사·정규식 의존 검색 결과, 이번 번역 문자열을 직접 비교하는 테스트 기대값은 발견하지 못해 **테스트 기대값 수정이 필요하지 않다**. 테스트 파일은 변경하지 않았다. 이는 테스트 실행 결과가 아니라 정적 대조이다.
- **잔여 분류**: 兜底 보류 해소 후 범위(1018–1216) 내 번역 대상 중국어·전각 문장부호 잔여는 없다. 번역 대상 누락·미분류 중국어 DNT 잔여, 추가 계약 보류·후속 단계 보류는 발견하지 못했다.
- **단위 상태·실행 검증**: `tools.go`의 worker write tools 이후 함수가 남아 U5 전체 완료가 아니다. 기존 팀원 기록·U1·U2 보류·U6 상태와 완료된 첫·두 번째 묶음을 유지한다. 설치·빌드·타입 검사·lint·테스트·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 네 번째 묶음(worker write tools: 취약점·사실 기록)  [번역·정적 검토 완료 / 테스트 기대값 조정 완료 / 실행 미검증]
- **실제 범위**: `agent/tools.go`의 `// --- worker write tools ---`(1217행)부터 `recordFact()` 닫는 중괄호(1405행)까지. 함수·타입 경계 우선: `addFinding()`(1219) · `factItem`(1306) · `recordOneFact()`(1317) · `recordFact()`(1353)와 관련 주석. 다음 `type hintItem struct`(1408행)부터는 수정하지 않았다. 파일 원본/현재 모두 2054줄이며, 1–1216행과 1406행 이후(완료된 앞 묶음 포함)는 작업 전과 동일하다.
- **번역·소비 경로**: `report_finding`(addFinding)·`record_fact`(recordFact)의 모델용 도구 설명·스키마 설명, 자체 오류, 중국어 개발자 주석을 DNT 판단 후 번역했다. `report_finding`→`findingRecorder.Record`/`ts.RecordFinding`→`notifyFinding`/트리거, 결과 JSON(`evidence_status`/`evidence_note`)과 첫 줄 `finding recorded: %d`의 노드 ID 계약, `record_fact`→`recordOneFact`→`AddNode(KindFact)`·`Link(intent, RelYields)` 경로를 읽기 전용으로 확인했다. DB에서 받은 데이터와 외부 오류(`err.Error()`)는 변환하지 않았다.
- **의미·지시 강도 보존**: report_finding은 확인된 취약점과 검증 가능한 증거를 등록, `finding_id`와 `finding_node_id`는 서로 다른 ID(첫 줄은 노드 ID), HTTP 증거는 실제 요청·응답을 한 건씩 확인한 뒤 재현 순서로 참조, ID 추측·도메인/시간 기반 연관 추정·패킷 보충만을 위한 반복 탐지 금지(`~해서는 안 됩니다`), 트래픽 기록이 없어도(TCP/비 HTTP/미수집) 보고 가능하며 `[]`나 생략 허용, 기존 취약점의 증거는 `bind_finding_traffic`으로 보완하고 취약점 중복 생성 금지를 보존했다. record_fact는 일반 탐색 사실·결론이며 취약점·자산 등록과 구분, 한 번의 탐색 관찰은 가능한 한 하나의 사실로 통합, `facts` 일괄 입력과 기본 `intent_id` 처리, 실제 관찰·추론 및 부정 결론의 증거 조건, `observed`(직접 봄)/`inferred`(추론) 판정 조건과 방향 포기 판단이 planner(规划者)의 역할이라는 점을 보존했다.
- **계약·DNT 보존**: `finding recorded: %d\n%s` 첫 줄 노드 ID 계약(reporter 트리거 연결)과 `fact recorded: %d` 등 영어 성공 응답을 원문 보존했다. `evidence_status`의 `bound`/`not_bound`, `confidence`의 `observed`/`inferred`, `traffic_refs`의 role enum `baseline`/`proof`/`verification`/`supporting`, 모든 키·도구명(`report_finding`·`record_fact`·`add_task_hint`·`bind_finding_traffic`·`traffic_search`)·파라미터명·실제 값·`Agent`를 보존했다. 기록·검증·계보 연결·증거 바인딩·알림(`notifyFinding`)·콜백·조건 분기·반환 구조는 변경하지 않았다.
- **형식 보존**: 포맷 지정자 집합·순서·인수(범위 내 %d·%s·%d, HEAD와 동일), `\n` 리터럴 7개, 마크다운 `**` 2쌍, 이스케이프 따옴표 `\"` 2개(1360행 부정 결론 설명의 `\"관찰 + 시험적 해석\"`), Go 구조체 태그 백틱 26개를 HEAD와 동일하게 유지했다. 숫자(1.25·3·200 등 예시값)·들여쓰기·개행·bullet(`  · `)을 보존했다. 포맷 지정자 뒤 조사는 콜론·고정 명사로 처리했다. 산문 전각 문장부호만 반각·한글·작은따옴표·`[]`로 정리했고 코드 구분자·이스케이프는 유지했다.
- **관련 테스트 — 기대값 조정 완료(1곳)**: 사용자 허용에 따라 `server/finding_workflow_test.go:178`의 `strings.Contains(offRecord.EvidenceNote, "已关闭")` 기대 문자열만 `"비활성화"`로 변경했다. 대응 소스는 이번에 번역한 `agent/tools.go:1292`(binding 비활성화 `not_bound` note: "…자동 트래픽 바인딩이 **비활성화**되어 있어…")이며, note에 `비활성화`가 포함되어 `Contains(…, "비활성화")`가 참임을 텍스트로 대조했다. 테스트의 다른 조건(`FindingID`·`Bindings` 길이·`EvidenceStatus != "not_bound"`)·구조·입력·다른 기대값은 보존했다. 같은 파일 `:134`의 `独立漏洞记录 ID` 검사는 `server/finding_workflow.go:127`(범위 밖)의 서버 오류를 읽으므로 변경하지 않았다. `agent/finding_recorder_test.go:36,95`의 `finding recorded:` 기대는 영어 성공 응답을 보존해 영향이 없다. 상세는 `/tmp/artex-u5-tools-worker-write-holds.md`에 기록했다. 테스트 실행 결과가 아닌 정적 대조이다.
- **용어·계약 보류**: 없다. 미등록 용어·불명확한 계약으로 보류한 문장은 없다.
- **잔여 분류**: 범위(1217–1405) 내 번역 대상 중국어·전각 문장부호 잔여는 없다. 1245행 영어 주석이 인용하던 중국어 note 문구는 번역된 note와 일치하도록 갱신했다(영어 주석 구조 유지). 번역 누락·미분류 중국어 DNT 잔여, 추가 용어·후속 단계 보류는 없다.
- **단위 상태·실행 검증**: `tools.go`의 1406행 이후 함수(hint 등)가 남아 U5 전체 완료가 아니다. 기존 팀원 기록·U1·U2 보류·U6 상태와 완료된 앞 묶음을 유지한다. 설치·빌드·타입 검사·lint·테스트·`go test`·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 다섯 번째 묶음(hints/goals/constraints write tools)  [번역·정적 검토 완료 / 실행 미검증]
- **실제 범위**: `agent/tools.go`의 `type hintItem struct`(1408행)부터 `addHint()` 닫는 중괄호(1673행)까지. 함수·타입 경계 우선: `hintItem`(1408)·`addOneHint()`(1415)·`goalItem`(1441)·`addOneGoal()`(1449)·`setGoals()`(1474)·`constraintItem`(1542)·`addOneConstraint()`(1549)·`setConstraints()`(1567)·`addHint()`(1616)와 관련 주석. 다음 `killWorkTool` 설명 주석(1675행)부터는 수정하지 않았다. 파일 원본/현재 모두 2054줄이며, 1–1407행과 1674행 이후(완료된 앞 묶음 포함)는 작업 전과 동일하다.
- **번역·소비 경로**: `set_goals`·`set_constraints`·`add_hint`의 모델용 도구 설명·스키마 설명, 자체 오류, 중국어 개발자 주석을 DNT 판단 후 번역했다. `addOneHint`→`AddNode(KindHint, …, "active", "human", anchors)`, `addOneGoal`→`AddNode(KindGoal, …, "open", origin)`·`Link(origin fact, RelSpawns, goal)`, `addOneConstraint`→`AddConstraint(kind, text, worker)`, 일괄 기록 후 `notifyGoal`/`notifyHint`(없으면 `notify`)와 `resumeTask` 경로를 읽기 전용으로 확인했다. DB 데이터와 외부 오류(`err.Error()`)는 변환하지 않았다.
- **의미·지시 강도 보존**: 힌트의 자산 앵커와 `traffic_refs` 처리, 자동 트래픽 바인딩이 꺼졌을 때 `traffic_refs` 포함 힌트를 저장하지 않는 조건(1417), 힌트·목표 일괄 기록 후 planner 알림을 한 번만 호출(전체 묶음 1회), `notifyHint`/`notifyGoal` 우선·없으면 일반 `notify` 대체 경로, `resumeTask` 연결 조건(mainagent만)과 작업 재개(완료/일시 중지→running), 목표는 검증 가능한 최종 결과이며 공격 단계·정찰 동작과 구분, `set_goals`는 목표 추가만 담당하고 달성 판정과 구분, 동작 제약 조건은 allow/deny 규정이며 목표·공격 단계와 구분, 명시된 제약만 등록하고 임의 생성 금지(`지어내지 마세요`), type 미지정 시 `deny` 처리 및 불확실하면 `deny` 사용, 일괄/단건 입력·ids 길이·순서·실패 항목 `id=0`과 `errors`, round-0에서 planner 미시작 조건을 보존했다.
- **계약·DNT 보존**: 영어 성공 응답 `goal added: %d`·`constraint added: %d`·`hint added: %d`(server 테스트가 `Sscanf`로 파싱)와 enum·상태·출처 값 `active`/`human`/`open`/`system`/`goals`/`allow`/`deny`를 원문 보존했다. 도구명·파라미터명·JSON/스키마 키·`origin fact`·`rel spawns`·`task_constraints`·콜백명(`notifyGoal`·`notifyHint`·`notify`·`resumeTask`)·`ExplorationStore`를 보존했다. 검증·DB 쓰기·계보 연결·조건 분기·반환 구조는 변경하지 않았다. `提示`=힌트, `系统提示`는 문맥상 시스템 프롬프트로 구분 번역했다(상세는 보고서).
- **형식 보존**: 포맷 지정자 집합·순서·인수(범위 내 %d·%d·%d, HEAD와 동일), `\n` 리터럴 6개, Go 구조체 태그 백틱 20개를 HEAD와 동일하게 유지했다. 숫자·들여쓰기·개행·이스케이프를 보존했다. 포맷 지정자 뒤 조사는 콜론·고정 명사로 처리했다(해당 묶음의 오류는 숫자 `%d` 직접 인접 조사 없음). 산문 전각 문장부호(`，：；、。「」『』【】""`)만 반각·한글·작은따옴표·`[]`로 정리했고, `callback이 nil`인 설명을 ToolSet 전체가 nil인 것으로 바꾸지 않았다(분해기/worker에서 특정 콜백만 nil).
- **주석 인용 알림 문구(생성부 미수정)**: 1512-1513·1650-1651 주석의 인용 알림 문구 `人新增了 N 个目标`/`人新增了 N 条战略提示`를, 같은 콜백의 struct 주석 `tools.go:97,104`의 기존 한국어 표기에 맞춰 `사용자가 목표 N개 추가`/`사용자가 힌트 N개 추가`로 번역했다. 실제 생성부(`server/manager.go:1782,1798`·`server/orchestration.go:170`·`server/goals_api.go:71`·`agent/mainagent.go:122-123`)는 아직 중국어이며 **범위 밖이라 수정하지 않았다**(U11 등 소유 단위 번역 시 위 한국어와 맞춰야 함).
- **관련 테스트**: `agent/finding_workflow_test.go`·`server/finding_workflow_test.go`를 읽기 전용으로 확인했다. `server/finding_workflow_test.go:147`은 영어 `hint added: %d`를 `Sscanf`로 파싱하며 이 응답을 보존했으므로 영향이 없다. 이 묶음의 번역 중국어 문자열을 직접 비교하는 테스트 기대값은 발견하지 못했다. **테스트 기대값 조정 불필요**, 테스트 파일은 수정하지 않았다. 정적 대조이며 실행 결과가 아니다.
- **후속 표현 정리**: (1) `战略提示`를 용어집대로 "힌트"로 통일(1617행 `add_hint` 설명의 "전략 힌트"→"힌트"). (2) 이번 범위의 `操作约束`를 "동작 제약 조건"으로 통일(1547·1564·1569·1575행). 영어·키·식별자는 보존. (3) 함수명 뒤 조사를 고정 명사 "함수"로 정리: `addOneHint 함수는`(1414)·`addHint 함수가`(1436)·`addOneGoal 함수는`(1446)·`setGoals 함수가`(1448)·`setGoals 함수는`(1472)·`addOneConstraint 함수는`(1547)·`setConstraints 함수는`(1564). (4) `set_goals`/`set_constraints` 오류 안내를 "`set_goals` 도구를 사용할 수 없습니다" 형태로 정리(1486·1580행, `도구` 고정 명사 삽입).
- **용어 보류 해소(被动侦察)**: 직전 보류였던 `被动侦察`를 사용자 확정에 따라 "패시브 정찰"(passive reconnaissance)로 등록했다. 용어집 2장(mock UI·데모 용어)의 `侦察`=정찰 항목 바로 뒤에 `被动侦察=패시브 정찰` 행을 추가하되, 대상에 능동적으로 요청을 보내지 않는 정찰 **방식** 문맥에 한정하고 자동/수동(manual) 실행 방식과 구분하며 "수동 정찰"로 쓰지 않는다고 명시했다. `agent/tools.go:1570`의 `set_constraints` 설명 예시 문장을 번역했다: `"제약 조건='어떤 동작을 할 수 있는지/할 수 없는지'에 대한 규정(예: '현재 포트만 테스트하고 다른 포트는 스캔하지 않음' '운영 DB에 쓰기 작업 금지' '패시브 정찰만 허용')이며, 목표도 아니고 공격 단계도 아닙니다."` '패시브 정찰만 허용' 예시와 제약 조건이 목표·공격 단계와 다르다는 의미를 보존했다. 이번 묶음에 남은 용어·계약·후속 단계 보류는 없다.
- **잔여 분류**: 被动侦察 보류 해소 후 범위(1408–1673) 내 번역 대상 중국어·전각 문장부호 잔여는 없다. 번역 누락·미분류 중국어 DNT 잔여, 추가 용어·계약·후속 단계 보류는 없다.
- **단위 상태·실행 검증**: `tools.go`의 1674행 이후 함수(kill/steer work 등)가 남아 U5 전체 완료가 아니다. 기존 팀원 기록·U1·U2 보류·U6 상태와 완료된 앞 묶음을 유지한다. 설치·빌드·타입 검사·lint·테스트·`go test`·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 여섯 번째 묶음(work control: kill/steer/get_worker_output)  [번역·정적 검토 완료 / 실행 미검증]
- **실제 범위**: `agent/tools.go`의 `killWorkTool` 설명 주석(1675행)부터 `getWorkerOutput()` 닫는 중괄호(1798행)까지. 함수 경계 우선: `killWorkTool()`(1676)·`steerWorkTool()`(1706)·`getWorkerOutput()`(1740)와 각 설명 주석. 다음 `traceSteps` 설명 주석·함수(1800행~)부터는 수정하지 않았다. 파일 원본/현재 모두 2054줄이며, 1–1674행과 1799행 이후(완료된 앞 묶음 포함)는 작업 전과 동일하다.
- **번역·소비 경로**: `kill_work`·`steer_work`·`get_worker_output`의 모델용 도구 설명·스키마 설명, 자체 오류, 성공 반환 안내문, 중국어 주석과 영어 주석 속 중국어 조각을 DNT 판단 후 번역했다. `killWork`/`steerWork` 콜백, `GetNode`/`GetNodeWithSources`로 의도 검증, `ActivityListWithSources`로 활동 조회 후 `result`(result)·`text`(text) 중 선택, `inheritedMap`으로 상속 표시를 붙이는 경로를 읽기 전용으로 확인했다. 실제 worker 출력·DB `Summary`/`Detail`·외부 오류(`err.Error()`)는 변환하지 않았다.
- **의미·지시 강도 보존**: kill_work는 현재 작업의 실행 중인 의도를 중지하고 `stopped`로 표시하며 자동 재할당을 막고, 먼저 `get_worker_output`으로 확인하라는 지시를 유지. steer_work는 실행을 중단하거나 기존 진행 내용을 버리지 않고 다음 동작 전에 방향 조정 지시를 전달하며, 의도 내부 조정과 전체 방향 변경을 구분하고 전체 방향이 잘못됐으면 kill_work 후 새 의도를 생성하라는 지시를 유지. get_worker_output은 현재 작업 또는 직접 관련된 작업만 조회하고 관련 작업은 읽기 전용이며 `source_task_id`/`inherited` 상속·출처 표시를 유지, 정상 종료 요약과 중지(stopped)/오류 시 마지막 출력을 구분하는 의미를 보존했다.
- **계약·DNT 보존**: 도구명·`intent_id`·`message`·`final_text`·`summary`·`is_error`·`source_task_id`·`inherited`·`stopped`·`result`·`text` 등 키·enum·계약값·영어·숫자를 보존했다. 조건 분기·검증·콜백·DB 조회·결과 선택(`result`/`text` 우선순위)·반환 구조는 변경하지 않았다. 포맷 지정자 `%d`(범위 내 2개, kill/steer 성공 반환) 개수·순서·인수를 보존했고, 숫자 `%d` 뒤에는 받침에 영향받지 않는 `의`를 사용했다. 들여쓰기·개행·이스케이프·줄 수를 유지했다. 실제 worker 출력·DB Summary/Detail·외부 오류는 번역하지 않았다.
- **출력 없음 안내문(두 경로)**: 일반 반환(1783)과 상속 JSON 반환(1780, `final_text` 값) 모두 `(이 work에는 아직 아무 출력도 없습니다)`로 번역했다. 두 문구와 `final_text` 값을 비교·파싱하는 프로덕션 소비부는 발견하지 못해(계약 의존 없음) 보류하지 않았다.
- **관련 테스트**: `agent/blackboard_inheritance_test.go:356`은 `final_text` **키**와 실제 출력이 있는 경로의 값(`shared trace full detail`, 번역 대상 아님)을 검사하므로 출력-없음 placeholder 번역과 무관하다. `server/engine_cancelcause_test.go:35`는 kill_work의 취소 범위 동작을 검사하며 `t.Fatal` 인자는 테스트 자체 영어 메시지다. 이 묶음의 번역 문자열을 직접 비교하는 기대값은 발견하지 못했다. **테스트 기대값 조정 불필요**, 테스트 파일은 수정하지 않았다. 정적 대조이며 실행 결과가 아니다.
- **용어·계약 보류**: 없다.
- **잔여 분류**: 범위(1675–1798) 내 번역 대상 중국어·전각 문장부호 잔여는 없다. 영어 주석 속 중국어 조각(1739 `截至中止时的`→`중지 시점까지의`, 1705 `"停做 X、聚焦 Y"`→`"X는 그만하고 Y에 집중"`)도 번역했다. 번역 누락·미분류 중국어 DNT 잔여, 추가 용어·계약·후속 단계 보류는 없다. 상세는 `/tmp/artex-u5-tools-work-control-holds.md`.
- **단위 상태·실행 검증**: `tools.go`의 `traceSteps` 이후(1800행~) 함수가 남아 U5 전체 완료가 아니다. 기존 팀원 기록·U1·U2 보류·U6 상태와 완료된 앞 묶음을 유지한다. 설치·빌드·타입 검사·lint·테스트·`go test`·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · tools.go 일곱 번째 묶음(trace 조회·검색·planner 도구 모음)  [번역·정적 검토 완료 / 실행 미검증]
- **실제 범위**: `agent/tools.go`의 `traceSteps` 설명 주석(1800행)부터 파일 끝(2054행)까지. 함수 경계 우선: `traceSteps()`(주석 1800/함수 1803)·`getWorkerTrace()`(1821)·`searchAllWorkerTraces()`(1930)·`listWorkerTraces()`(1979)·`PlannerTools()`(2033)와 관련 주석. 파일 끝의 도구 모음 반환 함수는 이 범위에서 `PlannerTools()` 하나뿐이다. 파일 원본/현재 모두 2054줄이며, 1–1799행(완료된 앞 묶음)은 작업 전과 동일하다.
- **번역·소비 경로**: `get_worker_trace`·`search_all_worker_traces`·`list_worker_traces`의 모델용 도구 설명·스키마 설명, 자체 오류, `notice` 안내문, 중국어 개발자 주석을 DNT 판단 후 번역했다. `PlannerTools()` 꼬리의 도구 주석도 번역했다. `getWorkerTrace`→`GetNodeWithSources`로 의도 검증→③`ActivityByIDsWithSources`(step_ids 상세, 중복/무효 제거 후 5개 상한)·①/②`ActivityTraceWithSources`/`ActivityTraceSearchWithSources`(요약 스트림/검색), `searchAllWorkerTraces`→`ActivityTraceSearchAllWithSources`(ownerNode 제외), `listWorkerTraces`→`ListByKindWithSources(KindIntent)`로 실행된 의도만 필터, `traceSteps`→`firstLine(summary,100)` 요약 재절단, `inheritedMap`으로 상속 표시를 붙이는 경로를 읽기 전용으로 확인했다. 실제 trace의 `summary`/`detail`/`result_summary`·DB 데이터·외부 오류(`err.Error()`, `nodeErr.Error()`)는 변환하지 않았다.
- **의미·지시 강도 보존**: get_worker_trace의 세 용법(①intent_id만=요약 스트림 / ②+q=키워드 매칭 요약 검색 / ③+step_ids=상세 조회)과 요약 검색·전체 상세 조회의 차이, 사고(thinking) 단계 제외, 직접 관련 작업의 과거 trace·`source_task_id`/`inherited=true`·읽기 전용 조건, summary≤100자 길이 제한, 상한(step_ids 1회 최대 5개)·`returned_step_ids`와 `omitted_step_ids` 구분 및 `notice` 안내를 보존했다. search_all_worker_traces는 다른 work의 실행 과정에서 fact에 기록되지 않은 것을 찾되 자신의 의도 단계를 자동 제외하고 매칭 요약만 intent_id와 함께 반환, list_worker_traces는 worker가 탐색 그래프 없이 실행된 work(의도) 색인을 발견하는 용도이며 실행된 상태(running/done/exhausted/blocked/stopped)만 나열하고 open 제외, 작업 경계는 자신이 배정받은 의도라는 범위 제한을 보존했다. PlannerTools의 도구 구성·순서·노출 범위(read + 의도 생성 + 목표 판정 도구 모음)는 변경하지 않았다.
- **계약·DNT 보존**: 도구명(`get_worker_trace`·`search_all_worker_traces`·`list_worker_traces`·`get_worker_output`·`report_finding`·`list_companies`·`list_assets`·`add_company_scope`·`add_task_scope`·`list_untested_assets` 등)·파라미터/키(`intent_id`·`q`·`step_ids`·`limit`·`step_id`·`source_task_id`·`inherited`·`returned_step_ids`·`omitted_step_ids`·`notice`·`summary`·`detail`·`is_error`·`kind`·`tool`·`worker`·`works`·`hits`·`query`·`state`·`scope`·`company_id`)·상태 enum(`running`/`done`/`exhausted`/`blocked`/`stopped`/`open`)·영어·숫자(100·5·50·100·500)·`thinking`·`①②③` 용법 번호를 보존했다. 조회·검색·필터·정렬·절단·상한·반환·도구 등록 로직, 조건 분기, `maxStepIDs=5` 상수를 변경하지 않았다. `work 句柄`→`work 핸들`(앞 묶음 표기와 일치)로 번역했다.
- **형식 보존**: `notice`의 포맷 지정자 `%d %d %v %d %v`(범위 내 %d 3개·%v 2개, HEAD와 동일) 종류·개수·순서·인수(maxStepIDs, len(ids), ids, len(omitted), omitted)를 보존했다. 숫자 `%d` 뒤에는 콜론·고정 명사·받침 무관 조사로 처리했다(해당 묶음의 `%d` 직접 인접 조사 없음). 산문 전각 문장부호(`【】（）：；，。""`)만 반각·한글·작은따옴표·`[]`로 정리했고(`notice`의 전각 `，`·`。`는 계약 파서 구분자가 아닌 산문 구두점이므로 정리), `\n` 4개(getWorkerTrace 설명)·들여쓰기·개행·이스케이프·줄 수(2054)를 HEAD와 동일하게 유지했다. PlannerTools 주석의 전각 콜론 `：`→`: `, 전각 인용 `""`→`"…"`로 정리했다.
- **표현 통일**: searchAllWorkerTraces 설명(1933)과 해당 주석(1948)의 `你自己这条意图的步骤`/`调用者自身这条意图的步骤`를 모두 `호출자 자신의 의도에 속한 단계`로 통일했다.
- **자유 서술 문자열의 계약 의존 확인**: `notice`·자체 오류(`intent_id는 필수입니다`·`q는 필수입니다`·`intent_id가 이 작업 또는 직접 관련된 작업에 속하지 않습니다`)·도구 설명을 비교·파싱·접두사·정규식으로 의존하는 소비부를 `agent/`·`server/`·`web/`·`db/`에서 검색했으나(producer=tools.go 제외) 발견하지 못했다. `returned_step_ids`/`omitted_step_ids`/`notice`/`result_summary`는 결과 **키**이며, trace 결과의 실제 값은 DB·worker 출력이다. `server/orchestration.go`의 `get_task_worker_trace`·`search_task_worker_traces` 등은 **별개의 서버측 중국어 도구 설명(U11 범위)**이며 이 문자열의 소비부가 아니다.
- **관련 테스트 — 조정 불필요**: `agent/blackboard_inheritance_test.go`가 세 trace 도구를 호출한다. `:306` `result_summary == "shared trace marker"`·`:356` `final_text == "shared trace full detail"`는 번역 대상이 아닌 영어 fixture 값과 키를 검사하고, `:406/:417` `returned_step_ids`/`omitted_step_ids`는 키·배열을 검사하며, `:421` `notice`는 **비어 있지 않은지만** 검사(`notice == ""`)하고 `:431` 상한 이하에서 `notice` 부재만 검사한다. 번역한 `notice`는 비어 있지 않으므로 기대값 영향이 없다. 이 묶음의 번역 문자열을 직접 비교하는 기대값은 없어 **테스트 기대값 조정 불필요**, 테스트 파일은 수정하지 않았다. 정적 대조이며 실행 결과가 아니다.
- **용어 확정·보류 해소(2건)**: 직전 보류 2건을 사용자 권장안으로 확정해 용어집에 등록하고 주석을 번역했다. (1) `态势研判`=현황 분석·판단(situation analysis and assessment) — 1장 그래프/개념 표의 `兜底` 행 뒤에 등록, 이번 `PlannerTools()`의 `report_finding` 주석(계획 수립 중 탐색 현황 분석·판단) 문맥에 한정. `tools.go:2040`을 "계획 수립 중 현황 분석·판단 시 자신이 이미 취약점을 확증했다면 직접 등록할 수 있다(worker와 같은 도구)."로 번역했다. (2) `认领`=편입(claim) — 2장 `范围 / 备案` 행 뒤에 등록, 회사 자산 범위 규칙에 매칭된 자산을 해당 회사에 편입하는 문맥에 한정하고 작업·의도를 맡거나 재할당하는 문맥(`server/server.go`의 `重新认领` 등)은 제외. `tools.go:2047`을 "계획 수립 시 도메인/IP/CIDR/ICP/키워드를 특정 회사의 자산 범위에 포함할 수 있다(매칭된 자산을 자동으로 편입)."로 번역했다. 상세 근거는 `/tmp/artex-u5-tools-traces-review.md`.
- **다른 파일 표기 통일(해소 완료)**: `agent/tools_insert.go:380`(이전에 커밋된 U5 묶음)은 원문 `会自动认领命中的资产`의 `认领`을 "귀속"으로 번역했었다. 사용자 추가 허용에 따라 이 한 곳의 "귀속"만 "편입"으로 변경해 `认领`=편입 확정과 일치시켰다("매칭되는 자산을 자동으로 편입하며 …"). 같은 파일의 `归属` 대응 "귀속"(170 `归属企业`·214 주석·383/388 `归属依据`·590 `已归属资产数`)은 다른 단어이므로 변경하지 않았다.
- **잔여 분류**: 번역 후 `tools.go` **파일 전체**(1–2054)의 중국어·전각 문장부호 잔여는 **0건**이다. 앞 여섯 묶음과 이번 묶음으로 파일 전체가 커버됐고, 번역 누락·미분류 중국어 DNT 잔여는 발견하지 못했다. (위 tools_insert.go:380 일관성 항목 외) 앞 묶음이나 다른 U5 파일에서 추가 수정 사항은 발견하지 못했다.
- **단위 상태·실행 검증**: 이번 묶음은 번역·정적 검토를 마쳤고 남은 용어 보류는 없다. 설치·빌드·타입 검사·lint·테스트·`go test`·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**했다. 정적 텍스트 대조를 실행·동작 검증 완료로 표시하지 않는다.

### U5 · 전체 종합  [번역·정적 검토 완료 / 실행 미검증]
- **네 소스 종합**: `agent/toolcatalog.go`·`agent/tools_digest.go`·`agent/tools_insert.go`·`agent/tools.go`의 번역·정적 검토를 모두 마쳤다(각 묶음 기록은 위에 보존). 일곱 묶음으로 나눠 작업한 `tools.go`를 포함해 네 소스의 번역 대상 중국어는 모두 처리했고, 남은 용어·계약 보류는 없다.
- **정적 잔여(네 소스)**: `toolcatalog.go`·`tools_digest.go`·`tools.go`의 중국어 잔여 0건. `tools_insert.go`의 중국어 잔여는 `备案`(§5.7 값 계약) 6곳뿐으로, 모두 표시 표기 `ICP 등록(备案)` 안에 보존한 의도적 DNT이다(번역 누락 아님).
- **표기 통일(해소 완료)**: 위 `tools_insert.go:380`의 `认领`→"귀속"을 "편입"으로 변경해 `认领`=편입 확정과 일치시켰다(사용자 추가 허용). U5 네 소스의 `认领` 표기가 "편입"으로 통일됐으며 남은 후속 정리 항목은 없다.
- **실행 검증**: U5 네 소스 모두 설치·빌드·타입 검사·lint·테스트·`go test`·실제 템플릿 렌더·`gofmt`·앱·화면·스모크·Git 훅은 **사용자 요청으로 미실행**이다. 따라서 U5 상태는 **번역·정적 검토 완료 / 실행 미검증**이며 실행·동작 검증 완료가 아니다. 기존 이력·팀원 기록·U1/U2 보류·U6 상태는 보존한다.


### U7 · 입력 경계 첫 파일럿  [파일럿 번역·정적 검토 완료 / 후속 호환성·테스트 검토 보류 / 실행 미검증]
- **실제 범위**: `intercept/prompt.go:10–18`의 `JudgeContextBoundary` 상수만 번역했다. `EffectiveJudgePrompt`·`JudgeOutputContract`·`DefaultJudgePrompt`·판정 파서 및 상수 밖 내용은 작업 전과 동일하다. 원본 줄 수·들여쓰기·영어·필드명·값·개행·이스케이프를 보존했다.
- **번역·의미 보존**: 현재 호출만 검토, 실제 사용자 메시지만 배경으로 선택, Worker 배경·의도 요약·상위 배경 및 실행 이력 제외, 누락된 이력으로 일반 읽기 전용 동작을 거부하지 않음, 미실행 동작의 성공 주장 금지, 배경 잘림과 전체 파라미터의 구분, 판정 규칙 추가·덮어쓰기 금지 및 숨겨진 사고 과정 요구·도구 실행·대체 파라미터 반환 금지를 유지했다.
- **용어 보류 해소**: 첫 부분 번역에서는 `生产`·`归属`·`提示注入`가 용어집에 없어 13행 마지막 문장, 15행 첫 문장, 16행 두 문장을 원문 유지했다. 이후 사용자 확정에 따라 운영·귀속·프롬프트 인젝션을 적용 문맥과 함께 등록하고 네 문장을 번역했다. `认领`=편입과 귀속을 구분하고 식별자·키·값·계약 문자열은 DNT로 보존한다. “이번 스케줄링의 전체 입력”·“배경에 적힌 주장”을 적용했으며 전체 상수의 조사·띄어쓰기·조건·금지 강도를 다시 대조했다. 원문과 해소 이력은 `/tmp/artex-u7-boundary-review.md`에 보존했다.
- **생성부·소비부·관련 테스트**: `BuildReviewInput`→`Interceptor.Judge`→`EffectiveJudgePrompt`→서버 `reviewCompletion`/`streamCollectText`의 모델 입력 경로를 읽기 전용으로 확인했다. `intercept/review_context_test.go`의 경계 중복 방지 검사는 상수를 직접 참조하며 `server/intercept_review_test.go`·`agent/review_context_test.go`는 입력 구조·배경·이력 제외를 검사하므로 직접 기대값 조정 필요는 발견하지 못했다. 파서 테스트의 중국어 앵커·입력은 그대로 보존했다. `server/intercept_live_test.go:126–136`은 모델 출력의 중국어 동사(创建·新建·写入·写)를 검사하므로 한국어/혼합 프롬프트의 출력 언어에 따라 영향받을 가능성이 있으며 후속 계약·테스트 검토 보류로 기록한다. 이번에는 테스트 파일을 수정하지 않았고 실행 영향은 확인하지 않았다.
- **후속 호환성**: 빈 저장값은 런타임 상수를 사용하므로 기존 DB에도 이번 경계 번역이 반영되며 새 DB는 필수 조건이 아니다. 옛 중국어 경계 블록이 포함된 커스텀 프롬프트에는 전체 문자열 `Contains` 판정 때문에 새 한국어 블록이 중복 부착될 수 있다. 저장·로드·중복 판정 로직과 DB는 변경하지 않았고 실제 사용자 DB 상태는 확인하지 않았다.
- **잔여·상태**: 이번 상수의 번역 대상 중국어·전각 문장부호 잔여, 번역 누락 및 추가 용어 보류는 발견하지 못했다. `prompt.go`는 원본과 같은 205줄, 상수는 9줄이며 상수 밖 내용은 동일하다. 키·값·영어·경로는 DNT로 보존하며 상수 밖 중국어는 이번 범위 밖이다. 커스텀 프롬프트 중복 부착 가능성과 live 테스트의 중국어 동사 의존은 후속 검토 보류로 유지하고 실제 DB·모델 출력은 미확인이다. 파일럿 번역·정적 검토 완료는 U7 전체 완료가 아니며 기존 U1·U2 보류와 U5·U6 기록은 보존한다.
- **실행 검증**: 설치·빌드·타입 검사·lint·테스트·gofmt·실제 렌더·앱·스모크·Git 훅은 **사용자 요청으로 미실행**이다. 정적 텍스트 검토를 실행·동작 검증 완료로 표시하지 않는다.


### U7 · 판정 출력 계약 두 번째 묶음  [대상 번역·정적 검토 완료 / 후속 호환성·테스트 검토 보류 / 실행 미검증]
- **실제 범위**: `intercept/prompt.go:32–40`의 `JudgeOutputContract` 상수 전체의 자연어 안내만 번역했다. 상수는 원본과 같은 9줄, 파일은 205줄이며 줄 배치·들여쓰기·개행·이스케이프를 유지했다. `JudgeContextBoundary`·`EffectiveJudgePrompt`·`DefaultJudgePrompt`의 출력 예시·파서·테스트·저장/로드 로직·DB는 변경하지 않았다.
- **계약·의미 보존**: JSON 객체 하나, 첫/마지막 문자 `{`/`}`, 문자열 필드 `decision`·`comment` 정확히 두 개, 큰따옴표 사용, YAML·코드 블록·사고 과정·서문·설명 및 JSON 앞뒤 다른 문자 금지, `allow`·`ask`·`deny`와 `A1–A6`·`D1–D6`·`ASK`·`DEFAULT`를 보존했다. 사유 세 부분 모두 필수·비어 있지 않음·각 한 문장·간결함, 현재 호출만 설명하고 미실행 성공 주장 및 판정 근거 조작을 금지하는 조건을 유지했다.
- **앵커 정적 대조**: `实际操作：`·`；成功后的后果：`·`；命中规则：`는 형식 안내·정의 라벨과 원문 예시·파서의 `HasPrefix`/`strings.Cut`에 쓰인 표기를 전각 콜론·세미콜론까지 보존했다. 정의 라벨의 `成功后的后果：`·`命中规则：`도 원래대로 유지했다. 앵커 밖 자연어만 번역했으며 새 용어·계약 보류나 번역 누락은 발견하지 못했다. 중국어 잔여는 계약 앵커·정의 라벨인 DNT로 구분한다.
- **길이 지시와 파서 제한**: “전체 comment는 120자를 넘지 않아야 합니다”는 모델의 출력 길이 지시 번역이다. 파서의 `len(reason) > 2400`은 Go 문자열의 2400바이트 제한이며 별개다. 파서가 120자를 검사하도록 변경하지 않았고 모델의 실제 길이·형식 준수도 확인하지 않았다.
- **관련 테스트**: `intercept/prompt_test.go`의 앵커·JSON·길이·불완전 응답·코드펜스 검증 데이터와 `intercept/review_context_test.go:132–138`·`server/intercept_review_test.go:42–52`의 상수 참조를 읽기 전용으로 대조했다. 번역한 안내 산문을 직접 비교하는 기대값은 발견하지 못해 직접 기대값 수정은 필요하지 않다. 테스트 파일은 수정·실행하지 않았다. 코드펜스 출력 금지 지시는 유지하고 기존 파서의 코드펜스 허용 동작도 변경하지 않았다.
- **후속 검토·단위 상태**: 기존 live 테스트(`server/intercept_live_test.go:125–134`)의 중국어 동사 의존과 옛 중국어 경계가 포함된 커스텀 프롬프트의 중복 부착 가능성은 유지한다. 출력 계약도 전체 상수 `Contains`로 중복 판정하므로 옛 중국어 출력 계약이 저장된 커스텀 프롬프트에는 번역 블록이 추가 부착될 수 있음을 후속 호환성 항목으로 기록한다. 실제 DB·모델 출력은 미확인이며 U7 전체 완료가 아니다. 기존 이력·U1·U2·U5·U6·입력 경계 기록은 보존한다.
- **실행 검증·산출물**: 설치·빌드·타입 검사·lint·테스트·go test·gofmt·렌더·앱·스모크·Git 훅은 **사용자 요청으로 미실행**이다. 정적 비교를 실행 검증 통과로 표시하지 않는다. 전체 diff는 `/tmp/artex-u7-output-contract.diff`, 상세 보고는 `/tmp/artex-u7-output-contract-review.md`에 저장했다.

### U7 · 기본 판정 정책 전체 상수  [부분 번역·정적 검토 / 용어 검토 보류 / 후속 호환성·테스트 검토 보류 / 실행 미검증]
- **실제 범위·시작 상태**: `intercept/prompt.go:50–128`의 `DefaultJudgePrompt` 전체를 검토했다. 기존 50–75행 부분 작업이나 대상 파일의 미커밋 변경은 없었다. 역할·두 판단 축·공통 제약·D1–D6·A1–A6·ASK·판정 순서·비교 예시·JSON 출력 안내에서 번역 가능한 자연어를 처리했다. 상수는 원본/현재 79줄, 파일은 205줄이다. 줄 배치·빈 줄·들여쓰기·개행·이스케이프·마크다운을 유지했다.
- **정책·계약 보존**: 직접 효과만 검토, 민감한 정보의 읽기 자체는 차단 사유가 아님, 귀속·영향 범위의 확인 조건, 운영 자산과 테스트 산출물의 구분, 필수 제약의 우선순위, 규칙 순서 및 매칭 시 중단, 정보 부족과 손상 조건을 함께 충족할 때만 ASK, 나머지 기본 ALLOW를 유지했다. 출력 언어 지시나 정책을 추가하지 않았다. 명령·경로·실제 예시 값·판정값·규칙 번호·숫자·JSON 구조·키를 보존했다. 예시의 `实际操作：`·`；成功后的后果：`·`；命中规则：`는 각각 4곳을 전각 부호까지 유지했으며 파서와 정적으로 일치한다.
- **용어 보류(15곳)**: 60–62·68(첫 문장)·71(첫 문장)·80·83·87·90–92·102·106–107·112행의 해당 문장을 원문 유지했다. 가역성/가역/비가역, 복합 명령, 미들웨어·시작 항목·예약 작업, 부하 테스트, SQLi의 불리언·시간 기반 블라인드 인젝션, 경로 순회·파일 포함, 크롤링·패킷 캡처·핑거프린트 식별, 거점 확보 이후의 용례는 최신 용어집에서 해당 전문 용어의 등록을 찾지 못했다. 기존 `立足点`=거점만으로 `落脚`을 임의 확정하지 않았다. 각 문장의 원문 전체·위치·문맥·권장 번역은 `/tmp/artex-u7-default-policy-review.md`에 기록했다. 보류를 DNT나 번역 완료로 처리하지 않는다.
- **잔여 분류**: 미분류 번역 누락은 발견하지 못했다. 중국어 DNT는 계약 앵커와 115행 명령 예시의 `WHERE 全表` 안에 있는 원문 표기다(실제 예시 값 보존). 15곳의 원문 문장은 용어 보류이며 계약 앵커와 구분한다. `密码喷洒`는 확정된 `Password Spraying`을 적용했고 기존 영어 토큰의 순서·내용은 유지했다. `JudgeContextBoundary`·`JudgeOutputContract`·`EffectiveJudgePrompt`·파서·테스트·저장/로드·DB 및 상수 밖 내용은 변경하지 않았다.
- **생성부·소비부·관련 테스트**: `intercept.go:346–349`의 빈 저장값→기본 상수 선택, `:400–405`의 기본 상수와 같은 설정은 빈 값으로 저장하는 비교, `:445–446`의 입력 생성·공통 블록 조립, `server/intercept_review_test.go:42–50` 및 `intercept/review_context_test.go:132–138`의 상수 참조, `intercept/prompt_test.go`의 파서 검증을 읽기 전용으로 대조했다. 번역된 기본 정책 산문을 직접 비교하는 테스트 기대값은 발견하지 못해 이번 직접 기대값 조정은 필요하지 않다. 단 `server/intercept_live_test.go:124–134`의 중국어 동사 검사 의존은 기존 후속 검토 보류로 유지하며 테스트를 수정·실행하지 않았다.
- **호환성·단위 상태**: 빈 저장값은 변경된 런타임 기본 상수를 사용하므로 새 DB가 필수는 아니다. 저장된 비어 있지 않은 커스텀 정책은 자동으로 번역되지 않으며 옛 중국어 입력 경계·출력 계약의 중복 부착 가능성은 유지한다. 저장/로드 및 중복 판정 로직은 변경하지 않았고 실제 DB·모델 출력은 미확인이다. `prompt.go`는 세 상수를 모두 검토했으나 이번 용어 보류가 남아 번역 완료가 아니며 U7 패키지 전체 완료도 아니다. 기존 U1·U2·U5·U6 및 U7 이력·보류를 보존한다.
- **실행 검증·산출물**: 설치·빌드·타입 검사·lint·테스트·go test·gofmt·렌더·앱·스모크·Git 훅은 **사용자 요청으로 미실행**이다. 정적 텍스트 비교와 예시 JSON 구조 대조는 실행·모델 동작 검증이 아니다. 두 파일 전체 diff는 `/tmp/artex-u7-default-policy.diff`, 상세 검토·보류 보고는 `/tmp/artex-u7-default-policy-review.md`에 저장했다. stage·commit·push·다음 파일은 진행하지 않았다.

### U7 · 기본 판정 정책 용어 보류 해소  [세 상수 번역·정적 검토 완료 / U7 후속 검토 보류 / 실행 미검증]
- **용어 등록·범위**: 사용자 확정 용어 16개를 중복·충돌 확인 후 용어집의 U7 판정 정책·공격 기법·시스템 설정 문맥으로 등록했다. `可逆性`=가역성, `可逆`=가역적, `不可逆`=비가역적, `复合命令`=복합 명령, `中间件`=미들웨어, `启动项`=시작 항목, `计划任务`=예약 작업, `压测`=부하 테스트, `布尔`=boolean, `时间盲注`=time-based blind, `路径遍历`=경로 순회, `文件包含`=파일 포함, `指纹识别`=핑거프린트 식별, `爬取`=크롤링, `抓包`=패킷 캡처, `落脚后`=초기 접근 성공 후를 적용했다. 时间盲注의 영문 전체 명칭은 time-based blind SQL injection이며 이번 SQLi 목록에서는 time-based blind로 표기한다. 기존 다른 문맥·항목은 보존했다.
- **보류 해소·표현**: 직전 60–62·68·71·80·83·87·90–92·102·106–107·112행의 보류 문장 15곳을 모두 번역했다. A1은 `SQLi(UNION/boolean/time-based blind/쓰기 구문을 포함한 인젝션)`으로 맞췄다. 53행의 “실제 사용자의 이용을 불가능하게 합니까?”와 97행의 “정말 판단을 확정할 수 없는 경우에만 ASK”를 적용했다. 그 밖의 직전 번역은 유지했다. 이전 보류 기록과 후보는 당시의 이력이다.
- **정책·정적 검토**: 가역성은 원문의 판단 축을 번역했으며 정책 자체를 재평가하지 않았다. 역할·두 판단 축·공통 제약 우선순위·D1–D6·A1–A6·ASK 조건·판정 순서·예시 판정과 근거를 원문 전체와 다시 대조했다. 금지·필수·허용·조건·예외 강도를 유지했고 새 설명이나 출력 언어 지시를 추가하지 않았다. `DefaultJudgePrompt`는 50–128행, 원본/현재 79줄이며 `prompt.go`는 205줄이다. 원본 줄 배치·빈 줄·들여쓰기·개행·이스케이프·마크다운·숫자·기존 영어·명령·경로·예시 값·JSON 구조·키·enum·규칙 번호를 유지했다. 확정된 boolean/time-based blind/Password Spraying 영문 적용은 중국어 용어 번역으로 별도 구분한다.
- **잔여 분류·계약**: 세 상수의 번역 대상 누락이나 추가 용어·계약 보류는 발견하지 못했다. 중국어 DNT는 형식 설명·정의·예시 및 파서의 `实际操作：`·`；成功后的后果：`·`；命中规则：`이며 전각 부호까지 보존했다. 115행 `WHERE 全表`는 **원문 예시 보존**으로 기록하며 파서 계약이라고 단정하지 않는다. 직전 기록의 포괄적인 DNT 분류와 구분한다.
- **테스트·범위 밖 보존**: 파서와 출력 예시의 앵커·JSON 키·판정값을 정적으로 대조했고 이번 번역에 직접 연결된 테스트 기대값 수정 필요는 발견하지 못했다. 기존 live 테스트(`server/intercept_live_test.go:124–134`)의 중국어 동사 의존은 후속 검토 보류로 유지한다. 완료된 `JudgeContextBoundary`·`JudgeOutputContract`·`EffectiveJudgePrompt`·파서·테스트·저장/로드·DB 및 이번 상수 밖 내용은 작업 전과 동일하다.
- **파일·단위 상태**: `prompt.go`의 세 상수는 **번역·정적 검토 완료 / 후속 호환성·테스트 검토 보류 / 실행 미검증**이다. 이것은 U7 패키지 전체 완료가 아니다. 옛 중국어 경계·출력 계약을 포함한 커스텀 프롬프트의 중복 부착 가능성과 저장된 커스텀 정책 미번역 상태를 유지하며 실제 DB·모델 출력은 미확인이다. 기존 U1·U2·U5·U6 및 U7 기록·보류를 보존한다.
- **실행 검증·산출물**: 저장소·의존성 코드·Git 훅·설치·빌드·타입 검사·lint·테스트·go test·gofmt·렌더·앱·화면·스모크는 **사용자 요청으로 미실행**이다. 누적 세 파일 전체 diff는 `/tmp/artex-u7-default-policy-final.diff`, 갱신된 정적 검토·보류 해소 보고는 `/tmp/artex-u7-default-policy-review.md`에 저장했다. 기존 프롬프트 변경·미추적 담당표·다른 파일·Git index를 보존하며 stage·commit·push·다음 파일은 진행하지 않았다.


### U7 · 인터셉트 런타임·검토 입력·감사 trace  [대상 번역·정적 검토 완료 / 패키지 잔여·후속 검토 보류 / 실행 미검증]
- **시작 상태·기준**: 로컬 `work/ko-translation`, HEAD `141d36ff7d78a4aa941f42f8cf8d711ef96c43b2`를 기준으로 진행했다. 진행 중인 병합·staged 변경·대상 파일의 기존 변경은 없었고, 직전 `DefaultJudgePrompt`와 두 공통 상수의 번역은 로컬 HEAD에 반영돼 있었다. fetch·병합 없이 기존 `TRANSLATION_PROMPT.ko.md` 변경·미추적 담당표·다른 파일·Git index를 보존했다.
- **파일별 결과**: `intercept/intercept.go` 전체(689줄)를 검토해 판정 실패·입력 확인 안내, 허용/차단/사람 승인으로 전환 라벨, 기본 규칙 안내, 취소·시간 초과·사람 판정·중복 승인 오류의 자연어 문자열 18곳을 번역했다. `intercept/review_context.go` 전체(70줄)는 유효하지 않은 JSON 파라미터의 자체 오류 1곳을 번역했다. `intercept/trace.go` 전체(214줄)는 실행 종료 시 결과를 받지 못한 경우의 자체 안내 1곳을 번역했다. 세 파일의 주석은 모두 기존 영어여서 변경하지 않았다.
- **계약·소비 경로**: `[模型]` 접두와 판별(`intercept.go:480–482`)은 DB의 SQL/HasPrefix 출처 분류와 프런트 승인 기록의 startsWith/정규식에 연결되므로 보존했다. `工具 %s 请求审批 (#%d)`(`:612`)는 `transcript.tsx:196–202`의 ID·도구명 정규식과 연결되어 문자열 전체를 보존했다. 그 밖의 이번 자연어 문구를 비교·파싱하는 소비부는 발견하지 못했다. 자체 메시지는 guard의 모델용 차단 안내, DB 감사 기록, 승인 화면·API 오류로 전달된다. 모델 지시의 금지·승인 필요·대기 조건을 유지했다. 외부 오류·DB 규칙 message·모델 판정 사유·사용자 입력·실제 도구 출력은 변환하지 않았다.
- **정적 보존**: 세 소스의 원본/현재 줄 수는 각각 689/689·70/70·214/214다. 문자열을 가린 텍스트의 동일성으로 모든 비문자열 코드와 주석 보존을 확인했고, 변경 문자열별 영어 토큰·포맷 지정자와 순서·이스케이프 및 줄별 들여쓰기도 대조했다. 규칙 우선순위·허용/승인/거부·실패 시 처리·승인 시간 제한·입력 선택·Worker 배경 제외·잘림·UTF-8·설정 저장/로드·동시성·반환 구조·키·enum·상태·숫자·정규식·코드 구분자는 변경하지 않았다. `prompt.go`·파서·테스트는 그대로다.
- **관련 테스트**: 패키지의 `prompt_test.go`·`review_context_test.go`·`trace_test.go`, `agent/capture_approval_test.go`·`agent/review_context_test.go`, `server/intercept_review_test.go`·`server/intercept_detail_test.go`·`server/intercept_live_test.go`, `db/intercept_detail_test.go`와 관련 소비부를 읽기 전용으로 대조했다. invalid JSON 검사는 오류 존재 여부, trace 종료 검사는 `unknown` 상태, 승인 수명주기 검사는 키·상태·출력·출처를 검사한다. `ErrAlreadyDecided` 소비부는 errors.Is로 분기하며 HTTP 테스트는 409 상태를 검사한다. 이번 번역에 직접 연결된 기대 문자열 변경 필요는 발견하지 못했다. DB 테스트의 `人工允许执行`은 독립 입력이며 이번 소스의 기대값이 아니다. 테스트를 수정·실행하지 않았다.
- **잔여·범위 밖**: 세 소스의 중국어 잔여는 위 `[模型]` 세 줄과 승인 요청 요약 한 줄의 DNT뿐이다. 미분류 번역 누락·추가 용어/계약 보류는 발견하지 못했다. `guard/guard.go`의 모델용 플랫폼 차단 안내, `server/intercept.go`의 판정 모델 설정·형식 오류, 저장된 규칙/커스텀 프롬프트·기존 감사 기록과 외부 provider 오류는 범위 밖이며 변경하지 않았다. U1·U2 기존 보류와 직접 연결되는 새 의존은 발견하지 못했고 기존 상태를 유지했다.
- **U7 범위·후속 보류**: 최신 추적 목록의 패키지 파일은 7개다. 네 비테스트 소스(`prompt.go`·이번 세 파일)는 번역·정적 검토 대상에 포함됐지만, 세 테스트(`prompt_test.go`·`review_context_test.go`·`trace_test.go`)는 관련 검증으로 읽었을 뿐 독립 번역 묶음은 미작업이다. 중국어 계약·길이 검증·실제 입력/출력 fixture를 임의로 번역하지 않았다. 옛 중국어 경계/출력 계약이 포함된 커스텀 프롬프트의 중복 부착 가능성, 저장된 커스텀 정책의 미번역 상태, live 테스트의 중국어 동사 의존은 기존 후속 검토 보류로 유지하며 실제 DB·모델 출력은 미확인이다. 따라서 U7 전체 완료·최종 동작 검증 완료로 표시하지 않는다. 기존 팀원 기록·다른 단위 상태는 보존한다.
- **실행 검증·산출물**: 저장소·의존성 코드·Git 훅·설치·빌드·타입 검사·lint·테스트·go test·gofmt·렌더·앱·화면·스모크는 **사용자 요청으로 미실행**이다. 허용 파일 전체 diff는 `/tmp/artex-u7-runtime.diff`, 상세 정적 검토·잔여 보고서는 `/tmp/artex-u7-runtime-review.md`에 저장했다. git add·commit·push·다음 단위는 진행하지 않았다.


### U7 · 패키지 테스트 분류·전체 파일 종합  [전체 파일 번역·정적 검토 완료 / 후속 호환성·테스트 검토 보류 / 실행 미검증]
- **시작 상태**: 로컬 `work/ko-translation`, HEAD `131d12b6fb078d3a12a0ca973e1a6914ba1a0876` 기준이다. 직전 런타임 번역 커밋과 운영 소스의 HEAD 일치를 확인했다. 진행 중 병합·staged 변경·대상 파일의 기존 변경은 없었으며 fetch·병합은 하지 않았다.
- **테스트 파일별 결과**: `intercept/prompt_test.go` 전체 75줄, `review_context_test.go` 전체 184줄, `trace_test.go` 전체 135줄을 읽고 중국어와 전각 문장부호를 검증 목적별로 모두 분류했다. 세 파일의 개발자 주석·테스트 실패 안내는 모두 영어여서 번역 대상이 없으며 세 테스트 파일은 바이트 단위로 변경하지 않았다.
- **보존 근거**: `prompt_test.go`의 세 판정 앵커·앵커 제거/접미사 검사·유효/무효 모델 응답·정확한 사유 반환 fixture와 중국어 반복 길이 검증을 보존했다. `review_context_test.go`의 사용자 배경·거부/부분 출력·인젝션 문구·커스텀 정책·판정 사유는 입력 선택/배제와 정확한 보존을 검증하는 fixture다. 중국어 반복 입력은 잘림·전체 arguments 보존·UTF-8을 검사한다. `trace_test.go:98,100`의 `中文`·`中`은 바이트 상한·잘림·UTF-8 경계 검증 데이터다. 앵커는 DNT이며 나머지 fixture는 검증 목적에 따라 원문 보존한 것으로 구분하고 중국어 잔여 전체를 계약 DNT로 단정하지 않는다.
- **기대값·정적 동일성**: 현재 파서 앵커·반환 사유, BuildReviewInput의 입력 구조·배경 선택·잘림, trace의 결과 상태와 대조했다. 소스 번역과 직접 연결된 기대값 불일치는 발견하지 못했다. invalid JSON 검사는 오류 존재 여부만, missing-result 검사는 출력 본문을 무시하고 `unknown` 상태만 검사하며 공통 블록은 상수 자체를 참조한다. 테스트 함수·조건·검증 방식·입력·기대값·순서·숫자·포맷·영어·들여쓰기·개행을 모두 보존했다. 미분류 번역 누락·추가 용어/계약 보류는 없다.
- **U7 전체 목록·현재 상태**: GUIDE의 패키지 전체 범위를 최신 추적 7개 파일과 대조했다. 네 운영 소스(`prompt.go`·`intercept.go`·`review_context.go`·`trace.go`)는 이전 번역·정적 검토와 로컬 커밋 반영을 확인했고, 세 테스트는 이번에 독립 분류 검토까지 마쳤다. 이전 기록의 테스트 독립 번역 묶음 미작업 상태는 이번 검토로 해소했으며 이력은 보존한다. 추적 목록에 미작업 파일은 없다. 이는 실행 검증이나 모든 후속 문제 해결을 뜻하지 않는다.
- **후속 보류 유지**: `server/intercept_live_test.go:124–134`의 중국어 동사 검사 의존은 U11 서버 테스트와 U7 프롬프트의 후속 협업 검토 항목이다. `prompt.go:20–27`의 전체 상수 Contains는 저장된 옛 중국어 공통 블록을 인식하지 못해 새 한국어 블록을 추가할 수 있다. `intercept.go:346–349`의 비어 있지 않은 저장 정책은 그대로 사용하므로 기존 커스텀 본문은 자동 번역되지 않는다. 호환성 정책·저장/로드·테스트 조정은 별도 범위 확정 후 검토하며 실제 DB·모델 출력은 미확인이다. 세 항목을 이번에 해결한 것으로 표시하지 않는다.
- **보존·실행 검증·산출물**: 기존 이력·다른 단위 상태·TRANSLATION_PROMPT 변경·미추적 담당표·모든 다른 파일·Git index를 보존했다. 저장소·의존성 코드·Git 훅·설치·빌드·타입 검사·lint·테스트·go test·gofmt·렌더·앱·화면·스모크는 **사용자 요청으로 미실행**이다. 허용 파일 전체 diff는 `/tmp/artex-u7-tests.diff`, 파일별 잔여 분류와 U7 종합 보고서는 `/tmp/artex-u7-final-review.md`에 저장했다. git add·commit·push·U8은 진행하지 않았다.


### U8 · report 패키지 전체  [전체 파일 번역·정적 검토 완료 / 실행 미검증]
- **범위·시작 상태**: GUIDE의 U8은 `report/` 패키지 전체이며 최신 추적·실제 파일은 `report/report.go`·`report/findings.go` 두 개다. `work/ko-translation`, 로컬 HEAD `e2b5ac04ba501a85dfb8a61fcfe3d83e073fafd5` 기준으로 작업했다. 시작 당시 대상 파일의 기존 변경·staged 변경·진행 중 병합은 없었고 fetch·병합은 하지 않았다.
- **파일별 결과·소비자**: `report.go`의 표시 문자열 12개와 `findings.go`의 표시 문자열 41개·중국어 주석 14개를 번역했다. 일반 작업 Markdown, 취약점 종합/단건 Markdown·CSV 열 제목·심각도 라벨·트래픽 증거 안내·첨부 링크 표시명이다. 이 패키지에는 모델용 프롬프트·도구/스키마 설명·자체 오류가 없다. `server/server.go`의 getReport·exportFindings와 `server/finding_traffic_export.go`의 ZIP 작성 경로에서 다운로드로 전달되고 프런트 `api.report`는 텍스트, `exportFindings`는 blob으로 받는다. DB의 이름·분류·요약·상세 보고서·증거·작업 설명·노트·HTTP 원문은 변환하지 않았다.
- **정적 보존**: 원본/현재 줄 수 `report.go` 99/99·`findings.go` 200/200과 각 줄의 들여쓰기·개행을 유지했다. 주석·문자열 토큰을 제외한 텍스트 동일성, 변경 토큰의 영어·숫자, 문자열별 포맷 지정자의 종류·개수·순서·이스케이프·Markdown 구조를 대조했다. 인수·정렬·필터·조회·반환 구조·enum·JSON 태그·정규식·날짜 형식·CSV 열 수/순서·UTF-8 BOM·ID 구분자·첨부 경로·원본 데이터 전달은 보존했다. AST·렌더·실행 검증을 대신하는 확인은 아니다.
- **잔여 분류·보류 해소**: `report.go:66`의 `b.WriteString("、")`는 정렬된 자산 유형 목록 사이 표시용 출력 구분자임을 확인해 최초에는 허용 위치 확인을 보류했다. 이후 사용자 명시 허용으로 `b.WriteString(", ")`로 변경해 구분자 보류를 해소했다(파서 계약 DNT가 아님). CSV 헤더의 `트래픽 증거ID`도 `트래픽 증거 ID`로 정리했다. `findings.go:125`의 `critical_SQL注入_#123.md`는 파일명 예시 원문 보존이다. 그 외 번역 누락·중국어 계약 DNT·미등록 용어·불명확한 계약·후속 프롬프트 보류는 발견하지 못했다. 전체 잔여 0이나 U8 최종 동작 완료로 표시하지 않는다.
- **관련 테스트**: report 패키지에 테스트 파일은 없다. `server/finding_traffic_test.go`의 ZIP/Markdown 첨부 경로·원본 바이트 검사(199–234), md-single/json/csv 증거 ID 보존 검사(467–475)를 읽기 전용으로 확인했다. 제목·CSV 헤더·한국어 표시 문구를 직접 비교하는 테스트나 저장소 파싱 소비자는 발견하지 못해 기대값 수정은 필요하지 않다. 기존 중국어 보고서 입력(456행)은 DB 보고서 원문 보존 fixture이며 변경하지 않았다. 외부 CSV 소비자는 미확인이다.
- **기존 상태·범위 밖**: 기존 팀원 기록과 다른 단위 상태, U7의 옛 중국어 공통 블록 중복 부착·live 테스트 중국어 동사 의존·커스텀 정책 미번역 보류를 그대로 보존했다. 기존 DB의 보고서 본문·증거·사용자 입력·외부 응답은 이번 소스 번역으로 한국어화되지 않는다. U1/U2 기록 사이 현재 상태와 후속 확인/실행 보고의 차이는 기존 기록이며 이번에 정정하지 않았다.
- **실행 검증·산출물**: 저장소·의존성 코드·Git 훅·설치·빌드·타입 검사·lint·테스트·go test·gofmt·템플릿 렌더·앱·화면·스모크는 **사용자 요청으로 미실행**이다. `/tmp/artex-u8-final.diff`와 `/tmp/artex-u8-review.md`에 누적 diff·상세 보류를 기록했다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존하고 git add·commit·push·다음 단위는 진행하지 않았다.
