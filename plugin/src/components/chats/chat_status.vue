<template>
  <n-popover
    placement="right"
    trigger="manual"
    :show="isVisible"
    @clickoutside="isVisible = false"
  >
    <template #trigger>
      <n-text class="status" @click.stop="isVisible = !isVisible">
        <span :style="({ color: statusColor } as StyleValue)">{{ getStatus(chat.status) }}</span>
        <n-icon> <swap-icon /> </n-icon>
      </n-text>
    </template>

    <div style="display: flex; gap: 5px">
      <n-select style="min-width: 150px" v-model:value="status" :options="statuses" />
      <n-button @click="updateChat">Ok</n-button>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { StyleValue, computed, defineAsyncComponent, ref } from 'vue'
import { NIcon, NSelect, NButton, NPopover, NText } from 'naive-ui'
import { Chat, Status } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting'

interface ChatStatusProps {
  chat: Chat,
  hiddenColor?: boolean
}

const props = defineProps<ChatStatusProps>()
const store = useCcStore()

const SwapIcon = defineAsyncComponent(() => import('@vicons/ionicons5/SwapHorizontal'));

const isVisible = ref(false)
const status = ref(props.chat.status)

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

function updateChat () {
  store.update_chat({ ...props.chat, status: status.value } as Chat)
  isVisible.value = false
}

const statusColor = computed(() => {
  if (props.hiddenColor) return null
  switch (props.chat.status) {
    case 0:
      return '#5084ff'
    case 1:
      return '#1ea01e'
    case 2:
      return '#ff8300'
    case 3:
      return '#e23535'
    default:
      return null
  }
})

function getStatus (statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
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
}
</style>
