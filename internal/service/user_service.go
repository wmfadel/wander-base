package service

import (
	"mime/multipart"

	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/repository"
)

type UserService struct {
	repo      *repository.UserRepository
	rolesRepo *repository.RoleRepository
}

func NewUserService(repo *repository.UserRepository, rolesRepo *repository.RoleRepository) *UserService {
	return &UserService{
		repo:      repo,
		rolesRepo: rolesRepo,
	}
}

func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(id)

	if err != nil {
		return nil, err
	}
	user.Roles, err = s.rolesRepo.GetRolesByUserId(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) ValidateCredintials(user *models.User) error {
	return s.repo.ValidateCredintials(user)
}

func (s *UserService) UpdatePhoto(userId int64, photo *multipart.FileHeader) (string, error) {
	user := &models.User{ID: userId}
	return s.repo.AddPhoto(user, photo)
}
