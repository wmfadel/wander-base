package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/core"
	"github.com/wmfadel/wander-base/internal/service"
)

type RegistrationHandler struct {
	service *service.RegistrationService
}

func NewRegistrationHandler(service *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{service: service}
}

func (h *RegistrationHandler) RegisterForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Cannot parse event ID", err))
		return
	}

	err = h.service.Register(userId, eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Cannot register for event", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Requested to register for event, please wait for approval",
	})
}

func (h *RegistrationHandler) CancelRegistrationEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Cannot parse event ID", err))
		return
	}

	event := models.Event{ID: eventId}
	err = h.service.CancelRegister(userId, event.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Cancelling registration failed", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Unregistration Successful",
	})
}
