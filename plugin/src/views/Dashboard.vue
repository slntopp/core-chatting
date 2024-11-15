<template>
    <n-layout>
        <n-layout-header>
            <cc-nav />
        </n-layout-header>
        <n-layout has-sider  style="max-height: 100dvh !important;">
            <router-view />
        </n-layout>
    </n-layout>
</template>

<script setup lang="ts">
import { NLayout, NLayoutHeader, } from 'naive-ui';
import CcNav from '../components/core/nav.vue';

import { useCcStore } from '../store/chatting';

useCcStore().load_me()

</script>