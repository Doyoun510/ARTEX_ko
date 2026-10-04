"use client";

// LLM 재시도 설정 공통 컴포넌트: 다섯 재시도 계층별 '횟수 + 간격'.
//
// 다섯 계층의 안쪽부터 바깥쪽 순서: 연결 수립(SDK) → 빈 응답(SDK) → 같은 provider의 안전 구간 → 순환 전환 회로 차단 → 의도 재실행.
// 앞의 세 계층은 엔드포인트별로 적용되므로 각 모델 설정에서 전역 기본값을 재정의할 수 있습니다. 뒤의 두 계층은 프로세스 수준으로 전역 설정만 있습니다.
//
// 모든 입력은 백엔드 db.RetryRule과 같은 '빈 값 = 설정 안 함' 의미를 따릅니다:
//   횟수   빈 값/0 = 내장 기본값 사용 | -1 = 이 계층 재시도 끄기 | >0 = 지정한 횟수 사용
//   간격   빈 값/0 = 이 계층의 기존 지수 백오프 사용 | >0 = 지정한 고정 밀리초 간격 사용

import * as React from "react";

import { Loader2Icon, SaveIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import type { LLMRetryOverride, LLMRetryPolicy, LLMRetryRule } from "@/lib/types";

export const ZERO_RULE: LLMRetryRule = { attempts: 0, interval_ms: 0 };
export const ZERO_OVERRIDE: LLMRetryOverride = {
  connect: ZERO_RULE,
  empty: ZERO_RULE,
  stream: ZERO_RULE,
};
const ZERO_POLICY: LLMRetryPolicy = {
  ...ZERO_OVERRIDE,
  breaker: ZERO_RULE,
  intent: ZERO_RULE,
};

type LayerMeta = {
  title: string;
  /** 이 계층의 재시도 위치와 실행 주체 */
  where: string;
  /** 이 계층에 해당하는 오류를 상태 코드까지 명시 */
  trigger: string;
  /** 비슷하지만 이 계층에 해당하지 않는 오류. 설정해도 반응이 없어 bug로 오해하지 않도록 설명 */
  skips?: string;
  desc: string;
  attemptsLabel: string;
  /** 횟수가 비어 있을 때의 기본값. 자리표시자에 사용 */
  defAttempts: number;
  /** 간격이 비어 있을 때의 기본 정책. 자리표시자에 사용 */
  defInterval: string;
  /** 횟수 -1의 의미 */
  offHint: string;
};

export const RETRY_LAYERS = {
  connect: {
    title: "연결 수립 재시도",
    where: "SDK · 200 수신 전",
    trigger:
      "연결할 수 없거나 아직 200을 받지 못한 경우: 연결 재설정 / 읽기·쓰기 타임아웃 / DNS 실패 등 네트워크 계층 오류 및 HTTP 408, 429, 500, 502, 503, 504.",
    skips: "나머지 상태 코드(400 / 401 / 403 / 404 / 413 / 422 등)는 확정적인 거부로 재전송해도 실패하므로 즉시 상위 계층으로 오류를 전달합니다.",
    desc: "같은 요청을 그대로 재전송합니다. 스트림이 시작된 후(200 수신 후)의 연결 끊김은 이 계층에서 처리하지 않습니다.",
    attemptsLabel: "재시도 횟수",
    defAttempts: 3,
    defInterval: "0.5s→1s→2s 지수적으로 증가(상한 8s)",
    offHint: "-1 = 재시도하지 않고 실패 즉시 상위 계층으로 오류 전달",
  },
  empty: {
    title: "빈 응답 재시도",
    where: "SDK · openai 형식만",
    trigger:
      "HTTP 200이며 finish_reason은 정상적인 stop이지만 응답에 내용 블록이 전혀 없는 경우입니다. 게이트웨이의 빈 프레임, 사고 필드 프레임 누락, 샘플링 일시 오류 등이 원인입니다.",
    skips: "max_tokens에 의해 잘려 내용이 없는 경우는 제외합니다(출력 상한을 높여야 해결되며 재전송하면 같은 문제가 반복됩니다).",
    desc: "전체 prompt를 재전송하므로 긴 컨텍스트에서는 비용이 높습니다. 횟수를 너무 높게 설정하지 마세요.",
    attemptsLabel: "재시도 횟수",
    defAttempts: 2,
    defInterval: "0.5s→1s→2s 지수적으로 증가(상한 8s)",
    offHint: "-1 = 빈 응답을 그대로 전달",
  },
  stream: {
    title: "같은 provider의 안전 구간 재시도",
    where: "이 프로젝트 · 출력 전달 전",
    trigger:
      "스트림 수립 후(200 수신 후) 발생한 문제: 연결 중단, 공급자 overloaded, 스트림 내부의 429 / 5xx 오류 이벤트이며, 호출자에게 아직 token을 하나도 전달하지 않은 경우입니다.",
    skips:
      "할당량 소진(402 / insufficient_quota, 순환 전환에서 설정 전환), 컨텍스트 길이 초과(413 / context length, 압축으로 처리), 400 / 401 / 403 / 404 / 422의 확정적인 거부는 모두 재시도하지 않습니다.",
    desc: "같은 설정에서 같은 요청을 다시 실행합니다. 출력을 아직 전달하지 않았으므로 모델 출력이나 도구 실행이 중복되지 않습니다.",
    attemptsLabel: "재시도 횟수",
    defAttempts: 2,
    defInterval: "0.5s→1s 지수적으로 증가(상한 4s)",
    offHint: "-1 = 스트림 중단을 바깥 계층의 의도 재실행에 바로 전달",
  },
  breaker: {
    title: "순환 전환 회로 차단",
    where: "이 프로젝트 · 프로세스 수준, 전역 설정만",
    trigger:
      "일시적 실패(429, 5xx, 네트워크 오류)가 연속으로 기준값에 도달하면 회로 차단이 적용됩니다. 잔액 부족(402), 키 무효(401 / 403), 모델 없음(404) 같은 확정적 실패는 기준값과 무관하게 첫 실패부터 회로 차단이 적용됩니다.",
    skips: "한 번 성공하면 누적 횟수를 초기화하므로 간헐적으로 오류가 발생하는 설정은 실패가 쌓여 회로 차단이 적용되지 않습니다.",
    desc: "회로 차단 후 재시도 대기 상태로 전환되며 재시도 대기 시간 동안에는 순환 전환에서 이 설정을 건너뜁니다. 상태를 저장하므로 재시작해도 유지됩니다.",
    attemptsLabel: "회로 차단 기준 연속 실패 횟수",
    defAttempts: 3,
    defInterval: "1min→5min→30min 단계",
    offHint: "-1 = 일시적 실패는 회로 차단하지 않음(확정적 실패는 여전히 회로 차단)",
  },
  intent: {
    title: "의도 재실행",
    where: "이 프로젝트 · 프로세스 수준, 전역 설정만",
    trigger:
      "앞의 계층에서 처리하지 못한 경우: worker가 model_error로 종료됩니다. 내부 재시도를 모두 소진했거나 스트림에서 출력 전달을 시작한 뒤 끊긴 경우입니다(이때 재실행은 안전하지 않으므로 전체 의도를 다시 실행해야 합니다).",
    skips: "할당량 소진은 이미 순환 전환에서 설정을 전환해 처리하므로 여기서 재실행하지 않습니다. 작업이 일시 중지 / 종료 / 마무리 상태가 되면 즉시 양보하여 백오프 시간을 사용하지 않습니다.",
    desc: "의도 전체를 처음부터 다시 실행합니다. 가장 바깥 계층이므로 한 번 재실행할 때마다 내부 계층의 횟수가 다시 곱해집니다.",
    attemptsLabel: "재실행 횟수",
    defAttempts: 2,
    defInterval: "고정 3s",
    offHint: "-1 = 재실행하지 않고 해당 의도를 바로 blocked로 판정",
  },
} satisfies Record<string, LayerMeta>;

type LayerKey = keyof typeof RETRY_LAYERS;

/** 밀리초를 읽기 쉬운 단위로 변환하여 입력란 옆에만 표시합니다. */
function humanMs(ms: number) {
  if (!Number.isFinite(ms) || ms <= 0) return "";
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60_000) return `${Number((ms / 1000).toFixed(2))}s`;
  return `${Number((ms / 60_000).toFixed(2))}min`;
}

