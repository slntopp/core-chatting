<template>
  <details class="step" :open="open" @toggle="onToggle">
    <summary class="step__sum">
      <span class="node" :style="{ background: tone }">{{ index + 1 }}</span>
      <span class="step__kind" :style="{ color: tone }">{{ kindLabel }}</span>
      <span v-if="step.name" class="step__name">{{ step.name }}</span>
      <span v-if="step.model" class="tag">{{ step.model }}</span>
      <span v-if="(step.tool_calls || []).length" class="tag tag--tool">
        {{ step.tool_calls!.length }}⚙
      </span>
      <span class="spacer" />
      <span class="step__bar" aria-hidden="true">
        <span class="step__barfill" :style="barStyle" />
      </span>
      <span class="step__dur">{{ fmtDur(durMs(step.started_at, step.finished_at)) }}</span>
      <svg class="chev chev--sm" viewBox="0 0 16 16" width="11" height="11" aria-hidden="true">
        <path d="M6 4l4 4-4 4" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
    </summary>

    <div class="step__detail">
      <trace-block v-if="step.prompt" label="Prompt" :content="step.prompt" />

      <!-- tools OFFERED this step (vs the timeline below = actually called) -->
      <trace-disclosure v-if="(step.tools || []).length">
        <template #summary>
          <span class="gsum__ico">⚒</span>
          <span class="gsum__lbl">Available tools</span>
          <span class="gsum__count">{{ step.tools!.length }}</span>
          <span class="spacer" />
        </template>
        <trace-disclosure v-for="(t, ti) in step.tools" :key="'t' + ti" dense>
          <template #summary>
            <span class="tdef__name">{{ t.name }}</span>
            <span v-if="t.description" class="tdef__desc">{{ t.description }}</span>
            <span class="spacer" />
          </template>
          <trace-block v-if="t.description" label="Description" tone="muted" :content="t.description" />
          <trace-block v-if="t.schema" label="Parameters" :content="prettyJson(t.schema)" />
        </trace-disclosure>
      </trace-disclosure>

      <!-- knowledge bases in scope for this step's RAG, per-kind -->
      <trace-disclosure v-if="(step.databases || []).length">
        <template #summary>
          <span class="gsum__ico">▤</span>
          <span class="gsum__lbl">Knowledge in scope</span>
          <span class="gsum__count">{{ step.databases!.length }}</span>
          <span class="spacer" />
        </template>
        <div v-for="(db, di) in step.databases" :key="'db' + di" class="kb">
          <div class="kb__hd">
            <span class="kb__name">{{ db.name || shortId(db.id) }}</span>
            <span class="kb__counts">{{ dbCounts(db) }}</span>
          </div>
          <div v-if="(db.files || []).length" class="kb__files">
            <span
              v-for="(f, fi) in db.files"
              :key="'f' + fi"
              class="kb__file"
              :class="{ 'is-off': !f.enabled || f.status !== 'READY' }"
              :title="`${f.status || ''}${f.enabled ? '' : ' · disabled'}`"
            >
              {{ f.name || "(file)" }}
              <span class="kb__fstat">{{ f.enabled ? f.status || "" : "disabled" }}</span>
            </span>
          </div>
        </div>
      </trace-disclosure>

      <template v-if="step.rag_query || (step.rag || []).length">
        <trace-block v-if="step.rag_query" label="RAG query" tone="accent" :content="step.rag_query" />
        <trace-block
          v-for="(c, ci) in step.rag || []"
          :key="'rag' + ci"
          :label="`RAG chunk ${ci + 1}`"
          :meta="`score ${c.score.toFixed(3)}${c.origin ? ' · ' + shortId(c.origin) : ''}`"
          :content="c.content || ''"
        />
      </template>

      <!-- everything the model received as input for this step -->
      <trace-disclosure v-if="contextMessages.length">
        <template #summary>
          <span class="gsum__ico">≡</span>
          <span class="gsum__lbl">Context the model saw</span>
          <span class="gsum__count">{{ contextMessages.length }} msg</span>
          <span class="spacer" />
        </template>
        <trace-block
          v-for="(m, mi) in contextMessages"
          :key="'ctx' + mi"
          :label="roleLabel(m.role)"
          :tone="roleTone(m.role)"
          :content="m.content"
        />
      </trace-disclosure>

      <!-- chronological timeline: tool calls (in call order) and the answers -->
      <template v-for="ev in events" :key="ev.key">
        <trace-disclosure v-if="ev.type === 'tool'" :accent="ev.tool!.error ? 'error' : 'tool'">
          <template #summary>
            <span class="call__ico">⚙</span>
            <span class="call__name">{{ ev.tool!.name }}</span>
            <span v-if="argPreview(ev.tool!.args)" class="call__args">{{ argPreview(ev.tool!.args) }}</span>
            <span class="spacer" />
            <span class="call__status" :class="ev.tool!.error ? 'is-err' : 'is-ok'">
              {{ ev.tool!.error ? "error" : "ok" }}
            </span>
          </template>
          <trace-block v-if="ev.tool!.args" label="Arguments" :content="ev.tool!.args" />
          <trace-block
            :label="ev.tool!.error ? 'Error' : 'Output'"
            :tone="ev.tool!.error ? 'error' : 'code'"
            :content="ev.tool!.error || ev.tool!.output || ''"
          />
        </trace-disclosure>

        <trace-block v-else-if="ev.type === 'reasoning'" label="Reasoning" tone="muted" :content="ev.text || ''" />
        <trace-block v-else-if="ev.type === 'answer'" label="Answer" :content="ev.text || ''" />
      </template>

      <trace-block v-if="hasVars" label="Vars" :content="varsText" />
    </div>
  </details>

  <!-- routing to the next step, in plain language, hung off the spine -->
  <div v-if="step.routing" class="route">
    <span class="route__hook">↳</span>
    <template v-if="step.routing.stop && !step.routing.goto">
      <span class="route__end">conversation ends here</span>
    </template>
    <template v-else>
      <span class="route__lead">next</span>
      <span class="route__arrow">→</span>
      <code class="route__goto">{{ step.routing.goto || "(next step)" }}</code>
      <span v-if="step.routing.value" class="route__why">
        because it decided “{{ step.routing.value }}”
      </span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { TraceStep, TraceToolCall, TraceDatabase } from "../../store/traces";
