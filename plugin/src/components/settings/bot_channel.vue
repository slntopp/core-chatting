<template>
  <div>
    <n-space style="padding: 6px 24px" align="center" justify="space-between">
      <n-h3>Bot channel links</n-h3>
      <n-button ghost type="success" @click="openCreate">
        <template #icon><plus-icon /></template>
        Add
      </n-button>
    </n-space>

    <n-table :bordered="false" :single-line="false">
      <thead>
        <tr>
          <th>Bot</th>
          <th>Linked account</th>
          <th>Custom name</th>
          <th>Skip review</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="link in links" :key="link.channelId">
          <td>{{ link.botName }}</td>
          <td>{{ link.accountUuid }}</td>
          <td>{{ link.customName || "-" }}</td>
          <td>{{ link.skipReview ? "Yes" : "No" }}</td>
          <td>
            <n-space :wrap-item="false">
              <n-button ghost text type="info" @click="openEdit(link)">
                <template #icon><save-icon /></template>
              </n-button>
              <n-popconfirm
                @positive-click="onDelete(link)"
                positive-text="Unlink"
                negative-text="Cancel"
              >
                <template #trigger>
                  <n-button
                    ghost
                    text
                    type="error"
                    :loading="deletingId === link.channelId"
                  >
                    <template #icon><delete-icon /></template>
                  </n-button>
                </template>
                This will remove the core_chatting channel from
                "{{ link.botName }}". The bot will stop receiving messages
                from this account.
              </n-popconfirm>
            </n-space>
          </td>
        </tr>
        <tr v-if="!links.length">
          <td colspan="5">
            <n-text depth="3">No bot channels linked yet.</n-text>
          </td>
        </tr>
      </tbody>
    </n-table>

    <n-modal v-model:show="showForm" preset="card" style="max-width: 600px" :title="editingLink ? 'Edit bot channel link' : 'Add bot channel link'">
      <n-space vertical :wrap-item="false">
        <div>
          <n-text depth="3">Bot instance</n-text>
          <n-select
            v-model:value="form.botId"
            :options="instanceOptions"
            :disabled="!!editingLink"
            filterable
            placeholder="Select a 'bots' instance"
          />
        </div>

        <div>
          <n-text depth="3">Linked account (bot identity in core-chatting)</n-text>
          <account-select-single v-model="form.accountUuid" />
        </div>

        <div>
          <n-text depth="3">Custom name</n-text>
          <n-input v-model:value="form.customName" placeholder="e.g. Kirill" />
        </div>

        <div class="skip-review-row">
          <n-switch v-model:value="form.skipReview">
            <template #checked>Skip review</template>
            <template #unchecked>Requires review</template>
          </n-switch>
          <span style="max-width: 60%;">
            When disabled, every bot reply must be approved by a human
            operator in core-chatting before the client sees it.
          </span>
        </div>
      </n-space>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showForm = false">Cancel</n-button>
          <n-button
            type="info"
            :loading="isSubmitting"
            :disabled="!form.botId || !form.accountUuid"
            @click="submit"
          >
            {{ editingLink ? "Save" : "Add" }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, reactive, ref } from "vue";
import {
  NButton,
  NInput,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NSwitch,
  NTable,
  NText,
  NH3,
  useNotification,
} from "naive-ui";
import AccountSelectSingle from "../users/account_select_single.vue";
import {
  BotChannelLink,
  useBotChannelStore,
} from "../../store/bot_channel";

const plusIcon = defineAsyncComponent(() => import("@vicons/ionicons5/Add"));
const saveIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SaveOutline")
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline")
);

const store = useBotChannelStore();
const notification = useNotification();

const emit = defineEmits(["refresh"]);

const links = computed(() => store.links);
const instanceOptions = computed(() =>
  store.instances.map((i) => ({ label: `${i.title} (${i.status})`, value: i.botId }))
);

const showForm = ref(false);
const editingLink = ref<BotChannelLink | null>(null);
const isSubmitting = ref(false);
const deletingId = ref<string | null>(null);

const form = reactive({
  botId: "",
  accountUuid: "",
  customName: "",
  skipReview: false,
});

onMounted(() => {
  refresh();
});

async function refresh() {
  await store.fetchAll();
  emit("refresh");
}

function openCreate() {
  editingLink.value = null;
  form.botId = "";
  form.accountUuid = "";
  form.customName = "";
  form.skipReview = false;
  showForm.value = true;
}

function openEdit(link: BotChannelLink) {
  editingLink.value = link;
  form.botId = link.botId;
  form.accountUuid = link.accountUuid;
  form.customName = link.customName;
  form.skipReview = link.skipReview;
  showForm.value = true;
}

async function submit() {
  isSubmitting.value = true;
  try {
    if (editingLink.value) {
      await store.editLink(
        editingLink.value,
        form.accountUuid,
        form.customName,
        form.skipReview
      );
    } else {
      await store.addLink(
        form.botId,
        form.accountUuid,
        form.customName,
        form.skipReview
      );
    }
    notification.success({ title: "Done", duration: 1500 });
    showForm.value = false;
    await refresh();
  } catch (error) {
    notification.error({ title: (error as Error).message ?? String(error) });
  } finally {
    isSubmitting.value = false;
  }
}

async function onDelete(link: BotChannelLink) {
  deletingId.value = link.channelId;
  try {
    await store.removeLink(link.botId, link.channelId);
    notification.success({ title: "Unlinked", duration: 1500 });
    await refresh();
  } catch (error) {
    notification.error({ title: (error as Error).message ?? String(error) });
  } finally {
    deletingId.value = null;
  }
}
</script>

<script lang="ts">
export default {
  name: "settings-bot-channel",
};
</script>

<style scoped>
.skip-review-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.skip-review-row span {
  font-size: 0.9rem;
  opacity: 0.8;
}
</style>
