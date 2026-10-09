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

### 묶음 D-2·D-3 (server/ 긴 꼬리 완료)  [정적검토 done · 실행 검증 일부]
- **D-2**: intercept·asset_intercept·task_intercept·constraints_api·intent_intervention·side_questions.
- **D-3**: assembly·assets·chatupload·findings_groups·llmpool·logsink·mcpdiscover·task_archive_package·task_assets·task_categories·task_metadata·task_resolution·task_templates·triggers·webui_embed·webui_stub·workspace·skill_zip.
- **계약/특수 처리**:
  - intercept.go 판정 형식 에러의 판정 앵커(§5.7)·server_mgmt `已存在`·task_archives `归档不存在`·chat_mentions 멘션 토큰 정규식/맵·conversations `新对话`/SDK 참조·finding_retests 프롬프트 버전 라벨(db/config U12 일치)·fallbackChat 접두·각종 설계문서 경로·orchestration 마이그레이션 상수 = **DNT 보존**.
  - **logsink.go**: 로그 분류 키워드 리스트에 한국어 키워드(실패·폐기·거부·불가·비활성·재시도) **추가**(중국어 키워드 유지) — 로그를 한국어로 번역해 분류가 깨지지 않게 보강. 기존 동작 유지.
  - 모델-읽기(assembly bash 상호작용 노트·shell 힌트, chatupload 첨부 노트, finding/traffic 도구 설명)는 지시 강도 보존 번역.
  - 회사/분류/템플릿 `已存在`은 프런트 매칭 없음(실제 페이지 기준) 확인 후 번역.
- **계약 테스트 동기화**: chat_mentions_test(`잘림`)·findings_groups_test(수동 제출 요약) 갱신.
- **실행 검증(Go 1.26.3)**: `gofmt -l` 무출력 ✅ / `go build ./...` OK ✅ / `go vet ./server/` clean ✅ / 순수·no-DB 서버 테스트 + agent 회귀 PASS ✅. DB 필요 서버 테스트는 미검증.

### 묶음 D 요약 — **server/ 전체 한국어화 완료**
server/ 비테스트 파일에 남은 중국어는 전부 **의도적 DNT**(설계문서 경로·§5.7 계약 앵커·멘션 토큰 파서·logsink 분류 키워드·마이그레이션 상수). U11 server/ 작업 종료. (trans 트랙 B의 notify/·db/ 등 타 단위는 범위 밖)
## U12 · 데이터/DB (트랙 B)  [완료 — 번역·정적·빌드/vet PASS]
사용자(사람1)가 U11(server/) 완료 후 U12 착수. 실행 환경에 **Go 1.26.3 존재** → `go build ./...`·`go vet ./db/`·`gofmt`·no-DB 테스트를 매 커밋 실행. DB 필요 테스트는 PostgreSQL 미연결로 미검증(U11과 동일).

### 완료 배치 (커밋됨)
- **배치1(주석)**: `triggers.go`·`constants.go`·`task_context.go` — 순수 주석. 계약 문자열 없음.
- **배치2(파서 오류)**: `constraints.go`·`asset_dsl.go` — fmt.Errorf 11건. §3 확인(프런트 매칭·코드 비교 0). allow/deny·task_id/company_id·token·%s 보존. DSL 파서 테스트 PASS.
- **배치3(중형 11파일)**: `nkey.go`·`task_scope.go`·`side_questions.go`·`task_assets.go`·`asset_intercept.go`·`task_intercept.go`·`tools.go`·`intercept.go`·`task_archives.go`·`intercept_execution.go`·`task_archives_restore.go`. §3 전수 확인.
  - **DNT 보존**: `[模型]` 접두(§5.7 계약) — `intercept.go` SQL `LIKE '[模型]%'`·`HasPrefix`, `task_archives_restore.go` `HasPrefix`. enum `'block'/'allow'`·scope kind 값(company/root_domain/…)·도구명·`%w/%d/%s`·`company store`.
  - `manualTaskScopeSummary`(상수): 테스트가 **상수 심볼** 비교(리터럴 아님)라 값 번역 무해(`task_assets_test.go:47`).
  - 용어: 旁路问题=보조 질문(U9 일치), 归档=아카이브, 范围未锚定=범위 앵커 미설정, 覆盖度=커버리지, 墙钟=실제 경과 시간.
- **배치4(설정+계약 쌍)**: `config.go` 주석 전부 + **프롬프트 버전 라벨 계약 해소**.
  - **⚠️ 계약(U11 이월) 해소**: `内置默认`→`내장 기본값`, `恢复为内置默认`→`내장 기본값으로 복원`. 이 값은 `agent_prompts.note`(순수 표시, §3: Go 비교·프런트 매칭 0). 같은 라벨을 쓰는 **`server/finding_retests.go:230`의 SQL 저장 값도 동일 문자열로 동기화**(두 writer 일관). 구 DB 레코드의 구 라벨은 note 미매칭이라 표시만 혼재(기능 무해, §7.1 새 DB 전제).
  - DNT: 설계문서 경로 `docs/交互式shell设计.md`(U11 `server/assembly.go:311` 선례로 보존 — 포크에 실물 없어도 경로 참조)·thinking.type/serial/parallel enum.

### 완료 배치 (대형·테스트)  [2026-10-06]
U15 완료 후 사용자가 U12 재개. wip/u12-db(소/중형 18파일) work에 머지 후 대형·테스트 전부 번역.
- **대형 생산**: notification·notification_delivery·finding_assets·db.go·exploration·asset_intercept_match·findings·tasks·llmretry·finding_retests·finding_traffic·company_scope·companies·chat_mentions·assets.
- **⚠️ db.go 안전장치 준수**: 내장 인터셉트 규칙은 `name`·`message`만 번역, `pattern`/`target`/`typ`/`action`/`priority` 전부 불변(diff에 pattern 라인 0건, 패턴 19개 유지). `[内置]`→`[내장]`.
- **⚠️ company_scope.go**: `备案`(L160 `strings.Contains` 파서)·`.．。` 입력 점·ICP 예시 DNT 보존.
- **계약 동기화**: tasks.go '작업 생성 시 회사 연결:' ↔ tasks_test.go:610, `内置默认`→`내장 기본값`(db/config.go+server/finding_retests.go), task_archives_restore 경고 ↔ task_archives_test.go:274 substring.
- **테스트(~20개 _test.go)**: §7.4대로 주석·t.* 진단·서브테스트명만 번역, 픽스처/계약/인코딩 데이터 보존(VulnClass·ICP·`界` 길이·`余额不足`/`额度不足` 할당량 계약·`[模型]`·프롬프트 라운드트립·FTS 쿼리쌍).
- **검증**: `go build ./...`·`go vet ./db/`·`gofmt`(기준선 대비) PASS. DB 필요 테스트만 PostgreSQL 미연결 미검증.
- **잔여 중국어(의도 DNT만)**: 설계문서 경로(`任务级超时与收尾设计.md`·`交互式shell设计.md`·`LLM重试设计.md`·`资产模型与自动关联设计.md`)·`备案`·`[模型]` 계약·테스트 픽스처.

## U13 · 그 외 백엔드 패키지 (트랙 B)  [완료 — 번역·정적·실행 일부 PASS]
사용자(사람1)가 U12 보류 후 착수(다형 권장 순서). 대상: traffic·selfupdate·llmrec·llmpool·guard·evidence·mcphttp·cmd·config·enrich. Go 1.26.3로 매 커밋 `go build ./...`·`go vet`·`gofmt`·no-DB 테스트 실행(PASS). DB 필요 테스트만 미검증.

- **소형 7파일**: evidence·mcphttp·cmd/artex·config·enrich·llmpool(pool/health). 로그·오류·CLI 배너·주석. §3 확인(프런트·코드 비교 0).
- **guard·llmrec**: guard는 모델 전달 차단 래퍼(`【ARTEX 平台管控…】`)·systemBlockMessage. §3: 배너를 파싱/매칭하는 **비테스트 소비부 0**(guard가 유일 생산자). intercept/server 테스트의 동일 배너는 입력 픽스처일 뿐 guard 출력과 미비교 → 번역 무영향, 그 픽스처는 U7/U11 소유라 미변경. `[内置]`→`[내장]`(§3 매칭 0).
- **traffic**(traffic.go+evidence.go): 모델 전달 도구 설명 traffic_search/traffic_get/traffic_read_blob + 로그·오류.
  - **⚠️ 마이그레이션 계약**: `TrafficSearchDescription`은 `server/finding_workflow.go`의 새 값($1), 비교 앵커 `legacy`($2, U11 소유·이미 한국어·구버전 설명)는 불변. 소스를 신버전 한국어로 번역해 구→신 업그레이드가 올바르게 동작.
- **selfupdate**(github·selfupdate·stage·bootstrap): 자가 업데이트 주석·SSE 진행·오류·State.Detail. §3 확인. DNT: Phase enum 값(idle/downloading/…)·exit 코드·경로(artex.new/.old/.sha256)·SHA256SUMS·버전 리터럴.
- **테스트 4파일**(traffic 3 + selfupdate_test): §7.4대로 주석·진단 메시지만 번역. **인코딩/계약 데이터 보존** — 중국어 FTS 입력(`内网测试账号`·`内网测试` 쿼리 페어·`中文正文`·`老数据正文`·`清空后仍可录制`)·버전·URL·parseSums 픽스처.
- **용어**: 轮询=순환 전환·故障转移=장애 조치·熔断=회로 차단(프런트 llm "즉시 복구"/"회로 차단" 일치)·流量=트래픽·记录代理=기록 프록시·透传=패스스루·压实=압축·冷/热=cold/hot.
- **U13 잔여(의도 DNT)**: 설계문서 경로 `docs/资产模型与自动关联设计.md`·`docs/交互式shell设计.md`, traffic 테스트 FTS/인코딩 픽스처. **전체 `go build ./...`·`go vet ./...`·`gofmt` PASS.**

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


