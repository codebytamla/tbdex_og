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
	getRoleByIDFn        func(context.Context, uuid.UUID) (Role, error)
	listRolesFn          func(context.Context) ([]Role, error)
	createRoleFn         func(context.Context, string, string) (Role, error)
	updateRoleFn         func(context.Context, uuid.UUID, *string, *string) (Role, error)
	listRolePermsFn      func(context.Context, uuid.UUID) ([]Permission, error)
	replaceRolePermsFn   func(context.Context, uuid.UUID, []string) ([]Permission, error)
	listAllPermsFn       func(context.Context) ([]Permission, error)
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

func (s stubAdminUsersStore) GetRoleByID(ctx context.Context, id uuid.UUID) (Role, error) {
	if s.getRoleByIDFn == nil {
		return Role{}, pgx.ErrNoRows
	}
	return s.getRoleByIDFn(ctx, id)
}

func (s stubAdminUsersStore) ListRoles(ctx context.Context) ([]Role, error) {
	if s.listRolesFn == nil {
		return []Role{}, nil
	}
	return s.listRolesFn(ctx)
}

func (s stubAdminUsersStore) CreateRole(ctx context.Context, name, description string) (Role, error) {
	if s.createRoleFn == nil {
		return Role{}, nil
	}
	return s.createRoleFn(ctx, name, description)
}

func (s stubAdminUsersStore) UpdateRole(ctx context.Context, id uuid.UUID, name, description *string) (Role, error) {
	if s.updateRoleFn == nil {
		return Role{ID: id}, nil
	}
	return s.updateRoleFn(ctx, id, name, description)
}

func (s stubAdminUsersStore) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error) {
	if s.listRolePermsFn == nil {
		return []Permission{}, nil
	}
	return s.listRolePermsFn(ctx, roleID)
}

func (s stubAdminUsersStore) ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, permissionNames []string) ([]Permission, error) {
	if s.replaceRolePermsFn == nil {
		return []Permission{}, nil
	}
	return s.replaceRolePermsFn(ctx, roleID, permissionNames)
}

func (s stubAdminUsersStore) ListAllPermissions(ctx context.Context) ([]Permission, error) {
	if s.listAllPermsFn == nil {
		return []Permission{}, nil
	}
	return s.listAllPermsFn(ctx)
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
