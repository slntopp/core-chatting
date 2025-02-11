<template>
  <div class="chat-options" v-if="!isDefaultLoading">
    <n-space v-if="!isEdit" justify="start" align="center">
      <n-button round quaternary @click="cancel">
        <template #icon>
          <n-icon :component="CloseSharp" />
        </template>
      </n-button>

      <n-text> Start a new Chat </n-text>
    </n-space>

    <n-space
      vertical
      justify="start"
      style="max-width: 800px; margin: auto; width: 100%"
    >
      <n-form :model="chat" ref="form" :rules="rules" label-placement="left">
        <n-form-item label="Topic" label-align="left" label-width="100">
          <n-input
            v-model:value="chat.topic"
            clearable
            placeholder="What are we chatting about?"
          />
        </n-form-item>

        <n-form-item label="Members" label-align="left" label-width="100">
          <member-select-pagination v-model:value="chat.users" />
        </n-form-item>

        <n-form-item label="Departament" label-align="left" label-width="100">
          <n-select
            v-model:value="chat.department"
            :options="departamentsOptions"
            filterable
            clearable
          />
        </n-form-item>

        <n-form-item label="Admins" label-align="left" label-width="100">
          <n-select
            v-model:value="chat.admins"
            multiple
            :loading="isUsersLoading"
            :options="adminsWithoutDuplicates"
            filterable
          />
        </n-form-item>

        <n-form-item label="Gateways" label-align="left" label-width="100">
          <n-select
            v-model:value="chat.gateways"
            multiple
            :options="gateways_options"
            filterable
          />
        </n-form-item>

        <n-form-item
          label-align="left"
          label-width="100"
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
          <n-button
            :loading="isEditLoading"
            ghost
            type="success"
            @click="submit"
          >
            {{ isEdit ? "Save" : "Start chat" }}
          </n-button>
        </n-space>
      </n-form>
    </n-space>
  </div>
  <n-spin
    style="width: 100%; height: 100%; margin: auto"
    size="large"
    v-else
  ></n-spin>
</template>

<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  onMounted,
  ref,
  toRefs,
  watch,
} from "vue";
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
} from "naive-ui";

import { useRouter } from "vue-router";
import { useCcStore } from "../../store/chatting.ts";
import { useAppStore } from "../../store/app.ts";
import { Chat, ChatMeta, Role } from "../../connect/cc/cc_pb";
import { ValueAtom } from "naive-ui/es/select/src/interface";
import { Value } from "@bufbuild/protobuf";
import { useDefaultsStore } from "../../store/defaults.ts";
import { storeToRefs } from "pinia";
import { useUsersStore } from "../../store/users.ts";
import memberSelectPagination from "../users/member_select_pagination.vue";

interface ChatOptionsProps {
  minHeight?: string;
  isEdit?: boolean;
  chat?: Chat;
}

const CloseSharp = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseSharp")
);

const props = defineProps<ChatOptionsProps>();
const { isEdit, chat: oldChat } = toRefs(props);

const emit = defineEmits(["close"]);

const router = useRouter();
const store = useCcStore();
const appStore = useAppStore();
const defaultsStore = useDefaultsStore();
const usersStore = useUsersStore();

const { admins, isDefaultLoading, gateways, metrics, departments } =
  storeToRefs(defaultsStore);
const { users, isUsersLoading } = storeToRefs(usersStore);

const form = ref<FormInst>();
const rules = {
  members: {
    required: true,
  },
};

const isEditLoading = ref(false);

const chat = ref<Chat>(
  new Chat({
    uuid: "",
    topic: "",
    users: [store.me.uuid],
    admins: [],
    gateways: [],
    role: Role.OWNER,
  })
);

window.addEventListener("message", ({ data, origin }) => {
  if (origin.includes("localhost")) return;
  if (data.type !== "user-uuid") return;
  chat.value.users.push(data.value);
});

const adminsWithoutDuplicates = computed(() =>
  admins_options.value.map((op) => {
    const disabled = !!chat.value.users.find((m) => m === op.value);
    return { ...op, disabled };
  })
);

const departamentsOptions = computed(() =>
  departments.value.map((departament) => ({
    label: departament.title,
    value: departament.key,
  }))
);

const gateways_options = computed(() =>
  gateways.value.map((gateway) => {
    return {
      label: gateway,
      value: gateway,
    };
  })
);

const admins_options = computed(() =>
  admins.value.map((admin) => {
    return {
      label: users.value.get(admin)?.title ?? "Unknown",
      value: admin,
    };
  })
);

onMounted(() => {
  if (isEdit?.value && oldChat?.value) {
    chat.value = { ...oldChat.value } as Chat;
  }
  window.top?.postMessage({ type: "get-user" }, "*");

  chat.value.gateways = gateways.value;
});

function submit() {
  form.value!.validate(async (errs) => {
    if (errs) {
      console.error("Errors", errs);
      return;
    }

    try {
      isEditLoading.value = true;

      if (isEdit?.value) {
        await store.update_chat(chat.value as Chat);

        emit("close");
      } else {
        let { uuid } = await store.create_chat(chat.value as Chat);

        router.push({ name: "Chat", params: { uuid } });
      }
    } finally {
      isEditLoading.value = false;
    }
  });
}

function getMetric(key: string) {
  return chat.value.meta?.data[key]?.kind.value as ValueAtom;
}

function setMetric(value: number, key: string) {
  if (!chat.value.meta) chat.value.meta = new ChatMeta({});

  chat.value.meta.data[key] = new Value({
    kind: { case: "numberValue", value },
  });
}

function cancel() {
  router.push({ name: "Empty Chat" });
  appStore.displayMode = "half";
}

watch(
  () => chat.value.department,
  () => {
    if (isEdit.value) {
      return;
    }

    if (chat.value.department) {
      const departament = departments.value.find(
        (departament) => departament.key === chat.value.department
      );
      chat.value.admins = departament?.admins || [];
    } else {
      chat.value.admins = [];
    }
  }
);

watch(gateways, () => {
  chat.value.gateways = gateways.value;
});
</script>

<style scoped>
.chat-options {
  display: flex;
  flex-direction: column;
  min-height: v-bind("minHeight");
}
</style>
