-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES (@user_id, @role_id)
ON CONFLICT DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = @user_id AND role_id = @role_id;

-- name: RemoveAllUserRoles :exec
DELETE FROM user_roles WHERE user_id = @user_id;

-- name: GetRoleByName :one
SELECT id, name, description FROM roles WHERE name = @name;

-- name: GetRoleByID :one
SELECT id, name, description FROM roles WHERE id = @id;

-- name: ListRoles :many
SELECT id, name, description FROM roles ORDER BY name;

-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES (@name, @description)
RETURNING id, name, description;

-- name: UpdateRole :one
UPDATE roles SET
    name        = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = @id
RETURNING id, name, description;

-- name: ListUserRoles :many
SELECT r.id, r.name, r.description
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = @user_id;

-- name: GetEffectivePermissions :many
SELECT DISTINCT p.name
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = @user_id;

-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM permissions p
    JOIN role_permissions rp ON rp.permission_id = p.id
    JOIN user_roles ur ON ur.role_id = rp.role_id
    WHERE ur.user_id = @user_id AND p.name = @permission_name
) AS has_permission;

-- name: ListRolePermissions :many
SELECT p.id, p.name
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = @role_id
ORDER BY p.name;

-- name: ListAllPermissions :many
SELECT id, name FROM permissions ORDER BY name;

-- name: DeleteRolePermissions :exec
DELETE FROM role_permissions WHERE role_id = @role_id;

-- name: InsertRolePermissionByName :exec
INSERT INTO role_permissions (role_id, permission_id)
SELECT @role_id, p.id FROM permissions p WHERE p.name = @permission_name
ON CONFLICT DO NOTHING;
