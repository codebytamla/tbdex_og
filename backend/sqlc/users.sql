-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES (lower(@email), @password_hash, @full_name)
RETURNING id, email, full_name, is_active, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, full_name, is_active, created_at
FROM users
WHERE lower(email) = lower(@email);

-- name: GetUserByID :one
SELECT id, email, full_name, is_active, created_at
FROM users
WHERE id = @id;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = now()
WHERE id = @id;

-- name: UpdateUser :one
UPDATE users SET
    full_name  = COALESCE(sqlc.narg('full_name'), full_name),
    updated_at = now()
WHERE id = @id
RETURNING id, email, full_name, is_active, created_at, updated_at;

-- name: SetUserActive :one
UPDATE users SET
    is_active  = @is_active,
    updated_at = now()
WHERE id = @id
RETURNING id, email, full_name, is_active, created_at, updated_at;

-- name: ListUsersPage :many
SELECT u.id, u.email, u.full_name, u.is_active, u.created_at
FROM users u
WHERE
    (
        @search::text = ''
        OR u.email ILIKE '%' || @search::text || '%' ESCAPE '\'
        OR u.full_name ILIKE '%' || @search::text || '%' ESCAPE '\'
    )
    AND (
        NOT @use_is_active::bool
        OR u.is_active = @is_active::bool
    )
    AND (
        NOT @use_role_filter::bool
        OR EXISTS (
            SELECT 1
            FROM user_roles ur
            JOIN roles r ON r.id = ur.role_id
            WHERE ur.user_id = u.id
              AND r.name = ANY(@role_names::text[])
        )
    )
    AND (
        NOT @use_site_filter::bool
        OR EXISTS (
            SELECT 1
            FROM user_sites us
            WHERE us.user_id = u.id
              AND us.site_id = ANY(@site_ids::uuid[])
        )
    )
    AND (
        NOT @use_cursor::bool
        OR (u.created_at, u.id) < (@cursor_created_at::timestamptz, @cursor_id::uuid)
    )
ORDER BY u.created_at DESC, u.id DESC
LIMIT @limit_rows;

-- name: CountUsersFiltered :one
SELECT COUNT(*)
FROM users u
WHERE
    (
        @search::text = ''
        OR u.email ILIKE '%' || @search::text || '%' ESCAPE '\'
        OR u.full_name ILIKE '%' || @search::text || '%' ESCAPE '\'
    )
    AND (
        NOT @use_is_active::bool
        OR u.is_active = @is_active::bool
    )
    AND (
        NOT @use_role_filter::bool
        OR EXISTS (
            SELECT 1
            FROM user_roles ur
            JOIN roles r ON r.id = ur.role_id
            WHERE ur.user_id = u.id
              AND r.name = ANY(@role_names::text[])
        )
    )
    AND (
        NOT @use_site_filter::bool
        OR EXISTS (
            SELECT 1
            FROM user_sites us
            WHERE us.user_id = u.id
              AND us.site_id = ANY(@site_ids::uuid[])
        )
    );

-- name: ListRolesByUserIDs :many
SELECT ur.user_id, r.id, r.name, r.description
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = ANY(@user_ids::uuid[])
ORDER BY ur.user_id, r.name;

-- name: ListSitesByUserIDs :many
SELECT us.user_id, s.id, s.name, s.code
FROM user_sites us
JOIN sites s ON s.id = us.site_id
WHERE us.user_id = ANY(@user_ids::uuid[])
ORDER BY us.user_id, s.name;
