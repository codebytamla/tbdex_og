<script setup lang="ts">
import type { Component } from 'vue'
import type { MenuItem } from 'primevue/menuitem'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ClipboardList,
  Frown,
  HeartPulse,
  LayoutDashboard,
  MapPin,
  MoreVertical,
  ShieldCheck,
  Users,
} from 'lucide-vue-next'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import ConfirmDialog from 'primevue/confirmdialog'
import Toast from 'primevue/toast'
import { useAuthStore } from '@/stores/auth'

type NavItem = {
  label: string
  to: string
  icon: Component
  section: 'main' | 'admin'
  visibleTo: Array<'user' | 'admin'>
  requiredPermissions?: string[]
  requiredRoles?: string[]
}

const roleAccessPermissions = [
  'roles.*',
  'roles.read',
  'roles.write',
  'roles.update',
  'roles.delete',
  'roles.manage',
  'users.*',
  'users.read',
  'users.write',
  'users.update',
  'users.delete',
  'users.manage',
]

const userAccessPermissions = [
  'users.*',
  'users.read',
  'users.write',
  'users.update',
  'users.delete',
  'users.manage',
]

const siteAccessPermissions = [
  'sites.*',
  'sites.read',
  'sites.write',
  'sites.update',
  'sites.delete',
  'sites.manage',
  'users.*',
  'users.read',
  'users.write',
  'users.update',
  'users.delete',
  'users.manage',
]

const patientAccessPermissions = [
  'patients.*',
  'patients.manage',
  'patients.read',
  'patients.create',
  'patients.update',
  'patients.add_site',
  'patients.disable',
  'patients.write',
]

const navItems: NavItem[] = [
  {
    label: 'Dashboard',
    to: '/dashboard',
    icon: LayoutDashboard,
    section: 'main',
    visibleTo: ['user'],
  },
  {
    label: 'Users',
    to: '/admin/users',
    icon: Users,
    section: 'admin',
    visibleTo: ['admin'],
    requiredPermissions: userAccessPermissions,
    requiredRoles: ['administrator', 'admin'],
  },
  {
    label: 'Patients',
    to: '/patients',
    icon: Frown,
    section: 'main',
    visibleTo: ['user', 'admin'],
    requiredPermissions: patientAccessPermissions,
  },
  {
    label: 'Encounter',
    to: '/encounters',
    icon: ClipboardList,
    section: 'main',
    visibleTo: ['user', 'admin'],
  },
  {
    label: 'Sites',
    to: '/admin/sites',
    icon: MapPin,
    section: 'admin',
    visibleTo: ['admin'],
    requiredPermissions: siteAccessPermissions,
  },
  {
    label: 'Roles',
    to: '/admin/roles',
    icon: ShieldCheck,
    section: 'admin',
    visibleTo: ['admin'],
    requiredPermissions: roleAccessPermissions,
  },
]

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
authStore.hydrate()
onMounted(() => {
  void authStore.ensureProfile(true)
})

const isAuthRoute = computed(() => route.path.startsWith('/auth') || route.path === '/reset-password')
const visibleNavItems = computed(() => navItems)
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
</script>

<template>
  <Toast position="top-right" />
  <ConfirmDialog />

  <RouterView v-if="isAuthRoute" />

  <div v-else class="layout-shell">
    <aside class="sidebar">
      <div class="brand">
        <Avatar shape="circle" class="brand-avatar">
          <HeartPulse :size="16" :stroke-width="2.25" />
        </Avatar>
        <div class="brand-copy">
          <h1>BorderDesk</h1>
        </div>
      </div>

      <nav class="nav">
        <RouterLink v-for="item in visibleNavItems" :key="item.to" :to="item.to" custom>
          <template #default="{ navigate, isActive }">
            <Button fluid text class="nav-button" :class="{ 'is-active': isActive }" @click="navigate">
              <span class="nav-item-inner">
                <component
                  :is="item.icon"
                  class="nav-item-icon"
                  :size="16"
                  :stroke-width="isActive ? 2.35 : 2.1"
                />
                <span class="nav-item-label">{{ item.label }}</span>
              </span>
            </Button>
          </template>
        </RouterLink>
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
          <MoreVertical :size="20" :stroke-width="2.35" />
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
@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,300;0,400;0,500;0,600;0,700;1,400;1,500;1,600&display=swap');

