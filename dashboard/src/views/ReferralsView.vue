<script setup lang="ts">
import { computed, ref, type Component } from 'vue'
import {
  Ban,
  ChevronRight,
  CircleCheckBig,
  Eye,
  Hospital,
  MapPin,
  MoreHorizontal,
} from 'lucide-vue-next'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Menu from 'primevue/menu'
import MultiSelect from 'primevue/multiselect'
import type { MenuItem } from 'primevue/menuitem'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import PrimaryActionDropdown from '@/components/ui/PrimaryActionDropdown.vue'
import DataTableToolbar from '@/components/ui/DataTableToolbar.vue'

/* ── Types ─────────────────────────────────────────────── */

type ReferralStatus = 'pending' | 'in_transit' | 'arrived' | 'completed' | 'canceled'

type Referral = {
  id: string
  refId: string
  patientName: string
  patientAge: number
  patientSex: string
  patientVillage: string
  patientPhone: string
  referredByName: string
  referredToName: string
  referredToSiteId: string
  reason: string
  status: ReferralStatus
  notes: string
  screeningSummary: string
  createdAt: string
}

type StatusTab = {
  label: string
  value: ReferralStatus | 'all'
}

type ActionMenuItem = MenuItem & {
  lucideIcon: Component
  tone?: 'danger'
}

/* ── Stores ────────────────────────────────────────────── */

const authStore = useAuthStore()
const catalogStore = useCatalogStore()
authStore.hydrate()

const toast = useToast()

/* ── Mock data ─────────────────────────────────────────── */

