<template>
  <n-space class="title" align="center" justify="space-between">
    <n-h1>Statistics</n-h1>
    <n-button @click="goBack">Back</n-button>
  </n-space>

  <n-space vertical class="chart_container open_chats_count">
    <div class="chart_options">
      <n-space align="center">
        <n-h3> Number of open chats during the </n-h3>

        <n-space class="current_duration">
          <n-button @click="openChatsCountOffset--">
            <n-icon><prev-icon /></n-icon>
          </n-button>
          <span class="current_duration_info"
            >{{ formatDate(openChatCountDurationDatesTuple[0]) }} -
            {{ formatDate(openChatCountDurationDatesTuple[1]) }}</span
          >
          <n-button @click="openChatsCountOffset++">
            <n-icon><next-icon /></n-icon>
          </n-button>
        </n-space>
      </n-space>

      <n-space>
        <n-select
          style="width: 200px"
          v-model:value="openChatsCountType"
          :options="openChatsCountTypes"
        />

        <n-select
          v-if="openChatsCountType !== 'all'"
          style="width: 200px"
          v-model:value="openChatsCountByClient"
          :options="clientByOptions"
        />
        <n-select
          style="width: 200px"
          v-model:value="openChatsCountDuration"
          :options="durationOptions"
        />
      </n-space>
    </div>
    <div class="chart">
      <apexchart
        type="bar"
        height="250"
        :options="openChatsCountOptions"
        :series="openChatsCountOptions.series"
      ></apexchart>
    </div>
  </n-space>

  <n-space vertical class="chart_container closed_chats_count">
    <div class="chart_options">
      <n-space align="center">
        <n-h3> Number of closed chats during the </n-h3>

        <n-space class="current_duration">
          <n-button @click="closedChatsCountOffset--">
            <n-icon><prev-icon /></n-icon>
          </n-button>
          <span class="current_duration_info"
            >{{ formatDate(closedChatCountDurationDatesTuple[0]) }} -
            {{ formatDate(closedChatCountDurationDatesTuple[1]) }}</span
          >
          <n-button @click="closedChatsCountOffset++">
            <n-icon><next-icon /></n-icon>
          </n-button>
        </n-space>
      </n-space>

      <n-space>
        <n-select
          style="width: 200px"
          v-model:value="closedChatsCountType"
          :options="closedChatsCountTypes"
        />

        <n-select
          v-if="closedChatsCountType !== 'all'"
          style="width: 200px"
          v-model:value="closedChatsCountByClient"
          :options="clientByOptions"
        />
        <n-select
          style="width: 200px"
          v-model:value="closedChatsCountDuration"
          :options="durationOptions"
        />
      </n-space>
    </div>
    <div class="chart">
      <apexchart
        type="bar"
        height="250"
        :options="closedChatsCountOptions"
        :series="closedChatsCountOptions.series"
      ></apexchart>
    </div>
  </n-space>

  <n-space vertical class="chart_container closed_chats_count">
    <div class="chart_options">
      <n-space align="center">
        <n-h3> Activity in chats during the </n-h3>

        <n-space class="current_duration">
          <n-button @click="usersActivityOffset--">
            <n-icon><prev-icon /></n-icon>
          </n-button>
          <span class="current_duration_info"
            >{{ formatDate(usersActivityDurationDatesTuple[0]) }} -
            {{ formatDate(usersActivityDurationDatesTuple[1]) }}</span
          >
          <n-button @click="usersActivityOffset++">
            <n-icon><next-icon /></n-icon>
          </n-button>
        </n-space>
      </n-space>

      <n-space>
        <n-select
          style="width: 200px"
          v-model:value="usersActivityType"
          :options="usersActivityTypes"
        />

        <n-select
          style="width: 200px"
          v-model:value="usersActivityCount"
          :options="clientByOptions"
        />
        <n-select
          style="width: 200px"
          v-model:value="usersActivityDuration"
          :options="durationOptions"
        />
      </n-space>
    </div>
    <div class="chart">
      <apexchart
        v-if="!isChatsMessagesLoading"
        type="bar"
        height="250"
        :options="usersActivityOptions"
        :series="usersActivityOptions.series"
      ></apexchart>
      <n-skeleton v-else height="250px" width="100%" :sharp="false" />
    </div>
  </n-space>
