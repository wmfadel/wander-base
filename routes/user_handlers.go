package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/models"
	"github.com/wmfadel/escape-be/utils"
)

func signupHanlder(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse user", err))
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Could not save user", err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func loginHandler(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.NewESError("Could not parse user", err))
		return
	}

	err = user.ValidateCredintials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, models.NewESError("Cannot verify identity", err))
		return
	}

	token, err := utils.GernerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.NewESError("Cannot create token", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User Validated",
		"token":   token,
	})
}
