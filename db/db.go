package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"github.com/wmfadel/escape-be/pkg/utils"
)

func InitDB(migrate, seed bool) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		utils.GetFromEnv("DB_HOST"),
		utils.GetFromEnv("DB_PORT"),
		utils.GetFromEnv("DB_USER"),
		utils.GetFromEnv("DB_PASSWORD"),
		utils.GetFromEnv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalf("Error pinging database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Run migrations
	if migrate {
		if err := runMigrations(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Seed data if necessary
	if seed {
		if err := seedData(db); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
	}
	log.Println("Database Connected...")
	return db
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to migration files
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func seedData(db *sql.DB) error {
	// Check if roles table is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P01" { // Table doesn't exist
			count = 0
		} else {
			return fmt.Errorf("failed to check roles table: %w", err)
		}
	}

	if count == 0 {
		seedQueries := []string{
			`INSERT INTO roles (name, description, default_role) VALUES ('admin', 'Admin role', FALSE) ON CONFLICT (name) DO NOTHING`,
			`INSERT INTO roles (name, description, default_role) VALUES ('organizer', 'Organizer role', FALSE) ON CONFLICT (name) DO NOTHING`,
			`INSERT INTO roles (name, description, default_role) VALUES ('photographer', 'Photographer role', FALSE) ON CONFLICT (name) DO NOTHING`,
			`INSERT INTO roles (name, description, default_role) VALUES ('user', 'User role', TRUE) ON CONFLICT (name) DO NOTHING`,
		}

		for _, query := range seedQueries {
			if _, err := db.Exec(query); err != nil {
				return fmt.Errorf("failed to seed data: %w", err)
			}
		}
		log.Println("Seed data inserted")
	} else {
		log.Println("Seed data already exists")
	}
	return nil
}
