package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/service"
	"github.com/wmfadel/wander-base/pkg/utils"
)

type AuthMiddleware struct {
	userService  *service.UserService
	eventService *service.EventService
}

func NewAuthMiddleware(userService *service.UserService, eventService *service.EventService) *AuthMiddleware {
	return &AuthMiddleware{
		userService:  userService,
		eventService: eventService,
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
		context.AbortWithStatusJSON(http.StatusUnauthorized, models.NewESError("Invalid token", err))
		return
	}

	if userId == 0 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, models.NewESError("user not found", err))
		return
	}

	user, err := amw.userService.GetUserByID(userId)

	if err != nil || user == nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, models.NewESError("user not found", err))
		return
	}
	log.Printf("user in middleware: %v", user)
	context.Set("userId", userId)
	context.Set("user", *user)
	context.Next()
}

func (amw *AuthMiddleware) RequiresAdmin(context *gin.Context) {

	user, err := utils.GetUserFromContext(context)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, models.NewESError("Could not find user", err))
		return
	}

	isAdmin := false
	for _, role := range user.Roles {
		if role.ID == 1 {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		context.AbortWithStatusJSON(http.StatusUnauthorized, models.NewESError("Unauthorized to create/edit events", nil))
		return
	}
	context.Next()
}
