<script setup lang="ts">
import { computed } from 'vue'
import InputText from 'primevue/inputtext'
import MultiSelect from 'primevue/multiselect'

const props = defineProps<{
  // Search query
  searchQuery?: string | null
  searchPlaceholder?: string
  searchDisabled?: boolean

  // Columns toggling
  columns: Array<{ field: string; header: string }>
  selectedColumns: Array<{ field: string; header: string }>
  columnsLabel?: string

}>()

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:selectedColumns', value: Array<{ field: string; header: string }>): void
  (e: 'search'): void
}>()

const localSearchQuery = computed({
  get: () => props.searchQuery ?? '',
  set: (val: string) => emit('update:searchQuery', val)
})

const localSelectedColumns = computed({
  get: () => props.selectedColumns,
  set: (val: Array<{ field: string; header: string }>) => emit('update:selectedColumns', val)
})
</script>

<template>
  <div class="datatable-toolbar">
    <div class="toolbar-left">
      <InputText
        v-model="localSearchQuery"
        :placeholder="searchPlaceholder || 'Search...'"
        :disabled="searchDisabled"
        class="toolbar-search"
        @keydown.enter.prevent="emit('search')"
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

.toolbar-left {
  flex: 1;
  min-width: 0;
}

.toolbar-search {
  width: min(100%, 24rem);
  max-width: none;
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
