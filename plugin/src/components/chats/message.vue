<template>
  <render @contextmenu="show_dropdown">
    <div v-html="content()">
    </div>
  </render>

  <n-dropdown placement="bottom-start" trigger="manual" :x="x" :y="y" :options="options" :show="show"
              @clickoutside=" show = false" @select="handle_select"/>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, h, nextTick, PropType, ref, toRefs} from 'vue';
import {NButton, NDivider, NDropdown, NGi, NGrid, NH2, NIcon, NSpace, NText, NTooltip, useThemeVars,} from 'naive-ui'

import {Kind, Message, Role} from '../../connect/cc/cc_pb'
import {useCcStore} from '../../store/chatting';

import hljs from 'highlight.js';
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import css from 'highlight.js/lib/languages/css'
import xml from 'highlight.js/lib/languages/xml'
import json from 'highlight.js/lib/languages/json'
import markdown from 'highlight.js/lib/languages/markdown'

import {marked, Renderer} from 'marked'

// @ts-ignore
import {mangle} from 'marked-mangle'

import DOMPurify from 'dompurify'
import UserAvatar from "../ui/userAvatar.vue";

const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));
const PencilOutline = defineAsyncComponent(() => import('@vicons/ionicons5/PencilOutline'));
const TrashOutline = defineAsyncComponent(() => import('@vicons/ionicons5/TrashOutline'));
const ReviewOutline = defineAsyncComponent(() => import('../../assets/icons/ReviewOutline.svg'));

const props = defineProps({
  message: {
    type: Object as PropType<Message>,
    required: true
  }
})
const {message} = toRefs(props)

const emit = defineEmits(['approve', 'convert', 'edit', 'delete'])

const theme = useThemeVars()

const store = useCcStore()

function sender(): string {
  return store.users.get(message.value.sender)?.title ?? 'Unknown'
}

const x = ref(0)
const y = ref(0)
const show = ref(false)
const options = computed(() => {

  let result = []

  const label = (text: string) => () => h('b', {}, text)
  const icon = (component: any) => () => h(NIcon, {size: 24, component: component})

  if (store.chats.get(message.value.chat!)?.role == Role.ADMIN) {
    result.push({
      label: label(message.value.underReview ? 'Approve' : 'Review'), key: 'approve',
      icon: icon(ReviewOutline)
    }, {
      label: label(message.value.kind == Kind.ADMIN_ONLY ? 'Convert to Message' : 'Convert to Admin Note'),
      key: 'convert', icon: icon(ClipboardOutline)
    }, {
      type: 'divider', key: 'd-pre-delete'
    })
  }

  if (store.chats.get(message.value.chat!)?.role == Role.ADMIN || message.value.sender == store.me.uuid) {
    result.unshift({
      label: label('Edit'), key: 'edit',
      icon: icon(PencilOutline)
    }, {
      type: 'divider', key: 'd-post-edit'
    })
    result.push({
      label: label('Delete'), key: 'delete',
      icon: icon(TrashOutline)
    })
  }

  return result
})

function handle_select(key: "approve" | "convert" | "edit" | "delete") {
  console.log(key)
  show.value = false

  switch (key) {
    case 'approve':
      emit('approve', message.value.underReview)
      break
    case 'convert':
      emit('convert', message.value.kind == Kind.ADMIN_ONLY ? Kind.DEFAULT : Kind.ADMIN_ONLY)
      break
    default:
      emit(key)
  }
}

function show_dropdown(e: MouseEvent) {
  e.preventDefault()
  show.value = false

  nextTick(() => {
    show.value = true
    x.value = e.clientX
    y.value = e.clientY
  })
}


function avatar() {

  let elements = [
    h(UserAvatar, {round: true, size: 64, avatar: sender()})
  ]

  if (message.value.kind == Kind.ADMIN_ONLY) {
    elements.push(h(NIcon, {
      color: theme.value.warningColor,
      size: 24,
      component: ClipboardOutline,
    }))
  }

  return h(NSpace, {vertical: true, justify: 'start', align: 'center'}, () => elements)
}

const container_style = computed(() => {
  let style = {
    padding: '12px 0 12px 12px', borderRadius: theme.value.borderRadius,
    maxWidth: '98%', border: `1px solid ${theme.value.borderColor}`,
    backgroundColor: ''
  }

  if (message.value.underReview)
    style = {...style, backgroundColor: theme.value.infoColor + '40', border: `1px solid ${theme.value.infoColor}`}
  else if (message.value.kind == Kind.ADMIN_ONLY)
    style = {
      ...style,
      backgroundColor: theme.value.warningColor + '40',
      border: `1px solid ${theme.value.warningColor}`
    }

  return style
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

  return `<div class="code"><code>${hljs.highlight(code, {language}).value}</code></div>`
}
marked.setOptions({renderer})

function content() {
  const parsed = marked.parse(message.value.content)
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
  if (message.value.edited) {
    result = 'edited, '
  }

  result += getRelativeTime(Number(message.value.edited ? message.value.edited : message.value.sent))

  let tooltip = [
    h(NDivider, {titlePlacement: 'left'}, () => 'Sent'),
    h(NText, {}, () => new Date(Number(message.value.sent)).toString())
  ]

  if (message.value.edited) {
    tooltip.push(
        h(NDivider, {titlePlacement: 'left'}, () => 'Edited'),
        h(NText, {}, () => new Date(Number(message.value.edited)).toString())
    )
  }

  return h(NTooltip, {
    placement: 'top'
  }, {
    trigger: () => h(NText, {depth: 3}, () => result),
    default: () => tooltip
  })
}

function render(_props: any, {slots}: any) {

  const is_sender = message.value.sender == store.me.uuid

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

  if (message.value.underReview) {
    title.push(h(NButton, {
      size: 'small', type: 'info',
      ghost: true, round: true,
      onClick: () => emit('approve', true)
    }, () => 'Approve'))
  }

  if (is_sender) {
    title = title.reverse()
  }

  let elements = [
    avatar_item,
    h(NGi, {span: 21},
        () => h(
            NSpace,
            {
              vertical: true,
              align: is_sender ? "end" : "start",
            },
            () => [
              h(NSpace, {align: 'center'}, () => title),
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