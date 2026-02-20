<script setup lang="ts">
import type { MenuItem } from 'primevue/menuitem'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PanelLeft, MoreVertical } from 'lucide-vue-next'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import PanelMenu from 'primevue/panelmenu'
import ConfirmDialog from 'primevue/confirmdialog'
import Toast from 'primevue/toast'
import { useAuthStore } from '@/stores/auth'
import { useNavItems } from '@/composables/useNavItems'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { visibleNavItems } = useNavItems()
const sidebarCollapseStorageKey = 'borderform.ui.sidebarCollapsed'
const sidebarCollapsed = ref(false)
authStore.hydrate()
onMounted(() => {
  sidebarCollapsed.value = localStorage.getItem(sidebarCollapseStorageKey) === '1'
})

const isAuthRoute = computed(() => route.path.startsWith('/auth') || route.path === '/reset-password')
const sidebarToggleLabel = computed(() =>
  sidebarCollapsed.value ? 'Expand navigation sidebar' : 'Collapse navigation sidebar',
)
const accountName = computed(() => authStore.accountName)
const accountRole = computed(() => {
  const primary = authStore.accountRole?.trim()
  if (primary && primary.toLowerCase() !== 'role') return primary
  return authStore.roles[0] || 'No role assigned'
})
const accountInitial = computed(() => accountName.value.trim().charAt(0).toUpperCase() || 'A')

const accountActions = computed<MenuItem[]>(() => [
  {
    label: 'Settings',
    command: async () => {
      await router.push('/settings')
    },
  },
  {
    label: 'Logout',
    command: async () => {
      authStore.logout()
      await router.push('/auth/login')
    },
  },
])
const accountMenu = ref()

function toggleAccountMenu(event: Event) {
  accountMenu.value?.toggle(event)
}

function isNavMenuItemActive(item: MenuItem) {
  return typeof item.class === 'string' && item.class.includes('nav-menu-item-active')
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem(sidebarCollapseStorageKey, sidebarCollapsed.value ? '1' : '0')
}
</script>

<template>
  <Toast position="top-right" />
  <ConfirmDialog />

  <RouterView v-if="isAuthRoute" />

  <div v-else class="layout-shell" :class="{ 'layout-shell-collapsed': sidebarCollapsed }">
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="brand-row">
        <div class="brand">
          <div class="brand-copy">
            <h1>TBDeX</h1>
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
          <PanelLeft :size="16" :stroke-width="2.35" />
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
        <Menu ref="accountMenu" :model="accountActions" popup />
      </div>
    </aside>

    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,300;0,400;0,500;0,600;0,700;1,400;1,500;1,600&family=Space+Grotesk:wght@600;700&display=swap');

:global(:root) {
  --app-font-family: 'IBM Plex Sans', 'Avenir Next', 'SF Pro Text', 'Segoe UI', 'Helvetica Neue', sans-serif;
  --app-font-size-base: 15px;
  --app-line-height-base: 1.5;
  --app-control-height: 2.4rem;
  --app-control-radius: 0.46rem;
  --app-control-font-size: 0.9rem;
  --page-section-margin-bottom: 1.15rem;
  --page-content-padding: 1.5rem;
  --page-content-padding-mobile: 1rem;
}

:global(body) {
  margin: 0;
  font-family: var(--app-font-family);
  font-size: var(--app-font-size-base);
  line-height: var(--app-line-height-base);
  background:
    radial-gradient(circle at top left, var(--p-surface-50) 0%, transparent 45%),
    radial-gradient(circle at bottom right, var(--p-surface-100) 0%, transparent 40%),
  var(--p-surface-100);
  color: var(--p-surface-900);
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
}

:global(#app) {
  min-height: 100vh;
}

:global(.p-datatable .p-datatable-thead > tr > th),
:global(.p-datatable .p-datatable-header-cell) {
  background: var(--p-surface-100);
  color: var(--p-surface-900);
  border-color: var(--p-surface-200);
  font-weight: 700;
}

:global(.p-datatable .p-datatable-tbody > tr > td) {
  background: var(--p-surface-0);
  color: var(--p-surface-900);
  border-color: var(--p-surface-200);
}

:global(.p-datatable .p-datatable-tbody > tr:not(.p-datatable-row-selected):hover > td) {
  background: var(--p-surface-100);
}

:global(.p-datatable .p-datatable-paginator-bottom) {
  border-top: 1px solid color-mix(in srgb, var(--p-surface-300) 56%, transparent);
  border-bottom-width: 0;
  background: transparent;
  padding: 0.5rem 0.75rem;
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
  min-width: 2.42rem;
  height: 2.42rem;
  border-radius: var(--app-control-radius);
  border: none;
  background: transparent;
  color: var(--p-surface-700);
  box-shadow: none;
  transition: border-color 0.15s ease, background-color 0.15s ease, color 0.15s ease;
  font-size: 0.95rem;
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
  background: color-mix(in srgb, var(--p-primary-color) 12%, var(--p-surface-0));
  color: var(--p-primary-700);
  font-weight: 600;
  box-shadow: none;
}

:global(.p-datatable .p-paginator .p-disabled) {
  opacity: 0.28;
}

:global(.p-datatable .p-paginator .p-paginator-content > :first-child) {
  margin-left: auto;
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

  :global(.p-datatable .p-paginator .p-paginator-rpp-dropdown) {
    margin-left: 0;
  }
}

.layout-shell {
  --sidebar-width: 240px;
  min-height: 100vh;
  height: 100vh;
  display: grid;
  grid-template-columns: var(--sidebar-width) minmax(0, 1fr);
  overflow: hidden;
  transition: grid-template-columns 180ms ease;
}

