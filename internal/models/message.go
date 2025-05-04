package models

import (
	"time"

	"github.com/google/uuid"
)

// структура сообщения
type Message struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ChatID       uuid.UUID `json:"chat_id" db:"chat_id"`
	SenderID     uuid.UUID `json:"sender_id" db:"sender_id"`
	ReceiverID   uuid.UUID `json:"receiver_id" db:"receiver_id"`
	Text         string    `json:"text" db:"text"`
	SendAt       time.Time `json:"send_at" db:"send_at"`
	Read         bool      `json:"read" db:"read"`
	IsAttachment bool      `json:"is_attachment" db:"is_attachment"`
	IsPinned     bool      `json:"is_pinned" db:"is_pinned"`
}
