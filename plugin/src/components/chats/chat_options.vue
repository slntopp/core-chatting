<template>
  <template v-if="!isDefaultLoading">
    <n-space v-if="!isEdit" justify="start" align="center">
      <n-button round quaternary @click="router.push({ name: 'Empty Chat' })">
        <template #icon>
          <n-icon :component="CloseSharp" />
        </template>
      </n-button>

      <n-text>
        Start a new Chat
      </n-text>
    </n-space>

    <n-space vertical justify="start" style="padding-left: 16px; padding-top: 10%;">
      <n-space justify="center">
        <n-form :model="chat" ref="form" :rules="rules" label-placement="left">
          <n-form-item label="Topic">
            <n-input v-model:value="chat.topic" clearable placeholder="What are we chatting about?" />
          </n-form-item>

          <n-form-item label="Members">
            <n-select :render-tag="renderTag" v-model:value="chat.users" multiple :options="membersWithoutDublicates"
              filterable />
          </n-form-item>

          <n-form-item label="Admins">
            <n-select v-model:value="chat.admins" multiple :options="adminsWithoutDublicates" filterable />
          </n-form-item>

          <n-form-item label="Gateways">
            <n-select v-model:value="chat.gateways" multiple :options="gateways_options" filterable />
          </n-form-item>
        </n-form>
      </n-space>

      <n-space justify="space-around" style="width: 80%;">
        <n-button :loading="isEditLoading" ghost type="success" @click="submit">
          {{ isEdit ? 'Edit chat' : 'Start chat' }}
        </n-button>
      </n-space>
    </n-space>
  </template>
  <n-spin style="width: 100%;height: 100%; margin: auto" size="large" v-else></n-spin>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, h, onMounted, ref, toRefs } from 'vue';
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
  NTag,
  NText,
  SelectOption,
  SelectRenderTag
} from 'naive-ui';

import { useRouter } from 'vue-router';
import { useCcStore } from "../../store/chatting.ts";
import { Chat, Role } from "../../connect/cc/cc_pb";

interface ChatOptionsProps {
  isEdit?: boolean
  chat?: Chat
}

const CloseSharp = defineAsyncComponent(() => import('@vicons/ionicons5/CloseSharp'));

const props = defineProps<ChatOptionsProps>()
const { isEdit, chat: oldChat } = toRefs(props)

const emit = defineEmits(['close'])

const router = useRouter();
const store = useCcStore();

const form = ref<FormInst>()
const rules = {
  members: {
    required: true
  }
}

const isDefaultLoading = ref(false)
const isEditLoading = ref(false)

const chat = ref<Chat>(new Chat({
  uuid: '', topic: '', users: [store.me.uuid],
  admins: [], gateways: [],
  role: Role.OWNER,
}))

const admins_options = ref<SelectOption[]>([])
const gateways_options = ref<SelectOption[]>([])
const members_options = ref<SelectOption[]>([])

const adminsWithoutDublicates = computed(() => admins_options.value.filter(op => !chat.value.users.find(m => m === op.value)))
const membersWithoutDublicates = computed(() => members_options.value.filter(op => !chat.value.admins.find(m => m === op.value)))

const me=computed(()=>store.me)

onMounted(() => {
  if (isEdit?.value && oldChat?.value) {
    chat.value = { ...oldChat.value } as Chat
  }
})

async function fetch_defaults() {
  try {
    isDefaultLoading.value = true

    let result = await store.resolve();

    let defaults = await store.fetch_defaults();
    console.debug("Defaults", defaults)

    admins_options.value = defaults.admins.map(admin => {
      return {
        label: store.users.get(admin)?.title ?? 'Unknown',
        value: admin,
      }
    })

    gateways_options.value = defaults.gateways.map(gateway => {
      return {
        label: gateway,
        value: gateway,
      }
    })
    chat.value.gateways = defaults.gateways

    members_options.value = result.users.map(user => {
      return {
        label: user.title,
        value: user.uuid,
        disabled:me.value.uuid!==user.uuid && defaults.admins.includes(user.uuid),
      }
    })
  } finally {
    isDefaultLoading.value = false
  }
}

fetch_defaults()

function submit() {
  form.value!.validate(async (errs) => {
    if (errs) {
      console.error("Errors", errs)
      return
    }

    try {
      isEditLoading.value = true

      if (isEdit?.value) {
        await store.update_chat(chat.value as Chat);

        emit('close')
      } else {
        let result = await store.create_chat(chat.value as Chat);

        router.push({ name: 'Chat', params: { uuid: result.uuid } })
      }
    } finally {
      isEditLoading.value = false
    }


  })
}

const renderTag: SelectRenderTag = ({ option, handleClose }) => {
  return h(
    NTag,
    {
      type: option.type as 'success' | 'warning' | 'error',
      closable: true,
      onMousedown: (e: FocusEvent) => {
        e.preventDefault()
      },
      onClose: (e: MouseEvent) => {
        e.stopPropagation()
        handleClose()
      }
    },
    { default: () => option.value !== option.label ? option.label : 'Deleted' }
  )
}
</script>

<style scoped></style>