package data

import (
	"testing"
	"time"
)

func TestTokenModel_New(t *testing.T) {
	userModel := UserModel{DB: testDB}
	tokenModel := TokenModel{DB: testDB}
	user := insertTestUser(t, userModel)

	token, err := tokenModel.New(user.Id, time.Hour, ScopeAuthentication)
	t.Cleanup(func() {
		tokenModel.DeleteAllForUser(ScopeAuthentication, user.Id)
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if len(token.Plaintext) != 26 {
		t.Errorf("New() Plaintext length = %d, want 26", len(token.Plaintext))
	}
	if len(token.Hash) != 32 {
		t.Errorf("New() Hash length = %d, want 32", len(token.Hash))
	}
	if token.UserID != user.Id {
		t.Errorf("New() UserID = %d, want %d", token.UserID, user.Id)
	}
	if token.Scope != ScopeAuthentication {
		t.Errorf("New() Scope = %q, want %q", token.Scope, ScopeAuthentication)
	}
	if !token.Expiry.After(time.Now()) {
		t.Error("New() Expiry should be in the future")
	}
}

func TestUserModel_GetForToken(t *testing.T) {
	userModel := UserModel{DB: testDB}
	tokenModel := TokenModel{DB: testDB}
	user := insertTestUser(t, userModel)

	token, err := tokenModel.New(user.Id, time.Hour, ScopeAuthentication)
	t.Cleanup(func() {
		tokenModel.DeleteAllForUser(ScopeAuthentication, user.Id)
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	fetched, err := userModel.GetForToken(ScopeAuthentication, token.Plaintext)
	if err != nil {
		t.Fatalf("GetForToken() error = %v", err)
	}
	if fetched.Id != user.Id {
		t.Errorf("GetForToken() user ID = %d, want %d", fetched.Id, user.Id)
	}
}

func TestUserModel_GetForToken_NotFound(t *testing.T) {
	userModel := UserModel{DB: testDB}

	_, err := userModel.GetForToken(ScopeAuthentication, "AAAAAAAAAAAAAAAAAAAAAAAAAA")
	if err != ErrRecordNotFound {
		t.Errorf("GetForToken() error = %v, want ErrRecordNotFound", err)
	}
}

func TestTokenModel_DeleteAllForUser(t *testing.T) {
	userModel := UserModel{DB: testDB}
	tokenModel := TokenModel{DB: testDB}
	user := insertTestUser(t, userModel)

	token, err := tokenModel.New(user.Id, time.Hour, ScopeAuthentication)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := tokenModel.DeleteAllForUser(ScopeAuthentication, user.Id); err != nil {
		t.Fatalf("DeleteAllForUser() error = %v", err)
	}

	_, err = userModel.GetForToken(ScopeAuthentication, token.Plaintext)
	if err != ErrRecordNotFound {
		t.Errorf("GetForToken() after DeleteAllForUser() error = %v, want ErrRecordNotFound", err)
	}
}
