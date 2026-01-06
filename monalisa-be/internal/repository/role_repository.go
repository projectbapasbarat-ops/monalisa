package repository

import "database/sql"

type RoleRepository struct {
	DB *sql.DB
}

func (r *RoleRepository) RoleExists(code string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS (SELECT 1 FROM roles WHERE code = $1)
	`, code).Scan(&exists)
	return exists, err
}

func (r *RoleRepository) ListRoleCodes() ([]string, error) {
	rows, err := r.DB.Query(`SELECT code FROM roles ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		roles = append(roles, code)
	}
	return roles, nil
}
