package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Mock utility functions for testing
	HashPasswordFunc = func(password string) (string, error) {
		return "hashed_" + password, nil
	}
	CheckHashPasswordFunc = func(hash, password string) bool {
		return hash == "hashed_"+password
	}
	GenerateTokensFunc = func(email string, id int64) (string, error) {
		return "mock-token", nil
	}

	// Mock user save and validation functions
	UserSaveFunc = func(u *models.User) error {
		return nil
	}
	ValidateUserFunc = func(u *models.User) error {
		if u.Email == "valid@example.com" && u.Password == "password" {
			return nil
		}
		return assert.AnError
	}
}

// Helper to perform HTTP requests
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Setup Gin router for user routes only
func setupUserRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/signup", signup)
	r.POST("/login", login)
	return r
}

func TestSignup(t *testing.T) {
	router := setupUserRouter()

	user := map[string]string{
		"Email":    "test@example.com",
		"Password": "password",
	}

	w := performRequest(router, "POST", "/signup", user)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User created successfully!", resp["Message"])
}

func TestLogin(t *testing.T) {
	router := setupUserRouter()

	// ✅ Valid login
	validUser := map[string]string{
		"Email":    "valid@example.com",
		"Password": "password",
	}
	w := performRequest(router, "POST", "/login", validUser)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Login successful!", resp["message"])
	assert.Equal(t, "mock-token", resp["token"])

	// ❌ Invalid login
	invalidUser := map[string]string{
		"Email":    "invalid@example.com",
		"Password": "wrong",
	}
	w = performRequest(router, "POST", "/login", invalidUser)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
