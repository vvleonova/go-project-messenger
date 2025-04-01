package config

import (
	"os"
	"strconv"

	"go-messenger/internal/logger"

	"github.com/joho/godotenv"
)

// структура конфига
type Config struct {
	Env       string
	Host      string
	Port      int
	SecretKey string
	DB        DBConfig
}

// структура базы данных
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// загрузка .env
func LoadConfig(logger logger.Logger) (*Config, error) {
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Error("error loading .env file", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Error("error getting PORT value", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Error("error getting DB_PORT value", err)
	}

	return &Config{
		Env:       os.Getenv("ENV"),
		Host:      os.Getenv("HOST"),
		Port:      port,
		SecretKey: os.Getenv("SECRET_KEY"),
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}, nil
}
