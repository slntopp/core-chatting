<template>
  <div class="chat-options" v-if="!isDefaultLoading">
    <n-space style="padding: 6px 24px" align="center">
      <n-text>Update config</n-text>
    </n-space>

    <n-space
      vertical
      justify="start"
      style="max-width: 800px; margin: auto; padding-left: 16px"
    >
      <template v-for="(option, key) in options" :key="option.key">
        <n-space v-if="option.headers" align="center" :wrap-item="false">
          {{ option.key }}
          <n-button ghost text type="success" size="small" @click="option.onClick">
            <template #icon> <plus-icon /> </template>
          </n-button>
        </n-space>
        <n-text v-else>{{ option.key }}</n-text>

        <n-table v-if="option.headers" :single-line="false">
          <thead>
            <tr>
              <th v-for="header of option.headers">{{ header.title }}</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in config[key]">
              <td v-for="header of option.headers">
                <n-input clearable v-model:value="item[header.value]" />
              </td>
              <td>
                <n-icon> <save-icon /> </n-icon>
              </td>
            </tr>
          </tbody>
        </n-table>
      </template>

      <!-- <n-form-item
        label-align="left"
        label-width="75"
        v-for="metric of metrics"
        :label="metric.title"
        :key="metric.key"
      >
        <n-select
          filterable
          label-field="key"
          :value="getMetric(metric.key)"
          :options="metric.options"
          @update:value="(value) => setMetric(value, metric.key)"
        />
      </n-form-item> -->

      <n-space justify="end">
        <n-button :loading="isEditLoading" ghost type="success" @click="submit">
          Send
        </n-button>
      </n-space>
    </n-space>
  </div>
  <n-spin v-else style="width: 100%; height: 100%; margin: auto" size="large" />
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, reactive, ref, watch } from 'vue'
import {
  FormInst,
  NButton,
  NIcon,
  NInput,
  NTable,
  NSelect,
  NSpace,
  NSpin,
  NText
} from 'naive-ui';
import { Department, Metric, User } from '../connect/cc/cc_pb'
import { useCcStore } from '../store/chatting.ts'
import useDefaults from '../hooks/useDefaults.ts'

const saveIcon = defineAsyncComponent(() =>
  import('@vicons/ionicons5/SaveOutline')
)
const plusIcon = defineAsyncComponent(() =>
  import('@vicons/ionicons5/Add')
)

const store = useCcStore()
const {
  fetch_defaults,
  users,
  metrics,
  admins,
  gateways,
  departments,
  isDefaultLoading
} = useDefaults()

fetch_defaults()

interface configType {
  admins: string[]
  departments: Department[]
  gateways: string[],
  metrics: Metric[]
}

const config = reactive<configType>({
  admins: [],
  departments: [],
  gateways: [],
  metrics: []
})

watch(isDefaultLoading, () => {
  config.metrics = metrics.value
  config.admins = admins.value
  config.gateways = gateways.value
  config.departments = departments.value
})

const form = ref<FormInst>()
const isEditLoading = ref(false)

const departmentsOptions = computed(() => ({
  key: 'Departments',
  headers: [
    { title: 'Key', value:'key' },
    { title: 'Title', value:'title' },
    { title: 'Description', value:'description' },
    { title: 'Admins', value:'admins' }
  ],
  onClick() {
    config.departments.push(new Department({
      key: '', title: '', description: '', admins: []
    }))
  }
}))

const adminsOptions = computed(() => ({
  key: 'Admins',
  items: users.value,
  onClick() {}
}))

const gatewaysOptions = computed(() => ({
  key: 'Gateways',
  headers: [{ title: 'Value', value: 'value' }],
  onClick() {
    config.gateways.push('')
  }
}))

interface optionsType {
  [key: string]: {
    key: string
    headers?: [{
      title: string
      value: string
    }]
    items?: User[]
    onClick: () => void
  }
}

/* @ts-ignore */
const options = computed<optionsType>(() => ({
  departments: departmentsOptions.value,
  admins: adminsOptions.value,
  gateways: gatewaysOptions.value
}))

function submit() {
  form.value!.validate(async (errs) => {
    if (errs) {
      console.error("Errors", errs)
      return
    }

    try {
      isEditLoading.value = true
      await store.update_chat(config)
    } finally {
      isEditLoading.value = false
    }
  })
}

function getMetric(key: string) {
  return key
}

function setMetric(value: number, key: string) {
  console.log(value, key)  
}
</script>
