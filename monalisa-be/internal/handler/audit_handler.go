package handler

import (
	"net/http"
	"strconv"

	"monalisa-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler(s *service.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

func (h *AuditHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))

	data, err := h.service.List(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
