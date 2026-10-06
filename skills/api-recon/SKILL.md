---
name: api-recon
description: 웹사이트 API 엔드포인트를 수집할 때 이 skill을 호출합니다.
---

# API Recon(프런트엔드 엔드포인트 정찰)

**권한을 부여받은** 상태에서 **백엔드 API**(경로·메서드·파라미터·응답 본문), **프런트엔드 라우트**, **UI 기능 트리거 지점**(Tab·팝업·표 동작 등)을 가능한 한 빠짐없이 발견합니다.

---

## 경계와 금지 사항(Agent 필독 · 위반 시 범위 이탈)

이 skill은 **API / 파라미터 영역 정찰만 수행**하며, 취약점 탐색이나 침투 공격 단계가 아닙니다.

### 작업 경계

| 범위 | 허용 | 금지 |
|---|---|---|
| **목표** | path·method·파라미터·라우트·UI 트리거 지점 열거 | SQLi/XSS/권한 우회/무차별 대입/fuzz 취약점·요청 변조 공격·파괴적 동작 |
| **인증** | Hook + stub/mock으로 **클라이언트** 로그인 관문 우회 | 사용자에게 계정/비밀번호 요구 또는 추측·실제 로그인 폼 제출 시도 |
| **런타임** | 자격 증명 없이 hook으로 엔드포인트를 가로채고 mock 응답으로 SPA의 로그인 후 기본 화면에 진입 | 실제 백엔드 세션에 의존해야 계속할 수 있는 흐름 |

### 자격 증명 없는 동적 분석(Phase 3 기본)

1. `preload.js` / `runtime_harvest.js`를 통해 **가로채고 stub 처리**합니다. 대상은 로그인·권한·메뉴 등의 bootstrap 엔드포인트입니다.
2. 업무 조회 엔드포인트에 **구조가 올바르고 애플리케이션 응답 코드가 성공이며 데이터는 비어 있어도 되는** mock body를 반환합니다.
3. 백엔드가 없거나 401이 반환되는 환경에서도 프런트엔드가 로그인 후 페이지를 렌더링하도록 하여 더 많은 XHR/fetch/WebSocket을 트리거합니다.
4. **빈 데이터·빈 표·자리만 표시된 UI는 모두 예상된 결과입니다**. 이 때문에 실제 로그인이나 취약점 테스트로 전환하지 마세요.

**요약**: mock으로 프런트엔드 라우트와 컴포넌트 마운트가 이루어지게 하고, **outbound 요청만 기록합니다**. 백엔드가 무엇을 반환하는지는 중요하지 않으며, 프런트엔드가 **어떤 엔드포인트를 추가로 요청하는지**가 중요합니다.

### 절차상 엄격한 금지 사항

| 금지 | 대안 |
|---|---|
| Phase 1 완료 전에 grep/curl/Read로 메인 entry `index-*.js`에서 API path 추출 | `OUTDIR/harvest_static.py` 실행 |
| `extract_apis.py` 등 harvest를 대체하는 스크립트를 직접 작성 | `OUTDIR/harvest_static.py`를 수정한 뒤 다시 실행 |
| 동일한 grep/명령이 ≥2회 실패했는데도 반복 | 전략 변경: tool_logs 읽기·harvest 수정·reference 확인 |
| 실행 전 필수 확인 A/B를 건너뛰고 `scripts/` 원본을 직접 실행 | OUTDIR에 복사한 뒤 대상에 맞게 수정 |
| 실제 사용자 이름/비밀번호·OTP·OAuth 등의 인증 | stub/mock(위 내용 참조) |
| '실제 데이터 확보'를 이유로 stub을 건너뛰고 권한 우회/인젝션 테스트 수행 | outbound만 기록, recon 경계에 해당 |
| 삭제·민감정보 내보내기·일괄 쓰기 등의 비가역적 동작 | coverage 클릭에도 동일하게 적용 |
| runtime + 동적 열거를 완료하지 않고 모든 페이지와 엔드포인트를 확보했다고 주장 | '완료 정의' 참조 또는 한계 명시 |
| 파라미터 트리거 매트릭스 + diff를 완료하지 않고 모든 파라미터를 파악했다고 주장 | Phase 3b 매트릭스 + Phase 5 diff |
| 단일 runtime 샘플로 필수/선택 사항 추론 | 여러 샘플의 diff 또는 검증 규칙/오류 기반 역추론 |

---

## 두 계층 모델 + 실행 모드

