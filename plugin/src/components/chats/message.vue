<template>
  <render
    :settingsAvatar="true"
    messagePlacement="left"
    @contextmenu="show_dropdown"
  >
    <message-content
      :message="message"
      :attach-files="attachFiles"
      @click:file="onFileClick"
    />

    <div v-if="attachFiles?.length" class="chat__files">
      <div
        @click="() => onFileClick(file)"
        class="files__preview"
        v-for="file in attachFiles"
      >
        <n-tooltip trigger="hover">
          <template #trigger>
            <img
              v-if="file.type === 'image'"
              :src="file.url"
              :alt="file.name"
            />
            <img
              class="document-icon"
              v-else-if="file.type === 'document'"
              src="/icons/file.svg"
              :alt="file.name"
            />
          </template>
          {{ file.name }}
        </n-tooltip>
      </div>
    </div>
  </render>

  <n-modal
    v-model:show="currentImage.visible"
    preset="card"
    style="max-width: 600px"
    :title="currentImage.alt"
  >
    <img style="width: 100%" :alt="currentImage.alt" :src="currentImage.src" />
  </n-modal>

  <n-dropdown
    placement="bottom-start"
    trigger="manual"
    :x="x"
    :y="y"
    :show="show"
    :options="options"
    @clickoutside="show = false"
    @select="handle_select"
  />
</template>

<script setup lang="ts">
import {
  Component,
  computed,
  defineAsyncComponent,
  h,
  nextTick,
  onMounted,
  onUnmounted,
  reactive,
  ref,
  toRefs,
} from "vue";
import {
  NButton,
  NDivider,
  NDropdown,
  NH2,
  NH3,
  NIcon,
  NModal,
  NSpace,
  NText,
  NTooltip,
  useThemeVars,
} from "naive-ui";
import messageContent from "./message/message_content.vue";

import { Kind, Message, Role, User } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { useAppStore } from "../../store/app";

import UserAvatar from "../ui/user_avatar.vue";
import UserItem from "../users/user_item.vue";
import { addToClipboard, getRelativeTime } from "../../functions.ts";
import { useDefaultsStore } from "../../store/defaults.ts";
import { storeToRefs } from "pinia";
import { useUsersStore } from "../../store/users.ts";

interface MessageProps {
  message: Message;
}

interface AttachFile {
  type: "image" | "document";
  name: string;
  url: string;
}

const CogIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CogOutline")
);
const DownloadIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/DownloadOutline")
);
const ClipboardOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/ClipboardOutline")
);
const PencilOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/PencilOutline")
);
const TrashOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/TrashOutline")
);
const EyeOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/EyeOutline")
);
const CopyOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/CopyOutline")
);
const ReviewOutline = defineAsyncComponent(
  () => import("../../assets/icons/ReviewOutline.svg")
);

const props = defineProps<MessageProps>();
const { message } = toRefs(props);

const emit = defineEmits(["approve", "convert", "edit", "delete"]);

const theme = useThemeVars();
const appStore = useAppStore();
const store = useCcStore();
const usersStore = useUsersStore();
const defaultsStore = useDefaultsStore();

const { admins } = storeToRefs(defaultsStore);
const { users } = storeToRefs(usersStore);

const sender = computed(
  () => users.value.get(message.value.sender)?.title ?? "Unknown"
);

const x = ref(0);
const y = ref(0);
const target = ref<HTMLImageElement>();
const show = ref(false);
const options = computed(() => {
  let result = [];

  const label = (text: string) => () => h("b", {}, text);
  const icon = (component: any) => () =>
    h(NIcon, { size: 24, component: component });
  if (store.chats.get(message.value.chat!)?.role == Role.ADMIN) {
    result.push(
      {
        label: label(message.value.underReview ? "Approve" : "Review"),
        key: "approve",
        icon: icon(ReviewOutline),
      },
      {
        label: label(
          message.value.kind == Kind.ADMIN_ONLY
            ? "Convert to Message"
            : "Convert to Admin Note"
        ),
        key: "convert",
        icon: icon(ClipboardOutline),
      },
      {
        type: "divider",
        key: "d-pre-delete",
      }
    );
  }

  if (
    store.chats.get(message.value.chat!)?.role == Role.ADMIN ||
    message.value.sender == store.me.uuid
  ) {
    result.unshift(
      {
        label: label("Edit"),
        key: "edit",
        icon: icon(PencilOutline),
      },
      {
        type: "divider",
        key: "d-post-edit",
      }
    );
    result.push({
      label: label("Delete"),
      key: "delete",
      icon: icon(TrashOutline),
    });
  }

  if (target.value?.tagName.toLowerCase() === "img") {
    result.push({
      label: label("Download"),
      icon: icon(DownloadIcon),
      key: "download",
    });
  }

  result.push({
    label: label("Readers"),
    icon: icon(EyeOutline),
    key: "readers",
    children: message.value!.readers.map((r) => ({
      key: r,
      type: "render",
      render: () =>
        h(UserItem, { user: users.value.get(r) as User, actions: false }),
    })),
  });

  return result;
});

