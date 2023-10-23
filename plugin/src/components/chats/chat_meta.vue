<template>
  <n-space vertical :wrap-item="false">
    <template v-for="option in metaOptions">
      <n-text v-if="option.type === 'date' && typeof option.value === 'number'">
        {{ option.title }}:
        <n-tag round :type="(option.value <= Date.now()) ? 'error' : 'success'">
          {{ new Date(option.value).toLocaleString() }}
        </n-tag>
      </n-text>

      <template v-else-if="option.headers">
        <n-text>{{ option.title }}:</n-text>
        <n-data-table :columns="option.headers" :data="(option.value as RowData[])" />
      </template>

      <n-text v-else-if="option.type === 'dateRange' && Array.isArray(option.value)">
        {{ option.title }}:
        <n-tag round type="info">
          {{ (option.value[0] as Date).toLocaleString() }}
        </n-tag>
        -
        <n-tag round type="info">
          {{ (option.value[1] as Date).toLocaleString() }}
        </n-tag>
      </n-text>

      <n-text v-else-if="option.type === 'time'">
        {{ option.title }}:
        <n-tag round>{{ new Date(option.value as number).toLocaleTimeString() }}</n-tag>
      </n-text>

      <n-text v-else>
        {{ option.title }}:
        <n-tag round size="small">{{ option.value }}</n-tag>
      </n-text>
    </template>
  </n-space>
</template>

<script setup lang="ts">
import { VNode, computed, h } from 'vue'
import { NDataTable, NSpace, NSwitch, NTag, NText } from 'naive-ui'
import { RowData } from 'naive-ui/es/data-table/src/interface';
import { JsonObject, JsonValue } from '@bufbuild/protobuf';
import { useCcStore } from '../../store/chatting.ts'
import { Chat } from '../../connect/cc/cc_pb'
import useDefaults from '../../hooks/useDefaults.ts';

const props = defineProps({
  chat: Chat
})

const store = useCcStore()
const { metrics, fetch_defaults } = useDefaults()

interface OptionType {
  type?: string
  title: string
  value: DeviceType[] | Date[] | JsonValue | undefined
  headers?: {
    title: string,
    key: string,
    render?: (row: DeviceType) => VNode
  }[]
}

interface DeviceType {
  uuid: string
  title: string
  enabled: boolean
}

const metaOptions = computed(() => {
  const result: OptionType[] = []

  Object.entries(props.chat?.meta?.data ?? {}).forEach(([key, value]) => {
    if (metrics.value.find((metric) => metric.key === key)) return

    const title = key.replace(/([a-z])([A-Z])/g, "$1 $2")
    const option: OptionType = {
      title: `${title[0].toUpperCase()}${title.slice(1)}`,
      value: value.toJson()
    }
    const isDate = Array.isArray(option.value) && option.value.every(
      (value) => typeof value === 'number'
    )

    if (key === 'devices') {
      const devices = (option.value as []).map(
        (uuid: string) => store.devices.get(uuid)
      )

      option.headers = [
        { title: 'Uuid', key: 'uuid', render(row) {
          return h(NTag, { round: true }, () => row.uuid)
        }, },
        { title: 'Title', key: 'title' },
        { title: 'Enabled', key: 'enabled', render(row) {
          return h(NSwitch, { value: row.enabled, disabled: true })
        } }
      ]

      option.value = Array.from(devices.values()).map((device) => ({
        uuid: device?.uuid ?? '',
        title: device?.title ?? '',
        enabled: device?.enabled ?? false
      }))
    } else if (isDate) {
      option.type = 'dateRange'
      option.value = (option.value as []).map(
        (ts: number) => new Date(ts * 1000)
      )
    } else if ((option.value as JsonObject).type === 'date') {
      option.value = ((option.value as JsonObject).value) as number * 1000
      option.type = 'date'
    } else if ((option.value as JsonObject).type === 'time') {
      const date = new Date(new Date().toISOString().split('T')[0])
      
      option.value = date.getTime() + (date.getTimezoneOffset() * 60 * 1000) + ((option.value as JsonObject)?.value as number) * 1000
      option.type = 'time'
    }

    result.push(option)
  })

  result.sort(({ headers }) => (headers) ? 1 : -1)
  return result
})

if (store.devices.size < 1) store.fetch_devices()
if (metrics.value.length < 1) fetch_defaults()
</script>
