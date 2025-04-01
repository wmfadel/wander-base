package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/core"
	"github.com/wmfadel/wander-base/internal/service"
)

type ActivityHandler struct {
	service *service.ActivityService
}

func NewActivityHandler(service *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) CreateActivity(context *gin.Context) {
	var activity models.Activity
	err := context.ShouldBindJSON(&activity)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Could not parse activity", err))
		return
	}

	err = h.service.Save(&activity)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Could not save activity", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Activity created",
	})
}

func (h *ActivityHandler) GetActivities(context *gin.Context) {
	activities, err := h.service.GetAllActivities()
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Failed to get activities data", err))
		return
	}
	context.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetActivity(context *gin.Context) {
	activityID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Failed to parse activity ID", err))
		return
	}
	activity, err := h.service.GetActivityById(activityID)
	if err != nil {
		context.JSON(http.StatusNotFound, core.NewESError("Failed to get activity data", err))
		return
	}
	context.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) GetActivityBySlug(context *gin.Context) {
	slug := context.Param("slug")
	activity, err := h.service.GetActivityBySlug(slug)
	if err != nil {
		context.JSON(http.StatusNotFound, core.NewESError("Failed to get activity data", err))
		return
	}
	context.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) DeleteActivity(context *gin.Context) {

	activityId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, core.NewESError("Failed to parse activity ID", err))
		return
	}

	err = h.service.Delete(activityId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Failed to delete activity", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
}
