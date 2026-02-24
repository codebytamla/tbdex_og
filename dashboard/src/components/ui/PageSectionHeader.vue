<script setup lang="ts">
import { computed, useSlots } from 'vue'

const props = withDefaults(
  defineProps<{
    title: string
    compact?: boolean
  }>(),
  {
    compact: false,
  },
)

const slots = useSlots()
const hasActions = computed(() => !!slots.actions)
</script>

<template>
  <section class="page-section-header" :class="{ compact: props.compact }">
    <div class="page-section-top">
      <h1 class="page-section-title">{{ props.title }}</h1>
      <div v-if="hasActions" class="page-section-actions">
        <slot name="actions" />
      </div>
    </div>
  </section>
</template>

<style scoped>
.page-section-header {
  margin-bottom: var(--page-section-margin-bottom, 1.25rem);
}

.page-section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.page-section-title {
  margin: 0;
  font-size: 1.65rem;
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.03em;
  color: #0d1117;
}

.page-section-header.compact .page-section-title {
  font-size: 1.35rem;
}

.page-section-actions {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.5rem;
  flex-wrap: wrap;
}

@media (max-width: 900px) {
  .page-section-top {
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;
  }

  .page-section-actions {
    margin-left: 0;
    justify-content: flex-end;
  }

  .page-section-title {
    font-size: 1.35rem;
  }
}
</style>
