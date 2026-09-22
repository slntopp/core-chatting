<template>
  <n-tooltip>
    <template #trigger>
      <n-button
        type="info"
        size="small"
        ghost
        circle
        @click="isBotSettingsOpen = true"
      >
        <template #icon> <bot-icon /> </template>
      </n-button>
    </template>
    Bot settings
  </n-tooltip>

  <n-tooltip>
    <template #trigger>
      <n-button
        size="small"
        ghost
        circle
        :color="appStore.isCopilotOpen ? COPILOT_COLOR : undefined"
        @click="appStore.isCopilotOpen = !appStore.isCopilotOpen"
      >
        <template #icon> <sparkles-icon /> </template>
      </n-button>
    </template>
    Ask the bot about this chat
  </n-tooltip>

  <n-tooltip>
    <template #trigger>
      <n-button
        type="default"
        size="small"
        ghost
        circle
        @click="isTracesOpen = true"
      >
        <template #icon> <bug-icon /> </template>
      </n-button>
    </template>
    Bot debug traces
  </n-tooltip>

  <n-tooltip>
    <template #trigger>
      <n-button
        type="warning"
        size="small"
        ghost
        circle
        @click="isLearnOpen = true"
      >
        <template #icon> <bulb-icon /> </template>
      </n-button>
    </template>
    Learn from this chat
  </n-tooltip>

  <n-modal v-model:show="isBotSettingsOpen">
    <n-card
      title="Bot settings"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
      style="width: 780px; max-width: 92vw; min-height: 300px"
    >
      <div class="bot_state_settings_field">
        <span>Bot state:</span>
        <n-switch v-model:value="botState.enabled">
          <template #checked> active </template>
          <template #unchecked> disabled (no bot) </template>
        </n-switch>
      </div>

      <div class="bot_state_settings_field">
        <span>Review:</span>
        <n-switch v-model:value="botState.review">
          <template #unchecked> open mode (skip preview) </template>
          <template #checked> review mode </template>
        </n-switch>
      </div>

      <!-- Operators only: this rung decides whether the bot may write to the
           client on its own, so the chat's owner - the client - must not see
           or reach it. The server refuses them the key as well. -->
      <div v-if="isChatAdmin" class="bot_state_settings_field bot_state_otus">
        <span>Otus mode:</span>
        <div>
          <n-radio-group v-model:value="otusMode" name="chat-otus-mode">
            <n-space>
              <n-radio
                v-for="option in OTUS_CHAT_MODES"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </n-radio>
            </n-space>
          </n-radio-group>
          <n-text depth="3" tag="div" style="margin-top: 4px">
            {{ otusHint }}
          </n-text>
          <n-text depth="3" tag="div">
            Install-wide now: {{ installOtusModeLabel }}.
          </n-text>
        </div>
      </div>

      <div class="bot_state_settings_actions">
        <n-button
          style="margin-right: 15px"
          ghost
          type="info"
          @click="isBotSettingsOpen = false"
        >
          Cancel
        </n-button>
        <n-button
          ghost
          type="warning"
          @click="saveBotState"
          :loading="isSaveBotStateLoading"
        >
          Save
        </n-button>
      </div>

      <n-divider />

      <process-chat-panel :chat="chat" />
    </n-card>
  </n-modal>

  <n-drawer v-model:show="isTracesOpen" :width="620" placement="right">
    <n-drawer-content title="Bot debug traces" closable :native-scrollbar="false">
      <trace-viewer v-if="isTracesOpen" :chat-uuid="chat.uuid" lane="customer" />
    </n-drawer-content>
  </n-drawer>

  <n-drawer v-model:show="isLearnOpen" :width="620" placement="right">
    <n-drawer-content title="Learn from this chat" closable :native-scrollbar="false">
      <learn-panel v-if="isLearnOpen" :chat-uuid="chat.uuid" />
    </n-drawer-content>
  </n-drawer>

  <n-tooltip>
    <template #trigger>
      <n-button type="success" size="small" ghost circle @click="copyLink">
        <template #icon> <copy-icon /> </template>
      </n-button>
    </template>
    Copy chat link
  </n-tooltip>

  <n-tooltip>
    <template #trigger>
      <n-button type="info" size="small" ghost circle @click="refresh">
        <template #icon> <refresh-icon /> </template>
      </n-button>
    </template>
    Refresh chat
  </n-tooltip>

  <n-tooltip>
    <template #trigger>
      <n-popconfirm @positive-click="deleteChat">
        <template #trigger>
          <n-button type="error" size="small" ghost circle>
            <template #icon> <delete-icon /> </template>
          </n-button>
        </template>
        Are you sure you want to delete the chat?
      </n-popconfirm>
    </template>
    Delete chat
  </n-tooltip>

  <n-tooltip v-if="buttonTitle">
    <template #trigger>
      <n-button
        ghost
        circle
        type="warning"
        size="small"
        @mouseenter="hoverInstancesButton"
      >
        <template #icon> <list-icon /> </template>
      </n-button>
    </template>
    Instances list
  </n-tooltip>

  <n-popover
    scrollable
    trigger="click"
    placement="bottom"
    content-style="padding: 0"
    v-if="commands.length > 0"
  >
    <template #trigger>
      <n-tooltip>
        <template #trigger>
          <n-button
            ref="commandsButton"
            type="success"
            size="small"
            ghost
            circle
          >
            <template #icon> <console-icon /> </template>
          </n-button>
        </template>
        Commands
      </n-tooltip>
    </template>

    <n-list hoverable clickable>
      <n-list-item
        v-for="command of commands"
        :key="command.key"
        @click="sendCommand(command.key)"
      >
        {{ command.key }} ({{ command.description }})
      </n-list-item>
    </n-list>
  </n-popover>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import {
  NButton,
  NList,
  NListItem,
  NPopover,
  NTooltip,
  NPopconfirm,
  useNotification,
  NModal,
  NCard,
  NSwitch,
  NDivider,
  NDrawer,
  NDrawerContent,
  NRadio,
  NRadioGroup,
  NSpace,
  NText,
} from "naive-ui";
import TraceViewer from "./trace_viewer.vue";
import LearnPanel from "./learn_panel.vue";

