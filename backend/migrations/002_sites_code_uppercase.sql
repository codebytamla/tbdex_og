-- +goose Up

-- Normalize existing data first so the new check constraint can be applied.
UPDATE sites
SET code = upper(code)
WHERE code <> upper(code);

-- Enforce uppercase-only site codes going forward.
ALTER TABLE sites
    ADD CONSTRAINT chk_sites_code_uppercase
    CHECK (code = upper(code));
