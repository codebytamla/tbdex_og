---
id: 512A1312-5611-4C01-8563-61486B71656C
title: database_tables
project: BE51B5CC-C223-4123-B4E0-C5AF82572C86
tags: [[]]
created: 2026-02-25T10:16:38Z
updated: 2026-02-25T10:17:39Z
---

# Database Tables

> Source: `backend/migrations/`

---

## users
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK, `gen_random_uuid()` |
| email | TEXT | Unique (case-insensitive index) |
| password_hash | TEXT | argon2id encoded string |
| full_name | TEXT | Default `''` |
| is_active | BOOLEAN | Default `true` |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

## sites
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | TEXT | |
| code | TEXT | Short unique code e.g. `CLI-01`, must be uppercase |
| site_type | TEXT | Default `'clinic'` |
| is_active | BOOLEAN | Default `true` |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

## user_sites _(junction)_
| Column | Type | Notes |
|---|---|---|
| user_id | UUID | FK → users.id |
| site_id | UUID | FK → sites.id |
| assigned_at | TIMESTAMPTZ | |

**PK:** `(user_id, site_id)`

---

## roles
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | TEXT | Unique. Built-in: `ADMIN`, `DOCTOR`, `CHW`, `DATA_ENTRY`, `SUPERVISOR` |
| description | TEXT | Default `''` |

---

## permissions
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| name | TEXT | Unique. Format `resource.action` |

**Known permissions:**
- `patients.create`, `patients.read`, `patients.update`, `patients.add_site`, `patients.disable`, `patients.write`
- `forms.create`, `forms.read`, `forms.update`, `forms.disable`
- `reports.export`
- `users.read`
- `audit.read`
- `sites.read` _(used in frontend router)_
- `roles.read` _(used in frontend router)_
- `referrals.read`, `referrals.create`, `referrals.cancel`
- `screenings.read`, `screenings.create`

---

## role_permissions _(junction)_
| Column | Type | Notes |
|---|---|---|
| role_id | UUID | FK → roles.id |
| permission_id | UUID | FK → permissions.id |

**PK:** `(role_id, permission_id)`

---

## user_roles _(junction)_
| Column | Type | Notes |
|---|---|---|
| user_id | UUID | FK → users.id |
| role_id | UUID | FK → roles.id |
| assigned_at | TIMESTAMPTZ | |

**PK:** `(user_id, role_id)`

---

## refresh_tokens
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → users.id |
| token_hash | TEXT | SHA-256 hex of the opaque token |
| device_name | TEXT | Default `'unknown'` |
| created_ip | TEXT | |
| user_agent | TEXT | |
| expires_at | TIMESTAMPTZ | |
| last_used_at | TIMESTAMPTZ | |
| revoked | BOOLEAN | Default `false` |
| created_at | TIMESTAMPTZ | |

---

## password_reset_tokens
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → users.id |
| token_hash | TEXT | SHA-256 hex |
| expires_at | TIMESTAMPTZ | Short-lived |
| used | BOOLEAN | Default `false`, single-use |
| created_at | TIMESTAMPTZ | |

---

## audit_events _(append-only)_
| Column | Type | Notes |
|---|---|---|
| id | BIGSERIAL | PK |
| event_type | TEXT | e.g. `login_success`, `login_fail`, `refresh`, `logout` |
| actor_id | UUID | FK → users.id (nullable for failed logins) |
| target_id | UUID | FK → users.id (user affected) |
| ip_address | TEXT | |
| user_agent | TEXT | |
| metadata | JSONB | Extra context |
| created_at | TIMESTAMPTZ | |

---

## patients
| Column | Type | Notes |
|---|---|---|
| id | UUID | PK |
| full_name | TEXT | |
| date_of_birth | DATE | Nullable |
| is_dob_estimated | BOOLEAN | Default `false` |
| sex | TEXT | `CHECK IN ('male','female','other')` |
| phone | TEXT | Nullable |
| national_id | TEXT | Nullable, unique |
| mrn | TEXT | Nullable, unique |
| status | TEXT | `CHECK IN ('provisional','active','inactive','deceased')`, Default `'active'` |
| village | TEXT | Nullable |
| is_active | BOOLEAN | Default `true` |
| created_by | UUID | FK → users.id |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

## patient_sites _(junction)_
| Column | Type | Notes |
|---|---|---|
| patient_id | UUID | FK → patients.id |
| site_id | UUID | FK → sites.id |
| assigned_at | TIMESTAMPTZ | |

**PK:** `(patient_id, site_id)`

---

## Migration History

| File | Description |
|---|---|
| `001_initial_schema.sql` | Core tables: users, sites, user_sites, roles, permissions, role_permissions, user_roles, refresh_tokens, password_reset_tokens, audit_events, patients |
| `002_sites_code_uppercase.sql` | Enforce uppercase site codes |
| `003_patients_site_and_mrn_immutability.sql` | MRN immutability constraints |
| `004_patient_sites.sql` | Add `patient_sites` junction table |
| `005_backfill_default_role_permissions.sql` | Seed default roles & permissions |
| `006_patients_status_and_village.sql` | Add `status` and `village` to patients |
| `007_patients_drop_legacy_site_and_metadata.sql` | Drop legacy patient site/metadata columns |
| `008_remove_encounter_permissions.sql` | Remove encounter-related permissions |
| `009_add_referrals_and_screenings_permissions.sql` | Add referral & screening permissions |