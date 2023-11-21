<template>
    <n-space justify="center" align="center" style="min-height: 30vh">
        <handler />
    </n-space>
</template>

<script setup lang="ts">
import { h } from 'vue';
import { useRoute, useRouter } from 'vue-router'
import { NAlert, NSpace } from 'naive-ui';

import { useAppStore, PluginConf } from "../store/app.ts";

const route = useRoute()
const router = useRouter()
const store = useAppStore()

function alert(title: string, content: string) {
    return h(NAlert, {
        type: 'error',
        title,
    }, {
        default: () => content
    })
}

function handler() {
    if (!route.query.a) {
        return alert('Falsy Entrypoint', 'Smartcountr is a plugin for infinimesh. It is not meant to be used as a standalone application.')
    }

    let json_data = atob(route.query.a as string)
    if (!json_data) {
        return alert('Invalid Entrypoint', 'The entrypoint is invalid. Please contact your administrator.')
    }
    console.debug('Entrypoint data(string):', json_data)

    let data: PluginConf | undefined

    try {
        data = JSON.parse(json_data)
    } catch (e) {
        console.error(e)
        data = undefined
    }
    if (!data) {
        return alert('Invalid Entrypoint', 'The entrypoint is invalid. Please contact your administrator.')
    }
    console.debug('Entrypoint data(object):', data)

    store.conf = data

    if (store.conf?.params?.redirect) {
        router.push(`/${store.conf.params.redirect}`)
    } else {
        router.push('/dashboard')
    }

    return h(NAlert, {
        title: 'Success', type: 'success', style: { minWidth: '30vw' }
    }, () => `Redirecting to ${store.conf?.params?.redirect || 'dashboard'}...`)
}
</script>
