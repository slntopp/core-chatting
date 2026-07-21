<template>
  <div class="tb" :class="`tb--${tone}`">
    <div class="tb__bar">
      <span class="tb__label">{{ label }}</span>
      <span class="tb__spacer" />
      <span v-if="meta" class="tb__meta">{{ meta }}</span>
      <span class="tb__chars">{{ chars }}</span>
      <button class="tb__copy" type="button" @click="copy" :title="copied ? 'Copied' : 'Copy'">
        {{ copied ? "copied" : "copy" }}
      </button>
    </div>
    <div class="tb__body" :class="{ 'tb__body--clamped': clampable && !expanded }">
      <pre class="tb__pre">{{ content }}</pre>
      <div v-if="clampable && !expanded" class="tb__fade" />
    </div>
    <button v-if="clampable" class="tb__more" type="button" @click="expanded = !expanded">
      {{ expanded ? "Show less" : `Show all · ${chars}` }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

const props = withDefaults(
  defineProps<{
    label: string;
    content: string;
    tone?: "code" | "muted" | "error" | "accent";
    meta?: string;
  }>(),
  { tone: "code", meta: "" },
);

const expanded = ref(false);
const copied = ref(false);

// Long values collapse behind a fade with a Show all toggle; short ones render
// whole. A char threshold avoids measuring layout and is stable enough here.
const clampable = computed(() => props.content.length > 480);

const chars = computed(() => {
  const n = props.content.length;
  return n >= 1000 ? `${(n / 1000).toFixed(1)}k chars` : `${n} chars`;
});

async function copy() {
  try {
    await navigator.clipboard.writeText(props.content);
    copied.value = true;
    setTimeout(() => (copied.value = false), 1200);
  } catch {
    /* clipboard blocked; ignore */
  }
}
</script>

<style scoped>
.tb {
  border: 1px solid var(--tb-border, rgba(128, 128, 128, 0.22));
  border-radius: 7px;
  overflow: hidden;
  background: rgba(128, 128, 128, 0.045);
}
.tb--error {
  border-color: rgba(208, 48, 80, 0.4);
  background: rgba(208, 48, 80, 0.06);
}
.tb--accent {
  border-color: rgba(124, 92, 255, 0.35);
}
.tb__bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 6px 4px 10px;
  border-bottom: 1px solid rgba(128, 128, 128, 0.16);
  background: rgba(128, 128, 128, 0.06);
}
.tb__label {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 10.5px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  opacity: 0.75;
}
.tb__spacer {
  flex: 1;
}
.tb__meta {
  font-size: 11px;
  opacity: 0.6;
}
.tb__chars {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  opacity: 0.4;
}
.tb__copy {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  letter-spacing: 0.04em;
  color: inherit;
  background: transparent;
  border: 1px solid rgba(128, 128, 128, 0.3);
  border-radius: 5px;
  padding: 1px 7px;
  cursor: pointer;
  opacity: 0.65;
  transition: opacity 0.15s, border-color 0.15s;
}
.tb__copy:hover {
  opacity: 1;
  border-color: rgba(128, 128, 128, 0.6);
}
.tb__copy:focus-visible {
  outline: 2px solid #2f6fed;
  outline-offset: 1px;
}
.tb__body {
  position: relative;
}
.tb__body--clamped {
  max-height: 200px;
  overflow: hidden;
}
.tb__pre {
  margin: 0;
  padding: 9px 11px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 460px;
  overflow: auto;
}
.tb--muted .tb__pre {
  opacity: 0.72;
  font-style: italic;
}
.tb--error .tb__pre {
  color: #e0607d;
}
.tb__fade {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 48px;
  background: linear-gradient(
    to bottom,
    transparent,
    var(--tb-fade, rgba(250, 250, 250, 0.9))
  );
  pointer-events: none;
}
@media (prefers-color-scheme: dark) {
  .tb__fade {
    background: linear-gradient(to bottom, transparent, rgba(24, 24, 28, 0.92));
  }
}
.tb__more {
  display: block;
  width: 100%;
  font-family: ui-monospace, monospace;
  font-size: 10.5px;
  letter-spacing: 0.04em;
  text-align: center;
  color: inherit;
  background: rgba(128, 128, 128, 0.08);
  border: none;
  border-top: 1px solid rgba(128, 128, 128, 0.16);
  padding: 4px;
  cursor: pointer;
  opacity: 0.7;
}
.tb__more:hover {
  opacity: 1;
  background: rgba(128, 128, 128, 0.14);
}
</style>
