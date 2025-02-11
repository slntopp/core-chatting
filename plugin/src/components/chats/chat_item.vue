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
          :style="{
            backgroundColor:
              chat.responsible === store.me.uuid &&
              'var(--main-app-primary-color)',
          }"
          :avatar="
            hovered || selected.includes(chat.uuid) ? ' ' : members.join(' ')
          "
        />

        <transition name="fade">
          <n-checkbox
            v-if="hovered || selected.includes(chat.uuid)"
            style="position: absolute; top: 21px"
            :style="{ left: appStore.isMobile ? '12px' : '22px' }"
            :checked="selected.includes(chat.uuid)"
            @update:checked="emits('select', chat.uuid)"
            @click.stop
          />
        </transition>
      </span>

      <div v-if="!hideMessage" class="preview">
        <n-space
          align="center"
          :wrap-item="false"
          style="gap: 5px; width: 100%"
        >
          <n-skeleton width="40%" text v-if="isUsersLoading" />
          <n-text
            v-else
            depth="3"
            class="sub"
            @click.stop="openUser(chat.owner)"
          >
            {{ allUsers.get(chat.owner)?.title }}
          </n-text>

          <n-icon v-if="chat.botState.escalated" size="20" color="red">
            <escalated-icon />
          </n-icon>

          <n-icon size="20" @click.stop="openChat(chat.owner)">
            <login-icon v-if="!appStore.isMobile" />
          </n-icon>
        </n-space>
        <template v-if="appStore.displayMode === 'full' && !onlyMainInfo">
          <n-skeleton width="40%" text v-if="isUsersLoading" />
          <n-text v-else class="responsible" depth="3">
            {{ allUsers.get(chat.responsible ?? "")?.title }}
          </n-text>

          <n-space class="department">
            <n-space
              v-if="chat.gateways?.length"
              :style="{ 'min-width': 24 * chat.gateways?.length + 'px' }"
            >
              <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
                <template #trigger>
                  <img height="24" :src="getImageUrl(gateway)" :alt="gateway" />
                </template>
                {{ gateway }}
              </n-tooltip>
            </n-space>
            <n-text italic v-if="chat.department" style="text-align: start">
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
            @mouseenter="(e:any) => emits('hover', e.clientX, e.clientY, chat.uuid)"
            @mouseleave="emits('hoverEnd')"
          >
            <mail-icon />
          </n-icon>

          <n-tooltip placement="bottom">
            <template #trigger>
              <n-icon
                size="18"
                style="grid-column: -1; grid-row: 2"
                @click.stop="copyLink"
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

      <div style="position: absolute; right: 5px; top: -20px">
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
  onUnmounted,
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
  NSkeleton,
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
import { useUsersStore } from "../../store/users.ts";
import { storeToRefs } from "pinia";

const CopyIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CopyOutline")
);
const MailIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChatbubbleEllipsesOutline")
);
const LoginIcon = defineAsyncComponent(
  () => import("../../assets/icons/LoginOutline.svg")
);

const escalatedIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/NotificationsOutline")
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
const usersStore = useUsersStore();
const { users: allUsers, isUsersLoading } = storeToRefs(usersStore);

const router = useRouter();
const notification = useNotification();

const users = computed(() =>
  chat.value.users.map((uuid) => allUsers.value.get(uuid)?.title ?? "Unknown")
);
const admins = computed(() =>
  chat.value.admins.map((uuid) => allUsers.value.get(uuid)?.title ?? "Unknown")
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

if (appStore.isMobile) hovered.value = true;

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
      resWidth = element.clientWidth + 10;
    }
  });

  departments.forEach((element) => {
    if (depWidth < element.clientWidth) {
      depWidth = element.clientWidth + 10;
    }
  });

  if (resWidth > 0) {
    responsibleColsWidth.value = resWidth > 150 ? "150px" : `${resWidth}px`;
  }
  if (depWidth > 0) {
    departmentColsWidth.value = depWidth > 200 ? "200px" : `${depWidth}px`;
  }
}

const previewColumns = computed(() =>
  appStore.displayMode === "full"
    ? `1fr ${responsibleColsWidth.value} ${departmentColsWidth.value} 210px`
    : ""
);

const subDecoration = computed(() => (window.top ? "underline" : "none"));

const chatRightColumn = computed(() =>
  appStore.displayMode === "full" ? 6 : 3
);
const avatarScale = computed(() => {
  if (appStore.isMobile) return 0;

  return props.selected.includes(chat.value.uuid) ? 0.5 : 1;
});

const chatAvatarWidth = computed(() => (appStore.isMobile ? "15px" : null));

const onlyMainInfo = computed(
  () => appStore.isMobile || appStore.displayMode === "half"
);

const interval = ref<number>();
const now = ref(Date.now());

onMounted(() => {
  interval.value = setInterval(() => (now.value = Date.now()), 1000);
});

onUnmounted(() => {
  clearInterval(interval.value);
});

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
  &.n-list-item {
    padding: 12px 10px;
  }

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
      max-width: calc(100% - 50px);

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
      max-width: calc(100% - 75px);
    }

    .responsible {
      grid-column: 2;
      grid-row: 1 / 3;
      align-self: center;
      padding-right: 10px;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .department {
      flex-direction: column;
      gap: 0;
      grid-row: 1 / 3;
      grid-column: 3;
      overflow: hidden;
      text-overflow: ellipsis;
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
    width: 130px;
    padding: 5px 10px;
    border: 1px solid var(--n-border-color-popover);
    border-radius: 15px;
  }
}

.mobile_right {
  position: absolute;
  right: 10px;
  top: 12px;
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
  width: v-bind("chatAvatarWidth");
}

.chat__avatar {
  transform: scale(v-bind(avatarScale));
  transition: 0.3s;
}

.chat__avatar-wrapper:hover .chat__avatar {
  transform: scale(0.5);
}
</style>
