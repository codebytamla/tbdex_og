package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrUserNotActive        = errors.New("user account is deactivated")
	ErrTokenRevoked         = errors.New("refresh token has been revoked")
	ErrTokenExpired         = errors.New("token has expired")
	ErrUserExists           = errors.New("user with this email already exists")
	ErrSiteExists           = errors.New("site with this code already exists")
	ErrSiteCodeUppercase    = errors.New("site code must be uppercase")
	ErrSiteNotFound         = errors.New("site not found")
	ErrInvalidResetToken    = errors.New("invalid or expired reset token")
	ErrRoleExists           = errors.New("role with this name already exists")
	ErrRoleNotFound         = errors.New("role not found")
	ErrUnknownPermissions   = errors.New("one or more permissions do not exist")
	ErrAdminRoleRequired    = errors.New("admin user must keep ADMIN role")
	ErrInvalidCursor        = errors.New("invalid cursor")
	ErrPatientNotFound      = errors.New("patient not found")
	ErrPatientDisableDeny   = errors.New("only admin or the creator can disable a patient")
	ErrInvalidPatientStatus = errors.New("invalid patient status")
	ErrDuplicateNationalID  = errors.New("patient with this national ID already exists")
	ErrDuplicateMRN         = errors.New("patient with this MRN already exists")
)

const (
	AdminUsersDefaultLimit   = 50
	AdminUsersMaxLimit       = 100
	adminRoleName            = "ADMIN"
	patientStatusProvisional = "provisional"
	patientStatusActive      = "active"
	patientStatusInactive    = "inactive"
	patientStatusDeceased    = "deceased"
)

// -----------------------------------------------------------------
// Store is the interface to the database layer. It abstracts the
// sqlc-generated queries so the service is testable without pgx.
// -----------------------------------------------------------------

type Store interface {
	// Users
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (UserWithPassword, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	UpdateUserPassword(ctx context.Context, id uuid.UUID, hash string) error
	UpdateUser(ctx context.Context, id uuid.UUID, params UpdateUserParams) (User, error)
	SetUserActive(ctx context.Context, id uuid.UUID, active bool) (User, error)
	ListUsers(ctx context.Context, p AdminListUsersRequest) (UserPage, error)
	CountUsers(ctx context.Context, p AdminListUsersRequest) (int, error)
	ListUserRolesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]Role, error)
	ListUserSitesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]SiteSummary, error)

	// Refresh tokens
	CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) (RefreshTokenID, error)
	GetRefreshTokenByHash(ctx context.Context, hash string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, params CreatePasswordResetTokenParams) error
	GetPasswordResetTokenByHash(ctx context.Context, hash string) (PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error
	ExpireOldPasswordResetTokens(ctx context.Context, userID uuid.UUID) error

	// RBAC
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	GetRoleByName(ctx context.Context, name string) (Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (Role, error)
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
	RemoveAllUserRoles(ctx context.Context, userID uuid.UUID) error

	// Role management
	ListRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, name, description string) (Role, error)
	UpdateRole(ctx context.Context, id uuid.UUID, name, description *string) (Role, error)
	ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, permissionNames []string) ([]Permission, error)
	ListAllPermissions(ctx context.Context) ([]Permission, error)

	// Sites
	AssignUserToSite(ctx context.Context, userID, siteID uuid.UUID) error
	RemoveUserFromSite(ctx context.Context, userID, siteID uuid.UUID) error
	ListUserSites(ctx context.Context, userID uuid.UUID) ([]SiteSummary, error)
	RemoveAllUserSites(ctx context.Context, userID uuid.UUID) error
	GetSiteSummaryByID(ctx context.Context, id uuid.UUID) (SiteSummary, error)

	// Sites
	CreateSite(ctx context.Context, name, code, siteType string) (Site, error)
	ListSites(ctx context.Context) ([]Site, error)
	GetSiteByID(ctx context.Context, id uuid.UUID) (Site, error)
	UpdateSite(ctx context.Context, id uuid.UUID, name, code, siteType string) (Site, error)
	DisableSite(ctx context.Context, id uuid.UUID) (Site, error)

	// Patients
	CreatePatient(ctx context.Context, params CreatePatientParams) (Patient, error)
	GetPatientByID(ctx context.Context, id uuid.UUID) (Patient, error)
	ListPatients(ctx context.Context, p ListPatientsRequest) (PatientPage, error)
	CountPatients(ctx context.Context, p ListPatientsRequest) (int, error)
	UpdatePatient(ctx context.Context, id uuid.UUID, params UpdatePatientParams) (Patient, error)
	DisablePatient(ctx context.Context, id uuid.UUID) (Patient, error)
	GetPatientCreatedBy(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UserIsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)

	// Transactions — for refresh rotation
	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) Store
}

