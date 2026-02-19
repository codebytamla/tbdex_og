package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Handler wires auth endpoints to the Service.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers public auth routes (no JWT required).
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.handleLogin)
	mux.HandleFunc("POST /auth/refresh", h.handleRefresh)
	mux.HandleFunc("POST /auth/logout", h.handleLogout)
	mux.HandleFunc("POST /auth/forgot-password", h.handleForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", h.handleResetPassword)
}

// HandleMe returns the http.HandlerFunc for GET /auth/me.
// Exported so main.go can wrap it with RequireAuth middleware.
func (h *Handler) HandleMe() http.HandlerFunc {
	return h.handleMe
}

// HandleAdminCreateUser returns handler for POST /admin/users.
func (h *Handler) HandleAdminCreateUser() http.HandlerFunc {
	return h.handleAdminCreateUser
}

// HandleAdminListUsers returns handler for GET /admin/users.
func (h *Handler) HandleAdminListUsers() http.HandlerFunc {
	return h.handleAdminListUsers
}

// HandleAdminUpdateUser returns handler for PATCH /admin/users/{id}.
func (h *Handler) HandleAdminUpdateUser() http.HandlerFunc {
	return h.handleAdminUpdateUser
}

// HandleAdminSetRoles returns handler for PUT /admin/users/{id}/roles.
func (h *Handler) HandleAdminSetRoles() http.HandlerFunc {
	return h.handleAdminSetRoles
}

// HandleAdminGetUserSitesAndRoles returns handler for GET /admin/users/{id}/sites.
func (h *Handler) HandleAdminGetUserSitesAndRoles() http.HandlerFunc {
	return h.handleAdminGetUserSitesAndRoles
}

// HandleAdminSetSites returns handler for PUT /admin/users/{id}/sites.
func (h *Handler) HandleAdminSetSites() http.HandlerFunc {
	return h.handleAdminSetSites
}

// HandleAdminForcePasswordReset returns handler for POST /admin/users/{id}/password-reset.
func (h *Handler) HandleAdminForcePasswordReset() http.HandlerFunc {
	return h.handleAdminForcePasswordReset
}

// HandleAdminDisableUser returns handler for POST /admin/users/{id}/disable.
func (h *Handler) HandleAdminDisableUser() http.HandlerFunc {
	return h.handleAdminDisableUser
}

// HandleAdminEnableUser returns handler for POST /admin/users/{id}/enable.
func (h *Handler) HandleAdminEnableUser() http.HandlerFunc {
	return h.handleAdminEnableUser
}

// HandleAdminCreateSite returns handler for POST /admin/sites.
func (h *Handler) HandleAdminCreateSite() http.HandlerFunc {
	return h.handleAdminCreateSite
}

// HandleAdminListSites returns handler for GET /admin/sites.
func (h *Handler) HandleAdminListSites() http.HandlerFunc {
	return h.handleAdminListSites
}

// HandleAdminGetSite returns handler for GET /admin/sites/{id}.
func (h *Handler) HandleAdminGetSite() http.HandlerFunc {
	return h.handleAdminGetSite
}

// HandleAdminUpdateSite returns handler for PUT /admin/sites/{id}.
func (h *Handler) HandleAdminUpdateSite() http.HandlerFunc {
	return h.handleAdminUpdateSite
}

// HandleAdminDisableSite returns handler for PATCH /admin/sites/{id}/disable.
func (h *Handler) HandleAdminDisableSite() http.HandlerFunc {
	return h.handleAdminDisableSite
}

// HandleAdminListRoles returns handler for GET /admin/roles.
func (h *Handler) HandleAdminListRoles() http.HandlerFunc {
	return h.handleAdminListRoles
}

// HandleAdminCreateRole returns handler for POST /admin/roles.
func (h *Handler) HandleAdminCreateRole() http.HandlerFunc {
	return h.handleAdminCreateRole
}

// HandleAdminGetRole returns handler for GET /admin/roles/{id}.
func (h *Handler) HandleAdminGetRole() http.HandlerFunc {
	return h.handleAdminGetRole
}

