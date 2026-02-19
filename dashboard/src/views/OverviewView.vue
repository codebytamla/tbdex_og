<script setup lang="ts">
import { computed, ref, type Component } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowRightLeft,
  CheckCircle2,
  ClipboardList,
  Clock3,
  FileText,
  RefreshCcw,
  Stethoscope,
  UserPlus,
} from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import { useAuthStore } from '@/stores/auth'

type EncounterType = 'all' | 'referral' | 'assessment' | 'clinic' | 'notes'
type QueueAction = 'accept' | 'assess' | 'approve' | 'review'
type QueueUrgency = 'urgent' | 'today' | 'routine'

type EncounterRow = {
  id: string
  patientName: string
  mrn: string
  type: Exclude<EncounterType, 'all'>
  action: QueueAction
  urgency: QueueUrgency
  site: string
  summary: string
  updatedAt: string
}

const authStore = useAuthStore()
authStore.hydrate()

const router = useRouter()

const selectedType = ref<EncounterType>('all')

const encounterRows = ref<EncounterRow[]>([
  {
    id: 'enc-001',
    patientName: 'Amina Ismail',
    mrn: 'MRN-20314',
    type: 'referral',
    action: 'accept',
    urgency: 'urgent',
    site: 'Kisumu Main Clinic',
    summary: 'Incoming referral from Nyeri East. Needs intake acceptance before assignment.',
    updatedAt: '8m ago',
  },
  {
    id: 'enc-002',
    patientName: 'David Otieno',
    mrn: 'MRN-17708',
    type: 'assessment',
    action: 'assess',
    urgency: 'urgent',
    site: 'Kisumu Main Clinic',
    summary: 'Assessment requested after triage. Vitals captured, clinician review pending.',
    updatedAt: '14m ago',
  },
  {
    id: 'enc-003',
    patientName: 'Grace Wanjiru',
    mrn: 'MRN-19021',
    type: 'referral',
    action: 'approve',
    urgency: 'today',
    site: 'County Referral Unit',
    summary: 'Referral draft completed and waiting for final approval.',
    updatedAt: '34m ago',
  },
  {
    id: 'enc-004',
    patientName: 'Peter Mboya',
    mrn: 'MRN-22190',
    type: 'clinic',
    action: 'assess',
    urgency: 'today',
    site: 'Mumias Satellite',
    summary: 'Clinic follow-up encounter is open with incomplete assessment notes.',
    updatedAt: '1h ago',
  },
  {
    id: 'enc-005',
    patientName: 'Mercy Akinyi',
    mrn: 'MRN-22441',
    type: 'notes',
    action: 'review',
    urgency: 'routine',
    site: 'Kisumu Main Clinic',
    summary: 'Notes-only encounter created by nursing team. Needs physician sign-off.',
    updatedAt: '2h ago',
  },
  {
    id: 'enc-006',
    patientName: 'Samuel Kariuki',
    mrn: 'MRN-21675',
    type: 'assessment',
    action: 'assess',
    urgency: 'routine',
    site: 'Outreach Camp A',
    summary: 'Field assessment received from mobile team and queued for review.',
    updatedAt: '3h ago',
  },
])

const primarySite = computed(() => authStore.sites[0] || 'No assigned site yet')
const siteCount = computed(() => authStore.sites.length)

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

const typeCounts = computed(() => {
  const counts = { referral: 0, assessment: 0, clinic: 0, notes: 0 }
  for (const row of encounterRows.value) {
    counts[row.type] += 1
  }
  return counts
})

const actionCounts = computed(() => {
  const counts = { accept: 0, assess: 0, approve: 0, review: 0 }
  for (const row of encounterRows.value) {
    counts[row.action] += 1
  }
  return counts
})

const urgentCount = computed(() => encounterRows.value.filter((row) => row.urgency === 'urgent').length)

