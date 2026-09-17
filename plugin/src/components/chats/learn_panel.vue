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
        Otus reads this ticket and proposes what is worth remembering. One
        ticket can hold several topics; keep the takes that are true, drop the
        rest.
      </n-text>
    </div>

    <n-spin :show="proposing">
      <n-empty
        v-if="!takes.length"
        class="learn-empty"
        :description="
          proposed
            ? 'Nothing worth remembering was found. Add a take, or save an empty list to skip.'
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

    <div v-if="result" class="learn-result">
      <n-text>{{ result.summary || result.why }}</n-text>
      <a v-if="result.prUrl" :href="result.prUrl" target="_blank" rel="noopener">
        Open the pull request
      </a>
    </div>
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
  type LearnTake,
} from "../../store/otus_learn";

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
  display: flex;
  flex-direction: column;
  gap: 6px;
}
</style>
