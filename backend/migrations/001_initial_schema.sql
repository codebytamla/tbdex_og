-- +goose Up

-- Extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL,
    password_hash   TEXT NOT NULL,       -- argon2id encoded string
    full_name       TEXT NOT NULL DEFAULT '',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_users_email ON users (lower(email));

-- ============================================================
-- SITES
-- ============================================================
CREATE TABLE sites (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    code            TEXT NOT NULL,       -- short unique code e.g. "CLI-01"
    site_type       TEXT NOT NULL DEFAULT 'clinic',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_sites_code ON sites (lower(code));

-- ============================================================
-- USER <-> SITE assignments
-- ============================================================
CREATE TABLE user_sites (
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, site_id)
);

CREATE INDEX idx_user_sites_site ON user_sites (site_id);

-- ============================================================
-- ROLES
-- ============================================================
CREATE TABLE roles (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    TEXT NOT NULL UNIQUE,           -- e.g. ADMIN, DOCTOR
    description TEXT NOT NULL DEFAULT ''
);

-- ============================================================
-- PERMISSIONS
-- ============================================================
CREATE TABLE permissions (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    TEXT NOT NULL UNIQUE            -- e.g. users.create, sites.read
);

-- ============================================================
-- ROLE <-> PERMISSIONS
-- ============================================================
CREATE TABLE role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_perm ON role_permissions (permission_id);

-- ============================================================
-- USER <-> ROLES
-- ============================================================
CREATE TABLE user_roles (
    user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id  UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_role ON user_roles (role_id);

-- ============================================================
-- REFRESH TOKENS
-- One row per device/session. The token_hash stores SHA-256 of
-- the opaque token (we don't need argon2 here because the token
-- is high-entropy random, not a human-chosen password).
-- ============================================================
CREATE TABLE refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      TEXT NOT NULL,          -- SHA-256 hex of opaque token
    device_name     TEXT NOT NULL DEFAULT 'unknown',
    created_ip      TEXT NOT NULL DEFAULT '',
    user_agent      TEXT NOT NULL DEFAULT '',
    expires_at      TIMESTAMPTZ NOT NULL,
    last_used_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked         BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens (user_id) WHERE NOT revoked;
CREATE UNIQUE INDEX idx_refresh_tokens_hash ON refresh_tokens (token_hash) WHERE NOT revoked;

-- ============================================================
-- PASSWORD RESET TOKENS
-- Short-lived, single-use. Store hashed like refresh tokens.
-- ============================================================
CREATE TABLE password_reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL,              -- SHA-256 hex
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_password_reset_tokens_user ON password_reset_tokens (user_id);
CREATE UNIQUE INDEX idx_password_reset_tokens_hash ON password_reset_tokens (token_hash) WHERE NOT used;

-- ============================================================
-- AUDIT EVENTS (append-only)
-- ============================================================
CREATE TABLE audit_events (
    id          BIGSERIAL PRIMARY KEY,
    event_type  TEXT NOT NULL,              -- login_success, login_fail, refresh, logout, etc.
    actor_id    UUID REFERENCES users(id),  -- NULL for failed login by unknown email
    target_id   UUID REFERENCES users(id),  -- the user affected (usually same as actor)
    ip_address  TEXT NOT NULL DEFAULT '',
    user_agent  TEXT NOT NULL DEFAULT '',
    metadata    JSONB NOT NULL DEFAULT '{}', -- extra context
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_events_actor ON audit_events (actor_id, created_at DESC);
CREATE INDEX idx_audit_events_type  ON audit_events (event_type, created_at DESC);

-- ============================================================
-- PATIENTS
-- ============================================================
CREATE TABLE patients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       TEXT NOT NULL,
    date_of_birth   DATE,
    is_dob_estimated BOOLEAN NOT NULL DEFAULT false,
    sex             TEXT CHECK (sex IN ('male', 'female', 'other')),
    phone           TEXT,
    national_id     TEXT,
    mrn             TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Prevent duplicates
CREATE UNIQUE INDEX idx_patients_national_id ON patients (national_id) WHERE national_id IS NOT NULL;
CREATE UNIQUE INDEX idx_patients_mrn ON patients (mrn) WHERE mrn IS NOT NULL;
-- Common lookups
CREATE INDEX idx_patients_full_name ON patients (full_name);
CREATE INDEX idx_patients_dob ON patients (date_of_birth);
CREATE INDEX idx_patients_is_active ON patients (is_active) WHERE is_active = true;

-- ============================================================
-- SEED GRANULAR PERMISSIONS
-- ============================================================

