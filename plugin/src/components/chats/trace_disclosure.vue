<template>
  <details
    class="disc"
    :class="[`disc--${accent}`, { 'disc--dense': dense }]"
    :open="open"
  >
    <summary class="disc__sum">
      <slot name="summary" />
      <svg class="chev" viewBox="0 0 16 16" width="11" height="11" aria-hidden="true">
        <path
          d="M6 4l4 4-4 4"
          fill="none"
          stroke="currentColor"
          stroke-width="1.7"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </summary>
    <div class="disc__body"><slot /></div>
  </details>
</template>

<script setup lang="ts">
// A single collapsible box: header (the #summary slot) + a chevron that rotates
// when open, and a body (default slot). Every foldable section in the trace
// viewer is built from this so the summary/chevron/reset boilerplate lives once.
withDefaults(
  defineProps<{
    accent?: "neutral" | "tool" | "error";
    open?: boolean;
    dense?: boolean;
  }>(),
  { accent: "neutral", open: false, dense: false },
);
</script>

<style scoped>
.disc {
  border: 1px solid rgba(128, 128, 128, 0.24);
  border-radius: 7px;
  overflow: hidden;
  background: rgba(128, 128, 128, 0.03);
}
.disc--tool {
  border-color: rgba(224, 145, 47, 0.3);
  background: rgba(224, 145, 47, 0.05);
}
.disc--error {
  border-color: rgba(208, 48, 80, 0.4);
  background: rgba(208, 48, 80, 0.06);
}
.disc--dense {
  border-radius: 6px;
}
.disc__sum {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 6px 9px;
  cursor: pointer;
  list-style: none;
  user-select: none;
}
.disc--dense .disc__sum {
  padding: 5px 8px;
}
.disc__sum::-webkit-details-marker {
  display: none;
}
.disc__sum:hover {
  background: rgba(128, 128, 128, 0.06);
}
.disc--tool .disc__sum:hover {
  background: rgba(224, 145, 47, 0.08);
}
.chev {
  opacity: 0.4;
  transition: transform 0.18s ease;
  flex: none;
}
.disc[open] > .disc__sum .chev {
  transform: rotate(90deg);
}
.disc__body {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 2px 9px 9px;
}
.disc--dense .disc__body {
  padding: 2px 8px 8px;
}
@media (prefers-reduced-motion: reduce) {
  .chev {
    transition: none;
  }
}
</style>