| 계층 | 산출물 | 한계 |
|---|---|---|
| **정적**(JS bundle) | 전체 endpoint 경로·라우트 초안·요청 구성 지점의 필드 후보 | HTTP 메서드는 없습니다. 파라미터는 Phase 1b가 필요하며, 런타임에 조합되는 URL은 놓칩니다 |
| **런타임**(활성 세션) | 메서드 + body + 응답 + 동적 URL + WS/SSE. 여러 샘플의 diff로 파라미터 보완 | 페이지가 실제로 렌더링되어야 요청이 발생합니다. 단일 샘플만으로 필수/선택 사항을 확정할 수 없습니다 |

| 실행 모드 | 엔진 | 용도 |
|---|---|---|
| **depth** | `runtime_harvest.js`(Puppeteer) | API 목록·METHOD/params/응답 본문·WS/SSE·재현 가능한 일괄 실행 |
| **coverage** | browser + `preload.js` | Tab/팝업/표 클릭으로 기능 지점을 더 깊이 확인 |
| **both** | depth 후 coverage | 가장 완전하며, 가장 오래 걸립니다 |

**파라미터 분석 방법**(범용 스크립트 없음): path는 harvest/정규 표현식으로 찾고, 파라미터는 **앵커 주변 검색 범위 확대 + UI 바인딩 연결 관계 + 여러 샘플의 diff + 오류 기반 역추론**으로 확인합니다(grep 방법은 [reference.md](reference.md) J절 참조).

---

## 완료 정의

다음을 모두 충족해야 recon 완료라고 주장할 수 있습니다:

- [ ] **정적**: Phase 1 harvest에서 `api_static.txt`·`routes.txt`·`js/` 생성
- [ ] **런타임**: depth 또는 coverage 중 적어도 하나. coverage/both는 **Hook 적용 + 동적 열거 반복 절차**가 필수
- [ ] **로그인 후 화면 진입**: 업무 path에 접근할 때 `/login`이 아님(hash 라우트 주의)
- [ ] **파라미터**: coverage/both에서 파라미터 트리거 매트릭스 + `param_samples.json` 완료. Phase 5에서 `params_merged.json` 병합
- [ ] **심층**(모듈 페이지가 비어 있는 경우): Phase 4 권한 트리 복원 후 다시 실행하며, **module 수준 API**가 나타날 때까지 반복(locale/bootstrap만으로는 부족)
- [ ] **제출**: Phase 5 산출물이 모두 갖춰져 있어야 합니다(Phase 5 산출물 표 참조). `insert_assets`로 서비스와 엔드포인트 자산 기록

---

## 스크립트와 실행 전 필수 확인

`scripts/`는 참조 템플릿일 뿐이며, 원본을 직접 실행하고 최종 결과로 취급하는 것은 **금지합니다**.

**규칙**: 먼저 읽기 → 대상에 맞게 수정 → `OUTDIR`(예: `recon/`)에 저장 → `CHANGES.md`에 기록. 맞지 않으면 방법론에 따라 다시 작성하며 구조만 참고합니다.

| 실행 전 필수 확인 | 시점 | 참조 스크립트 → OUTDIR 사본 | 일반적인 필수 수정 사항 |
|---|---|---|---|
| **A(정적)** | Phase 0 후, harvest/spider를 **처음** 실행하기 전 | `harvest_static.py` / `spider_mpa.py` | **대부분의 사이트는 기본 regex로 바로 실행 가능**. manifest/표현 방식이 맞지 않을 때만 endpoint 정규 표현식·webpack/Vite `publicPath`·MPA exclude/cookie 수정 |
| **B(런타임)** | Phase 2 후, depth/coverage 실행 전 | `runtime_harvest.js` / `preload.js` + `config.json` | Cookie/localStorage 키·neutralize 성공값·stubs·login 정규 표현식·api 접두사·hash/history |

**SPA 필수 순서**(순서 변경 불가. Phase 번호가 '먼저 탐색한 뒤 스크립트 실행'보다 우선):

| 단계 | 필수 | 금지 |
|---|---|---|
| Phase 0 완료 후 | 다음 Bash = `python3 OUTDIR/harvest_static.py <URL> OUTDIR` | curl/grep/Read로 메인 entry `index-*.js` 처리(보통 >500KB) |
| 실행 전 필수 확인 A | 스크립트 복사 → 필요에 따라 소폭 수정 → **즉시 실행** | 먼저 API를 직접 추출한 뒤 harvest 여부 결정 |
| Phase 1 완료 전 | `wc -l`로 산출물 확인. 404이면 harvest 수정 후 재시도 | extract 스크립트 직접 작성·다운로드하지 않은 URL에 반복 grep |
| Phase 1b부터 | grep은 `OUTDIR/js/*.js`에만 사용 | 메인 bundle로 harvest 대체 |

