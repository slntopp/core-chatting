<template>
  <n-tooltip>
    <template #trigger>
      <n-button type="success" size="small" ghost circle @click="copyLink">
        <template #icon> <copy-icon /> </template>
      </n-button>
    </template>
    Copy chat link
  </n-tooltip>

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
      <n-popconfirm @positive-click="deleteChat">
        <template #trigger>
          <n-button type="error" size="small" ghost circle>
          <template #icon> <delete-icon /> </template>
        </n-button>
        </template>
        Are you sure you want to delete the chat?
      </n-popconfirm>
    </template>
    Delete chat
  </n-tooltip>

  <n-tooltip v-if="buttonTitle">
    <template #trigger>
      <n-button
        ghost
        circle
        type="warning"
        size="small"
        @mouseenter="hoverInstancesButton"
      >
        <template #icon> <list-icon /> </template>
      </n-button>
    </template>
    Instances list
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
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NList, NListItem, NPopover, NTooltip, NPopconfirm, useNotification } from 'naive-ui'

import { ConnectError } from '@connectrpc/connect'
import { Chat, Kind, Message } from '../../connect/cc/cc_pb'
import { useAppStore } from '../../store/app.ts'
import { useCcStore } from '../../store/chatting.ts'
import { addToClipboard } from '../../functions.ts'

const copyIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/CopyOutline')
)
const refreshIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/RefreshOutline')
)
const deleteIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/TrashBinOutline')
)
const listIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/ListOutline')
)
const consoleIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/TerminalOutline')
)

interface ChatActionsProps {
  chat: Chat
}

const appStore = useAppStore()
const store = useCcStore()
const notification = useNotification()

const props = defineProps<ChatActionsProps>()
const router = useRouter()

function copyLink () {
  const url = (new URL(appStore.conf?.params.fullUrl))

  url.searchParams.set('chat', props.chat.uuid)
  addToClipboard(url.href, notification)
}

async function refresh () {
  try {
    await store.get_messages(props.chat, false)
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? '[Error]: Unknown'
    })
  }
}

async function deleteChat () {
  try {
    await store.delete_chat(props.chat)
    appStore.displayMode = 'half'
    
    router.push({name: 'Empty Chat'})
    notification.success({ title: 'Done' })
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? '[Error]: Unknown'
    })
  }
}

const buttonTitle = ref('')

window.addEventListener('message', ({ data, origin }) => {
  if (origin.includes('localhost:8081')) return
  if (data.type !== 'button-title') return
  buttonTitle.value = data.value
})

const hoverInstancesButton = (event: MouseEvent) => {
  const button = (event.target as HTMLElement).closest('.n-button') as HTMLButtonElement

  window.top?.postMessage({
    type: 'click-on-button', value: { x: button?.offsetLeft, y: button?.offsetTop }
  }, '*')
}

interface commandType {
  key: string
  description: string
}

const commandsButton = ref()

const commands = computed(() => {
  const result: commandType[] = []

  props.chat.admins.forEach((uuid) => {
    const bot = store.users.get(uuid)

    if (!bot?.ccIsBot) return
    Object.entries(bot.ccCommands).forEach(([key, description]) => {
      result.push({ key: `/${key}`, description })
    })
  })

  return result
})

async function sendCommand (content: string) {
  try {
    await store.send_message(new Message({
      content,
      kind: Kind.FOR_BOT,
      chat: router.currentRoute.value.params.uuid as string
    }))
    commandsButton.value.onClick()
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? '[Error]: Unknown'
    })
  }
}
</script>
