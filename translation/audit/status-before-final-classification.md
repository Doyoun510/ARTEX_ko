# ARTEX A001–A090 현재 로컬 감사 상태

기준 HEAD `d4d288b480fcde527c3302e2c324be6d7816176a`. 브랜치 `work/ko-translation`. 로컬 origin/work/ko-translation 참조 `d4d288b480fcde527c3302e2c324be6d7816176a`. 두 해시는 같지만 fetch하지 않았으므로 실제 원격 최신성을 보장하지 않는다.

Git index staged 변경 없음, MERGE_HEAD 없음. 현재 기존 로컬 변경:
```
 M translation/TRANSLATION_PROMPT.ko.md
?? "monitor=false diff --cached --no-ext-diff --no-textconv --stat"
?? translation/ARTEX_UI_ABC_SCOPE.ko.md
```

기존 TRANSLATION_PROMPT 변경·미추적 담당표와 기타 미추적 파일을 그대로 보존했다. 저장소 파일·index를 수정하지 않았다.

## 처리 상태별 수

| 상태 | A 항목 수 |
| --- | --- |
| 수정 완료 / 정적 대응 확인 | 70 |
| DNT 또는 원문 예시 보존 | 14 |
| 번역·표시 수정 잔여 | 1 |
| 용어·문맥 확인 잔여 | 2 |
| 원문 기능·정책·호환성 문제 | 3 |
| 실행 검증 필요 | 0 |

합계 90개. 실행 검증은 모든 항목 **사용자 요청으로 미실행**이다. 주 분류의 실행 검증 필요 0개는 실행 검증 완료라는 뜻이 아니다. 번역/정적 상태와 실행 상태는 독립적이다. A046은 실행 전 테스트 보완이 필요하여 호환성 문제로 분류했다.

## 기준·이력과 검토 한계

최신 로컬 GLOSSARY/TRANSLATION_PROMPT/GUIDE/PROGRESS, 보완 커밋 e7f74ef·92ad854·451bd13·478544d·5d58222·9262641·d4d288b 및 현재 소스를 대조했다. ID와 이력 연결을 유지한다. A025/A059, A075/A078/A089 같은 중복은 linked 열로 연결하고 ID를 삭제하지 않았다.

환경이 이어진 뒤 /tmp의 이전 감사 원본·중간 보고서·스냅샷이 모두 없어졌다. 앞서 읽은 integrated A001–A090 내용과 지정 위치, 현재 Git/PROGRESS 이력으로 상태표를 재구성하고 현재 줄을 다시 읽었다. 기존 감사의 모든 설명 열을 바이트 단위 복제한 산출물은 아니며, 없어진 파일을 이번에 삭제하거나 덮어쓰지 않았다. 신규 TSV에는 현재 처리 상태와 원문/용어·계약 근거 및 관련 이력을 기록했다.

추적 파일 603개: UTF-8 585개 전체를 NUL 구분 목록/바이트 읽기로 검색(기존 NUL도 유지), 비UTF-8 18개는 바이너리/미디어로 문구 해석 제외. 중국어/지정 전각 부호 2178행 / 124파일. 검색 행 수는 누락 건수가 아니다. U12는 제외하지 않았다.

README.md 원문·CHANGELOG 및 docs/漏洞流量证据.md는 검색했으나 GUIDE의 원문 유지/범위 보류에 따라 번역 누락으로 집계하지 않는다. 나머지 검색은 표시/모델 안내/주석/계약/fixture/예시/고유명으로 문맥 검토했다. 모든 한국어 문장을 Git 원문과 전수 의미 검증한 것은 아니다. 검색 후보가 없는 파일도 완역/동작 검증 완료로 단정하지 않는다.

## 주요 잔여와 수정 후보

### A046 — 원문 기능·정책·호환성 문제

위치 `server/intercept_live_test.go:126,132`. 중국어 생성/쓰기 동사 검사 2곳 유지. read의 과거 생성 차용 금지 및 report write의 현재 쓰기 식별 목적을 한국어에서도 유지하는 보완 필요. 검사를 없애거나 decision만 검사하면 검증 약화.

