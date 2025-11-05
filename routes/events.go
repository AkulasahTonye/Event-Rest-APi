package routes

import (
	"net/http"
	"strconv"

	"example.com/api-testing/models"
	"github.com/gin-gonic/gin"
)

func getEvents(ctx *gin.Context) {

	events, err := models.GetAllEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch events. Try again later. "})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func getEvent(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event id. "})
		return
	}
	event, err := models.GetEventByID(eventId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch event. "})
		return
	}
	ctx.JSON(http.StatusOK, event)
}

func createEvent(ctx *gin.Context) {

	var event models.Event
	err := ctx.ShouldBindJSON(&event)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message": "Could not parse request data",
			"Error":   err.Error(),
		})
		return
	}

	userId := ctx.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not create event. Try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"Message": "Event Created", "Event": event})
}

func updateEvent(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Could not Parse event id."})
		return
	}
	userId := ctx.GetInt64("userId")
	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Message": "Not Authorized to update event."})
		return
	}

	var updatedEvent models.Event

	err = ctx.ShouldBindJSON(&updatedEvent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"Message": "Could not parse request data",
				"Error": err.Error()})
		return
	}
	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not update event.", "Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Event Updated Successfully!"})
}

func deleteEvent(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Could not Parse event id."})
		return
	}
	userId := ctx.GetInt64("userId")
	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Message": "Not Authorized to Delete event.", "Error": err.Error()})
		return
	}

	err = event.Delete()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not delete  event. ", "Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Event deleted Successfully!"})
}
