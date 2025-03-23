package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/service"
)

type RegistrationHandler struct {
	service *service.EventService
}

func NewRegistrationHandler(service *service.EventService) *RegistrationHandler {
	return &RegistrationHandler{service: service}
}

func (h *RegistrationHandler) RegisterForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Cannot parse event ID", err))
		return
	}

	event, err := h.service.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, models.NewESError("Event not found", err))
		return
	}

	err = h.service.Register(userId, event.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cannot register for event", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Registration Successful",
	})
}

func (h *RegistrationHandler) CancelRegistrationEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Cannot parse event ID", err))
		return
	}

	event := models.Event{ID: eventId}
	err = h.service.CancelRegister(userId, event.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cancelling registration failed", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Unregistration Successful",
	})
}
