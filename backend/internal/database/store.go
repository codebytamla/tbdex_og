// Package database provides the concrete Store implementation that bridges
// sqlc-generated queries to the auth.Store interface.
package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/codebytamla/borderform-backend/internal/auth"
	"github.com/codebytamla/borderform-backend/internal/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGStore implements auth.Store by delegating to sqlc-generated Queries.
type PGStore struct {
	pool *pgxpool.Pool
	q    *db.Queries
	dbtx db.DBTX // pool or tx — used for seed helper SQL
}

func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool, q: db.New(pool), dbtx: pool}
}

// -----------------------------------------------------------------
// Type conversion helpers: uuid.UUID ↔ pgtype.UUID, time ↔ pgtype
// -----------------------------------------------------------------

func pgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func fromPgUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return p.Bytes
}

func pgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func fromPgTimestamptz(p pgtype.Timestamptz) time.Time {
	if !p.Valid {
		return time.Time{}
	}
	return p.Time
}

func pgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func pgUUIDSlice(ids []uuid.UUID) []pgtype.UUID {
	if len(ids) == 0 {
		return []pgtype.UUID{}
	}

	out := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		out[i] = pgUUID(id)
	}
	return out
}

func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// -----------------------------------------------------------------
// auth.Store implementation — Users
// -----------------------------------------------------------------

func (s *PGStore) CreateUser(ctx context.Context, p auth.CreateUserParams) (auth.User, error) {
	row, err := s.q.CreateUser(ctx, db.CreateUserParams{
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		FullName:     p.FullName,
	})
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:        fromPgUUID(row.ID),
		Email:     row.Email,
		FullName:  row.FullName,
		IsActive:  row.IsActive,
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
	}, nil
}

func (s *PGStore) GetUserByEmail(ctx context.Context, email string) (auth.UserWithPassword, error) {
	row, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return auth.UserWithPassword{}, err
	}
	return auth.UserWithPassword{
		User: auth.User{
			ID:        fromPgUUID(row.ID),
			Email:     row.Email,
			FullName:  row.FullName,
			IsActive:  row.IsActive,
			CreatedAt: fromPgTimestamptz(row.CreatedAt),
		},
		PasswordHash: row.PasswordHash,
	}, nil
}

func (s *PGStore) GetUserByID(ctx context.Context, id uuid.UUID) (auth.User, error) {
	row, err := s.q.GetUserByID(ctx, pgUUID(id))
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:        fromPgUUID(row.ID),
		Email:     row.Email,
		FullName:  row.FullName,
		IsActive:  row.IsActive,
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
	}, nil
}

func (s *PGStore) UpdateUserPassword(ctx context.Context, id uuid.UUID, hash string) error {
	return s.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		PasswordHash: hash,
		ID:           pgUUID(id),
	})
}

func (s *PGStore) UpdateUser(ctx context.Context, id uuid.UUID, p auth.UpdateUserParams) (auth.User, error) {
	row, err := s.q.UpdateUser(ctx, db.UpdateUserParams{
		FullName: pgText(p.FullName),
		ID:       pgUUID(id),
	})
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:        fromPgUUID(row.ID),
		Email:     row.Email,
		FullName:  row.FullName,
		IsActive:  row.IsActive,
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
	}, nil
}

func (s *PGStore) SetUserActive(ctx context.Context, id uuid.UUID, active bool) (auth.User, error) {
	row, err := s.q.SetUserActive(ctx, db.SetUserActiveParams{
		IsActive: active,
		ID:       pgUUID(id),
	})
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:        fromPgUUID(row.ID),
		Email:     row.Email,
		FullName:  row.FullName,
		IsActive:  row.IsActive,
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
	}, nil
}

