<template>
  <n-text v-if="!chat">Loading...</n-text>
  <template v-else>
    <n-space justify="start" align="center">
      <user-avatar round :avatar="members.join(' ')" />
      <n-text>{{ chat.topic ?? members }}</n-text>
      <n-button text @click="startEditChat">
        <n-icon size="20">
          <edit-icon />
        </n-icon>
      </n-button>
      <n-divider vertical />
      <n-dropdown trigger="hover" :options="membersOptions">
        <n-text>{{ members.length }} members</n-text>
      </n-dropdown>
      <n-divider vertical />
      <n-button type="info" size="small" ghost round @click="refresh">Refresh</n-button>
      <n-divider vertical />
      <n-button type="error" size="small" ghost round @click="deleteChat">Delete</n-button>
    </n-space>
  </template>

  <n-modal v-model:show="isEdit">
    <n-card title="Edit chat options" :bordered="false" size="huge" role="dialog" aria-modal="true"
      style="width: 500px; height: 500px">
      <chat-options @close="isEdit = false" is-edit :chat="chat" />
    </n-card>
  </n-modal>
</template>

// TODO:
//  - [ ] Wrap Topic Around Avatar
//  - [ ] Make menu draggable (increase width)

<script setup lang="ts">
import { computed, h, ref, toRefs } from "vue";
import { NButton, NCard, NDivider, NDropdown, NIcon, NModal, NSpace, NText } from "naive-ui";
import { Chat } from "../../../connect/cc/cc_pb.ts";
import { useCcStore } from "../../../store/chatting.ts";
import { useRouter } from "vue-router";
import ChatOptions from "../chatOptions.vue";
import UserAvatar from "../../ui/userAvatar.vue";
import { PencilSharp as EditIcon } from '@vicons/ionicons5'

interface ChatHeaderProps {
  chat: Chat
}

const props = defineProps<ChatHeaderProps>()
const { chat } = toRefs(props)

const store = useCcStore()
const router = useRouter()

const isEdit = ref<boolean>(false)

const members = computed(() => {
  return chat!.value.users.map((uuid: string) => store.users.get(uuid)?.title ?? 'Unknown').concat(chat.value.admins.map((uuid: string) => store.users.get(uuid)?.title ?? 'Unknown'))
})

const membersOptions = computed(() => {
  return [...new Set(members.value)].map((m: any, i: number) => ({
    key: m as string,
    label: m as string,
    icon: renderIcon(members.value.map((el: any) => el[i]).join(','))
  }))
})

const renderIcon = (icon: string) => {
  return () => {
    return h(UserAvatar, { round: true, size: "medium", avatar: icon })
  }
}

const refresh = () => {
  if (chat) {
    store.get_messages(chat.value as Chat, false)
  }
}

const deleteChat = async () => {
  if (chat) {
    await store.delete_chat(chat.value! as Chat)
    router.push({ name: 'Empty Chat' })
  }
}

const startEditChat = () => {
  isEdit.value = true
}
</script>

<style scoped></style>