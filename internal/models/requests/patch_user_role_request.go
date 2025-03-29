package requests

type PatchUserRoleRequest struct {
	UserIds []any `json:"user_ids" binding:"required"`
	RoleId  int64 `json:"role_id" binding:"required"`
}