import { durMs, fmtDur, shortId } from "./trace_util";
import TraceBlock from "./trace_block.vue";
import TraceDisclosure from "./trace_disclosure.vue";

const props = defineProps<{
  step: TraceStep;
  index: number;
  turnMs: number; // whole-turn duration, for the waterfall bar proportion
  open: boolean;
}>();
const emit = defineEmits<{ (e: "toggle", open: boolean): void }>();

function onToggle(e: Event) {
  emit("toggle", (e.target as HTMLDetailsElement).open);
}

// One timeline event: a tool call (matched to the round that requested it) or a
// piece of the model's answer. Reconstructed in true chronological order.
type StepEvent = { key: string } & (
  | { type: "tool"; tool: TraceToolCall }
  | { type: "answer"; text: string }
  | { type: "reasoning"; text: string }
);

// Plain-language label + colour, by what the step actually does:
//   decision (classifier) — picks which branch runs next   (purple)
//   reply    — the message the customer actually receives   (green)
//   action   — internal work: gather data, write admin note (blue)
const kindLabel = computed(() => {
  if (props.step.kind === "classifier") return "decision";
  if (props.step.reply) return "reply";
  return "action";
});
const tone = computed(() => {
  if (props.step.kind === "classifier") return "#7c5cff";
  if (props.step.reply) return "#18a058";
  return "#2f6fed";
});

const barStyle = computed(() => {
  const total = props.turnMs || 1;
  const w = Math.max(3, Math.round((durMs(props.step.started_at, props.step.finished_at) / total) * 100));
  return { width: `${Math.min(w, 100)}%`, background: tone.value };
});

