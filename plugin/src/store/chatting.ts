import { ref } from "vue";
import { defineStore } from "pinia";

import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import { useAppStore } from "./app";

import {
    Empty, Chat, Defaults, Users, User, Messages, Message
} from "../connect/cc/cc_pb"
import {
    ChatsAPI, MessagesAPI, UsersAPI
} from "../connect/cc/cc_connect"

export const useCcStore = defineStore('cc', () => {
    const app = useAppStore();

    const transport = createConnectTransport({
        baseUrl: 'http://localhost:8080',
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

    const chats = ref<Map<string, Chat>>(new Map())
    const users = ref<Map<string, User>>(new Map())

    const me = ref<User>(new User())

    async function list_chats() {
        let result = await chats_c.list(new Empty())
        console.debug('Got Chats', result.chats)
        chats.value = new Map(result.chats.map((chat) => [chat.uuid, chat]))
    }

    async function create_chat(chat: Chat) {
        chat = await chats_c.create(chat)

        chats.value.set(chat.uuid, chat)

        return chat
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

    function get_messages(chat: Chat): Promise<Messages> {
        return messages_c.get(chat)
    }
    function send_message(message: Message): Promise<Empty> {
        return messages_c.send(message)
    }

    async function load_me() {
        me.value = await users_c.me(new Empty())
    }

    return {
        users, load_me, me,

        chats, list_chats, create_chat,
        get_messages, send_message,

        fetch_defaults, resolve
    }
})

