package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/di"
)

func RegisterDestinationRoutes(r *gin.Engine, c di.DIContainer) {
	guarded := r.Group("/destination", c.AuthMiddleware.Authenticate, c.AuthMiddleware.RequiresAdmin)
	// Public event routes
	guarded.GET("/:id", c.DestinationHandler.GetDestinationById)
	guarded.GET("/", c.DestinationHandler.GetAllDestinations)
	guarded.POST("/", c.DestinationHandler.CreateDestination)
	guarded.DELETE("/:id", c.DestinationHandler.DeleteDestination)
}
