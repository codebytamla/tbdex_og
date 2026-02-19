import { createRouter, createWebHistory } from 'vue-router'
import { h } from 'vue'
import { RouterView } from 'vue-router'
import OverviewView from '@/views/OverviewView.vue'
import RolesView from '@/views/RolesView.vue'
import SitesView from '@/views/SitesView.vue'
import UsersView from '@/views/UsersView.vue'
import PatientsView from '@/views/PatientsView.vue'
import EncountersView from '@/views/EncountersView.vue'
import LoginView from '@/views/LoginView.vue'
import ResetPasswordView from '@/views/ResetPasswordView.vue'
import SettingsView from '@/views/SettingsView.vue'
import { useAuthStore } from '@/stores/auth'

const AdminRouteView = {
  render: () => h(RouterView),
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

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
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
      path: '/dashboard',
      name: 'dashboard',
      component: OverviewView,
      meta: { requiresAuth: true },
    },
    {
      path: '/overview',
      redirect: '/dashboard',
    },
    {
      path: '/admin',
      component: AdminRouteView,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/users',
        },
        {
          path: 'roles',
          name: 'admin-roles',
          component: RolesView,
          meta: {
            requiredPermissions: roleAccessPermissions,
          },
        },
        {
          path: 'users',
          name: 'admin-users',
          component: UsersView,
          meta: {
            requiredPermissions: userAccessPermissions,
            requiredRoles: ['administrator', 'admin'],
          },
        },
        {
          path: 'patients',
          name: 'admin-patients',
          component: PatientsView,
          meta: {
            requiredPermissions: patientAccessPermissions,
          },
        },
        {
          path: 'sites',
          name: 'admin-sites',
          component: SitesView,
          meta: {
            requiredPermissions: siteAccessPermissions,
          },
        },
      ],
    },
    {
      path: '/roles',
      redirect: '/admin/roles',
    },
    {
      path: '/patients',
      name: 'patients',
      component: PatientsView,
      meta: {
        requiresAuth: true,
        requiredPermissions: patientAccessPermissions,
      },
    },
    {
      path: '/encounters',
      name: 'encounters',
      component: EncountersView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/users',
      redirect: '/admin/users',
    },
    {
      path: '/sites',
      redirect: '/admin/sites',
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/dashboard',
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
    return { path: authStore.defaultDashboardRoute() }
  }

  return true
})

export default router
