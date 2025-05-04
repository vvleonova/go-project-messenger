package dto

import (
	"errors"
)

// обработка ошибок
var errNoText = errors.New("user login is not long enough")

// стуктура сообщения
type MessageRequest struct {
	ReceiverPhone string `json:"to"`
	Text          string `json:"text"`
}

// стуктура сообщения для показа пользователю
type MessageResponse struct {
	Text   string `json:"text" db:"text"`
	SendAt string `json:"send_at" db:"send_at"`
	Read   bool   `json:"read" db:"read"`
}

// валидация сообщения
func (msg MessageRequest) ValidateMessage() error {
	// текст не пустой
	if len(msg.Text) == 0 {
		return errNoText
	}

	return nil
}
