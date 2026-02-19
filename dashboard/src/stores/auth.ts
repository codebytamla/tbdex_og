import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

type ClaimList = string[] | string | null | undefined

interface BasicUserState {
  name: string
  role: string
}

export interface AuthProfile {
  id: string | null
  email: string | null
  name: string
  role: string
  roles: string[]
  permissions: string[]
  sites: string[]
  isAdmin: boolean
}

const LEGACY_AUTH_STORAGE_KEY = 'borderform.auth'
const ACCESS_TOKEN_STORAGE_KEY = 'borderform.auth.accessToken'
const REFRESH_TOKEN_STORAGE_KEY = 'borderform.auth.refreshToken'
const ACCESS_TOKEN_EXPIRES_AT_STORAGE_KEY = 'borderform.auth.accessTokenExpiresAt'
const BASIC_USER_STORAGE_KEY = 'borderform.auth.basicUser'
const PROFILE_STORAGE_KEY = 'borderform.auth.profile'
const SESSION_ACTIVE_STORAGE_KEY = 'borderform.auth.sessionActive'

const TOKEN_KEYS = ['accessToken', 'access_token', 'token', 'jwt'] as const
const REFRESH_TOKEN_KEYS = ['refreshToken', 'refresh_token'] as const
const EXPIRES_IN_KEYS = ['expires_in', 'expiresIn'] as const
const NAME_KEYS = [
  'fullName',
  'full_name',
  'name',
  'accountName',
  'account_name',
  'displayName',
  'display_name',
  'username',
  'userName',
] as const
const ROLE_KEYS = [
  'role',
  'roleName',
  'role_name',
  'accountRole',
  'account_role',
  'userRole',
  'user_role',
] as const
const EMAIL_KEYS = ['email', 'emailAddress', 'email_address'] as const
const ID_KEYS = ['id', 'userId', 'user_id', 'uuid'] as const
const IS_ADMIN_KEYS = ['isAdmin', 'admin'] as const
const ROLE_LIST_KEYS = ['roles', 'authorities', 'userRoles', 'user_roles'] as const
const PERMISSION_LIST_KEYS = [
  'permissions',
  'perms',
  'scopes',
  'abilities',
  'grants',
  'entitlements',
] as const
const SITE_LIST_KEYS = ['sites', 'facilities', 'locations', 'sites'] as const
const DEFAULT_ITEM_KEYS = ['name', 'slug', 'code', 'key', 'id', 'role', 'label', 'title'] as const
const SESSION_EXPIRED_MESSAGE = 'Session expired. Please sign in again.'
const API_V1_PREFIX = '/api/v1'

function toArray(value: ClaimList) {
  if (!value) return []
  return Array.isArray(value) ? value : [value]
}

function normalizeClaim(value: string) {
  return value.trim().toLowerCase()
}

function normalizeList(values: string[]) {
  const deduped: string[] = []
  const seen = new Set<string>()
  for (const value of values) {
    const normalized = value.trim()
    if (!normalized) continue
    const key = normalizeClaim(normalized)
    if (seen.has(key)) continue
    seen.add(key)
    deduped.push(normalized)
  }
  return deduped
}

function collectRecords(root: unknown, maxDepth = 4) {
  const records: Array<Record<string, unknown>> = []
  const queue: Array<[unknown, number]> = [[root, 0]]
  const seen = new Set<object>()

  while (queue.length > 0) {
    const [value, depth] = queue.shift() as [unknown, number]
    if (!value || typeof value !== 'object') continue
    if (seen.has(value)) continue
    seen.add(value)

    const record = value as Record<string, unknown>
    records.push(record)

    if (depth >= maxDepth) continue

    for (const nested of Object.values(record)) {
      if (nested && typeof nested === 'object') {
        queue.push([nested, depth + 1])
      }
    }
  }

  return records
}

function parseErrorMessage(payload: unknown, fallback = 'Unable to sign in. Please check your credentials.') {
  if (typeof payload === 'string') return payload

  if (payload && typeof payload === 'object') {
    const record = payload as Record<string, unknown>
    for (const key of ['message', 'error', 'detail']) {
      const value = record[key]
      if (typeof value === 'string' && value.trim()) return value.trim()
    }
  }

  return fallback
}

