# ARTEX 통합 소스 번역 감사 갱신(U12 포함)

## 기준과 보존
- 브랜치 `work/ko-translation`; 기준 HEAD `474300b6d355fc1c8cf7051d644c52a7abae1fd4`.
- U12 병합 커밋 474300b의 부모는 f412bf1 및 5cc5ec5. MERGE_HEAD 없음, staged 변경 없음. 병합 완료 확인 후 정적 감사를 진행했다.
```text
 M translation/TRANSLATION_PROMPT.ko.md
?? translation/ARTEX_UI_ABC_SCOPE.ko.md
```
- TRANSLATION_PROMPT 미커밋 변경·미추적 ARTEX_UI_ABC_SCOPE를 포함한 모든 시작 파일 해시, index 해시 및 staged/unstaged diff를 전후 비교했다.
- 내용 변경 0개, index 동일 True, 전체 기존 상태 보존 True. 기존 감사 .md/.tsv 동일 True.
- 최신 용어집/TRANSLATION_PROMPT/GUIDE는 이전 감사와 파일 해시가 같고 규칙을 재적용했다. PROGRESS의 새 U12 24행과 서버 변경을 읽어 반영했다. 프로젝트 AGENTS.md는 기존 확인과 이번 검색에서 발견하지 못했다.
- fetch/병합/stage/commit/push/소스·진행 기록·용어집 수정 없음. 실행 검증: **사용자 요청으로 미실행**.

## 요약
- 기존 A001–A074 재대조: {'유지': 62, '판단 정정': 3, '추가 확인': 9}.
- 새 A075–A090: 16그룹. schema.sql 주석 누락165행, SQL 생성 표시문구3곳, ICP 주석, 표시 구분자와 U12 용어 등. 이 숫자는 문자열/행/실패 테스트 개수가 아니라 항목 그룹 수다.
- 기존 직접 기대값 불일치 A001–A010은 모두 유지. U12 DB 내부에서 **추가 직접 기대값 불일치를 확정한 항목은 없다**. 이는 모든 테스트 일치/통과를 뜻하지 않는다.
- A062는 양쪽 writer가 내장 기본값으로 일치하여 **판단 정정**했다. 과거의 DNT 분류는 고정 SQL 값 자체를 계약으로 과도하게 취급했다. 소스 라벨 문제는 이미 해소됐고 옛 저장 note는 별도 이력이다.
- A020의 미번역 개발자 안내는 유지하지만 이미 한국어인 직접 기대값 위치를 누락 목록에서 제외했다. A059 新对话는 현재 비교 쌍이므로 보존하되 영구 외부 계약이라고 단정하지 않는다.

## 범위와 깊이
- 전체 추적 파일603개 중 텍스트585개 검색, PNG/ICO 이미지18개 내용 제외. 중국어/전각 후보137개 파일, 검색 행2623개. 기존 후보 목록을 재사용하고 현 HEAD 전체 검색을 갱신했다.
- U12 db/ 추적 파일96개를 모두 정적 텍스트 검색·인벤토리/관련 리터럴 분석에 포함했다. 후보23개 파일282행: schema.sql 168행(주석165+SQL 표시값3), 그 밖은 입력/fixture/계약/원문 경로 및 두 주석의 备案 등.
- U12 후보의 생성부·소비부/원문을 문맥 대조했다. CJK 없는73개 파일에도 재시도/알림/자산/seed 용어와 기대값·모델 안내 분석을 추가했다. 그러나 96개 모든 한국어 문장·모든 SQL/분기 의미를 전수 검증한 것은 아니다.
- 이전 비U12 문자 검색 및 상세 의미 감사 한계 유지. CJK 0건은 번역 완료 근거가 아니다. 이 보고서는 범위 목록과 직접 발견의 정적 감사이며 전체 한국어의 원문 의미 완결 감사가 아니다.

## 기존 74개 항목 재대조 장부

| ID | 현재 처리 | 분류 | 위치 | 재대조 결과 |
|---|---|---|---|---|
| A001 | 유지 | 계약·테스트 기대값 불일치 | `agent/finding_recorder_test.go:102` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A002 | 유지 | 계약·테스트 기대값 불일치 | `server/finding_traffic_test.go:355` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A003 | 유지 | 계약·테스트 기대값 불일치 | `server/finding_workflow_test.go:134` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A004 | 유지 | 계약·테스트 기대값 불일치 | `server/finding_workflow_test.go:249` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A005 | 유지 | 계약·테스트 기대값 불일치 | `server/intercept_detail_test.go:96` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A006 | 유지 | 계약·테스트 기대값 불일치 | `server/task_categories_test.go:48` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A007 | 유지 | 계약·테스트 기대값 불일치 | `server/chat_mentions_test.go:180` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A008 | 유지 | 계약·테스트 기대값 불일치 | `server/chat_mentions_test.go:307` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A009 | 유지 | 계약·테스트 기대값 불일치 | `server/skill_upload_test.go:192` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A010 | 유지 | 계약·테스트 기대값 불일치 | `server/skill_upload_test.go:213` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A011 | 유지 | 실제 번역 누락 | `web/next.config.mjs:3,4,6,11,14,15,16,23,30,34` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A012 | 유지 | 실제 번역 누락 | `agent/capture_usage_test.go:36,62` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A013 | 유지 | 실제 번역 누락 | `agent/coldgraph_test.go:17,40` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A014 | 유지 | 실제 번역 누락 | `agent/insert_assets_test.go:45,114,162,245,306,364,401,402` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A015 | 유지 | 실제 번역 누락 | `server/llmretry_test.go:11,12` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A016 | 유지 | 실제 번역 누락 | `server/task_archive_package_test.go:107` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A017 | 유지 | 실제 번역 누락 | `server/trigger_merge_test.go:12,95` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A018 | 유지 | 실제 번역 누락 | `server/engine_emptyturn_test.go:11,40,57,79,93,96,101,105,109,113,117,125,128,136,141,143` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A019 | 유지 | 실제 번역 누락 | `server/update_test.go:13,14,15,38,52,57,71,77,89,91,93,96,99,101,105,108,121,122,126,131,134,135,137,140,143,144,150,153` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A020 | 판단 정정 | 실제 번역 누락 | `server/notify_api_test.go:188,196,201,202,216,220,246,249,431,438,720,780` | 중국어 주석·자체 실패 안내 잔여는 유지한다. 기존 위치 목록의 719 상세 보기/779 나머지 등 한국어 기대값·문맥 행은 누락 위치에서 제거했다. 대응된 8개 기대값을 새 누락/불일치로 집계하지 않는다. |
| A021 | 유지 | 실제 번역 누락 | `server/skill_upload_test.go:57,76,84,85,131,154` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A022 | 유지 | 실제 번역 누락 | `sidequestion/README.md:1,3,5,7,9` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A023 | 유지 | 실제 번역 누락 | `sidequestion/CONTEXT_BUDGET.md:1,3,5,15,17` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A024 | 유지 | 실제 번역 누락 | `sidequestion/VALIDATION.md:1,3,5,7,9` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A025 | 유지 | 실제 번역 누락 | `web/src/app/(main)/chat/page.tsx:900,932,966,994` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A026 | 유지 | 실제 번역 누락 | `web/src/lib/mock/handler.ts:1231,1250,1310,1336,1804` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A027 | 유지 | 용어 불일치 | `server/llmpool.go:13,14,86,139` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A028 | 유지 | 용어 불일치 | `server/server_mgmt.go:1941,1943` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A029 | 유지 | 용어 불일치 | `server/llmpool.go:17,34,51` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A030 | 유지 | 용어 불일치 | `server/llmretry.go:13,16,52` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A031 | 유지 | 용어 불일치 | `server/server.go:155` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A032 | 유지 | 용어 불일치 | `server/task_llm.go:94,249,299,397` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A033 | 유지 | 용어 불일치 | `server/server_mgmt.go:1848,1849` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A034 | 유지 | 용어 불일치 | `server/dto.go:616` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A035 | 유지 | 용어 불일치 | `server/llmretry.go:11` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A036 | 유지 | 용어 불일치 | `agent/planner.go:304,311,333,338` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A037 | 유지 | 용어 불일치 | `agent/compaction.go:538` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A038 | 유지 | 용어 불일치 | `web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx:243,274,417,418` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A039 | 유지 | 용어 불일치 | `agent/retester.go:6` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A040 | 유지 | 용어 불일치 | `agent/finding_workflow.go:77` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A041 | 유지 | 용어 불일치 | `sidequestion/service.go:61` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A042 | 유지 | 용어 불일치 | `agent/promptcatalog.go:95` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A043 | 유지 | 용어 불일치 | `README.ko.md:307` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A044 | 유지 | 띄어쓰기·문장부호 | `web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx:2121,2122,2123` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A045 | 유지 | 띄어쓰기·문장부호 | `README.ko.md:61,126,182,186,297,300,305,307,354,359,371,425,426,427,428,442,451,452,479` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A046 | 유지 | 후속 검토 보류 | `server/intercept_live_test.go:126,132` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A047 | 유지 | 후속 검토 보류 | `intercept/prompt.go:20,21,22,23,24,25,26,27,28` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A048 | 유지 | 후속 검토 보류 | `skills/api-recon/SKILL.md:21,22,41,151` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A049 | 유지 | 후속 검토 보류 | `skills/api-recon/scripts/preload.js:192` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A050 | 유지 | 후속 검토 보류 | `web/src/lib/mock/data.ts:3067,3080,3093,3106,3136,3155,3156,3163,3168,3172,3218,3229,3239,3513,3515,3544,4165,4166,4196` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A051 | 추가 확인 | 용어·계약 판단 보류 | `web/src/components/mention-textarea.tsx:148,239,268,303` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A052 | 추가 확인 | 용어·계약 판단 보류 | `docker-compose.bench.yml:1` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A053 | 추가 확인 | 용어·계약 판단 보류 | `server/intercept.go:50` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A054 | 유지 | 계약상 DNT | `intercept/intercept.go:480,481,482` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A055 | 유지 | 계약상 DNT | `intercept/intercept.go:612` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A056 | 유지 | 계약상 DNT | `server/task_archives.go:461` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A057 | 유지 | 계약상 DNT | `server/server_mgmt.go:1361` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A058 | 유지 | 계약상 DNT | `intercept/prompt.go:36,37,38,39,124,125,126,127,193,196,200` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A059 | 판단 정정 | 계약상 DNT | `server/conversations.go:112,443` | 新对话는 저장 기본값/자동 제목 if의 같은 파일 쌍이다. 현재 비교값 보존은 필요하지만, 영구 외부 계약이라고 확정하지 않는다. 쌍 변경 허용 조건·옛 DB 제목·mock 및 SDK 경로를 확인하면 바꿀 수 있는 후속 후보. 현재 표시 fallback 문제 A025는 별도 유지. |
| A060 | 유지 | 계약상 DNT | `server/orchestration.go:699,700,701,719` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A061 | 유지 | 계약상 DNT | `skills/api-recon/reference.md:52` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A062 | 판단 정정 | 이미 해소된 표시 라벨 | `server/finding_retests.go:230` | server/finding_retests.go:230과 db/config.go:751의 writer 모두 내장 기본값. db/config.go:758은 내장 기본값으로 복원. note는 agent-editor.tsx:423/456에서 표시하며, 해당 값의 비교/파싱 소비자는 검색되지 않았다. 이전의 “계약상 DNT” 분류는 과도했다. 현재 소스 라벨 대응은 이미 해소됐으며 옛 DB note의 소급 번역은 하지 않은 것으로 구분한다. |
| A063 | 유지 | 원문 예시·placeholder·검증 데이터 | `report/findings.go:125` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A064 | 유지 | 원문 예시·placeholder·검증 데이터 | `agent/terminalreason_test.go:104,109,163,166` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A065 | 유지 | 원문 예시·placeholder·검증 데이터 | `skills/scopesentry/SKILL.md:15,16,32,34,41,93,120,133,137,158,169,188,262,334` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A066 | 유지 | 원문 예시·placeholder·검증 데이터 | `intercept/prompt.go:115` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A067 | 유지 | 계약상 DNT | `agent/tools_insert.go:142,169,380,382,387,437` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |
| A068 | 추가 확인 | 용어·계약 판단 보류 | `notify/email_test.go:206,207,208,209,210,264,265,266,267,268` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A069 | 추가 확인 | 용어·계약 판단 보류 | `notify/filter_test.go:13,14,15,16,17,77,78,79,80,81,82,83,84,105,106,107,108,109,110,111,113,118` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A070 | 추가 확인 | 용어·계약 판단 보류 | `notify/http_test.go:36,37,38,39,40,41,42,43,44,45` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A071 | 추가 확인 | 용어·계약 판단 보류 | `notify/mask_test.go:132,142,149,156,164,171,215,221,227,233,239,245` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A072 | 추가 확인 | 용어·계약 판단 보류 | `notify/pack_test.go:122,133,142,152` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A073 | 추가 확인 | 용어·계약 판단 보류 | `notify/redact_test.go:39,45,51,57,63,98,99,100,101,102` | 현재 후보/계약 관계 유지. 외부 test-name 선택, 전문명 표기 또는 표시·계약 분리의 추가 확인을 번역 누락 확정과 구분한다. |
| A074 | 유지 | 진행 기록·범위 확인 | `translation/PROMPT_GUIDE.ko.md:69` | 현재 인용 위치 문구가 기존 감사와 동일. 이 항목 파일은 U12 통합 변경 대상 밖이며 기존 연결 소비부·테스트의 해당 문구도 유지. |

