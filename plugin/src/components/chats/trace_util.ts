// Small pure helpers shared by the trace viewer and its step component.

// Milliseconds between two ISO timestamps, clamped to >= 0 (traces occasionally
// carry a zero finished_at for an aborted step).
export function durMs(start: string, end: string): number {
  const ms = new Date(end).getTime() - new Date(start).getTime();
  return isFinite(ms) && ms >= 0 ? ms : 0;
}

export function fmtDur(ms: number): string {
  if (!ms) return "—";
  return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(ms < 10000 ? 2 : 1)}s`;
}

export function fmtTime(iso: string): string {
  const d = new Date(iso);
  return isNaN(d.getTime()) ? "" : d.toLocaleString();
}

// A uuid trimmed to its first block, for compact display.
export const shortId = (s: string) => (s && s.length > 8 ? s.slice(0, 8) : s || "—");
