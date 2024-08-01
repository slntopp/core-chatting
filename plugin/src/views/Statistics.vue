<template>
  <span>Statistics</span>

  <n-space vertical class="chats_count">
    <div class="chart_options">
      <n-select
        style="width: 200px"
        v-model:value="chatsCountDuration"
        :options="chatsCountDurations"
      />
    </div>
    <div class="chart">
      <apexchart
        type="bar"
        height="350"
        :options="options"
        :series="options.series"
      ></apexchart>
    </div>
  </n-space>
</template>

<script lang="ts" setup>
import apexchart from "vue3-apexcharts";
import { useCcStore } from "../store/chatting";
import { computed } from "vue";
import { NSpace, NSelect } from "naive-ui";
import { ref } from "vue";

const store = useCcStore();

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

const chatsCountDuration = ref("weekly");
const chatsCountDurations = [
  { label: "Weekly", value: "weekly" },
  { label: "Monthly", value: "monthly" },
  { label: "Yearly", value: "yearly" },
];

function getWeekTuple() {
  const date = new Date();
  const today = date.getDate();
  const currentDay = date.getDay();
  const firstDay = date.setDate(today - (currentDay || 7));
  const lastDay = date.setDate(today - currentDay + 7);

  return [firstDay, lastDay];
}

const options = computed(() => {
  let categories: string[] = [];

  const todayDate = new Date(Date.now());
  const currentMonth = todayDate.getMonth();

  switch (chatsCountDuration.value) {
    case "weekly": {
      const [firstDayTs] = getWeekTuple();

      categories = Array.from({ length: 7 }, (_, i) => {
        const dayNumber = new Date(firstDayTs + 86400000 * i).getDate();

        return `${dayNumber} ${months[currentMonth]}`;
      });

      break;
    }
    case "monthly": {
      const daysInMonth = new Date(
        todayDate.getFullYear(),
        todayDate.getMonth() + 1,
        0
      ).getDate();

      categories = Array.from(
        { length: daysInMonth },
        (_, i) => `${i + 1} ${months[currentMonth]}`
      );
      break;
    }
    case "yearly":
    default: {
      categories = [...months];
      break;
    }
  }

  let series: number[] = [];

  switch (chatsCountDuration.value) {
    case "weekly": {
      const dayCountMap = new Map();
      const [firstDayTs, lastDayTs] = getWeekTuple();

      chats.value.forEach((chat) => {
        const createdDateTs = Number(chat.created);

        if (!(createdDateTs > firstDayTs) || !(createdDateTs < lastDayTs)) {
          return;
        }

        const createdDate = new Date(Number(chat.created));
        const currentDay = createdDate.getDate();
        let count = 0;
        if (dayCountMap.has(currentDay)) {
          count = dayCountMap.get(currentDay);
        }
        count++;
        dayCountMap.set(currentDay, count);
      });

      [...dayCountMap.keys()].forEach((key) => {
        series[key - 1] = dayCountMap.get(key);
      });

      break;
    }
    case "monthly": {
      const currentMonth = todayDate.getMonth();

      const dayCountMap = new Map();
      chats.value.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (createdDate.getMonth() !== currentMonth) {
          return;
        }

        const currentDay = createdDate.getDate();
        let count = 0;
        if (dayCountMap.has(currentDay)) {
          count = dayCountMap.get(currentDay);
        }
        count++;
        dayCountMap.set(currentDay, count);
      });

      [...dayCountMap.keys()].forEach((key) => {
        series[key - 1] = dayCountMap.get(key);
      });

      break;
    }
    case "yearly":
    default: {
      const currentYear = todayDate.getFullYear();
      const monthCountMap = new Map();
      chats.value.forEach((chat) => {
        const createdDate = new Date(Number(chat.created));

        if (createdDate.getFullYear() !== currentYear) {
          return;
        }

        const currentMonth = createdDate.getMonth();
        let count = 0;
        if (monthCountMap.has(currentMonth)) {
          count = monthCountMap.get(currentMonth);
        }
        count++;
        monthCountMap.set(currentMonth, count);
      });

      [...monthCountMap.keys()].forEach((key) => {
        series[key - 1] = monthCountMap.get(key);
      });

      break;
    }
  }

  series = series.map((count) => count || 0);

  const newSeries=[]
  for (let index = 0; index < series.length; index++) {
    newSeries[index]=series[index] || 0;
  }
  series=newSeries;
  console.log(series);

  return {
    series: [
      {
        name: "New tickets",
        data: series,
      },
    ],
    chart: {
      height: 350,
      type: "bar",
    },
    plotOptions: {
      bar: {
        borderRadius: 10,
        dataLabels: {
          position: "top", // top, center, bottom
        },
      },
    },
    dataLabels: {
      enabled: true,
      offsetY: -20,
      style: {
        fontSize: "12px",
        colors: ["#304758"],
      },
    },

    xaxis: {
      categories,
      position: "top",
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
      crosshairs: {
        fill: {
          type: "gradient",
          gradient: {
            colorFrom: "#D8E3F0",
            colorTo: "#BED1E6",
            stops: [0, 100],
            opacityFrom: 0.4,
            opacityTo: 0.5,
          },
        },
      },
      tooltip: {
        enabled: true,
      },
    },
    yaxis: {
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
      labels: {
        show: false,
      },
    },
    theme: {
      mode: "dark",
      palette: "palette1",
    },
    title: {
      text: "New tickets",
      floating: true,
      offsetY: 330,
      align: "center",
      style: {
        color: "#444",
      },
    },
  };
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

console.log(chats.value);
</script>

<style lang="scss">
span.apexcharts-legend-text {
  color: var(--v-primary-base) !important;
}

.chats_count {
  margin: 0px 5vw;
}

.chart {
}

.chart_options {
  display: flex;
  justify-content: end;
}
</style>
