package repository

import (
	"database/sql"
	"errors"

	"monalisa-be/internal/model"

	"github.com/lib/pq"
)

type UserRepository struct {
	DB *sql.DB
}

/* =========================
   AUTH
========================= */

func (r *UserRepository) GetUserAuthByNIP(nip string) (*model.UserAuth, error) {
	var userID string

	err := r.DB.QueryRow(`
		SELECT id FROM users
		WHERE nip = $1 AND is_active = true
	`, nip).Scan(&userID)

	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(`
		SELECT DISTINCT p.code
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}

	return &model.UserAuth{
		UserID:      userID,
		Permissions: perms,
	}, nil
}

/* =========================
   ADMIN
========================= */

func (r *UserRepository) ListUsersWithRoles() ([]model.UserWithRoles, error) {
	rows, err := r.DB.Query(`
		SELECT 
			u.id,
			u.nip,
			e.nama,
			e.jabatan,
			COALESCE(
				array_agg(r.code) FILTER (WHERE r.code IS NOT NULL),
				'{}'
			) AS roles
		FROM users u
		JOIN employees e ON e.nip = u.nip
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		GROUP BY u.id, e.nama, e.jabatan
		ORDER BY e.nama
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.UserWithRoles

	for rows.Next() {
		var u model.UserWithRoles
		var roles []string

		if err := rows.Scan(
			&u.ID,
			&u.NIP,
			&u.Nama,
			&u.Jabatan,
			pq.Array(&roles),
		); err != nil {
			return nil, err
		}

		u.Roles = roles
		result = append(result, u)
	}

	return result, nil
}

func (r *UserRepository) AssignRole(userID, roleCode string) error {
	_, err := r.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE code = $2
		ON CONFLICT DO NOTHING
	`, userID, roleCode)
	return err
}

func (r *UserRepository) RemoveRole(userID, roleCode string) error {
	res, err := r.DB.Exec(`
		DELETE FROM user_roles
		WHERE user_id = $1
		  AND role_id = (SELECT id FROM roles WHERE code = $2)
	`, userID, roleCode)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("role not assigned")
	}
	return nil
}
