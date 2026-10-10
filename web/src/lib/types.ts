// ARTEX domain model — types used across the UI.
// Derived from the functional spec (section 7: 핵심 데이터 형태).

export type TaskStatus = "created" | "queued" | "running" | "paused" | "done" | "failed" | "timeout";
export type EngineMode = "exploring" | "paused" | "stalled" | "idle";

export interface Task {
  id: string;
  name?: string; // 선택적 작업 이름; 빈 값/생략=이름 없음, 표시 시 설명으로 폴백
  category_id?: number;
  category_name?: string;
  pinned?: boolean;
  pinned_at?: string | null;
  description: string;
  goal: string;
  status: TaskStatus;
  created_at: string;
  created_unix?: number; // created_at as unix seconds (run-duration calc)
  completed_at?: string; // RFC3339 finish time (done/failed); "" if unfinished
  completed_unix?: number; // completed_at as unix seconds (0/undef if unfinished)
  last_activity_unix?: number; // unix seconds of the last activity (0/undef if none)
  paused?: boolean;
  queued?: boolean;
  active?: boolean;
  in_flight?: number;
  findings?: { critical: number; high: number; medium: number; low: number }; // 등록된 취약점 수(심각도별 구분)
  last_activity?: string;
  stalled?: boolean;
  goals_total?: number;
  goals_met?: number;
  engine_mode?: EngineMode;
  tokens?: TokenTotal; // whole-task token consumption
  llm_profile_id?: number; // LLM profile used; absent = default profile
  llm_profile_ids?: number[]; // ordered task-level failover chain
  active_llm_profile_id?: number; // profile used by the next LLM call
  llm_failover_state?: "default" | "ready" | "chain_exhausted" | string;
  llm_failover_reason?: string;
  source_task_ids?: string[]; // directly related tasks inherited as read-only context
  archive_blocked_by_task_id?: string; // live direct dependent that must be archived first
  company_ids?: number[]; // associated company scopes; current company assets join the task at creation
  coverage_enabled?: boolean; // 자산 커버리지 기능 스위치(생성 시 결정, 기본 켜짐); false=커버리지 계산/표시 안 함
}

export interface TaskCategory {
  id: number;
  name: string;
  task_count: number;
  created_at: string;
  updated_at: string;
}

export interface TaskTemplate {
  id: number;
  name: string;
  description: string;
  goal: string;
  category_id?: number | null; // 사전 설정 분류; null/생략=없음
  intercept_rules?: AssetInterceptRuleInput[]; // 사전 설정된 작업 수준 인터셉트/허용 규칙
  created_at: string;
  updated_at: string;
}

export interface DeleteTaskOptions {
  delete_assets: boolean;
  delete_traffic: boolean;
  delete_files: boolean;
  delete_findings: boolean;
  delete_llm_records: boolean;
}

export interface DeleteTaskResult {
  deleted: string;
  assets_deleted: number;
  assets_detached: number;
  traffic_deleted: number;
  files_deleted: boolean;
  findings_deleted: number;
  llm_records_deleted: number;
  cleanup_warning?: string;
}

export type TaskArchiveState =
  | "archive_queued"
  | "archiving"
  | "archive_failed"
  | "ready"
  | "restore_queued"
  | "restoring"
  | "restore_failed"
  | "delete_queued"
  | "deleting"
  | "delete_failed";

export interface TaskArchiveTokenStats {
  calls?: number;
  input_tokens?: number;
  output_tokens?: number;
  cache_read_tokens?: number;
  cache_write_tokens?: number;
}

export interface TaskArchive {
  id: number;
  task_id: number;
  state: TaskArchiveState;
  phase: string;
  progress: number;
  error?: string;
  warnings?: string[];
  format_version: number;
  sha256?: string;
  original_size: number;
  compressed_size: number;
  task_name: string;
  task_description: string;
  task_goal: string;
  original_status: TaskStatus;
  category_id?: number;
  category_name?: string;
  source_task_ids: number[];
  remaining_timeout_seconds: number;
  data_counts: Record<string, number>;
  aggregate_stats: {
    tokens?: TaskArchiveTokenStats;
    skills?: Record<string, number>;
    tools?: Record<string, number>;
    findings?: Record<string, number>;
  };
  archived_at?: string;
  requested_at: string;
  created_at: string;
  updated_at: string;
}

export interface TaskArchivePage {
  items: TaskArchive[];
  total: number;
  page: number;
  size: number;
}

export interface ArchiveBatchItem {
  id: string;
  archive_id?: number;
  ok: boolean;
  queued: boolean;
  error?: string;
}

// ---- Asset graph (global, shared across tasks) ----
export type AssetType =
  | "company"
  | "domain"
  | "ip"
  | "port"
  | "service"
  | "site"
  | "endpoint"
  | "parameter"
  | "tech"
  | "credential"
  | "data";

export type NodeState = "observed" | "confirmed" | "tombstoned";

export interface AssetNode {
  id: string;
  type: AssetType;
  name: string;
  key: string; // nkey
  value?: string;
  company_id?: string; // 귀속 회사 자산 id; 빈 값=미귀속
  state: NodeState;
  confidence: number; // 0..1
  attrs?: Record<string, unknown>;
  first_seen: string;
  last_seen: string;
}

export type AssetRel =
  | "owns"
  | "resolves"
  | "exposes"
  | "runs"
  | "serves"
  | "has_endpoint"
  | "has_param"
  | "fingerprinted"
  | "authenticates_as"
  | "reachable"
  | "has_subdomain";

export interface Edge {
  src: string;
  dst: string;
  rel: AssetRel | ExploreRel;
}

// Task asset view — server-side enriched, paginated.
export interface TaskAssetRef {
  id: string;
  name?: string;
  key: string;
  attrs?: Record<string, unknown>;
}

export interface TaskAssetItem extends AssetNode {
  techs?: TaskAssetRef[];
  auth?: TaskAssetRef[];
  params?: TaskAssetRef[];
}

export interface TaskAssetView {
  counts: Record<string, number>;
  total: number;
  items: TaskAssetItem[];
}

// ---- New unified asset model (new backend) ----
export type NewAssetType = "root_domain" | "ip" | "subdomain" | "app" | "service" | "endpoint";

export interface Asset {
  id: number;
  type: NewAssetType;
  company_id?: number;
  task_ids: number[];
  domain?: string;
  root_domain?: string;
  ip?: string;
  c_segment?: string;
  port?: number;
  icp?: string;
  bound_domains?: string[];
  open_ports?: { port: number; service?: string }[];
  record_type?: string;
  record_value?: string[] | string;
  bundle_id?: string;
  app_name?: string;
  category?: string;
  app_description?: string;
  app_icp?: string;
  url?: string;
  service_type?: string;
  service_name?: string;
  favicon_mmh3?: string;
  status_code?: number;
  content_length?: number;
  page_title?: string;
  technologies?: string[];
  auth?: Record<string, unknown>[];
  method?: string;
  params?: Record<string, unknown>[];
  extra?: Record<string, unknown>;
  last_seen: string;
  task_source?: string;
  task_source_summary?: string;
  task_source_node_id?: number;
}

