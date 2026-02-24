<script setup lang="ts">
import { computed, onUnmounted, reactive, ref, watch, type Component } from 'vue'
import { FilterMatchMode, FilterOperator } from '@primevue/core/api'
import {
  Ban,
  MapPinPlus,
  Copy,
  SquareAsterisk,
  Link2,
  MoreHorizontal,
  Pencil,
  Shield,
  UserPlus,
} from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Message from 'primevue/message'
import MultiSelect from 'primevue/multiselect'
import Password from 'primevue/password'
import type { MenuItem } from 'primevue/menuitem'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import { useCatalogStore } from '@/stores/catalog'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import PrimaryActionDropdown from '@/components/ui/PrimaryActionDropdown.vue'
import DataTableToolbar from '@/components/ui/DataTableToolbar.vue'

type AdminUser = {
  id: string
  email: string
  fullName: string
  isActive: boolean
  createdAt: string
  roles: string[]
  sites: string[]
}

type RoleOption = {
  id: string
  name: string
  description: string
}

type SiteOption = {
  id: string
  name: string
}

type ActionMenuItem = MenuItem & {
  lucideIcon: Component
  tone?: 'danger'
}

const authStore = useAuthStore()
const catalogStore = useCatalogStore()
authStore.hydrate()

const confirm = useConfirm()
const toast = useToast()

const canReadUsers = computed(() =>
  authStore.canAccess([
    'users.*',
    'users.read',
    'users.write',
    'users.update',
    'users.delete',
    'users.manage',
  ]),
)
const canManageUsers = computed(() =>
  authStore.canAccess(['users.*', 'users.write', 'users.update', 'users.delete', 'users.manage']),
)
const canDisableUsers = computed(() =>
  authStore.canAccess(['users.*', 'users.disable', 'users.manage']),
)
const canOpenUserRowActions = computed(() => canManageUsers.value || canDisableUsers.value)

const users = ref<AdminUser[]>([])
const selectedUsers = ref<AdminUser[]>([])
const usersLoading = ref(false)
const usersError = ref('')
const totalUsers = ref<number | null>(null)
const limit = ref(20)
const pageIndex = ref(0)
const cursorByPage = ref<Array<string | null>>([null])
const nextCursor = ref<string | null>(null)
const hasNextPage = ref(false)
const filters = ref({
  global: { value: null as string | null, matchMode: FilterMatchMode.CONTAINS },
  role: { operator: FilterOperator.AND, constraints: [{ value: null as string[] | null, matchMode: FilterMatchMode.IN }] },
  site: { operator: FilterOperator.AND, constraints: [{ value: null as string[] | null, matchMode: FilterMatchMode.IN }] },
})

const latestResetToken = ref('')
const latestResetEmail = ref('')
const latestResetLink = ref('')

const rolesCatalog = computed(() => catalogStore.roles)
const rolesLoading = computed(() => catalogStore.rolesLoading)
const sitesCatalog = computed(() => catalogStore.sites)
const sitesLoading = computed(() => catalogStore.sitesLoading)
const roleAssignmentsByUser = ref<Record<string, string[]>>({})
const siteAssignmentsByUser = ref<Record<string, string[]>>({})

const createDialogVisible = ref(false)
const createDialogBusy = ref(false)
const createDialogError = ref('')
const createForm = reactive({
  email: '',
  full_name: '',
  password: '',
  roles: [] as string[],
  sitesText: '',
})

const editDialogVisible = ref(false)
const editDialogBusy = ref(false)
const editDialogError = ref('')
const editTargetUser = ref<AdminUser | null>(null)
const editForm = reactive({
  full_name: '',
  is_active: true,
})

const rolesDialogVisible = ref(false)
const rolesDialogBusy = ref(false)
const rolesDialogError = ref('')
const rolesTargetUser = ref<AdminUser | null>(null)
const rolesForm = reactive({
  roles: [] as string[],
})

const sitesDialogVisible = ref(false)
const sitesDialogBusy = ref(false)
const sitesDialogLoading = ref(false)
const sitesDialogError = ref('')
const sitesTargetUser = ref<AdminUser | null>(null)
const sitesForm = reactive({
  siteIds: [] as string[],
})

const resetInFlightUserId = ref<string | null>(null)
const statusToggleInFlightUserId = ref<string | null>(null)
const rowActionMenu = ref()
const selectedMenuUser = ref<AdminUser | null>(null)
const dt = ref()

const columns = ref([
  { field: 'fullName', header: 'Name' },
  { field: 'email', header: 'Email' },
  { field: 'roles', header: 'Roles' },
  { field: 'sites', header: 'Sites' },
  { field: 'createdAt', header: 'Created' },
  { field: 'status', header: 'Status' }
])
const selectedColumns = ref([...columns.value])

const paginatorFirst = computed(() => pageIndex.value * limit.value)
const paginatorTotalRecords = computed(() => {
  if (typeof totalUsers.value === 'number' && Number.isFinite(totalUsers.value)) return totalUsers.value
  if (hasNextPage.value) return paginatorFirst.value + limit.value + 1
  return paginatorFirst.value + users.value.length
})
const roleFilterOptions = computed(() =>
  rolesCatalog.value.map((role) => role.name).sort((a, b) => a.localeCompare(b)),
)
const siteFilterOptions = computed(() =>
  sitesCatalog.value
    .map((site) => ({ label: site.name, value: site.id }))
    .sort((a, b) => a.label.localeCompare(b.label)),
)
const siteOptionsForSelect = computed(() => {
  const options = sitesCatalog.value
    .map((site) => ({ label: site.name, value: site.id }))
    .sort((a, b) => a.label.localeCompare(b.label))

  const seen = new Set(options.map((option) => option.value.toLowerCase()))
  for (const siteId of sitesForm.siteIds) {
    const normalized = siteId.trim()
    if (!normalized) continue
    const dedupeKey = normalized.toLowerCase()
    if (seen.has(dedupeKey)) continue
    seen.add(dedupeKey)
    options.push({ label: normalized, value: normalized })
  }

  return options
})
const rowActionItems = computed<ActionMenuItem[]>(() => [
  {
    label: 'Edit user',
    lucideIcon: Pencil,
    disabled: !canManageUsers.value,
    command: () => {
      if (!selectedMenuUser.value) return
      openEditDialog(selectedMenuUser.value)
    },
  },
  {
    label: 'Assign roles',
    lucideIcon: Shield,
    disabled: !canManageUsers.value,
    command: () => {
      if (!selectedMenuUser.value) return
      void openRolesDialog(selectedMenuUser.value)
    },
  },
  {
    label: 'Assign sites',
    lucideIcon: MapPinPlus,
    disabled: !canManageUsers.value,
    command: () => {
      if (!selectedMenuUser.value) return
      void openSitesDialog(selectedMenuUser.value)
    },
  },
  {
    label: 'Reset password',
    lucideIcon: SquareAsterisk,
    disabled: !canManageUsers.value,
    command: () => {
      if (!selectedMenuUser.value) return
      triggerPasswordReset(selectedMenuUser.value)
    },
  },
  {
    label: selectedMenuUser.value?.isActive ? 'Disable user' : 'Enable user',
    lucideIcon: selectedMenuUser.value?.isActive ? Ban : UserPlus,
    tone: 'danger',
    disabled: !canDisableUsers.value || Boolean(statusToggleInFlightUserId.value),
    command: () => {
      if (!selectedMenuUser.value) return
      toggleUserStatus(selectedMenuUser.value)
    },
  },
])

