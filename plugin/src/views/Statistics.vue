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
            >{{ formatDate(chatCountDurationDatesTuple[0]) }} -
            {{ formatDate(chatCountDurationDatesTuple[1]) }}</span
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
          :options="openChatsCountByClientTypes"
        />
        <n-select
          style="width: 200px"
          v-model:value="openChatsCountDuration"
          :options="openChatsCountDurations"
        />
      </n-space>
    </div>
    <div class="chart">
      <apexchart
        type="bar"
        height="250"
        :options="chatsCountOptions"
        :series="chatsCountOptions.series"
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
          :options="closedChatsCountByClientTypes"
        />
        <n-select
          style="width: 200px"
          v-model:value="closedChatsCountDuration"
          :options="closedChatsCountDurations"
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
          :options="usersActivityCountTypes"
        />
        <n-select
          style="width: 200px"
          v-model:value="usersActivityDuration"
          :options="usersActivityDurations"
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
import useDefaults from "../hooks/useDefaults";

const nextIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChevronForwardOutline")
);
const prevIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/ChevronBackOutline")
);

const store = useCcStore();
const router = useRouter();
const { admins, fetch_defaults } = useDefaults();

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

const openChatsCountDuration = ref("weekly");
const openChatsCountDurations = [
  { label: "Weekly", value: "weekly" },
  { label: "Monthly", value: "monthly" },
  { label: "Yearly", value: "yearly" },
];
const openChatsCountType = ref("all");
const openChatsCountTypes = [
  { label: "Clients", value: "clients" },
  { label: "All", value: "all" },
];
const openChatsCountByClient = ref("5");
const openChatsCountByClientTypes = [
  { label: "Top 5", value: "5" },
  { label: "Top 10", value: "10" },
  { label: "Top 20", value: "20" },
];
const openChatsCountOffset = ref(0);

const closedChatsCountDuration = ref("weekly");
const closedChatsCountDurations = [
  { label: "Weekly", value: "weekly" },
  { label: "Monthly", value: "monthly" },
  { label: "Yearly", value: "yearly" },
];
const closedChatsCountType = ref("all");
const closedChatsCountTypes = [
  { label: "Clients", value: "clients" },
  { label: "All", value: "all" },
];
const closedChatsCountByClient = ref("5");
const closedChatsCountByClientTypes = [
  { label: "Top 5", value: "5" },
  { label: "Top 10", value: "10" },
  { label: "Top 20", value: "20" },
];
const closedChatsCountOffset = ref(0);

const chatMessagesMap = ref(new Map<string, Messages>());
const isChatsMessagesLoading = ref(false);

const usersActivityDuration = ref("weekly");
const usersActivityDurations = [
  { label: "Weekly", value: "weekly" },
  { label: "Monthly", value: "monthly" },
  { label: "Yearly", value: "yearly" },
];
const usersActivityType = ref("only-admins");
const usersActivityTypes = [
  { label: "Clients", value: "only-clients" },
  { label: "Admins", value: "only-admins" },
];
const usersActivityCount = ref("5");
const usersActivityCountTypes = [
  { label: "Top 5", value: "5" },
  { label: "Top 10", value: "10" },
  { label: "Top 20", value: "20" },
];
const usersActivityOffset = ref(0);

