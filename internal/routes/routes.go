package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/di"
)

func RegisterRoutes(server *gin.Engine, c di.DIContainer) {
	RegisterAuthRoutes(server, c)
	RegisterAdminRoutes(server, c)
	RegisterProfileRoutes(server, c)
	RegisterEventRoutes(server, c)
}