// Minimal domain types so the service doesn't import db package directly.
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserWithPassword struct {
	User
	PasswordHash string `json:"-"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

type RefreshTokenID struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

type Permission struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RoleWithPermissions struct {
	Role
	Permissions []Permission `json:"permissions"`
}

type SiteSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

type Site struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	SiteType  string    `json:"site_type"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Patient struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"full_name"`
	DateOfBirth    *string   `json:"date_of_birth,omitempty"` // "YYYY-MM-DD"
	IsDobEstimated bool      `json:"is_dob_estimated"`
	Sex            *string   `json:"sex,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	NationalID     *string   `json:"national_id,omitempty"`
	MRN            *string   `json:"mrn,omitempty"`
	Address        *string   `json:"address,omitempty"`
	Village        *string   `json:"village,omitempty"`
	Origin         *string   `json:"origin,omitempty"`
	Status         string    `json:"status"`
	CreatedBy      uuid.UUID `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreatePatientParams struct {
	FullName       string
	DateOfBirth    *string
	IsDobEstimated bool
	Sex            *string
	Phone          *string
	NationalID     *string
	MRN            *string
	Address        *string
	Village        *string
	Origin         *string
	Status         string
	CreatedBy      uuid.UUID
}

type UpdatePatientParams struct {
	FullName       *string // nil = don't change
	DateOfBirth    *string
	IsDobEstimated *bool
	Sex            *string
	Phone          *string
	NationalID     *string
	Address        *string
	Village        *string
	Origin         *string
	Status         *string
}

type ListPatientsRequest struct {
	Limit        int
	Cursor       string
	Search       string
	Status       *string
	IncludeTotal bool
}

type PatientPage struct {
	Patients   []Patient `json:"patients"`
	HasNext    bool      `json:"has_next"`
	NextCursor *string   `json:"next_cursor,omitempty"`
	Total      *int      `json:"total,omitempty"`
}

type UpdateUserParams struct {
	FullName *string // nil = don't change
}

type UserProfile struct {
	User
	Roles       []Role        `json:"roles"`
	Permissions []string      `json:"permissions"`
	Sites       []SiteSummary `json:"sites"`
}

type CreateUserParams struct {
	Email        string
	PasswordHash string
	FullName     string
}

type CreateRefreshTokenParams struct {
	UserID     uuid.UUID
	TokenHash  string
	DeviceName string
	CreatedIP  string
	UserAgent  string
	ExpiresAt  time.Time
}

type CreatePasswordResetTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

type AdminListUsersRequest struct {
	Limit        int
	Cursor       string
	Search       string
	IsActive     *bool
	Roles        []string
	SiteIDs      []uuid.UUID
	IncludeTotal bool
}

type UserPage struct {
	Users      []User
	HasNext    bool
	NextCursor *string
}

// -----------------------------------------------------------------
// Service is the core authentication service.
// -----------------------------------------------------------------

type Service struct {
	store Store
	jwt   JWTConfig
	audit AuditWriter
}

func NewService(store Store, jwtCfg JWTConfig, audit AuditWriter) *Service {
	return &Service{store: store, jwt: jwtCfg, audit: audit}
}

// -----------------------------------------------------------------
// Request/Response types
// -----------------------------------------------------------------

type DeviceInfo struct {
	DeviceName string
	IP         string
	UserAgent  string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// -----------------------------------------------------------------
// Login
// -----------------------------------------------------------------

func (s *Service) Login(ctx context.Context, req LoginRequest, dev DeviceInfo) (*User, *TokenPair, error) {
	userWP, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		LogAudit(ctx, s.audit, AuditLoginFail, nil, nil, dev.IP, dev.UserAgent,
			map[string]any{"email": req.Email, "reason": "not_found"})
		return nil, nil, ErrInvalidCredentials
	}

	if !userWP.IsActive {
		LogAudit(ctx, s.audit, AuditLoginFail, &userWP.ID, &userWP.ID, dev.IP, dev.UserAgent,
			map[string]any{"reason": "deactivated"})
		return nil, nil, ErrUserNotActive
	}

	ok, err := VerifyPassword(req.Password, userWP.PasswordHash)
	if err != nil || !ok {
		LogAudit(ctx, s.audit, AuditLoginFail, &userWP.ID, &userWP.ID, dev.IP, dev.UserAgent,
			map[string]any{"reason": "bad_password"})
		return nil, nil, ErrInvalidCredentials
	}

	LogAudit(ctx, s.audit, AuditLoginSuccess, &userWP.ID, &userWP.ID, dev.IP, dev.UserAgent, nil)

	tokens, err := s.issueTokens(ctx, userWP.ID, dev)
	if err != nil {
		return nil, nil, err
	}

	user := userWP.User
	return &user, tokens, nil
}

// -----------------------------------------------------------------
// Refresh — atomic rotation
// -----------------------------------------------------------------

func (s *Service) Refresh(ctx context.Context, rawRefreshToken string, dev DeviceInfo) (*TokenPair, error) {
	tokenHash := HashToken(rawRefreshToken)

	// Begin transaction for atomic rotation.
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	txStore := s.store.WithTx(tx)

	existing, err := txStore.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrTokenRevoked
	}

	if existing.Revoked {
		// Possible token reuse attack — revoke ALL tokens for this user.
		_ = txStore.RevokeRefreshToken(ctx, existing.ID)
		LogAudit(ctx, s.audit, AuditRefresh, &existing.UserID, &existing.UserID, dev.IP, dev.UserAgent,
			map[string]any{"reason": "reuse_detected"})
		return nil, ErrTokenRevoked
	}

	if time.Now().After(existing.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Revoke old token.
	if err := txStore.RevokeRefreshToken(ctx, existing.ID); err != nil {
		return nil, err
	}

	// Issue new refresh token.
	newRaw, err := GenerateOpaqueToken(32)
	if err != nil {
		return nil, err
	}
	newHash := HashToken(newRaw)
	if _, err := txStore.CreateRefreshToken(ctx, CreateRefreshTokenParams{
		UserID:     existing.UserID,
		TokenHash:  newHash,
		DeviceName: dev.DeviceName,
		CreatedIP:  dev.IP,
		UserAgent:  dev.UserAgent,
		ExpiresAt:  time.Now().Add(s.jwt.RefreshTTL),
	}); err != nil {
		return nil, err
	}

	// Issue new access token.
	accessToken, err := IssueAccessToken(s.jwt, existing.UserID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, AuditRefresh, &existing.UserID, &existing.UserID, dev.IP, dev.UserAgent, nil)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRaw,
		ExpiresIn:    int(s.jwt.AccessTTL.Seconds()),
	}, nil
}

// -----------------------------------------------------------------
// Logout — revoke the current refresh token
// -----------------------------------------------------------------

func (s *Service) Logout(ctx context.Context, rawRefreshToken string, dev DeviceInfo) error {
	tokenHash := HashToken(rawRefreshToken)
	existing, err := s.store.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil // already revoked or doesn't exist — idempotent
	}

	if err := s.store.RevokeRefreshToken(ctx, existing.ID); err != nil {
		return err
	}

	LogAudit(ctx, s.audit, AuditLogout, &existing.UserID, &existing.UserID, dev.IP, dev.UserAgent, nil)
	return nil
}

// -----------------------------------------------------------------
// Forgot Password
// -----------------------------------------------------------------

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest, dev DeviceInfo) (string, error) {
	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal whether email exists.
		return "", nil
	}

	// Expire any old unused reset tokens for this user.
	_ = s.store.ExpireOldPasswordResetTokens(ctx, user.ID)

	raw, err := GenerateOpaqueToken(32)
	if err != nil {
		return "", err
	}
	tokenHash := HashToken(raw)

	if err := s.store.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}); err != nil {
		return "", err
	}

	LogAudit(ctx, s.audit, AuditPasswordResetReq, &user.ID, &user.ID, dev.IP, dev.UserAgent, nil)

	// Return the raw token. In production, this would be emailed.
	return raw, nil
}

// -----------------------------------------------------------------
// Reset Password
// -----------------------------------------------------------------

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest, dev DeviceInfo) error {
	if err := ValidatePasswordPolicy(req.NewPassword); err != nil {
		return err
	}

	tokenHash := HashToken(req.Token)
	resetToken, err := s.store.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidResetToken
	}

	if resetToken.Used || time.Now().After(resetToken.ExpiresAt) {
		return ErrInvalidResetToken
	}

	hash, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	if err := txStore.UpdateUserPassword(ctx, resetToken.UserID, hash); err != nil {
		return err
	}
	if err := txStore.MarkPasswordResetTokenUsed(ctx, resetToken.ID); err != nil {
		return err
	}
	if err := txStore.RevokeAllUserRefreshTokens(ctx, resetToken.UserID); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	LogAudit(ctx, s.audit, AuditPasswordResetDone, &resetToken.UserID, &resetToken.UserID, dev.IP, dev.UserAgent, nil)

	return nil
}

// -----------------------------------------------------------------
// GetMe — returns full user profile with roles/permissions/sites.
// -----------------------------------------------------------------

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.store.ListUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	perms, err := s.store.GetEffectivePermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	sites, err := s.store.ListUserSites(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		User:        user,
		Roles:       roles,
		Permissions: perms,
		Sites:       sites,
	}, nil
}

// -----------------------------------------------------------------
// AdminCreateUser — admin creates a user with optional roles/sites.
// -----------------------------------------------------------------

type AdminCreateUserRequest struct {
	Email    string   `json:"email"`
	FullName string   `json:"full_name"`
	Password string   `json:"password"` // if empty, user must reset
	Roles    []string `json:"roles"`    // role names
	Sites    []string `json:"sites"`    // site IDs as strings
}

func (s *Service) AdminCreateUser(ctx context.Context, req AdminCreateUserRequest, dev DeviceInfo) (*User, string, error) {
	// If no password provided, generate a random one (user must reset).
	password := req.Password
	resetToken := ""
	if password == "" {
		generated, err := GenerateOpaqueToken(16)
		if err != nil {
			return nil, "", err
		}
		password = generated
	}

	if err := ValidatePasswordPolicy(password); err != nil {
		return nil, "", err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user, err := s.store.CreateUser(ctx, CreateUserParams{
		Email:        req.Email,
		PasswordHash: hash,
		FullName:     req.FullName,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, "", ErrUserExists
		}
		return nil, "", err
	}

	// Assign roles.
	for _, roleName := range req.Roles {
		role, err := s.store.GetRoleByName(ctx, roleName)
		if err != nil {
			continue // skip unknown roles
		}
		_ = s.store.AssignRoleToUser(ctx, user.ID, role.ID)
	}

	// Assign sites.
	for _, siteIDStr := range req.Sites {
		siteID, err := uuid.Parse(siteIDStr)
		if err != nil {
			continue
		}
		_ = s.store.AssignUserToSite(ctx, user.ID, siteID)
	}

	// If no password was provided, generate a reset token so admin can share it.
	if req.Password == "" {
		resetToken, err = s.generateResetToken(ctx, user.ID)
		if err != nil {
			return &user, "", err
		}
	}

	LogAudit(ctx, s.audit, AuditRegister, nil, &user.ID, dev.IP, dev.UserAgent,
		map[string]any{"admin_created": true})

	return &user, resetToken, nil
}

// -----------------------------------------------------------------
// AdminListUsers
// -----------------------------------------------------------------

type UserListResponse struct {
	Users      []UserListItem     `json:"users"`
	Pagination UserListPagination `json:"pagination"`
	Total      *int               `json:"total,omitempty"`
}

type UserListItem struct {
	User
	Roles []Role        `json:"roles"`
	Sites []SiteSummary `json:"sites"`
}

type UserListPagination struct {
	NextCursor *string `json:"next_cursor"`
	HasNext    bool    `json:"has_next"`
	Limit      int     `json:"limit"`
}

func normalizeAdminListUsersRequest(req AdminListUsersRequest) AdminListUsersRequest {
	if req.Limit <= 0 {
		req.Limit = AdminUsersDefaultLimit
	}
	if req.Limit > AdminUsersMaxLimit {
		req.Limit = AdminUsersMaxLimit
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Roles = normalizeNonEmptyStrings(req.Roles)
	req.SiteIDs = uniqueUUIDs(req.SiteIDs)

	return req
}

func normalizeNonEmptyStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, raw := range values {
		v := strings.TrimSpace(raw)
		if v == "" {
			continue
		}
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func (s *Service) validateRolePermissionNames(ctx context.Context, permissionNames []string) ([]string, error) {
	normalized := normalizeNonEmptyStrings(permissionNames)
	if len(normalized) == 0 {
		return []string{}, nil
	}

	catalog, err := s.store.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	known := make(map[string]struct{}, len(catalog))
	for _, permission := range catalog {
		name := strings.TrimSpace(permission.Name)
		if name == "" {
			continue
		}
		known[name] = struct{}{}
	}

	for _, name := range normalized {
		if _, ok := known[name]; !ok {
			return nil, ErrUnknownPermissions
		}
	}

	return normalized, nil
}

func uniqueUUIDs(values []uuid.UUID) []uuid.UUID {
	if len(values) == 0 {
		return []uuid.UUID{}
	}

	seen := make(map[uuid.UUID]struct{}, len(values))
	result := make([]uuid.UUID, 0, len(values))
	for _, v := range values {
		if v == uuid.Nil {
			continue
		}
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func normalizePatientStatus(input *string, fallback string) (*string, error) {
	if input == nil {
		if fallback == "" {
			return nil, nil
		}
		v := fallback
		return &v, nil
	}

	status := strings.ToLower(strings.TrimSpace(*input))
	switch status {
	case patientStatusProvisional, patientStatusActive, patientStatusInactive, patientStatusDeceased:
		return &status, nil
	default:
		return nil, ErrInvalidPatientStatus
	}
}

func hasRoleName(roles []Role, name string) bool {
	for _, role := range roles {
		if strings.TrimSpace(role.Name) == name {
			return true
		}
	}
	return false
}

func (s *Service) findRoleByName(ctx context.Context, roleName string) (Role, error) {
	name := strings.TrimSpace(roleName)
	role, err := s.store.GetRoleByName(ctx, name)
	if err == nil {
		return role, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return Role{}, err
	}

	upper := strings.ToUpper(name)
	if upper == name {
		return Role{}, err
	}
	return s.store.GetRoleByName(ctx, upper)
}

func (s *Service) AdminListUsers(ctx context.Context, req AdminListUsersRequest) (*UserListResponse, error) {
	req = normalizeAdminListUsersRequest(req)

	page, err := s.store.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	var total *int
	if req.IncludeTotal {
		count, err := s.store.CountUsers(ctx, req)
		if err != nil {
			return nil, err
		}
		total = &count
	}

	response := &UserListResponse{
		Users: []UserListItem{},
		Pagination: UserListPagination{
			NextCursor: page.NextCursor,
			HasNext:    page.HasNext,
			Limit:      req.Limit,
		},
		Total: total,
	}

	if len(page.Users) == 0 {
		return response, nil
	}

	userIDs := make([]uuid.UUID, len(page.Users))
	for i, u := range page.Users {
		userIDs[i] = u.ID
	}

	rolesByUser, err := s.store.ListUserRolesByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	sitesByUser, err := s.store.ListUserSitesByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	items := make([]UserListItem, len(page.Users))
	for i, u := range page.Users {
		roles := rolesByUser[u.ID]
		if roles == nil {
			roles = []Role{}
		}

		sites := sitesByUser[u.ID]
		if sites == nil {
			sites = []SiteSummary{}
		}

		items[i] = UserListItem{
			User:  u,
			Roles: roles,
			Sites: sites,
		}
	}

	response.Users = items
	return response, nil
}

// -----------------------------------------------------------------
// AdminUpdateUser — update profile fields only (full_name)
// -----------------------------------------------------------------

type AdminUpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty"`
}

