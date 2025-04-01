package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (repo *ActivityRepository) Save(activity *models.Activity) error {
	result := repo.db.Create(activity)
	if result.Error != nil {
		return fmt.Errorf("failed to save activity: %w", result.Error)
	}
	return nil
}

func (repo *ActivityRepository) GetAllActivities() ([]models.Activity, error) {
	var activities []models.Activity
	result := repo.db.Find(&activities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all activities: %w", result.Error)
	}
	return activities, nil
}

func (repo *ActivityRepository) GetActivityById(id int64) (*models.Activity, error) {
	var activity models.Activity
	result := repo.db.First(&activity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get activity %d: %w", id, result.Error)
	}
	return &activity, nil
}

func (repo *ActivityRepository) GetActivityBySlug(slug string) (*models.Activity, error) {
	var activity models.Activity
	result := repo.db.Where("slug = ?", slug).First(&activity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get activity %s: %w", slug, result.Error)
	}
	return &activity, nil
}

func (repo *ActivityRepository) Delete(activityId int64) error {
	result := repo.db.Delete(&models.Activity{}, activityId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete activity %d: %w", activityId, result.Error)
	}
	return nil
}
