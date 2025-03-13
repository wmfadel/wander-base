package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse event", err))
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId
	err = h.service.CreateEvent(&event)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to create event", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
	})
}

func (h *EventHandler) GetEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}
	event, err := h.service.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusNotFound, models.NewESError("Failed to get event data", err))
		return
	}
	context.JSON(http.StatusOK, event)
}

func (h *EventHandler) GetEvents(context *gin.Context) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to get events data", err))
		return
	}
	context.JSON(http.StatusOK, events)
}

func (h *EventHandler) UpdateEvent(context *gin.Context) {
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
	err = h.service.UpdatePartially(eventId, patchEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to update event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func (h *EventHandler) DeleteEvent(context *gin.Context) {
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

	err := h.service.Delete(event.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to delete event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
