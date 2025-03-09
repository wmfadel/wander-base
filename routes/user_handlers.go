package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/models" // Adjust to your actual import path (e.g., "github.com/wmfadel/go_events/models")
	"githuv.com/wmfadel/go_events/utils"  // Adjust to your actual import path
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