import { ConnectError } from "@connectrpc/connect";
import {
  Chat,
  Kind,
  Message,
  Role,
  SetBotStateRequest,
} from "../../connect/cc/cc_pb";
import { useAppStore } from "../../store/app.ts";
import { useCcStore } from "../../store/chatting.ts";
import { addToClipboard } from "../../functions.ts";
import { storeToRefs } from "pinia";
import { useUsersStore } from "../../store/users.ts";
import { onUnmounted } from "vue";
import ProcessChatPanel from "./process_chat_panel.vue";
import { useDefaultsStore } from "../../store/defaults.ts";
import {
  OTUS_CHAT_MODES,
  OTUS_CHAT_MODE_INHERIT,
  OTUS_MODE_KEY,
  otusModeLabel,
} from "../../otus.ts";

const copyIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CopyOutline")
);
const refreshIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/RefreshOutline")
);
const botIcon = defineAsyncComponent(() => import("@vicons/fa/Robot"));
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/TrashBinOutline")
);
const listIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ListOutline")
);
const sparklesIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SparklesOutline")
);
// Same purple as the copilot panel and the trace viewer.
const COPILOT_COLOR = "#7c5cff";
const bugIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/BugOutline")
);
const bulbIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/BulbOutline")
);
const consoleIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/TerminalOutline")
);

interface ChatActionsProps {
  chat: Chat;
}

const appStore = useAppStore();
const store = useCcStore();
const usersStore = useUsersStore();
const defaultsStore = useDefaultsStore();
const notification = useNotification();

const { currentChat } = storeToRefs(store);
const { users } = storeToRefs(usersStore);

const props = defineProps<ChatActionsProps>();
const router = useRouter();

function copyLink() {
  const url = new URL(appStore.conf?.params.fullUrl);

  url.searchParams.set("chat", props.chat.uuid);
  addToClipboard(url.href, notification);
}