// Rounds and tool_calls are stored flat in call order; each round names the
// tools it requested, so we consume that many from the flat tool list per round,
// interleaving with that round's answer. Leftover tools are appended so nothing
// is ever hidden.
const events = computed<StepEvent[]>(() => {
  const out: StepEvent[] = [];
  const tools = props.step.tool_calls || [];
  let ti = 0;
  let n = 0;
  for (const r of props.step.rounds || []) {
    if (r.reasoning) out.push({ key: `e${n++}`, type: "reasoning", text: r.reasoning });
    const req = (r.tools_requested || []).length;
    for (let k = 0; k < req && ti < tools.length; k++) {
      out.push({ key: `e${n++}`, type: "tool", tool: tools[ti++] });
    }
    if (r.response) out.push({ key: `e${n++}`, type: "answer", text: r.response });
  }
  while (ti < tools.length) out.push({ key: `e${n++}`, type: "tool", tool: tools[ti++] });
  return out;
});

// The full input the model saw, minus the step's own prompt (already shown as
// "Prompt" above) — what's left is the genuinely extra context: earlier steps'
// forwarded output, chat history, injected knowledge.
const contextMessages = computed(() => {
  const req = props.step.request || [];
  let promptDropped = false;
  return req.filter((m) => {
    if (!promptDropped && props.step.prompt && m.role === "system" && m.content === props.step.prompt) {
      promptDropped = true;
      return false;
    }
    return (m.content || "").trim().length > 0;
  });
});

const hasVars = computed(
  () =>
    Object.keys(props.step.vars_in || {}).length > 0 ||
    Object.keys(props.step.vars_out || {}).length > 0,
);
const varsText = computed(() => {
  const lines: string[] = [];
  const inV = props.step.vars_in || {};
  const outV = props.step.vars_out || {};
  for (const k of Object.keys(inV)) lines.push(`${k}: ${inV[k]}`);
  for (const k of Object.keys(outV)) {
    if (inV[k] !== outV[k]) lines.push(`${k} → ${outV[k]}`);
  }
  return lines.join("\n");
});

function roleLabel(role: string): string {
  if (role === "user") return "customer";
  if (role === "assistant") return "bot";
  if (role === "system") return "system";
  return role || "message";
}
function roleTone(role: string): "code" | "muted" | "accent" {
  if (role === "system") return "muted";
  if (role === "user") return "accent";
  return "code";
}

function prettyJson(s: string): string {
  try {
    return JSON.stringify(JSON.parse(s), null, 2);
  } catch {
    return s;
  }
}

function dbCounts(db: TraceDatabase): string {
  const parts: string[] = [];
  if (db.simple) parts.push(`${db.simple} docs`);
  if ((db.files || []).length) parts.push(`${db.files!.length} files`);
  if (db.qa_pairs) parts.push(`${db.qa_pairs} Q&A`);
  if (db.urls) parts.push(`${db.urls} urls`);
  return parts.length ? parts.join(" · ") : "empty";
}

// Compact one-line preview of a tool call's JSON args for the collapsed row,
// e.g. {"action":"balance"} -> action: balance. Falls back to a trimmed string.
function argPreview(args?: string): string {
  if (!args) return "";
  let obj: unknown;
  try {
    obj = JSON.parse(args);
  } catch {
    return args.length > 60 ? args.slice(0, 60) + "…" : args;
  }
  if (obj && typeof obj === "object" && !Array.isArray(obj)) {
    const parts = Object.entries(obj as Record<string, unknown>).map(
      ([k, v]) => `${k}: ${typeof v === "string" ? v : JSON.stringify(v)}`,
    );
    const s = parts.join(", ");
    return s.length > 70 ? s.slice(0, 70) + "…" : s;
  }
  return "";
}
</script>

<style scoped>
.spacer {
  flex: 1;
}

