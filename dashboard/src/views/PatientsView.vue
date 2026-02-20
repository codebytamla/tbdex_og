<script setup lang="ts">
import { computed, reactive, ref, type Component } from 'vue'
import {
  Ban,
  MapPinPlus,
  MoreHorizontal,
  Pencil,
  RefreshCcw,
  SlidersHorizontal,
  X,
} from 'lucide-vue-next'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Message from 'primevue/message'
import MultiSelect from 'primevue/multiselect'
import Textarea from 'primevue/textarea'
import type { MenuItem } from 'primevue/menuitem'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import CreateActionButton from '@/components/ui/CreateActionButton.vue'
import FilterPanelField from '@/components/ui/FilterPanelField.vue'
import PageSectionHeader from '@/components/ui/PageSectionHeader.vue'
import RowActionButton from '@/components/ui/RowActionButton.vue'
import ToolbarIconButton from '@/components/ui/ToolbarIconButton.vue'
import {
  TABLE_PAGE_REPORT_TEMPLATE,
  TABLE_PAGINATOR_TEMPLATE,
  TABLE_ROWS_PER_PAGE_OPTIONS,
} from '@/components/ui/tablePagination'

type PatientStatus = 'provisional' | 'active' | 'inactive' | 'deceased'

type SiteSummary = {
  id: string
  name: string
  code: string
}

type Patient = {
  id: string
  fullName: string
  dateOfBirth: string
  isDobEstimated: boolean
  sex: string
  phone: string
  nationalId: string
  mrn: string
  address: string
  village: string
  origin: string
  metadata: unknown
  primarySiteId: string
  primarySiteName: string
  sites: SiteSummary[]
  status: PatientStatus
  isActive: boolean
  createdById: string
  createdByName: string
  createdAt: string
  updatedAt: string
}

type SiteOption = {
  id: string
  name: string
  code: string
}

type PatientColumnKey =
  | 'fullName'
  | 'ageSex'
  | 'village'
  | 'createdBy'
  | 'createdAt'
  | 'updatedAt'
  | 'status'

type PatientColumnOption = {
  label: string
  value: PatientColumnKey
  minWidthRem: number
}

type ActionMenuItem = MenuItem & {
  lucideIcon: Component
  tone?: 'danger'
}

const authStore = useAuthStore()
authStore.hydrate()

const confirm = useConfirm()
const toast = useToast()

const canReadPatients = computed(() =>
  authStore.canAccess(['patients.*', 'patients.manage', 'patients.read', 'patients.write']),
)
const canCreatePatients = computed(() =>
  authStore.canAccess(['patients.*', 'patients.manage', 'patients.create', 'patients.write']),
)
const canUpdatePatients = computed(() =>
  authStore.canAccess(['patients.*', 'patients.manage', 'patients.update', 'patients.write']),
)
const canAddPatientSites = computed(() =>
  authStore.canAccess([
    'patients.*',
    'patients.manage',
    'patients.add_site',
    'patients.update',
    'patients.write',
  ]),
)
const currentUserId = computed(() =>
  typeof authStore.profile?.id === 'string' ? authStore.profile.id : '',
)

const patients = ref<Patient[]>([])
const patientsLoading = ref(false)
const patientsError = ref('')
const totalPatients = ref<number | null>(null)
const limit = ref(20)
const pageIndex = ref(0)
const cursorByPage = ref<Array<string | null>>([null])
const nextCursor = ref<string | null>(null)
const hasNextPage = ref(false)

const searchText = ref('')
const appliedSearch = ref('')
const siteFilter = ref<string[]>([])
const appliedSiteFilter = ref<string[]>([])
const filtersExpanded = ref(false)

const sitesCatalog = ref<SiteOption[]>([])
const sitesLoading = ref(false)

const createDialogVisible = ref(false)
const createDialogBusy = ref(false)
const createDialogError = ref('')
const createForm = reactive({
  full_name: '',
  primary_site_id: '',
  date_of_birth: '',
  is_dob_estimated: false,
  status: 'provisional' as PatientStatus,
  sex: '',
  phone: '',
  national_id: '',
  mrn: '',
  address: '',
  village: '',
  origin: '',
  metadata_text: '',
})

const editDialogVisible = ref(false)
const editDialogBusy = ref(false)
const editDialogError = ref('')
const editTargetPatient = ref<Patient | null>(null)
const editForm = reactive({
  full_name: '',
  date_of_birth: '',
  is_dob_estimated: false,
  status: 'provisional' as PatientStatus,
  sex: '',
  phone: '',
  national_id: '',
  mrn: '',
  address: '',
  village: '',
  origin: '',
  metadata_text: '',
})

const sitesDialogVisible = ref(false)
const sitesDialogBusy = ref(false)
const sitesDialogLoading = ref(false)
const sitesDialogError = ref('')
const sitesTargetPatient = ref<Patient | null>(null)
const sitesForm = reactive({
  site_id: '',
})

const patientColumnOptions: PatientColumnOption[] = [
  { label: 'Full name', value: 'fullName', minWidthRem: 11 },
  { label: 'Age/Sex', value: 'ageSex', minWidthRem: 8.5 },
  { label: 'Village', value: 'village', minWidthRem: 10 },
  { label: 'Created by', value: 'createdBy', minWidthRem: 12 },
  { label: 'Created', value: 'createdAt', minWidthRem: 10.5 },
  { label: 'Updated', value: 'updatedAt', minWidthRem: 10.5 },
  { label: 'Status', value: 'status', minWidthRem: 8 },
]

const allPatientColumnKeys = patientColumnOptions.map((column) => column.value)
const fallbackPatientColumns: PatientColumnKey[] = ['fullName']
const visiblePatientColumns = ref<PatientColumnKey[]>([...allPatientColumnKeys])
const selectedPatients = ref<Patient[]>([])

const statusToggleInFlightPatientId = ref<string | null>(null)
const rowActionMenu = ref()
const selectedMenuPatient = ref<Patient | null>(null)

const paginatorFirst = computed(() => pageIndex.value * limit.value)
const paginatorTotalRecords = computed(() => {
  if (typeof totalPatients.value === 'number' && Number.isFinite(totalPatients.value)) return totalPatients.value
  if (hasNextPage.value) return paginatorFirst.value + limit.value + 1
  return paginatorFirst.value + patients.value.length
})
const hasPendingFilters = computed(
  () =>
    Boolean(searchText.value.trim()) ||
    siteFilter.value.length > 0,
)
const activeFilterCount = computed(() => {
  let count = 0
  if (appliedSiteFilter.value.length > 0) count++
  return count
})

type FilterChip = {
  key: string
  label: string
}

const activeFilterChips = computed<FilterChip[]>(() => {
  const chips: FilterChip[] = []
  if (appliedSiteFilter.value.length > 0) {
    const labels = appliedSiteFilter.value.map((siteId) => {
      const site = sitesCatalog.value.find((entry) => entry.id === siteId)
      return site ? site.name : siteId.slice(0, 8)
    })
    chips.push({
      key: 'site',
      label: labels.length === 1 ? `Site: ${labels[0]}` : `Sites: ${labels.join(', ')}`,
    })
  }
  return chips
})

