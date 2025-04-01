package models

import "github.com/google/uuid"

// структура чата
type Chat struct {
	ID       uuid.UUID   `json:"id"`
	Name     string      `json:"name"`
	UsersIDs []uuid.UUID `json:"user_ids"`
	Messages []uuid.UUID `json:"messages"`
	// Photo
}
