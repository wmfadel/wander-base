package service

import (
	"fmt"
	"mime/multipart"

	"github.com/wmfadel/wander-base/internal/repository"
)

type EventPhotoService struct {
	repo *repository.EventPhotoRepository
}

func NewEventPhotoService(repo *repository.EventPhotoRepository) *EventPhotoService {
	return &EventPhotoService{repo: repo}
}

func (s *EventPhotoService) AddPhotos(eventID int64, photos []*multipart.FileHeader) error {
	if len(photos) == 0 {
		return fmt.Errorf("no photos provided")
	}
	return s.repo.AddPhotos(eventID, photos)
}

func (s *EventPhotoService) GetPhotos(eventID int64) ([]string, error) {
	return s.repo.GetPhotos(eventID)
}

func (s *EventPhotoService) DeletEventPhotos(eventId int64, urls []string) error {
	return s.repo.DeletePhotos(eventId, urls)
}