async function requestJson<T>(path: string, init: RequestInit, fallbackError: string) {
  return authStore.authenticatedRequestJson<T>(path, init, fallbackError)
}

function parseCsvList(value: string) {
  return [...new Set(value.split(/[\n,]/).map((entry) => entry.trim()).filter(Boolean))]
}

function buildResetLink(token: string) {
  const encoded = encodeURIComponent(token)
  const path = `/reset-password?token=${encoded}`
  if (typeof window === 'undefined') return path
  return `${window.location.origin}${path}`
}

function normalizeEntityList(value: unknown, itemKeys: readonly string[]) {
  if (!Array.isArray(value)) return []

  const normalized: string[] = []
  const seen = new Set<string>()

  for (const item of value) {
    let entry = ''

    if (typeof item === 'string') {
      entry = item.trim()
    } else if (item && typeof item === 'object') {
      const record = item as Record<string, unknown>
      for (const key of itemKeys) {
        const candidate = record[key]
        if (typeof candidate === 'string' && candidate.trim()) {
          entry = candidate.trim()
          break
        }
        if (typeof candidate === 'number' && Number.isFinite(candidate)) {
          entry = String(candidate)
          break
        }
      }
    }

    if (!entry) continue
    const dedupeKey = entry.toLowerCase()
    if (seen.has(dedupeKey)) continue
    seen.add(dedupeKey)
    normalized.push(entry)
  }

  return normalized
}

function mergeUniqueLabels(existing: string[], incoming: string[]) {
  const merged: string[] = []
  const seen = new Set<string>()
  for (const value of [...existing, ...incoming]) {
    const normalized = value.trim()
    if (!normalized) continue
    const key = normalized.toLowerCase()
    if (seen.has(key)) continue
    seen.add(key)
    merged.push(normalized)
  }
  return merged
}