const patientStatusOptions: Array<{ label: string; value: PatientStatus }> = [
  { label: 'Provisional', value: 'provisional' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Deceased', value: 'deceased' },
]

const sexOptions = [
  { label: 'Not set', value: '' },
  { label: 'Male', value: 'male' },
  { label: 'Female', value: 'female' },
  { label: 'Other', value: 'other' },
]

const siteFilterOptions = computed(() =>
  sitesCatalog.value
    .map((site) => ({ label: site.name, value: site.id }))
    .sort((a, b) => a.label.localeCompare(b.label)),
)

const siteOptionsForSelect = computed(() =>
  sitesCatalog.value
    .map((site) => ({ label: site.name, value: site.id }))
    .sort((a, b) => a.label.localeCompare(b.label)),
)

const selectedPatientColumns = computed<PatientColumnKey[]>({
  get: () => visiblePatientColumns.value,
  set: (value: unknown) => {
    visiblePatientColumns.value = normalizePatientColumns(value)
  },
})

const patientTableMinWidth = computed(() => {
  const dataColumnsWidth = patientColumnOptions
    .filter((column) => visiblePatientColumns.value.includes(column.value))
    .reduce((total, column) => total + column.minWidthRem, 0)

  const selectionColumnWidth = 3
  const actionColumnWidth = 3.5
  return Math.max(24, dataColumnsWidth + selectionColumnWidth + actionColumnWidth)
})

const rowActionItems = computed<ActionMenuItem[]>(() => [
  {
    label: 'Update patient',
    lucideIcon: Pencil,
    disabled: !canUpdatePatients.value,
    command: () => {
      if (!selectedMenuPatient.value) return
      void openEditDialog(selectedMenuPatient.value)
    },
  },
  {
    label: 'Assign site',
    lucideIcon: MapPinPlus,
    disabled: !canAddPatientSites.value,
    command: () => {
      if (!selectedMenuPatient.value) return
      void openSitesDialog(selectedMenuPatient.value)
    },
  },
  {
    label: 'Deactivate patient',
    lucideIcon: Ban,
    tone: 'danger',
    disabled:
      !selectedMenuPatient.value?.isActive ||
      !canDisablePatient(selectedMenuPatient.value) ||
      Boolean(statusToggleInFlightPatientId.value),
    command: () => {
      if (!selectedMenuPatient.value) return
      deactivatePatient(selectedMenuPatient.value)
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

function toTrimmedString(value: unknown) {
  if (typeof value === 'string') return value.trim()
  if (typeof value === 'number' && Number.isFinite(value)) return String(value)
  return ''
}

function resolvePatientStatus(
  record: Record<string, unknown>,
  fallback: PatientStatus = 'provisional',
): PatientStatus {
  for (const key of ['is_deceased', 'isDeceased', 'deceased']) {
    const parsed = parseBooleanLike(record[key])
    if (parsed === true) return 'deceased'
  }

  const statusValue = record.status ?? record.state
  const status = typeof statusValue === 'string' ? statusValue.trim().toLowerCase() : ''
  if (['provisional', 'pending', 'new'].includes(status)) return 'provisional'
  if (['active', 'enabled'].includes(status)) return 'active'
  if (['inactive', 'disabled', 'suspended', 'blocked', 'archived'].includes(status)) return 'inactive'
  if (['deceased', 'dead'].includes(status)) return 'deceased'

  for (const key of ['is_active', 'isActive', 'active', 'enabled']) {
    const parsed = parseBooleanLike(record[key])
    if (parsed !== null) return parsed ? 'active' : 'inactive'
  }

  for (const key of ['is_disabled', 'isDisabled', 'disabled', 'inactive']) {
    const parsed = parseBooleanLike(record[key])
    if (parsed !== null) return parsed ? 'inactive' : 'active'
  }

  return fallback
}

function formatStatus(value: PatientStatus) {
  if (value === 'provisional') return 'Provisional'
  if (value === 'active') return 'Active'
  if (value === 'inactive') return 'Inactive'
  return 'Deceased'
}

function statusClass(value: PatientStatus) {
  if (value === 'active') return 'status-active'
  if (value === 'provisional') return 'status-provisional'
  if (value === 'deceased') return 'status-deceased'
  return 'status-inactive'
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
      const candidate = record[key]
      if (Array.isArray(candidate)) return candidate
      if (candidate && typeof candidate === 'object') queue.push(candidate)
    }

    for (const value of Object.values(record)) {
      if (Array.isArray(value)) return value
      if (value && typeof value === 'object') nested.push(value)
    }

    for (const value of nested) queue.push(value)
  }

  return []
}

function normalizeSite(value: unknown): SiteOption | null {
  if (typeof value === 'string') {
    const normalized = value.trim()
    if (!normalized) return null
    return { id: normalized, name: normalized, code: '' }
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

  const nameCandidate = record.name ?? record.label ?? record.code ?? record.slug ?? id
  const name =
    typeof nameCandidate === 'string'
      ? nameCandidate.trim()
      : typeof nameCandidate === 'number' && Number.isFinite(nameCandidate)
        ? String(nameCandidate)
        : ''

  const code =
    typeof record.code === 'string'
      ? record.code.trim()
      : typeof record.code === 'number' && Number.isFinite(record.code)
        ? String(record.code)
        : ''

  if (!id) return null
  return {
    id,
    name: name || id,
    code,
  }
}

function normalizeSiteList(payload: unknown) {
  const source = extractCollection(payload, ['sites', 'patient_sites', 'assignments'])
  const normalized = source.map(normalizeSite).filter((entry): entry is SiteOption => !!entry)
  const byId = new Map<string, SiteSummary>()
  for (const site of normalized) {
    byId.set(site.id.toLowerCase(), { id: site.id, name: site.name, code: site.code })
  }
  return [...byId.values()]
}

function resolveCreatedBy(record: Record<string, unknown>) {
  const createdByNameCandidate = record.created_by_name ?? record.createdByName
  const fallbackName = toTrimmedString(createdByNameCandidate)
  const direct = record.created_by ?? record.createdBy ?? record.created_by_id ?? record.createdById

  if (typeof direct === 'string' || (typeof direct === 'number' && Number.isFinite(direct))) {
    return { id: toTrimmedString(direct), name: fallbackName }
  }

  if (!direct || typeof direct !== 'object') {
    return { id: '', name: fallbackName }
  }

  const nested = direct as Record<string, unknown>
  const id = toTrimmedString(nested.id ?? nested.user_id ?? nested.userId ?? nested.uuid)
  const name = toTrimmedString(
    nested.full_name ??
      nested.fullName ??
      nested.name ??
      nested.username ??
      nested.display_name ??
      nested.displayName ??
      nested.email,
  )
  return { id, name: name || fallbackName }
}

function normalizePatient(value: unknown): Patient | null {
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>

  const idCandidate = record.id ?? record.patient_id ?? record.patientId
  const id =
    typeof idCandidate === 'string'
      ? idCandidate.trim()
      : typeof idCandidate === 'number' && Number.isFinite(idCandidate)
        ? String(idCandidate)
        : ''
  if (!id) return null

  const fullNameCandidate = record.full_name ?? record.fullName ?? record.name
  const fullName =
    typeof fullNameCandidate === 'string'
      ? fullNameCandidate.trim()
      : typeof fullNameCandidate === 'number' && Number.isFinite(fullNameCandidate)
        ? String(fullNameCandidate)
        : ''

  const dateOfBirthCandidate = record.date_of_birth ?? record.dateOfBirth ?? record.dob
  const dateOfBirth = typeof dateOfBirthCandidate === 'string' ? dateOfBirthCandidate.trim() : ''
  const isDobEstimated =
    parseBooleanLike(record.is_dob_estimated) ??
    parseBooleanLike(record.isDobEstimated) ??
    false

  const sexCandidate = record.sex ?? record.gender
  const sex = typeof sexCandidate === 'string' ? sexCandidate.trim() : ''
  const phoneCandidate = record.phone ?? record.phone_number ?? record.phoneNumber
  const phone = typeof phoneCandidate === 'string' ? phoneCandidate.trim() : ''
  const nationalIdCandidate = record.national_id ?? record.nationalId
  const nationalId = typeof nationalIdCandidate === 'string' ? nationalIdCandidate.trim() : ''
  const mrnCandidate = record.mrn
  const mrn = typeof mrnCandidate === 'string' ? mrnCandidate.trim() : ''
  const address = typeof record.address === 'string' ? record.address.trim() : ''
  const village = toTrimmedString(record.village ?? record.village_name ?? record.villageName)
  const origin = typeof record.origin === 'string' ? record.origin.trim() : ''

  const metadata = record.metadata ?? null

  const primarySiteIdCandidate =
    record.primary_site_id ??
    record.primarySiteId ??
    (record.primary_site && typeof record.primary_site === 'object'
      ? (record.primary_site as Record<string, unknown>).id
      : null)
  const primarySiteId =
    typeof primarySiteIdCandidate === 'string'
      ? primarySiteIdCandidate.trim()
      : typeof primarySiteIdCandidate === 'number' && Number.isFinite(primarySiteIdCandidate)
        ? String(primarySiteIdCandidate)
        : ''

  const primarySiteNameCandidate =
    record.primary_site_name ??
    record.primarySiteName ??
    (record.primary_site && typeof record.primary_site === 'object'
      ? (record.primary_site as Record<string, unknown>).name
      : null)
  const primarySiteName =
    typeof primarySiteNameCandidate === 'string'
      ? primarySiteNameCandidate.trim()
      : typeof primarySiteNameCandidate === 'number' && Number.isFinite(primarySiteNameCandidate)
        ? String(primarySiteNameCandidate)
        : ''

  const status = resolvePatientStatus(record, 'provisional')
  const isActive = status === 'active'
  const createdAt = typeof record.created_at === 'string' ? record.created_at : ''
  const updatedAt = typeof record.updated_at === 'string' ? record.updated_at : ''
  const createdBy = resolveCreatedBy(record)

  const sitesSource =
    record.sites ??
    record.patient_sites ??
    record.assignments ??
    record.site_assignments ??
    record.patientSites ??
    []
  const sites = normalizeSiteList(sitesSource)
  if (sites.length === 0 && primarySiteId) {
    sites.push({
      id: primarySiteId,
      name: primarySiteName || primarySiteId,
      code: '',
    })
  }

  return {
    id,
    fullName: fullName || 'Unnamed patient',
    dateOfBirth,
    isDobEstimated,
    sex,
    phone,
    nationalId,
    mrn,
    address,
    village,
    origin,
    metadata,
    primarySiteId,
    primarySiteName,
    sites,
    status,
    isActive,
    createdById: createdBy.id,
    createdByName: createdBy.name,
    createdAt,
    updatedAt,
  }
}

function mergeSites(existing: SiteSummary[], incoming: SiteSummary[]) {
  const merged: SiteSummary[] = []
  const seen = new Set<string>()
  for (const item of [...existing, ...incoming]) {
    const id = item.id.trim()
    if (!id) continue
    const key = id.toLowerCase()
    if (seen.has(key)) continue
    seen.add(key)
    merged.push({
      id,
      name: item.name.trim() || id,
      code: item.code.trim(),
    })
  }
  return merged
}

function mergePatientsById(entries: Patient[]) {
  const byId = new Map<string, Patient>()
  for (const entry of entries) {
    const existing = byId.get(entry.id)
    if (!existing) {
      byId.set(entry.id, {
        ...entry,
        sites: [...entry.sites],
      })
      continue
    }
    existing.sites = mergeSites(existing.sites, entry.sites)
    if (!existing.primarySiteId && entry.primarySiteId) existing.primarySiteId = entry.primarySiteId
    if (!existing.primarySiteName && entry.primarySiteName) existing.primarySiteName = entry.primarySiteName
    if (!existing.createdAt && entry.createdAt) existing.createdAt = entry.createdAt
    if (!existing.updatedAt && entry.updatedAt) existing.updatedAt = entry.updatedAt
    if (!existing.fullName || existing.fullName === 'Unnamed patient') existing.fullName = entry.fullName
    if (!existing.dateOfBirth && entry.dateOfBirth) existing.dateOfBirth = entry.dateOfBirth
    if (!existing.phone && entry.phone) existing.phone = entry.phone
    if (!existing.nationalId && entry.nationalId) existing.nationalId = entry.nationalId
    if (!existing.mrn && entry.mrn) existing.mrn = entry.mrn
    if (!existing.origin && entry.origin) existing.origin = entry.origin
    if (!existing.village && entry.village) existing.village = entry.village
    if (!existing.address && entry.address) existing.address = entry.address
    if (!existing.createdById && entry.createdById) existing.createdById = entry.createdById
    if (!existing.createdByName && entry.createdByName) existing.createdByName = entry.createdByName
    existing.status = entry.status
    existing.isActive = entry.isActive
  }
  return [...byId.values()]
}

function normalizeSelection(values: string[]) {
  return [...new Set(values.map((value) => value.trim()).filter(Boolean))].sort((a, b) =>
    a.localeCompare(b),
  )
}

function isPatientColumnKey(value: unknown): value is PatientColumnKey {
  return typeof value === 'string' && allPatientColumnKeys.includes(value as PatientColumnKey)
}

function normalizePatientColumns(values: unknown) {
  if (!Array.isArray(values)) return [...fallbackPatientColumns]

  const normalized: PatientColumnKey[] = []
  for (const value of values) {
    if (!isPatientColumnKey(value)) continue
    if (normalized.includes(value)) continue
    normalized.push(value)
  }

  if (normalized.length === 0) return [...fallbackPatientColumns]
  return normalized
}

function isPatientColumnVisible(column: PatientColumnKey) {
  return visiblePatientColumns.value.includes(column)
}

function stringifyMetadata(value: unknown) {
  if (value === null || typeof value === 'undefined' || value === '') return ''
  if (typeof value === 'string') return value.trim()
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return ''
  }
}

function parseMetadataText(value: string) {
  const normalized = value.trim()
  if (!normalized) return { hasValue: false as const, value: null }
  try {
    return { hasValue: true as const, value: JSON.parse(normalized) as unknown }
  } catch {
    throw new Error('Metadata must be valid JSON.')
  }
}

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

function formatDateTime(value: string) {
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

function formatAgeSex(dateOfBirth: string, sex: string) {
  const dob = dateOfBirth ? new Date(dateOfBirth) : null
  const now = new Date()
  let ageText = '?'

  if (dob && !Number.isNaN(dob.getTime())) {
    let age = now.getFullYear() - dob.getFullYear()
    const birthdayPassed =
      now.getMonth() > dob.getMonth() ||
      (now.getMonth() === dob.getMonth() && now.getDate() >= dob.getDate())
    if (!birthdayPassed) age--
    ageText = age >= 0 ? String(age) : '?'
  }

  if (!sex) return ageText === '?' ? 'N/A' : ageText
  const normalized = sex.trim().toLowerCase()
  const sexCode = normalized === 'male' ? 'M' : normalized === 'female' ? 'F' : normalized === 'other' ? 'O' : ''
  if (!sexCode) return ageText === '?' ? 'N/A' : ageText
  return `${ageText}${sexCode}`
}

function createdByLabel(patient: Patient) {
  if (patient.createdByName) return patient.createdByName
  if (patient.createdById) return patient.createdById
  return 'N/A'
}

function canDisablePatient(patient: Patient | null) {
  if (!patient) return false
  if (!patient.isActive) return false
  if (authStore.isAdmin) return true
  const current = currentUserId.value.trim()
  if (!current) return false
  return patient.createdById.trim() === current
}

function canOpenActionsForPatient(patient: Patient) {
  return canUpdatePatients.value || canAddPatientSites.value || canDisablePatient(patient)
}

function patchPatientListEntry(nextPatient: Patient) {
  patients.value = patients.value.map((entry) => (entry.id === nextPatient.id ? nextPatient : entry))
}

function resetPatientsPagination() {
  pageIndex.value = 0
  cursorByPage.value = [null]
  nextCursor.value = null
  hasNextPage.value = false
}

function resetCreateForm() {
  createForm.full_name = ''
  createForm.primary_site_id = ''
  createForm.date_of_birth = ''
  createForm.is_dob_estimated = false
  createForm.status = 'provisional'
  createForm.sex = ''
  createForm.phone = ''
  createForm.national_id = ''
  createForm.mrn = ''
  createForm.address = ''
  createForm.village = ''
  createForm.origin = ''
  createForm.metadata_text = ''
  createDialogError.value = ''
}

function resetEditForm() {
  editTargetPatient.value = null
  editForm.full_name = ''
  editForm.date_of_birth = ''
  editForm.is_dob_estimated = false
  editForm.status = 'provisional'
  editForm.sex = ''
  editForm.phone = ''
  editForm.national_id = ''
  editForm.mrn = ''
  editForm.address = ''
  editForm.village = ''
  editForm.origin = ''
  editForm.metadata_text = ''
  editDialogError.value = ''
}

function resetSitesForm() {
  sitesTargetPatient.value = null
  sitesForm.site_id = ''
  sitesDialogError.value = ''
}

async function loadSitesCatalog(options?: { silent?: boolean; force?: boolean }) {
  if (sitesLoading.value) return
  if (!options?.force && sitesCatalog.value.length > 0) return

  sitesLoading.value = true
  try {
    const payload = await requestJson<unknown>(
      '/admin/sites',
      { method: 'GET' },
      'Unable to load sites catalog.',
    )
    sitesCatalog.value = extractCollection(payload, ['sites'])
      .map(normalizeSite)
      .filter((entry): entry is SiteOption => !!entry)
  } catch (error) {
    if (!options?.silent) {
      const message = error instanceof Error ? error.message : 'Unable to load sites catalog.'
      toast.add({ severity: 'warn', summary: 'Sites unavailable', detail: message, life: 3500 })
    }
  } finally {
    sitesLoading.value = false
  }
}

async function loadPatients() {
  if (!canReadPatients.value) return

  patientsLoading.value = true
  patientsError.value = ''

  try {
    const requestedLimit = Math.max(1, Math.min(limit.value, 100))
    if (requestedLimit !== limit.value) limit.value = requestedLimit

    const currentCursor = cursorByPage.value[pageIndex.value] ?? null
    const params = new URLSearchParams({
      limit: String(requestedLimit),
      include_total: 'true',
    })

    if (currentCursor) params.set('cursor', currentCursor)
    if (appliedSearch.value.trim()) params.set('search', appliedSearch.value.trim())
    for (const siteId of appliedSiteFilter.value) params.append('site_id', siteId)

    const payload = await requestJson<Record<string, unknown>>(
      `/patients?${params.toString()}`,
      { method: 'GET' },
      'Unable to load patients.',
    )

    const collection =
      Array.isArray(payload.patients) ? payload.patients : extractCollection(payload, ['patients'])
    const normalized = collection.map(normalizePatient).filter((entry): entry is Patient => !!entry)
    patients.value = mergePatientsById(normalized)
    const visiblePatientIds = new Set(patients.value.map((patient) => patient.id))
    selectedPatients.value = selectedPatients.value.filter((patient) => visiblePatientIds.has(patient.id))

    const paginationSource =
      payload.pagination && typeof payload.pagination === 'object'
        ? (payload.pagination as Record<string, unknown>)
        : payload

    const apiNextCursorCandidate = paginationSource.next_cursor ?? payload.next_cursor
    const apiNextCursor =
      typeof apiNextCursorCandidate === 'string' && apiNextCursorCandidate.trim()
        ? apiNextCursorCandidate.trim()
        : null
    const apiHasNext =
      parseBooleanLike(paginationSource.has_next ?? payload.has_next) ?? Boolean(apiNextCursor)
    hasNextPage.value = apiHasNext
    nextCursor.value = apiNextCursor

    const apiLimitCandidate = paginationSource.limit
    if (
      typeof apiLimitCandidate === 'number' &&
      Number.isFinite(apiLimitCandidate) &&
      apiLimitCandidate > 0 &&
      apiLimitCandidate <= 100 &&
      apiLimitCandidate !== limit.value
    ) {
      limit.value = apiLimitCandidate
    }

    cursorByPage.value = cursorByPage.value.slice(0, pageIndex.value + 1)
    if (apiHasNext && apiNextCursor) {
      cursorByPage.value[pageIndex.value + 1] = apiNextCursor
    }

    const totalCandidate = payload.total ?? paginationSource.total
    if (typeof totalCandidate === 'number' && Number.isFinite(totalCandidate) && totalCandidate >= 0) {
      totalPatients.value = totalCandidate
    } else if (pageIndex.value === 0) {
      totalPatients.value = null
    }
  } catch (error) {
    patientsError.value = error instanceof Error ? error.message : 'Unable to load patients.'
    patients.value = []
    selectedPatients.value = []
    totalPatients.value = null
    hasNextPage.value = false
    nextCursor.value = null
  } finally {
    patientsLoading.value = false
  }
}

function onPageChange(event: { page: number; first: number; rows: number }) {
  const nextRows = Math.max(1, Math.min(event.rows, 100))
  if (nextRows !== limit.value) {
    limit.value = nextRows
    resetPatientsPagination()
    void loadPatients()
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
    void loadPatients()
    return
  }

  if (targetPage === pageIndex.value + 1 && hasNextPage.value && nextCursor.value) {
    cursorByPage.value[targetPage] = nextCursor.value
    pageIndex.value = targetPage
    void loadPatients()
  }
}

function refreshPatients() {
  void loadPatients()
}

let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

function applyFilters() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
    searchDebounceTimer = null
  }

  const nextSearch = searchText.value.trim()
  const nextSiteFilter = normalizeSelection(siteFilter.value)

  const changed =
    nextSearch !== appliedSearch.value ||
    nextSiteFilter.join('|') !== appliedSiteFilter.value.join('|')
  if (!changed) return

  siteFilter.value = [...nextSiteFilter]
  appliedSearch.value = nextSearch
  appliedSiteFilter.value = [...nextSiteFilter]
  resetPatientsPagination()
  void loadPatients()
}

function scheduleSearchFilter() {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    searchDebounceTimer = null
    applyFilters()
  }, 280)
}

