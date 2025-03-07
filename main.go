package main

import (
	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/db"
	"githuv.com/wmfadel/go_events/routes"
	"githuv.com/wmfadel/go_events/utils"
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