.layout-shell.layout-shell-collapsed {
  --sidebar-width: 88px;
}

.sidebar {
  background: linear-gradient(180deg, var(--p-surface-50) 0%, var(--p-surface-100) 100%);
  border-right: 1px solid var(--p-surface-200);
  padding: 2rem 1.25rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  min-height: 0;
  overflow-y: auto;
  transition: padding 180ms ease, gap 180ms ease;
}

.sidebar.collapsed {
  padding: 1.6rem 0.65rem 1rem;
  gap: 1.1rem;
}

.brand-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.45rem;
}

.brand {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  padding: 0.15rem 0;
}

.brand-copy h1 {
  margin: 0;
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1rem;
  font-weight: 700;
  color: var(--p-primary-700);
  letter-spacing: -0.03em;
}

.sidebar-toggle-btn {
  width: 2rem;
  height: 2rem;
  border: none;
  background: transparent;
  color: var(--p-surface-700);
  border-radius: 0.55rem;
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.sidebar-toggle-btn:hover {
  background: color-mix(in srgb, var(--p-primary-color) 10%, var(--p-surface-0));
}

.sidebar-toggle-btn:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--p-primary-color) 30%, transparent);
  outline-offset: 2px;
}

.nav {
  margin-top: 0.35rem;
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
  min-height: 2.7rem;
  padding-block: 0.58rem;
  column-gap: 0.65rem;
}

.nav :deep(.nav-panelmenu .p-panelmenu-header-label) {
  color: inherit;
  font-size: 1rem;
  font-weight: 400;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item-label) {
  font-size: 1rem;
  font-weight: 400;
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-link:hover) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content .p-panelmenu-item-link:hover) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content:hover) {
  background: var(--p-surface-200);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content:hover) {
  background: var(--p-surface-200);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content) {
  background: var(--p-surface-200);
  box-shadow: inset 0 0 0 1px var(--p-surface-300);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-link) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-label) {
  color: var(--p-primary-600);
  font-weight: 500;
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content) {
  background: var(--p-surface-200);
  box-shadow: inset 0 0 0 1px var(--p-surface-300);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-link) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-label) {
  color: var(--p-primary-600);
  font-weight: 500;
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-icon),
.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-icon) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content:hover) {
  background: var(--p-surface-200);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content:hover) {
  background: var(--p-surface-200);
}

.nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-link:hover) {
  color: var(--p-primary-600);
}

.nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content .p-panelmenu-item-link:hover) {
  color: var(--p-primary-600);
}

.content {
  padding: var(--page-content-padding);
  min-width: 0;
  height: 100vh;
  overflow-y: auto;
}


.account-panel {
  margin-top: auto;
  background: var(--p-surface-0);
  border: 1px solid var(--p-surface-200);
  border-radius: 12px;
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
  background: var(--p-surface-100);
  color: var(--p-surface-900);
  border: 1px solid var(--p-surface-200);
  flex: 0 0 auto;
}

.account-menu-avatar {
  width: 2rem;
  height: 2rem;
  font-size: 0.78rem;
}

.account-name {
  margin: 0;
  color: var(--p-surface-900);
  font-size: 0.88rem;
  font-weight: 500;
  line-height: 1.1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.account-role {
  margin: 0.14rem 0 0;
  color: var(--p-surface-600);
  font-size: 0.8rem;
  line-height: 1.1;
}

.account-menu-btn {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--app-control-radius);
  color: var(--p-surface-800);
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.account-menu-btn:hover {
  background: color-mix(in srgb, var(--p-primary-color) 10%, var(--p-surface-0));
  border-color: color-mix(in srgb, var(--p-primary-color) 28%, var(--p-surface-0));
}

.sidebar.collapsed .brand {
  display: none;
  padding-left: 0;
}

.sidebar.collapsed .brand-row {
  justify-content: center;
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
  padding-inline: 0.58rem;
  min-height: 2.7rem;
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel.nav-menu-item-active > .p-panelmenu-header .p-panelmenu-header-content),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item.nav-menu-item-active > .p-panelmenu-item-content) {
  width: 2.9rem;
  min-height: 2.7rem;
  margin-inline: auto;
  border-radius: 0.62rem;
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-panel > .p-panelmenu-header .p-panelmenu-header-content:hover),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item > .p-panelmenu-item-content:hover) {
  width: 2.9rem;
  min-height: 2.7rem;
  margin-inline: auto;
  border-radius: 0.62rem;
}

.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-header-label),
.sidebar.collapsed .nav :deep(.nav-panelmenu .p-panelmenu-item-label) {
  display: none;
}

.sidebar.collapsed .account-panel {
  justify-content: center;
  gap: 0;
  padding: 0.42rem;
}

.sidebar.collapsed .account-meta {
  display: none;
}

.sidebar.collapsed .account-menu-btn {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 999px;
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
    border-bottom: 1px solid var(--p-surface-200);
    overflow: visible;
  }

  .sidebar-toggle-btn {
    display: none;
  }

  .sidebar.collapsed {
    padding: 2rem 1.25rem 1.25rem;
    gap: 1.5rem;
  }

  .sidebar.collapsed .brand {
    display: flex;
    justify-content: flex-start;
    padding-left: 0;
  }

  .sidebar.collapsed .brand-row {
    justify-content: flex-start;
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
    padding-inline: initial;
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
    width: 2.3rem;
    height: 2.3rem;
  }

  .account-panel {
    margin-top: 0.3rem;
  }

  .content {
    padding: var(--page-content-padding-mobile);
    height: auto;
    overflow: visible;
  }
}
</style>
