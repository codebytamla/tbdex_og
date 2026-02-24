<script setup lang="ts">
import { computed, reactive, ref, type Component } from 'vue'
import { FilterMatchMode, FilterOperator } from '@primevue/core/api'
import { Eye, MoreHorizontal, Pencil, Trash2 } from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Message from 'primevue/message'
import MultiSelect from 'primevue/multiselect'
import type { MenuItem } from 'primevue/menuitem'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import PrimaryActionDropdown from '@/components/ui/PrimaryActionDropdown.vue'
import DataTableToolbar from '@/components/ui/DataTableToolbar.vue'

type DialogMode = 'create' | 'edit' | 'view'

type RoleSummary = {
  id: string
  name: string
  description: string
}

type PermissionItem = {
  id?: string
  name: string
}

type RoleDetails = RoleSummary & {
  permissions: PermissionItem[]
}

type RoleFormSnapshot = {
  name: string
  description: string
  permissions: string[]
}

type ActionMenuItem = MenuItem & {
  lucideIcon: Component
}

const authStore = useAuthStore()
const catalogStore = useCatalogStore()
authStore.hydrate()

const canManageRoles = computed(() =>
  authStore.canAccess([
    'roles.*',
    'roles.write',
    'roles.update',
    'roles.delete',
    'roles.manage',
    'users.*',
    'users.write',
    'users.update',
    'users.delete',
    'users.manage',
  ]),
)

const roles = ref<RoleSummary[]>([])
const selectedRoles = ref<RoleSummary[]>([])
const rolesLoading = ref(false)

const filters = ref({
  global: { value: null as string | null, matchMode: FilterMatchMode.CONTAINS },
  name: { operator: FilterOperator.AND, constraints: [{ value: null as string | null, matchMode: FilterMatchMode.CONTAINS }] },
  description: { operator: FilterOperator.AND, constraints: [{ value: null as string | null, matchMode: FilterMatchMode.CONTAINS }] },
})
const permissionsCatalog = computed(() => catalogStore.permissions)
const permissionsLoading = computed(() => catalogStore.permissionsLoading)
const pageError = ref('')
const confirm = useConfirm()
const toast = useToast()

const dialogVisible = ref(false)
const dialogMode = ref<DialogMode>('create')
const dialogBusy = ref(false)
const dialogError = ref('')
const roleDetailsLoading = ref(false)
const selectedRoleId = ref<string | null>(null)
const originalSnapshot = ref<RoleFormSnapshot | null>(null)
const rowActionMenu = ref()
const selectedMenuRole = ref<RoleSummary | null>(null)
const dt = ref()

const columns = ref([
  { field: 'name', header: 'Role' },
  { field: 'description', header: 'Description' }
])
const selectedColumns = ref([...columns.value])

const form = reactive({
  name: '',
  description: '',
  permissions: [] as string[],
})

const dialogTitle = computed(() => {
  if (dialogMode.value === 'create') return 'Create Role'
  if (dialogMode.value === 'edit') return 'Edit Role'
  return 'Role Details'
})

const canSubmit = computed(() => {
  if (dialogMode.value === 'view') return false
  return form.name.trim().length > 0 && !dialogBusy.value && canManageRoles.value
})

const selectedPermissionsPreview = computed(() =>
  [...new Set(form.permissions.map((item) => item.trim()).filter(Boolean))].sort(),
)
const rowActionItems = computed<ActionMenuItem[]>(() => [
  {
    label: 'View',
    lucideIcon: Eye,
    disabled: !canManageRoles.value,
    command: () => {
      if (!selectedMenuRole.value) return
      void openRoleDialog('view', selectedMenuRole.value)
    },
  },
  {
    label: 'Edit',
    lucideIcon: Pencil,
    disabled: !canManageRoles.value,
    command: () => {
      if (!selectedMenuRole.value) return
      void openRoleDialog('edit', selectedMenuRole.value)
    },
  },
  {
    label: 'Delete',
    lucideIcon: Trash2,
    disabled: !canManageRoles.value,
    command: () => {
      if (!selectedMenuRole.value) return
      deleteRole(selectedMenuRole.value)
    },
  },
])

async function requestJson<T>(path: string, init: RequestInit, fallbackError: string) {
  return authStore.authenticatedRequestJson<T>(path, init, fallbackError)
}

