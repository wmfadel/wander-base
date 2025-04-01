package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/core"
	"github.com/wmfadel/wander-base/internal/service"
)

type DestinationHandler struct {
	service *service.DestinationService
}

func NewDestinationHandler(service *service.DestinationService) *DestinationHandler {
	return &DestinationHandler{service: service}
}

func (h *DestinationHandler) CreateDestination(context *gin.Context) {
	var destination models.Destination
	err := context.ShouldBindJSON(&destination)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Could not parse destination", err))
		return
	}

	err = h.service.Save(&destination)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Could not save destination", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Destination created",
	})
}

func (h *DestinationHandler) GetAllDestinations(context *gin.Context) {
	destinations, err := h.service.GetAllDestinations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Failed to get destinations data", err))
		return
	}
	context.JSON(http.StatusOK, destinations)
}

func (h *DestinationHandler) GetDestinationById(context *gin.Context) {
	destinationID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Failed to parse destination ID", err))
		return
	}
	destination, err := h.service.GetDestinationById(destinationID)
	if err != nil {
		context.JSON(http.StatusNotFound, core.NewESError("Failed to get destination data", err))
		return
	}
	context.JSON(http.StatusOK, destination)
}

func (h *DestinationHandler) DeleteDestination(context *gin.Context) {

	destinationId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, core.NewESError("Failed to parse destination ID", err))
		return
	}

	err = h.service.DeleteDestination(destinationId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Failed to delete destination", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Destination deleted"})
}
