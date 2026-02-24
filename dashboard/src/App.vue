<script setup lang="ts">
import type { MenuItem } from 'primevue/menuitem'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MoreVertical, PanelLeftClose, PanelLeftOpen } from 'lucide-vue-next'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import PanelMenu from 'primevue/panelmenu'
import ConfirmDialog from 'primevue/confirmdialog'
import Toast from 'primevue/toast'
import { useAuthStore } from '@/stores/auth'
import { useNavItems } from '@/composables/useNavItems'
import { useThemeColors, type ThemePresetId } from '@/composables/useThemeColors'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { visibleNavItems } = useNavItems()
const themeColors = useThemeColors()
const sidebarCollapseStorageKey = 'borderform.ui.sidebarCollapsed'
const sidebarColorModeStorageKey = 'borderform.ui.sidebarDark'
const sidebarCollapsed = ref(false)
const sidebarDarkMode = ref(false)
authStore.hydrate()

const isAuthRoute = computed(() => route.path.startsWith('/auth') || route.path === '/reset-password')
const sidebarToggleLabel = computed(() =>
  sidebarCollapsed.value ? 'Expand navigation sidebar' : 'Collapse navigation sidebar',
)
const sidebarToggleIcon = computed(() => (sidebarCollapsed.value ? PanelLeftOpen : PanelLeftClose))
const accountName = computed(() => authStore.accountName)
const accountRole = computed(() => {
  const primary = authStore.accountRole?.trim()
  if (primary && primary.toLowerCase() !== 'role') return primary
  return authStore.roles[0] || 'No role assigned'
})
const accountInitial = computed(() => accountName.value.trim().charAt(0).toUpperCase() || 'A')
const themePresets = computed(() => themeColors.themePresets)
const activeThemePresetId = computed(() => themeColors.activeThemePresetId.value)

const accountActions = computed<MenuItem[]>(() => {
  const activePresetId = activeThemePresetId.value
  const appearanceActions: MenuItem[] = [
    {
      label: `Dark navbar: ${sidebarDarkMode.value ? 'On' : 'Off'}`,
      class: sidebarDarkMode.value
        ? 'account-menu-navbar-toggle account-menu-navbar-toggle-active'
        : 'account-menu-navbar-toggle',
      command: () => toggleSidebarColorMode(),
    },
    {
      separator: true,
    },
    ...themePresets.value.map((themePreset) => ({
      label: activePresetId === themePreset.id ? `${themePreset.name} (Active)` : themePreset.name,
      class:
        activePresetId === themePreset.id
          ? 'account-menu-theme-item account-menu-theme-item-active'
          : 'account-menu-theme-item',
      command: () => setThemePreset(themePreset.id),
    })),
  ]

  return [
    {
      label: 'Settings',
      command: async () => {
        await router.push('/settings')
      },
    },
    {
      separator: true,
    },
    {
      label: 'Appearance',
      items: appearanceActions,
    },
    {
      separator: true,
    },
    {
      label: 'Logout',
      command: async () => {
        authStore.logout()
        await router.push('/auth/login')
      },
    },
  ]
})
const accountMenu = ref()
const isNavbarDark = computed(() => sidebarDarkMode.value)

onMounted(() => {
  sidebarCollapsed.value = localStorage.getItem(sidebarCollapseStorageKey) === '1'
  sidebarDarkMode.value = localStorage.getItem(sidebarColorModeStorageKey) === '1'
})

function toggleAccountMenu(event: Event) {
  accountMenu.value?.toggle(event)
}

function setThemePreset(themePresetId: ThemePresetId) {
  themeColors.setThemePreset(themePresetId)
}

function isNavMenuItemActive(item: MenuItem) {
  return typeof item.class === 'string' && item.class.includes('nav-menu-item-active')
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem(sidebarCollapseStorageKey, sidebarCollapsed.value ? '1' : '0')
}

function toggleSidebarColorMode() {
  sidebarDarkMode.value = !sidebarDarkMode.value
  localStorage.setItem(sidebarColorModeStorageKey, sidebarDarkMode.value ? '1' : '0')
}
</script>

