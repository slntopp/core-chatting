import { defineStore } from "pinia";
import { useAppStore } from "./app";

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

  function base(): string {
    const override = (import.meta as any).env?.VITE_AI_BOT_MANAGER_API as
      | string
      | undefined;
    const b = (override || app.conf?.api || "").replace(/\/$/, "");
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

  // Databases of the bot whose core_chatting channel account is a member of this chat.
  async function basesForChat(adminUuids: string[]): Promise<KnowledgeBase[]> {
    const res = await fetch(`${base()}/get_bots`, { headers: headers() });
    await ok(res, "Failed to fetch bots");
    const data = await res.json();
    for (const bot of data.bots || []) {
      const linked = (bot.channels || []).some(
        (c: any) =>
          c.type === "core_chatting" && adminUuids.includes(c.external_id),
      );
      if (!linked) continue;
      return (bot.databases || []).map((d: any) => ({
        id: d.id,
        name: d.name,
      }));
    }
    return [];
  }

  async function process(
    database: string,
    messages: { author: string; text: string }[],
  ): Promise<QAPair[]> {
    const res = await fetch(`${base()}/process_chat`, {
      method: "POST",
      headers: headers(),
      body: JSON.stringify({ database, messages }),
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

  return { basesForChat, process, append };
});
