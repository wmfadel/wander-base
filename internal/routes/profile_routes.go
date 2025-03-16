package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterProfileRoutes(r *gin.Engine, c di.DIContainer) {
	r.Group("/users", c.AuthMiddleware.Authenticate)
	// Public event routes
	r.GET("/", c.ProfileHandler.GetProfile)        // Get Profile // TODO add correct handler
	r.PUT("/", c.ProfileHandler.UpdateProfile)     // Update Profile // TODO add correct handler
	r.POST("/photo", c.ProfileHandler.UpdatePhoto) // Update Profile
}
