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
                <n-list-item>
                    {{ store.chats }}
                </n-list-item>
            </n-list>
        </n-scrollbar>
    </n-layout-sider>
    <n-layout-content>
        <router-view />
    </n-layout-content>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue';
import {
    NLayoutSider, NLayoutContent,
    NScrollbar, NList, NListItem,
    NSpace, NButton, NIcon
} from 'naive-ui';

import { useCcStore } from '../../store/chatting.ts';

import { useRouter } from 'vue-router';

const ChatbubbleEllipsesOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ChatbubbleEllipsesOutline'));

const store = useCcStore();
const router = useRouter();

async function sync() {
    await store.list_chats();
}
sync()
</script>