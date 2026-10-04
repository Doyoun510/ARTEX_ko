// 채널 필드 표와 설정값 파싱 도구입니다.
//
// 이 파일은 화면이 아니라 **데이터**이므로 페이지와 분리했습니다. 채널별 필드와
// 각 필드에 사용할 입력 요소, 양식 텍스트와 설정값(JSON) 간의 양방향 변환을 정의합니다.
// 파일을 분리하여 새 채널을 추가할 때 페이지 대신 이 파일만 수정하면 됩니다.
// 채널 종류의 표시 이름과 소개입니다. 표시 문구에만 영향을 주므로 백엔드 대신 프런트엔드에 둡니다.
export const KIND_LABEL: Record<string, string> = {
  dingtalk: "DingTalk",
  feishu: "Feishu",
  wecom: "WeCom",
  webhook: "범용 Webhook",
  telegram: "Telegram",
  email: "메일",
};

// 채널별 설정 필드 정의입니다.
//
// 백엔드에서 schema를 전달하는 대신 프런트엔드 필드 표를 유지합니다. 백엔드는
// Validate(필수 입력/형식)만 담당하며 UI는 배치와 입력 요소 종류가 필요하므로 역할이 다릅니다.
// 연결 지점은 secret_keys뿐입니다. 비밀번호 입력란으로 표시할 필드는 백엔드가 지정합니다.
// 자격 증명에 해당하는 값은 채널 구현에서만 알기 때문입니다(WeCom은 Webhook 전체가 자격 증명이고
// DingTalk은 그중 secret만 해당). 새 채널의 항목이 누락되면 양식이 비어 있게 됩니다.
// 오류가 조용히 발생하지는 않습니다(아래의 hasFields에서 안내).
export type FieldKind = "text" | "password" | "number" | "select" | "textarea" | "switch" | "kv" | "list";
export interface FieldDef {
  key: string;
  label: string;
  kind: FieldKind;
  placeholder?: string;
  help?: string;
  options?: { value: string; label: string }[];
}
export const CHANNEL_FIELDS: Record<string, FieldDef[]> = {
  dingtalk: [
    {
      key: "webhook",
      label: "Webhook 주소",
      kind: "text",
      placeholder: "https://oapi.dingtalk.com/robot/send?access_token=...",
    },
    {
      key: "secret",
      label: "서명 키",
      kind: "password",
      help: "봇 보안 설정에서 '서명 추가'를 선택하면 입력하세요. '사용자 지정 키워드'를 선택하거나 보안 설정을 켜지 않았으면 비워 두세요",
    },
  ],
  feishu: [
    {
      key: "webhook",
      label: "Webhook 주소",
      kind: "text",
      placeholder: "https://open.feishu.cn/open-apis/bot/v2/hook/...",
    },
    { key: "secret", label: "서명 검증 키", kind: "password", help: "봇에서 '서명 검증'을 켜면 입력하고, 그 외에는 비워 두세요" },
  ],
  wecom: [
    {
      key: "webhook",
      label: "Webhook 주소",
      kind: "text",
      placeholder: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=...",
    },
  ],
  webhook: [
    { key: "url", label: "대상 URL", kind: "text", placeholder: "https://your-endpoint.example.com/hook" },
    {
      key: "method",
      label: "요청 메서드",
      kind: "select",
      options: [
        { value: "POST", label: "POST(요청 본문 포함)" },
        { value: "PUT", label: "PUT(요청 본문 포함)" },
        { value: "PATCH", label: "PATCH(요청 본문 포함)" },
        { value: "GET", label: "GET(요청 본문 없음)" },
      ],
    },
    { key: "headers", label: "사용자 지정 요청 헤더", kind: "kv", help: "한 줄에 KEY=VALUE, 예: Authorization=Bearer xxx" },
    {
      key: "body_template",
      label: "요청 본문 템플릿",
      kind: "textarea",
      help:
        "비워 두면 내장 기본 템플릿을 사용합니다. 변수: {{.Title}} {{.Batch}} {{.Count}} {{.HomeURL}} {{.SentAt}}, " +
        "그리고 range .Items 안의 .Name/.VulnClass/.Severity/.Summary/.Assets/.DetailURL/.StatusLabel. " +
        "문자열을 삽입할 때는 {{json .Xxx}}를 사용하고 {{.Xxx}}는 사용하지 마세요. 그렇지 않으면 제목의 따옴표로 JSON 형식이 깨집니다.",
    },
  ],
  telegram: [
    { key: "bot_token", label: "Bot Token", kind: "password", placeholder: "123456:ABC-DEF..." },
    { key: "chat_id", label: "Chat ID", kind: "text", placeholder: "-1001234567890" },
    {
      key: "base_url",
      label: "API 주소",
      kind: "text",
      placeholder: "https://api.telegram.org",
      help: "비워 두면 공식 주소를 사용합니다. 자체 Bot API 리버스 프록시를 사용할 때 입력하세요",
    },
  ],
  email: [
    { key: "host", label: "SMTP 서버", kind: "text", placeholder: "smtp.example.com" },
    {
      key: "port",
      label: "포트",
      kind: "number",
      placeholder: "587",
      help: "587은 STARTTLS를 사용합니다. 465는 '암시적 TLS'를 켜세요",
    },
    { key: "username", label: "계정", kind: "text" },
    { key: "password", label: "비밀번호 / 메일 인증 코드", kind: "password" },
    { key: "from", label: "보내는 사람", kind: "text", placeholder: "artex@example.com" },
    { key: "to", label: "받는 사람", kind: "list", help: "여러 주소는 쉼표로 구분하세요" },
    { key: "tls", label: "암시적 TLS", kind: "switch", help: "465 포트에서는 켜고, 587에서는 꺼 두세요(자동으로 STARTTLS 사용)" },
  ],
};