func (s *Service) AdminUpdateUser(ctx context.Context, userID uuid.UUID, req AdminUpdateUserRequest, dev DeviceInfo) (*User, error) {
	user, err := s.store.UpdateUser(ctx, userID, UpdateUserParams{
		FullName: req.FullName,
	})
	if err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "user_updated", nil, &userID, dev.IP, dev.UserAgent,
		map[string]any{"full_name": req.FullName})

	return &user, nil
}

// -----------------------------------------------------------------
// Site management (Admin)
// -----------------------------------------------------------------

type AdminCreateSiteRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	SiteType string `json:"site_type"`
}

type AdminUpdateSiteRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	SiteType string `json:"site_type"`
}

func (s *Service) AdminCreateSite(ctx context.Context, req AdminCreateSiteRequest, dev DeviceInfo) (*Site, error) {
	siteType := req.SiteType
	if siteType == "" {
		siteType = "site"
	}

	site, err := s.store.CreateSite(ctx, req.Name, req.Code, siteType)
	if err != nil {
		if isCheckConstraintError(err, "chk_sites_code_uppercase") {
			return nil, ErrSiteCodeUppercase
		}
		if isDuplicateKeyError(err) {
			return nil, ErrSiteExists
		}
		return nil, err
	}

	LogAudit(ctx, s.audit, "site_created", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"site_id": site.ID, "code": site.Code, "site_type": site.SiteType})

	return &site, nil
}

