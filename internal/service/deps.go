package service

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"
)

//go:generate moq -skip-ensure -out mock.go . storage

// хранилище данных
type storage interface {
	UserInsert(*models.User) error
	UserGet(phone string) (*models.User, error)
	UserDelete(phone string) error
	UserUpdate(phone string, u *dto.UserUpdate) error
}