function clearFilters() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
    searchDebounceTimer = null
  }

  const shouldReload =
    Boolean(appliedSearch.value.trim()) ||
    Boolean(searchText.value.trim()) ||
    appliedSiteFilter.value.length > 0 ||
    siteFilter.value.length > 0
  if (!shouldReload) return

  searchText.value = ''
  appliedSearch.value = ''
  siteFilter.value = []
  appliedSiteFilter.value = []
  resetPatientsPagination()
  void loadPatients()
}

async function toggleFilters() {
  filtersExpanded.value = !filtersExpanded.value
  if (!filtersExpanded.value) return
  await loadSitesCatalog({ silent: true })
}

function removeFilterChip(key: string) {
  if (key === 'site') {
    siteFilter.value = []
    appliedSiteFilter.value = []
  }
  resetPatientsPagination()
  void loadPatients()
}

function openRowActionMenu(event: Event, patient: Patient) {
  selectedMenuPatient.value = patient
  rowActionMenu.value?.toggle(event)
}

async function loadPatientById(patientId: string, fallbackError: string) {
  const payload = await requestJson<unknown>(`/patients/${patientId}`, { method: 'GET' }, fallbackError)
  const normalized = normalizePatient(payload)
  if (!normalized) throw new Error('Patient response is missing required fields.')
  return normalized
}

