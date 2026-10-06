# api-recon — 참조 매뉴얼

Grep 방법·`config.json` 템플릿·문제 해결 안내입니다. 모든 grep은 `js/` 디렉터리를 대상으로 실행합니다. bundle이 한 줄이면 먼저 `js-beautify` 또는 `sed 's/}/}\n/g'`를 사용할 수 있으며, 보통 컨텍스트 범위가 있는 raw grep이면 충분합니다.

## 스크립트 설명

`scripts/`의 모든 파일은 **참조 템플릿**이며, 실행 전에 반드시 대상 사이트에 맞게 수정해야 합니다. 대표적인 수정 사항:

| 스크립트 | 일반적인 수정 사항 |
|---|---|
| `harvest_static.py` | endpoint 정규 표현식·webpack/Vite manifest 파싱·micro-frontend publicPath·재시도/동시성 |
| `runtime_harvest.js` | neutralize 필드명과 성공값·stub 매칭 규칙과 body 구조·routes 출처·WS 기록·`waitUntil`/`routeTimeout`/`proxy` |
| `preload.js` | `loginPathRe`·L1 stubs·`neutralize.fields`·`apiPattern`·L3 활성화 여부·`recordDetail`·`observe.*`·`neutralizeVueRouter` |
| `spider_mpa.py` | `--exclude` 파괴적 링크·cookie·depth/max·동일 도메인 필터링 |
| `extract_route_map.py` | `routeMap` / `routeLink` 정규 표현식·KEY 이름 규칙 |
| `build_perm_tree.py` | `userRouteAuth` 파싱·`ROOTS`/`PREFIX_PARENT` 계층 구성 추론 규칙·stub 외부 필드명 |
| `config.json` | 위의 모든 사이트별 파라미터를 통합하는 진입점 |

수정한 파일은 작업 디렉터리(예: `recon/`)에 두는 것을 권장하며, 보고서에 참조 스크립트와 비교하여 구체적으로 무엇을 수정했는지 명시합니다.

---

## A. 세 관문 역분석

### A1. 렌더링 관문 — '로그인 여부를 어떻게 판단하는가?'

```bash
grep -rhoaE '.{0,40}(isLogin|isAuthenticated|loggedIn|hasLogin|requireAuth)\b.{0,80}' js | head
grep -rhoaE 'function (getUser|getToken|getAuth)[0-9]?\([^)]*\)\{.{0,200}' js | head
grep -rhoaE '(localStorage|sessionStorage)\.getItem\("[^"]+"\)' js | sort -u
grep -rhoaE '(Cookies?|cookie)\.(get|load)\("[^"]+"\)' js | sort -u
grep -rhoaE '\batob\(|JSON\.parse\(|jwt|decode' js | head
```

연결 관계 `isLogin = f(getUser())` → `getUser = decode(storage.read(KEY))`를 찾아 **저장 키**, **저장 위치**(Cookie vs localStorage), **인코딩**을 확인합니다:

| 인코딩 | config 위조 방식 |
|---|---|
| 평문 문자열 / `"1"` / token | `"value": "anything-truthy"` |
| `JSON.parse(x)` | `"value": "json:{\"id\":1,\"username\":\"admin\"}"` |
| `JSON.parse(atob(x))` | `"value": "b64json:{\"id\":1,\"username\":\"admin\"}"` |
| JWT | 서명 없는/`alg:none` JWT 또는 bundle 안의 키로 서명 |
| 암호화(SM2/AES/RSA) | 하드코딩된 키를 찾습니다. 렌더링 관문에서 디코딩 가능한 blob만 필요하면 forge할 수 있고, 그렇지 않으면 정적 분석으로 대체합니다 |

→ `cookies` / `localStorage`에 기록합니다.

### A2. 인터셉터 관문 — '무엇이 /login 이동을 트리거하는가?'