## 새 항목 상세
- 전체 현재 문구/이전 문구/원문·기준·소비부·수정 후보는 TSV에 기록했다. SQL 문구가 포함됐다는 이유로 SQL 전체를 번역 대상으로 삼지 않는다.

### A075 · 실제 번역 누락
- 위치/절: `db/schema.sql:1,2,5,11,12,13,14,24,151,196,285,343,374,375,379,380,381,384,386,387,389,390,391,392,393,395,396,398,399,413,416,418,419,424,425,428,435,438,439,440,441,442,443,450,456,457,458,459,460,461,476,477,478,481,482,483,489,541,543,545,546,547,554,619,621,626,732,733,734,735,736,737,750,755,764,791,795,797,814,831,842,851,857,858,859,875,897,898,899,910,911,917,918,933,952,991,1015,1031,1063,1070,1079,1086,1087,1094,1095,1096,1098,1108,1111,1131,1177,1232,1233,1234,1235,1236,1237,1238,1239,1260,1261,1262,1263,1264,1265,1282,1289,1291,1292,1293,1295,1296,1297,1300,1301,1302,1306,1309,1310,1318,1320,1322,1323,1324,1326,1327,1336,1337,1338,1339,1340,1347,1348,1355,1356,1357,1358,1363,1364,1365,1366,1369,1370,1371,1374` / schema/DDL developer comments; 담당 U12.
- 현재: 설명 주석165행 전체 현재 문구는 TSV와 아래 db 잔여 장부에 기록. 대표: schema.sql:1 “单一数据源”, :374–399 모델/재시도 설명, :1086–1098 심각도/상태 설명, :1291–1374 알림 설계.
- 원문·기준: 문서·SQL 주석도 U12 전체 범위; key/enum/SQL는 보존하고 산문만 번역
- 근거: 165개 후보 주석 행에 자산/의도/모델/재시도/알림 등의 중국어 설명이 남음. 역사적 문서 경로·[模型]·코드 값이 섞인 행은 해당 토큰을 유지해야 함. schema.sql은 f412bf1 이후 변경 없음.
- 연결: db/db.go:20–21 go:embed; :43 ExecContext(schemaSQL); schema migration tests
- 수정 후보: 중국어 설명 주석만 번역; SQL/키/enum/숫자/구문·원문 경로와 계약 토큰 유지
- 판정: 누락 확인; 일부 전문 용어는 후속 용어 판단.

### A076 · 실제 번역 누락
- 위치/절: `db/schema.sql:670` / sync_task_asset_links provenance summary; 담당 U12↔UI.
- 현재: 670: SELECT task.id, NEW.id, 'system', '任务执行期间自动关联'
- 원문·기준: source_summary는 표시 산문; source=system은 enum
- 근거: SQL이 생성하는 任务执行期间自动关联은 task_asset_links.source_summary에 저장됨. 정확한 문자열 비교/파싱 검색 결과 이 생산부만 발견. source/summary를 구분.
- 연결: db/task_assets.go:345–361 → Asset.TaskSourceSummary(db/assets.go:59) → assets-tab.tsx:136 SourceCell
- 수정 후보: 작업 실행 중 자동 연관 등으로 표시 문자열만 번역. system/source와 ON CONFLICT는 유지. 기존 행 소급 수정은 별도 범위.
- 판정: 생성/표시 연결 확인 / 실제 DB 미확인.

### A077 · 실제 번역 누락
- 위치/절: `db/schema.sql:688` / legacy task_asset_links migration summary; 담당 U12↔UI.
- 현재: 688: SELECT task.id, asset.id, 'legacy', '由历史任务资产关联迁移'
- 원문·기준: legacy enum과 description을 구분
- 근거: 由历史任务资产关联迁移는 마이그레이션 비교값이 아니라 INSERT의 source_summary. ON CONFLICT DO NOTHING은 기존 이력을 덮어쓰지 않음. 정확한 문구의 비교 소비자 없음.
- 연결: db/task_assets.go:345–361; assets-tab.tsx:136(한국어 fallback)
- 수정 후보: 과거 작업 자산 연관에서 마이그레이션됨 등으로 새 삽입 설명만 번역. source=legacy·SQL·기존 행 유지.
- 판정: 정적 표시값 확인 / 기존 데이터 미확인.

