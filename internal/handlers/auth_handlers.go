package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/core"
	"github.com/wmfadel/wander-base/internal/models/requests"
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
		context.JSON(http.StatusBadRequest, core.NewESError("Could not parse user", err))
		return
	}

	err = h.service.Create(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Could not save user", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (h *AuthHandler) LoginHandler(context *gin.Context) {
	var loginRequest requests.LoginRequest
	err := context.ShouldBindJSON(&loginRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, core.NewESError("Could not parse user", err))
		return
	}

	err = h.service.ValidateCredintials(&loginRequest)
	if err != nil {
		context.JSON(http.StatusUnauthorized, core.NewESError("Cannot verify identity", err))
		return
	}

	token, err := utils.GernerateToken(loginRequest.Phone, loginRequest.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, core.NewESError("Cannot create token", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User Validated",
		"token":   token,
	})
}

func (h *AuthHandler) LogoutHandler(context *gin.Context) {
	context.JSON(http.StatusNotExtended, core.NewESError("Not implemented", nil))
}
