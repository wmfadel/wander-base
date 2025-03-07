package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/models"
	"githuv.com/wmfadel/go_events/utils"
)

func signupHanlder(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Could not parse user %v", err),
		})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Could not save user %v", err),
		})
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
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Could not parse user %v", err),
		})
		return
	}

	err = user.ValidateCredintials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("cannot verify identity %v", err),
		})
		return
	}

	token, err := utils.GernerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("cannot create token %v", err),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User Validated",
		"token":   token,
	})

}
