<template>
  <n-list class="chat">
    <template #header>
      <chat-header :chat="chat!" style="height: 5vh"/>
    </template>

    <n-scrollbar
      ref="scrollbar"
      style="max-height: calc(100vh - 66px - 77px); margin-bottom: auto"
      v-if="isMessageLoading || messages.length>0"
    >
      <template v-if="isMessageLoading">
        <mock-message v-for="(_,index) in 5" :key="index+ (chat?.topic || '')"/>
      </template>

      <n-list-item v-else v-for="message in messages" :key="message.uuid">
        <!-- @vue-ignore -->
        <message-view
          :message="message"
          @approve="a => handle_approve(message, a)"
          @convert="kind => { store.updating = true; store.current_message = message; store.handle_send(kind, !messages['underReview']) }"
          @delete="handle_delete(message)"
          @edit="handle_edit(message)"
        />
      </n-list-item>
    </n-scrollbar>

    <n-space v-else style="height: 90vh; width: 100%;" justify="center" align="center">
      <n-alert type="info" title="No Messages yet">
        Use textarea below to send a message.
      </n-alert>
    </n-space>

    <template #footer>
      <div class="footer">
        <div class="avatar" v-if="chat?.role == Role.ADMIN">
          <n-tooltip placement="top-start">
            <template #trigger>
              <n-button
                ghost
                circle
                type="warning"
                size="small"
                @click="sendMode = (sendMode === 'admin') ? 'default' : 'admin'"
              >
                <template #icon>
                  <clipboard-icon />
                </template>
              </n-button>
            </template>
            Sent message as an Admin Note.
          </n-tooltip>

          <n-tooltip placement="top-start">
            <template #trigger>
              <n-button
                quaternary
                circle
                type="info"
                size="small"
                @click="sendMode = (sendMode === 'approve') ? 'default' : 'approve'"
              >
                <template #icon>
                  <n-icon :component="ReviewOutline" size="32"/>
                </template>
              </n-button>
            </template>
            Message will be put under moderation, thus not visible to user until one of the Admins approves it.
          </n-tooltip>
        </div>

        <n-space class="textarea" vertical justify="center">
          <n-alert
            closable
            type="info"
            title="Editing"
            style="min-width: 50vw"
            v-if="store.updating"
            @close="handle_stop_edit"
          />

          <n-tooltip placement="top-end">
            <template #trigger>

              <n-mention
                ref="input"
                size="small"
                type="textarea"
                style="width: 100%"
                placeholder="Type your message"
                v-model:value="store.current_message.content"
                :autosize="{ minRows: 3, maxRows: 15 }"
                :options="mentionsOptions"

                @blur="check_mentioned"
                @select="update_mentioned"
                @keypress.prevent.enter.exact="handle_new_line"
                @keypress.prevent.ctrl.enter.exact="store.handle_send(chat?.uuid ?? '')"
                @keypress.prevent.ctrl.shift.enter.exact="store.handle_send(chat?.uuid ?? '', Kind.ADMIN_ONLY)"
                @keyup.prevent.ctrl.up.exact="handle_begin_edit"
              />
            </template>

            <ul style="padding: 0 0 0 10px">
              <li><kbd>Ctrl</kbd> + <kbd>Enter</kbd> to send message</li>
              <li v-if="chat?.role == Role.ADMIN"><kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Enter</kbd> to
                send
                message as Admin Note
              </li>
              <li><kbd>Ctrl</kbd> + <kbd>↑</kbd> to edit last message</li>
            </ul>

          </n-tooltip>
        </n-space>

        <n-space class="actions" style="flex-flow: nowrap">
          <n-button
            ghost
            circle
            size="small"
            type="success"
            :color="(colorBySendMode as string)"
            @click="handle_send"
          >
            <template #icon>
              <n-icon :component="SendOutline" />
            </template>
          </n-button>
        </n-space>
      </div>
    </template>
  </n-list>
</template>

<script setup lang="ts">
import {Component, computed, defineAsyncComponent, nextTick, ref, watch} from 'vue';
import {
  NAlert,
  NButton,
  NIcon,
  NMention,
  NList,
  NListItem,
  NScrollbar,
  NSpace,
  NTooltip,
  MentionOption
} from 'naive-ui';

import {useRoute, useRouter} from 'vue-router';
import {useCcStore} from '../../../store/chatting';
import {Chat, Kind, Message, Role} from '../../../connect/cc/cc_pb';
import ChatHeader from "../../../components/chats/layouts/chat_header.vue";
import MockMessage from "../../../components/chats/mock_message.vue";

const SendOutline = defineAsyncComponent(() => import('@vicons/ionicons5/SendOutline'));
const ClipboardIcon = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));
const ReviewOutline = defineAsyncComponent(() => import('../../../assets/icons/ReviewOutline.svg')) as Component;

const MessageView = defineAsyncComponent(() => import('../../../components/chats/message.vue'));

const route = useRoute();
const router = useRouter();

const store = useCcStore()
const scrollbar = ref()
const input = ref()
const isMessageLoading = ref(false)

const chat = computed(() => {
  try {
    return store.chats.get(route.params.uuid as string)
  } catch (e) {
    console.log(e)

    router.push({name: 'Empty Chat'})
  }
})

