package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/service"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type ProfileHandler struct {
	UserService *service.UserService
}

func NewProfileHandler(service *service.UserService) *ProfileHandler {
	return &ProfileHandler{UserService: service}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get user from context", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get user from context", err))
		return
	}

	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))

}

func (h *ProfileHandler) UpdatePhoto(c *gin.Context) {
	userId := c.GetInt64("userId") // From auth middleware
	photo, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo required"})
		return
	}
	url, err := h.UserService.UpdatePhoto(userId, photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photo updated", "url": url})
}
