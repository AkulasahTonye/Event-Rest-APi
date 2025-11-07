package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mysecret123"

	hashed, err := HashPassword(password)

	// Check for no error
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed, "hashed password should not be empty")
	assert.NotEqual(t, password, hashed, "hashed password should not equal the plain password")

	// The hashed string should typically start with bcrypt prefix "$2"
	assert.Contains(t, hashed, "$2", "bcrypt hash should contain $2 prefix")
}

func TestCheckHashPassword(t *testing.T) {
	password := "supersecurepassword"

	// Generate a hash
	hashed, err := HashPassword(password)
	assert.NoError(t, err)

	//   Should return true for the correct password
	isValid := CheckHashPassword(password, hashed)
	assert.True(t, isValid, "expected password to match hash")

	// ❌ Should return false for an incorrect password
	isInvalid := CheckHashPassword("wrongpassword", hashed)
	assert.False(t, isInvalid, "expected password not to match hash")
}
