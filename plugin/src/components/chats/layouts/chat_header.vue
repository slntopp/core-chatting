<template>
  <n-text v-if="!chat">Loading...</n-text>
  <div class="grid" v-else>
    <n-button
      ghost
      v-if="appStore.displayMode === 'none'"
      @click="appStore.displayMode = 'full'"
    >
      <n-icon> <open-icon /> </n-icon>
    </n-button>

    <n-space justify="start" align="center">
      <n-tooltip>
        <template #trigger>
          <n-tag round @click="addToClipboard(chat.uuid, notification)">
            <code style="text-decoration: underline; cursor: pointer">
              {{ chat.uuid.slice(0, 8).toUpperCase() }}
            </code>
          </n-tag>
        </template>
        {{ chat.uuid }}
      </n-tooltip>

      <user-avatar round :avatar="members.map((m) => m?.title ?? '').join(' ')"/>
      <n-text>{{ chat.topic ?? members }}</n-text>
      <n-button text @click="startEditChat">
        <n-icon size="20">
          <edit-icon/>
        </n-icon>
      </n-button>

      <n-divider vertical/>
      <members-dropdown
        :admins="chat.admins"
        :members="members"
        @add="startAddMembers"
        @delete="deleteMember"
      />

      <n-divider vertical/>
      <n-tooltip>
        <template #trigger>
          <n-select
            filterable
            label-field="title"
            value-field="uuid"
            placeholder="Responsible"
            style="min-width: 200px; width: 100%"
            :value="chat.responsible"
            :options="users"
            @update:value="changeResponsible"
          />
        </template>
        Responsible
      </n-tooltip>

      <n-divider vertical/>
      <n-tooltip>
        <template #trigger>
          <n-select
            filterable
            label-field="title"
            value-field="key"
            placeholder="Department"
            style="min-width: 200px; width: 100%"
            :value="chat.department"
            :options="departments"
            @update:value="changeDepartment"
          />
        </template>
        Department
      </n-tooltip>

      <n-divider vertical/>
      <chat-status :chat="chat" />

      <n-divider vertical/>
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
          <n-button type="error" size="small" ghost circle @click="deleteChat">
            <template #icon> <delete-icon /> </template>
          </n-button>
        </template>
        Delete chat
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
              <n-button ref="commandsButton" type="success" size="small" ghost circle>
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
    </n-space>

    <n-divider vertical v-if="chat.gateways.length > 0" />
    <n-space :wrap-item="false" v-if="chat.gateways.length > 0">
      <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
        <template #trigger>
          <img height="24" :src="getImageUrl(gateway)" :alt="gateway">
        </template>
        {{ gateway }}
      </n-tooltip>
    </n-space>

    <n-divider
      vertical
      :style="{
        gridRow: (metricsOptions.length > 0) ? '1 / 3' : null, 
        gridColumn: (chat.gateways.length > 0) ? 5 : '3 / 5',
        height: 'calc(100% - 10px)',
        margin: '0 8px'
      }"
    />

    <n-space
      :wrap-item="false"
      :style="{ gridColumn: (appStore.displayMode !== 'half') ? 2 : 1 }"
    >
      <n-text v-for="metric in metricsOptions">
        {{ metric.title }}:
        <n-tag round size="small" type="error" :style="`filter: ${getTagColor(metric)}`">
          {{ metric.key }}
        </n-tag>
      </n-text>
    </n-space>

    <n-space
      vertical
      :wrap-item="false"
      :style="{
        gap: 0,
        whiteSpace: 'nowrap',
        paddingRight: '18px',
        gridRow: (metricsOptions.length > 0) ? '1 / 3' : null,
        gridColumn: (metricsOptions.length > 0) ? '4' : null
      }"
    >
      <n-text>
        Created:
        <n-tooltip>
          <template #trigger>
            <code style="text-decoration: underline">
              {{ getRelativeTime(Number(chat.created), now) }}
            </code>
          </template>
          {{ new Date(Number(chat.created)).toLocaleString() }}
        </n-tooltip>
      </n-text>

      <n-text>
        Last update:
        <n-tooltip>
          <template #trigger>
            <code style="text-decoration: underline">
              {{ getRelativeTime(lastUpdate, now) }}
            </code>
          </template>
          {{ new Date(lastUpdate).toLocaleString() }}
        </n-tooltip>
      </n-text>

      <n-text>
        Lifetime:
        <n-tooltip>
          <template #trigger>
            <code style="text-decoration: underline">
              {{ getRelativeTime(Number(chat.created), now, true) }}
            </code>
          </template>
          {{ new Date(Number(chat.created)).toLocaleString() }}
        </n-tooltip>
      </n-text>
    </n-space>

    <n-button
      ghost
      type="success"
      style="margin-right: 15px"
      v-if="buttonTitle"
      @click="sendMessage"
    >
      {{ buttonTitle }}
    </n-button>
  </div>

  <n-modal v-model:show="isEdit">
    <n-card title="Edit chat options" :bordered="false" size="huge" role="dialog" aria-modal="true"
            style="width: 500px; min-height: 500px">
      <chat-options @close="isEdit = false" is-edit :chat="chat" style="width: 100%" />
    </n-card>
  </n-modal>

  <n-modal v-model:show="isAddDialog">
    <n-card title="Add members" :bordered="false" size="huge" role="dialog" aria-modal="true"
            style="width: 400px;">
      <template v-if="!isDefaultLoading">
        <member-select v-model:value="chatWithNewMembers!.users" :options="availableMembersOptions"/>
        <n-space style="margin-top: 10px" vertical align="end" justify="end">
          <n-button :loading="isAddSaveLoading" @click="saveMembers">Save</n-button>
        </n-space>
      </template>
      <n-spin style="width: 100%;height: 100%; margin: auto" size="large" v-else></n-spin>
    </n-card>
  </n-modal>
