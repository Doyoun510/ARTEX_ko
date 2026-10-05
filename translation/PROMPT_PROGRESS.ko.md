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

## U5~U8·U10~U15  [todo]
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
