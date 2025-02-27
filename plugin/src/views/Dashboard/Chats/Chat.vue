<template>
  <n-list class="chat">
    <template #header>
      <chat-header :chat="chat" />
    </template>

    <n-scrollbar
      ref="scrollbar"
      style="max-height: calc(100dvh - 66px - 77px); margin-bottom: auto"
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
      style="height: 90vh; width: 100%"
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

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onMounted,
  ref,
  watch,
} from "vue";
import { NAlert, NList, NListItem, NScrollbar, NSpace } from "naive-ui";
import { useRoute, useRouter } from "vue-router";

import { useAppStore } from "../../../store/app";
import { useCcStore } from "../../../store/chatting";
import { Chat, Kind, Message } from "../../../connect/cc/cc_pb";

import ChatHeader from "../../../components/chats/layouts/chat_header.vue";
import ChatFooter from "../../../components/chats/layouts/chat_footer.vue";
import MockMessage from "../../../components/chats/mock_message.vue";

const MessageView = defineAsyncComponent(
  () => import("../../../components/chats/message.vue")
);

const route = useRoute();
const router = useRouter();

const appStore = useAppStore();
const store = useCcStore();
const scrollbar = ref();
const isMessageLoading = ref(false);

const chatUuid = computed(() => route.params.uuid);
const chat = computed(() => (store.currentChat as Chat) ?? null);

const messages = computed(() => {
  const chatMessages = store.chat_messages(chat.value! as Chat);

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
watch(
  messages,
  () => {
    if (chat.value?.meta) {
      chat.value.meta.unread = 0;
    }

    if (scrollbar.value) scrollToBottom();
  },
  { deep: true }
);

watch(messages, () => scrollToBottom());
watch(scrollbar, scrollToBottom);

const footer = ref();

function handle_edit(message: Message) {
  store.updating = true;
  store.current_message = JSON.parse(JSON.stringify(message));

  if (message.underReview) {
    footer.value.sendMode = "approve";
  } else if (message.kind === Kind.ADMIN_ONLY) {
    footer.value.sendMode = "admin";
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

onMounted(fetchChat);
watch(chatUuid, fetchChat);
</script>

<style scoped>
.chat {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  padding-left: v-bind("chatPaddingLeft");
}
</style>
