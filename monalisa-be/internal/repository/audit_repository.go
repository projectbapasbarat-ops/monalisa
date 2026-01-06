package repository

import (
	"database/sql"

	"monalisa-be/internal/model"
)

type AuditRepository struct {
	DB *sql.DB
}

/*
WRITE AUDIT LOG
Dipakai oleh AdminUserService
*/
func (r *AuditRepository) Log(actorID, action, target string) error {
	_, err := r.DB.Exec(`
		INSERT INTO audit_logs (actor_id, action, target)
		VALUES ($1, $2, $3)
	`, actorID, action, target)

	return err
}

/*
READ AUDIT LOG
Dipakai oleh AuditService
*/
func (r *AuditRepository) List(limit int) ([]model.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.DB.Query(`
		SELECT id, actor_id, action, target, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		if err := rows.Scan(
			&l.ID,
			&l.ActorID,
			&l.Action,
			&l.Target,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}
