package models

import (
	"time"

	"github.com/google/uuid"
)

// структура refresh JWT-token
type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedOn time.Time `db:"created_on"`
	Revoked   bool      `db:"revoked"`
}
