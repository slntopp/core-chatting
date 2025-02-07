<template>
  <div class="footer">
    <n-space
      v-if="chat?.role == Role.ADMIN"
      style="grid-column: 1 / 3; gap: 5px"
      :wrap-item="false"
    >
      <n-tooltip placement="top-start">
        <template #trigger>
          <n-button
            ghost
            size="small"
            :type="sendMode === 'admin' ? 'warning' : undefined"
            @click="sendMode = sendMode === 'admin' ? 'default' : 'admin'"
          >
            Notes
          </n-button>
        </template>
        Sent message as an Admin Note.
      </n-tooltip>

      <n-tooltip placement="top-start">
        <template #trigger>
          <n-button
            ghost
            size="small"
            :type="sendMode === 'approve' ? 'info' : undefined"
            @click="sendMode = sendMode === 'approve' ? 'default' : 'approve'"
          >
            Premoderate
          </n-button>
        </template>
        Message will be put under moderation, thus not visible to user until one
        of the Admins approves it.
      </n-tooltip>

      <n-tooltip placement="top-start">
        <template #trigger>
          <n-button
            ghost
            size="small"
            @click="handle_open_templates"
            :loading="isTemplatesLoading"
          >
            Templates
          </n-button>
        </template>
        Message templates created by our team. Ready to use.
      </n-tooltip>
    </n-space>

    <n-upload abstract v-model:file-list="fileList" list-type="image">
      <n-space class="textarea" vertical justify="center">
        <n-alert
          closable
          type="info"
          title="Editing"
          style="min-width: 50vw"
          v-if="store.updating"
          @close="handle_stop_edit"
        />

        <n-tooltip placement="top-end">
          <template #trigger>
            <!-- :value="get_mention_value(store.current_message.content)" -->
            <!-- @update:value="set_mention_value($event)" -->
            <n-mention
              ref="input"
              size="small"
              type="textarea"
              style="width: 100%"
              placeholder="Type your message"
              v-model:value="store.current_message.content"
              :autosize="{ minRows: 3, maxRows: 15 }"
              :options="mentionsOptions"
              :prefix="isMentionVisible ? '@' : ' @'"
              @input="onInput"
              @blur="check_mentioned"
              @keypress.prevent.enter.exact="handle_new_line"
              @keyup.prevent.ctrl.enter.exact="handle_send"
              @keyup.prevent.ctrl.shift.enter.exact="
                sendMode = 'admin';
                handle_send();
              "
              @keyup.prevent.ctrl.up.exact="handle_begin_edit"
            />
          </template>

          <ul v-if="!appStore.isMobile" style="padding: 0 0 0 10px">
            <li><kbd>Ctrl</kbd> + <kbd>Enter</kbd> to send message</li>
            <li v-if="chat?.role == Role.ADMIN">
              <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Enter</kbd> to send
              message as Admin Note
            </li>
            <li><kbd>Ctrl</kbd> + <kbd>↑</kbd> to edit last message</li>
          </ul>
        </n-tooltip>
      </n-space>

      <n-space class="actions" style="flex-flow: nowrap">
        <n-upload-trigger abstract #="{ handleClick }">
          <n-button ghost circle size="small" type="info" @click="handleClick">
            <template #icon>
              <n-icon :component="DownloadIcon" />
            </template>
          </n-button>
        </n-upload-trigger>

        <n-button
          ghost
          circle
          size="small"
          type="success"
          :color="(colorBySendMode as string)"
          @click="handle_send"
        >
          <template #icon>
            <n-icon :component="SendIcon" />
          </template>
        </n-button>
      </n-space>

      <n-upload-file-list v-if="fileList.length > 0" />
    </n-upload>

    <n-modal v-model:show="isTemplatesOpen">
      <n-card style="max-width: 800px; width: 90vw; min-height: 80vh">
        <templates-view
          @select="handle_select_template"
          v-bind="templatesOptions"
        />
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, nextTick, ref, watch } from "vue";
import {
  NAlert,
  NButton,
  NIcon,
  NMention,
  NSpace,
  NTooltip,
  NUpload,
  NUploadFileList,
  NUploadTrigger,
  NModal,
  NCard,
  MentionOption,
  UploadFileInfo,
} from "naive-ui";

import { useCcStore } from "../../../store/chatting";
import { Chat, Kind, Message, Role } from "../../../connect/cc/cc_pb";
import { useAppStore } from "../../../store/app";
import templatesView from "../../settings/templates.vue";
import { useDefaultsStore } from "../../../store/defaults";
import { storeToRefs } from "pinia";

const SendIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SendOutline")
);
const DownloadIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/DownloadOutline")
);

interface ChatFooterProps {
  chat: Chat;
  messages: Message[];
}
const props = defineProps<ChatFooterProps>();
const store = useCcStore();
const appStore = useAppStore();
const defaultsStore = useDefaultsStore();
const { admins, templates, users } = storeToRefs(defaultsStore);

interface FileInfo extends UploadFileInfo {
  uuid?: string;
}

const input = ref();
const sendMode = ref("default");
const isMentionVisible = ref(false);
const fileList = ref<FileInfo[]>([]);

const templatesOptions = ref<any>(null);
const isTemplatesOpen = ref(false);
const isTemplatesLoading = ref(false);

const mentionsOptions = computed(() => {
  if (!props.chat) return [];
  const uuids = new Set([...props.chat.users, ...props.chat.admins]);
  const result: MentionOption[] = [];

  uuids.forEach((uuid) => {
    const { title } = store.users.get(uuid) ?? {};

    result.push({
      label: title ?? uuid,
      value: title?.replace(" ", "_") ?? uuid,
    });
  });

  return result;
});

