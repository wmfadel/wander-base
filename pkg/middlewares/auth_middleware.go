package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/service"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type AuthMiddleware struct {
	service *service.EventService
}

func NewAuthMiddleware(service *service.EventService) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
	}
}

func (amw *AuthMiddleware) Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "UnAuthorized",
		})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("invalid token %v", err),
		})
		return
	}

	context.Set("userId", userId)
	context.Next()
}

func (amw *AuthMiddleware) AuthorizeForEventEdits(context *gin.Context) {

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, models.NewESError("Failed to parse event ID", err))
		return
	}

	event, err := amw.service.GetEventById(eventId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, models.NewESError("Could not find event", err))
		return
	}
	userId := context.GetInt64("userId")
	if event.UserID != userId {
		context.AbortWithStatusJSON(http.StatusUnauthorized, models.NewESError("Not authorized", err))
		return
	}

	context.Set("event", *event)
	context.Next()
}