const messages = computed(() => {
  const chatMessages = store.chat_messages(chat.value! as Chat)

  chatMessages.sort((a, b) => Number(a.sent - b.sent))
  return chatMessages
})

const mentionsOptions = computed(() => {
  if (!chat.value) return []
  const uuids = new Set([...chat.value.users, ...chat.value.admins])
  const result: MentionOption[] = []

  uuids.forEach((uuid) => {
    result.push({
      label: store.users.get(uuid)?.title ?? uuid,
      value: uuid
    })
  })

  return result
})

function update_mentioned(option: MentionOption) {
  if (!option.value) return
  const i = store.current_message.mentioned.indexOf(option.value)

  if (i !== -1) {
    store.current_message.mentioned.splice(i, 1)
  } else {
    store.current_message.mentioned.push(option.value)
  }
}

function check_mentioned() {
  const users = store.current_message.content.match(/\@[\w\d-]{1,}/g)

  users?.forEach((user) => {
    const mention = store.current_message.mentioned.find((uuid) => user.includes(uuid))

    if (mention) return
    store.current_message.mentioned.push(user.replace('@', '')) 
  })
}

async function handle_approve(msg: Message, approve: boolean) {
  msg.underReview = !approve
  await store.update_message(msg)
}

async function handle_delete(msg: Message) {
  await store.delete_message(msg)
}

async function scrollToBottom(smooth = false) {
  await nextTick()

  setTimeout(() => {
    if (!scrollbar.value) {
      console.warn('scrollbar not ready')
      return
    }
    
    const top = scrollbar.value.$el.nextSibling.firstChild.scrollHeight;

    scrollbar.value.scrollTo({ top, behavior: smooth ? 'smooth' : 'instant' })
  }, 300)
}

async function load_chat() {
  if (!chat.value) return
  try {
    isMessageLoading.value = true
    await Promise.all([
      store.resolve([...chat.value.users, ...chat.value.admins]),
      store.get_messages(chat.value as Chat)
    ])

    handle_stop_edit()
  } finally {
    isMessageLoading.value = false
  }
}

watch(chat, load_chat)
load_chat()

watch(messages, () => {
  if (chat.value?.meta) {
    chat.value.meta.unread = 0;
  }

  if (scrollbar.value) scrollToBottom()
}, { deep: true })

watch(isMessageLoading, scrollToBottom)
watch(scrollbar, scrollToBottom)

const sendMode = ref('default')
const textareaColors = computed(() => {
  const color = colorBySendMode.value

  if (!color) return ''
  return `
    --n-caret-color: ${color};
    --n-color-focus: ${color}1a;
    --n-border: 1px solid ${color};
    --n-border-hover: 1px solid ${color};
    --n-border-focus: 1px solid ${color};
  `
})

const colorBySendMode = computed(() => {
  switch (sendMode.value) {
    case 'admin':
      return '#f2c97d'
    case 'approve':
      return '#70c0e8'
    default:
      return null
  }
})

watch(textareaColors, (value) => {
  if (!input.value?.inputInstRef?.$el.style) return
  input.value.inputInstRef.$el.style.cssText += value
})

function handle_edit(message: Message) {
  store.updating = true
  store.current_message = message

  if (message.underReview) {
    sendMode.value = 'approve'
  } else if (message.kind === Kind.ADMIN_ONLY) {
    sendMode.value = 'admin'
  }
}

function handle_send() {
  switch (sendMode.value) {
    case 'admin':
      store.handle_send(chat.value?.uuid ?? '', Kind.ADMIN_ONLY)
      break;
    case 'approve':
      store.handle_send(chat.value?.uuid ?? '', Kind.DEFAULT, true)
      break;
    default:
      store.handle_send(chat.value?.uuid ?? '')
      break;
  }
}

function handle_begin_edit() {
  for (let i = messages.value.length - 1; i >= 0; i--) {
    let msg = messages.value[i]
    if (msg.sender == store.me.uuid) {
      store.updating = true
      store.current_message = msg
      nextTick(() => {
        input.value?.focus()
      })
      break
    }
  }
}

function handle_stop_edit() {
  store.updating = false
  store.current_message = new Message({
    content: '',
  })
}

function handle_new_line() {
  const value = store.current_message.content

  store.current_message.content = `${value}\n`
}
</script>

<style scoped lang="scss">
kbd {
  background-color: #eee;
  border-radius: 3px;
  border: 1px solid #b4b4b4;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.2), 0 2px 0 0 rgba(255, 255, 255, 0.7) inset;
  color: #333;
  display: inline-block;
  font-size: 0.85em;
  font-weight: 700;
  line-height: 1;
  padding: 2px 4px;
  white-space: nowrap;
}

textarea {
  overflow-wrap: anywhere;
}

.chat {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding-left: 16px;
}

.footer {
  display: flex;
  justify-content: center;
  align-items: center;
  width: calc(100% - 20px);
  gap: 15px;

  .avatar {
    max-width: 35px;

    *:not(:last-child) {
      margin-bottom: 10px;
    }
  }

  .textarea {
    min-width: 100px;
    width: 100%;
  }
  
  .actions {
    justify-content: space-between;
  }

  @media (max-width: 600px) {
    flex-wrap: wrap;
    .textarea {
      width: 60%;
    }
    .actions {
      margin: 20px auto;
    }
  }
}
</style>