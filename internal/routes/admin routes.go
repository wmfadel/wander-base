package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterAdminRoutes(r *gin.Engine, c di.DIContainer) {
	guared := r.Group("/admin", c.AuthMiddleware.Authenticate, c.AuthMiddleware.RequiresAdmin)
	handler := c.AdminHandler
	guared.GET("/all", handler.GetAllRoles)             // lists all roles
	guared.GET("/admins", handler.GetAllAdmins)         // lists admins
	guared.GET("/organizers", handler.GetAllOrganizers) // lists organizers
	guared.POST("/create", handler.AddRole)             // creates a new role
	guared.POST("/delete", handler.DeleteRole)          // deletes a role
	guared.POST("/block", handler.BlockUser)            // blocks a user
	guared.POST("/roles", handler.AssignRoleToUser)     // assigns a role to a user
	guared.DELETE("/roles", handler.RemoveRoleFromUser) // removes a role from a user

}
