<script setup lang="ts">
import { computed, ref } from 'vue'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'

type EncounterTabKey = 'all' | 'referral' | 'assessment' | 'clinic_notes'
type EncounterType = Exclude<EncounterTabKey, 'all'>

type EncounterRow = {
  id: string
  patient: string
  type: EncounterType
  site: string
  status: string
  updatedAt: string
}

const tabOptions: Array<{ label: string; value: EncounterTabKey }> = [
  { label: 'All', value: 'all' },
  { label: 'Referrals', value: 'referral' },
  { label: 'Assessments', value: 'assessment' },
  { label: 'Clinic Notes', value: 'clinic_notes' },
]

const selectedTab = ref<EncounterTabKey>('all')

const rows = ref<EncounterRow[]>([
  {
    id: 'enc-001',
    patient: 'Amina Hassan',
    type: 'referral',
    site: 'Kisumu Main Clinic',
    status: 'Pending acceptance',
    updatedAt: '2026-02-15 09:14',
  },
  {
    id: 'enc-002',
    patient: 'David Otieno',
    type: 'assessment',
    site: 'Kisumu Main Clinic',
    status: 'Awaiting assessment',
    updatedAt: '2026-02-15 08:51',
  },
  {
    id: 'enc-003',
    patient: 'Grace Wanjiru',
    type: 'referral',
    site: 'County Referral Unit',
    status: 'Awaiting approval',
    updatedAt: '2026-02-15 08:26',
  },
  {
    id: 'enc-004',
    patient: 'Peter Mboya',
    type: 'clinic_notes',
    site: 'Mumias Satellite',
    status: 'Notes drafted',
    updatedAt: '2026-02-15 07:58',
  },
  {
    id: 'enc-005',
    patient: 'Mercy Akinyi',
    type: 'assessment',
    site: 'Kisumu Main Clinic',
    status: 'In progress',
    updatedAt: '2026-02-15 07:22',
  },
])

const filteredRows = computed(() => {
  if (selectedTab.value === 'all') return rows.value
  return rows.value.filter((row) => row.type === selectedTab.value)
})

const tabCounts = computed(() => {
  return {
    all: rows.value.length,
    referral: rows.value.filter((row) => row.type === 'referral').length,
    assessment: rows.value.filter((row) => row.type === 'assessment').length,
    clinic_notes: rows.value.filter((row) => row.type === 'clinic_notes').length,
  }
})

function formatType(value: EncounterType) {
  if (value === 'referral') return 'Referral'
  if (value === 'assessment') return 'Assessment'
  return 'Clinic Notes'
}

function tabCount(value: EncounterTabKey) {
  return tabCounts.value[value]
}
</script>

<template>
  <PageSectionHeader
    title="Encounter"
    description="Track and triage encounter work from a single queue"
  />

  <Card class="encounters-card">
    <template #content>
      <div class="type-tabs-scroll">
        <div class="type-tabs-menu" role="tablist" aria-label="Encounter type filters">
          <button
            v-for="tab in tabOptions"
            :key="tab.value"
            class="type-tab-btn"
            :class="[{ active: selectedTab === tab.value }, `type-tab-btn--${tab.value}`]"
            role="tab"
            :aria-selected="selectedTab === tab.value"
            @click="selectedTab = tab.value"
          >
            <span class="type-tab-label">
              <span v-if="tab.value !== 'all'" class="type-tab-dot" :class="`type-tab-dot--${tab.value}`" />
              <span>{{ tab.label }}</span>
            </span>
            <span class="type-tab-count">{{ tabCount(tab.value) }}</span>
          </button>
        </div>
      </div>

      <div class="table-wrap">
        <DataTable
          :value="filteredRows"
          class="encounters-table"
          stripedRows
          rowHover
          :tableStyle="{ minWidth: '58rem' }"
        >
          <template #empty>
            <div class="table-empty">No encounters in this tab.</div>
          </template>
          <Column field="patient" header="Patient" style="min-width: 14rem" />
          <Column header="Type" style="min-width: 10rem">
            <template #body="{ data }">
              <span class="type-pill">{{ formatType(data.type) }}</span>
            </template>
          </Column>
          <Column field="site" header="Site" style="min-width: 14rem" />
          <Column field="status" header="Status" style="min-width: 12rem" />
          <Column field="updatedAt" header="Updated" style="min-width: 10rem" />
        </DataTable>
      </div>
    </template>
  </Card>
</template>

<style scoped>
.encounters-card {
  margin-bottom: 0.9rem;
}

.type-tabs-scroll {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  margin-bottom: 0.72rem;
}

.type-tabs-menu {
  display: flex;
  align-items: center;
  gap: 0.24rem;
  width: max-content;
  min-width: 0;
  padding: 0.2rem;
  border: 1px solid color-mix(in srgb, #d4d0e7 76%, var(--p-surface-200));
  border-radius: 0.72rem;
  background: color-mix(in srgb, #e8e5f2 78%, var(--p-surface-0));
}

.type-tab-btn {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: 0.52rem;
  height: var(--app-control-height, 2.2rem);
  border: 0;
  background: transparent;
  border-radius: var(--app-control-radius, 0.62rem);
  color: #85879a;
  padding: 0 0.78rem;
  font-size: var(--app-control-font-size, 0.83rem);
  font-weight: 650;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.14s ease, color 0.14s ease;
}

.type-tab-label {
  display: inline-flex;
  align-items: center;
  gap: 0.44rem;
}

.type-tab-dot {
  width: 0.58rem;
  height: 0.58rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.type-tab-dot--referral {
  background: #f4bf2a;
}

.type-tab-dot--assessment {
  background: #34d399;
}

.type-tab-dot--clinic_notes {
  background: #fb7185;
}

.type-tab-btn:hover {
  background: color-mix(in srgb, var(--p-surface-0) 72%, transparent);
  color: #696d84;
}

.type-tab-btn.active {
  background: #f8f8fc;
  color: #191a32;
  box-shadow: 0 1px 0 color-mix(in srgb, #d7d9e8 70%, transparent);
}

.type-tab-count {
  min-width: 1.26rem;
  height: 1.26rem;
  border-radius: 0.58rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #d8d8e4;
  color: #7f8194;
  font-size: 0.66rem;
  font-weight: 700;
  line-height: 1;
}

.type-tab-btn.active .type-tab-count {
  background: #1a1d38;
  color: #f5f7ff;
}

.table-wrap {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.type-pill {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  border: 1px solid color-mix(in srgb, #0ea5a4 32%, var(--p-surface-200));
  background: color-mix(in srgb, #0ea5a4 12%, var(--p-surface-0));
  color: #0f766e;
  padding: 0.15rem 0.48rem;
  font-size: 0.72rem;
  font-weight: 620;
  line-height: 1.2;
}

.table-empty {
  padding: 0.95rem;
  color: var(--p-surface-600);
  font-size: 0.82rem;
}

@media (max-width: 900px) {
  .type-tabs-menu {
    min-width: max-content;
  }

  .type-tab-btn {
    min-width: 0;
  }
}
</style>
