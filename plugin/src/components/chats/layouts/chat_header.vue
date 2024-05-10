<template>
  <n-text v-if="!chat">Loading...</n-text>
  <div class="grid" v-else>
    <n-button
      ghost
      v-if="appStore.displayMode === 'none'"
      @click="appStore.displayMode = 'full'"
    >
      <n-icon> <open-icon /> </n-icon>
    </n-button>

    <n-space class="main__info" justify="start" align="center">
      <n-tooltip>
        <template #trigger>
          <n-tag round @click="addToClipboard(chat.uuid, notification)">
            <code style="text-decoration: underline; cursor: pointer">
              {{ chat.uuid.slice(0, 8).toUpperCase() }}
            </code>
          </n-tag>
        </template>
        {{ chat.uuid }}
      </n-tooltip>

      <user-avatar
        round
        :avatar="members.map((m) => m?.title ?? '').join(' ')"
      />
      <n-text class="chat__topic">{{ chat.topic ?? members }}</n-text>
      <n-button text @click="startEditChat">
        <n-icon size="20">
          <edit-icon />
        </n-icon>
      </n-button>

      <template v-if="!appStore.isMobile">
        <n-divider vertical />
        <n-tooltip>
          <template #trigger>
            <n-select
              filterable
              ref="responsibleSelect"
              label-field="title"
              value-field="uuid"
              placeholder="Responsible"
              style="min-width: 200px; width: 100%"
              :value="chat.responsible"
              :options="adminsItems"
              @update:value="changeResponsible"
            />
          </template>
          Responsible
        </n-tooltip>
      </template>

      <n-collapse v-else>
        <n-collapse-item title="Responsible">
          <n-select
            filterable
            ref="responsibleSelect"
            label-field="title"
            value-field="uuid"
            placeholder="Responsible"
            style="min-width: 200px; width: 100%"
            :value="chat.responsible"
            :options="adminsItems"
            @update:value="changeResponsible"
          />
        </n-collapse-item>
      </n-collapse>

      <template v-if="!appStore.isMobile">
        <n-divider vertical />
        <n-tooltip>
          <template #trigger>
            <n-select
              filterable
              label-field="title"
              value-field="key"
              placeholder="Department"
              style="min-width: 200px; width: 100%"
              :value="chat.department"
              :options="departments"
              @update:value="changeDepartment"
            />
          </template>
          Department
        </n-tooltip>
      </template>
    </n-space>

    <template v-if="!appStore.isMobile">
      <n-space
        align="center"
        :wrap-item="false"
        v-if="chat.gateways.length > 0"
      >
        <n-divider vertical />
        <n-tooltip v-for="gateway of chat.gateways" placement="bottom">
          <template #trigger>
            <img height="24" :src="getImageUrl(gateway)" :alt="gateway" />
          </template>
          {{ gateway }}
        </n-tooltip>
      </n-space>

      <n-space align="center" justify="end" :wrap-item="false">
        <n-divider vertical />
        <div>
          <chat-status :chat="chat" />
          <members-dropdown
            :visible="isVisible"
            :admins="chat.admins"
            :members="members"
            @add="startAddMembers"
            @delete="deleteMember"
            @change="openResponsible"
          />
        </div>

        <n-divider vertical />
        <chat-actions :chat="chat" />
      </n-space>
    </template>

    <n-space v-if="!appStore.isMobile" :wrap-item="false">
      <n-divider vertical style="height: auto; margin: 0 8px" />
      <chat-dates :chat="chat" />
    </n-space>
    <n-collapse v-else>
      <n-collapse-item title="Dates">
        <n-divider vertical style="height: auto; margin: 0 8px" />
        <chat-dates :chat="chat" />
      </n-collapse-item>
    </n-collapse>

    <n-space
      v-if="!appStore.isMobile"
      style="justify-self: start"
      :wrap-item="false"
      :style="{ gridColumn: appStore.displayMode !== 'half' ? 2 : 1 }"
    >
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
    </n-space>
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
        <member-select
          v-model:value="chatWithNewMembers!.users"
          :options="availableMembersOptions"
        />
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

