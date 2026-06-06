package data

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func uniqueTestEmail(t *testing.T) string {
	t.Helper()

	name := strings.ToLower(t.Name())
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, name)

	return fmt.Sprintf("%s-%d@example.com", name, time.Now().UnixNano())
}

func newTestUser(t *testing.T) *User {
	t.Helper()

	user := &User{
		Photo:     "test-user.jpg",
		Name:      "Test User",
		Email:     uniqueTestEmail(t),
		Activated: true,
		Role:      "customer",
	}

	if err := user.Password.Set("password123"); err != nil {
		t.Fatalf("failed to set password: %v", err)
	}

	return user
}

func insertTestUser(t *testing.T, model UserModel) *User {
	t.Helper()

	user := newTestUser(t)
	if err := model.Insert(user); err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM users WHERE id = $1`, user.Id)
	})

	return user
}

func TestUserModel_Insert(t *testing.T) {
	model := UserModel{DB: testDB}
	user := newTestUser(t)

	err := model.Insert(user)
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM users WHERE id = $1`, user.Id)
	})

	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if user.Id == 0 {
		t.Error("Insert() did not set user.Id")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Insert() did not set CreatedAt")
	}
	if user.Version != 1 {
		t.Errorf("Insert() Version = %d, want 1", user.Version)
	}
}

func TestUserModel_Insert_DuplicateEmail(t *testing.T) {
	model := UserModel{DB: testDB}
	user := insertTestUser(t, model)

	duplicate := newTestUser(t)
	duplicate.Email = user.Email

	err := model.Insert(duplicate)
	if err != ErrDuplicateEmail {
		t.Errorf("Insert() error = %v, want ErrDuplicateEmail", err)
	}
}

func TestUserModel_GetByEmail(t *testing.T) {
	model := UserModel{DB: testDB}
	user := insertTestUser(t, model)

	fetched, err := model.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("GetByEmail() error = %v", err)
	}

	if fetched.Id != user.Id {
		t.Errorf("GetByEmail() Id = %d, want %d", fetched.Id, user.Id)
	}
	if fetched.Email != user.Email {
		t.Errorf("GetByEmail() Email = %q, want %q", fetched.Email, user.Email)
	}
	matches, err := fetched.Password.Matches("password123")
	if err != nil {
		t.Fatalf("Password.Matches() error = %v", err)
	}
	if !matches {
		t.Error("stored password hash did not match plaintext password")
	}
}

func TestUserModel_GetByEmail_NotFound(t *testing.T) {
	model := UserModel{DB: testDB}

	_, err := model.GetByEmail(uniqueTestEmail(t))
	if err != ErrRecordNotFound {
		t.Errorf("GetByEmail() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUserModel_Get(t *testing.T) {
	model := UserModel{DB: testDB}
	user := insertTestUser(t, model)

	fetched, err := model.Get(user.Id)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if fetched.Id != user.Id {
		t.Errorf("Get() Id = %d, want %d", fetched.Id, user.Id)
	}
	if fetched.Name != user.Name {
		t.Errorf("Get() Name = %q, want %q", fetched.Name, user.Name)
	}
}

func TestUserModel_Get_NotFound(t *testing.T) {
	model := UserModel{DB: testDB}

	_, err := model.Get(999999)
	if err != ErrRecordNotFound {
		t.Errorf("Get() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUserModel_Update(t *testing.T) {
	model := UserModel{DB: testDB}
	user := insertTestUser(t, model)

	user.Name = "Updated User"
	user.Email = uniqueTestEmail(t)
	user.Activated = false
	user.Role = "admin"
	if err := user.Password.Set("newpassword123"); err != nil {
		t.Fatalf("failed to set password: %v", err)
	}

	if err := model.Update(user); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	fetched, err := model.Get(user.Id)
	if err != nil {
		t.Fatalf("Get() after Update() error = %v", err)
	}

	if fetched.Name != "Updated User" {
		t.Errorf("Update() Name = %q, want %q", fetched.Name, "Updated User")
	}
	if fetched.Activated {
		t.Error("Update() Activated should be false")
	}
	if fetched.Role != "admin" {
		t.Errorf("Update() Role = %q, want admin", fetched.Role)
	}
	matches, err := fetched.Password.Matches("newpassword123")
	if err != nil {
		t.Fatalf("Password.Matches() error = %v", err)
	}
	if !matches {
		t.Error("updated password hash did not match plaintext password")
	}
}

func TestUserModel_Update_EditConflict(t *testing.T) {
	model := UserModel{DB: testDB}
	user := insertTestUser(t, model)

	stale := *user

	user.Name = "Current Version"
	if err := model.Update(user); err != nil {
		t.Fatalf("Update() current user error = %v", err)
	}

	stale.Name = "Stale Version"
	err := model.Update(&stale)
	if err != ErrEditConflict {
		t.Errorf("Update() stale user error = %v, want ErrEditConflict", err)
	}
}

func TestUserModel_Update_DuplicateEmail(t *testing.T) {
	model := UserModel{DB: testDB}
	first := insertTestUser(t, model)
	second := insertTestUser(t, model)

	second.Email = first.Email
	err := model.Update(second)
	if err != ErrDuplicateEmail {
		t.Errorf("Update() error = %v, want ErrDuplicateEmail", err)
	}
}
