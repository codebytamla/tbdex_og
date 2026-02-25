# Roles & Permissions

> Source: `backend/migrations/005_backfill_default_role_permissions.sql`, `009_add_referrals_and_screenings_permissions.sql`

---

## Built-in Roles

| Role Name | Description |
|---|---|
| `ADMIN` | Full system access |
| `DOCTOR` | Medical doctor / clinician |
| `CHW` | Community health worker |
| `DATA_ENTRY` | Data entry operator |
| `SUPERVISOR` | Clinical supervisor / manager |

> Custom roles can also be created via the admin API (`POST /admin/roles`).

---

## Permission Catalog

| Permission | Category |
|---|---|
| `patients.create` | Patients |
| `patients.read` | Patients |
| `patients.update` | Patients |
| `patients.add_site` | Patients |
| `patients.disable` | Patients |
| `patients.write` | Patients |
| `forms.create` | Forms |
| `forms.read` | Forms |
| `forms.update` | Forms |
| `forms.disable` | Forms |
| `reports.export` | Reports |
| `users.read` | Users |
| `sites.read` | Sites |
| `roles.read` | Roles |
| `audit.read` | Audit |
| `referrals.read` | Referrals |
| `referrals.create` | Referrals |
| `referrals.cancel` | Referrals |
| `screenings.read` | Screenings |
| `screenings.create` | Screenings |

---

## Default Role → Permission Mappings

| Role | Permissions |
|---|---|
| `ADMIN` | All permissions |
| `DOCTOR` | `patients.create`, `patients.read`, `patients.update`, `patients.add_site`, `patients.write`, `reports.export` |
| `CHW` | `patients.create`, `patients.read`, `patients.update`, `patients.add_site`, `patients.write` |
| `DATA_ENTRY` | `patients.create`, `patients.read`, `patients.update`, `patients.add_site`, `patients.write`, `forms.create`, `forms.read`, `forms.update`, `forms.disable` |
| `SUPERVISOR` | `patients.read`, `reports.export`, `users.read`, `audit.read` |

---

## Notes

- Permissions follow the `resource.action` naming convention.
- The ADMIN role automatically gets all new permissions added to the system.
- Roles can be updated/created by admin users via `PATCH /admin/roles/{id}` and `POST /admin/roles`.
- An ADMIN user must always retain the `ADMIN` role (enforced server-side with `ErrAdminRoleRequired`).
- Site code must be uppercase (`ErrSiteCodeUppercase` enforced on create/update).
