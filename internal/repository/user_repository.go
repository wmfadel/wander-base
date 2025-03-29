package repository

import (
	"fmt"
	"mime/multipart"

	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/requests"
	"github.com/wmfadel/wander-base/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db      *gorm.DB
	storage *utils.Storage
}

func NewUserRepository(db *gorm.DB, storage *utils.Storage) *UserRepository {
	return &UserRepository{db: db, storage: storage}
}

func (repo *UserRepository) Create(user *models.User) (*models.User, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash user password %w", err)
	}
	plainPassword := user.Password
	user.Password = hashedPassword

	result := repo.db.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed create new user %w", result.Error)
	}

	user, err = repo.GetUserByPhoneAndPassword(user.Phone, plainPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}

func (repo *UserRepository) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{ID: id}
	result := repo.db.
		Preload("Roles").
		Find(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	if user.Blocked() {
		return nil, fmt.Errorf("user blocked, assign \"user\" role to unblock")
	}

	return user, nil
}

func (repo *UserRepository) GetUserByPhoneAndPassword(phone, password string) (*models.User, error) {
	var user models.User

	// Find user by phone
	err := repo.db.Where("phone = ?", phone).Preload("Roles").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with phone %s not found", phone)
		}
		return nil, fmt.Errorf("failed to query user by phone: %w", err)
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password for user with phone %s", phone)
	}

	return &user, nil
}

func (repo *UserRepository) ValidateCredintials(loginRequest *requests.LoginRequest) error {

	var user models.User
	repo.db.Where("phone = ?", loginRequest.Phone).Find(&user)

	if user.ID == 0 {
		return fmt.Errorf("failed to find user: %w", gorm.ErrRecordNotFound)
	}
	loginRequest.ID = user.ID

	isValidPassword := utils.CheckPasswordHash(loginRequest.Password, user.Password)

	if !isValidPassword {
		return fmt.Errorf("invalid password, compared request password hash to existing password")
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
	user.Photo = url

	if err := repo.db.Save(user).Error; err != nil {
		return "", fmt.Errorf("failed to update user photo: %w", err)
	}
	return url, nil
}

func (r *UserRepository) UpdatePartially(userID int64, patch requests.PatchUser) error {
	if patch.IsEmpty() {
		return fmt.Errorf("no fields provided for update")
	}

	// Build a map of fields to update
	updates := make(map[string]interface{})
	if patch.FirstName != nil {
		updates["first_name"] = *patch.FirstName
	}
	if patch.LastName != nil {
		updates["last_name"] = *patch.LastName
	}

	// Perform the update
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user %d: %w", userID, result.Error)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", userID)
	}

	return nil
}
