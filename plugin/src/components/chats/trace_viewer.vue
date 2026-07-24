<template>
  <div class="tv">
    <!-- states -->
    <div v-if="loading" class="tv__state"><n-spin size="large" /></div>

    <n-alert v-else-if="error" type="error" title="Couldn't load traces">
      {{ error }}
      <template #action>
        <button class="lnk" @click="reload">Try again</button>
      </template>
    </n-alert>

    <n-empty
      v-else-if="!traces.length"
      class="tv__state"
      description="No traces recorded for this chat yet"
    >
      <template #extra>
        <button class="lnk" @click="reload">Refresh</button>
      </template>
    </n-empty>

    <template v-else>
      <!-- toolbar -->
      <div class="bar">
        <span class="bar__count">
          <b>{{ traces.length }}</b> turn{{ traces.length === 1 ? "" : "s" }}
          <span class="bar__dim">· {{ totalSteps }} steps · {{ totalTools }} tool calls</span>
        </span>
        <span class="bar__spacer" />
        <button class="lnk" @click="expandAll">expand all</button>
        <button class="lnk" @click="collapseAll">collapse all</button>
        <button class="lnk" @click="reload">refresh</button>
      </div>

      <!-- turns -->
      <details
        v-for="(t, ti) in traces"
        :key="t.id"
        class="turn"
        :open="!!openTurns[t.id]"
        @toggle="onTurnToggle(t.id, $event)"
      >
        <summary class="turn__sum">
          <span class="turn__idx">#{{ pad(ti + 1) }}</span>
          <span class="chip" :class="`chip--${t.mode}`">{{ t.mode }}</span>
          <span class="turn__input" :title="t.input">{{ t.input || "(no input)" }}</span>
          <span class="turn__metrics">
            <span v-if="turnTools(t)" class="metric metric--tool">{{ turnTools(t) }}⚙</span>
            <span class="metric">{{ t.steps.length }}▸</span>
            <span class="metric metric--dur">{{ fmtDur(dur(t.started_at, t.finished_at)) }}</span>
          </span>
          <svg class="chev" viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
            <path d="M6 4l4 4-4 4" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </summary>

        <div class="turn__body">
          <div class="turn__info">
            <span>{{ fmtTime(t.started_at) }}</span>
            <span class="dot">•</span>
            <span>bot {{ short(t.bot_id) }}</span>
            <span class="dot">•</span>
            <span>{{ t.channel }}</span>
          </div>

          <n-alert v-if="t.error" type="warning" :bordered="false" class="turn__err">
            {{ t.error }}
          </n-alert>

          <!-- timeline spine -->
          <div class="steps">
            <template v-for="(s, si) in t.steps" :key="si">
              <details
                class="step"
                :open="isStepOpen(t.id, si)"
                @toggle="onStepToggle(t.id, si, $event)"
              >
                <summary class="step__sum">
                  <span class="node" :style="{ background: stepTone(s) }">{{ si + 1 }}</span>
                  <span class="step__kind" :style="{ color: stepTone(s) }">{{ kindLabel(s) }}</span>
                  <span v-if="s.name" class="step__name">{{ s.name }}</span>
                  <span v-if="s.model" class="tag">{{ s.model }}</span>
                  <span v-if="(s.tool_calls || []).length" class="tag tag--tool">
                    {{ s.tool_calls!.length }}⚙
                  </span>
                  <span class="step__spacer" />
                  <span class="step__bar" aria-hidden="true">
                    <span class="step__barfill" :style="barStyle(s, t)" />
                  </span>
                  <span class="step__dur">{{ fmtDur(dur(s.started_at, s.finished_at)) }}</span>
                  <svg class="chev chev--sm" viewBox="0 0 16 16" width="11" height="11" aria-hidden="true">
                    <path d="M6 4l4 4-4 4" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </summary>

                <div class="step__detail">
                  <trace-block v-if="s.prompt" label="Prompt" :content="s.prompt" />

                  <template v-if="s.rag_query || (s.rag || []).length">
                    <trace-block
                      v-if="s.rag_query"
                      label="RAG query"
                      tone="accent"
                      :content="s.rag_query"
                    />
                    <trace-block
                      v-for="(c, ci) in s.rag || []"
                      :key="'rag' + ci"
                      :label="`RAG chunk ${ci + 1}`"
                      :meta="`score ${c.score.toFixed(3)}${c.origin ? ' · ' + short(c.origin) : ''}`"
                      :content="c.content || ''"
                    />
                  </template>

                  <!-- one sequential timeline: tool calls (in call order) and the
                       model's answers, interleaved as they actually happened -->
                  <template v-for="ev in stepEvents(s)" :key="ev.key">
                    <details v-if="ev.type === 'tool'" class="call" :class="{ 'call--err': !!ev.tool!.error }">
                      <summary class="call__sum">
                        <span class="call__ico">⚙</span>
                        <span class="call__name">{{ ev.tool!.name }}</span>
                        <span v-if="argPreview(ev.tool!.args)" class="call__args">{{ argPreview(ev.tool!.args) }}</span>
                        <span class="call__spacer" />
                        <span class="call__status" :class="ev.tool!.error ? 'is-err' : 'is-ok'">
                          {{ ev.tool!.error ? "error" : "ok" }}
                        </span>
                        <svg class="chev chev--sm" viewBox="0 0 16 16" width="11" height="11" aria-hidden="true">
                          <path d="M6 4l4 4-4 4" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                        </svg>
                      </summary>
                      <div class="call__body">
                        <trace-block v-if="ev.tool!.args" label="Arguments" :content="ev.tool!.args" />
                        <trace-block
                          :label="ev.tool!.error ? 'Error' : 'Output'"
                          :tone="ev.tool!.error ? 'error' : 'code'"
                          :content="ev.tool!.error || ev.tool!.output || ''"
                        />
                      </div>
                    </details>

                    <trace-block
                      v-else-if="ev.type === 'reasoning'"
                      label="Reasoning"
                      tone="muted"
                      :content="ev.text || ''"
                    />
                    <trace-block
                      v-else-if="ev.type === 'answer'"
                      label="Answer"
                      :content="ev.text || ''"
                    />
                  </template>

                  <trace-block
                    v-if="hasVars(s)"
                    label="Vars"
                    :content="varsText(s)"
                  />
                </div>
              </details>

              <!-- routing connector on the spine, in plain language -->
              <div v-if="s.routing" class="route">
                <span class="route__hook">↳</span>
                <template v-if="s.routing.stop && !s.routing.goto">
                  <span class="route__end">conversation ends here</span>
                </template>
                <template v-else>
                  <span class="route__lead">next</span>
                  <span class="route__arrow">→</span>
                  <code class="route__goto">{{ s.routing.goto || "(next step)" }}</code>
                  <span v-if="s.routing.value" class="route__why">
                    because it decided “{{ s.routing.value }}”
                  </span>
                </template>
              </div>
            </template>
          </div>
        </div>
      </details>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { NAlert, NEmpty, NSpin } from "naive-ui";
