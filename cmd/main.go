package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/internal/di"
	"github.com/wmfadel/escape-be/internal/routes"
	"github.com/wmfadel/escape-be/pkg/utils"
)

func main() {
	migrate := flag.Bool("migrate", true, "Run database migrations")
	seed := flag.Bool("seed", true, "Seed database with initial data")
	flag.Parse()

	err := utils.LoadEnv()
	if err != nil {
		panic("No .env file found")
	}
	dbConnection := db.InitDB(*migrate, *seed)

	server := gin.Default()
	server.Static("/user_photos", "./public/user_photos")
	server.Static("/event_photos", "./public/event_photos")

	container := di.NewDependencies(dbConnection)
	routes.RegisterRoutes(server, *container)

	server.Run(":8080")
}