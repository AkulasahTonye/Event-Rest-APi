package models

import (
	"database/sql"
	"testing"

	"example.com/api-testing/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// --- helper to create mock DB ---
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql db: %v", err)
	}
	db.DB = mockDB
	return mock, func() {
		mockDB.Close()
	}
}

// --- Test: Save (Success) ---
func TestUser_Save_Success(t *testing.T) {
	mock, closeFn := setupMockDB(t)
	defer closeFn()

	u := User{
		Email:    "test@example.com",
		Password: "password123",
	}

	mock.ExpectPrepare(`INSERT INTO users`).
		ExpectExec().
		WithArgs(u.Email, "hashedpass").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// inject a fake hash function
	hashFunc := func(p string) (string, error) {
		return "hashedpass", nil
	}

	err := u.Save(hashFunc)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
}

// --- Test: Save (DB error) ---
func TestUser_Save_DBError(t *testing.T) {
	mock, closeFn := setupMockDB(t)
	defer closeFn()

	u := User{
		Email:    "fail@example.com",
		Password: "secret",
	}

	mock.ExpectPrepare(`INSERT INTO users`).
		ExpectExec().
		WithArgs(u.Email, "hashedpass").
		WillReturnError(sql.ErrConnDone)

	hashFunc := func(p string) (string, error) {
		return "hashedpass", nil
	}

	err := u.Save(hashFunc)
	assert.Error(t, err)
}

// --- Test: ValidateCredentials (success) ---
func TestUser_ValidateCredentials_Success(t *testing.T) {
	mock, closeFn := setupMockDB(t)
	defer closeFn()

	u := &User{
		Email:    "test@example.com",
		Password: "password123",
	}

	rows := sqlmock.NewRows([]string{"id", "password"}).
		AddRow(1, "hashedpass")

	mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
		WithArgs(u.Email).
		WillReturnRows(rows)

	checkFunc := func(password, hashed string) bool {
		return true
	}

	err := u.ValidateCredentials(checkFunc)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
}

// --- Test: ValidateCredentials (invalid password) ---
func TestUser_ValidateCredentials_InvalidPassword(t *testing.T) {
	mock, closeFn := setupMockDB(t)
	defer closeFn()

	u := &User{
		Email:    "test@example.com",
		Password: "wrongpass",
	}

	rows := sqlmock.NewRows([]string{"id", "password"}).
		AddRow(1, "hashedpass")

	mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
		WithArgs(u.Email).
		WillReturnRows(rows)

	checkFunc := func(password, hashed string) bool {
		return false
	}

	err := u.ValidateCredentials(checkFunc)
	assert.Error(t, err)
	assert.EqualError(t, err, "Credentials invalid")
}

// --- Test: ValidateCredentials (no user found) ---
func TestUser_ValidateCredentials_NoUser(t *testing.T) {
	mock, closeFn := setupMockDB(t)
	defer closeFn()

	u := &User{
		Email:    "nouser@example.com",
		Password: "password123",
	}

	mock.ExpectQuery(`SELECT id, password FROM users WHERE email = \?`).
		WithArgs(u.Email).
		WillReturnError(sql.ErrNoRows)

	checkFunc := func(password, hashed string) bool {
		return true
	}

	err := u.ValidateCredentials(checkFunc)
	assert.Error(t, err)
	assert.EqualError(t, err, "Credentials invalid")
}
