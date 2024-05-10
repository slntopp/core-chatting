import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface PluginConf {
    api: string
    namespace: string
    theme: string
    title: string
    token: string
    params: {
        [key: string]: any
    }
    vars: {
        [key: string]: string
    },
}

export const useAppStore = defineStore('app', () => {

    const conf = ref<PluginConf>()
    const loading = ref(false)
    const displayMode = ref('full')
    const isMobile = ref(false)

    return {
        conf, loading, displayMode, isMobile
    }
})