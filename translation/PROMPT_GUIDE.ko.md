# ARTEX 한국어화 가이드 — 백엔드·프롬프트·인프라 (Stage 2·3+)

UI 1단계(프런트엔드 문자열, `web/`) 완료 후 **저장소 나머지 전체**를 한국어화하기 위한 작업 기준이다.
최종 목표는 **ARTEX 저장소 대부분의 중국어 한글화**이며, 빠지는 파일이 없도록 **§1 커버리지 맵**이 모든 영역을 단위로 배정한다.
용어·표기는 [`GLOSSARY.ko.md`](./GLOSSARY.ko.md)를 그대로 따른다.

작업은 성격에 따라 **3개 트랙**으로 나뉜다:
- **트랙 A · 프롬프트(U1~U9)** — 에이전트가 읽는 프롬프트. **번역이 모델 동작을 바꿀 수 있어** 검증 비중이 가장 크다(§2 규칙 전체 적용).
- **트랙 B · 백엔드 Go 비프롬프트(U10~U13)** — API 에러·DTO 주석·DB·알림 템플릿·로그. UI와 비슷한 난이도(§2 중 DNT·`%verb`·전각 규칙 적용, "모델 동작" 항목은 해당 없음).
- **트랙 C · 스크립트·인프라·설정(U14)** — `*.sh`·Dockerfile·compose·CI·설정 예시. CLI 출력·주석.

> 공통 규칙(§2)은 트랙 A에 **전부**, 트랙 B·C에는 **DNT·포맷 지정자·전각→반각·조사** 부분이 적용된다. "지시 강도·출력 언어" 같은 프롬프트 전용 항목은 트랙 A만.

---

## 0. 범위와 전제

- **Stage 2(출력 한국어화)**: 프롬프트 안의 "중국어로 답하라"류 지시를 "한국어로 답하라"로. 지시문이 프롬프트 본문에 박혀 있어 Stage 3와 **같은 파일에서 함께** 처리한다.
- **Stage 3(프롬프트 본문 번역)**: 프롬프트의 중국어 설명을 한국어로.
- **확인된 사실**: 프롬프트가 지시한 중국어를 Go 코드가 **직접 문자열 매칭하는 결합은 발견되지 않았다**(`Contains`/`HasPrefix`/`==` 중국어 리터럴 0건). 따라서 머지 충돌·파싱 깨짐 위험은 낮다.
- **그러나**: 표현·지시 강도·출력 형식이 바뀌면 모델의 **도구 선택·구조화 출력이 달라질 수 있다**. 그래서 "독립 번역 안전"이 아니라 **"단위별 동작 검증 필수"** 로 다룬다.

---

## 1. 작업 단위와 커버리지 맵 (U1 → U14)

한 번에 **한 단위씩** 번역 → 검토 → 검증 → 커밋한다. **각 단위는 "그 파일들 + 같은 이름의 `*_test.go`"를 함께** 처리한다(테스트가 번역된 문자열을 기대하면 같은 커밋에서 기댓값도 갱신하거나 사유 기록). 트랙 A는 **U1을 먼저** 끝내 기준을 깐 뒤 진행.

### 트랙 A · 프롬프트 (agent/intercept/report/sidequestion)

**U1 · 공용 기반 ⭐가장 먼저** — `agent/promptcatalog.go`·`prompt.go`·`assembly.go`·`chat.go`
모든 에이전트가 공유하는 프롬프트 토대. `promptcatalog.go`에 공통 시스템 프롬프트 상수(`DefaultAssistantPrompt` 등), `prompt.go`의 `renderSystem`이 `{{.Var}}` 템플릿 렌더링, `assembly.go`가 조립, `chat.go`가 소비. **여기서 "한국어로 답변" 표준 문구(§2.5)와 공통 용어를 확정**해 U2~U9가 참조.

**U2 · 플래너 계열** — `agent/planner.go`·`goals.go`·`constraints.go`
계획 수립 루프(의도 생성, 목표/제약 관리).

**U3 · 워커 실행** — `agent/worker.go`·`wrapup.go`·`compaction.go`·`noa.go`·`capture.go`·`coldgraph.go`
워커 1회 실행(run) 수명주기 — 실행 프롬프트 + 마무리 + 컨텍스트 압축 + 캡처/콜드그래프.

