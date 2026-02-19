-- name: InsertAuditEvent :exec
INSERT INTO audit_events (event_type, actor_id, target_id, ip_address, user_agent, metadata)
VALUES (@event_type, @actor_id, @target_id, @ip_address, @user_agent, @metadata);
