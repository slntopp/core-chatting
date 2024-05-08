<template>
  <div style="margin-top: 20px">
    <n-text>Filter by department:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="checkedDepartments"
      @update:value="emits('update:checkedDepartments', $event)"
    >
      <n-space :wrap-item="false">
        <n-checkbox
          v-for="dep of departments"
          :value="dep.value"
          :label="dep.label"
        />
      </n-space>
    </n-checkbox-group>
  </div>

  <div style="margin-top: 20px">
    <n-text>Filter by status:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="checkedStatuses"
      @update:value="emits('update:checkedStatuses', $event)"
    >
      <n-space :wrap-item="false">
        <n-checkbox
          v-for="status of statuses"
          :value="status.value"
          :label="status.label"
        />
      </n-space>
    </n-checkbox-group>
  </div>

  <div style="margin-top: 20px">
    <n-text>Filter by admin:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="checkedAdmins"
      @update:value="emits('update:checkedAdmins', $event)"
    >
      <n-space :wrap-item="false">
        <n-checkbox
          v-for="admin of admins"
          :value="admin.value"
          :label="admin.label"
        />
      </n-space>
    </n-checkbox-group>
  </div>

  <div
    style="margin-top: 20px"
    v-for="metric of store.metrics"
    :key="metric.key"
  >
    <n-text>Filter by {{ metric.title.toLowerCase() }}:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="metricsOptions[metric.key]"
      @update:value="emits('update:checkedMetrics', metric.key, $event)"
    >
      <n-space :wrap-item="false">
        <n-checkbox
          v-for="option of metric.options"
          :value="option.value"
          :label="option.key"
        />
      </n-space>
    </n-checkbox-group>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { NSpace, NDivider, NCheckboxGroup, NCheckbox, NText } from "naive-ui";

import { useCcStore } from "../../store/chatting";
import { Chat, Status } from "../../connect/cc/cc_pb";

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

interface adminsType {
  value: string;
  label: string;
}

const admins = computed(() => {
  const result: adminsType[] = [];

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
</script>
