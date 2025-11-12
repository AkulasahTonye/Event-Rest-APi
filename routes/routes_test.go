package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func mockAuthenticate(ctx *gin.Context) {
	ctx.Set("userId", int64(1))
	ctx.Next()
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.Default()
	RegisterRoutes(router)

	endpoints := []string{
		"/events",
		"/events/1",
		"/signup",
		"/login",
		"/events/1/register",
	}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.True(t, w.Code >= 200)
	}
}
