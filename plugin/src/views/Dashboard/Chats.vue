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

        <n-popover
          scrollable
          trigger="click"
          placement="bottom"
          style="max-height: 50vh"
        >
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
                <n-radio value="status" label="Status" />
              </n-space>
            </n-radio-group>
          </div>

          <div style="margin-top: 20px">
            <n-text>Status priority:</n-text>
            <n-divider style="margin: 5px 0" />
            <n-space :wrap-item="false" style="gap: 7px">
              <div v-for="(status, i) of statusSorters" style="display: flex; align-items: center">
                <n-icon v-if="i !== 0" @click="() => {
                  const i = statusSorters.indexOf(status)

                  if (!statusSorters[i - 1]) return
                  [statusSorters[i], statusSorters[i- 1]] = [statusSorters[i- 1], statusSorters[i]]
                }">
                  <left-icon />
                </n-icon>

                <n-tag round style="margin: 0 3px">
                  {{ statuses[status].label }}
                </n-tag>

                <n-icon v-if="i !== statusSorters.length - 1" @click="() => {
                  const i = statusSorters.indexOf(status)
                  
                  if (!statusSorters[i + 1]) return
                  [statusSorters[i], statusSorters[i + 1]] = [statusSorters[i + 1], statusSorters[i]]
                }">
                  <right-icon />
                </n-icon>
              </div>
            </n-space>
          </div>

          <chats-filters
            v-model:meta="meta"
            v-model:checked-statuses="checkedStatuses"
            v-model:checked-admins="checkedAdmins"
            v-model:metrics-options="metricsOptions"
            v-model:meta-options="metaOptions"
            :metrics="(metrics as MetricWithKey[])"
          />
        </n-popover>
      </n-space>

      <n-scrollbar ref="scrollbar" style="height: calc(100vh - 100px); min-width: 150px">
        <n-popover trigger="manual" :show="isFirstMessageVisible" :x="x" :y="y">
          {{ firstMessage }}
        </n-popover>

        <n-list hoverable clickable style="margin-bottom: 25px">
          <chat-item
            v-for="chat in viewedChats"
            :id="chat.uuid"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat"
            :class="{ active: chat.uuid === router.currentRoute.value.params.uuid }"
            @click="changeMode('none')"
            @hover="onMouseMove"
            @hover-end="isFirstMessageVisible = false"
          />
        </n-list>
        <n-spin ref="loading" style="display: flex; margin: 0 auto 20px" />
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
  NPopover, NRadioGroup, NRadio, NDivider, NText, NSpin, NTag
} from 'naive-ui';

import {useAppStore} from '../../store/app.ts';
import {useCcStore} from '../../store/chatting.ts';

import {useRouter} from 'vue-router';
import {Chat, Status} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";
import ChatsFilters from "../../components/chats/chats_filters.vue";
import useDraggable from "../../hooks/useDraggable.ts";
import useDefaults, { MetricWithKey } from '../../hooks/useDefaults.ts';

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const CloseIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));
const SortIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SettingsOutline'));
const SwitchIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SwapHorizontal'));
const rightIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const leftIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));

const appStore = useAppStore()
const store = useCcStore();
const router = useRouter();
const {makeDraggable} = useDraggable()
const {metrics, fetch_defaults} = useDefaults()

const searchParam = ref('')
const scrollbar = ref<InstanceType<typeof NScrollbar>>()
const loading = ref<InstanceType<typeof NSpin>>()
const page = ref(1)

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
    first: document.querySelector(".chats__panel")!,
    second: document.querySelector(".chat__item")!
  })

  const options = {
    root: scrollbar.value?.$el.nextElementSibling,
    threshold: 1
  }

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(async ({ isIntersecting }) => {
      await new Promise((resolve) => setTimeout(resolve, 200))
      if (isIntersecting) page.value += 1
    })
  }, options)

  observer.observe(loading.value?.$el)
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

const sortBy = ref<'sent' | 'created' | 'status'>('sent')

interface optionsIncludedType {
  [index: string]: boolean
}

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

const statusSorters = ref(statuses.value.map(({ value }) => value))

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
    const options = { ...metaOptions.value, ...metricsOptions.value }
    const isOptionsIncluded: optionsIncludedType = {}

    Object.entries(options).forEach(([key, value]) => {
      if (value.length < 1) {
        isOptionsIncluded[key] = true
        return
      }
      const metric = chat.meta?.data[key]?.toJson() as never

      if (value.some(({ type }) => type === 'date')) {
        const [start, end] = value as { value: number }[]
        const date = (metric as { value: number }).value * 1000

        isOptionsIncluded[key] = start.value <= date && date <= end.value
      } else {
        isOptionsIncluded[key] = value.includes(metric)
      }
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

  const sortable = (chat: Chat) => {
    if (sortBy.value === 'status') {
      return statusSorters.value.indexOf(chat.status)
    }
    if (chat.meta?.lastMessage && sortBy.value === 'sent') {
      return Number(chat.meta.lastMessage.sent)
    }
    return Number(chat.created)
  }

  return result.sort((a: Chat, b: Chat) =>
    sortable(b) - sortable(a)
  )
})

const viewedChats = computed(() =>
  chats.value.slice(0, 10 * page.value)
)

if (Object.keys(metrics.value).length < 1) fetch_defaults()

interface metricsOptionsType {
  [index: string]: []
}

const checkedStatuses = ref<Status[]>([])
const checkedAdmins = ref<string[]>([])

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

const meta = ref({})
const metaOptions = ref<metricsOptionsType>(
  Object.keys(meta.value).reduce((result, key) => ({
    ...result, [key]: []
  }), {})
)

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
  background-color: var(--n-color);
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
    padding: 10px;

    &.hide {
      flex-flow: column-reverse !important;
      margin: 0px;
      padding: 10px;
    }
  }

  .search {
    margin: 5px 9px 10px;

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