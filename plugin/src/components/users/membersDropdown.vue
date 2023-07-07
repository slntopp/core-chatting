<template>
  <n-dropdown :render-option="renderOption" trigger="hover" :options="membersOptions">
    <n-text>{{ members.length }} members</n-text>
  </n-dropdown>
</template>

<script setup lang="ts">
import {NDropdown, NText,DropdownOption} from "naive-ui";
import {User} from "../../connect/cc/cc_pb.ts";
import {computed, h, toRefs, VNode} from "vue";
import memberItem from './userItem.vue'

interface MembersDropdownProps{
  members:User[]
}

const props=defineProps<MembersDropdownProps>()
const {members}=toRefs(props)

const emit=defineEmits(['delete'])

const membersOptions=computed(()=>{
    return [...new Set(members.value)].map((m) => ({
      key: m.uuid,
      extra:m,
      label:m.title || 'Unknown',
    }))
})

const renderOption=({option}:{option:DropdownOption})=>{
  return h(memberItem,{user:option.extra,onDelete:()=>emit('delete',option.key)})
}
</script>

<style scoped>

</style>