async function openCreateDialog() {
  if (!canCreatePatients.value) return
  resetCreateForm()
  createDialogVisible.value = true
  await loadSitesCatalog({ silent: true })
}

async function submitCreatePatient() {
  if (!canCreatePatients.value || createDialogBusy.value) return
  if (!createForm.full_name.trim()) {
    createDialogError.value = 'Full name is required.'
    return
  }

  createDialogBusy.value = true
  createDialogError.value = ''

  try {
    const payload: Record<string, unknown> = {
      full_name: createForm.full_name.trim(),
      status: createForm.status,
    }
    if (createForm.primary_site_id.trim()) payload.primary_site_id = createForm.primary_site_id.trim()
    if (createForm.date_of_birth.trim()) payload.date_of_birth = createForm.date_of_birth.trim()
    if (createForm.is_dob_estimated) payload.is_dob_estimated = true
    if (createForm.sex.trim()) payload.sex = createForm.sex.trim()
    if (createForm.phone.trim()) payload.phone = createForm.phone.trim()
    if (createForm.national_id.trim()) payload.national_id = createForm.national_id.trim()
    if (createForm.mrn.trim()) payload.mrn = createForm.mrn.trim()
    if (createForm.address.trim()) payload.address = createForm.address.trim()
    if (createForm.village.trim()) payload.village = createForm.village.trim()
    if (createForm.origin.trim()) payload.origin = createForm.origin.trim()
    const parsedMetadata = parseMetadataText(createForm.metadata_text)
    if (parsedMetadata.hasValue) payload.metadata = parsedMetadata.value

    const response = await requestJson<unknown>(
      '/patients',
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      },
      'Unable to register patient.',
    )

    const normalized = normalizePatient(response)
    if (normalized) patchPatientListEntry(normalized)

    createDialogVisible.value = false
    toast.add({ severity: 'success', summary: 'Patient registered', life: 2500 })
    resetPatientsPagination()
    await loadPatients()
  } catch (error) {
    createDialogError.value = error instanceof Error ? error.message : 'Unable to register patient.'
  } finally {
    createDialogBusy.value = false
  }
}

