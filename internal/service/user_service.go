package service

import (
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Save(user *models.User) error {
	return s.repo.Save(user)
}

func (s *UserService) ValidateCredintials(user *models.User) error {
	return s.repo.ValidateCredintials(user)
}
