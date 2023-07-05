<template>
  <n-layout-sider>
    <n-scrollbar style="height: 90vh">
      <n-list hoverable clickable>
        <template #header>
          <n-space justify="center" align="center" style="height: 5vh">
            <n-button ghost type="success" @click="router.push({ name: 'Start Chat' })">
              <template #icon>
                <n-icon :component="ChatbubbleEllipsesOutline"/>
              </template>
              Start Chat
            </n-button>
          </n-space>
        </template>
        <chat-item v-for="chat in chats" :uuid="chat.uuid" :chat="chat"/>
      </n-list>
    </n-scrollbar>
  </n-layout-sider>
  <n-layout-content>
    <router-view/>
  </n-layout-content>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent} from 'vue';
import {NButton, NIcon, NLayoutContent, NLayoutSider, NList, NScrollbar, NSpace,} from 'naive-ui';

import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chatItem.vue";

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));

const store = useCcStore();
const router = useRouter();

async function sync() {
  await store.list_chats();
}

sync()

const chats = computed(() => {
  let res: Chat[] = []
  store.chats.forEach((chat) => {
    res.push(chat)
  })

  let sortable = (chat: Chat) => {
    if (chat.meta && chat.meta.lastMessage)
      return chat.meta.lastMessage.sent
    return chat.created
  }

  return res.sort((a: Chat, b: Chat) => {
    return Number(sortable(b) - sortable(a))
  })
})
</script>