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
