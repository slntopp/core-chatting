import { ref } from "vue";
import { defineStore } from "pinia";

import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import { useAppStore } from "./app";

import {
    Empty, Chat, Chats
} from "../connect/cc/cc_pb"
import { 
    ChatsAPI, MessagesAPI
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

    const chats = ref<Map<string, Chat>>(new Map())

    async function list_chats() {
        let result = await chats_c.list(new Empty())
        console.debug('Got Chats', result.chats)
        chats.value = new Map(result.chats.map((chat) => [chat.uuid, chat]))
    }

    async function create_chat(chat: Chat) {
        chat = await chats_c.create(chat)

        chats.value.set(chat.uuid, chat)
    }

    return {
        chats, list_chats, create_chat
    }
})

