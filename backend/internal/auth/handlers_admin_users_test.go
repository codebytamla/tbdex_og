package auth

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestParseAdminListUsersRequestValid(t *testing.T) {
	siteID := uuid.New()
	req := httptest.NewRequest("GET",
		"/admin/users?limit=200&search=%20alex%20&is_active=true&include_total=true&role=ADMIN&role=DOCTOR&site_id="+siteID.String(),
		nil,
	)

	parsed, err := parseAdminListUsersRequest(req)
	if err != nil {
		t.Fatalf("parse request: %v", err)
	}

	if parsed.Limit != AdminUsersMaxLimit {
		t.Fatalf("limit mismatch: got %d want %d", parsed.Limit, AdminUsersMaxLimit)
	}
	if parsed.Search != "alex" {
		t.Fatalf("search mismatch: got %q", parsed.Search)
	}
	if parsed.IsActive == nil || !*parsed.IsActive {
		t.Fatalf("expected is_active true, got %#v", parsed.IsActive)
	}
	if !parsed.IncludeTotal {
		t.Fatal("expected include_total true")
	}
	if len(parsed.Roles) != 2 || parsed.Roles[0] != "ADMIN" || parsed.Roles[1] != "DOCTOR" {
		t.Fatalf("roles mismatch: got %#v", parsed.Roles)
	}
	if len(parsed.SiteIDs) != 1 || parsed.SiteIDs[0] != siteID {
		t.Fatalf("site_ids mismatch: got %#v", parsed.SiteIDs)
	}
}

func TestParseAdminListUsersRequestInvalidSiteID(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin/users?site_id=not-uuid", nil)
	_, err := parseAdminListUsersRequest(req)
	if err == nil || err.Error() != "invalid site_id" {
		t.Fatalf("expected invalid site_id error, got %v", err)
	}
}

func TestParseAdminListUsersRequestInvalidBoolean(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin/users?include_total=yes", nil)
	_, err := parseAdminListUsersRequest(req)
	if err == nil || err.Error() != "invalid include_total" {
		t.Fatalf("expected invalid include_total error, got %v", err)
	}
}

func TestHandleAdminListUsersInvalidCursor(t *testing.T) {
	store := stubAdminUsersStore{
		listUsersFn: func(_ context.Context, _ AdminListUsersRequest) (UserPage, error) {
			return UserPage{}, ErrInvalidCursor
		},
	}
	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("GET", "/admin/users?cursor=bad", nil)
	rr := httptest.NewRecorder()

	h.handleAdminListUsers(rr, req)

	if rr.Code != 400 {
		t.Fatalf("status mismatch: got %d want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "invalid cursor") {
		t.Fatalf("expected invalid cursor message, got %s", rr.Body.String())
	}
}

func TestHandleAdminGetUserSitesAndRolesInvalidUserID(t *testing.T) {
	h := NewHandler(NewService(stubAdminUsersStore{}, JWTConfig{}, nil))

	req := httptest.NewRequest("GET", "/admin/users/not-uuid/sites", nil)
	req.SetPathValue("id", "not-uuid")
	rr := httptest.NewRecorder()

	h.handleAdminGetUserSitesAndRoles(rr, req)

	if rr.Code != 400 {
		t.Fatalf("status mismatch: got %d want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "invalid user id") {
		t.Fatalf("expected invalid user id message, got %s", rr.Body.String())
	}
}

func TestHandleAdminGetUserSitesAndRolesSuccess(t *testing.T) {
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

	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("GET", "/admin/users/"+userID.String()+"/sites", nil)
	req.SetPathValue("id", userID.String())
	rr := httptest.NewRecorder()

	h.handleAdminGetUserSitesAndRoles(rr, req)

	if rr.Code != 200 {
		t.Fatalf("status mismatch: got %d want 200", rr.Code)
	}

	var body UserSitesAndRolesResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Roles) != 1 || body.Roles[0].ID != roleID {
		t.Fatalf("roles mismatch: %#v", body.Roles)
	}
	if len(body.Sites) != 1 || body.Sites[0].ID != siteID {
		t.Fatalf("sites mismatch: %#v", body.Sites)
	}
}

func TestHandleAdminSetRolesRejectsRemovingAdminRole(t *testing.T) {
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

	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("PUT", "/admin/users/"+userID.String()+"/roles", strings.NewReader(`{"roles":["DOCTOR"]}`))
	req.SetPathValue("id", userID.String())
	rr := httptest.NewRecorder()

	h.handleAdminSetRoles(rr, req)

	if rr.Code != 409 {
		t.Fatalf("status mismatch: got %d want 409", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "admin user must keep ADMIN role") {
		t.Fatalf("expected protected admin error message, got %s", rr.Body.String())
	}
}

func TestHandleAdminCreateSiteRejectsLowercaseCode(t *testing.T) {
	store := stubAdminUsersStore{
		createSiteFn: func(_ context.Context, _, _, _ string) (Site, error) {
			return Site{}, &pgconn.PgError{
				Code:           "23514",
				ConstraintName: "chk_sites_code_uppercase",
			}
		},
	}
	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("POST", "/admin/sites", strings.NewReader(`{"name":"Main Site","code":"cli-01","site_type":"clinic"}`))
	rr := httptest.NewRecorder()

	h.handleAdminCreateSite(rr, req)

	if rr.Code != 400 {
		t.Fatalf("status mismatch: got %d want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "site code must be uppercase") {
		t.Fatalf("expected site code uppercase error message, got %s", rr.Body.String())
	}
}

func TestHandleAdminUpdateSiteRejectsLowercaseCode(t *testing.T) {
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
	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("PUT", "/admin/sites/"+siteID.String(), strings.NewReader(`{"name":"Main Site","code":"cli-01","site_type":"clinic"}`))
	req.SetPathValue("id", siteID.String())
	rr := httptest.NewRecorder()

	h.handleAdminUpdateSite(rr, req)

	if rr.Code != 400 {
		t.Fatalf("status mismatch: got %d want 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "site code must be uppercase") {
		t.Fatalf("expected site code uppercase error message, got %s", rr.Body.String())
	}
}

func TestHandleMeOmitsIsAdmin(t *testing.T) {
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

	h := NewHandler(NewService(store, JWTConfig{}, nil))

	req := httptest.NewRequest("GET", "/auth/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), CtxUserIDKey, userID))
	rr := httptest.NewRecorder()

	h.handleMe(rr, req)

	if rr.Code != 200 {
		t.Fatalf("status mismatch: got %d want 200", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := body["is_admin"]; ok {
		t.Fatal("expected /auth/me response to omit is_admin")
	}
}
