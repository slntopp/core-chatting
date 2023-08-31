<template>
  <n-list-item class="chat" @click="goToChat">
    <n-space :wrap-item="false" justify="start">
      <user-avatar round size="large" :avatar="members.join(' ')"/>
      <div v-if="!hideMessage" class="preview">
        <n-text class="topic">{{ chatTopic }}</n-text>
        <n-tag round :title="uuid">{{ uuid.slice(0, 4) }}...</n-tag>
        <n-text class="sub" depth="3">{{ sub }}</n-text>

        <n-popover
          placement="right"
          trigger="manual"
          :show="isVisible"
          @clickoutside="isVisible = false"
        >
          <template #trigger>
            <n-text class="status" @click.stop="isVisible = !isVisible">
              {{ getStatus(chat.status) }}
              <n-icon> <swap-horizontal /> </n-icon>
            </n-text>
          </template>

          <div style="display: flex; gap: 5px">
            <n-select style="min-width: 150px" v-model:value="status" :options="statuses" />
            <n-button @click="updateChat">Ok</n-button>
          </div>
        </n-popover>
      </div>

      <div style="position: absolute; right: 15px">
        <n-badge v-if="isUnreadMessages" :value="chat.meta!.unread" :max="99" size="24" :offset="[12, 12]"/>
      </div>
    </n-space>
  </n-list-item>
</template>

<script setup lang="ts">
import UserAvatar from "../ui/user_avatar.vue";
import {NBadge, NListItem, NSpace, NText, NTag, NIcon, NSelect, NButton, NPopover} from "naive-ui";
import {SwapHorizontal} from "@vicons/ionicons5";
import {Chat, Status} from "../../connect/cc/cc_pb";
import {computed, ref, toRefs} from "vue";
import {useCcStore} from "../../store/chatting.ts";
import {useRouter} from "vue-router";

interface ChatItemProps {
  chat: Chat
  uuid: string
  hideMessage:boolean
}

const props = defineProps<ChatItemProps>()
const {chat, uuid} = toRefs(props)

const store = useCcStore()
const router = useRouter()

const users = computed(() => chat.value.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))
const admins = computed(() => chat.value.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))

const members = computed(() => users.value.concat(admins.value))

const sub = computed(() => {
  if (chat.value.meta && chat.value.meta.lastMessage) {
    return chat.value.meta.lastMessage!.content.replace(/(<([^>]+)>)/ig, "")
  }

  return "No messages yet"
})

const chatTopic = computed(() => {
  const topic = chat.value.topic ?? members.value.join(', ')
  const words = topic.split(' ')

  if (topic.length > 512)
    return topic.slice(0, 509) + '...';
  else if (words.length > 16) return words.slice(0, 16).join(' ') + '...'

  return topic
})


const isUnreadMessages = computed(() => chat.value.meta && chat.value.meta.unread > 0)

const goToChat = () => {
  router.push({name: 'Chat', params: {uuid: uuid.value}})
}

const getStatus = (statusCode: Status | number) => {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

const isVisible = ref(false)
const status = ref(chat.value.status)

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

const updateChat = () => {
  store.update_chat({ ...chat.value, status: status.value } as Chat)
  isVisible.value = false
}
</script>

<style scoped lang="scss">
.chat {
  .preview {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 5px;
    width: calc(100% - 80px);

    *:nth-child(even) {
      justify-self: end;
    }

    .sub {
      display: block;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .topic {
      word-break: break-all;
    }

    .status {
      display: flex;
      align-items: center;
      gap: 5px;
    }
  }
}
</style>

<style>
.n-list .n-list-item .n-list-item__main {
  overflow: hidden;
}
</style>