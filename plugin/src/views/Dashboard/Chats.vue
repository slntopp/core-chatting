<template>
  <div
    v-if="appStore.displayMode !== 'none'"
    :class="{ chats__panel: true, closed: !isChatPanelOpen }"
    :style="appStore.displayMode === 'full' ? 'min-width: 100%' : undefined"
  >
    <n-layout-sider
      style="margin-bottom: 25px"
      collapse-mode="transform"
      :collapsed="!isChatPanelOpen"
    >
      <n-space
        align="center"
        :wrap="appStore.isMobile ? false : true"
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

          <n-button
            v-if="!appStore.isMobile"
            ghost
            type="info"
            @click="goToStatistics"
            :disabled="isLoading"
          >
            <n-icon :component="StatsChartOutline" />
            <span v-if="isChatPanelOpen" style="margin-left: 5px">
              Statistic
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
            clearable
            placeholder="Status"
            :value="selectedStatus"
            :options="
              Object.values(chatsCountByStatus).map(({ count, status }) => ({
                value: +status,
                label: `${getStatus(status)} (${count})`,
              }))
            "
            @update:value="selectStatus"
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
              :loading="isUsersLoading"
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
          :style="appStore.isMobile ? undefined : 'padding-left: 10px'"
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
          v-model:show="isFilterOpen"
        >
          <template #trigger>
            <n-icon size="24" style="vertical-align: middle; cursor: pointer">
              <sort-icon />
            </n-icon>
          </template>

          <chats-filters
            v-model:checkedDepartments="checkedDepartments"
            v-model:checkedStatuses="checkedStatuses"
            v-model:checkedAdmins="checkedAdmins"
            v-model:checkedResponsibles="checkedResponsibles"
            :update-date-range="updateDateRange"
            :create-date-range="createDateRange"
            @update:create-date-range="createDateRange = { ...$event }"
            @update:update-date-range="updateDateRange = { ...$event }"
            v-model:sort-by="sortBy"
            v-model:sort-type="sortType"
            @reset="resetFilters"
            @close="isFilterOpen = false"
            :metricsOptions="metricsOptions"
            @update:checked-metrics="metricsOptions = $event"
          />
        </n-popover>
      </n-space>

      <n-scrollbar
        ref="scrollbar"
        :class="{
          scrollBarClosed: appStore.displayMode !== 'half',
          scrollBarOpened: appStore.displayMode === 'half',
        }"
      >
        <n-popover trigger="manual" :show="isFirstMessageVisible" :x="x" :y="y">
          {{ firstMessage }}
        </n-popover>

        <span ref="topPanel"></span>

        <div v-if="isLoading || isDefaultLoading" class="chats_spin">
          <n-spin size="large" />
        </div>
        <div class="no_chats_message" v-else-if="chats.length === 0">
          <n-alert type="info" title="No chats available">
            Try another search request
          </n-alert>
        </div>
        <n-list v-else hoverable clickable style="margin-bottom: 25px">
          <chat-item
            v-for="chat in chats"
            :selected="selectedChats"
            :hide-message="!isChatPanelOpen"
            :uuid="chat.uuid"
            :chat="chat as Chat"
            :chats="chats as Chat[]"
            :class="{
              active: chat.uuid === router.currentRoute.value.params.uuid,
            }"
            :key="chat.uuid"
            @click="changeMode('none')"
            @hover="onMouseMove"
            @hoverEnd="isFirstMessageVisible = false"
            @select="selectChat"
          />
        </n-list>
      </n-scrollbar>

      <div class="chats_pagination" v-if="chats.length">
        <n-pagination
          :disabled="isLoading || isDefaultLoading"
          v-model:page="page"
          v-model:page-size="pageSize"
          :page-count="pageCount"
          show-size-picker
          :page-sizes="[10, 20, 30, 50]"
        />
      </div>
    </n-layout-sider>
  </div>

  <div id="separator" v-show="appStore.displayMode === 'half'"></div>
  <div class="chat__item" v-show="appStore.displayMode !== 'full'">
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
  NScrollbar,
  NSpace,
  NPopover,
  NTag,
  NText,
  NSpin,
  NPagination,
  NSelect,
  NCheckbox,
  NAlert,
  useNotification,
} from "naive-ui";

import { useAppStore } from "../../store/app.ts";
import { useCcStore } from "../../store/chatting.ts";

