<template>
    <render>
        <div v-html="content()">
        </div>
    </render>
</template>

<script setup lang="ts">
import { PropType, defineAsyncComponent, h, computed, ref } from 'vue';
import {
    NAvatar, NGrid,
    NGi, NIcon, NSpace,
    NText, NH2, useThemeVars, NTooltip, NDivider
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

    let elements = [
        h(NAvatar, { round: true, size: 64 }, () => h('span', { style: { fontSize: '1.5rem' } }, sender()[0]))
    ]

    if (message.kind == Kind.ADMIN_ONLY) {
        elements.push(h(
            NGi,
            {},
            () => h(NIcon, {
                color: theme.value.warningColor,
                size: 24,
                component: ClipboardOutline,
            }),
        ))
    }

    return h(NSpace, { vertical: true, justify: 'start', align: 'center' }, () => elements)
}

const container_style = computed(() => {
    if (message.kind == Kind.ADMIN_ONLY)
        return { maxWidth: '98%', backgroundColor: theme.value.warningColor + '40', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.warningColor}` }

    return { maxWidth: '98%', padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius, border: `1px solid ${theme.value.borderColor}` }
})

{
    hljs.registerLanguage('javascript', javascript)
    hljs.registerLanguage('typescript', typescript)
    hljs.registerLanguage('css', css)
    hljs.registerLanguage('xml', xml)
    hljs.registerLanguage('json', json)
    hljs.registerLanguage('markdown', markdown)
}

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

const now = ref(new Date())
setInterval(() => now.value = new Date(), 1000)

function getRelativeTime(timestamp: number) {
    const timeDifference = (now.value.getTime() - timestamp) / 1000;
    const minutesDifference = Math.floor(timeDifference / 60);

    if (minutesDifference >= 4320) {
        return new Date(timestamp).toLocaleDateString();
    } else if (minutesDifference >= 1440) {
        const daysDifference = Math.floor(minutesDifference / 1440);
        return `${daysDifference} days ago`;
    } else if (minutesDifference >= 60) {
        const hoursDifference = Math.floor(minutesDifference / 60);
        return `${hoursDifference} hours ago`;
    } else if (minutesDifference > 0) {
        return `${minutesDifference} minutes ago`;
    } else {
        return 'just now';
    }
}


function timestamp() {

    let result = ''
    if (message.edited) {
        result = 'edited, '
    }

    result += getRelativeTime(Number(message.edited ? message.edited : message.sent))

    let tooltip = [
        h(NDivider, { titlePlacement: 'left' }, () => 'Sent'),
        h(NText, {}, () => new Date(Number(message.sent)).toString())
    ]

    if (message.edited) {
        tooltip.push(
            h(NDivider, { titlePlacement: 'left' }, () => 'Edited'),
            h(NText, {}, () => new Date(Number(message.edited)).toString())
        )
    }

    return h(NTooltip, {
        placement: 'top'
    }, {
        trigger: () => h(NText, { depth: 3 }, () => result),
        default: () => tooltip
    })
}

function render(_props: any, { slots }: any) {

    const is_sender = message.sender == store.me.uuid

    const avatar_item = h(NGi, {
        span: 3,
    },
        () => avatar()
    )


    let title = [
        h(
            NH2,
            {
                style: "margin: 0",
            },
            () => h(NText, {}, () => sender())
        ),
        timestamp()
    ]

    if (is_sender) {
        title = title.reverse()
    }

    let elements = [
        avatar_item,
        h(NGi, { span: 21 },
            () => h(
                NSpace,
                {
                    vertical: true,
                    align: is_sender ? "end" : "start",
                },
                () => [
                    h(NSpace, { align: 'center' }, () => title),
                    slots.default(),
                ]
            )
        ),
    ]

    if (is_sender) {
        elements = elements.reverse()
    }

    return h(
        NGrid,
        {
            xGap: 12,
            style: container_style.value,
        },
        () => elements
    );
}
</script>

<style>
div.code {
    padding: 8px;
    background-color: black;
    border-radius: 6px;
    white-space: pre-wrap;
    overflow-wrap: anywhere;
}
</style>