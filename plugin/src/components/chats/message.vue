<template>
    <n-grid :x-gap="12" :style="container_style">
        <n-gi :span="3" v-if="message.sender != store.me.uuid">
            <avatar />
        </n-gi>

        <n-gi :span="21">
            <n-space vertical :align="message.sender == store.me.uuid ? 'end' : 'start'">
                <n-h2 style="margin: 0">
                    <n-text>
                        {{ sender() }}
                    </n-text>
                </n-h2>
                <n-text>
                    {{ message.content }}
                </n-text>
            </n-space>
        </n-gi>

        <n-gi :span="2" v-if="message.sender == store.me.uuid">
            <avatar />
        </n-gi>
        <n-gi v-if="message.kind == Kind.ADMIN_ONLY">
            <n-icon :color="theme.warningColor" size="24" :component="ClipboardOutline" />
        </n-gi>
    </n-grid>
</template>

<script setup lang="ts">
import { PropType, defineAsyncComponent, h, computed } from 'vue';
import {
    NAvatar, NGrid,
    NGi, NIcon, NSpace,
    NText, NH2, useThemeVars
} from 'naive-ui'

import { Message, Kind } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting';

const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));

const { message } = defineProps({
    message: {
        type: Object as PropType<Message>,
        required: true
    }
})

const theme = useThemeVars()

const store = useCcStore()

function sender(): string {
    return store.users.get(message.sender)?.title ?? 'Unknown'
}

function avatar() {
    console.log('avatar called', sender())
    return h(NAvatar, { round: true, size: 64 }, () => h('span', { style: { fontSize: '1.5rem' } }, sender()[0]))
}

const container_style = computed(() => {
    if (message.kind == Kind.ADMIN_ONLY)
        return { maxWidth: '98%', backgroundColor: theme.value.warningColor + '40', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.warningColor}` }

    return { maxWidth: '98%', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.borderColor}` }
})
</script>