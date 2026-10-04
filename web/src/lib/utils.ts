import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// 슬라이드 패널/대화상자(Sheet/Dialog)의 onInteractOutside 닫기 판정을 돕습니다.
//
// 배경: 슬라이드 패널 안의 Radix 오버레이(Select 드롭다운, DropdownMenu, Popover 등)는 portal을 통해 슬라이드 패널
// 밖으로 이동합니다. 열린 오버레이를 닫으려고 배경 오버레이/슬라이드 패널 바깥을 클릭하면 이 pointerdown은 Select와 Sheet의 두
// DismissableLayer가 동시에 처리합니다. Select가 먼저 닫히며 discrete 이벤트이므로 React가 동기적으로 flush합니다. 따라서
// Sheet가 처리할 차례에는 오버레이의 data-state가 이미 closed로 바뀌어 있습니다. 즉, "현재" 오버레이가 열려 있는지 감지하는 것은
// 본질적으로 신뢰할 수 없습니다(실측으로 확인됨).
//
// 올바른 방법: Radix의 pointerdown 감지는 이벤트 버블링 단계에서 이루어집니다. 그보다 앞선 capture 단계에서 먼저
// "지금 열린 오버레이가 있는지" 기록하고 onInteractOutside가 이 기록값을 읽어 닫기를 허용할지 결정합니다.
function isRadixOverlayOpenNow(): boolean {
  if (typeof document === "undefined") return false;
  return !!document.querySelector(
    [
      "[data-slot='select-trigger'][data-state='open']",
      "[data-slot='select-content'][data-state='open']",
      "[role='listbox'][data-state='open']",
      "[data-radix-popper-content-wrapper]",
      "[aria-expanded='true'][data-state='open']",
    ].join(","),
  );
}

let overlayOpenAtLastPointerDown = false;
if (typeof document !== "undefined") {
  document.addEventListener(
    "pointerdown",
    () => {
      overlayOpenAtLastPointerDown = isRadixOverlayOpenNow();
    },
    true, // capture: Radix의 이벤트 버블링 단계 pointerdown 처리보다 먼저 기록합니다.
  );
}

// radixOverlayWasOpenAtPointerDown은 "마지막 pointerdown 발생 시 Radix 오버레이가
// 열려 있었는지" 반환합니다. 슬라이드 패널/대화상자는 이를 기준으로 열린 오버레이가 있을 때 배경 오버레이를 클릭하면 → 오버레이만 닫고 자신은 닫지 않습니다.
export function radixOverlayWasOpenAtPointerDown(): boolean {
  return overlayOpenAtLastPointerDown;
}

// copyText는 텍스트를 클립보드에 기록하고 성공 여부를 반환합니다.
// 배경: navigator.clipboard는 보안 컨텍스트(HTTPS / localhost)에서만 사용할 수 있습니다. IP + HTTP로
// 접근하면 undefined이므로 이때 execCommand("copy")를 사용해 대체 방식으로 처리합니다.
export async function copyText(text: string): Promise<boolean> {
  if (navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      // 대체 방식으로 계속 진행합니다.
    }
  }
  try {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.style.position = "fixed";
    textarea.style.left = "-9999px";
    textarea.style.top = "0";
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    const ok = document.execCommand("copy");
    document.body.removeChild(textarea);
    return ok;
  } catch {
    return false;
  }
}

export const getInitials = (str: string): string => {
  if (typeof str !== "string" || !str.trim()) return "?";

  return (
    str
      .trim()
      .split(/\s+/)
      .filter(Boolean)
      .map((word) => word[0])
      .join("")
      .toUpperCase() || "?"
  );
};

export function formatCurrency(
  amount: number,
  opts?: {
    currency?: string;
    locale?: string;
    minimumFractionDigits?: number;
    maximumFractionDigits?: number;
    noDecimals?: boolean;
  },
) {
  const { currency = "USD", locale = "en-US", minimumFractionDigits, maximumFractionDigits, noDecimals } = opts ?? {};

  const formatOptions: Intl.NumberFormatOptions = {
    style: "currency",
    currency,
    minimumFractionDigits: noDecimals ? 0 : minimumFractionDigits,
    maximumFractionDigits: noDecimals ? 0 : maximumFractionDigits,
  };

  return new Intl.NumberFormat(locale, formatOptions).format(amount);
}