**U4 · 메인/리테스트/발견** — `agent/mainagent.go`·`retester.go`·`finding_workflow.go`·`finding_recorder.go`
대화형 메인 에이전트 + 재검증 + 발견 기록.

**U5 · 도구 설명(분량 최대)** — `agent/tools.go`·`tools_insert.go`·`toolcatalog.go`·`tools_digest.go`
에이전트가 읽는 도구 설명·카탈로그. DNT 최다(도구명·파라미터·스키마 키). 크면 `tools.go`/`tools_insert.go`로 2분할.

**U6 · 상태/사유 문자열** — `agent/terminalreason.go`·`cancelcause.go`·`provider.go`
종료/취소 사유 코드→설명 맵(표시용) + provider 주석/에러. 사유 **코드 키** DNT, 설명만 번역.

**U7 · 인터셉트** — `intercept/` **패키지 전체**(prompt.go·review_context.go·intercept.go·trace.go 등)
모델 대체 승인 심판 프롬프트 + 인터셉트 로직 주석·메시지. ⚠️ `[模型]` 접두사는 프런트(approval-records.tsx) 계약 — 바꾸면 핑.

**U8 · 리포트** — `report/` **패키지 전체**(findings.go·report.go)
취약점 리포트 생성 프롬프트 + 로직.

**U9 · 사이드 질문** — `sidequestion/` **패키지 전체**(context.go·request.go·service.go 등)
`/btw` 보조 질문 프롬프트 + 오케스트레이션 로직. ⚠️ 컨텍스트 라벨 `"작업: #<id>"`가 U1 ReporterDefaultPrompt와 연동(동일 표기 유지).

**U15 · 내장 스킬 번들** — `skills/` (`api-recon/SKILL.md`·`reference.md`·`scopesentry/SKILL.md` + `scripts/*`)
에이전트가 로드하는 **스킬 지시문**(SKILL.md)·참고 문서·헬퍼 스크립트. SKILL.md 본문은 프롬프트성(트랙 A 규칙 적용)이나, **YAML 프런트매터 키·스크립트 코드·명령/경로·셀렉터는 DNT**. 스크립트는 주석만 번역.

### 트랙 B · 백엔드 Go 비프롬프트

**U10 · 알림 시스템** — `notify/` **전체**(~26 파일)
알림 채널·템플릿·전송 메시지(사용자 노출 푸시 내용). ※ `notify.go` StatusLabel은 `web/src/lib/status.ts`와 **동일 용어**로.

**U11 · 서버 API** — `server/` **전체**(~76 파일, 테스트 포함)
HTTP 핸들러 에러 메시지·DTO 주석·테스트 기댓값. ⚠️ 프런트가 매칭하는 계약 문자열(`归档不存在`, `skill 已存在` 등)은 프런트(B)와 **동시 변경 or DNT 유지**. 크므로 하위 묶음(핸들러별)으로 쪼개도 됨.

**U12 · 데이터/DB** — `db/` **전체**(~53 파일)
DB 레이어 주석·기본값·마이그레이션·로그.

**U13 · 그 외 백엔드 패키지** — `traffic/`·`selfupdate/`·`llmrec/`·`llmpool/`·`guard/`·`evidence/`·`mcphttp/`·`cmd/`·`config/`·`enrich/`
각 패키지의 주석·로그·메시지. 작은 패키지라 묶어서 처리 가능.

### 트랙 C · 스크립트·인프라·설정

**U14 · 스크립트·인프라** — `*.sh`·`start.bat`·`Dockerfile`·`docker-compose*.yml`·`.github/workflows/*.yml`·`config.example.json`·`.env.example`·`.gitignore`
설치·실행 스크립트의 CLI 출력/주석, 설정 예시 주석. 사용자가 직접 실행하는 `install.sh`/`start.sh` 등 우선순위 중.

