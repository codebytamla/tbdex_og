<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import { ChevronDown, Columns2, MoreHorizontal } from 'lucide-vue-next'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import MultiSelect from 'primevue/multiselect'
import type { MenuItem } from 'primevue/menuitem'
import PrimaryActionDropdown from '@/components/ui/PrimaryActionDropdown.vue'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import {
  TABLE_PAGE_REPORT_TEMPLATE,
  TABLE_PAGINATOR_TEMPLATE,
  TABLE_ROWS_PER_PAGE_OPTIONS,
} from '@/components/ui/tablePagination'

type ReferralPriority = 'URGENT' | 'ROUTINE'
type ReferralStatus = 'In Progress' | 'Arrived' | 'Canceled'
type ReferralStatusTab = 'all' | 'in-progress' | 'arrived' | 'canceled'

type Referral = {
  id: string
  patient: string
  referredBy: string
  referredTo: string
  date: string
  priority: ReferralPriority
  status: ReferralStatus
}

const referrals = ref<Referral[]>([
  {
    id: 'REF-001',
    patient: 'Michael Chen',
    referredBy: 'Dr. Aw',
    referredTo: 'Pulmonology Dept',
    date: '02/20/2026',
    priority: 'URGENT',
    status: 'In Progress',
  },
  {
    id: 'REF-002',
    patient: 'Sarah Jenkins',
    referredBy: 'System Admin',
    referredTo: 'Dr. Smith (Cardiology)',
    date: '02/19/2026',
    priority: 'ROUTINE',
    status: 'Arrived',
  },
  {
    id: 'REF-003',
    patient: 'David Rodriguez',
    referredBy: 'Clinic North',
    referredTo: 'Main Hospital',
    date: '02/18/2026',
    priority: 'ROUTINE',
    status: 'Canceled',
  },
])

const statusTabs: Array<{ key: ReferralStatusTab; label: string }> = [
  { key: 'all', label: 'All' },
  { key: 'in-progress', label: 'In Progress' },
  { key: 'arrived', label: 'Arrived' },
  { key: 'canceled', label: 'Canceled' },
]
const activeStatusTab = ref<ReferralStatusTab>('in-progress')

const columns = ref([
  { field: 'id', header: 'Ref ID' },
  { field: 'patient', header: 'Patient' },
  { field: 'referredBy', header: 'Referred By' },
  { field: 'referredTo', header: 'Referred To' },
  { field: 'date', header: 'Date' },
  { field: 'priority', header: 'Priority' },
  { field: 'status', header: 'Status' },
])
const selectedColumns = ref([...columns.value])

const filters = ref({
  global: { value: null as string | null, matchMode: FilterMatchMode.CONTAINS },
})

const selectedReferrals = ref<Referral[]>([])
const statusTabMap: Record<ReferralStatusTab, ReferralStatus | null> = {
  'all': null,
  'in-progress': 'In Progress',
  'arrived': 'Arrived',
  'canceled': 'Canceled',
}
const filteredReferrals = computed(() => {
  const statusFilter = statusTabMap[activeStatusTab.value]
  if (!statusFilter) return referrals.value
  return referrals.value.filter((entry) => entry.status === statusFilter)
})

watch(filteredReferrals, (visibleRows) => {
  const visibleIds = new Set(visibleRows.map((entry) => entry.id))
  selectedReferrals.value = selectedReferrals.value.filter((entry) => visibleIds.has(entry.id))
})

const rowActionMenu = ref()
const selectedMenuReferral = ref<Referral | null>(null)

const rowActionItems: MenuItem[] = [
  { label: 'View details', icon: 'pi pi-eye' },
  { label: 'Edit referral', icon: 'pi pi-pencil' },
  { separator: true },
  { label: 'Cancel referral', icon: 'pi pi-times', class: 'menu-item-danger' },
]

function openRowActionMenu(event: MouseEvent, referral: Referral) {
  selectedMenuReferral.value = referral
  rowActionMenu.value?.toggle(event)
}
</script>