/** 제어형 숫자 입력: 빈 문자열 ↔ 0. 중간 입력 상태("-", "1e")는 로컬에 그대로 두어 부모에 영향을 주지 않습니다. */
function NumField({
  id,
  value,
  onChange,
  placeholder,
  min,
}: {
  id: string;
  value: number;
  onChange: (n: number) => void;
  placeholder: string;
  min: number;
}) {
  const [text, setText] = React.useState(value === 0 ? "" : String(value));
  // 부모의 값 전체가 바뀌면(정책 불러오기, 설정 전환) 따라 갱신합니다. 직접 입력할 때는 여기로 들어오지 않습니다.
  // 그때는 value가 이미 로컬 텍스트의 parse 결과와 같기 때문입니다.
  React.useEffect(() => {
    const incoming = value === 0 ? "" : String(value);
    setText((cur) => (Number(cur || 0) === value ? cur : incoming));
  }, [value]);
  return (
    <Input
      id={id}
      type="number"
      min={min}
      className="w-28 shrink-0"
      value={text}
      placeholder={placeholder}
      onChange={(e) => {
        setText(e.target.value);
        const n = Number(e.target.value);
        onChange(e.target.value.trim() === "" || !Number.isFinite(n) ? 0 : Math.trunc(n));
      }}
    />
  );
}

