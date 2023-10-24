<template>
  <chat-options min-height="100vh" />
</template>

<script setup lang="ts">

import ChatOptions from "../../../components/chats/chat_options.vue";
</script>

<style>
.n-layout:has(.chat-options) {
  overflow-y: scroll;
}
</style>
