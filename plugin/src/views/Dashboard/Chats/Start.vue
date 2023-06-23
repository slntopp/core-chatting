<template>
    <n-space justify="start" align="center">
        <n-button round quaternary @click="router.push({ name: 'Empty Chat' })">
            <template #icon>
                <n-icon :component="CloseSharp" />
            </template>
        </n-button>

        <n-text>
            Start a new chat
        </n-text>
    </n-space>

    <n-space vertical justify="start" style="padding-left: 16px; padding-top: 10%; min-width: 50vw;">
        <n-space justify="space-around" style="width: 80%;">
            <n-form :model="chat" ref="form" :rules="rules" label-placement="left">
                <n-form-item label="Topic">
                    <n-input v-model:value="chat.topic" clearable placeholder="What are we chatting about?" />
                </n-form-item>

                <n-form-item label="Members">
                    <n-select v-model:value="chat.users" multiple :options="members_options" filterable />
                </n-form-item>

                <n-form-item label="Admins">
                    <n-select v-model:value="chat.admins" multiple :options="admins_options" filterable />
                </n-form-item>

                <n-form-item label="Gateways">
                    <n-select v-model:value="chat.gateways" multiple :options="gateways_options" filterable />
                </n-form-item>
            </n-form>
        </n-space>

        <n-space justify="space-around" style="width: 80%;">
            <n-button ghost type="success" @click="submit">
                Start Chat
            </n-button>
        </n-space>
    </n-space>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent } from 'vue';
import {
    NSpace, NButton, NIcon, NText,
    NForm, NFormItem, NInput,
    SelectOption, NSelect, FormInst
} from 'naive-ui';

import { useRouter } from 'vue-router';

import { Chat, Role } from '../../../connect/cc/cc_pb.ts';
import { useCcStore } from '../../../store/chatting.ts';

const CloseSharp = defineAsyncComponent(() => import('@vicons/ionicons5/CloseSharp'));

const router = useRouter();

const form = ref<FormInst>()
const rules = {
    members: {
        required: true
    }
}

const chat = ref<Chat>(new Chat({
    uuid: '', topic: '', users: [],
    admins: [], gateways: [],
    role: Role.OWNER,
}))

const store = useCcStore();

const members_options = ref<SelectOption[]>()
async function resolve() {
    let result = await store.resolve();
    console.debug("Resolve", result)

    members_options.value = result.users.map(user => {
        return {
            label: user.title,
            value: user.uuid,
        }
    })

}
resolve()

const admins_options = ref<SelectOption[]>()
const gateways_options = ref<SelectOption[]>()
async function fetch_defaults() {
    let result = await store.fetch_defaults();
    console.debug("Defaults", result)

    admins_options.value = result.admins.map(admin => {
        return {
            label: admin.title,
            value: admin.uuid,
        }
    })

    gateways_options.value = result.gateways.map(gateway => {
        return {
            label: gateway,
            value: gateway,
        }
    })
    chat.value.gateways = result.gateways
}
fetch_defaults()

function submit() {
    form.value!.validate(async (errs) => {
        if (errs) {
            console.error("Errors", errs)
            return
        }

        let result = await store.create_chat(chat.value as Chat);

        router.push({ name: 'Chat', params: { uuid: result.uuid } })
    })
}
</script>
