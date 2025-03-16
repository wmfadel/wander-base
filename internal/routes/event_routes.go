package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterEventRoutes(r *gin.Engine, c di.DIContainer) {
	// Public event routes
	r.GET("/events", c.EventHandler.GetEvents)
	r.GET("/events/:id", c.EventHandler.GetEvent)

	// Guarded event routes (require authentication)
	guarded := r.Group("/", c.AuthMiddleware.Authenticate)

	// Event edit routes (require authentication + authorization)
	editGuarded := guarded.Group("/", c.AuthMiddleware.RequiresAdmin)
	editGuarded.POST("/events", c.EventHandler.CreateEvent)
	editGuarded.PUT("/events/:id", c.EventHandler.UpdateEvent)
	editGuarded.DELETE("/events/:id", c.EventHandler.DeleteEvent)
	editGuarded.POST("/events/photos/:id", c.EventHandler.AddPhotos)
	editGuarded.DELETE("/events/photos/:id", c.EventHandler.DeletePhotos)

	// Registration routes (authentication only)
	guarded.POST("/events/:id/register", c.RegistrationHandler.RegisterForEvent)
	guarded.DELETE("/events/:id/register", c.RegistrationHandler.CancelRegistrationEvent)
}