### 문서·메타 (단위 아님 — 결정/보류/DNT)
- `README.md` — **현재 중국어 원문**(readme-i18n 규약상 원문 유지, 한국어판은 `README.ko.md`). 포크를 한국어 기준으로 바꿀지는 **팀 결정 필요**.
- `CHANGELOG.md` — 과거 이력(중국어 대량). **보류/선택**(가치 낮음).
- `docs/漏洞流量证据.md` — 설계/증거 문서(1건, ~1940줄). **파일명도 중국어** → 내용 번역 + 파일명 리네임(`취약점-트래픽-증거.md` 등) 여부는 **결정 필요**. 우선순위 중하.
- `translation/*.md` — 우리 작업 문서. 중국어는 **용어집 원문(zh)열·예시 = DNT**(번역 대상 아님).
- `web/src/lib/mock/data.ts`·`handler.ts` — 샘플 프롬프트 **보류**(후속).

### ✅ 커버리지 체크리스트 (모든 top dir 배정 — 미할당 0)
| 영역 | 담당 |
|---|---|
| `web/` | **완료**(UI 1단계), 잔여는 DNT |
| `agent/` | U1~U6 (+ 각 `*_test.go` 동반) |
| `intercept/` | U7 |
| `report/` | U8 |
| `sidequestion/` | U9 |
| `notify/` | U10 |
| `server/` | U11 |
| `db/` | U12 |
| `traffic/ selfupdate/ llmrec/ llmpool/ guard/ evidence/ mcphttp/ cmd/ config/ enrich/` | U13 |
| `*.sh *.bat Dockerfile docker-compose*.yml .github/ config.example.json .env.example .gitignore` | U14 |
| `skills/` | U15 |
| `README.md` `CHANGELOG.md` `docs/` `translation/` `mock/data.ts·handler.ts` | 결정/보류/DNT (위 참고) |

> 신규 파일이 생기면 위 표에 먼저 배정하고 작업한다. 번역 전 `git ls-files | xargs ... grep` 로 **미할당 중국어 파일이 있는지** 재확인(§6 NUL 주의).

---

## 2. 공통 규칙 (이것만 지키면 안 깨진다)

### 2.1 템플릿 문법 **전체** 보존
Go `text/template`과 `fmt`를 쓴다. **치환자뿐 아니라 제어 구문·포맷 지정자 전부 그대로** 둔다.
- 변수: `{{.Goal}}` `{{.Now}}` `{{.DataDir}}` …
- 제어: `{{if .X}}` `{{else}}` `{{range .Items}}` `{{end}}`
- 공백 제거: `{{- ... -}}`
- 포맷 지정자: `%s` `%d` `%v` … — **개수·순서 그대로**. 한국어 어순상 바꾸고 싶어도 `Sprintf` 인자에 위치로 매핑되므로 **재배치 금지**.

```
원문 : fmt.Sprintf("已完成 %d 个目标，还剩 %d 个", done, left)
한국어: fmt.Sprintf("목표 %d개 완료, %d개 남음", done, left)   // %d 2개·순서 동일
```

### 2.2 DNT — 원문 유지 (프롬프트에서 확대 적용)
다음은 **번역하면 에이전트가 깨진다**. 반드시 원문 유지:
- 도구명: `traffic_search` `update_finding_report` `get_worker_trace` …
- 파라미터·옵션명: `body_contains` `host` `limit` `id` …
- JSON 출력 키·필드명: `summary` `state` `vulnclass` …
- enum 값: `exploring` `blocked` `exhausted` `open` `pause` …
- 멘션 토큰 종류명: `@[漏洞#id]`의 `漏洞`·`资产`·`企业` 등 (§5.7 값 계약)
- reason/cause **코드 키**: `cause("paused_by_user", …)`의 `"paused_by_user"` — **설명 문자열만 번역**
- 경로·명령·env·URL·`{{.Var}}`·백틱 코드블록 내부 식별자

> 판단 기준: *"모델이 이 글자 그대로 출력해야 하거나, 코드/다른 시스템이 이 글자로 식별하는가?"* → 그렇다면 DNT.

### 2.3 지시 강도 보존
프롬프트에선 어감이 곧 동작이다. **약화도 강화도 금지.**
- `必须/一定` → **반드시** (≠ ~하면 좋다)
- `尽量/建议` → **가능하면 / 권장** (≠ 반드시)
- `禁止/不要` → **금지 / ~하지 마라** (≠ 지양)
- `先…再…` 같은 **실행 순서**, `除非/如果` 같은 **조건·예외**는 그대로 보존.

