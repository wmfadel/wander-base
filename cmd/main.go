package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/db"
	"github.com/wmfadel/wander-base/internal/di"
	"github.com/wmfadel/wander-base/internal/routes"
	"github.com/wmfadel/wander-base/pkg/utils"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Run database seeder")
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
