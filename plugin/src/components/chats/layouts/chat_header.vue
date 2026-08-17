<template>
  <!-- Sized by its OWN width, not the window's: the copilot pane takes up to
       half the screen, and a viewport-based header keeps laying out for a width
       it no longer has (that is what wrapped the action icons onto two rows).

       Nothing is dropped when there is no room - the secondary controls move
       into the "more" popover. An operator who shrank the pane still has to be
       able to change the department. -->
  <div class="hdr" ref="root">
    <n-text v-if="!chat">Loading...</n-text>

    <template v-else>
      <n-button
        ghost
        size="small"
        v-if="appStore.displayMode === 'none'"
        @click="appStore.displayMode = 'full'"
      >
        <n-icon> <open-icon /> </n-icon>
      </n-button>

      <div class="hdr__id">
        <n-tooltip>
          <template #trigger>
            <n-tag round size="small" @click="addToClipboard(chat.uuid, notification)">
              <code style="text-decoration: underline; cursor: pointer">
                {{ chat.uuid.slice(0, 8).toUpperCase() }}
              </code>
            </n-tag>
          </template>
          {{ chat.uuid }}
        </n-tooltip>

        <user-avatar
          round
          style="flex: none"
          :avatar="members.map((m) => m?.title ?? '').join(' ')"
        />

        <n-text class="chat__topic">{{ chat.topic ?? members }}</n-text>
        <n-button text style="flex: none" @click="startEditChat">
          <n-icon size="18">
            <edit-icon />
          </n-icon>
        </n-button>
      </div>

      <!-- Inline while they fit; inside the popover otherwise. -->
      <template v-if="!compact">
        <n-tooltip>
          <template #trigger>
            <n-select
              filterable
              ref="responsibleSelect"
              placeholder="Responsible"
              size="small"
              class="hdr__select"
              :value="chat.responsible"
              :options="adminsItems"
              :loading="isUsersLoading"
              @update:value="changeResponsible"
            />
          </template>
          Responsible
        </n-tooltip>

        <n-tooltip>
          <template #trigger>
            <n-select
              filterable
              label-field="title"
              value-field="key"
              placeholder="Department"
              size="small"
              class="hdr__select"
              :value="chat.department"
              :options="departments"
              @update:value="changeDepartment"
            />
          </template>
          Department
        </n-tooltip>
      </template>

      <!-- Takes the slack so the free space lands between the two groups
           instead of inside the left one, which read as a hole in the header. -->
      <div class="hdr__spacer"></div>

      <div class="hdr__side">
        <template v-if="!compact && chat.gateways.length > 0">
          <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
            <template #trigger>
              <img height="22" :src="getImageUrl(gateway)" :alt="gateway" />
            </template>
            {{ gateway }}
          </n-tooltip>
          <n-divider vertical />
        </template>

        <chat-status :chat="chat" />
        <members-dropdown
          :visible="isVisible"
          :admins="chat.admins"
          :members="members"
          @add="startAddMembers"
          @delete="deleteMember"
          @change="openResponsible"
        />

        <n-divider vertical />

        <!-- nowrap: eight circle buttons breaking onto a second row is what
             made the header twice as tall and pushed the chat down. -->
        <div class="hdr__actions">
          <chat-actions :chat="chat" />
        </div>

        <template v-if="!compact">
          <n-divider vertical />
          <chat-dates :chat="chat" />
        </template>

        <n-popover v-else trigger="click" placement="bottom-end">
          <template #trigger>
            <n-button size="small" ghost circle>
              <template #icon>
                <n-icon><more-icon /></n-icon>
              </template>
            </n-button>
          </template>

          <div class="more">
            <div class="more__row">
              <span class="more__label">Responsible</span>
              <n-select
                filterable
                placeholder="Responsible"
                size="small"
                :value="chat.responsible"
                :options="adminsItems"
                :loading="isUsersLoading"
                @update:value="changeResponsible"
              />
            </div>

            <div class="more__row">
              <span class="more__label">Department</span>
              <n-select
                filterable
                label-field="title"
                value-field="key"
                placeholder="Department"
                size="small"
                :value="chat.department"
                :options="departments"
                @update:value="changeDepartment"
              />
            </div>

            <div class="more__row" v-if="chat.gateways.length > 0">
              <span class="more__label">Channels</span>
              <div class="more__gateways">
                <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
                  <template #trigger>
                    <img height="22" :src="getImageUrl(gateway)" :alt="gateway" />
                  </template>
                  {{ gateway }}
                </n-tooltip>
              </div>
            </div>

            <div class="more__row" v-if="metricsOptions.length > 0">
              <span class="more__label">Metrics</span>
              <div class="more__metrics">
                <n-text v-for="metric in metricsOptions">
                  {{ metric.title }}:
                  <n-tag
                    round
                    size="small"
                    type="error"
                    :style="`filter: ${getTagColor(metric)}`"
                  >
                    {{ metric.key }}
                  </n-tag>
                </n-text>
              </div>
            </div>

            <div class="more__row">
              <span class="more__label">Dates</span>
              <chat-dates :chat="chat" />
            </div>
          </div>
        </n-popover>
      </div>

      <div v-if="!compact && metricsOptions.length > 0" class="hdr__metrics">
        <n-text v-for="metric in metricsOptions">
          {{ metric.title }}:
          <n-tag
            round
            size="small"
            type="error"
            :style="`filter: ${getTagColor(metric)}`"
          >
            {{ metric.key }}
          </n-tag>
        </n-text>
      </div>
    </template>
  </div>

  <n-modal v-model:show="isEdit">
    <n-card
      title="Edit chat options"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
      style="width: 500px; min-height: 500px"
    >
      <chat-options
        @close="isEdit = false"
        is-edit
        :chat="chat"
        style="width: 100%"
      />
    </n-card>
  </n-modal>

  <n-modal v-model:show="isAddDialog">
    <n-card
      title="Add members"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
      style="width: 400px"
    >
      <template v-if="!isDefaultLoading">
        <member-select-pagination v-model:value="chatWithNewMembers!.users" />

        <n-space style="margin-top: 10px" vertical align="end" justify="end">
          <n-button :loading="isAddSaveLoading" @click="saveMembers"
            >Save</n-button
          >
        </n-space>
      </template>
      <n-spin
        style="width: 100%; height: 100%; margin: auto"
        size="large"
        v-else
      ></n-spin>
    </n-card>
  </n-modal>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  toRefs,
} from "vue";
import {
  NButton,
  NPopover,
  NCard,
  NDivider,
  NIcon,
  NModal,
  NSpace,
  NSpin,
  NTag,
  NText,
  NTooltip,
  NSelect,
  useNotification,
} from "naive-ui";
import { ConnectError } from "@connectrpc/connect";
import { Chat, User } from "../../../connect/cc/cc_pb";
import { useCcStore } from "../../../store/chatting.ts";
import { useAppStore } from "../../../store/app";
import ChatOptions from "../chat_options.vue";
import UserAvatar from "../../ui/user_avatar.vue";
import MembersDropdown from "../../users/members_dropdown.vue";
import { addToClipboard, getImageUrl } from "../../../functions.ts";
import ChatStatus from "../chat_status.vue";
import ChatActions from "../chat_actions.vue";
import ChatDates from "../chat_dates.vue";
import { useDefaultsStore } from "../../../store/defaults.ts";
import { storeToRefs } from "pinia";
import { useUsersStore } from "../../../store/users.ts";
import MemberSelectPagination from "../../users/member_select_pagination.vue";

const EditIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/PencilSharp"),
);
const OpenIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ArrowBack"),
);
const MoreIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/EllipsisHorizontal"),
);

interface ChatHeaderProps {
  chat: Chat;
}

interface Metric {
  title: string;
  value: number;
  key: string;
  max: number;
  min: number;
}

const props = defineProps<ChatHeaderProps>();
const { chat } = toRefs(props);

const appStore = useAppStore();
const store = useCcStore();
const notification = useNotification();
const defaultsStore = useDefaultsStore();
const { isDefaultLoading, admins, metrics, departments } =
  storeToRefs(defaultsStore);

const usersStore = useUsersStore();
const { users, isUsersLoading } = storeToRefs(usersStore);

const isEdit = ref<boolean>(false);
const isAddDialog = ref<boolean>(false);
const chatWithNewMembers = ref<Chat>();
const isAddSaveLoading = ref<boolean>(false);
// How much room the header itself has. 0 until the first observation, and both
// flags read false then, so the first paint is the full layout rather than a
// mobile one that snaps wide a frame later.
const root = ref<HTMLElement>();
const width = ref(0);
// Everything laid out inline needs roughly 1240px: two selects (380), the
// action row (270), status + members (150), dates (200) and the topic. Below
// that the secondary half moves into the popover instead of being cut off.
const compact = computed(() => width.value > 0 && width.value < 1240);

