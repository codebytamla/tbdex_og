<script setup lang="ts">
import { ref } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import { ChevronDown, MoreHorizontal } from 'lucide-vue-next'
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
        :value="referrals"
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
            <InputText
              v-model="filters.global.value"
              placeholder="Search patient, provider, or ID..."
              class="toolbar-search"
            />
            <MultiSelect
              v-model="selectedColumns"
              :options="columns"
              optionLabel="header"
              :maxSelectedLabels="0"
              selectedItemsLabel="Columns"
              placeholder="Columns"
              class="toolbar-columns"
            >
              <template #dropdownicon>
                <ChevronDown :size="13" />
              </template>
            </MultiSelect>
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

/* Match the generous py-4 row height from the design */
:deep(.p-datatable-tbody > tr > td) {
  padding: 1rem;
}

:deep(.p-datatable-header-cell) {
  padding: 0.8rem 1rem;
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

.toolbar-search {
  width: 20rem;
  font-size: var(--app-control-font-size, 0.9rem);
}

/* Make the MultiSelect look like a compact button */
.toolbar-columns :deep(.p-multiselect) {
  height: var(--app-control-height, 2.4rem);
  border-color: var(--p-surface-200);
  background: var(--p-surface-50);
  font-size: var(--app-control-font-size, 0.9rem);
  font-weight: 600;
  color: var(--p-surface-700);
  padding: 0 0.65rem;
  gap: 0.4rem;
}
.toolbar-columns :deep(.p-multiselect-label) {
  padding: 0;
  font-size: var(--app-control-font-size, 0.9rem);
  font-weight: 600;
}
.toolbar-columns :deep(.p-multiselect-dropdown) {
  padding: 0;
  width: auto;
  color: var(--p-surface-500);
}

/* ── Ref ID link ── */
.ref-id-link {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  font-weight: 700;
  font-size: inherit;
  font-family: inherit;
  color: var(--p-primary-600);
}
.ref-id-link:hover {
  color: var(--p-primary-700);
  text-decoration: underline;
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