async function openEditDialog(patient: Patient) {
  if (!canUpdatePatients.value) return
  resetEditForm()
  editTargetPatient.value = patient
  editForm.full_name = patient.fullName
  editForm.date_of_birth = patient.dateOfBirth
  editForm.is_dob_estimated = patient.isDobEstimated
  editForm.status = patient.status
  editForm.sex = patient.sex
  editForm.phone = patient.phone
  editForm.national_id = patient.nationalId
  editForm.mrn = patient.mrn
  editForm.address = patient.address
  editForm.village = patient.village
  editForm.origin = patient.origin
  editForm.metadata_text = stringifyMetadata(patient.metadata)
  editDialogVisible.value = true

  try {
    const latest = await loadPatientById(patient.id, 'Unable to load patient details.')
    editTargetPatient.value = latest
    editForm.full_name = latest.fullName
    editForm.date_of_birth = latest.dateOfBirth
    editForm.is_dob_estimated = latest.isDobEstimated
    editForm.status = latest.status
    editForm.sex = latest.sex
    editForm.phone = latest.phone
    editForm.national_id = latest.nationalId
    editForm.mrn = latest.mrn
    editForm.address = latest.address
    editForm.village = latest.village
    editForm.origin = latest.origin
    editForm.metadata_text = stringifyMetadata(latest.metadata)
    patchPatientListEntry(latest)
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load patient details.'
    editDialogError.value = message
  }
}

function normalizeOptionalString(value: string) {
  return value.trim()
}

async function submitEditPatient() {
  if (!canUpdatePatients.value || editDialogBusy.value) return
  if (!editTargetPatient.value) return
  if (!editForm.full_name.trim()) {
    editDialogError.value = 'Full name is required.'
    return
  }

  const target = editTargetPatient.value
  const payload: Record<string, unknown> = {}

  const nextName = normalizeOptionalString(editForm.full_name)
  if (nextName && nextName !== target.fullName) payload.full_name = nextName

  const nextDob = normalizeOptionalString(editForm.date_of_birth)
  if (nextDob && nextDob !== target.dateOfBirth) payload.date_of_birth = nextDob

  if (editForm.is_dob_estimated !== target.isDobEstimated) {
    payload.is_dob_estimated = editForm.is_dob_estimated
  }

  if (editForm.status !== target.status) payload.status = editForm.status

  const nextSex = normalizeOptionalString(editForm.sex)
  if (nextSex !== target.sex) payload.sex = nextSex

  const nextPhone = normalizeOptionalString(editForm.phone)
  if (nextPhone !== target.phone) payload.phone = nextPhone

  const nextNationalId = normalizeOptionalString(editForm.national_id)
  if (nextNationalId !== target.nationalId) payload.national_id = nextNationalId

  const nextAddress = normalizeOptionalString(editForm.address)
  if (nextAddress !== target.address) payload.address = nextAddress

  const nextVillage = normalizeOptionalString(editForm.village)
  if (nextVillage !== target.village) payload.village = nextVillage

  const nextOrigin = normalizeOptionalString(editForm.origin)
  if (nextOrigin !== target.origin) payload.origin = nextOrigin

  const existingMrn = target.mrn.trim()
  const nextMrn = normalizeOptionalString(editForm.mrn)
  if (!existingMrn && nextMrn) payload.mrn = nextMrn

  const currentMetadata = stringifyMetadata(target.metadata)
  const nextMetadataRaw = editForm.metadata_text.trim()
  if (nextMetadataRaw && nextMetadataRaw !== currentMetadata.trim()) {
    const parsedMetadata = parseMetadataText(editForm.metadata_text)
    if (parsedMetadata.hasValue) payload.metadata = parsedMetadata.value
  }

  if (Object.keys(payload).length === 0) {
    toast.add({ severity: 'info', summary: 'No changes to save', life: 2200 })
    return
  }

  editDialogBusy.value = true
  editDialogError.value = ''
  try {
    const response = await requestJson<unknown>(
      `/patients/${target.id}`,
      {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      },
      'Unable to update patient.',
    )
    const normalized = normalizePatient(response)
    if (normalized) patchPatientListEntry(normalized)
    toast.add({ severity: 'success', summary: 'Patient updated', life: 2500 })
    editDialogVisible.value = false
    await loadPatients()
  } catch (error) {
    editDialogError.value = error instanceof Error ? error.message : 'Unable to update patient.'
  } finally {
    editDialogBusy.value = false
  }
}

async function openSitesDialog(patient: Patient) {
  if (!canAddPatientSites.value) return
  resetSitesForm()
  sitesTargetPatient.value = patient
  sitesDialogVisible.value = true
  await loadSitesCatalog({ silent: true })

  sitesDialogLoading.value = true
  sitesDialogError.value = ''
  try {
    const latest = await loadPatientById(patient.id, 'Unable to load patient details.')
    sitesTargetPatient.value = latest
    patchPatientListEntry(latest)
  } catch (error) {
    sitesDialogError.value = error instanceof Error ? error.message : 'Unable to load patient details.'
  } finally {
    sitesDialogLoading.value = false
  }
}

async function submitSitesAssignment() {
  if (!canAddPatientSites.value || sitesDialogBusy.value) return
  if (!sitesTargetPatient.value) return

  const siteId = sitesForm.site_id.trim()
  if (!siteId) {
    sitesDialogError.value = 'Site is required.'
    return
  }

  sitesDialogBusy.value = true
  sitesDialogError.value = ''
  try {
    const payload = await requestJson<unknown>(
      `/patients/${sitesTargetPatient.value.id}/sites`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ site_id: siteId }),
      },
      'Unable to assign patient site.',
    )

    const assignedSites = normalizeSiteList(payload)
    if (assignedSites.length > 0) {
      const nextPatient: Patient = {
        ...sitesTargetPatient.value,
        sites: mergeSites(sitesTargetPatient.value.sites, assignedSites),
      }
      sitesTargetPatient.value = nextPatient
      patchPatientListEntry(nextPatient)
    } else {
      await loadPatients()
    }

    toast.add({ severity: 'success', summary: 'Site assigned', life: 2500 })
    sitesDialogVisible.value = false
  } catch (error) {
    sitesDialogError.value = error instanceof Error ? error.message : 'Unable to assign patient site.'
  } finally {
    sitesDialogBusy.value = false
  }
}

