package service

import (
	"fmt"
	"go-messenger/internal/api/errs"

	"github.com/google/uuid"
)

// получение нового JWT-token
func (s *Service) RefreshTokenUpdate(refreshToken string) (string, *errs.HTTPError) {
	// проверка наличия активного refresh JWT-token в БД
	rt, err := s.storage.RefreshTokenGet(refreshToken)
	if err != nil {
		return "", errs.NotFound(fmt.Errorf("unable to find refresh token: %w", err).Error())
	}

	// получение данных о пользователе по id
	user, err := s.UserGetID(rt.UserID)
	if err != nil {
		return "", errs.NotFound(err.Error())
	}

	// генерация нового JWT-token
	accessToken, httpError := s.UserGenerateToken(user.ID)
	if httpError != nil {
		return "", httpError
	}

	// revoke старого refresh JWT-token
	err = s.storage.RefreshTokenRevoke(refreshToken)
	if err != nil {
		return "", errs.InternalServerError(fmt.Errorf("unable to revoke refresh token %s: %w", refreshToken, err).Error())
	}

	return accessToken, nil
}

// revoke нового JWT-token
func (s *Service) RefreshTokenRevoke(id uuid.UUID) *errs.HTTPError {
	// revoke старого refresh JWT-token
	err := s.storage.RefreshTokenRevokeByID(id)
	if err != nil {
		return errs.InternalServerError(fmt.Errorf("unable to revoke refresh token for user %s: %w", id, err).Error())
	}

	return nil
}