async function refresh() {
  try {
    await store.get_messages(props.chat, false);
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
}

async function deleteChat() {
  try {
    await store.delete_chat(props.chat);
    appStore.displayMode = "half";

    router.push({ name: "Empty Chat" });
    notification.success({ title: "Done" });
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
}

const buttonTitle = ref("");

const isBotSettingsOpen = ref(false);
const isTracesOpen = ref(false);
const isLearnOpen = ref(false);
const isSaveBotStateLoading = ref(false);
const botState = ref<{ [key: string]: any }>({});
const otusMode = ref(OTUS_CHAT_MODE_INHERIT);

// Role.ADMIN is an operator of this chat; its OWNER is the client it is about.
const isChatAdmin = computed(() => currentChat.value?.role === Role.ADMIN);
const otusHint = computed(
  () =>
    OTUS_CHAT_MODES.find((option) => option.value === otusMode.value)?.hint ?? ""
);
// What this chat falls back to. An override means nothing to an operator who
// cannot see what it overrides.
const installOtusModeLabel = computed(() =>
  otusModeLabel(defaultsStore.bot?.values[OTUS_MODE_KEY])
);

const onMessage = ({ data, origin }: any) => {
  if (origin.includes("localhost:8081")) return;
  if (data.type !== "button-title") return;
  buttonTitle.value = data.value;
};

onMounted(() => {
  window.addEventListener("message", onMessage);
});

onUnmounted(() => {
  window.removeEventListener("message", onMessage);
});

const hoverInstancesButton = (event: MouseEvent) => {
  const button = (event.target as HTMLElement).closest(
    ".n-button"
  ) as HTMLButtonElement;

  window.top?.postMessage(
    {
      type: "click-on-button",
      value: { x: button?.offsetLeft, y: button?.offsetTop },
    },
    "*"
  );
};

interface commandType {
  key: string;
  description: string;
}

const commandsButton = ref();

const commands = computed(() => {
  const result: commandType[] = [];

  props.chat.admins.forEach((uuid) => {
    const bot = users.value.get(uuid);

    if (!bot?.ccIsBot) return;
    Object.entries(bot.ccCommands || {}).forEach(([key, description]) => {
      result.push({ key: `/${key}`, description });
    });
  });

  return result;
});

async function sendCommand(content: string) {
  try {
    await store.send_message(
      new Message({
        content,
        kind: Kind.FOR_BOT,
        chat: router.currentRoute.value.params.uuid as string,
      })
    );
    commandsButton.value.onClick();
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
}

function setBotState() {
  botState.value = (currentChat.value?.toJson() as any as Chat).botState;
  botState.value.enabled = !botState.value.disabled;
  botState.value.review = !botState.value.skip_review;
  otusMode.value = botState.value[OTUS_MODE_KEY] || OTUS_CHAT_MODE_INHERIT;
}

async function saveBotState() {
  try {
    isSaveBotStateLoading.value = true;

    // The mode is sent only when it actually moved. The server refuses the key
    // to anyone below Role.ADMIN, and a client muting their own bot must not be
    // refused over a value they never touched. The rest of bot_state - the
    // escalated flag the bot writes on a handoff - survives either way: the
    // server merges what arrives into what the chat already holds.
    const stored = botState.value[OTUS_MODE_KEY] || OTUS_CHAT_MODE_INHERIT;
    const modeChanged = otusMode.value !== stored;

    await store.update_bot_state(
      SetBotStateRequest.fromJson({
        chat: currentChat.value!.uuid,
        disabled: !botState.value.enabled,
        skipReview: !botState.value.review,
        ...(modeChanged ? { state: { [OTUS_MODE_KEY]: otusMode.value } } : {}),
      })
    );

    isBotSettingsOpen.value = false;
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  } finally {
    isSaveBotStateLoading.value = false;
  }
}

setBotState();

watch(currentChat, () => {
  setBotState();
});
</script>

<style scoped>
.bot_state_otus {
  align-items: flex-start;
}

.bot_state_settings_field {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 20px 0px;
}

.bot_state_settings_actions {
  display: flex;
  justify-content: end;
  align-items: center;
  margin-top: 30px;
}
</style>
