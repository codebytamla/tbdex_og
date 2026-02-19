<script setup lang="ts">
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
</script>

<template>
  <PageSectionHeader
    title="Screenings"
    description="Review screening queue context and readiness for follow-up actions"
  />

  <Card>
    <template #content>
      <p class="helper-copy">
        Screening workflows will be surfaced on this page.
      </p>
    </template>
  </Card>
</template>

<style scoped>
.helper-copy {
  margin: 0;
  color: var(--p-surface-700);
}
</style>
