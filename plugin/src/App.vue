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
  console.log(to, from, next);

  if (from.name === "Start Chat") {
    store.displayMode = mode.value;
  }

  if (to.params.uuid) {
    store.displayMode = store.isMobile ? "none" : "half";
  }

  if (to.path === "/dashboard") {
    store.displayMode = "full";
  }
  next();
});

store.isMobile = screen.width <= 760;
</script>
