package repository

import (
	"database/sql"
	"fmt"
	"mime/multipart"

	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type UserRepository struct {
	db      *sql.DB
	storage *utils.Storage
}

func NewUserRepository(db *sql.DB, storage *utils.Storage) *UserRepository {
	return &UserRepository{db: db, storage: storage}
}

func (repo *UserRepository) Create(user *models.User) error {
	query := "INSERT INTO users (phone, password) VALUES ($1, $2) RETURNING id"

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
	err = stmt.QueryRow(user.Phone, hashedPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to execute query for saving user %w", err)
	}
	// Set the user's ID
	user.ID = userID
	return nil
}

func (repo *UserRepository) GetUserByID(id int64) (*models.User, error) {
	query := "SELECT id, phone, password FROM users WHERE id = $1"

	user := &models.User{}

	err := repo.db.QueryRow(query, id).Scan(&user.ID, &user.Phone, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}
func (repo *UserRepository) ValidateCredintials(user *models.User) error {
	query := "SELECT id, password FROM users WHERE Phone = $1"

	row := repo.db.QueryRow(query, user.Phone)
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

func (repo *UserRepository) AddPhoto(user *models.User, photo *multipart.FileHeader) (string, error) {
	// Upload new photo
	url, err := repo.storage.UploadFile(photo, "user_photos", user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to upload photo: %w", err)
	}

	// Delete old photo if it exists
	if user.Photo != "" {
		if err := repo.storage.DeleteFile(user.Photo); err != nil {
			fmt.Printf("warning: failed to delete old photo %s: %v\n", user.Photo, err)
		}
	}

	// Update database
	query := "UPDATE users SET photo = $1 WHERE id = $2"
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query for updating photo: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(url, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to execute query for updating photo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("user %d not found", user.ID)
	}

	user.Photo = url
	return url, nil
}
