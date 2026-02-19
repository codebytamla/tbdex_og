<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Drawer from 'primevue/drawer'
import Tag from 'primevue/tag'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import { useAuthStore } from '@/stores/auth'

type ReferralStatus = 'new' | 'arrived' | 'seen'
type ReferralPriority = 'routine' | 'urgent'
type BlueAnswer = 'yes' | 'no' | 'unknown'
type DiagnosticStatus = 'requested' | 'scheduled' | 'completed' | 'not_needed'

type TbBlueQuestion = {
  code: string
  prompt: string
  answer: BlueAnswer
}

type TbDiagnostic = {
  name: string
  status: DiagnosticStatus
  detail: string
}

type TbScreening = {
  method: string
  symptomScore: number
  threshold: number
  positiveSymptoms: string[]
  blueQuestions: TbBlueQuestion[]
  diagnostics: TbDiagnostic[]
}

type PatientProfile = {
  id: string
  name: string
  age: number
  sex: 'male' | 'female' | 'other'
}

type ReferralRow = {
  id: string
  patient: PatientProfile
  fromSite: string
  toSite: string
  referredAt: string
  status: ReferralStatus
  priority: ReferralPriority
  assignedDoctor: string
  referralReason: string
  tb: TbScreening
}

const authStore = useAuthStore()
authStore.hydrate()

const mySite = computed(() => authStore.sites[0] || 'Border Clinic Hub')

