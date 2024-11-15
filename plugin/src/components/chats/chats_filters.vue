<template>
  <div class="filter_item dates" style="margin-top: 20px">
    <n-text>Filter by created date:</n-text>

    <n-date-picker
      placeholder="From"
      size="small"
      :value="createDateRange.from"
      @update:value="
        emits('update:createDateRange', { ...createDateRange, from: $event })
      "
      style="margin-right: 10px"
      clearable
      :is-date-disabled="(ts:number)=>ts >= (createDateRange.to || Number.MAX_SAFE_INTEGER)"
    />
    <n-date-picker
      placeholder="To"
      size="small"
      :value="createDateRange.to"
      @update:value="
        emits('update:createDateRange', { ...createDateRange, to: $event })
      "
      clearable
      :is-date-disabled="(ts:number)=>ts <= (createDateRange.from || 0)"
    />
  </div>

  <div class="filter_item dates" style="margin-top: 20px">
    <n-text>Filter by updated date:</n-text>

    <n-date-picker
      placeholder="From"
      size="small"
      :value="updateDateRange.from"
      @update:value="
        emits('update:updateDateRange', { ...updateDateRange, from: $event })
      "
      :is-date-disabled="(ts:number)=>ts >= (updateDateRange.to || Number.MAX_SAFE_INTEGER)"
      style="margin-right: 10px"
      clearable
    />
    <n-date-picker
      placeholder="To"
      size="small"
      :value="updateDateRange.to"
      :is-date-disabled="(ts:number)=>ts <= (updateDateRange.from || 0)"
      @update:value="
        emits('update:updateDateRange', { ...updateDateRange, to: $event })
      "
      clearable
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by department:</n-text>
    <n-select
      filterable
      :options="departments"
      :value="checkedDepartments"
      @update:value="emits('update:checkedDepartments', $event)"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by status:</n-text>
    <n-select
      filterable
      :options="statuses"
      :value="checkedStatuses"
      @update:value="emits('update:checkedStatuses', $event)"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by admins:</n-text>
    <n-select
      filterable
      :options="admins"
      :value="checkedAdmins"
      @update:value="emits('update:checkedAdmins', $event)"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by responsible:</n-text>
    <n-select
      filterable
      :options="admins"
      :value="checkedResponsibles"
      @update:value="emits('update:checkedResponsibles', $event)"
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
      :options="getMetricOptions(metric)"
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
  checkedResponsibles: string[];
  metricsOptions: { [key: string]: [] };
  createDateRange: { from: null | number; to: null | number };
  updateDateRange: { from: null | number; to: null | number };
}
defineProps<ChatsFiltersProps>();

const emits = defineEmits([
  "update:checkedDepartments",
  "update:checkedStatuses",
  "update:checkedAdmins",
  "update:checkedResponsibles",
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

.filter_item.dates {
  display: grid;
  grid-template-columns: 150px 4fr 4fr;

  @media only screen and (max-width: 1024px) {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0px;
  }
}
</style>
