<template>
  <n-text v-if="!chat">Loading...</n-text>
  <template v-else>
    <n-space justify="start" align="center">
      <user-avatar round :avatar="members.map(m=>m.title).join(' ')"/>
      <n-text>{{ chat.topic ?? members }}</n-text>
      <n-button text @click="startEditChat">
        <n-icon size="20">
          <edit-icon/>
        </n-icon>
      </n-button>
      <n-divider vertical/>
      <members-dropdown @add="startAddMembers" @delete="deleteMember" :members="members"/>
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

  <n-modal v-model:show="isAddDialog">
    <n-card title="Add members" :bordered="false" size="huge" role="dialog" aria-modal="true"
            style="width: 400px;">
      <template v-if="!isDefaultLoading">
        <member-select v-model:value="chatWithNewMembers!.users" :options="availableMembersOptions"/>
        <n-space style="margin-top: 10px" vertical align="end" justify="end">
          <n-button :loading="isAddSaveLoading" @click="saveMembers">Save</n-button>
        </n-space>
      </template>
      <n-spin style="width: 100%;height: 100%; margin: auto" size="large" v-else></n-spin>
    </n-card>
  </n-modal>
</template>

// TODO:
//  - [ ] Wrap Topic Around Avatar
//  - [ ] Make menu draggable (increase width)

<script setup lang="ts">
import {computed, ref, toRefs} from "vue";
import {NButton, NCard, NDivider, NIcon, NModal, NSpace, NSpin, NText, SelectOption} from "naive-ui";
import {Chat, User} from "../../../connect/cc/cc_pb";
import {useCcStore} from "../../../store/chatting.ts";
import {useRouter} from "vue-router";
import ChatOptions from "../chat_options.vue";
import UserAvatar from "../../ui/user_avatar.vue";
import {PencilSharp as EditIcon} from '@vicons/ionicons5'
import MembersDropdown from "../../users/members_dropdown.vue";
import useDefaults from "../../../hooks/useDefaults.ts";
import MemberSelect from "../../users/member_select.vue";

interface ChatHeaderProps {
  chat: Chat
}

const props = defineProps<ChatHeaderProps>()
const {chat} = toRefs(props)

const store = useCcStore()
const router = useRouter()
const {fetch_defaults, isDefaultLoading, users, admins} = useDefaults()

const isEdit = ref<boolean>(false)
const isAddDialog = ref<boolean>(false)
const chatWithNewMembers = ref<Chat>()
const availableMembersOptions = ref<SelectOption[]>([])
const isAddSaveLoading = ref<boolean>(false)

const members = computed(() => {
  return chat!.value.users.map((uuid: string) => store.users.get(uuid) as User).concat(chat.value.admins.map((uuid: string) => store.users.get(uuid)  as User))
})
const me = computed(() => store.me)

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

const deleteMember = (uuid: string) => {
  const users = chat.value.users.filter((u) => u !== uuid)
  store.update_chat({...chat.value, users} as Chat)
}

const fetchAvailableUsers = async () => {
  await fetch_defaults()
  availableMembersOptions.value = users.value.map(user => {
    return {
      label: user.title,
      value: user.uuid,
      disabled: me.value.uuid !== user.uuid && admins.value.includes(user.uuid),
    }
  })
}

const startAddMembers = () => {
  if (!users.value.length) {
    fetchAvailableUsers()
  }
  chatWithNewMembers.value = {...chat.value} as Chat
  isAddDialog.value = true
}

const startEditChat = () => {
  isEdit.value = true
}

const saveMembers = async () => {
  try {
    isAddSaveLoading.value = true
    await store.update_chat(chatWithNewMembers.value as Chat)
    isAddDialog.value = false
  } finally {
    isAddSaveLoading.value = false
  }
}
</script>

<style scoped></style>