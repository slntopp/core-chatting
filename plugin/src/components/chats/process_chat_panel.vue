<template>
  <div class="pc-panel">
    <div class="pc-title">Process chat into knowledge base</div>

    <div v-if="basesLoading" class="pc-status">
      <n-spin size="small" />
      <span>Checking bot knowledge bases…</span>
    </div>

    <div v-else-if="!botAccount" class="pc-status pc-status--muted">
      This chat has no bot knowledge base.
    </div>

    <template v-else>
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
          <div v-if="p.match >= 0 && p.old" class="qa-old">
            <span class="qa-old-label">Was:</span>
            <del>{{ p.old.answer }}</del>
          </div>
        </div>
      </div>
    </n-spin>

    <div v-if="proposals.length" class="pc-save">
      <n-button
        type="success"
        :loading="saving"
        :disabled="!database || !includedCount"
        @click="save"
      >
        Save to knowledge base ({{ includedCount }})
      </n-button>
    </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, watch } from "vue";
import {
  NButton,
  NCheckbox,
  NEmpty,
  NInput,
  NSelect,
  NSpin,
  NTag,
  useNotification,
} from "naive-ui";
import { Chat, Kind } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { useUsersStore } from "../../store/users";
import { QAProposal, useChatProcessStore } from "../../store/chat_process";

interface ProposalRow extends QAProposal {
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

const database = ref("");
const bases = ref<{ id: string; name: string }[]>([]);
const basesLoading = ref(false);
const processing = ref(false);
const processed = ref(false);
const saving = ref(false);
const proposals = ref<ProposalRow[]>([]);
const botAccount = ref("");

const baseOptions = computed(() =>
  bases.value.map((b) => ({ label: b.name, value: b.id }))
);
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
// admins and users). Empty account -> no bot -> the panel stays hidden.
async function resolve() {
  botAccount.value = "";
  bases.value = [];
  proposals.value = [];
  processed.value = false;
  const accounts = [...(props.chat.admins || []), ...(props.chat.users || [])];
  if (!accounts.length) return;
  basesLoading.value = true;
  try {
    const r = await store.resolveBot(accounts);
    botAccount.value = r.account;
    bases.value = r.databases;
    database.value = bases.value.length === 1 ? bases.value[0].id : "";
  } catch {
    // no bot / unreachable -> keep hidden, don't spam the operator
  } finally {
    basesLoading.value = false;
  }
}

onMounted(resolve);
watch(() => props.chat.uuid, resolve);

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
    // Send only the approved proposals; the backend applies them against the
    // base (replace at match / append). The full base never leaves the server.
    const items = kept.value.map((p) => ({
      match: p.match,
      question: p.question,
      answer: p.answer,
    }));
    await store.save(database.value, items);
    notification.success({ title: "Saved to knowledge base", duration: 1500 });
    proposals.value = [];
    processed.value = false;
  } catch (e) {
    notification.error({ title: (e as Error).message });
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.pc-panel {
  margin-top: 8px;
}
.pc-title {
  font-weight: 600;
  margin-bottom: 12px;
}
.pc-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  font-size: 0.9rem;
}
.pc-status--muted {
  opacity: 0.6;
}
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
  padding: 24px 0;
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
  max-height: 45vh;
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
.pc-save {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
