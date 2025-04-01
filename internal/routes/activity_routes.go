package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/di"
)

func RegisterActivityRoutes(r *gin.Engine, c di.DIContainer) {
	guarded := r.Group("/activity", c.AuthMiddleware.Authenticate, c.AuthMiddleware.RequiresAdmin)
	// Public event routes
	guarded.GET("/:id", c.ActivityHandler.GetActivity)
	guarded.GET("/slug/:slug", c.ActivityHandler.GetActivityBySlug)
	guarded.GET("/", c.ActivityHandler.GetActivities)
	guarded.POST("/", c.ActivityHandler.CreateActivity)
	guarded.DELETE("/:id", c.ActivityHandler.DeleteActivity)
}
