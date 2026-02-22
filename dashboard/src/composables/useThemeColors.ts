import { computed, ref } from 'vue'
import { palette, updatePrimaryPalette, updateSurfacePalette } from '@primeuix/themes'

export type ThemePaletteName =
  | 'blue'
  | 'emerald'
  | 'teal'
  | 'cyan'
  | 'indigo'
  | 'violet'
  | 'purple'
  | 'rose'
  | 'orange'
  | 'amber'
  | 'slate'
  | 'zinc'
  | 'neutral'
  | 'stone'

export interface ThemeColors {
  primary: ThemePaletteName
  secondary: ThemePaletteName
}

export interface ThemePaletteOption {
  label: string
  value: ThemePaletteName
}

export type ThemePresetId =
  | 'clinical-calm'
  | 'premium-saas'
  | 'warm-accessible'
  | 'monochrome-edge'
  | 'vibrant-violet'
  | 'holistic-forest'
  | 'sunrise-coral'
  | 'royal-fintech'

export interface ThemePresetOption {
  id: ThemePresetId
  name: string
  description: string
  primary: ThemePaletteName
  secondary: ThemePaletteName
}

const THEME_STORAGE_KEY = 'borderform.theme.colors.v1'
const THEME_PRESET_STORAGE_KEY = 'borderform.theme.preset.v1'
const THEME_PRESET_ATTRIBUTE = 'data-dashboard-theme'
const PALETTE_ALIASES: Record<string, ThemePaletteName> = {
  bluegray: 'slate',
  'blue-gray': 'slate',
  blue_gray: 'slate',
}

export const THEME_PALETTE_OPTIONS: ThemePaletteOption[] = [
  { label: 'Blue', value: 'blue' },
  { label: 'Emerald', value: 'emerald' },
  { label: 'Teal', value: 'teal' },
  { label: 'Cyan', value: 'cyan' },
  { label: 'Indigo', value: 'indigo' },
  { label: 'Violet', value: 'violet' },
  { label: 'Purple', value: 'purple' },
  { label: 'Rose', value: 'rose' },
  { label: 'Orange', value: 'orange' },
  { label: 'Amber', value: 'amber' },
  { label: 'Slate', value: 'slate' },
  { label: 'Zinc', value: 'zinc' },
  { label: 'Neutral', value: 'neutral' },
  { label: 'Stone', value: 'stone' },
]

const PALETTE_NAME_SET = new Set<ThemePaletteName>(
  THEME_PALETTE_OPTIONS.map((option) => option.value),
)

export const DEFAULT_THEME_COLORS: ThemeColors = {
  primary: 'blue',
  secondary: 'slate',
}

export const DEFAULT_THEME_PRESET_ID: ThemePresetId = 'clinical-calm'

export const THEME_PRESET_OPTIONS: ThemePresetOption[] = [
  {
    id: 'clinical-calm',
    name: 'Clinical Calm',
    description: 'Professional, trustworthy slate-blues.',
    primary: 'blue',
    secondary: 'slate',
  },
  {
    id: 'premium-saas',
    name: 'Premium SaaS',
    description: 'High-contrast dark sidebar with indigo actions.',
    primary: 'indigo',
    secondary: 'zinc',
  },
  {
    id: 'warm-accessible',
    name: 'Warm & Accessible',
    description: 'Human-centric warm grays and deep teal.',
    primary: 'teal',
    secondary: 'stone',
  },
  {
    id: 'monochrome-edge',
    name: 'Monochrome Edge',
    description: 'Ultra-minimalist, stark black and white design.',
    primary: 'zinc',
    secondary: 'neutral',
  },
  {
    id: 'vibrant-violet',
    name: 'Vibrant Violet',
    description: 'Modern tech aesthetic with soft violet accents.',
    primary: 'violet',
    secondary: 'slate',
  },
  {
    id: 'holistic-forest',
    name: 'Holistic Forest',
    description: 'Calming sage and deep green natural tones.',
    primary: 'emerald',
    secondary: 'stone',
  },
  {
    id: 'sunrise-coral',
    name: 'Sunrise Coral',
    description: 'Warm, energetic beige and coral accents.',
    primary: 'orange',
    secondary: 'amber',
  },
  {
    id: 'royal-fintech',
    name: 'Royal Fintech',
    description: 'High-trust deep blue with striking cyan.',
    primary: 'cyan',
    secondary: 'slate',
  },
]

