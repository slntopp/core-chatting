// The Otus Learn window. Otus is guarded by a shared secret the browser must
// never hold, so these calls go to core-chatting's own /otus/learn* proxy
// (pkg/otus), which adds the bearer and checks the caller is an admin of the
// chat. Every other why value comes straight from Otus.
import { nocloudBase, authHeaders } from "./aiBotManager";

export interface LearnTake {
  // Absent on a take the operator added: Otus assigns the id.
  id?: string;
  question: string;
  answer: string;
}

export interface LearnAnswer {
  ok: boolean;
  why: string;
  chat?: string;
  takes?: LearnTake[];
  summary?: string;
  prUrl?: string;
}

// Propose and Save run a model over the whole ticket. Nothing here times out
// on its own; the proxy caps the wait at 180s.
async function call(path: string, init: RequestInit): Promise<LearnAnswer> {
  const res = await fetch(`${nocloudBase()}${path}`, {
    headers: authHeaders(),
    ...init,
  });
  const data = (await res.json().catch(() => ({}))) as Partial<LearnAnswer>;
  if (!res.ok) {
    throw new Error(data.why || `Learn failed (${res.status})`);
  }
  return { ok: true, why: "", ...data } as LearnAnswer;
}

export function proposeTakes(chat: string): Promise<LearnAnswer> {
  return call("/otus/learn", { method: "POST", body: JSON.stringify({ chat }) });
}

export function draftTakes(chat: string): Promise<LearnAnswer> {
  return call(`/otus/learn?chat=${encodeURIComponent(chat)}`, { method: "GET" });
}

export function saveTakes(chat: string, takes: LearnTake[]): Promise<LearnAnswer> {
  return call("/otus/learn/save", {
    method: "POST",
    body: JSON.stringify({ chat, takes }),
  });
}