<template>
  <Toast position="top-right" />
  <ConfirmDialog />

  <RouterView v-if="isAuthRoute" />

  <div v-else class="layout-shell" :class="{ 'layout-shell-collapsed': sidebarCollapsed }">
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed, 'sidebar-navbar-dark': isNavbarDark }">
      <div class="brand-row">
        <div class="brand">
          <div class="brand-copy">
            <div class="brand-wordmark" aria-label="TBDex">
              <span class="brand-wordmark-primary">TB</span>
              <span class="brand-wordmark-accent">DeX</span>
            </div>
          </div>
        </div>
        <button
          type="button"
          class="sidebar-toggle-btn"
          :aria-label="sidebarToggleLabel"
          :title="sidebarToggleLabel"
          :aria-pressed="sidebarCollapsed"
          @click="toggleSidebar"
        >
          <component :is="sidebarToggleIcon" :size="16" :stroke-width="2.35" />
        </button>
      </div>

      <nav class="nav" aria-label="Primary">
        <PanelMenu :model="visibleNavItems" class="nav-panelmenu">
          <template #headericon="{ item, class: iconClass }">
            <component
              :is="item.iconComponent"
              v-if="item.iconComponent"
              :class="[iconClass, 'nav-item-icon']"
              :size="20"
              :stroke-width="isNavMenuItemActive(item) ? 2.4 : 2.2"
            />
          </template>
          <template #itemicon="{ item, class: iconClass }">
            <component
              :is="item.iconComponent"
              v-if="item.iconComponent"
              :class="[iconClass, 'nav-item-icon']"
              :size="20"
              :stroke-width="isNavMenuItemActive(item) ? 2.4 : 2.2"
            />
          </template>
        </PanelMenu>
      </nav>

      <div class="account-panel">
        <div class="account-meta">
          <Avatar :label="accountInitial" shape="circle" class="account-avatar" />
          <div>
            <p class="account-name">{{ accountName }}</p>
            <p class="account-role">{{ accountRole }}</p>
          </div>
        </div>
        <Button text class="account-menu-btn" @click="toggleAccountMenu">
          <Avatar v-if="sidebarCollapsed" :label="accountInitial" shape="circle" class="account-menu-avatar" />
          <MoreVertical v-else :size="24" :stroke-width="2.45" />
        </Button>
        <Menu ref="accountMenu" :model="accountActions" popup class="account-menu" />
      </div>
    </aside>

    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@24..48,700;24..48,800&family=DM+Sans:ital,wght@0,400;0,500;0,600;0,700;1,400;1,500&family=IBM+Plex+Sans:ital,wght@0,400;0,500;0,600;0,700;1,400;1,500&family=Plus+Jakarta+Sans:ital,wght@0,300;0,400;0,500;0,600;0,700;0,800;1,400;1,500;1,600&display=swap');

:global(:root) {
  --app-font-family: 'IBM Plex Sans', 'DM Sans', 'Plus Jakarta Sans', 'SF Pro Display', 'Segoe UI', 'Helvetica Neue', sans-serif;
  --app-font-size-base: 14.5px;
  --app-line-height-base: 1.55;
  --app-control-height: 2.25rem;
  --app-control-radius: 8px;
  --app-control-font-size: 0.855rem;
  --page-section-margin-bottom: 1rem;
  --page-content-padding: 2rem;
  --page-content-padding-mobile: 1.25rem;
  --page-content-top-offset: 0.5rem;
  --sidebar-content-inline-offset: 0.75rem;

  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 14px;
  --radius-xl: 18px;
  --radius-full: 999px;

  --shadow-xs: 0 1px 2px 0 rgb(0 0 0 / 0.04);
  --shadow-sm: 0 1px 3px 0 rgb(0 0 0 / 0.06), 0 1px 2px -1px rgb(0 0 0 / 0.04);
  --shadow-md: 0 4px 8px -2px rgb(0 0 0 / 0.06), 0 2px 4px -2px rgb(0 0 0 / 0.04);
  --shadow-lg: 0 12px 24px -4px rgb(0 0 0 / 0.08), 0 4px 8px -4px rgb(0 0 0 / 0.04);
  --brand-accent: var(--p-primary-color);
  --brand-accent-strong: var(--p-primary-500);
  --brand-wordmark-primary: var(--p-surface-900);

  --sidebar-bg: color-mix(in srgb, var(--p-surface-100) 64%, transparent);
  --sidebar-text: var(--p-surface-500);
  --sidebar-text-hover: var(--p-surface-900);
  --sidebar-text-active: var(--p-surface-900);
  --sidebar-hover-bg: color-mix(in srgb, var(--p-surface-300) 46%, transparent);
  --sidebar-active-bg: color-mix(in srgb, var(--p-surface-300) 64%, var(--p-surface-100));
  --sidebar-active-icon: var(--p-primary-600);
  --sidebar-border: color-mix(in srgb, var(--p-surface-500) 32%, transparent);
  --sidebar-glass-edge: color-mix(in srgb, var(--p-surface-0) 72%, transparent);
  --sidebar-depth-shadow: color-mix(in srgb, var(--p-surface-950) 10%, transparent);

  --app-body-bg: var(--p-surface-50);
  --app-body-text: var(--p-surface-900);
  --content-bg: var(--p-surface-50);

  --account-panel-bg: color-mix(in srgb, var(--p-surface-0) 54%, transparent);
  --account-panel-border: color-mix(in srgb, var(--p-surface-500) 36%, transparent);
}

