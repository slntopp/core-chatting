<template>
    <n-alert :type="message.kind ? 'warning' : 'default'" :title="sender()" style="--n-icon-size: 36px">
        <template #icon>
            <n-avatar round :size="24">
                {{ sender()[0] }}
            </n-avatar>
        </template>

        {{ message.content }}

    </n-alert>
</template>

<script setup lang="ts">
import { PropType } from 'vue';
import { NAlert, NAvatar } from 'naive-ui'

import { Message } from '../../connect/cc/cc_pb'
import { useCcStore } from '../../store/chatting';

const { message } = defineProps({
    message: {
        type: Object as PropType<Message>,
        required: true
    }
})

const store = useCcStore()

function sender() {
    return store.users.get(message.sender)?.title ?? 'Unknown'
}
</script>