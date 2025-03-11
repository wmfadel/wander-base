package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/models" // Adjust to your actual import path (e.g., "github.com/wmfadel/go_events/models")
)

func getEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusNotFound, models.NewESError("Failed to get event data", err))
		return
	}
	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to get events data", err))
		return
	}
	context.JSON(http.StatusOK, events)
}

func creatEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse event", err))
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to create event", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
	})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}

	var patchEvent models.PatchEvent
	var updatedEvent models.Event

	err = context.ShouldBindJSON(&patchEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Missing event details", err))
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.UpdatePartially(patchEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to update event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func deleteEvent(context *gin.Context) {
	// Get event from context using the key "event"
	value, exists := context.Get("event")
	if !exists {
		context.JSON(http.StatusInternalServerError, models.NewESError("Event not found in context", nil))
		return
	}

	// Type assert the value to models.Event
	event, ok := value.(models.Event)
	if !ok {
		context.JSON(http.StatusInternalServerError, models.NewESError("Invalid event data in context", nil))
		return
	}

	err := event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to delete event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