function normalizePermissionName(value: unknown) {
  if (typeof value === 'string') {
    const normalized = value.trim()
    return normalized || null
  }
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  for (const key of ['name', 'permission', 'slug', 'code']) {
    const candidate = record[key]
    if (typeof candidate === 'string' && candidate.trim()) return candidate.trim()
  }
  return null
}

function normalizeRoleSummary(value: unknown): RoleSummary | null {
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  const id = typeof record.id === 'string' ? record.id : null
  const name = typeof record.name === 'string' ? record.name.trim() : ''
  const description = typeof record.description === 'string' ? record.description.trim() : ''
  if (!id || !name) return null
  return { id, name, description }
}

function normalizePermissions(value: unknown) {
  if (!Array.isArray(value)) return [] as PermissionItem[]
  const byName = new Map<string, PermissionItem>()
  for (const item of value) {
    const name = normalizePermissionName(item)
    if (!name) continue
    const id =
      item && typeof item === 'object' && typeof (item as { id?: unknown }).id === 'string'
        ? ((item as { id: string }).id as string)
        : undefined
    byName.set(name.toLowerCase(), { id, name })
  }
  return [...byName.values()]
}

function normalizeRoleDetails(value: unknown): RoleDetails {
  const base = normalizeRoleSummary(value) || { id: '', name: '', description: '' }
  const permissions =
    value && typeof value === 'object'
      ? normalizePermissions((value as { permissions?: unknown }).permissions)
      : []
  return { ...base, permissions }
}

function resetForm() {
  form.name = ''
  form.description = ''
  form.permissions = []
  originalSnapshot.value = null
  selectedRoleId.value = null
  dialogError.value = ''
}

function snapshotForm(): RoleFormSnapshot {
  return {
    name: form.name.trim(),
    description: form.description.trim(),
    permissions: [...new Set(form.permissions.map((item) => item.trim()).filter(Boolean))].sort(),
  }
}

function applyRoleDetails(role: RoleDetails) {
  form.name = role.name
  form.description = role.description
  form.permissions = role.permissions.map((permission) => permission.name)
  originalSnapshot.value = snapshotForm()
}

async function loadRoles() {
  if (!canManageRoles.value) return

  rolesLoading.value = true
  pageError.value = ''
  try {
    const payload = await requestJson<{ roles?: unknown[] }>(
      '/admin/roles',
      { method: 'GET' },
      'Unable to load roles.',
    )
    roles.value = (payload.roles || []).map(normalizeRoleSummary).filter((entry): entry is RoleSummary => !!entry)
    const visibleRoleIds = new Set(roles.value.map((role) => role.id))
    selectedRoles.value = selectedRoles.value.filter((role) => visibleRoleIds.has(role.id))
  } catch (error) {
    pageError.value = error instanceof Error ? error.message : 'Unable to load roles.'
    roles.value = []
    selectedRoles.value = []
  } finally {
    rolesLoading.value = false
  }
}

async function loadPermissionsCatalog() {
  try {
    await catalogStore.fetchPermissions()
  } catch (error) {
    dialogError.value =
      error instanceof Error ? error.message : 'Unable to load permissions catalog.'
  }
}

async function loadRoleDetails(roleId: string) {
  roleDetailsLoading.value = true
  dialogError.value = ''
  try {
    const payload = await requestJson<unknown>(
      `/admin/roles/${roleId}`,
      { method: 'GET' },
      'Unable to load role details.',
    )
    const details = normalizeRoleDetails(payload)
    if (!details.id) {
      throw new Error('Role details response is missing required fields.')
    }
    applyRoleDetails(details)
  } catch (error) {
    dialogError.value = error instanceof Error ? error.message : 'Unable to load role details.'
  } finally {
    roleDetailsLoading.value = false
  }
}

async function openCreateDialog() {
  if (!canManageRoles.value) return
  resetForm()
  dialogMode.value = 'create'
  dialogVisible.value = true
  await loadPermissionsCatalog()
}

async function openRoleDialog(mode: DialogMode, role: RoleSummary) {
  if (!canManageRoles.value || mode === 'create') return

  resetForm()
  selectedRoleId.value = role.id
  dialogMode.value = mode
  dialogVisible.value = true
  await Promise.all([loadPermissionsCatalog(), loadRoleDetails(role.id)])
}

