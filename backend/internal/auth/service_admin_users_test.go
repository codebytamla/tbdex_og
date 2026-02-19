package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestAdminListUsersSkipsCountWhenIncludeTotalFalse(t *testing.T) {
	countCalled := false
	store := stubAdminUsersStore{
		listUsersFn: func(_ context.Context, p AdminListUsersRequest) (UserPage, error) {
			if p.Limit != 10 {
				t.Fatalf("unexpected limit: got %d want 10", p.Limit)
			}
			return UserPage{
				Users:      []User{},
				HasNext:    false,
				NextCursor: nil,
			}, nil
		},
		countUsersFn: func(_ context.Context, _ AdminListUsersRequest) (int, error) {
			countCalled = true
			return 0, nil
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	resp, err := svc.AdminListUsers(context.Background(), AdminListUsersRequest{
		Limit:        10,
		IncludeTotal: false,
	})
	if err != nil {
		t.Fatalf("AdminListUsers error: %v", err)
	}

	if countCalled {
		t.Fatal("count should not be called when include_total=false")
	}
	if resp.Total != nil {
		t.Fatalf("expected total nil, got %v", *resp.Total)
	}
	if resp.Pagination.Limit != 10 {
		t.Fatalf("pagination limit mismatch: got %d", resp.Pagination.Limit)
	}
}

func TestAdminListUsersIncludeTotalAndEmptyRoleSiteSlices(t *testing.T) {
	next := "cursor-token"
	userID := uuid.New()
	countCalled := false
	store := stubAdminUsersStore{
		listUsersFn: func(_ context.Context, p AdminListUsersRequest) (UserPage, error) {
			if p.Limit != AdminUsersMaxLimit {
				t.Fatalf("limit should be clamped to max, got %d", p.Limit)
			}
			return UserPage{
				Users: []User{
					{
						ID:        userID,
						Email:     "alex@example.com",
						FullName:  "Alex",
						IsActive:  true,
						CreatedAt: time.Now().UTC(),
					},
				},
				HasNext:    true,
				NextCursor: &next,
			}, nil
		},
		countUsersFn: func(_ context.Context, _ AdminListUsersRequest) (int, error) {
			countCalled = true
			return 42, nil
		},
		listUserRolesByIDsFn: func(_ context.Context, _ []uuid.UUID) (map[uuid.UUID][]Role, error) {
			return map[uuid.UUID][]Role{}, nil
		},
		listUserSitesByIDsFn: func(_ context.Context, _ []uuid.UUID) (map[uuid.UUID][]SiteSummary, error) {
			return map[uuid.UUID][]SiteSummary{}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	resp, err := svc.AdminListUsers(context.Background(), AdminListUsersRequest{
		Limit:        999,
		IncludeTotal: true,
	})
	if err != nil {
		t.Fatalf("AdminListUsers error: %v", err)
	}

	if !countCalled {
		t.Fatal("expected count to be called when include_total=true")
	}
	if resp.Total == nil || *resp.Total != 42 {
		t.Fatalf("expected total 42, got %#v", resp.Total)
	}
	if !resp.Pagination.HasNext || resp.Pagination.NextCursor == nil || *resp.Pagination.NextCursor != next {
		t.Fatalf("pagination mismatch: %#v", resp.Pagination)
	}
	if len(resp.Users) != 1 {
		t.Fatalf("expected one user, got %d", len(resp.Users))
	}
	if resp.Users[0].Roles == nil || len(resp.Users[0].Roles) != 0 {
		t.Fatalf("expected empty roles slice, got %#v", resp.Users[0].Roles)
	}
	if resp.Users[0].Sites == nil || len(resp.Users[0].Sites) != 0 {
		t.Fatalf("expected empty sites slice, got %#v", resp.Users[0].Sites)
	}
}

func TestAdminGetUserSitesAndRoles(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	siteID := uuid.New()

	store := stubAdminUsersStore{
		listUserRolesFn: func(_ context.Context, gotUserID uuid.UUID) ([]Role, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return []Role{{ID: roleID, Name: "ADMIN"}}, nil
		},
		listUserSitesFn: func(_ context.Context, gotUserID uuid.UUID) ([]SiteSummary, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return []SiteSummary{{ID: siteID, Name: "Main Site", Code: "CLI-01"}}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	resp, err := svc.AdminGetUserSitesAndRoles(context.Background(), userID)
	if err != nil {
		t.Fatalf("AdminGetUserSitesAndRoles error: %v", err)
	}

	if len(resp.Roles) != 1 || resp.Roles[0].ID != roleID {
		t.Fatalf("roles mismatch: %#v", resp.Roles)
	}
	if len(resp.Sites) != 1 || resp.Sites[0].ID != siteID {
		t.Fatalf("sites mismatch: %#v", resp.Sites)
	}
}

func TestAdminGetUserSitesAndRolesNormalizesNilSlices(t *testing.T) {
	store := stubAdminUsersStore{
		listUserRolesFn: func(_ context.Context, _ uuid.UUID) ([]Role, error) {
			return nil, nil
		},
		listUserSitesFn: func(_ context.Context, _ uuid.UUID) ([]SiteSummary, error) {
			return nil, nil
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	resp, err := svc.AdminGetUserSitesAndRoles(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("AdminGetUserSitesAndRoles error: %v", err)
	}

	if resp.Roles == nil || len(resp.Roles) != 0 {
		t.Fatalf("expected empty roles slice, got %#v", resp.Roles)
	}
	if resp.Sites == nil || len(resp.Sites) != 0 {
		t.Fatalf("expected empty sites slice, got %#v", resp.Sites)
	}
}

func TestAdminSetRolesRejectsRemovingAdminRole(t *testing.T) {
	userID := uuid.New()

	store := stubAdminUsersStore{
		listUserRolesFn: func(_ context.Context, gotUserID uuid.UUID) ([]Role, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return []Role{{ID: uuid.New(), Name: "ADMIN"}}, nil
		},
		getRoleByNameFn: func(_ context.Context, name string) (Role, error) {
			if name == "DOCTOR" {
				return Role{ID: uuid.New(), Name: "DOCTOR"}, nil
			}
			return Role{}, pgx.ErrNoRows
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	_, err := svc.AdminSetRoles(context.Background(), userID, []string{"DOCTOR"}, DeviceInfo{})
	if err != ErrAdminRoleRequired {
		t.Fatalf("expected ErrAdminRoleRequired, got %v", err)
	}
}

func TestAdminCreateSiteRejectsLowercaseCode(t *testing.T) {
	store := stubAdminUsersStore{
		createSiteFn: func(_ context.Context, _, _, _ string) (Site, error) {
			return Site{}, &pgconn.PgError{
				Code:           "23514",
				ConstraintName: "chk_sites_code_uppercase",
			}
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	_, err := svc.AdminCreateSite(context.Background(), AdminCreateSiteRequest{
		Name:     "Main Site",
		Code:     "cli-01",
		SiteType: "clinic",
	}, DeviceInfo{})
	if err != ErrSiteCodeUppercase {
		t.Fatalf("expected ErrSiteCodeUppercase, got %v", err)
	}
}

func TestAdminUpdateSiteRejectsLowercaseCode(t *testing.T) {
	siteID := uuid.New()
	store := stubAdminUsersStore{
		updateSiteFn: func(_ context.Context, gotSiteID uuid.UUID, _, _, _ string) (Site, error) {
			if gotSiteID != siteID {
				t.Fatalf("unexpected site id: got %s want %s", gotSiteID, siteID)
			}
			return Site{}, &pgconn.PgError{
				Code:           "23514",
				ConstraintName: "chk_sites_code_uppercase",
			}
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	_, err := svc.AdminUpdateSite(context.Background(), siteID, AdminUpdateSiteRequest{
		Name:     "Main Site",
		Code:     "cli-01",
		SiteType: "clinic",
	}, DeviceInfo{})
	if err != ErrSiteCodeUppercase {
		t.Fatalf("expected ErrSiteCodeUppercase, got %v", err)
	}
}

func TestGetMeReturnsRolesAndPermissions(t *testing.T) {
	userID := uuid.New()
	store := stubAdminUsersStore{
		getUserByIDFn: func(_ context.Context, gotUserID uuid.UUID) (User, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return User{
				ID:       userID,
				Email:    "admin@example.com",
				FullName: "Admin User",
				IsActive: true,
			}, nil
		},
		listUserRolesFn: func(_ context.Context, gotUserID uuid.UUID) ([]Role, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return []Role{{ID: uuid.New(), Name: "ADMIN"}}, nil
		},
		getEffectivePermsFn: func(_ context.Context, gotUserID uuid.UUID) ([]string, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: got %s want %s", gotUserID, userID)
			}
			return []string{"users.read"}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, nil)
	profile, err := svc.GetMe(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetMe error: %v", err)
	}

	if len(profile.Roles) != 1 || profile.Roles[0].Name != "ADMIN" {
		t.Fatalf("roles mismatch: %#v", profile.Roles)
	}
	if len(profile.Permissions) != 1 || profile.Permissions[0] != "users.read" {
		t.Fatalf("permissions mismatch: %#v", profile.Permissions)
	}
}