function findFirstString(records: Array<Record<string, unknown>>, keys: readonly string[]) {
  for (const record of records) {
    for (const key of keys) {
      const value = record[key]
      if (typeof value !== 'string') continue
      const normalized = value.trim()
      if (normalized) return normalized
    }
  }
  return null
}

function findFirstNumber(records: Array<Record<string, unknown>>, keys: readonly string[]) {
  for (const record of records) {
    for (const key of keys) {
      const value = record[key]
      if (typeof value === 'number' && Number.isFinite(value) && value > 0) return value
      if (typeof value === 'string') {
        const parsed = Number.parseInt(value.trim(), 10)
        if (Number.isFinite(parsed) && parsed > 0) return parsed
      }
    }
  }
  return null
}

function parseBooleanValue(value: unknown) {
  if (typeof value === 'boolean') return value
  if (typeof value === 'number' && Number.isFinite(value)) return value !== 0
  if (typeof value !== 'string') return null

  const normalized = value.trim().toLowerCase()
  if (!normalized) return null
  if (['true', '1', 'yes', 'y', 'on', 'active', 'enabled'].includes(normalized)) return true
  if (['false', '0', 'no', 'n', 'off', 'inactive', 'disabled'].includes(normalized)) return false
  return null
}

function findFirstBoolean(records: Array<Record<string, unknown>>, keys: readonly string[]) {
  for (const record of records) {
    for (const key of keys) {
      const parsed = parseBooleanValue(record[key])
      if (parsed !== null) return parsed
    }
  }
  return null
}

function decodeJwtExpiryMs(token: string | null) {
  if (!token) return null
  const sections = token.split('.')
  if (sections.length < 2) return null

  try {
    const payloadPart = sections[1]
    if (!payloadPart) return null
    const payloadBase64 = payloadPart.replace(/-/g, '+').replace(/_/g, '/')
    const padding = '='.repeat((4 - (payloadBase64.length % 4)) % 4)
    if (typeof atob !== 'function') return null
    const payloadText = atob(`${payloadBase64}${padding}`)
    const payload = JSON.parse(payloadText) as { exp?: unknown }
    if (typeof payload.exp !== 'number' || !Number.isFinite(payload.exp) || payload.exp <= 0) return null
    return payload.exp * 1000
  } catch {
    return null
  }
}

function extractStringArrayValue(
  value: unknown,
  itemKeys: readonly string[] = DEFAULT_ITEM_KEYS,
): string[] {
  if (typeof value === 'string') {
    const normalized = value.trim()
    return normalized ? [normalized] : []
  }

  if (!Array.isArray(value)) return []

  const output: string[] = []

  for (const item of value) {
    if (typeof item === 'string') {
      const normalized = item.trim()
      if (normalized) output.push(normalized)
      continue
    }

    if (!item || typeof item !== 'object') continue

    const record = item as Record<string, unknown>
    for (const key of itemKeys) {
      const candidate = record[key]
      if (typeof candidate === 'string' && candidate.trim()) {
        output.push(candidate.trim())
        break
      }
      if (typeof candidate === 'number') {
        output.push(String(candidate))
        break
      }
    }
  }

  return normalizeList(output)
}

function findFirstStringArray(
  records: Array<Record<string, unknown>>,
  keys: readonly string[],
  itemKeys: readonly string[] = DEFAULT_ITEM_KEYS,
) {
  for (const record of records) {
    for (const key of keys) {
      const values = extractStringArrayValue(record[key], itemKeys)
      if (values.length > 0) return values
    }
  }
  return []
}

function resolvePrimaryRole(profile: AuthProfile | null, fallbackRole = 'Role') {
  if (!profile) return fallbackRole
  const primary = profile.role.trim()
  if (primary && normalizeClaim(primary) !== 'role') return primary
  return profile.roles[0] || fallbackRole
}

function inferIsAdminFromRoles(role: string, roles: string[]) {
  const normalizedRoles = normalizeList([role, ...roles]).map(normalizeClaim)
  return normalizedRoles.some((entry) => entry === 'admin' || entry === 'administrator')
}

