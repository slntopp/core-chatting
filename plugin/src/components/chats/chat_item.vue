<template>
  <n-list-item class="chat" @click="goToChat">
    <n-space :wrap-item="false" justify="start">
      <user-avatar round size="large" :avatar="members.join(' ')"/>
      <n-space class="preview" vertical>
        <n-text class="topic">{{ chatTopic }}</n-text>
        <n-text depth="3">{{ sub }}</n-text>
      </n-space>
      <div style="    position: absolute;
    right: 15px;">
        <n-badge v-if="isUnreadMessages" :value="chat.meta!.unread" :max="99" size="24" :offset="[12, 12]"/>
      </div>
    </n-space>
  </n-list-item>
</template>

<script setup lang="ts">
import UserAvatar from "../ui/user_avatar.vue";
import {NBadge, NListItem, NSpace, NText} from "naive-ui";
import {Chat} from "../../connect/cc/cc_pb";
import {computed, toRefs} from "vue";
import {useCcStore} from "../../store/chatting.ts";
import {useRouter} from "vue-router";

interface ChatItemProps {
  chat: Chat
  uuid: string
}

const props = defineProps<ChatItemProps>()
const {chat, uuid} = toRefs(props)

const store = useCcStore()
const router = useRouter()

const users = computed(() => chat.value.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))
const admins = computed(() => chat.value.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))

const members = computed(() => users.value.concat(admins.value))

const sub = computed(() => {
  if (chat.value.meta && chat.value.meta.lastMessage)
  {
    const lastMessage= chat.value.meta.lastMessage!.content.length > 20 ? chat.value.meta!.lastMessage!.content.slice(0, 16) + '...' : chat.value.meta!.lastMessage!.content
    return lastMessage.replace(/(<([^>]+)>)/ig,"")
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
</script>

<style scoped lang="scss">
.chat {
  .preview {
    width: calc(100% - 80px);

    .topic {
      word-break: break-all;
    }
  }
}
</style>