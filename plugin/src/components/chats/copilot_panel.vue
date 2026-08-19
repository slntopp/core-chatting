<template>
  <!-- Deliberately the same chat as the main one: same n-list shell, same
       message component, same composer shape. Only the message lane and the
       header differ. -->
  <n-list class="chat">
    <template #header>
      <div class="head">
        <n-icon size="18" :color="COPILOT_COLOR"><sparkles-icon /></n-icon>
        <div class="head__titles">
          <div class="head__title">Ask the bot</div>
          <div class="head__sub">Operator-only — the customer never sees this</div>
        </div>

        <n-tooltip>
          <template #trigger>
            <n-button type="default" size="small" ghost circle @click="isTracesOpen = true">
              <template #icon><bug-icon /></template>
            </n-button>
          </template>
          Copilot debug traces
        </n-tooltip>

        <n-tooltip>
          <template #trigger>
            <n-button
              type="default"
              size="small"
              ghost
              circle
              @click="appStore.isCopilotOpen = false"
            >
              <template #icon><close-icon /></template>
            </n-button>
          </template>
          Close
        </n-tooltip>
      </div>
    </template>

    <n-scrollbar ref="scrollbar" class="chat__scroll" v-if="messages.length > 0">
      <n-list-item v-for="message in messages" :key="message.uuid">
        <!-- Notes render ONLY here now, so this pane has to carry every action
             the message menu offers. Wiring just @delete left "Edit" and
             "Convert to Message" dead app-wide - and converting a note into a
             customer-visible reply is the escape hatch for a misfiled message. -->
        <message-view
          :message="message"
          @delete="store.delete_message(message)"
          @edit="startEdit(message)"
          @convert="
            (kind) => {
              store.updating = true;
              store.current_message = message;
              store.handle_send(chat.uuid, kind);
            }
          "
        />
      </n-list-item>
    </n-scrollbar>

    <n-space
      v-else
      class="chat__scroll"
      style="width: 100%"
      justify="center"
      align="center"
    >
      <n-alert type="info" title="No Notes yet">
        Ask the bot about this chat, or leave a note for the team.
      </n-alert>
    </n-space>

    <template #footer>
      <div class="footer">
        <!-- Editing has to be escapable: without a way out, the composer stays
             bound to the note forever and the next thing typed silently
             overwrites it instead of asking the bot. -->
        <div v-if="editing" class="editing">
          <span>Editing a note</span>
          <n-button text size="tiny" @click="cancelEdit">cancel</n-button>
        </div>

        <n-space class="textarea" vertical justify="center">
          <n-tooltip placement="top-end">
            <template #trigger>
              <n-input
                v-model:value="draft"
                type="textarea"
                size="small"
                style="width: 100%"
                :placeholder="
                  editing
                    ? 'Edit the note'
                    : 'Note for the team — or ask the bot about this chat'
                "
                :autosize="{ minRows: 3, maxRows: 15 }"
                :disabled="sending"
                @keyup.prevent.ctrl.enter.exact="send"
              />
            </template>
            <ul v-if="!appStore.isMobile" style="padding: 0 0 0 10px">
              <li><kbd>Ctrl</kbd> + <kbd>Enter</kbd> to send note</li>
            </ul>
          </n-tooltip>
        </n-space>

        <n-space class="actions" style="flex-flow: nowrap">
          <n-button
            ghost
            circle
            size="small"
            :color="COPILOT_COLOR"
            :loading="sending"
            :disabled="!draft.trim()"
            @click="send"
          >
            <template #icon>
              <n-icon :component="SendIcon" />
            </template>
          </n-button>
        </n-space>
      </div>
    </template>

    <n-drawer v-model:show="isTracesOpen" :width="620" placement="right">
      <n-drawer-content title="Copilot debug traces" closable :native-scrollbar="false">
        <trace-viewer v-if="isTracesOpen" :chat-uuid="chat.uuid" lane="copilot" />
      </n-drawer-content>
    </n-drawer>
  </n-list>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, nextTick, ref, watch } from "vue";
