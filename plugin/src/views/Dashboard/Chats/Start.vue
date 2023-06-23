<template>
    <n-space vertical justify="start">
        <n-space justify="start">
            <n-button round quaternary @click="router.push({ name: 'Empty Chat' })">
                <template #icon>
                    <n-icon :component="CloseSharp" />
                </template>
            </n-button>
        </n-space>

        <n-space justify="start" style="padding-left: 16px">
            <n-form :model="chat" ref="form" label-placement="left">
                <n-form-item label="Topic">
                    <n-input v-model="chat.topic" clearable placeholder="What are we chatting about?" />
                </n-form-item>

                <n-form-item label="Members">
                    <n-select v-model:value="chat.users" multiple :options="options" filterable />
                </n-form-item>
            </n-form>
        </n-space>
    </n-space>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent } from 'vue';
import { 
    NSpace, NButton, NIcon,
    NForm, NFormItem, NInput,
    SelectOption
 } from 'naive-ui';

import { useRouter } from 'vue-router';

import { Chat, Role } from '../../../connect/cc/cc_pb.ts';

const CloseSharp = defineAsyncComponent(() => import('@vicons/ionicons5/CloseSharp'));

const router = useRouter();

const form = ref()
const chat = ref<Chat>({
    role: Role.ROLE_OWNER,
})

const options = ref<SelectOption[]>()
async function get_users() {

}
</script>
