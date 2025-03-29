package service

import (
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/repository"
)

type RegistrationService struct {
	repo *repository.RegistrationRepository
}

func NewRegistrationService(repo *repository.RegistrationRepository) *RegistrationService {
	return &RegistrationService{repo: repo}
}

func (s *RegistrationService) Register(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}

func (s *RegistrationService) ApproveRegistration(userId int64, eventId int64) error {
	return s.repo.ApproveRegistration(userId, eventId)
}

func (s *RegistrationService) CancelRegister(userId int64, eventId int64) error {
	return s.repo.Register(userId, eventId)
}

func (s *RegistrationService) ApproveCancelRegister(userId int64, eventId int64) error {
	return s.repo.ApproveCancelRegister(userId, eventId)
}

func (s *RegistrationService) GetUsersWithStatus(eventID int64, status models.RegistrationStatus) ([]models.User, error) {
	return s.repo.GetUsersWithStatus(eventID, status)
}
