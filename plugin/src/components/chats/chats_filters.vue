<template>
  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by department:</n-text>
    <n-select
      filterable
      :items="departments"
      :value="checkedDepartments"
      @update:value="emits('update:checkedDepartments', $event)"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by status:</n-text>
    <n-select
      filterable
      :items="statuses"
      :value="checkedStatuses"
      @update:value="emits('update:checkedStatuses', $event)"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by admins:</n-text>
    <n-select
      filterable
      :items="admins"
      :value="checkedAdmins"
      @update:value="emits('update:checkedAdmins', $event)"
      multiple
    />
  </div>

  <div
    class="filter_item"
    v-for="metric of store.metrics"
    style="margin-top: 20px"
    :key="metric.key"
  >
    <n-text>Filter by {{ metric.title.toLowerCase() }}:</n-text>

    <n-select
      filterable
      :items="getMetricOptions(metric)"
      :value="metricsOptions[metric.key]"
      @update:value="emits('update:checkedMetrics', metric.key, $event)"
      multiple
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { NText, NSelect } from "naive-ui";

import { useCcStore } from "../../store/chatting";
import { Chat, Status } from "../../connect/cc/cc_pb";
import { SelectMixedOption } from "naive-ui/es/select/src/interface";
import { MetricWithKey } from "../../hooks/useDefaults";

interface ChatsFiltersProps {
  checkedDepartments: string[];
  checkedStatuses: Status[];
  checkedAdmins: string[];
  metricsOptions: { [key: string]: [] };
}
defineProps<ChatsFiltersProps>();

const emits = defineEmits([
  "update:checkedDepartments",
  "update:checkedStatuses",
  "update:checkedAdmins",
  "update:checkedMetrics",
]);
const store = useCcStore();

const departments = computed(() =>
  store.departments.map(({ key, title }) => ({ label: title, value: key }))
);

const statuses = computed(() =>
  Object.keys(Status)
    .filter((key) => isFinite(+key))
    .map((key) => ({
      label: getStatus(+key),
      value: +key,
    }))
);

const admins = computed(() => {
  const result: SelectMixedOption[] = [];

  (store.chats as Map<string, Chat>).forEach((chat: Chat) => {
    chat.admins.forEach((uuid) => {
      const element = result.find(({ value }) => uuid === value);

      if (element) return;
      else
        result.push({
          value: uuid,
          label: store.users.get(uuid)?.title ?? uuid,
        });
    });
  });

  return result;
});

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace("_", " ");

  return `${status[0].toUpperCase()}${status.slice(1)}`;
}

function getMetricOptions(metric: MetricWithKey) {
  return metric.options.map((option) => ({
    label: option.key,
    value: option.value,
  }));
}
</script>

<style>
.filter_item {
  display: grid;
  grid-template-columns: 150px 8fr;
}

.filter_item span {
  display: flex;
  justify-content: start;
  align-items: center;
}
</style>