</template>

<script lang="ts" setup>
import apexchart from "vue3-apexcharts";
import { useCcStore } from "../store/chatting";
import { computed, defineAsyncComponent, onMounted, watch } from "vue";
import { NSpace, NSelect, NButton, NIcon, NH1, NH3, NSkeleton } from "naive-ui";
import { ref } from "vue";
import { Chat, Messages } from "../connect/cc/cc_pb";
import { ApexOptions } from "apexcharts";
import { useRouter } from "vue-router";
import { useAppStore } from "../store/app";
import { useDefaultsStore } from "../store/defaults";
import { storeToRefs } from "pinia";
import { useUsersStore } from "../store/users";

const nextIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChevronForwardOutline")
);
const prevIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChevronBackOutline")
);

const store = useCcStore();
const appStore = useAppStore();
const router = useRouter();
const defaultsStore = useDefaultsStore();
const usersStore = useUsersStore();

const { admins } = storeToRefs(defaultsStore);
const { users } = storeToRefs(usersStore);

const months = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];
const apexchartColors = [
  "#FF6633",
  "#FFB399",
  "#FF33FF",
  "#FFFF99",
  "#00B3E6",
  "#E6B333",
  "#3366E6",
  "#999966",
  "#99FF99",
  "#B34D4D",
  "#80B300",
  "#809900",
  "#E6B3B3",
  "#6680B3",
  "#66991A",
  "#FF99E6",
  "#CCFF1A",
  "#FF1A66",
  "#E6331A",
  "#33FFCC",
  "#66994D",
];

const durationOptions = [
  { label: "Week", value: "week" },
  { label: "Month", value: "month" },
  { label: "Month by week", value: "month-week" },
  { label: "Year", value: "year" },
  { label: "Year by month", value: "year-month" },
];

const clientByOptions = [
  { label: "Top 5", value: "5" },
  { label: "Top 10", value: "10" },
  { label: "Top 25", value: "25" },
  { label: "Top 50", value: "50" },
];
const openChatsCountDuration = ref("month");
const openChatsCountType = ref("all");
const openChatsCountTypes = [
  { label: "Clients", value: "clients" },
  { label: "All", value: "all" },
];
const openChatsCountByClient = ref("5");
const openChatsCountOffset = ref(0);

const closedChatsCountDuration = ref("month");
const closedChatsCountType = ref("all");
const closedChatsCountTypes = [
  { label: "Clients", value: "clients" },
  { label: "All", value: "all" },
];
const closedChatsCountByClient = ref("5");
const closedChatsCountOffset = ref(0);

const chatMessagesMap = ref(new Map<string, Messages>());
const isChatsMessagesLoading = ref(false);

const usersActivityDuration = ref("month");
const usersActivityType = ref("only-admins");
const usersActivityTypes = [
  { label: "Clients", value: "only-clients" },
  { label: "Admins", value: "only-admins" },
];
const usersActivityCount = ref("5");
const usersActivityOffset = ref(0);

onMounted(async () => {
  try {
    isChatsMessagesLoading.value = true;
    await Promise.all(
      chats.value.map(async (chat) => {
        const messagesData = await store.get_messages(chat as Chat, false);
        chatMessagesMap.value.set(chat.uuid, messagesData);
      })
    );
  } finally {
    isChatsMessagesLoading.value = false;
  }
});

