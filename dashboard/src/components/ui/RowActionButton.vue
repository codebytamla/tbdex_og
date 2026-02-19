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
  color: #52525b !important;
}

.row-action-btn:hover:not(:disabled) {
  color: #3f3f46 !important;
  background: #f4f4f5 !important;
}

.row-action-btn:focus-visible {
  outline: 2px solid #d4d4d8;
  outline-offset: 1px;
}

.row-action-btn :deep(.p-button-label) {
  display: inline-flex;
  align-items: center;
  gap: 0.28rem;
}
</style>