const referrals = computed<ReferralRow[]>(() => [
  {
    id: 'RF-22031',
    patient: {
      id: 'PT-10241',
      name: 'Amina Yusuf',
      age: 29,
      sex: 'female',
    },
    fromSite: 'Lia Border Post Unit',
    toSite: mySite.value,
    referredAt: '2026-02-19T13:15:00Z',
    status: 'new',
    priority: 'urgent',
    assignedDoctor: 'Dr. Mary Atim',
    referralReason: 'High TB symptom score, CXR requested.',
    tb: {
      method: 'Systematic TB symptom screening',
      symptomScore: 7,
      threshold: 5,
      positiveSymptoms: ['Cough > 2 weeks', 'Night sweats', 'Weight loss'],
      blueQuestions: [
        { code: 'B1', prompt: 'Any previous TB treatment?', answer: 'no' },
        { code: 'B2', prompt: 'Close contact with known TB case?', answer: 'yes' },
        { code: 'B3', prompt: 'HIV or other immunosuppression risk?', answer: 'unknown' },
      ],
      diagnostics: [
        { name: 'Chest X-ray', status: 'requested', detail: 'Priority slot required.' },
        { name: 'QT follow-up', status: 'scheduled', detail: 'Schedule after clinician review.' },
      ],
    },
  },
  {
    id: 'RF-22035',
    patient: {
      id: 'PT-10289',
      name: 'Samir Okello',
      age: 41,
      sex: 'male',
    },
    fromSite: 'Yumbe Outreach Team',
    toSite: mySite.value,
    referredAt: '2026-02-19T11:35:00Z',
    status: 'arrived',
    priority: 'routine',
    assignedDoctor: 'Dr. John Wani',
    referralReason: 'TB score at threshold; confirm with imaging.',
    tb: {
      method: 'TB risk questionnaire',
      symptomScore: 5,
      threshold: 5,
      positiveSymptoms: ['Persistent cough'],
      blueQuestions: [
        { code: 'B1', prompt: 'Any previous TB treatment?', answer: 'yes' },
        { code: 'B2', prompt: 'Close contact with known TB case?', answer: 'no' },
        { code: 'B3', prompt: 'HIV or other immunosuppression risk?', answer: 'no' },
      ],
      diagnostics: [
        { name: 'Chest X-ray', status: 'completed', detail: 'Uploaded by radiology.' },
        { name: 'QT follow-up', status: 'not_needed', detail: 'Not indicated after review.' },
      ],
    },
  },
  {
    id: 'RF-22041',
    patient: {
      id: 'PT-10308',
      name: 'Nabirye Stella',
      age: 34,
      sex: 'female',
    },
    fromSite: 'Koboko CHW Zone 3',
    toSite: mySite.value,
    referredAt: '2026-02-18T16:10:00Z',
    status: 'seen',
    priority: 'routine',
    assignedDoctor: 'Dr. Mary Atim',
    referralReason: 'Follow-up TB review after prior abnormal imaging.',
    tb: {
      method: 'Systematic TB symptom screening',
      symptomScore: 4,
      threshold: 5,
      positiveSymptoms: ['Weight loss'],
      blueQuestions: [
        { code: 'B1', prompt: 'Any previous TB treatment?', answer: 'yes' },
        { code: 'B2', prompt: 'Close contact with known TB case?', answer: 'yes' },
        { code: 'B3', prompt: 'HIV or other immunosuppression risk?', answer: 'no' },
      ],
      diagnostics: [
        { name: 'Chest X-ray', status: 'scheduled', detail: 'Follow-up in 3 days.' },
        { name: 'QT follow-up', status: 'scheduled', detail: 'Clinical check requested.' },
      ],
    },
  },
  {
    id: 'RF-22044',
    patient: {
      id: 'PT-10328',
      name: 'Hassan Dadi',
      age: 53,
      sex: 'male',
    },
    fromSite: 'Lia Border Post Unit',
    toSite: 'Arua Regional Referral',
    referredAt: '2026-02-19T09:40:00Z',
    status: 'new',
    priority: 'urgent',
    assignedDoctor: 'Dr. Peter Ladu',
    referralReason: 'Escalated for high-risk TB workup and specialist review.',
    tb: {
      method: 'Systematic TB symptom screening',
      symptomScore: 8,
      threshold: 5,
      positiveSymptoms: ['Cough > 2 weeks', 'Fever', 'Night sweats', 'Hemoptysis'],
      blueQuestions: [
        { code: 'B1', prompt: 'Any previous TB treatment?', answer: 'unknown' },
        { code: 'B2', prompt: 'Close contact with known TB case?', answer: 'yes' },
        { code: 'B3', prompt: 'HIV or other immunosuppression risk?', answer: 'yes' },
      ],
      diagnostics: [
        { name: 'Chest X-ray', status: 'requested', detail: 'Urgent request submitted.' },
        { name: 'QT follow-up', status: 'requested', detail: 'Needs same-day slot.' },
      ],
    },
  },
  {
    id: 'RF-22047',
    patient: {
      id: 'PT-10344',
      name: 'Grace Tendo',
      age: 22,
      sex: 'female',
    },
    fromSite: 'Bidi Bidi Mobile Team',
    toSite: 'Koboko HC III',
    referredAt: '2026-02-19T07:55:00Z',
    status: 'new',
    priority: 'routine',
    assignedDoctor: 'Dr. Grace Atim',
    referralReason: 'TB symptom persistence; outpatient diagnostics advised.',
    tb: {
      method: 'TB risk questionnaire',
      symptomScore: 6,
      threshold: 5,
      positiveSymptoms: ['Persistent cough', 'Weight loss'],
      blueQuestions: [
        { code: 'B1', prompt: 'Any previous TB treatment?', answer: 'no' },
        { code: 'B2', prompt: 'Close contact with known TB case?', answer: 'unknown' },
        { code: 'B3', prompt: 'HIV or other immunosuppression risk?', answer: 'unknown' },
      ],
      diagnostics: [
        { name: 'Chest X-ray', status: 'scheduled', detail: 'Pending clinic check-in.' },
        { name: 'QT follow-up', status: 'not_needed', detail: 'Hold until X-ray review.' },
      ],
    },
  },
])

const viewMode = ref<'mine' | 'network'>('mine')
const selectedSite = ref('')
const detailVisible = ref(false)
const activeReferral = ref<ReferralRow | null>(null)

