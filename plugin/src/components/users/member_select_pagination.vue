<template>
  <n-select
    :consistent-menu-width="false"
    :render-tag="renderTag"
    :value="internalValue"
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
import { onMounted, ref, toRefs, watch, h, computed } from "vue";
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
const internalValue = computed<string[]>({
  get() {
    return Array.from(new Set((props.value || []).map((v: any) => v?.toString())));
  },
  set(v: any) {
    const uniq = Array.from(new Set((v || []).map((x: any) => x?.toString())));
    emit("update:value", uniq);
  },
});
const emit = defineEmits(["update:value", "on-search"]);

const ccStore = useCcStore();
const usersStore = useUsersStore();

const options = ref<SelectOption[]>([]);
const isMembersLoading = ref(false);
const cachedEmpty = ref<User[]>([]);

// Debug helpers — enable by setting `localStorage.setItem('msp_debug','1')`
function maskId(id?: any) {
  try {
    const s = id?.toString() || "";
    if (!s) return "";
    return s.length <= 8 ? s : `${s.slice(0, 4)}…${s.slice(-4)}`;
  } catch (e) {
    return "";
  }
}

function dbg(...args: any[]) {
    // eslint-disable-next-line no-console
    console.debug("[member_select_pagination]", ...args);
}

onMounted(() => {
  dbg("mounted props.value", (props.value || []).map(maskId));
  fetchExisted();
});

function startSearch(val: string) {
  isMembersLoading.value = true;
  onSearchDebounced(val);
}

function onUpdate(value: any) {
  dbg("onUpdate incoming", (value || []).map(maskId));
  const uniq = Array.from(new Set((value || []).map((v: any) => v?.toString())));
  dbg("onUpdate deduped", uniq.map(maskId));
  emit("update:value", uniq);

  onSearch();
}

const onSearchDebounced = debounce(onSearch, 600);

async function onSearch(value?: string) {
  isMembersLoading.value = true;

  if (!value) {
    value = "";
  }

  value = value.trim();

  dbg("onSearch called", { value });

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
        }),
      )
    ).toJson() as any as Users;

    if (value === "") {
      cachedEmpty.value = data.users;
    }

    dbg("onSearch result count", data.users?.length);
    setOptions(data.users);
  } finally {
    isMembersLoading.value = false;
  }
}

function onUpdateShow() {
  if (options.value.length === 0 || cachedEmpty.value.length === 0) {
    dbg("onUpdateShow triggering search, options.length", options.value.length, "cachedEmpty.length", cachedEmpty.value.length);
    onSearch("");
  }
}

function setOptions(users?: User[]) {
  dbg("setOptions start", { usersCount: (users || []).length, propsValueCount: props.value.length, optionsBefore: options.value.length });

  const mapped = (users || []).map((user) => ({
    label: user.title,
    value: user.uuid?.toString(),
  }));
  dbg("mapped sample", mapped.slice(0, 5).map((m) => ({ label: m.label, value: maskId(m.value) }))); 

  const byId = new Map<string, SelectOption>();

  // Preserve selected items in the order of props.value
  for (const uuid of props.value) {
    const key = uuid?.toString() || "";
    if (!key) continue;
    let opt = options.value.find((o) => o.value?.toString() === key);
    if (!opt) opt = mapped.find((m) => m.value === key);
    if (!opt && usersStore.users.get(key)) {
      const u = usersStore.users.get(key) as User;
      opt = { label: u.title, value: key } as SelectOption;
    }
    if (opt) byId.set(key, opt);
  }

  // Add fetched/mapped items after selected ones
  for (const opt of mapped) {
    const key = opt.value?.toString() || "";
    if (!byId.has(key)) byId.set(key, opt);
  }

  dbg("byId size before set", byId.size, "sample keys", Array.from(byId.keys()).slice(0, 8).map(maskId));

  options.value = Array.from(byId.values());

  dbg("setOptions end", { optionsAfter: options.value.length, sample: options.value.slice(0, 8).map((o) => ({ label: o.label, value: maskId(o.value) })) });
}

async function fetchExisted() {
  const notExisted: string[] = [];
  const inCache: User[] = [];

  dbg("fetchExisted start props.value", (props.value || []).map(maskId));

  for (const uuid of props.value) {
    if (options.value.findIndex((option) => option.value === uuid) == -1) {
      if (usersStore.users.get(uuid)) {
        inCache.push(usersStore.users.get(uuid) as User);
      }
      notExisted.push(uuid);
    }
  }

  dbg("fetchExisted computed", { notExisted: notExisted.map(maskId), inCacheCount: inCache.length });

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
      }),
    )
  ).toJson() as any as Users;

  dbg("fetchExisted remote returned", data?.users?.length);

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
    },
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
