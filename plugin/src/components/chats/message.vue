<template>
  <render messagePlacement="left" @contextmenu="show_dropdown">
    <span v-html="content()"></span>
  </render>

  <n-dropdown
    placement="bottom-start"
    trigger="manual"
    :x="x"
    :y="y"
    :show="show"
    :options="options"
    @clickoutside="show = false"
    @select="handle_select"
  />
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, h, nextTick, ref, toRefs} from 'vue';
import {NButton, NDivider, NDropdown, NH2, NIcon, NSpace, NText, NTooltip, useThemeVars} from 'naive-ui'

import {Kind, Message, Role, User} from '../../connect/cc/cc_pb'
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
// @ts-ignore
import DOMPurify from 'dompurify'
import UserAvatar from "../ui/user_avatar.vue";
import UserItem from "../users/user_item.vue";
import {getRelativeTime} from '../../functions.ts';

interface MessageProps {
  message: Message
}

const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));
const PencilOutline = defineAsyncComponent(() => import('@vicons/ionicons5/PencilOutline'));
const TrashOutline = defineAsyncComponent(() => import('@vicons/ionicons5/TrashOutline'));
const EyeOutline = defineAsyncComponent(() => import('@vicons/ionicons5/EyeOutline'));
const ReviewOutline = defineAsyncComponent(() => import('../../assets/icons/ReviewOutline.svg'));


const props = defineProps<MessageProps>()
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
  result.push({
    label: label('Readers'),
    icon: icon(EyeOutline),
    key: 'readers',
    children: message.value!.readers.map(r => ({key: r,type: 'render', render:()=>h(UserItem,{user: store.users.get(r) as User,actions:false})}))
  })

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
  const is_sender = message.value.sender == store.me.uuid
  let style = {
    padding: '12px',
    borderRadius: theme.value.borderRadius,
    maxWidth: 'calc(100% - 45px)',
    border: `1px solid ${theme.value.borderColor}`,
    backgroundColor: (is_sender) ? 'var(--n-color-hover)' : null
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
// @ts-ignore
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

  return sanitized.replace(/^<p>/, '').replace(/<\/p>$/, '')
}

const now = ref(Date.now())
setInterval(() => now.value = Date.now(), 1000)

function timestamp() {

  let result = ''
  if (message.value.edited) {
    result = 'edited, '
  }

  result += getRelativeTime(Number(message.value.edited ? message.value.edited : message.value.sent), now.value)

  let tooltip = [
    h(NText, {}, () =>
      `Sent: ${new Date(Number(message.value.sent)).toLocaleString()}`
    )
  ]

  if (message.value.edited) {
    tooltip.push(
      h(NDivider, { style: { margin: '5px 0' } }),
      h(NText, {}, () =>
        `Edited: ${new Date(Number(message.value.edited)).toLocaleString()}`
      )
    )
  }

  return h(NTooltip, {
    placement: 'top'
  }, {
    trigger: () => h(NText, {depth: 3}, () => result),
    default: () => tooltip
  })
}

const subStyle = computed(() =>
  (window.top) ? {
    decoration: 'underline',
    cursor: 'pointer'
  } : {
    decoration: 'none',
    cursor: 'default'
  }
)

interface RenderProps {
  messagePlacement?: 'left' | 'right' | 'auto'
}

function render(props: RenderProps, { slots }: any) {
  let is_sender: boolean

  switch (props.messagePlacement) {
    case 'left':
      is_sender = false
      break
    case 'right':
      is_sender = true
      break
    case 'auto':
    default:
      is_sender = message.value.sender === store.me.uuid
  }
  let title = [
    h(
        NH2,
        {
          style: "margin: 0",
        },
        () => h(
          NText,
          {
            class: 'sub',
            onClick() {
              const { uuid } = store.me

              window.top?.postMessage({ type: 'open-user', value: { uuid } }, '*')
            }
          },
          () => sender()
        )
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
    avatar(),
    h(NSpace,
      {
        vertical: true,
        align: is_sender ? "end" : "start",
      },
      () => [
        h(NSpace, {align: 'center'}, () => title),
        (slots as any).default(),
      ]
    ),
  ]

  if (is_sender) {
    elements = elements.reverse()
  }

  return h('div', { style: {
    display: 'grid',
    gridTemplateColumns: (is_sender) ? '1fr auto' : 'auto 1fr',
    gap: '15px',
    ...container_style.value
  } }, elements);
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

.n-space span img {
  max-width: 50vw;
}

.sub {
  cursor: v-bind('subStyle.cursor');
}

.sub:hover {
  text-decoration: v-bind('subStyle.decoration');
}
</style>