```bash
grep -rhoaE '.{0,60}(interceptors\.response|axios|request\.use).{0,120}' js | head
grep -rhoaE '.{0,40}(response_code|errcode|errno|\bcode\b|\bret\b|\bstatus\b)\s*[=!]==?\s*[\-0-9]{1,4}.{0,60}' js | head -20
grep -rhoaE '.{0,40}(未登录|请重新登录|登录已过期|unauthorized|登录失效|授权|token.{0,10}invalid).{0,40}' js | head
grep -rhoaE '.{0,30}(location\.href|router\.(push|replace)|navigate)\([^)]*login[^)]*\)' js | head
```

확인할 내용: **필드명**, **성공값**(보통 `0` 또는 `200`), **리디렉션을 트리거하는 실패값**. junk session으로 검증합니다:

```bash
curl -sk -X POST -H 'Cookie: <fakekey>=junk' https://target/api/<protected> -d '{}' -H 'Content-Type: application/json'
```

→ `neutralize.fields` + `neutralize.success`에 기록합니다.

### A3. 콘텐츠 관문 — '메뉴/권한은 어디에서 오는가?'

```bash
grep -rhoaE '"/api[^"]*(permission|perm|role|menu|acl|resource|nav)[^"]*"' js | sort -u
grep -rhoaE '.{0,30}(menus|permissions|menuList|routeList|authList|role_permissions)\b.{0,120}' js | head
grep -rhoaE 'userRouteAuth|getResultTree|routeMap|routeLink|hasPermission|checkAuth' js | head
grep -rhoaE '([A-Z_][A-Z0-9_]*):\{name:"[^"]*",link:"/[^"]+"\}' js | head
```

**두 계층 데이터**(일반적인 기업 관리 화면):

| API | 일반적인 payload | 소비부 |
|---|---|---|
| `.../role_permissions` | `{ permissions: string[], role_type }` | 라우트 인증 가드, 버튼 단위 ACL |
| `.../permissions/all` | `tree[{ code, position, children }]` | 사이드바 메뉴 렌더링 |
| bundle 안의 `userRouteAuth` | `{ CODE: { url, name? } }` | code → 프런트엔드 path |
| bundle 안의 `routeMap` | `{ KEY: { name, link } }` | 별칭 해석(webpack `o.DASHBOARD`) |

소비부 코드를 읽어 확인합니다: `getResultTree(tree, permissions)`가 어떻게 필터링하는지, `v-if` / `hasAuth(code)`가 어떤 필드를 검사하는지.

**직접 forge**(소규모 사이트): permissive payload 구성 → `stubs`.

**전체 권한 트리 복원**(대규모 사이트, 사이드바/하위 모듈이 여전히 비어 있는 경우): **I절** 참조.

---

## B. config.json 템플릿

```json
{
  "baseUrl": "https://target/",
  "runtimeMode": "both",
  "chromium": "/usr/bin/chromium",

  "cookies": [
    { "name": "auth", "value": "b64json:{\"id\":1,\"username\":\"admin\",\"role\":\"admin\",\"func\":{},\"permissions\":[\"*\"]}" }
  ],
  "localStorage": { "token": "faketoken", "isLogin": "1" },

  "neutralize": {
    "fields": ["response_code", "code", "errno", "ret", "status"],
    "success": 0,
    "flags": { "success": true, "message": "ok" }
  },
  "forward": true,
  "loginUrlPattern": "/login",
  "apiPattern": "/api/|/rest/|/graphql",

  "mockTier": "L1+L2",
  "recordDetail": true,
  "observe": {
    "storageReads": false,
    "cookieReads": false,
    "xhrHeaders": true
  },
  "neutralizeVueRouter": true,
  "stubs": [
    {
      "match": "permissions/all|/menu|role_permissions",
      "body": {
        "response_code": 0, "code": 0,
        "data": {
          "permissions": ["*"],
          "menus": [
            { "name": "dashboard", "path": "/dashboard", "show": true, "children": [] },
            { "name": "alert", "path": "/alert", "show": true, "children": [] }
          ]
        }
      }
    }
  ],

  "explore": {
    "clickTabs": true,
    "clickTables": true,
    "pushStateFallback": true,
    "maxMenuItems": 50
  },

  "routes": ["/dashboard", "/alert", "/asset", "/device", "/report", "/config", "/system"],
  "waitMs": 1500, "perRouteMs": 900, "headless": true,
  "waitUntil": "domcontentloaded",
  "routeTimeout": 12000,
  "proxy": "",

  "captureResponses": true, "recordWs": true, "respMax": 600
}
```