func (s *Service) AdminListSites(ctx context.Context) ([]Site, error) {
	return s.store.ListSites(ctx)
}

func (s *Service) AdminGetSite(ctx context.Context, siteID uuid.UUID) (*Site, error) {
	site, err := s.store.GetSiteByID(ctx, siteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}
	return &site, nil
}

func (s *Service) AdminUpdateSite(ctx context.Context, siteID uuid.UUID, req AdminUpdateSiteRequest, dev DeviceInfo) (*Site, error) {
	siteType := req.SiteType
	if siteType == "" {
		siteType = "site"
	}

	site, err := s.store.UpdateSite(ctx, siteID, req.Name, req.Code, siteType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSiteNotFound
		}
		if isCheckConstraintError(err, "chk_sites_code_uppercase") {
			return nil, ErrSiteCodeUppercase
		}
		if isDuplicateKeyError(err) {
			return nil, ErrSiteExists
		}
		return nil, err
	}

	LogAudit(ctx, s.audit, "site_updated", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"site_id": site.ID, "code": site.Code, "site_type": site.SiteType})

	return &site, nil
}

func (s *Service) AdminDisableSite(ctx context.Context, siteID uuid.UUID, dev DeviceInfo) (*Site, error) {
	site, err := s.store.DisableSite(ctx, siteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	LogAudit(ctx, s.audit, "site_disabled", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"site_id": site.ID, "code": site.Code})

	return &site, nil
}

