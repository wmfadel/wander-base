package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterRoutes(server *gin.Engine, c di.DIContainer) {
	// Public user routes
	server.POST("/signup", c.UserHandler.SignupHandler) // Fixed typo
	server.POST("/login", c.UserHandler.LoginHandler)

	// Guarded user routes (require authentication)
	guarded := server.Group("/", c.AuthMiddleware.Authenticate)
	guarded.POST("/users/photo", c.UserHandler.UpdatePhoto)

	// Event routes (delegated to separate function)
	RegisterEventRoutes(server, c, c.AuthMiddleware)
}
