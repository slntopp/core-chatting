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
          @click="startChat"
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

        <n-popover trigger="click" placement="bottom" :width="520">
          <template #trigger>
            <n-icon size="24" style="vertical-align: middle; cursor: pointer">
              <sort-icon />
            </n-icon>
          </template>

          <div>
            <n-text>Sort By:</n-text>
            <n-divider style="margin: 5px 0" />
            <n-radio-group v-model:value="sortBy">
              <n-space :wrap-item="false">
                <n-radio value="created" label="Created" />
                <n-radio value="sent" label="Sent" />
              </n-space>
            </n-radio-group>
          </div>

          <div style="margin-top: 20px">
            <n-text>Filter by status:</n-text>
            <n-divider style="margin: 5px 0" />
            <n-checkbox-group v-model:value="checkedStatuses">
              <n-space :wrap-item="false">
                <n-checkbox v-for="status of statuses" :value="status.value" :label="status.label" />
              </n-space>
            </n-checkbox-group>
          </div>

          <div style="margin-top: 20px">
            <n-text>Filter by admin:</n-text>
            <n-divider style="margin: 5px 0" />
            <n-checkbox-group v-model:value="checkedAdmins">
              <n-space :wrap-item="false">
                <n-checkbox v-for="admin of admins" :value="admin.value" :label="admin.label" />
              </n-space>
            </n-checkbox-group>
          </div>

          <div style="margin-top: 20px" v-for="metric of metrics" :key="metric.key">
            <n-text>Filter by {{ metric.title.toLowerCase() }}:</n-text>
            <n-divider style="margin: 5px 0" />
            <n-checkbox-group v-model:value="metricsOptions[metric.key]">
              <n-space :wrap-item="false">
                <n-checkbox v-for="option of metric.options" :value="option.value" :label="option.key" />
              </n-space>
            </n-checkbox-group>
          </div>
        </n-popover>
      </n-space>

      <n-scrollbar style="height: calc(100vh - 100px); min-width: 150px">
        <n-popover trigger="manual" :show="isFirstMessageVisible" :x="x" :y="y">
          {{ firstMessage }}
        </n-popover>

        <n-list hoverable clickable style="margin-bottom: 25px">
          <chat-item
            v-for="chat in chats"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat"
            :class="{ active: chat.uuid === router.currentRoute.value.params.uuid }"
            @click="changeMode('none')"
            @hover="onMouseMove"
            @hoverEnd="isFirstMessageVisible = false"
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
import {computed, defineAsyncComponent, onMounted, ref, watch} from 'vue';
import {
  NButton, NIcon, NInput, NLayoutContent,
  NLayoutSider, NList, NScrollbar, NSpace,
  NPopover, NRadioGroup, NRadio, NDivider,
  NCheckboxGroup, NCheckbox, NText
} from 'naive-ui';

import {useAppStore} from '../../store/app.ts';
import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat, Status} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";
import useDraggable from "../../hooks/useDraggable.ts";
import useDefaults from '../../hooks/useDefaults.ts';

const emits = defineEmits(['hover', 'hoverEnd'])

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const CloseIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));
const SortIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SettingsOutline'));
const SwitchIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SwapHorizontal'));

const appStore = useAppStore()
const store = useCcStore();
const router = useRouter();
const {makeDraggable} = useDraggable()
const {metrics, fetch_defaults} = useDefaults()

const searchParam = ref('')

async function sync() {
  try {
    await store.list_chats()
    const { users } = await store.get_members()

    users.forEach((user) => {
      store.users.set(user.uuid, user)
    })
  } catch (error) {
    console.log(error);
  }
}

function startChat() {
  router.push({ name: 'Start Chat' })
  appStore.displayMode = 'none'
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

interface optionsIncludedType {
  [index: string]: boolean
}

const chats = computed(() => {
  let result: Chat[] = []
  store.chats.forEach((chat) => {
    result.push(chat)
  })

  result = result.filter((chat) => {
    let isIncluded = checkedStatuses.value.includes(chat.status)
    let isAdminsExist = !!checkedAdmins.value.find((uuid) =>
      chat.admins.includes(uuid)
    )
    let isOptionsIncluded: optionsIncludedType = {}

    Object.entries(metricsOptions.value).forEach(([key, value]) => {
      if (value.length < 1) {
        isOptionsIncluded[key] = true
        return
      }
      const metric = chat.meta?.data[key]?.kind.value as never

      isOptionsIncluded[key] = value.includes(metric)
    })
    
    if (checkedStatuses.value.length < 1) {
      isIncluded = true
    }
    if (checkedAdmins.value.length < 1) {
      isAdminsExist = true
    }

    return filterChat(chat, searchParam.value) &&
      isIncluded && isAdminsExist &&
      Object.values(isOptionsIncluded).every((value) => value)
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

interface adminsType {
  value: string
  label: string
}

const checkedAdmins = ref<string[]>([])

const admins = computed(() => {
  const result: adminsType[] = []

  store.chats.forEach((chat: Chat) => {
    chat.admins.forEach((uuid) => {
      const element = result.find(({ value }) => uuid === value)

      if (element) return
      else result.push({
        value: uuid, label: store.users.get(uuid)?.title ?? uuid
      })
    })
  })

  return result
})

if (Object.keys(metrics.value).length < 1) fetch_defaults()

interface metricsOptionsType {
  [index: string]: []
}

const metricsOptions = ref<metricsOptionsType>(
  Object.keys(metrics.value).reduce((result, key) => ({
    ...result, [key]: []
  }), {})
)

watch(metrics, (value) => {
  Object.keys(value).reduce((result, key) => ({
    ...result, [key]: []
  }), {})
})

const isChatPanelOpen = ref(true)

function changePanelOpen() {
  isChatPanelOpen.value = !isChatPanelOpen.value

  if (isChatPanelOpen) appStore.displayMode = 'half'
  else appStore.displayMode = 'none'
}

function changeMode(mode: string | null) {
  if (mode === 'none' && appStore.displayMode === 'half') {
    return
  }

  if (mode) appStore.displayMode = mode
  else if (appStore.displayMode === 'half') {
    appStore.displayMode = 'full'
    isChatPanelOpen.value = true
  } else {
    appStore.displayMode = 'half'
  }
}

const firstMessage = ref('gg')
const isFirstMessageVisible = ref(false)
const x = ref(0)
const y = ref(0)

function onMouseMove(clientX: number, clientY: number, chatId: string) {
  const chat = chats.value.find(({ uuid }) => uuid === chatId)

  x.value = clientX
  y.value = clientY - 10
  firstMessage.value = chat?.meta?.firstMessage?.content ?? ''
  isFirstMessageVisible.value = true
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