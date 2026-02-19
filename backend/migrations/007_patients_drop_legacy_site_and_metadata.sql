-- +goose Up

ALTER TABLE patients
    DROP CONSTRAINT IF EXISTS fk_patients_primary_site;

DROP INDEX IF EXISTS idx_patients_primary_site_id;
DROP INDEX IF EXISTS idx_patients_is_active;

ALTER TABLE patients
    DROP COLUMN IF EXISTS primary_site_id,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS metadata;

DROP TABLE IF EXISTS patient_sites;

DELETE FROM permissions
WHERE name = 'patients.add_site';
