<template>
  <div id="left">
    <n-layout-sider>
      <n-scrollbar style="height: 100vh">
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
  </div>
  <div id="separator"></div>
  <div id="right">
    <n-layout-content>
      <router-view/>
    </n-layout-content>
  </div>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, onMounted} from 'vue';
import {NButton, NIcon, NLayoutContent, NLayoutSider, NList, NScrollbar, NSpace,} from 'naive-ui';

import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));

const store = useCcStore();
const router = useRouter();

async function sync() {
  await store.list_chats();
}

function dragElement(element: any, direction: string = 'H') {
  let md: any;
  const first = document.getElementById("left");
  const second = document.getElementById("right");
  element.onmousedown = onMouseDown;

  function onMouseDown(e: any) {
    md = {
      e,
      offsetLeft: element.offsetLeft,
      offsetTop: element.offsetTop,
      firstWidth: first!.offsetWidth,
      secondWidth: second!.offsetWidth
    };

    document.onmousemove = onMouseMove;
    document.onmouseup = () => {
      document.onmousemove = document.onmouseup = null;
    }
  }

  function onMouseMove(e: any) {
    var delta = {
      x: e.clientX - md.e.clientX,
      y: e.clientY - md.e.clientY
    };

    if (direction === "H") // Horizontal
    {
      delta.x = Math.min(Math.max(delta.x, -md.firstWidth),
          md.secondWidth);

      element.style.left = md.offsetLeft + delta.x + "px";
      first!.style.width = (md.firstWidth + delta.x) + "px";
    }
  }
}

onMounted(() => {
  dragElement(document.getElementById("separator"), "H");
})

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

#left {
  min-width: 275px;

  .n-layout-sider n-layout-sider--static-positioned n-layout-sider--left-placement n-layout-sider--show-content {
    width: 100%;
    max-width: 100%;
  }

  aside {
    width: 100% !important;
    max-width: 100% !important;
  }
}

#right {
  height: 100%;
  min-width: 50%;
  width: 100%;
}
</style>