/** 한 재시도 계층의 두 설정 항목입니다. idPrefix는 한 페이지에 여러 번 표시될 때 label의 htmlFor를 유지합니다. */
export function RetryRuleFields({
  layer,
  idPrefix,
  value,
  onChange,
  compact,
}: {
  layer: LayerKey;
  idPrefix: string;
  value: LLMRetryRule;
  onChange: (r: LLMRetryRule) => void;
  /** true = 설정 패널의 간략 표시: 상세 설명을 생략하고 이 계층에 해당하는 오류만 표시 */
  compact?: boolean;
}) {
  const meta = RETRY_LAYERS[layer];
  const human = humanMs(value.interval_ms);
  return (
    <div className={compact ? "grid gap-2" : "grid gap-3 rounded-lg border p-3"}>
      <div className="grid gap-0.5">
        <div className="flex flex-wrap items-baseline gap-2">
          <Label className="text-sm">{meta.title}</Label>
          <span className="text-muted-foreground text-xs">{meta.where}</span>
        </div>
        {/* 이 계층에 해당하는 오류를 상태 코드까지 명시합니다. 설정해도 효과가 없다면 대부분 해당 계층의 오류가 아닙니다. */}
        <p className="text-muted-foreground text-xs">
          <span className="font-medium text-foreground">발동 조건</span>: {meta.trigger}
        </p>
        {!compact && meta.skips && (
          <p className="text-muted-foreground text-xs">
            <span className="font-medium text-foreground">이 계층 제외</span>: {meta.skips}
          </p>
        )}
        {!compact && <p className="text-muted-foreground text-xs">{meta.desc}</p>}
      </div>
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-n`} className="text-muted-foreground text-xs">
            {meta.attemptsLabel}
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-n`}
            min={-1}
            value={value.attempts}
            placeholder={`기본값 ${meta.defAttempts}`}
            onChange={(n) => onChange({ ...value, attempts: n })}
          />
        </div>
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-ms`} className="text-muted-foreground text-xs">
            간격 ms
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-ms`}
            min={0}
            value={value.interval_ms}
            placeholder="기본 백오프"
            onChange={(n) => onChange({ ...value, interval_ms: n })}
          />
          <span className="text-muted-foreground text-xs">{human ? `고정 ${human}` : meta.defInterval}</span>
        </div>
      </div>
      {!compact && <p className="text-muted-foreground text-xs">빈 값 = 기본값 사용; {meta.offHint}.</p>}
    </div>
  );
}

