package models

import (
	"time"

	"github.com/google/uuid"
)

// структура чата
type Chat struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedOn time.Time `db:"created_on"`
	// Photo
}

// структура чата для показа пользователю
type ChatPreview struct {
	ID            uuid.UUID `db:"id" json:"id"`
	CompanionID   uuid.UUID `db:"companion_id" json:"companion_id"`
	LastMessage   string    `db:"last_message_text" json:"last_message_text"`
	LastMessageAt string    `db:"last_message_time" json:"last_message_time"`
}
