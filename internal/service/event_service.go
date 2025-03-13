package service

import (
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/repository"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	return s.repo.Save(event)
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

func (s *EventService) Delete(eventId int64) error {
	return s.repo.Delete(eventId)
}

func (s *EventService) Register(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}

func (s *EventService) CancelRegister(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}
