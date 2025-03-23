package di

import (
	"database/sql"

	"github.com/wmfadel/wander-base/internal/handlers"
	"github.com/wmfadel/wander-base/internal/repository"
	"github.com/wmfadel/wander-base/internal/service"
	middleware "github.com/wmfadel/wander-base/pkg/middlewares"
	"github.com/wmfadel/wander-base/pkg/utils"
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
	AuthHandler         *handlers.AuthHandler
	AdminHandler        *handlers.AdmingHandler
	ProfileHandler      *handlers.ProfileHandler
	EventHandler        *handlers.EventHandler
	RegistrationHandler *handlers.RegistrationHandler

	// Middlewares
	AuthMiddleware *middleware.AuthMiddleware
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
	userService := service.NewUserService(userRepo, rolesRepo)
	rolesService := service.NewRoleService(rolesRepo)
	// Handlers initialization

	authHandler := handlers.NewAuthHandler(userService)
	adminHandler := handlers.NewAdmingHandler(rolesService, userService)
	profileHandler := handlers.NewProfileHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService, eventPhotosService)
	registrationHandler := handlers.NewRegistrationHandler(eventService)

	// Middlewares initialization
	authMiddleware := middleware.NewAuthMiddleware(userService, eventService)

	return &DIContainer{
		DB:                  db,
		Storage:             storage,
		EventPhotosService:  eventPhotosService,
		EventService:        eventService,
		UserService:         userService,
		RolesService:        rolesService,
		AuthHandler:         authHandler,
		AdminHandler:        adminHandler,
		ProfileHandler:      profileHandler,
		EventHandler:        eventHandler,
		RegistrationHandler: registrationHandler,
		AuthMiddleware:      authMiddleware,
	}
}