:global(body) {
  margin: 0;
  font-family: var(--app-font-family);
  font-size: var(--app-font-size-base);
  line-height: var(--app-line-height-base);
  background: var(--app-body-bg);
  color: var(--app-body-text);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

:global(#app),
:global(.p-component),
:global(button),
:global(input),
:global(select),
:global(textarea) {
  font-family: var(--app-font-family);
}

:global(.p-button) {
  border-radius: var(--app-control-radius);
  box-sizing: border-box;
  transform: none;
  font-weight: 600;
  letter-spacing: -0.005em;
  transition: all 0.15s ease;
}

:global(.p-button:hover),
:global(.p-button:active),
:global(.p-button:focus-visible) {
  transform: none;
}

:global(.p-inputtext) {
  border-radius: var(--app-control-radius);
  font-size: var(--app-control-font-size);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

:global(.p-inputtext:focus) {
  border-color: var(--p-primary-400);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--p-primary-color, #64748b) 24%, transparent);
}

:global(.p-dialog) {
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg), 0 0 0 1px rgb(0 0 0 / 0.03);
  border: none;
}

:global(.p-dialog .p-dialog-header) {
  border-bottom: 1px solid #f0f1f3;
  padding: 1.25rem 1.5rem;
}

:global(.p-dialog .p-dialog-header .p-dialog-title) {
  font-weight: 700;
  font-size: 1.05rem;
  letter-spacing: -0.01em;
}

:global(.p-dialog .p-dialog-content) {
  padding: 1.5rem;
}

:global(.p-dialog .p-dialog-footer) {
  border-top: 1px solid #f0f1f3;
  padding: 1rem 1.5rem;
}

:global(.p-select) {
  border-radius: var(--app-control-radius);
}

:global(.p-confirmdialog) {
  border-radius: var(--radius-xl);
}

:global(.p-toast-message) {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
}

:global(.p-menu.p-menu-overlay) {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg), 0 0 0 1px rgb(0 0 0 / 0.04);
  border: none;
  padding: 6px;
}

:global(.p-menu .p-menu-item-link) {
  border-radius: var(--radius-sm);
  padding: 0.55rem 0.75rem;
  font-size: 0.855rem;
  font-weight: 500;
  gap: 0.5rem;
}

:global(.p-menu .p-menu-item-link:hover) {
  background: #f4f5f7;
}

:global(.p-menu.account-menu .p-menu-submenu-label) {
  font-size: 0.67rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--p-surface-500);
  padding: 0.55rem 0.75rem 0.35rem;
}

:global(.p-menu.account-menu .account-menu-navbar-toggle-active .p-menu-item-link),
:global(.p-menu.account-menu .account-menu-theme-item-active .p-menu-item-link) {
  color: var(--p-primary-700);
  background: color-mix(in srgb, var(--p-primary-color) 14%, transparent);
  font-weight: 600;
}

:global(#app) {
  min-height: 100vh;
}

:global(.p-datatable .p-datatable-thead > tr > th),
:global(.p-datatable .p-datatable-header-cell) {
  background: #f8f9fb;
  color: #6b7280;
  border-color: #f0f1f3;
  border-bottom: 1px solid #e8eaed;
  font-weight: 600;
  font-size: 0.78rem;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  padding: 0.38rem 1rem;
}

:global(.p-datatable .p-datatable-tbody > tr > td) {
  background: #ffffff;
  color: #1a1d23;
  border-bottom: 1px solid #eef0f3;
  border-color: #eef0f3;
  padding: 0.38rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
}

:global(.p-datatable .p-datatable-tbody > tr:not(.p-datatable-row-selected):hover > td) {
  background: #f8f9fb;
}

:global(.p-datatable .p-datatable-tbody > tr.p-datatable-row-selected > td) {
  background: rgb(100 116 139 / 0.06);
}

:global(.p-datatable .p-datatable-thead .p-checkbox),
:global(.p-datatable .p-datatable-tbody .p-checkbox) {
  width: 1rem;
  height: 1rem;
}

:global(.p-datatable .p-datatable-thead .p-checkbox .p-checkbox-box),
:global(.p-datatable .p-datatable-tbody .p-checkbox .p-checkbox-box) {
  width: 1rem;
  height: 1rem;
  border-radius: 4px;
}

:global(.p-datatable .p-datatable-thead .p-checkbox .p-checkbox-icon),
:global(.p-datatable .p-datatable-tbody .p-checkbox .p-checkbox-icon) {
  width: 0.75rem;
  height: 0.75rem;
  font-size: 0.75rem;
}

:global(.p-datatable .p-datatable-paginator-bottom) {
  border-top: 1px solid #f0f1f3;
  border-bottom-width: 0;
  background: transparent;
  padding: 0.625rem 1rem;
  min-height: 3.25rem;
  box-sizing: border-box;
}

