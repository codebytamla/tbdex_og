-- +goose Up

-- Add referrals and screenings permissions to the catalog.
INSERT INTO permissions (name) VALUES
    ('referrals.read'),
    ('referrals.create'),
    ('referrals.cancel'),
    ('screenings.read'),
    ('screenings.create')
ON CONFLICT (name) DO NOTHING;

-- Ensure ADMIN has access to the new permissions.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
    'referrals.read',
    'referrals.create',
    'referrals.cancel',
    'screenings.read',
    'screenings.create'
)
WHERE r.name = 'ADMIN'
ON CONFLICT DO NOTHING;
