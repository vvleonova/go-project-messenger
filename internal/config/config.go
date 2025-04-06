package config

import (
	"fmt"
	"os"
	"strconv"

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
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("error getting PORT value: %w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("error getting DB_PORT value: %w", err)
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