/* step row on the timeline spine */
.step {
  position: relative;
  margin: 0 0 2px;
}
.node {
  position: absolute;
  left: -28px;
  top: 6px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-family: ui-monospace, monospace;
  font-size: 10px;
  font-weight: 700;
  box-shadow: 0 0 0 3px var(--node-ring, rgba(255, 255, 255, 0.9));
}
@media (prefers-color-scheme: dark) {
  .node {
    box-shadow: 0 0 0 3px rgba(24, 24, 28, 1);
  }
}
.step__sum {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 6px 4px;
  cursor: pointer;
  list-style: none;
  border-radius: 6px;
}
.step__sum::-webkit-details-marker {
  display: none;
}
.step__sum:hover {
  background: rgba(128, 128, 128, 0.07);
}
.step__kind {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}
.step__name {
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}
.step__bar {
  width: 54px;
  height: 4px;
  border-radius: 3px;
  background: rgba(128, 128, 128, 0.18);
  overflow: hidden;
  flex: none;
}
.step__barfill {
  display: block;
  height: 100%;
  border-radius: 3px;
  opacity: 0.75;
}
.step__dur {
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  opacity: 0.6;
  min-width: 44px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.chev {
  opacity: 0.4;
  transition: transform 0.18s ease;
  flex: none;
}
.chev--sm {
  opacity: 0.35;
}
.step[open] > .step__sum .chev {
  transform: rotate(90deg);
}
.step__detail {
  display: flex;
  flex-direction: column;
  gap: 7px;
  padding: 4px 0 10px 2px;
}

/* shared group-header bits (slotted into trace-disclosure summaries) */
.gsum__ico {
  opacity: 0.55;
  font-size: 13px;
  flex: none;
}
.gsum__lbl {
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  opacity: 0.7;
}
.gsum__count {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  opacity: 0.45;
}

/* one available-tool row header */
.tdef__name {
  font-family: ui-monospace, monospace;
  font-size: 11.5px;
  font-weight: 700;
  flex: none;
}
.tdef__desc {
  font-size: 11px;
  opacity: 0.6;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

/* one knowledge base card */
.kb {
  border: 1px solid rgba(128, 128, 128, 0.2);
  border-radius: 6px;
  padding: 6px 9px;
  background: rgba(128, 128, 128, 0.03);
}
.kb__hd {
  display: flex;
  align-items: baseline;
  gap: 9px;
}
.kb__name {
  font-weight: 600;
  font-size: 12px;
}
.kb__counts {
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  opacity: 0.55;
}
.kb__files {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-top: 6px;
}
.kb__file {
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  padding: 1px 7px;
  border-radius: 5px;
  background: rgba(24, 160, 88, 0.13);
  color: #18a058;
}
.kb__file.is-off {
  background: rgba(128, 128, 128, 0.14);
  color: inherit;
  opacity: 0.55;
}
.kb__fstat {
  opacity: 0.6;
  font-size: 9.5px;
  margin-left: 3px;
}

/* one tool-call row header */
.call__ico {
  color: #e0912f;
  font-size: 12px;
  flex: none;
}
.call__name {
  font-family: ui-monospace, monospace;
  font-size: 12px;
  font-weight: 700;
  flex: none;
}
.call__args {
  font-family: ui-monospace, monospace;
  font-size: 11px;
  opacity: 0.6;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}
.call__status {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 1px 6px;
  border-radius: 20px;
  flex: none;
}
.call__status.is-ok {
  color: #18a058;
  background: rgba(24, 160, 88, 0.14);
}
.call__status.is-err {
  color: #d03050;
  background: rgba(208, 48, 80, 0.14);
}

/* tags */
.tag {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 5px;
  background: rgba(128, 128, 128, 0.15);
  white-space: nowrap;
  flex: none;
}
.tag--tool {
  color: #e0912f;
  background: rgba(224, 145, 47, 0.14);
}

/* routing connector */
.route {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 1px 0 8px -2px;
  padding: 3px 0;
  font-family: ui-monospace, monospace;
  font-size: 11px;
}
.route__hook {
  opacity: 0.4;
  margin-left: -20px;
  font-size: 13px;
}
.route__lead {
  opacity: 0.5;
  text-transform: uppercase;
  font-size: 9.5px;
  letter-spacing: 0.06em;
}
.route__arrow {
  opacity: 0.45;
}
.route code {
  padding: 1px 6px;
  border-radius: 4px;
  background: rgba(128, 128, 128, 0.14);
  font-size: 11px;
}
.route__goto {
  color: #2f6fed;
  font-weight: 600;
}
.route__why {
  opacity: 0.55;
  font-style: italic;
}
.route__end {
  color: #d03050;
  opacity: 0.8;
  font-weight: 600;
}
</style>
