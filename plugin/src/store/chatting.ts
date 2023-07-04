import { ref } from "vue";
import { defineStore } from "pinia";

import { createPromiseClient } from "@bufbuild/connect";
import { createGrpcWebTransport } from "@bufbuild/connect-web";

import { useAppStore } from "./app";

import {
    Empty, Chat, Defaults, Users, User, Messages, Message, Event, EventType, ChatMeta
} from "../connect/cc/cc_pb"
import {
    ChatsAPI, MessagesAPI, StreamService, UsersAPI
} from "../connect/cc/cc_connect"
import { useRoute, useRouter } from "vue-router";

export const useCcStore = defineStore('cc', () => {
    const app = useAppStore();

    let baseUrl = '/'
    if (import.meta.env.VITE_API_URL) {
        baseUrl = import.meta.env.VITE_API_URL
    } else if (import.meta.env.DEV) {
        baseUrl = 'http://localhost:8080'
    }

    const transport = createGrpcWebTransport({
        baseUrl: baseUrl,
        interceptors: [
            (next) => async (req) => {
                req.header.set("Authorization", `Bearer ${app.conf?.token}`);
                return next(req);
            },
        ],
    })
    const chats_c = createPromiseClient(ChatsAPI, transport);
    const messages_c = createPromiseClient(MessagesAPI, transport);
    const users_c = createPromiseClient(UsersAPI, transport);
    const streaming = createPromiseClient(StreamService, transport)

    const route = useRoute()
    const router = useRouter()

    const chats = ref<Map<string, Chat>>(new Map())
    const users = ref<Map<string, User>>(new Map())

    const me = ref<User>(new User())

    async function list_chats() {
        let result = await chats_c.list(new Empty())
        console.debug('Got Chats', result.chats)
        chats.value = new Map(result.chats.map((chat) => {
            if (!chat.meta) {
                chat.meta = new ChatMeta({
                    unread: 0,
                })
            }

            return [chat.uuid, chat]
        }))
    }

    async function create_chat(chat: Chat) {
        chat = await chats_c.create(chat)

        chats.value.set(chat.uuid, chat)

        return chat
    }

    async function update_chat(chat:Chat){
        chat = await chats_c.update(chat)

        chats.value.set(chat.uuid, chat)

        return chat
    }

    async function delete_chat(chat: Chat) {
        await chats_c.delete(chat)

        chats.value.delete(chat.uuid)
    }

    function fetch_defaults(): Promise<Defaults> {
        return users_c.fetchDefaults(new Empty())
    }

    async function resolve(req_users: string[] = []): Promise<Users> {
        let result = await users_c.resolve(new Users({ users: req_users.map((uuid) => new User({ uuid })) }))
        console.debug('Got Users', result.users)

        for (let user of result.users) {
            users.value.set(user.uuid, user)
        }

        return result
    }


    const messages = ref<Map<string, Message[]>>(new Map())

    function chat_messages(chat: Chat | undefined): Message[] {
        if (!chat) return []
        return messages.value.get(chat.uuid) || []
    }

    async function get_messages(chat: Chat, cache = true): Promise<Messages> {
        if (cache && messages.value.get(chat.uuid) && messages.value.get(chat.uuid)!.length > 0) return new Messages({ messages: messages.value.get(chat.uuid)! })

        let res = await messages_c.get(chat)
        console.debug('Got Messages', res.messages)

        messages.value.set(chat.uuid, res.messages)
        return res
    }
    function send_message(message: Message): Promise<Empty> {
        return messages_c.send(message)
    }
    function update_message(message: Message): Promise<Message> {
        return messages_c.update(message)
    }
    function delete_message(message: Message): Promise<Message> {
        return messages_c.delete(message)
    }

    async function load_me() {
        me.value = await users_c.me(new Empty())
    }

    const msg_handler = (event: Event) => {
        console.log('Received Message Event', event)

        let msg: Message = event.item.value as Message
        let idx: number

        if (!msg.chat) return console.warn('message without chat', msg)

        if (!messages.value.get(msg.chat)) messages.value.set(msg.chat, [])

        switch (event.type) {
            case EventType.MESSAGE_SEND:
                messages.value.get(msg.chat)!.push(msg)
                chats.value.get(msg.chat)!.meta = new ChatMeta({
                    unread: chats.value.get(msg.chat)!.meta!.unread + 1,
                    lastMessage: msg,
                })
                break
            case EventType.MESSAGE_UPDATED:
                idx = messages.value.get(msg.chat)!.findIndex(el => el.uuid == msg.uuid)
                messages.value.get(msg.chat)![idx] = msg
                chats.value.get(msg.chat)!.meta!.unread++
                break
            case EventType.MESSAGE_DELETED:
                idx = messages.value.get(msg.chat)!.findIndex(el => el.uuid == msg.uuid)
                messages.value.get(msg.chat)!.splice(idx, 1)
                break
            default:
                console.warn('unknown event type', event.type)
        }
    }

    const chat_handler = (event: Event) => {
        console.log('Received Chat Event', event)

        let chat: Chat = event.item.value as Chat

        if (!chat.meta) {
            chat.meta = new ChatMeta({
                unread: 0,
            })
        }

        switch (event.type) {
            case EventType.CHAT_CREATED:
                chats.value.set(chat.uuid, chat)
                break
            case EventType.CHAT_UPDATED:
                chats.value.set(chat.uuid, chat)
                break
            case EventType.CHAT_DELETED:
                if (route.name == 'Chat' && route.params.uuid == chat.uuid) {
                    router.push({ name: 'Empty Chat' })
                }

                chats.value.delete(chat.uuid)
                break
            default:
                console.warn('unknown event type', event.type)
        }
    }

    (async () => {
        console.log("Subscribing to state updates");

        while (true) {
            await new Promise((resolve) => setTimeout(resolve, 1000));
            if (!app.conf) continue;

            try {
                const stream = streaming.stream(new Empty())
                console.log("Subscribed");
                for await (const event of stream) {
                    console.debug('Received Event', event)
                    if (event.type == EventType.PING) continue
                    else if (event.type >= EventType.MESSAGE_SEND) {
                        msg_handler(event)
                    } else {
                        chat_handler(event)
                    }
                }
            } catch (e) {
                console.debug("Disconnected", e);
            }
        }
    })();

    return {
        users, load_me, me,

        chats, list_chats, create_chat, delete_chat, update_chat,

        messages, chat_messages, get_messages,
        send_message, update_message, delete_message,

        fetch_defaults, resolve
    }
})

