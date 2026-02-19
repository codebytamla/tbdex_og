package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type noopAuditWriter struct{}

func (noopAuditWriter) InsertAuditEvent(_ context.Context, _ InsertAuditEventParams) error {
	return nil
}

func TestAdminUpdateRoleReplacesPermissionsWhenProvided(t *testing.T) {
	roleID := uuid.New()
	requestedPermissions := []string{"patients.read", "users.read"}

	replaceCalled := false
	listCalled := false

	store := stubAdminUsersStore{
		getRoleByIDFn: func(_ context.Context, id uuid.UUID) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR", Description: "Existing"}, nil
		},
		updateRoleFn: func(_ context.Context, id uuid.UUID, _, _ *string) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR", Description: "Existing"}, nil
		},
		replaceRolePermsFn: func(_ context.Context, id uuid.UUID, permissionNames []string) ([]Permission, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			replaceCalled = true
			if !reflect.DeepEqual(permissionNames, requestedPermissions) {
				t.Fatalf("permissions mismatch: got %#v want %#v", permissionNames, requestedPermissions)
			}
			return []Permission{{ID: uuid.New(), Name: "patients.read"}}, nil
		},
		listRolePermsFn: func(_ context.Context, _ uuid.UUID) ([]Permission, error) {
			listCalled = true
			return []Permission{{ID: uuid.New(), Name: "legacy.permission"}}, nil
		},
		listAllPermsFn: func(_ context.Context) ([]Permission, error) {
			return []Permission{
				{ID: uuid.New(), Name: "patients.read"},
				{ID: uuid.New(), Name: "users.read"},
			}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, noopAuditWriter{})
	resp, err := svc.AdminUpdateRole(context.Background(), roleID, AdminUpdateRoleRequest{
		Permissions: &requestedPermissions,
	}, DeviceInfo{})
	if err != nil {
		t.Fatalf("AdminUpdateRole error: %v", err)
	}

	if !replaceCalled {
		t.Fatal("expected ReplaceRolePermissions to be called")
	}
	if listCalled {
		t.Fatal("expected ListRolePermissions not to be called when permissions are provided")
	}
	if len(resp.Permissions) != 1 || resp.Permissions[0].Name != "patients.read" {
		t.Fatalf("returned permissions mismatch: %#v", resp.Permissions)
	}
}

func TestAdminUpdateRoleListsPermissionsWhenPermissionsOmitted(t *testing.T) {
	roleID := uuid.New()

	replaceCalled := false
	listCalled := false

	store := stubAdminUsersStore{
		getRoleByIDFn: func(_ context.Context, id uuid.UUID) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR"}, nil
		},
		updateRoleFn: func(_ context.Context, id uuid.UUID, _, _ *string) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR"}, nil
		},
		replaceRolePermsFn: func(_ context.Context, _ uuid.UUID, _ []string) ([]Permission, error) {
			replaceCalled = true
			return []Permission{}, nil
		},
		listRolePermsFn: func(_ context.Context, id uuid.UUID) ([]Permission, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			listCalled = true
			return []Permission{{ID: uuid.New(), Name: "roles.read"}}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, noopAuditWriter{})
	resp, err := svc.AdminUpdateRole(context.Background(), roleID, AdminUpdateRoleRequest{}, DeviceInfo{})
	if err != nil {
		t.Fatalf("AdminUpdateRole error: %v", err)
	}

	if replaceCalled {
		t.Fatal("expected ReplaceRolePermissions not to be called when permissions are omitted")
	}
	if !listCalled {
		t.Fatal("expected ListRolePermissions to be called when permissions are omitted")
	}
	if len(resp.Permissions) != 1 || resp.Permissions[0].Name != "roles.read" {
		t.Fatalf("returned permissions mismatch: %#v", resp.Permissions)
	}
}

func TestAdminUpdateRoleAllowsClearingPermissions(t *testing.T) {
	roleID := uuid.New()
	emptyPermissions := []string{}

	replaceCalled := false

	store := stubAdminUsersStore{
		getRoleByIDFn: func(_ context.Context, id uuid.UUID) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR"}, nil
		},
		updateRoleFn: func(_ context.Context, id uuid.UUID, _, _ *string) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "SUPERVISOR"}, nil
		},
		replaceRolePermsFn: func(_ context.Context, id uuid.UUID, permissionNames []string) ([]Permission, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			replaceCalled = true
			if len(permissionNames) != 0 {
				t.Fatalf("expected empty permissions list, got %#v", permissionNames)
			}
			return []Permission{}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, noopAuditWriter{})
	resp, err := svc.AdminUpdateRole(context.Background(), roleID, AdminUpdateRoleRequest{
		Permissions: &emptyPermissions,
	}, DeviceInfo{})
	if err != nil {
		t.Fatalf("AdminUpdateRole error: %v", err)
	}

	if !replaceCalled {
		t.Fatal("expected ReplaceRolePermissions to be called")
	}
	if resp.Permissions == nil {
		t.Fatal("expected permissions slice, got nil")
	}
	if len(resp.Permissions) != 0 {
		t.Fatalf("expected no permissions, got %#v", resp.Permissions)
	}
}

func TestAdminUpdateRoleRejectsUnknownPermissions(t *testing.T) {
	roleID := uuid.New()
	requestedPermissions := []string{"patients.read", "reference.read"}

	replaceCalled := false

	store := stubAdminUsersStore{
		getRoleByIDFn: func(_ context.Context, id uuid.UUID) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "DOCTOR"}, nil
		},
		updateRoleFn: func(_ context.Context, id uuid.UUID, _, _ *string) (Role, error) {
			if id != roleID {
				t.Fatalf("unexpected role id: got %s want %s", id, roleID)
			}
			return Role{ID: roleID, Name: "DOCTOR"}, nil
		},
		replaceRolePermsFn: func(_ context.Context, _ uuid.UUID, _ []string) ([]Permission, error) {
			replaceCalled = true
			return []Permission{}, nil
		},
		listAllPermsFn: func(_ context.Context) ([]Permission, error) {
			return []Permission{
				{ID: uuid.New(), Name: "patients.read"},
				{ID: uuid.New(), Name: "users.read"},
			}, nil
		},
	}

	svc := NewService(store, JWTConfig{}, noopAuditWriter{})
	_, err := svc.AdminUpdateRole(context.Background(), roleID, AdminUpdateRoleRequest{
		Permissions: &requestedPermissions,
	}, DeviceInfo{})
	if !errors.Is(err, ErrUnknownPermissions) {
		t.Fatalf("expected ErrUnknownPermissions, got %v", err)
	}
	if replaceCalled {
		t.Fatal("expected ReplaceRolePermissions not to be called for unknown permission names")
	}
}
