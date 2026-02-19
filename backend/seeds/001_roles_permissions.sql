-- Seed: default roles
INSERT INTO roles (name, description) VALUES
    ('ADMIN',       'Administrator')
ON CONFLICT (name) DO NOTHING;

-- Seed: default permissions
INSERT INTO permissions (name) VALUES
    ('users.create'),
    ('users.read'),
    ('users.update'),
    ('users.disable'),
    ('users.assign_roles'),
    ('users.assign_sites'),
    ('users.reset_password'),
    ('sites.create'),
    ('sites.read'),
    ('sites.update'),
    ('sites.disable'),
    ('roles.create'),
    ('roles.read'),
    ('roles.update'),
    ('roles.disable'),
    ('patients.create'),
    ('patients.read'),
    ('patients.update'),
    ('patients.disable'),
    ('permissions.read')
ON CONFLICT (name) DO NOTHING;

-- Seed: role <-> permission mappings
-- ADMIN gets everything
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r CROSS JOIN permissions p
WHERE r.name = 'ADMIN'
ON CONFLICT DO NOTHING;
