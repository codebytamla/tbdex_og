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
:deep(.p-card) {
  border: 1px solid #e5e7eb;
  box-shadow: var(--shadow-sm, 0 1px 3px 0 rgb(0 0 0 / 0.06));
  border-radius: var(--radius-xl, 18px) !important;
  overflow: hidden;
}

:deep(.p-card .p-card-content) {
  padding: 1.5rem;
}

.account-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
}

.account-item {
  border: 1px solid #e5e7eb;
  border-radius: var(--radius-lg, 14px);
  padding: 1rem 1.125rem;
  background: #f9fafb;
  transition: border-color 0.15s ease;
}

.account-item:hover {
  border-color: #d1d5db;
}

.account-item h3 {
  margin: 0;
  font-size: 0.7rem;
  color: #9ca3af;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.account-item p {
  margin: 0.375rem 0 0;
  color: #111827;
  font-weight: 700;
  font-size: 1rem;
  letter-spacing: -0.01em;
  overflow-wrap: anywhere;
}

.claim-grid {
  margin-top: 1rem;
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.claim-item {
  border: 1px solid #e5e7eb;
  border-radius: var(--radius-lg, 14px);
  padding: 1rem 1.125rem;
  background: #ffffff;
}

.claim-item h3 {
  margin: 0 0 0.625rem;
  font-size: 0.78rem;
  color: #6b7280;
  font-weight: 700;
  letter-spacing: -0.005em;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  border-radius: var(--radius-sm, 6px);
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  color: #374151;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  letter-spacing: -0.005em;
}

.chip.muted {
  color: #9ca3af;
  background: transparent;
  border-style: dashed;
}

@media (max-width: 900px) {
  .account-grid {
    grid-template-columns: 1fr;
  }
}
</style>
