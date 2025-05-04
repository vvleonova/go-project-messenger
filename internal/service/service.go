package service

import (
	"go-messenger/internal/config"
	"go-messenger/internal/logger"

	"golang.org/x/crypto/bcrypt"
)

// структура сервиса
type Service struct {
	storage
	logger       logger.Logger
	config       *config.Config
	hashComparer func(hashedPassword, password []byte) error
}

// создание структуры сервиса
func New(config *config.Config, db storage, logger logger.Logger) *Service {
	return &Service{
		storage:      db,
		logger:       logger,
		config:       config,
		hashComparer: bcrypt.CompareHashAndPassword,
	}
}
