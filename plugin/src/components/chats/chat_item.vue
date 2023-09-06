<template>
  <n-list-item class="chat" @click="goToChat">
    <n-space :wrap-item="false" justify="start">
      <user-avatar round size="large" :avatar="members.join(' ')"/>
      <div v-if="!hideMessage" class="preview">
        <n-text class="topic">{{ chatTopic }}</n-text>
        <div class="chat__right">
          <n-tooltip>
            <template #trigger>
              <code style="text-decoration: underline" @click.stop="addToClipboard(uuid, notification)">
                {{ uuid.slice(0, 8).toUpperCase() }}
              </code>
            </template>
            {{ uuid }}
          </n-tooltip>

          <chat-status :chat="chat" />
        </div>
        <n-text class="sub" depth="3">{{ sub }}</n-text>
      </div>

      <div style="position: absolute; right: 18px; top: -10px">
        <n-badge v-if="isUnreadMessages" :value="chat.meta!.unread" :max="99" size="24" :offset="[12, 12]"/>
      </div>
    </n-space>
  </n-list-item>
</template>

<script setup lang="ts">
import {computed, toRefs} from "vue";
import {useRouter} from "vue-router";
import {NBadge, NListItem, NSpace, NText, NTooltip, useNotification} from "naive-ui";
import {Chat} from "../../connect/cc/cc_pb";
import {useCcStore} from "../../store/chatting.ts";
import {addToClipboard} from "../../functions.ts";
import UserAvatar from "../ui/user_avatar.vue";
import ChatStatus from "./chat_status.vue";

interface ChatItemProps {
  chat: Chat
  uuid: string
  hideMessage: boolean
}

const props = defineProps<ChatItemProps>()
const { chat, uuid } = toRefs(props)

const store = useCcStore()
const router = useRouter()
const notification = useNotification()

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

  return topic || '-'
})

const isUnreadMessages = computed(() => chat.value.meta && chat.value.meta.unread > 0)

const goToChat = () => {
  router.push({ name: 'Chat', params: { uuid: uuid.value }, state: { displayMode: 'none' } })
}
</script>

<style scoped lang="scss">
.chat {
  .preview {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 5px;
    width: calc(100% - 52px);

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
  }

  &__right {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    grid-row: 1 / 3;
    grid-column: 2;
    padding: 5px 10px;
    border: 1px solid var(--n-border-color-popover);
    border-radius: 15px;
  }
}
</style>

<style>
.n-list .n-list-item .n-list-item__main {
  overflow: hidden;
}
</style>