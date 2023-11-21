<template>
  <div class="chat-options" v-if="!isDefaultLoading">
    <n-space style="padding: 6px 24px" align="center">
      <n-text>Update config</n-text>
    </n-space>

    <n-space
      vertical
      justify="start"
      style="max-width: 800px; margin: auto; padding-left: 16px; gap: 20px"
      :wrap-item="false"
    >
      <div v-for="(option, key) in options" :key="option.key">
        <n-space v-if="option.headers || option.isInput" align="center" :wrap-item="false">
          {{ option.key }}
          <n-button
            ghost
            text
            size="small"
            :type="(inputs[key]?.editMode) ? 'error' : 'success'"
            @click="option.onClick"
          >
            <template #icon>
              <delete-icon v-if="inputs[key]?.editMode" />
              <plus-icon v-else />
            </template>
          </n-button>

          <n-space
            v-if="option.isInput && inputs[key].editMode"
            style="flex-wrap: nowrap"
            :wrap-item="false"
          >
            <n-input v-model:value="(inputs[key].value as string)" size="small" />
            <n-button ghost text type="success" @click="option.onSave">
              <template #icon> <save-icon /> </template>
            </n-button>
          </n-space>
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
            <!-- @vue-ignore -->
            <tr v-for="(item, i) of config[key]">
              <td v-for="header of option.headers">
                <n-space v-if="header.isEditable" style="max-width: 350px" :wrap-item="false">
                  <n-tag
                    round
                    closable
                    type="info"
                    v-for="(value, optionIndex) of item[header.value]"
                    :key="i"
                    @close="header.onClose && header.onClose(i, optionIndex)"
                  >
                    {{ value.key }} ({{ value.value }})
                  </n-tag>

                  <n-space
                    align="center"
                    style="flex-wrap: nowrap"
                    :wrap-item="false"
                  >
                    <template v-if="inputs[key].index === i && header.onSave">
                      <n-input
                        v-model:value="inputs[key].key"
                        size="small"
                        placeholder="Key"
                      />
                      <n-input-number
                        v-model:value="(inputs[key].value as number)"
                        size="small"
                        placeholder="Value"
                      />
                      <n-button ghost text type="success" @click="header.onSave(i)">
                        <template #icon> <save-icon /> </template>
                      </n-button>
                    </template>

                    <n-button
                      ghost
                      text
                      size="small"
                      :type="(inputs[key]?.index === i) ? 'error' : 'success'"
                      @click="header.onClick && header.onClick(i)"
                    >
                      <template #icon>
                        <delete-icon v-if="inputs[key]?.index === i" />
                        <plus-icon v-else />
                      </template>
                    </n-button>
                  </n-space>
                </n-space>

                <n-select
                  tag
                  multiple
                  v-else-if="Array.isArray(item[header.value])"
                  v-model:value="item[header.value]"
                  value-field="uuid"
                  label-field="title"
                  :options="header.options"
                  style="min-width: 200px"
                />

                <n-input v-else clearable v-model:value="item[header.value]" />
              </td>
              <td>
                <n-space justify="center" :wrap-item="false">
                  <n-button
                    ghost
                    text
                    type="error"
                    size="large"
                    @click="option.onClose && option.onClose(i)"
                  >
                    <template #icon> <delete-icon /> </template>
                  </n-button>
                </n-space>
              </td>
            </tr>
          </tbody>
        </n-table>

        <n-select
          v-else-if="option.items"
          tag
          multiple
          v-model:value="(config[key as keyof configType] as Value)"
          value-field="uuid"
          label-field="title"
          :options="option.items"
          style="min-width: 200px"
        />

        <n-space v-else :wrap-item="false">
          <n-tag
            round
            closable
            type="info"
            v-for="(value, i) of config[key as keyof configType]"
            :key="i"
            @close="option.onClose && option.onClose(i)"
          >
            {{ value }}
          </n-tag>
        </n-space>
      </div>

      <n-space justify="end">
        <n-button :loading="isEditLoading" ghost type="success" @click="submit">
          Update
        </n-button>
      </n-space>
    </n-space>
  </div>
  <n-spin v-else style="width: 100%; height: 100vh; margin: auto" size="large" />
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, reactive, ref, watch } from 'vue'
import {
  NButton,
  NInput,
  NInputNumber,
  NTable,
  NSelect,
  NSpace,
  NSpin,
  NText,
  NTag
} from 'naive-ui';
import { Value } from 'naive-ui/es/select/src/interface';
import { Defaults, Department, Metric, Option } from '../connect/cc/cc_pb'
import { useCcStore } from '../store/chatting.ts'
import useDefaults, { MetricWithKey } from '../hooks/useDefaults.ts'

