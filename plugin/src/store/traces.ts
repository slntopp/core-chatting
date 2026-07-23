import { defineStore } from "pinia";
import { aiBotManagerBase, authHeaders, throwIfNotOk } from "./aiBotManager";

// Mirrors ai-bot-manager's core.Trace JSON (pkg/core/trace.go). One Trace is one
// AI turn; a chat's traces are all its turns, oldest first.
export interface TraceMessage {
  role: string;
  content: string;
}
export interface TraceLLMRound {
  model?: string;
  response?: string;
  reasoning?: string;
  tools_requested?: string[];
}
export interface TraceToolCall {
  name: string;
  args?: string;
  output?: string;
  error?: string;
}
export interface TraceChunk {
  score: number;
  origin?: string;
  content?: string;
}
export interface TraceRouting {
  when_var?: string;
  value?: string;
  equals?: string;
  goto?: string;
  matched: boolean;
  stop: boolean;
}
export interface TraceStep {
  name?: string;
  kind: string; // "single" | "classifier" | "content"
  model?: string;
  prompt?: string;
  vars_in?: Record<string, string>;
  vars_out?: Record<string, string>;
  rag_query?: string;
  rag?: TraceChunk[];
  request?: TraceMessage[];
  rounds?: TraceLLMRound[];
  tool_calls?: TraceToolCall[];
  output?: string;
  routing?: TraceRouting;
  reply: boolean;
  started_at: string;
  finished_at: string;
}
export interface Trace {
  id: string;
  chat_id: string;
  external: string;
  channel: string;
  bot_id: string;
  account: string;
  mode: string; // "flow" | "single"
  input: string;
  steps: TraceStep[];
  error?: string;
  started_at: string;
  finished_at: string;
}

export const useTracesStore = defineStore("traces", () => {
  // Fetch every turn trace for a chat. `external` is the core-chatting chat
  // uuid, which is also ai-bot-manager's external chat id for the core_chatting
  // channel — that shared identity is the whole correlation, no per-message
  // plumbing needed. Admin-only on the server side.
  async function getChatTraces(external: string): Promise<Trace[]> {
    const res = await fetch(
      `${aiBotManagerBase()}/get_chat_traces/${encodeURIComponent(external)}`,
      { headers: authHeaders() },
    );
    await throwIfNotOk(res, "Failed to load traces");
    const data = await res.json();
    return (data.traces as Trace[]) || [];
  }

  return { getChatTraces };
});
