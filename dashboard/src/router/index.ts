import { createRouter, createWebHistory } from 'vue-router'
import OverviewView from '@/views/OverviewView.vue'
import ReferralsView from '@/views/ReferralsView.vue'
import ScreeningsView from '@/views/ScreeningsView.vue'
import RolesView from '@/views/RolesView.vue'
import SitesView from '@/views/SitesView.vue'
import UsersView from '@/views/UsersView.vue'
import PatientsView from '@/views/PatientsView.vue'
import LoginView from '@/views/LoginView.vue'
import ResetPasswordView from '@/views/ResetPasswordView.vue'
import SettingsView from '@/views/SettingsView.vue'
import { useAuthStore } from '@/stores/auth'

const referralsReadPermissions = ['referrals.read']
const screeningsReadPermissions = ['screenings.read']
const patientsReadPermissions = ['patients.read']
const sitesReadPermissions = ['sites.read']
const usersReadPermissions = ['users.read']
const rolesReadPermissions = ['roles.read']

function parseClaimMeta(value: unknown) {
  if (Array.isArray(value)) {
    return value.filter((entry): entry is string => typeof entry === 'string' && entry.trim().length > 0)
  }
  if (typeof value === 'string' && value.trim().length > 0) return [value]
  return []
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/overview',
    },
    {
      path: '/auth/login',
      name: 'login',
      component: LoginView,
      meta: { guestOnly: true },
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: ResetPasswordView,
    },
    {
      path: '/overview',
      name: 'overview',
      component: OverviewView,
      meta: { requiresAuth: true },
    },
    {
      path: '/dashboard',
      redirect: '/overview',
    },
    {
      path: '/referrals',
      name: 'referrals',
      component: ReferralsView,
      meta: {
        requiresAuth: true,
        requiredPermissions: referralsReadPermissions,
      },
    },
    {
      path: '/screenings',
      name: 'screenings',
      component: ScreeningsView,
      meta: {
        requiresAuth: true,
        requiredPermissions: screeningsReadPermissions,
      },
    },
    {
      path: '/patients',
      name: 'patients',
      component: PatientsView,
      meta: {
        requiresAuth: true,
        requiredPermissions: patientsReadPermissions,
      },
    },
    {
      path: '/sites',
      name: 'sites',
      component: SitesView,
      meta: {
        requiresAuth: true,
        requiredPermissions: sitesReadPermissions,
      },
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
      meta: {
        requiresAuth: true,
        requiredPermissions: usersReadPermissions,
      },
    },
    {
      path: '/roles',
      name: 'roles',
      component: RolesView,
      meta: {
        requiresAuth: true,
        requiredPermissions: rolesReadPermissions,
      },
    },
    {
      path: '/admin',
      redirect: '/overview',
    },
    {
      path: '/admin/referrals',
      redirect: '/referrals',
    },
    {
      path: '/admin/screenings',
      redirect: '/screenings',
    },
    {
      path: '/admin/patients',
      redirect: '/patients',
    },
    {
      path: '/admin/sites',
      redirect: '/sites',
    },
    {
      path: '/admin/users',
      redirect: '/users',
    },
    {
      path: '/admin/roles',
      redirect: '/roles',
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/overview',
    },
  ],
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()
  authStore.hydrate()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { path: '/auth/login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAuth) {
    const loaded = await authStore.ensureProfile()
    if (!loaded) {
      authStore.logout()
      return { path: '/auth/login', query: { redirect: to.fullPath } }
    }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { path: authStore.defaultOverviewRoute() }
  }

  const requiredPermissions = parseClaimMeta(to.meta.requiredPermissions)
  const requiredRoles = parseClaimMeta(to.meta.requiredRoles)
  if (
    (requiredPermissions.length > 0 || requiredRoles.length > 0) &&
    !authStore.canAccess(requiredPermissions, requiredRoles)
  ) {
    return { path: authStore.defaultOverviewRoute() }
  }

  return true
})

export default router
