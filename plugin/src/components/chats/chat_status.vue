<template>
  <n-popover
    placement="right"
    trigger="manual"
    :show="isVisible"
    @clickoutside="clickOutside"
  >
    <template #trigger>
      <n-text class="status" @click.stop="isVisible = !isVisible">
        <span :style="{ color: statusColor(chat.status) }">
          {{ getStatus(chat.status) }}
        </span>
        <n-icon> <swap-icon /> </n-icon>
      </n-text>
    </template>

    <n-space style="width: 370px; gap: 5px" :wrap-item="false">
      <n-button
        ghost
        v-for="item of statuses"
        :key="item.value"
        :disabled="chat.status === item.value"
        :loading="(selectStatus === item.value) ? isLoading : false"
        :color="statusColor(item.value)"
        @click="updateChat(item.value)"
      >
        {{ item.label.replaceAll('_', ' ') }}
      </n-button>
    </n-space>
  </n-popover>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref } from 'vue'
import { NIcon, NButton, NPopover, NSpace, NText, useNotification } from 'naive-ui'
import { ConnectError } from '@connectrpc/connect'

import { Chat, Status } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting.ts'
import { getStatusColor } from '../../functions.ts'

interface ChatStatusProps {
  chat: Chat,
  hiddenColor?: boolean
}

const props = defineProps<ChatStatusProps>()
const store = useCcStore()
const notification = useNotification()

const SwapIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/SwapHorizontal')
);

const selectStatus = ref(-1)
const isVisible = ref(false)
const isLoading = ref(false)

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

async function updateChat (status: Status) {
  selectStatus.value = status
  isLoading.value = true

  try {
    const chat = await store.change_status({ ...props.chat, status } as Chat)

    store.chats.set(chat.uuid, chat)
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message
    })
  } finally {
    selectStatus.value = -1
    isLoading.value = false
    isVisible.value = false
  }
}

function statusColor (status: Status) {
  if (props.hiddenColor) return undefined
  return getStatusColor(status)
}

function getStatus (statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

function clickOutside () {
  setTimeout(() => { isVisible.value = false })
}
</script>

<script lang="ts">
export default { name: 'chat-status' }
</script>

<style scoped>
.status {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  gap: 5px;
  cursor: pointer;
}

.status__change {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  max-width: 60vw;
}
</style>
