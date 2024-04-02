<template>
  <div class="user__item">
    <div class="user__data">
      <user-avatar class="avatar" round :avatar="user.title"/>
      <n-text style="color: inherit">{{ user.title }}</n-text>
    </div>

    <div class="actions" v-if="actions">
      <n-popover v-if="action === 'change'" trigger="click" placement="right">
        <template #trigger>
          <n-button text style="color: inherit" @click="emits(action)">
            <n-icon size="20"> <sync-icon /> </n-icon>
          </n-button>
        </template>
        <slot name="popover-content">Do you want to do this?</slot>
      </n-popover>

      <n-button v-else text style="color: inherit" @click="emits(action ?? 'delete')">
        <n-icon size="20">
          <trash-icon />
        </n-icon>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue'
import { NButton, NIcon, NPopover, NText } from 'naive-ui'
import { User } from "../../connect/cc/cc_pb.ts"
import userAvatar from '../ui/user_avatar.vue'

const trashIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/Trash')
)
const syncIcon = defineAsyncComponent(
  () => import('@vicons/ionicons5/Sync')
)

interface UserItemProps {
  user: User
  actions?: boolean,
  action?: 'delete' | 'change'
}

defineProps<UserItemProps>()
const emits = defineEmits(['delete', 'change'])
</script>

<style scoped lang="scss">
.user__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;

  .user__data {
    display: flex;
    justify-content: start;
    align-items: center;

    .avatar {
      margin-right: 10px;
    }
  }

  .actions {
    display: flex;
    margin-left: 5px;
    margin-right: 5px;
  }
}
</style>