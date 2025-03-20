package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
	"github.com/wmfadel/escape-be/internal/service"
)

type AdmingHandler struct {
	RolesService *service.RoleService
	UserService  *service.UserService
}

func NewAdmingHandler(rolesService *service.RoleService, userService *service.UserService) *AdmingHandler {
	return &AdmingHandler{RolesService: rolesService, UserService: userService}
}

func (h *AdmingHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.RolesService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *AdmingHandler) GetRoleById(c *gin.Context) {
	roleId, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to get role", nil))
	}

	roleIdInt, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse role ID", err))
		return
	}

	role, err := h.RolesService.GetRoleById(roleIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *AdmingHandler) SetDefaultRole(c *gin.Context) {
	roleId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to get role", err))
	}

	err = h.RolesService.SetDefaultRole(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role set as default"})
}

func (h *AdmingHandler) GetDefaultRoles(c *gin.Context) {
	role, err := h.RolesService.GetDefaultRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *AdmingHandler) AddRole(c *gin.Context) {

	var role *models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse role", err))
		return
	}

	role, err := h.RolesService.Save(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to add role", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *AdmingHandler) AssignRoleToUser(c *gin.Context) {
	var assignRoleRequest models.UserRoleRequest
	if err := c.ShouldBindJSON(&assignRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse role", err))
		return
	}

	err := h.RolesService.AssignRoleToUser(assignRoleRequest.UserId, assignRoleRequest.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to add role", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user"})
}

func (h *AdmingHandler) GetRolesByUserId(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse user ID", err))
		return
	}

	roles, err := h.RolesService.GetRolesByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *AdmingHandler) GetUsersByRoleId(c *gin.Context) {
	roleId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse user ID", err))
		return
	}

	users, err := h.RolesService.GetUsersByRoleId(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdmingHandler) GetAllAdmins(c *gin.Context) {

	users, err := h.RolesService.GetUsersByRoleId(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdmingHandler) GetAllOrganizers(c *gin.Context) {
	users, err := h.RolesService.GetUsersByRoleId(2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdmingHandler) RemoveRoleFromUser(c *gin.Context) {
	var deleteRoleRequest models.UserRoleRequest
	if err := c.ShouldBindJSON(&deleteRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse role", err))
		return
	}

	err := h.RolesService.RemoveRoleFromUser(deleteRoleRequest.UserId, deleteRoleRequest.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to add role", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed from user"})
}

func (h *AdmingHandler) DeleteRole(c *gin.Context) {
	roleId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to get role", err))
		return
	}

	err = h.RolesService.DeleteRole(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to get roles", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})

}

func (h *AdmingHandler) PatchAssignRoleToUsers(c *gin.Context) {
	var patchUsersRoleRequest models.PatchUserRoleRequest

	if err := c.ShouldBindJSON(&patchUsersRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse role", err))
		return
	}

	err := h.RolesService.PatchAssignRoleToUsers(patchUsersRoleRequest.UserIds, patchUsersRoleRequest.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("Failed to add role", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to users"})

}

func (h *AdmingHandler) BlockUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewESError("Failed to parse user ID", err))
		return
	}

	err = h.RolesService.DeleteUserRoles(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewESError("failed to bloc user", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user blocked"})

}
