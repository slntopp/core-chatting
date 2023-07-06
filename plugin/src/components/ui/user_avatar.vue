<template>
  <n-avatar :round="round" :size="props.size!">{{ avatarPreview }}</n-avatar>
</template>

<script setup lang="ts">
import {NAvatar} from "naive-ui";
import {computed, toRefs} from "vue";

interface AvatarProps {
  size?: 'medium' | 'large' | 'small' | number
  avatar: string
  round: boolean
}

const props = defineProps<AvatarProps>()

// @ts-ignore
const {size='medium',round=true,avatar}=toRefs(props)

const avatarPreview = computed(() => {
  return avatar.value.split(" ").map((s: string) => s?.[0]).join('').toUpperCase()
})
</script>

<style scoped>

</style>