:global(.p-datatable .p-paginator) {
  border: none;
  background: transparent;
  padding: 0;
  gap: 0;
  justify-content: center;
  align-items: center;
  flex-wrap: nowrap;
}

:global(.p-datatable .p-paginator .p-paginator-content) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  gap: 0.24rem;
}

:global(.p-datatable .p-paginator .p-paginator-pages) {
  display: inline-flex;
  align-items: center;
  gap: 0.18rem;
  margin: 0 0.12rem;
}

:global(.p-datatable .p-paginator .p-paginator-first),
:global(.p-datatable .p-paginator .p-paginator-last),
:global(.p-datatable .p-paginator .p-paginator-prev),
:global(.p-datatable .p-paginator .p-paginator-next) {
  min-width: 2rem;
  height: 2rem;
  border-radius: var(--app-control-radius);
  border: none;
  background: transparent;
  color: color-mix(in srgb, var(--p-surface-700) 92%, black 8%);
  box-shadow: none;
  transition: color 0.15s ease, background-color 0.15s ease;
}

:global(.p-datatable .p-paginator .p-paginator-page) {
  min-width: 2rem;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: #6b7280;
  box-shadow: none;
  transition: background-color 0.15s ease, color 0.15s ease;
  font-size: 0.82rem;
  font-weight: 500;
}

:global(.p-datatable .p-paginator .p-paginator-first:hover),
:global(.p-datatable .p-paginator .p-paginator-last:hover),
:global(.p-datatable .p-paginator .p-paginator-prev:hover),
:global(.p-datatable .p-paginator .p-paginator-next:hover),
:global(.p-datatable .p-paginator .p-paginator-page:hover) {
  border-color: transparent;
  background: color-mix(in srgb, var(--p-surface-900) 6%, transparent);
  color: var(--p-surface-900);
}

:global(.p-datatable .p-paginator .p-paginator-first:focus-visible),
:global(.p-datatable .p-paginator .p-paginator-last:focus-visible),
:global(.p-datatable .p-paginator .p-paginator-prev:focus-visible),
:global(.p-datatable .p-paginator .p-paginator-next:focus-visible),
:global(.p-datatable .p-paginator .p-paginator-page:focus-visible) {
  outline: none;
  border-color: transparent;
  box-shadow:
    0 0 0 2px color-mix(in srgb, var(--p-primary-color) 14%, transparent);
}

:global(.p-datatable .p-paginator .p-paginator-page.p-paginator-page-selected),
:global(.p-datatable .p-paginator .p-paginator-page[aria-current='page']) {
  border: none;
  background: var(--p-surface-100, #f1f5f9);
  color: var(--p-surface-900, #0f172a);
  font-weight: 600;
  box-shadow: none;
}

:global(.p-datatable .p-paginator .p-disabled) {
  opacity: 0.28;
}

:global(.p-datatable .p-paginator .p-paginator-current) {
  font-size: 0.78rem;
  color: #9ca3af;
  white-space: nowrap;
  margin-left: 0;
  margin-right: auto;
}

:global(.p-datatable .p-paginator .p-paginator-content > :first-child) {
  margin-left: 0;
}

:global(.p-datatable .p-paginator .p-paginator-rpp-dropdown) {
  margin-left: auto;
  height: 2rem;
  min-width: 0;
  border: 1px solid var(--p-surface-200);
  border-radius: var(--app-control-radius);
  background: var(--p-surface-0);
  box-shadow: none;
  font-size: 0.82rem;
  color: var(--p-surface-700);
  transition: border-color 0.15s ease, background-color 0.15s ease;
}

:global(.p-datatable .p-paginator .p-paginator-rpp-dropdown:hover) {
  border-color: var(--p-surface-300);
  background: var(--p-surface-50);
}

:global(.p-datatable .p-paginator .p-paginator-rpp-dropdown:focus-within) {
  border-color: color-mix(in srgb, var(--p-primary-color) 40%, var(--p-surface-300));
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-primary-color) 10%, transparent);
}

:global(.p-datatable .p-paginator .p-paginator-rpp-dropdown .p-select-label) {
  padding: 0 0.55rem;
  font-size: 0.82rem;
  color: var(--p-surface-700);
  line-height: 2rem;
}

:global(.p-datatable .p-paginator .p-paginator-rpp-dropdown .p-select-dropdown) {
  width: 1.6rem;
  color: var(--p-surface-500);
}

/* ── DataTable sort indicators ────────────────────────── */

:global(.p-datatable .p-datatable-sortable-column .p-datatable-sort-icon) {
  color: var(--p-surface-400);
  transition: color 0.15s ease;
  margin-left: 0.3rem;
}