### 2.4 한국어 출력 범위
- **자연어 설명·지시 → 한국어.**
- **JSON 키 · enum · 도구명 · 필수 리터럴 → 원문.**
- **예시 출력(JSON 샘플, 도구 호출 예)도 위 원칙에 맞춘다** — 키·값 구조는 원문, 설명 주석만 한국어.

```
원문 예시 출력:
  {"summary": "对 example.com 做目录爆破", "vulnclass": "信息泄露"}
한국어화:
  {"summary": "example.com 디렉터리 무차별 대입", "vulnclass": "정보 유출"}
  // 키(summary/vulnclass)=원문, 값=자유 텍스트라 한국어. 단 enum성 값이면 원문 유지.
```

### 2.5 출력 언어 지시 통일 (U1에서 확정)
`请用简洁、准确的中文回答` / `用人话简洁回复` / `全程中文` / `使用简洁中文答复` 류 → **"한국어로 간결·정확하게 답변"** 로 통일. U1에서 표준 문구를 정하고 전 단위가 동일 표현을 쓴다.
※ 주의: `worker.go`의 "支持子串和中文(검색)" 처럼 **출력 언어가 아닌 기능 설명**의 "中文"은 출력지시가 아니다 — 문맥 보고 그냥 번역.

### 2.6 전각→반각·조사
GLOSSARY §6 동일. 산문의 `，：；（）「」、。`→반각, 변수 뒤 조사 금지(`{{.X}}` 뒤 은/는/을/를 금지, `을(를)`·`(으)로` 병기).

---

## 3. 검증 (단위 완료 기준)

| 단계 | 내용 | 필수 |
| --- | --- | --- |
| 빌드 | `go build ./...` 통과 | ✅ |
| 테스트 | 해당 단위 관련 `go test ./<pkg>/...` 통과(기댓값이 중국어면 동시 조정 or 사유 기록) | ✅ |
| 템플릿 렌더 | 프롬프트가 렌더 오류 없이 생성되는지(치환자·제어구문 보존 확인) | ✅ |
| 1 run 스모크 | 해당 단위 에이전트를 실제 1회 실행해 **도구 호출·구조화 출력**이 정상인지 | 권장(최소선) |

- 1 run은 **최소 확인**이다. 실행 못 한 부분은 **"미검증"** 으로 기록에 남긴다.
- 기호 parity(따옴표·백틱·`{{`·`}}` 개수)와 중국어 잔여(NUL 포함 파일은 `tr -d '\000'` 후) 스캔은 UI와 동일하게.

---

## 4. 단위별 작업 기록

각 단위 완료 시 아래를 `translation/PROMPT_PROGRESS.ko.md`(또는 PR 본문)에 남긴다.

```
## U5 · 도구 설명  [done]
- 대상 파일: agent/tools.go, tools_insert.go, toolcatalog.go, tools_digest.go
- 검증: build ✅ / test ✅ / 템플릿 렌더 ✅ / 1 run 스모크 ✅(traffic_search·update_finding_report 호출 정상)
- 남긴 중국어 + 보존 이유:
  - `@[漏洞#id]` 종류명 — 멘션 토큰 값 계약(§5.7)
  - 도구명/파라미터명 전부 — DNT
```

---

## 5. 커밋 규칙
GLOSSARY/TRANSLATION_PROMPT와 동일. 한국어 메시지, 한자 금지, 타입 `translation(prompt):`(코드)·`translation(glossary):`(용어집). Claude 트레일러 금지. 푸시 전 `git pull` 대신 **머지**(`git merge origin/work/ko-translation --no-edit`).

---

## 6. NUL 바이트 주의
`exploration-graph.tsx`처럼 소스에 의도적 NUL이 있으면 `file`이 바이너리로 보고 `grep -P`가 중국어를 통째로 놓친다. 검수 스캔은 `tr -d '\000' < f | grep -nP '[\x{4e00}-\x{9fff}]'`로, 대상 디렉터리에 바이너리 오인 파일이 있는지 `find … | xargs file | grep -v text`로 먼저 확인한다. NUL은 코드이므로 **제거 금지**.
