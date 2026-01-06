package handler

import (
	"monalisa-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	service *service.AdminUserService
}

func NewAdminUserHandler(s *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{s}
}

func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users})
}

func (h *AdminUserHandler) AssignRole(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RoleCode string `json:"role_code"`
	}
	c.ShouldBindJSON(&req)

	adminID := c.GetString("user_id")
	if err := h.service.AssignRole(adminID, id, req.RoleCode); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "role assigned"})
}

func (h *AdminUserHandler) RemoveRole(c *gin.Context) {
	adminID := c.GetString("user_id")
	id := c.Param("id")
	role := c.Param("role")

	if err := h.service.RemoveRole(adminID, id, role); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "role removed"})
}