- ✅ `harvest_static.py` 복사 → (선택 사항) regex 수정 → **즉시 실행**
- ❌ curl로 메인 bundle 요청 → grep 여러 번 → 임시 extract 작성 → 마지막에 harvest
- **MPA**: Phase 0 후 다음 Bash = `python3 OUTDIR/spider_mpa.py ...`

---

## 도구와 출력 제약 조건

| 제약 조건 | 설명 |
|---|---|
| 대용량 파일 | >100KB인 `index-*.js`를 Read/grep으로 컨텍스트에 넣는 것은 **금지합니다**. OUTDIR 스크립트로 일괄 처리합니다 |
| grep 출력 | 반드시 `\| head -20` 또는 `-m 5`를 사용해야 합니다. 대화에는 path 요약만 남기고 bundle 조각을 붙이지 마세요 |
| 확인 | `wc -l`·`ls \| wc -l`을 사용합니다. 디렉터리 전체를 Read하지 마세요 |
| regex 초기 탐색 | 선택 사항, ≤1회, ≤50KB인 작은 chunk 또는 HTML에만 적용. 정식 정적 분석은 harvest를 기준으로 합니다 |
| reference | 방법/템플릿/문제 해결은 [reference.md](reference.md) 참조. inline으로 전문을 반복하지 마세요 |

---

## 실행 로드맵

```
Phase 0 분류 + OUTDIR
  → 실행 전 필수 확인 A → Phase 1 harvest(★ 즉시 실행 ★)
  → Phase 1b 파라미터 역분석
  → Phase 2 인증·인가 확인의 세 관문 → config.json
  → 실행 전 필수 확인 B → Phase 3 런타임 + 파라미터 매트릭스
  → Phase 4 권한 트리(필요한 경우)→ Phase 3 다시 실행
  → Phase 5 보고서 병합 + insert_assets로 발견한 모든 서비스·엔드포인트 api 자산을 일괄 등록하며, 어떤 경우에도 등록 시 발견한 자산을 누락해서는 안 됩니다
```

순서대로 체크합니다. **앞 항목을 완료하지 않으면 다음 Phase로 넘어가서는 안 됩니다**.

