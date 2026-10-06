# `/btw` 검증 기록

날짜: 2026-09-10. 브랜치: `codex/btw-side-question`. 기준: `8dae851b9b622f2ff2631f332fde9719d0b16fba`.

독립 PostgreSQL 테스트 DB 및 데이터 디렉터리를 사용했습니다. 실제 모델 자격 증명은 독립 테스트 환경에만 넣었으며 코드나 이 기록에 저장하지 않았고 제품 기본 모델도 바꾸지 않았습니다. Go 1.26.3, norma v0.3.6, Next.js 16.2.9.

실제 모델 대화·반환 객체·개발 검증문·Qwen 원문 검토 텍스트는 [validation-2026-09-10.json](validation-2026-09-10.json)에 저장했으며 API 자격 증명은 포함하지 않습니다.

## 개발 검사

| 범위 | 결과 | 증거 |
| --- | --- | --- |
| 구조화 메시지, 도구 파라미터 깊은 복사 | 통과 | `TestCheckpointDeepCopyAndBoundaries` |
| 요약 / 압축 요청이 덮어쓰지 않음, 완전한 응답과 종료 상태 공개, 일부 응답 제외 | 통과 | `TestCheckpointDeepCopyAndBoundaries`·`TestSnapshotExcludesPartialStreamAndSelectsPoolMember` |
| 실제 모델 풀 구성원 식별 정보 | 통과 | `TestSnapshotExcludesPartialStreamAndSelectsPoolMember` |
| 도구 쌍, 20쌍 재생, 예산에 맞춘 축소와 상한 초과 오류 | 통과 | `TestBuildRequestCompactionToolPairingAndBudget` |
| 메인·보조 질문 병렬 실행, 양방향 취소 격리 | 통과 | 차단하는 Provider, `TestMainSideConcurrencyAndIndependentCancellation` |
| 도구 실행 없음, 스트리밍 / 비스트리밍, 실패 시 이미 얻은 사용량 | 통과 | `TestServiceNoToolsAndUsageOnFailure` |
| 실제 norma ChatAgent + 로컬 Read 도구, 메인 transcript / 활동 격리 | 통과 | `TestSideActualChatCheckpointToolResultAndTranscriptIsolation`, 스트리밍 및 비스트리밍 하위 사례 |
| 영구 저장, 페이지네이션, 멱등성, 재시작 후 일부 답변 유지 | 통과 | `TestSideHistoryIdempotencyPagingAndRecovery` |
| 비우기와 늦은 쓰기의 경쟁, 부모 리소스 삭제, 버전 비교 | 통과 | `TestSideClearLateWritersAndDeletedParent` |
| MainAgent / Worker 아카이브 및 복구, v1/v2/v3 | 통과 | `TestSideTaskArchiveVersions` |
| 세 가지 부모 API, 인증, 리소스 귀속, Worker 논리 삭제 | 통과 | `TestSideHTTPGlobalLimitTaskWorkerAndDeletion`·`TestSideCheckpointPersistsBeforeAdmissionAndRestart` |
| 사용 중인 메인 세션에서도 보조 질문 가능, 독립 SSE 재연결 / 연결 해제, 취소, 비우기 | 통과 | `TestSideHTTPBusyIsolationClearAndReconnect` |
| 부모 세션마다 1 / 전역 4 동시 실행 | 통과 | 두 `TestSideHTTP…` 사례 |
| 제출 전 스냅샷 저장, 재시작 후 후속 질문, 기존 세션의 스냅샷 위조 불가 | 통과 | `TestSideCheckpointPersistsBeforeAdmissionAndRestart` |
| 캐시한 설정이 삭제되거나 모델이 바뀌면 계속 실행 거부 | 통과 | `TestSideRejectsDeletedOrChangedCachedProfile` |
| 아카이브 전 취소하고 최종 답변 및 사용량 저장 대기 | 통과 | `TestSideTaskDrainPersistsBeforeArchive` |
| 스트리밍 소비자가 일찍 취소하면 사용량 및 보조 질문 귀속을 한 번만 기록 | 통과 | `TestSideUsageRecordedOnceOnConsumerCancellation` |
| 재시작 시 자동 복구한 Worker / deadline 실행 컨텍스트가 새 스냅샷을 계속 공개 | 통과 | `TestSideRestoredWorkerRuntimePublishesNewCheckpoint` |
| 관련 패키지 race 검사 | 통과 | 아래 명령 |
| TypeScript 및 운영 빌드 | 통과 | `npx tsc --noEmit`·`npm run build` |
| 새 프런트엔드 모듈 Biome | 통과 | `biome check`, 새 모듈 3개 |

