package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/models"
)

func getEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to get event data: %v", err)})
		return
	}
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusNotExtended, gin.H{"message": fmt.Sprintf("Failed to get event data: %v", err)})
		return
	}
	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to get events data: %v", err)})
		return
	}
	context.JSON(http.StatusOK, events)
}

func creatEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Could not parse event %v", err),
		})
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId
	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to create event: %v", err)})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
	})
}

func updateEvent(context *gin.Context) {
	eventId, _ := strconv.ParseInt(context.Param("id"), 10, 64)

	var updatedEvent models.Event
	err := context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Missing event details: %v", err)})
		return
	}
	updatedEvent.ID = eventId
	err = updatedEvent.Update()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to update event: %v", err)})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func deleteEvent(context *gin.Context) {

	var event models.Event

	// Get event from context using the key "event"
	value, exists := context.Get("event")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Event not found in context"})
		return
	}

	// Type assert the value to models.Event
	var ok bool
	if event, ok = value.(models.Event); !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid event data in context"})
		return
	}

	err := event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to delete event: %v", err)})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