function deactivatePatient(patient: Patient) {
  if (!canDisablePatient(patient) || statusToggleInFlightPatientId.value) return

  confirm.require({
    header: 'Deactivate patient',
    message: `Deactivate ${patient.fullName}?`,
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Deactivate patient',
    rejectLabel: 'Cancel',
    accept: async () => {
      statusToggleInFlightPatientId.value = patient.id
      patientsError.value = ''
      try {
        const payload = await requestJson<unknown>(
          `/patients/${patient.id}/disable`,
          { method: 'POST' },
          'Unable to deactivate patient.',
        )
        const normalized = normalizePatient(payload)
        if (normalized) {
          patchPatientListEntry(normalized)
        } else {
          await loadPatients()
        }
        toast.add({ severity: 'success', summary: 'Patient deactivated', life: 2500 })
      } catch (error) {
        const response = error instanceof Error ? error.message : 'Unable to deactivate patient.'
        patientsError.value = response
        toast.add({
          severity: 'error',
          summary: 'Deactivate failed',
          detail: response,
          life: 3500,
        })
      } finally {
        statusToggleInFlightPatientId.value = null
      }
    },
  })
}

void loadPatients()
</script>

<template>
  <PageSectionHeader
    title="Patients"
    description="Patient registration, lookup, updates, site assignment, and deactivation"
  >
    <template #actions>
      <ToolbarIconButton
        :disabled="!canReadPatients"
        title="Refresh patients"
        ariaLabel="Refresh patients"
        @click="refreshPatients"
      >
        <RefreshCcw :size="16" />
      </ToolbarIconButton>
      <CreateActionButton
        label="Register patient"
        :disabled="!canCreatePatients"
        @click="openCreateDialog"
      />
    </template>
  </PageSectionHeader>

  <Message v-if="!canReadPatients" severity="warn" size="small" variant="simple">
    You need <code>patients.read</code> permission to view patients.
  </Message>
  <Message v-else-if="!canCreatePatients && !canUpdatePatients && !canAddPatientSites" severity="warn" size="small" variant="simple">
    You can view patients, but <code>patients.create</code>, <code>patients.update</code>, or
    <code>patients.add_site</code> is required for changes.
  </Message>
  <Message v-if="patientsError" severity="error" size="small" variant="simple">
    {{ patientsError }}
  </Message>

  <Card class="patients-card">
    <template #content>
      <div class="toolbar">
        <div class="toolbar-left">
          <div class="search-wrap">
            <InputText
              v-model="searchText"
              placeholder="Search full name or village..."
              :disabled="!canReadPatients"
              @input="scheduleSearchFilter"
              @keyup.enter="applyFilters"
            />
          </div>
          <button
            class="filter-toggle"
            :class="{ active: filtersExpanded, 'has-filters': activeFilterCount > 0 }"
            :disabled="!canReadPatients"
            @click="toggleFilters"
          >
            <SlidersHorizontal :size="15" />
            <span>Filters</span>
            <span v-if="activeFilterCount > 0" class="filter-badge">{{ activeFilterCount }}</span>
          </button>
        </div>
        <div class="toolbar-right">
          <div class="column-toggle-wrap">
            <MultiSelect
              v-model="selectedPatientColumns"
              :options="patientColumnOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Columns"
              :disabled="!canReadPatients"
              :maxSelectedLabels="1"
              selectedItemsLabel="{0} columns"
            />
          </div>
        </div>
      </div>

      <Transition name="panel-slide">
        <div v-if="filtersExpanded" class="filter-panel">
          <div class="filter-panel-grid">
            <FilterPanelField label="Site">
              <MultiSelect
                v-model="siteFilter"
                :options="siteFilterOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="All sites"
                :loading="sitesLoading"
                :disabled="!canReadPatients"
                :maxSelectedLabels="1"
                selectedItemsLabel="{0} selected"
                @change="applyFilters"
              />
            </FilterPanelField>
          </div>
          <button
            v-if="hasPendingFilters"
            class="clear-all-btn"
            :disabled="!canReadPatients"
            @click="clearFilters"
          >
            <X :size="13" />
            <span>Clear all filters</span>
          </button>
        </div>
      </Transition>

      <div v-if="activeFilterChips.length > 0" class="active-chips">
        <TransitionGroup name="chip-fade">
          <button
            v-for="chip in activeFilterChips"
            :key="chip.key"
            class="filter-chip"
            @click="removeFilterChip(chip.key)"
          >
            <span>{{ chip.label }}</span>
            <X :size="12" />
          </button>
        </TransitionGroup>
      </div>

      <DataTable
        v-model:selection="selectedPatients"
        :value="patients"
        dataKey="id"
        class="patients-table"
        :tableStyle="{ minWidth: `${patientTableMinWidth}rem` }"
        :loading="patientsLoading"
        paginator
        :rows="limit"
        :first="paginatorFirst"
        :totalRecords="paginatorTotalRecords"
        :paginatorTemplate="TABLE_PAGINATOR_TEMPLATE"
        :currentPageReportTemplate="TABLE_PAGE_REPORT_TEMPLATE"
        :rowsPerPageOptions="[...TABLE_ROWS_PER_PAGE_OPTIONS]"
        scrollable
        @page="onPageChange"
        stripedRows
        rowHover
      >
        <template #empty>
          <div class="table-empty">No patients match the current filters.</div>
        </template>
        <Column selectionMode="multiple" style="width: 3rem" />
        <Column v-if="isPatientColumnVisible('fullName')" field="fullName" header="Patient" style="min-width: 11rem" />
        <Column v-if="isPatientColumnVisible('ageSex')" header="Age/Sex" style="min-width: 8.5rem">
          <template #body="{ data }">
            <span>{{ formatAgeSex(data.dateOfBirth, data.sex) }}</span>
          </template>
        </Column>
        <Column v-if="isPatientColumnVisible('village')" header="Village" style="min-width: 10rem">
          <template #body="{ data }">
            <span class="cell-truncate" :title="data.village || 'N/A'">{{ data.village || 'N/A' }}</span>
          </template>
        </Column>
        <Column v-if="isPatientColumnVisible('createdBy')" header="Created by" style="min-width: 12rem">
          <template #body="{ data }">
            <span class="cell-truncate" :title="createdByLabel(data)">{{ createdByLabel(data) }}</span>
          </template>
        </Column>
        <Column v-if="isPatientColumnVisible('createdAt')" header="Created" style="min-width: 10.5rem">
          <template #body="{ data }">
            <span class="cell-date">{{ formatDateTime(data.createdAt) }}</span>
          </template>
        </Column>
        <Column v-if="isPatientColumnVisible('updatedAt')" header="Updated" style="min-width: 10.5rem">
          <template #body="{ data }">
            <span class="cell-date">{{ formatDateTime(data.updatedAt) }}</span>
          </template>
        </Column>
        <Column v-if="isPatientColumnVisible('status')" header="Status" style="width: 8rem">
          <template #body="{ data }">
            <span class="status-cell" :class="statusClass(data.status)">
              <span class="status-dot"></span>
              {{ formatStatus(data.status) }}
            </span>
          </template>
        </Column>
        <Column header="" style="width: 3.5rem">
          <template #body="{ data }">
            <div class="actions">
              <RowActionButton
                :disabled="!canOpenActionsForPatient(data)"
                :loading="statusToggleInFlightPatientId === data.id"
                ariaLabel="Open patient actions"
                title="Open patient actions"
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
    header="Register patient"
    :style="{ width: 'min(50rem, 96vw)' }"
    @hide="resetCreateForm"
  >
    <div class="dialog-grid">
      <Message v-if="createDialogError" severity="error" size="small" variant="simple">
        {{ createDialogError }}
      </Message>

      <label class="field">
        <span>Full name</span>
        <InputText v-model="createForm.full_name" autocomplete="off" />
      </label>

      <label class="field">
        <span>Primary site (optional)</span>
        <Select
          v-if="siteOptionsForSelect.length > 0"
          v-model="createForm.primary_site_id"
          :options="siteOptionsForSelect"
          optionLabel="label"
          optionValue="value"
          placeholder="Select primary site"
          filter
        />
        <InputText v-else v-model="createForm.primary_site_id" placeholder="Site UUID" autocomplete="off" />
      </label>

      <div class="field-inline">
        <label class="field">
          <span>Date of birth</span>
          <InputText v-model="createForm.date_of_birth" type="date" />
        </label>
        <label class="checkbox-field">
          <Checkbox v-model="createForm.is_dob_estimated" binary />
          <span>DOB is estimated</span>
        </label>
      </div>

      <label class="field">
        <span>Sex</span>
        <Select
          v-model="createForm.sex"
          :options="sexOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Not set"
        />
      </label>

      <label class="field">
        <span>Status</span>
        <Select
          v-model="createForm.status"
          :options="patientStatusOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Select status"
        />
      </label>

      <label class="field">
        <span>Phone</span>
        <InputText v-model="createForm.phone" autocomplete="off" />
      </label>

      <label class="field">
        <span>National ID</span>
        <InputText v-model="createForm.national_id" autocomplete="off" />
      </label>

      <label class="field">
        <span>MRN</span>
        <InputText v-model="createForm.mrn" autocomplete="off" />
      </label>

      <label class="field">
        <span>Address</span>
        <InputText v-model="createForm.address" autocomplete="off" />
      </label>

      <label class="field">
        <span>Village</span>
        <InputText v-model="createForm.village" autocomplete="off" />
      </label>

      <label class="field">
        <span>Origin</span>
        <InputText v-model="createForm.origin" autocomplete="off" />
      </label>

      <label class="field">
        <span>Metadata JSON</span>
        <Textarea v-model="createForm.metadata_text" rows="4" autoResize placeholder="{ &quot;key&quot;: &quot;value&quot; }" />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="createDialogVisible = false" />
        <Button label="Register patient" :loading="createDialogBusy" @click="submitCreatePatient" />
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="editDialogVisible"
    modal
    header="Update patient"
    :style="{ width: 'min(50rem, 96vw)' }"
    @hide="resetEditForm"
  >
    <div class="dialog-grid">
      <Message v-if="editDialogError" severity="error" size="small" variant="simple">
        {{ editDialogError }}
      </Message>

      <label class="field">
        <span>Full name</span>
        <InputText v-model="editForm.full_name" autocomplete="off" />
      </label>

      <div class="field-inline">
        <label class="field">
          <span>Date of birth</span>
          <InputText v-model="editForm.date_of_birth" type="date" />
        </label>
        <label class="checkbox-field">
          <Checkbox v-model="editForm.is_dob_estimated" binary />
          <span>DOB is estimated</span>
        </label>
      </div>

      <label class="field">
        <span>Sex</span>
        <Select
          v-model="editForm.sex"
          :options="sexOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Not set"
        />
      </label>

      <label class="field">
        <span>Status</span>
        <Select
          v-model="editForm.status"
          :options="patientStatusOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Select status"
        />
      </label>

      <label class="field">
        <span>Phone</span>
        <InputText v-model="editForm.phone" autocomplete="off" />
      </label>

      <label class="field">
        <span>National ID</span>
        <InputText v-model="editForm.national_id" autocomplete="off" />
      </label>

      <label class="field">
        <span>MRN</span>
        <InputText v-model="editForm.mrn" :disabled="Boolean(editTargetPatient?.mrn)" autocomplete="off" />
      </label>
      <p v-if="editTargetPatient?.mrn" class="help-text">MRN is immutable once set.</p>

      <label class="field">
        <span>Address</span>
        <InputText v-model="editForm.address" autocomplete="off" />
      </label>

      <label class="field">
        <span>Village</span>
        <InputText v-model="editForm.village" autocomplete="off" />
      </label>

      <label class="field">
        <span>Origin</span>
        <InputText v-model="editForm.origin" autocomplete="off" />
      </label>

      <label class="field">
        <span>Metadata JSON</span>
        <Textarea v-model="editForm.metadata_text" rows="4" autoResize placeholder="{ &quot;key&quot;: &quot;value&quot; }" />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="editDialogVisible = false" />
        <Button label="Save changes" :loading="editDialogBusy" @click="submitEditPatient" />
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="sitesDialogVisible"
    modal
    header="Assign site"
    :style="{ width: 'min(38rem, 96vw)' }"
    @hide="resetSitesForm"
  >
    <div class="dialog-grid">
      <Message v-if="sitesDialogError" severity="error" size="small" variant="simple">
        {{ sitesDialogError }}
      </Message>

      <p v-if="sitesTargetPatient" class="dialog-subtitle">{{ sitesTargetPatient.fullName }}</p>
      <p class="help-text">This adds a site assignment in addition to the primary site.</p>

      <label class="field">
        <span>Site</span>
        <Select
          v-if="siteOptionsForSelect.length > 0"
          v-model="sitesForm.site_id"
          :options="siteOptionsForSelect"
          optionLabel="label"
          optionValue="value"
          placeholder="Select site"
          filter
          :loading="sitesLoading || sitesDialogLoading"
        />
        <InputText
          v-else
          v-model="sitesForm.site_id"
          placeholder="Site UUID"
          autocomplete="off"
          :disabled="sitesDialogLoading"
        />
      </label>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <Button label="Cancel" text @click="sitesDialogVisible = false" />
        <Button
          label="Assign site"
          :loading="sitesDialogBusy"
          :disabled="!canAddPatientSites"
          @click="submitSitesAssignment"
        />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.patients-card {
  margin-bottom: 0.9rem;
  min-width: 0;
  width: 100%;
  max-width: 100%;
}

