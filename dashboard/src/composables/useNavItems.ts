import type { MenuItem } from 'primevue/menuitem'
import type { Component } from 'vue'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Frown, LayoutDashboard, MapPin, ShieldCheck, Users } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

type NavItemDefinition = {
  key: string
  label: string
  to?: string
  iconComponent?: Component
  items?: NavItemDefinition[]
  requiredPermissions?: string[]
  requiredRoles?: string[]
}

type BuiltNavItem = {
  item: MenuItem
  active: boolean
}

const usersReadPermissions = ['users.read']
const patientsReadPermissions = ['patients.read']
const sitesReadPermissions = ['sites.read']
const rolesReadPermissions = ['roles.read']

const navDefinitions: NavItemDefinition[] = [
  {
    key: 'overview',
    label: 'Overview',
    to: '/overview',
    iconComponent: LayoutDashboard,
  },
  {
    key: 'users',
    label: 'Users',
    to: '/users',
    iconComponent: Users,
    requiredPermissions: usersReadPermissions,
  },
  {
    key: 'patients',
    label: 'Patients',
    to: '/patients',
    iconComponent: Frown,
    requiredPermissions: patientsReadPermissions,
  },
  {
    key: 'sites',
    label: 'Sites',
    to: '/sites',
    iconComponent: MapPin,
    requiredPermissions: sitesReadPermissions,
  },
  {
    key: 'roles',
    label: 'Roles',
    to: '/roles',
    iconComponent: ShieldCheck,
    requiredPermissions: rolesReadPermissions,
  },
]

function isRouteMatch(currentPath: string, targetPath: string) {
  return currentPath === targetPath || currentPath.startsWith(`${targetPath}/`)
}

export function useNavItems() {
  const route = useRoute()
  const router = useRouter()
  const authStore = useAuthStore()

  const visibleNavItems = computed<MenuItem[]>(() => {
    const buildMenuItems = (definitions: NavItemDefinition[]): BuiltNavItem[] =>
      definitions.flatMap((definition) => {
        if (!authStore.canAccess(definition.requiredPermissions, definition.requiredRoles)) return []

        const childItems = definition.items ? buildMenuItems(definition.items) : []
        const hasChildren = childItems.length > 0
        const routeActive = definition.to ? isRouteMatch(route.path, definition.to) : false
        const childActive = childItems.some((child) => child.active)
        const active = routeActive || childActive

        const item: MenuItem = {
          key: definition.key,
          label: definition.label,
          iconComponent: definition.iconComponent,
          class: active ? 'nav-menu-item-active' : undefined,
          items: hasChildren ? childItems.map((child) => child.item) : undefined,
        }

        if (definition.to && !hasChildren) {
          const toPath = definition.to
          item.command = () => {
            if (route.path !== toPath) {
              void router.push(toPath)
            }
          }
        }

        return [{ item, active }]
      })

    return buildMenuItems(navDefinitions).map(({ item }) => item)
  })

  return { visibleNavItems }
}
