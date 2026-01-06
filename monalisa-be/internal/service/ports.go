package service

import "monalisa-be/internal/model"

type UserRepo interface {
	GetUserAuthByNIP(nip string) (*model.UserAuth, error)
	ListUsersWithRoles() ([]model.UserWithRoles, error)
	AssignRole(userID, roleCode string) error
	RemoveRole(userID, roleCode string) error
}

type RoleRepo interface {
	RoleExists(code string) (bool, error)
	ListRoleCodes() ([]string, error)
}

type AuditRepo interface {
	Log(actorID, action, target string) error
}
