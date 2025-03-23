package service

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/repository"
)

type DestinationService struct {
	repo *repository.DestinationRepository
}

func NewDestinationService(repo *repository.DestinationRepository) *DestinationService {
	return &DestinationService{repo: repo}
}

func (s *DestinationService) Create(destination *models.Destination) error {
	if destination.Name == "" || destination.Location == nil {
		return fmt.Errorf("name and location are required")
	}
	return s.repo.Save(destination)
}

func (s *DestinationService) GetByID(id int64) (*models.Destination, error) {
	return s.repo.GetByID(id)
}

func (s *DestinationService) GetAll() ([]models.Destination, error) {
	return s.repo.GetAll()
}
