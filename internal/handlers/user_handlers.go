package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/service"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SignupHanlder(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse user", err))
		return
	}

	err = h.service.Save(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Could not save user", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (h *UserHandler) LoginHandler(context *gin.Context) {
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
