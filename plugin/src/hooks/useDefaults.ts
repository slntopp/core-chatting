import { ref } from 'vue'
import { Defaults, Department, Metric, User, Users } from '../connect/cc/cc_pb.ts'
import { useCcStore } from '../store/chatting.ts'

export interface MetricWithKey extends Metric {
  key: string
}

export default () => {
  const store = useCcStore()

  const isDefaultLoading = ref(false)
  const admins = ref<string[]>([])
  const users = ref<User[]>([])
  const gateways = ref<string[]>([])
  const metrics = ref<MetricWithKey[]>([])
  const departments = ref<Department[]>([])

  async function fetch_defaults() {
    try {
      isDefaultLoading.value = true

      const data = await Promise.all([store.get_members(), store.fetch_defaults()])
      const members = data[0] as Users
      const defaults = data[1] as Defaults

      console.debug("Defaults", defaults)

      admins.value = defaults.admins
      gateways.value = defaults.gateways
      departments.value = defaults.departments
      users.value = members.users
      metrics.value = Object.entries(defaults.metrics).map(
        ([key, value]) => ({ key, ...value })
      ) as MetricWithKey[]

      store.departments = defaults.departments
      store.metrics = metrics.value
    } finally {
      isDefaultLoading.value = false
    }
  }

  return {
    fetch_defaults,
    isDefaultLoading,
    admins,
    gateways,
    users,
    metrics,
    departments
  }
}