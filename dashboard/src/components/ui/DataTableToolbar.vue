<script setup lang="ts">
import { computed } from 'vue'
import { ChevronDown, Columns2 } from 'lucide-vue-next'
import InputText from 'primevue/inputtext'
import MultiSelect from 'primevue/multiselect'

const props = defineProps<{
  searchQuery?: string | null
  searchPlaceholder?: string
  searchDisabled?: boolean
  columns: Array<{ field: string; header: string }>
  selectedColumns: Array<{ field: string; header: string }>
  columnsLabel?: string
}>()

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string | null): void
  (e: 'update:selectedColumns', value: Array<{ field: string; header: string }>): void
  (e: 'search'): void
}>()

const searchValue = computed({
  get: () => props.searchQuery ?? '',
  set: (val) => emit('update:searchQuery', val || null),
})

const columnSelection = computed({
  get: () => props.selectedColumns,
  set: (val) => emit('update:selectedColumns', val),
})

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter') emit('search')
}
</script>

<template>
  <div class="datatable-toolbar">
    <InputText
      v-model="searchValue"
      :placeholder="searchPlaceholder ?? 'Search...'"
      :disabled="searchDisabled"
      class="toolbar-search"
      @keydown="onKeydown"
    />
    <MultiSelect
      v-model="columnSelection"
      :options="columns"
      optionLabel="header"
      :placeholder="columnsLabel ?? 'Columns'"
      class="toolbar-columns"
      :maxSelectedLabels="0"
      selectedItemsLabel=""
    >
      <template #value>
        <span class="columns-trigger-content">
          <Columns2 :size="15" class="columns-trigger-icon" aria-hidden="true" />
          <span>{{ columnsLabel ?? 'Columns' }}</span>
        </span>
      </template>
      <template #dropdownicon>
        <ChevronDown :size="13" />
      </template>
    </MultiSelect>
  </div>
</template>

<style scoped>
.datatable-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  width: 100%;
}

.toolbar-search {
  width: 20rem;
  height: var(--app-control-height, 2.25rem);
  border-radius: var(--app-control-radius, 8px);
  padding: 0 0.875rem;
  font-size: var(--app-control-font-size, 0.855rem);
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  box-shadow: none;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}
.toolbar-search:hover {
  border-color: #d1d5db;
  background: #ffffff;
}
.toolbar-search:focus {
  border-color: var(--p-primary-400, #2dd4bf);
  background: #ffffff;
  box-shadow: 0 0 0 3px rgb(13 148 136 / 0.08);
}

:deep(.toolbar-columns.p-multiselect) {
  height: var(--app-control-height, 2.25rem);
  border: 1px solid #e5e7eb;
  border-radius: var(--app-control-radius, 8px);
  background: #ffffff;
  font-size: var(--app-control-font-size, 0.855rem);
  font-weight: 600;
  color: #374151;
  padding: 0 0.75rem;
  gap: 0.375rem;
  box-shadow: var(--shadow-xs, 0 1px 2px 0 rgb(0 0 0 / 0.04));
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
:deep(.toolbar-columns.p-multiselect:hover) {
  border-color: #d1d5db;
}
:deep(.toolbar-columns.p-multiselect.p-focus),
:deep(.toolbar-columns.p-multiselect:focus-within) {
  border-color: var(--p-primary-400, #2dd4bf);
  box-shadow: 0 0 0 3px rgb(13 148 136 / 0.08);
}
:deep(.toolbar-columns .p-multiselect-label-container) {
  display: flex;
  align-items: center;
}
:deep(.toolbar-columns .p-multiselect-label) {
  padding: 0;
  display: flex;
  align-items: center;
}
:deep(.toolbar-columns .p-multiselect-dropdown) {
  padding: 0;
  width: auto;
  color: #9ca3af;
}

.columns-trigger-content {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  white-space: nowrap;
  color: #374151;
  line-height: 1;
}
.columns-trigger-icon {
  color: #9ca3af;
}
</style>