const siteFilters = computed(() => {
  const countBySite = new Map<string, number>()
  for (const referral of referrals.value) {
    countBySite.set(referral.toSite, (countBySite.get(referral.toSite) ?? 0) + 1)
  }

  const ordered = [...countBySite.keys()]
  const mySiteIndex = ordered.indexOf(mySite.value)
  if (mySiteIndex > -1) {
    ordered.splice(mySiteIndex, 1)
    ordered.unshift(mySite.value)
  } else {
    ordered.unshift(mySite.value)
    countBySite.set(mySite.value, 0)
  }

  return ordered.map((site) => ({
    site,
    count: countBySite.get(site) ?? 0,
  }))
})

watch(
  siteFilters,
  (sites) => {
    if (!selectedSite.value || !sites.some((entry) => entry.site === selectedSite.value)) {
      selectedSite.value = sites[0]?.site ?? mySite.value
    }
  },
  { immediate: true },
)

const activeSite = computed(() => {
  if (viewMode.value === 'mine') return mySite.value
  return selectedSite.value || mySite.value
})

watch(activeSite, () => {
  activeReferral.value = null
  detailVisible.value = false
})

const visibleReferrals = computed(() =>
  referrals.value.filter((entry) => entry.toSite === activeSite.value),
)

function setMineMode() {
  viewMode.value = 'mine'
  selectedSite.value = mySite.value
}

function setNetworkMode() {
  viewMode.value = 'network'
  if (!selectedSite.value) selectedSite.value = mySite.value
}

function openDetails(row: ReferralRow) {
  activeReferral.value = row
  detailVisible.value = true
}

function handleRowClick(event: { data?: ReferralRow }) {
  if (!event.data) return
  openDetails(event.data)
}

const dateFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  hour: 'numeric',
  minute: '2-digit',
})

function formatDate(value: string) {
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return value
  return dateFormatter.format(parsed)
}

function sexShort(value: PatientProfile['sex']) {
  if (value === 'female') return 'F'
  if (value === 'male') return 'M'
  return 'Other'
}

function statusLabel(value: ReferralStatus) {
  if (value === 'new') return 'New'
  if (value === 'arrived') return 'Arrived'
  return 'Seen'
}

function statusSeverity(value: ReferralStatus) {
  if (value === 'new') return 'warn'
  if (value === 'arrived') return 'info'
  return 'success'
}

function prioritySeverity(value: ReferralPriority) {
  if (value === 'urgent') return 'warn'
  return 'secondary'
}

function blueAnswerSeverity(value: BlueAnswer) {
  if (value === 'yes') return 'warn'
  if (value === 'no') return 'success'
  return 'secondary'
}

function diagnosticSeverity(value: DiagnosticStatus) {
  if (value === 'requested') return 'warn'
  if (value === 'scheduled') return 'info'
  if (value === 'completed') return 'success'
  return 'secondary'
}

function diagnosticLabel(value: DiagnosticStatus) {
  if (value === 'requested') return 'Requested'
  if (value === 'scheduled') return 'Scheduled'
  if (value === 'completed') return 'Completed'
  return 'Not needed'
}

function screeningColumnTags(row: ReferralRow) {
  const tags = ['TB Symptoms', 'Blue Questions']
  const diagnosticNames = row.tb.diagnostics
    .filter((entry) => entry.status !== 'not_needed')
    .map((entry) => entry.name)
    .slice(0, 2)
  return [...tags, ...diagnosticNames]
}

const positiveSymptomText = computed(() => {
  if (!activeReferral.value) return ''
  if (activeReferral.value.tb.positiveSymptoms.length === 0) return 'No positive symptoms recorded.'
  return activeReferral.value.tb.positiveSymptoms.join(', ')
})
</script>

