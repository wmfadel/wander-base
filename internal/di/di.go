package di

import (
	"database/sql"

	"github.com/wmfadel/escape-be/internal/handlers"
	"github.com/wmfadel/escape-be/internal/repository"
	"github.com/wmfadel/escape-be/internal/service"
	"github.com/wmfadel/escape-be/pkg/middlewares"
)

// DIContainer holds all shared app dependencies
type DIContainer struct {
	DB                  *sql.DB
	EventService        *service.EventService
	UserService         *service.UserService
	EventHandler        *handlers.EventHandler
	UserHandler         *handlers.UserHandler
	RegistrationHandler *handlers.RegistrationHandler
	AuthMiddleware      *middlewares.AuthMiddleware
}

// NewDependencies initializes the dependency container
func NewDependencies(db *sql.DB) *DIContainer {
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db)
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)
	eventHandler := handlers.NewEventHandler(eventService)
	userHandler := handlers.NewUserHandler(userService)
	registrationHandler := handlers.NewRegistrationHandler(eventService)
	authMiddleware := middlewares.NewAuthMiddleware(eventService)

	return &DIContainer{
		DB:                  db,
		EventService:        eventService,
		UserService:         userService,
		EventHandler:        eventHandler,
		UserHandler:         userHandler,
		RegistrationHandler: registrationHandler,
		AuthMiddleware:      authMiddleware,
	}
}
