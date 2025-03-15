package service

import (
	"fmt"
	"mime/multipart"

	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/repository"
)

type EventService struct {
	repo         *repository.EventRepository
	photoService *EventPhotoService
}

func NewEventService(repo *repository.EventRepository, photoService *EventPhotoService) *EventService {
	return &EventService{repo: repo, photoService: photoService}
}

func (s *EventService) CreateEvent(event *models.Event, photos []*multipart.FileHeader) error {
	if err := s.repo.Save(event); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	if len(photos) > 0 {
		if err := s.photoService.AddPhotos(event.ID, photos); err != nil {
			return fmt.Errorf("failed to add photos to event: %w", err)
		}
		// Fetch and set photo URLs
		photoURLs, err := s.photoService.GetPhotos(event.ID)
		if err != nil {
			return fmt.Errorf("failed to get photo URLs: %w", err)
		}
		event.Photos = photoURLs
	}
	return nil
}

func (s *EventService) GetEventById(eventId int64) (*models.Event, error) {
	return s.repo.GetEventById(eventId)
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.repo.GetAllEvents()
}
func (s *EventService) UpdatePartially(eventId int64, patch models.PatchEvent) error {
	return s.repo.UpdatePartially(eventId, patch)
}

func (s *EventService) Delete(eventId int64, photos []string) error {
	go s.photoService.DeletEventPhotos(eventId, photos)
	return s.repo.Delete(eventId)
}

func (s *EventService) Register(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}

func (s *EventService) CancelRegister(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}
