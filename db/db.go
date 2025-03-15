package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
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
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			phone TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime TIMESTAMP NOT NULL,
			user_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS registrations (
			id SERIAL PRIMARY KEY,
			event_id INTEGER,
			user_id INTEGER,
			FOREIGN KEY (event_id) REFERENCES events(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description Text NOT NULL,
			default_role BOOLEAN NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS user_roles (
  			user_id INTEGER NOT NULL,
    		role_id INTEGER NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES users(id),
   			FOREIGN KEY (role_id) REFERENCES roles(id),
    		PRIMARY KEY (user_id, role_id)
		)`,
	}

	seedData := []string{
		`INSERT INTO roles (name, description, default_role) VALUES ('admin', 'Admin role', FALSE)`,
		`INSERT INTO roles (name, description, default_role) VALUES ('organizer', 'Admin role', FALSE)`,
		`INSERT INTO roles (name, description, default_role) VALUES ('photographer', 'Admin role', FALSE)`,
		`INSERT INTO roles (name, description, default_role) VALUES ('user', 'User role', TRUE)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
	for _, data := range seedData {
		if _, err := db.Exec(data); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}

			log.Fatalf("Failed to add seed data to table: %v", err)
		}
	}
}
