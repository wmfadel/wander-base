package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/internal/di"
	"github.com/wmfadel/escape-be/internal/routes"
	"github.com/wmfadel/escape-be/pkg/utils"
)

func main() {
	err := utils.LoadEnv()
	if err != nil {
		panic("No .env file found")
	}
	dbConnection := db.InitDB()

	container := di.NewDependencies(dbConnection)

	server := gin.Default()
	routes.RegisterRoutes(server, *container)

	server.Run(":8080")
}
