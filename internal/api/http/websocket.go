package http

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// преобразование HTTP-запроса в WebSocket-соединение
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// структура, хранящая соединение и канал для отправки сообщений
type UserConn struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan string
}

// создание WebSocket-соединения
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		c.Abort()
		return
	}
	defer conn.Close()

	// получение id пользователя
	tokenID, exists := c.Get("id")
	if !exists {
		c.Error(errs.Unauthorized("unable to get user ID"))
		c.Abort()
		return
	}

	userID, err := uuid.Parse(tokenID.(string))
	if err != nil {
		c.Error(errs.Unauthorized("unable to get user ID"))
		c.Abort()
		return
	}

	// сохранение подключения пользователя
	u := &UserConn{
		Conn:   conn,
		Send:   make(chan string, 10),
		UserID: userID,
	}

	s.clients.Store(userID, u)

	// отправка сообщений
	go s.writeMessage(u)

	// отправка непрочитанных сообщений
	go s.sendMessageUnread(u)

	// чтение сообщений
	s.readMessage(u)
}

// отправка данных клиенту
func (s *Server) writeMessage(u *UserConn) {
	// закрытие соединения
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()

	// обработка сообщений
	for {
		select {
		case msg, ok := <-u.Send:
			// закрытый канал
			if !ok {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // тайм-аут на запись

			// отправка сообщения
			if err := u.Conn.WriteJSON(msg); err != nil {
				s.logger.Error("writing message error", err)
				return
			}
		case <-ticker.C: // отправка ping - поддержка активности соединения
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.logger.Error("ping error", err)
				return
			}
		}
	}
}

// чтение непрочитанных сообщений клиенту
func (s *Server) sendMessageUnread(u *UserConn) {
	s.service.HandleMessageUnread(u)
}

// чтение потока данных от клиента
func (s *Server) readMessage(u *UserConn) {
	// закрытие соединения
	defer func() {
		u.Conn.Close()
		s.clients.Delete(u.UserID)
	}()

	// параметры соединения
	u.Conn.SetReadLimit(512)                                 // максимальный размер одного входящего сообщения
	u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // таймер на чтение сообщения
	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	}) // обработчик для входящих Pong-сообщений

	// обработка сообщений
	for {
		// чтение сообщения
		var msg dto.MessageRequest
		err := u.Conn.ReadJSON(&msg)
		if err != nil {
			s.logger.Error("reading message error", err)
			break
		}

		// проверка переданного запроса
		if err := msg.ValidateMessage(); err != nil {
			s.logger.Error(err.Error(), nil)
			return
		}

		// обработка и сохранение сообщения
		go s.service.HandleMessage(&msg, u, &s.clients)
	}
}
