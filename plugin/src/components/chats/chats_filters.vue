<template>
  <div>
    <n-space align="center" style="gap: 5px" :wrap-item="false">
      <n-text>Sort By:</n-text>
      <n-icon
        size="18"
        style="cursor: pointer"
        :style="{
          transform: localFilter.sortType === 'desc' ? 'rotate(180deg)' : null,
        }"
        @click="
          localFilter.sortType === 'desc'
            ? (localFilter.sortType = 'asc')
            : (localFilter.sortType = 'desc')
        "
      >
        <up-icon />
      </n-icon>
    </n-space>

    <n-divider style="margin: 5px 0" />
    <n-radio-group v-model:value="localFilter.sortBy">
      <n-space :wrap-item="false">
        <n-radio value="unread" label="None" />
        <n-radio value="created" label="Created" />
        <n-radio value="updated " label="Sent" />
        <n-radio value="status" label="Status" />
      </n-space>
    </n-radio-group>
  </div>

  <div class="filter_item dates" style="margin-top: 20px">
    <n-text>Filter by created date:</n-text>

    <n-date-picker
      placeholder="From"
      size="small"
      v-model="localFilter.createDateRange.from"
      style="margin-right: 10px"
      clearable
      :is-date-disabled="(ts:number)=>ts >= (localFilter.createDateRange.to || Number.MAX_SAFE_INTEGER)"
    />
    <n-date-picker
      placeholder="To"
      size="small"
      v-model:value="localFilter.createDateRange.to"
      clearable
      :is-date-disabled="(ts:number)=>ts <= (localFilter.createDateRange.from || 0)"
    />
  </div>

  <div class="filter_item dates" style="margin-top: 20px">
    <n-text>Filter by updated date:</n-text>

    <n-date-picker
      placeholder="From"
      size="small"
      v-model:value="localFilter.updateDateRange.from"
      :is-date-disabled="(ts:number)=>ts >= (localFilter.updateDateRange.to || Number.MAX_SAFE_INTEGER)"
      style="margin-right: 10px"
      clearable
    />
    <n-date-picker
      placeholder="To"
      size="small"
      v-model:value="localFilter.updateDateRange.to"
      :is-date-disabled="(ts:number)=>ts <= (localFilter.updateDateRange.from || 0)"
      clearable
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by department:</n-text>
    <n-select
      filterable
      :options="departments"
      v-model:value="localFilter.checkedDepartments"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by status:</n-text>
    <n-select
      filterable
      :options="statuses"
      v-model:value="localFilter.checkedStatuses"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by admins:</n-text>
    <n-select
      filterable
      :options="admins"
      v-model:value="localFilter.checkedAdmins"
      multiple
    />
  </div>

  <div class="filter_item" style="margin-top: 20px">
    <n-text>Filter by responsible:</n-text>
    <n-select
      filterable
      :options="admins"
      v-model:value="localFilter.checkedResponsibles"
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
      v-model:value="localFilter.metricsOptions[metric.key]"
      multiple
    />
  </div>

  <div style="margin: 15px 10px 0px 0px; display: flex; justify-content: end">
    <!-- <n-button
      ghost
      type="warning"
      @click="downloadReport"
      :loading="isReportLoading"
    >
      Report
    </n-button> -->

    <div>
      <n-button
        style="margin-right: 10px"
        ghost
        type="error"
        @click="emits('reset')"
      >
        Reset
      </n-button>

      <n-button ghost type="primary" :disabled="!isEdited" @click="onSearch">
        Search
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref, toRefs, watch } from "vue";
import {
  NText,
  NSelect,
  NDatePicker,
  NButton,
  NRadioGroup,
  NRadio,
  NSpace,
  NIcon,
  NDivider,
} from "naive-ui";

import { useCcStore } from "../../store/chatting";
import { Chat, Status } from "../../connect/cc/cc_pb";
import { SelectMixedOption } from "naive-ui/es/select/src/interface";
import { MetricWithKey } from "../../hooks/useDefaults";
import { getStatusItems } from "../../functions";

const upIcon = defineAsyncComponent(() => import("@vicons/ionicons5/ArrowUp"));

interface ChatsFiltersProps {
  checkedDepartments: string[];
  checkedStatuses: Status[];
  checkedAdmins: string[];
  checkedResponsibles: string[];
  metricsOptions: { [key: string]: [] };
  createDateRange: { from: null | number; to: null | number };
  updateDateRange: { from: null | number; to: null | number };
  sortBy: string;
  sortType: string;
}
const props = defineProps<ChatsFiltersProps>();

const emits = defineEmits([
  "update:checkedDepartments",
  "update:checkedStatuses",
  "update:checkedAdmins",
  "update:checkedResponsibles",
  "update:checkedMetrics",
  "update:createDateRange",
  "update:updateDateRange",
  "update:sortBy",
  "update:sortType",
  "close",
  "reset",
]);
const store = useCcStore();

const {
  checkedAdmins,
  checkedDepartments,
  checkedResponsibles,
  checkedStatuses,
  createDateRange,
  metricsOptions,
  sortBy,
  sortType,
  updateDateRange,
} = toRefs(props);

const localFilter = ref<ChatsFiltersProps>({} as ChatsFiltersProps);
const startFilter = ref<ChatsFiltersProps>({} as ChatsFiltersProps);

const departments = computed(() =>
  store.departments.map(({ key, title }) => ({ label: title, value: key }))
);

const metrics = computed<MetricWithKey[]>(
  () => store.metrics as MetricWithKey[]
);

const statuses = computed(() => getStatusItems());
const isEdited = computed(
  () => JSON.stringify(startFilter.value) !== JSON.stringify(localFilter.value)
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

function getMetricOptions(metric: MetricWithKey) {
  return metric.options.map((option) => ({
    label: option.key,
    value: option.value,
  }));
}

function onSearch() {
  emits("update:checkedAdmins", localFilter.value.checkedAdmins);
  emits("update:checkedDepartments", localFilter.value.checkedDepartments);
  emits("update:checkedResponsibles", localFilter.value.checkedResponsibles);
  emits("update:checkedStatuses", localFilter.value.checkedStatuses);
  emits("update:createDateRange", localFilter.value.createDateRange);
  emits("update:updateDateRange", localFilter.value.updateDateRange);
  emits("update:checkedMetrics", localFilter.value.metricsOptions);
  emits("update:sortBy", localFilter.value.sortBy);
  emits("update:sortType", localFilter.value.sortType);
  emits("close");
}

function setLocalFilters() {
  const data = {
    checkedAdmins: checkedAdmins.value,
    checkedDepartments: checkedDepartments.value,
    checkedResponsibles: checkedResponsibles.value,
    checkedStatuses: checkedStatuses.value,
    createDateRange: createDateRange.value,
    updateDateRange: updateDateRange.value,
    sortBy: sortBy.value,
    sortType: sortType.value,
    metricsOptions: metricsOptions.value,
  };
  localFilter.value = JSON.parse(JSON.stringify(data));
  startFilter.value = JSON.parse(JSON.stringify(data));
}
setLocalFilters();

watch(
  [
    checkedAdmins,
    checkedDepartments,
    checkedResponsibles,
    checkedStatuses,
    createDateRange,
    metricsOptions,
    sortBy,
    sortType,
    updateDateRange,
  ],
  () => {
    setLocalFilters();
  },
  { deep: true }
);
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
