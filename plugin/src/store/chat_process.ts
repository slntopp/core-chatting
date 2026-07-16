import { defineStore } from "pinia";
import {
  aiBotManagerBase,
  authHeaders,
  throwIfNotOk,
} from "./aiBotManager";

export interface QAPair {
  question: string;
  answer: string;
}

export interface KnowledgeBase {
  id: string;
  name: string;
}

export const useChatProcessStore = defineStore("chat_process", () => {
  // Try each chat participant against the (deployed) single-account endpoint;
  // the one that returns 200 is the core_chatting bot. Non-bot accounts return
  // 500 and are skipped. Empty account means no bot in the chat.
  async function resolveBot(
    accounts: string[],
  ): Promise<{ account: string; databases: KnowledgeBase[] }> {
    for (const account of accounts.filter(Boolean)) {
      const res = await fetch(
        `${aiBotManagerBase()}/bot_databases?account=${encodeURIComponent(account)}`,
        { headers: authHeaders() },
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
    const res = await fetch(`${aiBotManagerBase()}/process_chat`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ database, account, messages }),
    });
    await throwIfNotOk(res, "Failed to process chat");
    const data = await res.json();
    return data.pairs || [];
  }

  async function append(database: string, pairs: QAPair[]): Promise<void> {
    const res = await fetch(`${aiBotManagerBase()}/append_qa`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ database, pairs }),
    });
    await throwIfNotOk(res, "Failed to save Q&A");
  }

  return { resolveBot, process, append };
});
