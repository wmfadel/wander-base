package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) GetAllRoles() ([]models.Role, error) {
	roles := []models.Role{}
	results := repo.db.Find(&roles)

	if results.Error != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", results.Error)
	}
	return roles, nil
}

func (repo *RoleRepository) GetRoleById(roleId int64) (*models.Role, error) {
	role := models.Role{ID: roleId}
	result := repo.db.First(&role)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role with id %d not found", roleId)
		}
		return nil, fmt.Errorf("failed to get role with id %d: %w", roleId, result.Error)
	}
	return &role, nil
}
func (repo *RoleRepository) SetDefaultRole(roleID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Reset all roles to non-default
		err := tx.Model(&models.Role{}).
			Where("default_role = ?", true).
			Update("default_role", false).Error
		if err != nil {
			return fmt.Errorf("failed to reset default roles: %w", err)
		}

		// Set the specified role as default
		result := tx.Model(&models.Role{}).
			Where("id = ?", roleID).
			Update("default_role", true)
		if result.Error != nil {
			return fmt.Errorf("failed to set default role for role %d: %w", roleID, result.Error)
		}

		// Verify the role exists and was updated
		if result.RowsAffected == 0 {
			return fmt.Errorf("role %d not found", roleID)
		}

		return nil
	})
}

func (repo *RoleRepository) GetDefaultRole() (*models.Role, error) {
	var role models.Role
	results := repo.db.Where("default_role = ?", true).First(&role)
	if results.Error != nil {
		return nil, fmt.Errorf("failed to get default role: %w", results.Error)
	}

	return &role, nil
}

func (repo *RoleRepository) Save(role *models.Role) (*models.Role, error) {
	results := repo.db.Create(role)
	if results.Error != nil {
		return nil, fmt.Errorf("failed to save role: %w", results.Error)
	}
	results = repo.db.Where("name = ?", role.Name).First(&role)
	if results.Error != nil {
		return nil, fmt.Errorf("failed to get saved role: %w", results.Error)
	}
	return role, nil
}

func (repo *RoleRepository) AssignRoleToUser(userID, roleID int64) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := repo.db.Create(&userRole).Error; err != nil {
		return fmt.Errorf("failed to assign role %d to user %d: %w", roleID, userID, err)
	}

	return nil
}

func (repo *RoleRepository) GetRolesByUserId(userID int64) ([]models.Role, error) {
	var user models.User

	err := repo.db.
		Preload("Roles").
		First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.Role{}, nil // Return empty slice if user not found
		}
		return nil, fmt.Errorf("failed to get roles for user %d: %w", userID, err)
	}

	return user.Roles, nil
}

func (repo *RoleRepository) GetUsersByRoleId(roleID int64) ([]models.User, error) {
	var users []models.User

	err := repo.db.
		Joins("JOIN user_roles ur ON users.id = ur.user_id").
		Where("ur.role_id = ?", roleID).
		Preload("Roles"). // Load roles for each user
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users for role %d: %w", roleID, err)
	}

	return users, nil
}

func (repo *RoleRepository) RemoveRoleFromUser(userID, roleID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Check if the user has other roles
		var roleCount int64
		err := tx.Model(&models.UserRole{}).
			Where("user_id = ? AND role_id != ?", userID, roleID).
			Count(&roleCount).Error
		if err != nil {
			return fmt.Errorf("failed to count user roles: %w", err)
		}
		if roleCount == 0 {
			return fmt.Errorf("cannot remove role %d from user %d: user would be left with no roles", roleID, userID)
		}

		// Step 2: Remove the role from the user
		result := tx.Where("user_id = ? AND role_id = ?", userID, roleID).
			Delete(&models.UserRole{})
		if result.Error != nil {
			return fmt.Errorf("failed to remove role %d from user %d: %w", roleID, userID, result.Error)
		}

		// Verify the role was removed
		if result.RowsAffected == 0 {
			return fmt.Errorf("role %d not assigned to user %d", roleID, userID)
		}

		return nil
	})
}
func (repo *RoleRepository) DeleteRole(roleID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Check if role is default
		var defaultRole models.Role
		err := tx.Where("default_role = ?", true).First(&defaultRole).Error
		if err != nil {
			return fmt.Errorf("failed to get default role: %w", err)
		}
		if defaultRole.ID == roleID {
			return fmt.Errorf("cannot delete default role")
		}

		// Find users with only this role
		var userRoles []models.UserRole
		err = tx.Raw(`
			SELECT ur.*
			FROM user_roles ur
			WHERE ur.role_id = ?
			AND NOT EXISTS (
				SELECT 1
				FROM user_roles ur2
				WHERE ur2.user_id = ur.user_id
				AND ur2.role_id != ?
			)`,
			roleID, roleID).
			Scan(&userRoles).Error
		if err != nil {
			return fmt.Errorf("failed to get users with only this role: %w", err)
		}

		// Migrate users to default role if needed
		if len(userRoles) > 0 {
			userIDs := make([]interface{}, len(userRoles))
			for i, ur := range userRoles {
				userIDs[i] = ur.UserID
			}
			for _, userID := range userIDs {
				err = tx.Create(&models.UserRole{UserID: userID.(int64), RoleID: defaultRole.ID}).Error
				if err != nil {
					return fmt.Errorf("failed to migrate user %v to default role: %w", userID, err)
				}
			}
		}

		// Delete the role
		result := tx.Where("id = ?", roleID).Delete(&models.Role{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete role %d: %w", roleID, result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("role %d not found", roleID)
		}

		return nil
	})
}

func (repo *RoleRepository) PatchAssignRoleToUsers(userIDs []any, roleID int64) error {
	if len(userIDs) == 0 {
		return nil // Nothing to do
	}

	// Convert userIDs to int64 slice for GORM compatibility
	ids := make([]int64, len(userIDs))
	for i, id := range userIDs {
		if idInt64, ok := id.(int64); ok {
			ids[i] = idInt64
		} else {
			return fmt.Errorf("invalid user ID type at index %d", i)
		}
	}

	// Batch create user_roles entries
	var userRoles []models.UserRole
	for _, userID := range ids {
		userRoles = append(userRoles, models.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
	if err != nil {
		return fmt.Errorf("failed to assign role %d to users: %w", roleID, err)
	}

	return nil
}

func (repo *RoleRepository) DeleteUserRoles(userID int64) error {
	err := repo.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user roles for user %d: %w", userID, err)
	}
	return nil
}
