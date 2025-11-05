package routes

import (
	"net/http"
	"strconv"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.Register(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Registered!"})

}

func cancelRegistration(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}
