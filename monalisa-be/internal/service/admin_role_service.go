package service

type AdminRoleService struct {
	roleRepo RoleRepo
}

func NewAdminRoleService(r RoleRepo) *AdminRoleService {
	return &AdminRoleService{r}
}

func (s *AdminRoleService) ListRoles() ([]string, error) {
	return s.roleRepo.ListRoleCodes()
}
