<script setup lang="ts">
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
</script>

<template>
  <PageSectionHeader
    title="Screenings"
    description="Page temporarily paused for redesign"
  />

  <Card>
    <template #content>
      <p class="helper-copy">Screenings redesign is scheduled for tomorrow.</p>
    </template>
  </Card>
</template>

<style scoped>
.helper-copy {
  margin: 0;
  color: var(--p-surface-700);
}
</style>
