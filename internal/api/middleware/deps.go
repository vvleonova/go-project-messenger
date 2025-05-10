package middleware

import (
	"go-messenger/internal/models"
)

// хранилище
type storage interface {
	UserGetPhone(phone string) (*models.User, error)
}
