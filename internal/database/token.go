package database

import (
	"go-messenger/internal/models"

	"github.com/google/uuid"
)

// сохранение refresh token
func (pq *PgStorage) RefreshTokenSave(rt *models.RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens (id, token, user_id, expires_at, created_on, revoked)
	VALUES (:id, :token, :user_id, :expires_at, :created_on, :revoked)
	`
	_, err := pq.NamedExec(query, rt)

	return err
}

// получение refresh token
func (pq *PgStorage) RefreshTokenGet(refreshToken string) (*models.RefreshToken, error) {
	query := `
	SELECT * FROM refresh_tokens
	WHERE 1=1
		AND (token = $1)
		AND (expires_at > NOW())
		AND (revoked = false)
	`

	rt := &models.RefreshToken{}
	if err := pq.Get(rt, query, refreshToken); err != nil {
		return nil, err
	}

	return rt, nil
}

// revoke refresh token
func (pq *PgStorage) RefreshTokenRevoke(refreshToken string) error {
	query := `
	UPDATE refresh_tokens SET 
		revoked = TRUE
	WHERE token = :token
	`

	params := map[string]any{
		"token": refreshToken,
	}

	_, err := pq.NamedExec(query, params)

	return err
}

// revoke refresh token by user id
func (pq *PgStorage) RefreshTokenRevokeByID(id uuid.UUID) error {
	query := `
	UPDATE refresh_tokens SET 
		revoked = TRUE
	WHERE user_id = :id
	`

	params := map[string]any{
		"id": id.String(),
	}

	_, err := pq.NamedExec(query, params)

	return err
}
