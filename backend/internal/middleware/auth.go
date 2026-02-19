package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/codebytamla/borderform-backend/internal/auth"
	"github.com/google/uuid"
)

// UserIDFromContext extracts the authenticated user's ID from context.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(auth.CtxUserIDKey).(uuid.UUID)
	return id, ok
}

// RequireAuth validates the JWT in the Authorization header and injects
// the user ID into the request context.
func RequireAuth(cfg auth.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				writeErr(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			claims, err := auth.VerifyAccessToken(cfg, token)
			if err != nil {
				writeErr(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), auth.CtxUserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// PermissionChecker resolves whether a user has a given permission.
type PermissionChecker interface {
	UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
}

// RequirePermission checks that the authenticated user has the specified
// permission string (e.g. "patients.read"). Must be chained after RequireAuth.
func RequirePermission(checker PermissionChecker, permission string, auditW auth.AuditWriter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				writeErr(w, http.StatusUnauthorized, "not authenticated")
				return
			}

			has, err := checker.UserHasPermission(r.Context(), userID, permission)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, "internal error")
				return
			}
			if !has {
				auth.LogAudit(r.Context(), auditW, auth.AuditPermissionDenied, &userID, nil,
					r.RemoteAddr, r.UserAgent(),
					map[string]any{"permission": permission, "path": r.URL.Path})
				writeErr(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSite checks that the authenticated user is assigned to the
// site_id found in the request (query param or path). Must be chained
// after RequireAuth.
type SiteChecker interface {
	UserBelongsToSite(ctx context.Context, userID, siteID uuid.UUID) (bool, error)
}

func RequireSite(checker SiteChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				writeErr(w, http.StatusUnauthorized, "not authenticated")
				return
			}

			siteStr := r.URL.Query().Get("site_id")
			if siteStr == "" {
				siteStr = r.PathValue("site_id")
			}
			if siteStr == "" {
				// No site scope requested — pass through (global scope).
				next.ServeHTTP(w, r)
				return
			}

			siteID, err := uuid.Parse(siteStr)
			if err != nil {
				writeErr(w, http.StatusBadRequest, "invalid site_id")
				return
			}

			belongs, err := checker.UserBelongsToSite(r.Context(), userID, siteID)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, "internal error")
				return
			}
			if !belongs {
				writeErr(w, http.StatusForbidden, "not assigned to this site")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

func extractBearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(h, "Bearer ")
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + msg + `"}`)) //nolint:errcheck
}