export interface IntentAsset {
  intent_id: number | string;
  asset_id: number;
  type: NewAssetType;
  label: string;
  source: string;
  source_summary: string;
  source_node_id?: number;
  source_task_id: number;
  inherited: boolean;
}

export interface TaskAssetMutation {
  requested: number;
  attached: number;
  existing: number;
}

export interface TaskAssetScopeMutation {
  requested: number;
  assets_linked: number;
  assets_existing: number;
  scopes_added: number;
  scopes_existing: number;
}

// ---- Asset coverage graph (per task) ----
// 포스 디렉티드 '자산 커버리지 그래프'의 한 노드. key는 고유: 자산="a:<id>", 회사="c:<id>",
// 자산 행이 없는 루트 도메인="r:<domain>". in_scope=false는 연결선 용도의 회색 컨텍스트 노드.
export interface CoverageGraphNode {
  key: string;
  kind: "company" | "root_domain" | "subdomain" | "ip" | "service" | "app" | "endpoint";
  label: string;
  tested: boolean;
  in_scope: boolean;
  asset_id?: number;
  company_id?: number;
  domain?: string;
  root_domain?: string;
  ip?: string;
  url?: string;
  port?: number;
  service_type?: string;
  app_name?: string;
  page_title?: string;
  status_code?: number;
}

export interface CoverageGraphEdge {
  src: string;
  dst: string;
}

export interface CoverageGraphData {
  nodes: CoverageGraphNode[];
  edges: CoverageGraphEdge[];
}

// 특정 자산이 이 작업 탐색 그래프에서 연관된 의도/사실/발견(커버리지 그래프 노드 드로어용).
export interface CoverageAssetRef {
  id: number;
  kind: string;
  state: string;
  summary: string;
  source_task_id?: string;
  inherited?: boolean;
}
export interface CoverageAssetRefs {
  intents: CoverageAssetRef[];
  facts: CoverageAssetRef[];
  findings: CoverageAssetRef[];
}

// ---- Workspace file manager (workDir) ----
export interface WorkspaceEntry {
  name: string;
  path: string; // workspace-relative, forward slashes
  dir: boolean;
  size: number;
  mtime: number; // unix millis
}
export interface WorkspaceListing {
  path: string;
  entries: WorkspaceEntry[];
}
export interface WorkspaceFile {
  path: string;
  size: number;
  binary: boolean;
  too_large?: boolean;
  content?: string;
}

// 작업 테스트 범위 1건(커버리지 분모 + 권한 경계).
export interface TaskScopeRow {
  id: number;
  task_id: number;
  kind: "company" | "root_domain" | "subdomain" | "ip" | "cidr" | "icp" | "keyword";
  company_id?: number;
  company_name?: string; // 백엔드가 companies를 JOIN해 해석, kind=company일 때만 값 있음
  domain?: string;
  net?: string;
  value?: string;
  source: "auto" | "agent" | "manual";
  reason?: string;
}

export type CompanyScopeKind = "domain" | "ip" | "cidr" | "icp" | "keyword";

// 회사 추가 시 제출하는 구조화된 자산 범위 규칙.
export interface CompanyScopeRule {
  kind: CompanyScopeKind;
  value: string;
}

// 자산 범위 쓰기 결과. errors는 이번 제출에서 잘못된 행; warnings는 이번 제출과 무관하지만
// 귀속 결과를 예상과 다르게 만드는 기존 데이터 문제(예: ip 필드에 호스트명이 저장된 자산).
export interface CompanyScopeMutation {
  added: number;
  skipped: number;
  invalid: number;
  errors?: string[];
  warnings?: string[];
}

// 회사 자산 범위 규칙 1건(귀속의 유일한 진실 출처).
export interface ScopeRow {
  id: number;
  company_id: number;
  kind: CompanyScopeKind;
  domain?: string; // kind=domain일 때 값 있음
  net?: string; // kind=ip|cidr일 때 값 있음
  value?: string; // kind=icp|keyword일 때 백엔드가 직접 반환할 수 있음
  raw: string; // 원본 사용자 입력, 표시와 되채우기에 사용
  reason?: string;
}

// 회사: type=company인 자산 노드 + 아이콘 + 자산 개수 + 자산 범위 규칙.
export interface Company {
  id: number;
  name: string;
  logo?: string; // 원격 아이콘 URL; 비면 프런트가 이름 첫 글자 사용
  asset_count: number;
  scope?: ScopeRow[];
}

// ---- Exploration graph (per task) ----
export type ExploreKind = "task" | "begin" | "goal" | "intent" | "fact" | "finding" | "hint" | "digest";
export type GoalState = "open" | "met" | "abandoned";
export type IntentState = "open" | "running" | "paused" | "done" | "blocked" | "exhausted" | "stopped";
export type FindingState = "confirmed" | "dismissed";
export type HintState = "active" | "consumed";
export type ExploreRel = "spawns" | "derived_from" | "yields" | "proves" | "covers";

export interface TaskNode {
  id: string;
  type: ExploreKind;
  payload?: string;
  priority: number; // 0..10
  state: string; // GoalState | IntentState | FindingState | HintState
  origin: string;
  ts: string;
  source_task_id?: string;
  inherited?: boolean;
  delete_reason?: string; // 의도 논리 삭제(state='deleted') 시의 삭제 사유
}

// 브로드캐스트 보드 한 페이지: 생성 순서로 페이지네이션한 노드 + 이 페이지가 관련된 엣지 + 엣지 반대편 노드(refs, id로 인덱싱),
// 이렇게 하면 각 브로드캐스트가 '어디서 왔고 무엇을 만들었는지'를 명확히 할 수 있어 전체 그래프를 가져올 필요가 없다.
export interface ExplorationNodePage {
  items: TaskNode[];
  total: number;
  page: number;
  size: number;
  edges: Edge[];
  refs: Record<string, TaskNode>;
  // 노드 id → 그 노드가 앵커한 자산(브로드캐스트 보드 펼칠 때 함께 표시, 이 페이지 노드와 그 이웃 포함).
  assets: Record<string, FindingAsset[]>;
}

export interface ExplorationNodeQuery {
  page?: number;
  size?: number;
  kinds?: ExploreKind[];
  states?: string[];
  q?: string;
  order?: "asc" | "desc";
}

// 목표 관리 카드용 목표(백엔드가 payload를 text/vulnclass로 분해함).
export interface TaskGoal {
  id: string;
  text: string;
  vulnclass?: string;
  state: string; // GoalState
  origin?: string;
  ts: string;
}