별도의 폐기 가능한 데이터베이스에 `ARTEX_PG_DSN`을 설정한 뒤 자동화 검사를 재현할 수 있습니다(운영 DB를 지정하지 마세요):

```sh
go test -race ./agent ./db ./server ./sidequestion ./llmrec ./llmpool \
  -run 'Test(Side|Checkpoint|Snapshot|BuildRequest|Service|MainSide|CaptureRun|TaskArchive|CompleteForwards|StopIntent|CancelIntent)' -count=1
cd web
npx tsc --noEmit
npx biome check src/lib/side-questions.ts src/hooks/use-side-questions.ts src/components/side-question-workspace.tsx
npm run build
```

전체 Go 회귀 검사는 모두 통과한 것이 아닙니다: `server` 패키지의 기존 테스트 두 개가 임시 디렉터리 정리 단계에서 실패했으며, 모두 `TempDir RemoveAll … directory not empty`를 보고했습니다:

- `TestInheritedActivityDetailAndRelationDeletion`
- `TestTaskMetadataPatchReturnsRenameAndPin`

위의 수정하지 않은 기준에서 소스를 내보낸 뒤 같은 격리 환경에서 `server` 패키지를 다시 실행해 정리 실패 두 건을 재현했습니다. 기준 실행에서는 `TestCoreTaskLifecyclePG`의 대상 노드 개수 검사도 실패했으나, 최종 수정 후 `server` 회귀 검사에서는 이 검사가 실패하지 않았습니다. 다른 패키지는 통과했으며 이번 보조 질문 관련 사례와 race 검사도 통과했습니다. 기준의 문제를 이번 검증의 통과 항목으로 처리하지 않았고 문제를 숨기기 위해 기존 검증문을 바꾸지도 않았습니다.

Next.js 빌드 출력에는 기존 다중 lockfile / workspace root 추론 경고가 있으며, 빌드는 완료했고 모든 페이지가 성공적으로 생성되었습니다.

## 브라우저 검사

Codex In-app Browser로 독립 로컬 Go 서비스와 Next.js 개발 서버에 연결했습니다. 데스크톱과 390 × 844 좁은 화면에서 다음 수동 자동화 동작을 완료하고 스크린샷과 브라우저 로그를 검사했습니다:

- 일반 채팅 실행 중 `/btw`를 입력하면 메인 내용과 보조 질문이 동시에 표시됩니다. 데스크톱 사이드바는 정상입니다.
- 연속 후속 질문. 보조 질문 중지 후 이미 생성한 부분은 유지되며 메인 흐름은 계속됩니다.
- 패널을 닫아도 요청은 계속되며 다시 열면 완료된 답변이 복구됩니다. 페이지 새로고침 후 질문이 없는 `/btw`로 이력을 복구합니다.
- 좁은 화면의 Drawer 입력·버튼·이력·닫기 동작이 정상이며 가로 넘침은 없습니다.
- 비우기는 확인 팝업을 사용합니다. 비운 뒤 이력은 사라지고 메인 transcript와 스냅샷은 유지됩니다.
- 작업 MainAgent 및 두 Worker에 각각 질문하고 전환했습니다. Agent 라벨과 이력이 섞이지 않았습니다.
- 차단하는 로컬 모델 fixture로 Worker 실행을 유지했습니다. Worker 메인 입력란에서 `/btw`를 제출하고 보조 질문을 중지한 뒤에도 Worker는 실시간 실행과 자신의 일시 중지 버튼을 표시하며 보조 질문은 일부 답변을 저장합니다.
- 브라우저 오류 / 경고 로그는 비어 있습니다.

제어 가능한 fixture로 동시 실행 시점을 정확히 검증하며 실제 모델의 출력 속도에 의존하지 않습니다. 디버깅 중 두 번의 Worker 실행 상태 검사에서 유효한 동시 실행 구간이 형성되지 않았습니다(작업 종료 / 답변 조기 종료). fixture를 수정한 뒤 다시 수행해 통과했으며 초기 동작은 유효한 통과로 기록하지 않습니다.

