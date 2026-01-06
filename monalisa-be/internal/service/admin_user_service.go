package service

import "errors"

type AdminUserService struct {
	userRepo  UserRepo
	roleRepo  RoleRepo
	auditRepo AuditRepo
}

func NewAdminUserService(u UserRepo, r RoleRepo, a AuditRepo) *AdminUserService {
	return &AdminUserService{
		userRepo:  u,
		roleRepo:  r,
		auditRepo: a,
	}
}

func (s *AdminUserService) ListUsers() (interface{}, error) {
	return s.userRepo.ListUsersWithRoles()
}

func (s *AdminUserService) AssignRole(actorID, userID, role string) error {
	ok, err := s.roleRepo.RoleExists(role)
	if err != nil || !ok {
		return errors.New("role not found")
	}

	// === PRIMARY ACTION (WAJIB SUKSES) ===
	if err := s.userRepo.AssignRole(userID, role); err != nil {
		return err
	}

	// === SECONDARY ACTION (TIDAK BOLEH GAGALKAN) ===
	if s.auditRepo != nil {
		_ = s.auditRepo.Log(
			actorID,
			"assign_role",
			"user:"+userID+" role:"+role,
		)
	}

	return nil
}

func (s *AdminUserService) RemoveRole(actorID, userID, role string) error {
	// proteksi admin-self
	if actorID == userID && role == "admin" {
		return errors.New("cannot remove admin role from yourself")
	}

	if err := s.userRepo.RemoveRole(userID, role); err != nil {
		return err
	}

	if s.auditRepo != nil {
		_ = s.auditRepo.Log(
			actorID,
			"remove_role",
			"user:"+userID+" role:"+role,
		)
	}

	return nil
}