func (s *PGStore) ListUsers(ctx context.Context, p auth.AdminListUsersRequest) (auth.UserPage, error) {
	limit := p.Limit
	if limit <= 0 {
		limit = auth.AdminUsersDefaultLimit
	}
	if limit > auth.AdminUsersMaxLimit {
		limit = auth.AdminUsersMaxLimit
	}

	search := escapeLikePattern(strings.TrimSpace(p.Search))
	useIsActive := p.IsActive != nil
	isActive := false
	if useIsActive {
		isActive = *p.IsActive
	}

	params := db.ListUsersPageParams{
		Search:          search,
		UseIsActive:     useIsActive,
		IsActive:        isActive,
		UseRoleFilter:   len(p.Roles) > 0,
		RoleNames:       p.Roles,
		UseSiteFilter:   len(p.SiteIDs) > 0,
		SiteIds:         pgUUIDSlice(p.SiteIDs),
		UseCursor:       false,
		CursorCreatedAt: pgtype.Timestamptz{},
		CursorID:        pgtype.UUID{},
		LimitRows:       int32(limit + 1),
	}

	if p.Cursor != "" {
		cursor, err := auth.DecodeAdminUsersCursor(p.Cursor)
		if err != nil {
			return auth.UserPage{}, err
		}
		params.UseCursor = true
		params.CursorCreatedAt = pgTimestamptz(cursor.CreatedAt)
		params.CursorID = pgUUID(cursor.ID)
	}

	rows, err := s.q.ListUsersPage(ctx, params)
	if err != nil {
		return auth.UserPage{}, err
	}

	hasNext := len(rows) > limit
	if hasNext {
		rows = rows[:limit]
	}

	users := make([]auth.User, len(rows))
	for i, row := range rows {
		users[i] = auth.User{
			ID:        fromPgUUID(row.ID),
			Email:     row.Email,
			FullName:  row.FullName,
			IsActive:  row.IsActive,
			CreatedAt: fromPgTimestamptz(row.CreatedAt),
		}
	}

	var nextCursor *string
	if hasNext && len(rows) > 0 {
		last := rows[len(rows)-1]
		cursor, err := auth.EncodeAdminUsersCursor(fromPgTimestamptz(last.CreatedAt), fromPgUUID(last.ID))
		if err != nil {
			return auth.UserPage{}, err
		}
		nextCursor = &cursor
	}

	return auth.UserPage{
		Users:      users,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}, nil
}

func (s *PGStore) CountUsers(ctx context.Context, p auth.AdminListUsersRequest) (int, error) {
	search := escapeLikePattern(strings.TrimSpace(p.Search))
	useIsActive := p.IsActive != nil
	isActive := false
	if useIsActive {
		isActive = *p.IsActive
	}

	count, err := s.q.CountUsersFiltered(ctx, db.CountUsersFilteredParams{
		Search:        search,
		UseIsActive:   useIsActive,
		IsActive:      isActive,
		UseRoleFilter: len(p.Roles) > 0,
		RoleNames:     p.Roles,
		UseSiteFilter: len(p.SiteIDs) > 0,
		SiteIds:       pgUUIDSlice(p.SiteIDs),
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *PGStore) ListUserRolesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]auth.Role, error) {
	result := make(map[uuid.UUID][]auth.Role, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	rows, err := s.q.ListRolesByUserIDs(ctx, pgUUIDSlice(userIDs))
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		userID := fromPgUUID(row.UserID)
		role := auth.Role{
			ID:          fromPgUUID(row.ID),
			Name:        row.Name,
			Description: row.Description,
		}
		result[userID] = append(result[userID], role)
	}

	return result, nil
}

func (s *PGStore) ListUserSitesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]auth.SiteSummary, error) {
	result := make(map[uuid.UUID][]auth.SiteSummary, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	rows, err := s.q.ListSitesByUserIDs(ctx, pgUUIDSlice(userIDs))
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		userID := fromPgUUID(row.UserID)
		site := auth.SiteSummary{
			ID:   fromPgUUID(row.ID),
			Name: row.Name,
			Code: row.Code,
		}
		result[userID] = append(result[userID], site)
	}

	return result, nil
}

// -----------------------------------------------------------------
// auth.Store implementation — Refresh tokens
// -----------------------------------------------------------------