## 실제 모델 대화

`grok-4.6`을 먼저 확인했으며 OpenAI 호환 API는 `http://127.0.0.1:12580/tingly/openai`입니다. 확인 요청은 HTTP 200으로 모델명 `grok-4.6` 및 `READY`를 반환했고 2.82초가 걸렸습니다. 첫 선택을 사용할 수 있어 Tingly `glm`이나 智谱 `glm-5.3` 대체 체인을 활성화하지 않았으며, 이 두 대체 서비스는 이번에 검증하지 않았습니다.

| 상황 | 실제 결과 |
| --- | --- |
| 메인 세션 실행 중 자산·목표·표시 질문 | `redhaze.top`, 첫 페이지 읽기와 요약 목표, `BTW-REAL-0910` 반환. 보조 질문 완료, 16.97초 |
| 메인 세션이 첫 페이지를 읽은 뒤 도구 근거 질문 | WebFetch 200, curl 리디렉션 301 → 302 → 200, 페이지 제목을 정확히 인용. 7.24초 |
| 보조 질문에 Bash로 테스트 파일 생성 요청 | 실행 거부, 대상 파일 생성되지 않음. 7.74초 |
| 완료된 보조 질문이 메인 컨텍스트를 바꾸지 않음 | 메인 transcript SHA-256 및 메인 활동 기록이 일치. 보조 질문 도구 실행 횟수 0 |
| Go 서비스를 실제로 중지 / 재시작한 뒤 후속 질문 | 앞선 보조 질문 이력 3개 유지. 영구 저장한 스냅샷에서 바로 자산·표시·제목 답변, 메인 Agent 재실행 없음 |
| 새 세션에 Grok 비스트리밍 설정 사용 | 자산 및 `ATOMIC-0910`에 정확히 답변. 사용량 반환 및 저장:input 11734, output 138, cache_read 11520 |

자산 사례의 메인 세션은 WebFetch 및 Bash/curl로 공개 첫 페이지를 읽었으며 도착 페이지는 `https://id.redhaze.top/home`, 제목은 “红幕科技 RedHaze Group · 全球综合集团门户”입니다. Bash는 응답을 로컬 테스트 파일에 임시 저장했고 원격에는 쓰기를 수행하지 않았습니다. 이 사실은 “보조 질문은 도구를 실행하지 않음”과 별도로 확인했습니다.

메인 transcript 검증값: `e7e61f135a4a120954b539f357e8c4205d7d5cd7460dcaf3dc0fd066463e1d00`.

**사용량 제한:** Tingly의 Grok 스트리밍 응답은 usage를 반환하지 않았습니다. 별도로 `stream_options.include_usage=true`를 직접 전송해 검증했으며 HTTP 200, 데이터 프레임 12개, usage 프레임 0개였습니다. 따라서 스트리밍 테스트의 0은 엔드포인트가 사용량을 제공하지 않았다는 뜻이며 과금이 없었다는 뜻으로 해석할 수 없습니다. 비스트리밍 사용량 및 fixture의 실패 / 취소 사용량은 모두 올바르게 저장되었습니다.

## Qwen 검토

검토 모델은 `qwen-flash`, OpenAI 호환 API는 `https://dashscope.aliyuncs.com/compatible-mode/v1`이며 HTTP 200입니다. 앞선 세 가지 실제 보조 질문 대화·메인 세션 도구 근거·개발 검증문을 제공했습니다. `verdict: accept` 및 `concerns: []`를 반환했으며 답변이 자산·표시·페이지 읽기 증거와 일치하고 보조 질문의 도구 거부가 제약 조건에 맞는다고 판단했습니다. 검토 사용량:prompt 6625, completion 312, total 6937.

이번 Qwen 검토 범위에는 이후 추가한 서비스 재시작 및 비스트리밍 테스트가 포함되지 않습니다. Qwen의 “쓰기 없음” 요약은 범위가 지나치게 넓습니다. 메인 세션의 curl은 실제로 로컬 응답 임시 파일을 만들었으며 위에서 명시적으로 기록했습니다. 동시 실행·도구 실행 없음·transcript 격리는 개발 검증문으로 판단하며 모델 검토는 답변 품질 평가를 보조할 뿐입니다.