.toolbar {
  --toolbar-control-height: var(--app-control-height, 2.2rem);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.55rem;
  flex-wrap: wrap;
  padding-bottom: 0.8rem;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1 1 0;
  min-width: 0;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  flex-shrink: 0;
}

.column-toggle-wrap {
  min-width: 11rem;
  width: min(15rem, 46vw);
}

.column-toggle-wrap :deep(.p-multiselect) {
  width: 100%;
  height: var(--toolbar-control-height);
  min-height: var(--toolbar-control-height);
  border-radius: var(--app-control-radius, 0.62rem);
  border: 1px solid var(--p-surface-200);
  background: var(--p-surface-50);
}

.column-toggle-wrap :deep(.p-multiselect:hover) {
  border-color: var(--p-surface-400);
  background: var(--p-surface-100);
}

.column-toggle-wrap :deep(.p-multiselect-label) {
  font-size: var(--app-control-font-size, 0.83rem);
  padding-block: 0.42rem;
}

.column-toggle-wrap :deep(.p-multiselect-trigger) {
  width: 2rem;
}

.search-wrap {
  flex: 1 1 14rem;
  max-width: 28rem;
  min-width: 0;
}

.search-wrap :deep(.p-inputtext) {
  width: 100%;
  height: var(--toolbar-control-height);
  min-height: var(--toolbar-control-height);
  box-sizing: border-box;
  border-radius: var(--app-control-radius, 0.62rem);
  padding: 0 0.88rem;
  font-size: var(--app-control-font-size, 0.83rem);
  line-height: 1.2;
  border: 1px solid var(--p-surface-200);
  background: var(--p-surface-50);
  transition: border-color 0.15s, background 0.15s;
}

