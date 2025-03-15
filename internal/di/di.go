package di

import (
	"database/sql"

	"github.com/wmfadel/escape-be/internal/handlers"
	"github.com/wmfadel/escape-be/internal/repository"
	"github.com/wmfadel/escape-be/internal/service"
	"github.com/wmfadel/escape-be/pkg/middlewares"
	"github.com/wmfadel/escape-be/pkg/utils"
)

// DIContainer holds all shared app dependencies
type DIContainer struct {
	// DB
	DB      *sql.DB
	Storage *utils.Storage
	// Services
	EventService *service.EventService
	UserService  *service.UserService
	RolesService *service.RoleService

	// Handlers
	EventHandler        *handlers.EventHandler
	UserHandler         *handlers.UserHandler
	RegistrationHandler *handlers.RegistrationHandler

	// Middlewares
	AuthMiddleware *middlewares.AuthMiddleware
}

// NewDependencies initializes the dependency container
func NewDependencies(db *sql.DB) *DIContainer {
	storage := utils.NewStorage("http://localhost:8080")
	// Repositories initialization
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db, storage)
	rolesRepo := repository.NewRoleRepository(db)
	// Services initialization
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)
	rolesService := service.NewRoleService(rolesRepo)
	// Handlers initialization
	eventHandler := handlers.NewEventHandler(eventService)
	userHandler := handlers.NewUserHandler(userService)
	registrationHandler := handlers.NewRegistrationHandler(eventService)

	// Middlewares initialization
	authMiddleware := middlewares.NewAuthMiddleware(eventService)

	return &DIContainer{
		DB:                  db,
		Storage:             storage,
		EventService:        eventService,
		UserService:         userService,
		RolesService:        rolesService,
		EventHandler:        eventHandler,
		UserHandler:         userHandler,
		RegistrationHandler: registrationHandler,
		AuthMiddleware:      authMiddleware,
	}
}
