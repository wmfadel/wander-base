package routes

import (
	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signupHanlder)
	server.POST("/login", loginHandler)

	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	guardedRoutes := server.Group("/", middlewares.Authenticate)

	guardedRoutes.POST("/events", creatEvent)
	guardedRoutes.PUT("/events/:id", middlewares.AuthorizeForEventEdits, updateEvent)
	guardedRoutes.DELETE("/events/:id", middlewares.AuthorizeForEventEdits, deleteEvent)

	guardedRoutes.POST("events/:id/register", registerForEvent)
	guardedRoutes.DELETE("events/:id/register", cancelRegistrationEvent)

}
