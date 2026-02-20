import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

export type SiteCatalogEntry = {
  id: string
  name: string
  code: string
}

export type RoleCatalogEntry = {
  id: string
  name: string
  description: string
}

export type PermissionCatalogEntry = {
  id?: string
  name: string
}

const STALE_MS = 60_000

function isStale(timestamp: number | null): boolean {
  if (!timestamp) return true
  return Date.now() - timestamp > STALE_MS
}

function extractCollection(payload: unknown, keys: string[]): unknown[] {
  if (Array.isArray(payload)) return payload
  if (!payload || typeof payload !== 'object') return []
  const record = payload as Record<string, unknown>
  for (const key of keys) {
    const value = record[key]
    if (Array.isArray(value)) return value
  }
  for (const value of Object.values(record)) {
    if (Array.isArray(value)) return value
  }
  return []
}

function normalizeSiteEntry(value: unknown): SiteCatalogEntry | null {
  if (!value || typeof value !== 'object') return null
  const r = value as Record<string, unknown>
  const id = typeof r.id === 'string' ? r.id : null
  const name = typeof r.name === 'string' ? r.name.trim() : ''
  const code = typeof r.code === 'string' ? r.code.trim() : ''
  if (!id || !name) return null
  return { id, name, code }
}

function normalizeRoleEntry(value: unknown): RoleCatalogEntry | null {
  if (!value || typeof value !== 'object') return null
  const r = value as Record<string, unknown>
  const id = typeof r.id === 'string' ? r.id : null
  const name = typeof r.name === 'string' ? r.name.trim() : ''
  const description = typeof r.description === 'string' ? r.description.trim() : ''
  if (!id || !name) return null
  return { id, name, description }
}

function normalizePermissionName(value: unknown): string | null {
  if (typeof value === 'string') return value.trim() || null
  if (!value || typeof value !== 'object') return null
  const r = value as Record<string, unknown>
  for (const key of ['name', 'permission', 'slug', 'code']) {
    const candidate = r[key]
    if (typeof candidate === 'string' && candidate.trim()) return candidate.trim()
  }
  return null
}

function normalizePermissionEntry(value: unknown): PermissionCatalogEntry | null {
  const name = normalizePermissionName(value)
  if (!name) return null
  const id =
    value && typeof value === 'object' && typeof (value as Record<string, unknown>).id === 'string'
      ? (value as Record<string, unknown>).id as string
      : undefined
  return { id, name }
}

export const useCatalogStore = defineStore('catalog', () => {
  const sites = ref<SiteCatalogEntry[]>([])
  const sitesLoading = ref(false)
  const sitesFetchedAt = ref<number | null>(null)

  const roles = ref<RoleCatalogEntry[]>([])
  const rolesLoading = ref(false)
  const rolesFetchedAt = ref<number | null>(null)

  const permissions = ref<PermissionCatalogEntry[]>([])
  const permissionsLoading = ref(false)
  const permissionsFetchedAt = ref<number | null>(null)

  let sitesFetchPromise: Promise<void> | null = null
  let rolesFetchPromise: Promise<void> | null = null
  let permissionsFetchPromise: Promise<void> | null = null

  async function fetchSites(options?: { force?: boolean }) {
    if (sitesFetchPromise && !options?.force) return sitesFetchPromise
    if (!options?.force && !isStale(sitesFetchedAt.value) && sites.value.length > 0) return

    sitesFetchPromise = (async () => {
      const authStore = useAuthStore()
      sitesLoading.value = true
      try {
        const payload = await authStore.authenticatedRequestJson<unknown>(
          '/admin/sites',
          { method: 'GET' },
          'Unable to load sites catalog.',
        )
        sites.value = extractCollection(payload, ['sites'])
          .map(normalizeSiteEntry)
          .filter((e): e is SiteCatalogEntry => !!e)
        sitesFetchedAt.value = Date.now()
      } finally {
        sitesLoading.value = false
        sitesFetchPromise = null
      }
    })()

    return sitesFetchPromise
  }

  async function fetchRoles(options?: { force?: boolean }) {
    if (rolesFetchPromise && !options?.force) return rolesFetchPromise
    if (!options?.force && !isStale(rolesFetchedAt.value) && roles.value.length > 0) return

    rolesFetchPromise = (async () => {
      const authStore = useAuthStore()
      rolesLoading.value = true
      try {
        const payload = await authStore.authenticatedRequestJson<unknown>(
          '/admin/roles',
          { method: 'GET' },
          'Unable to load roles catalog.',
        )
        roles.value = extractCollection(payload, ['roles'])
          .map(normalizeRoleEntry)
          .filter((e): e is RoleCatalogEntry => !!e)
        rolesFetchedAt.value = Date.now()
      } finally {
        rolesLoading.value = false
        rolesFetchPromise = null
      }
    })()

    return rolesFetchPromise
  }

  async function fetchPermissions(options?: { force?: boolean }) {
    if (permissionsFetchPromise && !options?.force) return permissionsFetchPromise
    if (!options?.force && !isStale(permissionsFetchedAt.value) && permissions.value.length > 0) return

    permissionsFetchPromise = (async () => {
      const authStore = useAuthStore()
      permissionsLoading.value = true
      try {
        const payload = await authStore.authenticatedRequestJson<{ permissions?: unknown[] }>(
          '/admin/permissions',
          { method: 'GET' },
          'Unable to load permissions catalog.',
        )
        const byName = new Map<string, PermissionCatalogEntry>()
        for (const item of payload.permissions || []) {
          const entry = normalizePermissionEntry(item)
          if (entry) byName.set(entry.name.toLowerCase(), entry)
        }
        permissions.value = [...byName.values()]
        permissionsFetchedAt.value = Date.now()
      } finally {
        permissionsLoading.value = false
        permissionsFetchPromise = null
      }
    })()

    return permissionsFetchPromise
  }

  function invalidateSites() {
    sitesFetchedAt.value = null
  }

  function invalidateRoles() {
    rolesFetchedAt.value = null
  }

  function invalidatePermissions() {
    permissionsFetchedAt.value = null
  }

  return {
    sites,
    sitesLoading,
    roles,
    rolesLoading,
    permissions,
    permissionsLoading,
    fetchSites,
    fetchRoles,
    fetchPermissions,
    invalidateSites,
    invalidateRoles,
    invalidatePermissions,
  }
})