// -----------------------------------------------------------------
// AdminGetUserSitesAndRoles — fetch assignments for a single user
// -----------------------------------------------------------------

type UserSitesAndRolesResponse struct {
	Roles []Role        `json:"roles"`
	Sites []SiteSummary `json:"sites"`
}

func (s *Service) AdminGetUserSitesAndRoles(ctx context.Context, userID uuid.UUID) (*UserSitesAndRolesResponse, error) {
	roles, err := s.store.ListUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	if roles == nil {
		roles = []Role{}
	}

	sites, err := s.store.ListUserSites(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sites == nil {
		sites = []SiteSummary{}
	}

	return &UserSitesAndRolesResponse{
		Roles: roles,
		Sites: sites,
	}, nil
}

// -----------------------------------------------------------------
// AdminSetRoles — replace all roles for a user
// -----------------------------------------------------------------

type SetRolesRequest struct {
	Roles []string `json:"roles"` // role names
}

func (s *Service) AdminSetRoles(ctx context.Context, userID uuid.UUID, roleNames []string, dev DeviceInfo) ([]Role, error) {
	currentRoles, err := s.store.ListUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	resolvedRoles := make([]Role, 0, len(roleNames))
	seenRoleIDs := make(map[uuid.UUID]struct{}, len(roleNames))
	for _, rawName := range roleNames {
		name := strings.TrimSpace(rawName)
		if name == "" {
			continue
		}

		role, err := s.findRoleByName(ctx, name)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue // skip unknown
			}
			return nil, err
		}
		if _, exists := seenRoleIDs[role.ID]; exists {
			continue
		}
		seenRoleIDs[role.ID] = struct{}{}
		resolvedRoles = append(resolvedRoles, role)
	}

	if hasRoleName(currentRoles, adminRoleName) && !hasRoleName(resolvedRoles, adminRoleName) {
		return nil, ErrAdminRoleRequired
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	// Remove all existing roles.
	if err := txStore.RemoveAllUserRoles(ctx, userID); err != nil {
		return nil, err
	}

	// Assign new roles.
	assigned := make([]Role, 0, len(resolvedRoles))
	for _, role := range resolvedRoles {
		if err := txStore.AssignRoleToUser(ctx, userID, role.ID); err != nil {
			return nil, err
		}
		assigned = append(assigned, role)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, AuditRoleChange, nil, &userID, dev.IP, dev.UserAgent,
		map[string]any{"roles": roleNames})

	return assigned, nil
}