<template>
  <PageSectionHeader
    title="Referrals"
    description="TB referral queue for your site. Open a row for TB symptom, blue-question, and diagnostic details."
  />

  <div class="referrals-page">
    <section class="scope-row">
      <div class="mode-toggle" role="group" aria-label="Referral scope">
        <button
          type="button"
          class="mode-button"
          :class="{ 'mode-button-active': viewMode === 'mine' }"
          @click="setMineMode"
        >
          My Site
        </button>
        <button
          type="button"
          class="mode-button"
          :class="{ 'mode-button-active': viewMode === 'network' }"
          @click="setNetworkMode"
        >
          Network
        </button>
      </div>

      <p class="scope-copy">
        Showing incoming referrals for <strong>{{ activeSite }}</strong>
      </p>
    </section>

    <section v-if="viewMode === 'network'" class="site-chip-row" aria-label="Site filters">
      <button
        v-for="entry in siteFilters"
        :key="entry.site"
        type="button"
        class="site-chip"
        :class="{ 'site-chip-active': selectedSite === entry.site }"
        @click="selectedSite = entry.site"
      >
        <span>{{ entry.site }}</span>
        <span class="site-chip-count">{{ entry.count }}</span>
      </button>
    </section>

    
    <section class="queue-shell">
      <DataTable
        class="referrals-table"
        :value="visibleReferrals"
        dataKey="id"
        size="small"
        @row-click="handleRowClick"
      >
        <template #empty>
          <p class="empty-state">No incoming referrals for this site.</p>
        </template>

        <Column header="Patient" style="min-width: 14rem">
          <template #body="{ data }">
            <div class="patient-cell">
              <p class="patient-name">{{ data.patient.name }}</p>
              <p class="patient-id">{{ data.id }}</p>
            </div>
          </template>
        </Column>

        <Column header="Age / Sex" style="min-width: 6.8rem">
          <template #body="{ data }">
            <span class="small-copy">{{ data.patient.age }} / {{ sexShort(data.patient.sex) }}</span>
          </template>
        </Column>

        <Column header="From Site" style="min-width: 11rem">
          <template #body="{ data }">
            <span class="small-copy">{{ data.fromSite }}</span>
          </template>
        </Column>

        <Column header="TB Screening" style="min-width: 13rem">
          <template #body="{ data }">
            <div class="screening-chip-row">
              <Tag
                v-for="tag in screeningColumnTags(data)"
                :key="`${data.id}-${tag}`"
                :value="tag"
                severity="secondary"
              />
            </div>
          </template>
        </Column>

        <Column header="Referred" style="min-width: 8rem">
          <template #body="{ data }">
            <span class="small-copy">{{ formatDate(data.referredAt) }}</span>
          </template>
        </Column>

        <Column header="Status" style="min-width: 6.5rem">
          <template #body="{ data }">
            <Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" />
          </template>
        </Column>

        <Column header="" style="width: 5.2rem">
          <template #body="{ data }">
            <Button
              text
              size="small"
              class="open-button"
              label="Open"
              @click.stop="openDetails(data)"
            />
          </template>
        </Column>
      </DataTable>
    </section>
  </div>

  <Drawer
    v-model:visible="detailVisible"
    position="right"
    :style="{ width: 'min(48rem, 94vw)' }"
    class="referral-drawer"
    :header="activeReferral ? `${activeReferral.patient.name} (${activeReferral.id})` : 'Referral details'"
  >
    <div v-if="activeReferral" class="drawer-body">
      <p class="drawer-context">
        TB-only detail view for referral <strong>{{ activeReferral.id }}</strong> from
        <strong>{{ activeReferral.fromSite }}</strong> to
        <strong>{{ activeReferral.toSite }}</strong>.
      </p>

      <div class="row-details">
        <section class="detail-block">
          <h3>Referral Summary</h3>
          <div class="summary-grid">
            <p><span class="label">Patient</span><span>{{ activeReferral.patient.name }}</span></p>
            <p><span class="label">Age / Sex</span><span>{{ activeReferral.patient.age }} / {{ sexShort(activeReferral.patient.sex) }}</span></p>
            <p><span class="label">Assigned doctor</span><span>{{ activeReferral.assignedDoctor }}</span></p>
            <p><span class="label">Referred</span><span>{{ formatDate(activeReferral.referredAt) }}</span></p>
            <p><span class="label">Status</span><Tag :value="statusLabel(activeReferral.status)" :severity="statusSeverity(activeReferral.status)" /></p>
            <p><span class="label">Priority</span><Tag :value="activeReferral.priority === 'urgent' ? 'Urgent' : 'Routine'" :severity="prioritySeverity(activeReferral.priority)" /></p>
          </div>
          <p class="reason-copy">{{ activeReferral.referralReason }}</p>
        </section>

        <section class="detail-block">
          <h3>TB Symptom Screening</h3>
          <div class="summary-grid">
            <p><span class="label">Method</span><span>{{ activeReferral.tb.method }}</span></p>
            <p><span class="label">Score</span><span>{{ activeReferral.tb.symptomScore }} / {{ activeReferral.tb.threshold }}</span></p>
          </div>
          <p class="reason-copy">{{ positiveSymptomText }}</p>
        </section>

        <section class="detail-block">
          <h3>Blue Questions (Primary Referral)</h3>
          <ul class="question-list">
            <li
              v-for="question in activeReferral.tb.blueQuestions"
              :key="`${activeReferral.id}-${question.code}`"
            >
              <div class="question-top-row">
                <p>{{ question.code }}</p>
                <Tag :value="question.answer" :severity="blueAnswerSeverity(question.answer)" />
              </div>
              <p class="question-prompt">{{ question.prompt }}</p>
            </li>
          </ul>
        </section>

        <section class="detail-block">
          <h3>TB Diagnostics (X-ray / QT)</h3>
          <ul class="diagnostic-list">
            <li
              v-for="diagnostic in activeReferral.tb.diagnostics"
              :key="`${activeReferral.id}-${diagnostic.name}`"
            >
              <div class="question-top-row">
                <p>{{ diagnostic.name }}</p>
                <Tag
                  :value="diagnosticLabel(diagnostic.status)"
                  :severity="diagnosticSeverity(diagnostic.status)"
                />
              </div>
              <p class="question-prompt">{{ diagnostic.detail }}</p>
            </li>
          </ul>
        </section>
      </div>
    </div>

    <template #footer>
      <div class="drawer-actions">
        <Button label="Open Full Record" text disabled />
        <Button label="Mark Arrived" disabled />
      </div>
    </template>
  </Drawer>