function handle_select(
  key: "approve" | "convert" | "edit" | "delete" | "download"
) {
  show.value = false;

  switch (key) {
    case "approve":
      emit("approve", message.value.underReview);
      break;
    case "convert":
      emit(
        "convert",
        message.value.kind == Kind.ADMIN_ONLY ? Kind.DEFAULT : Kind.ADMIN_ONLY
      );
      break;
    case "download":
      if (!target.value) return;
      downloadFile(
        target.value.alt,
        attachFiles.value.find((file) => target.value?.alt === file.name)
          ?.url || target.value.src
      );
      break;
    default:
      emit(key);
  }
}

function show_dropdown(e: MouseEvent) {
  e.preventDefault();
  show.value = false;

  nextTick(() => {
    show.value = true;
    target.value = e.target as HTMLImageElement;
    x.value = e.clientX;
    y.value = e.clientY;
  });
}

function avatar() {
  const elements = [
    h(UserAvatar, { round: true, size: 64, avatar: sender.value }),
  ];

  if (message.value.kind == Kind.ADMIN_ONLY) {
    elements.push(
      h(NIcon, {
        color: theme.value.warningColor,
        size: 24,
        component: ClipboardOutline,
      })
    );
  }

  return h(
    NSpace,
    { vertical: true, justify: "start", align: "center" },
    () => elements
  );
}

const container_style = computed(() => {
  const is_sender = message.value.sender == store.me.uuid;
  const is_admin = admins.value.includes(message.value.sender);

  const paddingRight = appStore.isMobile ? "35px" : "45px";

  let style = {
    padding: "12px",
    borderRadius: theme.value.borderRadius,
    maxWidth: `calc(100% - ${paddingRight})`,
    border: `1px solid ${theme.value.borderColor}`,
    backgroundColor: is_sender
      ? "var(--n-color-hover)"
      : is_admin
      ? "color-mix(in hsl, var(--n-color-hover), white 10%)"
      : null,
  };

  if (message.value.underReview)
    style = {
      ...style,
      backgroundColor: theme.value.infoColor + "40",
      border: `1px solid ${theme.value.infoColor}`,
    };
  else if (message.value.kind == Kind.ADMIN_ONLY)
    style = {
      ...style,
      backgroundColor: theme.value.warningColor + "40",
      border: `1px solid ${theme.value.warningColor}`,
    };

  return style;
});

const attachFiles = computed<AttachFile[]>(() => {
  if (!message.value?.attachments) {
    return [];
  }

  const attachments: { url: string; title: string }[] =
    message.value?.attachments
      .map((uuid) => store.attachments.get(uuid))
      .filter((v) => !!v && v !== true);

  if (!attachments || !attachments.length) {
    return [];
  }

  const files: AttachFile[] = attachments?.map(({ url, title }) => {
    return {
      url: `https://${url}`,
      name: title,
      type: !!title.match(/\.(jpg|jpeg|png|gif|svg)$/i) ? "image" : "document",
    };
  });

  return files;
});

const now = ref(Date.now());
const interval = ref<number>();

onMounted(() => {
  interval.value = setInterval(() => (now.value = Date.now()), 1000);
});

onUnmounted(() => {
  clearInterval(interval.value);
});

function timestamp() {
  let result = "";
  if (message.value.edited) {
    result = "edited, ";
  }

  result += getRelativeTime(
    Number(message.value.edited ? message.value.edited : message.value.sent),
    now.value
  );

  let tooltip = [
    h(() => `Sent: ${new Date(Number(message.value.sent)).toLocaleString()}`),
  ];

  if (message.value.edited) {
    tooltip.push(
      h(NDivider, { style: { margin: "5px 0" } }),
      h(
        () =>
          `Edited: ${new Date(Number(message.value.edited)).toLocaleString()}`
      )
    );
  }

  return h(
    NTooltip,
    {
      placement: "top",
    },
    {
      trigger: () => h(NText, { depth: 3 }, () => result),
      default: () => tooltip,
    }
  );
}