// -----------------------------------------------------------------
// AdminSetSites — replace all site assignments for a user.
// -----------------------------------------------------------------

type SetSitesRequest struct {
	Sites []string `json:"sites"` // site IDs
}

func (s *Service) AdminSetSites(ctx context.Context, userID uuid.UUID, siteIDs []string, dev DeviceInfo) ([]SiteSummary, error) {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	// Remove all existing assignments.
	if err := txStore.RemoveAllUserSites(ctx, userID); err != nil {
		return nil, err
	}

	var assigned []SiteSummary
	for _, idStr := range siteIDs {
		siteID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		if err := txStore.AssignUserToSite(ctx, userID, siteID); err != nil {
			return nil, err
		}
		site, err := txStore.GetSiteSummaryByID(ctx, siteID)
		if err != nil {
			return nil, err
		}
		assigned = append(assigned, site)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "sites_changed", nil, &userID, dev.IP, dev.UserAgent,
		map[string]any{"sites": siteIDs})

	return assigned, nil
}

// -----------------------------------------------------------------
// AdminForcePasswordReset — generate reset token for a target user
// -----------------------------------------------------------------

func (s *Service) AdminForcePasswordReset(ctx context.Context, targetUserID uuid.UUID, dev DeviceInfo) (string, error) {
	// Verify user exists.
	_, err := s.store.GetUserByID(ctx, targetUserID)
	if err != nil {
		return "", err
	}

	token, err := s.generateResetToken(ctx, targetUserID)
	if err != nil {
		return "", err
	}

	LogAudit(ctx, s.audit, AuditPasswordResetReq, nil, &targetUserID, dev.IP, dev.UserAgent,
		map[string]any{"admin_forced": true})

	return token, nil
}

// -----------------------------------------------------------------
// AdminDisableUser — set is_active=false and revoke all sessions
// -----------------------------------------------------------------

func (s *Service) AdminDisableUser(ctx context.Context, targetUserID uuid.UUID, dev DeviceInfo) (*User, error) {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	user, err := txStore.SetUserActive(ctx, targetUserID, false)
	if err != nil {
		return nil, err
	}

	if err := txStore.RevokeAllUserRefreshTokens(ctx, targetUserID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "user_disabled", nil, &targetUserID, dev.IP, dev.UserAgent, nil)

	return &user, nil
}

// -----------------------------------------------------------------
// AdminEnableUser — set is_active=true
// -----------------------------------------------------------------

func (s *Service) AdminEnableUser(ctx context.Context, targetUserID uuid.UUID, dev DeviceInfo) (*User, error) {
	user, err := s.store.SetUserActive(ctx, targetUserID, true)
	if err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "user_enabled", nil, &targetUserID, dev.IP, dev.UserAgent, nil)

	return &user, nil
}

// -----------------------------------------------------------------
// AdminListRoles — list all roles
// -----------------------------------------------------------------

func (s *Service) AdminListRoles(ctx context.Context) ([]Role, error) {
	return s.store.ListRoles(ctx)
}

// -----------------------------------------------------------------
// AdminGetRole — get a single role with its permissions
// -----------------------------------------------------------------

