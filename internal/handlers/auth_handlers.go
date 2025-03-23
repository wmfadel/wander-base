package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/service"
	"github.com/wmfadel/wander-base/pkg/utils"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SignupHandler(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse user", err))
		return
	}

	err = h.service.Create(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Could not save user", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (h *AuthHandler) LoginHandler(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse user", err))
		return
	}

	err = h.service.ValidateCredintials(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, models.NewESError("Cannot verify identity", err))
		return
	}

	token, err := utils.GernerateToken(user.Phone, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cannot create token", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User Validated",
		"token":   token,
	})
}

func (h *AuthHandler) LogoutHandler(context *gin.Context) {
	context.JSON(http.StatusNotExtended, models.NewESError("Not implemented", nil))
}
