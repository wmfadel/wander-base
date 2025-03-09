package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/routes"
	"github.com/wmfadel/escape-be/utils"
)

func main() {
	err := utils.LoadEnv()
	if err != nil {
		panic("No .env file found")
	}
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8080")
}
