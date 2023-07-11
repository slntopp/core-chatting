<template>
  <n-list style="padding-left: 16px">
    <template #header>
      <chat-header :chat="chat!" style="height: 5vh"/>
    </template>

    <n-scrollbar style="height: 80vh; max-width: 80%;" v-if="isMessageLoading || messages.length>0" ref="scrollbar">
      <template  v-if="isMessageLoading">
          <mock-message v-for="(_,index) in 5" :key="index+ (chat?.topic || '')"/>
      </template>

      <n-list-item v-else v-for="message in messages" :key="message.uuid">
        <!-- @vue-ignore -->
        <message-view :message="message" @approve="a => handle_approve(message, a)"
                      @convert="kind => { updating = true; current_message = message; handle_send(kind, !messages['underReview']) }"
                      @delete="handle_delete(message)" @edit="() => { updating = true; current_message = message; }"/>
      </n-list-item>
    </n-scrollbar>

    <n-space v-else style="height: 80vh; max-width: 80%;" justify="center" align="center">
      <n-alert type="info" title="No Messages yet">
        Use textarea below to send a message.
      </n-alert>
    </n-space>

    <template #footer>
      <n-space style="min-height: 10vh" align="end">
        <user-avatar round size="medium" avatar="M E"></user-avatar>

        <n-space vertical justify="center">
          <n-alert style="min-width: 50vw;" title="Editing" type="info" closable v-if="updating"
                   @close="handle_stop_edit"/>
          <n-tooltip placement="top-end">
            <template #trigger>

              <n-input type="textarea" size="small" :autosize="{
                                minRows: 2,
                                maxRows: 5
                            }
                                " style="min-width: 50vw; max-width: 60vw;" placeholder="Type your message"
                       v-model:value="current_message.content" ref="input"
                       @keypress.ctrl.enter.exact="e => { e.preventDefault(); handle_send() }"
                       @keypress.ctrl.shift.enter.exact="e => { e.preventDefault(); handle_send(Kind.ADMIN_ONLY) }"
                       @keyup.ctrl.up.exact="e => { e.preventDefault(); handle_begin_edit() }"/>
            </template>

            <ul>
              <li><kbd>Ctrl</kbd> + <kbd>Enter</kbd> to send message</li>
              <li v-if="chat?.role == Role.ADMIN"><kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Enter</kbd> to
                send
                message as Admin Note
              </li>
              <li><kbd>Ctrl</kbd> + <kbd>↑</kbd> to edit last message</li>
            </ul>

          </n-tooltip>
        </n-space>

        <n-space>
          <n-button type="success" ghost circle size="small" @click="handle_send()">
            <template #icon>
              <n-icon :component="SendOutline"/>
            </template>
          </n-button>
          <n-tooltip v-if="chat?.role == Role.ADMIN">
            <template #trigger>
              <n-button type="warning" ghost circle size="small" @click="handle_send(Kind.ADMIN_ONLY)">
                <template #icon>
                  <n-icon :component="ClipboardOutline"/>
                </template>
              </n-button>
            </template>
            Message will be sent as an Admin Note, thus visible only to Admins.
          </n-tooltip>
          <n-tooltip v-if="chat?.role == Role.ADMIN">
            <template #trigger>
              <n-button type="info" quaternary circle size="small" @click="handle_send(Kind.DEFAULT, true)">
                <template #icon>
                  <n-icon :component="ReviewOutline" size="32"/>
                </template>
              </n-button>
            </template>
            Message will be put under moderation, thus not visible to user until one of the Admins approves it.
          </n-tooltip>
        </n-space>
      </n-space>
    </template>
  </n-list>
</template>

<script setup lang="ts">
import {Component, computed, defineAsyncComponent, nextTick, ref, watch} from 'vue';
import {
  InputInst,
  NAlert,
  NButton,
  NIcon,
  NInput,
  NList,
  NListItem,
  NScrollbar,
  NSpace,
  NTooltip,
  useNotification
} from 'naive-ui';

import {useRoute, useRouter} from 'vue-router';
import {useCcStore} from '../../../store/chatting';
import {Chat, Kind, Message, Role} from '../../../connect/cc/cc_pb';
import {ConnectError} from '@bufbuild/connect';
import ChatHeader from "../../../components/chats/layouts/chat_header.vue";
import UserAvatar from "../../../components/ui/user_avatar.vue";
import MockMessage from "../../../components/chats/mock_message.vue";

const SendOutline = defineAsyncComponent(() => import('@vicons/ionicons5/SendOutline'));
const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));
const ReviewOutline = defineAsyncComponent(() => import('../../../assets/icons/ReviewOutline.svg')) as Component;

const MessageView = defineAsyncComponent(() => import('../../../components/chats/message.vue'));

const route = useRoute();
const router = useRouter();

const store = useCcStore()
const scrollbar = ref()
const input = ref<InputInst>()
const isMessageLoading = ref(false)

const chat = computed(() => {
  try {
    return store.chats.get(route.params.uuid as string)
  } catch (e) {
    console.log(e)

    router.push({name: 'Empty Chat'})
  }
})

const notify = useNotification()

const messages = computed(() => {
  return store.chat_messages(chat.value! as Chat)
})

const updating = ref(false)
const current_message = ref<Message>(new Message({
  chat: route.params.uuid as string,
  content: '',
}))

async function handle_approve(msg: Message, approve: boolean) {
  msg.underReview = !approve
  await store.update_message(msg)
}

async function handle_send(kind = Kind.DEFAULT, review = false) {
  if (current_message.value.content == '') {
    return
  }

  try {
    current_message.value.underReview = review
    current_message.value.kind = kind

    if (updating.value) {
      await store.update_message(current_message.value as Message)
      updating.value = false
    } else {
      current_message.value.chat = route.params.uuid as string
      await store.send_message(current_message.value as Message)
    }
  } catch (e) {
    if (e instanceof ConnectError) {
      notify.error({
        title: 'Error',
        description: e.message,
      })
    }

    store.get_messages(chat.value! as Chat, false)
  }


  updating.value = false
  current_message.value = new Message({
    content: '',
  })
}

async function handle_delete(msg: Message) {
  await store.delete_message(msg)
}

function scrollToBottom(smooth=false) {
  setTimeout(() => {
    if (!scrollbar.value) {
      console.warn('scrollbar not ready')
      return
    }
    nextTick(() => {
      scrollbar.value.scrollTo({top: Number.MAX_SAFE_INTEGER,behavior:smooth?'smooth':undefined})
    })
  }, 0)
}

async function load_chat() {
  if (!chat.value) return
  try {
    isMessageLoading.value = true
    await Promise.all([store.resolve([...chat.value.users, ...chat.value.admins]),
      store.get_messages(chat.value as Chat)])
    scrollToBottom()
  } finally {
    isMessageLoading.value = false
  }
}

watch(chat, load_chat)
load_chat()

watch(messages, ()=>{
  if(chat.value?.meta){
    chat.value.meta.unread=0;
  }
  scrollToBottom(true)
},{deep:true})

function handle_begin_edit() {
  for (let i = messages.value.length - 1; i >= 0; i--) {
    let msg = messages.value[i]
    if (msg.sender == store.me.uuid) {
      updating.value = true
      current_message.value = msg
      nextTick(() => {
        input.value?.focus()
      })
      break
    }
  }
}

function handle_stop_edit() {
  updating.value = false
  current_message.value = new Message({
    content: '',
  })
}

</script>

<style>
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
</style>