package service

import "errors"

type AdminUserService struct {
	userRepo  UserRepo
	roleRepo  RoleRepo
	auditRepo AuditRepo
}

func NewAdminUserService(u UserRepo, r RoleRepo, a AuditRepo) *AdminUserService {
	return &AdminUserService{u, r, a}
}

func (s *AdminUserService) ListUsers() (interface{}, error) {
	return s.userRepo.ListUsersWithRoles()
}

func (s *AdminUserService) AssignRole(actorID, userID, role string) error {
	ok, err := s.roleRepo.RoleExists(role)
	if err != nil || !ok {
		return errors.New("role not found")
	}
	if err := s.userRepo.AssignRole(userID, role); err != nil {
		return err
	}
	return s.auditRepo.Log(actorID, "assign_role", role)
}

func (s *AdminUserService) RemoveRole(actorID, userID, role string) error {
	return s.userRepo.RemoveRole(userID, role)
}
