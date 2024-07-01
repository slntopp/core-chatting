<template>
  <div
    v-if="appStore.displayMode !== 'none'"
    :class="{ chats__panel: true, closed: !isChatPanelOpen }"
    :style="
      appStore.displayMode === 'full'
        ? 'min-width: calc(100% - 8px)'
        : undefined
    "
  >
    <n-layout-sider
      style="margin-bottom: 25px"
      collapse-mode="transform"
      :collapsed="!isChatPanelOpen"
    >
      <n-space
        align="center"
        :wrap-item="false"
        :justify="
          isChatPanelOpen || appStore.isMobile ? 'space-between' : 'center'
        "
        :class="{ chat__actions: true, hide: !isChatPanelOpen }"
      >
        <template v-if="!isSelectedChats">
          <n-button v-if="!appStore.isMobile" ghost @click="changeMode(null)">
            <n-icon> <switch-icon /> </n-icon>
          </n-button>
          <n-button ghost type="success" @click="startChat">
            <n-icon :component="ChatbubbleEllipsesOutline" />
            <span v-if="isChatPanelOpen" style="margin-left: 5px">
              Start Chat
            </span>
          </n-button>

          <n-space
            v-if="
              isChatPanelOpen &&
              appStore.displayMode === 'full' &&
              !appStore.isMobile
            "
            :wrap-item="false"
            :style="isChatPanelOpen ? 'margin-right: auto' : null"
          >
            <span v-for="{ count, status } in chatsCountByStatus">
              <n-text
                :style="{ color: getStatusColor(+status), cursor: 'pointer' }"
                @click="selectStatus(+status)"
              >
                {{ getStatus(+status) }}:
              </n-text>
              <n-tag
                :type="selectedStatus === +status ? 'success' : undefined"
                round
                size="small"
                >{{ count }}</n-tag
              >
            </span>
          </n-space>

          <n-select
            v-else
            :value="selectedStatus"
            @change="selectStatus"
            clearable
            placeholder="Status"
            :options="
              Object.values(chatsCountByStatus).map(({ count, status }) => ({
                value: +status,
                label: `${getStatus(status)} (${count})`,
              }))
            "
          />
        </template>

        <n-space v-else>
          <n-button ghost @click="resetSelectedChats"> Cancel </n-button>

          <n-button
            ghost
            type="warning"
            v-if="isMergeVisible"
            :loading="mergeLoading"
            @click="mergeChats"
          >
            <n-icon :component="mergeIcon" />
            <span v-if="isChatPanelOpen" style="margin-left: 5px">
              Merge Chats
            </span>
          </n-button>

          <n-button
            ghost
            type="error"
            v-if="isSelectedChats"
            :loading="deleteLoading"
            @click="deleteChats"
          >
            <n-icon :component="deleteIcon" />
            <span v-if="isChatPanelOpen" style="margin-left: 5px">
              Delete Chats
            </span>
          </n-button>
          <n-space>
            <n-select
              style="width: 200px"
              v-model:value="newStatus"
              clearable
              placeholder="Status"
              :options="getStatusItems()"
            />
            <n-button
              :disabled="isNaN(newStatus) || newStatus === null"
              @click="changeChatsStatus"
              type="warning"
              ghost
              :loading="isChangeStatusLoading"
            >
              Change
            </n-button>
          </n-space>

          <n-space>
            <n-select
              style="width: 200px"
              v-model:value="newDepartment"
              clearable
              placeholder="Departament"
              :options="departments"
            />
            <n-button
              :disabled="!newDepartment"
              @click="changeDepartment"
              type="warning"
              ghost
              :loading="isChangeDepartmentLoading"
            >
              Change
            </n-button>
          </n-space>

          <n-space>
            <n-select
              style="width: 200px"
              v-model:value="newResponsible"
              clearable
              placeholder="Responsible"
              :options="responsibles"
              label-field="title"
              value-field="uuid"
            />
            <n-button
              :disabled="!newResponsible"
              @click="changeResponsible"
              type="warning"
              ghost
              :loading="isChangeResponsibleLoading"
            >
              Change
            </n-button>
          </n-space>
        </n-space>

        <n-button v-if="appStore.isPC" ghost @click="changePanelOpen">
          <n-icon>
            <close-icon v-if="isChatPanelOpen" />
            <open-icon v-else />
          </n-icon>
        </n-button>
      </n-space>

      <n-space
        class="search"
        align="center"
        justify="space-between"
        v-if="isChatPanelOpen"
        :wrap="false"
      >
        <n-checkbox
          style="padding-left: 20px"
          v-model:checked="isAllChatsSelected"
        />

        <n-input
          v-model:value="searchParam"
          type="text"
          placeholder="Search..."
        />

        <n-popover
          scrollable
          trigger="click"
          placement="bottom"
          class="chats__filters"
        >
          <template #trigger>
            <n-icon size="24" style="vertical-align: middle; cursor: pointer">
              <sort-icon />
            </n-icon>
          </template>

          <div>
            <n-space align="center" style="gap: 5px" :wrap-item="false">
              <n-text>Sort By:</n-text>
              <n-icon
                size="18"
                style="cursor: pointer"
                :style="{
                  transform: sortType === 'desc' ? 'rotate(180deg)' : null,
                }"
                @click="
                  sortType === 'desc' ? (sortType = 'asc') : (sortType = 'desc')
                "
              >
                <up-icon />
              </n-icon>
            </n-space>

            <n-divider style="margin: 5px 0" />
            <n-radio-group v-model:value="sortBy">
              <n-space :wrap-item="false">
                <n-radio value="none" label="None" />
                <n-radio value="created" label="Created" />
                <n-radio value="sent" label="Sent" />
                <n-radio value="status" label="Status" />
              </n-space>
            </n-radio-group>
          </div>

          <chats-filters
            v-model:checkedDepartments="checkedDepartments"
            v-model:checkedStatuses="checkedStatuses"
            v-model:checkedAdmins="checkedAdmins"
            v-model:checkedResponsibles="checkedResponsibles"
            v-model:update-date-range="updateDateRange"
            v-model:create-date-range="createDateRange"
            :metricsOptions="metricsOptions"
            @update:checked-metrics="
              (key, value) => (metricsOptions[key] = value)
            "
          />

          <n-flex style="margin: 15px 10px 0px 0px;" justify="end">
            <n-button
              ghost
              type="primary"
              @click="downloadReport"
              :loading="isReportLoading"
            >
              Report
            </n-button>

            <n-button ghost type="primary" @click="resetFilters">
              Reset
            </n-button>
          </n-flex>
        </n-popover>
      </n-space>

      <n-scrollbar
        ref="scrollbar"
        style="height: calc(100vh - 100px); min-width: 150px"
      >
        <n-popover trigger="manual" :show="isFirstMessageVisible" :x="x" :y="y">
          {{ firstMessage }}
        </n-popover>

        <n-list
          v-if="!(isLoading || isDefaultLoading)"
          hoverable
          clickable
          style="margin-bottom: 25px"
        >
          <chat-item
            v-for="chat in viewedChats"
            :selected="selectedChats"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat"
            :chats="viewedChats"
            :class="{
              active: chat.uuid === router.currentRoute.value.params.uuid,
            }"
            @click="changeMode('none')"
            @hover="onMouseMove"
            @hoverEnd="isFirstMessageVisible = false"
            @select="selectChat"
          />
        </n-list>
        <n-spin
          v-if="isLoading || viewedChats.length < chats.length"
          ref="loading"
          style="display: flex; margin: 0 auto 20px"
        />
      </n-scrollbar>
    </n-layout-sider>
  </div>

  <div id="separator" v-show="isChatPanelOpen"></div>
  <div class="chat__item" v-if="appStore.displayMode !== 'full'">
    <n-layout-content>
      <router-view />
    </n-layout-content>
  </div>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onMounted,
  ref,
  watch,
} from "vue";
import {
  NButton,
  NIcon,
  NInput,
  NLayoutContent,
  NLayoutSider,
  NList,
  NFlex,
  NScrollbar,
  NSpace,
  NPopover,
  NRadioGroup,
  NRadio,
  NDivider,
  NTag,
  NText,
  NSpin,
  NSelect,
  NCheckbox,
  useNotification,
} from "naive-ui";

