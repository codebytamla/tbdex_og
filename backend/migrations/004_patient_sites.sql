-- +goose Up

-- ============================================================
-- PATIENT <-> SITE assignments
-- A patient can be seen/treated at multiple sites beyond their
-- primary site.
-- ============================================================
CREATE TABLE patient_sites (
    patient_id  UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (patient_id, site_id)
);

CREATE INDEX idx_patient_sites_site ON patient_sites (site_id);
