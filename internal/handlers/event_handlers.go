package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/service"
)

type EventHandler struct {
	service      *service.EventService
	photoService *service.EventPhotoService
}

func NewEventHandler(service *service.EventService, photoService *service.EventPhotoService) *EventHandler {
	return &EventHandler{service: service, photoService: photoService}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event
	// Bind form data to event struct, excluding photos
	if err := c.ShouldBind(&event); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Could not parse event", err))
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewESError("User not authenticated", nil))
		return
	}
	event.UserID = userID.(int64)

	// Create event with photos
	if err := h.service.CreateEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to create event", err))
		return
	}

	// Fetch updated event with photo URLs (optional, if service doesn't populate Photos)
	updatedEvent, err := h.service.GetEventById(event.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to fetch created event", err))
		return
	}

	c.JSON(http.StatusCreated, updatedEvent)
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

	err = context.ShouldBindJSON(&patchEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Missing event details", err))
		return
	}

	err = h.service.UpdatePartially(eventId, patchEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to update event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func (h *EventHandler) AddPhotos(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse form", err))
		return
	}
	photos := form.File["photos"]

	if err := h.photoService.AddPhotos(eventID, photos); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to save event photos", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photos added"})
}

func (h *EventHandler) DeletePhotos(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse photo ID", err))
		return
	}
	type TempPhotos struct {
		Photos []string `json:"photos"`
	}

	var p TempPhotos
	err = c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Missing photo details", err))
		return
	}
	var photos []string
	photos = append(photos, p.Photos...)
	err = h.photoService.DeletEventPhotos(eventId, photos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to delete photo", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted"})

}

func (h *EventHandler) DeleteEvent(context *gin.Context) {

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}

	event, err := h.service.GetEventById(eventId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, models.NewESError("Could not find event", err))
		return
	}

	err = h.service.Delete(event.ID, event.Photos)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Failed to delete event", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
