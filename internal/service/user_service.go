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
func (s *UserService) ValidateCredintials(loginRequest *models.LoginRequest) error {
	return s.repo.ValidateCredintials(loginRequest)
}

func (s *UserService) UpdatePhoto(userId int64, photo *multipart.FileHeader) (string, error) {
	user := &models.User{ID: userId}
	return s.repo.AddPhoto(user, photo)
}

func (s *UserService) UpdateUser(user *models.User, patch *models.PatchUser) error {
	err := s.repo.UpdatePartially(user.ID, *patch)
	if err != nil {
		return err
	}
	patch.Apply(user)
	roles, err := s.rolesRepo.GetRolesByUserId(user.ID)
	if err != nil {
		return err
	}
	user.Roles = roles
	return nil
}
