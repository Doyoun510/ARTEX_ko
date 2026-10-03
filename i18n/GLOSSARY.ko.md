# ARTEX 한국어 번역 단어집 (Glossary)

> README.ko.md 번역을 기준으로 추출한 **공용 용어집**입니다. UI → 에이전트 응답 → 내부 프롬프트·문서 전 단계에서 **하나의 단어집만** 공유해 용어 일관성을 유지합니다.
> 번역 담당자/번역 프롬프트 모두 이 파일을 기준으로 삼습니다.

---

## 0. 최상위 원칙 (반드시 먼저 읽기)

1. **번역 금지(DNT)가 번역보다 우선** — 어떤 문자열을 만나면 "번역할까?"보다 **"이게 로직에 쓰이나?"를 먼저** 판단한다. 비교(`==`/`switch`/`if`), 파싱(`split`/정규식 앵커), 저장(DB 값·직렬화 키), 에이전트 출력 계약(도구명·파라미터·enum)에 쓰이면 → **번역 금지 + 플래그**. ([4. DNT 목록](#4-번역-금지dnt-목록) 참고)
2. **중국어만 번역, 영어는 원문 유지** — 코드·명령·식별자에 섞인 영어는 그대로 둔다.
3. **자리표시자 보존** — `%s` `%d` `%v` `{count}` 등과 그 순서, `\n`, 마크다운(`**`, 백틱), HTML 태그는 보존.
4. **같은 개념 = 같은 번역** — 이 표에 있으면 무조건 표의 번역을 쓴다.
5. **애매하면 번역하지 말고 사람에게** — 틀린 번역보다 미번역 + 사유 플래그가 안전.

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
| 真值库 | single source of truth (SSOT) | 단일 신뢰 소스(SSOT) | 리뷰 반영(기존 "진실 저장소"에서 변경) — README 353·371줄 동반 수정 필요 |
| 攻击链(路) | attack chain | 공격 체인 | |
| 利用链 | exploit chain | 익스플로잇 체인 | |
| 过程级信息交换 | process-level info exchange | 프로세스 수준 정보 교환 | |
| 人在环路 | human-in-the-loop | human-in-the-loop | 영어 그대로 유지(음역 금지) |

### 그래프 노드/개념 (프로즈에서의 번역 — 식별자 자체는 DNT)

| 원문(zh) | 영문(식별자, DNT) | 한국어(개념어) | 비고 |
| --- | --- | --- | --- |
| 目标 | `goal` | 목표 | 노드 kind 식별자 `goal`은 그대로 |
| 意图 | `intent` | 의도 | 식별자 `intent`/`intent_id`는 DNT |
| 事实 | `fact` | 사실 | |
| 发现 / 漏洞 | `finding` / vulnerability | 발견 / 취약점 | `finding`은 "발견", `漏洞`는 "취약점"으로 구분 |
| 提示 | `hint` | 힌트 | |
| 规划者 | planner | 플래너 | 코드/프로즈 모두 `planner` 원어 유지 권장, 개념 설명 시 "플래너" |
| 执行者 | worker | 실행자 | `worker`는 원어 유지, 개념 설명 시 "실행자" |

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
| 审批 | approval | 승인 | |
| 审批记录 | approval records | 승인 기록 | |
| 拦截(기능) / 拦截审批 | intercept / intercept approval | 인터셉트 / 인터셉트 승인 | 코드 `intercept`와 일치 |
| 拦截(동작/결과, 被拦截) | blocked | 차단 | **확정**: 요청이 막힘(동사/결과). 기능명 "인터셉트"와 구분해 사용 |
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
| 域名 / 子域 | domain / subdomain | 도메인 / 서브도메인 | 식별자 `root_domain`/`subdomain`은 DNT |
| 端口 / 站点 / 端点 | port / site / endpoint | 포트 / 사이트 / 엔드포인트 | |
| 公司 / 企业 | company / enterprise | 회사 / (기업) | 公司=회사 확정; 企业는 UI 번역 때 결정 |
| 公司资产范围 | company asset scope | 회사 자산 범위 | 식별자 `scope`/`task_scope`는 DNT |

---

## 3. 인프라·운영 용어

| 원문(zh) | 영문 | 한국어 | 비고 |
| --- | --- | --- | --- |
| 反向代理 | reverse proxy | 리버스 프록시 | |
| 同源 | same-origin | 동일 출처 | 괄호로 `(same-origin)` 병기 권장 |
| 缓冲 | buffering | 버퍼링 | |
| 长连接 | long connection | 장기 연결 | |
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
| 闭环 | closed loop | 폐루프 | ⚠️ "폐곡선"(기하학 용어)은 오역 — README 387·389·411줄 동반 수정 필요 |
| 唤醒 | wake | 깨우기 | |
| 无状态会话 | stateless session | 무상태 세션 | |
| 发布包 / Release 压缩包 | release package | 배포 패키지 | `Release`/`Releases` 고유명은 DNT |

### 공격 기법 용어

| 원문(zh) | 영문 | 한국어 |
| --- | --- | --- |
| 注入点 | injection point | 주입점 |
| 凭据 | credential | 자격 증명 |
| 横向 | lateral movement | 횡적 이동 |
| 提权 | privilege escalation | 권한 상승 |

---

## 4. 상태·열거값 (⚠️ 표시 레이어만 번역, 저장/비교 레이어는 DNT)

finding 결론·상태 등은 **화면 라벨은 번역**하되, 코드/DB에서 그 값으로 비교·저장되면 **그 인스턴스는 건드리지 말 것**. (실제로 enum인지 생문자열인지 **코드 확인 후** 확정)

| 원문(zh) | 한국어(표시 라벨) | 비고 |
| --- | --- | --- |
| 仍可复现 | 여전히 재현 가능 | 재검증 결론 |
| 已修复 | 수정 완료 | **확정**. 화면 라벨만 번역 — 내부 enum `fixed`는 DNT |
| 无法确认 | 확인 불가 | 재검증 결론 |
| 复测中 | 재검증 중 | 진행 표시 |
| 复现 | 재현 | |

---

## 5. 번역 금지(DNT) 목록

아래는 **절대 번역하지 않고 원문 그대로** 둔다. (영어 식별자는 영어 그대로)

### 5.1 제품·고유명
`ARTEX` · `ScopeSentry` · `AegisHook` · `norma` · `SecSentry` · `Cairn` · `Kali` · `PostgreSQL` · `SQLite` · `Next.js` · `Go` · `Docker` · `Nginx` · `UPX` · `playwright` · `nmap` · `ripgrep` · `curl` · `vim` · `npm`

### 5.2 약어·프로토콜
`MITM` · `CA` · `JWT` · `SSE` · `MCP` · `HTTP(S)` · `REST` · `DNS` · `TCP` · `JSON-RPC` · `API` · `LLM` · `CWD`

### 5.3 도구명·함수·파라미터 (에이전트 출력 계약 — 번역 시 즉시 파손)
`prove_goal` · `add_intent` · `add_hint` · `add_task_hint` · `claimNext` · `graph_overview` · `search_all_worker_traces` · `list_worker_traces` · `get_worker_trace` · `report_finding` · `update_finding_report` · `bind_finding_traffic` · `get_finding_traffic` · `traffic_search` · `traffic_get` · `list_findings` · `list_task_findings` · `node_detail` · `get_task_node_detail` · `TodoWrite` · `plannerLoop` · `ToolSet` · `retester`
파라미터/키: `intent_id` · `asset_id(s)` · `node_id` · `step_ids` · `finding_id` · `finding_node_id` · `evidence_hint_id` · `traffic_refs` · `evidence_version` · `version` · `role`

### 5.4 그래프 노드 kind·엣지·필드 (식별자)
노드/엣지: `goal` · `intent` · `fact` · `finding` · `hint` · `spawns` · `derived_from` · `yields` · `proves` · `anchor` · `frontier` · `activity`
DB/스키마: `assets` · `companies` · `task_scope` · `exploration_nodes` · `exploration_anchors` · `anchors` · `root_domain` · `subdomain` · `ip` · `service` · `app` · `endpoint`
role enum: `baseline` · `proof` · `verification` · `supporting`

### 5.5 norma SDK 모듈
`agentcore` · `tool` · `permission` · `harness` · `memory` · `transcript`

### 5.6 환경변수·경로·명령·파일
`ANTHROPIC_API_KEY` · `OPENAI_API_KEY` · `ARTEX_PG_DSN` · `ARTEX_LLM_PROVIDER` · `ARTEX_LLM_MODEL` · `ARTEX_LLM_BASE_URL` · `ARTEX_LLM_PROXY` · `ARTEX_TAG` · `ARTEX_TARGETS` · `NEXT_PUBLIC_SSE_BASE` · `NEXT_PUBLIC_MOCK` · `CGO_ENABLED` · `embedui`
파일/경로: `config.json` · `config.example.json` · `.env` · `start.sh`/`start.bat` · `install.sh` · `update.sh` · `build.sh` · `dev.sh` · `schema.sql` · `jwt.key` · `pgdata` · `./data` · `./skills` · `artex(.new/.old/.failed)` · `SHA256SUMS`
모든 **명령·CLI 플래그·코드 블록 내 코드**(예: `docker compose up -d`, `go test ./...`, `-tags embedui`, `-addr`, `-proxy`, `--release`, `--upx`)는 DNT. (코드 블록 안 **중국어 주석만** 번역 대상)

---

## 6. 표기 규칙

- **설명 문장**: 자연스러운 한국어 존댓말(`~합니다`).
- **버튼·메뉴·라벨**: 짧고 명확하게(`재검증`, `롤백`, `업데이트`). 조사 생략 가능.
- **조건·제한·주의**: 원문의 조건/제한/경고를 **생략하거나 강화하지 않는다**.
- **외래어 표기**: 국립국어원 외래어 표기 기준(프런트엔드, 프록시, 디렉터리). 널리 쓰는 기술어는 관용 표기 허용(워크스페이스).
- **영문 유지(음역 금지)**: `agent`/`Agent`, `worker`, `planner`, `human-in-the-loop` 는 원문도 영문이므로 영문 그대로 둔다.
- **괄호 병기**: 처음 등장하는 핵심 개념은 영문 병기 권장 — 예) 이중 그래프 아키텍처, 동일 출처(same-origin).
- **`方式 N` 제목**: README는 "방식 N"으로 번역(리뷰는 "방법 N" 권장 — 선택 사항, 확정 시 일괄 변경).
- **문장부호 정리**(리뷰 반영, README 동반 수정 필요):
  - 목록·문장 끝에 남은 중국어 `；` → 한국어 쉼표/마침표로 교체.
  - 전각 괄호 `（）`·전각 제목(예: `자산 동기화（ScopeSentry）`) → 반각 `()`로.
  - UI 라벨 인용 `「」` → **결정 대기**: README는 유지하되, UI/코드 번역 시 작은따옴표 `'...'`로 통일할지 팀 결정.
- **중국 법률명**: 통용 표기 + 원문 병기 권장 — 예) 사이버보안법《网络安全法》, 데이터보안법《数据安全法》, 개인정보보호법《个人信息保护法》. (기존 음역 《네트워크안전법》에서 변경 권고 — README 501줄 동반 수정 필요)

---

## 7. 유지보수

- 번역이 UI → 에이전트 응답 → 내부 프롬프트로 깊어질수록 용어가 늘어난다. **이 파일 하나만** 갱신해 전 단계 일관성을 유지한다.
- 새 용어 추가 시: `원문(zh) | 영문 | 한국어 | 비고` 4열을 지키고, 로직 결합이 의심되면 **5장 DNT**에도 함께 등록한다.
- 변경은 PR 리뷰로 반영한다.
