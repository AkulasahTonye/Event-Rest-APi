package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.GET("/events/:id", getEvent)
	router.POST("/events", createEvent)
	router.PUT("/events/:id", updateEvent)
	router.DELETE("/events/:id", deleteEvent)
	return router
}

func TestGetEvents(t *testing.T) {
	// Setup
	router := setupTestRouter()

	// Mock the database function
	originalGetAllEvents := models.GetAllEventsFunc
	defer func() { models.GetAllEventsFunc = originalGetAllEvents }()

	models.GetAllEventsFunc = func() ([]models.Event, error) {
		return []models.Event{
			{ID: 1, Name: "Test Event", Description: "Test", Location: "Test", DateTime: time.Now(), UserID: 1},
		}, nil
	}

	// Create a request
	req, _ := http.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Event")
}

func TestCreateEvent(t *testing.T) {
	// Setup
	router := setupTestRouter()

	// Mock the database function
	originalSave := models.EventSaveFunc
	defer func() { models.EventSaveFunc = originalSave }()

	models.EventSaveFunc = func(e *models.Event) error {
		e.ID = 1
		return nil
	}

	// Create test event
	event := models.Event{
		Name:        "New Event",
		Description: "New Description",
		Location:    "New Location",
		DateTime:    time.Now(),
		UserID:      1,
	}
	jsonEvent, _ := json.Marshal(event)

	// Create a request
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(jsonEvent))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Event Created")
}
