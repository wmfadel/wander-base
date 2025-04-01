package service

import (
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/repository"
)

type DestinationService struct {
	repo *repository.DestinationRepository
}

func NewDestinationService(repo *repository.DestinationRepository) *DestinationService {
	return &DestinationService{repo: repo}
}

func (service *DestinationService) Save(destination *models.Destination) error {
	return service.repo.Save(destination)
}

func (service *DestinationService) GetAllDestinations() ([]models.Destination, error) {
	return service.repo.GetAllDestinations()
}

func (service *DestinationService) GetDestinationById(id int64) (*models.Destination, error) {
	return service.repo.GetDestinationById(id)
}

func (service *DestinationService) DeleteDestination(destinationId int64) error {
	return service.repo.DeleteDestination(destinationId)

}
