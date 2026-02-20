<script setup lang="ts">
import { computed, ref } from 'vue'
import { MoreHorizontal } from 'lucide-vue-next'
import InputText from 'primevue/inputtext'
import MultiSelect from 'primevue/multiselect'
import Button from 'primevue/button'
import Menu from 'primevue/menu'

const props = defineProps<{
  // Search query
  searchQuery?: string | null
  searchPlaceholder?: string
  searchDisabled?: boolean

  // Columns toggling
  columns: Array<{ field: string; header: string }>
  selectedColumns: Array<{ field: string; header: string }>
  columnsLabel?: string

  // Options Menu toggles
  showImport?: boolean
  showExport?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:selectedColumns', value: Array<{ field: string; header: string }>): void
  (e: 'import'): void
  (e: 'export'): void
}>()

const optionsMenu = ref()

const optionsMenuItems = computed(() => {
  const items = []
  if (props.showImport) {
    items.push({
      label: 'Import Data',
      command: () => emit('import')
    })
  }
  if (props.showExport) {
    items.push({
      label: 'Export Data',
      command: () => emit('export')
    })
  }
  return items
})

const localSearchQuery = computed({
  get: () => props.searchQuery ?? '',
  set: (val: string) => emit('update:searchQuery', val)
})

const localSelectedColumns = computed({
  get: () => props.selectedColumns,
  set: (val: Array<{ field: string; header: string }>) => emit('update:selectedColumns', val)
})

function toggleOptionsMenu(event: Event) {
  optionsMenu.value.toggle(event)
}
</script>

<template>
  <div class="datatable-toolbar">
    <div class="toolbar-left">
      <InputText
        v-model="localSearchQuery"
        :placeholder="searchPlaceholder || 'Search...'"
        :disabled="searchDisabled"
        class="toolbar-search"
      />
    </div>

    <div class="toolbar-right">
      <MultiSelect
        v-model="localSelectedColumns"
        :options="columns"
        optionLabel="header"
        :placeholder="columnsLabel || 'Columns'"
        :maxSelectedLabels="1"
        selectedItemsLabel="{0} columns"
        class="toolbar-columns-select"
        display="comma"
      />

      <Button
        v-if="optionsMenuItems.length > 0"
        type="button"
        icon="pi pi-ellipsis-h"
        class="toolbar-options-btn p-button-text p-button-secondary"
        aria-haspopup="true"
        aria-controls="options_menu"
        @click="toggleOptionsMenu"
      >
        <template #icon>
          <MoreHorizontal :size="18" />
        </template>
      </Button>

      <Menu ref="optionsMenu" id="options_menu" :model="optionsMenuItems" :popup="true" />
    </div>
  </div>
</template>

<style scoped>
.datatable-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  width: 100%;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toolbar-search {
  width: 100%;
  max-width: 17.33rem;
  height: var(--app-control-height, 2.25rem);
  box-sizing: border-box;
  border-radius: var(--app-control-radius, 0.62rem);
  padding: 0 0.88rem;
  font-size: var(--app-control-font-size, 0.83rem);
  line-height: 1.2;
  border: 1px solid var(--p-surface-200);
  background: var(--p-surface-50);
  transition: border-color 0.15s, background 0.15s;
}

.toolbar-search:focus {
  background: var(--p-surface-0);
  border-color: var(--p-surface-400);
}

.toolbar-columns-select {
  min-width: 10rem;
  height: var(--app-control-height, 2.25rem);
}

:deep(.toolbar-columns-select .p-multiselect-label) {
  padding: 0.4rem 0.88rem;
  font-size: var(--app-control-font-size, 0.83rem);
}

.toolbar-options-btn {
  height: var(--app-control-height, 2.25rem);
  width: var(--app-control-height, 2.25rem);
  padding: 0;
  color: var(--p-surface-600);
}

.toolbar-options-btn:hover {
  background: var(--p-surface-100);
  color: var(--p-surface-800);
}

@media (max-width: 640px) {
  .datatable-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .toolbar-left,
  .toolbar-right {
    justify-content: space-between;
    width: 100%;
  }

  .toolbar-search {
    max-width: 100%;
  }
}
</style>