const allReferrals = ref<Referral[]>([
  {
    id: '1',
    refId: 'REF-0042',
    patientName: 'Samuel Okafor',
    patientAge: 45,
    patientSex: 'Male',
    patientVillage: 'Kivu Hills',
    patientPhone: '+256-772-123456',
    referredByName: 'CHW Sarah Jones',
    referredToName: 'City Central Clinic',
    referredToSiteId: 'site-001',
    reason: 'TB Screening',
    status: 'arrived',
    notes: 'Patient presented with persistent cough and fever for 3 weeks. Requires full TB workup.',
    screeningSummary: 'Chest X-ray ordered. Sputum collection pending. Previous history negative.',
    createdAt: '2023-10-25T09:30:00Z',
  },
  {
    id: '2',
    refId: 'REF-0043',
    patientName: 'Maria Garcia',
    patientAge: 32,
    patientSex: 'Female',
    patientVillage: 'Lakeside',
    patientPhone: '+256-782-654321',
    referredByName: 'CHW David Lee',
    referredToName: 'Northside Health Center',
    referredToSiteId: 'site-002',
    reason: 'Follow-up',
    status: 'pending',
    notes: 'Scheduled follow-up after initial positive screening. Needs confirmatory test.',
    screeningSummary: 'Initial GeneXpert positive. Awaiting culture results.',
    createdAt: '2023-10-24T14:15:00Z',
  },
  {
    id: '3',
    refId: 'REF-0044',
    patientName: 'Kwame Nkrumah',
    patientAge: 58,
    patientSex: 'Male',
    patientVillage: 'Eastgate',
    patientPhone: '+256-701-987654',
    referredByName: 'CHW Fatima Hassan',
    referredToName: 'Regional Hospital',
    referredToSiteId: 'site-003',
    reason: 'TB Screening',
    status: 'in_transit',
    notes: 'High-risk contact of confirmed TB case. Priority referral for immediate evaluation.',
    screeningSummary: 'Symptom screening positive. Contact tracing initiated.',
    createdAt: '2023-10-24T11:00:00Z',
  },
  {
    id: '4',
    refId: 'REF-0045',
    patientName: 'Priya Patel',
    patientAge: 27,
    patientSex: 'Female',
    patientVillage: 'Westfield',
    patientPhone: '+256-773-111222',
    referredByName: 'CHW John Kim',
    referredToName: 'East End Clinic',
    referredToSiteId: 'site-001',
    reason: 'Symptom Check',
    status: 'canceled',
    notes: 'Patient reported symptoms resolved before referral appointment.',
    screeningSummary: 'Symptoms self-resolved. No further action.',
    createdAt: '2023-10-23T16:45:00Z',
  },
  {
    id: '5',
    refId: 'REF-0046',
    patientName: 'Chen Wei',
    patientAge: 63,
    patientSex: 'Male',
    patientVillage: 'Harbor View',
    patientPhone: '+256-704-333444',
    referredByName: 'CHW Elizabeth Chen',
    referredToName: 'Community Health Post',
    referredToSiteId: 'site-002',
    reason: 'TB Screening',
    status: 'arrived',
    notes: 'Elderly patient with chronic cough and weight loss. Diabetic comorbidity.',
    screeningSummary: 'Chest X-ray suggestive. Sputum collected. Awaiting lab results.',
    createdAt: '2023-10-22T08:20:00Z',
  },
  {
    id: '6',
    refId: 'REF-0047',
    patientName: 'Amara Diop',
    patientAge: 19,
    patientSex: 'Female',
    patientVillage: 'Greendale',
    patientPhone: '+256-788-555666',
    referredByName: 'CHW Rose Mwangi',
    referredToName: 'City Central Clinic',
    referredToSiteId: 'site-001',
    reason: 'Contact Tracing',
    status: 'completed',
    notes: 'Household contact of index case. Completed contact tracing protocol.',
    screeningSummary: 'Screening negative. Preventive therapy initiated.',
    createdAt: '2023-10-21T10:00:00Z',
  },
  {
    id: '7',
    refId: 'REF-0048',
    patientName: 'James Ochieng',
    patientAge: 41,
    patientSex: 'Male',
    patientVillage: 'Kivu Hills',
    patientPhone: '+256-752-777888',
    referredByName: 'CHW Sarah Jones',
    referredToName: 'Regional Hospital',
    referredToSiteId: 'site-003',
    reason: 'TB Screening',
    status: 'pending',
    notes: 'Referred for sputum collection and chest imaging. Presenting with night sweats.',
    screeningSummary: 'Pending intake at Regional Hospital.',
    createdAt: '2023-10-20T13:30:00Z',
  },
  {
    id: '8',
    refId: 'REF-0049',
    patientName: 'Fatou Sow',
    patientAge: 35,
    patientSex: 'Female',
    patientVillage: 'Lakeside',
    patientPhone: '+256-706-999000',
    referredByName: 'CHW David Lee',
    referredToName: 'Northside Health Center',
    referredToSiteId: 'site-002',
    reason: 'Follow-up',
    status: 'in_transit',
    notes: 'Follow-up appointment for treatment monitoring. Currently on month 2 of RHZE.',
    screeningSummary: 'On treatment. Month 2 sputum needed for conversion check.',
    createdAt: '2023-10-19T09:00:00Z',
  },
])

/* ── Status tabs ───────────────────────────────────────── */

const statusTabs: StatusTab[] = [
  { label: 'All', value: 'all' },
  { label: 'Pending', value: 'pending' },
  { label: 'In Transit', value: 'in_transit' },
  { label: 'Arrived', value: 'arrived' },
  { label: 'Completed', value: 'completed' },
  { label: 'Canceled', value: 'canceled' },
]

const activeStatusTab = ref<ReferralStatus | 'all'>('all')

/* ── Columns ───────────────────────────────────────────── */

const columns = ref([
  { field: 'refId', header: 'Ref ID' },
  { field: 'patientName', header: 'Patient' },
  { field: 'referredByName', header: 'Referred By' },
  { field: 'referredToName', header: 'Referred To' },
  { field: 'reason', header: 'Reason' },
  { field: 'createdAt', header: 'Date' },
  { field: 'status', header: 'Status' },
])
const selectedColumns = ref([...columns.value])

/* ── Search & filters ──────────────────────────────────── */

const searchQuery = ref<string | null>(null)

const sitesCatalog = computed(() => catalogStore.sites)
const sitesLoading = computed(() => catalogStore.sitesLoading)
const siteFilterOptions = computed(() =>
  sitesCatalog.value
    .map((site) => ({ label: site.name, value: site.id }))
    .sort((a, b) => a.label.localeCompare(b.label)),
)
const selectedSiteIds = ref<string[]>([])

