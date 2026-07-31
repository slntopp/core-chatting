<template>
  <div class="tv">
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

      <!-- one card per turn -->
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
            <span class="metric metric--dur">{{ fmtDur(durMs(t.started_at, t.finished_at)) }}</span>
          </span>
          <svg class="chev" viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
            <path d="M6 4l4 4-4 4" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </summary>

        <div class="turn__body">
          <div class="turn__info">
            <span>{{ fmtTime(t.started_at) }}</span>
            <span class="dot">•</span>
            <span>bot {{ shortId(t.bot_id) }}</span>
            <span class="dot">•</span>
            <span>{{ t.channel }}</span>
          </div>

          <n-alert v-if="t.error" type="warning" :bordered="false" class="turn__err">
            {{ t.error }}
          </n-alert>

          <div class="steps">
            <trace-step
              v-for="(s, si) in t.steps"
              :key="si"
              :step="s"
              :index="si"
              :turn-ms="durMs(t.started_at, t.finished_at)"
              :open="isStepOpen(t.id, si)"
              @toggle="(v) => setStepOpen(t.id, si, v)"
            />
          </div>
        </div>
      </details>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { NAlert, NEmpty, NSpin } from "naive-ui";
import { useTracesStore, type Trace } from "../../store/traces";
import { durMs, fmtDur, fmtTime, shortId } from "./trace_util";
import TraceStep from "./trace_step.vue";

const props = defineProps<{ chatUuid: string }>();

const store = useTracesStore();
const traces = ref<Trace[]>([]);
const loading = ref(true);
const error = ref("");

// Open state is held centrally (not in the child) so expand-all / collapse-all
// can drive every turn and step from here.
const openTurns = reactive<Record<string, boolean>>({});
const openSteps = reactive<Record<string, boolean>>({});

const totalSteps = computed(() => traces.value.reduce((n, t) => n + t.steps.length, 0));
const totalTools = computed(() =>
  traces.value.reduce((n, t) => n + t.steps.reduce((m, s) => m + (s.tool_calls?.length || 0), 0), 0),
);
const turnTools = (t: Trace) => t.steps.reduce((m, s) => m + (s.tool_calls?.length || 0), 0);

onMounted(reload);

async function reload() {
  loading.value = true;
  error.value = "";
  try {
    const data = await store.getChatTraces(props.chatUuid);
    traces.value = data;
    // Compact by default: only the latest turn expanded, all steps folded.
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
const setStepOpen = (turnId: string, si: number, open: boolean) =>
  (openSteps[stepKey(turnId, si)] = open);

function onTurnToggle(id: string, e: Event) {
  openTurns[id] = (e.target as HTMLDetailsElement).open;
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

const pad = (n: number) => String(n).padStart(2, "0");
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

/* timeline spine (the numbered step nodes live in trace-step) */
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

/* turn-mode chip */
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

@media (prefers-reduced-motion: reduce) {
  .chev {
    transition: none;
  }
}
</style>