const typeOptions = computed(() => [
  { value: 'all' as const, label: 'All', count: encounterRows.value.length },
  { value: 'referral' as const, label: 'Referral', count: typeCounts.value.referral },
  { value: 'assessment' as const, label: 'Assessment', count: typeCounts.value.assessment },
  { value: 'clinic' as const, label: 'Clinic', count: typeCounts.value.clinic },
  { value: 'notes' as const, label: 'Notes-only', count: typeCounts.value.notes },
])

const visibleRows = computed(() => {
  const source =
    selectedType.value === 'all'
      ? encounterRows.value
      : encounterRows.value.filter((row) => row.type === selectedType.value)

  const urgencyRank: Record<QueueUrgency, number> = { urgent: 0, today: 1, routine: 2 }
  return [...source].sort((a, b) => urgencyRank[a.urgency] - urgencyRank[b.urgency])
})

function encounterTypeLabel(type: Exclude<EncounterType, 'all'>) {
  if (type === 'referral') return 'Referral'
  if (type === 'assessment') return 'Assessment'
  if (type === 'clinic') return 'Clinic'
  return 'Notes-only'
}

function encounterActionLabel(action: QueueAction) {
  if (action === 'accept') return 'To accept'
  if (action === 'assess') return 'To assess'
  if (action === 'approve') return 'To approve'
  return 'To review'
}

function primaryActionLabel(action: QueueAction) {
  if (action === 'accept') return 'Accept'
  if (action === 'assess') return 'Start'
  if (action === 'approve') return 'Approve'
  return 'Review'
}

function urgencyLabel(urgency: QueueUrgency) {
  if (urgency === 'urgent') return 'Urgent'
  if (urgency === 'today') return 'Today'
  return 'Routine'
}

function urgencyClass(urgency: QueueUrgency) {
  if (urgency === 'urgent') return 'urgency-urgent'
  if (urgency === 'today') return 'urgency-today'
  return 'urgency-routine'
}

function typeIcon(type: Exclude<EncounterType, 'all'>): Component {
  if (type === 'referral') return ArrowRightLeft
  if (type === 'assessment') return ClipboardList
  if (type === 'clinic') return Stethoscope
  return FileText
}

function setType(next: EncounterType) {
  selectedType.value = next
}

function openPatients() {
  void router.push('/patients')
}

function previewOnly() {}
</script>