import { useAppStore } from "../../store/app.ts";
import { useCcStore } from "../../store/chatting.ts";

import { useRoute, useRouter } from "vue-router";
import { Chat, Status, User } from "../../connect/cc/cc_pb";
import ChatItem from "../../components/chats/chat_item.vue";
import ChatsFilters from "../../components/chats/chats_filters.vue";
import useDraggable from "../../hooks/useDraggable.ts";
import useDefaults from "../../hooks/useDefaults.ts";
import { getStatusColor, getStatusItems } from "../../functions.ts";
import { ConnectError } from "@connectrpc/connect";
import ChatReportService from "../../services/ChatReportService.ts";

defineEmits(["hover", "hoverEnd"]);

const ChatbubbleEllipsesOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChatbubbleEllipsesOutline")
);
const OpenIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ArrowForward")
);
const CloseIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ArrowBack")
);
const SortIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/EllipsisVertical")
);
const SwitchIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SwapHorizontal")
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseCircle")
);
const mergeIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/GitMerge")
);
const upIcon = defineAsyncComponent(() => import("@vicons/ionicons5/ArrowUp"));

const appStore = useAppStore();
const store = useCcStore();
const router = useRouter();
const route = useRoute();
const { makeDraggable } = useDraggable();
const { metrics, isDefaultLoading, fetch_defaults, admins, users } =
  useDefaults();
