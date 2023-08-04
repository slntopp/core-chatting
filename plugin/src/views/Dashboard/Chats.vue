<template>
  <div :class="{'chats__panel':true,'closed':!isChatPanelOpen}">
    <n-layout-sider style="margin-bottom: 25px" collapse-mode="transform" :collapsed="!isChatPanelOpen">
      <n-space :justify="isChatPanelOpen?'space-between':'center'" align="center"
               :class="{'chat__actions':true,hide:!isChatPanelOpen}">
        <n-button ghost type="success" @click="router.push({ name: 'Start Chat' })">
          <n-icon :component="ChatbubbleEllipsesOutline"/>
          <span v-if="isChatPanelOpen" style="margin-left: 5px">Start Chat</span>
        </n-button>
        <n-button ghost @click="isChatPanelOpen=!isChatPanelOpen">
          <n-icon>
            <close-icon v-if="isChatPanelOpen"/>
            <open-icon v-else/>
          </n-icon>
        </n-button>
      </n-space>
      <n-space class="search" v-if="isChatPanelOpen" align="center" justify="center">
        <n-input v-model:value="searchParam" type="text" placeholder="Search..."/>
      </n-space>
      <n-scrollbar style="height: 100vh;min-width: 150px">
        <n-list style="margin-bottom: 25px" hoverable clickable>
          <chat-item :hide-message="!isChatPanelOpen" v-for="chat in chats" :uuid="chat.uuid" :chat="chat"/>
        </n-list>
      </n-scrollbar>
    </n-layout-sider>
  </div>
  <div id="separator" v-show="isChatPanelOpen"></div>
  <div class="chat__item">
    <n-layout-content>
      <router-view/>
    </n-layout-content>
  </div>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, onMounted, ref} from 'vue';
import {NButton, NIcon, NInput, NLayoutContent, NLayoutSider, NList, NScrollbar, NSpace} from 'naive-ui';

import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";
import useDraggable from "../../hooks/useDraggable.ts";

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const CloseIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));

const store = useCcStore();
const router = useRouter();
const {makeDraggable} = useDraggable()

const isChatPanelOpen = ref(true)
const searchParam = ref('')

async function sync() {
  await store.list_chats();
}

onMounted(() => {
  makeDraggable({
    resizer: document.getElementById("separator")!,
    first: document.getElementsByClassName("chats__panel").item(0) as HTMLElement,
    second: document.getElementsByClassName("chat__item").item(0) as HTMLElement
  });
})

sync()

const filterChat = (chat: Chat, val: string): boolean => {
  if (!val) {
    return true
  }
  val = val.toLowerCase()

  const startsWithKeys = ['topic', 'uuid']

  for (const key of startsWithKeys) {
    if ((chat as any)[key]?.toLowerCase().startsWith(val)) {
      return true
    }
  }

  return !!chat.users.find(u => u.startsWith(val) || store.users.get(u)?.title.toLowerCase().startsWith(val)) || !!chat.admins.find(u => u.startsWith(val) || store.users.get(u)?.title.toLocaleLowerCase().startsWith(val))
}

const chats = computed(() => {
  let res: Chat[] = []
  store.chats.forEach((chat) => {
    res.push(chat)
  })

  res = res.filter((c) => filterChat(c, searchParam.value))

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

<style lang="scss">
#separator {
  cursor: col-resize;
  background-color: #18181C;
  background-repeat: no-repeat;
  background-position: center;
  width: 10px;
  /* Prevent the browser's built-in drag from interfering */
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.chats__panel {
  min-width: 275px;
  width: 275px;

  .chat__actions {
    padding: 10px 0px 10px 10px;

    &.hide {
      flex-flow: column-reverse !important;
      margin: 0px;
    }
  }

  .search {
    margin-top: 5px;
    margin-bottom: 10px;
    margin-left: 9px;

    div {
      width: 100%;
    }
  }

  &.closed {
    min-width: 70px !important;
    width: 70px !important;
    background-color: #18181C;

    .n-layout-sider-scroll-container {
      min-width: 70px !important;
      width: 70px !important;
    }
  }

  aside {
    width: 100% !important;
    max-width: 100% !important;
  }
}

.chat__item {
  height: 100%;
  min-width: 50%;
  width: 100%;
}
</style>