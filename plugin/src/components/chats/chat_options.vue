<template>
  <div class="chat-options" v-if="!isDefaultLoading">
    <n-space v-if="!isEdit" justify="start" align="center">
      <n-button round quaternary @click="cancel">
        <template #icon>
          <n-icon :component="CloseSharp"/>
        </template>
      </n-button>

      <n-text>
        Start a new Chat
      </n-text>
    </n-space>

    <n-space vertical justify="start" style="max-width: 800px; margin: auto; padding-left: 16px">
      <n-form :model="chat" ref="form" :rules="rules" label-placement="left">
        <n-form-item label="Topic" label-align="left" label-width="75">
          <n-input v-model:value="chat.topic" clearable placeholder="What are we chatting about?"/>
        </n-form-item>

        <n-form-item label="Members" label-align="left" label-width="75">
          <member-select v-model:value="chat.users" :options="membersWithoutDuplicates"/>
        </n-form-item>

        <n-form-item label="Admins" label-align="left" label-width="75">
          <n-select v-model:value="chat.admins" multiple :options="adminsWithoutDuplicates" filterable/>
        </n-form-item>

        <n-form-item label="Gateways" label-align="left" label-width="75">
          <n-select v-model:value="chat.gateways" multiple :options="gateways_options" filterable/>
        </n-form-item>

        <n-form-item
          label-align="left"
          label-width="75"
          v-for="metric of metrics"
          :label="metric.title"
          :key="metric.key"
        >
          <n-select
            filterable
            label-field="key"
            :value="getMetric(metric.key)"
            :options="metric.options"
            @update:value="(value) => setMetric(value, metric.key)"
          />
        </n-form-item>

        <n-space justify="end">
          <n-button :loading="isEditLoading" ghost type="success" @click="submit">
            {{ isEdit ? 'Edit chat' : 'Start chat' }}
          </n-button>
        </n-space>
      </n-form>
    </n-space>
  </div>
  <n-spin style="width: 100%;height: 100%; margin: auto" size="large" v-else></n-spin>
</template>

<script setup lang="ts">
import {computed, defineAsyncComponent, onMounted, ref, toRefs} from 'vue';
import {
  FormInst,
  NButton,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NSelect,
  NSpace,
  NSpin,
  NText,
  SelectOption,
} from 'naive-ui';

import {useRouter} from 'vue-router';
import {useCcStore} from "../../store/chatting.ts";
import {useAppStore} from '../../store/app.ts';
import {Chat, ChatMeta, Role} from "../../connect/cc/cc_pb";
import useDefaults from "../../hooks/useDefaults.ts";
import MemberSelect from "../users/member_select.vue";
import { ValueAtom } from 'naive-ui/es/select/src/interface';
import { Value } from '@bufbuild/protobuf';

interface ChatOptionsProps {
  minHeight?: string
  isEdit?: boolean
  chat?: Chat
}

const CloseSharp = defineAsyncComponent(() => import('@vicons/ionicons5/CloseSharp'));

const props = defineProps<ChatOptionsProps>()
const {isEdit, chat: oldChat} = toRefs(props)

const emit = defineEmits(['close'])

const router = useRouter()
const store = useCcStore()
const appStore = useAppStore()
const {admins, fetch_defaults, isDefaultLoading, users, gateways, metrics} = useDefaults()

const form = ref<FormInst>()
const rules = {
  members: {
    required: true
  }
}

const isEditLoading = ref(false)

const chat = ref<Chat>(new Chat({
  uuid: '', topic: '', users: [store.me.uuid],
  admins: [], gateways: [],
  role: Role.OWNER,
}))

const admins_options = ref<SelectOption[]>([])
const gateways_options = ref<SelectOption[]>([])
const members_options = ref<SelectOption[]>([])

const adminsWithoutDuplicates = computed(() => admins_options.value.map(op => {
  const disabled = !!chat.value.users.find(m => m === op.value)
  return {...op, disabled}
}))

const membersWithoutDuplicates = computed(() => members_options.value.map(op => {
  const disabled = !!chat.value.admins.find(m => m === op.value)
  return {...op, disabled}
}))

const me = computed(() => store.me)

onMounted(() => {
  if (isEdit?.value && oldChat?.value) {
    chat.value = {...oldChat.value} as Chat
  }
})

onMounted(async () => {
  await fetch_defaults()

  admins_options.value = admins.value.map(admin => {
    return {
      label: store.users.get(admin)?.title ?? 'Unknown',
      value: admin,
    }
  })

  gateways_options.value = gateways.value.map(gateway => {
    return {
      label: gateway,
      value: gateway,
    }
  })
  chat.value.gateways = gateways.value

  members_options.value = users.value.map(user => {
    return {
      label: user.title,
      value: user.uuid,
      disabled: me.value.uuid !== user.uuid && admins.value.includes(user.uuid),
    }
  })
})

function submit() {
  form.value!.validate(async (errs) => {
    if (errs) {
      console.error("Errors", errs)
      return
    }

    try {
      isEditLoading.value = true

      if (isEdit?.value) {
        await store.update_chat(chat.value as Chat)

        emit('close')
      } else {
        let { uuid } = await store.create_chat(chat.value as Chat)

        router.push({name: 'Chat', params: {uuid}})
      }
    } finally {
      isEditLoading.value = false
    }
  })
}

function getMetric(key: string) {
  return chat.value.meta?.data[key].kind.value as ValueAtom
}

function setMetric(value: number, key: string) {
  if (!chat.value.meta) chat.value.meta = new ChatMeta({})

  chat.value.meta.data[key] = new Value({
    kind: { case: 'numberValue', value }
  })
}

function cancel() {
  router.push({ name: 'Empty Chat' })
  appStore.displayMode = 'half'
}
</script>

<style scoped>
.chat-options {
  display: flex;
  flex-direction: column;
  min-height: v-bind('minHeight');
}
</style>