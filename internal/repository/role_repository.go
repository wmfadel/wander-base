package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/wmfadel/escape-be/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) GetAllRoles() ([]models.Role, error) {
	query := "SELECT * FROM roles"

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for roles: %w", err)
	}
	defer stmt.Close()

	var roles []models.Role
	rows, err := stmt.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for roles: %w", err)
	}

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role value: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (repo *RoleRepository) GetRoleById(roleId int64) (*models.Role, error) {
	query := "SELECT * FROM roles WHERE id = $1"

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for role: %w", err)
	}
	defer stmt.Close()

	var role models.Role
	err = stmt.QueryRow(roleId).Scan(&role.ID, &role.Name, &role.Description, &role.Default)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("roles not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan role value: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) SetDefaultRole(roleId int64) error {
	// Begin a transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback unless commit succeeds
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Reset all roles to non-default
	resetQuery := "UPDATE roles SET default_role = FALSE WHERE default_role = TRUE"
	resetStmt, err := tx.Prepare(resetQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare reset default role query: %w", err)
	}
	defer resetStmt.Close()

	_, err = resetStmt.Exec()
	if err != nil {
		return fmt.Errorf("failed to reset default roles: %w", err)
	}

	// Set the specified role as default
	setQuery := "UPDATE roles SET default_role = TRUE WHERE id = $1"
	setStmt, err := tx.Prepare(setQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare set default role query: %w", err)
	}
	defer setStmt.Close()

	result, err := setStmt.Exec(roleId)
	if err != nil {
		return fmt.Errorf("failed to set default role for role %d: %w", roleId, err)
	}

	// Verify the role exists and was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role %d not found", roleId)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo *RoleRepository) GetDefaultRole() (*models.Role, error) {
	query := "SELECT * FROM roles WHERE default_role = TRUE"

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for default role: %w", err)
	}
	defer stmt.Close()

	var role models.Role
	err = stmt.QueryRow().Scan(&role.ID, &role.Name, &role.Description, &role.Default)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("default role not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan default role value: %w", err)
	}
	return &role, nil
}

func (repo *RoleRepository) Save(role *models.Role) (*models.Role, error) {
	query := "INSERT INTO roles (name, description, default_role) VALUES ($1, $2, $3) RETURNING id"
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for saving role: %w", err)
	}
	defer stmt.Close()

	var roleId int64
	err = stmt.QueryRow(role.Name, role.Description, role.Default).Scan(&roleId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for saving role: %w", err)
	}
	role.ID = roleId
	return role, nil
}

func (repo *RoleRepository) AssignRoleToUser(userId int64, roleId int64) error {
	query := "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)"
	stmt, err := repo.db.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare query for assigning role to user: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, roleId)

	if err != nil {
		return fmt.Errorf("failed to execute query for assigning role to user: %w", err)
	}

	return nil
}

func (repo *RoleRepository) GetRolesByUserId(userId int64) ([]models.Role, error) {
	query := "SELECT r.* FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1"

	stmt, err := repo.db.Prepare(query)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for getting roles by user: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting roles by user: %w", err)
	}

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role value: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil

}

func (repo *RoleRepository) GetUsersByRoleId(roleId int64) ([]models.User, error) {
	query := "SELECT u.* FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.role_id = $1"

	stmt, err := repo.db.Prepare(query)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for getting users by role: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(roleId)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting users by role: %w", err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Phone, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user value: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *RoleRepository) RemoveRoleFromUser(userId, roleId int64) error {
	// Begin transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Step 1: Check if the user has other roles
	var roleCount int
	checkQuery := `
        SELECT COUNT(*) 
        FROM user_roles 
        WHERE user_id = $1 AND role_id != $2`
	err = tx.QueryRow(checkQuery, userId, roleId).Scan(&roleCount)
	if err != nil {
		return fmt.Errorf("failed to count user roles: %w", err)
	}
	if roleCount == 0 {
		return fmt.Errorf("cannot remove role %d from user %d: user would be left with no roles", roleId, userId)
	}

	// Step 2: Remove the role from the user
	deleteQuery := `
        DELETE FROM user_roles 
        WHERE user_id = $1 AND role_id = $2`
	deleteStmt, err := tx.Prepare(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete role query: %w", err)
	}
	defer deleteStmt.Close()

	result, err := deleteStmt.Exec(userId, roleId)
	if err != nil {
		return fmt.Errorf("failed to remove role %d from user %d: %w", roleId, userId, err)
	}

	// Verify the role was removed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("role %d not assigned to user %d", roleId, userId)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo *RoleRepository) DeleteRole(roleId int64) error {
	defaultRole, err := repo.GetDefaultRole()
	if err != nil {
		return fmt.Errorf("failed to get default role: %w", err)
	}
	if defaultRole.ID == roleId {
		return fmt.Errorf("cannot delete default role")
	}

	users, err := repo.usersWithOnlyThisRole(roleId)
	if err != nil {
		return fmt.Errorf("failed to get users with only this role: %w", err)
	}

	if len(users) > 0 {
		err = repo.PatchAssignRoleToUsers(users, defaultRole.ID)
		if err != nil {
			return fmt.Errorf("failed to migrate all users from this role: %w", err)
		}
	}

	query := "DELETE FROM roles WHERE id = $1"
	stmt, err := repo.db.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare query for deleting role: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleId)

	if err != nil {
		return fmt.Errorf("failed to execute query for deleting role: %w", err)
	}

	return nil
}

func (repo *RoleRepository) usersWithOnlyThisRole(roleId int64) ([]models.User, error) {
	query := `
		SELECT u.*
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			WHERE ur.role_id = $1
			AND NOT EXISTS (
  				SELECT 1
    			FROM user_roles ur2
    			WHERE ur2.user_id = u.id
    			AND ur2.role_id != $1
		);
	`

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for getting users with only this role: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(roleId)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting users with only this role: %w", err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Phone, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user value: %w", err)
		}
		users = append(users, user)
	}

	return users, nil

}

func (repo *RoleRepository) PatchAssignRoleToUsers(users []models.User, roleId int64) error {
	if len(users) == 0 {
		return nil // Nothing to do
	}

	// Build a list of user IDs
	userIDs := make([]interface{}, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	// Construct the placeholders for the VALUES clause
	placeholders := make([]string, len(users))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Single query to assign the specified role to all users
	query := fmt.Sprintf(`
        INSERT INTO user_roles (user_id, role_id)
        SELECT u.user_id, $%d
        FROM (VALUES %s) AS u(user_id)
        ON CONFLICT (user_id, role_id) DO NOTHING`, len(users)+1, strings.Join(placeholders, ","))

	// Append the role ID to the parameters
	params := append(userIDs, roleId)

	// Execute the query directly with repo.db
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare assign role query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		return fmt.Errorf("failed to execute assign role query: %w", err)
	}

	return nil
}