const THEME_PRESET_ID_SET = new Set<ThemePresetId>(
  THEME_PRESET_OPTIONS.map((themePreset) => themePreset.id),
)

const primaryColorRef = ref<ThemePaletteName>(DEFAULT_THEME_COLORS.primary)
const secondaryColorRef = ref<ThemePaletteName>(DEFAULT_THEME_COLORS.secondary)
const themePresetIdRef = ref<ThemePresetId | null>(DEFAULT_THEME_PRESET_ID)
let initialized = false

function normalizePaletteName(value: unknown, fallback: ThemePaletteName): ThemePaletteName {
  if (typeof value !== 'string') return fallback
  const normalized = value.trim().toLowerCase()
  if (!normalized) return fallback

  const aliased = PALETTE_ALIASES[normalized] ?? normalized
  return PALETTE_NAME_SET.has(aliased as ThemePaletteName)
    ? (aliased as ThemePaletteName)
    : fallback
}

function normalizeThemePresetId(value: unknown, fallback: ThemePresetId | null): ThemePresetId | null {
  if (typeof value !== 'string') return fallback
  const normalized = value.trim().toLowerCase()
  if (!normalized) return fallback

  return THEME_PRESET_ID_SET.has(normalized as ThemePresetId)
    ? (normalized as ThemePresetId)
    : fallback
}

function resolveThemePreset(themePresetId: ThemePresetId | null): ThemePresetOption | null {
  if (!themePresetId) return null
  return THEME_PRESET_OPTIONS.find((themePreset) => themePreset.id === themePresetId) ?? null
}

function findThemePresetByColors(colors: ThemeColors): ThemePresetOption | null {
  return (
    THEME_PRESET_OPTIONS.find(
      (themePreset) =>
        themePreset.primary === colors.primary && themePreset.secondary === colors.secondary,
    ) ?? null
  )
}

export function resolveThemePalette(name: ThemePaletteName) {
  return palette(`{${name}}`) as Record<string, string>
}

export function buildDarkSurfacePalette(lightSurfacePalette: Record<string, string>): Record<string, string> {
  return {
    0: lightSurfacePalette['950']!,
    50: lightSurfacePalette['900']!,
    100: lightSurfacePalette['800']!,
    200: lightSurfacePalette['700']!,
    300: lightSurfacePalette['600']!,
    400: lightSurfacePalette['500']!,
    500: lightSurfacePalette['400']!,
    600: lightSurfacePalette['300']!,
    700: lightSurfacePalette['200']!,
    800: lightSurfacePalette['100']!,
    900: lightSurfacePalette['50']!,
    950: lightSurfacePalette['0'] ?? lightSurfacePalette['50']!,
  }
}

function persistThemeColors(colors: ThemeColors) {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(THEME_STORAGE_KEY, JSON.stringify(colors))
  } catch {
    // Ignore write failures (private mode, disabled storage).
  }
}

function persistThemePresetId(themePresetId: ThemePresetId | null) {
  if (typeof window === 'undefined') return
  try {
    if (themePresetId) {
      window.localStorage.setItem(THEME_PRESET_STORAGE_KEY, themePresetId)
    } else {
      window.localStorage.removeItem(THEME_PRESET_STORAGE_KEY)
    }
  } catch {
    // Ignore write failures (private mode, disabled storage).
  }
}

function readStoredThemeColors(): ThemeColors | null {
  if (typeof window === 'undefined') return null

  try {
    const raw = window.localStorage.getItem(THEME_STORAGE_KEY)
    if (!raw) return null

    const parsed = JSON.parse(raw) as Partial<ThemeColors>
    return {
      primary: normalizePaletteName(parsed.primary, DEFAULT_THEME_COLORS.primary),
      secondary: normalizePaletteName(parsed.secondary, DEFAULT_THEME_COLORS.secondary),
    }
  } catch {
    return null
  }
}

function readStoredThemePresetId(): ThemePresetId | null {
  if (typeof window === 'undefined') return null
  try {
    const raw = window.localStorage.getItem(THEME_PRESET_STORAGE_KEY)
    if (!raw) return null
    return normalizeThemePresetId(raw, null)
  } catch {
    return null
  }
}

function applyThemePresetAttribute(themePresetId: ThemePresetId | null) {
  if (typeof document === 'undefined') return
  const resolvedThemePresetId = themePresetId ?? DEFAULT_THEME_PRESET_ID
  document.documentElement.setAttribute(THEME_PRESET_ATTRIBUTE, resolvedThemePresetId)
}

