package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"gorm.io/gorm"
)

type DestinationRepository struct {
	db *gorm.DB
}

func NewDestinationRepository(db *gorm.DB) *DestinationRepository {
	return &DestinationRepository{db: db}
}

func (repo *DestinationRepository) Save(destination *models.Destination) error {
	result := repo.db.Create(destination)
	if result.Error != nil {
		return fmt.Errorf("failed to save destination: %w", result.Error)
	}
	return nil
}

func (repo *DestinationRepository) GetAllDestinations() ([]models.Destination, error) {
	var destinations []models.Destination
	result := repo.db.Find(&destinations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all destinations: %w", result.Error)
	}
	return destinations, nil
}

func (repo *DestinationRepository) GetDestinationById(id int64) (*models.Destination, error) {
	var destination models.Destination
	result := repo.db.First(&destination, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get destination %d: %w", id, result.Error)
	}
	return &destination, nil
}

func (repo *DestinationRepository) DeleteDestination(destinationId int64) error {

	result := repo.db.Delete(&models.Destination{}, destinationId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete destination %d: %w", destinationId, result.Error)
	}
	return nil

}
