package database

import (
	"go-messenger/internal/models"

	"github.com/google/uuid"
)

// сохранение сообщения
func (pq *PgStorage) MessageInsert(u *models.Message) error {
	query := `
	INSERT INTO messages (id, chat_id, sender_id, receiver_id, text, send_at, read, is_attachment, is_pinned)
	VALUES (:id, :chat_id, :sender_id, :receiver_id, :text, :send_at, :read, :is_attachment, :is_pinned)
	`
	_, err := pq.NamedExec(query, u)

	return err
}

// обновление статуса прочитано
func (pq *PgStorage) MessageRead(messageID uuid.UUID) error {
	query := "UPDATE messages SET read = True WHERE id = $1"

	_, err := pq.Exec(query, messageID)

	return err
}

// получение непрочитанных сообщений
func (pq *PgStorage) MessageGetUnread(receiverID uuid.UUID) ([]models.Message, error) {
	query := "SELECT * FROM messages WHERE receiver_id = $1 AND read = False"

	messages := []models.Message{}
	if err := pq.Select(&messages, query, receiverID); err != nil {
		return nil, err
	}

	return messages, nil
}