import {
  NAlert,
  NButton,
  NDrawer,
  NDrawerContent,
  NIcon,
  NInput,
  NList,
  NListItem,
  NScrollbar,
  NSpace,
  NTooltip,
  useNotification,
} from "naive-ui";

import { Chat, Kind, Message } from "../../connect/cc/cc_pb";
import { splitChatMessages } from "../../functions";
import { useCcStore } from "../../store/chatting";
import { useAppStore } from "../../store/app";
import TraceViewer from "./trace_viewer.vue";

const MessageView = defineAsyncComponent(() => import("./message.vue"));
const SendIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SendOutline"),
);
const BugIcon = defineAsyncComponent(() => import("@vicons/ionicons5/BugOutline"));
const CloseIcon = defineAsyncComponent(() => import("@vicons/ionicons5/Close"));
const SparklesIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SparklesOutline"),
);

// Same purple the trace viewer uses for the AI's own internals.
const COPILOT_COLOR = "#7c5cff";

const props = defineProps<{ chat: Chat }>();

const store = useCcStore();
const appStore = useAppStore();
const notification = useNotification();

// Its own draft: sharing store.current_message with the main composer would
// make text typed in one appear in the other.
const draft = ref("");
const sending = ref(false);
const isTracesOpen = ref(false);
const scrollbar = ref();

// The operator's lane: notes, plus the bot's superseded drafts. The main chat
// renders the exact complement - see splitChatMessages.
const messages = computed(
  () => splitChatMessages(store.chat_messages(props.chat)).copilot,
);

// The note currently being edited, or null when composing a new one. Local to
// this pane for the same reason `draft` is: the main composer must not react to
// an edit started here.
const editing = ref<Message | null>(null);

function startEdit(message: Message) {
  editing.value = message;
  draft.value = message.content;
}

function cancelEdit() {
  editing.value = null;
  draft.value = "";
}

async function send() {
  const content = draft.value.trim();
  if (!content || sending.value) return;
  sending.value = true;
  try {
    if (editing.value) {
      const updated = Message.fromJson(editing.value.toJson() as any);
      updated.content = content;
      await store.update_message(updated);
      editing.value = null;
    } else {
      await store.send_message(
        new Message({ content, kind: Kind.COPILOT, chat: props.chat.uuid }),
      );
    }
    draft.value = "";
  } catch (e) {
    notification.error({
      title: editing.value ? "Couldn't save the note" : "Couldn't send the note",
      content: (e as Error).message,
      duration: 5000,
    });
  } finally {
    sending.value = false;
  }
}

async function scrollToBottom() {
  await nextTick();
  setTimeout(() => {
    if (!scrollbar.value) return;
    const top = scrollbar.value.$el.nextSibling.firstChild.scrollHeight;
    scrollbar.value.scrollTo({ top, behavior: "smooth" });
  }, 300);
}
watch(messages, scrollToBottom, { deep: true });
</script>

<style scoped>
/* Mirrors Chat.vue's .chat so both panes lay out identically. */
.chat {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  padding-left: 16px;
  border-left: 1px solid var(--n-border-color);
}
.chat__scroll {
  flex: 1 1 auto;
  min-height: 0;
}

/* The pane is dragged down to 20% of the window; below that the header has to
   give up the subtitle rather than push the buttons out of view. */
@container (max-width: 320px) {
  .head__sub {
    display: none;
  }
}

.head {
  display: flex;
  align-items: center;
  gap: 10px;
  container-type: inline-size;
}
.head__titles {
  flex: 1;
  min-width: 0;
}
.head__title {
  font-weight: 600;
  line-height: 1.25;
}
.head__sub {
  font-size: 11px;
  opacity: 0.5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.editing {
  grid-column: 1 / 3;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  opacity: 0.6;
}

/* Same grid as chat_footer's .footer */
.footer {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  width: calc(100% - 20px);
  gap: 5px 15px;
}
.footer .textarea {
  min-width: 100px;
  width: 100%;
}

kbd {
  background-color: #eee;
  border-radius: 3px;
  border: 1px solid #b4b4b4;
  color: #333;
  display: inline-block;
  font-size: 0.85em;
  font-weight: 700;
  line-height: 1;
  padding: 2px 4px;
  white-space: nowrap;
}
</style>
