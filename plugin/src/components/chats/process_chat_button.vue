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
        {{ pairs.length ? "Process again" : "Process" }}
      </n-button>
    </div>

    <n-spin :show="processing">
      <n-empty
        v-if="!pairs.length"
        class="pc-empty"
        :description="
          processed ? 'Nothing worth adding was found' : 'Press “Process”'
        "
      />

      <div v-else class="pc-list">
        <div class="pc-count">
          Pairs found: <b>{{ pairs.length }}</b>
        </div>
        <div v-for="(pair, i) in pairs" :key="i" class="qa-item">
          <div class="qa-head">
            <span class="qa-num">{{ i + 1 }}</span>
            <n-button
              text
              type="error"
              size="small"
              title="Remove"
              @click="pairs.splice(i, 1)"
            >
              <template #icon><delete-icon /></template>
            </n-button>
          </div>
          <n-input
            v-model:value="pair.question"
            placeholder="Question"
            class="qa-question"
          />
          <n-input
            v-model:value="pair.answer"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Answer"
          />
        </div>
      </div>
    </n-spin>

    <template #footer>
      <n-space justify="end">
        <n-button @click="show = false">Cancel</n-button>
        <n-button
          type="success"
          :loading="saving"
          :disabled="!database || !pairs.length"
          @click="save"
        >
          Add to knowledge base ({{ pairs.length }})
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, watch } from "vue";
import {
  NButton,
  NEmpty,
  NInput,
  NModal,
  NSelect,
  NSpace,
  NSpin,
  NTooltip,
  useNotification,
} from "naive-ui";
import { Chat, Kind } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { useUsersStore } from "../../store/users";
import { QAPair, useChatProcessStore } from "../../store/chat_process";

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
const pairs = ref<QAPair[]>([]);

const botAccount = ref("");

const baseOptions = computed(() =>
  bases.value.map((b) => ({ label: b.name, value: b.id }))
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
  pairs.value = [];
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
    pairs.value = await store.process(database.value, botAccount.value, transcript());
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
    await store.append(
      database.value,
      pairs.value.filter((p) => p.question.trim() && p.answer.trim())
    );
    notification.success({ title: "Added to knowledge base", duration: 1500 });
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
.qa-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.qa-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--main-app-primary-color, #6b7280);
  color: #fff;
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
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
