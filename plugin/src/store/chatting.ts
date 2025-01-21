import { ref } from "vue";
import { defineStore } from "pinia";

import { ConnectError, createPromiseClient } from "@connectrpc/connect";
import { createGrpcWebTransport } from "@connectrpc/connect-web";

import { useAppStore } from "./app";

import {
  Empty,
  Chat,
  Defaults,
  Users,
  User,
  Messages,
  Message,
  Event,
  EventType,
  ChatMeta,
  Kind,
  StreamRequest,
  Department,
  Merge,
  FetchDefaultsRequest,
  GetMembersRequest,
  ListChatsRequest,
  CountChatsRequest,
} from "../connect/cc/cc_pb";
import {
  ChatsAPI,
  MessagesAPI,
  StreamService,
  UsersAPI,
} from "../connect/cc/cc_connect";
import { useRoute, useRouter } from "vue-router";
import { useNotification } from "naive-ui";
import { MetricWithKey } from "../hooks/useDefaults";

export const useCcStore = defineStore("cc", () => {
  const app = useAppStore();

  let baseUrl = "/";
  if (import.meta.env.VITE_API_URL) {
    baseUrl = import.meta.env.VITE_API_URL;
  } else if (import.meta.env.DEV) {
    baseUrl = "http://localhost:8080";
  }

  const transport = createGrpcWebTransport({
    baseUrl: baseUrl,
    useBinaryFormat: true,
    interceptors: [
      (next) => async (req) => {
        req.header.set("Authorization", `Bearer ${app.conf?.token}`);
        return next(req);
      },
    ],
  });
  const chats_c = createPromiseClient(ChatsAPI, transport);
  const messages_c = createPromiseClient(MessagesAPI, transport);
  const users_c = createPromiseClient(UsersAPI, transport);
  const streaming = createPromiseClient(StreamService, transport);

  const route = useRoute();
  const router = useRouter();
  const notification = useNotification();

  const chats = ref<Map<string, Chat>>(new Map());
  const chats_count = ref<Map<string, number>>(new Map());
  const totalChats = ref<number>(0);
  const currentChat = ref<Chat | null>(null);
  const users = ref<Map<string, User>>(new Map());
  const departments = ref<Department[]>([]);
  const metrics = ref<MetricWithKey[]>([]);

  const lastChatsParams = ref<ListChatsRequest>();

  const updating = ref(false);
  const current_message = ref<Message>(
    new Message({
      chat: route.params.uuid as string,
      content: "",
    })
  );

  const me = ref<User>(new User());

  async function list_chats(data: ListChatsRequest) {
    let result = await chats_c.list(data);
    console.debug("Got Chats", result.pool);
    chats.value = new Map(
      result.pool.map((chat) => {
        if (!chat.meta) {
          chat.meta = new ChatMeta({
            unread: 0,
          });
        }

        return [chat.uuid, chat];
      })
    );

    lastChatsParams.value = data;
    totalChats.value = Number(result.total);

    resolve(result.pool.map((chat) => [...chat.admins, ...chat.users]).flat());
  }

  async function list_chats_count(data: CountChatsRequest) {
    let result = await chats_c.count(data);

    chats_count.value = new Map<string, number>(
      Object.keys(result.statuses).map((key) => {
        return [key, Number(result.statuses[+key])];
      })
    );
  }

  async function create_chat(chat: Chat) {
    chat = await chats_c.create(chat);

    chats.value.set(chat.uuid, chat);

    return chat;
  }

  async function update_chat(chat: Chat) {
    Object.entries(chat.meta?.data ?? {}).forEach(([key, value]) => {
      if (!value.kind.value) delete chat.meta?.data[key];
    });
    chat = await chats_c.update(chat);

    chats.value.set(chat.uuid, chat);

    return chat;
  }

  async function delete_chat(chat: Chat) {
    await chats_c.delete(chat);

    chats.value.delete(chat.uuid);
  }

  function fetch_defaults(fetchTemplates: boolean = false): Promise<Defaults> {
    return users_c.fetchDefaults(new FetchDefaultsRequest({ fetchTemplates }));
  }

  function update_defaults(config: Defaults): Promise<Defaults> {
    return users_c.setConfig(config);
  }

  function change_department(chat: Chat): Promise<Chat> {
    return chats_c.changeDepartment(chat);
  }

  function change_status(chat: Chat): Promise<Chat> {
    return chats_c.changeStatus(chat);
  }

  function merge_chats(chats: string[]) {
    return chats_c.mergeChats(new Merge({ chats }));
  }

  async function resolve(req_users: string[] = []): Promise<Users> {
    let result = await users_c.resolve(
      new Users({ users: req_users.map((uuid) => new User({ uuid })) })
    );
    console.debug("Got Users", result.users);

    for (let user of result.users) {
      users.value.set(user.uuid, user);
    }

    return result;
  }

  const messages = ref<Map<string, Message[]>>(new Map());

  function chat_messages(chat: Chat | undefined): Message[] {
    if (!chat) return [];
    return (messages.value.get(chat.uuid) as Message[]) || [];
  }

  async function get_messages(chat: Chat, cache = true): Promise<Messages> {
    if (
      cache &&
      messages.value.get(chat.uuid) &&
      messages.value.get(chat.uuid)!.length > 0
    )
      return new Messages({ messages: messages.value.get(chat.uuid)! });

    let res = await messages_c.get(chat);
    console.debug("Got Messages", res.messages);

    messages.value.set(chat.uuid, res.messages);
    return res;
  }
  function send_message(message: Message): Promise<Empty> {
    return messages_c.send(message);
  }
  function update_message(message: Message): Promise<Message> {
    return messages_c.update(message);
  }
  function delete_message(message: Message): Promise<Message> {
    return messages_c.delete(message);
  }

  async function handle_send(
    uuid: string,
    kind = Kind.DEFAULT,
    review = false
  ) {
    if (
      current_message.value.content == "" &&
      !current_message.value.attachments.length
    ) {
      return;
    }

    if (!current_message.value.chat) {
      current_message.value.chat = route.params.uuid as string;
    }

    try {
      current_message.value.underReview = review;
      current_message.value.kind = kind;

      if (updating.value) {
        await update_message(Message.fromJson(current_message.value as any));
        updating.value = false;
      } else {
        current_message.value.chat = route.params.uuid as string;
        await send_message(current_message.value as Message);
      }
    } catch (e) {
      if (e instanceof ConnectError) {
        notification.error({
          title: "Error",
          description: e.message,
        });
      }

      get_messages(chats.value.get(uuid) as Chat, false);
    }

    updating.value = false;
    current_message.value = new Message({
      content: "",
    });
  }

  async function load_me() {
    me.value = await users_c.me(new Empty());
  }

  function get_members(): Promise<Users> {
    return users_c.getMembers(new GetMembersRequest());
  }

  const msg_handler = (event: Event) => {
    console.log("Received Message Event", event);

    let msg: Message = event.item.value as Message;
    let idx: number;

    if (!msg.chat) return console.warn("message without chat", msg);
    const chat = currentChat.value || chats.value.get(msg.chat);
    if (!chat) return console.warn("no chat");

    if (!messages.value.get(msg.chat)) messages.value.set(msg.chat, []);

    switch (event.type) {
      case EventType.MESSAGE_SENT:
        if (messages.value.get(msg.chat)?.find(mesage => mesage.uuid === msg.uuid)) {
          break;
        }
        messages.value.get(msg.chat)!.push(msg);
        chat!.meta = new ChatMeta({
          unread: chat!.meta!.unread + 1,
          lastMessage: msg,
        });
        break;
      case EventType.MESSAGE_UPDATED:
        idx = messages.value
          .get(msg.chat)!
          .findIndex((el) => el.uuid == msg.uuid);
        messages.value.get(msg.chat)![idx] = msg;
        chat!.meta!.unread++;
        break;
      case EventType.MESSAGE_DELETED:
        idx = messages.value
          .get(msg.chat)!
          .findIndex((el) => el.uuid == msg.uuid);
        messages.value.get(msg.chat)!.splice(idx, 1);
        break;
      default:
        console.warn("unknown event type", event.type);
    }
  };

  const chat_handler = (event: Event) => {
    if (app.conf?.params?.filterByAccount) {
      return;
    }

    console.log("Received Chat Event", event);

    const chat = event.item.value as Chat;

    if (!chat.meta) {
      chat.meta = new ChatMeta({
        unread: 0,
      });
    }

    switch (event.type) {
      case EventType.CHAT_CREATED:
        const localChats = [...chats.value.values()];
        chats.value.clear()
        chats.value.set(chat.uuid, chat);
        localChats.forEach(c =>
          chats.value.set(c.uuid, c)
        )
        break;
      case EventType.CHAT_DEPARTMENT_CHANGED:
      case EventType.CHAT_STATUS_CHANGED:
      case EventType.CHAT_UPDATED:
        chats.value.set(chat.uuid, chat);
        break;
      case EventType.CHAT_DELETED:
        if (route.name == "Chat" && route.params.uuid == chat.uuid) {
          router.push({ name: "Empty Chat" });
        }

        chats.value.delete(chat.uuid);
        break;
      default:
        console.warn("unknown event type", event.type);
    }
  };

  const get_chat = async (uuid: string) => {
    currentChat.value = null;

    try {
      const data = await chats_c.get(Chat.fromJson({ uuid }));

      if (!data.meta) {
        data.meta = new ChatMeta({
          unread: 0,
        });
      }

      currentChat.value = data;
    } catch (e) {
      console.warn(e);
    }
  };

  (async () => {
    console.log("Subscribing to state updates");

    while (true) {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      if (!app.conf) continue;

      try {
        const stream = streaming.stream(new StreamRequest());
        console.log("Subscribed");
        for await (const event of stream) {
          console.debug("Received Event", event);
          if (event.type == EventType.PING) continue;
          if (
            event.type === +EventType.CHAT_DEPARTMENT_CHANGED ||
            event.type === +EventType.CHAT_STATUS_CHANGED
          ) {
            chat_handler(event);
          } else if (event.type === +EventType.CHATS_MERGED) {
            const chat = event.item.value as Chat;

            await list_chats(lastChatsParams.value!);
            await get_messages(chat, false);
            router.replace({ name: "Chat", params: { uuid: chat.uuid } });
          } else if (event.type >= EventType.MESSAGE_SENT) {
            msg_handler(event);
          } else {
            chat_handler(event);
          }
        }
      } catch (e) {
        console.debug("Disconnected", e);
      }
    }
  })();

  return {
    users,
    load_me,
    me,
    get_members,
    baseUrl,
    departments,
    metrics,

    chats,
    currentChat,
    totalChats,
    list_chats,
    create_chat,
    delete_chat,
    update_chat,
    merge_chats,
    get_chat,
    list_chats_count,
    chats_count,

    current_message,
    updating,
    handle_send,
    messages,
    chat_messages,
    get_messages,
    send_message,
    update_message,
    delete_message,

    fetch_defaults,
    update_defaults,
    change_department,
    change_status,
    resolve,
  };
});