import { useRoute, useRouter } from "vue-router";
import {
  Chat,
  CountChatsRequest,
  ListChatsRequest,
  Status,
} from "../../connect/cc/cc_pb";
import ChatItem from "../../components/chats/chat_item.vue";
import ChatsFilters from "../../components/chats/chats_filters.vue";
import useDraggable from "../../hooks/useDraggable.ts";
import {
  cleanObject,
  debounce,
  getStatusColor,
  getStatusItems,
} from "../../functions.ts";
import { ConnectError } from "@connectrpc/connect";
import { storeToRefs } from "pinia";
import { useDefaultsStore } from "../../store/defaults.ts";
import { useUsersStore } from "../../store/users.ts";

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
const StatsChartOutline = defineAsyncComponent(
  () => import("@vicons/ionicons5/StatsChartOutline")
);
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseCircle")
);
const mergeIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/GitMerge")
);

const SORTING_KEY = "chats-sorting";
const FILTERS_KEY = "chats-filters";

const appStore = useAppStore();
const store = useCcStore();
const router = useRouter();
const route = useRoute();
const { makeDraggable } = useDraggable();
const defaultsStore = useDefaultsStore();
const { metrics, isDefaultLoading, admins } = storeToRefs(defaultsStore);
const usersStore = useUsersStore();
const { isUsersLoading, users } = storeToRefs(usersStore);
const notification = useNotification();

const topPanel = ref<any>(null);

const selectedChats = ref<string[]>([]);
const deleteLoading = ref(false);
const mergeLoading = ref(false);
const searchParam = ref("");
const listSearchParam = ref("");
const scrollbar = ref<InstanceType<typeof NScrollbar>>();
const page = ref(1);
const pageSize = ref(10);
const isLoading = ref(false);
const newStatus = ref();
const isChangeStatusLoading = ref(false);
const newDepartment = ref();
const isChangeDepartmentLoading = ref(false);
const newResponsible = ref();
const isChangeResponsibleLoading = ref(false);
const isAllChatsSelected = ref(false);
const isFilterOpen = ref(false);

const isMergeVisible = computed(() => {
  return isSelectedChats.value;
});

const isSelectedChats = computed(() => selectedChats.value.length > 0);

const departments = computed(() =>
  store.departments.map(({ key, title }) => ({ label: title, value: key }))
);

const responsibles = computed(() =>
  admins.value
    .map((admin) => users.value.get(admin))
    .filter((u) => !!u)
    .map((user) => ({
      label: user?.title,
      value: user?.uuid,
    }))
);

async function sync() {
  try {
    isLoading.value = true;
    fetch_chats_debounced();
  } catch (error) {
    console.log(error);
  } finally {
    isLoading.value = false;
  }
}

const fetch_chats = async () => {
  try {
    isLoading.value = true;
    await store.list_chats(ListChatsRequest.fromJson(chatListOptions.value));
  } catch (error) {
    console.log(error);
  } finally {
    isLoading.value = false;
  }
};

const fetch_chats_debounced = debounce(fetch_chats, 100);

function startChat() {
  router.push({ name: "Start Chat" });
  appStore.displayMode = "none";
}

