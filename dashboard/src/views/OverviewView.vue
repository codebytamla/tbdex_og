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
        <UserPlus :size="14" :stroke-width="2.25" />
        <span>Register patient</span>
      </Button>
      <Button class="top-cta" severity="secondary" @click="openSettings">
        <Shield :size="14" :stroke-width="2.25" />
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
            <MapPin :size="14" :stroke-width="2.2" />
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
            <CheckCircle2 :size="14" :stroke-width="2.2" />
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
  gap: 0.9rem;
}

.top-cta {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  height: 2.2rem;
  border-radius: 0.68rem;
  padding: 0 0.92rem;
  font-size: 0.79rem;
  font-weight: 620;
}

.top-cta-primary {
  background: #0f172a;
  border-color: #0f172a;
  color: #ffffff;
}

.top-cta-primary:hover:not(:disabled) {
  background: #020617;
  border-color: #020617;
  color: #ffffff;
}

.hero-card {
  border-radius: 1rem;
  border: 1px solid color-mix(in srgb, #0891b2 28%, var(--p-surface-200));
  background:
    radial-gradient(circle at 10% 12%, rgb(6 182 212 / 18%) 0%, transparent 48%),
    radial-gradient(circle at 88% 0%, rgb(34 197 94 / 12%) 0%, transparent 56%),
    linear-gradient(146deg, var(--p-surface-0), color-mix(in srgb, #0ea5a4 6%, var(--p-surface-0)));
}

.hero-card :deep(.p-card-content) {
  display: grid;
  gap: 0.72rem;
}

.hero-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(0, 1fr);
  gap: 0.72rem;
}

.hero-eyebrow {
  margin: 0;
  font-size: 0.66rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 700;
  color: #0e7490;
}

.hero-copy h2 {
  margin: 0.26rem 0 0;
  font-size: 1.2rem;
  line-height: 1.16;
  letter-spacing: -0.012em;
  color: var(--p-surface-900);
}

.hero-copy p {
  margin: 0.48rem 0 0;
  font-size: 0.81rem;
  line-height: 1.42;
  color: var(--p-surface-700);
  max-width: 52ch;
}

.hero-note {
  display: inline-flex;
  width: fit-content;
  padding: 0.4rem 0.52rem;
  border-radius: 0.6rem;
  border: 1px solid var(--p-surface-200);
  background: color-mix(in srgb, #0891b2 8%, var(--p-surface-0));
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
}

.stat-card {
  border: 1px solid var(--p-surface-200);
  border-radius: 0.72rem;
  background: color-mix(in srgb, var(--p-surface-0) 92%, #0ea5a4 8%);
  padding: 0.58rem 0.68rem;
}

.stat-value {
  margin: 0;
  font-size: 1.03rem;
  font-weight: 760;
  color: var(--p-surface-900);
}

.stat-label {
  margin: 0.2rem 0 0;
  font-size: 0.72rem;
  color: var(--p-surface-700);
}

.context-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.72rem;
}

.context-card {
  border-radius: 0.9rem;
}

.context-label {
  margin: 0;
  font-size: 0.73rem;
  color: var(--p-surface-500);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 650;
}

.context-value {
  margin: 0.22rem 0 0.62rem;
  font-size: 0.88rem;
  color: var(--p-surface-900);
  font-weight: 620;
}

.context-inline {
  display: inline-flex;
  align-items: center;
  gap: 0.32rem;
}

.queue-health {
  display: inline-flex;
  align-items: center;
  gap: 0.32rem;
  color: #0f766e;
}

.reveal {
  opacity: 0;
  transform: translateY(10px);
  animation: rise-in 390ms ease forwards;
}

.reveal.delay-1 {
  animation-delay: 20ms;
}

.reveal.delay-2 {
  animation-delay: 95ms;
}

@keyframes rise-in {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.p-card) {
  border: 1px solid var(--p-surface-200);
  box-shadow: 0 14px 30px -28px color-mix(in srgb, #0ea5a4 30%, transparent);
}

:deep(.p-card .p-card-content) {
  padding: 0.92rem;
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
    font-size: 1.08rem;
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
