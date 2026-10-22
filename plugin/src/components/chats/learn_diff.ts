// A unified diff, just enough of it to draw the GitHub files view.
// Hunk headers, additions, deletions and context. Anything else is skipped.

export type DiffKind = "add" | "del" | "ctx" | "hunk";

export type DiffLine = {
  kind: DiffKind;
  oldNo: string;
  newNo: string;
  sign: string;
  text: string;
};

export type DiffFile = {
  path: string;
  lines: DiffLine[];
};

export function parseUnifiedDiff(patch: string): DiffFile[] {
  const files: DiffFile[] = [];
  let current: DiffFile | undefined;
  let oldNo = 0;
  let newNo = 0;

  for (const raw of patch.split("\n")) {
    if (raw.startsWith("diff --git ")) {
      const match = / b\/(.+)$/.exec(raw);
      current = { path: match?.[1] ?? "file", lines: [] };
      files.push(current);
      oldNo = 0;
      newNo = 0;
      continue;
    }
    if (current === undefined) continue;
    if (
      raw.startsWith("index ") ||
      raw.startsWith("--- ") ||
      raw.startsWith("+++ ") ||
      raw.startsWith("\\")
    ) {
      continue;
    }
    if (raw.startsWith("@@")) {
      const hunk = /@@ -(\d+)(?:,\d+)? \+(\d+)/.exec(raw);
      oldNo = hunk ? Number(hunk[1]) : 0;
      newNo = hunk ? Number(hunk[2]) : 0;
      current.lines.push({ kind: "hunk", oldNo: "", newNo: "", sign: "", text: raw });
      continue;
    }
    if (raw.startsWith("+")) {
      current.lines.push({
        kind: "add",
        oldNo: "",
        newNo: String(newNo),
        sign: "+",
        text: raw.slice(1),
      });
      newNo += 1;
      continue;
    }
    if (raw.startsWith("-")) {
      current.lines.push({
        kind: "del",
        oldNo: String(oldNo),
        newNo: "",
        sign: "-",
        text: raw.slice(1),
      });
      oldNo += 1;
      continue;
    }
    if (raw.startsWith(" ") || raw === "") {
      const text = raw.startsWith(" ") ? raw.slice(1) : "";
      current.lines.push({
        kind: "ctx",
        oldNo: String(oldNo),
        newNo: String(newNo),
        sign: " ",
        text,
      });
      oldNo += 1;
      newNo += 1;
    }
  }

  return files.filter((file) => file.lines.length > 0);
}
