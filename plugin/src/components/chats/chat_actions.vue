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

  <n-modal v-model:show="isBotSettingsOpen">
    <n-card
      title="Bot settings"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
      style="width: 500px; min-height: 300px"
    >
      <div class="bot_state_settings_field">
        <span>Bot state:</span>
        <n-switch v-model:value="botState.disabled">
          <template #checked> disabled (no bot)</template>
          <template #unchecked> active </template>
        </n-switch>
      </div>

      <div class="bot_state_settings_field">
        <span>Review:</span>
        <n-switch v-model:value="botState.skip_review">
          <template #checked> open mode (skip preview) </template>
          <template #unchecked> review mode </template>
        </n-switch>
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
    </n-card>
  </n-modal>

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
import { computed, defineAsyncComponent, ref, watch } from "vue";
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
} from "naive-ui";

import { ConnectError } from "@connectrpc/connect";
import {
  Chat,
  Kind,
  Message,
  SetBotStateRequest,
} from "../../connect/cc/cc_pb";
import { useAppStore } from "../../store/app.ts";
import { useCcStore } from "../../store/chatting.ts";
import { addToClipboard } from "../../functions.ts";
import { storeToRefs } from "pinia";

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
const consoleIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/TerminalOutline")
);

interface ChatActionsProps {
  chat: Chat;
}

const appStore = useAppStore();
const store = useCcStore();
const notification = useNotification();

const { currentChat } = storeToRefs(store);

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
const isSaveBotStateLoading = ref(false);
const botState = ref<{ [key: string]: any }>({});

window.addEventListener("message", ({ data, origin }) => {
  if (origin.includes("localhost:8081")) return;
  if (data.type !== "button-title") return;
  buttonTitle.value = data.value;
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
    const bot = store.users.get(uuid);

    if (!bot?.ccIsBot) return;
    Object.entries(bot.ccCommands).forEach(([key, description]) => {
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
}

async function saveBotState() {
  try {
    isSaveBotStateLoading.value = true;

    await store.update_bot_state(
      SetBotStateRequest.fromJson({
        chat: currentChat.value!.uuid,
        disabled: !!botState.value.disabled,
        skipReview: !!botState.value.skip_review,
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
