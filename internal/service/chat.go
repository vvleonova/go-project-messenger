package service

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"
	"time"

	"github.com/google/uuid"
)

func (s *Service) getChatID(senderID, receiverID uuid.UUID) (uuid.UUID, error) {
	// проверка, что чат уже существует
	chat, err := s.storage.ChatUsersGet(senderID, receiverID)
	if err != nil {
		return uuid.UUID{}, err
	}

	// создание чата, если он не существует
	if chat == nil {
		chat = &models.Chat{}

		chat.ID = uuid.New()
		chat.CreatedOn = time.Now()

		err = s.storage.ChatInsert(chat)
		if err != nil {
			return uuid.UUID{}, err
		}

		err = s.storage.ChatUsersInsert(chat.ID, senderID)
		if err != nil {
			return uuid.UUID{}, err
		}

		err = s.storage.ChatUsersInsert(chat.ID, receiverID)
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	return chat.ID, nil
}

func (s *Service) GetUserChats(userID uuid.UUID) ([]models.ChatPreview, error) {
	chats, err := s.storage.ChatsUserGet(userID)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *Service) GetUsersChatMessages(userID, companionID uuid.UUID) ([]dto.MessageResponse, error) {
	messages, err := s.storage.ChatMessagesGet(userID, companionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