// TODO: // - [ ] Wrap Topic Around Avatar // - [ ] Make menu draggable
(increase width)

<script setup lang="ts">
import { computed, defineAsyncComponent, nextTick, ref, toRefs } from "vue";
import {
  NButton,
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
  SelectOption,
  useNotification,
  NCollapse,
  NCollapseItem,
} from "naive-ui";
import { ConnectError } from "@connectrpc/connect";
import { Chat, User } from "../../../connect/cc/cc_pb";
import { useCcStore } from "../../../store/chatting.ts";
import { useAppStore } from "../../../store/app";
import ChatOptions from "../chat_options.vue";
import UserAvatar from "../../ui/user_avatar.vue";
import MembersDropdown from "../../users/members_dropdown.vue";
import useDefaults from "../../../hooks/useDefaults.ts";
import MemberSelect from "../../users/member_select.vue";
import { addToClipboard, getImageUrl } from "../../../functions.ts";
import ChatStatus from "../chat_status.vue";
import ChatActions from "../chat_actions.vue";
import ChatDates from "../chat_dates.vue";

const EditIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/PencilSharp")
);
const OpenIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ArrowBack")
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
const {
  fetch_defaults,
  isDefaultLoading,
  users,
  admins,
  metrics,
  departments,
} = useDefaults();

const isEdit = ref<boolean>(false);
const isAddDialog = ref<boolean>(false);
const chatWithNewMembers = ref<Chat>();
const availableMembersOptions = ref<SelectOption[]>([]);
const isAddSaveLoading = ref<boolean>(false);

const members = computed(() => {
  const uuids = new Set([
    ...chat!.value.users,
    ...chat.value.admins,
    chat.value.responsible,
  ]);
  const result: User[] = [];

  uuids.forEach((uuid) => {
    if (!uuid) return;
    result.push(store.users.get(uuid) as User);
  });

  return result;
});
const me = computed(() => store.me);

fetch_defaults();

const metricsOptions = computed(() => {
  const metricsEntries = Object.entries(chat.value.meta?.data ?? {});
  const result: Metric[] = [];

  metricsEntries.forEach(([keyMetric, { kind }]) => {
    const { title, options } =
      metrics.value.find((metric) => metric.key === keyMetric) ?? {};

    const { key, value } = options?.find(
      (option) => option.value === kind.value
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
      })
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
      })
    );

    notification.success({ title: "Done", duration: 3000 });
  } catch (error) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
    });
  }
};

const adminsItems = computed(() =>
  admins.value.map(
    (admin) => users.value.find(({ uuid }) => uuid === admin) ?? { uuid: admin }
  )
);
// const responsible = computed(() =>
//   users.value.find(({ uuid }) => uuid === chat.value.responsible) as User
// )

const isVisible = ref<boolean>();
const responsibleSelect = ref<any>();

const openResponsible = async () => {
  await nextTick();
  await new Promise((resolve) => setTimeout(resolve, 100));
  responsibleSelect.value?.handleTriggerClick();
};

const deleteMember = (uuid: string) => {
  const users = chat.value.users.filter((userId) => userId !== uuid);

  store.update_chat({ ...chat.value, users } as Chat);
};

const fetchAvailableUsers = () => {
  availableMembersOptions.value = users.value.map((user) => {
    return {
      label: user.title,
      value: user.uuid,
      disabled: me.value.uuid !== user.uuid && admins.value.includes(user.uuid),
    };
  });
};

const startAddMembers = () => {
  fetchAvailableUsers();

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

const gridColumns = computed(
  () => `repeat(${chat.value?.gateways.length > 0 ? 3 : 2}, auto) 1fr auto`
);
</script>

<style scoped lang="scss">
.grid {
  display: grid;
  grid-template-columns: v-bind("gridColumns");
  align-items: center;
  gap: 10px;

  @media only screen and (max-width: 1024px) {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0px;
  }
}

.main__info {
  @media only screen and (max-width: 1024px) {
    display: grid !important;
    grid-template-columns: 80px 30px auto 25px;
  }
}

.chat__topic {
  display: inline-block;
  max-width: 200px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
