<template>
  <n-config-provider :theme="theme">
    <n-loading-bar-provider>
      <n-notification-provider>
        <n-global-style />

        <router-view />
      </n-notification-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { watch, ref } from "vue";
import { useRouter } from "vue-router";
import {
  NConfigProvider,
  NGlobalStyle,
  NLoadingBarProvider,
  NNotificationProvider,
  darkTheme,
  lightTheme,
} from "naive-ui";
import { useAppStore } from "./store/app.ts";

const router = useRouter();
const store = useAppStore();

const theme = ref(darkTheme);
const mode = ref("");

watch(
  () => store.conf?.theme,
  () => {
    if (store.conf?.theme === "light") theme.value = lightTheme;
    else theme.value = darkTheme;
  }
);

watch(
  () => store.displayMode,
  (_, value) => {
    if (value === "none") return;
    mode.value = value;
  }
);

window.addEventListener("message", ({ data, origin }) => {
  if (origin.includes("localhost")) return;
  if (data.type !== "start-page") return;
  router.push({ path: "/dashboard" });
});

router.beforeResolve((to, from, next) => {
  if (from.name === "Start Chat") {
    store.displayMode = mode.value;
  }

  if (to.params.uuid) {
    if (store.isMobile) {
      store.displayMode = "none";
    } else if (store.isTablet && store.displayMode === "full") {
      store.displayMode = "none";
    } else if (store.isPC) {
      store.displayMode = "half";
    }
  }

  if (to.path === "/dashboard") {
    store.displayMode = "full";
  }
  next();
});

function setDevice() {
  const screenWidth = document.body.clientWidth;
  if (screenWidth <= 900) {
    store.device = "phone";
  } else if (store.conf?.fullscreen) {
    store.device = "pc";
  } else {
    store.device = "tablet";
  }
}

setDevice();
window.addEventListener("resize", setDevice, true);
watch(() => store.conf?.fullscreen, setDevice);
</script>
