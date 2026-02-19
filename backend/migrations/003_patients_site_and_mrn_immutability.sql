-- +goose Up

ALTER TABLE patients
    ADD COLUMN IF NOT EXISTS primary_site_id UUID,
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS origin TEXT,
    ADD COLUMN IF NOT EXISTS metadata JSONB;

ALTER TABLE patients
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb;

-- +goose StatementBegin
DO $$
DECLARE
    default_site_id UUID;
    missing_count BIGINT;
BEGIN
    SELECT id
    INTO default_site_id
    FROM sites
    ORDER BY created_at ASC, id ASC
    LIMIT 1;

    IF default_site_id IS NOT NULL THEN
        UPDATE patients
        SET primary_site_id = default_site_id
        WHERE primary_site_id IS NULL;
    END IF;

    SELECT COUNT(*)
    INTO missing_count
    FROM patients
    WHERE primary_site_id IS NULL;

    IF missing_count > 0 THEN
        RAISE EXCEPTION
            'Cannot set patients.primary_site_id NOT NULL: % row(s) have no site and no fallback site exists',
            missing_count;
    END IF;
END
$$;
-- +goose StatementEnd

ALTER TABLE patients
    ALTER COLUMN primary_site_id SET NOT NULL;

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_patients_primary_site'
          AND conrelid = 'patients'::regclass
    ) THEN
        ALTER TABLE patients
            ADD CONSTRAINT fk_patients_primary_site
            FOREIGN KEY (primary_site_id) REFERENCES sites(id);
    END IF;
END
$$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_patients_primary_site_id ON patients (primary_site_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_mrn ON patients (mrn) WHERE mrn IS NOT NULL;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION prevent_patient_mrn_update()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.mrn IS DISTINCT FROM NEW.mrn THEN
        IF OLD.mrn IS NULL THEN
            RETURN NEW;
        END IF;

        RAISE EXCEPTION 'patients.mrn is immutable once set';
    END IF;

    RETURN NEW;
END
$$;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_patients_mrn_immutable ON patients;

CREATE TRIGGER trg_patients_mrn_immutable
BEFORE UPDATE OF mrn ON patients
FOR EACH ROW
EXECUTE FUNCTION prevent_patient_mrn_update();