export const SEVERITY_OPTIONS = [
  { value: "", label: "제한 없음" },
  { value: "low", label: "낮음 이상" },
  { value: "medium", label: "중간 이상" },
  { value: "high", label: "높음 이상" },
  { value: "critical", label: "심각만" },
];

export type ChannelForm = {
  name: string;
  kind: string;
  mode: "realtime" | "digest";
  enabled: boolean;
  ratePerMin: string;
  config: Record<string, unknown>;
  minSeverity: string;
  includeText: string;
  excludeText: string;
  taskIDsText: string;
  assetIDsText: string;
  onStatusChange: boolean;
};

export const emptyForm = (kind: string): ChannelForm => ({
  name: "",
  kind,
  mode: "realtime",
  enabled: true,
  ratePerMin: "",
  config: {},
  minSeverity: "",
  includeText: "",
  excludeText: "",
  taskIDsText: "",
  assetIDsText: "",
  onStatusChange: false,
});

// parseKV는 '한 줄에 KEY=VALUE' 형식의 텍스트 영역을 파싱합니다.
export function parseKV(text: string): Record<string, string> {
  const out: Record<string, string> = {};
  for (const line of text.split("\n")) {
    const t = line.trim();
    if (!t) continue;
    const i = t.indexOf("=");
    if (i > 0) out[t.slice(0, i).trim()] = t.slice(i + 1).trim();
  }
  return out;
}
// parseIDs는 쉼표/공백으로 구분된 id 목록을 파싱합니다.
export function parseIDs(text: string): number[] {
  return text
    .split(/[\s,，]+/)
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => Number(s))
    .filter((n) => Number.isFinite(n) && n > 0);
}
// parseKeywords는 줄/쉼표로 구분된 키워드 목록을 파싱합니다(취약점 유형 이름에 공백이 있을 수 있으므로 줄이나 쉼표로 구분).
export function parseKeywords(text: string): string[] {
  return text
    .split(/[\n,，]+/)
    .map((s) => s.trim())
    .filter(Boolean);
}