let observer: ResizeObserver | undefined;
onMounted(() => {
  if (!root.value) return;
  observer = new ResizeObserver(([entry]) => {
    width.value = entry.contentRect.width;
  });
  observer.observe(root.value);
});
onBeforeUnmount(() => observer?.disconnect());

const members = computed(() => {
  const uuids = new Set([
    ...chat!.value.users,
    ...chat.value.admins,
    chat.value.responsible,
  ]);
  const result: User[] = [];

  uuids.forEach((uuid) => {
    if (!uuid) return;
    result.push(users.value.get(uuid) as User);
  });

  return result;
});

const metricsOptions = computed(() => {
  const metricsEntries = Object.entries(chat.value.meta?.data ?? {});
  const result: Metric[] = [];

  metricsEntries.forEach(([keyMetric, { kind }]) => {
    const { title, options } =
      metrics.value.find((metric) => metric.key === keyMetric) ?? {};

    const { key, value } = options?.find(
      (option) => option.value === kind.value,
    ) ?? { key: "", value: 0 };

    const optionsValues = options?.map(({ value }) => value) ?? [];
    const min = Math.min(...optionsValues);
    const max = Math.max(...optionsValues);

    if (title) result.push({ title, value, key, min, max });
  });

  return result;
});

const getTagColor = (metric: Metric) =>
  `hue-rotate(${
    220 - (220 * (metric.value - metric.min)) / (metric.max - metric.min)
  }deg)`;

const changeResponsible = async (uuid: string) => {
  try {
    await store.update_chat(
      new Chat({
        ...chat.value,
        responsible: uuid,
      }),
    );

    notification.success({ title: "Done", duration: 3000 });
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
};

const changeDepartment = async (key: string) => {
  try {
    await store.change_department(
      new Chat({
        ...chat.value,
        department: key,
      }),
    );

    notification.success({ title: "Done", duration: 3000 });
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
};

const adminsItems = computed(() =>
  admins.value
    .map((admin) => users.value.get(admin))
    .filter((u) => !!u)
    .map((user) => ({ label: user!.title, value: user!.uuid })),
);

const isVisible = ref<boolean>();
const responsibleSelect = ref<any>();

const openResponsible = async () => {
  await nextTick();
  await new Promise((resolve) => setTimeout(resolve, 100));
  responsibleSelect.value?.handleTriggerClick();
};

const deleteMember = (uuid: string) => {
  const users = chat.value.users.filter((userId) => userId !== uuid);
  const admins = chat.value.admins.filter((adminId) => adminId !== uuid);

  store.update_chat({ ...chat.value, users, admins } as Chat);
};

const startAddMembers = () => {
  chatWithNewMembers.value = { ...chat.value } as Chat;
  isAddDialog.value = true;
};

const startEditChat = () => {
  isEdit.value = true;
};

const saveMembers = async () => {
  try {
    isAddSaveLoading.value = true;
    await store.update_chat(chatWithNewMembers.value as Chat);
    isAddDialog.value = false;
  } finally {
    isAddSaveLoading.value = false;
  }
};


</script>

<style scoped lang="scss">
.hdr {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
  /* The pane only pads its left edge; without this the last control sits flush
     against the copilot divider and looks clipped. */
  padding-right: 12px;
}

.hdr__id {
  display: flex;
  align-items: center;
  gap: 6px;
  /* Shrinks but never grows: growing put the slack between the topic and the
     selects, which is the gap in the middle of the header. */
  flex: 0 1 auto;
  min-width: 0;
}

.hdr__spacer {
  flex: 1 1 0;
  min-width: 0;
}

.hdr__select {
  flex: 0 1 190px;
  min-width: 120px;
}

.hdr__side {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: none;
}

.hdr__actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: nowrap;
}

.hdr__metrics {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  flex-basis: 100%;
}

.chat__topic {
  flex: 0 1 auto;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.more {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 260px;
}
.more__row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.more__label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  opacity: 0.55;
}
.more__gateways {
  display: flex;
  align-items: center;
  gap: 8px;
}
.more__metrics {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
