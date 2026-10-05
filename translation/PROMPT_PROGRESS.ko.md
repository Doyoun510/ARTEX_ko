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
- **관련 테스트 확인 보류**: 기존 확인 소스의 `agent/prompt_test.go`는 `plannerSystem` 결과에 `中间产物输出规约`가 포함되는지 검사한다. 최신 브랜치에서 해당 생성 문자열·번역 결과·기대값을 정적으로 대조하고, 필요한 기대값 수정 여부와 편집 담당을 기록한다. 이 기록은 테스트 실패를 확정한 것이 아니다. 테스트 실행은 별도 허용 전에는 하지 않는다.

## U3~U15  [todo]
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