function goToStatistics() {
  window.top?.postMessage({ type: "open-chats-statistics", value: {} }, "*");
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

  const sorting = JSON.parse(
    localStorage.getItem(SORTING_KEY) ?? '{"sortBy":"unread","sortType":"desc"}'
  );
  const filters = JSON.parse(
    localStorage.getItem(FILTERS_KEY) ??
      `{"admins":[],"departments":[],"responsibles":[],"statuses":[0,1,4,5,7,8]}`
  );

  if (sorting) {
    sortBy.value = sorting.sortBy;
    sortType.value = sorting.sortType;
    pageSize.value = sorting.pageSize ?? 10;
  }
  if (filters) {
    checkedDepartments.value = filters.departments ?? [];
    checkedStatuses.value = filters.statuses ?? [];
    checkedAdmins.value = filters.admins ?? [];
    checkedResponsibles.value = filters.responsibles ?? [];
    updateDateRange.value = filters.updated ?? { from: null, to: null };
    createDateRange.value = filters.created ?? { from: null, to: null };
    metricsOptions.value = filters.metrics ?? {};
  }

  store.list_chats_count(
    appStore.conf?.params?.filterByAccount
      ? CountChatsRequest.fromJson({
          filters: { account: appStore.conf?.params?.filterByAccount },
        })
      : new CountChatsRequest()
  );
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

const sortBy = ref<"none" | "sent" | "created" | "status">("sent");
const sortType = ref<"asc" | "desc">("desc");

const checkedStatuses = ref<Status[]>([]);
const selectedStatus = ref<Status>();
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

const chatFilters = computed(() => {
  let filters = {
    status:
      selectedStatus.value === undefined
        ? checkedStatuses.value
        : [selectedStatus.value],
    department: checkedDepartments.value,
    admins: checkedAdmins.value,
    responsible: checkedResponsibles.value,
    account: appStore.conf?.params?.filterByAccount,
    updated: updateDateRange.value,
    created: createDateRange.value,
    search_param: listSearchParam.value,
  };

  return cleanObject(filters);
});

const chatListOptions = computed(() => {
  return {
    limit: pageSize.value,
    page: page.value,
    filters: chatFilters.value,
    field: sortBy.value,
    sort: sortType.value,
  };
});

const chats = computed(() => {
  return [...store.chats.values()];
});

const totalChats = computed(() => {
  return store.totalChats;
});

const pageCount = computed(() => Math.ceil(totalChats.value / pageSize.value));

const chatsCountByStatus = computed(() => {
  const order = [0, 1, 5, 8, 4, 7, 3];

  return [...store.chats_count.keys()]
    .map<{ count: number; status: number }>((key) => ({
      count: store.chats_count.get(key) || 0,
      status: +key,
    }))
    .sort(
      (a, b) =>
        order.findIndex((v) => v === a.status) -
        order.findIndex((v) => v === b.status)
    );
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

const createDateRange = ref<{ from: null | number; to: null | number }>({
  from: null,
  to: null,
});
const updateDateRange = ref<{ from: null | number; to: null | number }>({
  from: null,
  to: null,
});

if (Object.keys(metrics.value).length < 1) defaultsStore.fetch_defaults();

interface metricsOptionsType {
  [index: string]: [];
}

const metricsOptions = ref<metricsOptionsType>(
  metrics.value.reduce((result, { key }) => ({ ...result, [key]: [] }), {})
);

const saveSortingAndFilters = debounce(() => {
  localStorage.setItem(
    FILTERS_KEY,
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
  localStorage.setItem(
    SORTING_KEY,
    JSON.stringify({
      sortBy: sortBy.value,
      sortType: sortType.value,
      pageSize: pageSize.value,
    })
  );
}, 200);

watch(metrics, (value) => {
  const filters = JSON.parse(localStorage.getItem(FILTERS_KEY) ?? "null");

  metricsOptions.value = filters
    ? filters.metrics
    : value.reduce((result, { key }) => ({ ...result, [key]: [] }), {});
});

watch(
  [sortBy, sortType],
  () => {
    saveSortingAndFilters();
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
    selectedStatus.value = undefined;

    saveSortingAndFilters();
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
  updateDateRange.value = { from: null, to: null };
  createDateRange.value = { from: null, to: null };
  metricsOptions.value = metrics.value.reduce(
    (result, { key }) => ({ ...result, [key]: [] }),
    {}
  );

  await nextTick();
  localStorage.removeItem(SORTING_KEY);
  localStorage.removeItem(FILTERS_KEY);
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
  selectedChats.value = chats.value.map((chat) => chat.uuid);
}

const changeListSearchParam = debounce(
  () => (listSearchParam.value = searchParam.value),
  600
);

watch(isAllChatsSelected, () => {
  if (isAllChatsSelected.value) {
    choseAllChats();
  } else {
    resetSelectedChats();
  }
});

watch([chatListOptions], fetch_chats_debounced, { deep: true });
watch([pageSize, chatFilters], () => (page.value = 1), { deep: true });
watch(searchParam, changeListSearchParam);

watch(page, () => {
  topPanel.value.scrollIntoView({ behavior: "smooth" });
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

.chats_spin {
  height: 80vh;
  width: 100%;
  display: flex;
  justify-items: center;
  align-items: center;
  justify-content: center;
}

.chats_pagination {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-bottom: 10px;
}

.no_chats_message {
  height: 80vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.scrollBarOpened {
  height: calc(100vh - 220px);
  min-width: 150px;
}

.scrollBarClosed {
  height: calc(100vh - 155px);
  min-width: 150px;
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