### U10 · 알림 시스템 전체 [부분 번역·전체 파일 정적 검토 완료 / 용어 및 범위 밖 테스트 검토 보류 / 실행 미검증]
- **범위·시작 상태**: 사용자 명시 지정으로 GUIDE의 담당 미정 U10을 진행했다. 로컬 `work/ko-translation`, HEAD `e7a36d8d1ad797c5e18a54226aa29283a2e199d8` 기준이며 fetch·병합은 하지 않았다. `notify/` 추적 파일 26개(소스 15개·테스트 11개)를 모두 읽고 검토했다. 별도 템플릿 파일은 없으며 템플릿은 Go 문자열 안에 있다. 시작 당시 대상 파일·진행 기록의 기존 변경, staged 변경, 진행 중 병합은 없었다.
- **번역 범위**: 채널 어댑터·설정 검증·자체 오류·메일/Markdown/HTML/카드/Telegram 본문·상세 링크·모아 보내기 안내·주석·테스트 실패 안내를 번역했다. 심각도는 emoji를 유지한 심각/높음/중간/낮음으로 맞췄다. 기존 StatusLabel 함수의 9개 한국어 표시 문자열과 흐름은 그대로 보존했으며 프런트 finding과 일치한다. 표시 목록 구분자는 용어집 7.2의 명시 위치인 mask.go와 render.go에서만 변경했다.
- **소비부·계약**: server/notifier.go의 send/itemFor, server/notify_api.go의 저장 검증·테스트 발송, db의 알림 필터·이벤트 소비부와 관련 테스트를 읽기 전용으로 확인했다. 오류 처리는 IsPermanent/error wrapper와 HTTP/SMTP 숫자 코드에 의존하며 번역한 자체 설명을 파싱하는 운영 분기는 확인하지 못했다. URL·서명·채널/이벤트 식별자·JSON/템플릿 키·사용자 템플릿·마스킹 표식·DB 데이터·외부 API/SMTP 오류 본문은 유지했다. 실제 모델용 프롬프트는 발견하지 못했다.
- **테스트 대응**: notify 내부의 직접 연결된 기대 문자열 12개를 channels_test.go 4개·email_test.go 1개·pack_test.go 4개·render_test.go 2개·webhook_template_test.go 1개에서 조정했다. 테스트 입력·사용자 템플릿·서명 기준값·HTTP/SMTP 응답·Unicode/UTF-8/잘림 입력·테스트 식별명·검증 조건은 보존했다. 정적 대조 중 생긴 URL 입력 및 따옴표 편집 오류는 복구했고 최종 입력/코드 비교에 포함했다.
- **범위 밖 기대값 보류**: server/notify_api_test.go:244,353,430,491,719,779의 표시 문구 의존은 이번 수정 범위 밖이므로 원문 유지했다. 특히 719행의 부정 검사는 한국어 상세 링크를 놓칠 수 있어 후속 조정이 필요하다. 491행의 已修复는 이전 StatusLabel 번역에 따른 기존 불일치다. 같은 파일의 351·354행도 기존 서버 번역과 불일치하며 U10 신규 변경과 구분한다. 테스트 실행 없이 통과/실패를 단정하지 않는다.
- **용어 보류**: 叶子包·알림 문맥 兜底·脱敏·哨兵·握手·SMTP 信封·灰名单·环回/链路本地/组播·DNS 重绑定·네트워크 문맥 守卫·半盲读原语·反代·静默降级·SSRF 문맥 跳板을 포함한 문장/주석 블록을 원문 유지했다. 사용자 노출 문자열 중 HTTP/DingTalk의 링크로컬 거부 오류 두 곳과 SMTP 핸드셰이크 오류 한 곳도 용어 보류다. 기존 다른 문맥의 확정 용어를 확대 적용하지 않았다. 원문 전체·위치·문맥·권장안은 상세 보고서에 모았다.
- **잔여 분류**: 보류 설명과 테스트 실패 안내, 테스트 식별명, 문자열 매칭/UTF-8/출력 보존 fixture, 악성 Markdown·URL 원문 예시를 구분했다. 잔여 전체를 DNT나 번역 완료로 표시하지 않는다. 기존 DB의 이름·취약점 유형·요약 및 외부 API 오류는 이 소스 번역으로 한국어화되지 않는다. 미분류 누락은 발견하지 못했지만 용어 보류와 범위 밖 기대값 조정이 남아 U10 전체 번역 완료로 표시하지 않는다.
- **정적 검토**: 26개 Go 파일 총 5055줄의 원본/현재 줄 수·각 줄 들여쓰기가 일치한다. 주석/문자열 토큰 밖 텍스트, 문자열별 포맷 지정자의 종류·순서·개수와 이스케이프를 대조했다. 상태 9개·심각도 4개가 프런트와 일치하며 HTML/Markdown 구조·템플릿 키·서명 기준값·정규식·조건·숫자·반환 흐름을 보존했다. 이는 컴파일·AST·렌더·동작 검증을 대신하지 않는다.
- **기존 상태·실행·산출물**: U7 후속 3건 및 기존 모든 이력·다른 단위 상태를 보존했다. 기존 TRANSLATION_PROMPT 로컬 변경·미추적 담당표·다른 파일·Git index를 보존했다. 설치·빌드·타입 검사·lint·테스트·gofmt·렌더·앱·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-u10.diff`, 상세 검토·보류 `/tmp/artex-u10-review.md`. git add·commit·push·다음 단위는 진행하지 않았다.


### U10 · 확정 용어 반영 및 최종 정적 검토 [notify 대상 번역·정적 검토 완료 / 범위 밖 서버 테스트 조정 보류 / 실행 미검증]
- **이번 범위·기존 작업 보존**: 직전 notify 전체 부분 번역을 이어서 사용자 확정 H01–H59를 모두 해소했다. 소스·테스트 26개 전체의 줄 수는 5055/5055이며 기존 부분 번역을 유지했다. 용어집의 跳板 항목에 SSRF 문맥, 守卫 항목에 네트워크 연결 문맥을 보충하고 나머지 20개 문맥별 확정 항목을 등록했다. 기존 다른 문맥의 확정 번역은 보존했다.
- **문맥별 적용**: Greylisting·민감정보 제거·DNS 리바인딩·link-local 주소·루프백·멀티캐스트·핸드셰이크·SMTP 봉투·연결 보호 검사·센티널 값·최하위 패키지·리버스 프록시·경유 서버·오류 응답을 통한 부분 정보 읽기·오류를 알리지 않고 필터링이 해제되는 동작을 적용했다. 兜底는 빈 kind=기본값, 미지정 본문=기본 요청 본문, 잘못된 필터=기본 처리, deadline=최종 보호 조치, 길이 제한=최종 잘림 처리로 보호, 구조화되지 않은 오류=대체 처리, 제한 후 재시도=보완 수단으로 구분했다.
- **수정·테스트 대응**: H01–H59의 주석 49블록·자체 오류 3개·테스트 실패 안내 7개를 번역하고 channel.go의 '설정 읽기 helper' 공백을 정리했다. 새로 번역한 자체 오류 세 개를 직접 비교하는 notify 내부 기대값은 발견하지 못해 추가 기대값 수정은 없다. 이전에 맞춘 직접 기대 문자열 12개는 그대로 유지했다. 모든 테스트 입력·fixture·조건·로직·실행 순서는 보존했다.
- **원문 불일치**: feishu.go:25–26의 초당 5회·분당 100회는 원문에도 있는 환산 불일치다. 사용자 요청에 따라 숫자·DefaultRatePerMin 반환값을 변경하지 않았으며 후속 기술 확인 사항으로 기록했다.
- **잔여·후속**: 지정 용어 보류 H01–H59는 모두 해소했으며 새 미등록 용어나 계약 보류는 발견하지 못했다. 중국어 잔여는 테스트 식별명, 필터/UTF-8/출력 보존 입력·기대값, 사용자 템플릿·외부 문구 fixture 및 원문 공격 예시로 분류했다. server/notify_api_test.go는 팀원의 작업을 보존하며 읽기 전용으로 재확인했다. 현재 표시 문구 기대값 의존 6곳과 별도 기존 서버 오류 기대값은 후속 조정 보류다. 팀원의 서버 테스트까지 해결했다고 표시하지 않는다.
- **정적 보존**: HEAD와 직전 부분 번역 두 기준에서 각 Go 파일의 줄 수·들여쓰기·주석/문자열 밖 코드, 리터럴 수·숫자·포맷 지정자 종류/순서/개수·이스케이프와 기존 영어를 대조했다. StatusLabel 9개 및 심각도 4개·서명·템플릿 키·API 계약·조건·반환 구조는 보존됐다. 기존 진행 기록 및 U7 후속 3건을 그대로 유지했다. 컴파일·렌더·실제 전송이나 테스트 통과를 주장하지 않는다.
- **보존·실행·산출물**: 다른 파일·기존 TRANSLATION_PROMPT 변경·미추적 담당표·Git index를 보존했다. 설치·빌드·테스트·lint·gofmt·렌더·앱·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**이다. 누적 전체 diff `/tmp/artex-u10-terms-final.diff`, 갱신된 검토 보고서 `/tmp/artex-u10-review.md`. fetch·병합·stage·commit·push·다음 작업은 진행하지 않았다.


### U10 · 연결된 서버 테스트 기대값 대응 [번역·정적 대응 검토 완료 / 실행 미검증]
- **허용 범위·수정**: 사용자와 U11 담당자의 명시 허용으로 server/notify_api_test.go의 지정 검사 8곳(244·351·353·354·430·491·719·779)에서 기대 문자열 11개만 수정했다. U10 출력에 직접 연결된 6곳(244·353·430·491·719·779)과 기존 서버 번역의 불일치 2곳(351·354)을 구분한다. 앞선 대응 보류 기록은 당시 이력으로 유지하며 지정 8곳의 기대값 대응 보류는 이번 정적 수정으로 해소했다.
- **소스 대조·검증 목적**: 편집 직전에 notify의 심각도·Webhook 검증 오류·모아 보내기 시간/건수·상태 변경·상세 링크·나머지 안내와 server/notify_api.go의 채널 이름/모드 오류를 다시 읽었다. 244행의 SQL注入은 그대로이며 요약 검사는 fixture의 '的摘要'가 아니라 생성 Markdown 라벨 '**요약**:'를 검사한다. 430행의 두 조건, 491행의 두 조건, 719행의 상세 보기 부재 검사, 779행의 부정 조건을 그대로 유지했다.
- **정적 보존**: 서버 테스트 원본/현재 951/951줄, 실제 변경 줄 8개를 대조했다. 문자열을 제외한 전체 텍스트와 줄별 들여쓰기·숫자·포맷·이스케이프가 같다. 이미 일치하는 350·541행·테스트 입력/fixture·조건·구조·순서·실패 안내·다른 기대값·팀원의 기존 변경을 보존했다. notify/의 최종 번역은 추가 수정하지 않았다.
- **상태·한계**: 지정 8곳의 소스/기대값 정적 대응은 완료했으나 서버 테스트 전체의 검토 완료나 테스트 통과를 주장하지 않는다. Feishu 원문 환산 불일치와 U7 후속 3건 및 기존 기록은 유지했다. 기존 U10·로컬 변경·미추적 담당표·Git index는 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 diff `/tmp/artex-u10-final.diff`, 검토 보고서 `/tmp/artex-u10-review.md`. fetch·병합·stage·commit·push·다음 작업은 진행하지 않았다.


### U14 · 루트 스크립트·인프라·설정 예시 전체 [부분 번역·전체 파일 정적 검토 완료 / 용어 보류 / 실행 미검증]
- **범위·시작 상태**: 사용자 지정으로 GUIDE U14 전체를 로컬 `work/ko-translation`, HEAD `1c5c69e3f359d9389015522fc36b57e5e8c5ff04` 기준으로 진행했다. 추적 파일 14개와 실제 파일 목록이 일치한다. 대상 파일·진행 기록의 기존 변경, staged 변경 및 진행 중 병합은 없었고, 기록에서 기존 U14 완료분이나 다른 담당자의 겹치는 작업은 발견하지 못했다. fetch·병합 없이 기존 변경을 보존했다.
- **전체 대상**: `build.sh`·`dev.sh`·`install.sh`·`reset-password.sh`·`start.sh`·`update.sh`, `start.bat`, `Dockerfile`, `docker-compose.yml`·`docker-compose.bench.yml`, `.github/workflows/release.yml`, `config.example.json`, `.env.example`, 루트 `.gitignore`. 다른 소유 디렉터리의 스크립트와 벤치마크 참조 경로의 파일은 포함하지 않았다.
- **번역 범위**: 14개 파일의 중국어 주석·CLI 도움말/질문/자체 오류/상태 안내·설정 예시 `_comment*` 값 257줄을 번역했다. compose의 `POSTGRES_PASSWORD:?`는 필수 변수 오류의 자연어 안내부만 번역하고 변수명·연산자·환경 설정값은 유지했다. 영어만 있는 부분은 그대로 보존했다. 모델용 프롬프트는 발견하지 못했다.
- **용어 보류**: 交叉编译·裁剪·优雅关闭·持久化点·Windows ping 문맥의 兜底·热更新·网络栈/靶机·Git 快进의 13개 문장/주석 블록(19줄)은 전체 원문 유지했다. 확정된 优雅收尾 및 다른 兜底 문맥을 확대 적용하지 않았다. 원문 전체·위치·문맥·권장안·근거는 `/tmp/artex-u14-review.md`에 모았으며 용어집은 수정하지 않았다. 보류를 DNT나 번역 완료로 분류하지 않는다.
- **계약·원문 보존**: 명령·옵션·경로·URL·이미지/서비스 이름·설정 키/값·환경 변수·정규식·출력 구조·SQL·here-doc의 JSON/Python·CI 표현식과 작업 순서·조건을 보존했다. `y/n`·`y/Y`·메뉴 숫자·종료 코드 0/75와 지연값·변수 확장·포맷 지정자·인수는 동일하다. `dev.sh:5`의 README 제목 `单二进制`는 문서 참조 원문, `reset-password.sh:21`의 `pg容器名`은 명령 인수 예시 원문이며, `docker-compose.bench.yml:1`의 `腾讯`는 고유명 원문으로 유지했다.
- **소비부·관련 테스트**: CI의 build.sh 호출·산출물 업로드/다운로드, start.sh/bat의 숫자 종료 코드 판정과 `selfupdate.ExitRestart`, install/update의 `.env` 키 치환·사용자 응답 분기, reset-password의 설정 필드 파싱과 SQL을 읽기 전용으로 대조했다. 도움말은 원래의 `sed -n '2,40p'` 범위를 유지한다. `config/config_test.go`는 DSN·우선순위·오류 존재 여부를, `selfupdate/selfupdate_test.go`는 배포 파일명·바이너리 내용·교체 흐름을 검사하며 번역한 안내 문구의 직접 기대값 의존은 발견하지 못해 테스트 수정은 없다. 저장소 밖의 소비자는 미확인이다.
- **정적 보존**: 14개 파일은 개행 기준 1203/1203줄(논리 행 1204/1204)이며 각 줄 들여쓰기·빈 줄·개행 종류·파일 끝 개행 유무를 유지했다. start.bat의 CRLF 54개, `.gitignore`의 마지막 개행 없음, 기존 영어 토큰·숫자 순서·이스케이프와 변수를 대조했다. 표시 문자열/주석 밖 코드·YAML/JSON 구조·설정값은 보존했으며 배치 출력에 이스케이프 없는 괄호를 추가하지 않았다. 이는 셸 구문 검사·컴파일·CI·실행 검증을 대신하지 않는다.
- **잔여·범위 밖·실행**: 용어 보류와 위 원문 참조/예시/고유명 외 미분류 번역 누락·새 계약 보류는 발견하지 못했다. 범위 밖 `config/`·`selfupdate/`의 중국어 안내·오류/주석과 README는 수정하지 않았다. U7 후속 3건과 모든 기존 단위 이력, TRANSLATION_PROMPT 로컬 변경·미추적 담당표·다른 파일·Git index를 보존했다. 스크립트·설치·빌드·테스트·lint·포맷터·Docker/compose·CI·앱·Git 훅 등 실행 검증은 **사용자 요청으로 미실행**이다. `/tmp/artex-u14.diff`·`/tmp/artex-u14-review.md`에 결과를 저장하고 stage·commit·push·다음 단위는 진행하지 않았다.


### U14 · 확정 용어 반영·보류 해소 최종 검토 [번역·정적 검토 완료 / 실행 미검증]
- **용어 등록·보류 해소**: 交叉编译=교차 컴파일, 裁剪=불필요한 정보 제거, 优雅关闭=정상 종료 절차, 持久化点=영구 저장 위치, Windows ping 문맥의 兜底=대체 대기 수단, 热更新=Hot Reload, 网络栈=네트워크 스택, 靶机=테스트 대상 시스템, 快进=fast-forward를 중복·충돌 확인 후 U14 문맥의 `[확정]` 항목으로 등록했다. 기존 优雅收尾=정상적인 마무리와 다른 兜底 문맥은 그대로 유지했다. H01~H13의 13개 블록/19줄을 모두 번역했으며 이전 보류 기록은 당시 이력으로 보존한다.
- **최종 범위·정적 보존**: 앞서 검토한 U14 추적 파일 14개 전체의 누적 번역은 276줄이다. 이번 소스 추가 변경은 H01~H13의 19줄뿐이고 나머지 직전 번역은 동일하다. 원본/현재 개행 기준 1203/1203줄(논리 행 1204/1204), 들여쓰기·빈 줄·줄 배치·개행·이스케이프·인용부호·기존 영어·숫자·변수·조건·명령·설정값·코드/출력 구조를 유지했다. start.bat의 CRLF 54개와 괄호 이스케이프, .gitignore의 마지막 개행 없음도 유지했다. 추가 영어 Hot Reload/fast-forward는 사용자가 확정한 중국어 용어 번역이며 기존 영어 변경이 아니다.
- **SQLite 설명 읽기 전용 확인**: PostgreSQL은 작업/자산 데이터의 저장소지만 트래픽 인덱스에는 SQLite가 여전히 사용된다(server/manager.go:377~378 → traffic/traffic.go:176~194). 기본 data/traffic/_index/index.sqlite 경로가 확인되어 Dockerfile:41·docker-compose.yml:34의 SQLite 포함 설명은 근거가 있다. 다만 같은 원문 설명의 jwt.key 위치는 현재 기본 경로와 불일치한다: cmd/artex/main.go:115는 config.BaseDir()를 keyDir로 전달하고 server/auth.go:25~44는 실행 파일 옆 keyDir/jwt.key를 사용하며 옛 dataDir/jwt.key를 이전한다. Docker 기본 실행에서 /app/jwt.key는 /app/data 볼륨 밖이라는 정적 추론이며 실제 실행/파일 존재는 미확인이다. 원문 설명 문제로만 기록하고 SQLite·jwt.key 설명과 볼륨/인증 코드는 임의 수정하지 않았다.
- **잔여·테스트 상태**: 14개 파일의 추가 용어·계약 보류 및 미분류 번역 누락은 발견하지 못했다. 중국어 잔여는 腾讯 고유명, README의 单二进制 문서 참조, pg容器名 명령 인수 예시 원문 3곳뿐이다. 앞선 소비부·관련 테스트 확인 결과와 맞춰 직접 연결된 기대값 추가 수정 필요는 발견하지 못했으며 테스트는 수정·실행하지 않았다. U14 상태는 번역·정적 검토 완료이며 실행/동작 검증 완료가 아니다. 원문 저장 설명 문제는 후속 기술 검토 항목이다.
- **보존·실행·산출물**: 기존 U7 후속 3건과 모든 단위 기록·TRANSLATION_PROMPT 변경·미추적 담당표·그 밖의 파일·Git index를 보존했다. 스크립트 실행·설치·빌드·테스트·lint·포맷터·Docker/compose·CI·앱·Git 훅과 실행 검증은 **사용자 요청으로 미실행**이다. 누적 전체 diff `/tmp/artex-u14-final.diff`, 갱신된 상세 보고서 `/tmp/artex-u14-review.md`. fetch·병합·stage·commit·push·다음 단위는 진행하지 않았다.


### U15 · 첫 묶음(API 정찰 경계·준비 규칙) [부분 번역·정적 검토 완료 / 용어·제목 연결 검토 보류 / 실행 미검증]
- **범위·시작 상태**: 로컬 HEAD `a537cbc6005c9a7fafdcf1b3fc82d2ece5183fa9` 기준으로 `skills/api-recon/SKILL.md` 파일 시작부터 `## 执行路线图` 직전(1–117행), `skills/api-recon/reference.md` 시작부터 A절 직전(1–22행)만 처리했다. 대상 파일 기존 변경·staged 변경·MERGE_HEAD는 없었고 fetch·병합하지 않았다.
- **번역·잔여**: 모델용 자연어/description 46행과 참조 안내 11행을 번역했다. 미등록 전문 용어가 있는 문장/표 행 21개와 범위 밖 fragment가 연결된 제목 3개(1개 중복), 총 23행은 전체 원문 유지했다. 改包攻击는 사전 보고의 표 구분선 오인용을 실제 SKILL.md:20의 목표/금지 표 행 전체로 정정했다. 후보·원문·위치·문맥은 상세 보고서에 모았으며 용어집은 수정하지 않았다.
- **정책·링크 보류**: SKILL.md:21–22,41–43의 실제 인증/세션 의존 금지와 reference.md:288–297의 hardened 대상 조건부 실제 세션 안내를 앞뒤 문맥과 대조했다. 허용 범위의 번역 가능한 문장은 원문 강도로 번역했고, 정책 우선순위/예외를 새로 만들지 않았다. reference의 안내는 범위 밖 원문 유지했다. SKILL.md:12·80·106은 범위 밖 링크 의존으로 제목 원문 유지했다.
- **정적 보존**: 전체 줄 수 SKILL.md 425/425, reference.md 500/500; LF/CRLF(500개), 줄별 들여쓰기·빈 줄·영어/숫자 순서·inline code·강조·표·링크·name 키/값과 description 키를 유지했다. SKILL.md:118 이후/reference.md:23 이후는 바이트 동일하다. 금지/필수/권고/선택 조건·완료 정의·단계 순서를 대조했으며 새 출력 언어 지시는 추가하지 않았다.
- **관련 테스트·범위 밖**: 번들 본문을 직접 검사하는 테스트 기대값은 확인하지 못해 테스트 수정은 없다. U11의 skill_upload_test.go 압축 방식/암호화 오류 기대값 불일치 두 곳은 별도 공유 항목으로 유지하며 수정·해결 완료로 표시하지 않는다. ScopeSentry frontmatter, U13 및 다른 파일, 기존 U7 후속 보류와 모든 단위 이력은 보존했다. U15 나머지 범위가 남으므로 전체 완료가 아니다.
- **보존·실행·산출물**: 기존 TRANSLATION_PROMPT 변경·미추적 담당표·Git index는 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-u15-api-boundary.diff`, 상세 보고서 `/tmp/artex-u15-api-boundary-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 첫 묶음 확정 용어·제목 링크 보류 해소 [첫 묶음 번역·정적 검토 완료 / 정책 충돌 후속 검토 보류 / 실행 미검증]
- **용어 등록**: 사용자 확정 18개를 중복·충돌 확인 후 용어집의 U15 API 정찰 문맥 항목으로 등록했다. 壳层=로그인 후 기본 화면·进壳=로그인 후 화면 진입은 mock 표시이며 실제 인증 성공을 뜻하지 않는다. 业务码=애플리케이션 응답 코드는 HTTP 상태 코드와 구분한다. 다른 문맥의 용어와 코드·키·실제 값·계약 문자열은 보존했다.
- **번역 보류 해소**: SKILL.md 실행路线图 직전(1–117행) 및 reference.md A절 직전(1–22행)의 기존 원문 보류 23행을 모두 번역했다. 12·80·106행 제목도 지정 표기로 맞췄다. 이 범위의 중국어/지정 전각 문장부호 잔여·추가 미등록 용어·번역 보류는 발견하지 못했다. 이전 부분 번역 기록은 당시 이력으로 유지한다.
- **연결부 예외**: SKILL.md:160의 2개·178의 1개·281의 1개·425의 1개, 총 5개 fragment만 새 제목에 맞췄다. 기존 앵커의 소문자·문장부호 제거·공백 하이픈 규칙과 새 제목을 정적으로 대조했고 중복 제목 앵커가 없음을 확인했다. 지정 행의 링크 표시 문구·주변 중국어/영어·다른 목적지는 변경하지 않았다. 추가 연결부는 검색에서 발견하지 못했다. 렌더 실행은 하지 않았다.
- **정책 후속 보류**: SKILL.md의 실제 인증/백엔드 세션 의존 금지와 reference.md:290–293의 hardened 대상 조건부 실제 세션 사용 안내는 원문 정책 차이로 유지했다. 각 금지·허용 조건·기본 Phase 3 흐름·지시 강도를 번역했으며 우선순위·예외·새 설명/출력 언어 지시를 추가하지 않았다. 범위 밖 reference 안내는 원문 그대로다. 이 정책 확인은 용어·번역 보류 해소와 별개다.
- **최종 정적 보존**: SKILL.md 425/425줄·LF, reference.md 500/500줄·CRLF 500개와 줄 배치·들여쓰기·빈 줄·마지막 개행·Markdown 구조·명령·inline code·기존 영어/숫자 순서를 대조했다. micro-frontend만 사용자가 확정한 중국어 번역으로 추가됐고 기존 영어는 동일하다. 원래 허용 절 이후는 지정 4행의 fragment 변경만 있다. 기존 첫 묶음 번역은 유지했다.
- **잔여·범위·보존**: 첫 묶음의 번역·제목 연결 보류는 해소했지만 U15 나머지는 미작업이므로 전체 완료가 아니다. ScopeSentry frontmatter·U11 테스트 기대값 불일치 2건은 범위 밖 공유/확인 항목으로 유지하며 수정하지 않았다. 기존 U7 후속 보류·모든 팀원 이력·TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 4개 파일 전체 diff `/tmp/artex-u15-api-boundary-final.diff`, 검토 보고서 `/tmp/artex-u15-api-boundary-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 두 번째 묶음(API 정찰 정적 수집·파라미터 추적·클라이언트 조건) [부분 번역·정적 검토 완료 / 용어 검토·정책 충돌 후속 보류 / 실행 미검증]
- **범위·시작 상태**: 로컬 `work/ko-translation`, HEAD `70d6125921aa169076b3d5fd96ed4822fe082a64` 기준. 첫 묶음 커밋과 확정 용어·제목/fragment 반영을 확인했다. 시작 시 대상 파일의 기존 변경·staged 변경·진행 중 병합은 없었다. `SKILL.md` 실행 로드맵부터 Phase 3 직전(118–278행), `reference.md` A/B절(23–198행)·E절(252–262행)·J절(388–500행), 두 Python 파일의 모듈 docstring만 처리했다.
- **파일별 번역**: `SKILL.md` 자연어 64행, `reference.md` 자연어/설명용 코드 주석 68행, `harvest_static.py` 모듈 docstring 4행, `spider_mpa.py` 모듈 docstring 1행을 번역했다. text 흐름도의 자연어만 번역했고 명령·JSON/정규식·예시 값·영어·화살표·줄 배치를 보존했다. Phase 0/1 제목과 같은 허용 범위의 fragment 2개만 쌍으로 변경했다. 첫 묶음의 fragment와 범위 밖 제목 연결은 유지했다.
- **용어 보류·잔여**: 미등록 전문 용어가 있는 문장/표 행/제목 32개(SKILL 16·reference 16)는 전체 원문 유지했다. 参数逆向·鉴权三道门/三门·渲染门/拦截器门/内容门·绑定层/绑定源·包装层·置信度·静态兜底·校验门·发起程序·待触发/联动 select의 원문 전체·위치·문맥·권장안을 상세 보고서에 모았다. 용어 보류는 DNT·완료로 처리하지 않는다. reference:52의 중국어 로그인 검색 정규식과 SKILL:136/138/139의 범위 밖 Phase 3/4/5 제목 fragment는 계약/연결 보존으로 별도 분류했다. 추가 미분류 누락은 발견하지 못했다.
- **생성부·소비부·테스트**: harvest의 HTML→manifest→chunk→정적 API/라우트 파일, spider의 form/동일 출처 링크/inline API 수집과 config의 runtime/preload 소비를 읽기 전용으로 대조했다. 모듈 docstring의 `__doc__`/getdoc/pydoc 비교·파싱 의존과 변경된 번들 안내를 직접 검사하는 테스트 기대값은 검색에서 발견하지 못해 테스트를 수정하지 않았다. SDK skill 로더 내부·외부 소비자는 미확인이다. U11의 skill_upload_test.go 압축 방식·암호화 오류 기대값 불일치 두 곳은 기존 공유 항목으로 유지한다.
- **정적 보존**: 원본/현재 줄 수 SKILL 425/425·reference 500/500·harvest 207/207·spider 123/123. LF/CRLF·들여쓰기·빈 줄·개행·Markdown 구조·코드펜스·inline code·숫자·기존 영어 순서·이스케이프를 대조했다. 추가 micro-frontend는 확정된 중국어 용어 적용이다. 허용 절 밖 문서와 Python 모듈 docstring 밖 모든 코드/주석/문자열은 바이트 동일하다. 기존 진행 기록은 변경 없이 보존하고 이번 기록만 추가했다.
- **후속·상태·보존**: 로그인/세션 금지와 reference G절의 조건부 실제 세션 안내 간 원문 정책 차이를 보존했다. 정책 우선순위·예외·새 출력 언어 지시는 추가하지 않았다. ScopeSentry frontmatter·U11 공유 항목·U7 후속 3건·다른 단위 이력은 유지한다. 이번 묶음은 용어 보류가 남은 부분 번역이며 U15 나머지도 미작업이므로 U15 전체 완료가 아니다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index는 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-u15-api-static-params.diff`, 상세 검토·보류 보고서 `/tmp/artex-u15-api-static-params-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 두 번째 묶음 확정 용어·제목 링크 보류 해소 [두 번째 묶음 번역·정적 검토 완료 / 정책 충돌 후속 검토 보류 / 실행 미검증]
- **용어·보류 해소**: API 정찰 문맥의 사용자 확정 용어 15개를 중복·충돌 확인 후 `[확정]`으로 등록했다. 静态兜底는 런타임에서 진행할 수 없어 정적 분석으로 대신 확인하는 경우, 发起程序는 DevTools의 요청 시작 지점·호출 스택 확인으로 한정했다. 코드·JSON 키 confidence와 다른 문맥의 기존 용어는 보존했다. 앞선 H01~H32(SKILL 16·reference 16) 전체 문장/표 행/제목을 번역했으며 이전 보류 기록은 당시 이력으로 유지한다.
- **추가 표현·연결부**: SKILL.md:191을 “Hook으로 암호화 함수의 입력 파라미터 관찰”로 정리했다. 변경 제목 13개(SKILL 4·reference 9)의 연결부를 skills/ 전체에서 검색했다. Phase 1b:176·Phase 2:249 제목과 :134·135의 fragment를 쌍으로 변경했으며 추가 연결 fragment는 발견하지 못했다. 첫 묶음 링크와 아직 번역하지 않는 Phase 3·4·5 제목/fragment는 유지했다. Markdown 앵커를 정적으로 대조했고 렌더는 실행하지 않았다.
- **최종 정적 보존**: 허용 범위는 SKILL 118–278행, reference A/B 23–198·E 252–262·J 388–500행과 두 Python 모듈 docstring이다. 누적 변경은 SKILL 80행·reference 84행·harvest docstring 4행·spider docstring 1행이다. 원본/현재 줄 수 425/425·500/500·207/207·123/123, reference CRLF 500개와 나머지 LF·들여쓰기·개행·Markdown 구조·화살표·명령·설정 예시·기존 영어·숫자·inline code를 유지했다. 이번 용어 적용으로 추가된 Initiator는 확정 표기이며 기존 영어 변경이 아니다. 허용 범위 밖 문서와 Python docstring 밖 코드는 바이트 동일하다.
- **잔여·테스트**: 이번 범위의 용어 보류 32곳은 모두 해소했으며 추가 용어/계약 보류·미분류 번역 누락은 발견하지 못했다. reference:52의 중국어 grep 검색 정규식은 코드 DNT, SKILL:136/138/139의 중국어 fragment는 범위 밖 제목 연결 보존으로 분류한다. 본문/docstring의 직접 비교·파싱·테스트 기대값 의존은 검색에서 발견하지 못해 테스트 수정은 없다. 외부 SDK 소비자·실제 모델 해석은 미확인이다.
- **후속·보존**: 로그인/세션 금지와 reference G절의 조건부 실제 세션 안내 간 원문 정책 충돌, U11 기대값 불일치 공유 2건, ScopeSentry frontmatter 확인 항목, U7 후속 보류와 모든 기존 기록을 유지한다. 정책 우선순위·예외·새 출력 언어 지시는 추가하지 않았다. U15 나머지가 남으므로 전체 완료가 아니다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 6개 파일 전체 diff `/tmp/artex-u15-api-static-params.diff`, 갱신된 보고서 `/tmp/artex-u15-api-static-params-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 세 번째 묶음(API 정찰 런타임·권한 트리·보고서) [부분 번역·정적 검토 완료 / 용어·원문 정책/설명 후속 보류 / 실행 미검증]
- **범위·시작 상태**: 로컬 `work/ko-translation`, HEAD `19c3ceb4813588cadc84a3afa451a78b92f9b870`의 두 번째 묶음 커밋을 확인했다. 대상 파일의 기존 변경·staged 변경·진행 중 병합은 없었다. `SKILL.md` Phase 3부터 끝(279–425행), `reference.md` C/D 199–251·F/G/H/I 263–387행, preload.js·runtime_harvest.js·build_perm_tree.py·extract_route_map.py 전체를 검토했다. 스크립트는 중국어 주석만 번역했다.
- **파일별 변경·연결부**: SKILL 67행의 자연어와 앞부분 136·138·139행의 Phase 3/4/5 fragment 3개, reference 82행, preload 중국어 주석 14행, runtime 중국어 주석 2행을 번역했다. 두 Python 파일은 전체 영어이므로 수정하지 않았다. 변경 제목 22개와 skills/ 전체의 fragment를 대조했으며 추가 연결부는 발견하지 못했다. 앞 묶음은 허용된 fragment 3곳 외 바이트 동일하다.
- **용어·계약·잔여**: 미등록 负向修正, L3 stub 문맥 兜底, 中和, 反调试가 포함된 7개 문장/표 행/주석은 전체 원문 유지했다(SKILL 290·291, reference 244·248, preload 29·82·122). 원문 전체·위치·문맥·후보·근거는 상세 보고서에 모았다. preload:192의 중국어 NEGATIVE_RE 및 완료된 reference:52의 grep 검색 정규식은 코드 DNT로 보존했다. 중국어 실행/표시 문자열 수정 후보는 발견하지 못했다. 명령·JSON/config·예시·실제 출력은 변환하지 않았다. 미분류 번역 누락·새 계약 불명확 항목은 발견하지 못했으며 용어집은 수정하지 않았다.
- **정적 보존·테스트**: 원본/현재 줄 수 SKILL 425/425·reference 500/500·preload 348/348·runtime 228/228·build_perm_tree 210/210·extract_route_map 42/42. reference/preload/두 Python의 CRLF와 SKILL/runtime의 LF, 들여쓰기·개행·Markdown 구조·기존 영어·숫자·포맷·이스케이프·명령·정규식·설정 키/값을 대조했다. 스크립트 주석 밖 코드는 바이트 동일하고 reference 완료 A/B/E/J절도 동일하다. 관련 생성/소비 경로와 이름이 다른 서버 skill 테스트를 읽기 전용으로 확인했으며 변경한 자연어를 직접 비교하는 기대값은 검색에서 발견하지 못했다. 테스트 수정은 없고 SDK/실제 모델 해석은 미확인이다.
- **원문 문제·후속**: 실제 인증/세션 의존 금지와 reference G절의 조건부 실제 세션 안내·I6의 실제 session 비교 지시를 원문 의미대로 유지했다. 정책 우선순위나 예외를 추가하지 않았다. reference:267의 D절 정규식 참조(실제 안내는 E절), :350의 routes 확장 설명(build_perm_tree.py:196–206은 stubs만 갱신), depth/coverage 공통 안내와 preload의 forward=true 시 L1/L3 mock 미적용 경로 차이는 원문 설명 후속 검토로 기록했다. 번역에서 코드/안내를 임의 수정하지 않았다.
- **상태·보존·산출물**: ScopeSentry frontmatter·U11 공유 2건·U7 후속 보류·모든 이전 기록·기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존했다. 용어 보류와 U15 나머지 범위가 남으므로 U15 전체 완료가 아니다. 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-u15-api-runtime.diff`, 상세 보고서 `/tmp/artex-u15-api-runtime-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 세 번째 묶음 확정 용어 보류 해소 [세 번째 묶음 번역·정적 검토 완료 / 원문 정책·설명 후속 검토 보류 / 실행 미검증]
- **용어·번역**: 최신 용어집의 중복·충돌을 확인하고 负向修正=실패 응답 보정(L2), 兜底=기본 응답 처리(L3), 中和=무력화(클라이언트 로그인 이동·라우트 가드 억제), 反调试=안티 디버깅을 `[확정]`으로 등록했다. 다른 兜底 문맥과 서버 인증·인가 계약은 변경하지 않았다. H01~H07의 표 행·주석 7곳을 모두 번역했으며 H01은 지정 표기 그대로다. 앞선 부분 번역 기록은 이력으로 보존한다.
- **최종 정적 검토**: 이번 추가 소스 변경은 SKILL:290·291, reference:244·248, preload:29·82·122의 7행뿐이다. 누적 변경은 SKILL 72행(자연어 69·fragment 3), reference 84행, preload 주석 17행, runtime 주석 2행이다. 원본/현재 줄 수 425/425·500/500·348/348·228/228·210/210·42/42와 LF/CRLF, 줄 배치·들여쓰기·기존 영어·포맷·이스케이프·Markdown·명령·키·설정값·정규식은 동일하다. 기존 제목/fragment와 완료된 앞 묶음, 스크립트 주석 밖 코드는 유지했다.
- **잔여·테스트**: 이번 묶음의 용어 보류 7곳은 해소했고 추가 용어/계약 보류·미분류 번역 누락은 발견하지 못했다. preload:192의 중국어 NEGATIVE_RE와 완료된 reference:52의 로그인 검색 정규식은 코드 DNT로 보존한다. 직접 연결된 테스트 기대값 변경 필요는 발견하지 못했고 테스트는 수정·실행하지 않았다. 외부 SDK와 실제 모델 동작은 미확인이다.
- **후속·보존**: 로그인·세션 원문 정책 충돌, D/E절 참조 문제, routes 확장 설명 문제, depth/coverage 동작 차이는 후속 항목으로 유지하며 이번에 수정하지 않았다. U11 공유 2건·U7 후속 보류·ScopeSentry 확인 항목·기존 기록·로컬 변경·미추적 담당표·Git index를 보존한다. U15 전체 완료는 아니다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 전체 diff `/tmp/artex-u15-api-runtime.diff`, 갱신된 보고서 `/tmp/artex-u15-api-runtime-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 네 번째 묶음(ScopeSentry 문서) [부분 번역·정적 검토 완료 / 용어·원문 형식 검토 보류 / 실행 미검증]
- **범위·시작 상태**: `work/ko-translation`, HEAD `41efa3c448667c1e15291b0271f9745c2376e2f4`의 세 번째 API 정찰 묶음 커밋을 확인했다. 실제 추적 파일 `skills/scopesentry/SKILL.md` 전체(1–366행)를 검토했다. 대상 파일 기존 변경·staged 변경·진행 중 병합은 없었다. fetch·병합하지 않았다.
- **번역·링크**: description 값과 모델용 설명·조건·금지·선택 사항·표·제목 등 141행을 변경했다. 216·267행 제목과 이를 가리키는 192·194행의 fragment를 함께 변경했다. skills/ 내 추가 직접 참조는 검색에서 발견하지 못했다. 모든 JSON/mermaid 코드 블록과 inline code·작업명/검색값 쌍은 원문 보존했다. 模糊匹配는 용어집의 자산 DSL 문맥에 따라 부분 일치로 적용하되 원문의 regex 설명·연산자·인덱스 조건을 유지했다.
- **용어 보류·잔여**: 미등록 流水线·分布式·子域名接管·资产测绘가 포함된 59·143·145·152·159행의 문장/표 행 5개는 전체 원문 유지했다. 원문 전체·문맥·권장안은 상세 보고서에 모았으며 용어집은 수정하지 않았다. 중국어 URL/키 placeholder·작업명·DSL 매칭값·mermaid 예시는 코드/원문 예시 보존으로 별도 분류했다. 추가 미분류 번역 누락·새 계약 보류는 발견하지 못했다.
- **frontmatter 원문 문제**: 3행의 `## name: scopesentry-mcp`와 1–44행 구분선 안의 본문/JSON 배치를 확인하고 그대로 유지했다. 키·식별 값·구분선은 수정하지 않았다. 서버의 `skill.LoadDir`·디렉터리명 기준 필터·Skill 등록 경로를 정적으로 확인했지만 외부 SDK 로더 내부/실제 로드 결과는 미확인이다. 표준 YAML 파싱 시 형식 문제가 생길 가능성은 추론이며 로드 성공/실패를 확정하지 않는다.
- **정적 보존·테스트**: SKILL 원본/현재 366/366줄·LF 366개·CRLF 0개, 들여쓰기·빈 줄·마지막 개행·Markdown 구조·숫자·기존 영어 순서·inline code·코드 블록 바이트가 동일하다. 내부 fragment 두 곳은 새 제목과 정적으로 대조했다. MCP 동기화 소비부의 도구명·파라미터/응답 키와 관련 이름이 다른 skill 테스트를 확인했으며 이번 문구를 직접 검사하는 기대값 변경 필요는 발견하지 못했다. 테스트·외부 MCP를 실행/접속하지 않았다.
- **후속·보존**: API 정찰 로그인/세션 정책 충돌·D/E절 참조·routes 확장 설명·depth/coverage 동작 차이, U11 공유 2건·U7 후속 보류·기존 기록은 유지한다. playwright-cli와 다른 단위는 이번 범위 밖이며 U15 전체 완료가 아니다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·완료된 API 정찰 파일·다른 파일·Git index를 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-u15-scopesentry.diff`, 상세 검토·보류 보고서 `/tmp/artex-u15-scopesentry-review.md`. stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 네 번째 묶음 확정 용어·Mermaid 자연어 보류 해소 [ScopeSentry 번역·정적 검토 완료 / 원문 형식 후속 검토 보류 / 실행 미검증]
- **용어·보류 해소**: 최신 용어집에서 중복·충돌이 없음을 확인하고 流水线=파이프라인·分布式=분산 처리·子域名接管=서브도메인 탈취·资产测绘=자산 매핑을 ScopeSentry 문맥의 `[확정]`으로 등록했다. 子域名接管은 취약점·검사 기능 설명에 한정하며 외부 모듈 구현·실제 탈취 동작은 확인하지 않았다. H01~H05의 59·143·145·152·159행을 모두 번역했고 원문의 강조·금지·권고·단계·같은 노드 실행 및 부하/속도/오류 설명을 유지했다. 앞선 보류 기록은 당시 이력으로 보존한다.
- **추가 허용·정적 보존**: Mermaid 166–169행의 사람이 읽는 중국어 노드 라벨만 번역했다. `task==阶段1任务名`·노드 ID·화살표·배치·구문·도구명·DSL·실제 검색값은 동일하다. 이번 추가 SKILL 변경은 9행, 누적 변경은 150행이다. 366/366줄·LF 366개·CRLF 0개·들여쓰기·빈 줄·개행·강조·Markdown/Mermaid 구조·JSON·inline code·기존 영어·숫자 순서·기존 제목/fragment 대응을 보존했다.
- **잔여·후속·상태**: 이번 문서의 용어 보류 5곳은 해소했고 추가 번역 누락·용어/계약 보류는 발견하지 못했다. 남은 중국어는 작업명·규칙명·ObjectID·URL/헤더/인증키 placeholder·JSON 예시·DSL 매칭값으로 분류했다. frontmatter 기존 형식과 MCP 이름, 로더 영향 미확인·API 정찰 정책/설명 후속·U11 공유·U7 후속 보류는 유지한다. 관련 테스트 수정 필요는 발견하지 못했으며 실제 모델/외부 서비스 동작은 미확인이다. U15 전체 완료가 아니다.
- **보존·실행·산출물**: 완료된 다른 파일·기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index와 이전 진행 기록을 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. 누적 전체 diff `/tmp/artex-u15-scopesentry.diff`, 갱신된 보고서 `/tmp/artex-u15-scopesentry-review.md`. fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### U15 · 마지막 묶음(playwright-cli) 및 전체 종합 [전체 파일 번역·정적 검토 완료 / 원문 정책·형식·설명 후속 검토 보류 / 실행 미검증]
- **시작·커버리지**: `work/ko-translation`, HEAD `729e4f8c9f012c7c24fe1a6a30a1654454589b37`의 ScopeSentry 커밋을 확인했다. 대상 기존 변경·staged 변경·MERGE_HEAD는 없었다. git의 skills/ 추적 파일 21개(API 정찰 10·ScopeSentry 1·playwright-cli 10)를 모두 읽고 기존 묶음 기록과 대조했다. 기존 21개와 차이가 없으며 미검토 추적 파일은 없다. fetch·병합하지 않았다.
- **playwright-cli**: SKILL.md 및 references 9개 전체는 영어 설명·명령·옵션·예시뿐이며 중국어 자연어가 없어 모두 “검토 완료 / 번역 대상 없음”이다. 영어·코드·식별자·메타데이터·제목·링크를 수정하지 않았다. 아래 목록은 파일별 최종 검토 상태이다.
| 추적 파일 | HEAD/현재 줄 수 | 개행 | 검토·번역 상태 | 잔여 근거 |
| --- | --- | --- | --- | --- |
| `skills/api-recon/SKILL.md` | 425/425 | LF | 기존 번역 완료 / 전체 재검토 완료 | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/reference.md` | 500/500 | CRLF | 기존 번역 완료 / 전체 재검토 완료 | 52행 검색 정규식 DNT |
| `skills/api-recon/scripts/build_perm_tree.py` | 210/210 | CRLF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/extract_route_map.py` | 42/42 | CRLF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/harvest_static.py` | 207/207 | LF | 기존 번역 완료 / 전체 재검토 완료 | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/package-lock.json` | 1011/1011 | LF | 검토 완료 / 번역 대상 없음 / 메타데이터 DNT | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/package.json` | 9/9 | LF | 검토 완료 / 번역 대상 없음 / 메타데이터 DNT | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/preload.js` | 348/348 | CRLF | 기존 번역 완료 / 전체 재검토 완료 | 192행 판별 정규식 DNT |
| `skills/api-recon/scripts/runtime_harvest.js` | 228/228 | LF | 기존 번역 완료 / 전체 재검토 완료 | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/api-recon/scripts/spider_mpa.py` | 123/123 | LF | 기존 번역 완료 / 전체 재검토 완료 | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/SKILL.md` | 420/420 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/element-attributes.md` | 23/23 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/playwright-tests.md` | 39/39 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/request-mocking.md` | 87/87 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/running-code.md` | 241/241 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/session-management.md` | 225/225 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/storage-state.md` | 275/275 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/test-generation.md` | 433/433 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/tracing.md` | 139/139 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/playwright-cli/references/video-recording.md` | 143/143 | LF | 검토 완료 / 번역 대상 없음(영어 유지) | 중국어·지정 전각 문장부호 없음; 전체 문맥 확인 |
| `skills/scopesentry/SKILL.md` | 366/366 | LF | 기존 번역 완료 / 전체 재검토 완료 | 21행의 DSL/검색값·원문 예시/placeholder(아래 전부 분류) |
- **잔여·연결부**: 중국어/지정 전각 잔여는 reference:52·preload:192의 검색/판별 정규식과 ScopeSentry의 21행 DSL/검색값·원문 예시·placeholder, 총 23행이다. 번역 누락·미확정 용어·새 계약 보류는 발견하지 못했다. 문서의 코드 블록 밖 Markdown 링크 33곳의 실제 파일/제목 fragment 대응을 정적으로 확인했다. 코드 안 정규식은 링크로 취급하지 않았으며 렌더는 실행하지 않았다.
- **후속 미해결**: API 정찰 로그인/세션 정책 충돌·reference D/E절 참조·routes 확장 설명·depth/coverage stub 차이·ScopeSentry frontmatter 형식과 로더 영향 미확인·U11 skill_upload_test 기대값 불일치 2곳은 유지한다. U7의 live 중국어 동사 의존·옛 공통 블록 중복·기존 커스텀 정책 본문 보존은 별도 단위 보류다. 추가로 package.json의 api-recon-runtime과 lock 루트 spa-api-recon-runtime 이름 차이, playwright test-generation:293의 병렬 금지와 :352의 병렬 가능 안내는 원문 메타데이터/설명 후속 확인 항목으로 기록하며 수정하지 않는다.
- **정적 보존·실행**: skills/ 21개 파일은 모두 HEAD 및 시작 상태와 바이트 동일하다. 줄 수·LF/CRLF·들여쓰기·개행·Markdown·Mermaid·코드·JSON·DSL·명령·옵션·설정값·영어를 유지했다. PROMPT_PROGRESS의 이전 이력은 그대로 두고 이번 기록만 추가했다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 해시로 대조해 보존했다. 관련 테스트 기대값 변경 필요는 발견하지 못했고 테스트는 수정하지 않았다. 설치·빌드·테스트·lint·포맷터·스크립트·브라우저·앱·Git 훅·외부 MCP/API와 실행 검증은 **사용자 요청으로 미실행**이다. 실제 로더/모델/외부 서비스 동작 검증 완료를 뜻하지 않는다.
- **산출물·종료**: 전체 diff `/tmp/artex-u15-final-audit.diff`, 파일별 잔여·후속 종합 보고서 `/tmp/artex-u15-final-review.md`. fetch·병합·stage·commit·push·다음 단위는 진행하지 않았다.


