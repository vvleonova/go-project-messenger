package service

import (
	"errors"
	"go-messenger/internal/config"
	"go-messenger/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	refreshToken    = "refresh token"
	refreshTokenNew = "new refresh token"

	errTokenRevoke = "unable to revoke refresh token"
	errTokenSave   = "unable to save refresh token to database"
	errTokenGet    = "unable to find refresh token"
)

func TestService_RefreshTokenUpdate(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name          string
		storage       *storageMock
		input         string
		wantErr       bool
		wantErrText   string
		wantErrStatus int
	}{
		{
			name: "token updated",
			storage: &storageMock{
				RefreshTokenGetFunc: func(rt string) (*models.RefreshToken, error) {
					assert.Equal(t, refreshToken, rt)
					return &models.RefreshToken{
						Token:  rt,
						UserID: correctID,
					}, nil
				},
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					assert.Equal(t, correctID, id)
					return &models.User{
						ID: id,
					}, nil
				},
				RefreshTokenSaveFunc: func(rt *models.RefreshToken) error {
					assert.NotEmpty(t, rt.ID)
					assert.NotEmpty(t, rt.Token)
					assert.Equal(t, correctID, rt.UserID)
					assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), rt.ExpiresAt, 2*time.Second)
					assert.False(t, rt.Revoked)
					return nil
				},
				RefreshTokenRevokeFunc: func(refreshToken string) error {
					return nil
				},
			},
			input: refreshToken,
		},
		{
			name: "old token not revoked",
			storage: &storageMock{
				RefreshTokenGetFunc: func(rt string) (*models.RefreshToken, error) {
					assert.Equal(t, refreshToken, rt)
					return &models.RefreshToken{
						Token:  rt,
						UserID: correctID,
					}, nil
				},
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					assert.Equal(t, correctID, id)
					return &models.User{
						ID: id,
					}, nil
				},
				RefreshTokenSaveFunc: func(rt *models.RefreshToken) error {
					assert.NotEmpty(t, rt.ID)
					assert.NotEmpty(t, rt.Token)
					assert.Equal(t, correctID, rt.UserID)
					assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), rt.ExpiresAt, 2*time.Second)
					assert.False(t, rt.Revoked)
					return nil
				},
				RefreshTokenRevokeFunc: func(refreshToken string) error {
					return errors.New(errTokenRevoke)
				},
			},
			input:         refreshToken,
			wantErr:       true,
			wantErrText:   errTokenRevoke,
			wantErrStatus: 500,
		},
		{
			name: "new token not created",
			storage: &storageMock{
				RefreshTokenGetFunc: func(rt string) (*models.RefreshToken, error) {
					assert.Equal(t, refreshToken, rt)
					return &models.RefreshToken{
						Token:  rt,
						UserID: correctID,
					}, nil
				},
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					assert.Equal(t, correctID, id)
					return &models.User{
						ID: id,
					}, nil
				},
				RefreshTokenSaveFunc: func(rt *models.RefreshToken) error {
					return errors.New(errSaveDB)
				},
				RefreshTokenRevokeFunc: nil,
			},
			input:         refreshToken,
			wantErr:       true,
			wantErrText:   errTokenSave,
			wantErrStatus: 500,
		},
		{
			name: "new token not created",
			storage: &storageMock{
				RefreshTokenGetFunc: func(rt string) (*models.RefreshToken, error) {
					assert.Equal(t, refreshToken, rt)
					return &models.RefreshToken{
						Token:  rt,
						UserID: correctID,
					}, nil
				},
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					assert.Equal(t, id, id)
					return nil, errors.New(errGetDB)
				},
				RefreshTokenSaveFunc:   nil,
				RefreshTokenRevokeFunc: nil,
			},
			input:         refreshToken,
			wantErr:       true,
			wantErrText:   errGetDB,
			wantErrStatus: 404,
		},
		{
			name: "token not found",
			storage: &storageMock{
				RefreshTokenGetFunc: func(rt string) (*models.RefreshToken, error) {
					assert.Equal(t, refreshToken, rt)
					return nil, errors.New(errTokenGet)
				},
				UserGetIDFunc:          nil,
				RefreshTokenSaveFunc:   nil,
				RefreshTokenRevokeFunc: nil,
			},
			input:         refreshToken,
			wantErr:       true,
			wantErrText:   errTokenGet,
			wantErrStatus: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
				config:  &config.Config{SecretKey: "secretkey"},
			}

			// вызов тестируемой функции
			refreshTokenUpd, err := s.RefreshTokenUpdate(tt.input)

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.wantErrText)
				assert.Equal(t, tt.wantErrStatus, err.Status)
				assert.Empty(t, refreshTokenUpd)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, refreshTokenUpd)
			}
		})
	}
}

func TestService_RefreshTokenRevoke(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name          string
		storage       *storageMock
		input         uuid.UUID
		wantErr       bool
		wantErrText   string
		wantErrStatus int
	}{
		{
			name: "token revoked",
			storage: &storageMock{
				RefreshTokenRevokeByIDFunc: func(id uuid.UUID) error {
					return nil
				},
			},
			input: correctID,
		},
		{
			name: "token not revoked",
			storage: &storageMock{
				RefreshTokenRevokeByIDFunc: func(id uuid.UUID) error {
					return errors.New(errTokenRevoke)
				},
			},
			input:         correctID,
			wantErr:       true,
			wantErrText:   errTokenRevoke,
			wantErrStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			err := s.RefreshTokenRevoke(tt.input)

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.wantErrText)
				assert.Equal(t, tt.wantErrStatus, err.Status)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
