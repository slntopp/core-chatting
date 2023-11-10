<template>
  <chat-options min-height="100vh" />
</template>

<script setup lang="ts">
import ChatOptions from "../../../components/chats/chat_options.vue";
</script>

<style>
.chat__item .n-layout:has(.chat-options) {
  overflow-y: auto;
}
</style>
