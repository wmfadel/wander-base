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
	EventService       *service.EventService
	UserService        *service.UserService
	RolesService       *service.RoleService
	EventPhotosService *service.EventPhotoService

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
	eventPhotosRepository := repository.NewEventPhotoRepository(db, storage)
	eventRepo := repository.NewEventRepository(db, eventPhotosRepository)
	userRepo := repository.NewUserRepository(db, storage)
	rolesRepo := repository.NewRoleRepository(db)
	// Services initialization
	eventPhotosService := service.NewEventPhotoService(eventPhotosRepository)
	eventService := service.NewEventService(eventRepo, eventPhotosService)
	userService := service.NewUserService(userRepo)
	rolesService := service.NewRoleService(rolesRepo)
	// Handlers initialization
	eventHandler := handlers.NewEventHandler(eventService, eventPhotosService)
	userHandler := handlers.NewUserHandler(userService)
	registrationHandler := handlers.NewRegistrationHandler(eventService)

	// Middlewares initialization
	authMiddleware := middlewares.NewAuthMiddleware(eventService)

	return &DIContainer{
		DB:                  db,
		Storage:             storage,
		EventPhotosService:  eventPhotosService,
		EventService:        eventService,
		UserService:         userService,
		RolesService:        rolesService,
		EventHandler:        eventHandler,
		UserHandler:         userHandler,
		RegistrationHandler: registrationHandler,
		AuthMiddleware:      authMiddleware,
	}
}
