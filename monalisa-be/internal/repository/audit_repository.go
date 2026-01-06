package repository

import "database/sql"

type AuditRepository struct {
	DB *sql.DB
}

func (r *AuditRepository) Log(actorID, action, target string) error {
	_, err := r.DB.Exec(`
		INSERT INTO audit_logs (actor_id, action, target)
		VALUES ($1, $2, $3)
	`, actorID, action, target)

	return err
}