const chats = computed(() => {
  const chats = [...store.chats.values()];

  chats.sort((a, b) => {
    if (a.created < b.created) {
      return -1;
    }
    if (a.created > b.created) {
      return 1;
    }
    return 0;
  });

  return chats;
});

const openChatCountDurationDatesTuple = computed(() => {
  const [startDate, endDate] = getDurationTuple(
    openChatsCountDuration.value,
    openChatsCountOffset.value
  );

  return [startDate, endDate];
});
const openChatsCountOptions = computed<ApexOptions>(() => {
  const options = getChartOptions({
    duration: openChatsCountDuration.value,
    durationTuples: openChatCountDurationDatesTuple.value,
    title: "New tickets",
    type: openChatsCountType.value,
    chatsFilter:
      openChatsCountType.value === "clients"
        ? (params) =>
            filterChatsByTop({
              ...params,
              count: +openChatsCountByClient.value,
            })
        : filterChats,
    setToChartMap:
      openChatsCountType.value === "all"
        ? setToChatCountMap
        : setToChatAccountMap,
  });
  return options;
});

const closedChatCountDurationDatesTuple = computed(() => {
  const [startDate, endDate] = getDurationTuple(
    closedChatsCountDuration.value,
    closedChatsCountOffset.value
  );

  return [startDate, endDate];
});
const closedChatsCountOptions = computed<ApexOptions>(() => {
  const options = getChartOptions({
    duration: closedChatsCountDuration.value,
    durationTuples: closedChatCountDurationDatesTuple.value,
    title: "Closed tickets",
    type: closedChatsCountType.value,
    chatsFilter:
      closedChatsCountType.value === "clients"
        ? (params) =>
            filterChatsByTop({
              ...params,
              count: +closedChatsCountByClient.value,
            })
        : filterChats,
    filter: (chat) => chat.status === 3 || chat.status === 2,
    setToChartMap:
      closedChatsCountType.value === "all"
        ? setToChatCountMap
        : setToChatAccountMap,
  });
  return options;
});

const usersActivityDurationDatesTuple = computed(() => {
  const [startDate, endDate] = getDurationTuple(
    usersActivityDuration.value,
    usersActivityOffset.value
  );

  return [startDate, endDate];
});
const usersActivityOptions = computed<ApexOptions>(() => {
  const options = getChartOptions({
    duration: usersActivityDuration.value,
    durationTuples: usersActivityDurationDatesTuple.value,
    title: "Activity in chats",
    type: usersActivityType.value,
    filter: (chat) => chat.status === 3 || chat.status === 2,
    chatsFilter: filterChats,
    seriesFilter: (series) =>
      filterUsersActivitySeries(series, +usersActivityCount.value),
    setToChartMap: (params) =>
      setToChatActivityMap({
        ...params,
        type: usersActivityType.value,
      }),
  });
  return options;
});

