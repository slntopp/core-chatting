<template>
  <n-config-provider :theme="darkTheme" >
    <n-loading-bar-provider>
      <n-notification-provider>

        <n-global-style />

        <router-view />

      </n-notification-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">

import {
  NConfigProvider, NGlobalStyle, NLoadingBarProvider, NNotificationProvider,
  darkTheme
} from "naive-ui"

</script>