import { useTracesStore, type Trace, type TraceStep, type TraceToolCall } from "../../store/traces";
import TraceBlock from "./trace_block.vue";

// One step's activity, flattened into the order it actually happened: each tool
// call (matched to the round that requested it, by order) and the model's
// answers. Replaces the raw "LLM round N" dump - the tool calls are the story.
type StepEvent = { key: string } & (
  | { type: "tool"; tool: TraceToolCall }
  | { type: "answer"; text: string }
  | { type: "reasoning"; text: string }
);

const props = defineProps<{ chatUuid: string }>();

const store = useTracesStore();
const traces = ref<Trace[]>([]);
const loading = ref(true);
const error = ref("");

const openTurns = reactive<Record<string, boolean>>({});
const openSteps = reactive<Record<string, boolean>>({});

const totalSteps = computed(() => traces.value.reduce((n, t) => n + t.steps.length, 0));
const totalTools = computed(() =>
  traces.value.reduce((n, t) => n + t.steps.reduce((m, s) => m + (s.tool_calls?.length || 0), 0), 0),
);

onMounted(reload);

async function reload() {
  loading.value = true;
  error.value = "";
  try {
    const data = await store.getChatTraces(props.chatUuid);
    traces.value = data;
    // Everything collapsed by default (the panel is compact and scannable);
    // only the most recent turn is expanded, its steps still folded.
    for (const [ti, t] of data.entries()) {
      openTurns[t.id] = ti === data.length - 1;
      t.steps.forEach((_, si) => (openSteps[stepKey(t.id, si)] = false));
    }
  } catch (e) {
    error.value = (e as Error).message ?? "Unknown error";
  } finally {
    loading.value = false;
  }
}

