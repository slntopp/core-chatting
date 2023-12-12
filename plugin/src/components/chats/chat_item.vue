<template>
  <n-list-item class="chat" @click="goToChat">
    <n-space :wrap-item="false" justify="start">
      <user-avatar round size="large" :avatar="members.join(' ')"/>
      <div v-if="!hideMessage" class="preview">
        <n-text
          depth="3"
          class="sub"
          @click.stop="openUser(chat.owner)"
        >
          {{ store.users.get(chat.owner)?.title }}
        </n-text>

        <n-space
          v-if="appStore.displayMode === 'full'"
          style="flex-direction: column; gap: 0; grid-row: 1 / 3; grid-column: 2"
          :wrap-item="false"
        >
          <n-space :wrap-item="false">
            <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
              <template #trigger>
                <img height="24" :src="`/cc.ui/assets/${gateway}.png`" :alt="gateway">
              </template>
              {{ gateway }}
            </n-tooltip>
          </n-space>
          <span>
            Department:
            <n-text italic style="margin-left: 4px">{{ chat.department || 'none' }}</n-text>
          </span>
        </n-space>

        <div class="time" v-show="appStore.displayMode === 'full'">
          <div>Created: {{ new Date(Number(chat.created)).toLocaleDateString() }}</div>
          <div v-if="lastUpdate">Last update: {{ new Date(Number(lastUpdate)).toLocaleDateString() }}</div>
        </div>

        <n-icon
          size="18"
          @mouseenter="(e) => emits('hover', e.clientX, e.clientY, chat.uuid)"
          @mouseleave="emits('hoverEnd')"
        >
          <mail-icon />
        </n-icon>

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
        <n-text class="topic">{{ chatTopic }}</n-text>
      </div>

      <div style="position: absolute; right: 18px; top: -10px">
        <n-badge v-if="isUnreadMessages" :value="chat.meta!.unread" :max="99" size="24" :offset="[12, 12]"/>
      </div>
    </n-space>
  </n-list-item>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, toRefs} from "vue";
import {useRouter} from "vue-router";
import {NBadge, NIcon, NListItem, NSpace, NText, NTooltip, useNotification} from "naive-ui";
import {Chat} from "../../connect/cc/cc_pb";
import {useCcStore} from "../../store/chatting.ts";
import {addToClipboard} from "../../functions.ts";
import UserAvatar from "../ui/user_avatar.vue";
import ChatStatus from "./chat_status.vue";
import {useAppStore} from "../../store/app.ts";

const MailIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));

interface ChatItemProps {
  chat: Chat
  uuid: string
  hideMessage: boolean
}

const props = defineProps<ChatItemProps>()
const emits = defineEmits(['hover', 'hoverEnd'])
const { chat, uuid } = toRefs(props)

const store = useCcStore()
const appStore = useAppStore()
const router = useRouter()
const notification = useNotification()

const users = computed(() => chat.value.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))
const admins = computed(() => chat.value.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))

const members = computed(() => users.value.concat(admins.value))

// const sub = computed(() => {
//   if (chat.value.meta && chat.value.meta.lastMessage) {
//     return chat.value.meta.lastMessage!.content.replace(/(<([^>]+)>)/ig, "")
//   }

//   return "No messages yet"
// })

const chatTopic = computed(() => {
  const topic = chat.value.topic ?? members.value.join(', ')
  const words = topic.split(' ')

  if (topic.length > 512)
    return topic.slice(0, 509) + '...';
  else if (words.length > 16) return words.slice(0, 16).join(' ') + '...'

  return topic || '-'
})

const isUnreadMessages = computed(() => chat.value.meta && chat.value.meta.unread > 0)

const previewColumns = computed(() =>
  (appStore.displayMode === 'full') ? '1fr auto auto auto' : '1fr auto auto'
)

const subDecoration = computed(() =>
  (window.top) ? 'underline' : 'none'
)

const chatRightColumn = computed(() =>
  (appStore.displayMode === 'full') ? 5 : 3
)

const lastUpdate = computed(() =>
  Number(chat.value.meta?.lastMessage?.edited || chat.value.meta?.lastMessage?.sent)
)

function goToChat() {
  router.push({ name: 'Chat', params: { uuid: uuid.value } })

  window.top?.postMessage({
    type: 'get-user',
    value: { uuid: chat.value.owner }
  }, '*')
}

function openUser(uuid: string) {
  window.top?.postMessage({ type: 'open-user', value: { uuid } }, '*')
}
</script>

<style scoped lang="scss">
.chat {
  .preview {
    display: grid;
    grid-template-columns: v-bind(previewColumns);
    align-items: center;
    gap: 5px 15px;
    width: calc(100% - 52px);

    .sub {
      display: block;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      width: fit-content;
      overflow: hidden;
      text-overflow: ellipsis;

      &:hover {
        text-decoration: v-bind(subDecoration);
      }
    }

    .topic {
      word-break: break-all;
    }

    .time {
      grid-row: 1 / 3;
      grid-column: 3;
    }
  }

  &__right {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    grid-row: 1 / 3;
    grid-column: v-bind(chatRightColumn);
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