const notification = useNotification();

const selectedChats = ref<string[]>([]);
const deleteLoading = ref(false);
const mergeLoading = ref(false);
const searchParam = ref("");
const scrollbar = ref<InstanceType<typeof NScrollbar>>();
const loading = ref<InstanceType<typeof NSpin>>();
const page = ref(1);
const isLoading = ref(false);
const isReportLoading = ref(false);
const newStatus = ref();
const isChangeStatusLoading = ref(false);
const newDepartment = ref();
const isChangeDepartmentLoading = ref(false);
const newResponsible = ref();
const isChangeResponsibleLoading = ref(false);
const isAllChatsSelected = ref(false);

const isMergeVisible = computed(() => {
  const list = selectedChats.value.map(
    (uuid) => new Chat(store.chats.get(uuid))
  );

  const dep = list.at(0)?.department;
  const isDepsEqual = list.every(({ department }) => department === dep);

  return isSelectedChats.value && isDepsEqual;
});

const isSelectedChats = computed(() => selectedChats.value.length > 0);

const departments = computed(() =>
  store.departments.map(({ key, title }) => ({ label: title, value: key }))
);

const responsibles = computed(() =>
  admins.value.map(
    (admin) => users.value.find(({ uuid }) => uuid === admin) ?? { uuid: admin }
  )
);

async function sync() {
  try {
    isLoading.value = true;
    await store.list_chats();
    const { users } = await store.get_members();

    users.forEach((user) => {
      store.users.set(user.uuid, user);
    });
  } catch (error) {
    console.log(error);
  } finally {
    isLoading.value = false;
  }
}

function startChat() {
  router.push({ name: "Start Chat" });
  appStore.displayMode = "none";
}

async function deleteChats() {
  deleteLoading.value = true;
  try {
    const current = route.params.uuid as string;
    if (selectedChats.value.includes(current)) {
      await router.push({ name: "Dashboard" });
      changePanelOpen();
    }

    const promises = selectedChats.value.map((uuid) =>
      store.delete_chat(new Chat(store.chats.get(uuid)))
    );

    await Promise.all(promises);
    notification.success({
      title: "Chats successfully deleted",
      duration: 5000,
    });

    resetSelectedChats();
  } catch (error: any) {
    notification.error({
      title: error.response?.data.message ?? error.message ?? error,
      duration: 5000,
    });
    console.error(error);
  } finally {
    deleteLoading.value = false;
  }
}