function mergeUsersById(entries: AdminUser[]) {
  const byId = new Map<string, AdminUser>()

  for (const entry of entries) {
    const existing = byId.get(entry.id)
    if (!existing) {
      byId.set(entry.id, {
        ...entry,
        roles: [...entry.roles],
        sites: [...entry.sites],
      })
      continue
    }

    existing.roles = mergeUniqueLabels(existing.roles, entry.roles)
    existing.sites = mergeUniqueLabels(existing.sites, entry.sites)
    if (!existing.createdAt && entry.createdAt) existing.createdAt = entry.createdAt
    if (!existing.fullName || existing.fullName === 'Unnamed user') existing.fullName = entry.fullName
    if (!existing.email && entry.email) existing.email = entry.email
    existing.isActive = entry.isActive
  }

  return [...byId.values()]
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

function parseErrorMessage(payload: unknown, fallback: string) {
  if (typeof payload === 'string' && payload.trim()) return payload.trim()
  if (payload && typeof payload === 'object') {
    const record = payload as Record<string, unknown>
    for (const key of ['message', 'error', 'detail']) {
      const value = record[key]
      if (typeof value === 'string' && value.trim()) return value.trim()
    }
  }
  return fallback
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

const ROLE_ITEM_KEYS = ['name', 'role', 'slug', 'code', 'id', 'label'] as const

function normalizeRoleValues(payload: unknown) {
  if (Array.isArray(payload)) {
    return normalizeEntityList(payload, ROLE_ITEM_KEYS)
  }

  if (typeof payload === 'string') {
    return parseCsvList(payload)
  }

  if (!payload || typeof payload !== 'object') return []
  const record = payload as Record<string, unknown>

  const roleSource = Array.isArray(record.roles)
    ? record.roles
    : Array.isArray(record.user_roles)
      ? record.user_roles
      : Array.isArray(record.userRoles)
        ? record.userRoles
        : Array.isArray(record.authorities)
          ? record.authorities
          : []

  const rolesFromArray = normalizeEntityList(roleSource, ROLE_ITEM_KEYS)
  const rolesFromDelimited = [record.roles, record.user_roles, record.userRoles, record.authorities].flatMap((value) =>
    typeof value === 'string' ? parseCsvList(value) : [],
  )
  const fallbackRoles = normalizeEntityList(
    [record.role, record.primary_role, record.primaryRole, record.accountRole, record.account_role],
    ROLE_ITEM_KEYS,
  )

  return mergeUniqueLabels(rolesFromArray, mergeUniqueLabels(rolesFromDelimited, fallbackRoles))
}

function normalizeUser(value: unknown): AdminUser | null {
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>

  const id = typeof record.id === 'string' ? record.id : null
  const email = typeof record.email === 'string' ? record.email.trim() : ''
  const fullName =
    typeof record.full_name === 'string'
      ? record.full_name.trim()
      : typeof record.fullName === 'string'
        ? record.fullName.trim()
        : typeof record.name === 'string'
          ? record.name.trim()
          : ''
  const createdAt =
    typeof record.created_at === 'string'
      ? record.created_at
      : typeof record.createdAt === 'string'
        ? record.createdAt
        : ''
  const isActive = resolveActiveState(record, true)

  const roles = normalizeRoleValues(record)

  const siteSource = Array.isArray(record.sites)
    ? record.sites
    : Array.isArray(record.sites)
      ? record.sites
      : Array.isArray(record.locations)
        ? record.locations
        : Array.isArray(record.facilities)
          ? record.facilities
          : []
  const sites = normalizeEntityList(siteSource, ['name', 'label', 'code', 'slug', 'id'])
  if (sites.length === 0) {
    const fallbackSite =
      typeof record.site === 'string'
        ? record.site.trim()
        : typeof record.site === 'string'
          ? record.site.trim()
          : ''
    if (fallbackSite) sites.push(fallbackSite)
  }

  if (!id || !email) return null

  return {
    id,
    email,
    fullName: fullName || 'Unnamed user',
    isActive,
    createdAt,
    roles,
    sites,
  }
}

function normalizeRole(value: unknown): RoleOption | null {
  if (typeof value === 'string') {
    const name = value.trim()
    if (!name) return null
    return { id: name, name, description: '' }
  }

  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  const idCandidate = record.id ?? record.role_id ?? record.roleId ?? record.slug ?? record.code ?? record.name
  const id =
    typeof idCandidate === 'string'
      ? idCandidate.trim()
      : typeof idCandidate === 'number' && Number.isFinite(idCandidate)
        ? String(idCandidate)
        : ''

  const nameCandidates = [record.name, record.label, record.role, record.slug, record.code]
  let name = ''
  for (const candidate of nameCandidates) {
    if (typeof candidate !== 'string') continue
    const normalized = candidate.trim()
    if (!normalized) continue
    name = normalized
    break
  }
  if (!name) name = id

  const description = typeof record.description === 'string' ? record.description.trim() : ''
  if (!id || !name) return null
  return { id, name, description }
}

function normalizeSite(value: unknown): SiteOption | null {
  if (typeof value === 'string') {
    const normalized = value.trim()
    if (!normalized) return null
    return { id: normalized, name: normalized }
  }

  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  const idCandidate = record.id ?? record.site_id ?? record.siteId ?? record.uuid ?? record.code
  const id =
    typeof idCandidate === 'string'
      ? idCandidate.trim()
      : typeof idCandidate === 'number' && Number.isFinite(idCandidate)
        ? String(idCandidate)
        : ''

  const nameCandidates = [record.name, record.label, record.code, record.slug]
  let name = ''
  for (const candidate of nameCandidates) {
    if (typeof candidate !== 'string') continue
    const normalized = candidate.trim()
    if (!normalized) continue
    name = normalized
    break
  }
  if (!name) name = id

  if (!id) return null
  return { id, name }
}

function extractCollection(payload: unknown, keys: readonly string[]) {
  const queue: unknown[] = [payload]
  const seen = new Set<object>()

  while (queue.length > 0) {
    const current = queue.shift()
    if (Array.isArray(current)) return current
    if (!current || typeof current !== 'object') continue
    if (seen.has(current)) continue
    seen.add(current)

    const record = current as Record<string, unknown>
    const nested: unknown[] = []

    for (const key of keys) {
      const candidate = record[key]
      if (Array.isArray(candidate)) return candidate
      if (candidate && typeof candidate === 'object') nested.push(candidate)
    }

    for (const key of ['data', 'result', 'payload', 'items', 'results', 'rows', 'records']) {
      const nested = record[key]
      if (Array.isArray(nested)) return nested
      if (nested && typeof nested === 'object') queue.push(nested)
    }

    for (const value of Object.values(record)) {
      if (Array.isArray(value)) return value
      if (value && typeof value === 'object') nested.push(value)
    }

    for (const value of nested) {
      queue.push(value)
    }
  }

  return []
}

function normalizeSiteValues(payload: unknown) {
  const source = extractCollection(payload, ['sites', 'sites', 'site_ids', 'siteIds', 'assignments'])
  const normalized = normalizeEntityList(source, [
    'id',
    'site_id',
    'siteId',
    'site_id',
    'siteId',
    'name',
    'label',
    'code',
    'slug',
  ])
  if (normalized.length > 0) return normalized

  if (!payload || typeof payload !== 'object') return []
  const record = payload as Record<string, unknown>
  return normalizeEntityList(
    [record.site, record.site_id, record.siteId, record.site, record.site_id, record.siteId],
    ['id', 'site_id', 'siteId', 'site_id', 'siteId', 'name', 'label', 'code', 'slug'],
  )
}

function resolveSiteSelectionIds(values: string[]) {
  const idsByLower = new Map(sitesCatalog.value.map((site) => [site.id.toLowerCase(), site.id]))
  const idsByName = new Map(sitesCatalog.value.map((site) => [site.name.toLowerCase(), site.id]))
  const resolved: string[] = []
  const seen = new Set<string>()

  for (const value of values) {
    const normalized = value.trim()
    if (!normalized) continue
    const lookup = normalized.toLowerCase()
    const resolvedValue = idsByLower.get(lookup) || idsByName.get(lookup) || normalized
    const dedupeKey = resolvedValue.toLowerCase()
    if (seen.has(dedupeKey)) continue
    seen.add(dedupeKey)
    resolved.push(resolvedValue)
  }

  return resolved
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

function roleOptionsForSelect() {
  return rolesCatalog.value.map((role) => role.name)
}

function normalizeSelection(values: string[]) {
  return [...new Set(values.map((value) => value.trim()).filter(Boolean))].sort((a, b) =>
    a.localeCompare(b),
  )
}

function resetUsersPagination() {
  pageIndex.value = 0
  cursorByPage.value = [null]
  nextCursor.value = null
  hasNextPage.value = false
}

function resetCreateForm() {
  createForm.email = ''
  createForm.full_name = ''
  createForm.password = ''
  createForm.roles = []
  createForm.sitesText = ''
  createDialogError.value = ''
}

function resetEditForm() {
  editTargetUser.value = null
  editForm.full_name = ''
  editForm.is_active = true
  editDialogError.value = ''
}

function resetRolesForm() {
  rolesTargetUser.value = null
  rolesForm.roles = []
  rolesDialogError.value = ''
}

function resetSitesForm() {
  sitesTargetUser.value = null
  sitesForm.siteIds = []
  sitesDialogError.value = ''
}

async function loadRolesCatalog(options?: { silent?: boolean; force?: boolean }) {
  if (!canReadUsers.value) return
  try {
    await catalogStore.fetchRoles({ force: options?.force })
  } catch (error) {
    if (!options?.silent) {
      const message = error instanceof Error ? error.message : 'Unable to load roles catalog.'
      toast.add({ severity: 'error', summary: 'Roles unavailable', detail: message, life: 3500 })
    }
  }
}

async function loadSitesCatalog(options?: { silent?: boolean; force?: boolean }) {
  if (!canReadUsers.value) return
  try {
    await catalogStore.fetchSites({ force: options?.force })
  } catch (error) {
    if (!options?.silent) {
      const message = error instanceof Error ? error.message : 'Unable to load sites catalog.'
      toast.add({ severity: 'error', summary: 'Sites unavailable', detail: message, life: 3500 })
    }
  }
}

async function loadUsers() {
  if (!canReadUsers.value) return

  usersLoading.value = true
  usersError.value = ''
  try {
    const requestedLimit = Math.max(1, Math.min(limit.value, 100))
    if (requestedLimit !== limit.value) limit.value = requestedLimit
    const currentCursor = cursorByPage.value[pageIndex.value] ?? null
    const params = new URLSearchParams({
      limit: String(requestedLimit),
      include_total: 'true',
    })
    if (currentCursor) params.set('cursor', currentCursor)
    const searchValue = filters.value.global.value?.trim()
    if (searchValue) {
      params.set('search', searchValue)
    }
    const roleConstraint = filters.value.role.constraints[0]
    const roleValues = roleConstraint?.value
    if (roleValues && roleValues.length > 0) {
      for (const roleName of roleValues) {
        params.append('role', roleName)
      }
    }
    const siteConstraint = filters.value.site.constraints[0]
    const siteValues = siteConstraint?.value
    if (siteValues && siteValues.length > 0) {
      for (const siteId of siteValues) {
        params.append('site_id', siteId)
      }
    }
    const payload = await requestJson<{
      users?: unknown[]
      total?: unknown
      pagination?: {
        next_cursor?: unknown
        has_next?: unknown
        limit?: unknown
      }
    }>(
      `/admin/users?${params.toString()}`,
      { method: 'GET' },
      'Unable to load users.',
    )

    const normalizedUsers = (payload.users || [])
      .map(normalizeUser)
      .filter((entry): entry is AdminUser => !!entry)
    users.value = mergeUsersById(normalizedUsers)
    const visibleUserIds = new Set(users.value.map((user) => user.id))
    selectedUsers.value = selectedUsers.value.filter((user) => visibleUserIds.has(user.id))
    for (const user of users.value) {
      roleAssignmentsByUser.value[user.id] = mergeUniqueLabels(
        roleAssignmentsByUser.value[user.id] || [],
        user.roles,
      )
    }

    const pagination =
      payload.pagination && typeof payload.pagination === 'object'
        ? (payload.pagination as Record<string, unknown>)
        : {}
    const apiNextCursor =
      typeof pagination.next_cursor === 'string' && pagination.next_cursor.trim()
        ? pagination.next_cursor
        : null
    const apiHasNext = parseBooleanLike(pagination.has_next) ?? Boolean(apiNextCursor)
    hasNextPage.value = apiHasNext
    nextCursor.value = apiNextCursor

    const apiLimit = pagination.limit
    if (
      typeof apiLimit === 'number' &&
      Number.isFinite(apiLimit) &&
      apiLimit > 0 &&
      apiLimit <= 100 &&
      apiLimit !== limit.value
    ) {
      limit.value = apiLimit
    }

    cursorByPage.value = cursorByPage.value.slice(0, pageIndex.value + 1)
    if (apiHasNext && apiNextCursor) {
      cursorByPage.value[pageIndex.value + 1] = apiNextCursor
    }

    if (typeof payload.total === 'number' && Number.isFinite(payload.total) && payload.total >= 0) {
      totalUsers.value = payload.total
    } else if (pageIndex.value === 0) {
      totalUsers.value = null
    }
  } catch (error) {
    usersError.value = error instanceof Error ? error.message : 'Unable to load users.'
    users.value = []
    selectedUsers.value = []
    totalUsers.value = null
    hasNextPage.value = false
    nextCursor.value = null
  } finally {
    usersLoading.value = false
  }
}

function onPageChange(event: { page: number; first: number; rows: number }) {
  const nextRows = Math.max(1, Math.min(event.rows, 100))
  if (nextRows !== limit.value) {
    limit.value = nextRows
    resetUsersPagination()
    void loadUsers()
    return
  }

  const targetPage =
    typeof event.page === 'number' && Number.isFinite(event.page)
      ? event.page
      : Math.floor(event.first / nextRows)

  if (targetPage === pageIndex.value) return

  if (targetPage < pageIndex.value) {
    if (targetPage < 0 || targetPage >= cursorByPage.value.length) return
    if (typeof cursorByPage.value[targetPage] === 'undefined') return
    pageIndex.value = targetPage
    void loadUsers()
    return
  }

  if (targetPage === pageIndex.value + 1 && hasNextPage.value && nextCursor.value) {
    cursorByPage.value[targetPage] = nextCursor.value
    pageIndex.value = targetPage
    void loadUsers()
  }
}

function refreshUsers() {
  void loadUsers()
}

watch(() => filters.value.global.value, (val) => {
  if (!val) {
    resetUsersPagination()
    void loadUsers()
  }
})

function onSearchEnter() {
  resetUsersPagination()
  void loadUsers()
}

function onFilter() {
  resetUsersPagination()
  void loadUsers()
}

function openRowActionMenu(event: Event, user: AdminUser) {
  selectedMenuUser.value = user
  rowActionMenu.value?.toggle(event)
}

function toggleUserStatus(user: AdminUser) {
  if (!canDisableUsers.value || statusToggleInFlightUserId.value) return

  const shouldDisable = user.isActive
  const actionLabel = shouldDisable ? 'Disable user' : 'Enable user'
  const actionVerb = shouldDisable ? 'disable' : 'enable'
  const successSummary = shouldDisable ? 'User disabled' : 'User enabled'
  const caution =
    shouldDisable ? ` This will revoke all refresh tokens for ${user.fullName}.` : ''

  confirm.require({
    header: actionLabel,
    message: `${actionLabel} ${user.fullName}?${caution}`,
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: actionLabel,
    rejectLabel: 'Cancel',
    accept: async () => {
      statusToggleInFlightUserId.value = user.id
      usersError.value = ''
      try {
        await requestJson(
          `/admin/users/${user.id}/${shouldDisable ? 'disable' : 'enable'}`,
          { method: 'POST' },
          `Unable to ${actionVerb} user.`,
        )
        toast.add({ severity: 'success', summary: successSummary, life: 2500 })
        await loadUsers()
      } catch (error) {
        usersError.value = error instanceof Error ? error.message : `Unable to ${actionVerb} user.`
        toast.add({ severity: 'error', summary: `${actionLabel} failed`, detail: usersError.value, life: 3500 })
      } finally {
        statusToggleInFlightUserId.value = null
      }
    },
  })
}

async function openCreateDialog() {
  if (!canManageUsers.value) return
  resetCreateForm()
  createDialogVisible.value = true
  await loadRolesCatalog()
}

async function submitCreateUser() {
  if (!canManageUsers.value || createDialogBusy.value) return
  if (!createForm.email.trim() || !createForm.full_name.trim()) {
    createDialogError.value = 'Email and full name are required.'
    return
  }

  createDialogBusy.value = true
  createDialogError.value = ''

  try {
    const payload: Record<string, unknown> = {
      email: createForm.email.trim(),
      full_name: createForm.full_name.trim(),
    }
    const password = createForm.password.trim()
    if (password) payload.password = password
    if (createForm.roles.length > 0) payload.roles = [...createForm.roles]
    const sites = parseCsvList(createForm.sitesText)
    if (sites.length > 0) payload.sites = sites

    const response = await requestJson<{ user?: unknown; reset_token?: string }>(
      '/admin/users',
      { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) },
      'Unable to create user.',
    )

    const createdUser = normalizeUser(response.user)
    if (createdUser) {
      roleAssignmentsByUser.value[createdUser.id] = [...createForm.roles]
      siteAssignmentsByUser.value[createdUser.id] = sites
    }

    if (response.reset_token) {
      latestResetEmail.value = createForm.email.trim()
      latestResetToken.value = response.reset_token
      latestResetLink.value = buildResetLink(response.reset_token)
      toast.add({
        severity: 'warn',
        summary: 'Reset token generated',
        detail: 'Share the reset link with the invited user.',
        life: 4200,
      })
    } else {
      latestResetEmail.value = ''
      latestResetToken.value = ''
      latestResetLink.value = ''
      toast.add({ severity: 'success', summary: 'User created', life: 2500 })
    }

    createDialogVisible.value = false
    await loadUsers()
  } catch (error) {
    createDialogError.value = error instanceof Error ? error.message : 'Unable to create user.'
  } finally {
    createDialogBusy.value = false
  }
}

function openEditDialog(user: AdminUser) {
  if (!canManageUsers.value) return
  resetEditForm()
  editTargetUser.value = user
  editForm.full_name = user.fullName
  editForm.is_active = user.isActive
  editDialogVisible.value = true
}

async function submitEditUser() {
  if (!canManageUsers.value || editDialogBusy.value) return
  if (!editTargetUser.value) return
  if (!editForm.full_name.trim()) {
    editDialogError.value = 'Full name is required.'
    return
  }

  const payload: Record<string, unknown> = {}
  if (editForm.full_name.trim() !== editTargetUser.value.fullName) {
    payload.full_name = editForm.full_name.trim()
  }
  if (editForm.is_active !== editTargetUser.value.isActive) {
    payload.is_active = editForm.is_active
  }

  if (Object.keys(payload).length === 0) {
    toast.add({ severity: 'info', summary: 'No changes to save', life: 2200 })
    return
  }

  editDialogBusy.value = true
  editDialogError.value = ''
  try {
    await requestJson(
      `/admin/users/${editTargetUser.value.id}`,
      { method: 'PATCH', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) },
      'Unable to update user.',
    )
    toast.add({ severity: 'success', summary: 'User updated', life: 2500 })
    editDialogVisible.value = false
    await loadUsers()
  } catch (error) {
    editDialogError.value = error instanceof Error ? error.message : 'Unable to update user.'
  } finally {
    editDialogBusy.value = false
  }
}

async function openRolesDialog(user: AdminUser) {
  if (!canManageUsers.value) return
  resetRolesForm()
  rolesTargetUser.value = user
  rolesForm.roles = [...(roleAssignmentsByUser.value[user.id] || user.roles || [])]
  rolesDialogVisible.value = true
  await loadRolesCatalog()
}

async function openSitesDialog(user: AdminUser) {
  if (!canManageUsers.value) return
  resetSitesForm()
  sitesTargetUser.value = user
  sitesForm.siteIds = resolveSiteSelectionIds(siteAssignmentsByUser.value[user.id] || [])
  sitesDialogVisible.value = true

  await loadSitesCatalog()

  sitesDialogLoading.value = true
  sitesDialogError.value = ''
  try {
    const response = await authStore.authenticatedFetch(`/admin/users/${user.id}/sites`, { method: 'GET' })
    const payload = (await response.json().catch(() => ({}))) as unknown
    if (!response.ok) {
      if (response.status === 404 || response.status === 405) return
      throw new Error(parseErrorMessage(payload, 'Unable to load user sites.'))
    }

    const normalizedSites = resolveSiteSelectionIds(normalizeSiteValues(payload))
    siteAssignmentsByUser.value[user.id] = normalizedSites
    sitesForm.siteIds = [...normalizedSites]
  } catch (error) {
    sitesDialogError.value =
      error instanceof Error ? error.message : 'Unable to load user sites.'
  } finally {
    sitesDialogLoading.value = false
  }
}

async function submitRolesAssignment() {
  if (!canManageUsers.value || rolesDialogBusy.value) return
  if (!rolesTargetUser.value) return

  rolesDialogBusy.value = true
  rolesDialogError.value = ''
  try {
    const payload = {
      roles: [...new Set(rolesForm.roles.map((role) => role.trim()).filter(Boolean))],
    }
    const response = await requestJson<{ roles?: unknown[] }>(
      `/admin/users/${rolesTargetUser.value.id}/roles`,
      { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) },
      'Unable to update user roles.',
    )

    const normalizedRolesFromResponse = normalizeRoleValues(response)
    const normalizedRoles = mergeUniqueLabels(
      [],
      normalizedRolesFromResponse.length > 0 ? normalizedRolesFromResponse : payload.roles,
    )

    roleAssignmentsByUser.value[rolesTargetUser.value.id] = normalizedRoles
    rolesForm.roles = [...normalizedRoles]
    users.value = users.value.map((entry) =>
      entry.id === rolesTargetUser.value?.id ? { ...entry, roles: [...normalizedRoles] } : entry,
    )
    toast.add({ severity: 'success', summary: 'Roles updated', life: 2500 })
  } catch (error) {
    rolesDialogError.value =
      error instanceof Error ? error.message : 'Unable to update user roles.'
  } finally {
    rolesDialogBusy.value = false
  }
}

