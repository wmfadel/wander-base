// Package main implements the entry point for the Escape backend service.
//
// The service provides a RESTful API for the Escape application, handling:
// - User authentication and authorization
// - Event management and booking
// - Static file serving for user and event photos
// - Database operations with migration and seeding capabilities
//
// Configuration is handled through environment variables loaded from a .env file.
// The server uses the Gin framework for routing and middleware management.
package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/internal/di"
	"github.com/wmfadel/escape-be/internal/routes"
	"github.com/wmfadel/escape-be/pkg/utils"
)

// main initializes and starts the Escape backend service.
//
// Server Configuration:
// - Default port: :8080
// - Static file paths:
//   - User photos: /user_photos -> ./public/user_photos
//   - Event photos: /event_photos -> ./public/event_photos
//
// Command-line Flags:
// - migrate: Run database migrations (usage: -migrate)
// - seed: Seed the database with initial data (usage: -seed)
//
// Example usage:
//
//	# Start server normally
//	$ go run cmd/main.go
//
//	# Start server with migrations
//	$ go run cmd/main.go -migrate
//
//	# Start server with migrations and seeding
//	$ go run cmd/main.go -migrate -seed
func main() {
	// Command-line flags for database management
	// migrate: When true, runs all pending database migrations
	// seed: When true, populates the database with initial test data
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Seed database with initial data")
	flag.Parse()

	// LoadEnv reads configuration from .env file
	// Required environment variables:
	// - DB_HOST: Database host address
	// - DB_PORT: Database port
	// - DB_USER: Database username
	// - DB_PASSWORD: Database password
	// - DB_NAME: Database name
	err := utils.LoadEnv()
	if err != nil {
		panic("No .env file found")
	}

	// Initialize database connection and optionally run migrations/seeding
	// InitDB handles:
	// - Establishing database connection
	// - Running migrations if migrate flag is true
	// - Seeding initial data if seed flag is true
	dbConnection := db.InitDB(*migrate, *seed)

	// Create a new Gin server instance with default middleware
	// Default middleware includes:
	// - Logger
	// - Recovery (crashes recovery)
	server := gin.Default()

	// Configure static file serving for uploaded media
	// These directories must exist and be writable by the server process
	server.Static("/user_photos", "./public/user_photos")
	server.Static("/event_photos", "./public/event_photos")

	// Initialize dependency injection container
	// This sets up all service dependencies including:
	// - Repositories
	// - Services
	// - Utilities
	container := di.NewDependencies(dbConnection)

	// Register all API routes with their respective handlers
	// Routes are grouped by feature and include:
	// - Authentication routes
	// - User management routes
	// - Event management routes
	// - Booking routes
	routes.RegisterRoutes(server, *container)

	// Start the HTTP server on port 8080
	// The server runs until it receives a shutdown signal
	server.Run(":8080")
}