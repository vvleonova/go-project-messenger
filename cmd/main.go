package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-messenger/internal/api/http"
	"go-messenger/internal/config"
	"go-messenger/internal/database"
	"go-messenger/internal/logger"
	"go-messenger/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	// контекст
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// инициализация логгера
	logger, err := logger.NewLogger()
	if err != nil {
		panic("error creating logger: " + err.Error())
	}

	// конфиг
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("error loading config", err)
	}

	// подключение к БД
	pgStorage, err := database.New(&config.DB)
	if err != nil {
		logger.Fatal("error connecting to database", err)
	}

	// сервис
	service := service.New(config, pgStorage, logger)

	// запуск сервера
	server := http.New(config, service, logger)
	server.Start(ctx)
}
