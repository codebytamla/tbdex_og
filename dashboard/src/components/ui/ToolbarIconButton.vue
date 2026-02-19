<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    disabled?: boolean
    title?: string
    ariaLabel?: string
  }>(),
  {
    disabled: false,
    title: '',
    ariaLabel: '',
  },
)

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

function onClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <button
    class="toolbar-icon-btn"
    type="button"
    :disabled="props.disabled"
    :title="props.title || props.ariaLabel || undefined"
    :aria-label="props.ariaLabel || props.title || 'Action button'"
    @click="onClick"
  >
    <slot />
  </button>
</template>

<style scoped>
.toolbar-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: var(--toolbar-control-height, var(--app-control-height, 2.2rem));
  height: var(--toolbar-control-height, var(--app-control-height, 2.2rem));
  border: 1px solid var(--p-surface-200);
  border-radius: var(--app-control-radius, 0.62rem);
  background: var(--p-surface-50);
  color: var(--p-surface-700);
  cursor: pointer;
  transition: all 0.15s;
}

.toolbar-icon-btn:hover:not(:disabled) {
  border-color: var(--p-surface-400);
  background: var(--p-surface-100);
}

.toolbar-icon-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
