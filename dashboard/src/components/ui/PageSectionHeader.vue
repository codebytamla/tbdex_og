<script setup lang="ts">
import { computed, useSlots } from 'vue'

const props = withDefaults(
  defineProps<{
    title: string
    description?: string
    compact?: boolean
  }>(),
  {
    description: '',
    compact: false,
  },
)

const slots = useSlots()
const hasActions = computed(() => !!slots.actions)
</script>

<template>
  <section class="page-section-header" :class="{ compact: props.compact }">
    <div class="page-section-top">
      <div class="page-section-copy">
        <h1 class="page-section-title">{{ props.title }}</h1>
        <p v-if="props.description" class="page-section-description">{{ props.description }}</p>
      </div>
      <div v-if="hasActions" class="page-section-actions">
        <slot name="actions" />
      </div>
    </div>
  </section>
</template>

<style scoped>
.page-section-header {
  margin-bottom: var(--page-section-margin-bottom, 1.5rem);
}

.page-section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.page-section-copy {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  min-width: 0;
}

.page-section-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.025em;
  color: #111827;
}

.page-section-description {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 400;
  line-height: 1.5;
  color: #6b7280;
}

.page-section-actions {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.page-section-header.compact {
  margin-bottom: 1rem;
}

.page-section-header.compact .page-section-title {
  font-size: 1.25rem;
}

.page-section-header.compact .page-section-description {
  font-size: 0.825rem;
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
    font-size: 1.25rem;
  }

  .page-section-description {
    font-size: 0.825rem;
  }
}
</style>