func (s *PGStore) CreateRefreshToken(ctx context.Context, p auth.CreateRefreshTokenParams) (auth.RefreshTokenID, error) {
	row, err := s.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     pgUUID(p.UserID),
		TokenHash:  p.TokenHash,
		DeviceName: p.DeviceName,
		CreatedIp:  p.CreatedIP,
		UserAgent:  p.UserAgent,
		ExpiresAt:  pgTimestamptz(p.ExpiresAt),
	})
	if err != nil {
		return auth.RefreshTokenID{}, err
	}
	return auth.RefreshTokenID{
		ID:        fromPgUUID(row.ID),
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
	}, nil
}

func (s *PGStore) GetRefreshTokenByHash(ctx context.Context, hash string) (auth.RefreshToken, error) {
	row, err := s.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return auth.RefreshToken{}, err
	}
	return auth.RefreshToken{
		ID:        fromPgUUID(row.ID),
		UserID:    fromPgUUID(row.UserID),
		ExpiresAt: fromPgTimestamptz(row.ExpiresAt),
		Revoked:   row.Revoked,
	}, nil
}

func (s *PGStore) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	return s.q.RevokeRefreshToken(ctx, pgUUID(id))
}

func (s *PGStore) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return s.q.RevokeAllUserRefreshTokens(ctx, pgUUID(userID))
}

// -----------------------------------------------------------------
// auth.Store implementation — Password reset
// -----------------------------------------------------------------

func (s *PGStore) CreatePasswordResetToken(ctx context.Context, p auth.CreatePasswordResetTokenParams) error {
	return s.q.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID:    pgUUID(p.UserID),
		TokenHash: p.TokenHash,
		ExpiresAt: pgTimestamptz(p.ExpiresAt),
	})
}

func (s *PGStore) GetPasswordResetTokenByHash(ctx context.Context, hash string) (auth.PasswordResetToken, error) {
	row, err := s.q.GetPasswordResetTokenByHash(ctx, hash)
	if err != nil {
		return auth.PasswordResetToken{}, err
	}
	return auth.PasswordResetToken{
		ID:        fromPgUUID(row.ID),
		UserID:    fromPgUUID(row.UserID),
		ExpiresAt: fromPgTimestamptz(row.ExpiresAt),
		Used:      row.Used,
	}, nil
}

func (s *PGStore) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	return s.q.MarkPasswordResetTokenUsed(ctx, pgUUID(id))
}

func (s *PGStore) ExpireOldPasswordResetTokens(ctx context.Context, userID uuid.UUID) error {
	return s.q.ExpireOldPasswordResetTokens(ctx, pgUUID(userID))
}

// -----------------------------------------------------------------
// auth.Store implementation — RBAC
// -----------------------------------------------------------------

func (s *PGStore) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.q.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: pgUUID(userID),
		RoleID: pgUUID(roleID),
	})
}

func (s *PGStore) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.q.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: pgUUID(userID),
		RoleID: pgUUID(roleID),
	})
}

func (s *PGStore) RemoveAllUserRoles(ctx context.Context, userID uuid.UUID) error {
	return s.q.RemoveAllUserRoles(ctx, pgUUID(userID))
}

