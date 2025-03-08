package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/models"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64(("userId"))
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("cannot parse eventId %v", err),
		})
		return
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("event not found %v", err),
		})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("cannot register in event: %v", err),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Registration Successfull",
	})
}

func cancelRegistrationEvent(context *gin.Context) {
	userId := context.GetInt64(("userId"))
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("cannot parse eventId %v", err),
		})
		return
	}

	event := models.Event{ID: eventId}

	err = event.CancelRegister(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("cancelling registration failed: %v", err),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Unregistration Successfully",
	})
}
