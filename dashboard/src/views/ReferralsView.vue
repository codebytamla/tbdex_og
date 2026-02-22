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
import CreateActionButton from '@/components/ui/CreateActionButton.vue'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import TableActionDrawerButton from '@/components/ui/TableActionDrawerButton.vue'
import {
  TABLE_PAGE_REPORT_TEMPLATE,
  TABLE_PAGINATOR_TEMPLATE,
  TABLE_ROWS_PER_PAGE_OPTIONS,
} from '@/components/ui/tablePagination'

type ReferralPriority = 'URGENT' | 'ROUTINE'
type ReferralStatus = 'Pending' | 'Accepted' | 'Completed'
type ReferralStatusTab = 'all' | 'pending' | 'accepted' | 'completed'

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
    status: 'Pending',
  },
  {
    id: 'REF-002',
    patient: 'Sarah Jenkins',
    referredBy: 'System Admin',
    referredTo: 'Dr. Smith (Cardiology)',
    date: '02/19/2026',
    priority: 'ROUTINE',
    status: 'Accepted',
  },
  {
    id: 'REF-003',
    patient: 'David Rodriguez',
    referredBy: 'Clinic North',
    referredTo: 'Main Hospital',
    date: '02/18/2026',
    priority: 'ROUTINE',
    status: 'Completed',
  },
])

const statusTabs: Array<{ key: ReferralStatusTab; label: string }> = [
  { key: 'all', label: 'All' },
  { key: 'pending', label: 'Pending' },
  { key: 'accepted', label: 'Accepted' },
  { key: 'completed', label: 'Completed' },
]
const activeStatusTab = ref<ReferralStatusTab>('all')

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
const filteredReferrals = computed(() => {
  if (activeStatusTab.value === 'all') return referrals.value
  return referrals.value.filter((entry) => entry.status.toLowerCase() === activeStatusTab.value)
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
      <CreateActionButton label="New Referral" />
      <TableActionDrawerButton
        title="Referral table actions"
        ariaLabel="Open referral table actions"
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
            <div class="toolbar-left">
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
            <InputText
              v-model="filters.global.value"
              placeholder="Search patient, provider, or ID..."
              class="toolbar-search"
            />
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
                :class="`status-${data.status.toLowerCase()}`"
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
  border: 1px solid var(--p-surface-200);
  box-shadow: none;
}

:deep(.p-card-body) {
  padding: 0.25rem 0;
}

:deep(.p-card-content) {
  padding: 0;
}

:deep(.p-datatable-header-cell) {
  background: var(--p-surface-50);
  color: var(--p-surface-700);
}

/* ── Toolbar ── */
.table-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.1rem 0;
}
.toolbar-left {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  min-width: 0;
}

.toolbar-search {
  width: 20rem;
  height: var(--app-control-height, 2.4rem);
  border-radius: var(--app-control-radius, 0.46rem);
  padding: 0 0.78rem;
  font-size: var(--app-control-font-size, 0.9rem);
}

/* Make the MultiSelect look like a compact button */
:deep(.toolbar-columns.p-multiselect) {
  height: var(--app-control-height, 2.4rem);
  border: 1px solid var(--p-surface-300);
  border-radius: var(--app-control-radius, 0.46rem);
  background: var(--p-surface-0);
  font-size: var(--app-control-font-size, 0.9rem);
  font-weight: 600;
  color: var(--p-surface-900);
  padding: 0 0.78rem;
  gap: 0.45rem;
  box-shadow: 0 1px 1px rgb(15 23 42 / 0.04);
}
:deep(.toolbar-columns.p-multiselect:hover) {
  border-color: var(--p-surface-400);
  background: var(--p-surface-50);
}
:deep(.toolbar-columns.p-multiselect.p-focus),
:deep(.toolbar-columns.p-multiselect:focus-within) {
  border-color: color-mix(in srgb, var(--p-surface-400) 70%, transparent);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-surface-300) 28%, transparent);
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
  color: var(--p-surface-600);
}
.columns-trigger-content {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  white-space: nowrap;
  color: var(--p-surface-900);
  line-height: 1;
}
.columns-trigger-icon {
  color: var(--p-surface-600);
  font-size: 0.82rem;
}

.status-tabs {
  display: inline-flex;
  align-items: center;
  gap: 0.18rem;
  padding: 0.16rem;
  border: 1px solid var(--p-surface-300);
  border-radius: var(--app-control-radius, 0.46rem);
  background: var(--p-surface-50);
}
.status-tab {
  border: none;
  background: transparent;
  color: var(--p-surface-600);
  height: calc(var(--app-control-height, 2.4rem) - 0.32rem);
  border-radius: calc(var(--app-control-radius, 0.46rem) - 0.08rem);
  padding: 0 0.65rem;
  font-size: 0.84rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: background-color 0.15s ease, color 0.15s ease;
}
.status-tab:hover {
  background: color-mix(in srgb, var(--p-surface-900) 8%, transparent);
  color: var(--p-surface-800);
}
.status-tab:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-primary-color) 16%, transparent);
}
.status-tab-active {
  background: var(--p-surface-0);
  color: var(--p-surface-950);
  box-shadow: 0 1px 1px rgb(15 23 42 / 0.06);
}

/* ── Ref ID link ── */
.ref-id-link {
  display: inline-flex;
  align-items: center;
  background: none;
  border: none;
  padding: 0.08rem 0.12rem;
  margin: -0.08rem -0.12rem;
  cursor: pointer;
  border-radius: var(--app-control-radius, 0.46rem);
  font-weight: 600;
  letter-spacing: 0.01em;
  font-size: inherit;
  font-family: inherit;
  color: var(--p-primary-600);
  text-decoration: underline;
  text-decoration-color: transparent;
  text-underline-offset: 0.14em;
  text-decoration-thickness: 1.2px;
  transition:
    color 0.15s ease,
    text-decoration-color 0.15s ease,
    background-color 0.15s ease;
}
.ref-id-link:hover {
  color: var(--p-primary-700);
  text-decoration-color: currentColor;
}
.ref-id-link:focus-visible {
  outline: none;
  color: var(--p-primary-700);
  text-decoration-color: currentColor;
  background: color-mix(in srgb, var(--p-primary-color) 10%, transparent);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-primary-color) 20%, transparent);
}

/* ── Patient name ── */
.patient-name {
  font-weight: 400;
  color: var(--p-surface-900);
}

/* ── Priority badge ── */
.priority-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.2rem 0.55rem;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  white-space: nowrap;
}
.priority-urgent {
  background: #fff7ed;
  color: #c2410c;
}
.priority-routine {
  background: var(--p-surface-100);
  color: var(--p-surface-600);
}

/* ── Status badge ── */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.38rem;
  padding: 0.24rem 0.6rem;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
}
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-pending {
  background: #fffbeb;
  color: #92400e;
}
.status-pending .status-dot {
  background: #f59e0b;
}
.status-accepted {
  background: #ecfdf5;
  color: #065f46;
}
.status-accepted .status-dot {
  background: #10b981;
}
.status-completed {
  background: #eff6ff;
  color: #1e40af;
}
.status-completed .status-dot {
  background: #3b82f6;
}

/* ── Row actions ── */
.row-actions {
  display: flex;
  justify-content: center;
}

/* ── Empty state ── */
.table-empty {
  padding: 1.5rem;
  text-align: center;
  color: var(--p-surface-500);
  font-size: 0.9rem;
}
</style>
