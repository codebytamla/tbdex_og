<script setup lang="ts">
const props = defineProps<{
  label: string
}>()
</script>

<template>
  <label class="filter-field">
    <span class="filter-field-label">{{ props.label }}</span>
    <slot />
  </label>
</template>

<style scoped>
.filter-field {
  display: flex;
  flex-direction: column;
  gap: 0.32rem;
}

.filter-field-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--p-surface-600);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.filter-field :deep(.p-dropdown),
.filter-field :deep(.p-multiselect) {
  width: 100%;
  height: var(--filter-control-height, var(--toolbar-control-height, var(--app-control-height, 2.2rem)));
  min-height: var(--filter-control-height, var(--toolbar-control-height, var(--app-control-height, 2.2rem)));
  box-sizing: border-box;
  border-radius: var(--app-control-radius, 0.62rem);
  font-size: var(--app-control-font-size, 0.83rem);
}

.filter-field :deep(.p-dropdown-label),
.filter-field :deep(.p-multiselect-label) {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 0.88rem;
  font-size: var(--app-control-font-size, 0.83rem);
  line-height: 1.2;
  white-space: nowrap;
}

.filter-field :deep(.p-dropdown-trigger),
.filter-field :deep(.p-multiselect-dropdown) {
  width: 2.5rem;
  min-width: 2.5rem;
  height: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>