필드 설명:
- `runtimeMode`: `depth`(Puppeteer), `coverage`(browser MCP), `both`
- `cookies[].value` 접두: `b64json:` → base64(JSON), `json:` → 원본 JSON, 접두 없음 → 리터럴
- `forward: true`는 실제 요청을 전달하고 코드 필드를 변경하며, `false`는 완전히 오프라인인 stub
- `mockTier`: coverage 모드의 preload에서 활성화할 계층. 예: `L1+L2`, `L1+L2+L3`
- `routes`는 `routes.txt`에서 가져옵니다. 메뉴를 forge한 뒤 harness가 `<a href>`를 자동으로 추가합니다
- `captureResponses` / `recordWs`는 depth 모드에서만 유효합니다
- `waitUntil`: 대규모 SPA에는 `domcontentloaded`를 사용하고 `networkidle2`로 대기가 끝나지 않는 상황을 피합니다
- `routeTimeout`: 단일 라우트의 `page.goto` 시간 초과(밀리초)
- `proxy`: Puppeteer `--proxy-server`. `HTTP_PROXY` / `HTTPS_PROXY`를 설정할 수도 있습니다

### B1. 두 stub 템플릿(role_permissions + permissions/all)

```json
"stubs": [
  {
    "match": "role_permissions",
    "body": {
      "response_code": 0,
      "data": {
        "permissions": ["MONITOR", "MONITOR_ALERT", "THREAT", "ASSETS_RISK"],
        "role_type": "SUPER_ADMIN"
      }
    }
  },
  {
    "match": "permissions/all",
    "body": {
      "response_code": 0,
      "data": [
        {
          "code": "MONITOR",
          "position": 1,
          "children": [
            { "code": "MONITOR_ALERT", "position": 1, "children": [] }
          ]
        }
      ]
    }
  }
]
```

바깥 필드명(`response_code` / `code` / `data`)은 반드시 A2 인터셉터 관문과 일치해야 하고, `permissions`는 반드시 tree의 모든 leaf code를 포함해야 합니다.

---

## C. coverage 모드: preload 설정

`scripts/preload.js` 위쪽의 `CONFIG` 객체를 편집하거나 CDP를 통한 추가 전에 교체합니다:

```javascript
const CONFIG = {
  loginPathRe: /\/(login|signin)(\/|$|\?)/i,
  mockTier: 'L1+L2',
  forward: true,
  recordDetail: true,
  extractUrlsFromResponse: true,
  neutralizeVueRouter: true,
  observe: { storageReads: false, cookieReads: false, xhrHeaders: true },
  neutralize: { fields: ['response_code', 'code'], success: 0 },
  stubs: [ /* config.json stubs와 동일 */ ],
  apiPattern: /\/(api|apis|v\d+|dev|internal|graphql)\//i,
};
```

검증: `window.__API_RECON_PRELOAD__ === true`이고 pathname이 안정적으로 유지돼야 합니다.

레코딩 결과 내보내기:

```javascript
JSON.stringify({
  apis: [...window.__API_RECON_LOG__],
  detail: window.__API_RECON_DETAIL__,
  routes: [...(window.__API_RECON_ROUTES__ || [])],
  observe: window.__API_RECON_OBSERVE__,
}, null, 2)
```

---

## D. preload / runtime Hook 기능

preload(coverage)와 runtime_harvest(depth)에 내장된 브라우저 Hook 기능과 적용 범위:

