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

	server := gin.Default()
	server.Static("/user_photos", "./public/user_photos")
	server.Static("/event_photos", "./public/event_photos")

	container := di.NewDependencies(dbConnection)
	routes.RegisterRoutes(server, *container)

	server.Run(":8080")
}
