<template>
  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by created date:</n-text>

    <n-date-picker
      size="small"
      :value="createDateRange"
      @update:value="emits('update:createDateRange', $event)"
      type="daterange"
      clearable
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by updated date:</n-text>

    <n-date-picker
      size="small"
      :value="updateDateRange"
      @update:value="emits('update:updateDateRange', $event)"
      type="daterange"
      clearable
    />
  </div>

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
    v-for="metric of metrics"
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
import { NText, NSelect, NDatePicker } from "naive-ui";

import { useCcStore } from "../../store/chatting";
import { Chat, Status } from "../../connect/cc/cc_pb";
import { SelectMixedOption } from "naive-ui/es/select/src/interface";
import { MetricWithKey } from "../../hooks/useDefaults";
import { getStatusItems } from "../../functions";

interface ChatsFiltersProps {
  checkedDepartments: string[];
  checkedStatuses: Status[];
  checkedAdmins: string[];
  metricsOptions: { [key: string]: [] };
  createDateRange: [number, number] | null;
  updateDateRange: [number, number] | null;
}
defineProps<ChatsFiltersProps>();

const emits = defineEmits([
  "update:checkedDepartments",
  "update:checkedStatuses",
  "update:checkedAdmins",
  "update:checkedMetrics",
  "update:createDateRange",
  "update:updateDateRange",
]);
const store = useCcStore();

const departments = computed(() =>
  store.departments.map(({ key, title }) => ({ label: title, value: key }))
);

const metrics = computed<MetricWithKey[]>(
  () => store.metrics as MetricWithKey[]
);

const statuses = computed(() => getStatusItems());

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

  @media only screen and (max-width: 1024px) {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0px;
  }
}

.filter_item span {
  display: flex;
  justify-content: start;
  align-items: center;
}
</style>
