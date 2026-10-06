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

## C. coverage 模式：preload 配置

编辑 `scripts/preload.js` 顶部 `CONFIG` 对象，或通过 CDP 注入前替换：

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
  stubs: [ /* 同 config.json stubs */ ],
  apiPattern: /\/(api|apis|v\d+|dev|internal|graphql)\//i,
};
```

验证：`window.__API_RECON_PRELOAD__ === true` 且 pathname 稳定。

导出录制结果：

```javascript
JSON.stringify({
  apis: [...window.__API_RECON_LOG__],
  detail: window.__API_RECON_DETAIL__,
  routes: [...(window.__API_RECON_ROUTES__ || [])],
  observe: window.__API_RECON_OBSERVE__,
}, null, 2)
```

---

## D. preload / runtime Hook 能力

preload（coverage）与 runtime_harvest（depth）内置的浏览器 Hook 能力及覆盖范围：

| Hook 能力 | 对 API 发现的价值 | 覆盖 |
|---|---|---|
| Hook fetch / XHR.open | 录请求 URL/方法 | ✅ `recordDetail` + `__API_RECON_LOG__` |
| Hook XHR.setRequestHeader | 发现 Authorization 等头 | ✅ `observe.xhrHeaders` |
| Hook localStorage/cookie 读 | 确认会话键名 | ⚠️ 可选 `observe.storageReads/cookieReads` |
| Vue 获取路由 | 补全 frontendRoutes | ✅ `__API_RECON_ROUTES__`（已加载路由） |
| Vue 路由守卫中和 / 登录跳转阻断 | 撑开模块触发 API | ✅ `neutralizeVueRouter` + 原生跳转中和 |
| React 获取路由 | 补路由 | ⚠️ 静态 + 点击；无专用 Hook |
| 页面跳转阻断（登录 path） | 留页分析 | ⚠️ 仅阻断登录 path，避免挡业务导航 |
| Hook 加密库（CryptoJS/SM 等） | 加密参数 → 明文 API body | ❌ 须手工 Hook 加密函数入参；结论写 config |
| 反调试 bypass | 否则 runtime 录不到 API | ❌ 须手工处理；静态仍可用 |

---

## E. Endpoint 추출 정규 표현식(정적 결과가 너무 적을 때)

`harvest_static.py`의 `extract_endpoints`를 넓히거나 직접 실행합니다:

```bash
grep -rhoaE '"/[a-z][A-Za-z0-9_/\-]{3,}"' js | sort -u
grep -rhoaE '/api/[a-zA-Z0-9_./-]+' js | sort -u
```

---

## F. 排障

| 现象 | 原因 → 处理 |
|---|---|
| 静态 API 很少 | endpoint 方言不匹配 → 放宽正则（D 节） |
| chunk 数 ≪ manifest | CSS-only 或未部署 chunk；404 已重试 |
| runtime 仍显示登录页 | 渲染门错误 → 复查 A1：键名、容器、编码、domain |
| 进壳但模块空白 | 内容门 → forge 菜单（A3）；`routes` path 可能不对 |
| 每路由只有 bootstrap/locale | 权限码不全 → I 节权限树还原；检查 `role_permissions` + `permissions/all` 双 stub |
| 侧栏有项但子页空白 | tree 缺 intermediate 节点或 code 与 `userRouteAuth` 不一致 |
| 每个 API 都跳登录 | 拦截器门 → 确认 `neutralize`；嵌套字段需扩展 walk 逻辑 |
| WS 帧为 0 | 需用户交互后才 subscribe；加长 `perRouteMs` |
| 响应体空 | 仅 `forward: true` 时有真实响应 |
| Chromium 缺失 | 安装 chromium 或设置 `config.chromium` / `CHROMIUM` |
| Mock 很多仍回登录 | Hook 太晚或缺 `location.href` setter → document-start + preload |
| 列表全空 | L3 空数组正常；继续点 Tab/设置/详情 |
| 误把 Redux action 当路由 | 过滤含 get/set/change/clear/toggle/upload 的内部 path |
| Vue 仍跳登录 | preload 非 document-start → 改注入时机；或 `neutralizeVueRouter: false` 时手动清守卫 |
| 响应里有 URL 但未进 log | 开 `extractUrlsFromResponse`；或从 `__API_RECON_DETAIL__` 人工提取 |
| 不知 Authorization 头名 | 开 `observe.xhrHeaders` 或 DevTools 查看请求头 |
| runtime 极慢 / 超时 | 改 `waitUntil: domcontentloaded`；降 `routeTimeout`；勿用 `networkidle2` |
| 代理连接失败 | 检查 `proxy` / 环境变量；Puppeteer 与 curl 代理端口一致 |

---

## G. hardened 目标

服务端逐步校验会话（不可 forge 的签名 cookie、服务端渲染且不可 stub 的菜单）时，runtime 会在 shell 处卡住。预期行为：

- **静态足够做 endpoint 枚举** — 模块 path 在代码里
- 若授权允许，用**真实会话**跑同一 harness：`forward: true`、无需 neutralize，捕获真实 methods/params/responses

---

## H. 单次任务清单

1. 确认授权范围
2. **阅读** `scripts/harvest_static.py` → 按目标调整 → 运行 → 审 `api_static.txt`、`routes.txt`
3. **Phase 1b**：path 锚点扩窗 + 绑定层 → `param_candidates.json`（J 节）
4. 逆向 A1/A2/A3 → 写站点专属 `config.json`
5. **阅读并调整** `runtime_harvest.js` / `preload.js` 后再执行
6. `runtimeMode=depth`：`npm install` → 运行调整后的 harvest 脚本
7. `runtimeMode=coverage/both`：document-start 注入调整后的 preload → browser MCP 动态枚举 + **参数触发矩阵**
8. 模块不渲染 → **I 节权限树还原** → patch stubs → 重跑
9. 参数多样本 diff + 错误反推 → `params_merged.json`
10. 合并 → `site_map.json` + `api_merged.txt`，诚实标注覆盖、缺口及脚本改动点

---

## I. 权限树还原（Phase 4 深化）

当 forge 简单 `menus: [{ path, show: true }]` 无效、子模块仍不 mount 时使用。

### I1. 定位 auth 模块

```bash
grep -l 'userRouteAuth' js/*.js
grep -l 'routeMap\|routeLink' js/*.js
grep -rhoaE 'getResultTree|role_permissions|permissions/all' js | head
```

记录：**权限 API path**、**响应字段名**、**消费 chunk 文件名**。

### I2. 提取 routeMap

```bash
python3 scripts/extract_route_map.py recon/js recon/
# 产出 recon/route_map.json
```

若 `[!] no routeMap pattern found`：放宽 `extract_route_map.py` 中正则，或手工 grep：

```bash
grep -rhoaE '([A-Z_][A-Z0-9_]*):\{name:"[^"]*",link:"/[^"]+"\}' js | head -20
```

### I3. 构建权限树 + stub

```bash
python3 scripts/build_perm_tree.py recon/js recon/ --config recon/config.json
```

脚本逻辑：
1. 解析 `userRouteAuth={MONITOR:{url:...},...}`（含 webpack 别名 `He=o.DASHBOARD`）
2. 用 `route_map.json` 解析 alias → 真实 path
3. 按 code 前缀推断 parent（`MONITOR_ALERT` → `MONITOR`）
4. 输出 `permissions_tree.json`、`permissions_all_stub.json`、`role_permissions_stub.json`
5. `--config` 时自动写入 `config.json` 的 `stubs` 与扩展 `routes`

**按目标调整**（在脚本顶部）：
- `DEFAULT_ROOTS`：顶级模块 code 列表
- `DEFAULT_PREFIX_PARENT`：`PREFIX_` → parent 映射
- `DEFAULT_EXTRA_PARENT`：非前缀关系的 orphan 节点

### I4. 校验 stub 一致性

```bash
# permissions 数量应 ≈ userRouteAuth 条目数
wc -l recon/perm_codes_all.txt
# routes 应覆盖 route_map 全部 link
python3 -c "import json; m=json.load(open('recon/route_map.json')); r=set(json.load(open('recon/config.json'))['routes']); print('missing', [v['link'] for v in m.values() if v['link'] not in r])"
```

### I5. 重跑 runtime 并对比

```bash
node recon/runtime_harvest.js recon/config.json
# 对比 forge 前后 runtime_api.json 条数；检查 /attack、/asset 等是否出现模块 API
```

| forge 前 | forge 后（成功） |
|---|---|
| 每路由相同 3–5 条 bootstrap | 不同路由触发不同 module API |
| 仅 `/api/locale/language` | 出现 `/api/web/...` 模块 endpoint |
| `routes.txt` 个位路由 | `routes` 80–110+ 来自 route_map |

### I6. 仍失败时

- **coverage 模式**：点击侧栏 + Tab，权限 gating 可能在交互后才请求
- **stub 字段**：对比真实 API（curl + 真实 session）与 stub 的 nesting
- **额外守卫**：grep `hasPermission|checkRole|func.` 等按钮级检查，扩展 `role_permissions.permissions`
- **静态兜底**：模块 API path 仍在 `api_static.txt`，runtime 仅补 METHOD/body；参数保留 `param_candidates.json` + 已录样本

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
