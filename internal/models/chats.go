package models

import "github.com/google/uuid"

// структура списка чатов
type ChatList struct {
	UserID uuid.UUID   `json:"user_id"`
	Chats  []uuid.UUID `json:"chats"`
}
