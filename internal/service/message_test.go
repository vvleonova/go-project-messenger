package service

import (
	"errors"
	"fmt"
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/http"
	"go-messenger/internal/logger"
	"go-messenger/internal/models"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	receiverPhone     = "79999999999"
	messageText       = "text message"
	messageSecondText = "new text message"

	errSendMessage = "message was not sent"
)

var (
	senderID        = uuid.New()
	receiverID      = uuid.New()
	messageFirstID  = uuid.New()
	messageSecondID = uuid.New()
	senderSend      = make(chan string, 10)
	receiverSend    = make(chan string, 10)
)

func TestService_HandleMessage(t *testing.T) {
	// входные данные для теста
	msgRequest := &dto.MessageRequest{
		ReceiverPhone: receiverPhone,
		Text:          messageText,
	}

	// возможные кейсы
	tests := []struct {
		name               string
		storage            *storageMock
		clients            *sync.Map
		input              *dto.MessageRequest
		receiverMessageGet bool
		senderMessageGet   bool
		senderMessageText  string
	}{
		{
			name: "message sent, receiver online",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, receiverPhone, phone)
					return &models.User{
						ID:    receiverID,
						Phone: receiverPhone,
					}, nil
				},
				MessageInsertFunc: func(u *models.Message) error {
					assert.NotEmpty(t, u.ID)
					assert.Equal(t, senderID, u.SenderID)
					assert.Equal(t, receiverID, u.ReceiverID)
					assert.Equal(t, messageText, u.Text)
					assert.WithinDuration(t, time.Now(), u.SendAt, 2*time.Second)
					assert.False(t, u.Read)
					assert.False(t, u.IsAttachment)
					assert.False(t, u.IsPinned)
					return nil
				},
				MessageReadFunc: func(messageID uuid.UUID) error {
					assert.NotEmpty(t, messageID)
					return nil
				},
			},
			clients: func() *sync.Map {
				m := &sync.Map{}
				m.Store(receiverID, &http.UserConn{
					UserID: receiverID,
					Send:   receiverSend,
				})
				return m
			}(),
			input:              msgRequest,
			receiverMessageGet: true,
		},
		{
			name: "message sent, receiver not online",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, receiverPhone, phone)
					return &models.User{
						ID:    receiverID,
						Phone: receiverPhone,
					}, nil
				},
				MessageInsertFunc: func(u *models.Message) error {
					assert.NotEmpty(t, u.ID)
					assert.Equal(t, senderID, u.SenderID)
					assert.Equal(t, receiverID, u.ReceiverID)
					assert.Equal(t, messageText, u.Text)
					assert.WithinDuration(t, time.Now(), u.SendAt, 2*time.Second)
					assert.False(t, u.Read)
					assert.False(t, u.IsAttachment)
					assert.False(t, u.IsPinned)
					return nil
				},
				MessageReadFunc: nil,
			},
			clients: func() *sync.Map {
				m := &sync.Map{}
				m.Store(senderID, &http.UserConn{
					UserID: senderID,
					Send:   senderSend,
				})
				return m
			}(),
			input:             msgRequest,
			senderMessageGet:  true,
			senderMessageText: fmt.Sprintf(infoUserNotOnline, receiverPhone),
		},
		{
			name: "receiver doesn't exist",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					return nil, fmt.Errorf(errNoUser, receiverPhone)
				},
				MessageInsertFunc: nil,
				MessageReadFunc:   nil,
			},
			clients:           &sync.Map{},
			input:             msgRequest,
			senderMessageGet:  true,
			senderMessageText: fmt.Sprintf(errNoUser, receiverPhone),
		},
		{
			name: "message is not saved to database",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, receiverPhone, phone)
					return &models.User{
						ID:    receiverID,
						Phone: receiverPhone,
					}, nil
				},
				MessageInsertFunc: func(u *models.Message) error {
					return errors.New(errSaveDB)
				},
				MessageReadFunc: nil,
			},
			clients:           &sync.Map{},
			input:             msgRequest,
			senderMessageGet:  true,
			senderMessageText: errSaveDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
				logger:  logger.NewDummyLogger(),
			}

			// создание подключения пользователя
			userConn := &http.UserConn{
				UserID: senderID,
				Send:   senderSend,
			}

			// вызов тестируемой функции
			s.HandleMessage(tt.input, userConn, tt.clients)

			// проверка результата
			if tt.receiverMessageGet {
				receiverConn, ok := tt.clients.Load(receiverID)
				require.True(t, ok)
				receiver := receiverConn.(*http.UserConn)
				select {
				case messageReceived := <-receiver.Send:
					assert.Equal(t, messageText, messageReceived)
				case <-time.After(time.Second):
					t.Fatal(errSendMessage)
				}
			}

			if tt.senderMessageGet {
				select {
				case messageReceived := <-userConn.Send:
					assert.Equal(t, tt.senderMessageText, messageReceived)
				case <-time.After(time.Second):
					t.Fatal(errSendMessage)
				}
			}
		})
	}
}

func TestService_HandleMessageUnread(t *testing.T) {
	// список сообщений
	msg1 := models.Message{
		ID:   messageFirstID,
		Text: messageText,
	}
	msg2 := models.Message{
		ID:   messageSecondID,
		Text: messageSecondText,
	}

	// возможные кейсы
	tests := []struct {
		name             string
		storage          *storageMock
		input            uuid.UUID
		wantError        bool
		expectedMessages []string
	}{
		{
			name: "sent unread messages",
			storage: &storageMock{
				MessageGetUnreadFunc: func(id uuid.UUID) ([]models.Message, error) {
					assert.Equal(t, receiverID, id)
					return []models.Message{msg1, msg2}, nil
				},
				MessageReadFunc: func(messageID uuid.UUID) error {
					assert.NotEmpty(t, messageID)
					return nil
				},
			},
			input:            receiverID,
			expectedMessages: []string{messageText, messageSecondText},
		},
		{
			name: "error getting unread messages",
			storage: &storageMock{
				MessageGetUnreadFunc: func(id uuid.UUID) ([]models.Message, error) {
					return nil, errors.New(errMsgUnread)
				},
				MessageReadFunc: nil,
			},
			input:            receiverID,
			wantError:        true,
			expectedMessages: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
				logger:  logger.NewDummyLogger(),
			}

			// создание подключения пользователя
			userConn := &http.UserConn{
				UserID: receiverID,
				Send:   receiverSend,
			}

			// вызов тестируемой функции
			s.HandleMessageUnread(userConn)

			// проверка результата
			if tt.wantError {
				select {
				case messageReceived := <-userConn.Send:
					assert.Equal(t, errMsgUnread, messageReceived)
				case <-time.After(time.Second):
					t.Fatal(errSendMessage)
				}
			} else {
				received := []string{}
				for range tt.expectedMessages {
					select {
					case msg := <-userConn.Send:
						received = append(received, msg)
					case <-time.After(time.Second):
						t.Fatal(errSendMessage)
					}
				}

				assert.ElementsMatch(t, tt.expectedMessages, received)
			}
		})
	}
}