const stepKey = (turnId: string, si: number) => `${turnId}#${si}`;
const isStepOpen = (turnId: string, si: number) => openSteps[stepKey(turnId, si)] === true;

function onTurnToggle(id: string, e: Event) {
  openTurns[id] = (e.target as HTMLDetailsElement).open;
}
function onStepToggle(turnId: string, si: number, e: Event) {
  openSteps[stepKey(turnId, si)] = (e.target as HTMLDetailsElement).open;
}

function expandAll() {
  for (const t of traces.value) {
    openTurns[t.id] = true;
    t.steps.forEach((_, si) => (openSteps[stepKey(t.id, si)] = true));
  }
}
function collapseAll() {
  for (const t of traces.value) {
    openTurns[t.id] = false;
    t.steps.forEach((_, si) => (openSteps[stepKey(t.id, si)] = false));
  }
}

// Plain-language step label + colour, by what the step actually does:
//   decision (classifier) — picks which branch runs next   (purple)
//   reply    — the message the customer actually receives   (green)
//   action   — internal work: gather data, write admin note (blue)
function kindLabel(s: TraceStep): string {
  if (s.kind === "classifier") return "decision";
  if (s.reply) return "reply";
  return "action";
}
function stepTone(s: TraceStep): string {
  if (s.kind === "classifier") return "#7c5cff";
  if (s.reply) return "#18a058";
  return "#2f6fed";
}

const turnTools = (t: Trace) =>
  t.steps.reduce((m, s) => m + (s.tool_calls?.length || 0), 0);