:global(:root) {
  --app-font-family: 'IBM Plex Sans', 'Avenir Next', 'SF Pro Text', 'Segoe UI', 'Helvetica Neue', sans-serif;
  --app-font-size-base: 14px;
  --app-line-height-base: 1.45;
  --app-control-height: 2.2rem;
  --app-control-radius: 0.62rem;
  --app-control-font-size: 0.83rem;
  --page-section-margin-bottom: 0.95rem;
  --page-content-padding: 1.4rem;
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

:global(#app) {
  min-height: 100vh;
}

:global(.p-datatable .p-datatable-thead > tr > th),
:global(.p-datatable .p-datatable-header-cell) {
  background: var(--p-surface-50);
  color: var(--p-surface-700);
  border-color: var(--p-surface-200);
}

:global(.p-datatable .p-datatable-tbody > tr > td) {
  background: var(--p-surface-0);
  color: var(--p-surface-900);
  border-color: var(--p-surface-200);
}

:global(.p-datatable .p-datatable-tbody > tr:not(.p-datatable-row-selected):hover > td) {
  background: var(--p-surface-100);
}

.layout-shell {
  min-height: 100vh;
  height: 100vh;
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  overflow: hidden;
}

.sidebar {
  background: linear-gradient(180deg, var(--p-surface-50) 0%, var(--p-surface-100) 100%);
  border-right: 1px solid var(--p-surface-200);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  min-height: 0;
  overflow-y: auto;
}

.brand {
  display: flex;
  gap: 0.6rem;
  align-items: center;
  padding: 0.15rem 0;
}

.brand-copy h1 {
  margin: 0;
  font-size: 0.92rem;
  font-weight: 650;
  color: var(--p-surface-900);
  letter-spacing: -0.01em;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 0.35rem;
}

.nav-button {
  justify-content: flex-start;
  border-radius: 0.7rem;
  color: var(--p-surface-700);
  border: 1px solid transparent;
}

.nav-item-inner {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-item-icon {
  flex: 0 0 auto;
}

.nav-button.is-active {
  background: var(--p-surface-200);
  color: var(--p-surface-900);
  border-color: var(--p-surface-200);
}

.nav-button.is-active .nav-item-label {
  font-weight: 600;
}

.nav :deep(.p-button.nav-button:not(:disabled):hover) {
  background: color-mix(in srgb, var(--p-primary-color) 14%, var(--p-surface-0));
  border-color: color-mix(in srgb, var(--p-primary-color) 36%, var(--p-surface-0));
  color: var(--p-primary-800);
}

.nav :deep(.p-button.nav-button.is-active:not(:disabled):hover) {
  background: color-mix(in srgb, var(--p-primary-color) 22%, var(--p-surface-0));
  border-color: color-mix(in srgb, var(--p-primary-color) 44%, var(--p-surface-0));
  color: var(--p-primary-900);
}

.content {
  padding: var(--page-content-padding);
  min-width: 0;
  height: 100vh;
  overflow-y: auto;
}

.brand-avatar {
  background: var(--p-surface-900);
  color: var(--p-surface-0);
  border: 1px solid var(--p-surface-950);
}

.brand-avatar :deep(svg) {
  display: block;
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

.account-avatar {
  background: var(--p-surface-100);
  color: var(--p-surface-900);
  border: 1px solid var(--p-surface-200);
  flex: 0 0 auto;
}

.account-name {
  margin: 0;
  color: var(--p-surface-900);
  font-size: 0.82rem;
  font-weight: 600;
  line-height: 1.1;
}

.account-role {
  margin: 0.14rem 0 0;
  color: var(--p-surface-600);
  font-size: 0.73rem;
  line-height: 1.1;
}

.account-menu-btn {
  width: 2.3rem;
  height: 2.3rem;
  border-radius: 0.6rem;
  color: var(--p-surface-700);
  flex: 0 0 auto;
}

.account-menu-btn:hover {
  background: var(--p-surface-100);
  border-color: var(--p-surface-200);
}

@media (max-width: 900px) {
  .layout-shell {
    height: auto;
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    overflow: visible;
  }

  .sidebar {
    border-right: none;
    border-bottom: 1px solid var(--p-surface-200);
    overflow: visible;
  }

  .nav {
    flex-direction: row;
    overflow-x: auto;
  }

  .nav :deep(.p-button) {
    white-space: nowrap;
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