:global(.p-datatable .p-datatable-sortable-column.p-datatable-column-sorted .p-datatable-sort-icon),
:global(.p-datatable .p-datatable-thead > tr > th[data-p-sorted='true'] .p-datatable-sort-icon) {
  color: var(--brand-accent);
}

:global(.p-datatable .p-datatable-sortable-column:hover .p-datatable-sort-icon) {
  color: var(--p-surface-600);
}

:global(.p-datatable .p-datatable-sortable-column.p-datatable-column-sorted:hover .p-datatable-sort-icon),
:global(.p-datatable .p-datatable-thead > tr > th[data-p-sorted='true']:hover .p-datatable-sort-icon) {
  color: var(--brand-accent-strong);
}

/* ── DataTable column filter overlay ─────────────────── */

:global(.p-datatable-filter-overlay) {
  border: none;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg), 0 0 0 1px rgb(0 0 0 / 0.04);
  background: #ffffff;
  font-family: var(--app-font-family);
  min-width: 14rem;
}

:global(.p-datatable-filter-overlay .p-datatable-filter-constraint) {
  padding: 0.75rem;
  border-color: color-mix(in srgb, var(--p-surface-300) 35%, transparent);
}

:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar) {
  padding: 0.55rem 0.75rem;
  display: flex !important;
  width: 100%;
  justify-content: flex-end !important;
  align-items: center;
  gap: 0.4rem;
}

:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-button) {
  font-size: 0.82rem;
  padding: 0.35rem 0.8rem;
  min-width: 0;
  min-height: 2rem;
  line-height: 1;
  margin: 0 !important;
  border: none;
  box-shadow: none;
}

:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-button:hover),
:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-button:focus-visible) {
  border: none;
  box-shadow: none;
}

:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-datatable-filter-clear-button),
:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-button:first-child) {
  color: var(--p-red-600);
}

:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-datatable-filter-clear-button:hover),
:global(.p-datatable-filter-overlay .p-datatable-filter-buttonbar .p-button:first-child:hover) {
  color: var(--p-red-700);
  background: color-mix(in srgb, var(--p-red-500) 12%, transparent);
}

:global(.p-datatable-filter-overlay .p-select),
:global(.p-datatable-filter-overlay .p-dropdown),
:global(.p-datatable-filter-overlay .p-multiselect) {
  border: 1px solid color-mix(in srgb, var(--p-surface-300) 45%, transparent) !important;
  border-color: color-mix(in srgb, var(--p-surface-300) 45%, transparent) !important;
  box-shadow: none;
  background: var(--p-surface-50);
}

:global(.p-datatable-filter-overlay .p-select:hover),
:global(.p-datatable-filter-overlay .p-dropdown:hover),
:global(.p-datatable-filter-overlay .p-multiselect:hover) {
  border-color: color-mix(in srgb, var(--p-surface-300) 65%, transparent) !important;
}

:global(.p-datatable-filter-overlay .p-select.p-focus),
:global(.p-datatable-filter-overlay .p-dropdown.p-focus),
:global(.p-datatable-filter-overlay .p-multiselect.p-focus),
:global(.p-datatable-filter-overlay .p-select:focus-within),
:global(.p-datatable-filter-overlay .p-dropdown:focus-within),
:global(.p-datatable-filter-overlay .p-multiselect:focus-within) {
  border-color: color-mix(in srgb, var(--p-surface-400) 70%, transparent) !important;
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--p-surface-300) 28%, transparent);
}

:global(.p-datatable-filter-overlay .p-select.p-filled),
:global(.p-datatable-filter-overlay .p-dropdown.p-filled),
:global(.p-datatable-filter-overlay .p-multiselect.p-filled),
:global(.p-datatable-filter-overlay .p-inputwrapper-filled .p-select),
:global(.p-datatable-filter-overlay .p-inputwrapper-filled .p-dropdown),
:global(.p-datatable-filter-overlay .p-inputwrapper-filled .p-multiselect) {
  border-color: color-mix(in srgb, var(--p-surface-300) 45%, transparent) !important;
}

:global(.p-select-overlay),
:global(.p-dropdown-panel),
:global(.p-multiselect-overlay) {
  border: none;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg), 0 0 0 1px rgb(0 0 0 / 0.04);
}

:global(.p-datatable-column-filter-button) {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: var(--app-control-radius);
  color: var(--p-surface-500);
  transition: color 0.15s, background-color 0.15s;
}

:global(.p-datatable-column-filter-button:hover) {
  color: var(--p-surface-700);
  background: color-mix(in srgb, var(--p-surface-900) 6%, transparent);
}

:global(.p-datatable-column-filter-button[data-p-active="true"]) {
  color: var(--brand-accent);
  background: color-mix(in srgb, var(--brand-accent) 14%, transparent);
}