const saveIcon = defineAsyncComponent(() =>
  import('@vicons/ionicons5/SaveOutline')
)
const plusIcon = defineAsyncComponent(() =>
  import('@vicons/ionicons5/Add')
)
const deleteIcon = defineAsyncComponent(() =>
  import('@vicons/ionicons5/CloseOutline')
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
  metrics: MetricWithKey[]
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

interface inputsType {
  [key: string]: {
    value: string | number
    editMode?: boolean
    index?: number
    key?: string
  }
}

const isEditLoading = ref(false)
const inputs = reactive<inputsType>({
  gateways: { editMode: false, value: '' },
  metrics: { index: -1, key: '', value: 0 }
})

const departmentsOptions = computed(() => ({
  key: 'Departments',
  headers: [
    { title: 'Key', value:'key' },
    { title: 'Title', value:'title' },
    { title: 'Description', value:'description' },
    {
      title: 'Admins',
      value:'admins',
      options: config.admins.map((uuid) =>
        users.value.find((user) => user.uuid === uuid)
      )
    }
  ],
  onClick() {
    config.departments.push(new Department({
      key: '', title: '', description: '', admins: []
    }))
  },
  onClose(i: number) {
    config.departments.splice(i, 1)
  }
}))

const adminsOptions = computed(() => ({
  key: 'Admins', items: users.value
}))

const gatewaysOptions = computed(() => ({
  key: 'Gateways',
  isInput: true,
  onClick() {
    inputs.gateways.editMode = !inputs.gateways.editMode
    inputs.gateways.value = ''
  },
  onSave() {
    config.gateways.push(inputs.gateways.value as string)
    inputs.gateways.editMode = false
    inputs.gateways.value = ''
  },
  onClose(i: number) {
    config.gateways.splice(i, 1)
  }
}))

const metricsOptions = computed(() => ({
  key: 'Metrics',
  headers: [
    { title: 'Key', value:'key' },
    { title: 'Title', value:'title' },
    {
      title: 'Options',
      value:'options',
      isEditable: true,

      onClick(metricIndex: number) {
        if (inputs.metrics.index === metricIndex) {
          inputs.metrics = { index: -1, key: '', value: 0 }
        } else {
          inputs.metrics.index = metricIndex
        }
      },
      onSave(i: number) {
        config.metrics[i].options.push(new Option({
          key: inputs.metrics.key,
          value: inputs.metrics.value as number
        }))
        inputs.metrics = { index: -1, key: '', value: 0 }
      },
      onClose(metricIndex: number, i: number) {
        config.metrics[metricIndex].options.splice(i, 1)
      }
    }
  ],
  onClick() {
    const metric = new Metric({ options: [], title: '' }) as MetricWithKey

    metric.key = ''
    config.metrics.push(metric)
  },
  onClose(i: number) {
    config.metrics.splice(i, 1)
  }
}))

interface optionsType {
  [key: string]: {
    key: string
    headers?: [{
      title: string
      value: string
      options?: any[]
      isEditable: boolean

      onClick?: (metricIndex: number) => void
      onSave?: (i: number) => void
      onClose?: (metricIndex: number, i: number) => void 
    }]
    items?: any[]
    isInput?: boolean

    onClick?: () => void
    onSave?: (i?: number) => void
    onClose?: (i: number) => void 
  }
}

/* @ts-ignore */
const options = computed<optionsType>(() => ({
  admins: adminsOptions.value,
  departments: departmentsOptions.value,
  gateways: gatewaysOptions.value,
  metrics: metricsOptions.value
}))

async function submit() {
  try {
    isEditLoading.value = true
    await store.update_defaults(new Defaults({
      ...config, metrics: config.metrics.reduce((result, metric) => ({
        ...result, [metric.key]: metric
      }), {})
    }))
  } finally {
    isEditLoading.value = false
  }
}
</script>
