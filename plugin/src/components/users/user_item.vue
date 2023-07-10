<template>
  <div class="user__item">
    <div class="user__data">
      <user-avatar class="avatar" round :avatar="props.user.title"/>
      <n-text>{{ props.user.title }}</n-text>
    </div>
    <div class="actions">
      <n-button text @click="emit('delete')">
        <n-icon size="20">
          <trash-icon/>
        </n-icon>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {NButton, NIcon} from 'naive-ui'
import userAvatar from '../ui/user_avatar.vue'
import {User} from "../../connect/cc/cc_pb.ts";
import {Trash as TrashIcon} from '@vicons/ionicons5'

interface UserItemProps {
  user: User
}

const props = defineProps<UserItemProps>()
const emit = defineEmits(['delete'])
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
    margin-left: 5px;
    margin-right: 5px;
  }
}
</style>