function createEmptyProfileFallback() {
  return { name: 'Account name', role: 'Role' }
}

function buildProfile(payload: unknown, fallback: BasicUserState): AuthProfile {
  const records = collectRecords(payload, 5)
  const name = findFirstString(records, NAME_KEYS) || fallback.name || 'Account name'
  const explicitRole = findFirstString(records, ROLE_KEYS)
  const id = findFirstString(records, ID_KEYS)
  const email = findFirstString(records, EMAIL_KEYS)
  const roles = findFirstStringArray(records, ROLE_LIST_KEYS, ['name', 'role', 'slug', 'code'])
  const role = explicitRole || roles[0] || fallback.role || 'Role'
  if (role && !roles.some((entry) => normalizeClaim(entry) === normalizeClaim(role))) {
    roles.unshift(role)
  }
  const normalizedRoles = normalizeList(roles)
  const isAdmin = findFirstBoolean(records, IS_ADMIN_KEYS) ?? inferIsAdminFromRoles(role, normalizedRoles)

  return {
    id,
    email,
    name,
    role,
    roles: normalizedRoles,
    permissions: findFirstStringArray(records, PERMISSION_LIST_KEYS),
    sites: findFirstStringArray(records, SITE_LIST_KEYS, ['name', 'label', 'slug', 'code', 'id']),
    isAdmin,
  }
}

function matchesClaim(granted: string, required: string) {
  const normalizedGranted = normalizeClaim(granted)
  const normalizedRequired = normalizeClaim(required)

  if (!normalizedGranted || !normalizedRequired) return false
  if (normalizedGranted === normalizedRequired) return true
  if (normalizedGranted === '*' || normalizedGranted === 'all') return true
  if (normalizedGranted.endsWith('.*')) {
    const prefix = normalizedGranted.slice(0, -1)
    if (normalizedRequired.startsWith(prefix)) return true
  }
  if (normalizedRequired.endsWith('.*')) {
    const prefix = normalizedRequired.slice(0, -1)
    if (normalizedGranted.startsWith(prefix)) return true
  }
  return false
}

function parseStoredProfile(raw: string | null) {
  if (!raw) return null
  try {
    const parsed = JSON.parse(raw) as Partial<AuthProfile>
    if (!parsed || typeof parsed !== 'object') return null

    const parsedRecord = parsed as Record<string, unknown>
    const isAdmin = parseBooleanValue(parsedRecord.isAdmin) ?? inferIsAdminFromRoles(
      typeof parsed.role === 'string' ? parsed.role : 'Role',
      Array.isArray(parsed.roles)
        ? parsed.roles.filter((entry): entry is string => typeof entry === 'string')
        : [],
    )

    const profile: AuthProfile = {
      id: typeof parsed.id === 'string' ? parsed.id : null,
      email: typeof parsed.email === 'string' ? parsed.email : null,
      name: typeof parsed.name === 'string' && parsed.name.trim() ? parsed.name.trim() : 'Account name',
      role: typeof parsed.role === 'string' && parsed.role.trim() ? parsed.role.trim() : 'Role',
      roles: Array.isArray(parsed.roles)
        ? normalizeList(parsed.roles.filter((entry): entry is string => typeof entry === 'string'))
        : [],
      permissions: Array.isArray(parsed.permissions)
        ? normalizeList(parsed.permissions.filter((entry): entry is string => typeof entry === 'string'))
        : [],
      sites: Array.isArray(parsed.sites)
        ? normalizeList(parsed.sites.filter((entry): entry is string => typeof entry === 'string'))
        : [],
      isAdmin,
    }

    if (
      !profile.roles.some((entry) => normalizeClaim(entry) === normalizeClaim(profile.role)) &&
      profile.role !== 'Role'
    ) {
      profile.roles.unshift(profile.role)
    }

    return profile
  } catch {
    return null
  }
}

function parseStoredBasicUser(raw: string | null) {
  if (!raw) return createEmptyProfileFallback()

  try {
    const parsed = JSON.parse(raw) as Partial<BasicUserState>
    return {
      name: typeof parsed.name === 'string' && parsed.name.trim() ? parsed.name.trim() : 'Account name',
      role: typeof parsed.role === 'string' && parsed.role.trim() ? parsed.role.trim() : 'Role',
    }
  } catch {
    return createEmptyProfileFallback()
  }
}

