<template>
  <div style="margin-top: 20px">
    <n-text>Filter by status:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="checkedStatuses"
      @update:value="(value) => emits('update:checkedStatuses', value)"
    >
      <n-space :wrap-item="false">
        <n-checkbox v-for="status of statuses" :value="status.value" :label="status.label" />
      </n-space>
    </n-checkbox-group>
  </div>

  <div style="margin-top: 20px">
    <n-text>Filter by admin:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="checkedAdmins"
      @update:value="(value) => emits('update:checkedAdmins', value)"
    >
      <n-space :wrap-item="false">
        <n-checkbox v-for="admin of admins" :value="admin.value" :label="admin.label" />
      </n-space>
    </n-checkbox-group>
  </div>

  <div style="margin-top: 20px" v-for="metric of metrics" :key="metric.key">
    <n-text>Filter by {{ metric.title.toLowerCase() }}:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-checkbox-group
      :value="metricsOptions[metric.key]"
      @update:value="(value) => emits('update:metricsOptions', { ...metricsOptions, [metric.key]: value })"
    >
      <n-space :wrap-item="false">
        <n-checkbox v-for="option of metric.options" :value="option.value" :label="option.key" />
      </n-space>
    </n-checkbox-group>
  </div>

  <div style="margin-top: 20px" v-for="item of meta" :key="item.key">
    <n-text>Filter by {{ item.title.toLowerCase() }}:</n-text>
    <n-divider style="margin: 5px 0" />
    <n-date-picker
      v-if="item.type === 'date'"
      type="datetimerange"
      @update:value="(value) =>  changeDateValue('update:metaOptions', value, item.key, metaOptions)"
    />
    <n-checkbox-group
      v-else
      :value="metaOptions[item.key]"
      @update:value="(value) => emits('update:metaOptions', { ...metaOptions, [item.key]: value })"
    >
      <n-space :wrap-item="false">
        <n-checkbox v-for="option of item.options" :value="option.value" :label="option.label" />
      </n-space>
    </n-checkbox-group>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { NSpace, NDivider, NCheckboxGroup, NCheckbox, NText, NDatePicker } from 'naive-ui';
import { Chat, Status } from '../../connect/cc/cc_pb';
import { useCcStore } from '../../store/chatting.ts';
import { MetricWithKey } from '../../hooks/useDefaults';
import { Value } from 'naive-ui/es/date-picker/src/interface';

interface optionsType {
  [index: string]: []
}

interface filtersProps {
  checkedStatuses: Status[]
  checkedAdmins: string[]
  metricsOptions: optionsType
  metaOptions: optionsType,
  metrics: MetricWithKey[]
}

const props = defineProps<filtersProps>()
const emits = defineEmits([
  'update:meta',
  'update:checkedStatuses',
  'update:checkedAdmins',
  'update:metricsOptions',
  'update:metaOptions'
])

const store = useCcStore()

const statuses = computed(() =>
  Object.keys(Status).filter((key) => isFinite(+key)).map((key) => ({
    label: getStatus(+key), value: +key
  }))
)

function getStatus(statusCode: Status | number) {
  const status = Status[statusCode].toLowerCase().replace('_', ' ')

  return `${status[0].toUpperCase()}${status.slice(1)}`
}

interface adminsType {
  value: string
  label: string
}

const admins = computed(() => {
  const result: adminsType[] = []

  store.chats.forEach((chat: Chat) => {
    chat.admins.forEach((uuid) => {
      const element = result.find(({ value }) => uuid === value)

      if (element) return
      else result.push({
        value: uuid, label: store.users.get(uuid)?.title ?? uuid
      })
    })
  })

  return result
})

interface metaOptionType {
  label: string
  value: string | number
}

interface metaType {
  title: string
  key: string
  type?: string
  options?: metaOptionType[]
}

const meta = computed<metaType[]>(() => {
  const result: metaType[] = []

  store.chats.forEach(({ meta }) => {
    if (!meta?.data) return
    Object.entries(meta.data).forEach(([key, rawValue]) => {
      const value = rawValue.toJson()
      const isObject = typeof value === 'object' && !Array.isArray(value)
      const isValid = typeof value === 'string' || typeof value === 'number'
      const isMetric = props.metrics.find((metric) => metric.key === key)
      const isExist = result.find((el) => el.key === key)

      if (['dateStamp'].includes(key) || isMetric) return
      if (isObject && (value?.type === 'date' || value?.type === 'time')) {
        if (isExist) return
        result.push({ title: toKebabCase(key), key, type: value.type })
      } else if (isValid) {
        const i = result.findIndex((el) => el.key === key)
        const option = { label: `${value}`, value }

        if (i === -1 && !isExist) {
          result.push({ title: toKebabCase(key), key, options: [option] })
        }
        else if (!result[i].options?.find((option) => option.value === value)) {
          result[i].options?.push(option)
        }
      }
    })
  })

  result.sort((a, b) => {
    const isATypeDate = a.type === 'date'
    const isBTypeDate = b.type === 'date'

    if (isATypeDate && !isBTypeDate) return 1
    if (!isATypeDate && isBTypeDate) return -1
    return 0
  })
  return result
})

function toKebabCase (text: string) {
  return text.replace(/([a-z])([A-Z])/g, '$1 $2').toLowerCase()
}

function changeDateValue(type: any, value: Value, key: string, options: optionsType) {
  emits(type, {
    ...options,
    [key]: (value as number[]).map((num) => ({ value: num, type: 'date' }))
  })
}

watch(meta, (value) => emits('update:meta', value))
</script>
