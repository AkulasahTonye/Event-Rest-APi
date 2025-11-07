package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api-testing/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Token", func(t *testing.T) {
		// Backup the original function
		originalFunc := verifyTokenFunc
		defer func() { verifyTokenFunc = originalFunc }()

		// Mock VerifyToken to succeed
		verifyTokenFunc = func(token string) (int64, error) {
			return 42, nil
		}

		// Setup Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "valid-token")
		c.Request = req

		router := gin.New()
		router.Use(Authenticate)
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ok")
	})

	t.Run("Missing Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		c.Request = req

		router := gin.New()
		router.Use(Authenticate)
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Not Authorized")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		originalFunc := utils.VerifyToken
		defer func() { verifyTokenFunc = originalFunc }()

		// Mock VerifyToken to fail
		verifyTokenFunc = func(token string) (int64, error) {
			return 0, errors.New("invalid token")
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "bad-token")
		c.Request = req

		router := gin.New()
		router.Use(Authenticate)
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Not Authorized")
	})
}