function getChartOptions({
  type = "all",
  duration = "monthly",
  durationTuples = [new Date(), new Date()],
  title = "Tickets",
  filter = (chat: Chat) => !!chat,
  chatsFilter = (params: ChatsFilterParams) => params.chats,
  seriesFilter = (series: ApexAxisChartSeries) => series,
  setToChartMap = setToChatCountMap,
}): ApexOptions {
  const isTypeAll = type === "all";
  let filtredChats = chats.value;
  const startDate = new Date(durationTuples[0].getTime());
  const endDate = new Date(durationTuples[1].getTime());

  filtredChats = chatsFilter({
    chats: filtredChats as Chat[],
    endDate,
    startDate,
  });

  let categories: string[] = [];

  switch (duration) {
    case "week": {
      const firstDayTs = startDate.getTime();

      categories = Array.from({ length: 7 }, (_, i) => {
        const date = new Date(firstDayTs + 86400000 * i);
        const dayNumber = date.getDate();
        const currentMonth = date.getMonth();

        return `${dayNumber} ${months[currentMonth]}`;
      });

      break;
    }
    case "month-week": {
      const daysInMonth = new Date(
        startDate.getFullYear(),
        startDate.getMonth() + 1,
        0
      ).getDate();

      categories = Array.from(
        { length: Math.ceil(daysInMonth / 7) },
        (_, i) => {
          let startDay: any = new Date(startDate.getTime());
          startDate.setDate(i * 7 + 1);
          startDay = startDate.getDate();
          let endDay: any = new Date(startDate.getTime());
          endDay.setDate(Math.min((i + 1) * 7, endDate.getDate()));
          endDay = endDay.getDate();

          return `${startDay}-${endDay} ${months[startDate.getMonth()]}`;
        }
      );
      break;
    }
    case "month": {
      categories = [months[startDate.getMonth() + 1]];
      break;
    }
    case "year-month": {
      categories = [...months];
      break;
    }

    case "year":
    default: {
      categories = [startDate.getFullYear().toString()];
      break;
    }
  }
  let series: ApexAxisChartSeries = [];

  if (isTypeAll) {
    series.push({
      name: title,
      data: Array.from({ length: categories.length }, () => 0),
    });
  }

  const seriesMap = new Map();
  switch (duration) {
    case "week": {
      const firstDayTs = startDate.getTime();
      const lastDayTs = endDate.getTime();

      filtredChats.forEach((chat) => {
        const createdDateTs = Number(chat.created);

        if (
          !(createdDateTs > firstDayTs) ||
          !(createdDateTs < lastDayTs) ||
          !filter(chat as Chat)
        ) {
          return;
        }

        const currentDay = new Date(
          Number(chat.created) - firstDayTs
        ).getDate();

        setToChartMap({ map: seriesMap, key: currentDay, chat: chat as Chat });
      });

      break;
    }
    case "month-week": {
      const currentMonth = startDate.getMonth();

      filtredChats.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (createdDate.getMonth() !== currentMonth || !filter(chat as Chat)) {
          return;
        }

        const currentDay = Math.ceil(createdDate.getDate() / 7);
        setToChartMap({ map: seriesMap, key: currentDay, chat: chat as Chat });
      });
      break;
    }
    case "month": {
      filtredChats.forEach((chat) => {
        setToChartMap({ map: seriesMap, key: 1, chat: chat as Chat });
      });
      break;
    }
    case "year-month": {
      const currentYear = startDate.getFullYear();
      filtredChats.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (
          createdDate.getFullYear() !== currentYear ||
          !filter(chat as Chat)
        ) {
          return;
        }

        const currentMonth = createdDate.getMonth();
        setToChartMap({
          map: seriesMap,
          key: currentMonth,
          chat: chat as Chat,
        });
      });

      break;
    }
    case "year":
    default: {
      filtredChats.forEach((chat) => {
        setToChartMap({
          map: seriesMap,
          key: 1,
          chat: chat as Chat,
        });
      });

      break;
    }
  }

  [...seriesMap.keys()].forEach((key) => {
    if (isTypeAll) {
      series[0].data[key - 1] = seriesMap.get(key);
    } else {
      [...seriesMap.get(key).keys()].forEach((ownerUuid: string) => {
        let existsIndex = series.findIndex((serie) => serie.name === ownerUuid);

        if (existsIndex === -1) {
          existsIndex = series.length;
          series.push({
            name: ownerUuid,
            data: Array.from({ length: categories.length }, () => 0),
            color: apexchartColors[existsIndex],
          });
        }

        series[existsIndex].data[key - 1] = seriesMap.get(key).get(ownerUuid);
      });
    }
  });

  series = seriesFilter(series);

  if (!isTypeAll) {
    series = series.map((s) => ({
      ...s,
      name: users.value.get(s.name || "")?.title || "Unknown",
    }));
  }

  return {
    series,
    xaxis: {
      categories,
      position: "top",
      tickPlacement: "on",
    },
    yaxis: {
      labels: {
        rotate: 0,
        formatter: function (val: number) {
          return val?.toFixed(0);
        },
      },
      decimalsInFloat: 0,
    },
    theme: {
      mode: "dark",
      palette: "palette4",
    },
    dataLabels: {
      enabled: isTypeAll,
    },
    legend: {
      onItemClick: {
        toggleDataSeries: false,
      },
      show: true,
      showForSingleSeries: true,
      formatter: function (val: any, opts: any) {
        const total = series[opts.seriesIndex].data.reduce((a, acc) => {
          return (acc as number) + (a as number);
        }, 0);
        return `${val}: ${total}`;
      },
    },
    chart: {
      animations: {
        enabled: isTypeAll,
      },
      zoom: {
        type: "x",
        enabled: !isTypeAll,
        autoScaleYaxis: true,
      },
      toolbar: {
        tools: {
          download: true,
          zoom: !isTypeAll,
          zoomin: !isTypeAll,
          zoomout: !isTypeAll,
        },
      },
    },
    title: {
      text: "New tickets",
      floating: false,
      offsetY: 330,
      align: "center",
      style: {
        color: "#444",
      },
    },
  };
}

