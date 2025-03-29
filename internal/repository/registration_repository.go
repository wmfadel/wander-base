package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (repo *RegistrationRepository) Register(userId, eventId int64) error {
	result := repo.db.Create(&models.Registration{UserID: userId, EventID: eventId})
	if result.Error != nil {
		return fmt.Errorf("failed to register user %d for event %d: %w", userId, eventId, result.Error)
	}
	return nil
}

func (repo *RegistrationRepository) ApproveRegistration(userId, eventId int64) error {

	eventRegistration := models.Registration{UserID: userId, EventID: eventId}
	result := repo.db.Find(&eventRegistration)

	if result.Error != nil {
		return fmt.Errorf("failed to find registration for user %d and event %d: %w", userId, eventId, result.Error)
	}

	if eventRegistration.Status != models.PendingRegistration {
		return fmt.Errorf("registration for user %d and event %d is already approved", userId, eventId)
	}

	eventRegistration.Status = models.Registered

	result = repo.db.Save(&eventRegistration)
	if result.Error != nil {
		return fmt.Errorf("failed to register user %d for event %d: %w", userId, eventId, result.Error)
	}
	return nil
}

func (repo *RegistrationRepository) CancelRegister(userId int64, eventId int64) error {
	eventRegistration := models.Registration{UserID: userId, EventID: eventId}
	result := repo.db.Find(&eventRegistration)

	if result.Error != nil {
		return fmt.Errorf("failed to find registration for user %d and event %d: %w", userId, eventId, result.Error)
	}

	if eventRegistration.Status == models.Cancelled {
		return fmt.Errorf("registration for user %d and event %d is already cancelled", userId, eventId)
	}

	eventRegistration.Status = models.PendingCancellation

	result = repo.db.Save(&eventRegistration)
	if result.Error != nil {
		return fmt.Errorf("failed to register user %d for event %d: %w", userId, eventId, result.Error)
	}
	return nil
}

func (repo *RegistrationRepository) ApproveCancelRegister(userId int64, eventId int64) error {
	eventRegistration := models.Registration{UserID: userId, EventID: eventId}
	result := repo.db.Find(&eventRegistration)

	if result.Error != nil {
		return fmt.Errorf("failed to find registration for user %d and event %d: %w", userId, eventId, result.Error)
	}

	if eventRegistration.Status != models.PendingCancellation {
		return fmt.Errorf("user %d didn't request cancellation for event %d", userId, eventId)
	}

	eventRegistration.Status = models.Cancelled

	result = repo.db.Save(&eventRegistration)
	if result.Error != nil {
		return fmt.Errorf("failed to register user %d for event %d: %w", userId, eventId, result.Error)
	}
	return nil
}

func (repo *RegistrationRepository) GetUsersWithStatus(eventID int64, status models.RegistrationStatus) ([]models.User, error) {
	var registrations []models.Registration

	// Query registrations with the specified eventID and status, preloading User
	err := repo.db.
		Where("event_id = ? AND status = ?", eventID, status).
		Preload("User").
		Find(&registrations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get registrations for event %d with status %s: %w", eventID, status, err)
	}

	// Extract users from registrations
	users := make([]models.User, 0, len(registrations))
	for _, reg := range registrations {
		users = append(users, reg.User)
	}

	return users, nil
}
