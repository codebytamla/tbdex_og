<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { CheckCircle2, MapPin, Shield, UserPlus } from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
authStore.hydrate()

const router = useRouter()

const primarySite = computed(() => authStore.sites[0] || 'No assigned site yet')
const siteCount = computed(() => authStore.sites.length)
const hasAssignedSites = computed(() => siteCount.value > 0)

const userDisplayName = computed(() => {
  const profile = authStore.profile
  if (!profile || typeof profile !== 'object') return 'Current user'
  const record = profile as Record<string, unknown>
  for (const key of ['full_name', 'fullName', 'name', 'email']) {
    const candidate = record[key]
    if (typeof candidate === 'string' && candidate.trim()) return candidate.trim()
  }
  return 'Current user'
})

const accessRoleLabel = computed(() => {
  const role = authStore.accountRole?.trim()
  if (role) return role
  if (authStore.roles.length > 0) return authStore.roles.join(', ')
  return 'No role assigned'
})

function openPatients() {
  void router.push('/patients')
}

function openSettings() {
  void router.push('/settings')
}
</script>

<template>
  <PageSectionHeader title="Overview" description="Dashboard summary for account access, sites, and quick actions">
    <template #actions>
      <Button class="top-cta top-cta-primary" @click="openPatients">
        <UserPlus :size="16" :stroke-width="2.25" />
        <span>Register patient</span>
      </Button>
      <Button class="top-cta" severity="secondary" @click="openSettings">
        <Shield :size="16" :stroke-width="2.25" />
        <span>Account settings</span>
      </Button>
    </template>
  </PageSectionHeader>

  <div class="overview-shell">
    <Card class="hero-card reveal delay-1">
      <template #content>
        <div class="hero-layout">
          <div class="hero-copy">
            <p class="hero-eyebrow">Workspace Snapshot</p>
            <h2>Welcome back, {{ userDisplayName }}</h2>
            <p>Use this overview to confirm your account context and move into patient management quickly.</p>
            <p class="hero-note">This screen shows local UI state and account profile data.</p>
          </div>

          <div class="hero-stats">
            <article class="stat-card">
              <p class="stat-value">{{ siteCount }}</p>
              <p class="stat-label">Assigned sites</p>
            </article>
            <article class="stat-card">
              <p class="stat-value">{{ hasAssignedSites ? 'Active' : 'Pending' }}</p>
              <p class="stat-label">Site assignment</p>
            </article>
          </div>
        </div>
      </template>
    </Card>

    <div class="context-grid reveal delay-2">
      <Card class="context-card">
        <template #content>
          <p class="context-label">Primary site</p>
          <p class="context-value context-inline">
            <MapPin :size="16" :stroke-width="2.2" />
            <span>{{ primarySite }}</span>
          </p>

          <p class="context-label">Sites linked</p>
          <p class="context-value">{{ siteCount }}</p>
        </template>
      </Card>

      <Card class="context-card">
        <template #content>
          <p class="context-label">Access role</p>
          <p class="context-value">{{ accessRoleLabel }}</p>

          <p class="context-label">Session status</p>
          <p class="context-value queue-health">
            <CheckCircle2 :size="16" :stroke-width="2.2" />
            <span>Authenticated</span>
          </p>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.overview-shell {
  display: grid;
  gap: 1.25rem;
}

.top-cta {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  height: var(--app-control-height, 2.25rem);
  border-radius: var(--app-control-radius, 8px);
  padding: 0 1rem;
  font-size: 0.855rem;
  font-weight: 600;
  letter-spacing: -0.005em;
  transition: all 0.15s ease;
}

.top-cta-primary {
  background: #0f172a;
  border-color: #0f172a;
  color: #ffffff;
  box-shadow: 0 1px 3px 0 rgb(15 23 42 / 0.2);
}

.top-cta-primary:hover:not(:disabled) {
  background: #1e293b;
  border-color: #1e293b;
  color: #ffffff;
  box-shadow: 0 4px 12px -2px rgb(15 23 42 / 0.25);
}

.hero-card {
  border-radius: var(--radius-xl, 18px) !important;
  border: 1px solid #e2e8f0;
  background:
    radial-gradient(ellipse at 0% 0%, rgb(13 148 136 / 0.06) 0%, transparent 50%),
    radial-gradient(ellipse at 100% 100%, rgb(45 212 191 / 0.04) 0%, transparent 50%),
    #ffffff;
  overflow: hidden;
}

.hero-card :deep(.p-card-content) {
  display: grid;
  gap: 0;
}

.hero-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.5fr) minmax(0, 1fr);
  gap: 2rem;
}

.hero-eyebrow {
  margin: 0;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-weight: 700;
  color: #0d9488;
}

.hero-copy h2 {
  margin: 0.5rem 0 0;
  font-size: 1.375rem;
  line-height: 1.2;
  letter-spacing: -0.02em;
  font-weight: 800;
  color: #111827;
}

.hero-copy p {
  margin: 0.5rem 0 0;
  font-size: 0.875rem;
  line-height: 1.5;
  color: #6b7280;
  max-width: 48ch;
}

.hero-note {
  display: inline-flex;
  width: fit-content;
  padding: 0.375rem 0.625rem;
  border-radius: var(--radius-sm, 6px);
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  font-size: 0.825rem;
  color: #6b7280;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.75rem;
  align-content: center;
}

.stat-card {
  border: 1px solid #e5e7eb;
  border-radius: var(--radius-lg, 14px);
  background: #f9fafb;
  padding: 1rem 1.125rem;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.stat-card:hover {
  border-color: #d1d5db;
  box-shadow: 0 2px 8px -2px rgb(0 0 0 / 0.06);
}

.stat-value {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: #111827;
}

.stat-label {
  margin: 0.25rem 0 0;
  font-size: 0.78rem;
  font-weight: 500;
  color: #9ca3af;
  letter-spacing: 0.01em;
}

.context-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.context-card {
  border-radius: var(--radius-xl, 18px) !important;
}

.context-label {
  margin: 0;
  font-size: 0.7rem;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 700;
}

.context-value {
  margin: 0.25rem 0 0.75rem;
  font-size: 0.935rem;
  color: #111827;
  font-weight: 700;
  letter-spacing: -0.005em;
}

.context-inline {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
}

.queue-health {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  color: #0d9488;
}

.reveal {
  opacity: 0;
  transform: translateY(8px);
  animation: rise-in 450ms cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.reveal.delay-1 {
  animation-delay: 40ms;
}

.reveal.delay-2 {
  animation-delay: 120ms;
}

@keyframes rise-in {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.p-card) {
  border: 1px solid #e5e7eb;
  box-shadow: var(--shadow-sm, 0 1px 3px 0 rgb(0 0 0 / 0.06));
}

:deep(.p-card .p-card-content) {
  padding: 1.25rem 1.5rem;
}

@media (max-width: 1180px) {
  .hero-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .top-cta {
    width: 100%;
    justify-content: center;
  }

  .hero-copy h2 {
    font-size: 1.2rem;
  }

  .hero-stats {
    grid-template-columns: 1fr;
  }

  .context-grid {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal {
    animation: none;
    opacity: 1;
    transform: none;
  }
}
</style>
