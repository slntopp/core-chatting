<template>
  <n-list-item class="chat" @click="goToChat">
    <n-space :wrap-item="false" justify="start">
      <span
        class="chat__avatar-wrapper"
        @mouseenter="hovered = true"
        @mouseleave="hovered = false"
      >
        <user-avatar
          round
          size="large"
          class="chat__avatar"
          :avatar="hovered ? ' ' : members.join(' ')"
        />

        <transition name="fade">
          <n-checkbox
            v-if="hovered"
            style="position: absolute; left: 32px; top: 20px"
            :checked="selected.includes(chat.uuid)"
            @update:checked="emits('select', chat.uuid)"
            @click.stop
          />
        </transition>
      </span>

      <div v-if="!hideMessage" class="preview">
        <n-space align="center" :wrap-item="false" style="gap: 5px">
          <n-text depth="3" class="sub" @click.stop="openUser(chat.owner)">
            {{ store.users.get(chat.owner)?.title }}
          </n-text>

          <n-icon size="20" @click.stop="openChat(chat.owner)">
            <login-icon v-if="!appStore.isMobile" />
          </n-icon>
        </n-space>
        <template v-if="appStore.displayMode === 'full' && !onlyMainInfo">
          <n-text class="responsible" depth="3">
            {{ store.users.get(chat.responsible ?? "")?.title }}
          </n-text>

          <n-space
            class="department"
            style="
              flex-direction: column;
              gap: 0;
              grid-row: 1 / 3;
              grid-column: 3;
            "
            :wrap-item="false"
            :wrap="false"
          >
            <n-space
              :wrap="false"
              :wrap-item="false"
              :style="{ 'min-width': 24 * chat.gateways?.length + 'px' }"
            >
              <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
                <template #trigger>
                  <img height="24" :src="getImageUrl(gateway)" :alt="gateway" />
                </template>
                {{ gateway }}
              </n-tooltip>
            </n-space>
            <n-text italic v-if="chat.department" style="white-space: nowrap">
              {{ department }}
            </n-text>
          </n-space>

          <div class="time">
            <div>
              Created:
              <code style="margin-left: 5px">
                {{
                  Number(chat.created) > 864e5
                    ? new Date(Number(chat.created)).toLocaleDateString()
                    : getRelativeTime(Number(chat.created), now)
                }}
              </code>
            </div>
            <div v-if="lastUpdate">
              Last update:
              <code style="margin-left: 5px">
                {{ getRelativeTime(Number(lastUpdate), now, true) }}
                {{ now - lastUpdate > 6e4 ? "ago" : "" }}
              </code>
            </div>
          </div>
        </template>

        <template v-if="!onlyMainInfo">
          <n-icon
            size="18"
            @mouseenter="(e) => emits('hover', e.clientX, e.clientY, chat.uuid)"
            @mouseleave="emits('hoverEnd')"
          >
            <mail-icon />
          </n-icon>

          <n-tooltip placement="bottom">
            <template #trigger>
              <n-icon
                size="18"
                style="grid-column: -1; grid-row: 2"
                @click="copyLink"
              >
                <copy-icon />
              </n-icon>
            </template>
            Copy chat link
          </n-tooltip>

          <div class="chat__right">
            <n-tooltip>
              <template #trigger>
                <code
                  style="text-decoration: underline"
                  @click.stop="addToClipboard(uuid, notification)"
                >
                  {{ uuid.slice(0, 8).toUpperCase() }}
                </code>
              </template>
              {{ uuid }}
            </n-tooltip>

            <chat-status :chat="chat" />
          </div>
        </template>

        <n-text class="topic">{{ chatTopic }}</n-text>
      </div>

      <div style="position: absolute; right: 18px; top: -10px">
        <n-badge
          v-if="isUnreadMessages"
          :value="chat.meta!.unread"
          :max="99"
          size="24"
          :offset="[12, 12]"
        />
      </div>
    </n-space>

    <div class="mobile_right" v-if="onlyMainInfo">
      <chat-status :chat="chat" />

      <div v-if="lastUpdate">
        <code>
          {{ getRelativeTime(Number(lastUpdate), now, true) }}
          {{ now - lastUpdate > 6e4 ? "ago" : "" }}
        </code>
      </div>
    </div>
  </n-list-item>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onMounted,
  ref,
  toRefs,
  watch,
} from "vue";
import { useRouter } from "vue-router";
import {
  NBadge,
  NCheckbox,
  NIcon,
  NListItem,
  NSpace,
  NText,
  NTooltip,
  useNotification,
} from "naive-ui";
import { Chat } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting.ts";
import {
  addToClipboard,
  getImageUrl,
  getRelativeTime,
} from "../../functions.ts";
import UserAvatar from "../ui/user_avatar.vue";
import ChatStatus from "./chat_status.vue";
import { useAppStore } from "../../store/app.ts";

const CopyIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CopyOutline")
);
const MailIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChatbubbleEllipsesOutline")
);
const LoginIcon = defineAsyncComponent(
  () => import("../../assets/icons/LoginOutline.svg")
);

interface ChatItemProps {
  chat: Chat;
  uuid: string;
  hideMessage: boolean;
  selected: string[];
  chats: Chat[];
}

const props = defineProps<ChatItemProps>();
const emits = defineEmits(["hover", "hoverEnd", "select"]);
const { chat, uuid, chats } = toRefs(props);

