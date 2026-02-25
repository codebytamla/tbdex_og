---
id: 55222E0C-5C17-460F-925A-062B64665F52
title: backend_api_routes
project: BE51B5CC-C223-4123-B4E0-C5AF82572C86
tags: [[]]
created: 2026-02-25T10:16:48Z
updated: 2026-02-25T10:17:19Z
---

# Backend API Routes

> Source: `backend/internal/auth/handlers.go`  
> All routes are registered via Go's `net/http.ServeMux`.  
> Protected routes require a valid JWT via `RequireAuth` middleware.

---

## Auth (Public — no JWT required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/auth/login` | Login with email + password, returns user + tokens |
| `POST` | `/auth/refresh` | Refresh access token using a refresh token |
| `POST` | `/auth/logout` | Logout, revokes the refresh token |
| `POST` | `/auth/forgot-password` | Send password reset token (always returns 200) |
| `POST` | `/auth/reset-password` | Reset password using a reset token |

---

## Auth (Protected — JWT required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/auth/me` | Get current user profile with roles, permissions, and site assignments |

---

## Admin — Users (Protected + Permission gated)

| Method | Path | Description |
|---|---|---|
| `POST` | `/admin/users` | Create a new user |
| `GET` | `/admin/users` | List users (cursor pagination, filterable by `search`, `role`, `site_id`, `is_active`) |
| `PATCH` | `/admin/users/{id}` | Update user (full_name, email, etc.) |
| `PUT` | `/admin/users/{id}/roles` | Set roles for a user (replaces all) |
| `GET` | `/admin/users/{id}/sites` | Get user's site assignments and roles per site |
| `PUT` | `/admin/users/{id}/sites` | Set site assignments for a user (replaces all) |
| `POST` | `/admin/users/{id}/password-reset` | Force a password reset for a user |
| `POST` | `/admin/users/{id}/disable` | Disable a user account |
| `POST` | `/admin/users/{id}/enable` | Enable a disabled user account |

---

## Admin — Sites (Protected + Permission gated)

| Method | Path | Description |
|---|---|---|
| `POST` | `/admin/sites` | Create a new site (code must be uppercase) |
| `GET` | `/admin/sites` | List all sites |
| `GET` | `/admin/sites/{id}` | Get a single site |
| `PUT` | `/admin/sites/{id}` | Update a site |
| `PATCH` | `/admin/sites/{id}/disable` | Disable a site |

---

## Admin — Roles & Permissions (Protected + Permission gated)

| Method | Path | Description |
|---|---|---|
| `GET` | `/admin/roles` | List all roles |
| `POST` | `/admin/roles` | Create a new role |
| `GET` | `/admin/roles/{id}` | Get a single role |
| `PATCH` | `/admin/roles/{id}` | Update a role (name, description, permissions) |
| `GET` | `/admin/permissions` | List all available permissions |

---

## Patients (Protected + Permission gated)

| Method | Path | Description |
|---|---|---|
| `POST` | `/patients` | Create a new patient |
| `GET` | `/patients` | List patients (cursor pagination, filterable by `search`, `status`) |
| `GET` | `/patients/{id}` | Get a single patient |
| `PATCH` | `/patients/{id}` | Update a patient |
| `POST` | `/patients/{id}/disable` | Disable a patient (admin or creator only) |

---

## Query Parameters

### GET `/admin/users`
| Param | Type | Description |
|---|---|---|
| `cursor` | string | Base64 cursor for pagination |
| `limit` | int | Page size (default/max: server-defined) |
| `search` | string | Text search on name/email |
| `role` | string[] | Filter by role name (repeatable) |
| `site_id` | UUID[] | Filter by site assignment (repeatable) |
| `is_active` | bool | `true` or `false` |
| `include_total` | bool | Include total count in response |

### GET `/patients`
| Param | Type | Description |
|---|---|---|
| `cursor` | string | Base64 cursor for pagination |
| `limit` | int | Page size |
| `search` | string | Text search on name |
| `status` | string | `provisional`, `active`, `inactive`, `deceased` |
| `include_total` | bool | Include total count in response |

---

## Error Conventions

Errors are returned as JSON: `{ "error": "message" }`

| Status | Meaning |
|---|---|
| 400 | Bad request / validation failure |
| 401 | Not authenticated |
| 403 | Forbidden |
| 404 | Not found |
| 409 | Conflict (duplicate) |
| 500 | Internal server error |