</template>

<style scoped>
.referrals-page {
  display: grid;
  gap: 0.75rem;
}

.scope-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.7rem;
  flex-wrap: wrap;
}

.mode-toggle {
  display: inline-flex;
  align-items: center;
  border: 1px solid var(--p-surface-300);
  border-radius: 999px;
  padding: 0.16rem;
  background: var(--p-surface-0);
}

.mode-button {
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--p-surface-700);
  font-size: 0.78rem;
  font-weight: 620;
  padding: 0.36rem 0.72rem;
  cursor: pointer;
}

.mode-button-active {
  background: var(--p-surface-900);
  color: var(--p-surface-0);
}

.scope-copy {
  margin: 0;
  font-size: 0.82rem;
  color: var(--p-surface-600);
}

.site-chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}

.site-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  border: 1px solid var(--p-surface-300);
  border-radius: 999px;
  background: var(--p-surface-0);
  color: var(--p-surface-700);
  font-size: 0.76rem;
  font-weight: 560;
  padding: 0.3rem 0.58rem;
  cursor: pointer;
}

.site-chip-active {
  border-color: #0369a1;
  background: color-mix(in srgb, #0284c7 9%, var(--p-surface-0));
  color: #0c4a6e;
}

.site-chip-count {
  border-radius: 999px;
  padding: 0.06rem 0.34rem;
  font-size: 0.7rem;
  background: var(--p-surface-100);
}

.hint-copy {
  margin: 0;
  font-size: 0.76rem;
  color: var(--p-surface-500);
}

.queue-shell {
  min-width: 0;
  border: 1px solid var(--p-surface-200);
  border-radius: 0.86rem;
  overflow: hidden;
  background: var(--p-surface-0);
}

.empty-state {
  margin: 0.2rem 0;
  font-size: 0.82rem;
  color: var(--p-surface-600);
}

.patient-cell {
  display: grid;
  gap: 0.16rem;
}

.patient-name {
  margin: 0;
  font-size: 0.84rem;
  font-weight: 650;
  color: var(--p-surface-900);
}

.patient-id {
  margin: 0;
  font-size: 0.72rem;
  color: var(--p-surface-500);
}

.small-copy {
  font-size: 0.78rem;
  color: var(--p-surface-700);
}

.screening-chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.open-button {
  padding-inline: 0.3rem;
}

:deep(.referrals-table .p-datatable-tbody > tr) {
  cursor: pointer;
}

:deep(.referrals-table .p-datatable-thead > tr > th) {
  font-size: 0.74rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--p-surface-500);
}

