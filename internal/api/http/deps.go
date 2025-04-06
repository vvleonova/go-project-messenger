package http

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"go-messenger/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

//go:generate moq -skip-ensure -out mock.go . service

// сервис
type service interface {
	HealthCheck() error
	UserCreate(*dto.UserSignUp) (*models.User, error)
	UserCheck(*dto.UserSignIn) (*models.User, *errs.HTTPError)
	UserGenerateToken(string, uuid.UUID) (string, *errs.HTTPError)
	RefreshTokenUpdate(string) (string, *errs.HTTPError)
	UserGetPhone(string) (*models.User, error)
	UserDelete(string) error
	UserUpdate(string, *dto.UserUpdate) error
	// TokenGenerate(phone string) (string, error)
	TokenVerify(string) (jwt.MapClaims, error)
}
