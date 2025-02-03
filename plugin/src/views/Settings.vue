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
          :users="users as User[]"
          :metrics="metrics as MetricWithKey[]"
          :departments="departments as Department[]"
          :admins="admins"
          :gateways="gateways"
          :templates="templates"
        />
      </n-tab-pane>
      <n-tab-pane name="templates" tab="Templates">
        <n-space style="padding: 6px 24px" align="center">
          <n-h3>Update templates</n-h3>
        </n-space>
        <!-- @vue-ignore -->
        <templates-view
          @refresh="fetchData"
          :users="users as User[]"
          :metrics="metrics as MetricWithKey[]"
          :departments="departments as Department[]"
          :admins="admins"
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
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { MetricWithKey, useDefaultsStore } from "../store/defaults";
import { Department } from "../connect/cc/cc_pb";

const defaultsStore = useDefaultsStore();
const {
  isDefaultLoading,
  admins,
  departments,
  gateways,
  metrics,
  users,
  templates,
} = storeToRefs(defaultsStore);

const isRefresh = ref(false);

async function fetchData() {
  await defaultsStore.fetch_defaults(true);
  isRefresh.value = true;
}

fetchData();
</script>
