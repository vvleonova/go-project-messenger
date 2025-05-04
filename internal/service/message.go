package service

import (
	"fmt"
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/http"
	"go-messenger/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// обработка ответов сервера
const (
	errNoUser         = "user with phone %s doesn't exist"
	errNoChat         = "unable to get chat for users %s and %s"
	errCreateID       = "unable to create message uuid"
	infoUserNotOnline = "user with phone %s is not online, message is saved to database"
	errMsgUnread      = "unable to load unread messages"
)

// обработка и сохранение сообщения
func (s *Service) HandleMessage(msgRequest *dto.MessageRequest, userConn *http.UserConn, clients *sync.Map) {
	// проверка, что receiver существует в базе данных
	userReceiver, err := s.storage.UserGetPhone(msgRequest.ReceiverPhone)
	if err != nil {
		s.logger.Info(fmt.Sprintf(errNoUser, msgRequest.ReceiverPhone), err)
		userConn.Send <- fmt.Sprintf(errNoUser, msgRequest.ReceiverPhone)

		return
	}

	// получение id чата
	chatID, err := s.getChatID(userConn.UserID, userReceiver.ID)
	if err != nil {
		s.logger.Info(fmt.Sprintf(errNoChat, userConn.UserID, userReceiver.ID), err)
		userConn.Send <- fmt.Sprintf(errNoChat, userConn.UserID, userReceiver.ID)

		return
	}

	// обработка сообщения
	msg := &models.Message{}

	msg.ID, err = uuid.NewRandom()
	if err != nil {
		s.logger.Info(errCreateID, err)
		userConn.Send <- errCreateID
		return
	}

	msg.ChatID = chatID
	msg.SenderID = userConn.UserID
	msg.ReceiverID = userReceiver.ID
	msg.Text = msgRequest.Text
	msg.SendAt = time.Now()
	msg.Read = false
	msg.IsAttachment = false
	msg.IsPinned = false

	// сохранение сообщения в базу данных
	err = s.storage.MessageInsert(msg)
	if err != nil {
		s.logger.Error(err.Error(), nil)
		userConn.Send <- err.Error()
		return
	}

	// отправка сообщения, если пользователь онлайн
	if receiverConn, ok := clients.Load(userReceiver.ID); ok {
		// отправка сообщения в канал
		receiver := receiverConn.(*http.UserConn)
		receiver.Send <- msgRequest.Text

		// обновление статуса прочитано
		s.storage.MessageRead(msg.ID)
	} else {
		s.logger.Info(fmt.Sprintf(infoUserNotOnline, msgRequest.ReceiverPhone), nil)
		userConn.Send <- fmt.Sprintf(infoUserNotOnline, msgRequest.ReceiverPhone)
	}
}

// обработка и сохранение сообщения
func (s *Service) HandleMessageUnread(userConn *http.UserConn) {
	messages, err := s.storage.MessageGetUnread(userConn.UserID)
	if err != nil {
		s.logger.Info(errMsgUnread, err)
		userConn.Send <- fmt.Sprintf("%s: %s", errMsgUnread, err.Error())
		return
	}

	// отправка сообщения
	for _, msg := range messages {
		// отправка сообщения в канал
		userConn.Send <- msg.Text

		// обновление статуса прочитано
		s.storage.MessageRead(msg.ID)
	}
}
