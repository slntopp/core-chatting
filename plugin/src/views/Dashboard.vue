<template>
    <n-layout class="dash">
        <n-layout-header v-if="!appStore.isMobile">
            <cc-nav />
        </n-layout-header>
        <n-layout has-sider class="dash__body" :style="appStore.isMobile ? undefined : 'max-height: 100dvh !important;'">
            <router-view />
        </n-layout>
    </n-layout>
</template>

<script setup lang="ts">
import { NLayout, NLayoutHeader, } from 'naive-ui';
import CcNav from '../components/core/nav.vue';

import { useCcStore } from '../store/chatting';
import { useAppStore } from '../store/app';

useCcStore().load_me()
const appStore = useAppStore()

</script>

<style>
@media (max-width: 900px) {
  .dash {
    height: 100%;
    min-height: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .dash__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow: hidden;
  }
}
</style>