function copyMessage() {
  let tooltip = [h(() => `Copy message`)];

  const copyIcon = () => h(NIcon, { size: 18, component: CopyOutline });

  return h(
    NTooltip,
    {
      placement: "right",
    },
    {
      trigger: () =>
        h(
          NButton,
          {
            circle: true,
            quaternary: true,
            onClick: () => addToClipboard(message.value.content),
          },
          () => copyIcon()
        ),
      default: () => tooltip,
    }
  );
}

const subStyle = computed(() =>
  window.top
    ? {
        decoration: "underline",
        cursor: "pointer",
      }
    : {
        decoration: "none",
        cursor: "default",
      }
);

interface RenderProps {
  messagePlacement?: "left" | "right" | "auto";
  settingsAvatar?: boolean;
}

function render(props: RenderProps, { slots }: any) {
  let is_sender: boolean;

  switch (props.messagePlacement) {
    case "left":
      is_sender = false;
      break;
    case "right":
      is_sender = true;
      break;
    case "auto":
    default:
      is_sender = message.value.sender === store.me.uuid;
  }
  let title = [
    h(appStore.isMobile ? NH3 : NH2, { style: "margin: 0" }, () =>
      h(
        NText,
        {
          class: "sub",
          onClick() {
            const uuid = message.value.sender;

            window.top?.postMessage(
              { type: "open-user", value: { uuid } },
              "*"
            );
          },
        },
        () => sender.value
      )
    ),
    h(
      NSpace,
      { wrapItem: false, wrap: false, align: "center", size: 4 },
      () => [timestamp(), copyMessage()]
    ),
  ];

  if (message.value.underReview) {
    title.push(
      h(
        NButton,
        {
          size: "small",
          type: "info",
          ghost: true,
          round: true,
          onClick: () => emit("approve", true),
        },
        () => "Approve"
      )
    );
  }

  if (is_sender) {
    title = title.reverse();
  }

  let elements = [
    props.settingsAvatar
      ? h(UserAvatar as Component, {
          id: "chat-message",
          style: {
            cursor: "pointer",
            marginTop: appStore.isMobile ? "4px" : undefined,
          },
          round: true,
          size: appStore.isMobile ? 26 : 50,
          iconSize: appStore.isMobile ? 21 : 28,
          onClick: (e: MouseEvent) => {
            let avatar = e.target as HTMLElement;

            if (avatar.id !== "chat-message") {
              avatar = avatar.closest("#chat-message") as HTMLElement;
            }
            const isOpen = JSON.parse(avatar.dataset.open ?? "false");

            if (show.value || isOpen) {
              show.value = false;
              delete avatar.dataset.open;
            } else {
              show_dropdown(e);
              avatar.dataset.open = "true";
            }
          },
          avatar: CogIcon,
        })
      : avatar(),
    h(NSpace, { align: "center", size: [12, 0] }, () => title),
    (slots as any).default(),
  ];

  if (is_sender) elements.reverse();

  return h(
    "div",
    {
      style: {
        display: "grid",
        gridTemplateColumns: is_sender ? "1fr auto" : "auto 1fr",
        gap: appStore.isMobile ? "0 10px" : "15px",
        ...container_style.value,
      },
    },
    elements
  );
}

const currentImage = reactive({ src: "", alt: "", visible: false });
function onFileClick(file: AttachFile) {
  if (file.type === "image") {
    currentImage.src = file.url;
    currentImage.alt = file.name;
    currentImage.visible = true;
  } else {
    downloadFile(file.name, file.url);
  }
}

async function downloadFile(name: string, link: string) {
  const element = document.createElement("a");

  const response = await fetch(link);
  const blob = await response.blob();
  const url = URL.createObjectURL(blob);

  element.setAttribute("href", url);
  element.setAttribute("download", name);
  element.click();
}
</script>

<style>
.n-space span img {
  max-width: 50vw;
}

.sub {
  cursor: v-bind("subStyle.cursor");
}

.sub:hover {
  text-decoration: v-bind("subStyle.decoration");
}

.chat__files {
  display: flex;
  justify-self: start;
  flex-wrap: wrap;
  grid-column: 1 / 3;
}

.chat__files .files__preview {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 114px;
  height: 100px;
  padding: 5px;
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
}

.chat__files .files__preview > img {
  height: 100%;
  width: auto;
  max-width: 100%;
  object-fit: cover;
}

.chat__files .files__preview .document-icon {
  -webkit-filter: invert(0.75); /* safari 6.0 - 9.0 */
  filter: invert(0.75);
}
</style>
