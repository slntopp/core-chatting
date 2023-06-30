import { ref } from "vue";
import { defineStore } from "pinia";

import { createPromiseClient } from "@bufbuild/connect";
import { createGrpcWebTransport } from "@bufbuild/connect-web";

import { useAppStore } from "./app";

import {
    Empty, Chat, Defaults, Users, User, Messages, Message, Event, EventType
} from "../connect/cc/cc_pb"
import {
    ChatsAPI, MessagesAPI, StreamService, UsersAPI
} from "../connect/cc/cc_connect"

export const useCcStore = defineStore('cc', () => {
    const app = useAppStore();

    const transport = createGrpcWebTransport({
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
    const streaming = createPromiseClient(StreamService, transport)

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

    function get_messages(chat: Chat): Promise<Messages> {
        return messages_c.get(chat)
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

    let msg_handler = (event: Event) => {
        console.log('Received Message Event', event)
    }

    function set_msg_handler(handler: (event: Event) => void) {
        msg_handler = (event: Event) => {
            console.log('Received Message Event', event)
            handler(event)
        }
    }

    // (async () => {
    //     console.log("Subscribing to state updates");

    //     while (true) {
    //         await new Promise((resolve) => setTimeout(resolve, 1000));
    //         if (!app.conf) continue;

    //         try {
    //             const stream = streaming.stream(new Empty())
    //             console.log("Subscribed");
    //             for await (const event of stream) {
    //                 console.debug('Received Event', event)
    //                 if(event.type >= EventType.MESSAGE_SEND) {
    //                     msg_handler(event)
    //                 }
    //             }
    //         } catch (e) {
    //             console.debug("Disconnected", e);
    //         }
    //     }
    // })();

    (async () => {
        console.log("Subscribing to state updates");

        let socket = new WebSocket('ws://localhost:8080/', [app.conf?.token ?? 'nope'])
        socket.onerror = (e) => {
            console.error('Socket Error', e)
        }
        socket.onopen = (e) => {
            console.log('Socket Open', e)
        }
        socket.onclose = (e) => {
            console.log('Socket Close', e)
        }
        socket.onmessage = (e) => {
            console.log('Socket Message', e)
        }
    })();

    return {
        users, load_me, me,

        chats, list_chats, create_chat, delete_chat,
        get_messages, send_message, update_message, delete_message,

        set_msg_handler,

        fetch_defaults, resolve
    }
})

