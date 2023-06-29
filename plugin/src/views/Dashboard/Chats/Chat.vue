<template>
    <n-list style="padding-left: 16px">
        <template #header>
            <chat-header style="height: 5vh" />
        </template>

        <n-scrollbar style="height: 80vh; max-width: 80%;" v-if="messages.length > 0" ref="scrollbar">
            <n-list-item v-for="message in messages" :key="message.uuid">
                <!-- @vue-ignore -->
                <message-view :message="message" @approve="a => handle_approve(message, a)"
                    @convert="kind => { updating = true; current_message = message; handle_send(kind, !messages.underReview) }"
                    @delete="handle_delete(message)" @edit="() => { updating = true; current_message = message; }" />
            </n-list-item>
        </n-scrollbar>

        <n-space v-else style="height: 80vh; max-width: 80%;" justify="center" align="center">
            <n-alert type="info" title="No Messages yet">
                Use textarea below to send a message.
            </n-alert>
        </n-space>

        <template #footer>
            <n-space style="min-height: 10vh" align="end">
                <n-avatar round size="medium">Me</n-avatar>

                <n-space vertical justify="center">
                    <n-alert style="min-width: 50vw;" title="Editing" type="info" closable v-if="updating"
                        @close="handle_stop_edit" />
                    <n-tooltip placement="top-end">
                        <template #trigger>

                            <n-input type="textarea" size="small" :autosize="{
                                minRows: 2,
                                maxRows: 5
                            }
                                " style="min-width: 50vw; max-width: 60vw;" placeholder="Type your message"
                                v-model:value="current_message.content"
                                @keypress.ctrl.enter.exact="e => { e.preventDefault(); handle_send() }"
                                @keypress.ctrl.shift.enter.exact="e => { e.preventDefault(); handle_send(Kind.ADMIN_ONLY) }"
                                @keyup.ctrl.up.exact="e => { e.preventDefault(); handle_begin_edit() }" />
                        </template>

                        <ul>
                            <li><kbd>Ctrl</kbd> + <kbd>Enter</kbd> to send message</li>
                            <li v-if="chat?.role == Role.ADMIN"><kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Enter</kbd> to
                                send
                                message as Admin Note</li>
                            <li><kbd>Ctrl</kbd> + <kbd>↑</kbd> to edit last message</li>
                        </ul>

                    </n-tooltip>
                </n-space>

                <n-space>
                    <n-button type="success" ghost circle size="small" @click="handle_send()">
                        <template #icon>
                            <n-icon :component="SendOutline" />
                        </template>
                    </n-button>
                    <n-tooltip v-if="chat?.role == Role.ADMIN">
                        <template #trigger>
                            <n-button type="warning" ghost circle size="small" @click="handle_send(Kind.ADMIN_ONLY)">
                                <template #icon>
                                    <n-icon :component="ClipboardOutline" />
                                </template>
                            </n-button>
                        </template>
                        Message will be sent as an Admin Note, thus visible only to Admins.
                    </n-tooltip>
                    <n-tooltip v-if="chat?.role == Role.ADMIN">
                        <template #trigger>
                            <n-button type="info" quaternary circle size="small" @click="handle_send(Kind.DEFAULT, true)">
                                <template #icon>
                                    <n-icon :component="ReviewOutline" size="32" />
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
import { defineAsyncComponent, h, ref, computed, watch, nextTick } from 'vue';
import {
    NSpace, NButton, NIcon, NTooltip,
    NAvatar, NText, NInput, NAlert,
    NList, NListItem, NScrollbar, NDivider,
    useNotification
} from 'naive-ui';

import { useRoute, useRouter } from 'vue-router';
import { useCcStore } from '../../../store/chatting';
import { Chat, Message, Kind, Role, Event, EventType } from '../../../connect/cc/cc_pb';
import { ConnectError } from '@bufbuild/connect';

const SendOutline = defineAsyncComponent(() => import('@vicons/ionicons5/SendOutline'));
const ClipboardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ClipboardOutline'));
const ReviewOutline = defineAsyncComponent(() => import('../../../assets/icons/ReviewOutline.svg'));

const MessageView = defineAsyncComponent(() => import('../../../components/chats/message.vue'));

const route = useRoute();
const router = useRouter();

