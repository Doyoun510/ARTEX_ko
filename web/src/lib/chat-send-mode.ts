"use client";

import * as React from "react";

import { getLocalStorageValue, setLocalStorageValue } from "@/lib/local-storage.client";

// 세션 입력창의 전송/줄바꿈 키 설정. 순수 프런트엔드 환경설정: localStorage에만 저장하고 DB에 넣지 않으며 계정과 동기화하지 않으므로
// 브라우저를 바꾸면 다시 설정해야 한다. issue #39 참고 — 0.3.2에서 Ctrl+Enter 전송을 Enter 전송으로 바꿨으므로
// 여기서 기존 키 설정을 선택 항목으로 되살린다.
export type ChatSendMode = "enter" | "ctrl-enter";

export const CHAT_SEND_MODE_KEY = "artex_chat_send_mode";
export const DEFAULT_CHAT_SEND_MODE: ChatSendMode = "enter";

export const CHAT_SEND_MODE_OPTIONS: { value: ChatSendMode; label: string }[] = [
  { value: "enter", label: "Enter 전송, Shift+Enter 줄바꿈" },
  { value: "ctrl-enter", label: "Ctrl+Enter 전송, Enter 줄바꿈" },
];

function parseMode(raw: string | null): ChatSendMode {
  return raw === "ctrl-enter" || raw === "enter" ? raw : DEFAULT_CHAT_SEND_MODE;
}

// 같은 탭 안의 구독자 집합. localStorage의 storage 이벤트는 '다른' 탭에서만 발생하므로
// 이 페이지의 설정에서 바꾼 뒤에는 emit으로 같은 페이지의 입력창에 알려야 한다. 그렇지 않으면 새로고침해야 반영된다.
const listeners = new Set<() => void>();

function subscribe(listener: () => void) {
  listeners.add(listener);
  window.addEventListener("storage", listener);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", listener);
  };
}

// 문자열 리터럴을 반환하므로 Object.is가 값으로 비교해 useSyncExternalStore가 루프에 빠지지 않는다.
function getSnapshot(): ChatSendMode {
  return parseMode(getLocalStorageValue(CHAT_SEND_MODE_KEY));
}

// 서버에는 localStorage가 없으므로 먼저 기본값을 렌더링하고, hydrate 후 getSnapshot이 바로잡는다.
function getServerSnapshot(): ChatSendMode {
  return DEFAULT_CHAT_SEND_MODE;
}

export function useChatSendMode(): ChatSendMode {
  return React.useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
}

export function setChatSendMode(mode: ChatSendMode) {
  setLocalStorageValue(CHAT_SEND_MODE_KEY, mode);
  for (const listener of listeners) listener();
}

// shouldSubmitOnKey: 한 번의 키 입력이 전송이어야 하는지 판단한다.
// isComposing / keyCode 229는 중국어 등 입력기(IME)가 글자를 조합 중인 상태이므로 반드시 통과시켜야 한다. 그렇지 않으면 Enter로 단어를 고를 때 잘못 전송된다.
// enter 모드는 Shift만 제외하며 0.3.2의 동작과 완전히 같다 — 설정을 바꾸지 않은 사용자의 사용감은 그대로다.
// ctrl-enter 모드는 Ctrl과 Cmd(macOS)를 모두 받는다.
export function shouldSubmitOnKey(e: React.KeyboardEvent, mode: ChatSendMode): boolean {
  if (e.key !== "Enter") return false;
  if (e.nativeEvent.isComposing || e.nativeEvent.keyCode === 229) return false;
  if (mode === "ctrl-enter") return e.ctrlKey || e.metaKey;
  return !e.shiftKey;
}
