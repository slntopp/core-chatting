<template>
  <n-select
    :consistent-menu-width="false"
    :render-tag="renderTag"
    :value="value"
    @update-value="onUpdate"
    multiple
    :options="options"
    filterable
    :on-search="startSearch"
    :on-update:show="onUpdateShow"
    :loading="isMembersLoading"
    placeholder="Start typing..."
    :filter="filterOption"
  >
    <template #empty>
      <div>
        <n-spin v-if="isMembersLoading"></n-spin>
        <span v-else>No data</span>
      </div>
    </template>
  </n-select>
</template>

<script setup lang="ts">
import { onMounted, ref, toRefs, watch, h } from "vue";
import { GetMembersRequest, User, Users } from "../../connect/cc/cc_pb";
import { useCcStore } from "../../store/chatting";
import { NSelect, NTag, NSpin, SelectOption, SelectRenderTag } from "naive-ui";
import { debounce } from "../../functions";
import { useUsersStore } from "../../store/users";

interface MemberSelectPaginationProps {
  value: string[];
}

const props = defineProps<MemberSelectPaginationProps>();
const { value } = toRefs(props);
const emit = defineEmits(["update:value", "on-search"]);

const ccStore = useCcStore();
const usersStore = useUsersStore();

const options = ref<SelectOption[]>([]);
const isMembersLoading = ref(false);
const cachedEmpty = ref<User[]>([]);

onMounted(() => fetchExisted());

function startSearch(val: string) {
  isMembersLoading.value = true;
  onSearchDebounced(val);
}

function onUpdate(value: any) {
  emit("update:value", value);

  onSearch();
}

const onSearchDebounced = debounce(onSearch, 600);

async function onSearch(value?: string) {
  isMembersLoading.value = true;

  if (!value) {
    value = "";
  }

  value = value.trim();

  if (value === "" && cachedEmpty.value.length !== 0) {
    isMembersLoading.value = false;
    return setOptions(cachedEmpty.value as User[]);
  }

  try {
    const data = (
      await ccStore.users_c.getMembers(
        GetMembersRequest.fromJson({
          limit: 20,
          page: 1,
          filters: { search_param: value, exclude_uuids: [...props.value] },
        })
      )
    ).toJson() as any as Users;

    if (value === "") {
      cachedEmpty.value = data.users;
    }

    setOptions(data.users);
  } finally {
    isMembersLoading.value = false;
  }
}

function onUpdateShow() {
  if (options.value.length === 0 || cachedEmpty.value.length === 0) {
    onSearch("");
  }
}

function setOptions(users?: User[]) {
  const newOptions = options.value.filter((option) =>
    props.value.includes(option?.value?.toString() || "")
  );

  newOptions.push(
    ...(users || []).map((user) => ({
      label: user.title,
      value: user.uuid,
    }))
  );

  newOptions.sort((a) =>
    props.value.includes(a.value?.toString() || "") ? 1 : -1
  );

  options.value = newOptions.filter(
    (option, index) =>
      newOptions.findIndex((first) => first.value == option.value) === index
  );
}

async function fetchExisted() {
  const notExisted = [];
  const inCache = [];

  for (const uuid of props.value) {
    if (options.value.findIndex((option) => option.value === uuid) == -1) {
      if (usersStore.users.get(uuid)) {
        inCache.push(usersStore.users.get(uuid) as User);
      }
      notExisted.push(uuid);
    }
  }

  if (notExisted.length === 0) {
    return;
  }

  if (notExisted.length === inCache.length) {
    return setOptions(inCache);
  }

  const data = (
    await ccStore.users_c.getMembers(
      GetMembersRequest.fromJson({
        uuids: notExisted || [],
      })
    )
  ).toJson() as any as Users;

  setOptions(data.users);
}

const renderTag: SelectRenderTag = ({ option, handleClose }) => {
  return h(
    NTag,
    {
      type: option.type as "success" | "warning" | "error",
      closable: true,
      onMousedown: (e: FocusEvent) => {
        e.preventDefault();
      },
      onClose: (e: MouseEvent) => {
        e.stopPropagation();
        handleClose();
      },
    },
    {
      default: () => (option.value !== option.label ? option.label : "Deleted"),
    }
  );
};

function filterOption(pattern: string, option: SelectOption) {
  return (
    option.label!.toString().toLowerCase().includes(pattern.toLowerCase()) &&
    !props.value.includes(option.value!.toString())
  );
}

watch(value, fetchExisted);
</script>

<style scoped></style>
