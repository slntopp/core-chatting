<template>
  <div
    v-if="appStore.displayMode !== 'none'"
    :class="{ chats__panel: true, closed: !isChatPanelOpen }"
    :style="(appStore.displayMode === 'full') ? 'min-width: calc(100% - 8px)' : undefined"
  >
    <n-layout-sider
      style="margin-bottom: 25px"
      collapse-mode="transform"
      :collapsed="!isChatPanelOpen"
    >
      <n-space
        align="center"
        :wrap-item="false"
        :justify="(isChatPanelOpen) ? 'space-between' : 'center'"
        :class="{ chat__actions: true, hide: !isChatPanelOpen }"
      >
        <n-button ghost @click="changeMode(null)">
          <n-icon> <switch-icon /> </n-icon>
        </n-button>

        <n-button
          ghost
          type="success"
          :style="(isChatPanelOpen) ? 'margin-right: auto' : null"
          @click="router.push({ name: 'Start Chat' })"
        >
          <n-icon :component="ChatbubbleEllipsesOutline" />
          <span v-if="isChatPanelOpen" style="margin-left: 5px">
            Start Chat
          </span>
        </n-button>

        <n-button ghost @click="changePanelOpen">
          <n-icon>
            <close-icon v-if="isChatPanelOpen" />
            <open-icon v-else />
          </n-icon>
        </n-button>
      </n-space>

      <n-space
        class="search"
        align="center"
        justify="space-between"
        v-if="isChatPanelOpen"
        :wrap="false"
      >
        <n-input v-model:value="searchParam" type="text" placeholder="Search..."/>

        <n-popover trigger="click" placement="bottom">
          <template #trigger>
            <n-icon size="24" style="vertical-align: middle; cursor: pointer">
              <sort-icon />
            </n-icon>
          </template>

          <div>
            <div>Sort By:</div>
            <n-divider style="margin: 5px 0" />
            <n-radio-group v-model:value="sortBy">
              <n-space>
                <n-radio value="created" label="Created" />
                <n-radio value="sent" label="Sent" />
              </n-space>
            </n-radio-group>
          </div>

          <div style="margin-top: 10px">
            <div>Filter by status:</div>
            <n-divider style="margin: 5px 0" />
            <n-checkbox-group v-model:value="checkedStatuses">
              <n-space>
                <n-checkbox v-for="status of statuses" :value="status.value" :label="status.label" />
              </n-space>
            </n-checkbox-group>
          </div>
        </n-popover>
      </n-space>

      <n-scrollbar style="height: calc(100vh - 100px); min-width: 150px">
        <n-list hoverable clickable style="margin-bottom: 25px">
          <chat-item
            v-for="chat in chats"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat"
            :class="{ active: chat.uuid === router.currentRoute.value.params.uuid }"
            @click="changeMode('none')"
          />
        </n-list>
      </n-scrollbar>
    </n-layout-sider>
  </div>

  <div id="separator" v-show="isChatPanelOpen"></div>
  <div class="chat__item" v-if="appStore.displayMode !== 'full'">
    <n-layout-content>
      <router-view/>
    </n-layout-content>
  </div>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, onMounted, ref} from 'vue';
import {NButton, NIcon, NInput, NLayoutContent, NLayoutSider, NList, NScrollbar, NSpace, NPopover, NRadioGroup, NRadio, NDivider, NCheckboxGroup, NCheckbox} from 'naive-ui';

import {useAppStore} from '../../store/app.ts';
import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat, Status} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";
import useDraggable from "../../hooks/useDraggable.ts";

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const CloseIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));
const SortIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SettingsOutline'));
const SwitchIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SwapHorizontal'));

const appStore = useAppStore()
const store = useCcStore();
const router = useRouter();
const {makeDraggable} = useDraggable()

const searchParam = ref('')

async function sync() {
  await store.list_chats();
}

onMounted(() => {
  makeDraggable({
    resizer: document.getElementById("separator")!,
    first: document.getElementsByClassName("chats__panel").item(0) as HTMLElement,
    second: document.getElementsByClassName("chat__item").item(0) as HTMLElement
  })
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

const sortBy = ref<'sent' | 'created'>('sent')

const chats = computed(() => {
  let result: Chat[] = []
  store.chats.forEach((chat) => {
    result.push(chat)
  })

  result = result.filter((chat) => {
    const isIncluded = checkedStatuses.value.includes(chat.status)

    if (checkedStatuses.value.length > 0) {
      return filterChat(chat, searchParam.value) && isIncluded
    } else {
      return filterChat(chat, searchParam.value)
    }
  })

  let sortable = (chat: Chat) => {
    if (chat.meta?.lastMessage && sortBy.value === 'sent') {
      return chat.meta.lastMessage.sent
    }
    return chat.created
  }

  return result.sort((a: Chat, b: Chat) => {
    return Number(sortable(b) - sortable(a))
  })
})

const checkedStatuses = ref<Status[]>([])

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

const isChatPanelOpen = ref(true)

function changePanelOpen() {
  isChatPanelOpen.value = !isChatPanelOpen.value

  if (isChatPanelOpen) appStore.displayMode = 'half'
  else appStore.displayMode = 'none'
}

function changeMode(mode: string | null) {
  if (mode) appStore.displayMode = mode
  else if (appStore.displayMode === 'half') {
    appStore.displayMode = 'full'
    isChatPanelOpen.value = true
  } else {
    appStore.displayMode = 'half'
  }
}
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
  min-width: 450px;
  width: 450px;
  transition: .4s;

  .chat__actions {
    padding: 10px 0px 10px 10px;

    &.hide {
      flex-flow: column-reverse !important;
      margin: 0px;
      padding: 10px;
    }
  }

  .search {
    margin-top: 5px;
    margin-bottom: 10px;
    margin-left: 9px;

    div:first-child {
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

  .n-list-item.active {
    background: var(--n-color-hover-modal);
  }
}

.chat__item {
  height: 100%;
  min-width: 50%;
  width: 100%;
}
</style>