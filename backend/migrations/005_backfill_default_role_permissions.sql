-- +goose Up

-- Ensure default non-admin roles exist.
INSERT INTO roles (name, description) VALUES
    ('DOCTOR', 'Medical doctor / clinician'),
    ('CHW', 'Community health worker'),
    ('DATA_ENTRY', 'Data entry operator'),
    ('SUPERVISOR', 'Clinical supervisor / manager')
ON CONFLICT (name) DO NOTHING;

-- Ensure permission catalog includes baseline clinical permissions.
INSERT INTO permissions (name) VALUES
    ('patients.create'),
    ('patients.read'),
    ('patients.update'),
    ('patients.add_site'),
    ('patients.disable'),
    ('patients.write'),
    ('forms.create'),
    ('forms.read'),
    ('forms.update'),
    ('forms.disable'),
    ('reports.export'),
    ('users.read'),
    ('audit.read')
ON CONFLICT (name) DO NOTHING;

-- Backfill default role -> permission mappings used by the app.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM (
    VALUES
        ('DOCTOR', 'patients.create'),
        ('DOCTOR', 'patients.read'),
        ('DOCTOR', 'patients.update'),
        ('DOCTOR', 'patients.add_site'),
        ('DOCTOR', 'patients.write'),
        ('DOCTOR', 'reports.export'),

        ('CHW', 'patients.create'),
        ('CHW', 'patients.read'),
        ('CHW', 'patients.update'),
        ('CHW', 'patients.add_site'),
        ('CHW', 'patients.write'),

        ('DATA_ENTRY', 'patients.create'),
        ('DATA_ENTRY', 'patients.read'),
        ('DATA_ENTRY', 'patients.update'),
        ('DATA_ENTRY', 'patients.add_site'),
        ('DATA_ENTRY', 'patients.write'),
        ('DATA_ENTRY', 'forms.create'),
        ('DATA_ENTRY', 'forms.read'),
        ('DATA_ENTRY', 'forms.update'),
        ('DATA_ENTRY', 'forms.disable'),

        ('SUPERVISOR', 'patients.read'),
        ('SUPERVISOR', 'reports.export'),
        ('SUPERVISOR', 'users.read'),
        ('SUPERVISOR', 'audit.read')
) AS rp(role_name, permission_name)
JOIN roles r ON r.name = rp.role_name
JOIN permissions p ON p.name = rp.permission_name
ON CONFLICT DO NOTHING;