function filterChatsByTop({
  chats,
  startDate,
  endDate,
  count,
}: ChatsFilterParams & { count: number }) {
  const mapMostPopularOwner = new Map();

  chats
    .filter((chat) => {
      const created = Number(chat.created);
      return created > startDate.getTime() && created < endDate.getTime();
    })
    .forEach((chat) => {
      let count = 0;
      if (mapMostPopularOwner.has(chat.owner)) {
        count = mapMostPopularOwner.get(chat.owner);
      }
      count++;
      mapMostPopularOwner.set(chat.owner, count);
    });

  const arrayMostPopularOwner = [...mapMostPopularOwner.keys()]
    .map((key) => ({ owner: key, count: mapMostPopularOwner.get(key) }))
    .sort((a, b) => b.count - a.count)
    .slice(0, +count)
    .map(({ owner }) => owner);

  return chats.filter((chat) => arrayMostPopularOwner.includes(chat.owner));
}

function filterChats({ chats, endDate, startDate }: ChatsFilterParams) {
  return chats.filter((chat) => {
    const created = Number(chat.created);
    return created > startDate.getTime() && created < endDate.getTime();
  });
}

interface ChatsFilterParams {
  chats: Chat[];
  startDate: Date;
  endDate: Date;
}

function filterUsersActivitySeries(series: ApexAxisChartSeries, count: number) {
  const userActivityMap = new Map();
  series.forEach((s) => {
    userActivityMap.set(
      s.name,
      s.data.reduce((acc, a) => (acc as number) + (a as number), 0)
    );
  });

  const mostActiveUsers = [...userActivityMap.keys()]
    .map((key) => ({ user: key, activity: userActivityMap.get(key) }))
    .sort((a, b) => b.activity - a.activity)
    .slice(0, +count)
    .map(({ user }) => user);

  return series.filter((s) => mostActiveUsers.includes(s.name));
}

function getDurationTuple(type = "week", offset = 1) {
  let startDate, endDate;

  switch (type) {
    case "week": {
      let [firstDayTs, lastDayTs] = getWeekTuple();

      firstDayTs += offset * 86400000 * 7;
      lastDayTs += offset * 86400000 * 7;

      startDate = new Date(firstDayTs);
      endDate = new Date(lastDayTs);

      break;
    }
    case "month-week":
    case "month": {
      const todayDate = new Date(Date.now());
      todayDate.setMonth(todayDate.getMonth() + offset);

      startDate = new Date(todayDate.setDate(1));
      endDate = new Date(
        todayDate.setDate(
          new Date(
            todayDate.getFullYear(),
            todayDate.getMonth() + 1,
            0
          ).getDate()
        )
      );

      break;
    }
    case "year":
    case "year-month":
    default: {
      const todayDate = new Date(Date.now());
      const year = todayDate.getFullYear() + offset;

      startDate = new Date(todayDate.setFullYear(year, 0, 1));
      endDate = new Date(todayDate.setFullYear(year, 11, 31));
      break;
    }
  }
  return [startDate, endDate];
}

