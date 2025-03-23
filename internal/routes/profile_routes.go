package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/di"
)

func RegisterProfileRoutes(r *gin.Engine, c di.DIContainer) {
	guarded := r.Group("/", c.AuthMiddleware.Authenticate)
	// Public event routes
	guarded.GET("/users", c.ProfileHandler.GetProfile)    // Get Profile
	guarded.PUT("/users", c.ProfileHandler.UpdateProfile) // Update Profile
	guarded.POST("/photo", c.ProfileHandler.UpdatePhoto)  // Update Profile
}