// 제약 관리 카드용 동작 제약(allow=허용 / deny=금지).
export type ConstraintKind = "allow" | "deny";
export interface TaskConstraint {
  id: string;
  kind: ConstraintKind;
  text: string;
  origin?: string;
  ts?: string;
}

// ---- Findings ----
export type Severity = "critical" | "high" | "medium" | "low";

// 취약점 처리 상태: 처리 대기 / 처리 중 / 확인됨 / 처리됨 / 수정 완료 / 오탐 / 무시 / 중복 / 위험 수용.
export type FindingStatus =
  | "pending"
  | "in_progress"
  | "confirmed"
  | "resolved"
  | "fixed"
  | "false_positive"
  | "ignored"
  | "duplicate"
  | "risk_accepted";

// FindingAsset은 한 취약점에 바인딩된 자산(백엔드에서 label을 미리 렌더링함).
export interface FindingAsset {
  id: string;
  type: string;
  label: string;
}

export interface Finding {
  traffic_count?: number;
  evidence_version?: number;
  report_evidence_version?: number;
  report_stale?: boolean;
  id: string;
  finding_id?: string; // 독립 findings 테이블의 행 id, 상태 업데이트 핸들(작업 내 오래된 노드는 없을 수 있음)
  vulnclass: string;
  name?: string; // 취약점 이름; 비면 표시 시 vulnclass로 폴백
  severity: Severity;
  status: FindingStatus;
  summary: string;
  evidence: string;
  report?: string; // 상세 보고서(Markdown); 상세 API만 반환, 리스트에서는 빈 값
  intent_id?: string;
  param_id?: string;
  task_id?: string;
  task_description?: string;
  source_task_id?: string;
  inherited?: boolean;
  assets?: FindingAsset[];
  ts: string;
}

// FindingsPage는 발견 리스트의 서버 측 페이지네이션 응답.
export interface FindingsPage {
  items: Finding[];
  total: number;
  page: number;
  page_size: number;
}

export interface FindingGroup {
  task_id: string | number | null;
  task_name?: string; // 선택적 작업 이름; 빈 값/생략=이름 없음
  task_description: string;
  task_status: string;
  count: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  last_found_at: string;
}

export interface FindingGroupsPage {
  items: FindingGroup[];
  total: number;
  finding_total: number;
  page: number;
  page_size: number;
}

export interface FindingDeepenResponse {
  task_id: string;
  intent_id: string;
  state: IntentState;
  queued: boolean;
}

// FindingStats는 발견 전체 테이블 집계(통계 카드 + 취약점 유형 드롭다운), 서버 측 계산, 페이지네이션 영향 없음.
export interface FindingStats {
  total: number;
  pending: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  vulnclasses: string[];
  tasks: FindingTaskOption[];
}

// FindingTaskOption은 발견 페이지 '작업별' 필터 드롭다운의 한 항목: 취약점이 있는 작업(설명이 비면 작업이 삭제됨을 의미,
// 프런트는 id로 폴백 표시)과 그 취약점 건수.
export interface FindingTaskOption {
  id: string | number;
  name?: string; // 선택적 작업 이름; 빈 값/생략=이름 없음
  description: string;
  count: number;
}

// FindingQuery는 발견 리스트 페이지네이션/필터/정렬 파라미터.
export interface FindingQuery {
  page: number;
  pageSize: number;
  severity?: "all" | Severity;
  status?: "all" | FindingStatus;
  vulnclass?: string;
  task?: string; // 작업 id; "all"/빈 값 = 작업으로 필터 안 함
  query?: string;
  sort?: "severity" | "time";
  // 자산 트리 노드 key; 노드 하나 선택 = 그 서브트리 전체 선택. 빈 값 = 자산으로 필터 안 함.
  assetScope?: string;
}

// ---- Findings by asset (자산 뷰) ----
export type FindingAssetKind = "company" | "root_domain" | "subdomain" | "ip" | "service" | "app" | "endpoint" | "none";

// FindingAssetNode는 자산 트리의 한 노드. key는 a:<id>(자산), c:<id>(회사),
// r:<domain>(DB에 자산 행이 없는 루트 도메인), __none__(미연관 자산) 형태.
export interface FindingAssetNode {
  key: string;
  parent?: string;
  kind: FindingAssetKind;
  label: string;
  asset_id?: number;
  company_id?: number;
  self: number; // 해당 자산에 직접 매달린 발견 수
  total: number; // 자손 포함, 발견 기준 중복 제거
  critical: number;
  high: number;
  medium: number;
  low: number;
  last_found_at: string;
}

export interface FindingAssetTree {
  nodes: FindingAssetNode[];
  finding_total: number;
  truncated: boolean;
  dropped_kinds?: string[];
}

// FINDING_UNASSIGNED_ASSET은 백엔드 db.FindingUnassignedAsset에 대응.
export const FINDING_UNASSIGNED_ASSET = "__none__";

// ---- Activity / sessions ----
export type ActivityKind =
  | "tool_use"
  | "tool_result"
  | "text"
  | "thinking"
  | "result"
  | "user"
  | "intent" // LLM-generated exploration objective leading a worker session (UI-synthesized)
  | "round" // planner round boundary marker (engine-emitted)
  | "usage" // live cumulative token usage (per model turn); not rendered
  | "llm_switch" // automatic/manual task-level LLM switch
  | "llm_failover" // task-level provider switch / chain exhaustion audit event
  | "intercept_request"; // user-approval request from the intercept layer

// ChatAttachment은 한 번 업로드된 파일: path는 해당 세션/작업 작업 디렉터리(즉 agent의 CWD) 기준 상대 경로.
export interface ChatAttachment {
  name: string;
  path: string;
  size: number;
  abs?: string; // 절대 경로(scope=staging 임시 업로드 시 반환; 작업 생성 전 설명에 기록)
}

export interface Activity {
  seq: number;
  intent_id?: string;
  worker: string; // session owner: planner | mainagent | work#1 ...
  ts: string;
  kind: ActivityKind;
  tool?: string;
  tool_use_id?: string;
  is_error?: boolean;
  summary: string;
  detail?: string;
  metadata?: {
    llm_transition?: LLMTransition;
  };
  source_task_id?: string;
  inherited?: boolean;
  main_seg?: number; // main-agent conversation segment (present only on worker="mainagent" rows)
  // token usage (present only on kind='result')
  input_tokens?: number;
  output_tokens?: number;
  cache_read_tokens?: number;
  cache_write_tokens?: number;
}

export interface LLMAuditProfile {
  id: number;
  name: string;
  format: string;
  model: string;
}

