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

        <n-button ghost type="success" @click="startChat">
          <n-icon :component="ChatbubbleEllipsesOutline" />
          <span v-if="isChatPanelOpen" style="margin-left: 5px">
            Start Chat
          </span>
        </n-button>

        <n-button
          ghost
          type="error"
          v-if="selectedChats.length > 0"
          :loading="deleteLoading"
         
          @click="deleteChats"
        >
          <n-icon :component="deleteIcon" />
          <span v-if="isChatPanelOpen" style="margin-left: 5px">
            Delete Chats
          </span>
        </n-button>

        <n-space :wrap-item="false" :style="(isChatPanelOpen) ? 'margin-right: auto' : null">
          <span v-for="(count, status) in chatsCountByStatus">
          <n-text :style="{ color: getStatusColor(+status) }">
            {{ getStatus(+status) }}:
          </n-text>
          <n-tag round size="small">{{ count }}</n-tag>
        </span>
        </n-space>

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
          style="max-height: 50vh; max-width: 768px"
        >
          <template #trigger>
            <n-icon size="24" style="vertical-align: middle; cursor: pointer">
              <sort-icon />
            </n-icon>
          </template>

          <div>
            <n-space align="center" style="gap: 5px" :wrap-item="false">
              <n-text>Sort By:</n-text>
              <n-icon
                size="18"
                style="cursor: pointer"
                :style="{ transform: (sortType === 'desc') ? 'rotate(180deg)' : null }"
                @click="(sortType === 'desc') ? sortType = 'asc' : sortType = 'desc'"
              >
                <up-icon />
              </n-icon>
            </n-space>

            <n-divider style="margin: 5px 0" />
            <n-radio-group v-model:value="sortBy">
              <n-space :wrap-item="false">
                <n-radio value="none" label="None" />
                <n-radio value="created" label="Created" />
                <n-radio value="sent" label="Sent" />
                <n-radio value="status" label="Status" />
              </n-space>
            </n-radio-group>
          </div>

          <chats-filters
            v-model:checkedDepartments="checkedDepartments"
            v-model:checkedStatuses="checkedStatuses"
            v-model:checkedAdmins="checkedAdmins"
            :metricsOptions="metricsOptions"
            @update:checked-metrics="(key, value) => metricsOptions[key] = value"
          />

          <n-button
            ghost
            type="primary"
            style="display: block; margin: 10px 0 5px auto"
            @click="resetFilters"
          >
            Reset
          </n-button>
        </n-popover>
      </n-space>

      <n-scrollbar ref="scrollbar" style="height: calc(100vh - 100px); min-width: 150px">
        <n-popover trigger="manual" :show="isFirstMessageVisible" :x="x" :y="y">
          {{ firstMessage }}
        </n-popover>

        <n-list v-if="!(isLoading || isDefaultLoading)" hoverable clickable style="margin-bottom: 25px">
          <chat-item
            v-for="chat in viewedChats"
            :selected="selectedChats"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat"
            :chats="viewedChats"
            :class="{ active: chat.uuid === router.currentRoute.value.params.uuid }"
            @click="changeMode('none')"
            @hover="onMouseMove"
            @hoverEnd="isFirstMessageVisible = false"
            @select="selectChat"
          />
        </n-list>
        <n-spin
          v-if="isLoading || viewedChats.length < chats.length"
          ref="loading"
          style="display: flex; margin: 0 auto 20px"
        />
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
import {computed, defineAsyncComponent, nextTick, onMounted, ref, watch} from 'vue';
import {
  NButton, NIcon, NInput, NLayoutContent,
  NLayoutSider, NList, NScrollbar, NSpace,
  NPopover, NRadioGroup, NRadio, NDivider,
  NTag, NText, NSpin, useNotification
} from 'naive-ui';

import {useAppStore} from '../../store/app.ts';
import {useCcStore} from '../../store/chatting.ts';

import {useRoute, useRouter} from 'vue-router';
import {Chat, Status} from '../../connect/cc/cc_pb';
import ChatItem from "../../components/chats/chat_item.vue";
import ChatsFilters from "../../components/chats/chats_filters.vue"
import useDraggable from "../../hooks/useDraggable.ts";
import useDefaults from '../../hooks/useDefaults.ts';

defineEmits(['hover', 'hoverEnd'])

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowForward'));
const CloseIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));
const SortIcon = defineAsyncComponent(() => import('@vicons/ionicons5/EllipsisVertical'));
const SwitchIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SwapHorizontal'));
const deleteIcon = defineAsyncComponent(() => import('@vicons/ionicons5/CloseCircle'));
const upIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowUp'));

