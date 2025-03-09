package models

import (
	"errors"

	"githuv.com/wmfadel/go_events/db"
	"githuv.com/wmfadel/go_events/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Hash the password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Execute the query and retrieve the ID
	var userID int64
	err = stmt.QueryRow(u.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return err
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
		return errors.New("wrong credintials")
	}

	isValidPassword := utils.CheckPasswordHash(u.Password, storedPassword)

	if !isValidPassword {
		return errors.New("wrong credintials")
	}
	return nil
}