:global(.p-datatable-column-filter-button[data-p-active="true"] .p-button-icon),
:global(.p-datatable-column-filter-button[data-p-active="true"] .p-icon) {
  color: inherit;
}

:global(.p-datatable-column-filter-button[data-p-active="true"]:hover) {
  color: var(--brand-accent-strong);
  background: color-mix(in srgb, var(--brand-accent) 22%, transparent);
}

:global(.p-datatable .p-datatable-header-cell .p-datatable-popover-filter),
:global(.p-datatable .p-datatable-thead > tr > th .p-datatable-popover-filter) {
  margin-inline-start: 0.28rem;
}

/* ── DataTable rounded-corner inheritance inside Cards ── */

:global(.p-card .p-datatable) {
  border-radius: var(--radius-xl, 18px);
  overflow: hidden;
}

:global(.p-card .p-datatable .p-datatable-table-container) {
  border-radius: inherit;
}

:global(.p-card .p-datatable .p-datatable-header) {
  border-top-left-radius: inherit;
  border-top-right-radius: inherit;
}

:global(.p-card .p-datatable .p-datatable-paginator-bottom) {
  border-bottom-left-radius: inherit;
  border-bottom-right-radius: inherit;
}

/* ── DataTable header slot ───────────────────────────── */

:global(.p-datatable .p-datatable-header) {
  background: #ffffff;
  border: none;
  padding: 0.625rem;
  border-radius: inherit;
}

@media (max-width: 760px) {
  :global(.p-datatable .p-datatable-paginator-bottom) {
    padding: 0.5rem 0.3rem;
    min-height: 3.25rem;
  }

  :global(.p-datatable .p-paginator) {
    gap: 0.18rem;
    flex-wrap: wrap;
  }

  :global(.p-datatable .p-paginator .p-paginator-page) {
    min-width: 2.1rem;
    height: 2.1rem;
    font-size: 0.88rem;
  }

  :global(.p-datatable .p-paginator .p-paginator-current) {
    margin-right: 0;
  }

  :global(.p-datatable .p-paginator .p-paginator-rpp-dropdown) {
    margin-left: 0;
  }
}

.layout-shell {
  --sidebar-width-val: 260px;
  min-height: 100vh;
  height: 100vh;
  display: grid;
  grid-template-columns: var(--sidebar-width-val) minmax(0, 1fr);
  overflow: hidden;
  transition: grid-template-columns 200ms cubic-bezier(0.4, 0, 0.2, 1);
}

.layout-shell.layout-shell-collapsed {
  --sidebar-width-val: 76px;
}

.sidebar {
  --sidebar-bg: color-mix(in srgb, var(--p-surface-100) 64%, transparent);
  --sidebar-text: var(--p-surface-500);
  --sidebar-text-hover: var(--p-surface-900);
  --sidebar-text-active: var(--p-surface-900);
  --sidebar-hover-bg: color-mix(in srgb, var(--p-surface-300) 46%, transparent);
  --sidebar-active-bg: color-mix(in srgb, var(--p-surface-300) 64%, var(--p-surface-100));
  --sidebar-active-icon: var(--p-primary-600);
  --sidebar-border: color-mix(in srgb, var(--p-surface-500) 32%, transparent);
  --sidebar-glass-edge: color-mix(in srgb, var(--p-surface-0) 72%, transparent);
  --sidebar-depth-shadow: color-mix(in srgb, var(--p-surface-950) 10%, transparent);
  --account-panel-bg: color-mix(in srgb, var(--p-surface-0) 54%, transparent);
  --account-panel-border: color-mix(in srgb, var(--p-surface-500) 36%, transparent);
  --profile-avatar-bg: var(--p-surface-900);
  --profile-avatar-text: var(--p-surface-0);
  --profile-avatar-border: color-mix(in srgb, var(--p-surface-950) 30%, transparent);
  --brand-wordmark-primary: var(--p-surface-900);
  --brand-accent: var(--p-primary-400);

  background: var(--sidebar-bg);
  backdrop-filter: blur(32px) saturate(180%) brightness(0.97);
  -webkit-backdrop-filter: blur(32px) saturate(180%) brightness(0.97);
  border-right: 1px solid var(--sidebar-border);
  box-shadow: inset -1px 0 0 var(--sidebar-glass-edge), 4px 0 32px var(--sidebar-depth-shadow);
  padding: 1.5rem 1rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  transition: padding 200ms cubic-bezier(0.4, 0, 0.2, 1), gap 200ms ease;
}

