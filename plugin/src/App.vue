<template>
  <n-config-provider :theme="theme" >
    <n-loading-bar-provider>
      <n-notification-provider>

        <n-global-style />

        <router-view />

      </n-notification-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NConfigProvider, NGlobalStyle, NLoadingBarProvider, NNotificationProvider,
  darkTheme, lightTheme
} from "naive-ui"

const theme = ref(darkTheme)

onMounted(() => {
  window.top?.postMessage({ type: 'get-theme' }, '*')
  window.addEventListener('message', ({ data, origin }) => {
    if (origin.includes(location.host)) return
    if (data.theme === 'light') theme.value = lightTheme
    else theme.value = darkTheme
  })
})
</script>