.search-wrap :deep(.p-inputtext:focus) {
  background: var(--p-surface-0);
  border-color: var(--p-surface-400);
}

.filter-toggle {
  display: inline-flex;
  align-items: center;
  height: var(--toolbar-control-height);
  gap: 0.38rem;
  padding: 0 0.78rem;
  font-size: var(--app-control-font-size, 0.83rem);
  font-weight: 550;
  color: var(--p-surface-700);
  border: 1px solid var(--p-surface-200);
  border-radius: var(--app-control-radius, 0.62rem);
  background: var(--p-surface-50);
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
  line-height: 1.4;
}

.filter-toggle:hover:not(:disabled) {
  border-color: var(--p-surface-400);
  background: var(--p-surface-100);
}

.filter-toggle.active {
  border-color: var(--p-surface-400);
  background: var(--p-surface-100);
  color: var(--p-surface-800);
}

.filter-toggle.has-filters {
  border-color: var(--p-surface-500);
  color: var(--p-surface-800);
}

.filter-toggle:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.filter-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.2rem;
  height: 1.2rem;
  padding: 0 0.3rem;
  border-radius: 0.6rem;
  background: var(--p-surface-800);
  color: var(--p-surface-0);
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
}

.filter-panel {
  --filter-control-height: var(--toolbar-control-height, var(--app-control-height, 2.2rem));
  border: 1px solid var(--p-surface-200);
  border-radius: 0.72rem;
  background: var(--p-surface-50);
  padding: 0.9rem 1rem;
  margin-bottom: 0.8rem;
}

.filter-panel-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(12rem, 1fr));
  gap: 0.75rem;
}

.clear-all-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  margin-top: 0.7rem;
  padding: 0.3rem 0.55rem;
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--p-primary-700);
  border: 1px solid color-mix(in srgb, var(--p-primary-color) 32%, var(--p-surface-0));
  border-radius: 0.45rem;
  background: color-mix(in srgb, var(--p-primary-color) 10%, var(--p-surface-0));
  cursor: pointer;
  transition: all 0.15s;
}

.clear-all-btn:hover:not(:disabled) {
  background: color-mix(in srgb, var(--p-primary-color) 18%, var(--p-surface-0));
  border-color: color-mix(in srgb, var(--p-primary-color) 42%, var(--p-surface-0));
}

.clear-all-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.panel-slide-enter-active,
.panel-slide-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}

.panel-slide-enter-from,
.panel-slide-leave-to {
  opacity: 0;
  max-height: 0;
  margin-bottom: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.panel-slide-enter-to,
.panel-slide-leave-from {
  opacity: 1;
  max-height: 18rem;
}

.cell-stack {
  display: grid;
  gap: 0.1rem;
}

.cell-subtle,
.cell-date {
  font-size: 0.82rem;
  color: var(--p-surface-600);
  white-space: nowrap;
}

.cell-truncate {
  display: inline-block;
  max-width: 13rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.82rem;
  color: var(--p-surface-600);
}

.cell-list {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  max-width: 13rem;
  min-width: 0;
}

.cell-pill {
  display: inline-flex;
  align-items: center;
  max-width: 11rem;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  border: 1px solid var(--p-surface-200);
  background: var(--p-surface-100);
  color: var(--p-surface-900);
  font-size: 0.8rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cell-pill-overflow {
  padding: 0.15rem 0.42rem;
  font-size: 0.78rem;
}

.cell-empty {
  color: var(--p-surface-400);
}

.status-cell {
  display: inline-flex;
  align-items: center;
  gap: 0.38rem;
  border-radius: 999px;
  padding: 0.14rem 0.48rem;
  font-size: 0.8rem;
  font-weight: 600;
  border: 1px solid transparent;
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

.status-provisional {
  color: #005f82;
  border-color: #b9e6f2;
  background: #ebf8fd;
}

.status-provisional .status-dot {
  background: #06b6d4;
  box-shadow: 0 0 0 3px rgb(6 182 212 / 16%);
}

.status-inactive {
  color: #9a4214;
  border-color: #f8d2bd;
  background: #fff2e8;
}

.status-inactive .status-dot {
  background: #f97316;
  box-shadow: 0 0 0 3px rgb(249 115 22 / 16%);
}

.status-deceased {
  color: #5f6671;
  border-color: #d8dbe1;
  background: #f4f5f7;
}

.status-deceased .status-dot {
  background: #7c8798;
  box-shadow: 0 0 0 3px rgb(124 135 152 / 16%);
}

.active-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  padding-bottom: 0.7rem;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 0.25rem 0.55rem;
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--p-surface-800);
  border: 1px solid var(--p-surface-300);
  border-radius: 0.5rem;
  background: var(--p-surface-100);
  cursor: pointer;
  transition: all 0.15s;
  line-height: 1.3;
}

.filter-chip:hover {
  background: var(--p-surface-200);
  border-color: var(--p-surface-400);
}

.filter-chip:hover > :last-child {
  color: var(--p-primary-700);
}

.chip-fade-enter-active,
.chip-fade-leave-active {
  transition: all 0.18s ease;
}

.chip-fade-enter-from,
.chip-fade-leave-to {
  opacity: 0;
  transform: scale(0.85);
}

.actions {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  flex-wrap: wrap;
}

:deep(.patients-card .p-card-body) {
  min-width: 0;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
}

:deep(.patients-card .p-card-content) {
  min-width: 0;
  width: 100%;
  max-width: 100%;
}

.table-empty {
  padding: 1rem;
  color: var(--p-surface-600);
  font-size: 0.82rem;
}

.dialog-grid {
  display: grid;
  gap: 0.82rem;
}

.dialog-subtitle {
  margin: 0;
  font-size: 0.82rem;
  color: var(--p-surface-600);
  overflow-wrap: anywhere;
}

.field {
  display: grid;
  gap: 0.36rem;
}

.field span {
  font-size: 0.82rem;
  color: var(--p-surface-700);
  font-weight: 600;
}

.field-inline {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: end;
  gap: 0.85rem;
}

.checkbox-field {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  min-height: 2.1rem;
  font-size: 0.8rem;
  color: var(--p-surface-700);
}

.help-text {
  margin: 0;
  font-size: 0.78rem;
  color: var(--p-surface-600);
}

.dialog-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.45rem;
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
  color: var(--p-red-600);
}

:deep(.p-menu-item-link.menu-action-item-danger:hover) {
  color: var(--p-red-700);
}

:deep(.p-card) {
  border: 1px solid var(--p-surface-200);
  box-shadow: none;
}

:deep(.p-datatable-header-cell) {
  background: var(--p-surface-50);
  color: var(--p-surface-700);
  white-space: nowrap;
}

:deep(.p-datatable-tbody > tr > td) {
  color: var(--p-surface-900);
  white-space: nowrap;
}

@media (max-width: 900px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-left {
    flex-wrap: wrap;
  }

  .toolbar-right {
    justify-content: flex-end;
    flex-wrap: wrap;
  }

  .column-toggle-wrap {
    width: 100%;
    min-width: 0;
  }

  .search-wrap {
    flex: 1 1 100%;
    max-width: 100%;
  }

  .filter-panel-grid {
    grid-template-columns: 1fr;
  }

  .table-footer {
    gap: 0.5rem;
  }

  .field-inline {
    grid-template-columns: 1fr;
    align-items: stretch;
    gap: 0.55rem;
  }
}
</style>