export interface LLMTransition {
  mode: "automatic" | "manual" | "exhausted";
  reason: string;
  previous?: LLMAuditProfile;
  next?: LLMAuditProfile;
}

export interface TaskLLMResolution {
  profile_id?: number;
  name: string;
  format: string;
  model: string;
  source: "task_chain" | "agent_binding" | "global_profile" | "environment" | "global";
  available: boolean;
  reason?: string;
}

export interface TaskLLMResolutions {
  mainagent: TaskLLMResolution;
  planner: TaskLLMResolution;
  worker: TaskLLMResolution;
}

// ---- Agent triggers (P3 스케줄, 커스텀 agent만) ----
export interface AgentTrigger {
  id: number;
  agent_key: string;
  enabled: boolean;
  interval_sec: number; // 정기: N초마다(0=비정기)
  on_finding: boolean; // 임의 작업이 finding을 발견할 때 트리거
  on_goal_met: boolean; // 임의 작업이 목표를 달성할 때 트리거
  on_task_timeout: boolean; // 임의 작업이 타임아웃될 때 트리거
  on_tool_call: boolean; // 선택한 도구가 호출(실행 완료)될 때 트리거
  on_task_create: boolean; // 임의 작업이 생성될 때 트리거
  interval_message: string; // 각 트리거 조건의 개별 사용자 메시지
  finding_message: string;
  goal_message: string;
  task_timeout_message: string;
  tool_call_message: string;
  task_create_message: string;
  tool_names: string[]; // on_tool_call에서 선택한 도구 key(최소 1개)
  last_fire?: string;
}

// ---- Conversations (chat page) ----
export interface ActiveFindingRetest {
  id: number;
  finding_id: string;
  conversation_id: number;
  status: "pending" | "running";
}

export interface FindingRetest {
  id: number;
  finding_id: number;
  conversation_id: number | null;
  status: "pending" | "running" | "completed" | "failed" | "stopped";
  verdict: "" | "reproduced" | "fixed" | "inconclusive";
  notes: string;
  summary: string;
  evidence: string;
  error: string;
  created_at: string;
  started_at: string | null;
  finished_at: string | null;
}

export interface Conversation {
  id: number;
  running?: boolean; // live server state, returned with the conversation list
  agent_key: string;
  title: string;
  llm_profile_id?: number;
  pinned?: boolean;
  pinned_at?: string | null;
  created_at: string;
  updated_at: string;
}

// ---- Backend logs (/logs page) ----
export interface LogLine {
  seq: number;
  db_id?: number; // server_logs.id; present for DB-persisted lines
  ts: string;
  level: "info" | "warn" | "error";
  tag: string;
  text: string;
}

export type SessionRole = "mainagent" | "planner" | "worker" | "system";
export type SessionStatus = "running" | "paused" | "done" | "blocked" | "exhausted" | "pending" | "stopped" | "deleted";