function parseStoredSessionFlag(raw: string | null) {
  return raw === '1'
}

function parseStoredString(raw: string | null) {
  if (typeof raw !== 'string') return null
  const normalized = raw.trim()
  return normalized || null
}

function parseStoredNumber(raw: string | null) {
  if (typeof raw !== 'string') return null
  const parsed = Number.parseInt(raw.trim(), 10)
  if (!Number.isFinite(parsed) || parsed <= 0) return null
  return parsed
}

function parseLegacyAuthState(raw: string | null) {
  if (!raw) return null

  try {
    const parsed = JSON.parse(raw) as {
      sessionActive?: unknown
      accessToken?: unknown
      accountName?: unknown
      accountRole?: unknown
    }
    return {
      sessionActive: Boolean(parsed.sessionActive),
      accessToken:
        typeof parsed.accessToken === 'string' && parsed.accessToken.trim()
          ? parsed.accessToken.trim()
          : null,
      accountName:
        typeof parsed.accountName === 'string' && parsed.accountName.trim()
          ? parsed.accountName.trim()
          : 'Account name',
      accountRole:
        typeof parsed.accountRole === 'string' && parsed.accountRole.trim()
          ? parsed.accountRole.trim()
          : 'Role',
    }
  } catch {
    return null
  }
}

function apiUrl(path: string) {
  const baseUrl = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/+$/, '')
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const versionedPath =
    normalizedPath === API_V1_PREFIX || normalizedPath.startsWith(`${API_V1_PREFIX}/`)
      ? normalizedPath
      : `${API_V1_PREFIX}${normalizedPath}`
  if (baseUrl) return `${baseUrl}${versionedPath}`
  return versionedPath
}

function createUnauthorizedResponse(message = SESSION_EXPIRED_MESSAGE) {
  return new Response(JSON.stringify({ message }), {
    status: 401,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
  })
}

