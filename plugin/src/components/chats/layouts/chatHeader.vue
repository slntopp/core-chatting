<template>
  <n-text v-if="!chat">Loading...</n-text>
  <template v-else>
    <n-space justify="start" align="center">
      <n-avatar round size="medium">{{ avatar_title }}</n-avatar>
      <n-text>{{ chat.topic ?? members }}</n-text>
      <n-divider vertical/>
      <n-text depth="3">{{ `${members.length} members` }}</n-text>
      <n-divider vertical/>
      <n-button type="info" size="small" ghost round @click="refresh">Refresh</n-button>
      <n-divider vertical/>
      <n-button type="error" size="small" ghost round @click="deleteChat">Delete</n-button>
    </n-space>
  </template>
</template>
<script setup lang="ts">
import {computed, toRefs} from "vue";
import {NAvatar, NButton, NDivider, NSpace, NText} from "naive-ui";
import {Chat} from "../../../connect/cc/cc_pb.ts";
import {useCcStore} from "../../../store/chatting.ts";
import {useRouter} from "vue-router";

const props = defineProps(['chat'])
const {chat} = toRefs(props)

const store = useCcStore()
const router = useRouter()

const members = computed(() => {
  if (!chat || !chat.value) {
    return []
  }

  return chat.value.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown').concat(chat.value.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown'))
})

const avatar_title = computed(() => members.value.map(el => el[0]).join(','))

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
</script>

<style scoped>

</style>