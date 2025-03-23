package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/service"
)

type AdmingHandler struct {
	UserService *service.UserService
}

func NewAdmingHandler(service *service.UserService) *AdmingHandler {
	return &AdmingHandler{UserService: service}
}

func (h *AdmingHandler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) BlockUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) AddRoleToUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) RemoveRoleFromUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) GetAllAdmins(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}

func (h *AdmingHandler) GetAllOrganizers(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.NewESError("This action is not supported yet", nil))
}