| Hook 기능 | API 발견에 주는 가치 | 지원 범위 |
|---|---|---|
| Hook fetch / XHR.open | 요청 URL/메서드 기록 | ✅ `recordDetail` + `__API_RECON_LOG__` |
| Hook XHR.setRequestHeader | Authorization 등의 헤더 발견 | ✅ `observe.xhrHeaders` |
| Hook localStorage/cookie 읽기 | 세션 키 이름 확인 | ⚠️ 선택 사항 `observe.storageReads/cookieReads` |
| Vue 라우트 가져오기 | frontendRoutes 보완 | ✅ `__API_RECON_ROUTES__`(로드된 라우트) |
| Vue 라우트 인증 가드 무력화 / 로그인 이동 차단 | 모듈을 표시해 API 트리거 | ✅ `neutralizeVueRouter` + 원래 이동 기능 무력화 |
| React 라우트 가져오기 | 라우트 보완 | ⚠️ 정적 분석 + 클릭. 전용 Hook 없음 |
| 페이지 이동 차단(로그인 path) | 해당 페이지에서 분석 유지 | ⚠️ 로그인 path만 차단하여 애플리케이션 탐색을 막지 않도록 함 |
| Hook 암호화 라이브러리(CryptoJS/SM 등) | 암호화 파라미터 → 평문 API body | ❌ 반드시 직접 Hook으로 암호화 함수의 입력 파라미터를 관찰해야 함. 결론은 config에 기록 |
| 안티 디버깅 bypass | 그렇지 않으면 runtime에서 API를 기록할 수 없음 | ❌ 반드시 직접 처리해야 함. 정적 분석은 여전히 사용 가능 |

---

## E. Endpoint 추출 정규 표현식(정적 결과가 너무 적을 때)

`harvest_static.py`의 `extract_endpoints`를 넓히거나 직접 실행합니다:

```bash
grep -rhoaE '"/[a-z][A-Za-z0-9_/\-]{3,}"' js | sort -u
grep -rhoaE '/api/[a-zA-Z0-9_./-]+' js | sort -u
```

---

## F. 문제 해결

| 현상 | 원인 → 처리 |
|---|---|
| 정적 API가 매우 적음 | endpoint 표현 방식이 맞지 않음 → 정규 표현식을 넓힘(E절) |
| chunk 개수 ≪ manifest | CSS-only 또는 배포되지 않은 chunk. 404는 이미 재시도함 |
| runtime에서 여전히 로그인 화면이 표시됨 | 렌더링 관문 오류 → A1 재확인: 키 이름, 저장 위치, 인코딩, domain |
| 로그인 후 화면에 진입했으나 모듈이 비어 있음 | 콘텐츠 관문 → 메뉴 forge(A3). `routes` path가 잘못됐을 수 있음 |
| 각 라우트에 bootstrap/locale만 있음 | 권한 코드가 불완전함 → I절 권한 트리 복원. `role_permissions` + `permissions/all`의 두 stub 확인 |
| 사이드바에는 항목이 있으나 하위 화면이 비어 있음 | tree에 intermediate 노드가 없거나 code가 `userRouteAuth`와 일치하지 않음 |
| 모든 API가 로그인으로 이동함 | 인터셉터 관문 → `neutralize` 확인. 중첩 필드는 walk 로직 확장이 필요함 |
| WS 프레임이 0임 | 사용자 상호작용 후에만 subscribe함. `perRouteMs`를 늘림 |
| 응답 본문이 비어 있음 | `forward: true`일 때만 실제 응답이 있음 |
| Chromium이 없음 | chromium 설치 또는 `config.chromium` / `CHROMIUM` 설정 |
| Mock이 많으나 여전히 로그인으로 돌아감 | Hook이 너무 늦거나 `location.href` setter가 없음 → document-start + preload |
| 목록이 모두 비어 있음 | L3의 빈 배열은 정상. Tab/설정/상세를 계속 클릭 |
| Redux action을 라우트로 잘못 취급함 | get/set/change/clear/toggle/upload가 포함된 내부 path를 필터링 |
| Vue가 여전히 로그인으로 이동함 | preload가 document-start가 아님 → 추가 시점 변경. 또는 `neutralizeVueRouter: false`이면 인증 가드를 직접 제거 |
| 응답에 URL이 있으나 log에 기록되지 않음 | `extractUrlsFromResponse` 활성화 또는 `__API_RECON_DETAIL__`에서 직접 추출 |
| Authorization 헤더 이름을 모름 | `observe.xhrHeaders` 활성화 또는 DevTools에서 요청 헤더 확인 |
| runtime이 매우 느림 / 시간 초과 | `waitUntil: domcontentloaded`로 변경. `routeTimeout`을 낮춤. `networkidle2`를 사용하지 않음 |
| 프록시 연결 실패 | `proxy` / 환경 변수 확인. Puppeteer와 curl 프록시 포트가 일치해야 함 |