function applyPrimeThemeColors(colors: ThemeColors) {
  const primaryPalette = resolveThemePalette(colors.primary)
  const lightSurfacePalette = resolveThemePalette(colors.secondary)

  updatePrimaryPalette(primaryPalette)
  updateSurfacePalette({
    light: {
      0: lightSurfacePalette['0'] ?? lightSurfacePalette['50']!,
      ...lightSurfacePalette,
    },
    dark: buildDarkSurfacePalette(lightSurfacePalette),
  } as any)
}

function applyThemeColors(colors: ThemeColors, persist: boolean) {
  primaryColorRef.value = colors.primary
  secondaryColorRef.value = colors.secondary
  applyPrimeThemeColors(colors)
  if (persist) persistThemeColors(colors)
}

export function initializeThemeColors() {
  if (initialized) return
  initialized = true

  const storedThemePresetId = readStoredThemePresetId()
  const storedThemePreset = resolveThemePreset(storedThemePresetId)

  if (storedThemePreset) {
    themePresetIdRef.value = storedThemePreset.id
    applyThemeColors(
      {
        primary: storedThemePreset.primary,
        secondary: storedThemePreset.secondary,
      },
      false,
    )
    applyThemePresetAttribute(storedThemePreset.id)
    return
  }

  const storedColors = readStoredThemeColors()
  const resolvedColors = storedColors ?? DEFAULT_THEME_COLORS
  applyThemeColors(resolvedColors, false)

  const matchedThemePreset = findThemePresetByColors(resolvedColors)
  themePresetIdRef.value = matchedThemePreset?.id ?? null
  applyThemePresetAttribute(themePresetIdRef.value)
}

export function useThemeColors() {
  function setThemeColors(next: Partial<ThemeColors>) {
    const colors: ThemeColors = {
      primary: normalizePaletteName(next.primary ?? primaryColorRef.value, DEFAULT_THEME_COLORS.primary),
      secondary: normalizePaletteName(
        next.secondary ?? secondaryColorRef.value,
        DEFAULT_THEME_COLORS.secondary,
      ),
    }
    applyThemeColors(colors, true)
    const matchedThemePreset = findThemePresetByColors(colors)
    themePresetIdRef.value = matchedThemePreset?.id ?? null
    persistThemePresetId(themePresetIdRef.value)
    applyThemePresetAttribute(themePresetIdRef.value)
  }

  function setThemePreset(value: unknown) {
    const themePresetId = normalizeThemePresetId(value, null)
    const themePreset = resolveThemePreset(themePresetId)
    if (!themePreset) return

    themePresetIdRef.value = themePreset.id
    applyThemeColors(
      {
        primary: themePreset.primary,
        secondary: themePreset.secondary,
      },
      true,
    )
    persistThemePresetId(themePreset.id)
    applyThemePresetAttribute(themePreset.id)
  }

  function setPrimaryColor(value: unknown) {
    setThemeColors({ primary: normalizePaletteName(value, primaryColorRef.value) })
  }

  function setSecondaryColor(value: unknown) {
    setThemeColors({ secondary: normalizePaletteName(value, secondaryColorRef.value) })
  }

  function resetThemeColors() {
    const defaultThemePreset = resolveThemePreset(DEFAULT_THEME_PRESET_ID)
    if (!defaultThemePreset) {
      applyThemeColors(DEFAULT_THEME_COLORS, true)
      themePresetIdRef.value = null
      persistThemePresetId(null)
      applyThemePresetAttribute(null)
      return
    }

    themePresetIdRef.value = defaultThemePreset.id
    applyThemeColors(
      {
        primary: defaultThemePreset.primary,
        secondary: defaultThemePreset.secondary,
      },
      true,
    )
    persistThemePresetId(defaultThemePreset.id)
    applyThemePresetAttribute(defaultThemePreset.id)
  }

  return {
    paletteOptions: THEME_PALETTE_OPTIONS,
    themePresets: THEME_PRESET_OPTIONS,
    primaryColor: computed(() => primaryColorRef.value),
    secondaryColor: computed(() => secondaryColorRef.value),
    activeThemePresetId: computed(() => themePresetIdRef.value),
    themeColors: computed<ThemeColors>(() => ({
      primary: primaryColorRef.value,
      secondary: secondaryColorRef.value,
    })),
    setThemeColors,
    setThemePreset,
    setPrimaryColor,
    setSecondaryColor,
    resetThemeColors,
  }
}
