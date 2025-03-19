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
	user, err := utils.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get user from context", err))
		return
	}
	var patch models.PatchUser
	err = c.ShouldBindJSON(&patch)

	if err != nil || patch.IsEmpty() {
		c.JSON(http.StatusBadRequest, models.NewESError("Invalid request body", err))
		return
	}

	err = h.UserService.UpdateUser(user, &patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to update user", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})

}

func (h *ProfileHandler) UpdatePhoto(c *gin.Context) {
	userId := c.GetInt64("userId") // From auth middleware
	photo, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("photo required", err))
		return
	}
	url, err := h.UserService.UpdatePhoto(userId, photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to update photo", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photo updated", "url": url})
}
