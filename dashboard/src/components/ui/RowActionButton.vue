<script setup lang="ts">
import Button from 'primevue/button'

const props = withDefaults(
  defineProps<{
    disabled?: boolean
    loading?: boolean
    title?: string
    ariaLabel?: string
  }>(),
  {
    disabled: false,
    loading: false,
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
  <Button
    class="row-action-btn"
    text
    size="small"
    :disabled="props.disabled"
    :loading="props.loading"
    :title="props.title || props.ariaLabel || 'Row actions'"
    :aria-label="props.ariaLabel || props.title || 'Row actions'"
    @click="onClick"
  >
    <slot />
  </Button>
</template>

<style scoped>
.row-action-btn {
  color: #9ca3af !important;
  border-radius: var(--radius-sm, 6px) !important;
  transition: color 0.15s ease, background 0.15s ease !important;
}

.row-action-btn:hover:not(:disabled) {
  color: #374151 !important;
  background: #f4f5f7 !important;
}

.row-action-btn:focus-visible {
  outline: 2px solid rgb(13 148 136 / 0.3);
  outline-offset: 1px;
}

.row-action-btn :deep(.p-button-label) {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}
</style>
