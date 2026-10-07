# 감사 원본 근거 복구와 현재 상태

현재 기준 HEAD: `f3611044ad4e29acee8878fd4f3a0963c3bd573b`. 현재 브랜치 `work/ko-translation`, 로컬 origin 추적 참조 `f3611044ad4e29acee8878fd4f3a0963c3bd573b`. fetch하지 않았으며 원격 최신성이나 push 통신 성공은 확인하지 않았다. f361104의 mock 안내·ProxyAddr 설명 변경은 로컬 HEAD에 포함된다.

## 원본 보관

아래 파일은 전달본과 바이트/해시가 동일하다. 기존 목적지가 없음을 확인한 뒤 복사했으며 원본을 편집하지 않았다. `audit-input/`과 Windows `:Zone.Identifier` 4개는 전달용 원본/첨부 메타데이터로 그대로 남긴다. stage·commit에 포함하지 않았다.

| 전달 파일 | 보관 파일 | SHA256 | 바이트 |
|---|---|---|---:|
| `artex-translation-audit-integrated.md` | [integrated-original.md](integrated-original.md) | `c53cd7ecc7d2c6e6d448f271982fb03f4456569a7e846c0af109fb3e8f5b9a79` | 108277 |
| `artex-translation-audit-integrated.tsv` | [integrated-original.tsv](integrated-original.tsv) | `fc493a3613ca0b1873e157082d2116781a5aa1dfc791aef246413b40341709c0` | 154096 |
| `artex-translation-audit-status(3).md` | [status-before-final-classification.md](status-before-final-classification.md) | `01cf053494e985b34eb5ef2d6b296caf394023e6d78ba5e12da15de649d65711` | 63059 |
| `artex-translation-audit-status(3).tsv` | [status-before-final-classification.tsv](status-before-final-classification.tsv) | `94d239c38f91e63bd50ffe67ad5d0032fa09644f5a691052de83e364ebd679da` | 111754 |

integrated 원본 기준은 `474300b6d355fc1c8cf7051d644c52a7abae1fd4`, status 원본의 선언 기준은 `d4d288b480fcde527c3302e2c324be6d7816176a`이며 S 항목 후속에는 `564c80439114a88c58cb1c2dcaf3593199c48a41` 기준 기록도 있다. 명시적인 정확 작성 시각은 확보한 원본에 없어 미확인으로 남긴다. 파일 전달 시각·mtime을 작성 시각으로 단정하지 않는다. 원본은 과거 기준이며 현재 미해결 목록으로 그대로 옮기지 않는다.

## 복구 범위와 집계

A001–A090 90개와 S001–S006 6개, 총 96개 ID의 원본 근거를 복구했다. 원본 ID·발견 문구·위치·판단·이력·연결은 current-status.tsv의 `integrated_*`와 `previous_*` 열에 보존한다. `original_location`/`previous_location`은 과거 위치, `current_location`은 이번 소스 대조로 확인한 현재 위치다. 원본 근거 복구 여부와 소스 재대조 범위를 별도 열로 기록했다.

| 주 상태 | A 항목 | S 항목 | 합계 |
|---|---:|---:|---:|
| 수정 완료 / 정적 대응 확인 | 71 | 6 | 77 |
| DNT 또는 원문 예시 보존 | 15 | 0 | 15 |
| 원문 기능·정책·호환성 문제 | 4 | 0 | 4 |

이 96개 ID 범위의 실제 번역·표기 수정 잔여 및 근거 미복구 ID는 0개다. 이 수치는 저장소 전체에 새 누락이 전혀 없다는 주장이나 테스트 통과를 뜻하지 않는다. 96개 모두 실행 검증은 **사용자 요청으로 미실행**이다.

A049는 과거 status의 DNT 주 분류에서 원문 기능 문제 주 분류로 정정했다. 원본 integrated와 status의 secondary_status 양쪽에 남아 있던 depth/coverage 문제를 반영한 것이며, 정규식을 번역하거나 기능 문제를 해결한 것이 아니다. 이에 따라 단순 과거 집계의 DNT 16/원문 문제 3 대신 최종 A 집계는 보존 15/원문 문제 4다.

## 주요 정정과 유지