func (s *Service) AdminGetRole(ctx context.Context, roleID uuid.UUID) (*RoleWithPermissions, error) {
	role, err := s.store.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	perms, err := s.store.ListRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &RoleWithPermissions{Role: role, Permissions: perms}, nil
}

// -----------------------------------------------------------------
// AdminCreateRole — create a new role with optional permissions
// -----------------------------------------------------------------

type AdminCreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"` // permission names
}

func (s *Service) AdminCreateRole(ctx context.Context, req AdminCreateRoleRequest, dev DeviceInfo) (*RoleWithPermissions, error) {
	normalizedPermissions, err := s.validateRolePermissionNames(ctx, req.Permissions)
	if err != nil {
		return nil, err
	}

	role, err := s.store.CreateRole(ctx, req.Name, req.Description)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrRoleExists
		}
		return nil, err
	}

	var perms []Permission
	if len(normalizedPermissions) > 0 {
		perms, err = s.store.ReplaceRolePermissions(ctx, role.ID, normalizedPermissions)
		if err != nil {
			return nil, err
		}
	}
	if perms == nil {
		perms = []Permission{}
	}

	LogAudit(ctx, s.audit, "role_created", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"role": role.Name, "permissions": normalizedPermissions})

	return &RoleWithPermissions{Role: role, Permissions: perms}, nil
}

// -----------------------------------------------------------------
// AdminUpdateRole — update name, description, and/or permissions
// -----------------------------------------------------------------

type AdminUpdateRoleRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"` // nil => unchanged, empty => clear all
}

func (s *Service) AdminUpdateRole(ctx context.Context, roleID uuid.UUID, req AdminUpdateRoleRequest, dev DeviceInfo) (*RoleWithPermissions, error) {
	// Verify role exists.
	_, err := s.store.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	role, err := s.store.UpdateRole(ctx, roleID, req.Name, req.Description)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrRoleExists
		}
		return nil, err
	}

	var perms []Permission
	if req.Permissions != nil {
		normalizedPermissions, err := s.validateRolePermissionNames(ctx, *req.Permissions)
		if err != nil {
			return nil, err
		}
		perms, err = s.store.ReplaceRolePermissions(ctx, role.ID, normalizedPermissions)
		if err != nil {
			return nil, err
		}
	} else {
		perms, err = s.store.ListRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}
	}

	LogAudit(ctx, s.audit, "role_updated", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"role_id": roleID, "name": req.Name, "permissions": req.Permissions})

	return &RoleWithPermissions{Role: role, Permissions: perms}, nil
}

// -----------------------------------------------------------------
// AdminListPermissions — list all available permissions
// -----------------------------------------------------------------

func (s *Service) AdminListPermissions(ctx context.Context) ([]Permission, error) {
	return s.store.ListAllPermissions(ctx)
}

// -----------------------------------------------------------------
// Patient endpoints
// -----------------------------------------------------------------

type PatientCreateRequest struct {
	FullName       string  `json:"full_name"`
	DateOfBirth    *string `json:"date_of_birth,omitempty"`
	IsDobEstimated bool    `json:"is_dob_estimated"`
	Sex            *string `json:"sex,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	NationalID     *string `json:"national_id,omitempty"`
	MRN            *string `json:"mrn,omitempty"`
	Address        *string `json:"address,omitempty"`
	Village        *string `json:"village,omitempty"`
	Origin         *string `json:"origin,omitempty"`
	Status         *string `json:"status,omitempty"`
}

func (s *Service) CreatePatient(ctx context.Context, req PatientCreateRequest, actorID uuid.UUID, dev DeviceInfo) (*Patient, error) {
	normalizedStatus, err := normalizePatientStatus(req.Status, patientStatusActive)
	if err != nil {
		return nil, err
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	patient, err := txStore.CreatePatient(ctx, CreatePatientParams{
		FullName:       req.FullName,
		DateOfBirth:    req.DateOfBirth,
		IsDobEstimated: req.IsDobEstimated,
		Sex:            req.Sex,
		Phone:          req.Phone,
		NationalID:     req.NationalID,
		MRN:            req.MRN,
		Address:        req.Address,
		Village:        req.Village,
		Origin:         req.Origin,
		Status:         *normalizedStatus,
		CreatedBy:      actorID,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicateNationalID
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "patient_created", &actorID, nil, dev.IP, dev.UserAgent,
		map[string]any{"patient_id": patient.ID})

	return &patient, nil
}

func (s *Service) GetPatient(ctx context.Context, patientID uuid.UUID) (*Patient, error) {
	patient, err := s.store.GetPatientByID(ctx, patientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	return &patient, nil
}

func (s *Service) ListPatients(ctx context.Context, req ListPatientsRequest) (*PatientPage, error) {
	normalizedStatus, err := normalizePatientStatus(req.Status, "")
	if err != nil {
		return nil, err
	}
	req.Status = normalizedStatus

	page, err := s.store.ListPatients(ctx, req)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			return nil, ErrInvalidCursor
		}
		return nil, err
	}

	if req.IncludeTotal {
		count, err := s.store.CountPatients(ctx, req)
		if err != nil {
			return nil, err
		}
		page.Total = &count
	}

	return &page, nil
}

type PatientUpdateRequest struct {
	FullName       *string `json:"full_name,omitempty"`
	DateOfBirth    *string `json:"date_of_birth,omitempty"`
	IsDobEstimated *bool   `json:"is_dob_estimated,omitempty"`
	Sex            *string `json:"sex,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	NationalID     *string `json:"national_id,omitempty"`
	Address        *string `json:"address,omitempty"`
	Village        *string `json:"village,omitempty"`
	Origin         *string `json:"origin,omitempty"`
	Status         *string `json:"status,omitempty"`
}

