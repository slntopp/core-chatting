import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface PluginConf {
    api: string
    namespace: string
    theme: string
    title: string
    token: string
    vars: {
        [key: string]: string
    }
}

export const useAppStore = defineStore('app', () => {

    const conf = ref<PluginConf>()
    const loading = ref(false)

    return {
        conf, loading
    }
})