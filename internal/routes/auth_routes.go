package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterAuthRoutes(r *gin.Engine, c di.DIContainer) {
	r.POST("/signup", c.AuthHandler.SignupHandler)
	r.POST("/login", c.AuthHandler.LoginHandler)
	r.POST("/logout", c.AuthMiddleware.Authenticate, c.AuthHandler.LogoutHandler)
}