onMounted(async () => {
  try {
    isChatsMessagesLoading.value = true;

    fetch_defaults();
    await Promise.all(
      chats.value.map(async (chat) => {
        const messagesData = await store.get_messages(chat, false);

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

const chatCountDurationDatesTuple = computed(() => {
  const [startDate, endDate] = getDurationTuple(
    openChatsCountDuration.value,
    openChatsCountOffset.value
  );

  return [startDate, endDate];
});
const chatsCountOptions = computed<ApexOptions>(() => {
  const options = getChartOptions({
    duration: openChatsCountDuration.value,
    durationTuples: chatCountDurationDatesTuple.value,
    title: "New tickets",
    type: openChatsCountType.value,
    chatsFilter:
      closedChatsCountType.value === "clients"
        ? (params) =>
            filterChatsByTop({ ...params, count: +closedChatsCountType.value })
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
            filterChatsByTop({ ...params, count: +closedChatsCountType.value })
        : filterChats,
    stacked: true,
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
    stacked: true,
    filter: (chat) => chat.status === 3 || chat.status === 2,
    chatsFilter: filterChats,
    setToChartMap: (params) =>
      setToChatActivityMap({
        ...params,
        type: usersActivityType.value,
      }),
  });
  return options;
});

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

function getChartOptions({
  type = "all",
  duration = "monthly",
  durationTuples = [new Date(), new Date()],
  title = "Tickets",
  filter = (chat: Chat) => !!chat,
  stacked = false,
  chatsFilter = (params: ChatsFilterParams) => params.chats,
  setToChartMap = setToChatCountMap,
}): ApexOptions {
  const isTypeAll = type === "all";
  let filtredChats = chats.value;
  const startDate = new Date(durationTuples[0].getTime());
  const endDate = new Date(durationTuples[1].getTime());

  filtredChats = chatsFilter({
    chats: chats.value,
    endDate,
    startDate,
  });

  let categories: string[] = [];

  switch (duration) {
    case "weekly": {
      const firstDayTs = startDate.getTime();

      categories = Array.from({ length: 7 }, (_, i) => {
        const date = new Date(firstDayTs + 86400000 * i);
        const dayNumber = date.getDate();
        const currentMonth = date.getMonth();

        return `${dayNumber} ${months[currentMonth]}`;
      });

      break;
    }
    case "monthly": {
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
    case "yearly":
    default: {
      categories = [...months];
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
    case "weekly": {
      const firstDayTs = startDate.getTime();
      const lastDayTs = endDate.getTime();

      filtredChats.forEach((chat) => {
        const createdDateTs = Number(chat.created);

        if (
          !(createdDateTs > firstDayTs) ||
          !(createdDateTs < lastDayTs) ||
          !filter(chat)
        ) {
          return;
        }

        const currentDay = new Date(
          Number(chat.created) - firstDayTs
        ).getDate();

        setToChartMap({ map: seriesMap, key: currentDay, chat });
      });

      break;
    }
    case "monthly": {
      const currentMonth = startDate.getMonth();

      filtredChats.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (createdDate.getMonth() !== currentMonth || !filter(chat)) {
          return;
        }

        const currentDay = Math.ceil(createdDate.getDate() / 7);
        setToChartMap({ map: seriesMap, key: currentDay, chat });
      });
      break;
    }
    case "yearly":
    default: {
      const currentYear = startDate.getFullYear();
      filtredChats.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (createdDate.getFullYear() !== currentYear || !filter(chat)) {
          return;
        }

        const currentMonth = createdDate.getMonth();
        setToChartMap({
          map: seriesMap,
          key: currentMonth,
          chat,
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

  if (!isTypeAll) {
    series = series.map((s) => ({
      ...s,
      name: store.users.get(s.name || "")?.title || "Unknown",
    }));
  }

  return {
    series,
    xaxis: {
      categories,
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
      show: isTypeAll,
    },
    chart: {
      stacked: stacked,
      animations: {
        enabled: isTypeAll,
      },
      toolbar: {
        tools: {
          download: true,
          zoom: true,
          zoomin: true,
          zoomout: true,
          pan: true,
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

function getDurationTuple(type = "week", offset = 1) {
  let startDate, endDate;

  switch (type) {
    case "weekly": {
      let [firstDayTs, lastDayTs] = getWeekTuple();

      firstDayTs += offset * 86400000 * 7;
      lastDayTs += offset * 86400000 * 7;

      startDate = new Date(firstDayTs);
      endDate = new Date(lastDayTs);

      break;
    }
    case "monthly": {
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
    case "yearly":
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

watch(openChatsCountDuration, () => {
  openChatsCountOffset.value = 0;
});
watch(closedChatsCountDuration, () => {
  closedChatsCountOffset.value = 0;
});
watch(usersActivityDuration, () => {
  usersActivityOffset.value = 0;
});
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
</style>
