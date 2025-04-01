package models

import (
	"time"

	"github.com/google/uuid"
)

// структура сообщения
type Message struct {
	ID           uuid.UUID `json:"id"`
	ChatID       uuid.UUID `json:"chat_id"`
	SenderID     uuid.UUID `json:"sender_id"`
	ReceiverID   uuid.UUID `json:"receiver_id"`
	Text         string    `json:"text"`
	SendAt       time.Time `json:"send_at"`
	Read         bool      `json:"read"`
	IsAttachment bool      `json:"is_attachment"`
	IsPinned     bool      `json:"is_pinned"`
}

// проверка на корректность
func (m *Message) IsValid() bool {
	// непустое сообщения
	if len(m.Text) == 0 {
		return false
	}

	return true
}
