package service

import (
	"mime/multipart"

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

func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *UserService) ValidateCredintials(user *models.User) error {
	return s.repo.ValidateCredintials(user)
}

func (s *UserService) UpdatePhoto(userId int64, photo *multipart.FileHeader) (string, error) {
	user := &models.User{ID: userId}
	return s.repo.AddPhoto(user, photo)
}
