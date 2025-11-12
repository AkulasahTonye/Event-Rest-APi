package routes

import (
	"net/http"
	"strconv"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
)

// -------------------------------
// 🔹 Mockable Function Hooks
// -------------------------------

// These can be overridden in tests to prevent touching the real DB
var (
	GetEventByIDForRegisterFunc = models.GetEventByIDFunc
	RegisterEventFunc           = func(e *models.Event, userID int64) error { return e.Register(userID) }
	CancelRegistrationFunc      = func(e *models.Event, userID int64) error { return e.CancelRegistration(userID) }
)

// -------------------------------
// 🔹 Route Handlers
// -------------------------------

// registerForEvent handles POST /events/:id/register
func registerForEvent(ctx *gin.Context) {
	eventID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid event ID"})
		return
	}

	userID := ctx.GetInt64("userId")
	event, err := GetEventByIDForRegisterFunc(eventID)
	if err != nil || event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Message": "Event not found"})
		return
	}

	if err := RegisterEventFunc(event, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not register for event", "Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"Message": "Registration successful"})
}

// cancelRegistration handles DELETE /events/:id/register

func cancelRegistration(ctx *gin.Context) {
	eventID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid event ID"})
		return
	}

	userID := ctx.GetInt64("userId")
	event, err := GetEventByIDForRegisterFunc(eventID)
	if err != nil || event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Message": "Event not found"})
		return
	}

	if err := CancelRegistrationFunc(event, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not cancel registration", "Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "Registration canceled successfully"})
}
