"use client";

import { CheckIcon } from "lucide-react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import type { NotificationFilter } from "@/lib/types";

// asText / inputType은 이 파일의 값 변환 보조 함수입니다(입력 요소 표시와 밀접한 관계). channel-fields에 두지 않습니다.
import { type FieldDef, type FieldKind, SEVERITY_OPTIONS } from "./channel-fields";

// asText는 설정값을 입력란에 사용할 문자열로 변환합니다.
// config는 JSON에서 가져오며 값은 string / number / boolean / array / null일 수 있습니다.
// 여기서는 텍스트 입력란에 넣을 수 있는지만 처리하며 구체적인 직렬화는 buildConfig가 담당합니다.
function asText(v: unknown): string {
  if (typeof v === "string") return v;
  if (v === null || v === undefined) return "";
  return String(v);
}

// inputType은 필드 종류를 input의 type 속성에 매핑합니다.
function inputType(kind: FieldKind): "text" | "password" | "number" {
  if (kind === "password") return "password";
  if (kind === "number") return "number";
  return "text";
}

// ConfigField는 필드 정의에 맞는 입력 요소를 표시합니다.
//
// 마스킹된 필드 처리 시 입력란에 마스킹된 값 자체는 **표시하지 않고**, 한 줄의
// '저장됨' 안내만 표시합니다. 입력란에 텍스트가 있으면 사용자가 입력한 값이고,
// 빈 입력란은 빈 값입니다. "__masked__:…abc123"을 입력란에 넣으면 사용자는 직접
// 지워야 할 자리표시자로 오해하여 자격 증명을 잘못 지우기 쉽습니다.
export function ConfigField({
  def,
  value,
  isSecret,
  onChange,
}: {
  def: FieldDef;
  value: unknown;
  isSecret: boolean;
  onChange: (v: unknown) => void;
}) {
  const id = `n-cfg-${def.key}`;
  const raw = asText(value);
  // 백엔드에서 반환한 마스킹된 값: "__masked__:…abc123" 형식이며 끝부분은 원래 값을 식별할 수 있는 일부입니다.
  const masked = isSecret && raw.startsWith("__masked__");
  const maskedTail = masked ? (raw.split("…")[1] ?? "") : "";

  if (def.kind === "switch") {
    return (
      <div className="flex items-center gap-2 text-sm">
        <Switch checked={value === true} onCheckedChange={onChange} aria-label={def.label} />
        {def.label}
        {def.help && <span className="text-muted-foreground">({def.help})</span>}
      </div>
    );
  }

  if (def.kind === "select") {
    return (
      <div className="grid gap-2">
        <Label>{def.label}</Label>
        <Select value={raw || def.options?.[0]?.value} onValueChange={onChange}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {(def.options ?? []).map((o) => (
              <SelectItem key={o.value} value={o.value}>
                {o.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    );
  }

  // 필드 종류에 따라 입력 요소를 선택합니다. 네 종류를 구분해야 하므로 중첩 삼항식 대신 if를 이어 사용합니다.
  // 삼항식을 세 겹 중첩하면 괄호를 세어 가며 읽어야 하기 때문입니다.
  function control() {
    if (def.kind === "textarea" || def.kind === "kv") {
      return (
        <Textarea
          id={id}
          className="font-mono"
          placeholder={def.placeholder}
          value={masked ? "" : raw}
          onChange={(e) => onChange(e.target.value)}
        />
      );
    }
    if (def.kind === "list") {
      return (
        <Input
          id={id}
          value={Array.isArray(value) ? (value as string[]).join(", ") : raw}
          onChange={(e) => onChange(e.target.value)}
          placeholder={def.placeholder}
        />
      );
    }
    return (
      <Input
        id={id}
        className={def.kind === "text" ? "font-mono" : ""}
        type={inputType(def.kind)}
        placeholder={def.placeholder}
        value={masked ? "" : raw}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }

  const hint = masked ? (
    <p className="text-muted-foreground flex items-center gap-1 text-xs">
      <CheckIcon className="size-3" />
      저장됨{maskedTail ? `(끝부분 ${maskedTail})` : ""} · 새 값을 입력하면 덮어쓰고, 비우면 해당 항목을 삭제합니다
    </p>
  ) : (
    def.help && <p className="text-muted-foreground text-xs">{def.help}</p>
  );

  return (
    <div className="grid gap-2">
      <Label htmlFor={id}>{def.label}</Label>
      {control()}
      {hint}
    </div>
  );
}

// FilterSummary는 필터 조건을 한 줄로 요약하여 카드를 펼치지 않아도 채널의 알림 전송 대상을 보여 줍니다.
export function FilterSummary({ filter }: { filter: NotificationFilter }) {
  const parts: string[] = [];
  if (filter.min_severity) {
    parts.push(SEVERITY_OPTIONS.find((o) => o.value === filter.min_severity)?.label ?? filter.min_severity);
  }
  if (filter.vulnclass_include?.length) parts.push(`유형 포함 키워드 ${filter.vulnclass_include.length}개`);
  if (filter.vulnclass_exclude?.length) parts.push(`제외 키워드 ${filter.vulnclass_exclude.length}개`);
  if (filter.task_ids?.length) parts.push(`${filter.task_ids.length}개 작업`);
  if (filter.asset_ids?.length) parts.push(`${filter.asset_ids.length}개 자산`);
  if (filter.on_status_change) parts.push("상태 변경 포함");
  if (parts.length === 0) {
    return <p className="text-muted-foreground text-sm">모든 취약점</p>;
  }
  return <p className="text-muted-foreground text-sm">{parts.join(" · ")}</p>;
}
