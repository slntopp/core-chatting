<template>
  <n-button :bordered="false" @click="emit('click')"><slot></slot>
    <n-icon style="margin-left: 5px" size="20">
      <add-icon></add-icon>
    </n-icon>
  </n-button>
</template>

<script setup lang="ts">
import {NButton,NIcon} from "naive-ui";
import {AddCircle as AddIcon} from '@vicons/ionicons5'

const emit=defineEmits(['click'])
</script>

<style scoped>

</style>