// HandleAdminUpdateRole returns handler for PATCH /admin/roles/{id}.
func (h *Handler) HandleAdminUpdateRole() http.HandlerFunc {
	return h.handleAdminUpdateRole
}

// HandleAdminListPermissions returns handler for GET /admin/permissions.
func (h *Handler) HandleAdminListPermissions() http.HandlerFunc {
	return h.handleAdminListPermissions
}

// HandleCreatePatient returns handler for POST /patients.
func (h *Handler) HandleCreatePatient() http.HandlerFunc {
	return h.handleCreatePatient
}

// HandleGetPatient returns handler for GET /patients/{id}.
func (h *Handler) HandleGetPatient() http.HandlerFunc {
	return h.handleGetPatient
}

// HandleListPatients returns handler for GET /patients.
func (h *Handler) HandleListPatients() http.HandlerFunc {
	return h.handleListPatients
}

// HandleUpdatePatient returns handler for PATCH /patients/{id}.
func (h *Handler) HandleUpdatePatient() http.HandlerFunc {
	return h.handleUpdatePatient
}

// HandleDisablePatient returns handler for POST /patients/{id}/disable.
func (h *Handler) HandleDisablePatient() http.HandlerFunc {
	return h.handleDisablePatient
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	user, tokens, err := h.svc.Login(r.Context(), req, dev)
	if err != nil {
		switch err {
		case ErrInvalidCredentials, ErrUserNotActive:
			writeError(w, http.StatusUnauthorized, "invalid email or password")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	dev := extractDeviceInfo(r)
	tokens, err := h.svc.Refresh(r.Context(), req.RefreshToken, dev)
	if err != nil {
		switch err {
		case ErrTokenRevoked, ErrTokenExpired:
			writeError(w, http.StatusUnauthorized, "refresh token invalid or expired")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	if err := h.svc.Logout(r.Context(), req.RefreshToken, dev); err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	token, err := h.svc.ForgotPassword(r.Context(), req, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	// Always return 200 to avoid email enumeration.
	// In dev mode, include the token; in prod, you'd send it via email.
	resp := map[string]string{"status": "ok"}
	if token != "" {
		resp["reset_token"] = token // remove in production
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		writeError(w, http.StatusBadRequest, "token and new_password are required")
		return
	}

	dev := extractDeviceInfo(r)
	if err := h.svc.ResetPassword(r.Context(), req, dev); err != nil {
		switch err {
		case ErrPasswordTooShort:
			writeError(w, http.StatusBadRequest, err.Error())
		case ErrInvalidResetToken:
			writeError(w, http.StatusBadRequest, "invalid or expired reset token")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// -----------------------------------------------------------------
// GET /auth/me
// -----------------------------------------------------------------

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	profile, err := h.svc.GetMe(r.Context(), userID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// -----------------------------------------------------------------
// POST /admin/users
// -----------------------------------------------------------------

func (h *Handler) handleAdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var req AdminCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.FullName == "" {
		writeError(w, http.StatusBadRequest, "email and full_name are required")
		return
	}

	dev := extractDeviceInfo(r)
	user, resetToken, err := h.svc.AdminCreateUser(r.Context(), req, dev)
	if err != nil {
		switch err {
		case ErrPasswordTooShort:
			writeError(w, http.StatusBadRequest, err.Error())
		case ErrUserExists:
			writeError(w, http.StatusConflict, "user already exists")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	resp := map[string]any{"user": user}
	if resetToken != "" {
		resp["reset_token"] = resetToken
	}
	writeJSON(w, http.StatusCreated, resp)
}

// -----------------------------------------------------------------
// GET /admin/users
// -----------------------------------------------------------------

func (h *Handler) handleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	req, err := parseAdminListUsersRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.AdminListUsers(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			writeError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func parseAdminListUsersRequest(r *http.Request) (AdminListUsersRequest, error) {
	q := r.URL.Query()

	req := AdminListUsersRequest{
		Limit:  AdminUsersDefaultLimit,
		Cursor: strings.TrimSpace(q.Get("cursor")),
		Search: strings.TrimSpace(q.Get("search")),
		Roles:  normalizeNonEmptyStrings(q["role"]),
	}

	if limitRaw := strings.TrimSpace(q.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return AdminListUsersRequest{}, errors.New("invalid limit")
		}
		if limit <= 0 {
			limit = AdminUsersDefaultLimit
		}
		if limit > AdminUsersMaxLimit {
			limit = AdminUsersMaxLimit
		}
		req.Limit = limit
	}

	if isActiveRaw := strings.TrimSpace(q.Get("is_active")); isActiveRaw != "" {
		isActive, ok := parseStrictBool(isActiveRaw)
		if !ok {
			return AdminListUsersRequest{}, errors.New("invalid is_active")
		}
		req.IsActive = &isActive
	}

	if includeTotalRaw := strings.TrimSpace(q.Get("include_total")); includeTotalRaw != "" {
		includeTotal, ok := parseStrictBool(includeTotalRaw)
		if !ok {
			return AdminListUsersRequest{}, errors.New("invalid include_total")
		}
		req.IncludeTotal = includeTotal
	}

	siteIDs := q["site_id"]
	req.SiteIDs = make([]uuid.UUID, 0, len(siteIDs))
	for _, raw := range siteIDs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		siteID, err := uuid.Parse(raw)
		if err != nil {
			return AdminListUsersRequest{}, errors.New("invalid site_id")
		}
		req.SiteIDs = append(req.SiteIDs, siteID)
	}

	return req, nil
}

func parseStrictBool(value string) (bool, bool) {
	switch value {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

// -----------------------------------------------------------------
// PATCH /admin/users/{id}
// -----------------------------------------------------------------

func (h *Handler) handleAdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req AdminUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	user, err := h.svc.AdminUpdateUser(r.Context(), userID, req, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// -----------------------------------------------------------------
// PUT /admin/users/{id}/roles
// -----------------------------------------------------------------

func (h *Handler) handleAdminSetRoles(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req SetRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	roles, err := h.svc.AdminSetRoles(r.Context(), userID, req.Roles, dev)
	if err != nil {
		switch err {
		case ErrAdminRoleRequired:
			writeError(w, http.StatusConflict, "admin user must keep ADMIN role")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"roles": roles})
}

// -----------------------------------------------------------------
// GET /admin/users/{id}/sites
// -----------------------------------------------------------------

func (h *Handler) handleAdminGetUserSitesAndRoles(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	assignments, err := h.svc.AdminGetUserSitesAndRoles(r.Context(), userID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, assignments)
}

// -----------------------------------------------------------------
// PUT /admin/users/{id}/sites
// -----------------------------------------------------------------

func (h *Handler) handleAdminSetSites(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req SetSitesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	sites, err := h.svc.AdminSetSites(r.Context(), userID, req.Sites, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"sites": sites})
}

// -----------------------------------------------------------------
// POST /admin/users/{id}/password-reset
// -----------------------------------------------------------------

func (h *Handler) handleAdminForcePasswordReset(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	dev := extractDeviceInfo(r)
	token, err := h.svc.AdminForcePasswordReset(r.Context(), userID, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"reset_token": token})
}

// -----------------------------------------------------------------
// POST /admin/users/{id}/disable
// -----------------------------------------------------------------

func (h *Handler) handleAdminDisableUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	dev := extractDeviceInfo(r)
	user, err := h.svc.AdminDisableUser(r.Context(), userID, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// -----------------------------------------------------------------
// POST /admin/users/{id}/enable
// -----------------------------------------------------------------

func (h *Handler) handleAdminEnableUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	dev := extractDeviceInfo(r)
	user, err := h.svc.AdminEnableUser(r.Context(), userID, dev)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// -----------------------------------------------------------------
// Site management (Admin)
// -----------------------------------------------------------------

func (h *Handler) handleAdminCreateSite(w http.ResponseWriter, r *http.Request) {
	var req AdminCreateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "name and code are required")
		return
	}

	dev := extractDeviceInfo(r)
	site, err := h.svc.AdminCreateSite(r.Context(), req, dev)
	if err != nil {
		switch err {
		case ErrSiteExists:
			writeError(w, http.StatusConflict, "site with this code already exists")
		case ErrSiteCodeUppercase:
			writeError(w, http.StatusBadRequest, "site code must be uppercase")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, site)
}

func (h *Handler) handleAdminListSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.svc.AdminListSites(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"sites": sites})
}

func (h *Handler) handleAdminGetSite(w http.ResponseWriter, r *http.Request) {
	siteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid site id")
		return
	}

	site, err := h.svc.AdminGetSite(r.Context(), siteID)
	if err != nil {
		switch err {
		case ErrSiteNotFound:
			writeError(w, http.StatusNotFound, "site not found")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, site)
}

func (h *Handler) handleAdminUpdateSite(w http.ResponseWriter, r *http.Request) {
	siteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid site id")
		return
	}

	var req AdminUpdateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "name and code are required")
		return
	}

	dev := extractDeviceInfo(r)
	site, err := h.svc.AdminUpdateSite(r.Context(), siteID, req, dev)
	if err != nil {
		switch err {
		case ErrSiteNotFound:
			writeError(w, http.StatusNotFound, "site not found")
		case ErrSiteExists:
			writeError(w, http.StatusConflict, "site with this code already exists")
		case ErrSiteCodeUppercase:
			writeError(w, http.StatusBadRequest, "site code must be uppercase")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, site)
}

func (h *Handler) handleAdminDisableSite(w http.ResponseWriter, r *http.Request) {
	siteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid site id")
		return
	}

	dev := extractDeviceInfo(r)
	site, err := h.svc.AdminDisableSite(r.Context(), siteID, dev)
	if err != nil {
		switch err {
		case ErrSiteNotFound:
			writeError(w, http.StatusNotFound, "site not found")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, site)
}

// -----------------------------------------------------------------
// GET /admin/roles
// -----------------------------------------------------------------

func (h *Handler) handleAdminListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.AdminListRoles(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"roles": roles})
}

// -----------------------------------------------------------------
// POST /admin/roles
// -----------------------------------------------------------------

func (h *Handler) handleAdminCreateRole(w http.ResponseWriter, r *http.Request) {
	var req AdminCreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	dev := extractDeviceInfo(r)
	role, err := h.svc.AdminCreateRole(r.Context(), req, dev)
	if err != nil {
		switch err {
		case ErrRoleExists:
			writeError(w, http.StatusConflict, "role with this name already exists")
		case ErrUnknownPermissions:
			writeError(w, http.StatusBadRequest, "one or more permissions do not exist")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, role)
}

// -----------------------------------------------------------------
// GET /admin/roles/{id}
// -----------------------------------------------------------------

func (h *Handler) handleAdminGetRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid role id")
		return
	}

	role, err := h.svc.AdminGetRole(r.Context(), roleID)
	if err != nil {
		switch err {
		case ErrRoleNotFound:
			writeError(w, http.StatusNotFound, "role not found")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, role)
}

// -----------------------------------------------------------------
// PATCH /admin/roles/{id}
// -----------------------------------------------------------------

func (h *Handler) handleAdminUpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid role id")
		return
	}

	var req AdminUpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	role, err := h.svc.AdminUpdateRole(r.Context(), roleID, req, dev)
	if err != nil {
		switch err {
		case ErrRoleNotFound:
			writeError(w, http.StatusNotFound, "role not found")
		case ErrRoleExists:
			writeError(w, http.StatusConflict, "role with this name already exists")
		case ErrUnknownPermissions:
			writeError(w, http.StatusBadRequest, "one or more permissions do not exist")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, role)
}

// -----------------------------------------------------------------
// GET /admin/permissions
// -----------------------------------------------------------------

func (h *Handler) handleAdminListPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := h.svc.AdminListPermissions(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"permissions": perms})
}

// -----------------------------------------------------------------
// POST /patients
// -----------------------------------------------------------------

func (h *Handler) handleCreatePatient(w http.ResponseWriter, r *http.Request) {
	actorID, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req PatientCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FullName == "" {
		writeError(w, http.StatusBadRequest, "full_name is required")
		return
	}

	dev := extractDeviceInfo(r)
	patient, err := h.svc.CreatePatient(r.Context(), req, actorID, dev)
	if err != nil {
		switch err {
		case ErrInvalidPatientStatus:
			writeError(w, http.StatusBadRequest, "invalid status")
		case ErrDuplicateNationalID:
			writeError(w, http.StatusConflict, "patient with this national ID already exists")
		case ErrDuplicateMRN:
			writeError(w, http.StatusConflict, "patient with this MRN already exists")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, patient)
}

// -----------------------------------------------------------------
// GET /patients/{id}
// -----------------------------------------------------------------

func (h *Handler) handleGetPatient(w http.ResponseWriter, r *http.Request) {
	patientID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid patient id")
		return
	}

	patient, err := h.svc.GetPatient(r.Context(), patientID)
	if err != nil {
		switch err {
		case ErrPatientNotFound:
			writeError(w, http.StatusNotFound, "patient not found")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, patient)
}

// -----------------------------------------------------------------
// GET /patients
// -----------------------------------------------------------------

func (h *Handler) handleListPatients(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	req := ListPatientsRequest{
		Limit:  AdminUsersDefaultLimit,
		Cursor: strings.TrimSpace(q.Get("cursor")),
		Search: strings.TrimSpace(q.Get("search")),
	}

	if limitRaw := strings.TrimSpace(q.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if limit <= 0 {
			limit = AdminUsersDefaultLimit
		}
		if limit > AdminUsersMaxLimit {
			limit = AdminUsersMaxLimit
		}
		req.Limit = limit
	}

	if statusRaw := strings.TrimSpace(q.Get("status")); statusRaw != "" {
		status, err := normalizePatientStatus(&statusRaw, "")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid status")
			return
		}
		req.Status = status
	}

	if includeTotalRaw := strings.TrimSpace(q.Get("include_total")); includeTotalRaw != "" {
		v, ok := parseStrictBool(includeTotalRaw)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid include_total")
			return
		}
		req.IncludeTotal = v
	}

	result, err := h.svc.ListPatients(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			writeError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// -----------------------------------------------------------------
// PATCH /patients/{id}
// -----------------------------------------------------------------

func (h *Handler) handleUpdatePatient(w http.ResponseWriter, r *http.Request) {
	patientID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid patient id")
		return
	}

	var req PatientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dev := extractDeviceInfo(r)
	patient, err := h.svc.UpdatePatient(r.Context(), patientID, req, dev)
	if err != nil {
		switch err {
		case ErrPatientNotFound:
			writeError(w, http.StatusNotFound, "patient not found")
		case ErrInvalidPatientStatus:
			writeError(w, http.StatusBadRequest, "invalid status")
		case ErrDuplicateNationalID:
			writeError(w, http.StatusConflict, "patient with this national ID already exists")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, patient)
}

// -----------------------------------------------------------------
// POST /patients/{id}/disable
// -----------------------------------------------------------------

func (h *Handler) handleDisablePatient(w http.ResponseWriter, r *http.Request) {
	actorID, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	patientID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid patient id")
		return
	}

	dev := extractDeviceInfo(r)
	patient, err := h.svc.DisablePatient(r.Context(), patientID, actorID, dev)
	if err != nil {
		switch err {
		case ErrPatientNotFound:
			writeError(w, http.StatusNotFound, "patient not found")
		case ErrPatientDisableDeny:
			writeError(w, http.StatusForbidden, "only admin or the creator can disable this patient")
		default:
			writeInternalError(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, patient)
}

// -----------------------------------------------------------------
// Context helper — extracts user ID set by RequireAuth middleware.
// -----------------------------------------------------------------

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(CtxUserIDKey).(uuid.UUID)
	return id, ok
}

// CtxUserIDKey is the shared context key for the authenticated user ID.
// Exported so the middleware package can use the same key.
type ContextKey string

const CtxUserIDKey ContextKey = "user_id"

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

func extractDeviceInfo(r *http.Request) DeviceInfo {
	return DeviceInfo{
		DeviceName: r.Header.Get("X-Device-Name"),
		IP:         realIP(r),
		UserAgent:  r.UserAgent(),
	}
}

func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first (leftmost) IP — the client.
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeInternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error: method=%s path=%s err=%v", r.Method, r.URL.Path, err)
	writeError(w, http.StatusInternalServerError, "internal error")
}
