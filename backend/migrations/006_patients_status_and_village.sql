-- +goose Up

ALTER TABLE patients
    ADD COLUMN IF NOT EXISTS status TEXT,
    ADD COLUMN IF NOT EXISTS village TEXT;

UPDATE patients
SET status = CASE
    WHEN is_active THEN 'active'
    ELSE 'inactive'
END
WHERE status IS NULL;

ALTER TABLE patients
    ALTER COLUMN status SET DEFAULT 'active',
    ALTER COLUMN status SET NOT NULL;

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_patients_status'
          AND conrelid = 'patients'::regclass
    ) THEN
        ALTER TABLE patients
            ADD CONSTRAINT chk_patients_status
            CHECK (status IN ('provisional', 'active', 'inactive', 'deceased'));
    END IF;
END
$$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_patients_status ON patients (status);