- A048: 원본으로 U15 로그인·실제 세션 정책 충돌 연결 확인. 미커밋 PROGRESS의 연결 미확인 설명만 정정했다.
- A049.1: NEGATIVE_RE는 원문 계약/매칭 정규식 보존. A049.2: depth/coverage forward·stub·L2 차이는 미해결 기능 문제.
- A050.1: 승인 편집 문자열 13곳은 f361104로 수정. 이름은 `${a.name}` 편집 프롬프트에도 삽입된다. A050.2: 지정 8곳은 원문 기록 예시 보존. A050.3: mock/실제 api-recon 범위 차이는 별도 원문 문제.
- A052: 腾讯 TSec Benchmark 고유명 원문 보존. 공식 영문 표기는 미확인이며 번역 누락으로 집계하지 않는다.
- A090: 승인 ProxyAddr description 정정 완료. 메타데이터 seed upsert 경로와 사용자 프롬프트 본문을 구분한다.
- A059: 新对话 비교 계약 보존이며 UI 표시 문제 A025는 해소. A062: 내장 기본값 공유 표기는 이미 정적 대응 확인. 원본대로 ID를 유지한다.
- A089: 과거 U12 작업/PASS 이력과 후속 주석/표시 보완을 구분한다. 과거 실행 기록을 이번 결과로 인용하지 않는다.

## 확인 방법과 한계

원본 status의 번호가 붙은 전체 소스 인용을 현재 파일과 정확 텍스트로 대조했다. A050의 승인 13곳과 A090의 승인 1곳을 제외한 인용은 현재 소스에서 일치했다. 두 항목은 후속 커밋과 현재 문자열로 다시 확인했다. S 항목은 지정 설명·표기·구분자·주석 범위를 별도로 재대조했다. 이 대조는 인용에 없는 모든 문장이나 파일의 기능/정책 의미 전수 검증을 대체하지 않는다. 소비 경로와 테스트는 정적 근거이며 실제 실행 결과가 아니다.

기준 문서 4종의 최신 로컬 기준을 적용했다. 원본의 테스트 PASS 기록은 과거 이력으로 보관하며, 전체 한국어 문장의 원문 의미 전수 대조는 미완료다. 실제 SDK·DB·모델 출력과 브라우저/앱 동작은 미확인이다.

## 항목별 상태

전체 원문·현재 인용·소비부·처리 이력은 [current-status.tsv](current-status.tsv)를 참조한다. 원본 MD/TSV는 수정하지 않았으며 다음 표는 현재 결론만 요약한다.

