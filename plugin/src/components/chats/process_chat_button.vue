<template>
  <n-tooltip>
    <template #trigger>
      <n-button type="warning" size="small" ghost circle @click="open">
        <template #icon><sparkles-icon /></template>
      </n-button>
    </template>
    Обработать чат
  </n-tooltip>

  <n-modal v-model:show="show" preset="card" style="max-width: 640px" title="Обработать чат в базу знаний">
    <n-space vertical :wrap-item="false">
      <div>
        <n-text depth="3">База знаний</n-text>
        <n-select
          v-model:value="database"
          :options="baseOptions"
          :loading="basesLoading"
          placeholder="Выберите базу знаний бота"
        />
      </div>

      <n-button
        type="info"
        ghost
        :loading="processing"
        :disabled="!database"
        @click="runProcess"
      >
        {{ pairs.length ? "Обработать заново" : "Обработать" }}
      </n-button>

      <n-empty v-if="processed && !pairs.length" description="ИИ не нашёл, что добавить" />

      <div v-for="(pair, i) in pairs" :key="i" class="qa-item">
        <n-input v-model:value="pair.question" placeholder="Вопрос" />
        <n-input
          v-model:value="pair.answer"
          type="textarea"
          autosize
          placeholder="Ответ"
        />
        <n-button text type="error" @click="pairs.splice(i, 1)">
          <template #icon><delete-icon /></template>
        </n-button>
      </div>
    </n-space>

    <template #footer>
      <n-space justify="end">
        <n-button @click="show = false">Отмена</n-button>
        <n-button
          type="success"
          :loading="saving"
          :disabled="!database || !pairs.length"
          @click="save"
        >
          Добавить в базу ({{ pairs.length }})
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref } from "vue";
import {
  NButton,
  NEmpty,
  NInput,
  NModal,
  NSelect,
  NSpace,
  NText,
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

const baseOptions = computed(() =>
  bases.value.map((b) => ({ label: b.name, value: b.id }))
);

async function open() {
  show.value = true;
  pairs.value = [];
  processed.value = false;
  database.value = "";
  basesLoading.value = true;
  try {
    bases.value = await store.basesForChat(props.chat.admins);
    if (bases.value.length === 1) database.value = bases.value[0].id;
  } catch (e) {
    notification.error({ title: (e as Error).message });
  } finally {
    basesLoading.value = false;
  }
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
    pairs.value = await store.process(database.value, transcript());
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
    notification.success({ title: "Добавлено в базу", duration: 1500 });
    show.value = false;
  } catch (e) {
    notification.error({ title: (e as Error).message });
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.qa-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px;
  border: 1px solid rgba(128, 128, 128, 0.2);
  border-radius: 6px;
}
</style>