</template>

// TODO:
//  - [ ] Wrap Topic Around Avatar
//  - [ ] Make menu draggable (increase width)

<script setup lang="ts">
import {computed, defineAsyncComponent, ref, toRefs} from "vue";
import {
  NButton, NCard, NDivider, NIcon,
  NList, NListItem, NModal, NPopover,
  NSpace, NSpin, NTag, NText, NTooltip,
  NSelect, SelectOption, useNotification
} from "naive-ui";
import {Chat, Kind, Message, User} from "../../../connect/cc/cc_pb";
import {useCcStore} from "../../../store/chatting.ts";
import {useAppStore} from "../../../store/app";
import {useRouter} from "vue-router";
import ChatOptions from "../chat_options.vue";
import UserAvatar from "../../ui/user_avatar.vue";
import MembersDropdown from "../../users/members_dropdown.vue";
import useDefaults from "../../../hooks/useDefaults.ts";
import MemberSelect from "../../users/member_select.vue";
import {addToClipboard, getImageUrl, getRelativeTime} from "../../../functions.ts";
import ChatStatus from "../chat_status.vue";

const EditIcon = defineAsyncComponent(() => import('@vicons/ionicons5/PencilSharp'));
const OpenIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowBack'));
const RefreshIcon = defineAsyncComponent(() => import('@vicons/ionicons5/RefreshOutline'));
const DeleteIcon = defineAsyncComponent(() => import('@vicons/ionicons5/TrashBinOutline'));
const consoleIcon = defineAsyncComponent(() => import('@vicons/ionicons5/TerminalOutline'));

interface ChatHeaderProps {
  chat: Chat
}

interface Metric {
  title: string
  value: number
  key: string
  max: number
  min: number
}

const props = defineProps<ChatHeaderProps>()
const {chat} = toRefs(props)

const appStore = useAppStore()
const store = useCcStore()
const router = useRouter()
const notification = useNotification()
const {fetch_defaults, isDefaultLoading, users, admins, metrics, departments} = useDefaults()

const isEdit = ref<boolean>(false)
const isAddDialog = ref<boolean>(false)
const chatWithNewMembers = ref<Chat>()
const availableMembersOptions = ref<SelectOption[]>([])
const isAddSaveLoading = ref<boolean>(false)

