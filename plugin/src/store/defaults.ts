import { defineStore } from 'pinia'
import { Bot, Defaults, Department, FetchDefaultsRequest, Metric, } from '../connect/cc/cc_pb';
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
    const gateways = ref<string[]>([]);
    const metrics = ref<MetricWithKey[]>([]);
    const departments = ref<Department[]>([]);
    const templates = ref<MessageTemplate[]>([]);
    const bot = ref<Bot>();


    async function fetch_defaults(fetchTemplates = false) {
        if (isDefaultLoading.value && !fetchTemplates) {
            return;
        }

        try {
            isDefaultLoading.value = true;

            const data = await cc_store.users_c.fetchDefaults(new FetchDefaultsRequest({ fetchTemplates }))
            const defaults = data as Defaults;

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
            metrics.value = Object.entries(defaults.metrics).map(([key, value]) => ({
                key,
                ...value,
            })) as MetricWithKey[];

            bot.value = defaults.bot;

            cc_store.departments = defaults.departments;
            cc_store.metrics = metrics.value;
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
        gateways,
        metrics,
        departments,
        templates,
        bot,

        fetch_defaults,
        update_defaults

    }
})