package auth

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuditEventType constants for the event_type column.
const (
	AuditLoginSuccess       = "login_success"
	AuditLoginFail          = "login_fail"
	AuditRefresh            = "refresh"
	AuditLogout             = "logout"
	AuditRegister           = "register"
	AuditPasswordResetReq   = "password_reset_request"
	AuditPasswordResetDone  = "password_reset_success"
	AuditRoleChange         = "role_change"
	AuditPermissionDenied   = "permission_denied"
)

// AuditWriter is the interface for writing audit events.
// This decouples the audit logic from the concrete sqlc-generated code.
type AuditWriter interface {
	InsertAuditEvent(ctx context.Context, params InsertAuditEventParams) error
}

// InsertAuditEventParams mirrors the sqlc-generated params so we can
// use this in our service layer without importing the db package directly.
type InsertAuditEventParams struct {
	EventType  string
	ActorID    pgtype.UUID
	TargetID   pgtype.UUID
	IPAddress  string
	UserAgent  string
	Metadata   []byte // JSON
}

// LogAudit is a convenience helper to fire-and-forget an audit event.
func LogAudit(ctx context.Context, w AuditWriter, eventType string, actorID, targetID *uuid.UUID, ip, ua string, meta map[string]any) {
	params := InsertAuditEventParams{
		EventType: eventType,
		IPAddress: ip,
		UserAgent: ua,
	}
	if actorID != nil {
		params.ActorID = pgtype.UUID{Bytes: *actorID, Valid: true}
	}
	if targetID != nil {
		params.TargetID = pgtype.UUID{Bytes: *targetID, Valid: true}
	}
	if meta != nil {
		params.Metadata, _ = json.Marshal(meta)
	} else {
		params.Metadata = []byte("{}")
	}

	// Best-effort: don't fail the request if audit logging fails.
	_ = w.InsertAuditEvent(ctx, params)
}