async function mergeChats() {
  mergeLoading.value = true;
  try {
    const current = route.params.uuid as string;
    if (selectedChats.value.includes(current)) {
      await router.push({ name: "Dashboard" });
      changePanelOpen();
    }

    await store.merge_chats(selectedChats.value);
    notification.success({
      title: "Chats successfully merged",
      duration: 5000,
    });

    resetSelectedChats();
  } catch (error: any) {
    notification.error({
      title: error.response?.data.message ?? error.message ?? error,
      duration: 5000,
    });
    console.error(error);
  } finally {
    mergeLoading.value = false;
  }
}

onMounted(() => {
  makeDraggable({
    resizer: document.getElementById("separator")!,
    first: document.querySelector(".chats__panel")!,
    second: document.querySelector(".chat__item")!,
  });

  const sorting = JSON.parse(localStorage.getItem("sorting") ?? "null");
  const filters = JSON.parse(localStorage.getItem("filters") ?? "null");

  if (sorting) {
    sortBy.value = sorting.sortBy;
    sortType.value = sorting.sortType;
  }
  if (filters) {
    checkedDepartments.value = filters.departments;
    checkedStatuses.value = filters.statuses;
    checkedAdmins.value = filters.admins;
    checkedResponsibles.value = filters.responsibles;
    updateDateRange.value = filters.updated;
    createDateRange.value = filters.created;
  }

  const options = {
    root: scrollbar.value?.$el.nextElementSibling,
    threshold: 1,
  };

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(async ({ isIntersecting }) => {
      await new Promise((resolve) => setTimeout(resolve, 200));
      if (isIntersecting) page.value += 1;
    });
  }, options);

  if (!loading.value?.$el) return;
  observer.observe(loading.value?.$el);
});

sync();

function selectChat(uuid: string) {
  if (selectedChats.value.includes(uuid)) {
    const i = selectedChats.value.indexOf(uuid);

    selectedChats.value.splice(i, 1);
  } else {
    selectedChats.value.push(uuid);
  }
}

function filterChat(chat: Chat, val: string): boolean {
  if (!val) {
    return true;
  }
  val = val.toLowerCase();

  const startsWithKeys = ["topic", "uuid"];

  for (const key of startsWithKeys) {
    if ((chat as any)[key]?.toLowerCase().startsWith(val)) {
      return true;
    }
  }

  return (
    !!chat.users.find((u) => {
      const user = store.users.get(u)?.toJson() as any;
      return (
        u.startsWith(val) ||
        user?.title.toLowerCase().startsWith(val) ||
        user?.data?.email?.toLowerCase().startsWith(val)
      );
    }) ||
    !!chat.admins.find(
      (u) =>
        u.startsWith(val) ||
        store.users.get(u)?.title.toLocaleLowerCase().startsWith(val)
    ) ||
    !!chat.meta?.lastMessage?.content.toLowerCase().startsWith(val)
  );
}

const sortBy = ref<"none" | "sent" | "created" | "status">("sent");
const sortType = ref<"asc" | "desc">("desc");

interface optionsIncludedType {
  [index: string]: boolean;
}

const checkedStatuses = ref<Status[]>([]);
const selectedStatus = ref<Status>();
const statuses = computed(() =>
  Object.keys(Status)
    .filter((key) => isFinite(+key))
    .map((key) => ({
      label: getStatus(+key),
      value: +key,
    }))
    .reverse()
);

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace("_", " ");

  return `${status[0].toUpperCase()}${status.slice(1)}`;
}

function selectStatus(status: Status) {
  if (status === selectedStatus.value || status === null) {
    selectedStatus.value = undefined;
  } else {
    selectedStatus.value = status;
  }
}

