package di

import (
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/wmfadel/wander-base/internal/handlers"
	"github.com/wmfadel/wander-base/internal/repository"
	"github.com/wmfadel/wander-base/internal/service"
	middleware "github.com/wmfadel/wander-base/pkg/middlewares"
	"github.com/wmfadel/wander-base/pkg/utils"
	"gorm.io/gorm"
)

// DIContainer holds all shared app dependencies
type DIContainer struct {
	// DB
	DB      *gorm.DB
	Storage *utils.Storage
	// Services
	EventService        *service.EventService
	UserService         *service.UserService
	RolesService        *service.RoleService
	EventPhotosService  *service.EventPhotoService
	RegistrationService *service.RegistrationService
	ActivityService     *service.ActivityService
	DestinationService  *service.DestinationService
	CommentService      *service.CommentService

	// Handlers
	AuthHandler         *handlers.AuthHandler
	AdminHandler        *handlers.AdmingHandler
	ProfileHandler      *handlers.ProfileHandler
	EventHandler        *handlers.EventHandler
	RegistrationHandler *handlers.RegistrationHandler
	ActivityHandler     *handlers.ActivityHandler
	DestinationHandler  *handlers.DestinationHandler
	CommentHandler      *handlers.CommentHandler

	// Middlewares
	AuthMiddleware *middleware.AuthMiddleware
}

// NewDependencies initializes the dependency container
func NewDependencies(db *gorm.DB) *DIContainer {
	key := utils.GetFromEnv("OPENAI_API_KEY")
	log.Print("key: ", key)
	openaiClient := openai.NewClient(key)
	storage := utils.NewStorage("http://localhost:8080")
	// Repositories initialization
	eventPhotosRepository := repository.NewEventPhotoRepository(db, storage)
	eventRepo := repository.NewEventRepository(db, eventPhotosRepository)
	userRepo := repository.NewUserRepository(db, storage)
	rolesRepo := repository.NewRoleRepository(db)
	registrationRepo := repository.NewRegistrationRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	destinationRepo := repository.NewDestinationRepository(db)
	commentRepository := repository.NewCommentRepository(db)
	// Services initialization
	eventPhotosService := service.NewEventPhotoService(eventPhotosRepository)
	eventService := service.NewEventService(eventRepo, eventPhotosService)
	userService := service.NewUserService(userRepo, rolesRepo)
	rolesService := service.NewRoleService(rolesRepo)
	registrationService := service.NewRegistrationService(registrationRepo)
	activityService := service.NewActivityService(activityRepo)
	destinationService := service.NewDestinationService(destinationRepo)
	auditService := service.NewAuditService(openaiClient, 0.7)
	commentService := service.NewCommentService(commentRepository, auditService)
	// Handlers initialization

	authHandler := handlers.NewAuthHandler(userService)
	adminHandler := handlers.NewAdmingHandler(rolesService, userService)
	profileHandler := handlers.NewProfileHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService, eventPhotosService)
	registrationHandler := handlers.NewRegistrationHandler(registrationService)
	activityHandler := handlers.NewActivityHandler(activityService)
	destinationHandler := handlers.NewDestinationHandler(destinationService)
	commentHandler := handlers.NewCommentHandler(commentService)
	// Middlewares initialization
	authMiddleware := middleware.NewAuthMiddleware(userService, eventService)

	return &DIContainer{
		// DB
		DB:      db,
		Storage: storage,
		// Services
		EventPhotosService:  eventPhotosService,
		EventService:        eventService,
		UserService:         userService,
		RolesService:        rolesService,
		RegistrationService: registrationService,
		ActivityService:     activityService,
		DestinationService:  destinationService,
		CommentService:      commentService,
		// Handlers
		AuthHandler:         authHandler,
		AdminHandler:        adminHandler,
		ProfileHandler:      profileHandler,
		EventHandler:        eventHandler,
		RegistrationHandler: registrationHandler,
		ActivityHandler:     activityHandler,
		DestinationHandler:  destinationHandler,
		CommentHandler:      commentHandler,
		// Middlewares
		AuthMiddleware: authMiddleware,
	}
}
