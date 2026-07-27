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
                This will remove the core_chatting channel from "{{
                  link.botName
                }}". The bot will stop receiving messages from this account.
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

    <n-modal
      v-model:show="showForm"
      preset="card"
      style="max-width: 600px"
      :title="editingLink ? 'Edit bot channel link' : 'Add bot channel link'"
    >
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
          <n-text depth="3"
            >Linked account (bot identity in core-chatting)</n-text
          >
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
          <span style="max-width: 60%">
            When disabled, every bot reply must be approved by a human operator
            in core-chatting before the client sees it.
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

    <n-collapse style="padding: 6px 24px; margin-top: 24px; max-width: 90%">
      <n-collapse-item title="Bot behavior" name="bot-behavior">
        <div class="bots_config_switches">
          <div class="bots_config_switch">
            <n-switch class="switch" v-model:value="botConfig.enable">
              <template #checked> Active </template>
              <template #unchecked> Disabled </template>
            </n-switch>

            <span> Enable bot in new chats. </span>
          </div>

          <div class="bots_config_switch">
            <n-switch class="switch" v-model:value="botConfig.review">
              <template #checked> Review </template>
              <template #unchecked> No review </template>
            </n-switch>

            <span>
              Enable Pre-Moderation mode. New bot messages will be visible only
              to administrators.
            </span>
          </div>

          <div class="bots_config_switch">
            <n-switch class="switch" v-model:value="botConfig.initiator">
              <template #checked> Active </template>
              <template #unchecked> Disabled </template>
            </n-switch>

            <span>
              Enable Hybrid-Moderation mode. The first message from the bot will
              be published, and subsequent messages will be visible only to
              administrators.
            </span>
          </div>

          <div class="bots_config_switch">
            <n-switch class="switch" v-model:value="botConfig.emergency">
              <template #checked> Active </template>
              <template #unchecked> Disabled </template>
            </n-switch>

            <span>
              Enable EMERGENCY mode. The bot will respond according to the
              instructions given to it.
            </span>
          </div>
        </div>

        <n-space style="margin-top: 20px">
          <n-text>Promt</n-text>
        </n-space>
        <n-input
          v-model:value="botConfig.prompt"
          type="textarea"
          autosize
          placeholder="Bot Promt"
        />

        <n-space style="margin-top: 20px">
          <n-text>Custom values</n-text>
        </n-space>

        <n-table :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th>Key</th>
              <th>Value</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, index) in botCustomValues">
              <td>
                <n-input v-model:value="value.key" />
              </td>
              <td>
                <n-input v-model:value="value.value" />
              </td>
              <td>
                <n-button
                  ghost
                  @click="deleteBotCustomValue(index)"
                  type="warning"
                  >Delete</n-button
                >
              </td>
            </tr>

            <tr>
              <td></td>
              <td></td>
              <td>
                <n-space justify="end">
                  <n-button ghost @click="addBotCustomValue" type="success"
                    >Add</n-button
                  >
                </n-space>
              </td>
            </tr>
          </tbody>
        </n-table>

        <n-space justify="end" style="margin: 10px 0">
          <n-button
            :loading="isBotSaving"
            ghost
            type="info"
            @click="submitBotConfig"
          >
            Update
          </n-button>
        </n-space>
      </n-collapse-item>
    </n-collapse>
  </div>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  onMounted,
  reactive,
  ref,
  watch,
} from "vue";
import {
  NButton,
  NCollapse,
  NCollapseItem,
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
import { BotChannelLink, useBotChannelStore } from "../../store/bot_channel";
import { Bot, Defaults } from "../../connect/cc/cc_pb";
import { useDefaultsStore } from "../../store/defaults";

const plusIcon = defineAsyncComponent(() => import("@vicons/ionicons5/Add"));
const saveIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SaveOutline"),
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline"),
);

const store = useBotChannelStore();
const defaultsStore = useDefaultsStore();
const notification = useNotification();

const emit = defineEmits(["refresh"]);

const links = computed(() => store.links);
const instanceOptions = computed(() =>
  store.instances.map((i) => ({
    label: `${i.title} (${i.status})`,
    value: i.botId,
  })),
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

const botConfig = reactive({
  prompt: "",
  enable: false,
  review: false,
  initiator: false,
  emergency: false,
});
const botCustomValues = ref<{ key: string; value: string }[]>([]);
const isBotSaving = ref(false);

// core-chatting's department-wide bot config (Active/Review/Hybrid/Emergency,
// prompt, custom values) - lives on cc.Defaults.Bot, edited here as one form.
watch(
  () => defaultsStore.bot,
  (bot) => {
    if (!bot) return;
    botConfig.prompt = bot.prompt;
    botConfig.enable = bot.enable;
    botConfig.review = bot.review;
    botConfig.initiator = bot.initiator;
    botConfig.emergency = bot.emergency;
    botCustomValues.value = Object.entries(bot.values).map(([key, value]) => ({
      key,
      value,
    }));
  },
  { immediate: true },
);

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
        form.skipReview,
      );
    } else {
      await store.addLink(
        form.botId,
        form.accountUuid,
        form.customName,
        form.skipReview,
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

function addBotCustomValue() {
  botCustomValues.value.push({ key: "", value: "" });
}

function deleteBotCustomValue(index: number) {
  botCustomValues.value = botCustomValues.value.filter((_, i) => i !== index);
}

async function submitBotConfig() {
  isBotSaving.value = true;
  try {
    await defaultsStore.update_defaults(
      new Defaults({
        admins: defaultsStore.admins,
        departments: defaultsStore.departments,
        gateways: defaultsStore.gateways,
        metrics: defaultsStore.metrics.reduce(
          (result, metric) => ({ ...result, [metric.key]: metric }),
          {},
        ),
        templates: defaultsStore.templates.reduce(
          (result, template) => ({
            ...result,
            [template.name]: template.content,
          }),
          {},
        ),
        bot: new Bot({
          ...botConfig,
          values: botCustomValues.value.reduce<Record<string, string>>(
            (result, { key, value }) => ({ ...result, [key]: value }),
            {},
          ),
        }),
      }),
    );
    notification.success({ title: "Done", duration: 1000 });
  } catch (error) {
    notification.error({ title: (error as Error).message ?? String(error) });
  } finally {
    isBotSaving.value = false;
  }
}
</script>

<script lang="ts">
export default {
  name: "settings-bot-channel",
};
</script>

<style scoped>
.bots_config_switches {
  display: flex;
  flex-direction: column;
  justify-content: baseline;
}

.bots_config_switch {
  display: flex;
  margin: 5px 0px;
  align-items: center;
}

.bots_config_switch span {
  font-size: 1.2rem;
}

.bots_config_switches .switch {
  width: 200px;
  justify-content: normal;
}

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
