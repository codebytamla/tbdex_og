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

const THEME_STORAGE_KEY = 'borderform.theme.colors.v1'
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

const primaryColorRef = ref<ThemePaletteName>(DEFAULT_THEME_COLORS.primary)
const secondaryColorRef = ref<ThemePaletteName>(DEFAULT_THEME_COLORS.secondary)
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

export function resolveThemePalette(name: ThemePaletteName) {
  return palette(`{${name}}`) as Record<string, string>
}

export function buildDarkSurfacePalette(lightSurfacePalette: Record<string, string>): Record<string, string> {
  return {
    0: lightSurfacePalette['950'] || '#020617',
    50: lightSurfacePalette['900'] || '#0f172a',
    100: lightSurfacePalette['800'] || '#1e293b',
    200: lightSurfacePalette['700'] || '#334155',
    300: lightSurfacePalette['600'] || '#475569',
    400: lightSurfacePalette['500'] || '#64748b',
    500: lightSurfacePalette['400'] || '#94a3b8',
    600: lightSurfacePalette['300'] || '#cbd5e1',
    700: lightSurfacePalette['200'] || '#e2e8f0',
    800: lightSurfacePalette['100'] || '#f1f5f9',
    900: lightSurfacePalette['50'] || '#f8fafc',
    950: '#ffffff',
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

function applyPrimeThemeColors(colors: ThemeColors) {
  const primaryPalette = resolveThemePalette(colors.primary)
  const lightSurfacePalette = resolveThemePalette(colors.secondary)

  updatePrimaryPalette(primaryPalette)
  updateSurfacePalette({
    light: {
      0: '#ffffff',
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

  const stored = readStoredThemeColors()
  applyThemeColors(stored ?? DEFAULT_THEME_COLORS, false)
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
  }

  function setPrimaryColor(value: unknown) {
    setThemeColors({ primary: normalizePaletteName(value, primaryColorRef.value) })
  }

  function setSecondaryColor(value: unknown) {
    setThemeColors({ secondary: normalizePaletteName(value, secondaryColorRef.value) })
  }

  function resetThemeColors() {
    applyThemeColors(DEFAULT_THEME_COLORS, true)
  }

  return {
    paletteOptions: THEME_PALETTE_OPTIONS,
    primaryColor: computed(() => primaryColorRef.value),
    secondaryColor: computed(() => secondaryColorRef.value),
    themeColors: computed<ThemeColors>(() => ({
      primary: primaryColorRef.value,
      secondary: secondaryColorRef.value,
    })),
    setThemeColors,
    setPrimaryColor,
    setSecondaryColor,
    resetThemeColors,
  }
}