### A078 · 실제 번역 누락
- 위치/절: `db/schema.sql:1131,1134` / stop_deleted_conversation_retest; 담당 U12↔U11↔UI.
- 현재: 1131: -- 删除会话保留复测记录，同时解除尚未结束的复测占用。 || 1134: UPDATE finding_retests SET status='stopped', error='复测会话已删除', finished_at=now()
- 원문·기준: stopped는 상태값; error는 사용자 오류 안내
- 근거: 复测会话已删除는 삭제 trigger가 저장하는 error 산문. 문자열 비교/파싱 소비자는 검색되지 않음.
- 연결: db/finding_retests.go:29 Error, ListFindingRetests → server/finding_retests.go:53–58 → finding-retest-panel.tsx:151 item.error; db/finding_retests_test.go:125–151은 상태/참조 검사
- 수정 후보: 재검증 대화가 삭제되었습니다 등으로 error 값만 대응; 상태/트리거 조건 유지. 댓글 1131은 A075와 같은 누락.
- 판정: 소스/소비 경로 확인; 테스트 직접 문자열 변경 필요 발견 안 됨.

### A079 · 실제 번역 누락
- 위치/절: `db/company_scope.go:155` / ICP explanatory comment; 담당 U12.
- 현재: 155: // 备案 번호 자체에는 점이 없다(예: 京ICP备12345678号-1). 점이 있는 텍스트는 대개 도메인이나 버전 번호가 섞인 것이라,
- 원문·기준: 备案 표시=ICP 등록(备案), parser=备案 DNT
- 근거: 주석의 备案 번호는 산문이고 줄160의 strings.Contains(raw,备案)는 실제 입력 계약. 주석만 확정 표기로 맞출 수 있음.
- 연결: db/company_scope_test.go:63 동일 설명; :59–67/212–216 중국어 입력·normalization 쌍
- 수정 후보: 주석 두 곳을 ICP 등록(备案) 번호로 대응; parser/입력·예시/기대값은 그대로
- 판정: 산문/계약 구분 확인.

### A080 · 띄어쓰기·문장부호
- 위치/절: `db/companies.go:636` / malformed IP warning sample list; 담당 U12.
- 현재: 636: total, strings.Join(samples, "、"),
- 원문·기준: GLOSSARY:710에 이 strings.Join의 、→,  명시 허용
- 근거: 경고 표시용 목록 구분자가 아직 、. 값 비교가 아니라 경고 문자열 조합이며 logAttributionWarning으로도 노출.
- 연결: db/companies.go:579–584 및 recompute 호출 경로; malformed IP tests
- 수정 후보: strings.Join(samples, ", ")로 제한 대응; 다른 구분자/샘플/포맷/쿼리 유지
- 판정: 명시 허용 위치 확인.

### A081 · 용어 불일치
- 위치/절: `db/config.go:71` / RetryOverride documentation; 담당 U12.
- 현재: 71: // that are per-endpoint: 연결 수립(connect) / 빈 응답(empty) / 같은 provider 안전 윈도우
- 원문·기준: 安全窗口=안전 구간 (호출자 출력 전달 전 재시도)
- 근거: 같은 provider 안전 윈도우 표기. 같은 LLM 재시도 문맥.
- 연결: server/task_llm.go; 기존 A032/A034/A035; RetryOverride
- 수정 후보: 중국어에서 번역된 한국어 인용만 안전 구간으로 맞춤; 영어/구조 유지
- 판정: 용어집 문맥 일치 확인.

### A082 · 용어 불일치
- 위치/절: `db/llmretry.go:87` / stream retry defaults comment; 담당 U12.
- 현재: 87: // Stream: 같은 provider 안전 윈도우 재시도(출력 전달 전 끊긴 스트림 재생). 기본 2회·0.5s부터 지수(상한 4s).
- 원문·기준: 安全窗口=안전 구간
- 근거: 출력 전달 전 끊긴 스트림 재생을 안전 윈도우로 표기
- 연결: server/task_llm.go sameProvider retry; db RetryRule
- 수정 후보: 안전 구간으로 맞춤; 2회·0.5s·4s 유지
- 판정: 정적 문맥 확인.

### A083 · 용어 불일치
- 위치/절: `db/db.go:195,196` / Auto builtin description/comment; 담당 U12.
- 현재: 195: // Auto: 내장 '플랫폼 조작' agent. 침투 오케스트레이션 루프에 참여하지 않고, 대화 페이지로 구동되며, 도구로 플랫폼을 조작한다. || 196: {"auto", "Auto", "assistant", "플랫폼 조작 도우미: 도구로 작업(생성/조회/일시중지/힌트 제공)과 자산을 관리하고, skill·커스텀 도구·MCP를 생성/수정할 수 있다.", nil, false, nil},
- 원문·기준: 操作=동작
- 근거: 플랫폼 조작/플랫폼을 조작한다가 확정 표기와 다름. 원문 f412bf1:db/db.go 동일 위치 平台操作.
- 연결: seedBuiltins name/description upsert; agent UI
- 수정 후보: 플랫폼 동작 도우미/플랫폼 동작을 수행한다 등 최소 문장 대응. auto 키·기능 설명 유지.
- 판정: 원문/용어집 확인.

### A084 · 용어 불일치
- 위치/절: `db/db.go:198` / pentest builtin description; 담당 U12.
- 현재: 198: {"pentest", "침투 테스트", "assistant", "독립 침투 agent: 혼자 정찰→공격면 탐색→심화 익스플로잇→검증→마무리까지 전체 체인을 수행하고, 스스로 계획·실행·대립적으로 검증한다.", nil, true, intp(0)},
- 원문·기준: 攻击面=공격 표면
- 근거: 원문 找攻击面이 현재 공격면 탐색으로 번역됨
- 연결: seedBuiltins; agent UI
- 수정 후보: 공격 표면 탐색; 기존 agent key/명령·기능·순서 유지
- 판정: 동일 침투 테스트 문맥 확인.

### A085 · 용어 불일치
- 위치/절: `db/nkey.go:61` / NormalizeParamName explanatory comment; 담당 U12.
- 현재: 61: // '파라미터명으로 같은 회사 인터페이스 조회'가 재현 가능하도록 보장한다.
- 원문·기준: 接口=엔드포인트 (자산 endpoint 문맥)
- 근거: 같은 회사 인터페이스 조회는 endpoint.params/파라미터명 자산 조회 설명. 소프트웨어 interface 타입 문맥 아님.
- 연결: EndpointKey/ParameterKey/NormalizeParamName; endpoint asset query
- 수정 후보: 같은 회사 엔드포인트 조회로 맞춤; normalize 코드 변경 없음
- 판정: 자산 문맥 확인.

### A086 · 용어 불일치
- 위치/절: `db/notification.go:16,30,34` / notification delivery comments; 담당 U12↔U10.
- 현재: 16: // 이 파일은 IM 푸시의 채널 설정·이벤트 레이어다. 전달 작업의 획득과 상태 전이는 || 30: // 전달 상태. || 34: NotifyStateSent    = "sent"    // 전달됨
- 원문·기준: 投递=전송 (알림 문맥)
- 근거: 전달 작업/전달 상태/전달됨 표기는 알림 delivery 문맥의 확정 전송과 다름. 미전달 필드/false 전달처럼 generic data 전달은 제외.
- 연결: notify/ 전송 용어; notification_delivery.go; notification_test.go의 설명
- 수정 후보: delivery 주석을 전송으로 대응; 상태 enum pending/sent 등 불변. 모든 전달이라는 단어를 일괄 치환하지 않음.
- 판정: 문맥 확인.

### A087 · 용어 불일치
- 위치/절: `db/notification_delivery.go:12,15,18,23,24,38` / delivery lease/state documentation; 담당 U12.
- 현재: 12: // 이 파일은 전달 작업의 획득과 상태 전이다. || 15: // 리스 만료 시간으로 삼은 뒤, 트랜잭션을 커밋하고 나서 네트워크 전달을 한다. 이렇게 하면 전달 중 DB 락을 쥐지 않는다 —— || 18: // 대가는 프로세스가 전달 도중 크래시하면 행이 sending에 멈추는 것이다. 이건 **자가 치유 가능**하다: 리스 만료 후 || 23: // MaxNotifyAttempts는 한 전달의 최대 시도 횟수다(첫 시도 포함). || 24: // 전달 엔진이 아니라 여기에 정의한다: 상태 기계 자체의 정책이고, 엔진은 실행자일 뿐이다. || 38: // NotificationDelivery는 렌더링에 필요한 채널 설정과 이벤트 스냅샷을 포함한 전달 작업 하나다.
- 원문·기준: 投递=전송
- 근거: 알림 전달 작업/전달 엔진/네트워크 전달 등 동 문맥. 다음 파일 테스트의 전달 존재 여부 안내도 같은 용어 후보.
- 연결: db/notification_test.go:49,188,193,217; notify dispatcher
- 수정 후보: 알림 전송 문맥만 전송으로 대응. 네트워크 동작·lease·attempt 예산 설명 불변.
- 판정: 동일 알림 전송 문맥 확인.

