<script setup lang="ts">
import { computed, reactive, ref, type Component } from 'vue'
import { FilterMatchMode, FilterOperator } from '@primevue/core/api'
import { Ban, MoreHorizontal, Pencil } from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Message from 'primevue/message'
import type { MenuItem } from 'primevue/menuitem'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import PrimaryActionDropdown from '@/components/ui/PrimaryActionDropdown.vue'
import DataTableToolbar from '@/components/ui/DataTableToolbar.vue'

type Site = {
  id: string
  name: string
  code: string
  siteType: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

type SiteFormPayload = {
  name: string
  code: string
  site_type: string
}

type ActionMenuItem = MenuItem & {
  lucideIcon: Component
}

const authStore = useAuthStore()
const catalogStore = useCatalogStore()
authStore.hydrate()

const confirm = useConfirm()
const toast = useToast()

const canManageSites = computed(() =>
  authStore.canAccess([
    'sites.*',
    'sites.write',
    'sites.update',
    'sites.delete',
    'sites.manage',
    'users.*',
    'users.write',
    'users.update',
    'users.delete',
    'users.manage',
  ]),
)

const sites = ref<Site[]>([])
const selectedSites = ref<Site[]>([])
const loading = ref(false)
const pageError = ref('')

const filters = ref({
  global: { value: null as string | null, matchMode: FilterMatchMode.CONTAINS },
  name: { operator: FilterOperator.AND, constraints: [{ value: null as string | null, matchMode: FilterMatchMode.CONTAINS }] },
  code: { operator: FilterOperator.AND, constraints: [{ value: null as string | null, matchMode: FilterMatchMode.CONTAINS }] },
  siteType: { operator: FilterOperator.AND, constraints: [{ value: null as string | null, matchMode: FilterMatchMode.EQUALS }] },
})

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const dialogBusy = ref(false)
const dialogError = ref('')
const selectedSite = ref<Site | null>(null)

const rowActionMenu = ref()
const selectedMenuSite = ref<Site | null>(null)
const dt = ref()

const columns = ref([
  { field: 'name', header: 'Name' },
  { field: 'code', header: 'Code' },
  { field: 'siteType', header: 'Type' },
  { field: 'updatedAt', header: 'Updated' },
  { field: 'status', header: 'Status' }
])
const selectedColumns = ref([...columns.value])

const form = reactive({
  name: '',
  code: '',
  siteType: 'site',
})

const siteTypeOptions = [
  { label: 'Site', value: 'site' },
  { label: 'Border Post', value: 'border_post' },
  { label: 'Checkpoint', value: 'checkpoint' },
  { label: 'Mobile Unit', value: 'mobile_unit' },
  { label: 'Other', value: 'other' },
]

const dialogTitle = computed(() => (dialogMode.value === 'create' ? 'Create site' : 'Edit site'))
const canSubmit = computed(
  () =>
    canManageSites.value &&
    !dialogBusy.value &&
    form.name.trim().length > 0 &&
    form.code.trim().length > 0,
)

const rowActionItems = computed<ActionMenuItem[]>(() => [
  {
    label: 'Edit site',
    lucideIcon: Pencil,
    disabled: !canManageSites.value,
    command: () => {
      if (!selectedMenuSite.value) return
      openEditDialog(selectedMenuSite.value)
    },
  },
  {
    label: 'Disable site',
    lucideIcon: Ban,
    disabled: !canManageSites.value || !selectedMenuSite.value?.isActive,
    command: () => {
      if (!selectedMenuSite.value) return
      disableSite(selectedMenuSite.value)
    },
  },
])

async function requestJson<T>(path: string, init: RequestInit, fallbackError: string) {
  return authStore.authenticatedRequestJson<T>(path, init, fallbackError)
}

function parseBooleanLike(value: unknown): boolean | null {
  if (typeof value === 'boolean') return value
  if (typeof value === 'number' && Number.isFinite(value)) return value !== 0
  if (typeof value !== 'string') return null

  const normalized = value.trim().toLowerCase()
  if (!normalized) return null
  if (['true', '1', 'yes', 'y', 'on', 'active', 'enabled'].includes(normalized)) return true
  if (['false', '0', 'no', 'n', 'off', 'inactive', 'disabled'].includes(normalized)) return false
  return null
}

function resolveActiveState(record: Record<string, unknown>, fallback = true) {
  for (const key of ['is_active', 'isActive', 'active', 'enabled']) {
    const parsed = parseBooleanLike(record[key])
    if (parsed !== null) return parsed
  }

  for (const key of ['is_disabled', 'isDisabled', 'disabled', 'inactive']) {
    const parsed = parseBooleanLike(record[key])
    if (parsed !== null) return !parsed
  }

  const statusValue = record.status ?? record.state
  const status = typeof statusValue === 'string' ? statusValue.trim().toLowerCase() : ''
  if (['active', 'enabled'].includes(status)) return true
  if (['inactive', 'disabled', 'suspended', 'blocked', 'archived'].includes(status)) return false

  return fallback
}

function normalizeSite(value: unknown): Site | null {
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>

  const id = typeof record.id === 'string' ? record.id : null
  const name = typeof record.name === 'string' ? record.name.trim() : ''
  const code = typeof record.code === 'string' ? record.code.trim() : ''
  const siteType =
    typeof record.site_type === 'string'
      ? record.site_type.trim()
      : typeof record.siteType === 'string'
        ? record.siteType.trim()
        : 'site'
  const isActive = resolveActiveState(record, true)
  const createdAt =
    typeof record.created_at === 'string'
      ? record.created_at
      : typeof record.createdAt === 'string'
        ? record.createdAt
        : ''
  const updatedAt =
    typeof record.updated_at === 'string'
      ? record.updated_at
      : typeof record.updatedAt === 'string'
        ? record.updatedAt
        : ''

  if (!id || !name || !code) return null

  return {
    id,
    name,
    code,
    siteType: siteType || 'site',
    isActive,
    createdAt,
    updatedAt,
  }
}

function formatDate(value: string) {
  if (!value) return 'N/A'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'N/A'
  return date.toLocaleString(undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatSiteType(value: string) {
  if (!value) return 'site'
  return value
    .split(/[_\s-]+/)
    .filter(Boolean)
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join(' ')
}

function resetForm() {
  form.name = ''
  form.code = ''
  form.siteType = 'site'
  dialogError.value = ''
  selectedSite.value = null
}

async function loadSites() {
  if (!canManageSites.value) return

  loading.value = true
  pageError.value = ''
  try {
    const payload = await requestJson<{ sites?: unknown[] }>(
      '/admin/sites',
      { method: 'GET' },
      'Unable to load sites.',
    )
    sites.value = (payload.sites || []).map(normalizeSite).filter((entry): entry is Site => !!entry)
    const visibleSiteIds = new Set(sites.value.map((site) => site.id))
    selectedSites.value = selectedSites.value.filter((site) => visibleSiteIds.has(site.id))
  } catch (error) {
    pageError.value = error instanceof Error ? error.message : 'Unable to load sites.'
    sites.value = []
    selectedSites.value = []
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  if (!canManageSites.value) return
  resetForm()
  dialogMode.value = 'create'
  dialogVisible.value = true
}

function openEditDialog(site: Site) {
  if (!canManageSites.value) return
  resetForm()
  dialogMode.value = 'edit'
  selectedSite.value = site
  form.name = site.name
  form.code = site.code
  form.siteType = site.siteType || 'site'
  dialogVisible.value = true
}

function toPayload(): SiteFormPayload {
  return {
    name: form.name.trim(),
    code: form.code.trim(),
    site_type: (form.siteType || 'site').trim(),
  }
}

async function submitSite() {
  if (!canSubmit.value) return

  dialogBusy.value = true
  dialogError.value = ''
  try {
    const payload = toPayload()

    if (dialogMode.value === 'create') {
      await requestJson(
        '/admin/sites',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        },
        'Unable to create site.',
      )
      toast.add({ severity: 'success', summary: 'Site created', life: 2400 })
    } else {
      if (!selectedSite.value) throw new Error('Site context missing.')
      await requestJson(
        `/admin/sites/${selectedSite.value.id}`,
        {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        },
        'Unable to update site.',
      )
      toast.add({ severity: 'success', summary: 'Site updated', life: 2400 })
    }

    dialogVisible.value = false
    catalogStore.invalidateSites()
    await loadSites()
  } catch (error) {
    dialogError.value = error instanceof Error ? error.message : 'Unable to save site.'
  } finally {
    dialogBusy.value = false
  }
}

function disableSite(site: Site) {
  if (!canManageSites.value || !site.isActive) return

  confirm.require({
    header: 'Disable site',
    message: `Disable "${site.name}" (${site.code})?`,
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Disable',
    rejectLabel: 'Cancel',
    accept: async () => {
      pageError.value = ''
      try {
        await requestJson(
          `/admin/sites/${site.id}/disable`,
          { method: 'PATCH' },
          'Unable to disable site.',
        )
        toast.add({ severity: 'success', summary: 'Site disabled', life: 2400 })
        catalogStore.invalidateSites()
        await loadSites()
      } catch (error) {
        pageError.value = error instanceof Error ? error.message : 'Unable to disable site.'
        toast.add({ severity: 'error', summary: 'Disable failed', detail: pageError.value, life: 3600 })
      }
    },
  })
}

function openRowActionMenu(event: Event, site: Site) {
  selectedMenuSite.value = site
  rowActionMenu.value?.toggle(event)
}

void loadSites()
</script>

<template>
  <PageSectionHeader title="Sites">
    <template #actions>
      <PrimaryActionDropdown
        label="New Site"
        :disabled="!canManageSites"
        title="New Site options"
        importLabel="Import sites"
        exportLabel="Export sites"
        refreshLabel="Refresh sites"
        :importDisabled="true"
        :exportDisabled="true"
        :refreshDisabled="!canManageSites"
        @click="openCreateDialog"
        @refresh="loadSites"
      />
    </template>
  </PageSectionHeader>

  <Message v-if="!canManageSites" severity="warn" size="small" variant="simple">
    You can view sites, but <code>sites.write</code> or <code>sites.update</code> is required for changes.
  </Message>
  <Message v-if="pageError" severity="error" size="small" variant="simple">
    {{ pageError }}
  </Message>

  <Card class="sites-card">
    <template #content>
      <DataTable
        ref="dt"
        v-model:selection="selectedSites"
        v-model:filters="filters"
        :value="sites"
        dataKey="id"
        paginator
        :rows="10"
        :rowsPerPageOptions="[10, 20, 50]"
        paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        currentPageReportTemplate="Showing {first} to {last} of {totalRecords}"
        tableStyle="min-width: 44rem"
        :loading="loading"
        filterDisplay="menu"
        :globalFilterFields="['name', 'code', 'siteType']"
        removableSort
        stripedRows
        rowHover
      >
        <template #header>
          <DataTableToolbar
            v-model:searchQuery="filters.global.value"
            v-model:selectedColumns="selectedColumns"
            :columns="columns"
            searchPlaceholder="Search sites..."
          />
        </template>
        <Column selectionMode="multiple" style="width: 3rem" />
        <template v-for="col of selectedColumns" :key="col.field">
          <Column v-if="col.field === 'name'" field="name" header="Name" sortable :showFilterOperator="false" :showAddButton="false" />
          <Column v-if="col.field === 'code'" field="code" header="Code" sortable :showFilterOperator="false" :showAddButton="false" />
          <Column v-if="col.field === 'siteType'" field="siteType" header="Type" sortable :showFilterOperator="false" :showAddButton="false">
            <template #body="{ data }">
              <span>{{ formatSiteType(data.siteType) }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'updatedAt'" header="Updated">
            <template #body="{ data }">
              <span>{{ formatDate(data.updatedAt || data.createdAt) }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'status'" header="Status">
            <template #body="{ data }">
              <span class="status-badge" :class="data.isActive ? 'status-active' : 'status-disabled'">
                <span class="status-dot"></span>
                {{ data.isActive ? 'Active' : 'Disabled' }}
              </span>
            </template>
          </Column>
        </template>
        <Column header="" style="width: 3.5rem">
          <template #body="{ data }">
            <div class="actions">
              <RowActionButton
                :disabled="!canManageSites"
                ariaLabel="Open site actions"
                title="Open site actions"
                @click="openRowActionMenu($event, data)"
              >
                <MoreHorizontal :size="16" />
              </RowActionButton>
            </div>
          </template>
        </Column>
      </DataTable>
      <Menu ref="rowActionMenu" :model="rowActionItems" popup>
        <template #item="{ item, props }">
          <a class="p-menu-item-link menu-action-item" v-bind="props.action">
            <component :is="item.lucideIcon" class="menu-action-icon" :size="16" />
            <span class="menu-action-label">{{ item.label }}</span>
          </a>
        </template>
      </Menu>
    </template>
  </Card>

  <Dialog
    v-model:visible="dialogVisible"
    modal
    :header="dialogTitle"
    :style="{ width: 'min(34rem, 96vw)' }"
    @hide="resetForm"
  >
    <div class="dialog-grid">
      <Message v-if="dialogError" severity="error" size="small" variant="simple">
        {{ dialogError }}
      </Message>

      <label class="field">
        <span>Name</span>
        <InputText v-model="form.name" autocomplete="off" />
      </label>

      <label class="field">
        <span>Code</span>
        <InputText v-model="form.code" autocomplete="off" />
      </label>

      <label class="field">
        <span>Site type</span>
        <Select
          v-model="form.siteType"
          :options="siteTypeOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Select site type"
        />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" :class="{ 'create-cancel-btn': dialogMode === 'create' }" :text="dialogMode !== 'create'" @click="dialogVisible = false" />
        <Button :label="dialogMode === 'create' ? 'Create site' : 'Save changes'" :loading="dialogBusy" @click="submitSite" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.sites-card {
  margin-bottom: 0.9rem;
}

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  flex-wrap: wrap;
}

/* ── Status badge ── */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.38rem;
  padding: 0.24rem 0.6rem;
  border-radius: var(--radius-full, 999px);
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
  letter-spacing: 0.012em;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-active {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid rgb(4 120 87 / 0.1);
}

.status-active .status-dot {
  background: #10b981;
  box-shadow: 0 0 0 2px rgb(16 185 129 / 0.15);
}

.status-disabled {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid rgb(180 83 9 / 0.1);
}

.status-disabled .status-dot {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgb(245 158 11 / 0.15);
}

.dialog-grid {
  display: grid;
  gap: 1rem;
}

.field {
  display: grid;
  gap: 0.35rem;
}

.field span {
  font-size: 0.825rem;
  color: #374151;
  font-weight: 600;
  letter-spacing: -0.005em;
}

.dialog-actions {
  display: inline-flex;
  gap: 0.5rem;
  justify-content: flex-end;
  width: 100%;
}

:deep(.create-cancel-btn),
:deep(.create-cancel-btn.p-button) {
  border: none !important;
  border-color: transparent !important;
  box-shadow: none !important;
  background: transparent !important;
  color: #9ca3af !important;
}

:deep(.create-cancel-btn:hover) {
  background: transparent !important;
  color: #374151 !important;
  box-shadow: none !important;
}

.menu-action-item {
  display: inline-flex;
  align-items: center;
  gap: 0.48rem;
}

.menu-action-icon {
  flex: 0 0 auto;
}

.menu-action-label {
  line-height: 1.2;
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


</style>