1. [ ] **Phase 0**: SPA/MPA 초기 조사, `OUTDIR` 생성 → [Phase 0](#phase-0--분류)
2. [ ] **실행 전 필수 확인 A + Phase 1**: 스크립트 복사 → **즉시** harvest → `wc -l` 확인 → [Phase 1](#phase-1--정적)
3. [ ] **Phase 1b**: 앵커 주변 검색 범위 확대 + 바인딩 계층 → `param_candidates.json` → [Phase 1b](#phase-1b--파라미터-역분석)
4. [ ] **Phase 2**: 인증·인가 확인의 세 관문 → `config.json` → [Phase 2](#phase-2--인증인가-확인의-세-관문)
5. [ ] **실행 전 필수 확인 B**: runtime 스크립트 수정 → [Phase 3](#phase-3--运行时)
6. [ ] **Phase 3**: depth / coverage / both, 로그인 후 화면 진입 확인, 파라미터 트리거 매트릭스 → `param_samples.json`
7. [ ] **Phase 4**(필요한 경우): 권한 트리 → patch stubs → Phase 3 다시 실행 → [Phase 4](#phase-4--权限树还原)
8. [ ] **Phase 5**: 산출물 병합 + 보고서 + `insert_assets` → [Phase 5](#phase-5--合并与报告)

---

## Phase 0 — 분류

진입점 HTML을 가져오고 **`OUTDIR`을 생성합니다**(skill 안의 `scripts/`를 수정하지 마세요):

- **SPA**: 빈 기본 화면 + `<div id=app>` + chunk → Phase 1–5
- **MPA**: SSR + `<form>`, endpoint bundle 없음 → 실행 전 필수 확인 A 후:

```bash
python3 recon/spider_mpa.py <BASE_URL> <OUTDIR> [--cookie "session=..."] [--max 300] [--depth 5] [--exclude "logout|delete|destroy"]
```

산출물: `forms.txt`, `links.txt`, `api_inline.txt`. SPA에서 forms ≈ 0이면 → Phase 1로 전환합니다.

---

## Phase 1 — 정적

[스크립트와 실행 전 필수 확인](#스크립트와-실행-전-필수-확인) · [도구와 출력 제약 조건](#도구와-출력-제약-조건)을 준수합니다.

```bash
python3 recon/harvest_static.py <BASE_URL> <OUTDIR>
```

harvest: HTML script 파싱 → webpack/Vite manifest → 모든 lazy chunk 다운로드 → `js/`, `api_static.txt`, `routes.txt`, `chunkmap.txt` 생성.

```bash
wc -l OUTDIR/api_static.txt OUTDIR/routes.txt
ls OUTDIR/js | wc -l
```

- chunk 개수 vs manifest: 404이면 반드시 harvest를 수정해 재시도해야 합니다. curl로 각 chunk를 직접 가져오지 마세요
- `api_static.txt`가 너무 적으면 → OUTDIR 안의 endpoint 정규 표현식을 넓힌 뒤 다시 실행합니다(reference 참조)

### Phase 1b — 파라미터 역분석

path는 Phase 1에서 가져옵니다. 파라미터 필드는 별도로 recon해야 합니다. grep 규칙은 [도구와 출력 제약 조건](#도구와-출력-제약-조건)을 참조합니다.

**완료 기준**: 중요한 엔드포인트의 필드명, 전송 위치, 타입 추론, 필수 여부, 샘플값, 신뢰도를 설명할 수 있어야 합니다.

#### 1b.0 — 전송 형태

| 형태 | 파라미터 위치 | 정적 분석에서 먼저 확인할 것 |
|---|---|---|
| REST JSON | body + query | path 앵커 옆 `(params\|data\|body)\s*:\s*\{` |
| GraphQL | `variables` | gql 템플릿, `$page: Int` |
| 기존 form | urlencoded | `<form>`, `FormData` |
| 파일 업로드 | multipart | `FormData.append` |
| 경로 파라미터 | `/user/:id` | 라우트 표 + `useParams` / `$route.params` |
| 암호화/서명 | `sign`/`data`에 포함 | Hook으로 암호화 함수의 입력 파라미터 관찰(reference D절) |

산출물: 각 엔드포인트에 `transport: query|json|form|graphql|encrypted`를 표시합니다.

#### 1b.1 — 앵커 주변 검색 범위 확대

알려진 path를 앵커로 삼아 주변 검색 범위를 넓혀 요청 구성 객체를 찾습니다:

```bash
grep -n '"/api/user/list"' OUTDIR/js/*.js | head -20
grep -rhoaE '.{0,120}("/api[^"]+").{0,200}' OUTDIR/js/*.js | head -20
grep -rhoaE '(params|data|body|payload)\s*:\s*\{' OUTDIR/js/*.js | head -20
```

| 래퍼 계층 | 파라미터 단서 |
|---|---|
| axios 인스턴스 | `data` / `params` |
| 공통 request | 인터셉터가 전역 필드를 추가 |
| OpenAPI 클라이언트 | method 시그니처 생성 |
| React Query / SWR | hook의 두 번째 파라미터 |
| Vue composable | composable 입력 파라미터 |

남아 있는 타입 정보: `yup`/`zod`/rules, `Form.Item name=`, 내장 Swagger.

→ `param_candidates.json`: `{ path, fields[], source: "static-callsite", confidence }`

#### 1b.2 — 바인딩 계층

```
Form field → onFinish/handleSubmit → transform → API payload
```

| 바인딩 출처 | 방법 |
|---|---|
| 폼 submit | submit → transform → API를 추적 |
| 표 검색 | `getFieldsValue()` → `params` |
| 라우트 | `:id` / `?tab=` |
| 인터셉터 | 전역 `tenantId`, 페이지네이션, sign |
| 열거 select | `options` → API 열거값 |

DevTools call stack에서 `fetch`/`XHR.send`부터 거슬러 올라가 요청 구성 함수를 추적합니다.

#### 1b.3 — 요청 구성의 세 질문(≠ Phase 2 인증·인가 확인의 세 관문)

| 질문 | 확인할 내용 |
|---|---|
| **구성** | payload를 어디서 build하는지, transform 흔적 |
| **검증** | required, pattern, enum |
| **전송** | path / query / body / multipart / 헤더 |

인터셉터 관문(Phase 2)에서 전역 추가 필드(Authorization, `X-Tenant-Id`, sign)도 함께 읽습니다.

#### 1b.4 — Phase 3과 연결

후보 필드는 정적 분석/바인딩 계층에서 가져옵니다. **필수/선택/조건 의존성**은 반드시 Phase 3 파라미터 매트릭스 + diff + Phase 5 오류 기반 역추론으로 확인해야 합니다.

---

## Phase 2 — 인증·인가 확인의 세 관문

`OUTDIR/js/`에서 grep(`head` 포함)하고 `config.json`에 기록합니다(방법은 reference 참조):

| 관문 | 질문 | 키워드 |
|---|---|---|
| **렌더링 관문** | 로그인 여부를 어떻게 판단하는가? | `isLogin`, `getToken`, Cookie/localStorage |
| **인터셉터 관문** | 무엇이 `/login` 이동을 트리거하는가? | `response_code`, `errno`, axios interceptor |
| **콘텐츠 관문** | 메뉴/권한은 어디에서 오는가? | `menu`, `permission`, `role`, `acl`, `routes` |

localStorage 키 이름을 자격 증명으로 취급해서는 안 됩니다. 반드시 chunk/요청 연결 관계에서 확인해야 합니다.

**출구 = 실행 전 필수 확인 B**: 결론을 `config.json`에 반영하고 `OUTDIR/runtime_harvest.js` / `preload.js`를 수정합니다.

### Phase 2b — API 관찰(선택 사항)

OUTDIR 안의 `preload.js`로 세션 키 이름, Authorization, 중첩 API URL을 확인합니다:

| 설정 | 산출물 |
|---|---|
| `recordDetail: true` | `__API_RECON_DETAIL__` |
| `observe.xhrHeaders: true` | headers 관찰 |
| `extractUrlsFromResponse: true` | 응답 안의 하위 API |
| `observe.storageReads/cookieReads: true` | config에 반영 |
| `neutralizeVueRouter: true` | `__API_RECON_ROUTES__` |

coverage의 각 회차에서 내보냅니다: `__API_RECON_LOG__`, `__API_RECON_DETAIL__`, `__API_RECON_ROUTES__`, `__API_RECON_OBSERVE__`.

---

## Phase 3 — 运行时

须已过门禁 B；遵守 [边界与禁止](#경계와-금지-사항agent-필독--위반-시-범위-이탈) · 无凭据 mock 策略。

`config.json` 设置 `"runtimeMode": "depth" | "coverage" | "both"`（模板见 reference）。

### Hook 与 stub（depth + coverage 共用）

| 层 | 范围 | 目的 |
|---|---|---|
| L1 精确 | auth/权限/bootstrap stub | 过首屏鉴权 |
| L2 负向修正 | 所有 JSON 响应 | 未登录码 → 成功 |
| L3 兜底 | 未命中 L1 的 `/api` 等 | 空成功体，撑开 UI |

- **depth**：fake auth + `forward` 改业务码 + `stubs`；遍历 `routes`（hash/history）；产出 `runtime_api.json`
- **coverage**：**document-start** 注入 `preload.js`（CDP `addScriptToEvaluateOnNewDocument` 或 Userscript）

验证：`window.__API_RECON_PRELOAD__` 存在；业务 path 不回 `/login`。

```bash
cd recon && npm install
node runtime_harvest.js config.json
```

### 3b — coverage 动态枚举（必做）

1. 主导航/侧栏 — 每项点击，等网络 1–3s
2. Tab — `role=tab`、`.ant-tabs-tab`
3. 表格 — 首行查看/编辑/详情
4. 工具栏 — 导出、筛选、新建（**避免不可逆删除**）
5. 每进模块 — 合并 API/路由
6. SPA — 对 `routes.txt` 未覆盖 path 受控 `pushState`（MPA 禁止）

**参数触发矩阵**（必做）：每模块按操作类型各录一次，**diff 多样本**：

| 操作 | 通常多出的参数 |
|---|---|
| 列表首屏 | 分页 + 默认筛选 |
| 点搜索 | keyword、filter |
| 高级筛选 | 更多 optional |
| 新建/编辑 | 完整 entity |
| 批量/导出/排序 | `ids[]`、`exportType`、`sortField` |

**stub 下 outbound body/headers 仍真实**——以请求为准。录制 → `scan_raw.json`、`param_samples.json`、`api_detail.json`。

- **Vue**：`neutralizeVueRouter: true` + document-start preload
- **React**：`routes.txt` + 侧栏点击 + `pushState`
- **both**：先 3a depth，再 3b coverage

---

## Phase 4 — 权限树还原

**触发**：模块页空白 / 每路由仅 bootstrap（如 locale）→ 内容门未过。

| 现象 | 含义 |
|---|---|
| 进壳成功 | 渲染门 + 拦截器门已过 |
| 侧栏缺项/点击空白 | stub shape 或权限码不全 |
| 每路由 API 相同且极少 | `v-if permission` 未通过 |
| `routes.txt` 远少于 bundle | 须从 auth 模块补全 |

```bash
grep -rhoaE '"/api[^"]*(permission|perm|role|menu|acl)[^"]*"' OUTDIR/js/*.js | sort -u | head -30
grep -rhoaE 'userRouteAuth|getResultTree|routeMap|routeLink|menuList|authList' OUTDIR/js/*.js | head -20
```

典型链：`role_permissions`（flat codes）+ `permissions/all`（tree）→ `getResultTree` → `userRouteAuth[CODE].url`。

```bash
python3 recon/extract_route_map.py recon/js recon/
python3 recon/build_perm_tree.py recon/js recon/ --config recon/config.json
```

中间产出：`route_map.json`、`userRouteAuth.json`、`permissions_tree.json`、`*_stub.json`、`perm_codes_all.txt`。

stub 检查：外层 `response_code` 与拦截器门一致；flat codes 与 tree 对齐；`routes` 覆盖 `route_map` 全部 link。

更新 `config.json` 后**重跑 Phase 3**。大型 SPA 可调 `waitUntil`、`routeTimeout`、`perRouteMs`（见 reference A3/I 节）。

---

## Phase 5 — 合并与报告

### 产出表

| 文件 | 阶段 | 内容 |
|---|---|---|
| `js/`、`api_static.txt`、`routes.txt`、`chunkmap.txt` | 1 | 静态 bundle 与 path |
| `param_candidates.json` | 1b | 静态参数字段候选 |
| `config.json` | 2 | 三道门 + runtime 配置 |
| `runtime_api.json` | 3a | depth 详细录制（含 WS/SSE） |
| `param_samples.json`、`scan_raw.json`、`api_detail.json` | 3b | 多样本、点击日志、detail |
| `route_map.json` 等 | 4 | 权限树中间文件（若执行） |
| `params_merged.json` | 5 | 合并参数字段 + 置信度 |
| `api_merged.txt` | 5 | `METHOD /path [params] [static\|runtime\|both]` |
| `site_map.json` | 5 | 路由、API、params、功能点、局限 |
| **insert_assets** | 5 | 将所有服务、端点资产写入资产库 |

### 5b — 参数合并

从 `param_samples.json` diff，**无通用合并脚本**。置信度规则见 reference J7（高/中/低/待触发）。

### 5c — 错误反推

授权范围内可发不完整请求读 400（**属参数 recon，非漏洞测试**）：`field 'x' is required`、枚举错误等。注意 `data` 包装、`variables`、加密前 `bizData`。

报告须注明：runtimeMode、静态/运行时 API 数、参数置信度、未覆盖模块、相对参考脚本的 `CHANGES.md` 摘要。

`site_map.json` 建议结构：

```json
{
  "site": "https://example.com",
  "runtimeMode": "both",
  "appType": "vue-spa",
  "routeGuardStrategy": ["nav-neutralize", "L1-auth", "L2-patch", "forward"],
  "apisFromStatic": [],
  "apisFromRuntime": [],
  "apis": [],
  "params": [{ "method": "POST", "path": "/api/user/list", "transport": "json", "fields": [] }],
  "frontendRoutes": [],
  "routesVerifiedByClick": [],
  "featuresTriggered": [],
  "limitations": ""
}
```

更多字段与 grep 配方见 [reference.md](reference.md)。

---

## 通用说明

- **框架无关**：webpack/Vite/Angular lazy load 方法相同
- **传输**：REST/JSON、GraphQL、WebSocket、SSE；gRPC-web 不在范围
- **SSR**：客户端 fetch 可录；RSC/Server Actions 不完全可枚举
- **盲区**：JSVMP、WASM、HMAC/mTLS 强校验 → 静态 + 标注局限
- **参数盲区**：条件联动、hidden params、WASM 组包 → 「待触发」/「不可达」
- **静态是安全网**：runtime 被挡时静态仍能枚举 endpoint

---

## 附加资源

- Grep 配方、`config.json` 模板、排障、Hook、参数逆向 J 节、site_map 模板：**[reference.md](reference.md)**
- 参考脚本路径见 [脚本与门禁](#스크립트와-실행-전-필수-확인) 表
