package repository

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/pkg/utils"
	"gorm.io/gorm"
)

type EventPhotoRepository struct {
	db      *gorm.DB
	storage *utils.Storage
}

func NewEventPhotoRepository(db *gorm.DB, storage *utils.Storage) *EventPhotoRepository {
	return &EventPhotoRepository{db: db, storage: storage}
}

func (repo *EventPhotoRepository) AddPhotos(eventID int64, photos []*multipart.FileHeader) error {
	var eventPhotos []models.EventPhoto
	var uploadErrors []error

	// Process each photo, collecting successful uploads and logging failures
	for _, photo := range photos {
		url, err := repo.storage.UploadFile(photo, "events", eventID)
		if err != nil {
			// Log the error and continue
			log.Printf("Failed to upload photo for event %d: %v", eventID, err)
			uploadErrors = append(uploadErrors, fmt.Errorf("failed to upload photo: %w", err))
			continue
		}

		// Add successful upload to the list
		eventPhotos = append(eventPhotos, models.EventPhoto{
			EventID: eventID,
			URL:     url,
		})
	}

	// If no photos were successfully uploaded, return an error with details
	if len(eventPhotos) == 0 && len(uploadErrors) > 0 {
		return fmt.Errorf("no photos uploaded successfully: %v errors occurred", len(uploadErrors))
	}

	// Save successfully uploaded photos to the database
	if len(eventPhotos) > 0 {
		if err := repo.db.Create(&eventPhotos).Error; err != nil {
			return fmt.Errorf("failed to create event photos: %w", err)
		}
	}

	// If there were upload errors but some successes, log a warning but don’t fail
	if len(uploadErrors) > 0 {
		log.Printf("Partial success: %d photos uploaded, %d failed for event %d", len(eventPhotos), len(uploadErrors), eventID)
		// Optionally return nil or a custom error with details
		return nil // Proceed as success since some photos were added
	}

	return nil
}

func (repo *EventPhotoRepository) GetPhotos(eventID int64) ([]models.EventPhoto, error) {
	var photos []models.EventPhoto
	repo.db.Where("event_id = ?", eventID).Find(&photos)

	return photos, nil
}

func (repo *EventPhotoRepository) DeletePhotos(eventID int64, urls []string) error {
	if len(urls) == 0 {
		return nil // Nothing to delete
	}

	// Delete from database using GORM
	result := repo.db.Where("event_id = ? AND photo_url IN ?", eventID, urls).Delete(&models.EventPhoto{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete photos from event %d: %w", eventID, result.Error)
	}

	// Check if any rows were affected
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return fmt.Errorf("no photos found for event %d matching provided URLs", eventID)
	}

	// Delete files from storage, ignoring failures
	for _, url := range urls {
		if err := repo.storage.DeleteFile(url); err != nil {
			log.Printf("Warning: failed to delete file %s from storage: %v", url, err)
			// Continue despite failure—DB is already updated
		}
	}

	return nil
}
