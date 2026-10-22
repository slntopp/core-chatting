<template>
  <div class="learn-panel">
    <div class="learn-toolbar">
      <n-button
        type="info"
        :loading="proposing"
        :disabled="saving"
        @click="propose"
      >
        <template #icon><bulb-icon /></template>
        {{ proposed ? "Learn again" : "Learn" }}
      </n-button>
      <n-text depth="3">
        Otus reads the pages it already answers from, then this ticket. A take
        is a change to one of those pages, usually where you corrected it.
        Keep the changes that are true and drop the rest.
      </n-text>
    </div>

    <n-spin :show="proposing">
      <n-empty
        v-if="!takes.length"
        class="learn-empty"
        :description="
          proposed
            ? 'Nothing on the existing pages needs a change. Add a take, or save an empty list to skip.'
            : 'Press “Learn”'
        "
      />

      <div v-else class="learn-list">
        <div v-for="(take, i) in takes" :key="i" class="take">
          <div class="take-head">
            <n-text depth="3">Take {{ i + 1 }}</n-text>
            <n-button
              text
              type="error"
              size="small"
              title="Remove"
              @click="takes.splice(i, 1)"
            >
              <template #icon><delete-icon /></template>
            </n-button>
          </div>
          <n-input v-model:value="take.question" placeholder="Question" />
          <n-input
            v-model:value="take.answer"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Answer"
          />
        </div>
      </div>
    </n-spin>

    <div class="learn-actions">
      <n-button ghost :disabled="proposing || saving" @click="addTake">
        <template #icon><plus-icon /></template>
        Add take
      </n-button>
      <n-button
        type="success"
        :loading="saving"
        :disabled="proposing || !canSave"
        @click="save"
      >
        {{ kept.length ? `Save (${kept.length})` : "Save as skip" }}
      </n-button>
    </div>

    <p v-if="result && !result.prUrl" class="learn-result">{{ result.summary || result.why }}</p>

    <article v-else-if="result?.prUrl" class="pr">
      <header class="pr-top">
        <div class="pr-title-row">
          <h2 class="pr-title">Learn: {{ chatUuid }}</h2>
          <span class="pr-open">Open</span>
        </div>
        <p class="pr-merge">
          wants to merge into
          <span class="pr-branch">{{ result.base || "main" }}</span>
          from
          <span class="pr-branch">{{ result.head || "learn" }}</span>
        </p>
        <a class="pr-github" :href="result.prUrl" target="_blank" rel="noopener">
          View on GitHub
        </a>
      </header>

      <section v-if="result.summary" class="pr-body">
        <pre>{{ result.summary }}</pre>
      </section>

      <section class="pr-files">
        <div class="pr-files-bar">
          <span>Files changed</span>
          <span class="pr-count">{{ shownFiles.length }}</span>
          <span class="pr-add">+{{ totalAdd }}</span>
          <span class="pr-del">−{{ totalDel }}</span>
        </div>

        <p v-if="!shownFiles.length" class="pr-empty">The diff is on GitHub.</p>

        <div v-for="file in shownFiles" :key="file.path" class="pr-file">
          <div class="pr-file-head">
            <span class="pr-path">{{ file.path }}</span>
            <span class="pr-add">+{{ file.additions }}</span>
            <span class="pr-del">−{{ file.deletions }}</span>
          </div>
          <div v-if="file.lines.length" class="pr-diff-scroll">
            <table class="pr-diff">
              <tbody>
                <tr v-for="(line, i) in file.lines" :key="i" :class="line.kind">
                  <td class="gutter">{{ line.oldNo }}</td>
                  <td class="gutter">{{ line.newNo }}</td>
                  <td class="sign">{{ line.sign }}</td>
                  <td class="code">{{ line.text }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </article>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref } from "vue";
import { NButton, NEmpty, NInput, NSpin, NText, useNotification } from "naive-ui";
import {
  draftTakes,
  proposeTakes,
  saveTakes,
  type LearnAnswer,
  type LearnFile,
  type LearnTake,
} from "../../store/otus_learn";
import { parseUnifiedDiff, type DiffLine } from "./learn_diff";

// Otus caps a Save at 12 takes and a propose at 8.
const MAX_TAKES = 12;

const bulbIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/BulbOutline")
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline")
);
const plusIcon = defineAsyncComponent(() => import("@vicons/ionicons5/Add"));

const props = defineProps<{ chatUuid: string }>();
const notification = useNotification();

const takes = ref<LearnTake[]>([]);
const proposing = ref(false);
const saving = ref(false);
const proposed = ref(false);
const result = ref<LearnAnswer | null>(null);

const kept = computed(() =>
  takes.value.filter((t) => t.question.trim() && t.answer.trim())
);
// An empty list is a deliberate skip. A list of half-filled cards is not:
// Otus would reject it, so do not let them press Save on one.
const canSave = computed(
  () => takes.value.length === 0 || kept.value.length > 0
);

type ShownFile = LearnFile & { lines: DiffLine[] };

