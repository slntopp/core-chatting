<template>
  <div
    v-if="!isDefaultLoading || isRefresh"
    class="chat-options"
    style="max-width: 1200px; margin: auto; padding-left: 16px; gap: 20px"
  >
    <n-tabs
      class="card-tabs"
      default-value="settings"
      size="large"
      animated
      pane-wrapper-style="margin: 0 -4px"
      pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
    >
      <n-tab-pane name="settings" tab="Settings">
        <!-- @vue-ignore -->
        <config-view
          @refresh="fetchData"
          :metrics="metrics as MetricWithKey[]"
          :departments="departments as Department[]"
          :admins="admins"
          :gateways="gateways"
          :templates="templates"
          :bot="bot as Bot"
        />
      </n-tab-pane>
      <n-tab-pane name="templates" tab="Templates">
        <n-space style="padding: 6px 24px" align="center">
          <n-h3>Update templates</n-h3>
        </n-space>
        <!-- @vue-ignore -->
        <templates-view
          @refresh="fetchData"
          :metrics="metrics as MetricWithKey[]"
          :departments="departments as Department[]"
          :admins="adminsUuids"
          :gateways="gateways"
          :templates="templates"
          is-edit
        />
      </n-tab-pane>
    </n-tabs>
  </div>

  <n-spin
    v-else
    style="width: 100%; height: 100dvh; margin: auto"
    size="large"
  />
</template>

<script setup lang="ts">
import { NSpin, NTabs, NTabPane, NSpace, NH3 } from "naive-ui";
import templatesView from "../components/settings/templates.vue";
import configView from "../components/settings/config.vue";
import { computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { MetricWithKey, useDefaultsStore } from "../store/defaults";
import { Department } from "../connect/cc/cc_pb";
import { useUsersStore } from "../store/users";

const defaultsStore = useDefaultsStore();
const usersStore = useUsersStore();

const {
  isDefaultLoading,
  admins: adminsUuids,
  departments,
  gateways,
  metrics,
  templates,
  bot,
} = storeToRefs(defaultsStore);
const { users, isUsersLoading } = storeToRefs(usersStore);

const isRefresh = ref(false);

async function fetchData() {
  await defaultsStore.fetch_defaults(true);
  isRefresh.value = true;
}

const admins = computed(() => {
  if (isUsersLoading.value) {
    return [];
  }

  return adminsUuids.value
    .map((uuid) => users.value.get(uuid))
    .filter((v) => !!v);
});

setTimeout(() => fetchData(), 100);
</script>