<template>
  <PageSectionHeader
    title="Referrals"
    description="Manage incoming and outgoing patient referrals across all sites"
  >
    <template #actions>
      <PrimaryActionDropdown
        label="New Referral"
        title="Create new referral options"
        importLabel="Import referrals"
        exportLabel="Export to CSV"
        refreshLabel="Refresh data"
        :importDisabled="false"
        :exportDisabled="false"
        :refreshDisabled="false"
      />
    </template>
  </PageSectionHeader>

  <Card class="referrals-card">
    <template #content>
      <DataTable
        v-model:selection="selectedReferrals"
        v-model:filters="filters"
        :value="filteredReferrals"
        dataKey="id"
        filterDisplay="menu"
        removableSort
        paginator
        :rows="20"
        :rowsPerPageOptions="[...TABLE_ROWS_PER_PAGE_OPTIONS]"
        :paginatorTemplate="TABLE_PAGINATOR_TEMPLATE"
        :currentPageReportTemplate="TABLE_PAGE_REPORT_TEMPLATE"
        tableStyle="min-width: 44rem"
        stripedRows
        rowHover
      >
        <template #header>
          <div class="table-toolbar">
            <div class="toolbar-row toolbar-row-tabs">
              <div class="status-tabs" role="tablist" aria-label="Filter referrals by status">
                <button
                  v-for="tab in statusTabs"
                  :key="tab.key"
                  type="button"
                  role="tab"
                  class="status-tab"
                  :class="{ 'status-tab-active': activeStatusTab === tab.key }"
                  :aria-selected="activeStatusTab === tab.key"
                  @click="activeStatusTab = tab.key"
                >
                  {{ tab.label }}
                </button>
              </div>
            </div>
            <div class="toolbar-row toolbar-row-controls">
              <InputText
                v-model="filters.global.value"
                placeholder="Search patient, provider, or ID..."
                class="toolbar-search"
              />
              <MultiSelect
                v-model="selectedColumns"
                :options="columns"
                optionLabel="header"
                placeholder="Columns"
                class="toolbar-columns"
              >
                <template #value>
                  <span class="columns-trigger-content">
                    <Columns2 :size="15" class="columns-trigger-icon" aria-hidden="true" />
                    <span>Columns</span>
                  </span>
                </template>
                <template #dropdownicon>
                  <ChevronDown :size="13" />
                </template>
              </MultiSelect>
            </div>
          </div>
        </template>

        <template #empty>
          <div class="table-empty">No referrals match the current filters.</div>
        </template>

        <Column selectionMode="multiple" style="width: 3rem" />

        <template v-for="col of selectedColumns" :key="col.field">
          <Column
            v-if="col.field === 'id'"
            field="id"
            header="Ref ID"
            sortable
            style="width: 7.5rem"
          >
            <template #body="{ data }">
              <button class="ref-id-link">{{ data.id }}</button>
            </template>
          </Column>
          <Column
            v-else-if="col.field === 'patient'"
            field="patient"
            header="Patient"
            sortable
          >
            <template #body="{ data }">
              <span class="patient-name">{{ data.patient }}</span>
            </template>
          </Column>
          <Column
            v-else-if="col.field === 'referredBy'"
            field="referredBy"
            header="Referred By"
          />
          <Column
            v-else-if="col.field === 'referredTo'"
            field="referredTo"
            header="Referred To"
          />
          <Column
            v-else-if="col.field === 'date'"
            field="date"
            header="Date"
            sortable
          />
          <Column
            v-else-if="col.field === 'priority'"
            field="priority"
            header="Priority"
            style="width: 8rem"
          >
            <template #body="{ data }">
              <span
                class="priority-badge"
                :class="data.priority === 'URGENT' ? 'priority-urgent' : 'priority-routine'"
              >
                {{ data.priority }}
              </span>
            </template>
          </Column>
          <Column
            v-else-if="col.field === 'status'"
            field="status"
            header="Status"
            style="width: 9rem"
          >
            <template #body="{ data }">
              <span
                class="status-badge"
                :class="`status-${data.status.toLowerCase().replace(/\s+/g, '-')}`"
              >
                <span class="status-dot" />
                {{ data.status }}
              </span>
            </template>
          </Column>
        </template>

        <Column header="" style="width: 3.5rem">
          <template #body="{ data }">
            <div class="row-actions">
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
      </DataTable>

      <Menu ref="rowActionMenu" :model="rowActionItems" popup />
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

:deep(.p-card-body) {
  padding: 0;
}

:deep(.p-card-content) {
  padding: 0;
}

:deep(.p-datatable-header-cell) {
  background: #ffffff;
  color: #6b7280;
}

/* ── Toolbar ── */
.table-toolbar {
  display: flex;
  flex-direction: column;
  gap: 0;
}
.toolbar-row {
  display: flex;
  align-items: center;
}
.toolbar-row-tabs {
  border-bottom: 1px solid #f0f1f3;
  padding: 0 1rem 0;
}
.toolbar-row-controls {
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
}