const appStore = useAppStore()
const store = useCcStore()
const router = useRouter()
const route = useRoute()
const {makeDraggable} = useDraggable()
const {metrics, isDefaultLoading, fetch_defaults} = useDefaults()
const notification = useNotification()

const selectedChats = ref<string[]>([])
const deleteLoading = ref(false)
const searchParam = ref('')
const scrollbar = ref<InstanceType<typeof NScrollbar>>()
const loading = ref<InstanceType<typeof NSpin>>()
const page = ref(1)
const isLoading = ref(false)

async function sync() {
  try {
    isLoading.value = true
    await store.list_chats()
    const { users } = await store.get_members()

    users.forEach((user) => {
      store.users.set(user.uuid, user)
    })
  } catch (error) {
    console.log(error);
  } finally {
    isLoading.value = false
  }
}

function startChat() {
  router.push({ name: 'Start Chat' })
  appStore.displayMode = 'none'
}

async function deleteChats () {
  deleteLoading.value = true
  try {
    const current = route.params.uuid as string
    if (selectedChats.value.includes(current)) {
      await router.push({ name: 'Dashboard' })
      changePanelOpen()
    }

    const promises = selectedChats.value.map((uuid) =>
      store.delete_chat(new Chat(store.chats.get(uuid)))
    )

    await Promise.all(promises)
    notification.success({
      title: 'Chats successfully deleted'
    })
  } catch (error: any) {
    notification.error({
      title: error.response?.data.message ?? error.message ?? error
    })
    console.error(error)
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  makeDraggable({
    resizer: document.getElementById("separator")!,
    first: document.querySelector(".chats__panel")!,
    second: document.querySelector(".chat__item")!
  })

  const sorting = JSON.parse(localStorage.getItem('sorting') ?? 'null')
  const filters = JSON.parse(localStorage.getItem('filters') ?? 'null')

  if (sorting) {
    sortBy.value = sorting.sortBy
    sortType.value = sorting.sortType
  }
  if (filters) {
    checkedDepartments.value = filters.departments
    checkedStatuses.value = filters.statuses
    checkedAdmins.value = filters.admins
  }

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

  if (!loading.value?.$el) return
  observer.observe(loading.value?.$el)
})

sync()

function selectChat(uuid: string) {
  if (selectedChats.value.includes(uuid)) {
    const i = selectedChats.value.indexOf(uuid)

    selectedChats.value.splice(i, 1)
  } else {
    selectedChats.value.push(uuid)
  }
}

function filterChat (chat: Chat, val: string): boolean {
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

const sortBy = ref<'none' | 'sent' | 'created' | 'status'>('sent')
const sortType = ref<'asc' | 'desc'>('desc')

interface optionsIncludedType {
  [index: string]: boolean
}

const checkedStatuses = ref<Status[]>([])
const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  })).reverse()
)

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

function getStatusColor (status: Status) {
  switch (status) {
    case 0:
      return '#5084ff'
    case 1:
      return '#1ea01e'
    case 2:
      return '#ff8300'
    case 3:
      return '#e23535'
    default:
      return undefined
  }
}

const chats = computed(() => {
  const result = [...store.chats.values()].filter((chat) => {
    let isDepIncluded = checkedDepartments.value.includes(chat.department)
    let isIncluded = checkedStatuses.value.includes(chat.status)
    let isAdminsExist = !!checkedAdmins.value.find((uuid) =>
      chat.admins.includes(uuid)
    )
    let isOptionsIncluded: optionsIncludedType = {}
    let isAccountOwner = true

    Object.entries(metricsOptions.value).forEach(([key, value]) => {
      if (value.length < 1) {
        isOptionsIncluded[key] = true
        return
      }
      const metric = chat.meta?.data[key]?.kind.value as never

      isOptionsIncluded[key] = value.includes(metric)
    })

    if (checkedDepartments.value.length < 1) {
      isDepIncluded = true
    }
    if (checkedStatuses.value.length < 1) {
      isIncluded = true
    }
    if (checkedAdmins.value.length < 1) {
      isAdminsExist = true
    }

    if (appStore.conf?.params?.filterByAccount) {
      isAccountOwner = [...chat.users, ...chat.admins]
        .includes(appStore.conf.params.filterByAccount)
    }

    return filterChat(chat, searchParam.value) &&
      isDepIncluded && isIncluded && isAdminsExist && isAccountOwner &&
      Object.values(isOptionsIncluded).every((value) => value)
  })

  const sortable = (chat: Chat) => {
    if (sortBy.value === 'status') {
      return statuses.value.map(({ value }) => value).indexOf(chat.status)
    }
    if (chat.meta?.lastMessage && sortBy.value === 'sent') {
      return Number(chat.meta.lastMessage.sent)
    }
    if (sortBy.value === 'none') {
      return ((chat.meta?.unread ?? 0) > 0) ? Date.now() : Number(chat.created)
    }
    return Number(chat.created)
  }

  result.sort((a: Chat, b: Chat) =>
    (sortType.value === 'desc') ? sortable(b) - sortable(a) : sortable(a) - sortable(b)
  )

  return result
})

