-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, device_name, created_ip, user_agent, expires_at)
VALUES (@user_id, @token_hash, @device_name, @created_ip, @user_agent, @expires_at)
RETURNING id, created_at;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, expires_at, revoked
FROM refresh_tokens
WHERE token_hash = @token_hash AND NOT revoked;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked = true WHERE id = @id;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens SET revoked = true
WHERE user_id = @user_id AND NOT revoked;

-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES (@user_id, @token_hash, @expires_at);

-- name: GetPasswordResetTokenByHash :one
SELECT id, user_id, expires_at, used
FROM password_reset_tokens
WHERE token_hash = @token_hash AND NOT used;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens SET used = true WHERE id = @id;

-- name: ExpireOldPasswordResetTokens :exec
UPDATE password_reset_tokens SET used = true
WHERE user_id = @user_id AND NOT used;
