package repository

import (
	"database/sql"
	"time"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func (r *RefreshTokenRepository) Save(userID, token string, expiresAt time.Time) error {
	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)

	return err
}

func (r *RefreshTokenRepository) FindValid(token string) (string, error) {
	var userID string
	err := r.DB.QueryRow(`
		SELECT user_id
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false AND expires_at > now()
	`, token).Scan(&userID)

	return userID, err
}

func (r *RefreshTokenRepository) Revoke(token string) error {
	_, err := r.DB.Exec(`
		UPDATE refresh_tokens SET revoked = true WHERE token = $1
	`, token)
	return err
}
