<template>
  <n-tooltip v-if="grade" placement="bottom">
    <template #trigger>
      <n-tag
        size="tiny"
        :bordered="grade.manual"
        :color="tagColor"
        class="account-grade"
        @click.stop
      >
        {{ grade.title }}
      </n-tag>
    </template>
    <div>{{ spendingText }}</div>
    <div v-if="grade.manual">Grade is set manually</div>
  </n-tooltip>
</template>

<script setup lang="ts">
import { computed, toRefs } from "vue";
import { NTag, NTooltip } from "naive-ui";

// The users store keeps users as JSON (camelCase). The server sets grade only
// for operators, so clients never see it.
interface GradeJson {
  title?: string;
  color?: string;
  manual?: boolean;
  amount?: number;
  currency?: string;
  periodDays?: number;
}

const props = defineProps<{ user?: { grade?: GradeJson } | null }>();
const { user } = toRefs(props);

const grade = computed(() => {
  const g = user.value?.grade;
  return g?.title ? g : null;
});

const tagColor = computed(() => {
  const color = grade.value?.color || "#757575";
  return grade.value?.manual
    ? { color: "transparent", textColor: color, borderColor: color }
    : { color, textColor: "#fff" };
});

const spendingText = computed(() => {
  const g = grade.value;
  if (!g || g.amount === undefined) return "Not counted yet";
  const months = Math.round(((g.periodDays || 395) * 12) / 365);
  return `Pays ${(+g.amount).toFixed(2)} ${g.currency ?? ""} per month on average (payments of the last ${months} months)`;
});
</script>

<style scoped>
.account-grade {
  font-weight: 700;
  flex: none;
}
</style>