function dur(start: string, end: string): number {
  const ms = new Date(end).getTime() - new Date(start).getTime();
  return isFinite(ms) && ms >= 0 ? ms : 0;
}
function fmtDur(ms: number): string {
  if (!ms) return "—";
  return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(ms < 10000 ? 2 : 1)}s`;
}
function fmtTime(iso: string): string {
  const d = new Date(iso);
  return isNaN(d.getTime()) ? "" : d.toLocaleString();
}
const short = (s: string) => (s && s.length > 8 ? s.slice(0, 8) : s || "—");
const pad = (n: number) => String(n).padStart(2, "0");

// A slim proportional bar per step gives the whole turn a trace-waterfall read.
function barStyle(s: TraceStep, t: Trace) {
  const total = dur(t.started_at, t.finished_at) || 1;
  const w = Math.max(3, Math.round((dur(s.started_at, s.finished_at) / total) * 100));
  return { width: `${Math.min(w, 100)}%`, background: stepTone(s) };
}

function hasVars(s: TraceStep): boolean {
  return Object.keys(s.vars_in || {}).length > 0 || Object.keys(s.vars_out || {}).length > 0;
}
function varsText(s: TraceStep): string {
  const lines: string[] = [];
  const inV = s.vars_in || {};
  const outV = s.vars_out || {};
  for (const k of Object.keys(inV)) lines.push(`${k}: ${inV[k]}`);
  for (const k of Object.keys(outV)) {
    if (inV[k] !== outV[k]) lines.push(`${k} → ${outV[k]}`);
  }
  return lines.join("\n");
}

// Rebuild the true chronological sequence of a step. Rounds and tool_calls are
// both stored flat in call order; each round names the tools it requested, so
// we consume that many from the flat tool list per round, interleaving them
// with that round's answer text. Any tools not claimed by a round (shouldn't
// happen) are appended so nothing is ever hidden.
function stepEvents(s: TraceStep): StepEvent[] {
  const events: StepEvent[] = [];
  const tools = s.tool_calls || [];
  let ti = 0;
  let n = 0;
  for (const r of s.rounds || []) {
    if (r.reasoning) events.push({ key: `e${n++}`, type: "reasoning", text: r.reasoning });
    const req = (r.tools_requested || []).length;
    for (let k = 0; k < req && ti < tools.length; k++) {
      events.push({ key: `e${n++}`, type: "tool", tool: tools[ti++] });
    }
    if (r.response) events.push({ key: `e${n++}`, type: "answer", text: r.response });
  }
  while (ti < tools.length) events.push({ key: `e${n++}`, type: "tool", tool: tools[ti++] });
  return events;
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
.tv {
  font-size: 13px;
  padding-bottom: 32px;
}
.tv__state {
  display: flex;
  justify-content: center;
  padding-top: 72px;
}

/* toolbar */
.bar {
  position: sticky;
  top: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 2px 10px;
  margin-bottom: 6px;
  backdrop-filter: blur(6px);
}
.bar__count {
  font-size: 12px;
}
.bar__count b {
  font-variant-numeric: tabular-nums;
}
.bar__dim {
  opacity: 0.5;
}
.bar__spacer {
  flex: 1;
}
.lnk {
  font-family: ui-monospace, monospace;
  font-size: 11px;
  letter-spacing: 0.03em;
  color: inherit;
  background: transparent;
  border: none;
  border-bottom: 1px dotted rgba(128, 128, 128, 0.55);
  padding: 0 0 1px;
  cursor: pointer;
  opacity: 0.75;
}
.lnk:hover {
  opacity: 1;
  color: #2f6fed;
  border-color: #2f6fed;
}
.lnk:focus-visible {
  outline: 2px solid #2f6fed;
  outline-offset: 2px;
}

/* turn card */
.turn {
  border: 1px solid rgba(128, 128, 128, 0.22);
  border-radius: 10px;
  margin-bottom: 12px;
  overflow: hidden;
  background: rgba(128, 128, 128, 0.03);
}
.turn__sum {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 12px;
  cursor: pointer;
  list-style: none;
  user-select: none;
}
.turn__sum::-webkit-details-marker {
  display: none;
}
.turn__sum:hover {
  background: rgba(128, 128, 128, 0.05);
}
.turn__idx {
  font-family: ui-monospace, monospace;
  font-size: 12px;
  font-weight: 700;
  opacity: 0.45;
  font-variant-numeric: tabular-nums;
}
.turn__input {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 600;
}
.turn__metrics {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: ui-monospace, monospace;
  font-size: 11px;
  opacity: 0.7;
}
.metric {
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.metric--dur {
  min-width: 46px;
  text-align: right;
}
.metric--tool {
  color: #e0912f;
}
.chev {
  opacity: 0.4;
  transition: transform 0.18s ease;
  flex: none;
}
.turn[open] > .turn__sum .chev {
  transform: rotate(90deg);
}

.turn__body {
  padding: 2px 12px 14px;
}
.turn__info {
  display: flex;
  align-items: center;
  gap: 7px;
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  opacity: 0.5;
  padding: 2px 0 10px;
}
.turn__info .dot {
  opacity: 0.5;
}
.turn__err {
  margin-bottom: 10px;
}

/* timeline spine */
.steps {
  position: relative;
  padding-left: 28px;
}
.steps::before {
  content: "";
  position: absolute;
  left: 8px;
  top: 8px;
  bottom: 14px;
  width: 2px;
  background: rgba(128, 128, 128, 0.2);
  border-radius: 2px;
}
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
.step__spacer {
  flex: 1;
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

/* tool call: a compact collapsed row, name up front; expands to args+output */
.call {
  border: 1px solid rgba(224, 145, 47, 0.3);
  border-radius: 7px;
  overflow: hidden;
  background: rgba(224, 145, 47, 0.05);
}
.call--err {
  border-color: rgba(208, 48, 80, 0.4);
  background: rgba(208, 48, 80, 0.06);
}
.call__sum {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 6px 9px;
  cursor: pointer;
  list-style: none;
  user-select: none;
}
.call__sum::-webkit-details-marker {
  display: none;
}
.call__sum:hover {
  background: rgba(224, 145, 47, 0.08);
}
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
.call__spacer {
  flex: 1;
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
.call[open] > .call__sum .chev {
  transform: rotate(90deg);
}
.call__body {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 2px 9px 9px;
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

/* chips / tags */
.chip {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  padding: 1px 7px;
  border-radius: 20px;
  flex: none;
}
.chip--flow {
  color: #2f6fed;
  background: rgba(47, 111, 237, 0.13);
}
.chip--single {
  color: #18a058;
  background: rgba(24, 160, 88, 0.13);
}
.tag {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 5px;
  background: rgba(128, 128, 128, 0.15);
  white-space: nowrap;
  flex: none;
}
.tag--reply {
  color: #18a058;
  background: rgba(24, 160, 88, 0.14);
}
.tag--tool {
  color: #e0912f;
  background: rgba(224, 145, 47, 0.14);
}

@media (prefers-reduced-motion: reduce) {
  .chev {
    transition: none;
  }
}
</style>