### U15 · API 정찰 문서 후속 2건 정정 [문서 정정·정적 검토 완료 / 나머지 후속 검토 보류 / 실행 미검증]
- **확인 기준**: `work/ko-translation`, HEAD `79f71c6dbbfe2a49dbd673ef71d46ee548fb63dc`. 대상 기존 변경·staged 변경·진행 중 병합 없음. 최신 로컬 문서와 관련 코드를 읽기 전용으로 대조했다.
- **D/E절 참조 정정**: reference.md:267의 “정적 API가 매우 적음” 행에서 D절을 E절로 정정했다. 실제 Endpoint 추출 정규 표현식은 :252–262의 E절이며 D절(:234–248)은 Hook 기능 설명이다.
- **routes 자동 확장 주장 정정**: reference.md:350을 “`--config`로 지정한 `config.json`이 존재하면 `stubs`를 자동 갱신”으로 정정했다. build_perm_tree.py:155의 옵션 선언 및 :196–206의 설정 저장 경로를 확인했다. 기존 permissions/all·role_permissions stub을 제거하고 새 stub을 추가해 cfg['stubs']만 갱신하며 routes 변경은 없다. --config 생략 시 outdir/config.json을 선택하는 기존 구현도 보존한다.
- **보존·잔여**: reference 500/500줄·CRLF 500개·나머지 본문·Markdown·명령·코드·예시와 기존 번역을 보존했다. 이전 후속 기록은 이력으로 유지하며 위 두 건만 해소했다. 로그인/세션 정책 충돌·depth/coverage 차이·ScopeSentry 형식/로더·U11 공유·U7 후속 보류 및 다른 단위 상태는 변경하지 않았다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. `/tmp/artex-u15-doc-followup.diff`, `/tmp/artex-u15-doc-followup-review.md`. fetch·병합·stage·commit·push·다음 작업은 진행하지 않았다.


### U15 · ScopeSentry frontmatter 형식 정정 [형식 정정·정적 검토 완료 / 실제 SDK 로드·모델 동작 실행 미검증]
- **버전·근거**: go.mod의 replace 없음 및 norma v0.4.3 고정을 확인했다. Go 모듈 프록시 ZIP의 dirhash가 go.sum과 일치한다. SDK skill/skill.go:263–289는 YAML 성공 여부와 무관하게 먼저 본문을 분리하고, :231–234는 name 누락 시 디렉터리 이름을 fallback으로 쓴다. 기존 문서의 유효 등록 이름 scopesentry를 유지하는 수정이다.
- **참조 확인**: 저장소 내 scopesentry/scopesentry-mcp 참조를 읽기 전용으로 검색했다. scopesentry-mcp의 기존 SKILL 헤더 외 직접 skill 호출·하드코딩된 조회 참조는 발견하지 못했다(과거 진행 기록은 이력). server/assembly.go:35–51의 노출 키와 server/skill_usage.go:64–75의 사용 기록 키는 Dir basename scopesentry다. web/src/lib/mock/data.ts의 skill 예시 name도 scopesentry다. MCP 서버/동기화 이름은 다른 역할로 보존했다. 실제 DB/모델 호출은 미확인이다.
- **수정·보존**: SKILL.md:3을 name: scopesentry로 바꾸고 description 바로 뒤(:5)에 종료 ---를 추가했다. frontmatter에는 빈 줄 외 name/description만 있다. 기존 description 이후 본문 전체를 바이트 그대로 보존해 준비·연결·JSON 안내부터 파일 끝까지 본문에 포함된다. 기존 :44 구분선은 :45의 본문 절 구분선으로 유지했다. 366→367줄, LF 유지; JSON·명령·DSL·검색값·링크·MCP 이름·본문 제목·기존 번역/들여쓰기 동일.
- **후속·범위**: 원래 frontmatter 형식 문제는 정정했다. 실제 SDK 로드·모델 본문 전달은 실행 미검증이다. db/db.go:306–307의 mcps: ScopeSentry 선언 설명과 실제 문서의 mcps 필드 부재는 추가 원문 설명 불일치로 기록하며 mcps/코드/권한을 추가하지 않았다. 로그인·세션 정책·depth/coverage·메타데이터/병렬 안내·U11·U7 등 다른 후속 항목 및 이전 이력은 유지한다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff /tmp/artex-u15-scopesentry-frontmatter.diff, 검토 보고 /tmp/artex-u15-scopesentry-frontmatter-review.md. 기존 로컬 변경·미추적 담당표·다른 파일·Git index 보존. fetch·병합·stage·commit·push·다음 작업 없음.


### 통합 감사 보완 1 · A001–A010 [기대값 수정·정적 대조 완료 / 실행 미검증]
- **기준·범위**: HEAD `474300b6d355fc1c8cf7051d644c52a7abae1fd4` 및 통합 감사 A001–A010을 현재 생성부와 재대조했다. 지정된 테스트 7개에서 직접 연결된 기대 문자열 10곳만 수정했다. 감사 기준 이후 대상 내용 변경은 없었으며 이미 일치하는 항목은 없었다. 운영 소스·DB·프런트·용어집은 수정하지 않았다.
- **대응 결과**: A001 finding_recorder의 미등록 오류, A002 finding_traffic의 작업 접근 거부, A003 finding_workflow의 독립 취약점 기록 ID, A004 traffic_search 설명의 호스트 단독/호스트:포트/완전한 URL 지원, A005 intercept_detail의 삭제된 대화, A006 task_categories의 task_ids 개수 제약, A007–A008 chat_mentions의 서버 생성 레코드 스냅샷 헤더, A009–A010 skill_upload의 지원하지 않는 압축 방식/암호화 오류를 현재 한국어 출력과 맞췄다. A009–A010의 기존 U15 공유 기대값 보류는 이번 정적 대응으로 해소했으며 실행 결과는 미확인이다.
- **보존·정적 검토**: Contains 검사와 조건·HTTP 상태 코드·Deflate64·worker-hidden-proof·original proof·입력·fixture·호출·반환 구조·무관한 기대값·실패 안내를 보존했다. 테스트 7개의 줄 수·들여쓰기·개행·포맷 지정자는 동일하며 기대 문자열 10곳 외 바이트는 변경하지 않았다. 기존 이력·U7/U15의 다른 후속 보류·다른 단위 상태는 유지한다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index를 보존했다. A011 이후 감사 항목은 이번 작업 범위 밖이다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 테스트 통과 또는 실패를 관측하지 않았다. 전체 diff `/tmp/artex-audit-fix-tests.diff`, 검토 보고 `/tmp/artex-audit-fix-tests-review.md`. 기존 감사 보고서는 유지하고 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 2 · A011–A024·A075–A079 [부분 번역·정적 검토 완료 / 미등록 용어 보류 / 실행 미검증]
- **기준·범위**: HEAD `e7f74ef9ea9f82d2a4ae85ce1c34e1d8fc0cbfda`에서 직전 A001–A010 보완 커밋을 확인했다. 대상 기존 변경·staged 변경·MERGE_HEAD는 없었다. 지정 테스트·Next 설정의 주석/개발자 실패 안내, sidequestion 문서 3개 산문/Mermaid 설명 라벨, schema.sql 주석 및 허용 표시 문자열 3곳, company_scope.go의 ICP 설명을 처리했다. A025–A074 및 A080 이후는 작업하지 않았다.
- **처리·보류**: 소스/문서 14개에서 341행을 번역하거나 지정 안내를 정리했다. schema.sql 주석165행 중 143행을 번역하고 22행은 전문 용어 보류로 원문 유지했다. 전체 용어 보류는 35개 문장/주석 블록이며 원문·위치·문맥·후보를 상세 보고서에 모았다. capture_usage_test.go의 迭代器协议 및 task_archive_package_test.go의 符号链接 주석은 미등록 전문 용어가 있어 두 파일은 수정하지 않았다. sidequestion의 深拷贝·滚动摘要 등과 schema의 半开试探·调用账本·租约·DB 호환 兜底 등의 문맥도 확정하지 않았다. 용어집은 변경하지 않았다.
- **직접 표시 문자열**: A076 `작업 실행 중 자동 연관`, A077 `과거 작업 자산 연관에서 마이그레이션됨`, A078 `재검증 대화가 삭제되었습니다`를 적용했다. SQL 비주석 변경은 이 세 문자열뿐이며 system/legacy/stopped·테이블/열·조건·트리거·마이그레이션 비교값을 보존했다. 기존 저장 문구의 소급 번역 UPDATE나 마이그레이션은 추가하지 않았다. A079는 설명 주석의 备案를 `ICP 등록(备案)`으로 맞추고 파서의 备案 및 실제 ICP 예시는 보존했다.
- **기대값·문서·형식**: 테스트 입력·fixture·검사 조건·기대값과 A001–A010 보완값을 모두 유지했다. skill_upload의 지정 실패 안내 두 곳은 `want a Korean message naming Deflate64`와 `want 암호화 hint`로 맞췄다. 파일별 줄 수·LF·들여쓰기·개행·포맷/인수 순서·이스케이프를 정적으로 대조했다. 문서의 명령 코드 블록·inline 예시·링크 목적지·날짜/숫자·기존 검증 결과를 유지했으며 제목14개의 직접 연결 fragment 참조는 검색에서 발견하지 못했다. 과거 PASS/실제 모델·브라우저 기록은 과거 이력의 번역이며 이번 검증 결과가 아니다. 원문 문제는 보고서에 별도로 기록하고 기능·정책을 변경하지 않았다.
- **보존·실행·산출물**: 기존 이력·다른 단위와 U7/U15 후속 보류, 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index·기존 감사 보고서를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. `/tmp/artex-audit-fix-omissions.diff`, `/tmp/artex-audit-fix-omissions-review.md`. 전체 감사/전체 번역 완료 또는 테스트 통과를 주장하지 않으며 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 2 · H01–H35 확정 용어 반영 [이번 범위 번역·정적 검토 완료 / 실행 미검증]
- **용어·해소**: 사용자 확정 표기 21개 항목(쌍 표기 포함)을 중복·충돌 확인 후 이번 문맥의 `[확정]`으로 등록했다. 기존 다른 兜底 문맥은 유지한다. H01–H35의 모든 보류 문장·주석 블록을 번역했으며 H10 원문에는 回源이 없어 이전 보류 분류를 정정하고 기존 확정 용어로 번역했다. 이전 부분 번역 기록은 당시 이력으로 보존한다.
- **지정 표현**: README의 기록 빈도는 “250 ms당 최대 한 번 기록”, VALIDATION의 과거 검증 설명은 “이번 검증의 통과 항목으로 처리하지 않았고” 및 “Agent 라벨과 이력”으로 정리했다. 과거 통과 기록은 이번 실행 결과가 아니다.
- **정적 검토·보존**: 원본/현재 줄 수·LF·들여쓰기·코드·포맷 지정자/인수 순서·이스케이프·문서 코드 블록/링크/예시를 대조했다. 테스트 입력·fixture·기대값·로직과 A001–A010 보완값을 유지했다. SQL 비주석 변경은 앞서 허용한 표시 문자열 3곳뿐이며 추가 SQL 표시값·계약·마이그레이션 변경은 없다. 이번 범위의 번역 누락·새 용어/계약 보류는 발견하지 못했다. 범위 밖 문구와 기존 U7/U15 후속 항목·다른 단위 상태는 유지하며 전체 감사/전체 번역 완료를 뜻하지 않는다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 기존 로컬 변경·미추적 담당표·다른 파일·Git index와 기존 감사 보고서를 보존했다. 누적 diff `/tmp/artex-audit-fix-omissions.diff`, 검토 보고 `/tmp/artex-audit-fix-omissions-review.md`. fetch·병합·stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 3 · A027–A045·A080–A088 [이번 범위 용어·표기·기대값 정적 대조 완료 / 실행 미검증]
- **기준·현재 확인**: work/ko-translation, HEAD `92ad8541e268be689e052b9b648789209af7807f`에서 직전 번역 누락 보완·확정 용어 커밋을 확인했다. 통합 감사 28개 항목은 모두 현재 소스에 남아 있어 문맥별로 수정했다. 대상 기존 변경·staged 변경·MERGE_HEAD는 없었으며 fetch·병합하지 않았다.
- **수정**: LLM failover의 순환 전환·회로 차단·재시도 대기 시간·안전 구간, 의도/탐색 그래프의 상위·하위, 동작·심각도 등급·planner·공격 표면·자산 조회의 엔드포인트와 지정 알림 delivery/전송 속도 제한 문맥을 맞췄다. 알림 dispatcher의 lease 설명은 db/notification.go·notification_delivery.go·notification_test.go·server/notifier.go에서 함께 정리했다. 일반 데이터 전달·실제 폴링·영어·상태값은 유지했다.
- **기대값·표기**: A041 안내를 직접 검사하는 sidequestion/sidequestion_test.go:180 및 agent/side_questions_test.go:107의 Contains 기대 문자열만 대응시켰다. db/notification_test.go의 lease 주석/실패 안내는 같은 확정 표기로 정리하고 입력·기대 상태·조건은 유지했다. A044는 두 JSX 경계에 명시적 공백을 추가했고 A045는 지정 산문/Mermaid 라벨/명령 주석의 전각 괄호만 정리했다. A080은 db/companies.go:636의 표시 목록 strings.Join 구분자만 ", "로 변경했다.
- **정적 보존·잔여**: 소스/문서/테스트 26개·96행 변경, 파일별 줄 수·LF/CRLF·들여쓰기·개행·코드·키·enum·포맷/인수 순서·이스케이프·Markdown/Mermaid 구조·명령을 대조했다. 이번 28개 항목의 새 용어/계약 보류는 발견하지 못했다. 지정 위치 밖의 추가 전송/전송 속도 제한 및 플랫폼 동작 표기 후보는 상세 보고서에 별도로 남기며 이번에 확대 수정하지 않았다. A025–A026·A046 이후 정책/호환성·기존 U7/U15 후속 항목·다른 단위 이력은 유지한다. DB 데이터 소급 번역 마이그레이션이나 기능/정책 변경은 없다.
- **보존·실행·산출물**: 기존 TRANSLATION_PROMPT 변경·미추적 담당표·다른 파일·Git index와 기존 감사 보고서를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-audit-fix-terms.diff`, 항목별 정적 검토 `/tmp/artex-audit-fix-terms-review.md`. 전체 감사/전체 번역 완료나 테스트 통과를 주장하지 않으며 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 4 · 알림·플랫폼·심각도·등록 추가 후보 [이번 범위 용어·표기 정적 검토 완료 / 실행 미검증]
- **기준·범위**: HEAD `451bd13eb29310696ac6c782ef6bdd8288eb774a`의 직전 용어 정리 커밋을 확인하고, 기존 보고서의 “잔여·추가 후보와 범위 구분”을 최신 로컬 기준과 대조했다. 대상 기존 변경·staged 변경·MERGE_HEAD 없음. 알림 관련 5개 파일의 동일 문맥, 지정 플랫폼 설명 2곳·심각도 등급 주석 3곳·취약점 등록 안내 1곳을 처리했다.
- **수정·대응**: 알림 전송/전송 속도 제한의 주석·로그·자체 오류와 개발자 실패 안내를 정리했다. db 및 server의 음수 제한 오류는 같은 표기로 맞췄다. 변경 오류의 전체/부분 문자열·접두사·정규식 소비부와 관련 테스트를 검색했으며 직접 문구 의존은 발견하지 못해 기대값 변경은 없었다. db/notification_test.go의 변경 문자열 12곳은 실패 안내로, 입력·fixture·조건·기대 상태는 보존했다. finding_workflow.go:68의 원문 入库는 취약점 등록 문맥으로 확인하여 “등록한 뒤”로 맞췄다.
- **제외·잔여**: 일반 필드/인자/finding id 전달·실제 폴링·성공 상태 “전달됨”·송달 성공 설명은 유지했다. lease는 이미 확정 표기와 일치하여 추가 변경하지 않았다. 지정 후보는 모두 처리했고 새 용어/계약 보류는 발견하지 못했다. A025–A026·U7/U15 정책/호환성·다른 단위 이력은 유지한다. server/orchestration.go:734의 다른 “입고” 및 범위 밖 서버 테스트의 중국어 주석/안내는 수정하지 않았다. 기능·정책·DB 마이그레이션 변경은 없다.
- **정적 보존**: 소스/테스트 10개·120행 변경. 원본/현재 줄 수·LF·들여쓰기·코드·키·enum·영어·숫자·포맷 지정자/인수 순서·이스케이프를 정적으로 대조했다. 프런트 2개 파일은 지정 주석만 변경했다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표·기타 미추적 파일과 Git index를 보존했다. 이전 기록은 바이트 그대로 유지하고 이번 기록만 추가했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-audit-fix-notification-followup.diff`, 상세 정적 검토 `/tmp/artex-audit-fix-notification-followup-review.md`. 전체 감사/전체 번역 완료나 테스트 통과를 뜻하지 않으며 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 5 · A025–A026 표시 처리 [이번 범위 표시 보완·정적 검토 완료 / 실행 미검증]
- **기준·상태**: HEAD `478544deed2d0efc2277c5a1972505e171b83b3a`에서 직전 알림 표기 커밋을 확인했다. 대상 기존 변경·staged 변경·MERGE_HEAD 없음. 최신 로컬 기준과 통합 감사 A025–A026의 생성부·소비부·관련 테스트를 읽기 전용 대조했다.
- **A025**: ConversationItem에 표시용 displayTitle 한 줄을 추가하고, 화면 제목·선택/관리 접근성 이름·삭제 확인 네 곳에 사용했다. 빈 title 또는 정확히 新对话인 title만 새 대화로 표시하며 다른 사용자 지정 제목은 그대로 사용한다. 서버의 기본값·자동 제목 비교, DB/API 값 및 mock 기본 저장값은 동일하다. 이름 변경 입력은 원래 title을 사용한다. 사용자가 직접 新对话라고 입력한 경우 기본 제목과 구분할 메타데이터가 없어 같은 표시 변환이 적용되는 한계가 있다.
- **A026**: mock 분류 중복 2곳·템플릿 중복 2곳·회사 중복 1곳의 오류 문구만 번역했다. 분류/템플릿 UI는 오류를 toast로 표시한다. 회사 UI의 HTTP 409 접미사 정규식은 유지하며 이번 자체 오류는 일반 메시지 표시 경로를 따른다. skills의 已存在 검사는 /skills 업로드의 별도 계약으로 보존했다. 변경 문구에 직접 의존하는 테스트 기대값을 발견하지 못해 테스트는 수정하지 않았다.
- **추가 후보 제외**: orchestration.go:734의 원문 SeedTool 入库는 db/tools.go:31–43에서 tools 테이블에 도구를 저장하는 설명으로, 취약점 등록을 뜻하지 않는다. 요청의 취약점 등록 조건과 달라 수정하지 않았으며 근거를 상세 보고서에 기록했다.
- **보존·정적 검토**: chat/page.tsx는 표시용 상수 1줄 추가와 표시 참조 4곳만 변경(원본 대비 +1줄); handler.ts는 오류 리터럴 5곳만 변경(줄 수 동일). LF·들여쓰기·기존 JSX 구조·영어·포맷·개행·저장/전송/비교 로직·fixture·모델 응답·다른 계약값은 보존했다. 새 전문 용어/계약 보류 없음. 기존 이력·U7/U15 정책/호환성 보류·TRANSLATION_PROMPT 로컬 변경·미추적 담당표/다른 파일·Git index는 유지했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-audit-fix-display.diff`, 검토 보고 `/tmp/artex-audit-fix-display-review.md`. 테스트 통과/실패나 실제 화면 동작은 관측하지 않았고 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 6 · A068–A074·A089 [부분 정리·정적 검토 완료 / 사례 이름 1곳 용어 보류 / 실행 미검증]
- **기준·범위**: HEAD `5d58222`의 채팅 표시/mock 오류 보완 커밋을 확인했다. 대상 기존 변경·staged 변경·MERGE_HEAD 없음. 최신 기준과 현재 생성부·소비부·기존 해소 기록을 읽기 전용 대조했다.
- **A068–A073**: 알림 테스트 6개에서 name/t.Run 사례 설명 68곳 중 67곳을 번역했다. Greylisting·DingTalk·WeCom·Feishu·마스킹·자격 증명·전송 속도 제한 등 확정 표기를 적용했다. pack_test.go:133의 “标题里的图片信标”는 图片信标 미등록 전문 용어로 문장 전체를 유지하고 후보를 보고했다. 저장소의 정확한 이름 참조 및 -run 옵션을 검색했으며 연결된 선택/비교 참조는 발견하지 못했다. CHANGELOG의 清空输入框는 “空输入”의 우연한 부분 매칭이라 유지했다. 저장소 밖 개인 실행 명령은 미확인이다.
- **A074**: GUIDE의 U9 미확인 안내 한 줄만 현재 사실로 정정했다. ReporterDefaultPrompt의 실제 참조 “작업 #<id>”와 taskContextHeader의 “[작업 #%d %s (목표: %s)]”/“[작업 #%d %s]”가 일치하며 PROGRESS의 기존 라벨 해소 이력을 참조한다. 다른 U1 보류와 운영 코드/프롬프트는 변경하지 않았다.
- **A089**: 기존 통합 감사 보완 2의 A075–A079 및 H01–H35 해소 기록과 현재 schema.sql·company_scope.go·관련 알림 Go 보완 상태를 대조했다. schema의 주석/표시 문자열 3곳 및 기존 지정 Go 잔여 보완이 후속 기록으로 설명되므로 과거 U12 이력·빌드/vet PASS·초기 잔여 보고를 중복 정정하지 않았다. 이 확인은 현재 실행 검증이나 DB 저장 데이터의 소급 번역을 뜻하지 않는다.
- **원문 문제·보존**: mask_test.go:156의 원문 이름은 TLS를 꺼도 비밀번호 재명시를 요구한다고 설명하지만 입력은 tls:false→true다. 이름은 원문 의미대로 번역하고 fixture/조건은 유지했으며 번역과 별개인 원문 설명 문제로 기록했다. 6개 테스트는 이름 리터럴 외 주석·실패 안내·입력·중국어 필터/마스킹/Unicode 데이터·SMTP 응답·URL·비밀값·JSON·정규식·기대값·로직·영어·포맷·줄 수·LF·들여쓰기를 보존했다.
- **범위·실행·산출물**: 기존 U7/U15 후속 보류·다른 감사 항목·이력·TRANSLATION_PROMPT 로컬 변경·미추적 파일·Git index를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-audit-fix-test-labels.diff`, 항목별 보고 `/tmp/artex-audit-fix-test-labels-review.md`. 테스트 통과/실패는 관측하지 않았으며 fetch·병합·stage·commit·push·다음 묶음은 진행하지 않았다.