const chats = computed(() => {
  const result = [...store.chats.values()].filter((chat) => {
    let isCreatedInDate = true;
    if (createDateRange.value) {
      const createdDate = Number(chat.created);
      isCreatedInDate =
        createDateRange.value[0] <= createdDate &&
        createDateRange.value[1] >= createdDate;
    }

    let isUpdatedInDate = true;
    if (updateDateRange.value) {
      const updatedDate = Number(
        chat.meta?.lastMessage?.edited || chat.meta?.lastMessage?.sent
      );
      isUpdatedInDate =
        updateDateRange.value[0] <= updatedDate &&
        updateDateRange.value[1] >= updatedDate;
    }

    let isDepartamentIncluded = true;
    if (checkedDepartments.value.length > 0) {
      isDepartamentIncluded = checkedDepartments.value.includes(
        chat.department
      );
    }

    const selectedStatuses =
      selectedStatus.value === undefined
        ? checkedStatuses.value
        : [selectedStatus.value];
    let isStatusIncluded = true;
    if (selectedStatuses.length > 0) {
      isStatusIncluded = selectedStatuses.includes(chat.status);
    }

    let isAdminsExist = true;
    if (checkedAdmins.value.length > 0) {
      isAdminsExist = !!checkedAdmins.value.find((uuid) =>
        chat.admins.includes(uuid)
      );
    }

    let isResponsibleExist = true;
    if (checkedResponsibles.value?.length > 0) {
      isAdminsExist = !!checkedResponsibles.value.includes(
        chat.responsible || ""
      );
    }

    let isOptionsIncluded: optionsIncludedType = {};
    Object.entries(metricsOptions.value).forEach(([key, value]) => {
      if (value.length < 1) {
        isOptionsIncluded[key] = true;
        return;
      }
      const metric = chat.meta?.data[key]?.kind.value as never;

      isOptionsIncluded[key] = value.includes(metric);
    });

    let isAccountOwner = true;
    if (appStore.conf?.params?.filterByAccount) {
      isAccountOwner = [...chat.users, ...chat.admins].includes(
        appStore.conf.params.filterByAccount
      );
    }

    return (
      filterChat(chat as Chat, searchParam.value) &&
      isDepartamentIncluded &&
      isStatusIncluded &&
      isAdminsExist &&
      isResponsibleExist &&
      isAccountOwner &&
      isCreatedInDate &&
      isUpdatedInDate &&
      Object.values(isOptionsIncluded).every((value) => value)
    );
  }) as Chat[];

  const sortable = (chat: Chat) => {
    if (sortBy.value === "status") {
      return statuses.value.map(({ value }) => value).indexOf(chat.status);
    }
    if (chat.meta?.lastMessage && sortBy.value === "sent") {
      return Number(chat.meta.lastMessage.sent);
    }
    if (sortBy.value === "none") {
      return (chat.meta?.unread ?? 0) > 0 ? Date.now() : Number(chat.created);
    }
    return Number(chat.created);
  };

  result.sort((a: Chat, b: Chat) =>
    sortType.value === "desc"
      ? sortable(b) - sortable(a)
      : sortable(a) - sortable(b)
  );

  return result;
});

const viewedChats = computed(
  () => chats.value.slice(0, 10 * page.value) as Chat[]
);

const filteredChatsByAccount = computed(() =>
  [...store.chats.values()].filter((chat) => {
    let isAccountOwner = true;

    if (appStore.conf?.params?.filterByAccount) {
      isAccountOwner = [...chat.users, ...chat.admins].includes(
        appStore.conf.params.filterByAccount
      );
    }

    return isAccountOwner;
  })
);

