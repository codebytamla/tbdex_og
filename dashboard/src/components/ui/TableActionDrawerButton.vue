<script setup lang="ts">
import { computed, ref } from 'vue'
import { EllipsisVertical } from 'lucide-vue-next'
import Menu from 'primevue/menu'
import type { MenuItem } from 'primevue/menuitem'
import ToolbarIconButton from '@/components/ui/ToolbarIconButton.vue'

const props = withDefaults(
  defineProps<{
    disabled?: boolean
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
    title: 'Table actions',
    ariaLabel: 'Open table actions',
    importLabel: 'Import',
    exportLabel: 'Export',
    refreshLabel: 'Refresh',
    importDisabled: true,
    exportDisabled: true,
    refreshDisabled: false,
  },
)

const emit = defineEmits<{
  (e: 'refresh'): void
}>()

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

const triggerDisabled = computed(() =>
  props.disabled || (props.importDisabled && props.exportDisabled && props.refreshDisabled),
)

function toggleMenu(event: MouseEvent) {
  actionMenu.value?.toggle(event)
}
</script>

<template>
  <ToolbarIconButton
    :disabled="triggerDisabled"
    :title="props.title"
    :ariaLabel="props.ariaLabel"
    @click="toggleMenu"
  >
    <EllipsisVertical :size="16" />
  </ToolbarIconButton>

  <Menu ref="actionMenu" :model="actionItems" popup />
</template>