async function submitRole() {
  if (!canSubmit.value) return

  dialogBusy.value = true
  dialogError.value = ''

  try {
    if (dialogMode.value === 'create') {
      const payload: Record<string, unknown> = {
        name: form.name.trim(),
        description: form.description.trim(),
      }
      if (form.permissions.length > 0) payload.permissions = [...form.permissions]

      await requestJson(
        '/admin/roles',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        },
        'Unable to create role.',
      )
      toast.add({ severity: 'success', summary: 'Role created', life: 2500 })
    }

    if (dialogMode.value === 'edit') {
      const roleId = selectedRoleId.value
      if (!roleId) throw new Error('Role identifier missing.')

      const current = snapshotForm()
      const original = originalSnapshot.value
      const payload: Record<string, unknown> = {}

      if (!original || current.name !== original.name) payload.name = current.name
      if (!original || current.description !== original.description) {
        payload.description = current.description
      }

      const currentPermissions = current.permissions.join('|')
      const originalPermissions = original ? original.permissions.join('|') : null
      if (!original || currentPermissions !== originalPermissions) {
        payload.permissions = current.permissions
      }

      if (Object.keys(payload).length > 0) {
        await requestJson(
          `/admin/roles/${roleId}`,
          {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
          },
          'Unable to update role.',
        )
        toast.add({ severity: 'success', summary: 'Role updated', life: 2500 })
      } else {
        toast.add({ severity: 'info', summary: 'No changes to save', life: 2200 })
      }
    }

    dialogVisible.value = false
    catalogStore.invalidateRoles()
    catalogStore.invalidatePermissions()
    await loadRoles()
  } catch (error) {
    dialogError.value = error instanceof Error ? error.message : 'Unable to save role.'
    toast.add({ severity: 'error', summary: 'Save failed', detail: dialogError.value, life: 3500 })
  } finally {
    dialogBusy.value = false
  }
}

function deleteRole(role: RoleSummary) {
  if (!canManageRoles.value) return

  confirm.require({
    message: `Delete role "${role.name}"?`,
    header: 'Delete role',
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Delete',
    rejectLabel: 'Cancel',
    accept: async () => {
      pageError.value = ''

      try {
        await requestJson(
          `/admin/roles/${role.id}`,
          { method: 'DELETE' },
          'Unable to delete role.',
        )
        toast.add({ severity: 'success', summary: 'Role deleted', life: 2500 })
        catalogStore.invalidateRoles()
        await loadRoles()
      } catch (error) {
        pageError.value = error instanceof Error ? error.message : 'Unable to delete role.'
        toast.add({ severity: 'error', summary: 'Delete failed', detail: pageError.value, life: 3500 })
      }
    },
  })
}

function openRowActionMenu(event: Event, role: RoleSummary) {
  selectedMenuRole.value = role
  rowActionMenu.value?.toggle(event)
}

void loadRoles()
</script>

