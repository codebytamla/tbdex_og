---
id: 5724040A-5611-4C03-9678-6E5E6F6C646C
title: frontend_routes
project: BE51B5CC-C223-4123-B4E0-C5AF82572C86
tags: [[]]
created: 2026-02-25T10:16:50Z
updated: 2026-02-25T10:18:33Z
---

# Frontend Routes (Dashboard)

> Source: `dashboard/src/router/index.ts`  
> Built with Vue Router (`createWebHistory`).

---

## Route Table

| Path | Name | View Component | Auth | Required Permissions |
|---|---|---|---|---|
| `/` | — | — | — | Redirects → `/overview` |
| `/auth/login` | `login` | `LoginView.vue` | Guest only | — |
| `/reset-password` | `reset-password` | `ResetPasswordView.vue` | — | — |
| `/overview` | `overview` | `OverviewView.vue` | ✅ Required | — |
| `/dashboard` | — | — | — | Redirects → `/overview` |
| `/referrals` | `referrals` | `ReferralsView.vue` | ✅ Required | `referrals.read` |
| `/screenings` | `screenings` | `ScreeningsView.vue` | ✅ Required | `screenings.read` |
| `/patients` | `patients` | `PatientsView.vue` | ✅ Required | `patients.read` |
| `/sites` | `sites` | `SitesView.vue` | ✅ Required | `sites.read` |
| `/users` | `users` | `UsersView.vue` | ✅ Required | `users.read` |
| `/roles` | `roles` | `RolesView.vue` | ✅ Required | `roles.read` |
| `/settings` | `settings` | `SettingsView.vue` | ✅ Required | — |
| `/:pathMatch(.*)` | — | — | — | Redirects → `/overview` |

---

## Legacy `/admin/*` Redirects

| Old Path | Redirects To |
|---|---|
| `/admin` | `/overview` |
| `/admin/referrals` | `/referrals` |
| `/admin/screenings` | `/screenings` |
| `/admin/patients` | `/patients` |
| `/admin/sites` | `/sites` |
| `/admin/users` | `/users` |
| `/admin/roles` | `/roles` |

---

## Navigation Guard Logic

- Routes with `requiresAuth: true` → redirect to `/auth/login` with `?redirect=<path>` if not authenticated.
- After auth, `authStore.ensureProfile()` is called; on failure → force logout.
- Routes with `guestOnly: true` → redirect authenticated users to `authStore.defaultOverviewRoute()`.
- Routes with `requiredPermissions` → checked via `authStore.canAccess(permissions, roles)`.