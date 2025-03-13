package repository

import (
	"database/sql"
	"fmt"

	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(user *models.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	// Prepare the statement
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query for saving user %w", err)
	}
	defer stmt.Close()

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash user password %w", err)
	}

	// Execute the query and retrieve the ID
	var userID int64
	err = stmt.QueryRow(user.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to execute query for saving user %w", err)
	}

	// Set the user's ID
	user.ID = userID
	return nil
}

func (repo *UserRepository) ValidateCredintials(user *models.User) error {
	query := "SELECT id, password FROM users WHERE email = $1"

	row := repo.db.QueryRow(query, user.Email)
	var storedPassword string

	err := row.Scan(&user.ID, &storedPassword)

	if err != nil {
		return fmt.Errorf("failed to find user, wrong credintials: %w", err)
	}

	isValidPassword := utils.CheckPasswordHash(user.Password, storedPassword)

	if !isValidPassword {
		return fmt.Errorf("invalid password, compared request password hash to existing password: %w", err)
	}
	return nil
}
