import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

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

    // Side panel holding the operator's own conversation with the bot. The
    // button that opens it (chat_actions, inside the header) and the panel
    // itself (the Chat view) are two levels apart, so the flag lives here.
    //
    // Persisted: an operator who works out of the copilot lane had to reopen it
    // and re-drag it on every chat switch and every reload.
    const isCopilotOpen = ref(localStorage.getItem('cc.copilot.open') === '1')
    // Share of the width kept by the main chat; the operator drags it.
    const copilotSplit = ref(Number(localStorage.getItem('cc.copilot.split')) || 0.6)

    watch(isCopilotOpen, (v) => localStorage.setItem('cc.copilot.open', v ? '1' : '0'))
    watch(copilotSplit, (v) => localStorage.setItem('cc.copilot.split', String(v)))

    const isMobile = computed(() => device.value === 'phone')
    const isPC = computed(() => device.value === 'pc')
    const isTablet = computed(() => device.value === 'tablet')

    return {
        conf, loading, displayMode, device,

        isPC, isMobile, isTablet,

        isCopilotOpen, copilotSplit,

        isEmergencyMode
    }
})