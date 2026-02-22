<script setup lang="ts">
import { computed, ref } from 'vue'
import { Plus, ChevronDown } from 'lucide-vue-next'
import Menu from 'primevue/menu'
import type { MenuItem } from 'primevue/menuitem'
import Button from 'primevue/button'

const props = withDefaults(
  defineProps<{
    label: string
    disabled?: boolean
    loading?: boolean
    title?: string
    ariaLabel?: string
    importLabel?: string
    exportLabel?: string
    refreshLabel?: string
    importDisabled?: boolean
    exportDisabled?: boolean
    refreshDisabled?: boolean
  }>(),
  {
    disabled: false,
    loading: false,
    title: '',
    ariaLabel: '',
    importLabel: 'Import',
    exportLabel: 'Export',
    refreshLabel: 'Refresh',
    importDisabled: true,
    exportDisabled: true,
    refreshDisabled: false,
  },
)

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
  (e: 'refresh'): void
}>()

function onPrimaryClick(event: MouseEvent) {
  emit('click', event)
}

const actionMenu = ref()

const actionItems = computed<MenuItem[]>(() => [
  {
    label: props.importLabel,
    icon: 'pi pi-upload',
    disabled: props.importDisabled,
  },
  {
    label: props.exportLabel,
    icon: 'pi pi-download',
    disabled: props.exportDisabled,
  },
  {
    label: props.refreshLabel,
    icon: 'pi pi-refresh',
    disabled: props.refreshDisabled,
    command: () => {
      emit('refresh')
    },
  },
])

const menuDisabled = computed(() =>
  props.disabled || (props.importDisabled && props.exportDisabled && props.refreshDisabled),
)

function toggleMenu(event: MouseEvent) {
  actionMenu.value?.toggle(event)
}
</script>

<template>
  <div class="primary-action-dropdown">
    <Button
      class="primary-btn-main"
      :disabled="props.disabled"
      :loading="props.loading"
      :title="props.title || props.label"
      :aria-label="props.ariaLabel || props.label"
      @click="onPrimaryClick"
    >
      <Plus :size="16" :stroke-width="2.3" />
      <span>{{ props.label }}</span>
    </Button>
    <Button
      class="primary-btn-split"
      :disabled="menuDisabled"
      aria-label="More actions"
      title="More actions"
      @click="toggleMenu"
      aria-haspopup="true"
      aria-controls="primary_action_menu"
    >
      <ChevronDown :size="16" :stroke-width="2.3" />
    </Button>

    <Menu ref="actionMenu" id="primary_action_menu" :model="actionItems" popup />
  </div>
</template>

<style scoped>
.primary-action-dropdown {
  display: inline-flex;
  align-items: center;
  height: 2.375rem;
  border-radius: var(--radius-md, 10px);
  background: #0f172a;
  padding: 0;
  overflow: hidden;
  box-shadow: 0 1px 3px 0 rgb(15 23 42 / 0.2), 0 1px 2px -1px rgb(15 23 42 / 0.12);
  transition: box-shadow 0.15s ease;
}

.primary-action-dropdown:hover {
  box-shadow: 0 4px 12px -2px rgb(15 23 42 / 0.25), 0 2px 4px -2px rgb(15 23 42 / 0.15);
}

.primary-btn-main {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  height: 100%;
  padding: 0 1rem 0 1.125rem;
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 0.855rem;
  font-weight: 600;
  line-height: 1;
  letter-spacing: -0.005em;
  white-space: nowrap;
  border-radius: 0;
  border-right: 1px solid rgb(255 255 255 / 0.15);
  transition: background 0.15s ease;
}

.primary-btn-main:hover:not(:disabled) {
  background: rgb(255 255 255 / 0.08);
  color: #ffffff;
}

.primary-btn-main :deep(.p-button-label) {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
}

.primary-btn-split {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.375rem;
  height: 100%;
  background: transparent;
  border: none;
  color: rgb(255 255 255 / 0.7);
  padding: 0;
  border-radius: 0;
  transition: background 0.15s ease, color 0.15s ease;
}

.primary-btn-split:hover:not(:disabled) {
  background: rgb(255 255 255 / 0.08);
  color: #ffffff;
}
</style>
