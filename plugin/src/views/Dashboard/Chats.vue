<template>
    <n-layout-sider>
        <n-scrollbar style="height: 90vh">
            <n-list hoverable clickable>
                <template #header>
                    <n-space justify="center" align="center" style="height: 5vh">
                        <n-button ghost type="success" @click="router.push({ name: 'Start Chat' })">
                            <template #icon>
                                <n-icon :component="ChatbubbleEllipsesOutline" />
                            </template>
                            Start Chat
                        </n-button>
                    </n-space>
                </template>
                <chat-list-item v-for="chat in chats" :uuid="chat.uuid" :chat="chat" />
            </n-list>
        </n-scrollbar>
    </n-layout-sider>
    <n-layout-content>
        <router-view />
    </n-layout-content>
</template>

<script setup lang="ts">
import { defineAsyncComponent, h, computed } from 'vue';
import {
    NLayoutSider, NLayoutContent,
    NScrollbar, NList, NListItem,
    NSpace, NButton, NIcon,
    NAvatar, NText
} from 'naive-ui';

import { useCcStore } from '../../store/chatting.ts';

import { useRouter } from 'vue-router';
import { Chat } from '../../connect/cc/cc_pb';

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));

const store = useCcStore();
const router = useRouter();

async function sync() {
    await store.list_chats();
}
sync()

store.resolve()

const chats = computed(() => {
    let res: Chat[] = []
    store.chats.forEach((chat) => {
        res.push(chat)
    })

    let sortable = (chat: Chat) => {
        if (chat.meta && chat.meta.lastMessage)
            return chat.meta.lastMessage.sent
        return chat.created
    }

    return res.sort((a: Chat, b: Chat) => {
        return Number(sortable(b) - sortable(a))
    })
})

interface ChatListItemProps {
    uuid: string
    chat: Chat
}
function chatListItem(props: ChatListItemProps) {
    let { chat } = props;

    let users = chat.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')
    let admins = chat.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')

    let members = users.concat(admins)
    let avatar_title = members.map(el => el[0]).join(',')

    let sub = "No messages yet"
    if (chat.meta && chat.meta.lastMessage)
        sub = chat.meta!.lastMessage!.content.slice(0, 16) + '...'

    return h(NListItem, {
        "onClick": () => router.push({ name: 'Chat', params: { uuid: props.uuid } })
    }, () => h(
        NSpace, { justify: "start" }, () => [
            h(NAvatar, { round: true, size: "large" }, () => avatar_title),
            h(NSpace, { vertical: true }, () => [
                h(NText, {}, () => chat.topic ?? members),
                h(NText, { depth: "3" }, () => sub)
            ])
        ]
    ))
}

</script>