:deep(.referrals-table .p-datatable-tbody > tr > td) {
  padding-top: 0.52rem;
  padding-bottom: 0.52rem;
  vertical-align: middle;
}

.drawer-body {
  display: grid;
  gap: 0.64rem;
}

.drawer-context {
  margin: 0;
  font-size: 0.79rem;
  line-height: 1.42;
  color: var(--p-surface-700);
}

.drawer-context strong {
  color: var(--p-surface-900);
}

.row-details {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.6rem;
}

.detail-block {
  border: 1px solid var(--p-surface-200);
  border-radius: 0.72rem;
  background: var(--p-surface-0);
  padding: 0.62rem 0.68rem;
}

.detail-block h3 {
  margin: 0;
  font-size: 0.81rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--p-surface-600);
}

.summary-grid {
  margin-top: 0.52rem;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
}

.summary-grid p {
  margin: 0;
  display: grid;
  gap: 0.12rem;
  font-size: 0.8rem;
  color: var(--p-surface-900);
}

.label {
  font-size: 0.72rem;
  color: var(--p-surface-500);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.reason-copy {
  margin: 0.55rem 0 0;
  font-size: 0.8rem;
  line-height: 1.38;
  color: var(--p-surface-700);
}

.question-list,
.diagnostic-list {
  margin: 0.48rem 0 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 0.45rem;
}

.question-list li,
.diagnostic-list li {
  border: 1px solid var(--p-surface-200);
  border-radius: 0.62rem;
  background: var(--p-surface-0);
  padding: 0.45rem 0.52rem;
}

.question-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.question-top-row p {
  margin: 0;
  font-size: 0.81rem;
  font-weight: 650;
  color: var(--p-surface-900);
}

.question-prompt {
  margin: 0.3rem 0 0;
  font-size: 0.77rem;
  color: var(--p-surface-600);
}

.drawer-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.45rem;
}

:deep(.p-drawer-mask) {
  background: rgb(2 6 23 / 72%) !important;
}

:deep(.referral-drawer.p-drawer) {
  margin: 0;
  border-radius: 0;
  height: 100vh;
  max-height: 100vh;
}

:deep(.referral-drawer .p-drawer-header) {
  border-bottom: 1px solid var(--p-surface-200);
  padding: 0.72rem 0.82rem;
}

:deep(.referral-drawer .p-drawer-content) {
  overflow: auto;
  padding: 0.72rem 0.82rem;
}

:deep(.referral-drawer .p-drawer-footer) {
  border-top: 1px solid var(--p-surface-200);
  padding: 0.6rem 0.82rem 0.72rem;
}

@media (max-width: 1180px) {
  .row-details {
    grid-template-columns: 1fr;
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 740px) {
  .scope-row {
    align-items: flex-start;
  }

  .mode-toggle {
    width: 100%;
    justify-content: space-between;
  }

  .mode-button {
    flex: 1 1 0;
  }
}
</style>
