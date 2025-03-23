package models

type UserRoleRequest struct {
	UserId int64 `json:"user_id" binding:"required"`
	RoleId int64 `json:"role_id" binding:"required"`
}
