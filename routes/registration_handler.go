package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/models"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Cannot parse event ID", err))
		return
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, models.NewESError("Event not found", err))
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cannot register for event", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Registration Successful",
	})
}

func cancelRegistrationEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Cannot parse event ID", err))
		return
	}

	event := models.Event{ID: eventId}
	err = event.CancelRegister(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cancelling registration failed", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Unregistration Successful",
	})
}