### A088 · 용어 불일치
- 위치/절: `db/notification.go:131,135` / explicit zero rate comments; 담당 U12↔U10.
- 현재: 131: // 여기서는 일부러 0을 가공하지 **않는다**: 0은 유효한 설정이며 '레이트 리밋 없음'을 의미한다. || 135: // takeTokens 모두 0을 레이트 리밋 없음으로 해석하는데, 여기서만 몰래 20(DingTalk/WeCom/Telegram)
- 원문·기준: 限流=전송 속도 제한(알림), 원문 未限流
- 근거: 레이트 리밋 없음 표현은 알림 전송 rate_per_min 설명. 원문 f412bf1에서 不限流.
- 연결: db/notification_test.go:699,702 동일 설명; notify UI ratePerMin
- 수정 후보: 전송 속도 제한 없음으로 맞춤; 0·20·100과 원문 환산/과거 동작 기록 유지
- 판정: 용어 문맥 확인.

### A089 · 진행 기록·범위 확인
- 위치/절: `translation/PROMPT_PROGRESS.ko.md:186,202,208` / U12 completion claim; 담당 U12 기록.
- 현재: 186: ## U12 · 데이터/DB (트랙 B)  [완료 — 번역·정적·빌드/vet PASS] || 202: - **대형 생산**: notification·notification_delivery·finding_assets·db.go·exploration·asset_intercept_match·findings·tasks·llmretry·finding_retests·finding_traffic·company_scope·companies·chat_mentions·assets. || 208: - **잔여 중국어(의도 DNT만)**: 설계문서 경로(`任务级超时与收尾设计.md`·`交互式shell设计.md`·`LLM重试设计.md`·`资产模型与自动关联设计.md`)·`备案`·`[模型]` 계약·테스트 픽스처.
- 원문·기준: 완료 보고와 실제 산문 잔여를 구분
- 근거: U12 기록은 잔여 의도 DNT만이라고 보고하지만 schema.sql 설명165행 및 SQL 표시값3곳, Go 주석 잔여가 확인됨. 기록상의 빌드/vet PASS는 번역 완결 근거 아님.
- 연결: A075–A079; 실제 통합 변경 db Go46개 / schema.sql unchanged
- 수정 후보: 담당자 작업 이력·기존 PASS 기록은 보존하고 실제 잔여·정적 검토 범위를 별도 후속 기록
- 판정: 기록과 소스 차이 확인.

### A090 · 용어·계약 판단 보류
- 위치/절: `db/db.go:192` / ProxyAddr variable description; 담당 U12↔U3.
- 현재: 192: {"ProxyAddr", "기록 프록시 주소(if 이중 문구 구동)", "127.0.0.1:8080", "runtime"},
- 원문·기준: 원문 f412bf1: 记录代理地址(驱动 if 双文案)
- 근거: if 이중 문구 구동은 중국어 전문 표현 双文案을 풀어 옮긴 것으로 보이나 용어집에 직접 정의 없음. 두 문구/두 모드 무엇을 지칭하는지 확정 못함.
- 연결: workerTrafficBlock/egressOnly system rules; Agent variable UI
- 수정 후보: 원문 문장 전체 ProxyAddr: 记录代理地址(驱动 if 双文案). 문맥 확인 후 “기록 프록시 주소(if 두 문구 방식 구동)” 등 후보를 검토; 임의 코드 기능 추가 금지.
- 판정: 추가 확인 필요.

## U12 계약·테스트·형식 대조
- [模型]: db/intercept.go:231 LIKE 및 :251 HasPrefix, db/task_archives_restore.go:680 prefix와 기존 intercept 생성부/프런트 분류가 같은 중국어 접두사를 사용. 유지해야 하며 누락 아님.
- 内置默认/恢复为内置默认: 현재 두 소스 writer는 내장 기본값/내장 기본값으로 복원. agent-editor.tsx:423/456은 note를 표시한다. 해당 정확 문자열 비교/파싱을 검색하지 못했다. 기존 DB note의 실제 내용과 모델 동작은 미확인.
- ICP: company_scope.go:160 备案 판별, NormalizeICP와 company_scope_test.go:59–67/212–216의 중국어 입력/expected는 파서·정규화 검사로 보존. 주석:155/63의 설명만 A079 번역 대상.
- tasks.go:298의 작업 생성 시 회사 연결: + 이름은 tasks_test.go:610과 정적 일치. manualTaskScopeSummary는 task_assets_test.go:47이 상수 심볼 비교하므로 fixture를 임의 번역하지 않는다.
- task_archives_restore.go:219의 작업 회사 %d이(가) 삭제…는 task_archives_test.go:274의 같은 부분 문자열 포맷과 정적 일치. 형식 강도/동작 테스트는 미실행.
- db/schema.sql:1134 삭제 사유는 기존 finding_retests_test.go:125–151에서 상태/참조만 검사하므로 새 직접 기대값 수정이 필요하다고 단정하지 않는다. API→UI의 error 표시는 확인했고 새 번역 대상 A078로 구분.
- U12 변경 Go46개는 f412bf1 원문 대비 원본/현재 줄 수와 printf 포맷 토큰의 종류·순서·개수가 동일했다. 전체 인수 AST/실행 동등성을 증명한 것은 아니다.
- db/db.go 내장 규칙의 pattern/target/typ/action/priority 코드 필드91행은 f412bf1과 문자열 단위 동일했다. 실제 DB 안전장치 작동은 검증하지 않았다.
- FindingRetest.InitialMessage(db/finding_retests.go:126–131)의 원본 증거/제약 먼저 조회→표적 검증→결론 저장 순서와 %d는 원문과 대조했다. 모델의 실제 수행은 미확인이다.

## 다음 수정 묶음(이번에는 구현 안 함)

### 직접 연결된 테스트 기대값
A001, A002, A003, A004, A005, A006, A007, A008, A009, A010

### 번역 누락
A011, A012, A013, A014, A015, A016, A017, A018, A019, A020, A021, A022, A023, A024, A025, A026, A075, A076, A077, A078, A079

### 용어·표기·링크 정리
A027, A028, A029, A030, A031, A032, A033, A034, A035, A036, A037, A038, A039, A040, A041, A042, A043, A044, A045, A074, A080, A081, A082, A083, A084, A085, A086, A087, A088, A089

### 추가 확인이 필요한 계약
A046, A047, A050, A051, A052, A053, A068, A069, A070, A071, A072, A073, A090

### 번역과 별개인 원문 기능·정책 문제
A048, A049

- 직접 기대값 묶음: 검사 조건·입력·HTTP status/Deflate64/원 증거 검증을 유지하고 기대 문구만 대응. 수정 전 소스 최종 출력 재대조.
- 누락 묶음: U11 테스트 안내·U9 개발 문서/UI 잔여와 U12 schema 주석/자체 표시값을 구분해 처리. U12 SQL 구문·enum·migration 비교값·원문 fixture 보존.
- 용어·표기·링크 묶음: 안전 구간/회로 차단/순환 전환·계보/동작/엔드포인트/공격 표면·알림 전송 용어. 링크 새 불일치가 확인된 것은 없으며 이전 skills fragment 검사 한계를 유지한다. companies.go 구분자는 용어집 명시 허용 위치만 포함.
- 추가 확인 묶음: 테스트 이름 외부 선택 의존, mention UI displayLabel 분리, 新对话 비교 쌍/옛 저장값, 커스텀 공통 블록 중복·live 테스트 의미 검증, 双文案驱动. 사용자 정책 선택을 요구하거나 기능을 재설계하지 않았다.
- 원문 정책 문제: API 정찰 로그인/실제 세션 금지·조건부 세션 안내, depth/coverage의 forward/stub/L2 차이. 기존 U7/U15 후속 검토를 해결 완료로 표시하지 않는다.
- ScopeSentry 형식 정정, reference D→E 및 routes 설명 정정은 현재 소스에 반영된 상태로 유지. 새 수정 항목으로 다시 집계하지 않는다.

## 분류별 현재 장부 건수(보존/이력 포함)
| 분류 | 그룹 수 |
|---|---:|
| 계약·테스트 기대값 불일치 | 10 |
| 실제 번역 누락 | 21 |
| 용어 불일치 | 25 |
| 띄어쓰기·문장부호 | 3 |
| 후속 검토 보류 | 5 |
| 용어·계약 판단 보류 | 10 |
| 계약상 DNT | 9 |
| 이미 해소된 표시 라벨 | 1 |
| 원문 예시·placeholder·검증 데이터 | 4 |
| 진행 기록·범위 확인 | 2 |

