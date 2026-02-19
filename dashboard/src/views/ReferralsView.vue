<script setup lang="ts">
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
</script>

<template>
  <PageSectionHeader
    title="Referrals"
    description="View referral traffic and care handoffs between service points"
  />

  <Card>
    <template #content>
      <p class="helper-copy">
        Referral workflows will be surfaced on this page.
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
