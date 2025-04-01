package server

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// хранилище данных
type storage interface {
	Ping() error
}

// сервис
type service interface {
	UserCreate(*dto.UserSignUp) (*models.User, error)
	UserGet(phone string) (*models.User, error)
	UserDelete(phone string) error
	UserUpdate(string, *dto.UserUpdate) error
	TokenGenerate(phone string) (string, error)
	TokenVerify(accessToken string) (jwt.MapClaims, error)
}
