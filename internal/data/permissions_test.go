package data

import "testing"

func TestPermissions_Include(t *testing.T) {
	permissions := Permissions{"dishes:read", "orders:write"}

	if !permissions.Include("dishes:read") {
		t.Error("Include() returned false for existing permission")
	}
	if permissions.Include("restaurants:write") {
		t.Error("Include() returned true for missing permission")
	}
}

func TestPermissionModel_AddGetDeleteForUser(t *testing.T) {
	userModel := UserModel{DB: testDB}
	permissionModel := PermissionModel{DB: testDB}
	user := insertTestUser(t, userModel)

	if err := permissionModel.AddForUser(user.Id, "dishes:read", "orders:write"); err != nil {
		t.Fatalf("AddForUser() error = %v", err)
	}

	permissions, err := permissionModel.GetAllForUser(user.Id)
	if err != nil {
		t.Fatalf("GetAllForUser() error = %v", err)
	}
	if !permissions.Include("dishes:read") {
		t.Error("GetAllForUser() did not include dishes:read")
	}
	if !permissions.Include("orders:write") {
		t.Error("GetAllForUser() did not include orders:write")
	}

	if err := permissionModel.DeleteForUser(user.Id, "dishes:read"); err != nil {
		t.Fatalf("DeleteForUser() error = %v", err)
	}

	permissions, err = permissionModel.GetAllForUser(user.Id)
	if err != nil {
		t.Fatalf("GetAllForUser() after DeleteForUser() error = %v", err)
	}
	if permissions.Include("dishes:read") {
		t.Error("GetAllForUser() still included deleted permission dishes:read")
	}
	if !permissions.Include("orders:write") {
		t.Error("GetAllForUser() should still include orders:write")
	}
}

func TestPermissionModel_GetAllForUser_NoPermissions(t *testing.T) {
	userModel := UserModel{DB: testDB}
	permissionModel := PermissionModel{DB: testDB}
	user := insertTestUser(t, userModel)

	permissions, err := permissionModel.GetAllForUser(user.Id)
	if err != nil {
		t.Fatalf("GetAllForUser() error = %v", err)
	}
	if len(permissions) != 0 {
		t.Errorf("GetAllForUser() returned %d permissions, want 0", len(permissions))
	}
}