## db/ 96개 파일 검사·잔여 장부
| 파일 | 행 수 | CJK/전각 후보 행 | 검사 상태 |
|---|---:|---:|---|
| `db/activity_page_test.go` | 178 | 2 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/asset_dsl.go` | 637 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/asset_intercept.go` | 80 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/asset_intercept_match.go` | 250 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/asset_intercept_match_test.go` | 129 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/assets.go` | 1536 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/assets_test.go` | 1126 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/chat_mentions.go` | 113 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/commands.go` | 353 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/companies.go` | 774 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/companies_test.go` | 460 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/company_lock_test.go` | 15 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/company_scope.go` | 272 | 2 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/company_scope_consistency_test.go` | 250 | 2 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/company_scope_test.go` | 405 | 11 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/config.go` | 1087 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/config_max_tokens_test.go` | 100 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/config_test.go` | 196 | 3 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/constants.go` | 40 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/constraints.go` | 85 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/conversation.go` | 327 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/conversation_pin_test.go` | 69 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/customtool_test.go` | 82 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/db.go` | 602 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/db_test.go` | 114 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/digest.go` | 253 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/exploration.go` | 2021 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/exploration_sources.go` | 397 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/exploration_sources_test.go` | 431 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/exploration_test.go` | 296 | 2 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/exploration_tokens_test.go` | 130 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/finding_assets.go` | 676 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/finding_assets_test.go` | 233 | 9 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/finding_retests.go` | 259 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/finding_retests_test.go` | 242 | 3 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/finding_traffic.go` | 476 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/finding_traffic_archive.go` | 85 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/finding_traffic_archive_test.go` | 87 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/findings.go` | 794 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/findings_test.go` | 519 | 14 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/intent_admission.go` | 32 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/intent_control_test.go` | 285 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/intent_delete_test.go` | 170 | 22 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/intercept.go` | 374 | 3 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/intercept_detail.go` | 114 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/intercept_detail_test.go` | 167 | 4 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/intercept_execution.go` | 120 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/intercept_execution_test.go` | 150 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/intercept_filter_test.go` | 141 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/intercept_seed_test.go` | 46 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/jsonb_clean_test.go` | 44 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/llm_records_migrate_test.go` | 62 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/llm_usage.go` | 270 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/llmhealth.go` | 57 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/llmretry.go` | 127 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/llmretry_test.go` | 145 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/logs.go` | 78 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/nkey.go` | 134 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/notification.go` | 551 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/notification_delivery.go` | 446 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/notification_test.go` | 874 | 24 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/schema.sql` | 1386 | 169 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/settings.go` | 43 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/side_questions.go` | 288 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/side_questions_test.go` | 446 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/skill_usage.go` | 217 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/skill_usage_test.go` | 116 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_archive_aggregate_stats.go` | 91 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_archives.go` | 787 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_archives_restore.go` | 793 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/task_archives_test.go` | 468 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_assets.go` | 425 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_assets_context.go` | 260 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_assets_test.go` | 188 | 4 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/task_categories.go` | 282 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_categories_test.go` | 183 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_context.go` | 578 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_context_lock_test.go` | 527 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_context_unit_test.go` | 27 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/task_delete_concurrency_test.go` | 210 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_intercept.go` | 122 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_list_performance_test.go` | 117 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_metadata_test.go` | 104 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_queue_test.go` | 66 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_scope.go` | 644 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_scope_test.go` | 23 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_source_limit_test.go` | 47 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_templates.go` | 276 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/task_templates_test.go` | 131 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/tasks.go` | 654 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/tasks_test.go` | 708 | 1 | 후보 문맥 분류 및 관련 생성/소비 대조; 원문 의미 전수 미검증 |
| `db/testmain_test.go` | 32 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/tool_usage.go` | 58 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/tool_usage_test.go` | 44 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/tools.go` | 232 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |
| `db/triggers.go` | 299 | 0 | 전체 텍스트 검색·리터럴/용어 검사; 원문 의미 전수 미검증 |

### 중국어 잔여의 개별 위치/분류
- 아래 원문 fixture 행에는 중국어 설명 이름과 테스트 실제 입력/expected가 혼합돼 있다. 입력 값 자체를 교체하지 않으며 fixture 내용을 반복한 failure 진단의 이름도 원문값으로 유지할 수 있다.

#### db/activity_page_test.go
- 15: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("test", "分页历史")`
- 136: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("test", "意图分页")`

#### db/companies.go
- 636: **표시 구분자 정리(A080), 명시 허용** — `total, strings.Join(samples, "、"),`

#### db/company_scope.go
- 155: **주석 설명 번역 누락(A079)** — `// 备案 번호 자체에는 점이 없다(예: 京ICP备12345678号-1). 점이 있는 텍스트는 대개 도메인이나 버전 번호가 섞인 것이라,`
- 160: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `(strings.Contains(lower, "icp") || strings.Contains(raw, "备案")) {`

#### db/company_scope_consistency_test.go
- 175: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `boundary := strings.Repeat("界", MaxCompanyScopeRawRunes)`
- 180: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := ValidateCompanyScopeInputBounds([]ScopeInput{{Kind: "keyword", Value: boundary + "界"}}); !errors.As(err, &validationErr) {`

#### db/company_scope_test.go
- 59: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "icp latin", input: "京 ICP备 123号", kind: "icp", normalized: "京icp备123号"},`
- 60: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "icp chinese", input: "沪网备案 9988", kind: "icp", normalized: "沪网备案9988"},`
- 63: **주석 설명 번역 누락(A079)** — `// 备案 번호는 점을 포함하지 않음: 도메인/버전 번호가 섞인 설명 문구는 키워드로 분류, 안 그러면`
- 65: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "icp with domain text", input: "备案 www.example.com", kind: "keyword", normalized: "备案 www.example.com"},`
- 66: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "icp with version text", input: "某公司 ICP v1.0", kind: "keyword", normalized: "某公司 icp v1.0"},`
- 67: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "icp fullwidth dot", input: "备案 例．com", kind: "keyword", normalized: "备案 例．com"},`
- 212: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `icp, err := ParseScopeInput(ScopeInput{Kind: "ICP", Value: " 京ICP 备 123号-1\t"})`
- 216: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if icp.Kind != "icp" || icp.Value != "京icp备123号-1" {`
- 268: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{Kind: "icp", Value: "京 ICP备 998877号"},`
- 276: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Domain: fmt.Sprintf("icp-scope-%d.example", suffix), ICP: "京icp备998877号", TaskID: suffix,`
- 288: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Name: fmt.Sprintf("ICP Matched App %d", suffix), ICP: " 京 ICP备 998877号 ", TaskID: suffix,`

#### db/config.go
- 646: **원문 문서 경로 DNT** — `// (영속 PTY 세션) tool family + Bash 프롬프트 연동(docs/交互式shell设计.md §14.2 참조).`

