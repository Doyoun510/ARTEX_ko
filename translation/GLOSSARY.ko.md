# ARTEX 한국어 번역 용어집 · 번역 제외 목록 (통합본)

> 팀원 용어집(수정본)과 코드 분석 결과를 합친 **공용 기준 문서**입니다. UI → 에이전트 응답 → 내부 프롬프트·문서 전 단계에서 이 파일 하나만 공유합니다.
> 번역 담당자와 번역 프롬프트 모두 이 파일을 기준으로 삼습니다.
> 기준 코드: Autumn-27/ARTEX (얕은 클론 분석). 코드는 바꾸지 않고 **중국어 주석·UI 문구·에이전트 응답·프롬프트·문서의 언어만** 한국어로 바꿉니다.

**표시 규칙**
- `[확정]` 팀이 결정한 항목입니다.
- `[제안]` 코드 분석 중 추가한 항목으로, 팀 검토 전입니다.
- 경로와 줄 번호는 분석 시점의 코드 기준이며, 코드가 바뀌면 달라질 수 있습니다.

---

## 확정 결정 요약

| # | 항목 | 결정 |
| --- | --- | --- |
| 1 | human-in-the-loop | 영문 그대로 유지 (README.ko.md의 "휴먼 인 더 루프"는 영문으로 수정) |
| 2 | 拦截 | 기능명 = **인터셉트**, 판정(동작·결과) = **차단**, 대기 중 = **승인 대기** |
| 3 | 已修复 | **수정 완료** (화면 라벨만. 내부 enum `fixed`는 DNT) |
| 4 | 公司 / 企业 | 모두 **회사**. 예외: `企业微信`(제품명), `企业官网` |
| 5 | deny 용어 | 사용자·시스템이 승인 요청을 거절 = **거부**, 인터셉트가 막음 = **차단** (6장 "deny 맥락별 표기" 참고) |
| 6 | 「」 | 작은따옴표 `'…'`로 통일 (README 포함) |
| 7 | 方式 N 제목 | **방법 N** |
| 8 | DB | **새 DB로 시작** (기존 DB 데이터 없음) |

**새 DB 전제**: 에이전트 프롬프트와 내장 인터셉트 규칙은 **DB가 비어 있을 때 한 번만** 저장됩니다(`db/config.go:743` `SeedPromptIfEmpty`, 내장 규칙은 `intercept_default_rules_v1/v2/v3` 등 settings 플래그가 "done"이면 건너뜀). 번역을 고친 뒤에는 **DB를 새로 만들어야** 반영됩니다. 개발 중에는 DB를 지우고 재시작하거나, 화면의 "내장 기본값으로 복원" 기능(`ResetPromptToDefault`, `db/config.go:757`)을 쓰세요.

---

## 0. 최상위 원칙 (반드시 먼저 읽기)