<template>
  <PageSectionHeader
    title="Dashboard"
    description="Action inbox first: urgent referral, assessment, clinic, and notes-only work"
  >
    <template #actions>
      <Button class="top-cta top-cta-primary" @click="previewOnly">
        <ArrowRightLeft :size="14" :stroke-width="2.25" />
        <span>New encounter</span>
      </Button>
      <Button class="top-cta" severity="secondary" @click="openPatients">
        <UserPlus :size="14" :stroke-width="2.25" />
        <span>Register patient</span>
      </Button>
    </template>
  </PageSectionHeader>

  <div class="overview-shell">
    <Card class="hero-card reveal delay-1">
      <template #content>
        <div class="hero-layout">
          <div class="hero-copy">
            <p class="hero-eyebrow">Action Inbox</p>
            <h2>Put clinical work first when users land on dashboard</h2>
            <p>
              This prototype prioritizes queue-driven operations, so teams immediately see what to
              accept, assess, approve, or review.
            </p>
            <p class="hero-note">UI prototype only. No backend wiring in this version.</p>
          </div>
          <div class="hero-stats">
            <article class="stat-card">
              <p class="stat-value">{{ actionCounts.accept }}</p>
              <p class="stat-label">To accept</p>
            </article>
            <article class="stat-card">
              <p class="stat-value">{{ actionCounts.assess }}</p>
              <p class="stat-label">To assess</p>
            </article>
            <article class="stat-card">
              <p class="stat-value">{{ actionCounts.approve }}</p>
              <p class="stat-label">To approve</p>
            </article>
            <article class="stat-card">
              <p class="stat-value">{{ typeCounts.notes }}</p>
              <p class="stat-label">Notes-only</p>
            </article>
          </div>
        </div>
      </template>
    </Card>

    <Card class="queue-card reveal delay-2">
      <template #content>
        <div class="queue-head">
          <div>
            <h3>Needs attention now</h3>
            <p>Urgent items first, then today, then routine.</p>
          </div>
          <div class="queue-head-right">
            <span class="urgent-indicator">
              <Clock3 :size="14" :stroke-width="2.25" />
              <span>{{ urgentCount }} urgent</span>
            </span>
            <Button text class="refresh-btn" @click="previewOnly">
              <RefreshCcw :size="14" :stroke-width="2.25" />
            </Button>
          </div>
        </div>

        <div class="type-tabs">
          <button
            v-for="option in typeOptions"
            :key="option.value"
            class="type-tab"
            :class="{ active: selectedType === option.value }"
            @click="setType(option.value)"
          >
            <span>{{ option.label }}</span>
            <span class="type-tab-count">{{ option.count }}</span>
          </button>
        </div>

        <div class="queue-list">
          <article v-for="row in visibleRows" :key="row.id" class="queue-row">
            <div class="queue-main">
              <p class="queue-type">
                <component :is="typeIcon(row.type)" :size="13" :stroke-width="2.15" />
                <span>{{ encounterTypeLabel(row.type) }}</span>
              </p>
              <h4>
                {{ row.patientName }}
                <span class="queue-mrn">{{ row.mrn }}</span>
              </h4>
              <p class="queue-summary">{{ row.summary }}</p>
            </div>

            <div class="queue-meta">
              <span class="action-pill">{{ encounterActionLabel(row.action) }}</span>
              <span class="urgency-pill" :class="urgencyClass(row.urgency)">{{ urgencyLabel(row.urgency) }}</span>
              <p>{{ row.site }}</p>
              <p>Updated {{ row.updatedAt }}</p>
            </div>

            <div class="queue-actions">
              <Button class="row-btn row-btn-primary" size="small" :label="primaryActionLabel(row.action)" @click="previewOnly" />
              <Button class="row-btn" text size="small" label="Open" @click="previewOnly" />
            </div>
          </article>

          <p v-if="visibleRows.length === 0" class="queue-empty">
            No encounters in this filter.
          </p>
        </div>
      </template>
    </Card>

    <div class="context-grid reveal delay-3">
      <Card class="context-card">
        <template #content>
          <p class="context-label">Signed in as</p>
          <p class="context-value">{{ userDisplayName }}</p>
          <p class="context-label">Primary site</p>
          <p class="context-value">{{ primarySite }}</p>
        </template>
      </Card>

      <Card class="context-card">
        <template #content>
          <p class="context-label">Assigned sites</p>
          <p class="context-value">{{ siteCount }}</p>
          <p class="context-label">Queue health</p>
          <p class="context-value queue-health">
            <CheckCircle2 :size="14" :stroke-width="2.2" />
            <span>Prototype view active</span>
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

.queue-card {
  border-radius: 1rem;
}

.queue-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.7rem;
}

.queue-head h3 {
  margin: 0;
  font-size: 1.02rem;
  letter-spacing: -0.012em;
  color: var(--p-surface-900);
}

.queue-head p {
  margin: 0.3rem 0 0;
  font-size: 0.79rem;
  color: var(--p-surface-600);
}

