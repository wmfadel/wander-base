package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterRoutes(server *gin.Engine, c di.DIContainer) {
	server.POST("/signup", c.UserHandler.SignupHanlder)
	server.POST("/login", c.UserHandler.LoginHandler)

	server.GET("/events", c.EventHandler.GetEvents)
	server.GET("/events/:id", c.EventHandler.GetEvent)

	guardedRoutes := server.Group("/", c.AuthMiddleware.Authenticate)

	guardedRoutes.POST("/users/photo", c.UserHandler.UpdatePhoto)

	// Event Routes
	guardedRoutes.POST("/events", c.EventHandler.CreateEvent)
	guardedRoutes.PUT("/events/:id", c.AuthMiddleware.AuthorizeForEventEdits, c.EventHandler.UpdateEvent)
	guardedRoutes.DELETE("/events/:id", c.AuthMiddleware.AuthorizeForEventEdits, c.EventHandler.DeleteEvent)

	guardedRoutes.POST("/events/photos/:id", c.AuthMiddleware.AuthorizeForEventEdits, c.EventHandler.AddPhotos)
	guardedRoutes.DELETE("/events/photos/:id", c.AuthMiddleware.AuthorizeForEventEdits, c.EventHandler.DeletePhotos)

	guardedRoutes.POST("events/:id/register", c.RegistrationHandler.RegisterForEvent)
	guardedRoutes.DELETE("events/:id/register", c.RegistrationHandler.CancelRegistrationEvent)

}
