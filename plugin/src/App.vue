<template>
  <n-config-provider :theme="theme">
    <n-loading-bar-provider>
      <n-notification-provider>
        <n-global-style />

        <router-view />

        <div
          v-if="store.isEmergencyMode"
          style="position: fixed; top: 10px; right: 40px; z-index: 10000"
        >
          <n-alert
            closable
            title="There is an emergency situation happening now!!!"
            type="error"
          >
          </n-alert>
        </div>
      </n-notification-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { watch, ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import {
  NConfigProvider,
  NGlobalStyle,
  NLoadingBarProvider,
  NNotificationProvider,
  darkTheme,
  lightTheme,
  NAlert,
} from "naive-ui";
import { useAppStore } from "./store/app.ts";
import { onUnmounted } from "vue";

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

const onMessage = ({ data, origin }: any) => {
  if (origin.includes("localhost")) return;
  if (data.type !== "start-page") return;
  router.push({ path: "/dashboard" });
};

onMounted(() => {
  window.addEventListener("message", onMessage);

  window.addEventListener("resize", setDevice, true);
});

onUnmounted(() => {
  window.removeEventListener("message", onMessage);
  window.removeEventListener("resize", setDevice, true);
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
watch(() => store.conf?.fullscreen, setDevice);
</script>

<style>
:root {
  --main-app-primary-color: #ff00ff;
}
</style>