1. **번역 금지(DNT)가 번역보다 우선** — 어떤 문자열을 만나면 "번역할까?"보다 **"이게 로직에 쓰이나?"를 먼저** 판단한다. 비교(`==`/`switch`/`if`), 파싱(`split`/정규식 앵커/`strings.Cut`/`HasPrefix`), 저장(DB 값·직렬화 키), 백엔드↔프런트 문자열 계약, 에이전트 출력 계약(도구명·파라미터·enum)에 쓰이면 → **번역 금지 + 플래그**. ([5. 번역 금지(DNT) 목록](#5-번역-금지dnt-목록) 참고)
2. **중국어만 번역, 영어는 원문 유지** — 코드·명령·식별자에 섞인 영어는 그대로 둔다.
3. **자리표시자 보존** — `%s` `%d` `%v` `{count}` 등과 그 순서, `\n`, 마크다운(`**`, 백틱), HTML 태그는 보존.
4. **같은 개념 = 같은 번역** — 이 표에 있으면 무조건 표의 번역을 쓴다.
5. **애매하면 번역하지 말고 사람에게** — 틀린 번역보다 미번역 + 사유 플래그가 안전.
6. **코드 불변** — 문자열 리터럴 안의 언어만 바꾼다. 식별자·로직·구조는 건드리지 않는다. 비교값을 바꿔야만 하는 경우는 [5.7](#57-값-계약-문자열-중국어-값-자체가-계약인-것)·[5.8](#58-같은-파일-안의-쌍-둘-다-같게-바꾸면-허용)을 따른다.
7. **전각 문장부호도 대상** — `，：；、「」【】（）《》`는 중국어 범위 검색(`一-鿿`)에 안 잡힌다. 검색할 때 전각 문장부호 범위도 포함한다.

---

## 1. 핵심 도메인·아키텍처 용어

| 원문(zh) | 영문 canonical | 한국어 | 비고 |
| --- | --- | --- | --- |
| 渗透测试 | penetration testing | 침투 테스트 | |
| 自主 | autonomous | 자율 | |
| 后端 / 前端 | backend / frontend | 백엔드 / 프런트엔드 | |
| 单体后端 | monolithic backend | 모놀리식 백엔드 | |
| 内嵌 | embed | 임베드 | `go:embed` 식별자는 DNT |
| 单二进制 | single binary | 단일 바이너리 | |
| 跨平台 | cross-platform | 크로스 플랫폼 | |
| 静态导出 | static export | 정적 내보내기 | |
| 双图架构 | dual-graph architecture | 이중 그래프 아키텍처 | |
| 资产图 | Asset Graph | 자산 그래프 | |
| 探索图 | Exploration Graph | 탐색 그래프 | `探索图` ≠ `探索链路` |
| 探索链路 | exploration chain | 탐색 체인 | |
| 探索 | exploration | 탐색 | |
| 锚点 | anchor | 앵커 | 엣지명 `anchor`는 DNT |
| 血缘链 | lineage chain | 계보 체인 | "혈통" 쓰지 말 것 |
| 覆盖图 | coverage map | 커버리지 맵 | |
| 覆盖度 | coverage | 커버리지 | 식별자 `coverage`는 DNT |
| 真值库 | single source of truth (SSOT) | 단일 신뢰 소스(SSOT) | README 동반 수정(8장) |
| 攻击链(路) | attack chain | 공격 체인 | |
| 利用链 | exploit chain | 익스플로잇 체인 | |
| 过程级信息交换 | process-level info exchange | 프로세스 수준 정보 교환 | |
| 人在环路 | human-in-the-loop | human-in-the-loop | **[확정]** 영문 유지. 원문은 중국어(`人在环路`)지만 관용적으로 영문 표기를 쓰는 용어이므로 영문을 그대로 쓴다(음역 금지) |

### 그래프 노드/개념 (프로즈에서의 번역 — 식별자 자체는 DNT)

| 원문(zh) | 영문(식별자, DNT) | 한국어(개념어) | 비고 |
| --- | --- | --- | --- |
| 目标 | `goal` | 목표 | 노드 kind 식별자 `goal`은 그대로. 테스트 대상(URL·호스트)을 가리키는 **目标**는 "대상" `[확정]` — 문맥으로 구분 |
| 意图 | `intent` | 의도 | 식별자 `intent`/`intent_id`는 DNT |
| 事实 | `fact` | 사실 | |
| 发现 / 漏洞 | `finding` / vulnerability | 발견 / 취약점 | `finding`은 "발견", `漏洞`는 "취약점"으로 구분 |
| 提示 | `hint` | 힌트 | `战略提示` 등도 "힌트". **提示词(프롬프트)와 혼동 금지** |
| 提示词 | prompt | 프롬프트 | `[확정]` 에이전트 시스템 프롬프트. 提示(힌트)와 구분 |
| 规划者 | planner | planner | 에이전트명은 영문 유지(6장 영문 유지 규칙). 개념 설명 문장에서도 `planner`로 통일 |
| 执行者 | worker | worker | 위와 동일 |
| 主 agent / 主 Agent | main agent | 메인 agent | `[확정]` `agent` 영문 유지 규칙에 맞춤. 식별자 `mainagent`는 DNT |

---

## 2. 기능·화면(UI) 용어

| 원문(zh) | 영문 | 한국어 | 비고 |
| --- | --- | --- | --- |
| 仪表盘 | dashboard | 대시보드 | |
| 任务 | task | 작업 | 전역 통일 |
| 资产 | asset | 자산 | |
| 流量 | traffic | 트래픽 | |
| 流量录制 | traffic recording | 트래픽 레코딩 | |
| 记录代理 / 流量记录代理 | recording proxy | (트래픽) 레코딩 프록시 | |
| 工作空间 | workspace | 워크스페이스 | |
| 系统配置 / 系统设置 | system settings | 시스템 설정 | |
| 配置 / 默认 / 内置 | config / default / built-in | 설정 / 기본(값) / 내장 | `[확정]` `[内置]` 규칙 이름 접두는 "[내장]" |
| 模型 | model | 모델 | `[확정]` |
| 工具 / 工具调用 | tool / tool call | 도구 / 도구 호출 | `[확정]` |
| 对话 | conversation | 대화 | `[확정]` 会话(세션)과 구분 |
| 审批 | approval | 승인 | |
| 审批记录 | approval records | 승인 기록 | |
| 待审批 | pending approval | 승인 대기 | **[확정]** 대기 중인 승인 요청 |
| 转人工 / 转人工审批 | escalate to human | 사람 승인으로 전환 | `[확정]` 판정 `ask` 라벨 |
| 审批超时 | approval timeout | 승인 시간 초과 | `[확정]` |
| 拦截(기능) / 拦截审批 | intercept / intercept approval | 인터셉트 / 인터셉트 승인 | **[확정]** 코드 `intercept`와 일치 |
| 拦截(동작/결과, 被拦截) | blocked | 차단 | **[확정]** 요청이 막힘(동사/결과). 기능명 "인터셉트"와 구분 |
| 放行 / 允许 | allow | 허용 | `[확정]` enum `allow`의 라벨 |
| 拒绝 / 禁止 / 阻断 | deny | 6장 "deny 맥락별 표기" | **[확정]** 맥락별로 거부·차단·금지를 구분 |
| 审批卡片 | approval card | 승인 카드 | |
| 资产同步 | asset sync | 자산 동기화 | |
| 后端日志 | backend log | 백엔드 로그 | |
| 活动流 | activity stream | 활동 스트림 | 식별자 `activity`는 DNT |
| 页签 | tab | 탭 | |
| 复测 | retest / re-verification | 재검증 | "재테스트" 금지, "재검증" 통일 |
| 漏洞复测 | vulnerability retest | 취약점 재검증 | |
| 处置状态 | disposition state | 처리 상태 | |
| 证据 | evidence | 증거 | 식별자 `evidence`는 DNT |
| 会话 | session | 세션 | |
| 运行预算 | run budget | 실행 예산 | |
| 演示模式 | demo mode | 데모 모드 | |
| 域名 / 子域(名) | domain / subdomain | 도메인 / 서브도메인 | 식별자 `root_domain`/`subdomain`은 DNT |
| 端口 / 站点 / 端点 | port / site / endpoint | 포트 / 사이트 / 엔드포인트 | |
| 接口 | endpoint (API) | 엔드포인트 | `[확정]` 에셋 종류 `endpoint`. "인터페이스"로 쓰지 말 것. 단 멘션 토큰 `@[接口#id]`는 DNT(5.7) |
| 公司 / 企业 | company | 회사 | **[확정]** 公司·企业 모두 "회사"(DB `companies`). 예외: `企业微信`=제품명, `企业官网`=기업 공식 사이트 |
| 企业微信 | WeCom | WeCom(기업용 위챗) | **[확정]** 처음 한 번 "WeCom(기업용 위챗)", 이후 "WeCom". 파일명 `notify/wecom.go`와 일치 |
| 公司资产范围 | company asset scope | 회사 자산 범위 | 식별자 `scope`/`task_scope`는 DNT |
| 范围 / 备案 | scope / ICP filing | 범위 / ICP 등록(备案) | `[확정]` 표시 문구만. 자동 분류에 쓰는 "备案"은 DNT(5.7) |
| 收尾 | wrap-up | 마무리(wrap-up) | `[확정]` 예산·타임아웃 시 실행되는 정리 라운드 |
| 旁路(提问) | side question | 보조 질문 | `[확정]` /btw 곁다리 질문(by the way). "우회 질의" 아님. 코드 `sidequestion`·`/btw`는 DNT |
| 渠道 / 投递 / 推送 | channel / delivery / push | 채널 / 전송 / 알림 전송 | `[확정]` |
| 摘要 / 报告 | summary / report | 요약 / 보고서 | `[확정]` |
| 上下文 / 压缩 / 快照 / 归档 | context / compaction / snapshot / archive | 컨텍스트 / 압축 / 스냅샷 / 아카이브 | `[확정]` |
| 轮询 / 超时 / 重试 | polling / timeout / retry | 폴링 / 타임아웃 / 재시도 | `[확정]` |
| 约束 / 命中 | constraint / hit(match) | 제약 조건 / 매칭 | `[확정]` 규칙이 걸린 경우 "매칭" |

---

## 3. 인프라·운영 용어

| 원문(zh) | 영문 | 한국어 | 비고 |
| --- | --- | --- | --- |
| 反向代理 | reverse proxy | 리버스 프록시 | |
| 同源 | same-origin | 동일 출처 | 괄호로 `(same-origin)` 병기 권장 |
| 缓冲 | buffering | 버퍼링 | |
| 长连接 | long-lived connection | 지속 연결 | 팀원안 "장기 연결"에서 변경 `[확정]` |
| 持续推送 | continuous push | 지속적 푸시 | |
| 冒烟测试 | smoke test | 스모크 테스트 | |
| 回滚 | rollback | 롤백 | |
| (数据库)迁移 | migration | 마이그레이션 | |
| 幂等 | idempotent | 멱등 | |
| 全局代理 | global proxy | 전역 프록시 | |
| 并发 | concurrency | 동시성 | |
| 调度循环 | schedule loop | 스케줄 루프 | |
| 生命周期 | lifecycle | 생명주기 | |
| 守护脚本 | daemon script | 데몬 스크립트 | |
| 工具审批门 | tool approval gate | 도구 승인 게이트 | |
| 异步补全 | async enrichment | 비동기 보강 | |
| 事件驱动 | event-driven | 이벤트 기반 | |
| 闭环 | closed loop | 폐루프 | "폐곡선"(기하학 용어)은 오역 — README 동반 수정(8장) |
| 唤醒 | wake | 깨우기 | |
| 无状态会话 | stateless session | 무상태 세션 | |
| 发布包 / Release 压缩包 | release package | 배포 패키지 | `Release`/`Releases` 고유명은 DNT |

### 공격 기법 용어

| 원문(zh) | 영문 | 한국어 | 비고 |
| --- | --- | --- | --- |
| 注入点 | injection point | 인젝션 포인트(injection point) | `[확정]` 업계 관용어(직역 "주입점" 대신) |
| 凭据 | credential | 자격 증명 | |
| 横向 | lateral movement | 횡적 이동 | |
| 提权 | privilege escalation | 권한 상승 | |

---

## 4. 상태·열거값

화면 라벨은 번역하고, **enum 값(저장·비교 레이어)은 영어 그대로**입니다. 코드로 확인했습니다: 저장값은 영어 enum이고 중국어는 표시 라벨뿐입니다(`web/src/lib/status.ts`, `db/findings.go:39-49`, `agent/retester.go`, `web/src/components/finding-retest-panel.tsx:28`).

**같은 enum의 라벨이 TS와 Go 두 곳에 있으면 두 곳을 동일하게 번역**합니다. 취약점 처리 상태는 `web/src/lib/status.ts`의 `finding`과 `notify/notify.go:63`의 `StatusLabel`에 같은 라벨이 중복되어 있습니다.

### 4.1 취약점 재검증 결론 (verdict enum: `reproduced` / `fixed` / `inconclusive`)

| enum (DNT) | 원문(zh) | 한국어(표시 라벨) |
| --- | --- | --- |
| `reproduced` | 仍可复现 | 여전히 재현 가능 |
| `fixed` | 已修复 | 수정 완료 **[확정]** |
| `inconclusive` | 无法确认 | 확인 불가 |
| (진행 표시) | 复测中 | 재검증 중 |
| — | 复现 | 재현 |

### 4.2 취약점 처리 상태 (`findings.status`)

| enum (DNT) | 원문(zh) | 한국어 `[확정]` |
| --- | --- | --- |
| `pending` | 待处理 | 처리 대기 |
| `in_progress` | 处理中 | 처리 중 |
| `confirmed` | 已确认 | 확인됨 |
| `resolved` | 已处理 | 처리됨 |
| `fixed` | 已修复 | 수정 완료 **[확정]** |
| `false_positive` | 误报 | 오탐 |
| `ignored` | 忽略 | 무시 |
| `duplicate` | 重复 | 중복 |
| `risk_accepted` | 风险接受 | 위험 수용 |

### 4.3 기타 상태 라벨 (`web/src/lib/status.ts`) `[확정]`

| 도메인 | enum (DNT) → 원문(zh) → 한국어 |
| --- | --- |
| 위험도 `severity` | `critical` 严重 → 심각 · `high` 高危 → 높음 · `medium` 中危 → 중간 · `low` 低危 → 낮음 |
| 의도 `intent` | `open` 待领 → 대기 · `running` 执行中 → 실행 중 · `paused` 已暂停 → 일시 중지됨 · `done` 已完成 → 완료 · `blocked` 执行出错 → 실행 오류 · `exhausted` 预算耗尽 → 예산 소진 · `stopped` 已停止 → 중지됨 · `deleted` 已删除 → 삭제됨 |
| 작업 `task` | `created` 已创建 → 생성됨 · `queued` 排队中 → 대기열 · `running` 运行中 → 실행 중 · `paused` 已暂停 → 일시 중지됨 · `done` 已完成 → 완료 · `failed` 失败 → 실패 · `timeout` 已超时 → 시간 초과 |
| 엔진 `engine` | `exploring` 探索中 → 탐색 중 · `paused` 已暂停 → 일시 중지됨 · `stalled` 停滞 → 정체 · `idle` 空闲 → 유휴 |
| 목표 `goal` | `open` 进行中 → 진행 중 · `met` 已达成 → 달성 · `abandoned` 已放弃 → 포기됨 |
| 감사 `audit` | `allow` 放行 → 허용 · `block` 拦截 → 차단 |
| 노드 `node` | `observed` 观测 → 관측 · `confirmed` 确认 → 확인 · `tombstoned` 废弃 → 폐기 |
| 전송 `delivery` | `pending` 待发送 → 전송 대기 · `sending` 发送中 → 전송 중 · `sent` 已送达 → 전달됨 · `failed` 失败 → 실패 · `skipped` 已跳过 → 건너뜀 |

> 주의: 의도 상태 `blocked`는 주석에 "모델/API/네트워크 장애로 재시도가 소진됨(목표 차단 아님)"이라고 적혀 있습니다. 인터셉트의 "차단"과 **다른 개념**이므로 "실행 오류"로 번역하고 "차단"을 쓰지 않습니다.

---

## 5. 번역 금지(DNT) 목록

아래는 **절대 번역하지 않고 원문 그대로** 둔다. (영어 식별자는 영어 그대로)

### 5.1 제품·고유명
`ARTEX` · `ScopeSentry` · `AegisHook` · `norma` · `SecSentry` · `Cairn` · `Kali` · `PostgreSQL` · `SQLite` · `Next.js` · `Go` · `Docker` · `Nginx` · `UPX` · `playwright` · `nmap` · `ripgrep` · `curl` · `vim` · `npm`

### 5.2 약어·프로토콜
`MITM` · `CA` · `JWT` · `SSE` · `MCP` · `HTTP(S)` · `REST` · `DNS` · `TCP` · `JSON-RPC` · `API` · `LLM` · `CWD` · `CIDR` · `ICP`

### 5.3 도구명·파라미터 (에이전트 출력 계약 — 번역 시 즉시 파손)
코드의 도구 정의(`agent/*.go`, `traffic/traffic.go`)에서 추출한 전체 목록입니다.

`prove_goal` · `add_intent` · `add_hint` · `add_task_hint` · `set_goals` · `set_constraints` · `goal_met` · `list_goals` · `graph_overview` · `get_task_graph` · `node_detail` · `get_task_node_detail` · `record_fact` · `list_facts` · `report_finding` · `update_finding_report` · `list_findings` · `list_task_findings` · `bind_finding_traffic` · `get_finding_traffic` · `get_finding_retest_context` · `record_finding_retest_result` · `search_all_worker_traces` · `search_task_worker_traces` · `list_worker_traces` · `list_task_worker_traces` · `get_worker_trace` · `get_task_worker_trace` · `get_worker_output` · `expand_digest` · `traffic_search` · `traffic_get` · `traffic_blob` · `insert_assets` · `list_assets` · `list_untested_assets` · `delete_assets_by_host` · `add_company_scope` · `add_task_scope` · `list_companies` · `list_tasks` · `list_llm_profiles` · `spawn_task` · `pause_task` · `kill_work` · `steer_work` · `create_skill` · `update_skill_file` · `create_custom_tool` · `update_custom_tool` · `create_mcp` · `update_mcp` · `TodoWrite`

SDK 공통 도구: `Read` · `Write` · `Edit` · `MultiEdit` · `LS` · `Glob` · `Grep` · `Bash` · `Skill`

파라미터/키: `intent_id` · `asset_id(s)` · `node_id` · `step_ids` · `finding_id` · `finding_node_id` · `evidence_hint_id` · `traffic_refs` · `evidence_version` · `version` · `role` · `company_id` · `verdict` · `summary` · `evidence`

> 이 목록은 정규식 추출이라 완전하지 않을 수 있습니다. 3단계(프롬프트 번역)에서 프롬프트에 나온 도구명을 이 목록과 대조해 누락을 보강합니다. `run_nuclei`·`run_shell`은 테스트 파일에만 있는 가짜 이름이라 제외했습니다.

### 5.3b 함수·타입·에이전트 식별자
`claimNext` · `plannerLoop` (`server/engine.go`, README 다이어그램에 등장) · `ToolSet` (Go 타입) · `retester` · `planner` · `worker` · `mainagent` · `goals` · `reporter` · `auto` · `pentest` (에이전트 key)

### 5.4 그래프 노드 kind·엣지·필드 (식별자)
노드/엣지: `goal` · `intent` · `fact` · `finding` · `hint` · `spawns` · `derived_from` · `yields` · `proves` · `anchor` · `frontier` · `activity`
DB/스키마: `assets` · `companies` · `task_scope` · `exploration_nodes` · `exploration_anchors` · `anchors` · `root_domain` · `subdomain` · `ip` · `service` · `app` · `endpoint`
role enum: `baseline` · `proof` · `verification` · `supporting`
인터셉트 판정 enum: `allow` · `ask` · `deny` · `block`

### 5.5 norma SDK 모듈
`agentcore` · `tool` · `permission` · `harness` · `memory` · `transcript`

### 5.6 환경변수·경로·명령·파일
`ANTHROPIC_API_KEY` · `OPENAI_API_KEY` · `ARTEX_PG_DSN` · `ARTEX_LLM_PROVIDER` · `ARTEX_LLM_MODEL` · `ARTEX_LLM_BASE_URL` · `ARTEX_LLM_PROXY` · `ARTEX_TAG` · `ARTEX_TARGETS` · `NEXT_PUBLIC_SSE_BASE` · `NEXT_PUBLIC_MOCK` · `CGO_ENABLED` · `embedui`
파일/경로: `config.json` · `config.example.json` · `.env` · `start.sh`/`start.bat` · `install.sh` · `update.sh` · `build.sh` · `dev.sh` · `schema.sql` · `jwt.key` · `pgdata` · `./data` · `./skills` · `artex(.new/.old/.failed)` · `SHA256SUMS`
모든 **명령·CLI 플래그·코드 블록 내 코드**(예: `docker compose up -d`, `go test ./...`, `-tags embedui`, `-addr`, `-proxy`, `--release`, `--upx`)는 DNT. (코드 블록 안 **중국어 주석만** 번역 대상)

### 5.7 값 계약 문자열 (중국어 값 자체가 계약인 것)

**영어 식별자가 아니라 중국어 문자열 값이 백엔드·프런트·파서 사이의 약속**인 곳입니다. 한쪽만 번역하면 조용히 깨집니다(오류 없이 기능이 사라짐). 모두 **번역 금지**입니다.

| 문자열 | 백엔드 위치 | 짝(읽는 쪽) | 깨지는 이유 |
| --- | --- | --- | --- |
| `[模型]` | `intercept/intercept.go:479-483`, `db/intercept.go:229-231`(SQL `LIKE '[模型]%'`)·`:251`, `db/task_archives_restore.go:680`, `db/schema.sql:1076` | `web/src/components/approval-records.tsx:53,112,373,487` (+ mock `handler.ts:2436`, `data.ts`) | 백엔드가 사유에 쓰고 SQL·프런트가 접두로 읽음 |
| `工具 %s 请求审批 (#%d)` | `intercept/intercept.go:612` | `web/src/components/transcript.tsx:201` 정규식 `/工具\s+(\S+)\s+请求/`, `/\(#(\d+)\)/` | 프런트가 정규식으로 파싱 |
| `实际操作：` / `；成功后的后果：` / `；命中规则：` | `intercept/prompt.go:36-39`(형식 설명), `:124-127`(예시) | 같은 파일 파서 `:193-200` (`HasPrefix`, `strings.Cut`) | 프롬프트만 번역하면 모든 판정이 무효(`Verdict{}`)가 되어 FailAction으로 떨어짐. 설명문·예시 본문은 번역하되 **이 세 구분자는 원문 그대로** |
| `不超过 120 个汉字` / `len(reason) > 2400` | `intercept/prompt.go`, `:193-200` | 같은 파일 | 길이 제한. 구분자는 아니지만 번역 후 한국어 길이 여유를 확인할 것 |
| `@[漏洞#id]` 등 멘션 토큰 (종류명 `漏洞`·`资产`·`企业`·`接口`·`IP`·`应用`·`域名`·`子域名`·`服务`) | `server/chat_mentions.go:20-24` (정규식 + 맵 키) | `web/src/lib/chat-mentions.ts:55` | 서버가 종류명으로 타입을 판별. **사용자에게 보이는 토큰**(`chat_mentions.go:18` 주석)이라 "코드 불변"을 지키면 화면에 중국어 토큰이 남음. 표시 라벨(` 이름` 부분)만 번역 |
| `意图` / `提示` (명령 접두) | `server/server.go:3820-3836` (`fallbackChat`, `HasPrefix`) | 같은 함수 (영문 `intent`/`hint`도 허용) | 비교값 유지. 사용자 안내문은 영문 `intent`/`hint` 사용을 안내하는 쪽으로 번역 |
| `归档不存在` | `server/task_archives.go:461` | `web/src/app/(main)/function/tasks/page.tsx:1680` | 에러 메시지 부분 문자열 검사 |
| `已存在` | `server/server_mgmt.go:266` 등 | `web/src/app/(main)/system/skills/page.tsx:370` | 에러 메시지 부분 문자열 검사 |
| `新对话` | `server/conversations.go:443` | `web/src/app/(main)/chat/page.tsx` (표시 대체값 4곳) | 기본 제목 비교 |
| `备案` | `db/company_scope.go:160` (`strings.Contains(raw, "备案")`) | `web/src/lib/company-scope.ts:169` (`/icp\|备案/i`) | 범위 규칙을 ICP로 자동 분류 |
| 입력 구분자 `，` `．` `。` | `db/company_scope.go:159` (`".．。"`) | `notify/.../channel-fields.ts`, `notify/page.tsx` (`[\s,，]+`) | 사용자 입력 파싱. 정규식 안의 전각 문자는 유지 |

### 5.8 같은 파일 안의 쌍 (둘 다 같게 바꾸면 허용)

| 문자열 | 위치 | 규칙 |
| --- | --- | --- |
| `命令` | `web/src/components/transcript.tsx:394` (`label: "命令"` 생성) ↔ `:406` (`x.label === "命令"` 비교) | 같은 파일 안의 쌍이다. **두 곳을 완전히 같은 문자열로** 바꾸면 허용(예: `명령`). 한쪽만 바꾸면 도구 입력 표시가 깨진다. 확신이 없으면 둘 다 그대로 둔다 |

---

## 6. 표기 규칙

- **설명 문장**: 자연스러운 한국어 존댓말(`~합니다`).
- **버튼·메뉴·라벨**: 짧고 명확하게(`재검증`, `롤백`, `업데이트`). 조사 생략 가능.
- **조건·제한·주의**: 원문의 조건/제한/경고를 **생략하거나 강화하지 않는다**.
- **외래어 표기**: 국립국어원 외래어 표기 기준(프런트엔드, 프록시, 디렉터리). 널리 쓰는 기술어는 관용 표기 허용(워크스페이스).
- **영문 유지(음역 금지)**: `agent`/`Agent`, `worker`, `planner`, `human-in-the-loop`은 영문 그대로 둔다. (`agent`·`worker`·`planner`는 원문도 영문 표기를 쓰고, `human-in-the-loop`은 원문이 중국어지만 관용적으로 영문을 쓴다.)
- **영어 표기 범위(선별 · 확정)**: 다단어 전문용어·약어·관용 영어만 영어로 둔다 — `human-in-the-loop`, `SSOT`, `injection point`, `wrap-up` 등. 단일 외래어(대시보드·세션·프록시·트래픽·스냅샷 등)는 한글 외래어 표기를 유지한다.
- **괄호 병기**: 처음 등장하는 핵심 개념은 영문 병기 권장 — 예) 이중 그래프 아키텍처, 동일 출처(same-origin).
- **`方式 N` 제목**: **"방법 N"으로 번역(확정)**. 단 `使用方式` 등 일반 "방식"은 문맥대로(강제 치환 금지).
- **문장부호 정리**:
  - 목록·문장 끝에 남은 중국어 `；` → 한국어 쉼표/마침표로 교체.
  - 전각 괄호 `（）`·전각 제목(예: `자산 동기화（ScopeSentry）`) → 반각 `()`로.
  - UI 라벨 인용 `「」` → **작은따옴표 `'…'`로 통일(확정)**. README 포함 전부 변환. 코드를 확인한 결과 「」를 비교값이나 파싱에 쓰는 곳은 없었고, TS 작은따옴표 문자열 안에 「」가 있는 곳과 SQL 문자열 안에 있는 곳도 없어서(SQL은 주석만) 따옴표 충돌은 확인되지 않았다. 이 확인은 정규식 기반이다.
- **중국 법률명**: 통용 표기 + 원문 병기 권장 — 예) 사이버보안법《网络安全法》, 데이터보안법《数据安全法》, 개인정보보호법《个人信息保护法》. (기존 음역 《네트워크안전법》에서 변경 — README 동반 수정)

### deny 맥락별 표기 **[확정]**

소스 자체가 같은 enum `deny`를 화면마다 다르게 부르고 있습니다. 아래 표를 따릅니다.

| 소스 위치 | 원문(zh) | 한국어 | 비고 |
| --- | --- | --- | --- |
| `approval-records.tsx:167` `actionLabels` (`allow`/`ask`/`deny`) | 允许 / 转人工审批 / 拒绝 | 허용 / 사람 승인으로 전환 / **거부** | 승인 기록의 최종 결과 라벨 |
| `system/intercept/page.tsx:300,325` 규칙 선택 | 拦截 (`deny`) | **차단** | 규칙이 자동으로 막음 |
| `system/intercept/page.tsx:747` | 禁止 — 阻断，返回拒绝消息给模型 | 금지 — 차단, 거부 메시지를 모델에 반환 | 규칙 동작 설명 |
| `system/intercept/page.tsx:754` | 拒绝消息（返回给模型） | 거부 메시지(모델에 반환) | |
| `system/intercept/page.tsx:799` | 自动拒绝 | 자동 거부 | 승인 시간 초과 시 동작 |
| `overview-tab.tsx:575`, `lib/types.ts:459` (작업 제약 `allow`/`deny`) | 允许 / 禁止 | 허용 / **금지** | 작업의 조작 제약 관리 |
| `status.ts` `audit` (`allow`/`block`) | 放行 / 拦截 | 허용 / 차단 | 감사 로그 라벨 |

> 원칙: 사람이 거절하거나 시간 초과로 거절되면 **거부**, 규칙·인터셉트가 자동으로 막으면 **차단**, 작업 제약의 `deny`는 **금지**. 소스의 중국어가 일관되지 않아서 코드만으로는 구분되지 않으므로, 새 화면을 번역할 때는 이 표에서 가장 가까운 맥락을 찾아 따릅니다.

---

## 7. 번역해도 되는 것 (과제외 방지)

### 7.1 DB에 저장되는 표시용 문자열 (새 DB 전제로 번역 가능)

새 DB로 시작하므로 **번역된 값이 처음부터 저장**됩니다. 비교에 쓰이는 흔적은 코드에서 찾지 못했습니다(정규식·구문 트리 기반 확인이라 완전하지는 않음).

- 내장 인터셉트 규칙의 이름·메시지: `db/db.go:286-289`(자산 규칙 노트), `:333-518`(v1), `:528-567`(v2), `:586-602`(v3, `const name = "[내장] 삭제류 인터페이스 경로"`)
  - 규칙 이름은 승인 기록 화면의 `row.rule_name`(`approval-records.tsx:121`)에 표시되는 라벨입니다.
  - `name`에는 UNIQUE 제약이 없고 v1/v2의 `ON CONFLICT DO NOTHING`은 실질적 효과가 없습니다. 중복 생성은 settings 플래그가 막으므로, **플래그가 없는 새 DB에서만 안전하게 한 번 시드**됩니다.
- `db/schema.sql:670`('任务执行期间自动关联'), `:688`('由历史任务资产关联迁移'), `:1134`('复测会话已删除')
- `db/tasks.go:270,298`, `db/finding_retests.go:85`('未分类')·`:257`, `db/side_questions.go:286`, `server/finding_retests.go:222,230`('漏洞复测', '内置默认')
- `schema.sql`은 시작할 때마다 다시 실행됩니다. 번역한 트리거 함수의 문구는 기존 DB에서도 **새로 생기는 행**부터 적용되지만, 새 DB에서는 문제가 없습니다.

### 7.2 번역해도 안전한 것
- 화면 표시용 `strings.Join` 구분자 "、": `agent/planner.go:229`, `db/companies.go:636`, `notify/mask.go:73`, `notify/render.go:134-136`
- 5.7에 없는 에러 메시지·토스트·플레이스홀더
- `config.example.json`의 `_comment*` 값(키 이름은 유지), 쉘 스크립트의 안내 문구(응답 비교값은 `y`/`n`과 숫자이므로 유지)

### 7.3 프롬프트 (3단계)
- 프롬프트는 DB에 첫 삽입만 되므로 **새 DB 전제**가 필요합니다(위 "새 DB 전제").
- 출력 언어를 지정하는 곳은 3곳입니다: `agent/promptcatalog.go:70`(기본 어시스턴트, "请用简洁、准确的中文回答"), `:107`(리포터, "全程**中文**"), `agent/retester.go:15`(리테스터, "使用简洁中文答复"). 이 문구를 "한국어"로 번역하면 곧 응답 언어 전환입니다.
- planner·worker·mainagent에는 출력 언어 지시가 없습니다. 이들이 프롬프트 언어를 따라 한국어로 답할지는 **추정이며 아직 실행해서 확인하지 않았습니다**.
- 프롬프트별 출력 형식(`intercept/prompt.go` 외)은 3단계에서 프롬프트마다 점검합니다.
- `skills/*/SKILL.md`(`api-recon`, `scopesentry` 등)는 에이전트가 읽는 지시문이므로 문서가 아니라 **프롬프트로 취급**합니다. `skills/playwright-cli`의 중국어 여부는 미확인입니다.
- 한국어는 중국어보다 길어지는 경향이 있습니다. 길이 제한이 있는 곳(프롬프트의 "120자 이내", `len(reason) > 2400`)은 번역 후 확인하고, UI는 mock 모드(`NEXT_PUBLIC_MOCK`)로 버튼·배지가 깨지지 않는지 봅니다.

### 7.4 테스트
- Go 테스트 안에 중국어 비교가 약 36곳 있습니다. 비교 대상 코드 문자열을 번역할 때만 기대값을 같이 바꿉니다. 코드 문자열이 안 바뀌면 테스트는 건드리지 않습니다.

---

## 8. README 동반 수정 (확정 결정 반영)

README.ko.md(팀원 번역본)에서 아래를 고칩니다. 줄 번호는 팀원 번역본 기준입니다.

| 항목 | 현재 | 수정 |
| --- | --- | --- |
| human-in-the-loop | "휴먼 인 더 루프" (40·309·342행) | `human-in-the-loop` |
| 방식 N | "방식 1~5" 제목 10개 (75·90·107·119·131·153·166·178·189·198행) | "방법 N" (본문의 일반 "방식"은 유지) |
| 已修复 | "수정됨" 4곳 | "수정 완료" |
| 「」 | 27곳 | `'…'` |
| SSOT | "진실 저장소" 353·371행 | "단일 신뢰 소스(SSOT)" |
| 폐루프 | "폐곡선" 387·389·411행 | "폐루프" |
| 법률명 | 501행 음역 | 통용 표기 + 원문 병기 |
| 문장부호 | 중국어 `；`, 전각 `（）` | 6장 문장부호 정리 |
| 企业 | "기업" 0곳 | 영향 없음 |

---

## 9. 미확인·보류

- **프롬프트별 출력 형식 계약**: `intercept/prompt.go` 외 프롬프트는 아직 점검 전이다(3단계에서 확인).
- **변수·상수에 담긴 중국어 비교**: 구문 트리로 못 잡는다. 번역 중 비교 지점을 만나면 5.7에 추가한다.
- **객체 키·`value=` 속성에 쓰인 중국어(프런트)**: 찾지 못했지만 정규식 기반이라 완전하지 않다.
- **멘션 토큰 UX**: 종류명이 중국어로 남는다(5.7). "코드 불변" 방침과의 trade-off로 팀이 수용했는지 확인이 필요하다.
- **planner·worker·mainagent 출력 언어**: 프롬프트를 번역하면 한국어로 답할 것이라는 건 추정이다(7.3).
- **4.2·4.3의 한국어 라벨**: 일괄 `[확정]`(표준 번역). enum은 DNT라 화면 라벨만 번역한다.

---

## 10. 유지보수

- 번역이 UI → 에이전트 응답 → 내부 프롬프트로 깊어질수록 용어가 늘어난다. **이 파일 하나만** 갱신해 전 단계 일관성을 유지한다.
- 새 용어 추가 시: `원문(zh) | 영문 | 한국어 | 비고` 4열을 지키고, 로직 결합이 의심되면 **5장 DNT**에도 함께 등록한다. 중국어 값이 비교·파싱에 쓰이면 **5.7**에 짝 위치와 함께 등록한다.
- 번역 후 코드가 정말 그대로인지 확인하려면, 번역 전후의 Go 구문 트리(AST) 비교와 `git diff`로 비주석·비문자열 변경이 없는지 본다.
- 변경은 PR 리뷰로 반영한다.
- 라이선스: 업스트림은 AGPL-3.0이다. 포크는 라이선스 고지와 저작자 표시를 유지하고, 업스트림 README의 면책·사용 제한(개인 학습·로컬 검증 목적) 문구를 보존한다.
