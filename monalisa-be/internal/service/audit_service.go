package service

import "monalisa-be/internal/model"

type AuditRepository interface {
	List(limit int) ([]model.AuditLog, error)
}

type AuditService struct {
	repo AuditRepository
}

func NewAuditService(r AuditRepository) *AuditService {
	return &AuditService{repo: r}
}

func (s *AuditService) List(limit int) ([]model.AuditLog, error) {
	return s.repo.List(limit)
}
