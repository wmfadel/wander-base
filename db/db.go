package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/wmfadel/escape-be/pkg/utils"
)

func InitDB() *sql.DB {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		utils.GetFromEnv("DB_HOST"),
		utils.GetFromEnv("DB_PORT"),
		utils.GetFromEnv("DB_USER"),
		utils.GetFromEnv("DB_PASSWORD"),
		utils.GetFromEnv("DB_NAME"),
	)

	// Assign to the global DB variable directly
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Do NOT defer DB.Close() here; close it in main or when the app shuts down

	err = db.Ping()
	if err != nil {
		db.Close() // Clean up if ping fails
		log.Fatalf("Error pinging database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	createTables(db)
	log.Println("Database Connected...")
	return db
}

func createTables(db *sql.DB) {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	createEventsTable := `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime TIMESTAMP NOT NULL,
			user_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`
	_, err = db.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}

	createRegistrationTable := `
		CREATE TABLE IF NOT EXISTS registrations (
			id SERIAL PRIMARY KEY,
			event_id INTEGER,
			user_id INTEGER,
			FOREIGN KEY (event_id) REFERENCES events(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`
	_, err = db.Exec(createRegistrationTable)
	if err != nil {
		log.Fatalf("Failed to create registrations table: %v", err)
	}
}