const textareaColors = computed(() => {
  const color = colorBySendMode.value;

  return `
    --n-caret-color: ${color ?? "#63e2b7"};
    --n-color-focus: ${color ?? "#63e2b7"}1a;
    --n-border: 1px solid ${color};
    --n-border-hover: 1px solid ${color ?? "#63e2b7"};
    --n-border-focus: 1px solid ${color ?? "#63e2b7"};
  `;
});

const colorBySendMode = computed(() => {
  switch (sendMode.value) {
    case "admin":
      return "#f2c97d";
    case "approve":
      return "#70c0e8";
    default:
      return null;
  }
});

watch(textareaColors, async (value) => {
  await nextTick();
  const element = input.value?.$el.querySelector(".n-input");

  if (!element?.style) return;
  element.style.cssText += value;
});

watch(
  () => store.updating,
  (value) => {
    if (value) {
      const message = store.current_message;

      const attachments: { url: string; title: string; uuid: string }[] =
        message.attachments
          .map((uuid) =>
            store.attachments.has(uuid)
              ? { ...store.attachments.get(uuid), uuid }
              : null
          )
          .filter((v) => !!v && v !== true);
      if (!attachments || !attachments.length) {
        return;
      }

      attachments.forEach(({ uuid, title, url }) => {
        fileList.value.push({
          thumbnailUrl: `https://${url}`,
          name: title,
          uuid,
          id: uuid ?? "",
          status: "finished",
        });
      });
    } else {
      fileList.value = [];
    }
  }
);

function onInput({ target }: { target: HTMLTextAreaElement }) {
  if (target.selectionStart !== target.selectionEnd) return;
  const symbol = target.value[target.selectionEnd - 1];
  const prevSymbol = target.value[target.selectionEnd - 2];

  if (symbol === "@" && prevSymbol === " ") {
    isMentionVisible.value = true;
  } else {
    isMentionVisible.value = false;
  }
}

function check_mentioned() {
  const users = Array.from(store.users.values());
  const titles = store.current_message.content.match(/@[\w\d-]{1,}/g);

  store.current_message.mentioned = [];
  titles?.forEach((title) => {
    const { uuid } =
      users.find(
        (user) => user.title.replace(" ", "_") === title.replace("@", "")
      ) ?? {};
    const mention = store.current_message.mentioned.includes(uuid ?? title);

    if (mention) return;
    store.current_message.mentioned.push(uuid ?? title);
  });
}

async function handle_send() {
  if (fileList.value.length > 0) {
    await handle_send_files();

    store.current_message.attachments = fileList.value.map(
      (f) => f.uuid
    ) as string[];
  }

  switch (sendMode.value) {
    case "admin":
      await store.handle_send(props.chat?.uuid ?? "", Kind.ADMIN_ONLY);
      break;
    case "approve":
      await store.handle_send(props.chat?.uuid ?? "", Kind.DEFAULT, true);
      break;
    default:
      await store.handle_send(props.chat?.uuid ?? "");
      break;
  }

  fileList.value = [];
}

async function handle_send_files() {
  for await (const file of fileList.value) {
    if (file.uuid) continue;
    const baseUrl = store.baseUrl.endsWith("/")
      ? store.baseUrl.slice(0, -1)
      : store.baseUrl;

    const response = await fetch(`${baseUrl}/attachments`, {
      method: "PUT",
      body: JSON.stringify({ title: file.name, chat: props.chat?.uuid }),
    });
    const { uuid, object_id: objectId } = await response.json();
    await fetch(objectId, { method: "PUT", body: file.file });

    const responseGet = await fetch(`${baseUrl}/attachments/${uuid}`);
    const { url } = await responseGet.json();

    file.url = `https://${url}`;
    file.uuid = uuid;
  }
}

function handle_begin_edit() {
  for (let i = props.messages.length - 1; i >= 0; i--) {
    const msg = props.messages[i];

    if (msg.sender == store.me.uuid) {
      store.updating = true;
      store.current_message = JSON.parse(JSON.stringify(msg));

      nextTick(() => {
        input.value?.focus();
      });
      break;
    }
  }
}

function handle_stop_edit() {
  store.updating = false;
  store.current_message = new Message({
    content: "",
  });
}

function handle_new_line() {
  const textarea = input.value.$el.querySelector("textarea");

  textarea.setRangeText(
    "\n",
    textarea.selectionStart,
    textarea.selectionEnd,
    "end"
  );
}

async function handle_open_templates() {
  isTemplatesLoading.value = true;
  try {
    if (!templatesOptions.value) {
      templatesOptions.value = {
        admins: admins.value,
        users: users.value,
        templates: templates.value,
        isEdit: false,
      };
    }

    isTemplatesOpen.value = true;
  } finally {
    isTemplatesLoading.value = false;
  }
}

function handle_select_template(content: string) {
  store.current_message.content = content;
  isTemplatesOpen.value = false;
}

const deleteIconDisplayStyle = computed(() =>
  store.updating ? "none" : "flex"
);

defineExpose({ sendMode });
</script>

<style scoped lang="scss">
kbd {
  background-color: #eee;
  border-radius: 3px;
  border: 1px solid #b4b4b4;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.2),
    0 2px 0 0 rgba(255, 255, 255, 0.7) inset;
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

.footer {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  width: calc(100% - 20px);
  gap: 5px 15px;

  .textarea {
    min-width: 100px;
    width: 100%;
  }

  .actions {
    justify-content: space-between;
  }

  :deep(
      .n-upload-file-list
        .n-upload-file
        .n-upload-file-info
        .n-upload-file-info__action
    ) {
    display: v-bind("deleteIconDisplayStyle");
  }

  @media (max-width: 600px) {
    flex-wrap: wrap;
    .textarea {
      width: 100%;
    }
    .actions {
      display: flex;
      flex-direction: column !important;
    }
  }
}
</style>