const chatsCountByStatus = computed(() => {
  const result: { [key: string]: { status: number; count: number } } = {};

  const allowedStatuses = getStatusItems();

  allowedStatuses.forEach(
    (_, index) =>
      (result[index] = { status: allowedStatuses[index].value, count: 0 })
  );

  filteredChatsByAccount.value.forEach((chat) => {
    if (!allowedStatuses.find((status) => status.value === +chat.status)) {
      return;
    }
    const index = allowedStatuses.findIndex(
      (status) => status.value === +chat.status
    );
    result[index].count++;
  });

  return result;
});

const changeChatsStatus = async () => {
  isChangeStatusLoading.value = true;

  try {
    const promises = selectedChats.value
      .map((uuid) => chats.value.find((chat) => chat.uuid === uuid))
      .map(async (chat) => {
        const data = { ...chat, status: newStatus.value } as Chat;
        await store.change_status(data);
        store.chats.set(chat!.uuid, data);
      });

    await Promise.all(promises);

    notification.success({
      title: "Chats statuses successfully changed",
      duration: 5000,
    });
    newStatus.value = undefined;
    resetSelectedChats();
  } catch (error: any) {
    notification.error({
      title: error.response?.data.message ?? error.message ?? error,
      duration: 5000,
    });
  } finally {
    isChangeStatusLoading.value = false;
  }
};

const changeDepartment = async () => {
  isChangeDepartmentLoading.value = true;

  try {
    const promises = selectedChats.value
      .map((uuid) => chats.value.find((chat) => chat.uuid === uuid))
      .map(async (chat) => {
        const data = { ...chat, department: newDepartment.value } as Chat;

        await store.change_department(data);
        store.chats.set(chat!.uuid, data);
      });

    await Promise.all(promises);

    notification.success({
      title: "Department successfully changed",
      duration: 5000,
    });

    newDepartment.value = undefined;
    resetSelectedChats();
  } catch (error: any) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
      duration: 5000,
    });
  } finally {
    isChangeDepartmentLoading.value = false;
  }
};

const changeResponsible = async () => {
  isChangeResponsibleLoading.value = true;

  try {
    const promises = selectedChats.value
      .map((uuid) => chats.value.find((chat) => chat.uuid === uuid))
      .map(async (chat) => {
        const data = { ...chat, responsible: newResponsible.value } as Chat;

        await store.update_chat(data);
        store.chats.set(chat!.uuid, data);
      });

    await Promise.all(promises);

    notification.success({
      title: "Responsible successfully changed",
      duration: 5000,
    });

    newResponsible.value = undefined;
    resetSelectedChats();
  } catch (error: any) {
    notification.error({
      title: (error as ConnectError).message ?? "[Error]: Unknown",
      duration: 5000,
    });
  } finally {
    isChangeResponsibleLoading.value = false;
  }
};

const checkedAdmins = ref<string[]>([]);
const checkedResponsibles = ref<string[]>([]);
const checkedDepartments = ref<string[]>([]);

const createDateRange = ref<[number, number] | null>(null);
const updateDateRange = ref<[number, number] | null>(null);

if (Object.keys(metrics.value).length < 1) fetch_defaults();

interface metricsOptionsType {
  [index: string]: [];
}

const filters = JSON.parse(localStorage.getItem("filters") ?? "null");
const metricsOptions = ref<metricsOptionsType>(
  filters
    ? filters.metrics
    : metrics.value.reduce((result, { key }) => ({ ...result, [key]: [] }), {})
);

watch(metrics, (value) => {
  const filters = JSON.parse(localStorage.getItem("filters") ?? "null");

  metricsOptions.value = filters
    ? filters.metrics
    : value.reduce((result, { key }) => ({ ...result, [key]: [] }), {});
});

watch(
  [sortBy, sortType],
  () => {
    localStorage.setItem(
      "sorting",
      JSON.stringify({
        sortBy: sortBy.value,
        sortType: sortType.value,
      })
    );
  },
  { deep: true }
);