---

## G. hardened 대상

서버에서 세션을 단계적으로 검증하는 경우(forge할 수 없는 서명된 cookie, 서버에서 렌더링하며 stub 처리할 수 없는 메뉴), runtime이 shell에서 막힙니다. 예상 동작:

- **정적 분석으로 endpoint 열거를 충분히 수행할 수 있습니다** — 모듈 path는 코드에 있습니다
- 허용된 범위라면 **실제 세션**으로 같은 harness를 실행합니다: `forward: true`, neutralize 불필요, 실제 methods/params/responses 캡처

---

## H. 단일 작업 체크리스트

1. 허용된 범위 확인
2. **읽기** `scripts/harvest_static.py` → 대상에 맞게 조정 → 실행 → `api_static.txt`, `routes.txt` 검토
3. **Phase 1b**: path 앵커 주변 검색 범위 확대 + 바인딩 계층 → `param_candidates.json`(J절)
4. A1/A2/A3 역분석 → 사이트 전용 `config.json` 작성
5. **읽고 조정** `runtime_harvest.js` / `preload.js`, 그다음 실행
6. `runtimeMode=depth`: `npm install` → 조정한 harvest 스크립트 실행
7. `runtimeMode=coverage/both`: document-start에 조정한 preload 추가 → browser MCP 동적 열거 + **파라미터 트리거 매트릭스**
8. 모듈이 렌더링되지 않음 → **I절 권한 트리 복원** → patch stubs → 다시 실행
9. 파라미터의 여러 샘플 diff + 오류 기반 역추론 → `params_merged.json`
10. 병합 → `site_map.json` + `api_merged.txt`, 다룬 범위·누락·스크립트 변경점을 사실대로 표시

---

## I. 권한 트리 복원(Phase 4 심화)

단순한 형태로 forge한 `menus: [{ path, show: true }]`가 효과 없고 하위 모듈이 여전히 mount되지 않을 때 사용합니다.

### I1. auth 모듈 찾기

```bash
grep -l 'userRouteAuth' js/*.js
grep -l 'routeMap\|routeLink' js/*.js
grep -rhoaE 'getResultTree|role_permissions|permissions/all' js | head
```

기록: **권한 API path**, **응답 필드명**, **소비 chunk 파일명**.

### I2. routeMap 추출

```bash
python3 scripts/extract_route_map.py recon/js recon/
# 산출물 recon/route_map.json
```

`[!] no routeMap pattern found`이면: `extract_route_map.py`의 정규 표현식을 넓히거나 직접 grep합니다:

```bash
grep -rhoaE '([A-Z_][A-Z0-9_]*):\{name:"[^"]*",link:"/[^"]+"\}' js | head -20
```

### I3. 권한 트리 + stub 구성

```bash
python3 scripts/build_perm_tree.py recon/js recon/ --config recon/config.json
```

스크립트 로직:
1. `userRouteAuth={MONITOR:{url:...},...}` 파싱(webpack 별칭 `He=o.DASHBOARD` 포함)
2. `route_map.json`으로 alias → 실제 path 해석
3. code 접두로 parent 추론(`MONITOR_ALERT` → `MONITOR`)
4. `permissions_tree.json`, `permissions_all_stub.json`, `role_permissions_stub.json` 출력
5. `--config`로 지정한 `config.json`이 존재하면 `stubs`를 자동 갱신