async function submitSitesAssignment() {
  if (!canManageUsers.value || sitesDialogBusy.value) return
  if (!sitesTargetUser.value) return

  sitesDialogBusy.value = true
  sitesDialogError.value = ''
  try {
    const sites = normalizeSelection(sitesForm.siteIds)
    const payload = { sites }
    const response = await requestJson<unknown>(
      `/admin/users/${sitesTargetUser.value.id}/sites`,
      { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) },
      'Unable to update user sites.',
    )

    const normalizedFromResponse = resolveSiteSelectionIds(normalizeSiteValues(response))
    const normalizedSites = normalizedFromResponse.length > 0 ? normalizedFromResponse : sites

    siteAssignmentsByUser.value[sitesTargetUser.value.id] = normalizedSites
    sitesForm.siteIds = [...normalizedSites]
    toast.add({ severity: 'success', summary: 'Sites updated', life: 2500 })
  } catch (error) {
    sitesDialogError.value =
      error instanceof Error ? error.message : 'Unable to update user sites.'
  } finally {
    sitesDialogBusy.value = false
  }
}

function triggerPasswordReset(user: AdminUser) {
  if (!canManageUsers.value || resetInFlightUserId.value) return

  confirm.require({
    header: 'Force password reset',
    message: `Generate a password reset token for ${user.fullName}?`,
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Generate token',
    rejectLabel: 'Cancel',
    accept: async () => {
      resetInFlightUserId.value = user.id
      usersError.value = ''
      try {
        const payload = await requestJson<{ reset_token?: string }>(
          `/admin/users/${user.id}/password-reset`,
          { method: 'POST' },
          'Unable to generate reset token.',
        )
        latestResetEmail.value = user.email
        latestResetToken.value = payload.reset_token || ''
        latestResetLink.value = payload.reset_token ? buildResetLink(payload.reset_token) : ''
        toast.add({
          severity: 'success',
          summary: 'Reset token generated',
          detail: payload.reset_token ? 'Token ready to share securely.' : 'No token returned.',
          life: 3500,
        })
      } catch (error) {
        usersError.value = error instanceof Error ? error.message : 'Unable to generate reset token.'
      } finally {
        resetInFlightUserId.value = null
      }
    },
  })
}

