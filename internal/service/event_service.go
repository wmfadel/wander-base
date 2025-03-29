package service

import (
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/requests"
	"github.com/wmfadel/wander-base/internal/repository"
)

type EventService struct {
	repo         *repository.EventRepository
	photoService *EventPhotoService
}

func NewEventService(repo *repository.EventRepository, photoService *EventPhotoService) *EventService {
	return &EventService{repo: repo, photoService: photoService}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	return s.repo.Save(event)
}

func (s *EventService) SetDestinations(destinations []models.EventDestinationRequest, eventID int64) error {
	return s.repo.SetDestinations(destinations, eventID)
}

func (s *EventService) RemoveDestinations(destinationIDs []int64, eventID int64) error {
	return s.repo.RemoveDestinations(destinationIDs, eventID)
}

func (s *EventService) AddActivities(activityIDs []int64, eventID int64) error {
	return s.repo.AddActivities(activityIDs, eventID)
}

func (s *EventService) RemoveActivities(activityIDs []int64, eventID int64) error {
	return s.repo.RemoveActivities(activityIDs, eventID)
}

func (s *EventService) GetEventById(eventId int64) (*models.Event, error) {
	return s.repo.GetEventById(eventId)
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.repo.GetAllEvents()
}
func (s *EventService) UpdatePartially(eventId int64, patch requests.PatchEvent) error {
	return s.repo.UpdatePartially(eventId, patch)
}

func (s *EventService) Delete(eventId int64, photos []string) error {
	go s.photoService.DeletEventPhotos(eventId, photos)
	return s.repo.Delete(eventId)
}
