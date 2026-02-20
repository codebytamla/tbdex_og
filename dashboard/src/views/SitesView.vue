<script setup lang="ts">
import { computed, reactive, ref, type Component } from 'vue'
import { Ban, MoreHorizontal, Pencil, RefreshCcw } from 'lucide-vue-next'
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
import CreateActionButton from '@/components/ui/CreateActionButton.vue'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import ToolbarIconButton from '@/components/ui/ToolbarIconButton.vue'
import {
  TABLE_DEFAULT_ROWS,
  TABLE_PAGE_REPORT_TEMPLATE,
  TABLE_PAGINATOR_TEMPLATE,
  TABLE_ROWS_PER_PAGE_OPTIONS,
} from '@/components/ui/tablePagination'

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
const loading = ref(false)
const pageError = ref('')

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const dialogBusy = ref(false)
const dialogError = ref('')
const selectedSite = ref<Site | null>(null)

const rowActionMenu = ref()
const selectedMenuSite = ref<Site | null>(null)

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
  } catch (error) {
    pageError.value = error instanceof Error ? error.message : 'Unable to load sites.'
    sites.value = []
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
  <PageSectionHeader title="Sites" description="Create and manage active service sites">
    <template #actions>
      <ToolbarIconButton :disabled="!canManageSites" title="Refresh sites" ariaLabel="Refresh sites" @click="loadSites">
        <RefreshCcw :size="16" />
      </ToolbarIconButton>
      <CreateActionButton label="Create site" :disabled="!canManageSites" @click="openCreateDialog" />
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
        :value="sites"
        tableStyle="min-width: 44rem"
        :loading="loading"
        paginator
        :rows="TABLE_DEFAULT_ROWS"
        :rowsPerPageOptions="[...TABLE_ROWS_PER_PAGE_OPTIONS]"
        :paginatorTemplate="TABLE_PAGINATOR_TEMPLATE"
        :currentPageReportTemplate="TABLE_PAGE_REPORT_TEMPLATE"
      >
        <Column field="name" header="Name" />
        <Column field="code" header="Code" />
        <Column header="Type">
          <template #body="{ data }">
            <span>{{ formatSiteType(data.siteType) }}</span>
          </template>
        </Column>
        <Column header="Updated">
          <template #body="{ data }">
            <span>{{ formatDate(data.updatedAt || data.createdAt) }}</span>
          </template>
        </Column>
        <Column header="Status">
          <template #body="{ data }">
            <span class="status-cell" :class="data.isActive ? 'status-active' : 'status-disabled'">
              <span class="status-dot"></span>
              {{ data.isActive ? 'Active' : 'Disabled' }}
            </span>
          </template>
        </Column>
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
      <p class="table-meta table-meta-bottom">Total sites: {{ sites.length }}</p>
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
        <Button label="Cancel" text @click="dialogVisible = false" />
        <Button :label="dialogMode === 'create' ? 'Create site' : 'Save changes'" :loading="dialogBusy" @click="submitSite" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.sites-card {
  margin-bottom: 0.9rem;
}

.table-tools {
  --toolbar-control-height: var(--app-control-height, 2.2rem);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-bottom: 0.8rem;
}

.table-meta {
  margin: 0;
  color: var(--p-surface-600);
  font-size: 0.8rem;
}

.table-meta-bottom {
  margin-top: 0.7rem;
}

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  flex-wrap: wrap;
}

.status-cell {
  display: inline-flex;
  align-items: center;
  gap: 0.46rem;
  padding: 0.22rem 0.58rem 0.22rem 0.46rem;
  border-radius: 999px;
  border: 1px solid transparent;
  font-size: 0.74rem;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.012em;
}

.status-dot {
  width: 0.46rem;
  height: 0.46rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-active {
  color: #0f7a3f;
  border-color: #bfe9cd;
  background: #eaf8f0;
}

.status-active .status-dot {
  background: #16c35d;
  box-shadow: 0 0 0 3px rgb(22 195 93 / 18%);
}

.status-disabled {
  color: #9a4214;
  border-color: #f8d2bd;
  background: #fff2e8;
}

.status-disabled .status-dot {
  background: #f97316;
  box-shadow: 0 0 0 3px rgb(249 115 22 / 16%);
}

.dialog-grid {
  display: grid;
  gap: 0.8rem;
}

.field {
  display: grid;
  gap: 0.35rem;
}

.field span {
  font-size: 0.8rem;
  color: var(--p-surface-600);
}

.dialog-actions {
  display: inline-flex;
  gap: 0.45rem;
  justify-content: flex-end;
  width: 100%;
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
  border: 1px solid var(--p-surface-200);
  box-shadow: none;
}
</style>
