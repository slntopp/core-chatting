<template>
  <n-tooltip v-if="botAccount">
    <template #trigger>
      <n-button type="warning" size="small" ghost circle @click="open">
        <template #icon><sparkles-icon /></template>
      </n-button>
    </template>
    Process chat
  </n-tooltip>

  <n-modal
    v-model:show="show"
    preset="card"
    class="process-chat-modal"
    title="Process chat into knowledge base"
  >
    <div class="pc-toolbar">
      <n-select
        v-model:value="database"
        :options="baseOptions"
        :loading="basesLoading"
        placeholder="Select a knowledge base"
        class="pc-base-select"
      />
      <n-button
        type="info"
        :loading="processing"
        :disabled="!database"
        @click="runProcess"
      >
        <template #icon><sparkles-icon /></template>
        {{ processed ? "Process again" : "Process" }}
      </n-button>
    </div>

    <n-spin :show="processing">
      <n-empty
        v-if="!proposals.length"
        class="pc-empty"
        :description="
          processed ? 'Nothing worth adding was found' : 'Press “Process”'
        "
      />

      <div v-else class="pc-list">
        <div class="pc-count">
          New: <b>{{ addCount }}</b> · Replacements: <b>{{ replaceCount }}</b>
        </div>
        <div
          v-for="(p, i) in proposals"
          :key="i"
          class="qa-item"
          :class="p.match >= 0 ? 'qa-item--replace' : 'qa-item--add'"
        >
          <div class="qa-head">
            <n-tag size="small" :type="p.match >= 0 ? 'warning' : 'success'">
              {{ p.match >= 0 ? `Replace #${p.match + 1}` : "New" }}
            </n-tag>
            <n-checkbox v-model:checked="p.include">Include</n-checkbox>
            <n-button
              text
              type="error"
              size="small"
              title="Remove"
              @click="proposals.splice(i, 1)"
            >
              <template #icon><delete-icon /></template>
            </n-button>
          </div>
          <n-input
            v-model:value="p.question"
            placeholder="Question"
            class="qa-question"
          />
          <n-input
            v-model:value="p.answer"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Answer"
          />
          <div v-if="p.match >= 0 && existing[p.match]" class="qa-old">
            <span class="qa-old-label">Was:</span>
            <del>{{ existing[p.match].answer }}</del>
          </div>
        </div>
      </div>
    </n-spin>

    <template #footer>
      <n-space justify="end">
        <n-button @click="show = false">Cancel</n-button>
        <n-button
          type="success"
          :loading="saving"
          :disabled="!database || !includedCount"
          @click="save"
        >
          Save to knowledge base ({{ includedCount }})
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, watch } from "vue";
import {
  NButton,
  NCheckbox,
  NEmpty,
  NInput,
  NModal,
  NSelect,
  NSpace,
  NSpin,
  NTag,
  NTooltip,
  useNotification,
} from "naive-ui";
import { Chat, Kind } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { useUsersStore } from "../../store/users";
import {
  QAKnowledge,
  QAPair,
  useChatProcessStore,
} from "../../store/chat_process";

interface ProposalRow extends QAPair {
  match: number;
  include: boolean;
}

const sparklesIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SparklesOutline")
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline")
);

const props = defineProps<{ chat: Chat }>();

const ccStore = useCcStore();
const usersStore = useUsersStore();
const store = useChatProcessStore();
const notification = useNotification();

const show = ref(false);
const database = ref("");
const bases = ref<{ id: string; name: string }[]>([]);
const basesLoading = ref(false);
const processing = ref(false);
const processed = ref(false);
const saving = ref(false);
const proposals = ref<ProposalRow[]>([]);
const qaKnowledge = ref<QAKnowledge | null>(null);

const botAccount = ref("");

const baseOptions = computed(() =>
  bases.value.map((b) => ({ label: b.name, value: b.id }))
);
const existing = computed(() => qaKnowledge.value?.records || []);
const kept = computed(() =>
  proposals.value.filter(
    (p) => p.include && (p.question.trim() || p.answer.trim())
  )
);
const includedCount = computed(() => kept.value.length);
const addCount = computed(() => kept.value.filter((p) => p.match < 0).length);
const replaceCount = computed(
  () => kept.value.filter((p) => p.match >= 0).length
);

// Ask the backend which chat participant is a core_chatting bot (checking both
// admins and users). Sets botAccount + bases; empty account -> button hidden.
async function resolve() {
  botAccount.value = "";
  bases.value = [];
  const accounts = [...(props.chat.admins || []), ...(props.chat.users || [])];
  if (!accounts.length) return;
  basesLoading.value = true;
  try {
    const r = await store.resolveBot(accounts);
    botAccount.value = r.account;
    bases.value = r.databases;
  } catch {
    // no bot / unreachable -> keep button hidden, don't spam the operator
  } finally {
    basesLoading.value = false;
  }
}

onMounted(resolve);
watch(() => props.chat.uuid, resolve);

function open() {
  show.value = true;
  proposals.value = [];
  qaKnowledge.value = null;
  processed.value = false;
  database.value = bases.value.length === 1 ? bases.value[0].id : "";
}

function transcript() {
  return ccStore
    .chat_messages(props.chat)
    .filter((m) => !m.underReview && m.kind !== Kind.ADMIN_ONLY && m.content)
    .map((m) => ({
      author: usersStore.users.get(m.sender)?.title || "User",
      text: m.content,
    }));
}

async function runProcess() {
  processing.value = true;
  try {
    const res = await store.process(
      database.value,
      botAccount.value,
      transcript()
    );
    qaKnowledge.value = res.qa_knowledge;
    proposals.value = res.items.map((it) => ({ ...it, include: true }));
    processed.value = true;
  } catch (e) {
    notification.error({ title: (e as Error).message });
  } finally {
    processing.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    // Final record set: existing base, with approved replacements applied in
    // place and approved new pairs appended.
    const final: QAPair[] = existing.value.map((r) => ({
      question: r.question,
      answer: r.answer,
    }));
    for (const p of kept.value) {
      const pair = { question: p.question, answer: p.answer };
      if (p.match >= 0 && p.match < final.length) final[p.match] = pair;
      else final.push(pair);
    }
    await store.save(database.value, qaKnowledge.value, final);
    notification.success({ title: "Saved to knowledge base", duration: 1500 });
    show.value = false;
  } catch (e) {
    notification.error({ title: (e as Error).message });
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.pc-toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}
.pc-base-select {
  flex: 1;
}
.pc-empty {
  padding: 40px 0;
}
.pc-count {
  margin-bottom: 12px;
  opacity: 0.7;
  font-size: 0.9rem;
}
.pc-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  max-height: 60vh;
  overflow-y: auto;
  padding-right: 4px;
}
.qa-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border: 1px solid rgba(128, 128, 128, 0.2);
  border-radius: 10px;
  background: rgba(128, 128, 128, 0.04);
}
.qa-item--add {
  border-left: 3px solid #18a058;
}
.qa-item--replace {
  border-left: 3px solid #f0a020;
}
.qa-head {
  display: flex;
  align-items: center;
  gap: 12px;
}
.qa-head :deep(.n-button) {
  margin-left: auto;
}
.qa-old {
  font-size: 0.85rem;
  opacity: 0.7;
}
.qa-old-label {
  margin-right: 6px;
  opacity: 0.7;
}
.qa-question :deep(input) {
  font-weight: 600;
}
</style>

<style>
.process-chat-modal {
  width: 900px;
  max-width: 92vw;
}
</style>
