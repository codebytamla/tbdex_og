package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type stubAdminUsersStore struct {
	Store

	listUsersFn          func(context.Context, AdminListUsersRequest) (UserPage, error)
	countUsersFn         func(context.Context, AdminListUsersRequest) (int, error)
	getUserByIDFn        func(context.Context, uuid.UUID) (User, error)
	getEffectivePermsFn  func(context.Context, uuid.UUID) ([]string, error)
	listUserRolesByIDsFn func(context.Context, []uuid.UUID) (map[uuid.UUID][]Role, error)
	listUserSitesByIDsFn func(context.Context, []uuid.UUID) (map[uuid.UUID][]SiteSummary, error)
	listUserRolesFn      func(context.Context, uuid.UUID) ([]Role, error)
	listUserSitesFn      func(context.Context, uuid.UUID) ([]SiteSummary, error)
	getRoleByNameFn      func(context.Context, string) (Role, error)
	createSiteFn         func(context.Context, string, string, string) (Site, error)
	updateSiteFn         func(context.Context, uuid.UUID, string, string, string) (Site, error)
}

func (s stubAdminUsersStore) ListUsers(ctx context.Context, p AdminListUsersRequest) (UserPage, error) {
	if s.listUsersFn == nil {
		return UserPage{Users: []User{}}, nil
	}
	return s.listUsersFn(ctx, p)
}

func (s stubAdminUsersStore) CountUsers(ctx context.Context, p AdminListUsersRequest) (int, error) {
	if s.countUsersFn == nil {
		return 0, nil
	}
	return s.countUsersFn(ctx, p)
}

func (s stubAdminUsersStore) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	if s.getUserByIDFn == nil {
		return User{ID: id}, nil
	}
	return s.getUserByIDFn(ctx, id)
}

func (s stubAdminUsersStore) GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	if s.getEffectivePermsFn == nil {
		return []string{}, nil
	}
	return s.getEffectivePermsFn(ctx, userID)
}

func (s stubAdminUsersStore) ListUserRolesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]Role, error) {
	if s.listUserRolesByIDsFn == nil {
		return map[uuid.UUID][]Role{}, nil
	}
	return s.listUserRolesByIDsFn(ctx, userIDs)
}

func (s stubAdminUsersStore) ListUserSitesByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]SiteSummary, error) {
	if s.listUserSitesByIDsFn == nil {
		return map[uuid.UUID][]SiteSummary{}, nil
	}
	return s.listUserSitesByIDsFn(ctx, userIDs)
}

func (s stubAdminUsersStore) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	if s.listUserRolesFn == nil {
		return []Role{}, nil
	}
	return s.listUserRolesFn(ctx, userID)
}

func (s stubAdminUsersStore) ListUserSites(ctx context.Context, userID uuid.UUID) ([]SiteSummary, error) {
	if s.listUserSitesFn == nil {
		return []SiteSummary{}, nil
	}
	return s.listUserSitesFn(ctx, userID)
}

func (s stubAdminUsersStore) GetRoleByName(ctx context.Context, name string) (Role, error) {
	if s.getRoleByNameFn == nil {
		return Role{}, pgx.ErrNoRows
	}
	return s.getRoleByNameFn(ctx, name)
}

func (s stubAdminUsersStore) CreateSite(ctx context.Context, name, code, siteType string) (Site, error) {
	if s.createSiteFn == nil {
		return Site{}, nil
	}
	return s.createSiteFn(ctx, name, code, siteType)
}

func (s stubAdminUsersStore) UpdateSite(ctx context.Context, id uuid.UUID, name, code, siteType string) (Site, error) {
	if s.updateSiteFn == nil {
		return Site{}, nil
	}
	return s.updateSiteFn(ctx, id, name, code, siteType)
}
