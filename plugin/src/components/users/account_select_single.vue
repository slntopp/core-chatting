<template>
  <n-select
    :value="modelValue"
    @update-value="(v: string) => emit('update:modelValue', v)"
    :options="options"
    filterable
    :on-search="startSearch"
    :on-update:show="onUpdateShow"
    :loading="isLoading"
    placeholder="Start typing account name..."
  >
    <template #empty>
      <div>
        <n-spin v-if="isLoading"></n-spin>
        <span v-else>No data</span>
      </div>
    </template>
  </n-select>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { GetMembersRequest, Users } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { NSelect, NSpin, SelectOption } from "naive-ui";
import { debounce } from "../../functions";

interface Props {
  modelValue: string;
}

const props = defineProps<Props>();
const emit = defineEmits(["update:modelValue"]);

const ccStore = useCcStore();

const options = ref<SelectOption[]>([]);
const isLoading = ref(false);

onMounted(() => {
  if (props.modelValue) fetchExisting(props.modelValue);
});

function startSearch(val: string) {
  isLoading.value = true;
  onSearchDebounced(val);
}

const onSearchDebounced = debounce(onSearch, 600);

async function onSearch(value?: string) {
  isLoading.value = true;
  try {
    const data = (
      await ccStore.users_c.getMembers(
        GetMembersRequest.fromJson({
          limit: 20,
          page: 1,
          filters: { search_param: (value || "").trim() },
        })
      )
    ).toJson() as any as Users;
    options.value = (data.users || []).map((user) => ({
      label: user.title,
      value: user.uuid?.toString(),
    }));
  } finally {
    isLoading.value = false;
  }
}

async function fetchExisting(uuid: string) {
  const data = (
    await ccStore.users_c.getMembers(
      GetMembersRequest.fromJson({ uuids: [uuid] })
    )
  ).toJson() as any as Users;
  options.value = (data.users || []).map((user) => ({
    label: user.title,
    value: user.uuid?.toString(),
  }));
}

function onUpdateShow(show: boolean) {
  if (show && options.value.length === 0) onSearch("");
}
</script>

<style scoped></style>
