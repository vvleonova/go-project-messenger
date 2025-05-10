package service

import (
	"errors"
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_getChatID(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name          string
		storage       *storageMock
		inputSender   uuid.UUID
		inputReceiver uuid.UUID
		wantErr       bool
		wantErrText   string
	}{
		{
			name: "chat exists",
			storage: &storageMock{
				ChatUsersGetFunc: func(sID uuid.UUID, rID uuid.UUID) (*models.Chat, error) {
					assert.Equal(t, senderID, sID)
					assert.Equal(t, receiverID, rID)
					return &models.Chat{
						ID:        chatID,
						CreatedOn: time.Now(),
					}, nil
				},
				ChatInsertFunc:      nil,
				ChatUsersInsertFunc: nil,
			},
			inputSender:   senderID,
			inputReceiver: receiverID,
		},
		{
			name: "chat created",
			storage: &storageMock{
				ChatUsersGetFunc: func(sID uuid.UUID, rID uuid.UUID) (*models.Chat, error) {
					assert.Equal(t, senderID, sID)
					assert.Equal(t, receiverID, rID)
					return nil, nil
				},
				ChatInsertFunc: func(c *models.Chat) error {
					return nil
				},
				ChatUsersInsertFunc: func(chatID, userID uuid.UUID) error {
					return nil
				},
			},
			inputSender:   senderID,
			inputReceiver: receiverID,
		},
		{
			name: "error checking chat existence",
			storage: &storageMock{
				ChatUsersGetFunc: func(sID uuid.UUID, rID uuid.UUID) (*models.Chat, error) {
					assert.Equal(t, senderID, sID)
					assert.Equal(t, receiverID, rID)
					return nil, errors.New(errGetDB)
				},
				ChatInsertFunc:      nil,
				ChatUsersInsertFunc: nil,
			},
			inputSender:   senderID,
			inputReceiver: receiverID,
			wantErr:       true,
			wantErrText:   errGetDB,
		},
		{
			name: "error saving chat",
			storage: &storageMock{
				ChatUsersGetFunc: func(sID uuid.UUID, rID uuid.UUID) (*models.Chat, error) {
					assert.Equal(t, senderID, sID)
					assert.Equal(t, receiverID, rID)
					return nil, nil
				},
				ChatInsertFunc: func(c *models.Chat) error {
					return errors.New(errSaveDB)
				},
				ChatUsersInsertFunc: nil,
			},
			inputSender:   senderID,
			inputReceiver: receiverID,
			wantErr:       true,
			wantErrText:   errSaveDB,
		},
		{
			name: "error saving chat + users",
			storage: &storageMock{
				ChatUsersGetFunc: func(sID uuid.UUID, rID uuid.UUID) (*models.Chat, error) {
					assert.Equal(t, senderID, sID)
					assert.Equal(t, receiverID, rID)
					return nil, nil
				},
				ChatInsertFunc: func(c *models.Chat) error {
					return nil
				},
				ChatUsersInsertFunc: func(chatID, userID uuid.UUID) error {
					return errors.New(errSaveDB)
				},
			},
			inputSender:   senderID,
			inputReceiver: receiverID,
			wantErr:       true,
			wantErrText:   errSaveDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			got, err := s.getChatID(tt.inputSender, tt.inputReceiver)

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.wantErrText)
				assert.Empty(t, got)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestService_GetUserChats(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name        string
		storage     *storageMock
		input       uuid.UUID
		wantErr     bool
		wantErrText string
	}{
		{
			name: "got chats",
			storage: &storageMock{
				ChatsUserGetFunc: func(id uuid.UUID) ([]models.ChatPreview, error) {
					assert.Equal(t, correctID, id)
					return []models.ChatPreview{
						{
							ID:          senderID,
							CompanionID: receiverID,
							LastMessage: messageText,
						},
						{
							ID:          receiverID,
							CompanionID: senderID,
							LastMessage: messageSecondText,
						},
					}, nil
				},
			},
			input: correctID,
		},
		{
			name: "error getting chats",
			storage: &storageMock{
				ChatsUserGetFunc: func(id uuid.UUID) ([]models.ChatPreview, error) {
					assert.Equal(t, correctID, id)
					return nil, errors.New(errGetDB)
				},
			},
			input:       correctID,
			wantErr:     true,
			wantErrText: errGetDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			chats, err := s.GetUserChats(tt.input)

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.wantErrText)
				assert.Empty(t, chats)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, chats)
				assert.Len(t, chats, 2)
				assert.Equal(t, senderID, chats[0].ID)
				assert.Equal(t, receiverID, chats[0].CompanionID)
				assert.Equal(t, messageText, chats[0].LastMessage)
			}
		})
	}
}

func TestService_GetUsersChatMessages(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name             string
		storage          *storageMock
		inputUserID      uuid.UUID
		inputCompanionID uuid.UUID
		wantErr          bool
		wantErrText      string
	}{
		{
			name: "got chat messages",
			storage: &storageMock{
				ChatMessagesGetFunc: func(userID, companionID uuid.UUID) ([]dto.MessageResponse, error) {
					assert.Equal(t, senderID, userID)
					assert.Equal(t, receiverID, companionID)
					return []dto.MessageResponse{
						{Text: messageText},
						{Text: messageSecondText},
					}, nil
				},
			},
			inputUserID:      senderID,
			inputCompanionID: receiverID,
		},
		{
			name: "error getting chat messages",
			storage: &storageMock{
				ChatMessagesGetFunc: func(userID, companionID uuid.UUID) ([]dto.MessageResponse, error) {
					assert.Equal(t, senderID, userID)
					assert.Equal(t, receiverID, companionID)
					return nil, errors.New(errGetDB)
				},
			},
			inputUserID:      senderID,
			inputCompanionID: receiverID,
			wantErr:          true,
			wantErrText:      errGetDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			messages, err := s.GetUsersChatMessages(tt.inputUserID, tt.inputCompanionID)

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.wantErrText)
				assert.Empty(t, messages)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, messages)
				assert.Len(t, messages, 2)
				assert.Equal(t, messageText, messages[0].Text)
			}
		})
	}
}
