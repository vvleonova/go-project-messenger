package service

import (
	"go-messenger/internal/config"
	"go-messenger/internal/logger"
)

// структура сервиса
type Service struct {
	storage
	logger logger.Logger
	config *config.Config
}

// создание структуры сервиса
func New(config *config.Config, db storage, logger logger.Logger) *Service {
	return &Service{
		db,
		logger,
		config,
	}
}
