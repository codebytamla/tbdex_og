<script setup lang="ts">
defineProps<{
  searchQuery?: string | null
  searchPlaceholder?: string
  searchDisabled?: boolean
  columns: Array<{ field: string; header: string }>
  selectedColumns: Array<{ field: string; header: string }>
  columnsLabel?: string
}>()

defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:selectedColumns', value: Array<{ field: string; header: string }>): void
  (e: 'search'): void
}>()
</script>

<template>
  <div class="datatable-toolbar-placeholder">DataTable header placeholder</div>
</template>

<style scoped>
.datatable-toolbar-placeholder {
  width: 100%;
  min-height: 2.25rem;
  display: flex;
  align-items: center;
  color: var(--p-surface-500);
  font-size: 0.83rem;
}
</style>