const store = useCcStore();
const appStore = useAppStore();
const router = useRouter();
const notification = useNotification();

const users = computed(() =>
  chat.value.users.map((uuid) => store.users.get(uuid)?.title ?? "Unknown")
);
const admins = computed(() =>
  chat.value.admins.map((uuid) => store.users.get(uuid)?.title ?? "Unknown")
);

const members = computed(() => users.value.concat(admins.value));
const department = computed(
  () =>
    store.departments.find(({ key }) => key === chat.value.department)?.title ??
    chat.value.department
);

// const sub = computed(() => {
//   if (chat.value.meta && chat.value.meta.lastMessage) {
//     return chat.value.meta.lastMessage!.content.replace(/(<([^>]+)>)/ig, "")
//   }

//   return "No messages yet"
// })
const hovered = ref(false);

const chatTopic = computed(() => {
  const topic = chat.value.topic ?? members.value.join(", ");
  const words = topic.split(" ");

  if (topic.length > 512) return topic.slice(0, 509) + "...";
  else if (words.length > 16) return words.slice(0, 16).join(" ") + "...";

  return topic || "-";
});
const isUnreadMessages = computed(
  () => chat.value.meta && chat.value.meta.unread > 0
);

const onlyMainInfo = computed(
  () => appStore.isMobile || appStore.displayMode === "half"
);

const responsibleColsWidth = ref("auto");
const departmentColsWidth = ref("auto");

onMounted(setColumns);
watch(chats, setColumns);

async function setColumns() {
  if (chats.value.length < 1) return;
  await nextTick();

  const responsibles = document.querySelectorAll(".preview > .responsible");
  const departments = document.querySelectorAll(".preview > .department");
  let resWidth = 0;
  let depWidth = 0;

  responsibles.forEach((element) => {
    if (resWidth < element.clientWidth) {
      resWidth = element.clientWidth;
    }
  });

  departments.forEach((element) => {
    if (depWidth < element.clientWidth) {
      depWidth = element.clientWidth;
    }
  });

  if (resWidth > 0) {
    responsibleColsWidth.value = `${resWidth}px`;
  }
  if (depWidth > 0) {
    departmentColsWidth.value = `${depWidth}px`;
  }
}

const previewColumns = computed(() =>
  appStore.displayMode === "full"
    ? `1fr ${responsibleColsWidth.value} ${departmentColsWidth.value} 210px`
    : "1fr auto auto"
);

const subDecoration = computed(() => (window.top ? "underline" : "none"));

const chatRightColumn = computed(() =>
  appStore.displayMode === "full" ? 6 : 3
);
const now = ref(Date.now());

setInterval(() => (now.value = Date.now()), 1000);

const lastUpdate = computed(() =>
  Number(
    chat.value.meta?.lastMessage?.edited || chat.value.meta?.lastMessage?.sent
  )
);

function copyLink() {
  const url = new URL(appStore.conf?.params.fullUrl);

  url.searchParams.set("chat", chat.value.uuid);
  addToClipboard(url.href, notification);
}

function goToChat() {
  router.push({ name: "Chat", params: { uuid: uuid.value } });

  window.top?.postMessage(
    {
      type: "send-user",
      value: { uuid: chat.value.owner },
    },
    "*"
  );
}

function openUser(uuid: string) {
  window.top?.postMessage({ type: "open-user", value: { uuid } }, "*");
}

function openChat(user: string) {
  window.top?.postMessage(
    {
      type: "open-chat",
      value: { uuid: chat.value.uuid, user },
    },
    "*"
  );
}
</script>

<style scoped lang="scss">
.chat {
  .preview {
    display: grid;
    grid-template-columns: v-bind(previewColumns);
    align-items: center;
    gap: 5px 15px;
    width: calc(100% - 52px);

    @media only screen and (max-width: 1024px) {
      display: flex;
      flex-direction: column;
      align-items: start;
      gap: 0px;
    }

    .sub {
      display: block;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      width: fit-content;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: calc(100% - 25px);

      &:hover {
        text-decoration: v-bind(subDecoration);
      }
    }

    .topic {
      word-break: break-all;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      width: fit-content;
      overflow: hidden;
      text-overflow: ellipsis;
      min-width: 225px;
      max-width: calc(100% - 50px);
    }

    .responsible {
      grid-column: 2;
      grid-row: 1 / 3;
      align-self: center;
      padding-right: 10px;
      white-space: nowrap;
    }

    .time {
      grid-row: 1 / 3;
      grid-column: 4;
    }
  }

  &__right {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    grid-row: 1 / 3;
    grid-column: v-bind(chatRightColumn);
    padding: 5px 10px;
    border: 1px solid var(--n-border-color-popover);
    border-radius: 15px;
  }
}

.mobile_right {
  position: absolute;
  right: 5px;
  bottom: 0px;
  display: flex;
  align-items: end;
  flex-direction: column;
  code {
    font-size: 0.7rem;
  }
}
</style>

<style>
.n-list .n-list-item .n-list-item__main {
  overflow: hidden;
}

.chat__avatar-wrapper {
  height: fit-content;
}

.chat__avatar {
  transition: 0.3s;
}

.chat__avatar-wrapper:hover .chat__avatar {
  transform: scale(0.5);
}
</style>
