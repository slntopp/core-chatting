<template>
  <n-avatar :round="round" :size="size">{{ avatarPreview }}</n-avatar>
</template>

<script setup lang="ts">
import {NAvatar} from "naive-ui";
import {computed} from "vue";

interface AvatarProps{
  size:'medium' | 'large' | 'small'
  avatar:string
  round:boolean
}

const props = defineProps<AvatarProps>()

const avatarPreview=computed(()=>{
  return props.avatar?.split(" ").map((s:string)=>s?.[0]).join('').toUpperCase()
})
</script>

<style scoped>

</style>