// Shared helpers for the stores that call NoCloud's REST gateway and
// ai-bot-manager's REST API directly (bot_channel, chat_process). Every other
// store here only talks to core-chatting's own Connect-Web API.
//
// All three services (core-chatting, NoCloud gateway, ai-bot-manager) sit behind
// the same origin and validate JWTs signed with the same shared signing key, so
// we reuse core-chatting's base URL and the iframe bearer token for all of them.
import { useAppStore } from "./app";
import { useCcStore } from "./chatting";

export function nocloudBase(): string {
  return (useCcStore().baseUrl || "").replace(/\/$/, "");
}

export function aiBotManagerBase(): string {
  const override = (import.meta as any).env?.VITE_AI_BOT_MANAGER_API as
    | string
    | undefined;
  const base = (override || useCcStore().baseUrl || "").replace(/\/$/, "");
  return `${base}/agents/api`;
}

export function authHeaders(): HeadersInit {
  return {
    Authorization: `Bearer ${useAppStore().conf?.token}`,
    "Content-Type": "application/json",
  };
}

export async function throwIfNotOk(res: Response, action: string) {
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`${action} failed (${res.status}): ${text}`);
  }
}
