# Frontend Stores & Components

> Source: `dashboard/src/stores/`, `dashboard/src/components/ui/`

---

## Pinia Stores

### `auth.ts` — Auth Store
- Manages authentication state: user, access token, refresh token.
- Key actions:
  - `hydrate()` — restores auth state from localStorage on page load.
  - `ensureProfile()` — fetches current user profile from `GET /auth/me`.
  - `login(email, password)` — calls `POST /auth/login`, stores tokens.
  - `logout()` — calls `POST /auth/logout`, clears local state.
  - `canAccess(permissions, roles)` — checks if the current user has the required permissions/roles.
  - `defaultOverviewRoute()` — returns the best overview route for the current user based on permissions.
- Key getters: `isAuthenticated`, `currentUser`, `userPermissions`, `userRoles`.

### `catalog.ts` — Catalog Store
- Caches shared reference data such as sites and roles.
- Used to avoid re-fetching static data on every view.
- Populated on first access, shared across views.

---

## UI Components (`src/components/ui/`)

| Component | Description |
|---|---|
| `CreateActionButton.vue` | Primary "create" button used in page headers |
| `DataTableToolbar.vue` | Toolbar above DataTables (search, filters, actions) |
| `PageSectionHeader.vue` | Section heading with title and optional subtitle |
| `PrimaryActionDropdown.vue` | Dropdown for primary actions (e.g., Edit, Disable) |
| `RowActionButton.vue` | Per-row action button (three-dot menu) |
| `TableActionDrawerButton.vue` | Button that opens a side drawer for table row actions |
| `ToolbarIconButton.vue` | Icon-only button for toolbar actions |
| `tablePagination.ts` | Utility for table pagination state (composable/helper) |

---

## Views (`src/views/`)

| View File | Route | Notes |
|---|---|---|
| `LoginView.vue` | `/auth/login` | Login form, guest-only |
| `ResetPasswordView.vue` | `/reset-password` | Password reset form |
| `OverviewView.vue` | `/overview` | Dashboard overview / summary |
| `PatientsView.vue` | `/patients` | Patient list with search, edit drawer (largest view: ~60KB) |
| `UsersView.vue` | `/users` | Admin users management (~55KB) |
| `ReferralsView.vue` | `/referrals` | Referrals management |
| `RolesView.vue` | `/roles` | Role and permission management |
| `SitesView.vue` | `/sites` | Site management |
| `ScreeningsView.vue` | `/screenings` | Screenings (stub/basic) |
| `SettingsView.vue` | `/settings` | User profile settings |

---

## Composables (`src/composables/`)

> 2 composable files exist (names not listed). Likely include shared logic for API calls, pagination, or form handling.