#### db/config_test.go
- 109: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `v1, err := d.SavePrompt(ag.ID, "你是规划者 {{.Goal}}", "init", "test")`
- 113: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `v2, _ := d.SavePrompt(ag.ID, "你是规划者 v2 {{.Goal}} {{.Scope}}", "edit", "test")`
- 118: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if cur != "你是规划者 v2 {{.Goal}} {{.Scope}}" {`

#### db/exploration_test.go
- 15: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("test", "拿下测试目标")`
- 259: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("test", "id 搜索")`

#### db/finding_assets_test.go
- 58: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `tk, err := d.CreateTask("资产树测试", "目标", nil, 0, 0)`
- 81: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "XSS", "反射型 XSS", "high", "s", "e", "w", []int64{epID}); err != nil {`
- 85: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "Info", "信息泄露", "low", "s", "e", "w", []int64{svcID}); err != nil {`
- 89: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "Misc", "孤儿", "medium", "s", "e", "w", []int64{999000111}); err != nil {`
- 164: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `tk, err := d.CreateTask("资产筛选测试", "目标", nil, 0, 0)`
- 178: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "A", "子域名上的", "high", "s", "e", "w", []int64{subID}); err != nil {`
- 181: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "B", "别的根域名上的", "high", "s", "e", "w", []int64{otherID}); err != nil {`
- 184: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "C", "没有资产的", "high", "s", "e", "w", nil); err != nil {`
- 189: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.AddFinding(tk.ID, 0, "D", "资产已删除", "high", "s", "e", "w", []int64{999000333}); err != nil {`

#### db/finding_retests_test.go
- 18: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `fid, err := d.AddFinding(0, 0, "retest-test", "测试漏洞", "high", "original summary", "original evidence", "worker", nil)`
- 38: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `r, c, created, err := d.CreateFindingRetest(t.Context(), fid, "  修复版本 v2  ")`
- 95: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.RecordFindingRetestResult(ctx, *r.ConversationID, "fixed", "修复验证通过", "正常对照可用，原触发条件失效"); err != nil {`

#### db/findings_test.go
- 22: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `tk, err := d.CreateTask("删除漏洞测试", "目标", nil, 0, 0)`
- 34: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `fid, err := d.AddFinding(tk.ID, nodeID, "XSS", "反射型 XSS", "high", "summary", "poc", "worker", nil)`
- 82: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{"critical", "pending", "严重漏洞标题"},`
- 116: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if p1[0].Name != "严重漏洞标题" {`
- 117: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `t.Fatalf("name round-trip: want 严重漏洞标题, got %q", p1[0].Name)`
- 162: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if one == nil || one.ID != ids[0] || one.Severity != "critical" || one.Name != "严重漏洞标题" {`
- 169: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.Exec(ˋUPDATE findings SET report=$1 WHERE id=$2ˋ, "# 报告\n正文", ids[0]); err != nil {`
- 172: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if one, _ = d.GetFinding(ids[0]); one.Report != "# 报告\n正文" {`
- 180: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "name", query: "严重漏洞标题", want: 1},`
- 183: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{name: "report", query: "正文", want: 1},`
- 442: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `auditInput := Activity{Worker: "system", Kind: "text", Summary: "人工提交漏洞深入利用意图", Detail: "验证可利用性并形成证据链"}`
- 443: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intentID, audit, err := store.AddFindingFollowUpIntent(findingID, findingNodeID, "验证可利用性并形成证据链", auditInput)`
- 447: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `secondID, _, err := store.AddFindingFollowUpIntent(findingID, findingNodeID, "从另一条路径深入", Activity{`
- 448: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Worker: "system", Kind: "text", Summary: "人工提交漏洞深入利用意图", Detail: "从另一条路径深入",`

#### db/intent_delete_test.go
- 47: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("soft delete", "假删除")`
- 54: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intent := mustIntent(t, es, "待删意图")`
- 58: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `summary, err := es.SoftDeleteIntent(intent, "方向判断错误")`
- 62: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if summary != "待删意图" {`
- 63: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `t.Fatalf("summary=%q, want 待删意图", summary)`
- 69: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if n.State != StateIntentDeleted || n.DeleteReason != "方向判断错误" {`
- 70: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `t.Fatalf("state=%q delete_reason=%q, want deleted/方向判断错误", n.State, n.DeleteReason)`
- 73: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `openIntent := mustIntent(t, es, "待领意图")`
- 74: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := es.SoftDeleteIntent(openIntent, "方向不需要了"); err != nil {`
- 82: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := es.SoftDeleteIntent(intent, "再删"); err == nil {`
- 94: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("hard cascade", "级联删除")`
- 102: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intent1 := mustIntent(t, es, "根意图")`
- 103: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `fact1 := mustNode(t, es, KindFact, "事实1")`
- 105: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intent2 := mustIntent(t, es, "衍生意图")`
- 107: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `fact2 := mustNode(t, es, KindFact, "事实2")`
- 131: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `expID, err := d.CreateExploration("hard preserve", "保留共享/目标")`
- 138: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `goal, err := es.AddGoal(map[string]any{"text": "拿下后台"}, "human")`
- 143: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intent1 := mustIntent(t, es, "待删意图")`
- 144: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intentX := mustIntent(t, es, "旁路意图")`
- 145: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `finding := mustNode(t, es, KindFinding, "漏洞")`
- 148: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `shared := mustNode(t, es, KindFact, "共享事实")`
- 151: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `intent2 := mustIntent(t, es, "由共享事实衍生")`

#### db/intercept.go
- 40: **[模型] 계약 DNT** — `Reason         string          ˋjson:"reason"ˋ // 규칙 message 또는 모델 판정 사유(접두 [模型])`
- 231: **[模型] 계약 DNT** — `WHEN ip.reason LIKE '[模型]%' THEN 'model' ELSE 'unknown' END)ˋ`
- 251: **[模型] 계약 DNT** — `if strings.HasPrefix(reason, "[模型]") {`

#### db/intercept_detail_test.go
- 19: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `id, err := d.CreateInterceptPending(0, 0, "approval-detail-test", "test", "Write", []byte(ˋ{"path":"report.md"}ˋ), "[模型] 请确认", audit)`
- 52: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `ok, err := d.ResolveIntercept(id, "allowed", "allow", "人工允许执行")`
- 93: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.ResolveIntercept(id, "denied", "deny", "人工拒绝"); err != nil {`
- 104: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if _, err := d.ResolveIntercept(id, "timeout", "allow", "超时允许"); err != nil {`

#### db/intercept_filter_test.go
- 41: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `reason = "[模型] fixture"`

#### db/llmretry.go
- 8: **원문 문서 경로 DNT** — `// LLM 재시도 정책: 5계층 재시도의 '횟수 + 간격' 전역 설정. docs/LLM重试设计.md 참조.`

#### db/notification_test.go
- 33: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Name:       "测试渠道-" + t.Name(),`
- 159: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `highSQL := addTestEvent(t, d, notify.EventFindingCreated, 1001, notify.Snapshot{Severity: "high", VulnClass: "SQL注入"})`
- 208: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `pick := newTestChannel(t, d, notify.KindDingTalk, NotifyModeRealtime, ˋ{"vulnclass_include":["绝不匹配的类型"]}ˋ)`
- 396: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.RescheduleDeliveries(ctx, allIDs, time.Second, "模拟失败"); err != nil {`
- 429: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.RescheduleDeliveries(ctx, []int64{id}, time.Second, "网络抖动"); err != nil {`
- 436: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if state != NotifyStatePending || lastErr != "网络抖动" {`
- 440: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.FailDeliveries(ctx, []int64{id}, "重试耗尽"); err != nil {`
- 480: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `addTestEvent(t, d, notify.EventFindingCreated, int64(8000+i), notify.Snapshot{Severity: "high", Name: "分页测试"})`
- 526: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `tk, err := d.CreateTask("通知状态变更测试", "目标", nil, 0, 0)`
- 533: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `TaskID: tk.ID, Worker: "test", VulnClass: "SQL注入", Name: "状态变更用例",`
- 534: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Severity: "high", Summary: "摘要",`
- 586: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if snap.VulnClass != "SQL注入" || snap.Severity != "high" || snap.Name != "状态变更用例" {`
- 633: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Name:       "CRUD 往返",`
- 650: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if got.Mode != NotifyModeDigest || got.RatePerMin != 42 || got.Name != "CRUD 往返" {`
- 672: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `got.Name = "改名了"`
- 682: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if after.Name != "改名了" || after.IsEnabled() {`
- 712: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Name: "不限流", Kind: notify.KindDingTalk, RatePerMin: 0,`
- 733: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Name: "负限流", Kind: notify.KindDingTalk, RatePerMin: -1,`
- 814: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `tk, err := d.CreateTask("复测推送测试", "目标", nil, 0, 0)`
- 821: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `TaskID: tk.ID, Worker: "test", VulnClass: "SQL注入", Name: "复测目标",`
- 822: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `Severity: "high", Summary: "摘要",`
- 830: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `rt, _, _, err := d.CreateFindingRetest(ctx, f.FindingID, "复核")`
- 841: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.RecordFindingRetestResult(ctx, *rt.ConversationID, "fixed", "已修复", "证据"); err != nil {`
- 871: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if snap.Name != "复测目标" || snap.Severity != "high" {`

#### db/schema.sql
- 1: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- ARTEX PostgreSQL schema (单一数据源)`
- 2: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 幂等：可重复执行（IF NOT EXISTS / OR REPLACE / DROP TRIGGER IF EXISTS）。`
- 5: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 0. 通用：updated_at 触发器`
- 11: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 安全的 text→inet 转换：非法值返回 NULL 而不是抛 22P02。assets.ip 是自由文本`
- 12: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- (Agent / 资产 API 可能写进主机名)，裸转 a.ip::inet 会让单独一行脏数据把整条`
- 13: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 企业归属重算语句打挂。调用方用 try_inet(...) IS NULL 找出这些行并告警。`
- 14: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 不用 pg_input_is_valid 是因为那要 PG16+，这里要兼容更老的存量库。`
- 24: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- A. 资产层：companies / assets / company_scope`
- 151: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- B. 推理探索层`
- 196: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 意图假删除(soft delete):state='deleted' 时,delete_reason 记用户填写的删除原因。`
- 285: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- runtime by the main agent + 总览「约束管理」. Injected into the planner/worker system`
- 343: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- only the extra segments created by "新建会话" (seq >= 1). The current segment is`
- 374: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 思考参数拆成两个独立字段：thinking_type=思考开关(''/disabled/enabled)，`
- 375: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- reasoning_effort=思考强度(''/low/medium/high/xhigh/max)，互不牵连。`
- 379: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 轮询(故障转移)参数，见 docs/LLM轮询设计.md：`
- 380: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   priority     顺位，越大越先被选中；激活配置(is_default)永远排链首，与本值无关。`
- 381: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   pool_exclude true=不作为故障转移目标(仍可被 agent/任务显式绑定使用)。`
- 384: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- streaming=true(默认)走流式 SSE；false 走真·非流式(stream:false，一次性 JSON)。`
- 386: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 单次回复的输出上限(token)。0=不发送该字段，由服务端默认值决定——保持既有行为。`
- 387: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 与 context_window_k(模型总容量，仅本地用于压缩阈值)是两回事：本值会随请求发出。`
- 389: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 输出上限用哪个请求字段名，仅对 format='openai' 生效：`
- 390: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   ''                      = max_tokens(默认，兼容绝大多数网关)`
- 391: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   'max_completion_tokens' = 新字段；OpenAI 推理模型(o 系列/GPT-5)只认它，`
- 392: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--                             发 max_tokens 会被 unsupported_parameter 拒绝。`
- 393: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- anthropic(max_tokens 必填)与 openai-responses(max_output_tokens)自带字段名，不受此值影响。`
- 395: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 自定义会话头：非空时每次请求带一个该名字的 HTTP 头，头值=当前运行的 session id`
- 396: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- (chat 会话/worker 意图)。用于某些按 session-id 头做提示缓存/粘性路由的网关。''=不发送。`
- 398: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 重试覆盖：次数 0=用全局默认/-1=关闭/>0=该值；间隔 0=用默认指数退避/>0=固定毫秒。`
- 399: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 三组分别对应建连重试、空响应重试、同 provider 安全窗口重试，详见下方 ALTER 处注释。`
- 413: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 轮询顺位/排除标记；补旧库。默认 0 / false = 全部配置都参与轮询。`
- 416: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 流式开关；补旧库。默认 true = 保持既有的流式行为，旧配置无感升级。`
- 418: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 放开 format 约束以容纳 openai-responses(OpenAI Responses API)；补旧库。`
- 419: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 每次启动执行,幂等:先删旧 CHECK 再建含三值的新 CHECK。`
- 424: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 输出上限及其字段名；补旧库。默认 0 / '' = 不发送上限、沿用 max_tokens 字段名，`
- 425: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 旧配置行为完全不变。`
- 428: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 同 format：先删再建，保证每次启动幂等。`
- 435: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 自定义会话头名；补旧库。默认 '' = 不发送，旧配置行为不变。`
- 438: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 单配置的重试覆盖（见 docs/LLM重试设计.md）。三组各自一对「次数 + 固定间隔」，`
- 439: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 语义统一：次数 0=沿用全局默认、-1=关闭该层重试、>0=用该值；间隔 0=沿用该层的`
- 440: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 默认指数退避、>0=改用这个固定毫秒数。全部默认 0，所以旧库/旧配置行为不变。`
- 441: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   connect = 建连重试（SDK doStream：连接重置/超时/429/5xx，流开始前）`
- 442: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   empty   = 空响应重试（SDK：完成但没有任何 content block，仅 openai 格式）`
- 443: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   stream  = 同 provider 安全窗口重试（本项目 task_llm：未交付输出前的断流重放）`
- 450: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 同 format：先删再建，保证每次启动幂等。次数下限 -1(关闭)，间隔不能为负。`
- 456: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 思考开关字段 thinking_type，从旧的单一 reasoning_effort 语义一次性拆分而来。`
- 457: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- schema.sql 每次启动都执行，故迁移必须只跑一次：仅当该列尚不存在时才回填，`
- 458: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 否则每次启动都会把用户后来手动设的组合覆盖回去。旧 reasoning_effort 语义：`
- 459: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   'off'                    → 显式关闭  → thinking_type='disabled'，强度清空`
- 460: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   'low/medium/high/max'    → 开启+强度 → thinking_type='enabled'，强度保留`
- 461: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   ''                       → 不发送    → 两者皆空(默认)`
- 476: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- LLM 轮询熔断状态：某个配置连续失败(余额不足/key 失效/限流)后进入冷却，冷却期内`
- 477: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 轮询直接跳过它。内存态为准，这里落库只为重启后不丢冷却窗口——加载时只取尚未`
- 478: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 到期的行(open_until > now)，已到期的自然回到"正常"，等下一次调用半开试探。`
- 481: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `fails       INTEGER NOT NULL DEFAULT 0,  -- 当前连续失败次数(成功即清零)`
- 482: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `trips       INTEGER NOT NULL DEFAULT 0,  -- 累计熔断次数,用于冷却时间指数退避`
- 483: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `open_until  TIMESTAMPTZ,                 -- 冷却截止;NULL/过期 = 未熔断`
- 489: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- D. 任务层`
- 541: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- planner 心跳触发间隔(秒);补旧库。默认 300s(5min)。见 docs/planner-trigger-impl-plan.md`
- 543: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 并发上限挂起态;补旧库。true=因并发上限排队、等待空位自动启动。`
- 545: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 资产覆盖度功能开关;补旧库。true(默认)=计算/展示测试覆盖度、自动累积测试范围、`
- 546: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 给 agent 开放 add_task_scope/list_untested_assets;false=全部关闭(见 task_scope.go)。`
- 547: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 存量任务默认 true 保持原行为;company 关联(task_scope kind=company)不受此开关影响。`
- 554: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 可选的任务名称;补旧库。空串=未命名,前端展示时回退到描述。`
- 619: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 预设的任务分类；分类删除时置空（与 tasks.category_id 一致，不阻断）。`
- 621: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 预设的任务级拦截/允许规则快照(AssetInterceptRuleInput 数组)；应用模板时灌进新任务。`
- 626: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 补旧库(已发版,加列带 IF NOT EXISTS)。`
- 670: **표시용 SQL 산문 누락(A076–A078)** — `SELECT task.id, NEW.id, 'system', '任务执行期间自动关联'`
- 688: **표시용 SQL 산문 누락(A076–A078)** — `SELECT task.id, asset.id, 'legacy', '由历史任务资产关联迁移'`
- 732: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 任务测试范围（资产覆盖度的分母 + 授权边界）。`
- 733: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   自动填(source='auto')：insertAssets 顶层按 worker 显式插入的资产类型加保守范围`
- 734: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--     （root_domain→root_domain，subdomain/service/endpoint→subdomain(host)，ip→ip）；`
- 735: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--     side-effect 派生的资产不入范围（钩子在 handler 顶层，派生在 db 层内部）。`
- 736: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   agent 填(source='agent')：add_task_scope 加 company/root_domain/subdomain/ip。`
- 737: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 覆盖度 = 匹配 active 行的 assets（分母）中，被 fact 节点锚定过的占比（分子）。`
- 750: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 旧库升级：扩展任务范围，使其与企业范围的单文本框识别能力一致。`
- 755: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 去重：同一 task 的同一条范围只存一次（自动填批量插入靠它幂等）。`
- 764: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- E. Agents / 提示词模板 / 变量目录`
- 791: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 加列迁移(已发版,旧库升级补列;新库 CREATE 已含。迁移不带 CHECK:旧库存量安全 + 后端写入白名单兜底)。`
- 795: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- per-agent LLM 绑定(agent 级默认模型):列自初版即在上方 CREATE 中,此 ALTER 仅为极旧库兜底(幂等)。`
- 797: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- run_seconds 单次 run 墙钟默认 600→1200:只改列默认(影响将来新插入的行),不动旧库存量行。`
- 814: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 循环外键：agents.current_prompt_id → agent_prompts.id（需在两表创建后加）`
- 831: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- F. MCP 服务`
- 842: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `insecure    BOOLEAN NOT NULL DEFAULT false,  -- http: 跳过 TLS 证书校验(自签证书场景, issue #108)`
- 851: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 旧库补列(schema.sql 每次启动都会 Exec)。`
- 857: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 默认数据源占位：ScopeSentry 资产同步 MCP（地址与认证均留空、未启用）。`
- 858: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 供「资产同步」页检测数据源是否已配置；用户在页面填入 url 与 X-API-Key 后再启用。`
- 859: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 仅在缺失时插入，绝不覆盖用户已配置/已启用的服务器（schema.sql 每次启动都会 Exec）。`
- 875: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- G. 可见性：agent × mcp / skill`
- 897: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- Skill 调用账本（见 db/skill_usage.go）。一次 Skill() 调用一行，只记维度不记正文。`
- 898: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 刻意不设外键：任务/会话删除后统计仍要保留（与 llm_usage 同理），skill 本身也只是`
- 899: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 文件系统上的目录名，没有对应的表。`
- 910: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- false = 模型点名了一个不存在的 skill(未命中)。这类行同样保留：它反映"想用但没有"`
- 911: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 的缺口，是补 skill 的依据。`
- 917: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 工具调用账本（见 db/tool_usage.go）。一次实际 CoreTool.Call 一行，只记归属维度，`
- 918: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 不保存工具参数或返回内容。刻意不设外键，任务、会话或自定义工具删除后仍保留统计。`
- 933: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- H. 内置工具目录`
- 952: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- I. 会话（对话页）`
- 991: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- J. Agent 触发器`
- 1015: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 加列迁移(已发版,旧库升级补列;新库 CREATE 已含这些列,ALTER 为 no-op)。幂等,每次启动可重复执行。`
- 1031: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- K. 拦截规则`
- 1063: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 判定理由:规则命中时为规则 message;LLM 兜底判定时为模型给的简短理由(前缀 [模型])。`
- 1070: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 补旧库:reason 列(已发版,加列要带 IF NOT EXISTS)。`
- 1076: **표시용 SQL 산문 누락(A076–A078)** — `WHEN reason LIKE '[模型]%' THEN 'model' ELSE 'unknown' END WHERE decision_source='';`
- 1079: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- L. 漏洞发现持久化`
- 1086: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 漏洞名称(可读标题)；为空时前端回退展示 vulnclass。severity 取值：`
- 1087: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- critical 严重 / high 高 / medium 中 / low 低（不加 CHECK，与 status 一致由 server 白名单校验）。`
- 1094: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 处置状态：pending 待处理 / in_progress 处理中 / confirmed 已确认 / resolved 已处理 / fixed 已修复 /`
- 1095: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- false_positive 误报 / ignored 忽略 / duplicate 重复 / risk_accepted 风险接受。`
- 1096: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 取值不加 CHECK：旧库靠下面的 ALTER 补列,CHECK 无法回填,统一由 server 侧白名单校验。`
- 1098: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 漏洞详细报告(Markdown)；默认空,仅详情页读取/展示,不进列表接口以免 payload 膨胀。`
- 1108: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 「按资产」视图靠 asset_ids @> '[<id>]' 反查发现,没有这个 GIN 索引就是全表扫。`
- 1111: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 手动复测属于独立会话；结论与原漏洞处置状态分开保存。`
- 1131: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 删除会话保留复测记录，同时解除尚未结束的复测占用。`
- 1134: **표시용 SQL 산문 누락(A076–A078)** — `UPDATE finding_retests SET status='stopped', error='复测会话已删除', finished_at=now()`
- 1177: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- M. 后端日志持久化`
- 1232: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 资产拦截规则（全局黑名单）`
- 1233: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 独立于 §K 命令拦截(intercept_rules)：intercept_rules 匹配工具名/入参文本，`
- 1234: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 这张表匹配「目标资产」——全等/模糊的域名·IP·URL 以及 CIDR 网段。`
- 1235: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 仅存规则；具体的匹配/拦截逻辑在别处实现。`
- 1236: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- kind 七种：`
- 1237: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   exact_domain / exact_ip / exact_url  —— 全等匹配`
- 1238: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   fuzzy_domain / fuzzy_ip / fuzzy_url  —— 模糊匹配`
- 1239: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   cidr                                 —— CIDR 网段`
- 1260: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 任务级资产拦截/允许规则`
- 1261: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 与全局 asset_intercept_rules 同构（kind/pattern/note/enabled），但按 task_id`
- 1262: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 关联、随任务级联删除；创建任务时录入、任务详情里可编辑。`
- 1263: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- action: 'block'=拦截(禁止测试)  'allow'=允许(白名单)。`
- 1264: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 执行判定：先按 拦截规则(全局 ∪ 任务block) 匹配，命中即禁止；未命中且该任务存在`
- 1265: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 启用的 allow 规则时，须命中某条 allow 才放行，否则「不允许测试」。`
- 1282: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 补旧库(本会话早前建过该表、无 action 列)：加列(带 IF NOT EXISTS)。`
- 1289: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- M. 漏洞 IM 推送`
- 1291: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 三张表刻意分开，核心是**爆炸半径**：写漏洞的那个事务(RecordFindingTx，`
- 1292: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 持任务行锁)只允许做一次盲 INSERT，不读渠道表、不跑用户的过滤规则。否则`
- 1293: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 一条配错的 webhook 过滤条件就能污染/中止事务，导致漏洞存不进去。`
- 1295: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   notification_channels   渠道实例配置(可变、含凭据、UI 管理)`
- 1296: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   notification_events     事件事实(写漏洞事务内盲插，含渲染快照)`
- 1297: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   notification_deliveries 投递任务(事务外 fan-out 产生，承载状态/重试/批次)`
- 1300: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 渠道实例：同一 kind 可配任意多个(如「应急群」「日常群」各一个钉钉机器人)。`
- 1301: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- kind 取值由 server 侧白名单校验，不加 CHECK：与 findings.status 同理，`
- 1302: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 后续加渠道不应要求改表结构。`
- 1306: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- dingtalk 钉钉 / feishu 飞书 / wecom 企业微信 / webhook 通用 / telegram / email`
- 1309: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 凭据(明文存储，UI 掩码回显；见 server 侧 maskChannelSecrets)。六种渠道字段差异极大，`
- 1310: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 统一 JSONB + Go 侧按 kind 严格校验，避免为每渠道加一堆 NULL 列：`
- 1318: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 推送时机：realtime 命中即推 / digest 进批次按全局周期汇总成一条。`
- 1320: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 过滤条件，字段全部可选(缺省=不过滤)：`
- 1322: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   task_ids/asset_ids 空数组=不限；非空则须交集非空`
- 1323: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   vulnclass_include/exclude 关键词数组(大小写不敏感子串)；include 空=全收`
- 1324: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `--   on_status_change   bool，仅 realtime 模式有意义`
- 1326: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 每分钟投递上限；0=不限流。默认 20 对齐钉钉/企微官方硬限。`
- 1327: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 超限不丢消息，只把投递推迟到下一个 tick。`
- 1336: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 事件事实。由 RecordFindingTx / 状态变更事务**同事务**写入，保证「漏洞落库」`
- 1337: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 与「推送任务存在」原子一致——不存在提交成功但没入队、消息永久丢失的窗口。`
- 1338: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- snapshot 刻意冗余：漏洞事后会被改名/改级别/改状态，推送内容应反映「事发当时」，`
- 1339: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 且 fan-out 与渲染不必回查 findings/tasks/assets 多张表。`
- 1340: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- finding 删除后事件不级联删除：与 findings 表「任务删除仍独立留存」的语义一致。`
- 1347: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- fan-out 幂等标记：dispatcher 按此列取待分派事件，处理完置 true。`
- 1348: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 用列而非删行，以便投递历史能回溯到事件。`
- 1355: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 投递任务：一条事件 × 一个启用渠道 = 一行。fan-out 在事务外做，所以渠道`
- 1356: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 后开不会补历史(与 agent_triggers 的「迟开 trigger 不补历史」语义一致，`
- 1357: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 避免启用渠道时一次性刷屏历史积压)。`
- 1358: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- channel_id 级联删除：渠道配置都没了，其投递历史无意义。`
- 1363: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- pending 待发 / sent 已发 / failed 重试耗尽(可手动重发) / skipped 渠道停用或批次取消`
- 1364: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- pending 待发 / sending 已被某 dispatcher 领取(租约未到期) / sent 已发 /`
- 1365: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- failed 重试耗尽或永久失败(可手动重发) / skipped 渠道停用。取值不加 CHECK，`
- 1366: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 与 findings.status 同理，由 server 侧白名单校验。`
- 1369: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 兼作「下次可领取时间」与「租约到期时间」：领取时把它推到未来即构成租约，`
- 1370: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- 于是「租约未到期」与「未到重试时间」共用同一个条件表达，不需要额外的`
- 1371: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- lease_until 列。进程崩溃留下的 sending 行会因租约到期被下一轮重新领取。`
- 1374: **설명 주석 누락(A075), 내부 키/계약·경로는 보존** — `-- digest 模式同批次共享；realtime 恒为 NULL。整批渲染成一条消息后一起置 sent。`

#### db/task_archives_restore.go
- 680: **[模型] 계약 DNT** — `} else if reason, _ := row["reason"].(string); strings.HasPrefix(reason, "[模型]") {`

#### db/task_assets_test.go
- 29: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `{Kind: "icp", Value: " 京 ICP 备 12345678 号-1 "},`
- 107: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `mutation, err := d.Assets().AttachAssetsToTask(task.ID, []int64{assetID, assetID}, "授权资产清单第 3 项")`
- 118: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if assets[0].TaskSource != "manual" || assets[0].TaskSourceSummary != "授权资产清单第 3 项" {`
- 174: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `if err := d.Assets().SetTaskAssetSource(source.ID, assetID, "agent", "Worker 通过 insert_assets 登记", &nodeID); err != nil {`

#### db/task_context_unit_test.go
- 11: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `input := strings.Repeat("额度不足", 400)`

#### db/tasks.go
- 41: **원문 문서 경로 DNT** — `// 작업 수준 타임아웃(docs/任务级超时与收尾设计.md 참조).`

#### db/tasks_test.go
- 307: **입력/출력·Unicode·정규화·라운드트립 fixture/원문값 보존; 운영 출력 불일치로 보지 않음** — `last, err := d.MarkTaskLLMProfileQuotaExhausted(child.ID, profileIDs[2], "余额不足")`

## 전체 커버리지와 한계
- 전체603개 추적 파일 목록은 기존 감사 부록과 이번 /tmp/artex-audit-integrated/all-inventory.json을 대조했다. 경로 추가/삭제 없음. U12 제외를 해제하고96개 전부 장부에 포함했다. 이미지18개 제외 근거와 소유 미정 agent 목록은 이전과 동일하다.
- 비U12 파일의 기존 후보/수정 후보를 재대조했지만 CJK 없는374개 전체 및 db 한국어 문장의 Git 원문 의미 전수 대조는 여전히 미완료다. 고유명·예시·정규식·JSON·SQL이 섞인 문자 검색만으로 완료/누락을 판정하지 않았다.
- 새로운 U12 완료 기록은 번역 정확성이나 테스트 통과를 확증하지 않는다. 실제 DB·migration 실행·모델·UI/SDK 동작은 미확인. 기존 PASS 기록은 팀원 이력일 뿐 이 감사 결과가 아니다.
- 실행 검증: **사용자 요청으로 미실행**. 저장소·진행 문서·용어집·index 보존. 기존 감사 보고서를 덮어쓰지 않고 새 통합 .md/.tsv만 작성.
