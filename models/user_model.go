package models

import (
	"fmt"

	"github.com/wmfadel/escape-be/db"
	"github.com/wmfadel/escape-be/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=5"`
}

func (u *User) Save() error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query for saving user %w", err)
	}
	defer stmt.Close()

	// Hash the password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to hash user password %w", err)
	}

	// Execute the query and retrieve the ID
	var userID int64
	err = stmt.QueryRow(u.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to execute query for saving user %w", err)
	}

	// Set the user's ID
	u.ID = userID
	return nil
}

func (u *User) ValidateCredintials() error {
	query := "SELECT id, password FROM users WHERE email = $1"

	row := db.DB.QueryRow(query, u.Email)
	var storedPassword string

	err := row.Scan(&u.ID, &storedPassword)

	if err != nil {
		return fmt.Errorf("failed to find user, wrong credintials: %w", err)
	}

	isValidPassword := utils.CheckPasswordHash(u.Password, storedPassword)

	if !isValidPassword {
		return fmt.Errorf("invalid password, compared request password hash to existing password: %w", err)
	}
	return nil
}