.toolbar-search {
  width: 20rem;
  height: var(--app-control-height, 2.25rem);
  border-radius: var(--app-control-radius, 8px);
  padding: 0 0.875rem;
  font-size: var(--app-control-font-size, 0.855rem);
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}
.toolbar-search:hover {
  border-color: #d1d5db;
  background: #ffffff;
}
.toolbar-search:focus {
  border-color: var(--p-primary-400, #2dd4bf);
  background: #ffffff;
  box-shadow: 0 0 0 3px rgb(13 148 136 / 0.08);
}

/* Make the MultiSelect look like a compact button */
:deep(.toolbar-columns.p-multiselect) {
  height: var(--app-control-height, 2.25rem);
  border: 1px solid #e5e7eb;
  border-radius: var(--app-control-radius, 8px);
  background: #ffffff;
  font-size: var(--app-control-font-size, 0.855rem);
  font-weight: 600;
  color: #374151;
  padding: 0 0.75rem;
  gap: 0.375rem;
  box-shadow: var(--shadow-xs, 0 1px 2px 0 rgb(0 0 0 / 0.04));
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
:deep(.toolbar-columns.p-multiselect:hover) {
  border-color: #d1d5db;
}
:deep(.toolbar-columns.p-multiselect.p-focus),
:deep(.toolbar-columns.p-multiselect:focus-within) {
  border-color: var(--p-primary-400, #2dd4bf);
  box-shadow: 0 0 0 3px rgb(13 148 136 / 0.08);
}
:deep(.toolbar-columns .p-multiselect-label-container) {
  display: flex;
  align-items: center;
}
:deep(.toolbar-columns .p-multiselect-label) {
  padding: 0;
  display: flex;
  align-items: center;
}
:deep(.toolbar-columns .p-multiselect-dropdown) {
  padding: 0;
  width: auto;
  color: #9ca3af;
}
.columns-trigger-content {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  white-space: nowrap;
  color: #374151;
  line-height: 1;
}
.columns-trigger-icon {
  color: #9ca3af;
}

.status-tabs {
  display: inline-flex;
  align-items: center;
  gap: 0;
}
.status-tab {
  border: none;
  background: transparent;
  color: #9ca3af;
  height: 2.5rem;
  padding: 0 1rem;
  font-size: 0.855rem;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  position: relative;
  transition: color 0.15s ease;
  font-family: inherit;
}
.status-tab::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0.5rem;
  right: 0.5rem;
  height: 2px;
  border-radius: 2px 2px 0 0;
  background: transparent;
  transition: background 0.2s ease;
}
.status-tab:hover {
  color: #374151;
}
.status-tab:focus-visible {
  outline: none;
  color: #374151;
}
.status-tab-active {
  color: #111827;
  font-weight: 700;
}
.status-tab-active::after {
  background: #0f172a;
}

/* ── Ref ID link ── */
.ref-id-link {
  display: inline-flex;
  align-items: center;
  background: none;
  border: none;
  padding: 0.125rem 0.25rem;
  margin: -0.125rem -0.25rem;
  cursor: pointer;
  border-radius: var(--radius-sm, 6px);
  font-weight: 700;
  letter-spacing: -0.01em;
  font-size: inherit;
  font-family: inherit;
  color: #0d9488;
  text-decoration: none;
  transition: color 0.15s ease, background-color 0.15s ease;
}
.ref-id-link:hover {
  color: #0f766e;
  background: rgb(13 148 136 / 0.06);
}
.ref-id-link:focus-visible {
  outline: none;
  color: #0f766e;
  background: rgb(13 148 136 / 0.08);
  box-shadow: 0 0 0 2px rgb(13 148 136 / 0.15);
}

/* ── Patient name ── */
.patient-name {
  font-weight: 600;
  color: #111827;
  letter-spacing: -0.005em;
}

/* ── Priority badge ── */
.priority-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.1875rem 0.5rem;
  border-radius: var(--radius-sm, 6px);
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  white-space: nowrap;
}
.priority-urgent {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid rgb(220 38 38 / 0.12);
}
.priority-routine {
  background: #f9fafb;
  color: #6b7280;
  border: 1px solid #e5e7eb;
}

/* ── Status badge ── */
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
.status-in-progress {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid rgb(180 83 9 / 0.1);
}
.status-in-progress .status-dot {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgb(245 158 11 / 0.15);
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
.status-canceled {
  background: #f9fafb;
  color: #6b7280;
  border: 1px solid #e5e7eb;
}
.status-canceled .status-dot {
  background: #9ca3af;
}

/* ── Row actions ── */
.row-actions {
  display: flex;
  justify-content: center;
}

/* ── Empty state ── */
.table-empty {
  padding: 3rem 1.5rem;
  text-align: center;
  color: #9ca3af;
  font-size: 0.875rem;
}
</style>
