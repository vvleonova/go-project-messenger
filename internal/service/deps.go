package service

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"

	"github.com/google/uuid"
)

//go:generate moq -skip-ensure -out mock.go . storage

// хранилище данных
type storage interface {
	Ping() error
	UserInsert(*models.User) error
	UserGetPhone(string) (*models.User, error)
	UserGetID(uuid.UUID) (*models.User, error)
	UserDelete(string) error
	UserUpdate(string, *dto.UserUpdate) error
	RefreshTokenSave(*models.RefreshToken) error
	RefreshTokenGet(string) (*models.RefreshToken, error)
	RefreshTokenRevoke(string) error
}