const members = computed(() => {
  const uuids = new Set([...chat!.value.users, ...chat.value.admins])
  const result: User[] = []

  uuids.forEach((uuid) => {
    result.push(store.users.get(uuid) as User)
  })

  return result
})
const me = computed(() => store.me)

fetch_defaults()

const metricsOptions = computed(() => {
  const metricsEntries = Object.entries(chat.value.meta?.data ?? {})
  const result: Metric[] = []

  metricsEntries.forEach(([keyMetric, { kind }]) => {
    const { title, options } = metrics.value.find(
      (metric) => metric.key === keyMetric
    ) ?? {}

    const { key, value } = options?.find(
      (option) => option.value === kind.value
    ) ?? { key: '', value: 0 }

    const optionsValues = options?.map(({ value }) => value) ?? []
    const min = Math.min(...optionsValues)
    const max = Math.max(...optionsValues)
    
    if (title) result.push({ title, value, key, min, max })
  })

  return result
})

const getTagColor = (metric: Metric) => (
  `hue-rotate(${220 - 220 * (metric.value - metric.min) / (metric.max - metric.min)}deg)`
)

const changeResponsible = (uuid: string) => {
  store.update_chat(new Chat({
    ...chat.value, responsible: uuid
  }))
}

const changeDepartment = (key: string) => {
  store.change_department(new Chat({
    ...chat.value, department: key
  }))
}

const refresh = () => {
  if (chat) {
    store.get_messages(chat.value as Chat, false)
  }
}

const deleteChat = async () => {
  if (chat) {
    await store.delete_chat(chat.value! as Chat)
    appStore.displayMode = 'half'
    router.push({name: 'Empty Chat'})
  }
}

interface commandType {
  key: string
  description: string
}

const commandsButton = ref()

const commands = computed(() => {
  const result: commandType[] = []

  chat.value.admins.forEach((uuid) => {
    const bot = store.users.get(uuid)

    if (!bot?.ccIsBot) return
    Object.entries(bot.ccCommands).forEach(([key, description]) => {
      result.push({ key: `/${key}`, description })
    })
  })

  return result
})

const sendCommand = (content: string) => {
  store.send_message(new Message({
    content,
    kind: Kind.FOR_BOT,
    chat: router.currentRoute.value.params.uuid as string
  }))
  commandsButton.value.onClick()
}

const deleteMember = (uuid: string) => {
  const users = chat.value.users.filter((userId) => userId !== uuid)

  store.update_chat({ ...chat.value, users } as Chat)
}

const fetchAvailableUsers = () => {
  availableMembersOptions.value = users.value.map(user => {
    return {
      label: user.title,
      value: user.uuid,
      disabled: me.value.uuid !== user.uuid && admins.value.includes(user.uuid),
    }
  })
}

const startAddMembers = () => {
  fetchAvailableUsers()

  chatWithNewMembers.value = { ...chat.value } as Chat
  isAddDialog.value = true
}

const startEditChat = () => {
  isEdit.value = true
}

const saveMembers = async () => {
  try {
    isAddSaveLoading.value = true
    await store.update_chat(chatWithNewMembers.value as Chat)
    isAddDialog.value = false
  } finally {
    isAddSaveLoading.value = false
  }
}

const now = ref(Date.now())

setInterval(() => now.value = Date.now(), 1000)

const lastUpdate = computed(() =>
  Number(chat.value.meta?.lastMessage?.edited || chat.value.meta?.lastMessage?.sent)
)

const buttonTitle = ref('')
const gridColumns = computed(() =>
  `repeat(${(chat.value?.gateways.length > 0) ? 5 : 4}, auto) 1fr auto`
)
window.addEventListener('message', ({ data, origin }) => {
  if (origin.includes('localhost')) return
  buttonTitle.value = data
})

const sendMessage = (event: MouseEvent) => {
  const button = (event.target as HTMLElement).closest('.n-button') as HTMLButtonElement

  window.top?.postMessage({
    type: 'click-on-button', value: { x: button?.offsetLeft, y: button?.offsetTop }
  }, '*')
}
</script>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: v-bind('gridColumns');
  align-items: center;
  gap: 10px;
}
</style>
