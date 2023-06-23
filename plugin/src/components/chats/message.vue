<template>
    <n-grid :x-gap="12" :style="container_style">
        <n-gi :span="3" v-if="message.sender != store.me.uuid">
            <avatar />
        </n-gi>

        <n-gi :span="20">
            <n-space vertical :align="message.sender == store.me.uuid ? 'end' : 'start'">
                <n-h2 style="margin: 0">
                    <n-text>
                        {{ sender() }}
                    </n-text>
                </n-h2>

                <div v-html="content()">
                </div>

                <!-- <n-p v-for="line in message.content.split('\n')">
                    {{ line }}
                </n-p> -->
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
    NAvatar, NGrid, NCode,
    NGi, NIcon, NSpace,
    NText, NH2, useThemeVars
} from 'naive-ui'

import { Message, Kind } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting';

import hljs from 'highlight.js';
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import css from 'highlight.js/lib/languages/css'
import xml from 'highlight.js/lib/languages/xml'
import json from 'highlight.js/lib/languages/json'
import markdown from 'highlight.js/lib/languages/markdown'

import { marked, Renderer } from 'marked'

// @ts-ignore
import { mangle } from 'marked-mangle'

import DOMPurify from 'dompurify'

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
    return h(NAvatar, { round: true, size: 64 }, () => h('span', { style: { fontSize: '1.5rem' } }, sender()[0]))
}

const container_style = computed(() => {
    if (message.kind == Kind.ADMIN_ONLY)
        return { maxWidth: '98%', backgroundColor: theme.value.warningColor + '40', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.warningColor}` }

    return { maxWidth: '98%', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.borderColor}` }
})

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('css', css)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('json', json)
hljs.registerLanguage('markdown', markdown)

marked.use({
    async: false, gfm: true,
    breaks: true, mangle: false,
    headerIds: false, headerPrefix: '',
})
marked.use(mangle)

const renderer = new Renderer()
renderer.code = (code, language) => {
    if (!language) language = 'plaintext'

    return `<div class="code"><code>${hljs.highlight(code, { language }).value}</code></div>`
}
marked.setOptions({ renderer })

function content() {
    const parsed = marked.parse(message.content)
    const sanitized = DOMPurify.sanitize(parsed)

    return sanitized
}
</script>

<style>
div.code {
    padding: 8px;
    background-color: black;
    border-radius: 6px;
    white-space: pre;
}
</style>