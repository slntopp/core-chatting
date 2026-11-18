<template>
  <n-tooltip>
    <template #trigger>
      <n-text class="when">
        Updated
        <code>{{ updatedLabel }}</code>
      </n-text>
    </template>
    <div class="when__tip">
      <div>Created: {{ createdAbsolute }}</div>
      <div>Last update: {{ updatedAbsolute }}</div>
      <div>Lifetime: {{ lifetimeLabel }}</div>
    </div>
  </n-tooltip>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { NText, NTooltip } from "naive-ui";
import { Chat } from "../../connect/cc/cc_pb";
import { getRelativeTime } from "../../functions";

interface ChatDatesProps {
  chat: Chat;
}

const props = defineProps<ChatDatesProps>();
const now = ref(Date.now());
const interval = ref<number>();

const created = computed(() => Number(props.chat.created));
const lastUpdate = computed(() =>
  Number(
    props.chat.meta?.lastMessage?.edited || props.chat.meta?.lastMessage?.sent
  )
);
const updatedAt = computed(() => lastUpdate.value || created.value);

const updatedLabel = computed(() =>
  getRelativeTime(updatedAt.value, now.value)
);
const lifetimeLabel = computed(() =>
  getRelativeTime(created.value, now.value, true)
);
const createdAbsolute = computed(() =>
  new Date(created.value).toLocaleString()
);
const updatedAbsolute = computed(() =>
  updatedAt.value ? new Date(updatedAt.value).toLocaleString() : "—"
);

onMounted(() => {
  interval.value = setInterval(() => (now.value = Date.now()), 1000);
});

onUnmounted(() => {
  clearInterval(interval.value);
});
</script>

<style scoped>
.when {
  white-space: nowrap;
}
.when code {
  text-decoration: underline;
}
.when__tip {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
</style>
