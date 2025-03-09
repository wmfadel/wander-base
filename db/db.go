package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "pg-2ea894ef-escape-eg.e.aivencloud.com"
	port     = 14470
	user     = "avnadmin"
	password = "AVNS_prM8rH6cDM-1MBjZdxP"
	dbname   = "defaultdb"
)

func InitDB() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Assign to the global DB variable directly
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Do NOT defer DB.Close() here; close it in main or when the app shuts down

	err = DB.Ping()
	if err != nil {
		DB.Close() // Clean up if ping fails
		log.Fatalf("Error pinging database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
	log.Println("Database Connected...")
}

func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`
	_, err := DB.Exec(createUsersTable)
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
	_, err = DB.Exec(createEventsTable)
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
	_, err = DB.Exec(createRegistrationTable)
	if err != nil {
		log.Fatalf("Failed to create registrations table: %v", err)
	}
}
