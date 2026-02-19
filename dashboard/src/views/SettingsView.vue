<script setup lang="ts">
import { computed } from 'vue'
import Card from 'primevue/card'
import { useAuthStore } from '@/stores/auth'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'

const authStore = useAuthStore()
authStore.hydrate()

const accountName = computed(() => authStore.accountName)
const accountRole = computed(() => authStore.accountRole)
const accountEmail = computed(() => authStore.profile?.email || 'Not provided')
const roles = computed(() => authStore.roles)
const permissions = computed(() => authStore.permissions)
const sites = computed(() => authStore.sites)
</script>

<template>
  <PageSectionHeader title="Settings" description="Account and application preferences" />

  <Card>
    <template #content>
      <section class="account-grid">
        <article class="account-item">
          <h3>Account name</h3>
          <p>{{ accountName }}</p>
        </article>
        <article class="account-item">
          <h3>Role</h3>
          <p>{{ accountRole }}</p>
        </article>
        <article class="account-item">
          <h3>Email</h3>
          <p>{{ accountEmail }}</p>
        </article>
      </section>

      <section class="claim-grid">
        <article class="claim-item">
          <h3>Roles</h3>
          <div class="chips">
            <span v-for="role in roles" :key="role" class="chip">{{ role }}</span>
            <span v-if="!roles.length" class="chip muted">No role claims</span>
          </div>
        </article>
        <article class="claim-item">
          <h3>Permissions</h3>
          <div class="chips">
            <span v-for="permission in permissions" :key="permission" class="chip">{{ permission }}</span>
            <span v-if="!permissions.length" class="chip muted">No permission claims</span>
          </div>
        </article>
        <article class="claim-item">
          <h3>Sites</h3>
          <div class="chips">
            <span v-for="site in sites" :key="site" class="chip">{{ site }}</span>
            <span v-if="!sites.length" class="chip muted">No site assignments</span>
          </div>
        </article>
      </section>
    </template>
  </Card>
</template>

<style scoped>
.account-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
}

.account-item {
  border: 1px solid var(--p-surface-200);
  border-radius: 0.8rem;
  padding: 0.85rem;
  background: var(--p-surface-50);
}

.account-item h3 {
  margin: 0;
  font-size: 0.8rem;
  color: var(--p-surface-600);
  font-weight: 600;
}

.account-item p {
  margin: 0.5rem 0 0;
  color: var(--p-surface-900);
  font-weight: 700;
  font-size: 1rem;
  overflow-wrap: anywhere;
}

.claim-grid {
  margin-top: 0.75rem;
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.claim-item {
  border: 1px solid var(--p-surface-200);
  border-radius: 0.8rem;
  padding: 0.85rem;
  background: var(--p-surface-0);
}

.claim-item h3 {
  margin: 0 0 0.5rem;
  font-size: 0.82rem;
  color: var(--p-surface-600);
  font-weight: 600;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  border: 1px solid var(--p-surface-200);
  background: var(--p-surface-100);
  color: var(--p-surface-800);
  font-size: 0.73rem;
  font-weight: 600;
  padding: 0.24rem 0.58rem;
}

.chip.muted {
  color: var(--p-surface-600);
}

:deep(.p-card) {
  border: 1px solid var(--p-surface-200);
  box-shadow: none;
}

@media (max-width: 900px) {
  .account-grid {
    grid-template-columns: 1fr;
  }
}
</style>
