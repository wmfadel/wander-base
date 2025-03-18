// Package main is the entry point for the Escape backend application.
// It initializes the server, sets up database connections, and configures routes.
// The application supports database migrations and seeding through command-line flags.
package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/internal/di"
	"github.com/wmfadel/escape-be/internal/routes"
	"github.com/wmfadel/escape-be/pkg/utils"
)

// main is the entry point of the application. It performs the following tasks:
// - Parses command-line flags for database operations
// - Loads environment variables
// - Initializes database connection
// - Sets up static file serving
// - Configures dependency injection
// - Registers API routes
// - Starts the HTTP server
func main() {
	// Define command-line flags for database operations
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Seed database with initial data")
	flag.Parse()

	// Load environment variables from .env file
	err := utils.LoadEnv()
	if err != nil {
		panic("No .env file found")
	}

	// Initialize database connection and perform migrations/seeding if requested
	dbConnection := db.InitDB(*migrate, *seed)

	// Initialize Gin server with default middleware
	server := gin.Default()
	
	// Configure static file serving for user and event photos
	server.Static("/user_photos", "./public/user_photos")
	server.Static("/event_photos", "./public/event_photos")

	// Set up dependency injection container
	container := di.NewDependencies(dbConnection)
	
	// Register API routes with the server
	routes.RegisterRoutes(server, *container)

	// Start the server on port 8080
	server.Run(":8080")
}