// Daily token aggregate bucket (GET /api/tokens/daily).
export interface DailyTokenBucket {
  date: string; // "YYYY-MM-DD"
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

// Per-worker token usage (GET /api/exploration/tokens).
export interface TokenUsage {
  worker: string;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

export interface SessionTokenUsage {
  session: string;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

export interface BatchControlItem {
  id: string;
  ok: boolean;
  status?: string;
  queued?: boolean;
  error?: string;
}

// 일괄 분류 변경의 작업별 결과. 실패는 작업이 이미 삭제된 경우뿐이며, 분류 자체의 쓰기는 원자적이다.
export interface BatchCategoryItem {
  id: string;
  ok: boolean;
  error?: string;
}

// Whole-task (all agents) token aggregate.
export interface TokenTotal {
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

// Global per-profile token spend from the llm_usage ledger (GET /api/tokens/usage).
export interface ProfileUsage {
  profile_name: string;
  calls: number;
  tasks: number;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

// One (profile, UTC day) token bucket for the dashboard's daily chart (new source).
export interface ProfileDayUsage {
  profile_name: string;
  date: string; // YYYY-MM-DD
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
}

// Response of GET /api/tokens/usage — the dashboard's "new" (llm_usage) token view.
export interface UsageStats {
  by_profile: ProfileUsage[];
  daily: ProfileDayUsage[];
}

// Per-model token usage for one task (GET /api/llm/records/by-model), from the
// always-on llm_usage metering ledger. calls = number of LLM calls on this model.
export interface ModelTokenStat {
  model: string;
  calls: number;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

export interface Session {
  id: string;
  role: SessionRole;
  title: string;
  status: SessionStatus;
  live: boolean;
  last_activity: string;
  intent_id?: string;
  source_task_id?: string;
  inherited?: boolean;
  seg?: number; // main-agent session: which conversation segment (0 = original)
}

// ---- Security ----
export interface AuditEntry {
  ts: string;
  tool: string;
  action: "allow" | "block";
  reason?: string;
  command?: string;
}

export interface Audit {
  entries?: AuditEntry[];
  attributions?: Record<string, number>;
}

// ---- Traffic ----
export interface TrafficExchange {
  id: string;
  ts: string;
  host: string;
  method: string;
  url: string;
  status: number;
  content_type: string;
  resp_len: number;
}

export interface TrafficResp {
  enabled: boolean;
  proxy?: string;
  count?: number; // global total (unfiltered)
  total?: number; // rows matching the current filter (for pagination)
  page?: number;
  size?: number;
  exchanges?: TrafficExchange[];
}

// Full raw request/response of one exchange (lazy-loaded on row select).
export interface TrafficDetail {
  req: string;
  resp: string;
}

// One distinct recorded host with its exchange count (target picker).
export interface TrafficHost {
  host: string;
  count: number;
}

// ---- App settings (runtime toggles) ----
export interface Settings {
  traffic_capture: boolean;
  agent_traffic_binding: boolean; // Agent가 트래픽 증거를 자동 바인딩, 기본 꺼짐; 수동 바인딩에는 영향 없음
  llm_record: boolean; // LLM 레코딩 스위치(기본 꺼짐); 끄면 어떤 LLM 호출도 기록하지 않음
  // Web search. brave_key_set / tavily_key_set reflect whether a key is stored
  // (the values are never returned). On PUT, send the corresponding field to set/clear.
  web_search_enabled: boolean;
  web_search_backend: string; // "ddgs" | "brave-free" | "tavily" | "deepseek"
  brave_key_set: boolean;
  tavily_key_set: boolean;
  // write-only: only sent on PUT to store/clear the key.
  brave_search_api_key?: string;
  tavily_search_api_key?: string;
  // 독립 아웃바운드 프록시(http/https/socks5), 검색 엔드포인트 접근용; 트래픽을 기록하는 MITM 프록시와 무관. 빈 값=직접 연결.
  web_search_proxy?: string;
  // 전역 아웃바운드 프록시(http/https/socks5, user:pass 가능), 모든 대상 트래픽이 이를 경유. 트래픽 캡처 켜짐 시
  // MITM 업스트림으로 동작; 캡처 꺼짐 시 agent의 bash/WebFetch에 직접 주입. 빈 값=직접 연결.
  global_proxy?: string;
  python_interpreter?: string; // 커스텀 스크립트 도구의 python 인터프리터 경로(빈 값=런타임 감지)
  workers?: number; // 동시 작업 agent 수(기본 3); 이후 시작되는 작업에 적용
  // 작업 동시 실행 상한: 동시 '실행 중' 작업 수 상한. 꺼짐=무제한; 켜면 새 작업이 상한 초과 시 대기, 빈자리 생기면 자동 시작.
  task_concurrency_enabled?: boolean; // 기본 false
  task_concurrency_limit?: number; // 켜면 기본 5
  // LLM 순환 전환(장애 조치). 기본 꺼짐; 켜면 '모델 미지정' agent가 현재 설정 사용 불가일 때
  // (잔액 부족/key 무효/요청 제한/서비스 이상) 자동으로 다음 설정으로 전환.
  llm_pool_enabled?: boolean; // 기본 false
  // 지정 설정을 바인딩한 agent/작업이 실패할 때 순환 전환 체인을 사용해 대체 처리할지. 기본 false = 바인딩하면 독점.
  llm_pool_bind_fallback?: boolean;
  // 동작 제약 주입 범위(기본 모두 켜짐): 작업의 allow/deny 제약을 해당 agent의 시스템 프롬프트에 결합.
  constraints_inject_planner?: boolean;
  constraints_inject_worker?: boolean;
  // 실험 기능: noa 모델 기반 컨텍스트 압축(기본 꺼짐). 켜면 플랫폼에 연결된 4종 agent(planner/
  // worker/메인 agent/대화)의 컨텍스트 압축을 noa가 인계받아 내장 compaction을 대체; run마다 한 번 읽어,
  // 이후 시작되는 run에 적용.
  noa_compaction?: boolean;
  // ---- 취약점 IM 푸시(채널 자체는 독립 리소스, /api/notify/* 참고, 여기엔 전역 설정 3개만) ----
  notify_enabled?: boolean; // 푸시 마스터 스위치, 기본 켜짐; 유지보수 기간 원클릭 차단용
  notify_public_base_url?: string; // 취약점 상세 링크의 외부 접근 주소; 빈 값=메시지에 상세 링크 없음
  notify_digest_interval_min?: number; // 모아 보내기 모드 주기(분), 기본 30
}

// ---- 취약점 IM 푸시 ----

// NotificationFilter의 필드는 모두 선택적이며, 저장 시 min_severity는 빈 값 또는 등록된 심각도만 허용한다.
// 읽기 시 ParseFilter는 JSON 파싱 오류를 반환하지 않고 파싱된 값을 Match에 전달하며, Match는 이벤트 유형·심각도·작업/자산·유형 조건별로 판정한다.
export interface NotificationFilter {
  min_severity?: string; // "" | low | medium | high | critical
  task_ids?: number[]; // 빈 값=무제한; 비어 있지 않으면 취약점이 속한 작업과 교집합 필요
  asset_ids?: number[]; // 빈 값=무제한; 비어 있지 않으면 취약점이 앵커한 자산과 교집합 필요
  vulnclass_include?: string[]; // 빈 값=전부 수신; 비어 있지 않으면 취약점 유형이 키워드 중 하나에 매칭 필요(대소문자 무시 부분 문자열)
  vulnclass_exclude?: string[]; // 키워드 중 하나라도 매칭되면 제외(제외가 포함보다 우선)
  on_status_change?: boolean; // 취약점 처리 상태 변경 이벤트도 수신할지
}

// NotificationChannel은 하나의 채널 인스턴스. config의 필드는 kind에 따라 다르며,
// 자격 증명 필드는 읽을 때 "__masked__"로 시작하는 마스킹 값으로 대체됨——그대로 되돌려 보내면 '변경 안 함'을 뜻한다.
export interface NotificationChannel {
  id: number;
  name: string;
  kind: string;
  enabled: boolean;
  mode: "realtime" | "digest";
  config: Record<string, unknown>;
  filter: NotificationFilter;
  rate_per_min: number;
  created_at: string;
  updated_at: string;
  // secret_keys는 백엔드가 채널 유형별로 제공하며, 프런트는 이를 바탕으로 비밀번호 상자와 '비우면 변경 안 함' 안내를 렌더링하고,
  // 어떤 채널 지식도 하드코딩하지 않는다.
  secret_keys: string[];
}

// NotificationKind는 /api/notify/meta가 반환하는 채널 유형 메타데이터.
export interface NotificationKind {
  kind: string;
  default_rate_per_min: number;
  secret_keys: string[];
}

export interface NotificationMeta {
  kinds: NotificationKind[];
  enabled: boolean;
  public_base_url: string;
  digest_interval_min: string;
  defaults: { digest_interval_min: number };
  stats: {
    channels: number;
    channels_on: number;
    pending: number;
    failed: number;
    sent_today: number;
    backlog_age_ms: number;
  };
}

// NotificationDelivery는 하나의 전송 기록, 전송 이력과 실패 재전송에 사용.
export interface NotificationDelivery {
  id: number;
  finding_id: string;
  event_kind: string; // finding_created | finding_status_changed
  channel_id: number;
  channel_name: string;
  channel_kind: string;
  state: "pending" | "sending" | "sent" | "failed" | "skipped";
  attempts: number;
  last_error: string;
  batch_id?: number;
  created_at: string;
  sent_at?: string;
  next_attempt_at: string;
  title: string;
  severity: string;
}

// ---- LLM config ----
export interface LLMProfile {
  id: string;
  name: string;
  format: "openai" | "anthropic" | "openai-responses";
  base_url?: string;
  proxy?: string;
  model: string;
  api_key_hint?: string;
  rate_per_second: number;
  rate_per_minute: number;
  context_window_k?: number;
  // 사고 스위치(thinking.type): ""=전송 안 함(기본) | "disabled"=끔 | "enabled"=켬
  thinking_type?: string;
  // 사고 강도: ""=전송 안 함(기본) | "low"/"medium"/"high"/"xhigh"/"max"
  reasoning_effort?: string;
  is_default: boolean;
  // 순환 전환 순위: 클수록 먼저 선택됨. 활성 설정은 항상 체인 맨 앞이며 이 값과 무관.
  priority?: number;
  // true = 장애 조치 대상으로 쓰지 않음(여전히 agent/작업이 명시적으로 바인딩해 사용 가능).
  pool_exclude?: boolean;
  // true(기본) = 스트리밍(SSE) | false = 진짜 비스트리밍(stream:false, 한 번에 반환).
  streaming?: boolean;
  // 단일 응답의 출력 상한(token). 0 = 이 필드 전송 안 함, 서버 기본값으로 결정.
  // context_window_k와 구분 주의: 후자는 모델 총 용량이며 로컬에서 압축 임계값으로만 사용.
  max_tokens?: number;
  // 상한에 어떤 요청 필드명을 쓸지, format="openai"일 때만 의미 있음:
  // ""=max_tokens(기본) | "max_completion_tokens"(OpenAI 추론 모델은 이것만 인식)
  max_tokens_field?: string;
  // 커스텀 세션 헤더명: 비어 있지 않으면 매 요청에 이 HTTP 헤더 포함, 헤더값=현재 세션/의도의 session id.
  // ""=전송 안 함. session-id 헤더로 프롬프트 캐시/고정 라우팅을 하는 게이트웨이용.
  session_header_key?: string;
  // 이 설정의 재시도 오버라이드(연결/빈 응답/동일 provider 안전 구간). 비움/전부 0 = 전역 정책 따름.
  retry?: LLMRetryOverride;
}

// ---- LLM 재시도 정책 ----
// 한 계층 재시도의 두 손잡이. 둘 다 '0 = 미설정':
//   attempts    0=기본 횟수 사용 | -1=이 계층 재시도 끔 | >0=재시도 횟수
//   interval_ms 0=기본 지수 백오프 사용 | >0=이 고정 밀리초 간격 사용
export interface LLMRetryRule {
  attempts: number;
  interval_ms: number;
}

// 단일 LLM 설정이 오버라이드할 수 있는 3개 계층(모두 '엔드포인트를 따라가는' 재시도).
export interface LLMRetryOverride {
  connect: LLMRetryRule; // 연결 재시도: 연결 리셋/타임아웃/429/5xx, 스트림 시작 전
  empty: LLMRetryRule; // 빈 응답 재시도: 완료했지만 내용이 전혀 없음(openai 형식만)
  stream: LLMRetryRule; // 동일 provider 안전 구간 재시도: 출력 전달 전 스트림 끊김 재생
}

// 전역 정책 = 위 3개 계층의 기본값 + 전역에만 있는 2개 계층:
//   breaker 순환 전환 회로 차단(attempts=연속 몇 회 순간 실패 시 차단, interval_ms=고정 재시도 대기 시간)
//   intent  의도 재실행(worker가 model_error로 끝난 뒤 의도 전체 재실행)
export interface LLMRetryPolicy extends LLMRetryOverride {
  breaker: LLMRetryRule;
  intent: LLMRetryRule;
}

// ---- LLM 순환 전환(장애 조치) ----
// 한 설정의 순환 전환 체인 내 위치와 상태(health). state:
//   ok       정상
//   degraded 연속 실패 있으나 회로 차단 임계값 미달
//   tripped  회로 차단됨, 재시도 대기 기간 동안 건너뜀(cooldown_secs는 남은 초)
export interface LLMPoolMember {
  profile_id: string;
  name: string;
  model: string;
  format: string;
  priority: number;
  active: boolean; // 현재 활성 설정 여부(항상 체인 맨 앞)
  excluded: boolean; // pool_exclude: 순환 전환에 참여하지 않음
  state: "ok" | "degraded" | "tripped";
  fails: number;
  trips: number;
  cooldown_secs: number;
  last_error?: string;
  last_at?: string;
}

export interface LLMPoolStatus {
  enabled: boolean;
  bind_fallback: boolean;
  chain: LLMPoolMember[];
}

// ---- Agents ----
export interface Agent {
  id: string;
  key: string; // 내장은 goals/planner/mainagent/worker; 커스텀은 사용자 지정 key
  name: string;
  description?: string;
  role: string;
  builtin: boolean;
  enabled: boolean;
  llm_profile_id?: number | null; // 바인딩된 LLM 설정; null/absent = 작업/세션/전역 따름
  max_turns?: number; // 0 = 무제한
  run_seconds?: number; // worker 단일 실행 벽시계 상한(초); 0 = 무제한
  web_search?: boolean; // 네트워크 검색 활성화 여부(시스템 전역 스위치로 게이팅)
  interactive_shell?: boolean; // 인터랙티브 shell 활성화 여부(지속 PTY 세션 도구군)
  // P3 트리거 후 처리 정책(커스텀 agent만 의미 있음)
  trigger_run_mode?: "serial" | "parallel"; // 직렬 대기 / 트리거마다 각각 세션 하나씩 동시 실행
  trigger_merge_mode?: "by_task" | "all" | "none"; // serial만: 같은 작업 병합 / 전체 병합 / 병합 안 함
  trigger_max_parallel?: number; // parallel만: agent당 동시 실행 상한; 0=무제한
  // 바인딩 수(리스트 API만 반환): 노출 MCP / 노출 Skill / 바인딩 도구
  mcp_count?: number;
  skill_count?: number;
  tool_count?: number;
}

export interface PromptVar {
  name: string;
  description: string;
  example: string;
  source: "exploration" | "runtime" | "distilled";
}

export interface PromptVersion {
  version: number;
  ts: string;
  note: string;
  template_text: string;
}

export interface AgentDetail {
  agent: Agent;
  prompt: string;
  variables: PromptVar[];
  versions: PromptVersion[];
  visibility: { mcp: number[]; skill: string[] };
  // 바인딩 가능한 LLM 설정 후보('기본 모델' 드롭다운용); 현재 바인딩은 agent.llm_profile_id 참고
  llm_profiles?: { id: number; name: string; model: string; is_default: boolean }[];
  wrapup_prompt?: string; // 저장된 마무리 프롬프트(빈 값=내장 기본값 사용)
  wrapup_default?: string; // 내장 기본 마무리 프롬프트(플레이스홀더/기본값 복원)
  wrapup_max_turns?: number; // 저장된 마무리 횟수(0=내장 기본값 사용)
  wrapup_max_turns_default?: number; // 내장 기본 마무리 횟수("0=기본 N" 안내용)
  // 작업 수준 타임아웃 마무리 문구(worker/planner만, task_timeout_wrapup_supported=true일 때만 이 구역 표시)
  task_timeout_wrapup_supported?: boolean;
  task_timeout_wrapup_prompt?: string;
  task_timeout_wrapup_default?: string;
  task_timeout_wrapup_max_turns?: number;
  task_timeout_wrapup_max_turns_default?: number;
}

// ---- MCP ----
export interface MCPServer {
  id: number;
  name: string;
  transport: "stdio" | "http" | "sse";
  command?: string;
  args: string[];
  env: Record<string, string>;
  url?: string;
  enabled: boolean;
  insecure?: boolean; // http: skip TLS cert verification (self-signed servers)
  tools?: string[]; // mcp_tools_cache (names only, for the count)
}

export interface MCPTool {
  name: string;
  description: string;
}

// ---- Skills ----
// Fields align with the agentskills.io open specification.
// description covers both "what the skill does" and "when to use it".
export interface SkillItem {
  name: string; // unique key = directory name
  description?: string; // required per spec; covers what + when to use
  license?: string; // optional: SPDX identifier or free text
  compatibility?: string; // optional: environment requirements
  mcps?: string[]; // MCP server names this skill unlocks on load
  files: string[]; // files in the skill directory
  // 호출 통계(skill_usage 장부). 한 번도 호출된 적 없는 skill: calls=0, last_used 생략.
  calls: number;
  tasks: number; // 이를 로드한 작업 수(chat 세션은 제외)
  usage_agents: string[]; // 이를 로드한 agent key
  last_used?: string;
}

// SkillCall은 한 번의 Skill() 호출(단일 skill의 최근 호출 목록).
export interface SkillCall {
  ts: string;
  agent_key: string;
  task_id: number; // 0 = 비작업 상황(대화 세션)
  session_id: string;
  args_len: number;
}

// MissingSkill은 호명되었지만 존재하지 않는 skill —— "쓰려 했으나 없는" 공백.
export interface MissingSkill {
  skill: string;
  calls: number;
  agents: string[];
  last_used?: string;
}

// ---- Tools (내장 도구 카탈로그) ----
// key + handler live in Go; only these fields are page-editable. system tools lock
// the key and the parameter *structure* (name/type/required) — the per-param
// description/default and the agent binding are what move.
export interface Tool {
  key: string;
  system: boolean;
  description: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  schema: Record<string, any>; // full JSON-Schema (object with properties)
  agents: string[]; // bound agent keys
  enabled: boolean;
  kind?: "builtin" | "shell" | "command" | "script" | "http"; // 커스텀 도구 유형
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  exec?: Record<string, any>; // 커스텀 도구 실행 규격(kind!=builtin)
  deferred?: boolean; // schema 지연(SearchExtraTools/ExecuteExtraTool)
  calls?: number; // persistent runtime invocation count (older APIs may omit it)
}

// ---- Stats ----
export interface Stats {
  assets: number;
  engine_mode: EngineMode;
  llm_configured: boolean;
  roe_enabled: boolean;
  findings_confirmed: number;
  active_task?: Partial<Task>;
}

// ---- Intercept Rules ----
export type InterceptAction = "allow" | "deny" | "ask";
export type InterceptMatchTarget = "tool_name" | "tool_input";
export type InterceptMatchType = "string" | "regex";

export interface InterceptRule {
  id: number;
  name: string;
  enabled: boolean;
  priority: number;
  match_target: InterceptMatchTarget;
  match_type: InterceptMatchType;
  pattern: string;
  action: InterceptAction;
  message: string;
  timeout_enabled: boolean;
  timeout_seconds: number;
  timeout_action: "deny" | "allow";
  created_at: string;
  updated_at: string;
}

// ---- Asset Intercept Rules(자산 인터셉트: 전역 블랙리스트) ----
export type AssetInterceptKind =
  | "exact_domain"
  | "exact_ip"
  | "exact_url"
  | "fuzzy_domain"
  | "fuzzy_ip"
  | "fuzzy_url"
  | "cidr";

// action은 작업 수준 규칙에만 사용: block=차단(테스트 금지) allow=허용(화이트리스트).
export type AssetInterceptAction = "block" | "allow";

// 작업 수준 자산 인터셉트/허용 규칙의 입력 항목(작업 생성, 작업 상세 편집에 사용).
export interface AssetInterceptRuleInput {
  action: AssetInterceptAction;
  kind: AssetInterceptKind;
  pattern: string;
  note: string;
  enabled: boolean;
}

export interface AssetInterceptRule {
  id: number;
  enabled: boolean;
  action?: AssetInterceptAction; // 전역 규칙은 이 필드 없음(항상 차단); 작업 수준 규칙은 block/allow 구분
  kind: AssetInterceptKind;
  pattern: string;
  note: string;
  builtin: boolean;
  created_at: string;
  updated_at: string;
}

export interface InterceptPending {
  decision_source?: "rule" | "model" | "unknown" | "";
  id: number;
  rule_id?: number;
  conversation_id?: number;
  task_id?: string;
  agent_name: string;
  tool_name: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  tool_input: Record<string, any>;
  status: "pending" | "allowed" | "denied" | "timeout";
  reason: string; // 규칙 message 또는 모델 판정 사유(모델 판정은 '[模型]' 접두사 포함)
  decided_at?: string;
  created_at: string;
}

// JudgeConfig: 모델 보완 판정(인터셉트 규칙이 하나도 매칭되지 않을 때만 모델이 판단)의 전역 설정.
export interface JudgeConfig {
  enabled: boolean;
  profile_id: number; // 0 = 활성/기본 설정 따름
  prompt: string; // 판정 프롬프트; GET 시 미설정이면 백엔드가 내장 템플릿 전문을 채워 넣음
  timeout_seconds: number; // 모델 호출 타임아웃
  fail_action: "allow" | "ask" | "deny"; // 모델 오류/타임아웃/파싱 불가 시의 대체 처리
  ask_timeout_seconds: number; // 모델이 ask 판정해 수동 전환된 뒤의 승인 대기 타임아웃
  ask_timeout_action: "allow" | "deny"; // 승인 타임아웃 후의 기본 동작
}

export interface InterceptApprovalFilter {
  status?: InterceptPending["status"];
  decision_source?: "rule" | "model" | "unknown";
}

// InterceptApprovalRow enriches InterceptPending with conversation/task and rule context.
export interface InterceptApprovalRow extends InterceptPending {
  conv_title: string; // "" if no linked conversation
  conv_agent_key: string; // "" if no linked conversation
  rule_name: string; // "" if rule was deleted
}

// ── 자산 동기화 (ScopeSentry 데이터 소스) ──────────────────────────────────────────────
export interface SSProject {
  id: string; // MongoDB ObjectID — used as filter.project
  name: string;
  logo?: string;
  AssetCount?: number;
  tag?: string;
}

export interface SSTask {
  id: string;
  name: string; // used as filter.task
  status?: number;
  progress?: number;
  creatTime?: string;
  endTime?: string;
}

// ConvTokenSummary — one conversation's token total (+ profile/date) for merging
// chat usage into the dashboard token stats. GET /api/tokens/conversations.
export interface ConvTokenSummary {
  llm_profile_id: number | null;
  created_at: string;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens: number;
  cache_write_tokens: number;
}

// ---- Command recording (Bash execution history) ----
export interface CommandRecord {
  id: number;
  exploration_id: number;
  worker: string;
  tool: string;
  command: string; // raw tool input (JSON)
  output: string;
  is_error: boolean;
  created_at: string;
}

// 단일 도구의 호출 통계(/commands/stats); errors는 그중 실패 횟수.
export interface ToolStat {
  tool: string;
  total: number;
  errors: number;
}

// ---- LLM recording ----
export interface LLMRecordItem {
  id: number;
  ts: string;
  model: string;
  profile_name: string;
  session_id: string;
  task_id: string;
  worker: string;
  latency_ms: number;
  input_tokens: number;
  output_tokens: number;
  cache_read: number;
  cache_write: number;
  status: string;
  error?: string;
}

export interface LLMRecordDetail extends LLMRecordItem {
  request_body: string;
  response_body: string;
  // provider가 실제로 주고받은 HTTP 원문: 요청은 buildBody()가 보낸 전체 body(도구
  // schema 포함), 응답은 원본 SSE 프레임. 위의 request_body/response_body는 정규화 뷰이며,
  // 도구 schema와 tool_use 블록을 버림. 오래된 기록은 빈 값.
  raw_request?: string;
  raw_response?: string;
}

// One distinct task with its LLM-record count (task picker on the records page).
export interface LLMTask {
  task_id: string;
  count: number;
}

// The exact JSON sent to the review model, retained for all model verdicts.
export interface InterceptReviewInput {
  version: number;
  background?: {
    // worker_summary is retained only for immutable v2/v3 snapshots.
    source: "user_message" | "worker_summary";
    text: string;
    truncated?: boolean;
  };
  // Version 1 snapshots are immutable and remain readable in historical audits.
  task?: {
    task_id: number;
    description: string;
    goal: string;
    constraints: { id: number; kind: string; text: string; origin: string; created_at: number }[];
    truncated?: boolean;
  };
  working_directory?: string;
  worker_intent?: string;
  turn_input?: string;
  background_truncated?: boolean;
  // Legacy v1/v2 snapshots only; v3 never sends execution history.
  history?: {
    tool_use_id: string;
    tool: string;
    arguments_preview: string;
    result: string;
    status: "succeeded" | "failed";
    truncated?: boolean;
  }[];
  history_truncated?: boolean;
  correlation?: "exact" | "ambiguous" | "unavailable";
  tool_name: string;
  arguments: Record<string, unknown>;
}

// Immutable review snapshot plus separately recorded execution outcome.
export interface InterceptAudit {
  model_input?: InterceptReviewInput;
  model_input_digest?: string;
  run_id?: string;
  tool_use_id?: string;
  correlation: "exact" | "ambiguous" | "unavailable";
  input_digest: string;
  user_message: string;
  user_truncated?: boolean;
  context:
    | { kind: string; tool?: string; tool_use_id?: string; text: string; is_error?: boolean; truncated?: boolean }[]
    | null;
  context_truncated?: boolean;
  captured_at: string;
  model_fallback?: boolean;
  initial_action: "allow" | "ask" | "deny";
  initial_reason: string;
  effective_action?: "allow" | "deny";
  decision_reason?: string;
  rule_name?: string;
  config_digest?: string;
  profile_id?: number;
  execution_status: "not_started" | "not_executed" | "awaiting_result" | "succeeded" | "failed" | "unknown";
  output?: string;
  output_truncated?: boolean;
  execution_ended_at?: string;
}
export interface InterceptDetail extends InterceptApprovalRow {
  audit: InterceptAudit | null;
}

export type TrafficEvidenceRole = "baseline" | "proof" | "verification" | "supporting";
export interface TrafficEvidenceRef {
  traffic_id: string;
  role?: TrafficEvidenceRole;
  note?: string;
}
export interface TrafficEvidenceSnapshot {
  id: string;
  source_traffic_id: string;
  captured_at: number;
  url: string;
  method: string;
  status: number;
  content_type: string;
  req_head?: string;
  resp_head?: string;
  req_hash: string;
  resp_hash: string;
  req_len: number;
  resp_len: number;
}
export interface FindingTrafficBinding {
  id: string;
  finding_id: string;
  snapshot_id: string;
  role: TrafficEvidenceRole;
  note: string;
  position: number;
  created_at: string;
  snapshot: TrafficEvidenceSnapshot;
}
export interface FindingTraffic {
  finding_id: string;
  version: number;
  report_version: number;
  bindings: FindingTrafficBinding[];
}
export interface EvidenceBodyPreview {
  content: string;
  offset: number;
  total: number;
  next_offset: number;
  truncated: boolean;
  binary: boolean;
}
export interface FindingTrafficDetail {
  binding: FindingTrafficBinding;
  request: EvidenceBodyPreview;
  response: EvidenceBodyPreview;
}

/** GET /api/update/check —— 현재 버전과 GitHub 최신 정식 버전의 비교 결과. */
export interface UpdateCheck {
  /** 현재 실행 중인 버전; 개발 빌드는 "dev" 또는 git describe의 접미사 형태. */
  current: string;
  /** 실행 형태. docker에서는 교체가 컨테이너 쓰기 가능 레이어에만 적용되어, 컨테이너 재생성 시 이미지 버전으로 되돌아간다. */
  mode: "docker" | "binary";
  os: string;
  arch: string;
  repo: string;
  /** 롤백 가능한 이전 버전(artex.old) 존재 여부. */
  has_backup: boolean;
  /** 이번 시작 시 자가 업데이트 부트스트랩의 결론(교체 실패 / 롤백됨 등), 아무 일 없으면 빈 값. */
  boot_notice?: string;
  rolled_back?: boolean;
  /** GitHub 조회 실패 시 사유 제공, 이때 아래 필드는 모두 없음. */
  error?: string;
  latest?: string;
  notes?: string;
  html_url?: string;
  published_at?: string;
  /** 현재 플랫폼에 해당하는 릴리스 패키지명, 그리고 해당 Release가 실제로 그것을 포함하는지. */
  asset?: string;
  asset_available?: boolean;
  size?: number;
  has_update?: boolean;
  /** 양쪽 버전 번호 비교 가능 여부; 개발 빌드는 false이며 이때 원클릭 업데이트 비활성화. */
  comparable?: boolean;
  /** comparable이 false일 때의 설명. */
  reason?: string;
}

/** /api/update/stream이 푸시하는 업데이트 진행 1건. */
export interface UpdateProgress {
  phase: "idle" | "downloading" | "verifying" | "extracting" | "staged" | "failed";
  /** 다운로드 단계만 의미 있음(0-100); 나머지 단계는 -1. */
  percent: number;
  message: string;
  version?: string;
  error?: string;
}

// Original execution selected from an approval, never submitted to the reviewer.
export interface InterceptExecution {
  conversation_id: number | null;
  task_id: string | null;
  session: string;
  seq: number;
  items: Activity[];
}
