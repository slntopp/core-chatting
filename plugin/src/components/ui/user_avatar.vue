<template>
  <n-avatar v-if="avatarPreview" :round="round ?? true" :size="size">
    <transition name="fade">
      <span v-if="typeof avatarPreview === 'string'">{{ avatarPreview }}</span>
      <n-icon :size="iconSize" v-else>
        <component :is="avatarPreview" />
      </n-icon>
    </transition>
  </n-avatar>
</template>

<script setup lang="ts">
import { DefineComponent, computed } from 'vue'
import { NAvatar, NIcon } from 'naive-ui'

interface AvatarProps {
  size?: 'medium' | 'large' | 'small' | number
  iconSize?: number
  avatar: string | DefineComponent
  round: boolean
}

const props = defineProps<AvatarProps>()

const avatarPreview = computed(() => {
  if (typeof props.avatar !== 'string') {
    return props.avatar
  }

  return props.avatar.split(' ').map((s: string) => s?.[0]).join('').toUpperCase()
})
</script>
