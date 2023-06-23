<template>
    <n-list>
        <template #header>
            <chat-header style="height: 5vh" />
        </template>

        <n-scrollbar style="height: 80vh; max-width: 80%;">
            <n-list-item v-for="message in messages">
                <message-view :message="message" />
            </n-list-item>
        </n-scrollbar>

        <template #footer>
            <n-space style="min-height: 10vh" align="end">
                <n-avatar round size="medium">Me</n-avatar>

                <n-input type="textarea" size="small" :autosize="{
                    minRows: 2,
                    maxRows: 5
                }" style="min-width: 50vw;" placeholder="Type your message" v-model:value="current_message.content" />

                <n-space>
                    <n-button type="success" ghost circle size="small" @click="handle_send()">
                        <template #icon>
                            <n-icon :component="SendOutline" />
                        </template>
                    </n-button>
                    <n-button type="warning" ghost circle size="small" @click="handle_send(Kind.ADMIN_ONLY)">
                        <template #icon>
                            <n-icon :component="ClipboardOutline" />
                        </template>
                    </n-button>
                </n-space>
            </n-space>
        </template>
    </n-list>
</template>

<script setup lang="ts">
import { defineAsyncComponent, h, ref } from 'vue';
import {
    NSpace, NButton, NIcon,
    NAvatar, NText, NInput,
    NList, NListItem, NScrollbar, NDivider
} from 'naive-ui';

import { useRoute } from 'vue-router';
import { useCcStore } from '../../../store/chatting';
import { Chat, Message, Kind } from '../../../connect/cc/cc_pb';

const SendOutline = defineAsyncComponent(() => import('@vicons/ionicons5/SendOutline'));
const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));

const MessageView = defineAsyncComponent(() => import('../../../components/chats/message.vue'));

const route = useRoute();

const store = useCcStore()

const chat = ref(store.chats.get(route.params.uuid as string))

store.resolve([...chat.value!.users, ...chat.value!.admins])

function chatHeader() {
    let users = chat.value!.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')
    let admins = chat.value!.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')

    let members = users.concat(admins)
    let avatar_title = members.map(el => el[0]).join(',')

    return h(
        NSpace, { justify: "start", align: 'center' }, () => [
            h(NAvatar, { round: true, size: "medium" }, () => avatar_title),
            h(NText, {}, () => chat.value!.topic ?? members),
            h(NDivider, { vertical: true }),
            h(NText, { depth: "3" }, () => `${members.length} members`),
            h(NDivider, { vertical: true }),
            h(NButton, {
                type: 'info', size: 'small',
                ghost: true, round: true,
                onClick: get_messages,
            }, () => 'Refresh')
        ]
    )
}

const messages = ref<Message[]>([])
async function get_messages() {
    let res = await store.get_messages(chat.value! as Chat)
    messages.value = res.messages
}
get_messages()

const current_message = ref<Message>(new Message({
    chat: route.params.uuid as string,
    content: '',
}))

async function handle_send(kind = Kind.DEFAULT) {
    current_message.value.kind = kind
    await store.send_message(current_message.value as Message)

    current_message.value = new Message({
        chat: route.params.uuid as string,
        content: '',
    })
    await get_messages()
}
</script>