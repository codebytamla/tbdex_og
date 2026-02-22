<script setup lang="ts">
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
</script>

<template>
  <PageSectionHeader
    title="Referrals"
    description="Referral management is temporarily unavailable."
  />

  <Card class="referrals-placeholder-card">
    <template #content>
      <div class="referrals-placeholder">
        <h3 class="placeholder-title">Placeholder</h3>
        <p class="placeholder-copy">
          The previous referrals page has been removed.
        </p>
      </div>
    </template>
  </Card>
</template>

<style scoped>
.referrals-placeholder-card {
  margin-top: 1rem;
}

.referrals-placeholder {
  min-height: 15rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 0.5rem;
  padding: 1.5rem;
}

.placeholder-title {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #111827;
}

.placeholder-copy {
  margin: 0;
  font-size: 0.92rem;
  color: #4b5563;
}
</style>