func (s *PGStore) GetRoleByName(ctx context.Context, name string) (auth.Role, error) {
	row, err := s.q.GetRoleByName(ctx, name)
	if err != nil {
		return auth.Role{}, err
	}
	return auth.Role{
		ID:          fromPgUUID(row.ID),
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (s *PGStore) GetRoleByID(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	row, err := s.q.GetRoleByID(ctx, pgUUID(id))
	if err != nil {
		return auth.Role{}, err
	}
	return auth.Role{
		ID:          fromPgUUID(row.ID),
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (s *PGStore) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]auth.Role, error) {
	rows, err := s.q.ListUserRoles(ctx, pgUUID(userID))
	if err != nil {
		return nil, err
	}
	roles := make([]auth.Role, len(rows))
	for i, r := range rows {
		roles[i] = auth.Role{
			ID:          fromPgUUID(r.ID),
			Name:        r.Name,
			Description: r.Description,
		}
	}
	return roles, nil
}

func (s *PGStore) GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.q.GetEffectivePermissions(ctx, pgUUID(userID))
}

// -----------------------------------------------------------------
// auth.Store implementation — Role management
// -----------------------------------------------------------------

func (s *PGStore) ListRoles(ctx context.Context) ([]auth.Role, error) {
	rows, err := s.q.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	roles := make([]auth.Role, len(rows))
	for i, r := range rows {
		roles[i] = auth.Role{
			ID:          fromPgUUID(r.ID),
			Name:        r.Name,
			Description: r.Description,
		}
	}
	return roles, nil
}

func (s *PGStore) CreateRole(ctx context.Context, name, description string) (auth.Role, error) {
	row, err := s.q.CreateRole(ctx, db.CreateRoleParams{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return auth.Role{}, err
	}
	return auth.Role{
		ID:          fromPgUUID(row.ID),
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (s *PGStore) UpdateRole(ctx context.Context, id uuid.UUID, name, description *string) (auth.Role, error) {
	row, err := s.q.UpdateRole(ctx, db.UpdateRoleParams{
		Name:        pgText(name),
		Description: pgText(description),
		ID:          pgUUID(id),
	})
	if err != nil {
		return auth.Role{}, err
	}
	return auth.Role{
		ID:          fromPgUUID(row.ID),
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (s *PGStore) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	rows, err := s.q.ListRolePermissions(ctx, pgUUID(roleID))
	if err != nil {
		return nil, err
	}
	perms := make([]auth.Permission, len(rows))
	for i, p := range rows {
		perms[i] = auth.Permission{
			ID:   fromPgUUID(p.ID),
			Name: p.Name,
		}
	}
	return perms, nil
}

func (s *PGStore) ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, permissionNames []string) ([]auth.Permission, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.q.WithTx(tx)

	if err := qtx.DeleteRolePermissions(ctx, pgUUID(roleID)); err != nil {
		return nil, err
	}

	for _, name := range permissionNames {
		if err := qtx.InsertRolePermissionByName(ctx, db.InsertRolePermissionByNameParams{
			RoleID:         pgUUID(roleID),
			PermissionName: name,
		}); err != nil {
			return nil, err
		}
	}

	rows, err := qtx.ListRolePermissions(ctx, pgUUID(roleID))
	if err != nil {
		return nil, err
	}
	perms := make([]auth.Permission, len(rows))
	for i, p := range rows {
		perms[i] = auth.Permission{
			ID:   fromPgUUID(p.ID),
			Name: p.Name,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return perms, nil
}

func (s *PGStore) ListAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	rows, err := s.q.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}
	perms := make([]auth.Permission, len(rows))
	for i, p := range rows {
		perms[i] = auth.Permission{
			ID:   fromPgUUID(p.ID),
			Name: p.Name,
		}
	}
	return perms, nil
}

// -----------------------------------------------------------------
// auth.Store implementation — Sites / user-site assignments
// -----------------------------------------------------------------

func (s *PGStore) AssignUserToSite(ctx context.Context, userID, siteID uuid.UUID) error {
	return s.q.AssignUserToSite(ctx, db.AssignUserToSiteParams{
		UserID: pgUUID(userID),
		SiteID: pgUUID(siteID),
	})
}

func (s *PGStore) RemoveUserFromSite(ctx context.Context, userID, siteID uuid.UUID) error {
	return s.q.RemoveUserFromSite(ctx, db.RemoveUserFromSiteParams{
		UserID: pgUUID(userID),
		SiteID: pgUUID(siteID),
	})
}

func (s *PGStore) ListUserSites(ctx context.Context, userID uuid.UUID) ([]auth.SiteSummary, error) {
	rows, err := s.q.ListUserSites(ctx, pgUUID(userID))
	if err != nil {
		return nil, err
	}
	sites := make([]auth.SiteSummary, len(rows))
	for i, c := range rows {
		sites[i] = auth.SiteSummary{
			ID:   fromPgUUID(c.ID),
			Name: c.Name,
			Code: c.Code,
		}
	}
	return sites, nil
}

func (s *PGStore) RemoveAllUserSites(ctx context.Context, userID uuid.UUID) error {
	return s.q.RemoveAllUserSites(ctx, pgUUID(userID))
}

func (s *PGStore) GetSiteSummaryByID(ctx context.Context, id uuid.UUID) (auth.SiteSummary, error) {
	row, err := s.q.GetSiteSummaryByID(ctx, pgUUID(id))
	if err != nil {
		return auth.SiteSummary{}, err
	}
	return auth.SiteSummary{
		ID:   fromPgUUID(row.ID),
		Name: row.Name,
		Code: row.Code,
	}, nil
}

// -----------------------------------------------------------------
// auth.Store implementation — Sites
// -----------------------------------------------------------------

func (s *PGStore) CreateSite(ctx context.Context, name, code, siteType string) (auth.Site, error) {
	row, err := s.q.CreateSite(ctx, db.CreateSiteParams{
		Name:     name,
		Code:     code,
		SiteType: siteType,
	})
	if err != nil {
		return auth.Site{}, err
	}
	return siteFromDB(row), nil
}

func (s *PGStore) ListSites(ctx context.Context) ([]auth.Site, error) {
	rows, err := s.q.ListSites(ctx)
	if err != nil {
		return nil, err
	}
	sites := make([]auth.Site, len(rows))
	for i, r := range rows {
		sites[i] = siteFromDB(r)
	}
	return sites, nil
}

func (s *PGStore) GetSiteByID(ctx context.Context, id uuid.UUID) (auth.Site, error) {
	row, err := s.q.GetSiteByID(ctx, pgUUID(id))
	if err != nil {
		return auth.Site{}, err
	}
	return siteFromDB(row), nil
}

func (s *PGStore) UpdateSite(ctx context.Context, id uuid.UUID, name, code, siteType string) (auth.Site, error) {
	row, err := s.q.UpdateSite(ctx, db.UpdateSiteParams{
		Name:     name,
		Code:     code,
		SiteType: siteType,
		ID:       pgUUID(id),
	})
	if err != nil {
		return auth.Site{}, err
	}
	return siteFromDB(row), nil
}

func (s *PGStore) DisableSite(ctx context.Context, id uuid.UUID) (auth.Site, error) {
	row, err := s.q.DisableSite(ctx, pgUUID(id))
	if err != nil {
		return auth.Site{}, err
	}
	return siteFromDB(row), nil
}

func siteFromDB(s db.Site) auth.Site {
	return auth.Site{
		ID:        fromPgUUID(s.ID),
		Name:      s.Name,
		Code:      s.Code,
		SiteType:  s.SiteType,
		IsActive:  s.IsActive,
		CreatedAt: fromPgTimestamptz(s.CreatedAt),
		UpdatedAt: fromPgTimestamptz(s.UpdatedAt),
	}
}

// -----------------------------------------------------------------
// auth.Store implementation — Patients
// -----------------------------------------------------------------

func patientFromRow(id pgtype.UUID, fullName string, dob pgtype.Date, isDobEst bool,
	sex, phone, nationalID, mrn, address, village, origin pgtype.Text, status string,
	createdBy pgtype.UUID, createdAt, updatedAt pgtype.Timestamptz) auth.Patient {

	p := auth.Patient{
		ID:             fromPgUUID(id),
		FullName:       fullName,
		IsDobEstimated: isDobEst,
		Status:         status,
		CreatedBy:      fromPgUUID(createdBy),
		CreatedAt:      fromPgTimestamptz(createdAt),
		UpdatedAt:      fromPgTimestamptz(updatedAt),
	}
	if dob.Valid {
		s := dob.Time.Format("2006-01-02")
		p.DateOfBirth = &s
	}
	if sex.Valid {
		p.Sex = &sex.String
	}
	if phone.Valid {
		p.Phone = &phone.String
	}
	if nationalID.Valid {
		p.NationalID = &nationalID.String
	}
	if mrn.Valid {
		p.MRN = &mrn.String
	}
	if address.Valid {
		p.Address = &address.String
	}
	if village.Valid {
		p.Village = &village.String
	}
	if origin.Valid {
		p.Origin = &origin.String
	}
	return p
}

func pgDate(s *string) pgtype.Date {
	if s == nil {
		return pgtype.Date{}
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func pgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func (s *PGStore) CreatePatient(ctx context.Context, p auth.CreatePatientParams) (auth.Patient, error) {
	row, err := s.q.CreatePatient(ctx, db.CreatePatientParams{
		FullName:       p.FullName,
		DateOfBirth:    pgDate(p.DateOfBirth),
		IsDobEstimated: p.IsDobEstimated,
		Sex:            pgText(p.Sex),
		Phone:          pgText(p.Phone),
		NationalID:     pgText(p.NationalID),
		Mrn:            pgText(p.MRN),
		Address:        pgText(p.Address),
		Village:        pgText(p.Village),
		Origin:         pgText(p.Origin),
		Status:         p.Status,
		CreatedBy:      pgUUID(p.CreatedBy),
	})
	if err != nil {
		return auth.Patient{}, err
	}
	return patientFromRow(
		row.ID, row.FullName, row.DateOfBirth, row.IsDobEstimated,
		row.Sex, row.Phone, row.NationalID, row.Mrn,
		row.Address, row.Village, row.Origin, row.Status,
		row.CreatedBy, row.CreatedAt, row.UpdatedAt,
	), nil
}

func (s *PGStore) GetPatientByID(ctx context.Context, id uuid.UUID) (auth.Patient, error) {
	row, err := s.q.GetPatientByID(ctx, pgUUID(id))
	if err != nil {
		return auth.Patient{}, err
	}
	return patientFromRow(
		row.ID, row.FullName, row.DateOfBirth, row.IsDobEstimated,
		row.Sex, row.Phone, row.NationalID, row.Mrn,
		row.Address, row.Village, row.Origin, row.Status,
		row.CreatedBy, row.CreatedAt, row.UpdatedAt,
	), nil
}

func (s *PGStore) ListPatients(ctx context.Context, p auth.ListPatientsRequest) (auth.PatientPage, error) {
	limit := p.Limit
	if limit <= 0 {
		limit = auth.AdminUsersDefaultLimit
	}
	if limit > auth.AdminUsersMaxLimit {
		limit = auth.AdminUsersMaxLimit
	}

	search := escapeLikePattern(strings.TrimSpace(p.Search))
	useStatusFilter := p.Status != nil
	status := ""
	if useStatusFilter {
		status = *p.Status
	}

	params := db.ListPatientsParams{
		Search:          search,
		UseStatusFilter: useStatusFilter,
		Status:          status,
		UseCursor:       false,
		CursorCreatedAt: pgtype.Timestamptz{},
		CursorID:        pgtype.UUID{},
		LimitRows:       int32(limit + 1), // fetch one extra to detect has_next
	}

	if p.Cursor != "" {
		cursor, err := auth.DecodeAdminUsersCursor(p.Cursor)
		if err != nil {
			return auth.PatientPage{}, err
		}
		params.UseCursor = true
		params.CursorCreatedAt = pgTimestamptz(cursor.CreatedAt)
		params.CursorID = pgUUID(cursor.ID)
	}

	rows, err := s.q.ListPatients(ctx, params)
	if err != nil {
		return auth.PatientPage{}, err
	}

	hasNext := len(rows) > limit
	if hasNext {
		rows = rows[:limit]
	}

	patients := make([]auth.Patient, len(rows))
	for i, row := range rows {
		patients[i] = patientFromRow(
			row.ID, row.FullName, row.DateOfBirth, row.IsDobEstimated,
			row.Sex, row.Phone, row.NationalID, row.Mrn,
			row.Address, row.Village, row.Origin, row.Status,
			row.CreatedBy, row.CreatedAt, row.UpdatedAt,
		)
	}

	var nextCursor *string
	if hasNext && len(rows) > 0 {
		last := rows[len(rows)-1]
		c, err := auth.EncodeAdminUsersCursor(fromPgTimestamptz(last.CreatedAt), fromPgUUID(last.ID))
		if err != nil {
			return auth.PatientPage{}, err
		}
		nextCursor = &c
	}

	return auth.PatientPage{
		Patients:   patients,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}, nil
}

func (s *PGStore) CountPatients(ctx context.Context, p auth.ListPatientsRequest) (int, error) {
	search := escapeLikePattern(strings.TrimSpace(p.Search))
	useStatusFilter := p.Status != nil
	status := ""
	if useStatusFilter {
		status = *p.Status
	}

	count, err := s.q.CountPatients(ctx, db.CountPatientsParams{
		Search:          search,
		UseStatusFilter: useStatusFilter,
		Status:          status,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *PGStore) UpdatePatient(ctx context.Context, id uuid.UUID, p auth.UpdatePatientParams) (auth.Patient, error) {
	row, err := s.q.UpdatePatient(ctx, db.UpdatePatientParams{
		FullName:       pgText(p.FullName),
		DateOfBirth:    pgDate(p.DateOfBirth),
		IsDobEstimated: pgBool(p.IsDobEstimated),
		Sex:            pgText(p.Sex),
		Phone:          pgText(p.Phone),
		NationalID:     pgText(p.NationalID),
		Address:        pgText(p.Address),
		Village:        pgText(p.Village),
		Origin:         pgText(p.Origin),
		Status:         pgText(p.Status),
		ID:             pgUUID(id),
	})
	if err != nil {
		return auth.Patient{}, err
	}
	return patientFromRow(
		row.ID, row.FullName, row.DateOfBirth, row.IsDobEstimated,
		row.Sex, row.Phone, row.NationalID, row.Mrn,
		row.Address, row.Village, row.Origin, row.Status,
		row.CreatedBy, row.CreatedAt, row.UpdatedAt,
	), nil
}

func (s *PGStore) DisablePatient(ctx context.Context, id uuid.UUID) (auth.Patient, error) {
	row, err := s.q.DisablePatient(ctx, pgUUID(id))
	if err != nil {
		return auth.Patient{}, err
	}
	return patientFromRow(
		row.ID, row.FullName, row.DateOfBirth, row.IsDobEstimated,
		row.Sex, row.Phone, row.NationalID, row.Mrn,
		row.Address, row.Village, row.Origin, row.Status,
		row.CreatedBy, row.CreatedAt, row.UpdatedAt,
	), nil
}

func (s *PGStore) GetPatientCreatedBy(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	createdBy, err := s.q.GetPatientCreatedBy(ctx, pgUUID(id))
	if err != nil {
		return uuid.Nil, err
	}
	return fromPgUUID(createdBy), nil
}

func (s *PGStore) UserIsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.q.UserIsAdmin(ctx, pgUUID(userID))
}

// -----------------------------------------------------------------
// auth.Store implementation — Transactions
// -----------------------------------------------------------------

func (s *PGStore) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.pool.Begin(ctx)
}

func (s *PGStore) WithTx(tx pgx.Tx) auth.Store {
	return &PGStore{pool: s.pool, q: s.q.WithTx(tx), dbtx: tx}
}

// -----------------------------------------------------------------
// middleware.PermissionChecker implementation
// -----------------------------------------------------------------

func (s *PGStore) UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	has, err := s.q.UserHasPermission(ctx, db.UserHasPermissionParams{
		UserID:         pgUUID(userID),
		PermissionName: permission,
	})
	if err != nil || has {
		return has, err
	}

	for _, alias := range legacyPermissionAliases(permission) {
		allowed, aliasErr := s.q.UserHasPermission(ctx, db.UserHasPermissionParams{
			UserID:         pgUUID(userID),
			PermissionName: alias,
		})
		if aliasErr != nil {
			return false, aliasErr
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

func legacyPermissionAliases(permission string) []string {
	switch permission {
	case "patients.read", "patients.create", "patients.update", "patients.disable":
		return []string{"patients.write"}
	default:
		return nil
	}
}

// -----------------------------------------------------------------
// middleware.SiteChecker implementation
// -----------------------------------------------------------------

func (s *PGStore) UserBelongsToSite(ctx context.Context, userID, siteID uuid.UUID) (bool, error) {
	return s.q.UserBelongsToSite(ctx, db.UserBelongsToSiteParams{
		UserID: pgUUID(userID),
		SiteID: pgUUID(siteID),
	})
}

// -----------------------------------------------------------------
// auth.AuditWriter implementation
// -----------------------------------------------------------------

func (s *PGStore) InsertAuditEvent(ctx context.Context, p auth.InsertAuditEventParams) error {
	return s.q.InsertAuditEvent(ctx, db.InsertAuditEventParams{
		EventType: p.EventType,
		ActorID:   p.ActorID,
		TargetID:  p.TargetID,
		IpAddress: p.IPAddress,
		UserAgent: p.UserAgent,
		Metadata:  p.Metadata,
	})
}

// -----------------------------------------------------------------
// Seed — runs on startup, idempotent
// -----------------------------------------------------------------

func (s *PGStore) SeedRolesAndPermissions(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	exec := func(op, sql string, args ...any) error {
		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	// Roles
	roles := []struct{ name, desc string }{
		{"ADMIN", "Full system administrator"},
		{"DOCTOR", "Medical doctor / clinician"},
		{"CHW", "Community health worker"},
		{"DATA_ENTRY", "Data entry operator"},
		{"SUPERVISOR", "Clinical supervisor / manager"},
	}
	for _, r := range roles {
		if err := exec("insert role", `INSERT INTO roles (name, description) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`,
			r.name, r.desc); err != nil {
			return err
		}
	}

	// Permissions
	perms := []string{
		"users.create", "users.read", "users.update", "users.disable",
		"users.assign_roles", "users.assign_sites", "users.reset_password",
		"patients.create", "patients.read", "patients.update", "patients.disable",
		"patients.write", // legacy — kept for backward compatibility
		"encounters.read", "encounters.write",
		"forms.create", "forms.read", "forms.update", "forms.disable",
		"reports.export",
		"sites.create", "sites.read", "sites.update", "sites.disable",
		"roles.create", "roles.read", "roles.update",
		"permissions.read",
		"audit.read",
	}
	for _, p := range perms {
		if err := exec("insert permission", `INSERT INTO permissions (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, p); err != nil {
			return err
		}
	}

	// Role-permission mappings
	rolePerms := map[string][]string{
		"ADMIN":      perms, // all
		"DOCTOR":     {"patients.create", "patients.read", "patients.update", "patients.write", "encounters.read", "encounters.write", "reports.export"},
		"CHW":        {"patients.create", "patients.read", "patients.update", "patients.write", "encounters.read", "encounters.write"},
		"DATA_ENTRY": {"patients.create", "patients.read", "patients.update", "patients.write", "encounters.read", "encounters.write", "forms.create", "forms.read", "forms.update", "forms.disable"},
		"SUPERVISOR": {"patients.read", "encounters.read", "encounters.write", "reports.export", "users.read", "audit.read"},
	}
	for roleName, permNames := range rolePerms {
		for _, permName := range permNames {
			if err := exec("insert role permission", `INSERT INTO role_permissions (role_id, permission_id)
					 SELECT r.id, p.id FROM roles r CROSS JOIN permissions p
					 WHERE r.name = $1 AND p.name = $2
					 ON CONFLICT DO NOTHING`,
				roleName, permName); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}
	return nil
}

// Suppress unused import warnings.
var _ time.Duration