.sidebar.sidebar-navbar-dark {
  --sidebar-bg: var(--p-surface-900);
  --sidebar-text: var(--p-surface-400);
  --sidebar-text-hover: var(--p-surface-0);
  --sidebar-text-active: var(--p-surface-100);
  --sidebar-hover-bg: color-mix(in srgb, var(--p-surface-700) 56%, transparent);
  --sidebar-active-bg: color-mix(in srgb, var(--p-surface-700) 68%, var(--p-surface-600));
  --sidebar-active-icon: var(--p-primary-200);
  --sidebar-border: color-mix(in srgb, var(--p-surface-400) 25%, transparent);
  --sidebar-glass-edge: color-mix(in srgb, var(--p-surface-0) 15%, transparent);
  --sidebar-depth-shadow: color-mix(in srgb, var(--p-surface-0) 8%, transparent);
  --account-panel-bg: color-mix(in srgb, var(--p-surface-0) 6%, transparent);
  --account-panel-border: color-mix(in srgb, var(--p-surface-400) 24%, transparent);
  --profile-avatar-bg: var(--p-surface-0);
  --profile-avatar-text: var(--p-surface-900);
  --profile-avatar-border: color-mix(in srgb, var(--p-surface-0) 45%, var(--p-surface-500));
  --brand-wordmark-primary: var(--p-surface-0);
  --brand-accent: var(--p-primary-color);
}

.sidebar::-webkit-scrollbar {
  width: 4px;
}

.sidebar::-webkit-scrollbar-thumb {
  background: color-mix(in srgb, var(--p-surface-500) 45%, transparent);
  border-radius: 4px;
}

.sidebar.collapsed {
  padding: 1.5rem 0.5rem 1rem;
  gap: 1rem;
}

.brand-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  padding: 0.25rem 0.25rem 0;
}

.brand {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  padding: 0.15rem 0 0.15rem var(--sidebar-content-inline-offset);
  min-width: 0;
}

.brand-copy {
  min-width: 0;
}

.brand-wordmark {
  display: flex;
  align-items: baseline;
  gap: 0.02rem;
  line-height: 1;
  white-space: nowrap;
}

.brand-wordmark-primary,
.brand-wordmark-accent {
  margin: 0;
  font-family: 'Bricolage Grotesque', 'Avenir Next', 'Segoe UI', sans-serif;
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.045em;
  text-rendering: geometricPrecision;
}

.brand-wordmark-primary {
  color: var(--brand-wordmark-primary);
}

.brand-wordmark-accent {
  color: var(--brand-accent);
}

.sidebar-toggle-btn {
  width: 2rem;
  height: 2rem;
  border: none;
  background: transparent;
  color: var(--sidebar-text);
  border-radius: var(--radius-sm);
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: color 0.15s ease, background 0.15s ease;
}

.sidebar-toggle-btn:hover {
  background: var(--sidebar-hover-bg);
  color: var(--sidebar-text-hover);
}

.sidebar-toggle-btn:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--p-surface-500) 40%, transparent);
  outline-offset: 2px;
}

.nav {
  margin-top: 0.5rem;
  flex: 1;
}

.nav :deep(.nav-panelmenu) {
  border: none;
  background: transparent;
}

.nav :deep(.nav-panelmenu .p-panelmenu-header-icon) {
  color: inherit;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item-icon) {
  color: inherit;
}

.nav :deep(.nav-panelmenu .p-panelmenu-submenu-icon) {
  display: none;
}

.nav-item-icon {
  flex: 0 0 auto;
}

.nav :deep(.nav-panelmenu .p-panelmenu-header-link),
.nav :deep(.nav-panelmenu .p-panelmenu-item-link) {
  min-height: 2.5rem;
  padding-block: 0.5rem;
  padding-inline: var(--sidebar-content-inline-offset);
  column-gap: 0.75rem;
  transition: color 0.15s ease, background 0.15s ease;
}

.nav :deep(.nav-panelmenu .p-panelmenu-header-label) {
  color: inherit;
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: -0.005em;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item-label) {
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: -0.005em;
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-link) {
  color: var(--sidebar-text);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content .p-panelmenu-item-link) {
  color: var(--sidebar-text);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-link:hover) {
  color: var(--sidebar-text-hover);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content .p-panelmenu-item-link:hover) {
  color: var(--sidebar-text-hover);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-link:hover .p-panelmenu-header-icon),
.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content .p-panelmenu-item-link:hover .p-panelmenu-item-icon) {
  color: var(--sidebar-active-icon);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content) {
  background: transparent;
  border-radius: var(--radius-md);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content) {
  background: transparent;
  border-radius: var(--radius-md);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content:hover) {
  background: var(--sidebar-hover-bg);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content:hover) {
  background: var(--sidebar-hover-bg);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content) {
  background: var(--sidebar-active-bg);
  box-shadow: none;
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-link) {
  color: var(--sidebar-text-active);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-label) {
  color: var(--sidebar-text-active);
  font-weight: 600;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content) {
  background: var(--sidebar-active-bg);
  box-shadow: none;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-link) {
  color: var(--sidebar-text-active);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-label) {
  color: var(--sidebar-text-active);
  font-weight: 600;
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-icon),
.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-icon) {
  color: var(--sidebar-active-icon);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content:hover) {
  background: var(--sidebar-active-bg);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content:hover) {
  background: var(--sidebar-active-bg);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-link:hover) {
  color: var(--sidebar-text-active);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-link:hover) {
  color: var(--sidebar-text-active);
}

