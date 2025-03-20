package service

import (
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAllRoles()
}

func (s *RoleService) GetRoleById(roleId int64) (*models.Role, error) {
	return s.repo.GetRoleById(roleId)
}

func (s *RoleService) SetDefaultRole(roleId int64) error {
	return s.repo.SetDefaultRole(roleId)
}

func (s *RoleService) GetDefaultRole() (*models.Role, error) {
	return s.repo.GetDefaultRole()
}

func (s *RoleService) Save(role *models.Role) (*models.Role, error) {
	return s.repo.Save(role)
}

func (s *RoleService) AssignRoleToUser(userId int64, roleId int64) error {
	return s.repo.AssignRoleToUser(userId, roleId)
}

func (s *RoleService) GetRolesByUserId(userId int64) ([]models.Role, error) {
	return s.repo.GetRolesByUserId(userId)
}

func (s *RoleService) GetUsersByRoleId(roleId int64) ([]models.User, error) {
	return s.repo.GetUsersByRoleId(roleId)
}

func (s *RoleService) RemoveRoleFromUser(userId, roleId int64) error {
	return s.repo.RemoveRoleFromUser(userId, roleId)
}

func (s *RoleService) DeleteRole(roleId int64) error {
	return s.repo.DeleteRole(roleId)
}

func (s *RoleService) PatchAssignRoleToUsers(users []any, roleId int64) error {
	return s.repo.PatchAssignRoleToUsers(users, roleId)
}

func (s *RoleService) DeleteUserRoles(userId int64) error {
	return s.repo.DeleteUserRoles(userId)
}
