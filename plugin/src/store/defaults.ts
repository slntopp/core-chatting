import { defineStore } from 'pinia'
import { Defaults, Department, FetchDefaultsRequest, Metric, User, Users } from '../connect/cc/cc_pb';
import { ref } from 'vue';
import { useCcStore } from './chatting';


export interface MetricWithKey extends Metric {
    key: string;
}

export interface MessageTemplate {
    name: string;
    content: string;
}

export const useDefaultsStore = defineStore('defaults', () => {
    const cc_store = useCcStore()


    const isDefaultLoading = ref(false);
    const admins = ref<string[]>([]);
    const users = ref<User[]>([]);
    const gateways = ref<string[]>([]);
    const metrics = ref<MetricWithKey[]>([]);
    const departments = ref<Department[]>([]);
    const templates = ref<MessageTemplate[]>([]);


    async function fetch_defaults(fetchTemplates = false) {
        if (isDefaultLoading.value) {
            return;
        }

        try {
            isDefaultLoading.value = true;

            const data = await Promise.all([
                cc_store.get_members(),
                cc_store.users_c.fetchDefaults(new FetchDefaultsRequest({ fetchTemplates }))
            ]);
            const members = data[0] as Users;
            const defaults = data[1] as any as Defaults;

            console.debug("Defaults", defaults);

            admins.value = defaults.admins;
            gateways.value = defaults.gateways;
            departments.value = defaults.departments;
            templates.value = Object.entries(defaults.templates).map(
                ([key, value]) => ({
                    name: key,
                    content: value,
                })
            );
            users.value = members.users;
            metrics.value = Object.entries(defaults.metrics).map(([key, value]) => ({
                key,
                ...value,
            })) as MetricWithKey[];

            cc_store.departments = defaults.departments;
            cc_store.metrics = metrics.value;
        } catch (e) {
            console.log(e);
        } finally {
            isDefaultLoading.value = false;
        }
    }

    function update_defaults(config: Defaults): Promise<Defaults> {
        return cc_store.users_c.setConfig(config);
    }

    fetch_defaults();

    return {

        isDefaultLoading,

        admins,
        users,
        gateways,
        metrics,
        departments,
        templates,

        fetch_defaults,
        update_defaults

    }
})