/* ── Row expansion ─────────────────────────────────────── */

const expandedRows = ref<Record<string, boolean>>({})

/* ── Filtered data ─────────────────────────────────────── */

const filteredReferrals = computed(() => {
  let list = allReferrals.value

  if (activeStatusTab.value !== 'all') {
    list = list.filter((r) => r.status === activeStatusTab.value)
  }

  if (selectedSiteIds.value.length > 0) {
    const ids = new Set(selectedSiteIds.value)
    list = list.filter((r) => ids.has(r.referredToSiteId))
  }

  const query = searchQuery.value?.trim().toLowerCase()
  if (query) {
    list = list.filter(
      (r) =>
        r.patientName.toLowerCase().includes(query) ||
        r.refId.toLowerCase().includes(query) ||
        r.referredByName.toLowerCase().includes(query) ||
        r.referredToName.toLowerCase().includes(query) ||
        r.reason.toLowerCase().includes(query),
    )
  }

  return list
})

/* ── Status helpers ────────────────────────────────────── */

function formatStatus(value: ReferralStatus) {
  if (value === 'pending') return 'Pending'
  if (value === 'in_transit') return 'In Transit'
  if (value === 'arrived') return 'Arrived'
  if (value === 'completed') return 'Completed'
  return 'Canceled'
}

function statusClass(value: ReferralStatus) {
  if (value === 'arrived') return 'status-arrived'
  if (value === 'in_transit') return 'status-in-transit'
  if (value === 'completed') return 'status-completed'
  if (value === 'canceled') return 'status-canceled'
  return 'status-pending'
}

/* ── Date helpers ──────────────────────────────────────── */