현재 원문/코드: 126: 				for _, verb := range []string{"创建", "新建", "写入"} { || 132: 			if tc.name == "report_content_is_not_executed" && !strings.Contains(operation, "写") && !strings.Contains(operation, "新建") && !strings.Contains(operation, "创建") {

후속 후보: 기존 후속 검토 유지; 번역과 별도 설계/호환성 검토. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A047 — 원문 기능·정책·호환성 문제

위치 `intercept/prompt.go:20,21,22,23,24,25,26,27,28`. 현재 한국어 상수 전체를 Contains로만 검사. 정확한 옛 중국어 공통 블록이 들어 있으면 한국어 블록을 추가할 수 있다. 사용자 편집 블록을 정확한 옛 기본 블록과 분리 검토해야 함. 커스텀 정책 자동 번역 금지.

현재 원문/코드: 20: func EffectiveJudgePrompt(prompt string) string { || 21: 	if !strings.Contains(prompt, JudgeContextBoundary) { || 22: 		prompt += "\n\n" + JudgeContextBoundary || 23: 	} || 24: 	if !strings.Contains(prompt, JudgeOutputContract) { || 25: 		prompt += "\n\n" + JudgeOutputContract || 26: 	} || 27: 	return prompt || 28: }

후속 후보: 기존 후속 검토 유지; 번역과 별도 설계/호환성 검토. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A048 — 원문 기능·정책·호환성 문제

위치 `skills/api-recon/SKILL.md:21,22,41,151`. SKILL의 로그인/자격 증명/실제 세션 의존 금지와 reference G 실제 세션 안내 및 I6 비교 안내가 공존. 기존 세션 재사용과 로그인 제출은 별개. 원문 정책 충돌을 번역 예외로 확정할 수 없음.

현재 원문/코드: 21: | **인증** | Hook + stub/mock으로 **클라이언트** 로그인 관문 우회 | 사용자에게 계정/비밀번호 요구 또는 추측·실제 로그인 폼 제출 시도 | || 22: | **런타임** | 자격 증명 없이 hook으로 엔드포인트를 가로채고 mock 응답으로 SPA의 로그인 후 기본 화면에 진입 | 실제 백엔드 세션에 의존해야 계속할 수 있는 흐름 | || 41: | 실제 사용자 이름/비밀번호·OTP·OAuth 등의 인증 | stub/mock(위 내용 참조) | || 151: python3 recon/spider_mpa.py <BASE_URL> <OUTDIR> [--cookie "session=..."] [--max 300] [--depth 5] [--exclude "logout|delete|destroy"]

후속 후보: 기존 후속 검토 유지; 번역과 별도 설계/호환성 검토. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A049 — DNT 또는 원문 예시 보존

위치 `skills/api-recon/scripts/preload.js:192`. NEGATIVE_RE는 응답 비교 정규식이므로 DNT. 모드별 stub/forward/L2 동작 차이는 별도 원문 기능 후속 항목.

현재 원문/코드: 192:   const NEGATIVE_RE = /未登录|未授权|授权|not\s*login|unauthorized|forbidden/i;

후속 후보: 계약/원문 예시 유지. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A050 — 번역·표시 수정 잔여

위치 `web/src/lib/mock/data.ts:3067,3080,3093,3106,3136,3155,3156,3163,3168,3172,3218,3229,3239,3513,3515,3544,4165,4166,4196`. Agent 이름 삽입·template_text·prompt·wrapup·skill 선택 조건·LLM 모델용 예시의 후속 번역 검토 잔여. 감사 user_message/context와 실제 도구 입출력 예시는 번역 누락으로 일괄 집계하지 않음.

현재 원문/코드: 3067:     name: "漏洞复测", || 3080:     name: "目标拆解器", || 3093:     name: "规划者", || 3106:     name: "主 Agent", || 3136:     name: "侦察机器人（自定义）", || 3155:   { version: 3, ts: T("2026-07-25T10:00:00Z"), note: "권한 우회 탐지 안내 강화", template_text: "你是 ARTEX 的规划者……" }, || 3156:   { version: 2, ts: T("2026-07-20T10:00:00Z"), note: "초기 버전 소폭 조정", template_text: "你是 ARTEX 的规划者(v2)……" }, || 3163:     prompt: `你是 ARTEX 的「${a.name}」。\n目标：{{.Goal}}\n资产概览：{{.AssetSummary}}\n路线提示：{{.RouteHint}}\n请基于以上信息推进探索，并通过工具把结果写回图。`, || 3168:     wrapup_default: "时间/步数将尽，请总结已确认发现并标记意图终态。", || 3172:     task_timeout_wrapup_default: "任务超时，请立即收尾并落库当前结论。", || 3218:     description: "对 REST/GraphQL API 做侦察与越权面枚举；发现新 API 端点时使用。", || 3229:     description: "用 Playwright 驱动浏览器做动态爬取与截图；需要渲染 JS 站点时使用。", || 3239:     description: "从 ScopeSentry 拉取资产并归并到公司范围；批量导入资产时使用。", || 3513:     user_message: "请整理已经完成的检查，将结论和证据索引写入 reports/summary.md。\n保留待验证项，不修改业务数据。", || 3515:       { kind: "user", text: "汇总本轮已有证据，生成检查摘要。" }, || 3544:     user_message: "检查备份表的用途，涉及删除时需要先确认。", || 4165:         system: "你是一个授权渗透测试系统的「执行者」…（省略）", || 4166:         messages: [{ role: "user", content: "开始执行 system 提示里的这条意图：只做它、只产生事实、做完即停。" }], || 4196:           description: "在目标环境中执行一条 shell 命令并返回其输出。",

후속 후보: 모델용 mock 문자열을 감사 입력·원본 예시와 분리하여 번역 검토. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A052 — 용어·문맥 확인 잔여

위치 `docker-compose.bench.yml:1`. 腾讯 TSec 제품명은 원문 유지. 로컬 자료만으로 공식 영문/한국어 표기를 확정하지 않음. 후보 Tencent TSec; 외부 명칭 확인 별도 필요.

현재 원문/코드: 1: # ARTEX 腾讯 TSec Benchmark 평가 이미지 — 로컬 원클릭 build + run.

후속 후보: Tencent TSec 표기 확인 후 주석 정리. 연결: 해당 파일의 지정 함수/검사 및 PROGRESS 보완 기록; 현재 근거는 current 열 참조.

### A090 — 용어·문맥 확인 잔여

위치 `db/db.go:192`. 原文 记录代理地址(驱动 if 双文案). agent/prompt_test.go:45 dual-text 및 :51 if/else fixture가 사용자 템플릿 문구 분기 해석을 뒷받침. 직접 정의 없음. ProxyAddr는 캡처 OFF에서도 전역 프록시 반환. 설명은 미수정.

현재 원문/코드: 192: 		{"ProxyAddr", "기록 프록시 주소(if 이중 문구 구동)", "127.0.0.1:8080", "runtime"},

후속 후보: 프록시 주소(사용자 프롬프트의 if/else 안내 문구 선택에 사용). 연결: agent/prompt.go:19,38–64; worker.go:289; prompt_test.go:45–81; server/manager.go:740–763; config.go:713–725; agent-editor.tsx:377.

## 추가 검색 후보 (90개 집계와 별도)

| ID | 위치·현재 문구 | 판단·후보·소비 |
| --- | --- |
| S001 | server/notify_api_test.go:20–33,110,120,130,241,351,353,354,496 및 다수 | 중국어 주석/개발자 실패 안내/사례 설명 잔여. `记录漏洞失败: %v`→`취약점 기록 실패: %v`, `缺名称`→`이름 없음`. A020의 지정 12곳과 U10 운영 기대값은 완료; 같은 파일 전체 완료 아님. t.Fatal/t.Run 소비이며 변경 전 -run 참조 확인 필요. 216 검색행에는 fixture도 포함하므로 216개 누락이라고 집계하지 않는다. |
| S002 | server/engine_emptyturn_test.go:75,86,100,114,121,132,140 | 중국어 t.Run 설명 7곳. `空转回合注入续跑指令`→`무진행 턴에 실행 재개 지시 추가`. A018 주석 보완과 별개. 이름 -run 참조는 변경 전 확인 필요. |
| S003 | agent/promptcatalog.go:31,37,42,45 | `공격면`→확정 `공격 표면`. 0ef8561 원문 攻击面, 용어집:207과 같은 침투 문맥. 모델 본문/renderSystem/seed 소비. A084는 DB 안내 완료이며 이 본문까지 정리하지 않았음. |
| S004 | system/intercept/page.tsx:542–543; agent-editor.tsx:1140; agent/planner.go:185,191 | 표시/모델 목록의 join("、"), strings.Join(...,"；") 후보. 현재 명시 허용 범위와 소비 계약 추가 확인 필요. A080 허용을 전역으로 확대하지 않음. |
| S005 | server/conversations.go:561 | 영어 주석의 인용 `已手动停止`. 현재 생성부/SDK 실제 반환과 대조할 오래된 인용 후보; 외부 출력일 수 있어 미확인. 영어 주석 전체 재작성 아님. |
| S006 | README.ko.md:272 | 고정 라벨 뒤 `수정 완료으로`·`수정 완료을`→`수정 완료로`·`수정 완료를`. 문서 조사 정정 후보. 상태 enum/로직 불변. |

추가 확정 번역/표기 잔여는 S001/S002/S003/S006의 4개 묶음, S004/S005는 확인 후보 2개다. A-ID를 새 번호로 바꾸거나 known DNT를 중복 누락 집계하지 않는다.

## 계약 및 입력·예시 보존 확인

- [模型]: intercept 생성 → DB LIKE/HasPrefix → approval-records startsWith/replace. 불일치 발견 없음.
- 归档不存在: 서버 오류/프런트 tasks includes/mock 오류; 已存在: skill 업로드 오류/프런트 overwrite includes. 양쪽 유지.
- 实际操作：·；成功后的后果：·；命中规则：: 프롬프트/예시/오류 안내/파서 정적 일치. 120자 모델 지시와 2400바이트 파서 제한 별개.
- A051 한국어 displayLabel과 중국어 label 삽입/파싱 분리. selectedMentions 테스트 표시 기대값 취약점 확인. 입력창 중국어 계약 토큰은 남는다.
- 테스트 길이/UTF-8/개행/round-trip/필터/마스킹/모델 응답/실제 출력 fixture는 보존. fixture를 그대로 인용한 기대 실패 안내의 중국어는 그 자체가 잘못된 번역이라고 단정하지 않는다.
- docs 경로·ICP 번호/备案·ScopeSentry DSL·작업명·placeholder·WHERE 全表·파일명 critical_SQL注入_#123.md·법률 원문·과거 validation JSON은 원문 예시/계약. 腾讯 TSec는 고유명 확인 잔여 A052.

## 원문 문제 / 번역과 별개

- U7 live 검사 두 곳은 한국어에서도 과거 생성 차용/현재 보고서 쓰기 여부를 동등하게 잡도록 검토해야 한다. 검사 삭제/decision-only로 약화 금지.
- U7 정확한 옛 공통 블록 중복과 편집 커스텀 블록은 별도 식별. 기존 커스텀 정책 본문은 자동 번역하지 않는다. 기본값 복원은 기존 정책을 바꾸는 별도 사용자 선택이며 이번에 실행/DB 확인하지 않았다.
- U15 로그인·세션 정책 충돌은 기존 세션 재사용과 로그인 제출을 구분한 정책 정리 대상. 현재 예외가 이미 허용됐다고 주장하지 않는다.
- U15 depth: runtime_harvest.js:157–177은 명시 stub 우선이며 forward=true 미매칭 API를 실제 fetch/보정, false는 기본 본문. coverage: preload.js:214–267,314는 forward=true이면 stub 후보가 있어도 전송, false 미매칭도 전송 가능하며 L2 tier/응답 조건이 다르다. forward=false가 전체 오프라인이라는 보장 없음. 문서화 또는 동작 정렬은 별도 기능 결정이며 구현/정책을 수정하지 않았다.
- notify/mask_test.go:156–160의 TLS 끄기 설명과 tls:false→true 입력 방향 차이는 원문 설명 문제. fixture/검사를 번역 과정에서 정정하지 않음.
- U15 D→E 참조/routes 설명 정정, ScopeSentry name=scopesentry 및 frontmatter 종료선 정정은 이미 반영. 미해결 형식 문제로 재집계하지 않는다. 실제 SDK 로드/모델 동작은 미검증.

## 다음 묶음 제안

1. S001/S002 주석·실패 안내·테스트 이름만 번역(선택 참조 확인 및 fixture 보존).
2. S003/S006 용어/조사 정리. A052/A090은 근거 확인 후 별도 최소 정정.
3. S004/S005 구분자 소비/명시 허용 및 오래된 인용 확인.
4. A050 모델용 mock 안내를 감사 입력/실제 입출력 원문 예시와 분리한 번역 범위로 검토.
5. U7/U15 정책·호환성은 번역과 분리한 설계·테스트 검토로 팀에 전달. 이번 작업으로 해결 처리하지 않음.

## 항목별 처리 상태

| ID | 상태 | 파일·위치 | 근거 |
| --- | --- | --- | --- |
| A001 | 수정 완료 / 정적 대응 확인 | agent/finding_recorder_test.go:102 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A002 | 수정 완료 / 정적 대응 확인 | server/finding_traffic_test.go:355 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A003 | 수정 완료 / 정적 대응 확인 | server/finding_workflow_test.go:134 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A004 | 수정 완료 / 정적 대응 확인 | server/finding_workflow_test.go:249 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A005 | 수정 완료 / 정적 대응 확인 | server/intercept_detail_test.go:96 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A006 | 수정 완료 / 정적 대응 확인 | server/task_categories_test.go:48 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A007 | 수정 완료 / 정적 대응 확인 | server/chat_mentions_test.go:180 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A008 | 수정 완료 / 정적 대응 확인 | server/chat_mentions_test.go:307 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A009 | 수정 완료 / 정적 대응 확인 | server/skill_upload_test.go:192 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A010 | 수정 완료 / 정적 대응 확인 | server/skill_upload_test.go:213 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A011 | 수정 완료 / 정적 대응 확인 | web/next.config.mjs:3,4,6,11,14,15,16,23,30,34 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A012 | 수정 완료 / 정적 대응 확인 | agent/capture_usage_test.go:36,62 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A013 | 수정 완료 / 정적 대응 확인 | agent/coldgraph_test.go:17,40 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A014 | 수정 완료 / 정적 대응 확인 | agent/insert_assets_test.go:45,114,162,245,306,364,401,402 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A015 | 수정 완료 / 정적 대응 확인 | server/llmretry_test.go:11,12 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A016 | 수정 완료 / 정적 대응 확인 | server/task_archive_package_test.go:107 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A017 | 수정 완료 / 정적 대응 확인 | server/trigger_merge_test.go:12,95 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A018 | 수정 완료 / 정적 대응 확인 | server/engine_emptyturn_test.go:11,40,57,79,93,96,101,105,109,113,117,125,128,136,141,143 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A019 | 수정 완료 / 정적 대응 확인 | server/update_test.go:13,14,15,38,52,57,71,77,89,91,93,96,99,101,105,108,121,122,126,131,134,135,137,140,143,144,150,153 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A020 | 수정 완료 / 정적 대응 확인 | server/notify_api_test.go:188,196,201,202,216,220,246,249,431,438,720,780 | 지정 12곳은 완료. 같은 파일의 다른 중국어 주석/실패 안내/사례 이름은 별도 잔여 S001. |
| A021 | 수정 완료 / 정적 대응 확인 | server/skill_upload_test.go:57,76,84,85,131,154 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A022 | 수정 완료 / 정적 대응 확인 | sidequestion/README.md:1,3,5,7,9,36 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A023 | 수정 완료 / 정적 대응 확인 | sidequestion/CONTEXT_BUDGET.md:1,3,5,15,17 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A024 | 수정 완료 / 정적 대응 확인 | sidequestion/VALIDATION.md:1,3,5,7,9,51,64 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A025 | 수정 완료 / 정적 대응 확인 | web/src/app/(main)/chat/page.tsx:881,901,933,967,995 | displayTitle은 빈 제목 또는 정확히 新对话만 새 대화로 표시. 제목/aria/삭제 확인은 표시값, 이름 변경 입력/서버/API/DB는 원 title. 사용자가 동일한 新对话를 입력하면 구분 불가. |
| A026 | 수정 완료 / 정적 대응 확인 | web/src/lib/mock/handler.ts:1231,1250,1310,1336,1804 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A027 | 수정 완료 / 정적 대응 확인 | server/llmpool.go:13,14,86 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A028 | 수정 완료 / 정적 대응 확인 | server/server_mgmt.go:1941,1943 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A029 | 수정 완료 / 정적 대응 확인 | server/llmpool.go:17,34,51 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A030 | 수정 완료 / 정적 대응 확인 | server/llmretry.go:13,16,52 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A031 | 수정 완료 / 정적 대응 확인 | server/server.go:155 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A032 | 수정 완료 / 정적 대응 확인 | server/task_llm.go:94,249,299 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A033 | 수정 완료 / 정적 대응 확인 | server/server_mgmt.go:1848,1849 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A034 | 수정 완료 / 정적 대응 확인 | server/dto.go:616 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A035 | 수정 완료 / 정적 대응 확인 | server/llmretry.go:11 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A036 | 수정 완료 / 정적 대응 확인 | agent/planner.go:304,311,333 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A037 | 수정 완료 / 정적 대응 확인 | agent/compaction.go:538 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A038 | 수정 완료 / 정적 대응 확인 | web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx:243,274,417 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A039 | 수정 완료 / 정적 대응 확인 | agent/retester.go:6 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A040 | 수정 완료 / 정적 대응 확인 | agent/finding_workflow.go:77 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A041 | 수정 완료 / 정적 대응 확인 | sidequestion/service.go:61 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A042 | 수정 완료 / 정적 대응 확인 | agent/promptcatalog.go:95 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A043 | 수정 완료 / 정적 대응 확인 | README.ko.md:307 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A044 | 수정 완료 / 정적 대응 확인 | web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx:2121,2122,2123 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A045 | 수정 완료 / 정적 대응 확인 | README.ko.md:61,126,182,186,297,300,305,307,354,359,371,425,426,427,428,442,451,452,479 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A046 | 원문 기능·정책·호환성 문제 | server/intercept_live_test.go:126,132 | 중국어 생성/쓰기 동사 검사 2곳 유지. read의 과거 생성 차용 금지 및 report write의 현재 쓰기 식별 목적을 한국어에서도 유지하는 보완 필요. 검사를 없애거나 decision만 검사하면 검증 약화. |
| A047 | 원문 기능·정책·호환성 문제 | intercept/prompt.go:20,21,22,23,24,25,26,27,28 | 현재 한국어 상수 전체를 Contains로만 검사. 정확한 옛 중국어 공통 블록이 들어 있으면 한국어 블록을 추가할 수 있다. 사용자 편집 블록을 정확한 옛 기본 블록과 분리 검토해야 함. 커스텀 정책 자동 번역 금지. |
| A048 | 원문 기능·정책·호환성 문제 | skills/api-recon/SKILL.md:21,22,41,151 | SKILL의 로그인/자격 증명/실제 세션 의존 금지와 reference G 실제 세션 안내 및 I6 비교 안내가 공존. 기존 세션 재사용과 로그인 제출은 별개. 원문 정책 충돌을 번역 예외로 확정할 수 없음. |
| A049 | DNT 또는 원문 예시 보존 | skills/api-recon/scripts/preload.js:192 | NEGATIVE_RE는 응답 비교 정규식이므로 DNT. 모드별 stub/forward/L2 동작 차이는 별도 원문 기능 후속 항목. |
| A050 | 번역·표시 수정 잔여 | web/src/lib/mock/data.ts:3067,3080,3093,3106,3136,3155,3156,3163,3168,3172,3218,3229,3239,3513,3515,3544,4165,4166,4196 | Agent 이름 삽입·template_text·prompt·wrapup·skill 선택 조건·LLM 모델용 예시의 후속 번역 검토 잔여. 감사 user_message/context와 실제 도구 입출력 예시는 번역 누락으로 일괄 집계하지 않음. |
| A051 | 수정 완료 / 정적 대응 확인 | web/src/components/mention-textarea.tsx:148,239,268,303 | displayLabel로 선택/검색/배지 및 selectedMentions 표시를 분리. 실제 삽입 categories[index].label/mentionToken/정규식/서버 파서는 중국어 계약 보존. 표시 기대값은 취약점. |
| A052 | 용어·문맥 확인 잔여 | docker-compose.bench.yml:1 | 腾讯 TSec 제품명은 원문 유지. 로컬 자료만으로 공식 영문/한국어 표기를 확정하지 않음. 후보 Tencent TSec; 외부 명칭 확인 별도 필요. |
| A053 | 수정 완료 / 정적 대응 확인 | server/intercept.go:50 | 안내 문자열은 decision/comment 두 문자열 필드, allow/ask/deny 및 실제 comment 앵커를 명시. ParseVerdict의 현재 계약과 정적 일치. |
| A054 | DNT 또는 원문 예시 보존 | intercept/intercept.go:480,481,482 | [模型] 생성/DB LIKE·HasPrefix/프런트 startsWith·replace 일치. |
| A055 | DNT 또는 원문 예시 보존 | intercept/intercept.go:612 | 工具 %s 请求审批 (#%d)는 transcript.tsx:201 /工具\s+(\S+)\s+请求/ 소비 계약과 연결. |
| A056 | DNT 또는 원문 예시 보존 | server/task_archives.go:461 | 서버 오류의 归档不存在와 tasks/page.tsx:1680 includes 일치; mock 오류 계약도 보존. |
| A057 | DNT 또는 원문 예시 보존 | server/server_mgmt.go:1361 | skill 서버 오류 已存在와 skills/page.tsx:370 includes 일치. |
| A058 | DNT 또는 원문 예시 보존 | intercept/prompt.go:36,37,38,39,124,125,126,127,193,196,200 | 实际操作：, ；成功后的后果：, ；命中规则：가 프롬프트/예시/ParseVerdict에 전각 부호까지 일치. 모델 120자 지시와 2400바이트 제한 별개. |
| A059 | DNT 또는 원문 예시 보존 | server/conversations.go:112,443 | 新对话 저장/자동 제목 비교는 DNT. A025 화면 새 대화 표시와 계약 분리. |
| A060 | DNT 또는 원문 예시 보존 | server/orchestration.go:699,700,701,719 | reporterToolCallMessageV1는 저장된 옛 기본 메시지와 바이트 비교하는 마이그레이션 기준. |
| A061 | DNT 또는 원문 예시 보존 | skills/api-recon/reference.md:52 | 중국어 로그인/권한 응답 검색 정규식 보존. |
| A062 | 수정 완료 / 정적 대응 확인 | server/finding_retests.go:230 | 현재 SQL note는 내장 기본값. 과거 DB 저장 문구 소급 번역/실제 DB 미확인. |
| A063 | DNT 또는 원문 예시 보존 | report/findings.go:125 | critical_SQL注入_#123.md 원문 파일명 예시 유지. |
| A064 | DNT 또는 원문 예시 보존 | agent/terminalreason_test.go:104,109,163,166 | 실제 출력 보존, Unicode 길이/개행 입력 및 대응 기대값 보존. |
| A065 | DNT 또는 원문 예시 보존 | skills/scopesentry/SKILL.md:15,16,32,34,41,93,120,133,137,158,169,188,262,334 | ScopeSentry DSL 검색값·URL·키·placeholder·작업명 예시 유지. name=scopesentry와 description 다음 종료선 정정은 반영됨. |
| A066 | DNT 또는 원문 예시 보존 | intercept/prompt.go:115 | WHERE 全表는 원문 SQL 예시이며 파서 계약이라고 단정하지 않음. |
| A067 | DNT 또는 원문 예시 보존 | agent/tools_insert.go:142,169,380,382,387,437 | ICP 등록(备案) 확정 표기 및 별도 备案 파서 유지. |
| A068 | 수정 완료 / 정적 대응 확인 | notify/email_test.go:206,207,208 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A069 | 수정 완료 / 정적 대응 확인 | notify/filter_test.go:13,14,15 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A070 | 수정 완료 / 정적 대응 확인 | notify/http_test.go:36,37,38 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A071 | 수정 완료 / 정적 대응 확인 | notify/mask_test.go:132,142,149,156 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A072 | 수정 완료 / 정적 대응 확인 | notify/pack_test.go:122,133,142 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A073 | 수정 완료 / 정적 대응 확인 | notify/redact_test.go:39,45,51 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A074 | 수정 완료 / 정적 대응 확인 | translation/PROMPT_GUIDE.ko.md:69 | GUIDE는 작업 #<id> 참조와 taskContextHeader의 작업 #를 인용. 다른 U1 보류 해소로 확대하지 않음. |
| A075 | 수정 완료 / 정적 대응 확인 | db/schema.sql:1,2,5,11,12,13,14,24,151,196,285,343,374,375,379,380,381,384,386,387,389,390,391,392,393,395,396,398,399,413,416,418,419,424,425,428,435,438,439,440,441,442,443,450,456,457,458,459,460,461,476,477,478,481,482,483,489,541,543,545,546,547,554,619,621,626,732,733,734,735,736,737,750,755,764,791,795,797,814,831,842,851,857,858,859,875,897,898,899,910,911,917,918,933,952,991,1015,1031,1063,1070,1079,1086,1087,1094,1095,1096,1098,1108,1111,1131,1177,1232,1233,1234,1235,1236,1237,1238,1239,1260,1261,1262,1263,1264,1265,1282,1289,1291,1292,1293,1295,1296,1297,1300,1301,1302,1306,1309,1310,1318,1320,1322,1323,1324,1326,1327,1336,1337,1338,1339,1340,1347,1348,1355,1356,1357,1358,1363,1364,1365,1366,1369,1370,1371,1374 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A076 | 수정 완료 / 정적 대응 확인 | db/schema.sql:670 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A077 | 수정 완료 / 정적 대응 확인 | db/schema.sql:688 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A078 | 수정 완료 / 정적 대응 확인 | db/schema.sql:1131,1134 | 주석은 A075와 중복; 삭제된 재검증 오류 문자열만 대응 수정. |
| A079 | 수정 완료 / 정적 대응 확인 | db/company_scope.go:155,160 | 설명 주석 ICP 등록(备案), 실제 ICP 예시/备案 파서 보존. |
| A080 | 수정 완료 / 정적 대응 확인 | db/companies.go:636 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A081 | 수정 완료 / 정적 대응 확인 | db/config.go:71 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A082 | 수정 완료 / 정적 대응 확인 | db/llmretry.go:87 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A083 | 수정 완료 / 정적 대응 확인 | db/db.go:195,196 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A084 | 수정 완료 / 정적 대응 확인 | db/db.go:198 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A085 | 수정 완료 / 정적 대응 확인 | db/nkey.go:61 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A086 | 수정 완료 / 정적 대응 확인 | db/notification.go:16,30,34 | 알림 delivery 주석 전송. enum sent와 UI 전달됨은 별도 확정 상태 라벨. |
| A087 | 수정 완료 / 정적 대응 확인 | db/notification_delivery.go:12,15,18 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A088 | 수정 완료 / 정적 대응 확인 | db/notification.go:131,135 | 감사 지정 문구를 현재 줄 및 해당 보완 커밋/PROGRESS 기록과 재대조: 지정 번역·용어·기대값 반영. 파일 전체 완료를 의미하지 않음. |
| A089 | 수정 완료 / 정적 대응 확인 | translation/PROMPT_PROGRESS.ko.md:186,202,208 | 과거 U12/PASS 이력 보존. 이후 A075–A079/H01–H35·알림 표기 보완이 실제 잔여 해소 설명. 현재 실행 검증 근거가 아님. |
| A090 | 용어·문맥 확인 잔여 | db/db.go:192 | 原文 记录代理地址(驱动 if 双文案). agent/prompt_test.go:45 dual-text 및 :51 if/else fixture가 사용자 템플릿 문구 분기 해석을 뒷받침. 직접 정의 없음. ProxyAddr는 캡처 OFF에서도 전역 프록시 반환. 설명은 미수정. |

## 출력 언어·검색 커버리지

DefaultAssistantPrompt(promptcatalog.go:70), ReporterDefaultPrompt(:107), RetesterDefaultPrompt(retester.go:15)의 한국어 출력 지시 일치. 중국어 판정 계약 앵커는 출력 언어 충돌로 취급하지 않는다. A050 mock/저장 커스텀과 실제 모델 응답 언어는 별도 미확인.

검색 확인 UTF-8 파일 목록(문장 의미 전수 검증 목록 아님):
```
.dockerignore
.env.example
.gitattributes
.github/workflows/release.yml
.gitignore
CHANGELOG.md
Dockerfile
LICENSE
README.ko.md
README.md
agent/assembly.go
agent/blackboard_inheritance_test.go
agent/cancelcause.go
agent/capture.go
agent/capture_approval_test.go
agent/capture_usage_test.go
agent/chat.go
agent/coldgraph.go
agent/coldgraph_test.go
agent/compaction.go
agent/constraints.go
agent/deferred.go
agent/deferred_test.go
agent/finding_recorder.go
agent/finding_recorder_test.go
agent/finding_workflow.go
agent/finding_workflow_test.go
agent/goals.go
agent/insert_assets_test.go
agent/mainagent.go
agent/noa.go
agent/planner.go
agent/prompt.go
agent/prompt_now_test.go
agent/prompt_test.go
agent/promptcatalog.go
agent/provider.go
agent/provider_capture_e2e_test.go
agent/provider_capture_test.go
agent/provider_quota_test.go
agent/provider_responses_test.go
agent/proxyenv_test.go
agent/retester.go
agent/review_context_test.go
agent/runinfo.go
agent/session_header_test.go
agent/side_questions.go
agent/side_questions_test.go
agent/taskclock.go
agent/terminalreason.go
agent/terminalreason_test.go
agent/testmain_test.go
agent/toolcatalog.go
agent/toolcatalog_test.go
agent/tools.go
agent/tools_digest.go
agent/tools_insert.go
agent/tools_nil_store_test.go
agent/tools_overview_test.go
agent/worker.go
agent/worker_intervention_test.go
agent/wrapup.go
build.sh
cmd/artex/main.go
cmd/artex/main_test.go
config.example.json
config/config.go
config/config_test.go
db/activity_page_test.go
db/asset_dsl.go
db/asset_intercept.go
db/asset_intercept_match.go
db/asset_intercept_match_test.go
db/assets.go
db/assets_test.go
db/chat_mentions.go
db/commands.go
db/companies.go
db/companies_test.go
db/company_lock_test.go
db/company_scope.go
db/company_scope_consistency_test.go
db/company_scope_test.go
db/config.go
db/config_max_tokens_test.go
db/config_test.go
db/constants.go
db/constraints.go
db/conversation.go
db/conversation_pin_test.go
db/customtool_test.go
db/db.go
db/db_test.go
db/digest.go
db/exploration.go
db/exploration_sources.go
db/exploration_sources_test.go
db/exploration_test.go
db/exploration_tokens_test.go
db/finding_assets.go
db/finding_assets_test.go
db/finding_retests.go
db/finding_retests_test.go
db/finding_traffic.go
db/finding_traffic_archive.go
db/finding_traffic_archive_test.go
db/findings.go
db/findings_test.go
db/intent_admission.go
db/intent_control_test.go
db/intent_delete_test.go
db/intercept.go
db/intercept_detail.go
db/intercept_detail_test.go
db/intercept_execution.go
db/intercept_execution_test.go
db/intercept_filter_test.go
db/intercept_seed_test.go
db/jsonb_clean_test.go
db/llm_records_migrate_test.go
db/llm_usage.go
db/llmhealth.go
db/llmretry.go
db/llmretry_test.go
db/logs.go
db/nkey.go
db/notification.go
db/notification_delivery.go
db/notification_test.go
db/schema.sql
db/settings.go
db/side_questions.go
db/side_questions_test.go
db/skill_usage.go
db/skill_usage_test.go
db/task_archive_aggregate_stats.go
db/task_archives.go
db/task_archives_restore.go
db/task_archives_test.go
db/task_assets.go
db/task_assets_context.go
db/task_assets_test.go
db/task_categories.go
db/task_categories_test.go
db/task_context.go
db/task_context_lock_test.go
db/task_context_unit_test.go
db/task_delete_concurrency_test.go
db/task_intercept.go
db/task_list_performance_test.go
db/task_metadata_test.go
db/task_queue_test.go
db/task_scope.go
db/task_scope_test.go
db/task_source_limit_test.go
db/task_templates.go
db/task_templates_test.go
db/tasks.go
db/tasks_test.go
db/testmain_test.go
db/tool_usage.go
db/tool_usage_test.go
db/tools.go
db/triggers.go
dev.sh
docker-compose.bench.yml
docker-compose.yml
docs/漏洞流量证据.md
enrich/enrich.go
evidence/store.go
evidence/store_test.go
evidence/testmain_test.go
go.mod
go.sum
guard/guard.go
guard/guard_test.go
install.sh
intercept/intercept.go
intercept/prompt.go
intercept/prompt_test.go
intercept/review_context.go
intercept/review_context_test.go
intercept/trace.go
intercept/trace_test.go
llmpool/health.go
llmpool/health_policy_test.go
llmpool/pool.go
llmpool/pool_test.go
llmrec/capture.go
llmrec/capture_test.go
llmrec/llmrec.go
llmrec/llmrec_test.go
mcphttp/client.go
mcphttp/client_sse_test.go
notify/channel.go
notify/channels_test.go
notify/dingtalk.go
notify/email.go
notify/email_test.go
notify/event.go
notify/feishu.go
notify/filter.go
notify/filter_test.go
notify/html.go
notify/http.go
notify/http_test.go
notify/markdown.go
notify/mask.go
notify/mask_test.go
notify/notify.go
notify/pack_test.go
notify/redact_test.go
notify/render.go
notify/render_test.go
notify/sign_test.go
notify/ssrf_test.go
notify/telegram.go
notify/webhook.go
notify/webhook_template_test.go
notify/wecom.go
report/findings.go
report/report.go
reset-password.sh
screenshots/.gitkeep
selfupdate/bootstrap.go
selfupdate/github.go
selfupdate/selfupdate.go
selfupdate/selfupdate_test.go
selfupdate/stage.go
server/assembly.go
server/assembly_test.go
server/asset_intercept.go
server/assets.go
server/assets_scope_test.go
server/auth.go
server/broadcast.go
server/chat_llm_resolve_test.go
server/chat_mentions.go
server/chat_mentions_test.go
server/chatupload.go
server/commands.go
server/constraints_api.go
server/conversation_status_test.go
server/conversations.go
server/core_test.go
server/customtool.go
server/customtool_test.go
server/dto.go
server/engine.go
server/engine_cancelcause_test.go
server/engine_emptyturn_test.go
server/engine_llm_calls_test.go
server/engine_timeout.go
server/finding_retests.go
server/finding_retests_test.go
server/finding_traffic.go
server/finding_traffic_export.go
server/finding_traffic_test.go
server/finding_workflow.go
server/finding_workflow_test.go
server/findings_groups.go
server/findings_groups_test.go
server/goals.go
server/goals_api.go
server/goals_test.go
server/inheritance_api_test.go
server/inheritance_dto_test.go
server/intent_intervention.go
server/intercept.go
server/intercept_detail_test.go
server/intercept_filter_test.go
server/intercept_live_test.go
server/intercept_review_test.go
server/llmpool.go
server/llmpool_test.go
server/llmrec_raw_test.go
server/llmretry.go
server/llmretry_test.go
server/logsink.go
server/manager.go
server/manager_delete_files_test.go
server/manager_lifecycle_test.go
server/mcpdiscover.go
server/mgmt_test.go
server/notifier.go
server/notify_api.go
server/notify_api_test.go
server/orchestration.go
server/platform_tools.go
server/platform_tools_test.go
server/prompt_vars_test.go
server/scheduler.go
server/server.go
server/server_mgmt.go
server/side_questions.go
server/side_questions_test.go
server/skill_upload_test.go
server/skill_usage.go
server/skill_zip.go
server/sync_scopesentry.go
server/task_admission_test.go
server/task_archive_package.go
server/task_archive_package_test.go
server/task_archives.go
server/task_archives_test.go
server/task_assets.go
server/task_categories.go
server/task_categories_test.go
server/task_control.go
server/task_control_routes_test.go
server/task_delete_barrier_test.go
server/task_intercept.go
server/task_llm.go
server/task_llm_test.go
server/task_metadata.go
server/task_metadata_test.go
server/task_resolution.go
server/task_templates.go
server/task_templates_test.go
server/testmain_test.go
server/tool_usage.go
server/tool_usage_test.go
server/tools_wire_test.go
server/trigger_merge_test.go
server/triggers.go
server/update.go
server/update_test.go
server/webui_embed.go
server/webui_stub.go
server/worker_message_test.go
server/workspace.go
sidequestion/CONTEXT_BUDGET.md
sidequestion/README.md
sidequestion/VALIDATION.md
sidequestion/capture.go
sidequestion/context.go
sidequestion/context_test.go
sidequestion/request.go
sidequestion/service.go
sidequestion/sidequestion_test.go
sidequestion/validation-2026-09-10.json
skills/api-recon/SKILL.md
skills/api-recon/reference.md
skills/api-recon/scripts/build_perm_tree.py
skills/api-recon/scripts/extract_route_map.py
skills/api-recon/scripts/harvest_static.py
skills/api-recon/scripts/package-lock.json
skills/api-recon/scripts/package.json
skills/api-recon/scripts/preload.js
skills/api-recon/scripts/runtime_harvest.js
skills/api-recon/scripts/spider_mpa.py
skills/playwright-cli/SKILL.md
skills/playwright-cli/references/element-attributes.md
skills/playwright-cli/references/playwright-tests.md
skills/playwright-cli/references/request-mocking.md
skills/playwright-cli/references/running-code.md
skills/playwright-cli/references/session-management.md
skills/playwright-cli/references/storage-state.md
skills/playwright-cli/references/test-generation.md
skills/playwright-cli/references/tracing.md
skills/playwright-cli/references/video-recording.md
skills/scopesentry/SKILL.md
start.bat
start.sh
traffic/archive.go
traffic/archive_test.go
traffic/evidence.go
traffic/passthrough_test.go
traffic/proxy_test.go
traffic/reclaim_test.go
traffic/traffic.go
traffic/traffic_store_test.go
traffic/traffic_test.go
traffic/upgrade_test.go
translation/GLOSSARY.ko.md
translation/PROMPT_GUIDE.ko.md
translation/PROMPT_PROGRESS.ko.md
translation/TRANSLATION_PROMPT.ko.md
update.sh
web/.gitignore
web/.husky/pre-commit
web/LICENSE
web/biome.json
web/components.json
web/next.config.mjs
web/package-lock.json
web/package.json
web/postcss.config.mjs
web/public/logo.svg
web/src/app/(auth)/layout.tsx
web/src/app/(auth)/login/page.tsx
web/src/app/(auth)/setup/page.tsx
web/src/app/(external)/page.tsx
web/src/app/(main)/_components/main-content.tsx
web/src/app/(main)/_components/sidebar/account-switcher.tsx
web/src/app/(main)/_components/sidebar/app-sidebar.tsx
web/src/app/(main)/_components/sidebar/change-password-dialog.tsx
web/src/app/(main)/_components/sidebar/layout-controls.tsx
web/src/app/(main)/_components/sidebar/nav-documents.tsx
web/src/app/(main)/_components/sidebar/nav-main.tsx
web/src/app/(main)/_components/sidebar/nav-secondary.tsx
web/src/app/(main)/_components/sidebar/nav-user.tsx
web/src/app/(main)/_components/sidebar/search-dialog.tsx
web/src/app/(main)/_components/sidebar/sidebar-support-card.tsx
web/src/app/(main)/_components/sidebar/theme-switcher.tsx
web/src/app/(main)/_components/update-badge.tsx
web/src/app/(main)/chat/page.tsx
web/src/app/(main)/dashboard/page.tsx
web/src/app/(main)/function/assets/page.tsx
web/src/app/(main)/function/commands/page.tsx
web/src/app/(main)/function/findings/_components/asset-tree.tsx
web/src/app/(main)/function/findings/_components/findings-table.tsx
web/src/app/(main)/function/findings/detail/lineage.tsx
web/src/app/(main)/function/findings/detail/page.tsx
web/src/app/(main)/function/findings/page.tsx
web/src/app/(main)/function/llm-records/page.tsx
web/src/app/(main)/function/sync/page.tsx
web/src/app/(main)/function/tasks/detail/_tabs/assets-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/coverage-graph-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/findings-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/graph-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/intercept-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/overview-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/report-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/retests-tab.tsx
web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx
web/src/app/(main)/function/tasks/detail/page.tsx
web/src/app/(main)/function/tasks/page.tsx
web/src/app/(main)/function/traffic/page.tsx
web/src/app/(main)/function/workspace/page.tsx
web/src/app/(main)/layout.tsx
web/src/app/(main)/system/agents/detail/page.tsx
web/src/app/(main)/system/agents/page.tsx
web/src/app/(main)/system/intercept/approvals/page.tsx
web/src/app/(main)/system/intercept/assets/page.tsx
web/src/app/(main)/system/intercept/page.tsx
web/src/app/(main)/system/llm/_components/retry.tsx
web/src/app/(main)/system/llm/page.tsx
web/src/app/(main)/system/logs/page.tsx
web/src/app/(main)/system/mcp/page.tsx
web/src/app/(main)/system/notify/_components/channel-fields.ts
web/src/app/(main)/system/notify/_components/channel-form.tsx
web/src/app/(main)/system/notify/_components/delivery-list.tsx
web/src/app/(main)/system/notify/_components/stat-tile.tsx
web/src/app/(main)/system/notify/page.tsx
web/src/app/(main)/system/settings/_components/update-card.tsx
web/src/app/(main)/system/settings/page.tsx
web/src/app/(main)/system/skills/page.tsx
web/src/app/(main)/system/tools/page.tsx
web/src/app/globals.css
web/src/app/layout.tsx
web/src/app/not-found.tsx
web/src/components/agent-editor.tsx
web/src/components/approval-execution-focus.tsx
web/src/components/approval-records.tsx
web/src/components/asset-dsl-search.tsx
web/src/components/asset-intercept-rules-editor.tsx
web/src/components/calendar/event-calendar-views.tsx
web/src/components/copy-button.tsx
web/src/components/date-range-picker.tsx
web/src/components/exploration-graph.tsx
web/src/components/finding-retest-dialog.tsx
web/src/components/finding-retest-panel.tsx
web/src/components/finding-traffic-panel.tsx
web/src/components/http-code-block.tsx
web/src/components/link-traffic-dialog.tsx
web/src/components/markdown.tsx
web/src/components/mention-textarea.tsx
web/src/components/scope-text-editor.tsx
web/src/components/side-question-workspace.tsx
web/src/components/simple-icon.tsx
web/src/components/status-badge.tsx
web/src/components/table-pagination.tsx
web/src/components/task-llm-profile-chain.tsx
web/src/components/task-template-controls.tsx
web/src/components/todo-popover.tsx
web/src/components/traffic-evidence-viewer.tsx
web/src/components/traffic-picker-dialog.tsx
web/src/components/transcript.tsx
web/src/components/ui/accordion.tsx
web/src/components/ui/alert-dialog.tsx
web/src/components/ui/alert.tsx
web/src/components/ui/aspect-ratio.tsx
web/src/components/ui/attachment.tsx
web/src/components/ui/avatar.tsx
web/src/components/ui/badge.tsx
web/src/components/ui/breadcrumb.tsx
web/src/components/ui/bubble.tsx
web/src/components/ui/button-group.tsx
web/src/components/ui/button.tsx
web/src/components/ui/calendar.tsx
web/src/components/ui/card.tsx
web/src/components/ui/carousel.tsx
web/src/components/ui/chart.tsx
web/src/components/ui/checkbox.tsx
web/src/components/ui/collapsible.tsx
web/src/components/ui/combobox.tsx
web/src/components/ui/command.tsx
web/src/components/ui/context-menu.tsx
web/src/components/ui/dialog.tsx
web/src/components/ui/direction.tsx
web/src/components/ui/drawer.tsx
web/src/components/ui/dropdown-menu.tsx
web/src/components/ui/empty.tsx
web/src/components/ui/field.tsx
web/src/components/ui/hover-card.tsx
web/src/components/ui/input-group.tsx
web/src/components/ui/input-otp.tsx
web/src/components/ui/input.tsx
web/src/components/ui/item.tsx
web/src/components/ui/kbd.tsx
web/src/components/ui/label.tsx
web/src/components/ui/marker.tsx
web/src/components/ui/menubar.tsx
web/src/components/ui/message-scroller.tsx
web/src/components/ui/message.tsx
web/src/components/ui/native-select.tsx
web/src/components/ui/navigation-menu.tsx
web/src/components/ui/pagination.tsx
web/src/components/ui/popover.tsx
web/src/components/ui/progress.tsx
web/src/components/ui/radio-group.tsx
web/src/components/ui/resizable.tsx
web/src/components/ui/scroll-area.tsx
web/src/components/ui/select.tsx
web/src/components/ui/separator.tsx
web/src/components/ui/sheet.tsx
web/src/components/ui/sidebar.tsx
web/src/components/ui/skeleton.tsx
web/src/components/ui/slider.tsx
web/src/components/ui/sonner.tsx
web/src/components/ui/sortable-head.tsx
web/src/components/ui/spinner.tsx
web/src/components/ui/switch.tsx
web/src/components/ui/table.tsx
web/src/components/ui/tabs.tsx
web/src/components/ui/textarea.tsx
web/src/components/ui/toggle-group.tsx
web/src/components/ui/toggle.tsx
web/src/components/ui/tooltip.tsx
web/src/config/app-config.ts
web/src/hooks/use-current-user.ts
web/src/hooks/use-lg.ts
web/src/hooks/use-mobile.ts
web/src/hooks/use-side-questions.ts
web/src/lib/activity-merge.test.mjs
web/src/lib/activity-merge.ts
web/src/lib/api.ts
web/src/lib/auth.ts
web/src/lib/chat-mentions.test.mjs
web/src/lib/chat-mentions.ts
web/src/lib/chat-send-mode.ts
web/src/lib/company-scope.ts
web/src/lib/cookie.client.ts
web/src/lib/fonts/registry.ts
web/src/lib/local-storage.client.ts
web/src/lib/mock/data.ts
web/src/lib/mock/enabled.ts
web/src/lib/mock/handler.ts
web/src/lib/preferences/layout-utils.ts
web/src/lib/preferences/layout.ts
web/src/lib/preferences/preferences-config.ts
web/src/lib/preferences/preferences-storage.ts
web/src/lib/preferences/theme-utils.ts
web/src/lib/preferences/theme.ts
web/src/lib/side-questions.ts
web/src/lib/sort-preference.ts
web/src/lib/status.ts
web/src/lib/task-assets.ts
web/src/lib/types.ts
web/src/lib/utils.ts
web/src/navigation/sidebar/sidebar-items.ts
web/src/proxy.disabled.ts
web/src/proxy.ts
web/src/scripts/generate-theme-presets.ts
web/src/scripts/theme-boot.tsx
web/src/stores/preferences/preferences-provider.tsx
web/src/stores/preferences/preferences-store.ts
web/src/styles/flag-icons/flags.css
web/src/styles/presets/brutalist.css
web/src/styles/presets/soft-pop.css
web/src/styles/presets/tangerine.css
web/tsconfig.json
web/tsconfig.scripts.json
```

문구 해석 제외 비UTF-8 목록:
```
screenshots/agents.png
screenshots/assets.png
screenshots/assets_test.png
screenshots/chat.png
screenshots/dashboard.png
screenshots/findings.png
screenshots/graph.png
screenshots/intercept.png
screenshots/llm.png
screenshots/logs.png
screenshots/sessions.png
screenshots/tasks.png
screenshots/traffic.png
screenshots/wx.png
web/media/dashboard.png
web/public/logo.png
web/src/app/favicon.ico
web/src/app/icon.png
```

범위 보류/원문 유지: README.md, CHANGELOG.md, docs/漏洞流量证据.md. 미검토 범위는 전체 한국어 문장의 원문 의미 전수 대조와 실제 실행/모델/DB/브라우저 동작이다.

실행 검증: **사용자 요청으로 미실행**. fetch·병합·stage·commit·push 없음.

최종 보존: 이번 재확인·산출물 재작성 전후 추적/기존 미추적 파일 SHA-256 모두 동일, Git index SHA-256 동일. 기존 로컬 변경 포함 동일.

## 갱신: S001·S002·S003·S006 보완 결과

이 절은 앞선 검색 당시 후보 목록에 대한 최신 정정이다. 앞선 항목은 이력으로 보존하며 현재 판단은 본 절과 TSV 추가 열/행을 따른다. A001–A090 주 분류 수는 그대로이고, S항목은 별도 집계한다.

- S001: 부분 수정·정적 대응 확인. 174행(사례 이름 5개 포함)을 번역했으며 전문 용어 4개 블록은 보류. 전체 완료 아님. 입력·fixture·기대값·조건은 유지.
- S002/S003/S006: 지정 사례 7개·공격 표면 표기 6곳·조사 2곳 수정 완료/정적 대응 확인. 실행 미검증.
- A049: 정규식 DNT는 확정 보존, depth/coverage 기능 차이는 별도 미해결 후속(secondary_status). DNT를 이유로 기능 문제까지 해결 처리하지 않는다.
- A050: 표시값(Agent 이름과 모델 본문 삽입 연동), 모델용 mock prompt/skill 선택 조건, 감사 입력/원문 fixture를 구분한다. 감사 입력과 원본 도구 입출력은 중국어라는 이유만으로 누락 집계하지 않으며 실제 prompt/선택 안내는 후속 검토다. data.ts 미수정.
- A052/A090/S004/S005 및 U7/U15 후속은 이번 작업에서 제외하여 유지. 새 전문 용어는 확정 등록하지 않았다.

보류 원문과 후보: `/tmp/artex-audit-fix-final-text-review.md`. 기존 감사 원본은 발견되는 경우 변경하지 않았고, 이번에 읽은 상태 보고서 작업 전 판은 `/tmp/artex-final-text-review/status-before.md`와 `.tsv`로 보존했다. 실행 검증은 사용자 요청으로 미실행.


## 최신 갱신: S001 H01–H04 해소

앞선 S001의 부분 완료/용어 보류 상태는 이력이며 현재 상태는 **수정 완료 / 정적 대응 확인 / 실행 미검증**이다. 打桩·令牌桶·量纲를 문맥 제한과 함께 확정 등록하고 H01–H04를 번역했다. 전송 시점 설명 두 곳도 정정했다. notify_api_test.go 누적 186행 변경, 951줄·LF·코드·입력·fixture·기대값·조건·포맷을 보존했다. 「近 30 分钟新增 1 个漏洞」 및 검증용 원문 인용은 예시/인용 보존이다.

S002는 지정 t.Run 7개만 수정 완료이며, 범위 밖 table 사례 이름 8개는 남은 번역 작업으로 유지한다. S003/S006의 완료 상태와 A001–A090의 주 분류 수는 그대로다. 별도 S행 4개는 지정 범위 수정 완료/정적 대응 확인이며 S002의 추가 잔여는 secondary_status로 명시했다. 실행 검증은 전 항목 사용자 요청으로 미실행이다.

A049 정규식 DNT와 미해결 depth/coverage 기능 문제, A050의 표시값/모델용 mock 프롬프트/감사 입력·fixture 구분은 유지한다. A052·A090·S004·S005와 U7/U15 후속은 이번에 해결하지 않았다. 전체 감사·전체 번역 완료를 뜻하지 않는다. 상세 근거와 이전 보류 원문은 `/tmp/artex-audit-fix-final-text-review.md`에 기록했다. 원래 통합 감사 보고서는 수정하지 않았다.


## 최신 갱신: S002 잔여 사례 이름·S004 구분자·S005 주석 인용

현재 로컬 기준 HEAD `564c80439114a88c58cb1c2dcaf3593199c48a41`, work/ko-translation. 직전 감사 보완 커밋을 확인했고 fetch·병합하지 않았다. 원격 참조 최신성은 보장하지 않는다. 앞선 S002의 table 이름 8개 잔여 및 S004/S005 추가 확인 후보 기록은 이력으로 보존하며 현재 판단은 본 절과 TSV를 따른다.

- **S002 수정 완료 / 정적 대응 확인**: server/engine_emptyturn_test.go:29,30,31,35,39,41,45,46의 사례 name만 번역했다. 앞선 t.Run 7개 번역과 fixture·입력·기대값·thinking/guard 계약값·로직·146줄은 동일하다. 변경 전 이름의 저장소 내 직접 참조 없음. 외부 개인 -run 명령은 미확인이다.
- **S004 수정 완료 / 정적 대응 확인**: intercept/page.tsx:542–543 및 agent-editor.tsx:1140은 화면/툴팁 표시, planner.go:185·191은 모델 입력 자유 서술 목표/힌트 목록이다. join 결과에 비교·파싱·저장·모델 출력 계약 의존을 발견하지 못하여 지정 구분자 5곳만 ", "로 변경했다. 원래 API/DB 배열·항목·순서·지시 강도 보존. 다른 구분자를 전역 치환하지 않았다.
- **S005 수정 완료 / 정적 대응 확인**: conversations.go:561의 주석 인용만 현재 AbortChatStoppedByUser.Short인 “사용자가 이번 대화를 중지했습니다”로 맞췄다. pgStopConversation→Chat→captureRunSession→terminalText의 요약 포함 경로를 확인했다. 인용은 전체 요약이 아닌 포함된 사유 문구이며 영어 주석/운영 소스 출력은 보존했다.

A001–A090의 주 분류와 이력은 변경하지 않았다. 별도 S001–S006 6개는 지정 범위 수정 완료/정적 대응 확인(실행 미검증)이며 새 용어/계약 보류는 발견하지 못했다. A052·A090·A050 및 U7/U15 정책·기능·호환성 후속은 유지한다. 전체 감사/번역 완료나 원문 기능 문제 해결을 뜻하지 않는다.

전체 diff `/tmp/artex-audit-fix-labels-separators.diff`, 상세 근거와 파일별 줄 수 `/tmp/artex-audit-fix-labels-separators-review.md`. 상태 갱신 전 MD/TSV는 `/tmp/artex-labels-separators-review/status-before.md`와 `.tsv`에 보존했다. 원래 통합 감사 보고서·TSV는 수정하지 않았다. 실행 검증은 **사용자 요청으로 미실행**이다.
