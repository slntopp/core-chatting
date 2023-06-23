<template>
    <n-layout-sider>
        <n-scrollbar style="height: 90vh">
            <n-list hoverable clickable>
                <template #header>
                    <n-space justify="center">
                        <n-button ghost type="success" @click="router.push({ name: 'Start Chat' })">
                            <template #icon>
                                <n-icon :component="ChatbubbleEllipsesOutline" />
                            </template>
                            Start Chat
                        </n-button>
                    </n-space>
                </template>
                <chat-list-item v-for="[uuid, chat] in store.chats" :uuid="uuid" :chat="chat" />
            </n-list>
        </n-scrollbar>
    </n-layout-sider>
    <n-layout-content>
        <router-view />
    </n-layout-content>
</template>

<script setup lang="ts">
import { defineAsyncComponent, h } from 'vue';
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

interface ChatListItemProps {
    uuid: string
    chat: Chat
}
function chatListItem(props: ChatListItemProps) {
    let { chat} = props;

    let users = chat.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')
    let admins = chat.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')

    let members = users.concat( admins )
    let avatar_title = members.map(el => el[0]).join(',')

    return h(NListItem, {}, () => h(
        NSpace, {justify: "start"}, () => [
            h(NAvatar, {round: true, size: "large"}, () => avatar_title),
            h(NSpace, {vertical: true}, () => [
                h(NText, {}, () => chat.topic ?? members),
                h(NText, {depth: "3"}, () => "Last message placeholder")
            ])
        ]
    ))
}

</script>