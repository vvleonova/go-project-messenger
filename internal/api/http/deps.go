package http

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"go-messenger/internal/models"
	"sync"

	"github.com/google/uuid"
)

//go:generate moq -skip-ensure -out mock.go . service

// сервис
type service interface {
	HealthCheck() error
	UserCreate(*dto.UserSignUp) (*models.User, error)
	UserCheck(*dto.UserSignIn) (*models.User, *errs.HTTPError)
	UserGenerateToken(uuid.UUID) (string, *errs.HTTPError)
	RefreshTokenUpdate(string) (string, *errs.HTTPError)
	RefreshTokenRevoke(uuid.UUID) *errs.HTTPError
	UserGetPhone(string) (*models.User, error)
	UserDelete(string) error
	UserUpdate(string, *dto.UserUpdate) error
	HandleMessage(*dto.MessageRequest, *UserConn, *sync.Map)
	HandleMessageUnread(*UserConn)
	GetUserChats(uuid.UUID) ([]models.ChatPreview, error)
	GetUsersChatMessages(uuid.UUID, uuid.UUID) ([]dto.MessageResponse, error)
}
