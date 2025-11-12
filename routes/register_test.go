package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRegisterRouter() *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userId", int64(1)) // mock auth middleware
		c.Next()
	})
	r.POST("/events/:id/register", registerForEvent)
	r.DELETE("/events/:id/register", cancelRegistration)
	return r
}

func TestRegisterForEvent_Success(t *testing.T) {
	router := setupRegisterRouter()

	models.GetEventByIDFunc = func(id int64) (*models.Event, error) {
		return &models.Event{ID: id, Name: "Mock Event"}, nil
	}
	RegisterEventFunc = func(e *models.Event, userID int64) error { return nil }

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/events/1/register", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Registration successful")
}

//func TestCancelRegistration_Success(t *testing.T) {
//	router := setupRegisterRouter()
//
//	models.GetEventByIDFunc = func(id int64) (*models.Event, error) {
//		return &models.Event{ID: 1, UserID: 1}, nil
//	}
//	models.CancelRegistrationFunc = func(e *models.Event, userId int64) error { return nil }
//
//	req, _ := http.NewRequest("DELETE", "/events/1/register", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Contains(t, w.Body.String(), "Cancelled Successfully")
//}
