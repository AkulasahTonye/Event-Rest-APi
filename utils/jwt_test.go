package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	email := "test@example.com"
	userID := int64(12345)

	// Generate a valid token
	token, err := GenerateTokens(email, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, strings.HasPrefix(token, "ey"), "JWT should start with 'ey'")

	// Verify the token
	returnedID, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, returnedID)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	// A deliberately broken JWT string
	invalidToken := "invalid.token.string"

	userID, err := VerifyToken(invalidToken)
	assert.Error(t, err)
	assert.Equal(t, int64(0), userID)
	assert.Contains(t, err.Error(), "Could not parse token")
}

func TestVerifyToken_TamperedSignature(t *testing.T) {
	email := "tamper@example.com"
	userID := int64(1)

	validToken, err := GenerateTokens(email, userID)
	assert.NoError(t, err)

	// Tamper with the token by changing one character
	tamperedToken := validToken[:len(validToken)-2] + "xx"

	userIDResult, err := VerifyToken(tamperedToken)
	assert.Error(t, err)
	assert.Equal(t, int64(0), userIDResult)
}
