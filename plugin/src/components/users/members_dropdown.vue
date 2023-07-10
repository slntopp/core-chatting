<template>
  <n-dropdown :render-option="renderOption" trigger="hover" :options="membersOptions">
    <n-text>{{ members.length }} members</n-text>
  </n-dropdown>
</template>

<script setup lang="ts">
import {DropdownOption, NDropdown, NText} from "naive-ui";
import {User} from "../../connect/cc/cc_pb.ts";
import {computed, h, toRefs} from "vue";
import memberItem from './user_item.vue'
import AddButton from "../ui/add_button.vue";

interface MembersDropdownProps {
  members: User[]
}

const props = defineProps<MembersDropdownProps>()
const {members} = toRefs(props)

const emit = defineEmits(['delete', 'add'])

const membersOptions = computed(() => {
  const membersOptions = [...new Set(members.value)].map((m) => ({
    key: m.uuid,
    extra: m as any,
    label: m.title || 'Unknown',
  }))
  membersOptions.push({key: 'add', label: 'add', extra: {} as User})
  return membersOptions
})

const renderOption = ({option}: { option: DropdownOption }) => {
  if (option.key === 'add') {
    return h(AddButton, {style: {width: '100%'}, onClick: addMember}, 'Add')
  }
  return h(memberItem as any, {user: option.extra, onDelete: () => emit('delete', option.key)})
}

const addMember = () => {
  emit('add')
}
</script>

<style scoped>

</style>