.content {
  padding:
    calc(var(--page-content-padding) + var(--page-content-top-offset))
    var(--page-content-padding)
    var(--page-content-padding);
  min-width: 0;
  height: 100vh;
  overflow-y: auto;
  background: var(--content-bg);
}

.account-panel {
  margin-top: auto;
  background: var(--account-panel-bg);
  border: 1px solid var(--account-panel-border);
  border-radius: var(--radius-lg);
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.55rem;
}

.account-meta {
  display: flex;
  align-items: center;
  gap: 0.55rem;
  min-width: 0;
}

.account-meta > div {
  min-width: 0;
}

.account-avatar,
.account-menu-avatar {
  background: var(--profile-avatar-bg);
  color: var(--profile-avatar-text);
  border: 1px solid var(--profile-avatar-border);
  flex: 0 0 auto;
}

.account-menu-avatar {
  width: 2rem;
  height: 2rem;
  font-size: 0.78rem;
}

.account-name {
  margin: 0;
  color: var(--sidebar-text-active);
  font-size: 0.855rem;
  font-weight: 600;
  line-height: 1.1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.account-role {
  margin: 0.14rem 0 0;
  color: var(--sidebar-text);
  font-size: 0.78rem;
  line-height: 1.1;
}

.account-menu-btn {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: var(--radius-sm);
  color: var(--sidebar-text);
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: color 0.15s ease, background 0.15s ease;
}

.account-menu-btn:hover {
  background: var(--sidebar-hover-bg);
  color: var(--sidebar-text-hover);
}

.sidebar.collapsed .brand {
  display: none;
  padding-left: 0;
}

.sidebar.collapsed .brand-row {
  justify-content: center;
  padding: 0;
}

.sidebar.collapsed .brand-copy {
  display: none;
}

.sidebar.collapsed .nav {
  margin-top: 0;
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-header-link),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item-link) {
  justify-content: center;
  padding-inline: 0;
  min-height: 2.5rem;
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content) {
  width: 2.75rem;
  min-height: 2.5rem;
  margin-inline: auto;
  border-radius: var(--radius-md);
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content:hover),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content:hover) {
  width: 2.75rem;
  min-height: 2.5rem;
  margin-inline: auto;
  border-radius: var(--radius-md);
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-header-label),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item-label) {
  display: none;
}

.sidebar.collapsed .account-panel {
  justify-content: center;
  gap: 0;
  padding: 0.375rem;
}

.sidebar.collapsed .account-meta {
  display: none;
}

.sidebar.collapsed .account-menu-btn {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: var(--radius-full);
}

@media (max-width: 900px) {
  .layout-shell {
    height: auto;
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    overflow: visible;
  }

  .layout-shell.layout-shell-collapsed {
    grid-template-columns: 1fr;
  }

  .sidebar {
    border-right: none;
    border-bottom: 1px solid var(--sidebar-border);
    overflow: visible;
  }

  .sidebar-toggle-btn {
    display: none;
  }

  .sidebar.collapsed {
    padding: 1.5rem 1rem 1rem;
    gap: 1.25rem;
  }

  .sidebar.collapsed .brand {
    display: flex;
    justify-content: flex-start;
    padding-left: var(--sidebar-content-inline-offset);
  }

  .sidebar.collapsed .brand-row {
    justify-content: flex-start;
    padding: 0.25rem 0.25rem 0;
  }

  .sidebar.collapsed .brand-copy {
    display: block;
  }

  .sidebar.collapsed .nav {
    margin-top: 0.35rem;
  }

  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-header-link),
  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item-link) {
    justify-content: flex-start;
    padding-inline: var(--sidebar-content-inline-offset);
    min-height: initial;
  }

  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content),
  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content) {
    width: auto;
    min-height: initial;
    margin-inline: 0;
  }

  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content:hover),
  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content:hover) {
    width: auto;
    min-height: initial;
    margin-inline: 0;
  }

  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-header-label),
  .sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item-label) {
    display: initial;
  }

  .sidebar.collapsed .account-panel {
    flex-direction: row;
    gap: 0.55rem;
    padding: 0.5rem;
  }

  .sidebar.collapsed .account-menu-btn {
    width: 2.25rem;
    height: 2.25rem;
  }

  .account-panel {
    margin-top: 0.3rem;
  }

  .content {
    padding: var(--page-content-padding-mobile);
    height: auto;
    overflow: visible;
    background: var(--content-bg);
  }
}
</style>