func (s *Service) UpdatePatient(ctx context.Context, patientID uuid.UUID, req PatientUpdateRequest, dev DeviceInfo) (*Patient, error) {
	// Verify patient exists.
	if _, err := s.store.GetPatientByID(ctx, patientID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	status, err := normalizePatientStatus(req.Status, "")
	if err != nil {
		return nil, err
	}

	patient, err := txStore.UpdatePatient(ctx, patientID, UpdatePatientParams{
		FullName:       req.FullName,
		DateOfBirth:    req.DateOfBirth,
		IsDobEstimated: req.IsDobEstimated,
		Sex:            req.Sex,
		Phone:          req.Phone,
		NationalID:     req.NationalID,
		Address:        req.Address,
		Village:        req.Village,
		Origin:         req.Origin,
		Status:         status,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicateNationalID
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "patient_updated", nil, nil, dev.IP, dev.UserAgent,
		map[string]any{"patient_id": patientID})

	return &patient, nil
}

func (s *Service) DisablePatient(ctx context.Context, patientID, actorID uuid.UUID, dev DeviceInfo) (*Patient, error) {
	// Check if actor is admin OR the creator.
	isAdmin, err := s.store.UserIsAdmin(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		createdBy, err := s.store.GetPatientCreatedBy(ctx, patientID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrPatientNotFound
			}
			return nil, err
		}
		if createdBy != actorID {
			return nil, ErrPatientDisableDeny
		}
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	patient, err := txStore.DisablePatient(ctx, patientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	LogAudit(ctx, s.audit, "patient_disabled", &actorID, nil, dev.IP, dev.UserAgent,
		map[string]any{"patient_id": patientID})

	return &patient, nil
}

func (s *Service) generateResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	_ = s.store.ExpireOldPasswordResetTokens(ctx, userID)

	raw, err := GenerateOpaqueToken(32)
	if err != nil {
		return "", err
	}

	if err := s.store.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    userID,
		TokenHash: HashToken(raw),
		ExpiresAt: time.Now().Add(24 * time.Hour), // longer for admin-issued
	}); err != nil {
		return "", err
	}

	return raw, nil
}

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

func (s *Service) issueTokens(ctx context.Context, userID uuid.UUID, dev DeviceInfo) (*TokenPair, error) {
	accessToken, err := IssueAccessToken(s.jwt, userID)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := GenerateOpaqueToken(32)
	if err != nil {
		return nil, err
	}
	refreshHash := HashToken(rawRefresh)

	if _, err := s.store.CreateRefreshToken(ctx, CreateRefreshTokenParams{
		UserID:     userID,
		TokenHash:  refreshHash,
		DeviceName: dev.DeviceName,
		CreatedIP:  dev.IP,
		UserAgent:  dev.UserAgent,
		ExpiresAt:  time.Now().Add(s.jwt.RefreshTTL),
	}); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int(s.jwt.AccessTTL.Seconds()),
	}, nil
}

// isDuplicateKeyError checks for Postgres unique violation (23505).
func isDuplicateKeyError(err error) bool {
	// pgx wraps errors; check the code string.
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}

func isCheckConstraintError(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23514" && pgErr.ConstraintName == constraint
}

// BootstrapAdmin creates the first admin user if no users exist.
// Call this on server startup. If users already exist, it's a no-op.
func (s *Service) BootstrapAdmin(ctx context.Context, email, password, fullName string) error {
	// Check if any users exist (try to look up the email).
	_, err := s.store.GetUserByEmail(ctx, email)
	if err == nil {
		return nil // already exists
	}

	if err := ValidatePasswordPolicy(password); err != nil {
		return err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txStore := s.store.WithTx(tx)

	user, err := txStore.CreateUser(ctx, CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		FullName:     fullName,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil // race condition, fine
		}
		return err
	}

	// Assign ADMIN role.
	role, err := txStore.GetRoleByName(ctx, "ADMIN")
	if err != nil {
		return err
	}
	if err := txStore.AssignRoleToUser(ctx, user.ID, role.ID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Unexported — suppress unused import warning for pgtype.
var _ pgtype.UUID
