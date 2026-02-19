-- name: AssignUserToSite :exec
INSERT INTO user_sites (user_id, site_id)
VALUES (@user_id, @site_id)
ON CONFLICT DO NOTHING;

-- name: RemoveUserFromSite :exec
DELETE FROM user_sites
WHERE user_id = @user_id AND site_id = @site_id;

-- name: RemoveAllUserSites :exec
DELETE FROM user_sites WHERE user_id = @user_id;

-- name: ListUserSites :many
SELECT s.id, s.name, s.code
FROM sites s
JOIN user_sites us ON us.site_id = s.id
WHERE us.user_id = @user_id;

-- name: UserBelongsToSite :one
SELECT EXISTS (
    SELECT 1 FROM user_sites
    WHERE user_id = @user_id AND site_id = @site_id
) AS belongs;

-- name: GetSiteSummaryByID :one
SELECT id, name, code FROM sites WHERE id = @id;

-- name: CreateSite :one
INSERT INTO sites (name, code, site_type)
VALUES (@name, @code, @site_type)
RETURNING id, name, code, site_type, is_active, created_at, updated_at;

-- name: ListSites :many
SELECT id, name, code, site_type, is_active, created_at, updated_at
FROM sites
ORDER BY created_at DESC;

-- name: GetSiteByID :one
SELECT id, name, code, site_type, is_active, created_at, updated_at
FROM sites
WHERE id = @id;

-- name: UpdateSite :one
UPDATE sites
SET name = @name, code = @code, site_type = @site_type, updated_at = now()
WHERE id = @id
RETURNING id, name, code, site_type, is_active, created_at, updated_at;

-- name: DisableSite :one
UPDATE sites
SET is_active = false, updated_at = now()
WHERE id = @id
RETURNING id, name, code, site_type, is_active, created_at, updated_at;
