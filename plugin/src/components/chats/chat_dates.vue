<template>
  <n-space
    vertical
    :wrap-item="false"
    :style="{ gap: 0, whiteSpace: 'nowrap', paddingRight: '18px' }"
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
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { NSpace, NText, NTooltip } from 'naive-ui'
import { Chat } from '../../connect/cc/cc_pb'
import { getRelativeTime } from '../../functions'

interface ChatDatesProps {
  chat: Chat
}

const props = defineProps<ChatDatesProps>()
const now = ref(Date.now())

setInterval(() => now.value = Date.now(), 1000)

const lastUpdate = computed(() =>
  Number(props.chat.meta?.lastMessage?.edited || props.chat.meta?.lastMessage?.sent)
)
</script>