.queue-head-right {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

.urgent-indicator {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  border-radius: 999px;
  border: 1px solid #fecaca;
  background: #fff1f2;
  color: #be123c;
  padding: 0.2rem 0.46rem;
  font-size: 0.7rem;
  font-weight: 650;
  line-height: 1.3;
}

.refresh-btn {
  width: 2rem;
  height: 2rem;
  border-radius: 0.62rem;
  color: var(--p-surface-700);
}

.type-tabs {
  margin-top: 0.62rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.38rem;
}

.type-tab {
  display: inline-flex;
  align-items: center;
  gap: 0.38rem;
  border: 1px solid var(--p-surface-200);
  border-radius: 999px;
  background: var(--p-surface-50);
  color: var(--p-surface-700);
  font-size: 0.74rem;
  font-weight: 620;
  padding: 0.28rem 0.56rem;
  cursor: pointer;
  transition: all 0.15s;
}

.type-tab:hover {
  border-color: var(--p-surface-400);
  background: var(--p-surface-100);
}

.type-tab.active {
  border-color: #0891b2;
  background: color-mix(in srgb, #0891b2 14%, var(--p-surface-0));
  color: #0f172a;
}

.type-tab-count {
  min-width: 1.15rem;
  text-align: center;
  border-radius: 0.5rem;
  background: color-mix(in srgb, var(--p-surface-0) 84%, var(--p-surface-900) 16%);
  color: var(--p-surface-900);
  font-size: 0.66rem;
  padding: 0.1rem 0.28rem;
}

.queue-list {
  margin-top: 0.66rem;
  display: grid;
  gap: 0.52rem;
}

.queue-row {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) auto auto;
  gap: 0.62rem;
  align-items: center;
  border: 1px solid var(--p-surface-200);
  border-radius: 0.82rem;
  background: var(--p-surface-0);
  padding: 0.6rem 0.68rem;
}

.queue-main {
  min-width: 0;
}

.queue-type {
  margin: 0;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.69rem;
  font-weight: 670;
  text-transform: uppercase;
  letter-spacing: 0.055em;
  color: #0f766e;
}

.queue-main h4 {
  margin: 0.22rem 0 0;
  font-size: 0.91rem;
  color: var(--p-surface-900);
}

.queue-mrn {
  margin-left: 0.26rem;
  font-size: 0.74rem;
  color: var(--p-surface-500);
  font-weight: 500;
}

.queue-summary {
  margin: 0.2rem 0 0;
  font-size: 0.77rem;
  line-height: 1.4;
  color: var(--p-surface-600);
  max-width: 62ch;
}

.queue-meta {
  display: grid;
  justify-items: start;
  gap: 0.2rem;
}

.action-pill {
  border-radius: 999px;
  border: 1px solid color-mix(in srgb, #0ea5a4 32%, var(--p-surface-200));
  background: color-mix(in srgb, #0ea5a4 12%, var(--p-surface-0));
  color: #0f766e;
  font-size: 0.68rem;
  font-weight: 660;
  padding: 0.14rem 0.45rem;
}

.urgency-pill {
  border-radius: 999px;
  font-size: 0.66rem;
  font-weight: 700;
  padding: 0.12rem 0.43rem;
  border: 1px solid transparent;
}

.urgency-urgent {
  color: #be123c;
  background: #fff1f2;
  border-color: #fecaca;
}

.urgency-today {
  color: #b45309;
  background: #fffbeb;
  border-color: #fde68a;
}

.urgency-routine {
  color: #1d4ed8;
  background: #eff6ff;
  border-color: #bfdbfe;
}

.queue-meta p {
  margin: 0;
  color: var(--p-surface-600);
  font-size: 0.71rem;
  line-height: 1.3;
  white-space: nowrap;
}

.queue-actions {
  display: inline-flex;
  gap: 0.34rem;
}

.row-btn {
  height: 1.92rem;
  border-radius: 0.58rem;
  font-size: 0.72rem;
  font-weight: 650;
}

.row-btn-primary {
  background: #0f172a;
  border-color: #0f172a;
  color: #ffffff;
}

.row-btn-primary:hover:not(:disabled) {
  background: #020617;
  border-color: #020617;
  color: #ffffff;
}

.queue-empty {
  margin: 0;
  border: 1px dashed var(--p-surface-300);
  border-radius: 0.72rem;
  padding: 0.75rem;
  font-size: 0.8rem;
  color: var(--p-surface-600);
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

.reveal.delay-3 {
  animation-delay: 170ms;
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

  .queue-row {
    grid-template-columns: 1fr;
    align-items: start;
  }

  .queue-meta {
    grid-template-columns: repeat(2, minmax(0, auto));
    justify-content: start;
    align-items: center;
    column-gap: 0.4rem;
  }

  .queue-actions {
    width: 100%;
    justify-content: flex-start;
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
