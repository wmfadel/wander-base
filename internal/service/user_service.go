package service

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/requests"
	"github.com/wmfadel/wander-base/internal/repository"
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
	user, err := s.repo.Create(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	// TODO: handle and user create hook
	defaultRole, err := s.rolesRepo.GetDefaultRole()
	if err != nil {
		return fmt.Errorf("failed to get default role: %w", err)
	}
	err = s.rolesRepo.AssignRoleToUser(user.ID, defaultRole.ID)
	if err != nil {
		return fmt.Errorf("failed to assign default role to user: %w", err)
	}
	return nil
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

	if len(user.Roles) == 0 {
		return nil, errors.New("user blocked, assign \"user\" role to unblock")
	}
	return user, nil
}
func (s *UserService) ValidateCredintials(loginRequest *requests.LoginRequest) error {
	return s.repo.ValidateCredintials(loginRequest)
}

func (s *UserService) UpdatePhoto(userId int64, photo *multipart.FileHeader) (string, error) {
	user := &models.User{ID: userId}
	return s.repo.AddPhoto(user, photo)
}

func (s *UserService) UpdateUser(user *models.User, patch *requests.PatchUser) error {
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