async function copyLatestResetLink() {
  if (!latestResetLink.value) return
  try {
    if (typeof navigator === 'undefined' || !navigator.clipboard) {
      throw new Error('Clipboard is unavailable in this browser.')
    }
    await navigator.clipboard.writeText(latestResetLink.value)
    toast.add({ severity: 'success', summary: 'Reset link copied', life: 2500 })
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to copy reset link.'
    toast.add({ severity: 'error', summary: 'Copy failed', detail: message, life: 3500 })
  }
}

void loadUsers()
void loadRolesCatalog({ silent: true })
void loadSitesCatalog({ silent: true })
</script>

<template>
  <PageSectionHeader title="Users" description="Manage users, profile data, roles, sites, and password resets">
    <template #actions>
      <PrimaryActionDropdown
        label="New User"
        :disabled="!canManageUsers"
        title="New User options"
        importLabel="Import users"
        exportLabel="Export users"
        refreshLabel="Refresh users"
        :importDisabled="true"
        :exportDisabled="true"
        :refreshDisabled="!canReadUsers"
        @click="openCreateDialog"
        @refresh="refreshUsers"
      />
    </template>
  </PageSectionHeader>

  <Message v-if="!canReadUsers" severity="warn" size="small" variant="simple">
    You need <code>users.read</code> permission to view users.
  </Message>
  <Message v-else-if="!canManageUsers && !canDisableUsers" severity="warn" size="small" variant="simple">
    You can view users, but <code>users.write</code> or <code>users.update</code> is required for changes.
  </Message>
  <Message v-else-if="!canManageUsers && canDisableUsers" severity="warn" size="small" variant="simple">
    You can view users and toggle active status, but <code>users.write</code> or <code>users.update</code> is
    required for profile, role, and site changes.
  </Message>
  <Message v-if="usersError" severity="error" size="small" variant="simple">
    {{ usersError }}
  </Message>
  <Message v-if="latestResetLink" severity="warn" size="small" variant="simple">
    <div class="reset-link-box">
      <p>
        Reset link for <strong>{{ latestResetEmail }}</strong>
      </p>
      <a :href="latestResetLink" target="_blank" rel="noopener noreferrer" class="reset-link-anchor">
        {{ latestResetLink }}
      </a>
      <div class="reset-link-actions">
        <Button text size="small" @click="copyLatestResetLink">
          <Copy :size="16" />
          <span>Copy link</span>
        </Button>
        <a :href="latestResetLink" target="_blank" rel="noopener noreferrer" class="reset-link-open">
          <Link2 :size="16" />
          <span>Open reset page</span>
        </a>
      </div>
      <p class="reset-token-inline">Token: <code>{{ latestResetToken }}</code></p>
    </div>
  </Message>

  <Card class="users-card">
    <template #content>
      <DataTable
        ref="dt"
        v-model:selection="selectedUsers"
        v-model:filters="filters"
        :value="users"
        dataKey="id"
        lazy
        paginator
        :rows="limit"
        :first="paginatorFirst"
        :totalRecords="paginatorTotalRecords"
        @page="onPageChange"
        :rowsPerPageOptions="[10, 20, 50]"
        paginatorTemplate="CurrentPageReport FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        currentPageReportTemplate="Showing {first} to {last} of {totalRecords}"
        tableStyle="min-width: 44rem"
        :loading="usersLoading"
        filterDisplay="menu"
        removableSort
        @filter="onFilter"
        stripedRows
        rowHover
      >
        <template #header>
          <DataTableToolbar
            v-model:searchQuery="filters.global.value"
            v-model:selectedColumns="selectedColumns"
            :columns="columns"
            searchPlaceholder="Search name or email..."
            @search="onSearchEnter"
          />
        </template>
        <template #empty>
          <div class="table-empty">No users match the current filters.</div>
        </template>
        <Column selectionMode="multiple" style="width: 3rem" />
        <template v-for="col of selectedColumns" :key="col.field">
          <Column v-if="col.field === 'fullName'" field="fullName" header="Name" sortable />
          <Column v-if="col.field === 'email'" field="email" header="Email" sortable />
          <Column v-if="col.field === 'roles'" header="Roles" field="roles" filterField="role" :showFilterOperator="false" :showAddButton="false" :showFilterMatchModes="false">
            <template #filter="{ filterModel }">
              <MultiSelect
                v-model="filterModel.value"
                :options="roleFilterOptions"
                placeholder="All roles"
                :loading="rolesLoading"
                :maxSelectedLabels="1"
                selectedItemsLabel="{0} selected"
                style="min-width: 12rem"
              />
            </template>
            <template #body="{ data }">
              <div class="role-tags">
                <span v-for="role in data.roles" :key="role" class="role-tag" :title="role">
                  {{ role }}
                </span>
                <span v-if="!data.roles?.length" class="cell-empty">None</span>
              </div>
            </template>
          </Column>
          <Column v-if="col.field === 'sites'" header="Sites" field="sites" filterField="site" :showFilterOperator="false" :showAddButton="false" :showFilterMatchModes="false">
            <template #filter="{ filterModel }">
              <MultiSelect
                v-model="filterModel.value"
                :options="siteFilterOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="All sites"
                :loading="sitesLoading"
                :maxSelectedLabels="1"
                selectedItemsLabel="{0} selected"
                style="min-width: 12rem"
              />
            </template>
            <template #body="{ data }">
              <span v-if="data.sites.length" class="cell-list" :title="data.sites.join(', ')">
                <span class="cell-pill">{{ data.sites[0] }}</span>
                <span v-if="data.sites.length > 1" class="cell-pill cell-pill-overflow">+{{ data.sites.length - 1 }}</span>
              </span>
              <span v-else class="cell-empty">&mdash;</span>
            </template>
          </Column>
          <Column v-if="col.field === 'createdAt'" header="Created">
            <template #body="{ data }">
              <span class="cell-date">{{ formatDate(data.createdAt) }}</span>
            </template>
          </Column>
          <Column v-if="col.field === 'status'" header="Status" style="width: 8rem">
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
                :disabled="!canOpenUserRowActions"
                :loading="resetInFlightUserId === data.id || statusToggleInFlightUserId === data.id"
                ariaLabel="Open user actions"
                title="Open user actions"
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

  <Dialog
    v-model:visible="createDialogVisible"
    modal
    header="Create user"
    :style="{ width: 'min(42rem, 96vw)' }"
    @hide="resetCreateForm"
  >
    <div class="dialog-grid">
      <Message v-if="createDialogError" severity="error" size="small" variant="simple">
        {{ createDialogError }}
      </Message>

      <label class="field">
        <span>Email</span>
        <InputText v-model="createForm.email" type="email" autocomplete="off" />
      </label>

      <label class="field">
        <span>Full name</span>
        <InputText v-model="createForm.full_name" autocomplete="off" />
      </label>

      <label class="field">
        <span>Password (optional)</span>
        <Password v-model="createForm.password" :feedback="false" toggleMask fluid autocomplete="new-password" />
      </label>

      <label class="field">
        <span>Roles (optional)</span>
        <MultiSelect
          v-model="createForm.roles"
          :options="roleOptionsForSelect()"
          display="chip"
          filter
          fluid
          :loading="rolesLoading"
          placeholder="Select role names"
        />
      </label>

      <label class="field">
        <span>Site IDs (optional)</span>
        <InputText v-model="createForm.sitesText" placeholder="uuid1, uuid2" autocomplete="off" />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" class="create-cancel-btn" @click="createDialogVisible = false" />
        <Button label="Create user" :loading="createDialogBusy" @click="submitCreateUser" />
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="editDialogVisible"
    modal
    header="Edit user"
    :style="{ width: 'min(34rem, 96vw)' }"
    @hide="resetEditForm"
  >
    <div class="dialog-grid">
      <Message v-if="editDialogError" severity="error" size="small" variant="simple">
        {{ editDialogError }}
      </Message>

      <p v-if="editTargetUser" class="dialog-subtitle">{{ editTargetUser.email }}</p>

      <label class="field">
        <span>Full name</span>
        <InputText v-model="editForm.full_name" autocomplete="off" />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="editDialogVisible = false" />
        <Button label="Save changes" :loading="editDialogBusy" @click="submitEditUser" />
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="rolesDialogVisible"
    modal
    header="Assign Roles"
    :style="{ width: 'min(38rem, 96vw)' }"
    @hide="resetRolesForm"
  >
    <div class="dialog-grid">
      <Message v-if="rolesDialogError" severity="error" size="small" variant="simple">
        {{ rolesDialogError }}
      </Message>

      <p v-if="rolesTargetUser" class="dialog-subtitle">{{ rolesTargetUser.email }}</p>
      <p class="help-text">
        Role updates replace all existing role assignments for this user.
      </p>

      <label class="field">
        <span>Roles</span>
        <MultiSelect
          v-model="rolesForm.roles"
          :options="roleOptionsForSelect()"
          display="chip"
          filter
          fluid
          :loading="rolesLoading"
          placeholder="Select role names"
        />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="rolesDialogVisible = false" />
        <Button
          label="Save roles"
          :loading="rolesDialogBusy"
          :disabled="!canManageUsers"
          @click="submitRolesAssignment"
        />
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="sitesDialogVisible"
    modal
    header="Assign Sites"
    :style="{ width: 'min(38rem, 96vw)' }"
    @hide="resetSitesForm"
  >
    <div class="dialog-grid">
      <Message v-if="sitesDialogError" severity="error" size="small" variant="simple">
        {{ sitesDialogError }}
      </Message>

      <p v-if="sitesTargetUser" class="dialog-subtitle">{{ sitesTargetUser.email }}</p>
      <p class="help-text">
        Site updates replace all existing site assignments for this user.
      </p>

      <label class="field">
        <span>Sites</span>
        <MultiSelect
          v-model="sitesForm.siteIds"
          :options="siteOptionsForSelect"
          optionLabel="label"
          optionValue="value"
          display="chip"
          filter
          fluid
          :loading="sitesLoading || sitesDialogLoading"
          placeholder="Select site assignments"
          selectedItemsLabel="{0} selected"
        />
      </label>
      <p v-if="!sitesLoading && !sitesDialogLoading && siteOptionsForSelect.length === 0" class="help-text">
        No sites are available to assign yet.
      </p>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="sitesDialogVisible = false" />
        <Button
          label="Save sites"
          :loading="sitesDialogBusy"
          :disabled="!canManageUsers"
          @click="submitSitesAssignment"
        />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.users-card {
  margin-bottom: 0.9rem;
}

