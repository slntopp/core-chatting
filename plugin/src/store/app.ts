import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface PluginConf {
    api: string
    namespace: string
    theme: string
    title: string
    fullscreen: boolean
    token: string
    params: {
        [key: string]: any
    }
    vars: {
        [key: string]: string
    },
}

export const useAppStore = defineStore('app', () => {
    const isEmergencyMode = ref(false)

    const conf = ref<PluginConf>()
    const loading = ref(false)
    const displayMode = ref('full')

    const device = ref('pc')

    const isMobile = computed(() => device.value === 'phone')
    const isPC = computed(() => device.value === 'pc')
    const isTablet = computed(() => device.value === 'tablet')

    return {
        conf, loading, displayMode, device,

        isPC, isMobile, isTablet,

        isEmergencyMode
    }
})