/** 모델 설정 패널에서 재정의할 세 계층(엔드포인트별로 적용되는 계층). */
export function ProfileRetryFields({
  value,
  onChange,
}: {
  value: LLMRetryOverride;
  onChange: (o: LLMRetryOverride) => void;
}) {
  return (
    <div className="grid gap-3 rounded-lg border p-3">
      <div className="grid gap-0.5">
        <Label className="text-sm">재시도 재정의</Label>
        <p className="text-muted-foreground text-xs">
          이 설정에만 적용되며 '재시도 및 백오프'의 전역 기본값을 재정의합니다. 빈 값 = 전역 설정 사용; 횟수 -1 = 이 계층 재시도 끄기;
          간격을 입력하면 지수 백오프 대신 고정 간격을 사용합니다. 회로 차단과 의도 재실행은 프로세스 수준이므로 전역 설정 페이지에서만 조정할 수 있습니다.
        </p>
      </div>
      {(["connect", "empty", "stream"] as const).map((k) => (
        <div key={k} className="border-t pt-3 first:border-t-0 first:pt-0">
          <RetryRuleFields
            compact
            layer={k}
            idPrefix="pf"
            value={value[k]}
            onChange={(r) => onChange({ ...value, [k]: r })}
          />
        </div>
      ))}
    </div>
  );
}

/** '재시도 및 백오프' tab: 다섯 계층의 전역 기본값. */
export function RetryPolicyPanel() {
  const [policy, setPolicy] = React.useState<LLMRetryPolicy>(ZERO_POLICY);
  const [loading, setLoading] = React.useState(true);
  const [saving, setSaving] = React.useState(false);

  const load = React.useCallback(async () => {
    setLoading(true);
    try {
      const p = await api.llmRetryPolicy();
      setPolicy({ ...ZERO_POLICY, ...p });
    } catch (e) {
      toast.error(`재시도 정책 불러오기 실패: ${(e as Error).message}`);
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void load();
  }, [load]);

  async function save() {
    if (saving) return;
    setSaving(true);
    try {
      // 백엔드는 범위를 벗어난 값을 범위 안으로 조정하여 반환합니다. 반환값으로 갱신하므로 표시값과 저장값이 일치합니다.
      const saved = await api.saveLLMRetryPolicy(policy);
      setPolicy({ ...ZERO_POLICY, ...saved });
      toast.success("저장했습니다. 즉시 적용됩니다(진행 중인 호출은 기존 파라미터 사용)");
    } catch (e) {
      toast.error(`저장 실패: ${(e as Error).message}`);
    } finally {
      setSaving(false);
    }
  }

  const set = (k: LayerKey) => (r: LLMRetryRule) => setPolicy((p) => ({ ...p, [k]: r }));

  if (loading) {
    return (
      <div className="flex items-center gap-2 rounded-lg border border-dashed p-10 text-muted-foreground text-sm">
        <Loader2Icon className="size-4 animate-spin" /> 재시도 정책 불러오는 중…
      </div>
    );
  }

  return (
    <div className="grid gap-4">
      <div className="rounded-lg border bg-muted/30 p-3 text-muted-foreground text-xs leading-relaxed">
        모델 호출이 실패하면 안쪽부터 바깥쪽으로 다섯 재시도 계층을 거칩니다:
        <span className="text-foreground"> 연결 수립 → 빈 응답 → 같은 provider의 안전 구간 → 순환 전환 회로 차단 → 의도 재실행</span>
        . 내부 계층을 소진해야 바깥 계층으로 넘어가므로 횟수는
        <span className="text-foreground"> 서로 곱해집니다</span>
        . 각 계층을 최대로 설정하면 한 번의 일시적 장애에 수십 번의 요청을 소모할 수 있습니다.
        모두 비워 두면 현재 기본값을 사용하며 이 페이지가 없을 때와 동작이 같습니다. 앞의 세 계층은 각 모델 설정에서 따로 재정의할 수 있습니다.
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {(Object.keys(RETRY_LAYERS) as LayerKey[]).map((k) => (
          <RetryRuleFields key={k} layer={k} idPrefix="gl" value={policy[k]} onChange={set(k)} />
        ))}
      </div>

      <div className="flex gap-2">
        <Button onClick={save} disabled={saving}>
          {saving ? <Loader2Icon className="animate-spin" /> : <SaveIcon />}
          저장
        </Button>
        <Button variant="outline" onClick={() => setPolicy(ZERO_POLICY)} disabled={saving}>
          모두 기본값으로 복원
        </Button>
      </div>
    </div>
  );
}
