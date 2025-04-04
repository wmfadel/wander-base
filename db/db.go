package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(migrate, seed bool) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		utils.GetFromEnv("DB_HOST"),
		utils.GetFromEnv("DB_PORT"),
		utils.GetFromEnv("DB_USER"),
		utils.GetFromEnv("DB_PASSWORD"),
		utils.GetFromEnv("DB_NAME"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if migrate {
		// Create the registration_status ENUM type
		err = db.Exec(`CREATE TYPE registration_status AS ENUM (
			'registered',
			'pending_registration',
			'cancelled',
			'pending_cancellation'
		)`).Error
		if err != nil && !isAlreadyExistsError(err) {
			log.Fatalf("Failed to create registration_status ENUM: %v", err)
		}

		// Auto-migrate tables
		err = db.AutoMigrate(
			&models.User{},
			&models.Role{},
			&models.UserRole{},
			&models.Event{},
			&models.Destination{},
			&models.EventDestination{},
			&models.Activity{},
			&models.EventActivities{},
			&models.EventPhoto{},
			&models.Registration{},
			&models.Comment{},
		)
		if err != nil {
			log.Fatalf("Failed to auto-migrate: %v", err)
		}
		log.Println("Database migrated")
	}

	if seed {
		// Seed roles
		if err := seedRoles(db); err != nil {
			log.Fatalf("Failed to seed roles: %v", err)
		}
		log.Println("Roles seeded successfully")
	}

	log.Println("Database Connected...")
	return db
}

// Helper to check if error is due to type already existing
func isAlreadyExistsError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "42710" // PostgreSQL "duplicate_object" error code
	}
	return false
}

// seedRoles seeds the roles table with predefined roles
func seedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "admin", Description: "Admin role", Default: false},
		{Name: "organizer", Description: "Organizer role", Default: false},
		{Name: "photographer", Description: "Photographer role", Default: false},
		{Name: "user", Description: "User role", Default: true},
	}

	for _, role := range roles {
		// Check if role already exists by name (unique constraint)
		var existingRole models.Role
		err := db.Where("name = ?", role.Name).First(&existingRole).Error
		if err == gorm.ErrRecordNotFound {
			// Role doesn’t exist, create it
			if err := db.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Name, err)
			}
		} else if err != nil {
			// Unexpected error
			return fmt.Errorf("failed to check existing role %s: %w", role.Name, err)
		}
		// If role exists, do nothing (mimics ON CONFLICT DO NOTHING)
	}

	return nil
}