.table-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.table-search {
  width: 100%;
  max-width: 17.33rem;
  height: var(--app-control-height, 2.25rem);
  min-height: var(--app-control-height, 2.25rem);
  box-sizing: border-box;
  border-radius: var(--app-control-radius, 8px);
  padding: 0 0.88rem;
  font-size: var(--app-control-font-size, 0.83rem);
  line-height: 1.2;
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  transition: border-color 0.15s, background 0.15s;
}

.table-search:focus {
  background: #ffffff;
  border-color: var(--p-primary-400, #94a3b8);
  box-shadow: 0 0 0 3px rgb(100 116 139 / 0.12);
}

.cell-date {
  font-size: 0.825rem;
  color: #6b7280;
}

/* ── Role / Site list cells ──────────────────────────── */

.cell-list {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  max-width: 100%;
}

.cell-pill {
  display: inline-flex;
  align-items: center;
  max-width: 11rem;
  padding: 0.15rem 0.5rem;
  border-radius: var(--radius-sm, 6px);
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  color: #374151;
  font-size: 0.855rem;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cell-pill-overflow {
  justify-content: center;
  min-width: 2rem;
  max-width: none;
  border-style: dashed;
  background: #ffffff;
  color: #374151;
  font-weight: 700;
}

.cell-empty {
  color: #9ca3af;
  font-size: 0.875rem;
}

/* ── Status dot + text ───────────────────────────────── */

.status-cell {
  display: inline-flex;
  align-items: center;
  gap: 0.46rem;
  padding: 0.22rem 0.58rem 0.22rem 0.46rem;
  border-radius: var(--radius-full, 999px);
  border: 1px solid transparent;
  font-size: 0.8rem;
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


.reset-link-box {
  display: grid;
  gap: 0.42rem;
}

.reset-link-box p {
  margin: 0;
}

.reset-link-anchor {
  color: var(--p-primary-700);
  font-size: 0.78rem;
  text-decoration: underline;
  overflow-wrap: anywhere;
}

.reset-link-actions {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  flex-wrap: wrap;
}

.reset-link-actions :deep(.p-button-label) {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
}

.reset-link-open {
  display: inline-flex;
  align-items: center;
  gap: 0.28rem;
  color: var(--p-primary-700);
  font-size: 0.78rem;
  text-decoration: none;
}

.reset-link-open:hover {
  text-decoration: underline;
}

.reset-token-inline {
  font-size: 0.75rem;
  color: #6b7280;
}

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  flex-wrap: wrap;
}

.table-empty {
  padding: 3rem 1.5rem;
  text-align: center;
  color: #9ca3af;
  font-size: 0.875rem;
}

.dialog-grid {
  display: grid;
  gap: 1rem;
}

.dialog-subtitle {
  margin: 0;
  font-size: 0.825rem;
  color: #6b7280;
}

.field {
  display: grid;
  gap: 0.36rem;
}

.field span {
  font-size: 0.825rem;
  color: #374151;
  font-weight: 600;
  letter-spacing: -0.005em;
}

.help-text {
  margin: 0;
  font-size: 0.78rem;
  color: #6b7280;
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

.menu-action-item-danger {
  color: #dc2626;
}

:deep(.p-menu-item-link.menu-action-item-danger:hover) {
  color: #dc2626;
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

:deep(.p-datatable-header-cell) {
  background: #ffffff;
  color: #6b7280;
  font-weight: 600;
  font-size: 0.78rem;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  border-bottom: 1px solid #e8eaed;
  padding: 0.5rem 1rem;
  white-space: nowrap;
}

:deep(.p-datatable-tbody > tr > td) {
  color: #1a1d23;
  padding: 0.5rem 1rem;
  border-bottom: 1px solid #eef0f3;
  border-color: #eef0f3;
  font-size: 0.875rem;
  white-space: nowrap;
}

.role-tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.35rem;
}

.role-tag {
  display: inline-flex;
  align-items: center;
  max-width: 12rem;
  padding: 0.15rem 0.5rem;
  border-radius: var(--radius-sm, 6px);
  font-size: 0.72rem;
  font-weight: 500;
  line-height: 1.2;
  color: #374151;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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

</style>
