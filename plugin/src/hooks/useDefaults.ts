import {Defaults, User, Users} from "../connect/cc/cc_pb.ts";
import {ref} from "vue";
import {useCcStore} from "../store/chatting.ts";

export default () => {

    const store = useCcStore()

    const isDefaultLoading = ref(false)
    const admins = ref<string[]>([])
    const users = ref<User[]>([])
    const gateways = ref<string[]>([])

    async function fetch_defaults() {
        try {
            isDefaultLoading.value = true

            const data = await Promise.all([store.get_members(), store.fetch_defaults()])
            const members = data[0] as Users
            const defaults = data[1] as Defaults

            console.debug("Defaults", defaults)

            admins.value = defaults.admins
            gateways.value = defaults.gateways
            users.value = members.users
        } finally {
            isDefaultLoading.value = false
        }
    }

    return {fetch_defaults, isDefaultLoading, admins, gateways, users}
}