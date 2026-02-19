package main

import (
	"fmt"
	"net/http"

	"github.com/codebytamla/borderform-backend/internal/auth"
	"github.com/codebytamla/borderform-backend/internal/database"
	mw "github.com/codebytamla/borderform-backend/internal/middleware"
)

const apiV1Prefix = "/api/v1"

// registerRoutes wires every route in the application.
func registerRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	apiMux := http.NewServeMux()
	registerPublicRoutes(apiMux, h)
	registerAuthRoutes(apiMux, h, jwtCfg)
	registerUserRoutes(apiMux, h, jwtCfg, store)
	registerSiteRoutes(apiMux, h, jwtCfg, store)
	registerRoleRoutes(apiMux, h, jwtCfg, store)
	registerPermissionRoutes(apiMux, h, jwtCfg, store)
	registerPatientRoutes(apiMux, h, jwtCfg, store)

	mux.Handle(apiV1Prefix+"/", http.StripPrefix(apiV1Prefix, apiMux))
}

// ── Public (no auth) ────────────────────────────────────────────────

func registerPublicRoutes(mux *http.ServeMux, h *auth.Handler) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	// Auth endpoints (login, refresh, logout, forgot/reset password).
	h.RegisterRoutes(mux)
}

// ── Authenticated (auth required, no specific permission) ───────────

func registerAuthRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig) {
	mux.Handle("GET /auth/me", chain(
		h.HandleMe(),
		mw.RequireAuth(jwtCfg),
	))
}

// ── Admin - Users ───────────────────────────────────────────────────
func registerUserRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	mux.Handle("POST /admin/users", chain(
		h.HandleAdminCreateUser(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.create", store),
	))
	mux.Handle("GET /admin/users", chain(
		h.HandleAdminListUsers(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.read", store),
	))
	mux.Handle("PATCH /admin/users/{id}", chain(
		h.HandleAdminUpdateUser(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.update", store),
	))
	mux.Handle("PUT /admin/users/{id}/roles", chain(
		h.HandleAdminSetRoles(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.assign_roles", store),
	))
	mux.Handle("GET /admin/users/{id}/sites", chain(
		h.HandleAdminGetUserSitesAndRoles(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.read", store),
	))
	mux.Handle("PUT /admin/users/{id}/sites", chain(
		h.HandleAdminSetSites(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.assign_sites", store),
	))
	mux.Handle("POST /admin/users/{id}/password-reset", chain(
		h.HandleAdminForcePasswordReset(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.reset_password", store),
	))
	mux.Handle("POST /admin/users/{id}/disable", chain(
		h.HandleAdminDisableUser(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.disable", store),
	))
	mux.Handle("POST /admin/users/{id}/enable", chain(
		h.HandleAdminEnableUser(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "users.disable", store),
	))
}

// ── Sites ────────────────────────────────────────────────────────────
func registerSiteRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	mux.Handle("POST /admin/sites", chain(
		h.HandleAdminCreateSite(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "sites.create", store),
	))
	mux.Handle("GET /admin/sites", chain(
		h.HandleAdminListSites(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "sites.read", store),
	))
	mux.Handle("GET /admin/sites/{id}", chain(
		h.HandleAdminGetSite(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "sites.read", store),
	))
	mux.Handle("PUT /admin/sites/{id}", chain(
		h.HandleAdminUpdateSite(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "sites.update", store),
	))
	mux.Handle("PATCH /admin/sites/{id}/disable", chain(
		h.HandleAdminDisableSite(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "sites.disable", store),
	))
}

// ── Roles ────────────────────────────────────────────────────────────
func registerRoleRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	mux.Handle("GET /admin/roles", chain(
		h.HandleAdminListRoles(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "roles.read", store),
	))
	mux.Handle("POST /admin/roles", chain(
		h.HandleAdminCreateRole(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "roles.create", store),
	))
	mux.Handle("GET /admin/roles/{id}", chain(
		h.HandleAdminGetRole(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "roles.read", store),
	))
	mux.Handle("PATCH /admin/roles/{id}", chain(
		h.HandleAdminUpdateRole(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "roles.update", store),
	))
}

// ── Permissions ─────────────────────────────────────────────────────
func registerPermissionRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	mux.Handle("GET /admin/permissions", chain(
		h.HandleAdminListPermissions(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "permissions.read", store),
	))
}

// ── Patients ─────────────────────────────────────────────────────────
func registerPatientRoutes(mux *http.ServeMux, h *auth.Handler, jwtCfg auth.JWTConfig, store *database.PGStore) {
	mux.Handle("POST /patients", chain(
		h.HandleCreatePatient(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "patients.create", store),
	))
	mux.Handle("GET /patients", chain(
		h.HandleListPatients(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "patients.read", store),
	))
	mux.Handle("GET /patients/{id}", chain(
		h.HandleGetPatient(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "patients.read", store),
	))
	mux.Handle("PATCH /patients/{id}", chain(
		h.HandleUpdatePatient(),
		mw.RequireAuth(jwtCfg),
		mw.RequirePermission(store, "patients.update", store),
	))
	// Disable: RequireAuth only — admin/creator check is in the service layer.
	mux.Handle("POST /patients/{id}/disable", chain(
		h.HandleDisablePatient(),
		mw.RequireAuth(jwtCfg),
	))
}
