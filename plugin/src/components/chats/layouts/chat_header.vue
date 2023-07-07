<template>
  <n-text v-if="!chat">Loading...</n-text>
  <template v-else>
    <n-space justify="start" align="center">
      <user-avatar round :avatar="members.join(' ')"/>
      <n-text>{{ chat.topic ?? members }}</n-text>
      <n-button text @click="startEditChat">
        <n-icon size="20">
          <edit-icon/>
        </n-icon>
      </n-button>
      <n-divider vertical/>
      <members-dropdown @delete="deleteMember" :members="members"/>
      <n-divider vertical/>
      <n-button type="info" size="small" ghost round @click="refresh">Refresh</n-button>
      <n-divider vertical/>
      <n-button type="error" size="small" ghost round @click="deleteChat">Delete</n-button>
    </n-space>
  </template>

  <n-modal v-model:show="isEdit">
    <n-card title="Edit chat options" :bordered="false" size="huge" role="dialog" aria-modal="true"
            style="width: 500px; height: 500px">
      <chat-options @close="isEdit = false" is-edit :chat="chat"/>
    </n-card>
  </n-modal>
</template>

// TODO:
//  - [ ] Wrap Topic Around Avatar
//  - [ ] Make menu draggable (increase width)

<script setup lang="ts">
import {computed,  ref, toRefs} from "vue";
import {NButton, NCard, NDivider, NIcon, NModal, NSpace, NText} from "naive-ui";
import {Chat} from "../../../connect/cc/cc_pb";
import {useCcStore} from "../../../store/chatting.ts";
import {useRouter} from "vue-router";
import ChatOptions from "../chat_options.vue";
import UserAvatar from "../../ui/user_avatar.vue";
import {PencilSharp as EditIcon} from '@vicons/ionicons5'
import MembersDropdown from "../../users/membersDropdown.vue";

interface ChatHeaderProps {
  chat: Chat
}

const props = defineProps<ChatHeaderProps>()
const {chat} = toRefs(props)

const store = useCcStore()
const router = useRouter()

const isEdit = ref<boolean>(false)

const members = computed(() => {
  return chat!.value.users.map((uuid: string) => store.users.get(uuid)).concat(chat.value.admins.map((uuid: string) => store.users.get(uuid)))
})

const refresh = () => {
  if (chat) {
    store.get_messages(chat.value as Chat, false)
  }
}

const deleteChat = async () => {
  if (chat) {
    await store.delete_chat(chat.value! as Chat)
    router.push({name: 'Empty Chat'})
  }
}

const deleteMember=(uuid:string)=>{
  const users=chat.value.users.filter((u)=>u!==uuid)
  store.update_chat({...chat.value,users})
}

const startEditChat = () => {
  isEdit.value = true
}
</script>

<style scoped></style>