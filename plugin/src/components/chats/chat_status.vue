<template>
  <n-popover
    placement="right"
    trigger="manual"
    :show="isVisible"
    @clickoutside="clickOutside"
  >
    <template #trigger>
      <n-text class="status" @click.stop="isVisible = !isVisible">
        <span :style="{ color: getStatusColor(chat.status) }">
          {{ getStatus(chat.status) }}
        </span>
        <n-icon> <swap-icon /> </n-icon>
      </n-text>
    </template>

    <div style="display: flex; gap: 5px">
      <n-button
        ghost
        v-for="item of statuses"
        :key="item.value"
        :disabled="chat.status === item.value"
        :color="getStatusColor(item.value)"
        @click="updateChat(item.value)"
      >
        {{ item.label }}
      </n-button>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref } from 'vue'
import { NIcon, NButton, NPopover, NText } from 'naive-ui'
import { Chat, Status } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting.ts'

interface ChatStatusProps {
  chat: Chat,
  hiddenColor?: boolean
}

const props = defineProps<ChatStatusProps>()
const store = useCcStore()

const SwapIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/SwapHorizontal')
);

const isVisible = ref(false)

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

function updateChat (status: Status) {
  //@ts-ignore
  store.change_status({ ...props.chat, status } as Chat)
  isVisible.value = false
}

function getStatusColor (status: Status) {
  if (props.hiddenColor) return undefined
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

function getStatus (statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

function clickOutside () {
  setTimeout(() => { isVisible.value = false })
}
</script>

<script lang="ts">
export default {
  name: 'chat-status'
}
</script>

<style scoped>
.status {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}
</style>
