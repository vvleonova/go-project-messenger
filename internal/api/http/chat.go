package http

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение списка последних чатов пользователя из базы данных
func (s *Server) userGetChats(c *gin.Context) {
	// получение id пользователя
	tokenID, exists := c.Get("id")
	if !exists {
		c.Error(errs.Unauthorized("unable to get user ID"))
		return
	}
	userID := tokenID.(uuid.UUID)

	// получение списка последних чатов
	chats, err := s.service.GetUserChats(userID)
	if err != nil {
		c.Error(errs.InternalServerError("unable to get user's chats"))
		return
	}

	c.JSON(http.StatusOK, chats)
}

// получение списка последних сообщений в чате с другим пользователем
func (s *Server) userGetChatMessages(c *gin.Context) {
	// получение id пользователя
	tokenID, exists := c.Get("id")
	if !exists {
		c.Error(errs.Unauthorized("unable to get user ID"))
		return
	}
	userID := tokenID.(uuid.UUID)

	// получение id собеседника
	chatReq := &dto.ChatRequest{}
	if err := c.BindJSON(&chatReq); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}
	userCompanion, err := s.service.UserGetPhone(chatReq.CompanionPhone)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	// получение списка последних сообщений
	messages, err := s.service.GetUsersChatMessages(userID, userCompanion.ID)
	if err != nil {
		c.Error(errs.InternalServerError("unable to get users chat messages"))
		return
	}

	c.JSON(http.StatusOK, messages)
}
