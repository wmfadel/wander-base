package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterAdminRoutes(r *gin.Engine, c di.DIContainer) {
	r.Group("/admin", c.AuthMiddleware.Authenticate, c.AuthMiddleware.RequiresAdmin)
	handler := c.AdminHandler
	r.GET("/all", handler.GetAllAdmins)
	r.GET("/organizers", handler.GetAllOrganizers)
	r.POST("/create", handler.CreateUser)
	r.POST("/delete", handler.DeleteUser)
	r.POST("/block", handler.BlockUser)
	r.POST("/roles", handler.AddRoleToUser)
	r.DELETE("/roles", handler.RemoveRoleFromUser)

}