interface SetToChartMap {
  map: Map<string | number, any>;
  key: string | number;
  chat: Chat;
}

function setToChatCountMap({ key, map }: SetToChartMap) {
  let period = 0;
  if (map.has(key)) {
    period = map.get(key);
  }
  period++;
  map.set(key, period);
}

function setToChatAccountMap({ chat, key, map }: SetToChartMap) {
  let period = new Map();
  if (map.has(key)) {
    period = map.get(key);
  }

  let countForPeriod = 0;

  if (period.has(chat.owner)) {
    countForPeriod = period.get(chat.owner);
  }

  countForPeriod++;

  period.set(chat.owner, countForPeriod);
  map.set(key, period);
}

function setToChatActivityMap({
  chat,
  key,
  map,
  type,
}: SetToChartMap & { type: string }) {
  if (!chatMessagesMap.value.has(chat.uuid)) {
    return;
  }

  chatMessagesMap.value.get(chat.uuid)?.messages.map((message) => {
    const isAdmin = admins.value.includes(message.sender);

    if (
      (type === "only-clients" && isAdmin) ||
      (type === "only-admins" && !isAdmin)
    ) {
      return;
    }

    let activity = new Map();
    if (map.has(key)) {
      activity = map.get(key);
    }

    let countForPeriod = 0;

    if (activity.has(message.sender)) {
      countForPeriod = activity.get(message.sender);
    }

    countForPeriod++;

    activity.set(message.sender, countForPeriod);
    map.set(key, activity);
  });
}

function getWeekTuple() {
  const date1 = new Date();
  const date2 = new Date();
  const today = date1.getDate();
  const currentDay = date1.getDay();
  const firstDay = date1.setDate(today - (currentDay || 7));
  const lastDay = date2.setDate(today - currentDay + 7);

  return [firstDay, lastDay];
}

function formatDate(date: Date) {
  return date.toISOString().split("T")[0].replaceAll("-", "/");
}

function goBack() {
  router.go(-1);
}

function openUser(e: Event) {
  const name = (e.target as any).innerText.split(":")[0];
  const user = [...users.value.values()].find((u) => u.title === name);
  if (user) {
    window.open(
      `${appStore.conf?.params.fullUrl.split("plugin")[0]}accounts/${
        user.uuid
      }?tab=helpdesk`,
      "_blank"
    );
  }
}

watch(openChatsCountDuration, () => {
  openChatsCountOffset.value = 0;
});
watch(closedChatsCountDuration, () => {
  closedChatsCountOffset.value = 0;
});
watch(usersActivityDuration, () => {
  usersActivityOffset.value = 0;
});

watch(
  [openChatsCountOptions, closedChatsCountOptions, usersActivityOptions],
  () => {
    setTimeout(() => {
      document
        .querySelectorAll(".apexcharts-legend-series")
        .forEach((element) => {
          element.removeEventListener("click", openUser);
          element.addEventListener("click", openUser);
        });
    }, 200);
  },
  { deep: true }
);
</script>

<style lang="scss">
span.apexcharts-legend-text {
  color: var(--v-primary-base) !important;
}

.title {
  margin: 0px 20px;

  h1 {
    margin: 0px;
  }
}

.chart_container {
  margin: 0px 20px;

  h3 {
    margin: 0px;
  }

  .chart_options {
    display: flex;
    justify-content: space-between;
    .current_duration {
      .current_duration_info {
        font-size: 1.3rem;
        text-align: center;
      }
    }
  }
}

.apexcharts-legend-text {
  text-decoration: underline;
  cursor: pointer;
}
</style>