export const useAuthStore = defineStore('auth', () => {
  const hydrated = ref(false)
  const loading = ref(false)
  const refreshing = ref(false)
  const profileLoading = ref(false)
  const profileLoaded = ref(false)
  const errorMessage = ref('')

  const sessionActive = ref(false)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const accessTokenExpiresAt = ref<number | null>(null)
  const basicUser = ref<BasicUserState>(createEmptyProfileFallback())
  const profile = ref<AuthProfile | null>(null)

  let profileRequest: Promise<boolean> | null = null
  let refreshRequest: Promise<boolean> | null = null

  const accountName = computed(() => profile.value?.name || basicUser.value.name || 'Account name')
  const accountRole = computed(() => resolvePrimaryRole(profile.value, basicUser.value.role || 'Role'))
  const roles = computed(() => {
    if (profile.value) return normalizeList(profile.value.roles.length ? profile.value.roles : [accountRole.value])
    if (accountRole.value.trim()) return [accountRole.value]
    return []
  })
  const permissions = computed(() => profile.value?.permissions || [])
  const sites = computed(() => profile.value?.sites || [])
  const isAdmin = computed(() => {
    if (profile.value) return profile.value.isAdmin
    return inferIsAdminFromRoles(accountRole.value, roles.value)
  })
  const isAuthenticated = computed(
    () => sessionActive.value || Boolean(accessToken.value) || Boolean(refreshToken.value),
  )

  function defaultDashboardRoute() {
    return '/dashboard'
  }

  function persistSessionFlag() {
    if (sessionActive.value) {
      localStorage.setItem(SESSION_ACTIVE_STORAGE_KEY, '1')
      return
    }
    localStorage.removeItem(SESSION_ACTIVE_STORAGE_KEY)
  }

  function persistAccessToken() {
    if (accessToken.value) {
      sessionStorage.setItem(ACCESS_TOKEN_STORAGE_KEY, accessToken.value)
      return
    }
    sessionStorage.removeItem(ACCESS_TOKEN_STORAGE_KEY)
  }

  function persistRefreshToken() {
    if (refreshToken.value) {
      localStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, refreshToken.value)
      return
    }
    localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY)
  }

  function persistAccessTokenExpiry() {
    if (accessTokenExpiresAt.value && Number.isFinite(accessTokenExpiresAt.value)) {
      localStorage.setItem(ACCESS_TOKEN_EXPIRES_AT_STORAGE_KEY, String(accessTokenExpiresAt.value))
      return
    }
    localStorage.removeItem(ACCESS_TOKEN_EXPIRES_AT_STORAGE_KEY)
  }

  function persistBasicUser() {
    localStorage.setItem(BASIC_USER_STORAGE_KEY, JSON.stringify(basicUser.value))
  }

  function persistProfile() {
    if (profile.value) {
      localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify(profile.value))
      return
    }
    localStorage.removeItem(PROFILE_STORAGE_KEY)
  }

  function clearAuthState(options?: { keepError?: boolean }) {
    sessionActive.value = false
    accessToken.value = null
    refreshToken.value = null
    accessTokenExpiresAt.value = null
    basicUser.value = createEmptyProfileFallback()
    profile.value = null
    profileLoaded.value = false
    profileRequest = null
    refreshRequest = null

    persistSessionFlag()
    persistAccessToken()
    persistRefreshToken()
    persistAccessTokenExpiry()
    persistBasicUser()
    persistProfile()

    if (!options?.keepError) errorMessage.value = ''
  }

  function ingestLoginPayload(payload: unknown) {
    const records = collectRecords(payload, 5)
    const nextAccessToken = findFirstString(records, TOKEN_KEYS) || null
    const nextRefreshToken = findFirstString(records, REFRESH_TOKEN_KEYS) || null
    const expiresIn = findFirstNumber(records, EXPIRES_IN_KEYS)

    accessToken.value = nextAccessToken
    refreshToken.value = nextRefreshToken
    if (expiresIn && expiresIn > 0) {
      accessTokenExpiresAt.value = Date.now() + expiresIn * 1000
    } else {
      accessTokenExpiresAt.value = decodeJwtExpiryMs(nextAccessToken)
    }
    basicUser.value = {
      name: findFirstString(records, NAME_KEYS) || 'Account name',
      role: findFirstString(records, ROLE_KEYS) || 'Role',
    }
    sessionActive.value = true

    persistAccessToken()
    persistRefreshToken()
    persistAccessTokenExpiry()
    persistBasicUser()
    persistSessionFlag()
  }

  function accessTokenNeedsRefresh(bufferMs = 10_000) {
    if (!accessToken.value) return true
    if (!accessTokenExpiresAt.value) return false
    return Date.now() + bufferMs >= accessTokenExpiresAt.value
  }

  async function refreshAccessToken(force = false) {
    hydrate()

    if (refreshRequest) return refreshRequest
    if (!force && !accessTokenNeedsRefresh()) return true
    if (!refreshToken.value && !sessionActive.value && !accessToken.value) return false

    refreshRequest = (async () => {
      refreshing.value = true

      try {
        const requestBody: Record<string, string> = {}
        if (refreshToken.value) {
          requestBody.refresh_token = refreshToken.value
          requestBody.refreshToken = refreshToken.value
        }

        const response = await fetch(apiUrl('/auth/refresh'), {
          method: 'POST',
          headers: {
            Accept: 'application/json',
            ...(Object.keys(requestBody).length > 0 ? { 'Content-Type': 'application/json' } : {}),
          },
          credentials: 'include',
          body: Object.keys(requestBody).length > 0 ? JSON.stringify(requestBody) : undefined,
        })

        const payload = (await response.json().catch(() => ({}))) as unknown
        if (!response.ok) {
          const message = parseErrorMessage(payload, SESSION_EXPIRED_MESSAGE)
          const error = new Error(message) as Error & { status?: number }
          error.status = response.status
          throw error
        }

        const records = collectRecords(payload, 5)
        const nextAccessToken = findFirstString(records, TOKEN_KEYS)
        if (!nextAccessToken) {
          throw new Error('Refresh response missing access token.')
        }

        accessToken.value = nextAccessToken
        refreshToken.value = findFirstString(records, REFRESH_TOKEN_KEYS) || refreshToken.value
        const expiresIn = findFirstNumber(records, EXPIRES_IN_KEYS)
        accessTokenExpiresAt.value = expiresIn
          ? Date.now() + expiresIn * 1000
          : decodeJwtExpiryMs(nextAccessToken)
        sessionActive.value = true

        persistAccessToken()
        persistRefreshToken()
        persistAccessTokenExpiry()
        persistSessionFlag()
        return true
      } catch (error) {
        const status =
          typeof (error as { status?: unknown })?.status === 'number'
            ? ((error as { status: number }).status as number)
            : null
        errorMessage.value = error instanceof Error ? error.message : SESSION_EXPIRED_MESSAGE

        const shouldResetSession =
          status === 401 ||
          status === 403 ||
          status === 419 ||
          status === 440 ||
          (!accessToken.value && !refreshToken.value)

        if (shouldResetSession) {
          clearAuthState({ keepError: true })
        }
        return false
      } finally {
        refreshing.value = false
        refreshRequest = null
      }
    })()

    return refreshRequest
  }

  async function authenticatedFetch(
    path: string,
    init: RequestInit = {},
    options?: { retryOnUnauthorized?: boolean; forceRefresh?: boolean },
  ) {
    hydrate()

    const retryOnUnauthorized = options?.retryOnUnauthorized ?? true
    const forceRefresh = options?.forceRefresh ?? false

    if (isAuthenticated.value && (forceRefresh || accessTokenNeedsRefresh())) {
      const refreshed = await refreshAccessToken(forceRefresh)
      if (!refreshed && !accessToken.value) {
        return createUnauthorizedResponse()
      }
    }

    if (isAuthenticated.value && !accessToken.value) {
      const recovered = await refreshAccessToken(true)
      if (!recovered || !accessToken.value) {
        return createUnauthorizedResponse()
      }
    }

    const execute = async () => {
      const headers = new Headers(init.headers)
      if (!headers.has('Accept')) headers.set('Accept', 'application/json')
      if (accessToken.value) headers.set('Authorization', `Bearer ${accessToken.value}`)

      return fetch(apiUrl(path), {
        ...init,
        headers,
        credentials: init.credentials ?? 'include',
      })
    }

    let response = await execute()

    if (response.status === 401 && retryOnUnauthorized) {
      const refreshed = await refreshAccessToken(true)
      if (refreshed && accessToken.value) response = await execute()
    }

    return response
  }

  async function authenticatedRequestJson<T>(
    path: string,
    init: RequestInit,
    fallbackError: string,
    options?: { retryOnUnauthorized?: boolean; forceRefresh?: boolean },
  ) {
    const response = await authenticatedFetch(path, init, options)
    const payload = (await response.json().catch(() => ({}))) as unknown
    if (!response.ok) {
      throw new Error(parseErrorMessage(payload, fallbackError))
    }
    return payload as T
  }

  async function fetchProfile(options?: { silent?: boolean; force?: boolean }) {
    const silent = options?.silent ?? false
    const force = options?.force ?? false

    if (profileRequest && !force) return profileRequest
    if (profile.value && profileLoaded.value && !force) return true

    profileRequest = (async () => {
      profileLoading.value = true

      try {
        const payload = await authenticatedRequestJson<unknown>(
          '/auth/me',
          { method: 'GET' },
          'Unable to load account profile.',
        )

        const normalizedProfile = buildProfile(payload, basicUser.value)
        profile.value = normalizedProfile
        basicUser.value = {
          name: normalizedProfile.name,
          role: resolvePrimaryRole(normalizedProfile, normalizedProfile.role),
        }
        sessionActive.value = true
        profileLoaded.value = true

        persistProfile()
        persistBasicUser()
        persistSessionFlag()
        return true
      } catch (error) {
        profileLoaded.value = false
        if (!silent) {
          errorMessage.value =
            error instanceof Error ? error.message : 'Unable to load account profile.'
        }
        return false
      } finally {
        profileLoading.value = false
        profileRequest = null
      }
    })()

    return profileRequest
  }

  async function ensureProfile(force = false) {
    hydrate()
    if (!isAuthenticated.value) return false
    return fetchProfile({ silent: true, force })
  }

  function hydrate() {
    if (hydrated.value) return
    hydrated.value = true

    const legacy = parseLegacyAuthState(localStorage.getItem(LEGACY_AUTH_STORAGE_KEY))
    if (legacy) {
      if (legacy.accessToken) {
        accessToken.value = legacy.accessToken
        persistAccessToken()
      }
      basicUser.value = { name: legacy.accountName, role: legacy.accountRole }
      sessionActive.value = legacy.sessionActive
      persistBasicUser()
      persistSessionFlag()
      localStorage.removeItem(LEGACY_AUTH_STORAGE_KEY)
    }

    accessToken.value = parseStoredString(sessionStorage.getItem(ACCESS_TOKEN_STORAGE_KEY))
    refreshToken.value = parseStoredString(localStorage.getItem(REFRESH_TOKEN_STORAGE_KEY))
    accessTokenExpiresAt.value = parseStoredNumber(
      localStorage.getItem(ACCESS_TOKEN_EXPIRES_AT_STORAGE_KEY),
    )
    basicUser.value = parseStoredBasicUser(localStorage.getItem(BASIC_USER_STORAGE_KEY))
    profile.value = parseStoredProfile(localStorage.getItem(PROFILE_STORAGE_KEY))
    profileLoaded.value = Boolean(profile.value)
    sessionActive.value =
      parseStoredSessionFlag(localStorage.getItem(SESSION_ACTIVE_STORAGE_KEY)) ||
      Boolean(accessToken.value) ||
      Boolean(refreshToken.value)
  }

  function hasPermission(permission: string) {
    if (!permission.trim()) return false
    return permissions.value.some((entry) => matchesClaim(entry, permission))
  }

  function hasRole(role: string) {
    const normalizedRole = normalizeClaim(role)
    if (!normalizedRole) return false
    return roles.value.some((entry) => normalizeClaim(entry) === normalizedRole)
  }

  function canAccess(requiredPermissions?: ClaimList, requiredRoles?: ClaimList) {
    const permissionList = toArray(requiredPermissions)
    const roleList = toArray(requiredRoles)
    const hasPermissionRequirement = permissionList.length > 0
    const hasRoleRequirement = roleList.length > 0

    if (!hasPermissionRequirement && !hasRoleRequirement) {
      return true
    }

    const permissionAllowed = hasPermissionRequirement
      ? permissionList.some((permission) => hasPermission(permission))
      : false
    const roleAllowed = hasRoleRequirement ? roleList.some((role) => hasRole(role)) : false

    if (hasPermissionRequirement && hasRoleRequirement) {
      return permissionAllowed || roleAllowed
    }

    if (hasPermissionRequirement) return permissionAllowed
    return roleAllowed
  }

  async function login(email: string, password: string) {
    loading.value = true
    errorMessage.value = ''

    try {
      const response = await fetch(apiUrl('/auth/login'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ email, password }),
      })

      const payload = (await response.json().catch(() => ({}))) as unknown
      if (!response.ok) throw new Error(parseErrorMessage(payload))

      ingestLoginPayload(payload)

      const loaded = await fetchProfile({ force: true })
      if (!loaded) throw new Error(errorMessage.value || 'Unable to load account profile.')

      return true
    } catch (error) {
      errorMessage.value =
        error instanceof Error
          ? error.message
          : 'Unable to sign in. Please check your credentials.'
      clearAuthState({ keepError: true })
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    clearAuthState()
  }

  return {
    hydrated,
    loading,
    refreshing,
    profileLoading,
    profileLoaded,
    errorMessage,
    accessToken,
    refreshToken,
    accessTokenExpiresAt,
    profile,
    accountName,
    accountRole,
    roles,
    permissions,
    sites,
    isAdmin,
    isAuthenticated,
    defaultDashboardRoute,
    hydrate,
    ensureProfile,
    fetchProfile,
    refreshAccessToken,
    authenticatedFetch,
    authenticatedRequestJson,
    hasPermission,
    hasRole,
    canAccess,
    login,
    logout,
  }
})