const shownFiles = computed((): ShownFile[] => {
  const answer = result.value;
  if (answer === null) return [];
  const parsed = parseUnifiedDiff(answer.patch ?? "");
  const listed = answer.files ?? [];
  if (parsed.length === 0) {
    return listed.map((file) => ({ ...file, lines: [] }));
  }
  return parsed.map((file) => {
    const stat = listed.find((item) => item.path === file.path);
    const additions = stat?.additions ?? file.lines.filter((line) => line.kind === "add").length;
    const deletions = stat?.deletions ?? file.lines.filter((line) => line.kind === "del").length;
    return { path: file.path, additions, deletions, lines: file.lines };
  });
});

const totalAdd = computed(() => shownFiles.value.reduce((sum, file) => sum + file.additions, 0));
const totalDel = computed(() => shownFiles.value.reduce((sum, file) => sum + file.deletions, 0));

// A restart of Otus drops the draft, so an empty answer here is normal: it
// just means there is nothing to pick up and Learn has to run again.
onMounted(async () => {
  try {
    const answer = await draftTakes(props.chatUuid);
    if (answer.takes?.length) {
      takes.value = answer.takes;
      proposed.value = true;
    }
  } catch {
    // Otus down or not configured: the Learn button reports it properly.
  }
});

async function propose() {
  proposing.value = true;
  result.value = null;
  try {
    const answer = await proposeTakes(props.chatUuid);
    takes.value = answer.takes ?? [];
    proposed.value = true;
    if (!takes.value.length) {
      notification.info({ title: answer.why, duration: 3000 });
    }
  } catch (error) {
    notification.error({ title: (error as Error).message });
  } finally {
    proposing.value = false;
  }
}

function addTake() {
  if (takes.value.length >= MAX_TAKES) {
    notification.warning({ title: `At most ${MAX_TAKES} takes`, duration: 3000 });
    return;
  }
  takes.value.push({ question: "", answer: "" });
}

async function save() {
  saving.value = true;
  try {
    const answer = await saveTakes(props.chatUuid, kept.value);
    result.value = answer;
    notification.success({ title: answer.why, duration: 3000 });
  } catch (error) {
    notification.error({ title: (error as Error).message });
  } finally {
    saving.value = false;
  }
}
</script>

<script lang="ts">
export default {
  name: "otus-learn-panel",
};
</script>

<style scoped>
.learn-panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.learn-toolbar {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.learn-empty {
  padding: 24px 0;
}

.learn-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.take {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px;
  border: 1px solid var(--n-border-color, #3336);
  border-radius: 6px;
}

.take-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.learn-actions {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.learn-result {
  margin: 0;
}

/* Primer colors from GitHub's pull-request files view, dark. */
.pr {
  border: 1px solid #30363d;
  border-radius: 6px;
  background: #0d1117;
  color: #e6edf3;
  overflow: hidden;
}

.pr-top {
  padding: 16px;
  border-bottom: 1px solid #30363d;
}

.pr-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.pr-title {
  margin: 0;
  font-size: 20px;
  font-weight: 400;
  line-height: 1.25;
}

.pr-open {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: 2em;
  background: #238636;
  color: #ffffff;
  font-size: 12px;
  font-weight: 600;
}

.pr-merge {
  margin: 8px 0 0;
  color: #8b949e;
  font-size: 14px;
}

.pr-branch {
  display: inline-block;
  padding: 0 6px;
  border-radius: 2em;
  background: #6e768166;
  color: #e6edf3;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
}

.pr-github {
  display: inline-block;
  margin-top: 12px;
  color: #2f81f7;
  font-size: 14px;
  text-decoration: none;
}

.pr-github:hover {
  text-decoration: underline;
}

.pr-body {
  padding: 16px;
  border-bottom: 1px solid #30363d;
}

.pr-body pre {
  margin: 0;
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.5;
}

.pr-files-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #161b22;
  border-bottom: 1px solid #30363d;
  font-size: 14px;
  font-weight: 600;
}

.pr-count {
  min-width: 20px;
  padding: 0 6px;
  border-radius: 2em;
  background: #30363d;
  color: #e6edf3;
  font-size: 12px;
  text-align: center;
}

.pr-add {
  color: #3fb950;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
  font-weight: 600;
}

.pr-del {
  color: #f85149;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
  font-weight: 600;
}

.pr-empty {
  margin: 0;
  padding: 16px;
  color: #8b949e;
}

.pr-file {
  border-top: 1px solid #30363d;
}

.pr-file-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #161b22;
}

.pr-path {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
}

.pr-diff-scroll {
  overflow-x: auto;
}

.pr-diff {
  width: 100%;
  border-collapse: collapse;
  background: #0d1117;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
  line-height: 20px;
}

.pr-diff td {
  padding: 0 8px;
  white-space: pre;
  vertical-align: top;
}

.pr-diff .gutter {
  width: 1%;
  color: #6e7681;
  text-align: right;
  user-select: none;
}

.pr-diff .sign {
  width: 1%;
  color: #6e7681;
  user-select: none;
}

.pr-diff .add {
  background: #2ea04326;
}

.pr-diff .add .gutter,
.pr-diff .add .sign {
  background: #3fb95026;
  color: #3fb950;
}

.pr-diff .del {
  background: #f8514926;
}

.pr-diff .del .gutter,
.pr-diff .del .sign {
  background: #f8514933;
  color: #f85149;
}

.pr-diff .hunk {
  background: #161b22;
  color: #8b949e;
}
</style>
