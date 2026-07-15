import { defineStore } from "pinia";
import { useAppStore } from "./app";
import { useCcStore } from "./chatting";

// Talks directly to ai-bot-manager REST (same pattern/token as bot_channel store).
export interface QAPair {
  question: string;
  answer: string;
}

export interface KnowledgeBase {
  id: string;
  name: string;
}

export const useChatProcessStore = defineStore("chat_process", () => {
  const app = useAppStore();
  const chatting = useCcStore();

  function base(): string {
    const override = (import.meta as any).env?.VITE_AI_BOT_MANAGER_API as
      | string
      | undefined;
    const b = (override || chatting.baseUrl || "").replace(/\/$/, "");
    return `${b}/agents/api`;
  }

  function headers(): HeadersInit {
    return {
      Authorization: `Bearer ${app.conf?.token}`,
      "Content-Type": "application/json",
    };
  }

  async function ok(res: Response, action: string) {
    if (!res.ok) {
      const t = await res.text().catch(() => "");
      throw new Error(`${action} failed (${res.status}): ${t}`);
    }
  }

  // Try each chat participant against the (deployed) single-account endpoint;
  // the one that returns 200 is the core_chatting bot. Non-bot accounts return
  // 500 and are skipped. Empty account means no bot in the chat.
  async function resolveBot(
    accounts: string[],
  ): Promise<{ account: string; databases: KnowledgeBase[] }> {
    for (const account of accounts.filter(Boolean)) {
      const res = await fetch(
        `${base()}/bot_databases?account=${encodeURIComponent(account)}`,
        { headers: headers() },
      );
      if (!res.ok) continue;
      const data = await res.json();
      return { account, databases: data.databases || [] };
    }
    return { account: "", databases: [] };
  }

  async function process(
    database: string,
    account: string,
    messages: { author: string; text: string }[],
  ): Promise<QAPair[]> {
    const res = await fetch(`${base()}/process_chat`, {
      method: "POST",
      headers: headers(),
      body: JSON.stringify({ database, account, messages }),
    });
    await ok(res, "Failed to process chat");
    const data = await res.json();
    return data.pairs || [];
  }

  async function append(database: string, pairs: QAPair[]): Promise<void> {
    const res = await fetch(`${base()}/append_qa`, {
      method: "POST",
      headers: headers(),
      body: JSON.stringify({ database, pairs }),
    });
    await ok(res, "Failed to save Q&A");
  }

  return { resolveBot, process, append };
});
