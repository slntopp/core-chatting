<template>
  <!-- Two conversations over one core-chatting chat: the customer thread and
       the operator's own thread with the bot. While the copilot is closed the
       second pane is zero-width and the drag handle is gone, so the ordinary
       chat looks and behaves exactly as before. -->
  <n-split
    class="split"
    direction="horizontal"
    :size="appStore.isCopilotOpen ? appStore.copilotSplit : 1"
    :min="0.3"
    :max="0.8"
    :disabled="!appStore.isCopilotOpen"
    :resize-trigger-size="appStore.isCopilotOpen ? 3 : 0"
    @update:size="(v) => (appStore.copilotSplit = Number(v))"
  >
    <template #1>
      <n-list class="chat">
    <template #header>
      <chat-header :chat="chat" />
    </template>

    <n-scrollbar
      ref="scrollbar"
      class="chat__scroll"
      v-if="isMessageLoading || messages.length > 0"
    >
      <template v-if="isMessageLoading">
        <mock-message
          v-for="(_, index) in 5"
          :key="index + (chat?.topic || '')"
        />
      </template>

      <n-list-item v-else v-for="message in messages" :key="message.uuid">
        <message-view
          :message="message"
          @approve="(a) => handle_approve(message, a)"
          @convert="
            (kind) => {
              store.updating = true;
              store.current_message = message;
              store.handle_send(chat?.uuid, kind);
            }
          "
          @delete="handle_delete(message)"
          @edit="handle_edit(message)"
        />
      </n-list-item>
    </n-scrollbar>

    <n-space
      v-else
      class="chat__scroll"
      style="width: 100%"
      justify="center"
      align="center"
    >
      <n-alert type="info" title="No Messages yet">
        Use textarea below to send a message.
      </n-alert>
    </n-space>

    <template #footer>
      <chat-footer ref="footer" :chat="chat!" :messages="messages" />
    </template>
      </n-list>
    </template>

    <template #2>
      <!-- Side by side only where there is room for it. On a narrow layout the
           split would leave both panes unusable, so the copilot becomes a
           full-width drawer instead (below) and this pane stays empty. -->
      <copilot-panel
        v-if="appStore.isCopilotOpen && chat && !appStore.isMobile"
        :chat="chat as Chat"
      />
    </template>
  </n-split>

  <!-- Narrow layout: same panel, shown over the chat. Without this the admin
       lane is unreachable below 900px - the notes exist but nothing renders or
       composes them. -->
  <n-drawer
    v-if="appStore.isMobile"
    v-model:show="appStore.isCopilotOpen"
    placement="right"
    width="100%"
  >
    <n-drawer-content :native-scrollbar="false" body-content-style="padding: 0">
      <copilot-panel v-if="chat" :chat="chat as Chat" />
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onMounted,
  ref,
  watch,
} from "vue";
import {
  NAlert,
  NDrawer,
  NDrawerContent,
  NList,
  NListItem,
  NScrollbar,
  NSpace,
  NSplit,
} from "naive-ui";
import { useRoute, useRouter } from "vue-router";

import { useAppStore } from "../../../store/app";
import { useCcStore } from "../../../store/chatting";
import { Chat, Kind, Message } from "../../../connect/cc/cc_pb";

import ChatHeader from "../../../components/chats/layouts/chat_header.vue";
import ChatFooter from "../../../components/chats/layouts/chat_footer.vue";
import MockMessage from "../../../components/chats/mock_message.vue";
import CopilotPanel from "../../../components/chats/copilot_panel.vue";

const MessageView = defineAsyncComponent(
  () => import("../../../components/chats/message.vue"),
);

const route = useRoute();
const router = useRouter();

const appStore = useAppStore();
const store = useCcStore();
const scrollbar = ref();
const isMessageLoading = ref(false);

const chatUuid = computed(() => route.params.uuid);
const chat = computed(() => (store.currentChat as Chat) ?? null);

// The customer thread only. Admin notes are the copilot conversation and are
// rendered in the side pane instead, so the two never interleave.
const messages = computed(() => {
  const chatMessages = store
    .chat_messages(chat.value! as Chat)
    .filter((m) => m.kind !== Kind.ADMIN_ONLY);

  chatMessages.sort((a, b) => Number(a.sent - b.sent));
  return chatMessages;
});

async function handle_approve(msg: Message, approve: boolean) {
  msg.underReview = !approve;
  await store.update_message(msg);
}

async function handle_delete(msg: Message) {
  await store.delete_message(msg);
}

async function scrollToBottom(smooth = false) {
  await nextTick();

  setTimeout(() => {
    if (!scrollbar.value) {
      console.warn("scrollbar not ready");
      return;
    }

    const top = scrollbar.value.$el.nextSibling.firstChild.scrollHeight;

    scrollbar.value.scrollTo({ top, behavior: smooth ? "smooth" : "instant" });
  }, 300);
}

async function load_chat() {
  if (!chat.value) return;
  try {
    isMessageLoading.value = true;
    await Promise.all([store.get_messages(chat.value as Chat)]);
  } finally {
    isMessageLoading.value = false;
  }
}

watch(chat, load_chat);
// Unread counts every inbound message, admin notes included, but `messages` is
// now the customer thread alone - so watching it left the open chat showing
// unread badges that only a customer message could ever clear. Watch the full
// thread for the reset, and keep scrolling tied to what this pane renders.
watch(
  () => (chat.value ? store.chat_messages(chat.value as Chat).length : 0),
  () => {
    if (chat.value?.meta) {
      chat.value.meta.unread = 0;
    }
  },
  { immediate: true },
);
watch(
  messages,
  () => {
    if (scrollbar.value) scrollToBottom();
  },
  { deep: true },
);

watch(messages, () => scrollToBottom());
watch(scrollbar, scrollToBottom);

const footer = ref();

function handle_edit(message: Message) {
  store.updating = true;
  store.current_message = JSON.parse(JSON.stringify(message));

  if (message.underReview) {
    footer.value.sendMode = "approve";
  }
}

const chatPaddingLeft = computed(() => (appStore.isMobile ? "12px" : "16px"));

const fetchChat = async () => {
  try {
    await store.get_chat(route.params.uuid as string);
    if (chat.value == null) {
      throw new Error("Chat not found");
    }
    load_chat();
  } catch (e) {
    console.warn(e);
    router.push("/dashboard");
  }
};

function sendUser() {
  let userUuid = chat.value.owner;

  const isAdmin = chat.value.admins.includes(store.me.uuid);
  if (isAdmin) {
    userUuid =
      chat.value.users.find((u) => !chat.value.admins.includes(u)) ??
      chat.value.owner;
  }

  window.top?.postMessage(
    {
      type: "send-user",
      value: { uuid: userUuid },
    },
    "*",
  );
}

watch(
  () => chat.value?.owner,
  () => {
    sendUser();
  },
);

onMounted(fetchChat);
watch(chatUuid, fetchChat);
</script>

<style scoped>
/* The panes are 100dvh tall; pinning the container stops a dragged split from
   pushing the page into a horizontal scroll. */
.split {
  height: 100dvh;
  overflow: hidden;
}

.chat {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  padding-left: v-bind("chatPaddingLeft");
}

/* min-height:0 is what actually lets a flex child scroll instead of growing
   the column past the viewport. */
.chat__scroll {
  flex: 1 1 auto;
  min-height: 0;
}
</style>
