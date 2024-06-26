<template>
  <n-dropdown
    trigger="hover"
    :render-option="renderOption"
    :options="membersOptions"
    :show="visible"
  >
    <slot>
      <n-text v-if="admin?.title">{{ admin.title }}</n-text>
      <n-text v-else>{{ members.length }} members</n-text>
    </slot>
  </n-dropdown>
</template>

<script setup lang="ts">
import { Component, computed, h, useSlots } from 'vue'
import { DropdownOption, NDropdown, NText } from 'naive-ui'
import { User } from '../../connect/cc/cc_pb.ts'
import memberItem from './user_item.vue'
import AddButton from '../ui/add_button.vue'

interface MembersDropdownProps {
  visible?: boolean
  responsible?: User
  members: User[]
  admins: String[]
}

const props = defineProps<MembersDropdownProps>()
const emits = defineEmits(['delete', 'add', 'change'])
const slots = useSlots()

const membersOptions = computed(() => {
  const membersOptions: DropdownOption[] = []
  
  props.members.forEach((member) => {
    if (membersOptions.find(({ key }) => key === member?.uuid)) return
    membersOptions.push({
      key: member?.uuid,
      extra: member as any,
      label: member?.title ?? 'Unknown',
    })
  })

  membersOptions.push({ key: 'add', label: 'add' })
  return membersOptions
})

const admin = computed(() => (props.responsible ?? true)
  ? props.responsible
  : props.members.find(({ uuid } = {} as User) => props.admins.includes(uuid))
)

const renderOption = ({ option }: { option: DropdownOption }) => {
  if (option.key === 'add') {
    return h(
      AddButton,
      { style: 'width: 100%', onClick: () => emits('add') },
      () => 'Add'
    )
  }

  return h(memberItem as Component, {
    user: option.extra ?? { title: 'unknown' },
    style: ((option.key === props.responsible?.uuid)
      ? 'background: var(--n-text-color); color: var(--n-color)'
      : null
    ),
    onDelete: () => emits('delete', option.key),
    onChange: () => emits('change', option.key),
    action: (option.key === props.responsible?.uuid) ? 'change' : 'delete',
    actions: true
  }, slots)
}
</script>