**대상에 맞게 조정**(스크립트 위쪽):
- `DEFAULT_ROOTS`: 최상위 모듈 code 목록
- `DEFAULT_PREFIX_PARENT`: `PREFIX_` → parent 매핑
- `DEFAULT_EXTRA_PARENT`: 접두 관계가 아닌 orphan 노드

### I4. stub 일관성 확인

```bash
# permissions 개수 ≈ userRouteAuth 항목 개수여야 함
wc -l recon/perm_codes_all.txt
# routes는 route_map의 모든 link를 포함해야 함
python3 -c "import json; m=json.load(open('recon/route_map.json')); r=set(json.load(open('recon/config.json'))['routes']); print('missing', [v['link'] for v in m.values() if v['link'] not in r])"
```

### I5. runtime 다시 실행 및 비교

```bash
node recon/runtime_harvest.js recon/config.json
# forge 전후 runtime_api.json 개수를 비교. /attack, /asset 등에 모듈 API가 나타나는지 확인
```

| forge 전 | forge 후(성공) |
|---|---|
| 각 라우트에 동일한 3–5개 bootstrap | 라우트마다 다른 module API가 트리거됨 |
| `/api/locale/language`만 있음 | `/api/web/...` 모듈 endpoint가 나타남 |
| `routes.txt`의 라우트가 한 자릿수 | `routes` 80–110+개는 route_map에서 가져옴 |

### I6. 여전히 실패하는 경우

- **coverage 모드**: 사이드바 + Tab 클릭. 권한 gating 요청은 상호작용 후에만 발생할 수 있습니다
- **stub 필드**: 실제 API(curl + 실제 session)와 stub의 nesting 비교
- **추가 인증 가드**: grep으로 `hasPermission|checkRole|func.` 등의 버튼 단위 검사를 찾고 `role_permissions.permissions` 확장
- **정적 분석으로 대체**: 모듈 API path는 여전히 `api_static.txt`에 있고 runtime은 METHOD/body만 보완합니다. 파라미터는 `param_candidates.json` + 기록한 샘플 유지

---

## J. 파라미터 역분석(Phase 1b / 5b / 5c)

**방법론이며 범용 스크립트가 아닙니다.** path는 정규 표현식으로 찾고, 파라미터는 앵커 주변 검색 범위 확대 + UI 바인딩 연결 관계 + 여러 샘플의 diff + 오류 기반 역추론으로 확인합니다.

### J1. 앵커 주변 검색 범위 확대 — path에서 요청 구성 객체 찾기

```bash
# Phase 1의 알려진 path를 앵커로 사용
grep -n '"/api/user/list"' js/*.js
grep -rhoaE '.{0,120}("/api[^"]+").{0,200}' js | head
grep -rhoaE '(params|data|body|payload)\s*:\s*\{' js | head
grep -rhoaE '(get|post|put|delete|patch)\([^,]+,\s*\{' js | head
```

### J2. 래퍼 계층과 전송 형태

