<template>
  <n-select :consistent-menu-width="false" :render-tag="renderTag" :value="value"
            @update-value="emit('update:value',$event)" multiple
            :options="props.options"
            filterable/>
</template>

<script setup lang="ts">
import {NSelect, NTag, SelectOption, SelectRenderTag} from "naive-ui";
import {h} from "vue";

interface MemberSelectProps {
  options: SelectOption[]
  value: string[]
}

const props = defineProps<MemberSelectProps>()
const emit = defineEmits(['update:value'])

const renderTag: SelectRenderTag = ({option, handleClose}) => {
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
      {default: () => option.value !== option.label ? option.label : 'Deleted'}
  )
}
</script>

<style scoped>

</style>