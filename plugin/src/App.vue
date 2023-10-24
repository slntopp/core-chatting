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
import { watch, ref } from 'vue'
import {
  NConfigProvider, NGlobalStyle, NLoadingBarProvider, NNotificationProvider,
  darkTheme, lightTheme
} from 'naive-ui'
import { useAppStore } from './store/app.ts'

const store = useAppStore()
const theme = ref(darkTheme)

watch(() => store.conf?.theme, () => {
  if (store.conf?.theme === 'light') theme.value = lightTheme
  else theme.value = darkTheme
})
</script>