const store = useCcStore()
const scrollbar = ref()

const chat = computed(() => {
    return store.chats.get(route.params.uuid as string)
})

const notify = useNotification()

function chatHeader() {
    let users = chat.value!.users.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')
    let admins = chat.value!.admins.map(uuid => store.users.get(uuid)?.title ?? 'Unknown')

    let members = users.concat(admins)
    let avatar_title = members.map(el => el[0]).join(',')

    return h(
        NSpace, { justify: "start", align: 'center' }, () => [
            h(NAvatar, { round: true, size: "medium" }, () => avatar_title),
            h(NText, {}, () => chat.value!.topic ?? members),
            h(NDivider, { vertical: true }),
            h(NText, { depth: "3" }, () => `${members.length} members`),
            h(NDivider, { vertical: true }),
            h(NButton, {
                type: 'info', size: 'small',
                ghost: true, round: true,
                onClick: get_messages,
            }, () => 'Refresh'),
            h(NDivider, { vertical: true }),
            h(NButton, {
                type: 'error', size: 'small',
                ghost: true, round: true,
                onClick: async () => {
                    await store.delete_chat(chat.value! as Chat)
                    router.push({ name: 'Empty Chat' })
                },
            }, () => 'Delete')
        ]
    )
}

const messages = ref<Message[]>([])
async function get_messages() {
    let res = await store.get_messages(chat.value! as Chat)
    messages.value = res.messages
}

const updating = ref(false)
const current_message = ref<Message>(new Message({
    chat: route.params.uuid as string,
    content: '',
}))

async function handle_approve(msg: Message, approve: boolean) {
    msg.underReview = !approve
    msg = await store.update_message(msg)

    let idx = messages.value.findIndex(el => el.uuid == msg.uuid)
    messages.value.splice(idx, 1, msg)
}

async function handle_send(kind = Kind.DEFAULT, review = false) {
    if (current_message.value.content == '') {
        return
    }

    try {
        current_message.value.underReview = review
        current_message.value.kind = kind

        if (updating.value) {
            let msg = await store.update_message(current_message.value as Message)
            updating.value = false

            let idx = messages.value.findIndex(el => el.uuid == msg.uuid)
            messages.value.splice(idx, 1, msg)
        } else {
            current_message.value.chat = route.params.uuid as string
            await store.send_message(current_message.value as Message)

            get_messages()
        }
    } catch (e) {
        if (e instanceof ConnectError) {
            notify.error({
                title: 'Error',
                description: e.message,
            })
        }

        get_messages()
    }


    updating.value = false
    current_message.value = new Message({
        content: '',
    })
}
async function handle_delete(msg: Message) {
    await store.delete_message(msg)
    let idx = messages.value.findIndex(el => el.uuid == msg.uuid)
    messages.value.splice(idx, 1)
}

function scrollToBottom() {
    setTimeout(() => {
        if (!scrollbar.value) {
            console.warn('scrollbar not ready')
            return
        }

        nextTick(() => {
            scrollbar.value.scrollTo({ top: Number.MAX_SAFE_INTEGER, behavior: "smooth" })
        })
    }, 500)
}

async function load_chat() {
    store.resolve([...chat.value!.users, ...chat.value!.admins])
    await get_messages()

    store.set_msg_handler((event: Event) => {
        let msg: Message
        let idx: number

        switch (event.type) {
            case EventType.MESSAGE_SEND:
                messages.value.push(event.item.value as Message)
                break
            case EventType.MESSAGE_UPDATED:
                msg = event.item.value as Message
                idx = messages.value.findIndex(el => el.uuid == msg.uuid)
                messages.value.splice(idx, 1, msg)
                break
            case EventType.MESSAGE_DELETED:
                msg = event.item.value as Message
                idx = messages.value.findIndex(el => el.uuid == msg.uuid)
                messages.value.splice(idx, 1)
                break
            default:
                console.warn('unknown event type', event.type)
        }

        scrollToBottom()
    })
}

watch(chat, load_chat)
load_chat()

watch(messages, scrollToBottom)

function handle_begin_edit() {
    for (let i = messages.value.length - 1; i >= 0; i--) {
        let msg = messages.value[i]
        if (msg.sender == store.me.uuid) {
            updating.value = true
            current_message.value = msg
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