watch(
  [
    checkedDepartments,
    checkedStatuses,
    checkedAdmins,
    checkedResponsibles,
    metricsOptions,
    updateDateRange,
    createDateRange,
  ],
  () => {
    localStorage.setItem(
      "filters",
      JSON.stringify({
        departments: checkedDepartments.value,
        statuses: checkedStatuses.value,
        admins: checkedAdmins.value,
        responsibles: checkedResponsibles.value,
        metrics: metricsOptions.value,
        updated: updateDateRange.value,
        created: createDateRange.value,
      })
    );
  },
  { deep: true }
);

async function resetFilters() {
  sortBy.value = "sent";
  sortType.value = "desc";

  checkedDepartments.value = [];
  checkedStatuses.value = [];
  checkedAdmins.value = [];
  checkedResponsibles.value = [];
  updateDateRange.value = null;
  createDateRange.value = null;
  metricsOptions.value = metrics.value.reduce(
    (result, { key }) => ({ ...result, [key]: [] }),
    {}
  );

  await nextTick();
  localStorage.removeItem("sorting");
  localStorage.removeItem("filters");
}

async function downloadReport() {
  isReportLoading.value = true;

  try {
    await ChatReportService.generate(chats.value, responsibles.value as User[]);
  } finally {
    isReportLoading.value = false;
  }
}

const isChatPanelOpen = ref(true);

function changePanelOpen() {
  isChatPanelOpen.value = !isChatPanelOpen.value;

  if (isChatPanelOpen) appStore.displayMode = "half";
  else appStore.displayMode = "none";
}

function changeMode(mode: string | null) {
  if (mode === "none" && appStore.displayMode === "half") {
    return;
  }

  if (mode) appStore.displayMode = mode;
  else if (appStore.displayMode === "half") {
    appStore.displayMode = "full";
    isChatPanelOpen.value = true;
  } else {
    appStore.displayMode = "half";
  }
}

const firstMessage = ref("gg");
const isFirstMessageVisible = ref(false);
const x = ref(0);
const y = ref(0);

function onMouseMove(clientX: number, clientY: number, chatId: string) {
  const chat = chats.value.find(({ uuid }) => uuid === chatId);

  x.value = clientX;
  y.value = clientY - 10;
  firstMessage.value = chat?.meta?.firstMessage?.content ?? "";

  if (firstMessage.value.length > 99) {
    firstMessage.value = `${firstMessage.value.slice(0, 100)}...`;
  }
  isFirstMessageVisible.value = true;
}

const resetSelectedChats = () => {
  isAllChatsSelected.value = false;
  selectedChats.value = [];
};

function choseAllChats() {
  selectedChats.value = viewedChats.value.map((chat) => chat.uuid);
}

watch(isAllChatsSelected, () => {
  if (isAllChatsSelected.value) {
    choseAllChats();
  } else {
    resetSelectedChats();
  }
});
</script>

<style>
#separator {
  cursor: col-resize;
  background-color: var(--n-color);
  background-repeat: no-repeat;
  background-position: center;
  width: 10px;
  /* Prevent the browser's built-in drag from interfering */
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.chats__panel {
  min-width: 450px;
  width: 450px;
  transition: 0.4s;
}

.chats__panel.closed {
  min-width: 70px !important;
  width: 70px !important;
  background-color: #18181c;
}

.chats__panel.closed .n-layout-sider-scroll-container {
  min-width: 70px !important;
  width: 70px !important;
}

.chats__panel .search {
  margin: 5px 10px;
}

.chats__panel .search div:nth-child(2) {
  width: 100%;
}

.chats__panel aside {
  width: 100% !important;
  max-width: 100% !important;
}

.chats__panel .n-list-item.active {
  background: var(--n-color-hover-modal);
}

.chats__panel .chat__actions {
  padding: 10px;
}

.chat__actions.hide {
  flex-flow: column-reverse !important;
  margin: 0px;
  padding: 10px;
}

.chats__filters {
  max-height: 800px;
  min-width: 40vw;
  max-width: 90vw;
}

.chat__item {
  height: 100%;
  min-width: 50%;
  width: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
