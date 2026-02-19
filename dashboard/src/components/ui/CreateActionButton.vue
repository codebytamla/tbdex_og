<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import Button from 'primevue/button'

const props = withDefaults(
  defineProps<{
    label: string
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
    class="create-action-btn"
    :disabled="props.disabled"
    :loading="props.loading"
    :title="props.title || props.label"
    :aria-label="props.ariaLabel || props.label"
    @click="onClick"
  >
    <Plus :size="16" :stroke-width="2.3" />
    <span>{{ props.label }}</span>
  </Button>
</template>

<style scoped>
.create-action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  height: var(--toolbar-control-height, var(--app-control-height, 2.2rem));
  padding: 0 0.92rem;
  background: var(--p-primary-700);
  border-color: var(--p-primary-700);
  color: #ffffff;
  border-radius: var(--app-control-radius, 0.68rem);
  font-size: var(--app-control-font-size, 0.83rem);
  font-weight: 600;
  line-height: 1;
  letter-spacing: 0.002em;
  white-space: nowrap;
}

.create-action-btn:hover:not(:disabled) {
  background: var(--p-primary-800);
  border-color: var(--p-primary-800);
  color: #ffffff;
}

.create-action-btn :deep(.p-button-label) {
  display: inline-flex;
  align-items: center;
  gap: 0.36rem;
}
</style>