```bash
# axios / 공통 request
grep -rhoaE '(axios|request)\.(get|post|put|delete|patch)\(' js | head
grep -rhoaE 'interceptors\.(request|response)' js | head

# GraphQL
grep -rhoaE '(query|mutation)\s+\w+|gql`|graphql\(' js | head
grep -rhoaE '\$[a-zA-Z_]+\s*:\s*(Int|String|Boolean|\[)' js | head

# FormData / multipart
grep -rhoaE 'FormData|\.append\(' js | head

# 경로 파라미터
grep -rhoaE 'path:\s*"/[^"]*:[^"]+"' js | head
grep -rhoaE 'useParams|route\.params|\$route\.params' js | head
```

### J3. 검증 관문 — 필수 / 형식 / 열거

```bash
grep -rhoaE '(required|message|pattern|enum|validator)\s*:' js | head
grep -rhoaE 'yup\.|zod\.|async-validator|Form\.Item|a-form-item|el-form-item' js | head
grep -rhoaE 'rules\s*:\s*\[|name:\s*["\'][a-zA-Z_]+["\']' js | head
grep -rhoaE 'label.*value|options\s*:\s*\[' js | head
```

### J4. 바인딩 계층 — 폼 → API

```bash
grep -rhoaE 'onFinish|handleSubmit|getFieldsValue|validateFields' js | head
grep -rhoaE '(pick|omit|transform|dayjs|moment)\(' js | head
```

runtime에서 보완합니다: DevTools → Network → 요청 → **Initiator(요청 시작 지점)**(call stack)에서 `fetch`/`send`부터 거슬러 올라가 요청 구성 함수를 추적합니다.

### J5. 암호화 파라미터

```bash
grep -rhoaE 'encrypt|decrypt|sign|CryptoJS|sm2|sm3|sm4|RSA|AES' js | head
```

**암호문에서 필드를 추측하지 마세요** — Hook으로 암호화 함수의 **입력 파라미터**를 관찰하고, 암호화 전에 plaintext payload를 기록합니다. 결론은 `config.json` / `param_candidates.json`에 기록합니다.

### J6. 파라미터 트리거 매트릭스(Phase 3 필수)

각 모듈에서 동작별로 한 번씩 기록하고 diff로 요청 body/query를 대조합니다:

| 동작 | 확인할 내용 |
|---|---|
| 목록 첫 화면 | 페이지네이션 기본값 |
| 검색 | keyword, filters |
| 고급 필터링 | optional 필드 |
| 생성/편집 | 전체 entity |
| 일괄/내보내기 | `ids[]`, `exportType` |
| 정렬/페이지 이동 | `sortField`, `order` |

산출물 `param_samples.json`: `[{ "path", "method", "action": "search", "body", "query", "headers" }]`

### J7. 신뢰도 규칙

| 신뢰도 | 조건 |
|---|---|
| **높음** | 정적 callsite + runtime ≥2개 샘플 일치 |
| **중간** | 정적 결과만 있거나 runtime 1회만 있음 |
| **낮음** | 응답/오류 기반 역추론, 추가 검증하지 않음 |
| **트리거 대기** | 정적으로 알려진 필드이나 UI/권한 경로에 도달하지 못함 |

### J8. 상황별 빠른 설정

| 상황 | 순서 |
|---|---|
| REST 목록 페이지 | J1 요청 구성 객체 → J6 네 번의 diff → J3 rules |
| 생성/편집 폼 | J3 Form name → J4 submit 연결 관계 → runtime 제출 + 일부러 비워 두고 400 확인 |
| GraphQL | J2 variables 선언 → runtime에서 각 operation의 variables 기록 |
| 암호화 body | J5 Hook 입력 파라미터 → 암호화 전 필드가 실제 params |

### J9. api-recon 단계와 대응

| api-recon | 파라미터 recon |
|---|---|
| Phase 1 정적 | J1 앵커 주변 검색 범위 확대 |
| Phase 2 A2 인터셉터 | 전역 추가 필드(tenantId, sign) |
| Phase 3 runtime | J6 트리거 매트릭스 + `param_samples.json` |
| Phase 4 권한 트리 | 모듈마다 폼이 다름 → 권한이 충분해야 모든 필드를 트리거 |
| Phase 5 병합 | `params_merged.json` + 신뢰도. 단일 샘플로 필수 여부를 확정하지 마세요 |

### J10. 문제 해결

| 현상 | 처리 |
|---|---|
| 정적으로 필드명은 있으나 runtime에서 나타난 적 없음 | '트리거 대기'로 표시. 권한 트리 보완 / 고급 필터 클릭 / 연동 select의 각 option 선택 |
| 같은 path에서 body 형태가 다름 | 정상 — `action`별로 나누어 기록하고 schema를 강제로 병합하지 마세요 |
| stub 응답은 가짜지만 params를 확인하려 함 | **outbound 요청** body/headers를 확인합니다. stub 응답에서 역추론하지 마세요 |
| 400에서 nested field를 알림 | 바깥쪽 `data`/`bizData`/`variables` 구조에 유의합니다 |
| GraphQL에서 operation 이름만 보임 | `variables` JSON을 펼치고 정적 분석에서 `$var: Type`을 찾습니다 |

---
