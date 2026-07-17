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
// existing entry at that index.
export interface QAProposal {
  match: number;
  question: string;
  answer: string;
}

export interface QAKnowledge {
  id: string;
  database?: string;
  records: QAPair[];
}

export interface ProcessResult {
  qa_knowledge: QAKnowledge | null;
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
    return { qa_knowledge: data.qa_knowledge || null, items: data.items || [] };
  }

  // Persist the final record set. Reuses existing endpoints: update_knowledge
  // replaces the whole QA base (when it already exists), append_qa creates it
  // for a base that had no QA knowledge yet.
  async function save(
    database: string,
    qaKnowledge: QAKnowledge | null,
    records: QAPair[],
  ): Promise<void> {
    if (qaKnowledge?.id) {
      const res = await fetch(`${aiBotManagerBase()}/update_knowledge`, {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify({
          database,
          type: "question_answer",
          data: { qa_knowledge: { ...qaKnowledge, records } },
        }),
      });
      await throwIfNotOk(res, "Failed to save Q&A");
      return;
    }
    const res = await fetch(`${aiBotManagerBase()}/append_qa`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ database, pairs: records }),
    });
    await throwIfNotOk(res, "Failed to save Q&A");
  }

  return { resolveBot, process, save };
});
