package handler

import (
	"net/http"

	"monalisa-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminRoleHandler struct {
	service *service.AdminRoleService
}

func NewAdminRoleHandler(s *service.AdminRoleService) *AdminRoleHandler {
	return &AdminRoleHandler{s}
}

func (h *AdminRoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