<template>
  <PageSectionHeader title="Roles">
    <template #actions>
      <PrimaryActionDropdown
        label="New Role"
        :disabled="!canManageRoles"
        title="New Role options"
        importLabel="Import roles"
        exportLabel="Export roles"
        refreshLabel="Refresh roles"
        :importDisabled="true"
        :exportDisabled="true"
        :refreshDisabled="!canManageRoles"
        @click="openCreateDialog"
        @refresh="loadRoles"
      />
    </template>
  </PageSectionHeader>

  <Message v-if="!canManageRoles" severity="warn" size="small" variant="simple">
    You can view roles, but <code>roles.write</code> or <code>roles.update</code> is required for changes.
  </Message>
  <Message v-if="pageError" severity="error" size="small" variant="simple">
    {{ pageError }}
  </Message>

  <Card class="roles-card">
    <template #content>
      <DataTable
        ref="dt"
        v-model:selection="selectedRoles"
        v-model:filters="filters"
        :value="roles"
        dataKey="id"
        paginator
        :rows="10"
        :rowsPerPageOptions="[10, 20, 50]"
        paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        currentPageReportTemplate="Showing {first} to {last} of {totalRecords}"
        tableStyle="min-width: 32rem"
        :loading="rolesLoading"
        filterDisplay="menu"
        :globalFilterFields="['name', 'description']"
        removableSort
        stripedRows
        rowHover
      >
        <template #header>
          <DataTableToolbar
            v-model:searchQuery="filters.global.value"
            v-model:selectedColumns="selectedColumns"
            :columns="columns"
            searchPlaceholder="Search roles..."
          />
        </template>
        <Column selectionMode="multiple" style="width: 3rem" />
        <Column 
          v-for="col of selectedColumns" 
          :key="col.field" 
          :field="col.field" 
          :header="col.header" 
          sortable 
          :showFilterOperator="false" 
          :showAddButton="false" 
        />
        <Column header="" style="width: 3.5rem">
          <template #body="{ data }">
            <div class="actions">
              <RowActionButton
                :disabled="!canManageRoles"
                ariaLabel="Open role actions"
                title="Open role actions"
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
    :style="{ width: 'min(42rem, 96vw)' }"
    @hide="resetForm"
  >
    <div class="dialog-content">
      <Message v-if="dialogError" severity="error" size="small" variant="simple">
        {{ dialogError }}
      </Message>

      <p v-if="roleDetailsLoading" class="loading-copy">Loading role details...</p>

      <template v-else>
        <label class="field">
          <span>Name</span>
          <InputText v-model="form.name" :disabled="dialogMode === 'view'" />
        </label>

        <label class="field">
          <span>Description</span>
          <InputText v-model="form.description" :disabled="dialogMode === 'view'" />
        </label>

        <label class="field permissions-field">
          <span>Permissions</span>
          <MultiSelect
            v-model="form.permissions"
            class="role-permissions-select"
            :options="permissionsCatalog"
            optionLabel="name"
            optionValue="name"
            :maxSelectedLabels="1"
            selectedItemsLabel="{0} selected"
            filter
            fluid
            :loading="permissionsLoading"
            :disabled="dialogMode === 'view'"
            placeholder="Select permissions"
          />
          <div v-if="selectedPermissionsPreview.length > 0" class="permissions-preview">
            <p class="permissions-preview-meta">{{ selectedPermissionsPreview.length }} selected permissions</p>
            <div class="permissions-preview-list">
              <span
                v-for="permission in selectedPermissionsPreview"
                :key="permission"
                class="permissions-preview-chip"
              >
                {{ permission }}
              </span>
            </div>
          </div>
        </label>
      </template>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button
          label="Close"
          :class="{ 'create-cancel-btn': dialogMode === 'create' }"
          :text="dialogMode !== 'create'"
          @click="dialogVisible = false"
        />
        <Button
          v-if="dialogMode !== 'view'"
          :label="dialogMode === 'create' ? 'Create role' : 'Save changes'"
          :loading="dialogBusy"
          :disabled="!canSubmit"
          @click="submitRole"
        />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.roles-card {
  margin-bottom: 0.9rem;
}

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  flex-wrap: wrap;
}

.dialog-content {
  display: grid;
  gap: 1rem;
  min-width: 0;
}

.field {
  display: grid;
  gap: 0.38rem;
  min-width: 0;
}

.field span {
  font-size: 0.825rem;
  color: #374151;
  font-weight: 600;
  letter-spacing: -0.005em;
}

.loading-copy {
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.permissions-field {
  min-width: 0;
}

.permissions-preview {
  display: grid;
  gap: 0.35rem;
  min-width: 0;
}

.permissions-preview-meta {
  margin: 0;
  color: #6b7280;
  font-size: 0.77rem;
}

.permissions-preview-list {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 0.35rem;
  max-height: 7.4rem;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 0.42rem;
  border: 1px solid #e5e7eb;
  border-radius: var(--radius-sm, 6px);
  background: #ffffff;
}

.permissions-preview-chip {
  display: inline-flex;
  align-items: center;
  max-width: 100%;
  min-height: 1.25rem;
  padding: 0.14rem 0.48rem;
  border-radius: var(--radius-full, 999px);
  font-size: 0.74rem;
  line-height: 1.2;
  color: var(--p-primary-700);
  background: color-mix(in srgb, var(--p-primary-color) 10%, #ffffff);
  border: 1px solid color-mix(in srgb, var(--p-primary-color) 20%, #ffffff);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dialog-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.5rem;
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



:deep(.p-dialog-content) {
  overflow-x: hidden;
}

:deep(.role-permissions-select) {
  min-width: 0;
  max-width: 100%;
}

:deep(.role-permissions-select .p-multiselect-label-container) {
  min-width: 0;
}

:deep(.role-permissions-select .p-multiselect-label) {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

</style>