function formatDate(value: string) {
  if (!value) return 'N/A'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'N/A'
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

/* ── Row actions ───────────────────────────────────────── */

const rowActionMenu = ref()
const selectedMenuReferral = ref<Referral | null>(null)

const rowActionItems = computed<ActionMenuItem[]>(() => [
  {
    label: 'View details',
    lucideIcon: Eye,
    command: () => {
      if (!selectedMenuReferral.value) return
      const id = selectedMenuReferral.value.id
      expandedRows.value = { ...expandedRows.value, [id]: true }
    },
  },
  {
    label: 'Mark as arrived',
    lucideIcon: MapPin,
    disabled:
      !selectedMenuReferral.value ||
      selectedMenuReferral.value.status === 'arrived' ||
      selectedMenuReferral.value.status === 'completed' ||
      selectedMenuReferral.value.status === 'canceled',
    command: () => {
      if (!selectedMenuReferral.value) return
      selectedMenuReferral.value.status = 'arrived'
      toast.add({ severity: 'success', summary: 'Marked as arrived', life: 2500 })
    },
  },
  {
    label: 'Mark as completed',
    lucideIcon: CircleCheckBig,
    disabled:
      !selectedMenuReferral.value ||
      selectedMenuReferral.value.status === 'completed' ||
      selectedMenuReferral.value.status === 'canceled',
    command: () => {
      if (!selectedMenuReferral.value) return
      selectedMenuReferral.value.status = 'completed'
      toast.add({ severity: 'success', summary: 'Marked as completed', life: 2500 })
    },
  },
  {
    label: 'Cancel referral',
    lucideIcon: Ban,
    tone: 'danger',
    disabled:
      !selectedMenuReferral.value ||
      selectedMenuReferral.value.status === 'completed' ||
      selectedMenuReferral.value.status === 'canceled',
    command: () => {
      if (!selectedMenuReferral.value) return
      selectedMenuReferral.value.status = 'canceled'
      toast.add({ severity: 'warn', summary: 'Referral canceled', life: 2500 })
    },
  },
])

function openRowActionMenu(event: Event, referral: Referral) {
  selectedMenuReferral.value = referral
  rowActionMenu.value?.toggle(event)
}

function toggleRowExpansion(referral: Referral) {
  const next = { ...expandedRows.value }
  if (next[referral.id]) {
    delete next[referral.id]
  } else {
    next[referral.id] = true
  }
  expandedRows.value = next
}

function onSearchEnter() {
  /* search is reactive via computed, nothing extra needed */
}

/* ── Status tab count badges ───────────────────────────── */

const statusCounts = computed(() => {
  const counts: Record<string, number> = { all: allReferrals.value.length }
  for (const tab of statusTabs) {
    if (tab.value === 'all') continue
    counts[tab.value] = allReferrals.value.filter((r) => r.status === tab.value).length
  }
  return counts
})

/* ── Load sites catalog ────────────────────────────────── */

async function loadSitesCatalog() {
  try {
    await catalogStore.fetchSites()
  } catch {
    /* silent — optional filter */
  }
}

void loadSitesCatalog()
</script>

<template>
  <PageSectionHeader
    title="Referrals"
    description="Track and manage patient referrals across sites"
  >
    <template #actions>
      <PrimaryActionDropdown
        label="New Referral"
        title="New Referral options"
        ariaLabel="New Referral options"
        importLabel="Import referrals"
        exportLabel="Export referrals"
        refreshLabel="Refresh referrals"
        :importDisabled="true"
        :exportDisabled="true"
      />
    </template>
  </PageSectionHeader>

  <!-- Status filter tabs -->
  <div class="status-tabs">
    <button
      v-for="tab in statusTabs"
      :key="tab.value"
      type="button"
      class="status-tab"
      :class="{ 'status-tab-active': activeStatusTab === tab.value }"
      @click="activeStatusTab = tab.value"
    >
      {{ tab.label }}
      <span class="status-tab-count">{{ statusCounts[tab.value] ?? 0 }}</span>
    </button>
  </div>

  <!-- DataTable Card -->
  <Card class="referrals-card">
    <template #content>
      <DataTable
        v-model:expandedRows="expandedRows"
        :value="filteredReferrals"
        dataKey="id"
        paginator
        :rows="10"
        :rowsPerPageOptions="[10, 20, 50]"
        paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        currentPageReportTemplate="Showing {first} to {last} of {totalRecords}"
        class="referrals-table"
        tableStyle="min-width: 72rem"
        removableSort
        scrollable
        stripedRows
        rowHover
      >
        <template #header>
          <div class="toolbar-row">
            <DataTableToolbar
              v-model:searchQuery="searchQuery"
              v-model:selectedColumns="selectedColumns"
              :columns="columns"
              searchPlaceholder="Search referrals..."
              @search="onSearchEnter"
            >
              <template #beforeSearch>
                <MultiSelect
                  v-model="selectedSiteIds"
                  :options="siteFilterOptions"
                  optionLabel="label"
                  optionValue="value"
                  :loading="sitesLoading"
                  :maxSelectedLabels="0"
                  class="toolbar-site-filter"
                  :class="{ 'toolbar-site-filter--active': selectedSiteIds.length > 0 }"
                  :aria-label="selectedSiteIds.length > 0 ? `${selectedSiteIds.length} site filters active` : 'Filter by site'"
                  :title="selectedSiteIds.length > 0 ? `${selectedSiteIds.length} site filters active` : 'Filter by site'"
                >
                  <template #dropdownicon>
                    <span class="toolbar-site-filter-trigger">
                      <Hospital :size="14" />
                      <span class="toolbar-site-filter-label">sites</span>
                    </span>
                  </template>
                </MultiSelect>
              </template>
            </DataTableToolbar>
          </div>
        </template>

        <template #empty>
          <div class="table-empty">No referrals match the current filters.</div>
        </template>

        <!-- Expand toggle column -->
        <Column style="width: 3rem" class="expand-col">
          <template #body="{ data }">
            <button
              type="button"
              class="expand-toggle-btn"
              :class="{ 'expand-toggle-open': expandedRows[data.id] }"
              :aria-label="expandedRows[data.id] ? 'Collapse details' : 'Expand details'"
              :title="expandedRows[data.id] ? 'Collapse details' : 'Expand details'"
              @click="toggleRowExpansion(data)"
            >
              <ChevronRight :size="14" :stroke-width="2.5" />
            </button>
          </template>
        </Column>

        <!-- Dynamic columns -->
        <template v-for="col of selectedColumns" :key="col.field">
          <Column v-if="col.field === 'refId'" field="refId" header="Ref ID" sortable style="min-width: 8rem">
            <template #body="{ data }">
              <span class="cell-ref-id">{{ data.refId }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'patientName'" field="patientName" header="Patient" sortable style="min-width: 11rem" />
          <Column v-if="col.field === 'referredByName'" field="referredByName" header="Referred By" sortable style="min-width: 11rem">
            <template #body="{ data }">
              <span class="cell-truncate" :title="data.referredByName">{{ data.referredByName }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'referredToName'" field="referredToName" header="Referred To" sortable style="min-width: 12rem">
            <template #body="{ data }">
              <span class="cell-truncate" :title="data.referredToName">{{ data.referredToName }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'reason'" field="reason" header="Reason" sortable style="min-width: 9rem">
            <template #body="{ data }">
              <span class="cell-reason">{{ data.reason }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'createdAt'" header="Date" sortable field="createdAt" style="min-width: 8.5rem">
            <template #body="{ data }">
              <span class="cell-date">{{ formatDate(data.createdAt) }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'status'" header="Status" field="status" sortable style="width: 9rem">
            <template #body="{ data }">
              <span class="status-badge" :class="statusClass(data.status)">
                <span class="status-dot"></span>
                {{ formatStatus(data.status) }}
              </span>
            </template>
          </Column>
        </template>

        <!-- Actions column -->
        <Column header="" style="width: 3.5rem">
          <template #body="{ data }">
            <div class="actions">
              <RowActionButton
                ariaLabel="Open referral actions"
                title="Open referral actions"
                @click="openRowActionMenu($event, data)"
              >
                <MoreHorizontal :size="16" />
              </RowActionButton>
            </div>
          </template>
        </Column>

        <!-- Row expansion -->
        <template #expansion="{ data }">
          <div class="expansion-panel">
            <div class="expansion-grid">
              <div class="expansion-section">
                <h4 class="expansion-heading">Patient Demographics</h4>
                <div class="expansion-detail-list">
                  <div class="expansion-detail-item">
                    <span class="expansion-label">Age</span>
                    <span class="expansion-value">{{ data.patientAge }}</span>
                  </div>
                  <div class="expansion-detail-item">
                    <span class="expansion-label">Sex</span>
                    <span class="expansion-value">{{ data.patientSex }}</span>
                  </div>
                  <div class="expansion-detail-item">
                    <span class="expansion-label">Village</span>
                    <span class="expansion-value">{{ data.patientVillage }}</span>
                  </div>
                  <div class="expansion-detail-item">
                    <span class="expansion-label">Contact</span>
                    <span class="expansion-value">{{ data.patientPhone }}</span>
                  </div>
                </div>
              </div>
              <div class="expansion-section">
                <h4 class="expansion-heading">Referral Notes</h4>
                <p class="expansion-text">{{ data.notes }}</p>

                <h4 class="expansion-heading expansion-heading-mt">Screening Summary</h4>
                <p class="expansion-text">{{ data.screeningSummary }}</p>
              </div>
            </div>
          </div>
        </template>
      </DataTable>

      <Menu ref="rowActionMenu" :model="rowActionItems" popup>
        <template #item="{ item, props }">
          <a
            :class="[
              'p-menu-item-link',
              'menu-action-item',
              { 'menu-action-item-danger': item.tone === 'danger' },
            ]"
            v-bind="props.action"
          >
            <component :is="item.lucideIcon" class="menu-action-icon" :size="16" />
            <span class="menu-action-label">{{ item.label }}</span>
          </a>
        </template>
      </Menu>
    </template>
  </Card>
</template>

<style scoped>
/* ── Status tabs ───────────────────────────────────────── */

.status-tabs {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.status-tab {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: var(--radius-full, 999px);
  background: #ffffff;
  color: #6b7280;
  font-size: 0.82rem;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.3;
}

.status-tab:hover {
  border-color: #d1d5db;
  color: #374151;
  background: #f9fafb;
}

.status-tab-active {
  background: var(--p-primary-600);
  color: #ffffff;
  border-color: var(--p-primary-600);
}

.status-tab-active:hover {
  background: var(--p-primary-700);
  border-color: var(--p-primary-700);
  color: #ffffff;
}

.status-tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.25rem;
  height: 1.25rem;
  padding: 0 0.3rem;
  border-radius: var(--radius-full, 999px);
  font-size: 0.7rem;
  font-weight: 700;
  line-height: 1;
}

.status-tab:not(.status-tab-active) .status-tab-count {
  background: #f3f4f6;
  color: #9ca3af;
}

.status-tab-active .status-tab-count {
  background: rgb(255 255 255 / 0.18);
  color: #ffffff;
}

/* ── Card ──────────────────────────────────────────────── */

:deep(.referrals-card .p-card-body) {
  padding: 0;
  min-width: 0;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
}

:deep(.referrals-card .p-card-content) {
  padding: 0;
  min-width: 0;
  width: 100%;
  max-width: 100%;
}

:deep(.p-card) {
  overflow: hidden;
}

:deep(.p-card-body) {
  padding: 0;
}

:deep(.p-card-content) {
  padding: 0;
}

/* ── Toolbar ───────────────────────────────────────────── */

.toolbar-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
}

:deep(.toolbar-site-filter.p-multiselect) {
  width: auto;
  height: 2rem;
  min-width: 5rem;
  border: 1px solid #e5e7eb;
  border-radius: var(--app-control-radius, 8px);
  background: #ffffff;
  color: #64748b;
  padding: 0;
  box-shadow: var(--shadow-xs, 0 1px 2px 0 rgb(0 0 0 / 0.04));
  transition: border-color 0.15s ease, box-shadow 0.15s ease, background 0.15s ease, color 0.15s ease;
}

:deep(.toolbar-site-filter.p-multiselect:hover) {
  border-color: #d1d5db;
  background: #f8fafc;
}

:deep(.toolbar-site-filter.p-multiselect.p-focus),
:deep(.toolbar-site-filter.p-multiselect:focus-within) {
  border-color: var(--p-primary-400, #94a3b8);
  box-shadow: 0 0 0 3px rgb(100 116 139 / 0.12);
}

:deep(.toolbar-site-filter--active.p-multiselect) {
  border-color: var(--p-primary-500, #64748b);
  color: var(--p-primary-700, #334155);
  background: var(--p-primary-50, #f8fafc);
}

:deep(.toolbar-site-filter .p-multiselect-label-container) {
  display: none;
}

:deep(.toolbar-site-filter .p-multiselect-label) {
  display: none;
}

:deep(.toolbar-site-filter .p-multiselect-dropdown) {
  width: 100%;
  height: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  padding: 0;
}

.toolbar-site-filter-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.3rem;
  font-size: 0.8rem;
  font-weight: 600;
  line-height: 1;
}

.toolbar-site-filter-label {
  letter-spacing: 0.01em;
}

/* ── Table overrides ───────────────────────────────────── */

:deep(.p-datatable-header-cell) {
  background: #ffffff;
  color: #6b7280;
  font-weight: 600;
  font-size: 0.78rem;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  border-bottom: 1px solid #e8eaed;
  padding: 0.75rem 1rem;
  white-space: nowrap;
}

:deep(.p-datatable-tbody > tr > td) {
  color: #1a1d23;
  padding: 0.75rem 1rem;
  white-space: nowrap;
  border-color: #f4f5f7;
  font-size: 0.875rem;
}

/* ── Expand toggle ─────────────────────────────────────── */

.expand-toggle-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  border-radius: var(--radius-sm, 6px);
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s ease;
}

.expand-toggle-btn:hover {
  background: #f4f5f7;
  color: #374151;
}

.expand-toggle-btn:focus-visible {
  outline: 2px solid rgb(100 116 139 / 0.3);
  outline-offset: 1px;
}

.expand-toggle-open {
  color: #374151;
  transform: rotate(90deg);
}

/* ── Cell styles ───────────────────────────────────────── */

.cell-ref-id {
  font-family: 'SF Mono', ui-monospace, 'Cascadia Code', 'Fira Code', monospace;
  font-size: 0.8rem;
  color: #6b7280;
  font-weight: 500;
  letter-spacing: -0.01em;
}

.cell-truncate {
  max-width: 14rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}

.cell-reason {
  display: inline-flex;
  align-items: center;
  padding: 0.125rem 0.5rem;
  background: #f3f4f6;
  border-radius: var(--radius-full, 999px);
  font-size: 0.78rem;
  font-weight: 500;
  color: #4b5563;
}

.cell-date {
  color: #6b7280;
  font-size: 0.835rem;
}

/* ── Status badges ─────────────────────────────────────── */

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.1875rem 0.625rem;
  border-radius: var(--radius-full, 999px);
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
  letter-spacing: -0.005em;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-pending {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid rgb(180 83 9 / 0.1);
}
.status-pending .status-dot {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgb(245 158 11 / 0.15);
}

.status-in-transit {
  background: var(--p-primary-50, #f8fafc);
  color: var(--p-primary-700, #334155);
  border: 1px solid color-mix(in srgb, var(--p-primary-700, #334155) 10%, transparent);
}
.status-in-transit .status-dot {
  background: var(--p-primary-500, #64748b);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-primary-500, #64748b) 15%, transparent);
}

.status-arrived {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid rgb(4 120 87 / 0.1);
}
.status-arrived .status-dot {
  background: #10b981;
  box-shadow: 0 0 0 2px rgb(16 185 129 / 0.15);
}

.status-completed {
  background: #f0fdf4;
  color: #15803d;
  border: 1px solid rgb(21 128 61 / 0.1);
}
.status-completed .status-dot {
  background: #22c55e;
  box-shadow: 0 0 0 2px rgb(34 197 94 / 0.15);
}

.status-canceled {
  background: #f9fafb;
  color: #6b7280;
  border: 1px solid #e5e7eb;
}
.status-canceled .status-dot {
  background: #9ca3af;
}

/* ── Expansion panel ───────────────────────────────────── */

.expansion-panel {
  padding: 1.25rem 1.5rem 1.5rem;
  background: #f8fafb;
  border-top: 1px solid #f0f1f3;
}

.expansion-grid {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 2rem;
}

.expansion-heading {
  margin: 0 0 0.625rem;
  font-size: 0.82rem;
  font-weight: 700;
  color: #1a1d23;
  letter-spacing: -0.005em;
}

.expansion-heading-mt {
  margin-top: 1.25rem;
}

.expansion-detail-list {
  display: grid;
  gap: 0.5rem;
}

.expansion-detail-item {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
}

.expansion-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #9ca3af;
  min-width: 4.5rem;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.expansion-value {
  font-size: 0.855rem;
  color: #374151;
}

.expansion-text {
  margin: 0;
  font-size: 0.855rem;
  color: #4b5563;
  line-height: 1.55;
}

/* ── Actions ───────────────────────────────────────────── */

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.menu-action-item {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.menu-action-icon {
  flex: 0 0 auto;
  color: #9ca3af;
}

.menu-action-label {
  line-height: 1.2;
}

.menu-action-item-danger {
  color: #dc2626;
}

.menu-action-item-danger .menu-action-icon {
  color: #dc2626;
}

:deep(.p-menu-item-link.menu-action-item-danger:hover) {
  color: #b91c1c;
  background: #fef2f2;
}

/* ── Empty state ───────────────────────────────────────── */

.table-empty {
  padding: 3rem 1.5rem;
  color: #9ca3af;
  font-size: 0.875rem;
  text-align: center;
}

/* ── Responsive ────────────────────────────────────────── */

@media (max-width: 900px) {
  .toolbar-row {
    flex-direction: column;
    align-items: stretch;
  }

  .expansion-grid {
    grid-template-columns: 1fr;
    gap: 1.25rem;
  }

  .status-tabs {
    gap: 0.25rem;
  }
}
</style>
