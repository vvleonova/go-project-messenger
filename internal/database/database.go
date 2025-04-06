package database

import (
	"fmt"

	"go-messenger/internal/config"

	"github.com/jmoiron/sqlx"
)

// структура хранилища
type PgStorage struct {
	*sqlx.DB
}

// создание структуры хранилища
func New(config *config.DBConfig) (*PgStorage, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)

	return &PgStorage{db}, err
}