| ID | 현재 상태 | 대상 | 처리 커밋/기록 |
|---|---|---|---|
| A001 | 수정 완료 / 정적 대응 확인 | `agent/finding_recorder_test.go` | e7f74ef |
| A002 | 수정 완료 / 정적 대응 확인 | `server/finding_traffic_test.go` | e7f74ef |
| A003 | 수정 완료 / 정적 대응 확인 | `server/finding_workflow_test.go` | e7f74ef |
| A004 | 수정 완료 / 정적 대응 확인 | `server/finding_workflow_test.go` | e7f74ef |
| A005 | 수정 완료 / 정적 대응 확인 | `server/intercept_detail_test.go` | e7f74ef |
| A006 | 수정 완료 / 정적 대응 확인 | `server/task_categories_test.go` | e7f74ef |
| A007 | 수정 완료 / 정적 대응 확인 | `server/chat_mentions_test.go` | e7f74ef |
| A008 | 수정 완료 / 정적 대응 확인 | `server/chat_mentions_test.go` | e7f74ef |
| A009 | 수정 완료 / 정적 대응 확인 | `server/skill_upload_test.go` | e7f74ef |
| A010 | 수정 완료 / 정적 대응 확인 | `server/skill_upload_test.go` | e7f74ef |
| A011 | 수정 완료 / 정적 대응 확인 | `web/next.config.mjs` | 92ad854 |
| A012 | 수정 완료 / 정적 대응 확인 | `agent/capture_usage_test.go` | 92ad854 |
| A013 | 수정 완료 / 정적 대응 확인 | `agent/coldgraph_test.go` | 92ad854 |
| A014 | 수정 완료 / 정적 대응 확인 | `agent/insert_assets_test.go` | 92ad854 |
| A015 | 수정 완료 / 정적 대응 확인 | `server/llmretry_test.go` | 92ad854 |
| A016 | 수정 완료 / 정적 대응 확인 | `server/task_archive_package_test.go` | 92ad854 |
| A017 | 수정 완료 / 정적 대응 확인 | `server/trigger_merge_test.go` | 92ad854 |
| A018 | 수정 완료 / 정적 대응 확인 | `server/engine_emptyturn_test.go` | 92ad854 |
| A019 | 수정 완료 / 정적 대응 확인 | `server/update_test.go` | 92ad854 |
| A020 | 수정 완료 / 정적 대응 확인 | `server/notify_api_test.go` | 92ad854 |
| A021 | 수정 완료 / 정적 대응 확인 | `server/skill_upload_test.go` | 92ad854 |
| A022 | 수정 완료 / 정적 대응 확인 | `sidequestion/README.md` | 92ad854 |
| A023 | 수정 완료 / 정적 대응 확인 | `sidequestion/CONTEXT_BUDGET.md` | 92ad854 |
| A024 | 수정 완료 / 정적 대응 확인 | `sidequestion/VALIDATION.md` | 92ad854 |
| A025 | 수정 완료 / 정적 대응 확인 | `web/src/app/(main)/chat/page.tsx` | 5d58222 |
| A026 | 수정 완료 / 정적 대응 확인 | `web/src/lib/mock/handler.ts` | 5d58222 |
| A027 | 수정 완료 / 정적 대응 확인 | `server/llmpool.go` | 451bd13; 478544d (추가 문맥 보완) |
| A028 | 수정 완료 / 정적 대응 확인 | `server/server_mgmt.go` | 451bd13; 478544d (추가 문맥 보완) |
| A029 | 수정 완료 / 정적 대응 확인 | `server/llmpool.go` | 451bd13; 478544d (추가 문맥 보완) |
| A030 | 수정 완료 / 정적 대응 확인 | `server/llmretry.go` | 451bd13; 478544d (추가 문맥 보완) |
| A031 | 수정 완료 / 정적 대응 확인 | `server/server.go` | 451bd13; 478544d (추가 문맥 보완) |
| A032 | 수정 완료 / 정적 대응 확인 | `server/task_llm.go` | 451bd13; 478544d (추가 문맥 보완) |
| A033 | 수정 완료 / 정적 대응 확인 | `server/server_mgmt.go` | 451bd13; 478544d (추가 문맥 보완) |
| A034 | 수정 완료 / 정적 대응 확인 | `server/dto.go` | 451bd13; 478544d (추가 문맥 보완) |
| A035 | 수정 완료 / 정적 대응 확인 | `server/llmretry.go` | 451bd13; 478544d (추가 문맥 보완) |
| A036 | 수정 완료 / 정적 대응 확인 | `agent/planner.go` | 451bd13; 478544d (추가 문맥 보완) |
| A037 | 수정 완료 / 정적 대응 확인 | `agent/compaction.go` | 451bd13; 478544d (추가 문맥 보완) |
| A038 | 수정 완료 / 정적 대응 확인 | `web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx` | 451bd13; 478544d (추가 문맥 보완) |
| A039 | 수정 완료 / 정적 대응 확인 | `agent/retester.go` | 451bd13; 478544d (추가 문맥 보완) |
| A040 | 수정 완료 / 정적 대응 확인 | `agent/finding_workflow.go` | 451bd13; 478544d (추가 문맥 보완) |
| A041 | 수정 완료 / 정적 대응 확인 | `sidequestion/service.go` | 451bd13; 478544d (추가 문맥 보완) |
| A042 | 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go` | 451bd13; 478544d (추가 문맥 보완) |
| A043 | 수정 완료 / 정적 대응 확인 | `README.ko.md` | 451bd13; 478544d (추가 문맥 보완) |
| A044 | 수정 완료 / 정적 대응 확인 | `web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx` | 451bd13; 478544d (추가 문맥 보완) |
| A045 | 수정 완료 / 정적 대응 확인 | `README.ko.md` | 451bd13; 478544d (추가 문맥 보완) |
| A046 | 원문 기능·정책·호환성 문제 | `server/intercept_live_test.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A047 | 원문 기능·정책·호환성 문제 | `intercept/prompt.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A048 | 원문 기능·정책·호환성 문제 | `skills/api-recon/SKILL.md` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A049 | 원문 기능·정책·호환성 문제 | `skills/api-recon/scripts/preload.js` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A050 | DNT 또는 원문 예시 보존 | `web/src/lib/mock/data.ts` | f361104; 이번 미커밋 보완 11 (최종 분류) |
| A051 | 수정 완료 / 정적 대응 확인 | `web/src/components/mention-textarea.tsx` | d4d288b |
| A052 | DNT 또는 원문 예시 보존 | `docker-compose.bench.yml` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A053 | 수정 완료 / 정적 대응 확인 | `server/intercept.go` | d4d288b |
| A054 | DNT 또는 원문 예시 보존 | `intercept/intercept.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A055 | DNT 또는 원문 예시 보존 | `intercept/intercept.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A056 | DNT 또는 원문 예시 보존 | `server/task_archives.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A057 | DNT 또는 원문 예시 보존 | `server/server_mgmt.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A058 | DNT 또는 원문 예시 보존 | `intercept/prompt.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A059 | DNT 또는 원문 예시 보존 | `server/conversations.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A060 | DNT 또는 원문 예시 보존 | `server/orchestration.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A061 | DNT 또는 원문 예시 보존 | `skills/api-recon/reference.md` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A062 | 수정 완료 / 정적 대응 확인 | `server/finding_retests.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A063 | DNT 또는 원문 예시 보존 | `report/findings.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A064 | DNT 또는 원문 예시 보존 | `agent/terminalreason_test.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A065 | DNT 또는 원문 예시 보존 | `skills/scopesentry/SKILL.md` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A066 | DNT 또는 원문 예시 보존 | `intercept/prompt.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A067 | DNT 또는 원문 예시 보존 | `agent/tools_insert.go` | 원본 이력 및 PROGRESS 해당 U 기록; 이번 미커밋 보완 11 (근거 연결 정정) |
| A068 | 수정 완료 / 정적 대응 확인 | `notify/email_test.go` | 9262641 |
| A069 | 수정 완료 / 정적 대응 확인 | `notify/filter_test.go` | 9262641 |
| A070 | 수정 완료 / 정적 대응 확인 | `notify/http_test.go` | 9262641 |
| A071 | 수정 완료 / 정적 대응 확인 | `notify/mask_test.go` | 9262641 |
| A072 | 수정 완료 / 정적 대응 확인 | `notify/pack_test.go` | 9262641 |
| A073 | 수정 완료 / 정적 대응 확인 | `notify/redact_test.go` | 9262641 |
| A074 | 수정 완료 / 정적 대응 확인 | `translation/PROMPT_GUIDE.ko.md` | 9262641 |
| A075 | 수정 완료 / 정적 대응 확인 | `db/schema.sql` | 92ad854 |
| A076 | 수정 완료 / 정적 대응 확인 | `db/schema.sql` | 92ad854 |
| A077 | 수정 완료 / 정적 대응 확인 | `db/schema.sql` | 92ad854 |
| A078 | 수정 완료 / 정적 대응 확인 | `db/schema.sql` | 92ad854 |
| A079 | 수정 완료 / 정적 대응 확인 | `db/company_scope.go` | 92ad854 |
| A080 | 수정 완료 / 정적 대응 확인 | `db/companies.go` | 451bd13; 478544d (추가 문맥 보완) |
| A081 | 수정 완료 / 정적 대응 확인 | `db/config.go` | 451bd13; 478544d (추가 문맥 보완) |
| A082 | 수정 완료 / 정적 대응 확인 | `db/llmretry.go` | 451bd13; 478544d (추가 문맥 보완) |
| A083 | 수정 완료 / 정적 대응 확인 | `db/db.go` | 451bd13; 478544d (추가 문맥 보완) |
| A084 | 수정 완료 / 정적 대응 확인 | `db/db.go` | 451bd13; 478544d (추가 문맥 보완) |
| A085 | 수정 완료 / 정적 대응 확인 | `db/nkey.go` | 451bd13; 478544d (추가 문맥 보완) |
| A086 | 수정 완료 / 정적 대응 확인 | `db/notification.go` | 451bd13; 478544d (추가 문맥 보완) |
| A087 | 수정 완료 / 정적 대응 확인 | `db/notification_delivery.go` | 451bd13; 478544d (추가 문맥 보완) |
| A088 | 수정 완료 / 정적 대응 확인 | `db/notification.go` | 451bd13; 478544d (추가 문맥 보완) |
| A089 | 수정 완료 / 정적 대응 확인 | `translation/PROMPT_PROGRESS.ko.md` | 9262641 |
| A090 | 수정 완료 / 정적 대응 확인 | `db/db.go` | f361104; 이번 미커밋 보완 11 (최종 분류) |
| S001 | 수정 완료 / 정적 대응 확인 | `server/notify_api_test.go` | 564c804 |
| S002 | 수정 완료 / 정적 대응 확인 | `server/engine_emptyturn_test.go` | 564c804; a344782 |
| S003 | 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go` | 564c804 |
| S006 | 수정 완료 / 정적 대응 확인 | `README.ko.md` | 564c804 |
| S004 | 수정 완료 / 정적 대응 확인 | `web/src/app/(main)/system/intercept/page.tsx; web/src/components/agent-editor.tsx; agent/planner.go` | a344782 |
| S005 | 수정 완료 / 정적 대응 확인 | `server/conversations.go` | a344782 |

## 보존 결과

기존 TRANSLATION_PROMPT 변경·미추적 담당표·기타 미추적 파일·전달 원본·Git index를 보존한다. 변경 허용은 translation/audit의 감사 문서와 PROGRESS 보완 11의 근거 복구 설명 정정뿐이다. 소스·용어집·테스트·fixture·DB·설정은 수정하지 않았다. staged 및 진행 중 병합은 없다. 과거 커밋된 PROGRESS 기록은 보존한다. 상세 잔여는 [remaining-issues.md](remaining-issues.md) 참조.

최종 보존 확인: 작업 시작 해시 스냅샷의 기존 613개 파일/index 항목(PROGRESS 제외)은 모두 동일했다. audit-input의 4개 원본과 첨부 메타데이터도 동일하다. PROGRESS는 시작 상태 대비 731·732행의 설명 2줄만 정정했으며 HEAD에 커밋된 전체 바이트 접두부와 줄 수 732를 유지했다. 새 파일은 translation/audit의 7개 감사 문서뿐이다. staged diff는 비어 있고 MERGE_HEAD는 없다.


## M001–M007 핵심 프롬프트 승인 보완 (A/S 집계와 별도)

HEAD `b5468a7839ada1f322b4f1e4ff7d219481f15135` 이후의 현재 미커밋 작업이다. 최초 발견의 원문·수정 전 번역과 판단은 [prompt-semantic-review.md](prompt-semantic-review.md) 및 [prompt-semantic-findings.tsv](prompt-semantic-findings.tsv)에 이력으로 보존한다. 승인 표현의 소스 대응을 확인했다.

| 별도 ID | 처리 | 현재 위치 |
|---|---|---|
| M001 | 승인 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go:22` |
| M002 | 승인 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go:18` |
| M003 | 승인 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go:38` |
| M004 | 승인 수정 완료 / 정적 대응 확인 | `agent/promptcatalog.go:95` |
| M005 | 승인 수정 완료 / 정적 대응 확인 | `agent/planner.go:312` |
| M006 | 승인 수정 완료 / 정적 대응 확인 | `agent/planner.go:320` |
| M007 | 승인 수정 완료 / 정적 대응 확인 | `agent/planner.go:316,330` |

M은 7개 모두 승인 수정·정적 대응 완료이며 실행 미검증이다. M007은 316·330행 두 위치다. 기존 A/S 96개 집계(정적 대응 확인 77, 원문 보존 15, 정책·기능·호환성 문제 4)는 변경하지 않고 M과 합산하지 않는다. current-status.tsv의 기존 96개 행도 변경하지 않았다. A046–A049와 remaining-issues의 다른 후속 보류는 유지한다. 테스트 기대값 수정 필요는 확인한 범위에서 발견하지 못했다. 실제 DB/모델 동작·전체 저장소 의미 전수 대조는 미확인/미완료이고 실행 검증은 **사용자 요청으로 미실행**이다.
