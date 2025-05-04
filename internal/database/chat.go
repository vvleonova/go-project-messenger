package database

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"
)

// получение чата между двумя пользователями
func (pq *PgStorage) ChatUsersGet(senderID, receiverID uuid.UUID) (*models.Chat, error) {
	query := `
	SELECT chat_id AS id
	FROM chat_users
	GROUP BY chat_id
	HAVING
		COUNT(DISTINCT user_id) = 2
		AND SUM(CASE WHEN user_id IN ($1, $2) THEN 1 ELSE 0 END) = 2
	`

	chat := &models.Chat{}
	err := pq.Get(chat, query, senderID, receiverID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return chat, nil
}

// сохранение чата
func (pq *PgStorage) ChatInsert(c *models.Chat) error {
	query := `
	INSERT INTO chats (id, created_on)
	VALUES (:id, :created_on)
	`
	_, err := pq.NamedExec(query, c)

	return err
}

// сохранение чата + пользователей
func (pq *PgStorage) ChatUsersInsert(chatID, userID uuid.UUID) error {
	query := `
	INSERT INTO chat_users (chat_id, user_id)
	VALUES (:chat_id, :user_id)
	`
	params := map[string]any{
		"chat_id": chatID,
		"user_id": userID,
	}

	_, err := pq.NamedExec(query, params)

	return err
}

// получение последних чатов пользователя
func (pq *PgStorage) ChatsUserGet(userID uuid.UUID) ([]models.ChatPreview, error) {
	query := `
	WITH last_messages AS (
		SELECT
			chat_id,
			CASE 
				WHEN sender_id != $1 THEN sender_id 
				ELSE receiver_id 
			END AS companion_id,
			text AS last_message_text,
			TO_CHAR(send_at, 'DD.MM.YYYY HH24:MI:SS') AS last_message_time,
			ROW_NUMBER() OVER (PARTITION BY chat_id ORDER BY send_at DESC) AS rn
		FROM messages
	)
	SELECT
		cu.chat_id AS id,
		m.companion_id,
		m.last_message_text,
		m.last_message_time
	FROM chat_users cu
	JOIN last_messages m
		ON cu.chat_id=m.chat_id AND m.rn = 1
	WHERE cu.user_id = $1
	ORDER BY m.last_message_time DESC
	LIMIT 10
	`

	chats := []models.ChatPreview{}
	if err := pq.Select(&chats, query, userID); err != nil {
		return nil, err
	}

	return chats, nil
}

func (pq *PgStorage) ChatMessagesGet(userID, companionID uuid.UUID) ([]dto.MessageResponse, error) {
	query := `
	SELECT
		text,
		TO_CHAR(send_at, 'DD.MM.YYYY HH24:MI:SS') AS send_at,
		"read"
	FROM messages
	WHERE 
		(sender_id = $1 AND receiver_id = $2)
		OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY send_at DESC
	LIMIT 10
	`

	messages := []dto.MessageResponse{}
	if err := pq.Select(&messages, query, userID, companionID); err != nil {
		return nil, err
	}

	return messages, nil
}