const viewedChats = computed(() =>
  chats.value.slice(0, 10 * page.value)
)

const filteredChatsByAccount = computed(() =>
  [...store.chats.values()].filter((chat) => {
    let isAccountOwner = true

    if (appStore.conf?.params?.filterByAccount) {
      isAccountOwner = [...chat.users, ...chat.admins]
        .includes(appStore.conf.params.filterByAccount)
    }

    return isAccountOwner
  })
)

const chatsCountByStatus = computed(() => {
  const result: { [key: string]: number } = {}

  filteredChatsByAccount.value.forEach((chat) => {
    if (!result[chat.status]) {
      result[chat.status] = 0
    }
    result[chat.status]++
  })

  return result
})

const checkedAdmins = ref<string[]>([])
const checkedDepartments = ref<string[]>([])

if (Object.keys(metrics.value).length < 1) fetch_defaults()

interface metricsOptionsType {
  [index: string]: []
}

const filters = JSON.parse(localStorage.getItem('filters') ?? 'null')
const metricsOptions = ref<metricsOptionsType>((filters)
  ? filters.metrics
  : metrics.value.reduce((result, { key }) => ({ ...result, [key]: [] }), {})
)

watch(metrics, (value) => {
  const filters = JSON.parse(localStorage.getItem('filters') ?? 'null')

  metricsOptions.value = (filters)
    ? filters.metrics
    : value.reduce((result, { key }) => ({ ...result, [key]: [] }), {})
})

watch([sortBy, sortType], () => {
  localStorage.setItem('sorting', JSON.stringify({
    sortBy: sortBy.value,
    sortType: sortType.value
  }))
}, { deep: true })

watch([checkedDepartments, checkedStatuses, checkedAdmins, metricsOptions], () => {
  localStorage.setItem('filters', JSON.stringify({
    departments: checkedDepartments.value,
    statuses: checkedStatuses.value,
    admins: checkedAdmins.value,
    metrics: metricsOptions.value
  }))
}, { deep: true })

async function resetFilters() {
  sortBy.value = 'sent'
  sortType.value = 'desc'

  checkedDepartments.value = []
  checkedStatuses.value = []
  checkedAdmins.value = []
  metricsOptions.value = metrics.value.reduce(
    (result, { key }) => ({ ...result, [key]: [] }), {}
  )

  await nextTick()
  localStorage.removeItem('sorting')
  localStorage.removeItem('filters')
}

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
  firstMessage.value = chat?.meta?.firstMessage?.content?.replace(
    /<div class="chat__files">[\s\S]{1,}<\/div>$/g, ''
  ) ?? ''

  if (firstMessage.value.length > 99) {
    firstMessage.value = `${firstMessage.value.slice(0, 100)}...`
  }
  isFirstMessageVisible.value = true
}
</script>

<style>
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
}

.chats__panel.closed {
  min-width: 70px !important;
  width: 70px !important;
  background-color: #18181C;
}

.chats__panel.closed .n-layout-sider-scroll-container {
  min-width: 70px !important;
  width: 70px !important;
}

.chats__panel .search {
  margin: 5px 10px;
}

.chats__panel .search div:first-child {
  width: 100%;
}

.chats__panel aside {
  width: 100% !important;
  max-width: 100% !important;
}

.chats__panel .n-list-item.active {
  background: var(--n-color-hover-modal);
}

.chats__panel .chat__actions {
  padding: 10px;
}

.chat__actions.hide {
  flex-flow: column-reverse !important;
  margin: 0px;
  padding: 10px;
}

.chat__item {
  height: 100%;
  min-width: 50%;
  width: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>