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

// One proposal from the model: match === -1 -> add new, >= 0 -> replace the
// existing entry at that index. `old` is that one existing record (for a
// replace) so the UI can show "was" — the whole base is never sent here.
export interface QAProposal {
  match: number;
  question: string;
  answer: string;
  old?: QAPair;
}

export interface ProcessResult {
  items: QAProposal[];
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
  ): Promise<ProcessResult> {
    const res = await fetch(`${aiBotManagerBase()}/process_chat`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ database, account, messages }),
    });
    await throwIfNotOk(res, "Failed to process chat");
    const data = await res.json();
    return { items: data.items || [] };
  }

  // Apply the approved proposals. The backend reads the base and applies each
  // item in place (replace at match, or append) — the full base never travels.
  async function save(
    database: string,
    items: QAProposal[],
  ): Promise<void> {
    const res = await fetch(`${aiBotManagerBase()}/append_qa`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ database, items }),
    });
    await throwIfNotOk(res, "Failed to save Q&A");
  }

  return { resolveBot, process, save };
});