### 통합 감사 보완 6 · H01 해소 [이번 범위 설명·기록 정리 및 정적 검토 완료 / 실행 미검증]
- **확정·적용**: “标题里的图片信标”를 이번 Markdown 이스케이프 테스트 사례 설명에 한정해 “제목 안의 외부 이미지 삽입”으로 등록하고 notify/pack_test.go:133의 name만 변경했다. 图片信标의 일반 번역으로 확대하지 않는다. 직전 67개 이름 번역과 GUIDE 정정을 보존하여 이번 사례 설명 68곳의 보류는 모두 해소했다.
- **정적 보존·잔여**: 이미지 URL·Item.Name·must/wrong·입력·기대값·검사 로직·영어·포맷·들여쓰기·줄 수·LF는 그대로다. TLS 설명/fixture 방향 불일치는 원문 문제로 유지한다. A089의 과거 U12/PASS 기록 및 U7/U15·다른 감사 항목은 변경하지 않았다. 기존 로컬 변경·미추적 파일·Git index를 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 diff `/tmp/artex-audit-fix-test-labels.diff`, 검토 보고 `/tmp/artex-audit-fix-test-labels-review.md`를 갱신했다. fetch·병합·stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 7 · A051·A053·A090 [표시·안내 정정 및 정적 검토 완료 / ProxyAddr 설명 확인 보류 / 실행 미검증]
- **기준·범위**: HEAD `9262641`의 직전 알림 테스트 사례·작업 라벨 정리 커밋을 확인했다. 대상 기존 변경·staged 변경·MERGE_HEAD는 없었다. 최신 로컬 기준과 통합 감사 항목을 생성부·소비부·관련 테스트 및 Git 이력으로 재대조했다.
- **A051 표시/계약 분리**: mentionKinds의 기존 kind·label·alias에 한국어 displayLabel만 추가했다. 선택 목록·검색 안내·검색 결과 종류 배지와 selectedMentions의 선택된 참조 표시/접근성 이름은 표시값을 사용한다. categories[index].label의 실제 삽입 경로, mentionSearch의 기존 중국어/영어 검색 문법, mentionToken·선택 토큰 정규식·token/start 반환·제거 범위·서버 파서·API/DB 값은 보존했다. 알 수 없는 종류의 기존 fallback도 유지한다. 입력창의 @중국어 종류 및 @[중국어 종류#id …] 토큰은 DNT로 남으며 레코드 이름/설명과 서버 스냅샷은 변환하지 않는다.
- **기대값 대응**: web/src/lib/chat-mentions.test.mjs:28의 선택된 참조 표시 기대값에서 종류명만 취약점으로 맞췄다. 중국어 제목/설명 fixture·Unicode/개행/커서 입력·토큰·다른 기대값·검사 로직은 동일하다. 직접 소비자는 MentionTextarea와 이 테스트이며 서버 멘션 테스트는 중국어 계약 입력을 유지한다.
- **A053 안내 정정**: server/intercept.go:50의 자체 오류 안내만 decision/comment 두 문자열 필드·allow/ask/deny·실제 comment 앵커 안내로 정정했다. 裁决은 ParseVerdict의 앵커가 아니며 实际操作：·；成功后的后果：·；命中规则：는 정확히 유지한다. 오류 안내를 직접 비교하는 테스트/소비부는 검색에서 발견하지 못했다. 파서·프롬프트·정책은 변경하지 않았다.
- **A090 확인·유지**: db/db.go:192의 설명은 수정하지 않았다. 원문 记录代理地址(驱动 if 双文案), renderSystem의 Go 템플릿 렌더, WorkerVars.ProxyAddr와 prompt_test.go의 dual-text 주석 및 {{if .ProxyAddr}}…{{else}}…{{end}} fixture를 확인했다. 두 분기의 문구 선택을 가리킨다는 근거는 있으나 双文案의 직접 정의는 없고 현재 기본 worker 본문에는 해당 if가 없다. 후보는 “프록시 주소(사용자 프롬프트의 if/else 안내 문구 선택에 사용)”이다. ProxyAddr는 캡처 OFF에서도 전역 프록시를 반환할 수 있어 레코딩 여부와 동일시하지 않는다. 후보를 확정 용어로 등록하지 않았으며 변수·기본값·설정·DB 저장은 보존한다.
- **정적 보존·범위**: 표시 분리에 필요한 최소 참조 변경 외 코드·정규식·키·enum·영어·숫자·포맷/인수·줄 수·LF·들여쓰기는 보존했다. 진행 기록은 이번 이력만 추가했다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index 및 U7 live 테스트/커스텀 프롬프트와 U15 정책 후속 보류는 유지한다. 이번 안내 정정이 기존 후속 문제 해결을 뜻하지 않는다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 전체 diff `/tmp/artex-audit-fix-contract-display.diff`, 상세 보고 `/tmp/artex-audit-fix-contract-display-review.md`. fetch·병합·stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 8 · S001·S002·S003·S006 [부분 번역·정적 검토 완료 / 미등록 용어 4개 블록 보류 / 실행 미검증]
- **기준·범위**: 로컬 HEAD `d4d288b4`의 직전 A051·A053 커밋을 확인했다. 대상 기존 변경·staged·MERGE_HEAD 없음. 최신 네 기준 문서와 감사 상태를 적용했으며 fetch·병합하지 않았다.
- **S001**: server/notify_api_test.go 전체를 읽고 중국어 개발자 주석·실패 안내·사례 이름 5개를 번역했다(174행 변경). 입력·fixture·기대값·정규식·비교/검사 조건과 기존 U10 기대값은 그대로다. 실패 안내:496의 검증용 원문 인용 `状态变更 → 已修复`는 보존했다. 打桩(29–31), 令牌桶(855), 量纲 관련 주석(911–917)·실패 안내(931–932)의 4개 블록은 미등록 전문 용어로 전체 원문 보류했다.
- **S002·S003·S006**: engine_emptyturn_test.go의 지정 t.Run 사례 7개만 번역하고 무진행 턴을 적용했다. 앞부분의 다른 table 사례 이름·중국어 입력은 이번 범위 밖으로 보존했다. promptcatalog.go의 지정 4행에서 같은 침투 문맥의 공격면 6곳을 공격 표면으로 맞췄으며 모델 지시 강도·본문·코드는 유지했다. README 고정 라벨 뒤 조사 두 곳을 수정 완료로/수정 완료를로 정정했다.
- **참조·정적 보존**: 바뀐 사례 이름 12개와 저장소 -run 직접 참조를 검색했으며 정의 밖 연결 참조는 발견하지 못해 다른 참조 파일 수정은 없다. 외부 개인 실행 선택은 미확인이다. 원본/현재 줄 수 951/951·146/146·122/122·505/505, LF·들여쓰기·주석/허용 문자열 밖 코드·포맷 지정자 종류/개수/순서·인수·이스케이프·입력/기대값을 정적으로 대조했다. 줄 수의 실제 값은 상세 보고서와 함께 확인한다.
- **감사 정정·제외**: A049 정규식은 DNT, depth/coverage 동작 차이는 별도 미해결 기능 후속이다. A050은 표시값(Agent 이름 등)·모델용 mock 프롬프트/선택 조건·감사 입력/fixture를 나눠 기록했으며 data.ts는 수정하지 않았다. A052·A090·S004·S005와 U7/U15 정책·호환성은 제외하여 유지한다. 기존 이력·다른 단위 상태를 변경하지 않았다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·Git index를 보존했다. `/tmp/artex-audit-fix-final-text.diff`, `/tmp/artex-audit-fix-final-text-review.md` 및 갱신 감사 상태 MD/TSV. 전체 번역 완료/테스트 통과를 뜻하지 않으며 stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 8 · H01–H04 용어 보류 해소 [이번 범위 번역·정적 검토 완료 / 실행 미검증]
- **용어·범위**: 打桩=stub 처리(이번 테스트의 함수 반환값 지정), 令牌桶=토큰 버킷(알림 전송 속도 제한 알고리즘), 量纲=측정 단위(요청 횟수와 배치당 취약점 개수)를 중복·충돌 없이 문맥 제한과 함께 확정 등록했다. server/notify_api_test.go의 H01–H04 주석/실패 안내를 번역하고, 전송 시점 설명 두 곳을 정정했다. 직전 부분 번역 기록은 이력으로 유지한다.
- **정적 보존**: 이번 추가 수정 14행, 파일 원본/현재 951/951줄·LF·들여쓰기·줄 배치 동일. 입력·fixture·기대값·검사 조건·코드·포맷 지정자와 인수 순서·이스케이프를 보존했다. 주석의 옛 메시지 인용 「近 30 分钟新增 1 个漏洞」와 실패 안내의 검증용 인용 状态变更 → 已修复는 원문 예시/인용으로 유지한다. 기존 세 파일 보완·TRANSLATION_PROMPT 변경·미추적 파일·Git index도 보존했다.
- **잔여·실행**: S001 용어 보류 4개 블록 해소. S002의 범위 밖 table 사례 이름 8개는 남은 번역 작업으로 유지하며 A052·A090·S004·S005 및 U7/U15 후속 보류를 변경하지 않았다. 전체 감사/번역 완료나 테스트 통과를 뜻하지 않는다. 실행 검증은 **사용자 요청으로 미실행**이다. 누적 diff `/tmp/artex-audit-fix-final-text.diff`, 검토 보고 `/tmp/artex-audit-fix-final-text-review.md`, 감사 상태 MD/TSV를 갱신했다. fetch·병합·stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 9 · S002 잔여·S004·S005 [이번 범위 정리·정적 검토 완료 / 실행 미검증]
- **기준·상태**: HEAD `564c80439114a88c58cb1c2dcaf3593199c48a41`의 직전 서버 테스트 설명/잔여 용어 정리 커밋을 확인했다. work/ko-translation이며 대상 기존 변경·staged·MERGE_HEAD 없음. 최신 네 기준 문서를 적용하고 fetch·병합하지 않았다.
- **S002**: TestIsThinkingOnlyTurn의 앞부분 table 사례 이름 8개만 번역했다. 생각/도구/본문 구분·완전히 빈 assistant 턴·도구 결과·assistant 메시지 없음·빈 이력의 의미를 유지한다. 입력·fixture·thinking/guard 값·true/false 기대값·조건·t.Run(c.name)·기존 지정 t.Run 7개 번역은 동일하다. 전체 추적 파일에서 옛 이름의 정의 밖 직접 참조를 발견하지 못했으며 외부 개인 -run 명령은 미확인이다.
- **S004**: 인터셉트 페이지 :542–543의 텍스트/툴팁 및 agent-editor :1140의 도구 목록 표시 구분자 3곳을 ", "로 맞췄다. 저장은 원래 enabled_tools/tool_names 배열을 사용한다. planner :185·191은 renderTriggers가 모델 입력의 자유 서술 목표/힌트 목록을 만드는 경로로, 비교/파싱/저장/출력 계약 의존을 발견하지 못해 이번 사용자 지정 허용에 따라 구분자 2곳만 변경했다. 목표·힌트 원문 값·지시 강도·포맷/인수는 보존한다. 다른 구분자는 변경하지 않았다.
- **S005**: conversations.go:561의 영어 주석 안 인용만 “사용자가 이번 대화를 중지했습니다”로 맞췄다. pgStopConversation의 AbortChatStoppedByUser→Chat→captureRunSession→terminalText 경로 및 cancelcause.go:56의 Short와 대조했으며, 전체 결과 문구를 단독 문자열로 단정하지 않고 실제 요약에 포함되는 취소 사유를 인용했다. 운영 출력·취소 처리·영어 주석의 나머지는 동일하다.
- **정적 검토·잔여**: 5개 소스/테스트 파일의 총 14행만 변경했고 파일별 줄 수·LF·끝 개행·들여쓰기·허용 리터럴/인용 밖 코드·포맷 지정자/인수 순서·이스케이프 동일. 연결된 테스트 기대값 수정 필요는 발견하지 못했다. 새 용어/계약 보류 없음. A052·A090·A050 및 U7/U15 후속·다른 단위/과거 기록은 보존한다. 전체 감사/번역 완료나 테스트 통과를 뜻하지 않는다.
- **보존·산출물·실행**: 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·Git index 보존. 전체 diff `/tmp/artex-audit-fix-labels-separators.diff`, 검토 `/tmp/artex-audit-fix-labels-separators-review.md`, 감사 상태 MD/TSV 갱신. 실행 검증은 **사용자 요청으로 미실행**이다. stage·commit·push·다음 묶음 없음.


### 통합 감사 보완 10 · A050 승인 문자열·A090 [승인 범위 번역·정적 검토 완료 / 기록 fixture 후속 검토 보류 / 실행 미검증]
- **기준·범위**: HEAD `a344782784ec0a90e2eb9b5512094819b3738cbd`의 직전 사례 이름/목록 구분자 커밋을 확인했다. work/ko-translation, 대상 기존 변경·staged·MERGE_HEAD 없음. 최신 네 기준 문서와 앞선 읽기 전용 검토를 적용했으며 fetch·병합하지 않았다.
- **A050**: mock/data.ts의 Agent 이름 5곳(3067·3080·3093·3106·3136), 버전 template_text 2곳(3155·3156), 편집 prompt 1곳(3163), 마무리 기본값 2곳(3168·3172), skill 설명 3곳(3218·3229·3239)을 번역했다. `${a.name}`을 통한 이름 삽입, `{{.Goal}}`·`{{.AssetSummary}}`·`{{.RouteHint}}` 및 이스케이프를 보존했다. 시간/스텝 수의 제한이 곧 소진된다는 뜻과 즉시 마무리 지시를 유지했다. 실제 Agent 정책을 복제하거나 새 지시를 추가하지 않았다.
- **원문 문제·잔여**: mock api-recon의 권한 우회 영역 열거 설명은 실제 SKILL의 API/파라미터 정찰 전용·취약점 공격 금지 범위와 다르다. 번역과 별개인 원문 범위 문제로 기록하며 어느 쪽도 정책을 고치지 않았다. 감사 user_message/context·실제 도구 입출력·명령·LLM 요청/응답 기록 fixture·`[模型]` 5곳은 보존했다. A050의 기록 fixture 후속 검토는 유지하며 전체 A050 또는 전체 감사를 완료로 단정하지 않는다.
- **A090**: ProxyAddr description 1곳을 “현재 실행에 적용되는 프록시 주소(사용자 프롬프트에서 참조하거나 if/else 안내 문구 선택에 사용 가능)”으로 정정했다. 캡처 OFF의 전역 프록시 및 사용자 템플릿 if/else 소비와 대조했다. 변수·기본 예시·source·프록시 동작은 그대로다. 기존 seed upsert(db.go:215–220)는 이후 초기화에서 기존 DB 변수 메타데이터에도 설명을 반영할 수 있지만 사용자 프롬프트 본문을 바꾸지 않는다. 이번에는 DB를 열거나 갱신하지 않았다.
- **정적 보존·테스트**: 승인 13+1줄 밖 소스 바이트 동일, mock 4223/4223줄·db 602/602줄, LF·들여쓰기·끝 개행·포맷 지정자/변수/이스케이프 보존을 확인했다. 직접 연결된 테스트 기대값 의존은 발견하지 못해 테스트 수정은 없다. 별도 fixture의 规划者 및 ProxyAddr 템플릿 입력은 보존한다. A052·U7/U15 후속·기존 팀원 기록·다른 단위 상태를 유지한다.
- **보존·산출물·실행**: 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·Git index를 보존했다. `/tmp/artex-audit-fix-mock-proxy.diff`, `/tmp/artex-audit-fix-mock-proxy-review.md` 및 감사 상태 MD/TSV를 갱신했다. 실행 검증은 **사용자 요청으로 미실행**이다. stage·commit·push·다음 작업 없음.


### 통합 감사 보완 11 · A050 기록 예시·A052 고유명 최종 분류 [분류·정적 검토 완료 / 원문 정책·호환성 후속 유지 / 실행 미검증]
- **로컬 기준·커밋**: work/ko-translation, HEAD와 로컬 origin/work/ko-translation 참조가 모두 `f3611044ad4e29acee8878fd4f3a0963c3bd573b`이며 직전 mock 안내·ProxyAddr 설명 커밋이 반영됐다. 로컬 ahead/behind 0/0, staged·MERGE_HEAD 없음. fetch하지 않았으므로 현재 원격 최신성이나 별도 push 통신 성공을 확인했다고 주장하지 않는다.
- **A050 기록 경로**: llmRecordDetail(data.ts:4158–4223)의 request_body/raw_request는 mock handler:2254→api.llmRecordDetail→function/llm-records/page.tsx:140–156,239–249,514–545의 정규화/HTTP 원문 표시·복사에 쓰인다. 4165/4190의 system 및 4166/4192의 user 지시는 두 뷰에서 같은 기록을 표현하고, 4196은 raw_request의 모델용 도구 schema 설명 예시다. 활성 Agent 편집/모델 실행 입력으로 재사용되는 경로는 발견하지 못했다.
- **A050 감사 경로·해소**: interceptDetails:3513·3515·3544는 structuredClone→mock history detail→approval-records.tsx:441–478의 감사 user_message/context 표시로 소비된다. mock 승인 decide는 상태/결과 메타데이터만 바꾸며 도구를 실행하지 않는다. 지정 8곳을 “원문 기록 예시 보존”으로 분류하고 반복 번역 보류에서 분리했다. 이는 계약상 DNT라는 뜻이 아니며 실제 DB/모델 출력을 확인한 것도 아니다. 승인 편집용 13곳은 그대로 유지한다. 명령·HTTP/JSON/SSE·문자 수·다이제스트·[模型] 토큰을 보존했다.
- **A052**: docker-compose.bench.yml:1의 腾讯 TSec Benchmark는 원문 고유명 보존으로 분류한다. 공식 영문 표기는 미확인이며 임의 영문 이름을 만들지 않는다. 주석은 Compose 설정/계약에 쓰이지 않고 공식 영문 확인은 현재 기능·계약 유지의 필수 조건이 아니다. 원문 고유명 자체를 번역 누락으로 집계하지 않는다.
- **감사 잔여·한계**: A090 설명 정정은 직전 커밋에 포함됐다. 전달 감사 원본 4개를 내용 변경 없이 translation/audit/에 보관하고 A001–A090·S001–S006의 ID·과거 위치·발견 근거를 복구했다. 원본으로 A048은 U15 로그인·실제 세션 정책 충돌임을 확인했다. A049는 depth/coverage 기능 차이를 주 미해결 상태로, NEGATIVE_RE 원문 보존을 별도 하위 DNT로 구분한다. U7(A046 live 중국어 동사, A047 공통 블록 중복/커스텀 본문)과 mock api-recon/실제 SKILL 범위 차이는 유지한다. 승인 A050 13곳 수정과 지정 기록의 원문 기록 예시 보존, A052 고유명 보존, A090 설명 수정 완료를 유지한다. 원본 근거 복구와 현재 소스 인용 재대조 범위는 별도로 기록하며 새 의미 전수 감사나 실행 검증으로 단정하지 않는다. 전체 한국어 문장 의미 전수 대조는 미완료다.
- **보존·산출물·실행**: 저장소에서는 이번 미커밋 기록의 원본 확보·근거 연결 설명만 정정하고 translation/audit/에 원본 4개와 current-status.md·.tsv, remaining-issues.md를 보관했다. 소스·용어집·테스트·fixture·과거 커밋 기록·기존 TRANSLATION_PROMPT 변경·미추적 담당표/전달 파일·Git index는 보존했다. 이전 임시 산출물은 과거 이력으로 유지하며 지속 보관 결과는 translation/audit/를 기준으로 한다. 이번 전체 diff /tmp/artex-audit-recovery.diff. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.


### 핵심 프롬프트 의미 감사 · U1/U2/U7 [지정 본문 전체 정적 대조 완료 / 의미 후보·용어 보완 잔여 / 실행 미검증]
- **기준·상태**: HEAD `b5468a7839ada1f322b4f1e4ff7d219481f15135`의 직전 감사 문서 복구 커밋을 확인했다. work/ko-translation, staged·MERGE_HEAD 없음. 기준 문서 4종과 translation/audit 현재 기록을 적용했으며 fetch하지 않아 origin 최신성은 주장하지 않는다.
- **실제 대상**: U1 autoDefaultTmpl·pentestDefaultTmpl·DefaultAssistantPrompt·ReporterDefaultPrompt·chatWorkDirSpec, U2 goalsDefaultTmpl·goalsScopeTail·plannerDefaultTmpl·constraintBlock, U7 JudgeContextBoundary·JudgeOutputContract·DefaultJudgePrompt 전체를 대조했다. 실제 소비 경로에 연결된 U3 정의의 artifactSpec·genericWrapUpDefault·plannerWrapUpDefault·plannerTaskTimeoutDefault도 읽기 전용으로 포함해 총 16개 본문/조립 단위다. U3/U4 worker/mainagent 전체 본문·동적 user 지시·도구 설명 전체·DB 사용자 정책은 포함하지 않았다.
- **원문 기준**: U1은 7303d46, U2는 aa0edc0, 연결 마무리/산출물은 1b5e613, U7 세 상수는 cfef293·dac067a·141d36f 각각의 직접 부모에서 번역 직전 중국어 원문을 확보했다. 대상별 전체 SHA와 원문/현재 범위·전문을 보고서에 보관했다. 16개 모두 원문 확보·본문 끝까지 의미 대조를 했으며 대상 내부의 문자열/주석 밖 정적 구조·리터럴 개수·템플릿 변수 순서도 보조 대조했다. 실제 렌더/SDK 실행은 하지 않았다.
- **새 발견**: M001–M007 7개를 별도 기록했다. 확정 용어 불일치 4개는 힌트 추가, 인젝션 포인트, 공격의 인젝션/권한 우회, 목표 달성 검증·정량적 확인 문맥이다. 의미 추가 확인 후보 3개는 Auto의 少空转 정도 표현, Reporter 接管의 대상 한정, planner 定论의 확정 한정 표현이다. 실제 모델 행동 변화나 정책 오류를 관측했다고 단정하지 않으며 소스는 수정하지 않았다.
- **보존·기존 보류**: 승인된 中文→한국어 출력 지시, Reporter 작업 라벨 정정, 120자 지시와 별도 파서 2400바이트, U7 앵커/키/enum/예시·명령·원문 정책을 구분했다. A046–A049와 remaining-issues의 U7/U15 정책·호환성·기능 문제, 다른 단위의 기존 상태는 유지한다. 지정 본문 대조가 U1/U2 기존 테스트 보류나 U7 후속 보류 전체 해소를 뜻하지 않는다.
- **한계·산출물·실행**: translation/audit/prompt-semantic-review.md 및 prompt-semantic-findings.tsv. 지정 본문 중 원문 미확보·미독 범위는 없지만 전체 저장소 한국어 의미 전수 대조는 미완료다. 실제 DB/모델 출력·테스트는 미확인이고 실행 검증은 **사용자 요청으로 미실행**이다. 소스·용어집·테스트·정책·DB·기존 감사 문서·TRANSLATION_PROMPT 변경·미추적 파일·Git index를 보존했다. fetch·병합·stage·commit·push·다음 작업 없음.


### 핵심 프롬프트 의미 감사 승인 보완 · M001–M007 [승인 수정·정적 대응 완료 / 기존 정책·호환성 후속 유지 / 실행 미검증]
- **대상·수정**: agent/promptcatalog.go의 autoDefaultTmpl 18·22행, pentestDefaultTmpl 38행, ReporterDefaultPrompt 95행과 agent/planner.go의 plannerDefaultTmpl 312·316·320·330행에서 사용자가 지정한 8개 표현만 수정했다. 少空转는 진행 없는 동작 줄이기, 힌트는 추가, 注入点은 인젝션 포인트, 接管 예시는 장악, 定论는 확정된 결론, 공격은 인젝션/권한 우회, 목표 달성은 정량적 확인/달성 검증 항목으로 적용했다.
- **용어집**: 독립 接管 중복·충돌을 확인하고 Reporter 피해 예시처럼 대상을 특정하지 않은 문맥의 장악 1개만 [확정] 등록했다. 기존 계정·서브도메인·서버/도메인 장악 문맥과 다른 확정 항목은 유지했다.
- **정적 검토·테스트**: 지정 문장/상수 경계로 적용했다. 소스 2개는 각각 4줄 변경이며 원본/현재 줄 수·LF·끝 개행·코드·영어·도구명·키·enum·숫자·템플릿 변수·포맷·인수·금지/조건과 지정 표현 밖 내용은 동일하다. 관련 테스트의 직접 전체/부분 문자열 및 상수명 의존을 확인했으며 기대값 변경 필요는 발견하지 못해 테스트를 수정하지 않았다. 테스트 통과나 실제 모델 동작을 확인한 것은 아니다.
- **감사 이력·상태**: prompt-semantic-review.md와 prompt-semantic-findings.tsv에 원문·승인 전 번역·발견 근거를 이력으로 보존하고 승인 수정 후 번역/정적 대응 상태를 별도로 기록했다. current-status.md의 새 M 7개는 A/S 기존 96개 집계와 분리했다. A046–A049와 다른 정책·기능·호환성 문제 및 전체 저장소 의미 전수 대조 미완료는 유지한다.
- **보존·산출물·실행**: 기존 미커밋 감사 내용·TRANSLATION_PROMPT 변경·미추적 담당표/전달 파일·기타 로컬 변경·Git index를 보존했다. 누적 전체 diff `/tmp/artex-prompt-semantic-approved.diff`. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.


### 핵심 프롬프트 의미 감사 · U3 worker/U4 mainagent [지정 본문 전체 정적 대조 완료 / M008–M010 용어 잔여 / 실행 미검증]
- **기준·상태**: HEAD `1d8380d63580483235a637a9cebb9dc7599bf2da`(work/ko-translation), M001–M007 승인 보완 255304c와 감사 산출물 추적 제외 1d8380d 반영 확인. staged·MERGE_HEAD 없음. 기준 문서 4종과 로컬 감사 상태·잔여·이전 의미 감사 기록을 적용했다. fetch하지 않아 원격 최신성은 주장하지 않는다.
- **실제 범위**: workerDefaultTmpl(230–243), workerTrafficBlock(250–255), workerArtifactSpec(267–269), settleWrapUpPrompt(154), renderIntentTask(316–318), renderWorkerGraphOverview(325–338), workerTaskTimeoutDefault(wrapup:120), mainAgentDefaultTmpl(mainagent:87–101), mainAgentWrapUpDefault(wrapup:46) 전체와 Worker.Execute 고정 시작/자산/재개/수동 입력 안내(406·420·490·510·517·519–520·522–523)를 중국어 원문과 끝까지 대조했다. 완료된 artifactSpec·generic/planner 마무리·constraintBlock은 이전 감사 기록을 참조하고 신규 연결만 확인했다. mainagent의 전용 작업 타임아웃 기본 문구/선택 경로는 없다.
- **원문·이력**: U3는 1b5e613의 직접 부모 `0e10bfcde7cc32c51064b04fc5b3e14d56743309`, U4 본문은 a44df30의 직접 부모 `1b5e61362ddd81d11e1eb9cd39b384fbab12da5c`에서 확보했다. mainagent 마무리는 U3 부모 기준이다. 이후 42e1deb의 사실 분할 부정어 및 패시브 정찰 정정이 현재에 반영됨을 확인했다. 해당 대상 이후 기능 변경은 발견하지 못했다.
- **새 발견·테스트**: M008 공격 注入의 주입 2곳, M009 add_intent 의도 추가의 직접 주입 1곳, M010 指纹의 지문 2곳을 확정 용어 불일치로 보고한다. 후보는 각각 인젝션/추가/핑거프린트이며 이번에는 소스 미수정이다. 역할·조건·예외·도구 순서·실제 재현 증거·종료/재개·마무리 의무를 대조했고 추가적인 새 의미 변경은 발견하지 못했다. 관련 전체/부분 문자열·역할/상수 검색과 테스트 소비 대조에서 M 후보의 직접 기대값 의존은 발견하지 못했다. 산출물 제목/traffic_search와 overview 누출 금지 검사·override fixture를 구분했으며 테스트는 수정/실행하지 않았다.
- **보존·한계**: 9개 본문 단위의 문자열/주석 밖 구조·리터럴 수·영어 원문 토큰·템플릿/포맷/이스케이프/숫자/강조를 정적 보조 대조했다. 실제 DB·SDK·모델 동작은 미확인이다. M001–M007 이력, A/S 96개 집계, A046–A049 및 다른 정책·기능·호환성 후속은 유지한다. 지정 본문 원문 미확보/미독은 없지만 전체 저장소 의미 전수 대조는 미완료다.
- **산출물·실행**: Git 제외된 로컬 translation/audit/worker-mainagent-semantic-review.md·worker-mainagent-semantic-findings.tsv만 작성하고 이 기록만 추가했다. 기존 프롬프트 변경·미추적 담당표/audit-input/기타 파일·Git index를 보존한다. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.


### 핵심 프롬프트 의미 감사 승인 보완 · M008–M010 [승인 수정·정적 대응 완료 / 기존 후속 유지 / 실행 미검증]
- **승인 범위**: workerDefaultTmpl의 234·240행 공격 기법 주입 2곳을 인젝션으로, 238·240행 지문 2곳을 핑거프린트로 변경했다. mainAgentDefaultTmpl 92행의 직접 주입한다 1곳을 직접 추가한다로 변경했다. 현재 문장·상수 경계를 확인하여 정확히 5곳만 적용했다. 기존 확정 용어를 사용하며 용어집은 수정하지 않았다.
- **보존·정적 대조**: worker.go 535/535줄, mainagent.go 202/202줄. 지정 표현 밖 문장·정책·조건·증거 요구·코드·영어·도구명·숫자·템플릿/포맷·들여쓰기·개행을 보존했다. 240행의 기존 조사는 승인 범위 밖이므로 그대로 유지했다. 전체/부분 문자열 및 상수명의 테스트 의존을 읽기 전용 재검색하고 prompt_test.go의 산출물/트래픽/override 검사와 review_context_test.go의 overview 검사 경로를 대조했다. 직접 기대값 수정 필요는 발견하지 못했고 테스트는 미수정이다.
- **감사 기록**: Git 제외된 worker-mainagent-semantic-review.md·worker-mainagent-semantic-findings.tsv에 수정 전 원문/번역/발견 근거를 유지하고 승인 후 번역과 정적 대응 완료 상태를 별도로 추가했다. M001–M007·A/S 96개 집계·A046–A049 및 다른 정책·기능·호환성 후속은 유지한다. 전체 저장소 의미 전수 대조는 미완료다.
- **산출물·실행**: 커밋 검토용 /tmp/artex-worker-mainagent-approved.diff에는 소스 두 파일과 PROMPT_PROGRESS만 포함한다. 로컬 감사 보고서는 커밋용 diff에 포함하지 않는다. 기존 미커밋 PROGRESS 기록·TRANSLATION_PROMPT 변경·미추적 담당표/audit-input/기타 파일·Git index를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.

- **M010 조사 승인 보완**: workerDefaultTmpl의 “버전/핑거프린트이 CVE에 매칭됨”을 “버전/핑거프린트가 CVE에 매칭됨”으로 조사 한 곳만 정정했다. 앞선 조사 보존 설명은 당시 이력이며 이번 사용자 승인으로 해소했다. 코드·조건·증거 요구·535줄·LF·기존 변경·Git index를 보존했다. 로컬 감사 보고서/TSV에 승인 후 상태를 반영하며 커밋용 diff는 소스 두 파일과 PROGRESS만 포함한다. 실행 검증은 **사용자 요청으로 미실행**이다.


### 핵심 도구 설명 의미 감사 · U5 전체 [지정 자연어 전체 정적 대조 완료 / M011 용어 후보·M012/M013 원문 설명 후속 / 실행 미검증]
- **기준·범위**: HEAD `e29d563da4370aab6d49e93b72693759ae6965a2`(work/ko-translation)에서 직전 M008–M010 승인 보완 및 조사 보완 e29d563 반영 확인. staged·MERGE_HEAD 없음. 최신 기준 문서 4종과 로컬 감사 이력을 적용하고 tools.go/tools_insert.go/tools_digest.go/toolcatalog.go의 26개 도구 설명·파라미터/schema·자체 오류/성공/결과 자유 서술 231개 및 범위 이해에 필요한 주석을 원문과 끝까지 대조했다. toolcatalog에는 직접 중국어 모델용 리터럴이 없고 DB 설명/schema 재정의·기본값 처리 경로를 확인했다.
- **원문·대응**: 4caba7f부터 0406037까지 10개 번역 묶음의 직접 부모에서 해당 중국어 원문을 확보했다. 현재와 주석/문자열 밖 코드·키·enum·인수·분기·반환 구조가 같으며 네 파일의 줄 수는 2054/699/165/201로 동일하다. 포맷 지정자 종류/순서/개수·이스케이프·템플릿·숫자를 별도로 대조했다. 기존 승인 coverage.note/status·DNS 레코드 표현·연결 공백 등은 이력으로 보존했다.
- **결과**: M011은 steer_work 설명/성공 안내(tools.go:1707,1735)의 지시 추가 문맥에서 주입→추가 용어 후보 1개(두 곳)다. M012는 커버리지 비활성화 시 add_task_scope/auto-scope 처리의 오래된 주석, M013은 PlannerTools의 전체 저장소 조회 주석으로 둘 다 중국어/영어 원문부터 존재하는 설명·구현 불일치이며 번역 오류와 분리했다. 소스는 수정하지 않았다. 그 밖의 필수/금지/권고·조건·예외·범위·읽기 전용·증거·순서·종료/재개 의미 차이는 이번 대조에서 발견하지 못했다.
- **소비·테스트**: BuiltinToolSeeds→wireTools/SeedTool→ToolResolve/DecorateTool→AugmentTools 경로와 편집 Description/Schema/default, nil guard·finding 첫 줄/JSON·범위 membership·UTF-8 budget 관련 테스트를 읽기 전용 대조했다. DB 사용자 편집 설명과 외부 오류/summary/detail/evidence는 번역 기본값과 구분했다. 직접 기대값 수정 필요는 발견하지 못했으며 테스트는 수정/실행하지 않았다. 실제 DB 저장 설명·모델 출력은 미확인이다.
- **이력·한계·산출물**: A/S 96개 집계·M001–M010 승인 이력·A046–A049 및 기존 후속 보류는 유지한다. 이번 지정 자연어 미독/원문 미확보 범위는 없지만 외부 SDK/MCP 설명·DB 편집 본문 전체·다른 단위는 별도 범위이며 전체 저장소 의미 전수 대조는 미완료다. Git 제외 로컬 translation/audit/tool-semantic-review.md·tool-semantic-findings.tsv에 원문/현재/근거/후속을 기록했다. 커밋 검토용 /tmp/artex-tool-semantic-progress.diff에는 이 기록만 포함한다. 기존 로컬 변경·미추적 파일·Git index를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.


### 핵심 도구 설명 의미 감사 승인 보완 · M011 [승인 수정·정적 대응 완료 / M012·M013 원문 설명 후속 유지 / 실행 미검증]
- **승인 범위**: agent/tools.go의 steerWorkTool 설명(1707행)의 “방향 조정 지시를 실시간으로 주입하며”를 “방향 조정 지시를 실시간으로 추가하며”로, 성공 안내(1735행)의 “방향 조정 지시를 주입했습니다”를 “방향 조정 지시를 추가했습니다”로 정확히 두 곳만 수정했다. 기존 확정 용어를 적용하고 용어집은 수정하지 않았다.
- **정적 대응·보존**: 함수/문장 경계로 확인했다. 실행 중단 없음·기존 진행 보존·다음 동작 전 적용·kill_work 조건·%d와 인수 및 그 밖의 문장/주석/운영 코드는 그대로다. tools.go는 2054/2054줄이며 LF·끝 개행·들여쓰기·포맷을 보존했다. 추적 테스트/프런트에서 지정 전체 표현 및 대응 원문 특징 문자열을 재검색했으며 직접 기대값 수정 필요를 발견하지 못했다. 테스트는 수정하거나 실행하지 않았다.
- **이력·후속**: 로컬 tool-semantic-review.md·tool-semantic-findings.tsv는 수정 전 원문/번역/발견 근거를 보존하고 승인 후 번역·정적 대응 상태를 별도로 기록했다. M012·M013은 원문 자체 설명 문제로 미수정 유지한다. A/S 96개 집계·M001–M010 승인 이력·A046–A049 및 다른 후속 보류는 유지한다. 전체 저장소 의미 전수 대조와 실제 DB/모델 동작 확인은 미완료다.
- **산출물·실행**: 커밋 검토용 /tmp/artex-tool-semantic-approved.diff에는 agent/tools.go와 PROMPT_PROGRESS만 포함하며 감사 MD/TSV는 Git 제외 경로에 로컬로만 보관한다. 이전 미커밋 PROGRESS·TRANSLATION_PROMPT 변경·미추적 파일·Git index를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 작업 없음.


### U3/U4 보조 모델 안내 의미 감사 [확정 대상 전체 정적 대조 완료 / M014·M015 수정 후보 / 실행 미검증]
- **기준·범위**: HEAD 6d2160e의 M011 승인 보완 커밋을 확인했다. retester.go·compaction.go·noa.go·coldgraph.go·finding_workflow.go·finding_recorder.go 전체를 읽고 실제 생성→소비 경로를 확정했다. RetesterDefaultPrompt 전체, compressionSystemPrompt 및 압축 입력 네 안내, findingWorkflowTools/FindingGuidance·번호 규약·힌트 트래픽 스키마·자체 오류의 17개 모델용 자연어 리터럴을 끝까지 대조했다. FindingGuidance 생성부는 지정 파일 안에 있다.
- **원문·형식**: U3 번역 1b5e613의 직접 부모 0e10bfc, U4 번역 a44df30의 직접 부모 1b5e613에서 중국어 원문을 확보했다. 이후 451bd13·478544d 용어 승인 보완을 구분했다. 여섯 파일은 원문/현재 15/15·548/548·48/48·330/330·128/128·22/22줄이며 주석/리터럴 밖 코드가 동일하다. 숫자·포맷·템플릿·이스케이프와 도구·키·enum·인수·조회/쓰기 범위를 정적으로 대조했다.
- **결과**: M014는 retester.go:6의 发起时를 최초 보고 당시로 옮긴 시간 참조 차이다. DB는 재검증 시작 시점의 현재 취약점 스냅샷을 만든다. M015는 finding_workflow.go:58·73·77의 텍스트 인계/취약점 텍스트 기록/텍스트 증거를 문자로 옮긴 표현이다. 두 항목 네 위치의 원문 전체·현재 번역·수정 후보·근거를 로컬 감사 MD/TSV에 기록했고 소스는 수정하지 않았다. 그 밖의 확정 본문에서 지시 강도·조건·예외·재검증 결론/증거·인계 순서·계보 및 압축 보존 정보의 의미 차이를 발견하지 못했다.
- **제외·테스트·한계**: noa.go 오류는 로그이며 SDK 내부 압축 프롬프트는 미대조다. coldgraph.go는 알고리즘 연결 문맥으로 확인했고 finding_recorder.go의 findingTrafficGuidance는 정의 외 사용을 발견하지 못해 활성 모델 안내로 집계하지 않았다. 기본/마무리·U5 기검토 본문은 재감사하지 않았다. 관련 agent/server/db 테스트를 읽고 직접 기대값 의존을 검색했으며 M014/M015 후보의 기대값 변경 필요는 발견하지 못했다. 실제 DB·SDK·모델 동작/출력 및 전체 저장소 의미 전수 대조는 미확인이다.
- **보존·산출물**: A/S 96개·M001–M013 승인/발견 이력·M012/M013 원문 설명 후속·A046–A049 및 다른 후속 보류를 유지한다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·Git index를 보존한다. support-guidance-semantic-review.md·support-guidance-semantic-findings.tsv는 Git 제외 translation/audit/에 로컬로만 보관한다. 커밋용 diff에는 이번 PROGRESS 기록만 포함한다.
- **실행**: 사용자 요청으로 미실행. 소스·용어집·테스트 수정 및 fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### U3/U4 보조 모델 안내 승인 보완 · M014·M015 [승인 수정·정적 대응 완료 / 실행 미검증]
- **승인 범위**: RetesterDefaultPrompt의 retester.go:6에서 최초 보고 당시의→재검증 시작 당시의, findingWorkflowTools의 finding_workflow.go:58·73·77에서 문자만 인계하고→텍스트만 인계하고, 문자 취약점이 등록됐다는→취약점의 텍스트 기록이 등록됐다는, 문자/명령 증거→텍스트/명령 증거 네 표현만 수정했다.
- **정적 보존·테스트**: 현재 상수/함수·문장 경계로 확인했다. 소스는 15/15·128/128줄이며 LF·끝 개행·들여쓰기·지정 표현 밖 바이트가 동일하다. 정책·조건·금지·증거 요구·인계 순서·코드·키·영어·도구명·포맷/인수는 보존했다. 관련 agent/server/db 테스트에서 지정 전체/부분 표현과 상수/동적 안내 사용을 읽기 전용으로 확인했다. 직접 기대값 수정 필요는 발견하지 못했고 Worker 를 취소·보고 전 트래픽 자동 연관·바인딩에 성공한 뒤 검사는 유지한다. 테스트는 미수정·미실행이다.
- **감사·이력**: 로컬 support-guidance-semantic-review.md·support-guidance-semantic-findings.tsv에 중국어 원문·수정 전 번역·발견 근거를 유지하고 승인 후 번역 및 정적 대응 완료 상태를 별도로 기록했다. 기존 A/S 집계·M001–M013 승인/발견 이력·M012/M013 원문 설명 후속·A046–A049와 다른 정책·기능·호환성 보류는 유지한다. 실제 DB·SDK·모델 영향 및 전체 저장소 의미 전수 검증은 미확인이다.
- **산출물·실행**: 커밋용 diff는 agent/retester.go·agent/finding_workflow.go·PROMPT_PROGRESS만 포함한다. 감사 MD/TSV는 Git 제외 translation/audit/에 로컬로만 보관한다. 기존 로컬 변경·미추적 파일·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 작업 없음.


### U11 server 모델 안내 의미 감사 [확정 본문 정적 대조 완료 / M016–M021 후보 / 실행 미검증]
- **기준·범위**: HEAD 97fd0c3에서 M014/M015 승인 보완 커밋을 확인했다. server/ 추적 102개/운영 Go 55개를 검색하고 생성/등록/전달 경로로 확정한 12개 생성 파일의 모델용 자연어 165개와 모델 설명 교체에 연결된 SQL 비교값 1개를 끝까지 대조했다. orchestration·platform_tools·assembly·finding_retests·finding_workflow·finding_traffic 외 customtool·conversations·scheduler·chat_mentions·chatupload·engine도 포함했다. task_llm/manager/server_mgmt는 provider/도구/스냅샷/검증 연결부와 UI-only를 구분했다.
- **원문·보존**: bc7ca5e·db2c94c·d9768cf·b5aaf5e·18cd66c·bf425ef·9ccc5ae의 직접 부모, conversations 모델 header/병합은 별도 68704c8 직접 부모에서 대응 중국어 원문을 확보했다. 이후 용어/계약·승인 보완을 구분했다. 확정 파일들의 주석·리터럴 밖 코드/분기/인수와 숫자·포맷·이스케이프·schema 키/enum·출력 구조를 정적으로 대조했다. 코드가 같아도 리터럴 비교 계약이 보존됐다는 뜻은 아니다.
- **결과**: M016은 finding_workflow.go:22 legacy 설명을 SQL description=$2의 정확한 비교값인데 번역한 계약 불일치다. 옛 중국어 저장값과 매칭되지 않을 경로가 있으며 실제 DB 행/flag는 미확인이다. M017은 orchestration.go:367 힌트 주입→추가, M018은 :250 우아한 마무리→정상적인 마무리, M019는 :251 planner 하트비트→heartbeat, M020은 engine.go:562 공회전 알림→진행 없음 알림, M021은 finding_traffic.go:369 문자/명령 증거→텍스트/명령 증거 후보다. 총 6건이며 소스·용어집·테스트는 수정하지 않았다.
- **소비부·테스트**: 실제 hostTools/Spec→DB 재정의/필터→Options.Tools, trigger/attachment/mention→Chat, steer/empty-turn hook→모델 재개 안내를 대조했다. ControlWork는 API-only로 제외하고 SteerWork/KillWork는 모델 콜백 오류로 포함했다. 관련 server 테스트를 읽고 직접 기대값/프런트 의존을 검색했으며 5개 자연어 후보의 하드코딩 기대값 수정 필요는 발견하지 못했다. emptyTurnNudge 테스트는 상수 자체를 비교한다. 기존 migration 테스트가 정확한 중국어 역사 설명/host_search flag 회귀를 보장한다고 단정하지 않는다.
- **범위·한계·후속**: 확정 166개에는 원문 미확보 없음. 나머지 server 파일은 검색/연결 확인 수준이고 모든 한국어 본문 의미 대조 완료로 표시하지 않는다. SDK 내부·실제 DB 사용자 편집 본문·외부 실제 출력·모델 동작은 미확인이다. 기검토 agent 기본/마무리·U5·FindingGuidance는 재감사하지 않았다. A/S 96개·M001–M015 승인 이력·M012/M013 원문 설명·A046–A049 및 다른 정책·기능·호환성 보류는 유지한다.
- **산출물·실행**: server-model-semantic-review.md·server-model-semantic-findings.tsv는 Git 제외 translation/audit/에 로컬로만 보관한다. 커밋용 diff에는 이번 PROGRESS 추가만 포함한다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·기존 감사·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행. fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### U11 server 모델 안내 승인 보완 · M016–M021 [승인 수정·정적 대응 완료 / DB·테스트 후속 유지 / 실행 미검증]
- **승인 범위**: finding_workflow.go:22의 legacy SQL 비교 리터럴을 bf425ef 직접 부모 42e1deb의 중국어 원문 바이트로 복원하여 DNT 계약으로 보존했다. 현재 모델 설명 traffic.TrafficSearchDescription과 SQL·플래그·재시도·DB 저장값·마이그레이션 로직은 변경하지 않았다. orchestration.go:367의 힌트 주입→추가, :250 우아한 마무리→정상적인 마무리, :251 planner 하트비트→heartbeat, engine.go:562 공회전 알림→진행 없음 알림, finding_traffic.go:369 문자/명령 증거→텍스트/명령 증거 지정 표현만 수정했다.
- **정적 보존·테스트**: 네 소스는 172/172·920/920·1245/1245·414/414줄이며 지정 리터럴/표현 밖 바이트·LF·들여쓰기·조건·정책·영어·키·enum·도구명·숫자·포맷·인수를 보존했다. TestFindingWorkflowMigrationPreservesUserConfiguration은 reporter/hint 사용자 설명과 schema/바인딩 보존을 검사하지만 정확한 역사 중국어 traffic_search 설명 교체 및 해당 host_search 플래그 초기화를 직접 검사하지 않는다. 옛 설명과 flag=false로 현재 TrafficSearchDescription 교체를 검증하는 경우, 사용자 편집 traffic_search 설명이 그대로 유지되는 별도 경우를 테스트 보완 후보로 남긴다. emptyTurnNudge 검사는 상수 자체를 비교하며 나머지 표현의 직접 기대값 수정 필요는 발견하지 못했다. 테스트는 미수정·미실행이다.
- **후속·이력**: 역사 비교값 복원은 실제 DB 복구와 다르다. 이미 한국어 legacy 또는 완료 플래그가 저장된 실제 DB는 상태 미확인으로 별도 확인이 필요하며 이번에 변경하지 않았다. 프런트 function/tasks/page.tsx:3087·3333·3338·3353·3357 및 detail/_tabs/overview-tab.tsx:898의 하트비트/우아한 마무리는 별도 후속으로 미수정이다. 로컬 감사 MD/TSV에 수정 전 원문·번역·근거를 이력으로 보존하고 승인 후 상태를 별도 기록했다. 기존 A/S 96개·M001–M015·M012/M013 원문 문제·A046–A049와 다른 후속 보류를 유지한다.
- **산출물·실행**: 커밋용 /tmp/artex-server-model-approved.diff에는 소스 네 파일과 PROMPT_PROGRESS만 포함하며 감사 MD/TSV는 Git 제외 경로에 로컬로만 보관한다. 기존 미커밋 감사 기록·TRANSLATION_PROMPT 변경·미추적 파일·Git index를 보존한다. 실제 DB/모델 동작과 전체 저장소 의미 전수 대조는 미확인이다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 작업 없음.


### 작업 설정 프런트 용어 보완 · M018/M019 연결 표시 [승인 수정·정적 대응 완료 / 실행 미검증]
- **기준·승인 범위**: HEAD 09069c404c72ee3f0381a58854448f8808d5ee05에서 직전 M016–M021 보완 커밋 반영을 확인했다. 대상 파일 기존 변경·staged·MERGE_HEAD 없음. 최신 기준 문서와 확정 문맥을 적용하여 function/tasks/page.tsx:3087·3333·3338·3357의 하트비트→heartbeat, :3353의 우아한 마무리→정상적인 마무리, detail/_tabs/overview-tab.tsx:898의 하트비트→heartbeat만 수정했다. 조사나 별도 고정 명사 추가는 필요하지 않았다. 정상적인 마무리는 종료·정리 절차이며 목표 달성/성공 의미로 바꾸지 않았다.
- **연결·정적 검토**: 설정은 heartbeatMin→heartbeatSec→api.createTask(planHeartbeatSeconds)→plan_heartbeat_seconds 경로이며 UI 표시/개발자 주석만 수정했다. overview 제목은 엔진 상태/실행 중 Worker/최근 활동을 보여 주는 CardTitle이며 저장·파싱 키가 아니다. 지정 표현의 비교/파싱·테스트 직접 의존은 검색에서 발견하지 못했고 테스트 기대값 변경은 필요하지 않은 것으로 확인했다. web/ 내 같은 한글 용어 검색에서 지정 여섯 곳 외 추가 후보는 발견하지 못했다.
- **보존·후속**: 두 파일의 줄 수·LF·JSX 구조·들여쓰기·개행과 지정 표현 밖 바이트·코드·키/값·조건·숫자·포맷을 보존했다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index와 감사 이력을 보존한다. M016 실제 DB의 한국어 legacy/완료 플래그 상태 및 역사 설명/사용자 편집 설명 회귀 테스트 후보, 기존 정책·기능·호환성 후속은 그대로 유지한다. 과거 프런트 미수정 기록은 이력이며 이번 승인 범위만 해소했다.
- **산출물·실행**: /tmp/artex-frontend-heartbeat-approved.diff는 프런트 소스 두 파일과 PROMPT_PROGRESS만 포함한다. 로컬 감사 frontend-heartbeat-review.md·frontend-heartbeat-findings.tsv는 Git 제외 translation/audit/에만 작성한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### traffic/ 모델 안내 의미 감사 [확정 본문 전체 정적 대조 완료 / M022 용어 후보·M023 예시 추가 확인 / 실행 미검증]
- **기준·소유·범위**: HEAD 989d1f060cb1eb60d9cc252a96894c8f1ddf9fbc에서 직전 프런트 용어 승인 보완 989d1f0 반영 확인. staged·MERGE_HEAD·대상 기존 변경 없음. 최신 기준 문서와 로컬 감사 이력을 적용했다. 사용자 지정 traffic/은 GUIDE상 U13이며 U9는 sidequestion/으로 구분해 기록한다. 추적10개(운영3/테스트7)를 목록화하고 traffic_search/traffic_get/traffic_blob 세 도구의 설명/schema·결과 조립/자체 오류 및 evidence 오류를 포함한 모델용 자연어31개(traffic.go27/evidence.go4)를 끝까지 대조했다.
- **원문·결과**: 번역 a314c95 직접 부모 a255c7672bc3ba0c6cb23f3adafc4e246915c0e1에서 대응 중국어 원문을 확보했다. 해당 번역 이후 별도 파일 기능 변경은 발견하지 못했다. 모델 안내의 필수/금지·조회 범위·3/10/page0·요청/응답 확인 뒤 증거 ID 바인딩·8192/offset/hex512 조건을 대조했다. M022는 TrafficSearchDescription의 기록 프록시→레코딩 프록시 확정 용어 후보다. M023은 body_contains schema의 内网测试 예시를 내부망 테스트로 바꾼 부분으로, 원문 예시 복원 후보와 보존 기준 추가 확인을 기록하며 실제 계약 불일치로 단정하지 않는다. 소스·용어집·테스트는 미수정이다.
- **연결·보존·테스트**: HostTools→wireAgentAugment/ToolResolve/Options.Tools, SeedToolMetas→SeedTool/DB 재정의, ReadEvidence→evidence.prepare/Bind/Record 경로를 확인했다. 기검토 U5/U11 본문은 재감사하지 않았다. 요청/응답·blob 포인터/구간·host normalization·중국어 검색/UTF-8 fixture 관련 테스트를 읽고 직접 문자열 의존을 검색했으며 이번 후보의 기대값 변경 필요는 발견하지 못했다. 세 운영 파일의 원문/현재 줄 수 및 주석/리터럴 밖 텍스트가 동일하고 포맷/인수·키·enum·숫자·이스케이프를 대조했다. 실제 HTTP/외부 오류·중국어 검증 fixture·영어 형식 표식은 보존 대상이다.
- **잔여·한계**: 확정31개 원문 미확보/미독 없음. 일반 로그/API 전용 오류·전체 개발자 주석·테스트 전체 의미·사용자 DB 본문·SDK/모델 동작과 전체 저장소 의미 전수 검증은 미대조/미확인이다. A/S96개·M001–M021 승인/발견 이력·M016 실제 DB/회귀 테스트 후보·M012/M013·A046–A049 및 다른 정책/기능/호환성 후속은 유지한다. 과거 팀원 PASS 기록을 이번 실행 결과로 사용하지 않는다.
- **산출물·실행**: traffic-model-semantic-review.md·traffic-model-semantic-findings.tsv는 Git 제외 translation/audit/에 로컬로만 보관한다. 커밋용 /tmp/artex-traffic-model-semantic-progress.diff는 이번 PROGRESS 기록만 포함한다. 기존 로컬 변경·미추적 파일·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### traffic/ 모델 안내 승인 보완 · M022·M023 [승인 수정·정적 대응 완료 / 실행 미검증]
- **승인 범위**: traffic.go:134 TrafficSearchDescription의 “기록 프록시가 캡처한”→“레코딩 프록시가 캡처한”은 확정 용어 통일로, :2014 traffic_search body_contains schema의 마지막 예시 “내부망 테스트”→“内网测试”는 원문 검색 예시 복원으로 적용했다. M023을 고정 계약 불일치나 테스트 실패로 분류하지 않는다. server/finding_workflow.go의 역사 중국어 legacy 비교값과 SQL/플래그/마이그레이션 동작은 그대로다.
- **정적 보존·테스트**: 현재 상수/Tools schema·문장 경계로 확인했고 정확히 두 표현만 변경했다. traffic.go는 2130/2130줄이며 LF·끝 개행·들여쓰기·지정 표현 밖 바이트·설명/다른 예시·조건·키·숫자·포맷을 보존했다. 관련 traffic/server/agent 테스트에서 두 표현·TrafficSearchDescription 의존을 재검색했다. traffic_store_test.go:122·140의 중국어 원문 본문/검색 fixture는 schema 설명에서 추출하는 값이 아니며 그대로 보존한다. finding_workflow_test.go:249의 호스트/포트 안내 부분 검사는 이번 변경과 무관하다. 직접 기대값 수정 필요는 발견하지 못했으며 입력·fixture·기대값·테스트는 미수정이다.
- **기록·후속**: 로컬 traffic-model-semantic-review.md·traffic-model-semantic-findings.tsv는 수정 전 원문/번역/발견 근거를 이력으로 유지하고 승인 후 표현·분류·정적 대응 완료 상태를 별도 기록한다. A/S 96개·기존 M 이력·M016 실제 DB 상태/회귀 테스트 후보와 정책·기능·호환성 후속은 유지한다. 실제 DB/모델 동작·전체 저장소 의미 전수 대조는 미확인이다.
- **산출물·실행**: /tmp/artex-traffic-model-approved.diff에는 traffic/traffic.go와 누적 PROMPT_PROGRESS 변경만 포함한다. 감사 MD/TSV는 Git 제외 translation/audit/에 로컬로만 보관한다. 기존 미커밋 PROGRESS·TRANSLATION_PROMPT 변경·미추적 파일·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### U15 스킬 본문 의미 감사 [확정 세 문서 전체 정적 대조 완료 / M024 용어 후보 / 실행 미검증]
- **기준·범위**: HEAD 96439f9319f497d81cfd2cc3e4d33e60d5e275b5에서 M022·M023 승인 보완 커밋을 확인했다. staged·MERGE_HEAD·대상 기존 변경 없음. 최신 네 기준 문서와 로컬 감사 이력을 적용하고 api-recon/SKILL.md·reference.md와 scopesentry/SKILL.md의 번역 본문을 첫 줄부터 끝까지 원문과 대조했다. 현재 문서 전체는 425·500·367줄이다. Skill 등록/전달 및 reference 연결 경로를 읽기 전용으로 확인했다.
- **원문·승인 변경**: 70d6125·19c3ceb·41efa3c·729e4f8의 직접 부모에서 각 번역 묶음의 원문을 확보했다. 앞 묶음 fragment·후속 확정 용어, cd23b6d의 절 참조/config 설명 정정, f412bf1의 frontmatter 종료선/name=scopesentry 정정은 승인 이력으로 구분했다. 확정 본문의 원문 미확보/미독 범위는 없다.
- **결과·테스트**: M024는 scopesentry/SKILL.md:168 Mermaid의 子域名入库→서브도메인 저장을 자산 등록 문맥의 확정 표기 서브도메인 등록으로 맞추는 후보다. 원문 전체·현재 번역·근거·후보를 로컬 MD/TSV에 기록했고 소스는 미수정이다. 직접 테스트 기대값 의존은 검색/관련 테스트 읽기에서 발견하지 못했다. 나머지 대조 본문에서 새 의미·조건/지시 강도 오류를 발견하지 못했다. 명령·JSON/DSL·검색값·정규식·placeholder를 산문과 구분했고 fenced 블록/제목-fragment를 정적으로 대조했다. 실제 테스트/렌더 결과는 미확인이다.
- **한계·후속·보존**: A048 로그인/세션 정책 충돌과 A049.2 depth/coverage stub/forward/L2 차이는 원문 문제로 유지하며 A049.1 NEGATIVE_RE DNT와 분리한다. 기존 A/S96개·M001–M023·M012/M013·M016 DB/회귀 테스트 후보·U7/U15 및 다른 후속 보류를 보존한다. 스크립트는 필요한 동작 근거만 읽었고 영어 playwright-cli는 재감사하지 않았다. 실제 SDK/DB/모델/외부 동작 및 전체 저장소 의미 전수 대조는 미확인이다.
- **산출물·실행**: skills-semantic-review.md·skills-semantic-findings.tsv는 Git 제외 translation/audit/에 로컬로만 작성한다. 커밋용 /tmp/artex-skills-semantic-progress.diff에는 이번 PROGRESS 추가만 포함한다. 기존 TRANSLATION_PROMPT 변경·미추적 담당표/기타 파일·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 묶음 진행 없음.


### U15 ScopeSentry 확정 용어 승인 보완 · M024 [승인 수정·정적 대응 완료 / 실행 미검증]
- **승인 범위**: scopesentry/SKILL.md:168의 Mermaid 라벨 서브도메인 저장→서브도메인 등록 한 곳만 수정했다. 자산 등록 문맥의 기존 확정 용어를 적용하며 용어집·코드·테스트는 수정하지 않았다.
- **정적 보존**: 현재 3.3 흐름도의 C 노드로 위치를 확인했다. 367줄·LF·끝 개행·들여쓰기와 지정 표현 밖 바이트가 동일하다. 노드 ID C·화살표·배치·task==阶段1任务名 검색값·다른 라벨·frontmatter/name=scopesentry는 보존했다. 기존 감사에서 직접 테스트 기대값 의존은 발견하지 못했으며 테스트는 미수정·미실행이다.
- **기록·후속**: 로컬 skills-semantic-review.md·skills-semantic-findings.tsv에는 원문·수정 전 번역·발견 근거를 유지하고 승인 후 표현과 정적 대응 완료 상태를 별도로 기록했다. 기존 A/S96개·M 이력 및 A046–A049·M012/M013·M016 등 정책·기능·호환성 후속은 유지한다. 실제 SDK/모델 동작·전체 저장소 의미 전수 검증은 미확인이다.
- **산출물·실행**: 커밋용 /tmp/artex-skills-semantic-approved.diff에는 ScopeSentry SKILL.md와 누적 PROMPT_PROGRESS 변경만 포함한다. 감사 MD/TSV는 Git 제외 translation/audit/에만 로컬 보관한다. 기존 미커밋 PROGRESS·TRANSLATION_PROMPT 변경·미추적 파일·Git index를 보존한다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push 없음.


### 번역 작업 최종 종합 상태 [기록상 승인 번역 잔여 없음 / 원문·호환성 후속 유지 / 실행 미검증]
- **기준·상태**: HEAD 1ef79f0870c738c341703157d67bc10ef22b2d2d의 M024 승인 보완 커밋 1ef79f0 반영과 ScopeSentry C 노드의 서브도메인 등록을 확인했다. 브랜치 work/ko-translation, staged·MERGE_HEAD 없음. fetch·병합하지 않았으며 원격 최신성/push 성공을 별도로 주장하지 않는다.
- **단위·파일 근거**: 최신 기준 문서와 PROGRESS/로컬 감사 기록을 대조했다. 현재 GUIDE 명시/디렉터리 범위는 U1~U15 총 337개 파일이며 로컬 파일별 장부에 번역/정적 검토 보고·대상 없음·개별 근거 부족·의미 대조 범위를 구분했다. U11 테스트 및 U13 미명시 테스트의 전체 검토를 추정하지 않았다. U1 작업 라벨과 U2→U3 artifactSpec 기대값의 과거 이월은 후속 해소 기록을 따르되 다른 검토/실행까지 확대하지 않는다.
- **A/S·M 현재 분류**: A001–A090/S001–S006 96개는 수정·정적 대응 77, DNT/원문 예시 보존 15, 원문 정책·기능·호환성 문제 4(A046–A049)로 유지한다. M001–M024는 승인 수정·정적 대응 22, 원문 설명 문제 2(M012/M013)다. 두 집계는 섞지 않는다. 승인 커밋의 HEAD 포함을 확인했으며 이 ID 범위의 기록상 확인된 미해결 번역·표기 수정 잔여는 0이다. 미발견 오류가 없다는 전수 보증은 아니다.
- **의미 대조 범위·한계**: U1/U2 기본·연결 마무리/U7 세 상수, U3/U4 핵심·보조 안내, U5 231개 자연어/카탈로그 경로, U11 확정166개, U13 traffic31개, U15 번역 세 문서 본문 전체는 각 직접 부모 원문과 대조한 로컬 기록이 있다. 나머지 한국어 문장·전체 주석·테스트·DB 사용자 편집 본문·SDK 안내 및 전체 저장소 의미 전수 대조는 미완료/미확인으로 유지한다. 번역/검색/정적 검토를 본문 의미 전수 대조와 혼동하지 않는다.
- **잔여 분리**: A048 로그인/세션 원문 정책·A049.2 depth/coverage 차이와 A049.1 정규식 DNT, A046 live 중국어 동사·A047 옛 블록/기존 커스텀 정책, M012/M013 설명 문제, M016 실제 DB legacy/flag 및 회귀 테스트 후보, A050.3 mock 범위와 기타 원문 설명 후속을 유지한다. A050 승인13곳/기록 fixture 보존·A052 고유명·A090 설명 수정과 기존 해소 항목은 과거 보류를 현재 잔여로 재집계하지 않는다. 이미 해소한 문서 참조/routes 설명·frontmatter/name·공유 기대값·프런트 용어는 이력으로 보존한다. 원문 문제 해결을 번역 완료의 임의 전제조건으로 지정하지 않는다.
- **산출물·보존·실행**: final-translation-status.md·final-translation-file-status.tsv·final-translation-item-status.tsv는 Git 제외 translation/audit/에 로컬로만 작성한다. /tmp/artex-final-translation-status-progress.diff는 이번 PROGRESS 추가만 포함한다. 소스·용어집·테스트·기능 및 기존 로컬 변경·미추적 파일·Git index를 보존한다. 과거 팀원 PASS는 당시 이력이며 현재 보완의 실행 결과로 인용하지 않는다. 실행 검증은 사용자 요청으로 미실행이다. stage·commit·push·추가 감사 묶음 없음.


### M016 호환성 회귀 테스트 보완 [작성·정적 검토 완료 / 실제 DB 상태 미확인 / 실행 미검증]
- **기준·범위**: HEAD 27e4ff4c526dbb86f95f5f7206ee7e811dbd3c98에서 직전 병합 커밋 반영과 staged·MERGE_HEAD 부재를 확인했다. 최신 기준 및 M016 승인/후속 기록을 적용했고 기존 대상 변경·중복된 역사 fixture 회귀 검증은 없었다. 수정은 finding_workflow_test.go의 traffic import/신규 테스트와 이번 기록뿐이다.
- **직접 검증**: TestFindingWorkflowHostSearchDescriptionMigration의 순차 하위 사례 두 개를 추가했다. A는 bf425ef 직접 부모 42e1deb930e8e519a29380ed85a34cf53e92d63d의 정확한 중국어 역사 설명을 독립 const fixture로 준비하고 v3 host_search=false에서 seed 후 현재 traffic.TrafficSearchDescription과 바이트 정확 일치·완료 flag=true를 검사한다. B는 같은 조건에서 사용자 편집 설명(줄바꿈 포함)을 바이트 그대로 보존하며 현재 구현대로 완료 flag=true가 되는지 검사한다. 운영 legacy 변수 자체로 fixture를 만들지 않았다.
- **격리·기존 검사 보존**: 기존 trafficEvidenceServer/UpdateTool/t.Cleanup 관례로 system 행과 플래그를 준비한다. 사례마다 traffic_search description/schema/agents/enabled 및 v3/v2 플래그의 값·존재 여부를 저장해 cleanup에서 복원한다. v2=true로 별도 reporter/schema/binding 마이그레이션을 격리하고 사용자 schema/바인딩/비활성 설정의 보존도 검사한다. 기존 마이그레이션 테스트와 모든 입력·기대값·검사는 그대로다. updated_at 과거 시각 또는 외부 병렬 프로세스까지 복원/격리한다고 주장하지 않는다.
- **정적 기대·상태**: 독립 fixture는 현재 중국어 비교값과 같고 복원 전 한국어 비교값과 다름을 확인했다. 잘못 번역된 비교값에서는 UPDATE 매칭 없이 flag만 true가 될 수 있어 A의 description 정확 비교가 탐지할 것으로 예상한다. B는 무조건 덮어쓰기 회귀를 탐지한다. 실행 결과는 아니다. M016은 소스 비교값 복원 완료 / 실제 DB 상태 미확인 / 회귀 테스트 작성·정적 검토 완료 / 실행 미검증으로 구분한다. 이미 flag=true인 DB 자동 복구를 검증하지 않는다.
- **보존·산출물·실행**: 테스트는 328→420줄/LF이며 신규 함수/import를 제외한 기존 바이트 동일. A/S96개·M001–M024 승인/원문 문제 집계·기존 정책·기능·호환성 후속을 유지한다. 로컬 감사 m016-regression-review.md는 Git 제외 translation/audit/에만 추가하며 커밋용 /tmp/artex-m016-regression.diff에는 테스트와 PROGRESS만 포함한다. 기존 로컬 변경·미추적 파일·감사·Git index를 보존했다. 실행 검증은 사용자 요청으로 미실행이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### A046 · U7 live 판정 설명의 다국어 회귀 검사 보완 [테스트 작성·정적 검토 완료 / 실제 live 실행 미확인]
- **기준·범위**: HEAD 3f7b39d3d2d495a8da9798014f39410e0495048c의 M016 회귀 테스트 커밋 반영을 확인했다. work/ko-translation, staged·MERGE_HEAD 및 대상 기존 변경 없음. 최신 기준과 A046.1/A046.2를 적용했고 fetch·병합하지 않았다. 수정은 intercept_live_test.go의 두 설명 검사 교체, 신규 intercept_live_reason_test.go의 테스트 전용 검사기·표 기반 단위 테스트 및 이번 기록뿐이다.
- **검사 목적·방법**: ParseVerdict와 같은 실제 앵커로 comment의 实际操作 부분만 분리한다. 한국어·중국어·영어의 제한된 문장 형태를 절 전체로 대조해 현재 파일 읽기 또는 보고서 파일 생성/쓰기/저장 근거를 반드시 요구한다. 과거 생성·명시적 생성/실행 부정은 보조 설명으로만 인정하며 단독으로 통과 근거가 되지 않는다. 읽기 뒤 쓰기·현재 생성까지 포함하는 설명 및 보고서 속 rm/업로드를 실제 실행한 설명은 인정하지 않는다. 알 수 없거나 상충하는 절은 실제 comment·해당 절을 보이는 실패로 남긴다. 자유 자연어 전체를 완벽하게 판별하는 파서는 아니다.
- **단위 검증 작성**: 독립된 정상/오류 쌍을 세 언어로 작성했다. 읽기만/생성 후 읽기, 생성 부정/읽은 뒤 쓰기, 과거 생성·현재 읽기/현재 생성 포함, 실행 부정·보고서 저장/실행 후 저장, 문자열 기록/실제 rm·업로드, 모호·부정·과거만·조건부 설명 및 다른 사유 구간의 동사를 이용한 오판을 포함한다. 동작 근거 표 136개·구간 경계 표 24개 총 160개이며 실행 결과가 아닌 작성한 fixture 수다.
- **기존 검사·구성 보존**: 판정값·설명 필수·입력 필드/이력 제외·Worker 요약 제외·8개 live 입력 및 ARTEX_REVIEW_LIVE_CONFIG opt-in/45초 제한을 보존했다. 검사 교체 밖 기존 바이트 동일, live 파일 137→131줄/LF, 신규 파일 329줄/LF다. 운영 프롬프트·파서·정책·DB·다른 테스트는 변경하지 않았다. 두 파일에 build tag가 없고 신규 파일은 표준 라이브러리만 참조한다. 단위 테스트 본문은 네트워크/모델/DB를 사용하지 않지만 기존 server/TestMain의 패키지 실행 경로는 DB 연결을 시도할 수 있어 파일 지정 독립 실행과 구분한다. 어느 경로도 이번에는 실행하지 않았다.
- **질문형 보완**: 实际操作 구간에 반각 ? 또는 전각 ？가 있으면 문장 분리/끝 문장부호 정리 전에 실제 comment를 인용하는 명확한 실패로 처리한다. 두 문장부호 처리 경로에서도 물음표를 제거했다. 세 언어의 읽기/보고서 저장 × 두 물음표 × 단독/뒤 공백/후속 절 36개 실패 사례와 정상 마침표 12개, 다른 사유 구간의 물음표 6개 정상 사례를 추가했다. 기존 106개 fixture는 바이트 그대로 유지했다. 이번 54개 추가를 포함한 160개는 작성 수이며 실행 통과 수가 아니다.
- **상태·한계·보존**: A046은 테스트 작성·정적 검토 완료 / 실제 live 실행·표현 범위 확인 필요 / 실행 미검증이다. 정상 의미의 미지원 표현도 명시적 실패로 남을 수 있다. A/S96개·M001–M024 집계/승인 이력, M016 실제 DB 상태, A047–A049와 다른 정책·기능·호환성 후속은 유지한다. 로컬 감사 a046-live-reason-review.md는 Git 제외 translation/audit/에만 보관하며 커밋용 /tmp/artex-a046-live-reason.diff에는 허용 세 파일만 포함한다. 기존 로컬 변경·미추적 파일·감사 보고서·Git 제외 설정·Git index를 보존했다. 실행 검증은 **사용자 요청으로 미실행**이다. stage·commit·push·다음 묶음 없음.


### A047 · 사용자 지정 정책 안내·기본 템플릿 복원 확인창 [구현·정적 검토 완료 / 중복 부착 제한 유지 / 실행 미검증]
- **기준·승인 범위**: HEAD 95847dda21b12258c72da6e188e8711d054d05d9, work/ko-translation의 A046 보완 커밋과 staged·MERGE_HEAD 부재 및 대상 기존 변경 부재를 확인했다. 사용자 승인에 따라 system/intercept/page.tsx와 이번 기록만 수정했다. 기존 AlertDialog 구성 요소와 side-question-workspace.tsx의 취소/확인 패턴을 적용했다.
- **안내·확인창**: 편집기에 사용자 지정 정책은 자동 번역되지 않으며 비워서 저장하면 내장 기본 템플릿을 사용한다는 안내와 승인 placeholder를 추가했다. 복원 확인창은 사용자 정책 제거·즉시 저장·화면의 다른 판정 설정 동시 저장·필요한 정책의 사전 복사 보관을 명시한다. 버튼은 취소 / 복원하고 저장이다. 복원 버튼은 확인창만 열며 확인 동작만 기존 restorePrompt를 호출한다. 취소에는 저장 요청·cfg 변경 처리기를 연결하지 않았다.
- **중복 실행·저장 보존**: 요청 중 복원 트리거/확인 버튼은 saving으로 비활성화하고 restoreInFlight ref를 첫 await 전에 설정하여 React 상태 갱신 전 연속 확인 호출도 차단한다. finally에서 잠금을 해제한다. 기존 전체 cfg 저장({ ...cfg, prompt: "" })·조회 후 setCfg·성공/실패 안내·일반 저장·API/백엔드/정책/DB 경로는 유지한다. 자동 번역·옛 블록 교체·해시/바이트 승인 체계·백업 기능은 추가하지 않았다.
- **정적 검토·후속**: 취소/확인/요청 중 경로와 기존 구성 요소의 연결을 정적으로 대조했다. 안내는 미저장 다른 설정도 함께 저장되는 현재 동작에 맞춘 것이다. A047.1/.2 옛 공통 블록 중복 부착은 알려진 호환성 제한 / 미해결로 유지한다. 확인창은 정책 백업·저장 원자성·옛 블록 호환성 해결을 제공하지 않는다. 기존 A/S·M 이력과 A046 live 실행·M016 실제 DB 상태 및 U15 정책/기능 후속을 보존한다.
- **산출물·보존·실행**: 커밋용 /tmp/artex-a047-restore-guidance.diff에는 페이지와 PROMPT_PROGRESS만 포함한다. 구현 상세 보고서는 Git 제외 translation/audit/a047-restore-guidance-implementation-review.md에 로컬로만 보관한다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·로컬 감사·Git index·제외 설정을 보존한다. 실행 검증은 **사용자 요청으로 미실행**이며 빌드·타입 검사·테스트·렌더·앱·Git 훅·fetch·병합·stage·commit·push 없음.


### 알림·보고서 의미 감사 확정 용어 승인 보완 · M025–M028 [승인 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: work/ko-translation, HEAD `b0f523b9bac28a5eb9899713a5b1eb7d5eb75bf2`에서 staged·MERGE_HEAD 및 대상 기존 변경 부재를 확인했다. 로컬 notification-semantic-findings.tsv의 승인 위치를 현재 문장·함수 경계와 대조하여 소스 네 파일의 33개 행만 수정했다. 새 용어 등록·기능/정책/DB 변경은 없다.
- **승인 수정**: db/notification.go의 별도 조회 없이 한 번 삽입하고(M025) 1곳, db/notification_delivery.go의 측정 단위가 다르다(M026) 1곳, 두 DB/서버 파일의 지정 digest 모드·주기·배치·메시지·로그/오류의 모아 보내기(M027) 26곳, server/notifier.go·notify_api.go의 상세 링크(M028) 5곳이다. 기존 상세가 앞에 있으면 회신만 제거하여 상세 상세 링크를 만들지 않았다. 일반 요약·집계·종합 보고서, 영어 digest, Markdown **요약**: 라벨과 server/server.go 범위 밖 후보는 보존했다.
- **직접 의존·정적 보존**: 현재 생성부→FailDeliveries/로그·last_error→DTO/UI 표시 경로와 관련 테스트의 전체/부분 문자열 검사를 재확인했다. 지정 표현의 직접 테스트 기대값 변경 필요는 발견하지 못해 테스트를 수정하지 않았다. 모아 보내기 적용에 따른 조사 을→를 1곳을 함께 정리했다. 네 파일의 줄 수 551/446/555/515·LF·끝 개행·들여쓰기와 지정 표현 및 이 조사 밖 바이트, 코드·키·조건·SQL·계약값·숫자·포맷 지정자/인수·이스케이프를 보존했다. 동작/테스트 통과를 확인한 결과가 아니다.
- **기록·후속**: Git 제외 translation/audit/의 MD·TSV에는 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도 기록했다. 커버리지 장부는 해당 네 파일의 발견 처리 상태만 갱신했다. A/S96개·M001–M024 이력, A047 중복 부착 제한·A048/A049·M012/M013·M016 실제 DB 및 기존 정책/기능/호환성·실행 후속은 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.
- **산출물·보존·실행**: 커밋용 `/tmp/artex-notification-semantic-approved.diff`에는 소스 네 파일과 이번 PROGRESS 추가만 포함한다. 상세 보고서는 로컬 `translation/audit/notification-semantic-approved-review.md`에 보관하며 감사 문서는 커밋 대상에 포함하지 않는다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·다른 변경·Git index·제외 설정을 보존했다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push 없음.

- **M027 표현 보완**: renderBatch 주석 한 줄을 사용자 승인 표현 `// 그 한 건 때문에 모아 보내기 메시지 전체가 누락되지 않게 한다.`로 정리했다. 이전 승인 변경과 코드·포맷·555줄·LF·Git index를 보존했다. 주석 변경이며 테스트/실행 검증은 사용자 요청으로 미실행이다. 커밋용 diff는 기존 소스 네 파일과 PROGRESS만 포함하고 로컬 감사 보고서는 제외한다.


### M029·M030 · 설정 안내 의미 감사 승인 보완 [정적 대응 완료 / 실행 미검증]
- **기준·수정**: HEAD 1220165의 직전 알림 용어 보완 커밋, work/ko-translation 및 staged·MERGE_HEAD 부재를 확인했다. server/intercept.go의 지정 주석은 전역 모델 보완 판정으로, server/server_mgmt.go의 지정 로그는 메모리 상태를 사용해 배리어 복원으로 정리했다. 두 표현만 수정했으며 551/2232줄·LF·지정 표현 밖 바이트·%s/%v·taskID/getErr·조건과 코드를 보존했다.
- **직접 의존**: 지정 표현의 직접 테스트 기대값 의존은 검색에서 발견하지 못했다. TestAbortTaskDeleteUsesPersistedPauseAndQueueState는 일시중지/대기 상태 복원을 검사하며 로그 문자열은 검사하지 않는다. 테스트는 수정하거나 실행하지 않았다.
- **이력·한계·보존**: 로컬 감사 MD/TSV의 원문·수정 전 근거를 유지하고 승인 후 상태를 별도로 기록했으며 커버리지 장부는 해당 두 범위만 갱신했다. M031은 원문 설명 문제로 유지하고 UI·검증 코드·용어집은 변경하지 않았다. 기존 A/S 집계·M 이력·정책/기능/호환성 후속, 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다. 감사 문서는 Git 제외 translation/audit/에만 보관한다. 커밋용 /tmp/artex-settings-semantic-approved.diff에는 소스 두 파일과 이번 기록만 포함한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push 없음.


### M032 · 보조 질문 고정 안내 확정 용어 승인 보완 [정적 대응 완료 / 실행 미검증]
- **기준·수정**: HEAD a27a27c의 M029/M030 승인 커밋 반영과 staged·MERGE_HEAD 및 대상 기존 변경 부재를 확인했다. sidequestion/request.go:33 instruction의 조작을 실행하거나→동작을 실행하거나 한 표현만 수정했다. 지시 강도·금지/조건·코드·영어·포맷·들여쓰기·120줄·LF·끝 개행 및 지정 표현 밖 바이트를 보존했다. 용어집·테스트는 변경하지 않았다.
- **직접 의존·기록**: 지정 표현의 직접 테스트 기대값 의존은 재검색에서 발견하지 못했다. sidequestion/sidequestion_test.go:180과 agent/side_questions_test.go:107은 별도 service 반환 안내를 검사하며 현재 출력과 일치한다. 로컬 감사 MD/TSV의 수정 전 원문·번역·근거를 유지하고 승인 후 상태를 별도로 기록했다. 커버리지 장부는 request.go 해당 범위만 갱신했다. 기존 A/S 집계·M 이력·원문 문제·정책/기능/호환성 후속은 유지한다.
- **보존·산출물·실행**: 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다. 감사 문서는 Git 제외 translation/audit/에만 보관하고 /tmp/artex-m032-approved.diff에는 sidequestion/request.go와 이번 PROGRESS 추가만 포함한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push 없음.

### M033–M036 · 작업·채팅 UI 의미 감사 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: HEAD 954d1cd의 M032 승인 커밋 반영, work/ko-translation, staged·MERGE_HEAD 및 대상 기존 변경 부재를 확인했다. 작업 목록·개요·브로드캐스트 세 파일에서 지정 표현 6곳만 수정했다.
- **수정**: M033 archiveBlockReason에 직접 상속 관계와 의존 작업 우선 아카이브 순서를 명시하고 변수를 유지했다. M034 상태 필터 라벨 paused/done/timeout은 각각 일시 중지됨/완료/시간 초과로 맞췄다. 같은 파일의 추가 일시정지 표현 후보는 이번에 수정하지 않았다. M035 지정 placeholder의 월권만 권한 우회로, M036 지정 주석의 탐색 체인 그래프만 탐색 그래프로 변경했다.
- **정적 보존·직접 의존**: 세 파일의 3417/1322/734줄·LF·끝 개행과 승인 표현 밖 바이트를 역치환 비교로 확인했다. JSX·enum·필터 조건·목표 입력/등록·코드·키·변수·포맷과 인수는 유지했다. 지정 표현을 직접 검사하는 테스트 기대값은 재검색에서 발견하지 못했다. 테스트·용어집·기능·정책은 변경하지 않았다.
- **기록·후속·산출물**: Git 제외 translation/audit/의 원문·수정 전 번역·근거를 이력으로 유지하고 승인 후 상태를 별도 추가했다. 커버리지는 세 파일의 해당 범위만 갱신했다. 기존 A/S 96개·M 이력·원문 문제·정책/기능/호환성 보류를 유지한다. 커밋용 /tmp/artex-chat-tasks-semantic-approved.diff에는 소스 세 파일과 이번 PROGRESS 추가만 포함하며 로컬 감사 보고서는 제외한다.
- **실행·보존**: 기존 TRANSLATION_PROMPT 변경·미추적 파일·다른 변경·Git index·Git 제외 설정을 보존했다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push 없음.


### M037–M039 · 자산·취약점 UI 확정 용어 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation / upstream origin/ko-translation, HEAD 1159c85에서 M033–M036 보완·병합 반영, staged·MERGE_HEAD 및 대상 기존 변경 부재를 확인했다. 승인된 소스 세 파일의 네 표현만 수정했다.
- **수정**: assets/page.tsx:593 서비스 헤더 지문→핑거프린트(M037), :464 IP 목록 오픈 포트→열린 포트(M038), findings/page.tsx:756 취약점 통계 라벨과 types.ts:472 취약점 처리 상태 주석의 대기 중→처리 대기(M039). stats.pending·FindingStatus enum·필터·사용자 데이터는 유지하고, 자산 DSL의 기술 지문 및 다른 작업·큐의 대기 중은 변경하지 않았다. 용어집 추가 없음.
- **직접 의존·보존**: 표시 헤더·statCards와 DB pending 집계/기존 status 라벨 연결을 대조했다. 해당 네 표현의 직접 테스트 기대값 의존은 재검색에서 발견하지 못해 테스트를 수정하지 않았다. 세 소스의 1216/1263/1596줄·LF·끝 개행·JSX·코드·계약값 및 지정 표현 밖 바이트를 역치환 비교로 확인했다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index·제외 설정을 보존했다.
- **기록·한계**: 로컬 감사 MD/TSV에는 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도로 추가했다. 커버리지는 지정 세 파일/네 위치만 갱신하며 기존 A/S96개·M 승인 이력·원문 문제·정책/기능/호환성 후속을 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다. 감사 산출물은 Git 제외 translation/audit/에만 보관한다. 커밋용 /tmp/artex-m037-m039-approved.diff에는 소스 세 파일과 이번 PROGRESS 기록만 포함한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 묶음 없음.


### M040–M042 · 작업 제어 확정 용어 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD 4620fe2의 M037–M039 승인 커밋 반영과 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. 지정 10곳(db/tasks.go 2, server/task_control.go 7, server/task_archives.go 1) 및 조건부 승인 프런트 주석 2곳을 수정했다.
- **수정·추가 근거**: planner 하트비트→heartbeat 2곳, 일시정지→일시 중지 7곳(기존 활용형 유지), 가짜 삭제→논리 삭제 1곳. 추가 api.ts:344/types.ts:424 원문은 각각 번역 직접 부모 31e401b/b45957e의 假删除이며 deleted/delete_reason 문맥이다. applyIntentControl의 soft 경로와 db/exploration.go:389 SoftDeleteIntent는 의도 노드·계보를 보존하는 논리 삭제로 연결된다. 따라서 추가 2곳도 논리 삭제로 정리했으며 stopped/软删除 문맥은 변경하지 않았다.
- **직접 의존·보존**: 지정 표시 오류/로그를 직접 비교하는 테스트 기대값 의존은 재검색에서 발견하지 못했다. task_control_routes_test·task_archives_test·intent_delete_test의 상태/노드 보존 검사는 유지하며 테스트는 변경하지 않았다. 소스 654/288/728/1314/1596줄·LF·끝 개행 및 지정 표현 밖 바이트는 역치환 비교로 동일함을 확인했다. soft·paused·deleted·delete_reason·코드·조건·포맷·인수를 유지하고 M043의 300/600·5min/10min은 원문 설명 후속으로 그대로 보존한다. 용어집 수정 없음.
- **기록·한계**: 로컬 감사에는 수정 전 원문·번역·근거를 유지하고 승인 후 상태를 별도로 추가했다. 커버리지는 지정 5개 파일의 12곳만 갱신하며 기존 A/S 96개·M 이력·정책/기능/호환성 보류는 유지한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다. 감사 MD/TSV는 Git 제외 translation/audit/에만 보관한다. 커밋용 /tmp/artex-m040-m042-approved.diff는 소스 5개와 이번 PROGRESS 기록만 포함한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### M044–M046 · 자산 규칙·동기화 확정 용어 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD 033716f의 M040–M042 승인 커밋과 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. 소스 세 파일의 지정 11행/14표현만 수정했다.
- **승인 수정**: db/asset_intercept_match.go:12·35·37·44·103·118·179·180·197의 적중 12표현을 매칭으로 정리했다. 12의 매칭되면, 44/103의 매칭되는지·매칭되는/매칭된, 118의 매칭될, 180의 매칭되지 등 승인 활용형을 적용했다. server/task_intercept.go:11 action=block은 차단으로, server/sync_scopesentry.go:27 상한 주석은 자산 등록 보호 상한으로 정리했다. 기능명 인터셉트·block/allow 우선순위·조건·반환값·5000은 보존했다.
- **직접 의존·보존**: Reason→EvaluateAssetGate→agent 도구의 오류/설명 연결을 재확인했다. 지정 사유 문자열을 비교/파싱하는 소비부 및 직접 테스트 기대값 의존은 재검색에서 발견하지 못했다. asset_intercept_match_test는 Allowed·매칭 값·Reason 비어 있지 않음을 검사하며 유지했다. 세 파일의 250/189/513줄·LF·끝 개행·포맷 지정자/인수와 지정 표현 밖 바이트를 역치환 비교로 확인했다. 테스트 사례 이름·실패 안내의 추가 적중 후보, 센티넬·최종 방어선 문맥은 변경하지 않았다. 용어집 수정 없음.
- **기록·후속**: Git 제외 translation/audit/의 수정 전 원문·번역·근거를 유지하고 승인 후 상태를 별도 추가했다. 커버리지는 해당 세 파일의 지정 11행만 갱신하고 기존 A/S96·M 이력·원문 정책/기능/호환성 보류를 유지한다. 기존 변경·미추적 파일·Git index·제외 설정은 보존했다. 커밋용 /tmp/artex-m044-m046-approved.diff에는 소스 세 파일과 이번 PROGRESS 추가만 포함한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### M047–M053 · 실행·설정·등록 주석 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD 4a99b7b의 M044–M046 승인 커밋 반영과 staged·MERGE_HEAD 및 대상 기존 변경 부재를 확인했다. 승인된 소스 6개 파일의 7개 행을 수정했다. db/llmretry.go:90에는 M048·M053을 함께 적용했다.
- **승인 수정**: constants.go:28 논리 삭제(soft delete), db/llmretry.go:90 키 무효·고정 재시도 대기 시간, engine_timeout.go:21 정상적으로 마무리하는, server/llmretry.go:86 승인 괄호 문구와 :93 진행 없음 시 이어 실행 끔, task_llm.go:585 cold 상태의 노드, customtool.go:264 얇은 래퍼 도구로 정리했다. M049는 정리 절차를 뜻하며 목표 달성·성공 종료의 의미를 추가하지 않았다.
- **용어·직접 의존**: 薄壳工具의 중복·충돌을 확인한 뒤 custom tool의 빈 schema에 최소 args 포장을 제공하는 문맥에 한정해 용어집에 확정 등록했다. API 정찰의 包装层 적용 범위는 유지한다. TestEnsureSchema의 args/target 구조 검사, 삭제 상태/노드 보존, 재시도 설정/무진행 턴 검사를 읽기 전용으로 재대조했으며 지정 주석을 직접 검사하는 기대값 변경 필요는 발견하지 못했다. 테스트는 수정하지 않았다.
- **정적 보존·기록**: 소스 줄 수·LF·끝 개행·들여쓰기·영어·코드·조건·숫자·포맷 및 지정 표현 밖 바이트를 역치환으로 확인했다. 소스의 변경은 주석에 한정된다. 로컬 감사 MD/TSV에는 원문·수정 전 번역·발견 근거를 보존하고 승인 후 상태를 별도로 기록했다. 커버리지는 해당 6개 파일/7개 행만 갱신하며 추가 후보와 다른 파일의 동일 표현은 후속 목록에 유지한다. 기존 A/S96개·M 승인 이력·원문 문제·정책/기능/호환성 보류는 유지한다.
- **산출물·보존·실행**: 커밋용 /tmp/artex-m047-m053-approved.diff에는 소스 6개·GLOSSARY·이번 PROGRESS 기록만 포함한다. 상세 보고서는 Git 제외 translation/audit/execution-storage-semantic-approved-review.md에 로컬로만 보관한다. 기존 변경·미추적 파일·Git index·제외 설정은 보존한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### M054–M064 · 엔진·오케스트레이션·탐색 그래프 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD b31dd15592817fbf4e396e60af70b20c8961b3c8의 M047–M053 승인 커밋 반영과 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. 로컬 engine-graph-semantic-findings.tsv의 지정 소스 8개·47개 항목 위치를 현재 문장과 함수 경계로 재확인했다. db/exploration.go:385의 M056·M060을 함께 적용하여 실제 변경은 46개 소스 행이다.
- **승인 수정**: M054 무진행 턴/진행 없음 및 진행 없음 시 이어 실행·진행이 없는 상황의 활용형과 턴이 조사, M055 heartbeat, M056 일시 중지, M057 정상적인 마무리, M058 "%w: 의도 %d: %s 동작을 실행 중입니다", M059 "취소된 의도 #%d의 Token 계량", M060 지정 假删除 문맥의 논리 삭제, M061 LLM failover의 순환 전환, M062 목표 달성의 정량적 확인, M063 고정 라우팅, M064 이전 기록 조회를 적용했다. 실제 폴링·다른 문맥·기존 역사 비교 문자열은 변경하지 않았다. 새 용어 등록 없음.
- **직접 의존·보존**: TestControlWorkRejectsConcurrentController의 errors.Is 기반 충돌 검사와 db/intent_control_test.go의 Token 사용량·일자별 집계·rollup 보존 검사를 읽기 전용으로 재확인했다. 지정 오류/로그/제목을 직접 비교하는 테스트 기대값 또는 프런트 매칭 의존은 검색에서 발견하지 못하여 테스트를 수정하지 않았다. source 8개 줄 수·LF·끝 개행·들여쓰기·코드·조건·상태/SQL/JSON 계약값·포맷 지정자 종류/순서/개수 및 인수를 보존하고 지정 행 밖 바이트 동일성을 정적으로 확인했다. Token 계량·저장·집계 동작과 DB 기존 문자열은 변경하지 않았다.
- **기록·후속**: Git 제외 translation/audit/의 원문·수정 전 번역·발견 근거는 보존하며 승인 후 문구·정적 대응 상태를 별도 기록했다. 커버리지는 해당 8개 파일의 지정 46행만 갱신했다. 기존 A/S96개·M001–M053 승인 이력, M012·M013·M043 원문 설명 문제, A047 중복 부착 제한·A048/A049·M016 실제 DB 미확인 및 다른 정책/기능/호환성·실행 후속을 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.
- **산출물·보존·실행**: 커밋용 /tmp/artex-m054-m064-approved.diff에는 소스 8개와 이번 PROGRESS 추가만 포함한다. 로컬 보고서는 translation/audit/engine-graph-semantic-approved-review.md에 보관하고 감사 MD/TSV는 커밋용 diff에서 제외한다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index·제외 설정은 보존한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### M065–M073 · 서버 공통 안내·설정 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD b9944efad0a9a26f6f45fc1e76607fba8320a936의 M054–M064 보완 커밋 반영과 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. 로컬 remaining-candidates-semantic-findings.tsv에서 승인된 server/server.go의 22개 위치를 현재 문장·함수 경계로 확인하여 후보대로 수정했다. 실제 polling·일반 데이터 전달·요약·공격 인젝션 등 다른 문맥은 변경하지 않았다.
- **승인 수정**: M065 LLM 순차 시도/실패 전환 문맥의 순환 전환 6행, M066 알림 전송 엔진 3행, M067 LLM 보완 판정 1행, M068 논리 삭제 1행, M069 일시 중지 3행, M070 heartbeat 1행, M071 모아 보내기 주기 2행, M072 상세 링크 3행, M073 의도 추가 2행이다. 기존 intent/hint 입력 예시와 意图/提示 HasPrefix·TrimPrefix 파싱 계약은 그대로 보존했다. 새 용어 등록·기능·정책·DB 변경 없음.
- **직접 의존·정적 보존**: 오류 안내 3428/3439 두 곳과 지정 표시 표현의 전체/부분 문자열 비교 의존을 추적 파일에서 재검색했다. TestNotifyMetaAndSettingsRoundTrip의 설정값·끝 슬래시 정규화·HTTP400 검사와 직접 생성부를 읽기 전용으로 대조했으며, 수정 문구의 직접 테스트 기대값 변경 필요는 발견하지 못하여 테스트를 수정하지 않았다. 3907줄·LF·끝 개행·들여쓰기·코드·조건·enum/키/계약값·숫자·포맷 지정자 종류/순서/개수·인수·이스케이프 및 지정 표현 밖 바이트를 보존했다.
- **기록·후속**: Git 제외 translation/audit/의 기존 원문·수정 전 번역·발견 근거를 이력으로 보존하며 승인 후 상태를 별도 기록하고 커버리지는 해당 파일/22행만 갱신했다. 기존 A/S96개·M 승인 이력, M012/M013/M043 및 원문 설명 문제, A047 중복 부착 제한·A048/A049·M016 실제 DB 상태 미확인과 다른 정책/기능/호환성·실행 후속은 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.
- **산출물·보존·실행**: 커밋용 /tmp/artex-m065-m073-approved.diff에는 server/server.go와 이번 PROGRESS 추가만 포함하며 감사 MD/TSV는 제외한다. 로컬 상세 보고서는 translation/audit/remaining-candidates-semantic-approved-review.md에 보관한다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index·제외 설정은 보존한다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 감사 없음.


### M074·M075 · 설치·업데이트 안내 의미 감사 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·수정 범위**: ko-translation, HEAD d571c9a92df9d9c72aeeb39cb78f4b5483c7b242에서 staged·MERGE_HEAD 및 대상 파일 기존 변경 부재를 확인했다. selfupdate/selfupdate.go:140에 원문의 “0.3.7로 취급해서는 안 된다” 부정 절을 복원(M074)하고, selfupdate/stage.go:46의 종료 절차를 “정상 종료 절차를 수행하고 ExitRestart로 나간다”로 정리(M075)했다. 사용자 지정 주석 두 행만 수정했으며 원문 자체의 반환=완료 설명은 유지했다.
- **정적 보존·직접 의존**: 두 파일의 175/339줄·LF·끝 개행·들여쓰기·앞뒤 주석·코드·버전 숫자·ExitRestart 및 지정 표현 밖 바이트를 역치환 비교로 확인했다. TestCompareVersions의 개발 버전 비교 불가 사례 및 server/update.go의 Stage 성공 후 requestRestart 경로를 기존 감사 근거와 재대조했다. 지정 주석을 직접 검사하는 테스트 기대값 의존은 재검색에서 발견하지 못했으며 테스트·용어집은 수정하지 않았다.
- **기록·후속·보존**: 로컬 감사 MD/TSV의 중국어 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도 필드/기록으로 추가했다. 커버리지는 해당 두 주석만 갱신하고 A/S 집계·M 승인 이력·원문 설명 문제·정책/기능/DB/커스텀 프롬프트/live 후속을 유지한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다. 커밋용 /tmp/artex-m074-m075-approved.diff에는 소스 두 파일과 이번 PROGRESS 기록만 포함하고 감사 MD/TSV는 Git 제외 translation/audit/에 로컬로만 보관한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 묶음 없음.


### M076·M077 · guard 용어 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·수정 범위**: ko-translation, HEAD b0aa6bc4b7d54f30e213e3a13e974ebf434945b7의 M074·M075 커밋 반영 및 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. guard/guard.go:103의 “deny 적중은”을 “deny 매칭은”으로, :140의 “이 작업은 금지됩니다.”를 “이 동작은 금지됩니다.”로 승인된 두 표현만 수정했다. 기존 확정 용어를 적용하며 용어집·테스트는 수정하지 않았다.
- **정적 보존·직접 의존**: 197줄·LF·끝 개행·들여쓰기와 지정 표현 밖 바이트를 역치환 비교로 확인했다. reason 전달·차단 정책·승인 처리·hook.Result 반환 구조·코드·계약값을 보존했다. 지정 문구의 직접 문자열 비교·테스트 기대값 의존은 재검색에서 발견하지 못했다. guard/guard_test.go의 passthrough 검사는 기존 근거를 재사용하며 테스트를 실행하지 않았다.
- **기록·보존**: Git 제외 translation/audit/에 중국어 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도 기록했다. 커버리지는 해당 두 표현만 갱신했다. 기존 A/S 집계·M 승인 이력·원문 문제 및 정책·기능·호환성 후속을 유지한다. 기존 TRANSLATION_PROMPT 변경·미추적 파일·Git index·제외 설정을 보존했다. 커밋용 /tmp/artex-m076-m077-approved.diff에는 guard/guard.go와 이번 PROGRESS 추가만 포함한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push·다음 묶음 없음.


### M078–M083 · 공유 프런트·mock 주석 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD 3e3caac1edfc8b21b186f636b9c9c1602d869a6b의 M076·M077 보완 커밋 반영 및 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. 로컬 findings TSV와 현재 문장을 대조해 소스 3개 파일의 지정 주석만 수정했다. 집계는 **6개 ID / 항목별 16개 행 / 고유 소스 위치 15개**이며 M079·M083의 types.ts:937은 함께 반영했다. M084 별도 위치는 수정 대상 수에 포함하지 않는다.
- **승인 수정**: M078 테스트 메시지의 동기적으로 전송, M079 지정 LLM 순차 시도·실패 전환의 순환 전환, M080 key 무효, M081 모델 보완 판정, M082 회사 계층의 귀속, M083 지정 실패 처리의 대체 처리를 적용했다. types.ts:937은 승인된 전체 문장으로 정리했다. 실제 polling·다른 fallback 문맥·일반 데이터·입력·fixture는 보존했다. 용어집·테스트·mock/data.ts·status.ts·next.config.mjs는 변경하지 않았다.
- **정적 보존·직접 의존**: api.ts/types.ts/mock/handler.ts의 1314/1596/2484줄·LF·끝 개행·들여쓰기 및 지정 주석 밖 바이트를 행별 비교로 확인했다. 주석 앞 코드·조건·enum·키·반환 구조·포맷·입력·기대값은 유지했다. 지정 주석의 직접 소비부와 테스트 기대값 의존은 추적 파일 및 프런트 테스트 재검색에서 발견하지 못해 추가 기대값 변경은 하지 않았다.
- **기록·후속**: Git 제외 translation/audit/의 원문·수정 전 번역·발견 근거는 보존하고 승인 후 상태를 별도 기록했다. 커버리지는 해당 3개 파일/15개 위치만 갱신했다. M084는 원문 설명 문제로 유지하고 소스·검증 로직은 변경하지 않았다. 기존 A/S96개 집계·M 이력·A047 중복 부착 제한·A048/A049 및 정책/기능/호환성 후속은 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.
- **산출물·보존·실행**: 커밋용 /tmp/artex-m078-m083-approved.diff에는 소스 3개와 이번 PROGRESS 추가만 포함한다. 상세 보고서는 translation/audit/shared-frontend-mock-semantic-approved-review.md에 로컬로만 보관한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정은 보존했다. 실행 검증은 **사용자 요청으로 미실행**이며 fetch·병합·stage·commit·push·다음 묶음 없음.


### M085–M088 · README 의미·용어 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD e9507b3bff755b496f0a22266e4df4b58df434f7의 M078–M083 보완 커밋 및 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. README의 지정 5곳만 수정했다. M085 프런트 개발 설명은 Hot Reload, M086 engine 수명주기는 일시 중지, M087 불법 사용 금지는 갈취, M088 키워드 검색 결과는 “검색에 매칭된 항목”·“A/B의 매칭된 단계(자기 제외)”로 정리했다.
- **용어·보존**: 중복·충돌 확인 후 용어집의 기존 命中=매칭 항목에 키워드 검색 결과 문맥만 보충하고 勒索=갈취는 불법 사용 금지 안내 문맥에 한정해 등록했다. 기존 다른 문맥의 확정 용어는 유지했다. README 537줄·LF·끝 개행·명령·경로·강조·Mermaid ID/화살표/구조 및 지정 표현 밖 바이트를 보존했다. 네 스크립트·코드·테스트는 수정하지 않았다. 직접 테스트 기대값 변경 필요는 재검색에서 발견하지 못했다.
- **기록·후속**: Git 제외 translation/audit/의 원문·수정 전 번역·발견 근거를 유지하며 승인 후 상태를 별도 기록했고 커버리지는 README의 해당 5곳만 갱신했다. 기존 A/S 집계·M 이력·A048/A049 및 정책/기능/호환성 후속은 유지한다. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다. 기존 로컬 변경·미추적 파일·Git index·제외 설정은 보존했다. 커밋용 /tmp/artex-m085-m088-approved.diff에는 README·GLOSSARY·이번 PROGRESS 추가만 포함하고 감사 보고서는 제외한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push 없음.


### M089–M093 · 알림 테스트 의미 감사 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD aa5c082b443370cf89d69abbec1b918e544296af에서 staged·MERGE_HEAD·대상 기존 변경 부재를 확인했다. db/notification_test.go 7곳, notify/mask_test.go 1곳, server/notify_api_test.go 2곳의 승인된 주석·실패 안내만 수정했다(5개 ID / 고유 10개 위치).
- **승인 결과**: M089 채널 규칙의 매칭되지 않은, M090 알림 digest의 모아 보내기, M091 사용자가 필요한 조치를 취할 수 있도록 안내하는 전용 오류 타입, M092 상태 변경 → 수정 완료 인용, M093 채널별 토큰 버킷을 적용했다. 필요한 조사만 정리하고 일반 요약·잔여 토큰량 설명은 보존했다.
- **정적 검토·보존**: 세 테스트의 874/403/951줄·LF·들여쓰기·끝 개행 및 지정 표현 밖 바이트를 역치환 대조했다. 입력·fixture·기대값·검사 조건·로직·포맷 지정자·인수는 변경하지 않았다. M092는 현재 Contains 검사 및 알림 생성 출력과 대응하며 추가 기대값 수정 필요는 발견하지 못했다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존한다.
- **기록·후속**: Git 제외 translation/audit/에 중국어 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도 기록했다. 커버리지는 지정 10곳만 갱신하며 기존 A/S 집계·M 승인 이력 및 정책·기능·호환성 후속은 유지한다. M094는 원문 설명 문제로 남기고 수정하지 않았다. 커밋용 /tmp/artex-m089-m093-approved.diff에는 테스트 세 파일과 이번 PROGRESS 추가만 포함한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push 없음. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.


### M095–M099 · agent 잔여 자연어 의미 감사 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD 5cf6f569bed2d2520cb83e87619bd4f31b5efbe9의 원격 병합 커밋 반영과 대상 기존 변경 부재를 확인했다. agent/planner.go 4곳, mainagent.go 2곳, tools.go 1곳, wrapup.go 1곳의 승인된 8개 표현만 수정했다.
- **승인 결과**: M095 planner 정기 점검 안내·주석의 heartbeat 4곳, M096 목표 추가 재개 주석의 일시 중지된, M097 WebFetch 주석의 레코딩 프록시, M098 회사 범위 주석의 귀속 범위, M099 clamped=false 주석의 작업 종료까지 아직 시간이 남음을 적용했다. mainagent 기본 본문의 다른 일시정지 표현과 지정 범위 밖 표현은 유지한다. 용어집 수정 없음.
- **정적 검토·보존**: 소스 4개의 줄 수·LF·끝 개행·들여쓰기 및 지정 표현 밖 바이트를 역치환 비교로 확인했다. 코드·조건·지시 강도·영어·계약값·포맷/인수는 보존했다. 지정 문구의 직접 테스트 기대값 의존은 재검색에서 발견하지 못해 테스트를 수정하지 않았다. 기존 로컬 변경·미추적 파일·Git index·제외 설정은 보존한다.
- **기록·후속**: Git 제외 translation/audit/의 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도로 기록했다. 커버리지는 해당 8곳만 갱신했다. M100은 추가 확인·미수정으로 유지하며 M094·기존 A/S 집계·M 이력·정책/기능/호환성 후속은 유지한다. 커밋용 /tmp/artex-m095-m099-approved.diff에는 소스 4개와 이번 PROGRESS 기록만 포함하며 감사 MD/TSV는 제외한다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push 없음. 전체 저장소 의미 전수 검증 완료를 뜻하지 않는다.


### M101–M103 · M044 테스트 후속 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `701b79cc4125ea08617bfdde00d524236c948b99`. 대상 기존 변경·staged·MERGE_HEAD 없음과 최신 기준/감사 원문을 확인했다. db/asset_intercept_match_test.go 28행, company_scope_test.go 1행, intercept_seed_test.go 1행만 수정했다. **항목별 34행 / 고유 소스 위치 30곳**이며 M101/M044가 겹치는 4행은 한 번에 처리했다.
- **승인 결과**: M101의 사례 이름 4곳에 부분 일치, M102 설명 주석에 ICP 등록(备案) 번호, M103 진단에 매칭되어야 하는데 허용됨을 적용했다. M044의 테스트 후속 28행은 원문과 검사 문맥에 맞춰 매칭/매칭되지 않음·매칭되면 등 조사·활용형을 정리했다. 기존 운영 소스의 M044 승인 결과와 이번 테스트 후속을 구분하며 용어집은 수정하지 않았다.
- **정적 검토·보존**: 테스트 3개 129/405/46줄·LF·끝 개행·들여쓰기, 입력/fixture/기대값/bool/enum/SQL/검사 로직, 포맷 지정자와 인수를 유지했다. M103의 %s와 인수 s도 동일하다. 승인 위치를 수정 전 행으로 되돌린 정적 텍스트 비교가 작업 전 바이트와 일치한다. 사례 이름의 정의 밖 직접 선택 참조는 현재 추적 파일 검색에서 발견하지 못했으며 저장소 밖 개인 -run 명령은 미확인이다.
- **기록·후속**: Git 제외 translation/audit/의 원문·수정 전 근거는 유지하고 승인 후 상태/30곳 대응표를 별도로 추가했다. 커버리지는 해당 3개 테스트의 승인 범위만 갱신했다. db/asset_intercept_match.go:24/26/28/35의 퍼지는 원문·문맥을 추가 확인 후보로만 기록하며 미수정이다. 기존 A/S 집계·M 이력·M100·원문 정책/기능/호환성 보류와 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존한다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 커밋용 `/tmp/artex-m101-m103-m044-tests-approved.diff`에는 테스트 3개와 이번 PROGRESS만 포함한다. 감사 MD/TSV는 로컬 제외 경로에만 보관한다. fetch·병합·stage·commit·push 없음.


### M104–M106 · 삭제·검색 진단·조사 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `99998ac3c96acb616e1a37632235d2687de4942d`. 직전 M101–M103·M044 테스트 후속 커밋 반영과 대상 기존 변경·staged·MERGE_HEAD 부재를 확인했다. 로컬 감사 findings TSV의 승인 18곳만 수정했다: intent_delete_test 3곳, selfupdate_test 1곳, traffic_store_test 13곳, upgrade_test 1곳.
- **승인 결과**: M104의 假删除/deleted/delete_reason·노드 보존 주석 3곳은 논리 삭제, M105 검색 결과 진단 12곳은 매칭/매칭되어야 함, M106 진단 3곳은 고정 명사/콜론으로 정리했다. SoftDeleteIntent 식별자·다른 软删除 문맥·검색어·필터·fixture·기대값·검사 조건은 보존했다. 용어집은 수정하지 않았다.
- **정적 검토·보존**: 테스트 4개 170/459/528/164줄·LF·끝 개행·들여쓰기와 지정 표현 밖 바이트를 승인 행 역치환 비교로 확인했다. 포맷 지정자의 종류·개수·순서와 name/path/p.Dir·data/marker·n/got 등 인수는 동일하다. 진단 문자열의 저장소 소스/테스트 직접 비교·파싱·기대값 의존은 검색에서 발견하지 못해 추가 테스트 변경은 하지 않았다.
- **기록·후속**: Git 제외 translation/audit/에서 원문·수정 전 번역·발견 근거를 유지하고 승인 후 상태를 별도로 기록했다. 커버리지는 해당 테스트 4개/18곳만 갱신했다. M107 원문 설명 문제와 server/update_test.go, M100 추가 확인, 운영 퍼지 4곳 및 기존 A/S 집계·M 이력·정책/기능/호환성 보류는 유지한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 커밋용 `/tmp/artex-m104-m106-approved.diff`에는 테스트 4개와 이번 PROGRESS만 포함한다. 감사 MD/TSV는 Git 제외 경로에만 보관한다. fetch·병합·stage·commit·push 없음.


### M108–M109 · 증거 연계 문서·현재 UI 인용 누락 보완 [번역·정적 검토 / 용어 문장 보류 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `143589a835d1c93ba8c046603b4b15af6c5e65aa`. 이후 커밋 및 대상 기존 변경·staged·MERGE_HEAD 없음. docs/漏洞流量证据.md 전체 99줄을 읽고 제목·산문·표 설명 38행을 번역했다(일부 행에는 보류 원문 문장 유지). server/llmpool_test.go:5의 현재 UI 인용 한 곳은 순환 전환 순서로 정리했다. 파일명은 유지한다.
- **의미·계약**: 독립 finding_id/탐색 finding_node_id 및 update_finding_report의 기존 노드 ID, 선택적 바인딩·실패 시 롤백·상속 읽기 전용·최신 증거 버전 조건을 보존했다. JSON/text/sh 코드 블록·명령·SQL 키·enum·숫자·인라인 코드·URL·경로는 원문 그대로다. 설명용 ID/해시 접두사 placeholder는 실제 값이 아님을 기록하고 이번에는 원문 표기를 보존했다. 설정 제목은 현재 UI의 Agent 트래픽 자동 연결 표기를 인용한다.
- **보류·연결**: H01–H07의 8문장은 미등록 전문 용어/기존 문맥 제한 때문에 원문 보류하고 전체 원문·위치·문맥·후보를 로컬 보고서에 모았다. 기존 달성 검증/cold/연쇄 삭제의 제한 문맥을 확대하지 않았으며 용어집 수정 없음. 변경 제목으로 연결된 추적 파일의 직접 fragment 참조는 검색에서 발견하지 못했다. 테스트 주석 인용 외 입력·기대값·검사 로직은 동일하고, 추가 테스트 기대값 수정 필요는 발견하지 못했다.
- **정적 보존·후속**: 문서 99줄과 테스트 41줄의 LF·들여쓰기·끝 개행을 보존했다. 변경되지 않은 줄·코드 블록·인라인 코드·숫자·URL 및 테스트의 지정 인용 밖 바이트를 대조했다. M100·운영 퍼지 4곳·M107 및 A046–A049/M016 등 기존 원문 정책·기능·호환성 보류는 유지한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존한다. 전체 문서 번역 완료나 전체 저장소 의미 전수 검증 완료로 표시하지 않는다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 문서 속 테스트/빌드/DB 지시는 검토용 텍스트로만 읽었다. 감사 MD/TSV는 Git 제외 translation/audit/에 로컬로만 보관하고, `/tmp/artex-evidence-doc-omissions.diff`에는 허용 문서·테스트·이번 PROGRESS만 포함한다. fetch·병합·stage·commit·push 없음.


### M108 · H01–H07 증거 문서 용어 보류 해소 [번역·정적 대응 완료 / 실행 미검증]
- **범위·결과**: HEAD `143589a835d1c93ba8c046603b4b15af6c5e65aa`의 기존 미커밋 번역을 보존하고, H01–H07 8문장과 64행 메타데이터 제출 표현 한 곳만 보완했다. 영구 저장·연쇄 제거·정리 프로세스·활성 데이터·모델 fixture·종단 간 검증·활성 본문 데이터·전역 증거 조율 잠금을 적용했다. 이번 증거 저장/아카이브/테스트 문맥과 DB 提交=커밋을 용어집에 확정 등록했으며 기존 다른 문맥은 유지한다.
- **DB 커밋 근거**: evidence/store.go의 WithInstalledSnapshots가 본문 설치 후 복원 callback을 호출하며, server/task_archives.go:287–289/332–340에서 DB 복원으로 연결한다. db/task_archives_restore.go:301의 tx.Commit과 db/finding_traffic.go:27–39의 WithEvidenceTx commit 종료를 확인해 64행을 메타데이터 커밋으로 정리했다. 62행은 트랜잭션을 커밋한 뒤에만 planner에게 알립니다로 적용했다. 실제 DB 실행이나 커밋 성공을 관측한 것이 아니다.
- **정적 보존·상태**: 문서 99줄·LF·들여쓰기·코드 블록·인라인 코드·URL/경로·숫자·키·enum과 지정 표현 밖 바이트를 역치환 비교로 확인했다. 순서·조건·매시간 실행·최소 24시간 지연·공유 증거 보존 의미는 유지한다. 앞서 server/llmpool_test.go의 순환 전환 순서 인용 수정과 나머지 기존 변경은 그대로다. H01–H07의 과거 보류 근거는 이력으로 보존하고 로컬 감사에 승인 후 상태를 별도로 기록했다. 현재 문서의 번역 보류 0; 코드/예시/placeholder 중국어는 보존한다. M100·운영 퍼지 4곳·M107·기존 원문 정책/기능/호환성 후속은 유지한다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 누적 `/tmp/artex-evidence-doc-omissions.diff`에는 문서·테스트·용어집·PROGRESS만 포함하고, 감사 MD/TSV는 Git 제외 translation/audit/에만 보관한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다. fetch·병합·stage·commit·push 없음.


### M100 · M044 운영 후속 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `844731aabdb565ff7641db6e5027ea4d72e14fc3`. 대상 기존 변경·staged·MERGE_HEAD 부재를 확인하고 traffic/traffic.go:1541 오류 한 곳과 db/asset_intercept_match.go:24·26·28·35의 네 표현만 수정했다.
- **승인 결과**: M100 상태 파서→상태 확인 함수는 archiveCommitted의 실제 상태 조회 역할에 맞춘 번역 보완이다. 원문 状态解析器 명칭의 모호성은 감사 근거로 보존하며 원문 기능·콜백 정책을 변경하지 않는다. M044 운영 후속은 도메인/IP/URL 라벨 및 Reason 주석 인용을 부분 일치로 통일했다. 기존 확정 용어를 재사용하며 용어집 수정 없음.
- **정적 보존·연결**: 소스 두 파일의 줄 수·LF·들여쓰기 및 다섯 지정 행 밖 바이트가 수정 전과 동일함을 역치환 비교로 확인했다. %d/journal.ArchiveID·nil 검사·콜백·복원/정리 분기·fuzzy_* enum·pattern·Contains·출력 구조·기존 매칭 표현을 보존했다. 직접 테스트 기대값 의존은 재검색에서 발견하지 못해 테스트를 수정하지 않았다.
- **기록·후속**: 로컬 감사 MD/TSV의 원문·수정 전 문구·발견 근거·과거 추가 확인 상태는 이력으로 보존하고 승인 후 상태를 별도로 기록했다. 커버리지는 소스 두 파일의 승인 범위와 해시만 갱신하며 기존 A/S·M 이력 및 원문 정책/기능/호환성·실제 DB 미확인 사항은 유지한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존한다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. `/tmp/artex-m100-m044-operational-approved.diff`에는 소스 두 파일과 이번 PROGRESS만 포함한다. 감사 MD/TSV는 Git 제외 translation/audit/에 로컬로만 보관한다. fetch·병합·stage·commit·push 없음.


### M110–M111 · 논리 삭제 표시·README 제목 인용 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `5152f8c36246290b3150c7f5254483a1632f64f7`. 직전 M100·M044 운영 후속 커밋 반영 및 대상 기존 변경·staged·MERGE_HEAD 부재를 확인했다. 최신 기준 문서와 로컬 최종 검수 근거를 적용하고 승인 10곳만 수정했다.
- **승인 결과**: sessions-tab.tsx:184·314·428·429·539·619·1159·2127·2143의 주석·title·aria-label·표시 설명 9곳은 소프트 삭제→논리 삭제로 정리했다. soft/hard enum·deleted·delete_reason·조건·요청·상태 처리·JSX 구조를 보존했다. dev.sh:5의 옛 중국어 제목 인용만 현재 README:146의 방법 4: 소스에서 단일 바이너리 컴파일로 변경했으며 나머지 주석·명령·조건은 그대로다. 기존 확정 용어를 재사용하며 용어집 수정 없음.
- **정적 검토·보존**: 소스 두 파일의 줄 수·LF·끝 개행·들여쓰기와 지정 표현 밖 바이트를 수정 전/후 및 승인 행 역치환으로 확인했다. 직접 비교·파싱·테스트 기대값 의존은 재검색에서 발견하지 못해 테스트를 수정하지 않았다. status.ts의 stopped/软删除 설명, 추가 pause 표기 18행·적중 문맥 확인 후보 8행과 원문 정책/기능/호환성·실제 DB 미확인 후속은 유지한다.
- **기록·산출물**: 로컬 감사의 원문·수정 전 번역·발견 근거는 이력으로 보존하고 승인 후 상태를 별도로 추가했다. 커버리지는 지정 소스 두 파일의 승인 범위와 해시만 갱신하며 기존 A/S·M 집계를 혼합하거나 전체 저장소 의미 전수 검증 완료로 확대하지 않는다. 감사 MD/TSV는 Git 제외 translation/audit/에만 보관하고 커밋용 `/tmp/artex-m110-m111-approved.diff`에는 소스 두 파일과 이번 PROGRESS만 포함한다. 기존 로컬 변경·미추적 파일·Git index·제외 설정을 보존했다.
- **실행**: 실행 검증은 **사용자 요청으로 미실행**이다. fetch·병합·stage·commit·push 없음.


### M112–M115 · 일시 중지·매칭 문맥 승인 보완 [지정 수정·정적 대응 완료 / 실행 미검증]
- **기준·범위**: ko-translation, HEAD `3d7b3896836a993582be78e2839d663457d75d16`. 직전 M110·M111 커밋 반영과 대상 기존 변경·staged·MERGE_HEAD 부재를 확인하고 로컬 pause-match-followups-semantic-findings.tsv에 지정된 소스 8개 고유 26행만 보완했다.
- **승인 결과**: M112 暂停 문맥 18행은 일시 중지와 원문 활용형·부정·필수 강도를 유지했고, M113 채널 필터 3행은 매칭/매칭되지 않은으로 정리했다. M114 자산 조회 2행은 발견에 매칭된다·매칭된 자산 행, M115 트리 검색 3행은 매칭된 노드·매칭되지 않아도·검색에 매칭된 항목을 적용했다. asset-tree:203의 일치한 분기는 그대로다. 이전 같은 행의 힌트 추가 등 승인 표현을 보존했다.
- **용어·DB 영향**: 중복·충돌을 확인하고 기존 命中=매칭 항목에 취약점 자산 필터/조회 조건의 발견·자산 행 및 자산 트리 키워드 검색의 노드·항목 문맥만 추가했다. 일시 중지·채널 필터 항목은 중복 등록하지 않았고 캐시 hit·다른 문맥은 유지한다. db/db.go의 auto 설명은 기존 seedBuiltins ON CONFLICT(key) DO UPDATE의 description 갱신으로 기존 DB 메타데이터에 반영될 수 있다. DB 실행·마이그레이션·저장 로직·사용자 프롬프트 변경은 없다.
- **정적 검토·보존**: 소스 8개 지정 26행과 용어집 한 행의 줄 수·LF·들여쓰기·지정 밖 바이트를 수정 전/후 및 승인 행 역치환으로 확인했다. 도구명·enum·상태·ID·SQL·요청·JSX·변수·조건·포맷/인수는 그대로다. 직접 문자열 비교·파싱·테스트 기대값 의존은 재검색에서 발견하지 못해 테스트를 수정하지 않았다.
- **기록·후속**: 로컬 감사 원문·수정 전 근거·과거 문맥 확인 상태를 이력으로 보존하고 승인 후 상태를 별도로 추가했다. 커버리지는 소스 8개의 승인 scope·해시만 갱신한다. 기존 A/S 집계·M 승인 이력·원문 정책/기능/호환성·실제 DB 미확인 후속과 기존 로컬 변경·미추적 파일·Git index·제외 설정은 보존한다.
- **실행·산출물**: 실행 검증은 **사용자 요청으로 미실행**이다. 감사 MD/TSV는 Git 제외 translation/audit/에만 보관하고 `/tmp/artex-m112-m115-approved.diff`에는 소스 8개·용어집·이번 PROGRESS만 포함한다. fetch·병합·stage·commit·push 없음.
