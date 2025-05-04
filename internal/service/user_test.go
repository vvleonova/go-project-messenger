package service

import (
	"errors"
	"go-messenger/internal/api/dto"
	"go-messenger/internal/config"
	"go-messenger/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	login           = "testuser"
	correctPhone    = "79999999999"
	wrongPhone      = "00000000000"
	correctPassword = "mypassword"
	hashPassword    = "mypasswordhashed"
	wrongPassword   = "notmypassword"

	errSaveDB       = "unable to save to database"
	errGetDB        = "unable to get data from database"
	errDeleteDB     = "unable to delete data from database"
	errUpdateDB     = "unable to update data in database"
	errInvalidCreds = "invalid credentials"
)

var (
	loginVar     = "testuser"
	birth        = "1997-01-01"
	birthDate, _ = time.Parse("2006-01-02", birth)
	correctID    = uuid.New()
	wrongID      = uuid.New()
)

func TestService_UserCreate(t *testing.T) {
	// входные данные для теста
	signUp := &dto.UserSignUp{
		Login:     login,
		Phone:     correctPhone,
		BirthDate: birth,
		Password:  correctPassword,
	}

	// возможные кейсы
	tests := []struct {
		name    string
		storage *storageMock
		input   *dto.UserSignUp
		wantErr bool
	}{
		{
			name: "user creation successful",
			storage: &storageMock{
				UserInsertFunc: func(user *models.User) error {
					// проверка, что в UserInsert доходят правильные данные
					assert.Equal(t, login, user.Login)
					assert.Equal(t, correctPhone, user.Phone)
					assert.WithinDuration(t, birthDate, user.BirthDate, time.Second)
					assert.NotEmpty(t, user.PasswordHash)
					return nil
				},
			},
			input:   signUp,
			wantErr: false,
		},
		{
			name: "user creation fails",
			storage: &storageMock{
				UserInsertFunc: func(user *models.User) error {
					return errors.New(errSaveDB)
				},
			},
			input:   signUp,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			got, err := s.UserCreate(tt.input)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errSaveDB)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotEqual(t, uuid.Nil, got.ID)
				assert.Equal(t, login, got.Login)
				assert.Equal(t, correctPhone, got.Phone)
				assert.WithinDuration(t, birthDate, got.BirthDate, time.Second)
				assert.NotEmpty(t, got.PasswordHash)
			}
		})
	}
}

func TestService_UserGetPhone(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name    string
		storage *storageMock
		input   string
		wantErr bool
	}{
		{
			name: "user returned",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, correctPhone, phone)
					return &models.User{
						Phone: correctPhone,
					}, nil
				},
			},
			input:   correctPhone,
			wantErr: false,
		},
		{
			name: "user not returned",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					return nil, errors.New(errGetDB)
				},
			},
			input:   wrongPhone,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			got, err := s.UserGetPhone(tt.input)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errGetDB)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, correctPhone, got.Phone)
			}
		})
	}
}

func TestService_UserGetID(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name    string
		storage *storageMock
		input   uuid.UUID
		wantErr bool
	}{
		{
			name: "user returned",
			storage: &storageMock{
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					assert.Equal(t, correctID, id)
					return &models.User{
						ID: correctID,
					}, nil
				},
			},
			input:   correctID,
			wantErr: false,
		},
		{
			name: "user not returned",
			storage: &storageMock{
				UserGetIDFunc: func(id uuid.UUID) (*models.User, error) {
					return nil, errors.New(errGetDB)
				},
			},
			input:   wrongID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			got, err := s.UserGetID(tt.input)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errGetDB)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, correctID, got.ID)
			}
		})
	}
}

func TestService_UserCheck(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name               string
		storage            *storageMock
		hashComparer       func(hashedPassword, password []byte) error
		input              *dto.UserSignIn
		wantErrNotFound    bool
		wantErrInvalidCred bool
		wantErrStatus      int
	}{
		{
			name: "user is correct",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, correctPhone, phone)
					return &models.User{
						ID:           correctID,
						Login:        login,
						Phone:        correctPhone,
						PasswordHash: hashPassword,
					}, nil
				},
			},
			hashComparer: func(hashedPassword, password []byte) error {
				assert.Equal(t, []byte(hashPassword), hashedPassword)
				assert.Equal(t, []byte(correctPassword), password)
				return nil
			},
			input: &dto.UserSignIn{
				Phone:    correctPhone,
				Password: correctPassword,
			},
		},
		{
			name: "user not returned",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					return nil, errors.New(errGetDB)
				},
			},
			hashComparer: func(hashedPassword, password []byte) error {
				return nil
			},
			input: &dto.UserSignIn{
				Phone:    wrongPhone,
				Password: correctPassword,
			},
			wantErrNotFound: true,
			wantErrStatus:   500,
		},
		{
			name: "incorrect password",
			storage: &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, correctPhone, phone)
					return &models.User{
						ID:           correctID,
						Login:        login,
						Phone:        correctPhone,
						PasswordHash: hashPassword,
					}, nil
				},
			},
			hashComparer: func(hashedPassword, password []byte) error {
				return errors.New(errInvalidCreds)
			},
			input: &dto.UserSignIn{
				Phone:    correctPhone,
				Password: wrongPassword,
			},
			wantErrInvalidCred: true,
			wantErrStatus:      401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage:      tt.storage,
				hashComparer: tt.hashComparer,
			}

			// вызов тестируемой функции
			got, err := s.UserCheck(tt.input)

			// проверка результата
			if tt.wantErrNotFound {
				assert.ErrorContains(t, err, errGetDB)
				assert.Nil(t, got)
				assert.Equal(t, tt.wantErrStatus, err.Status)
			} else if tt.wantErrInvalidCred {
				assert.ErrorContains(t, err, errInvalidCreds)
				assert.Nil(t, got)
				assert.Equal(t, tt.wantErrStatus, err.Status)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, correctID, got.ID)
				assert.Equal(t, correctPhone, got.Phone)
				assert.Equal(t, hashPassword, got.PasswordHash)
			}
		})
	}
}

func TestService_UserDelete(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name    string
		storage *storageMock
		input   string
		wantErr bool
	}{
		{
			name: "user deleted",
			storage: &storageMock{
				UserDeleteFunc: func(phone string) error {
					assert.Equal(t, correctPhone, phone)
					return nil
				},
			},
			input:   correctPhone,
			wantErr: false,
		},
		{
			name: "user not deleted",
			storage: &storageMock{
				UserDeleteFunc: func(phone string) error {
					return errors.New(errDeleteDB)
				},
			},
			input:   wrongPhone,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			err := s.UserDelete(tt.input)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errDeleteDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UserUpdate(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name       string
		storage    *storageMock
		inputPhone string
		inputData  *dto.UserUpdate
		wantErr    bool
	}{
		{
			name: "user updated all",
			storage: &storageMock{
				UserUpdateFunc: func(phone string, ur *dto.UserUpdate) error {
					assert.Equal(t, correctPhone, phone)
					assert.Equal(t, &loginVar, ur.Login)
					assert.Equal(t, &loginVar, ur.Login)
					assert.Equal(t, &birth, ur.BirthDate)
					return nil
				},
			},
			inputPhone: correctPhone,
			inputData: &dto.UserUpdate{
				Login:     &loginVar,
				BirthDate: &birth,
			},
		},
		{
			name: "user updated login",
			storage: &storageMock{
				UserUpdateFunc: func(phone string, ur *dto.UserUpdate) error {
					assert.Equal(t, correctPhone, phone)
					assert.Equal(t, &loginVar, ur.Login)
					return nil
				},
			},
			inputPhone: correctPhone,
			inputData: &dto.UserUpdate{
				Login: &loginVar,
			},
		},
		{
			name: "user updated birth date",
			storage: &storageMock{
				UserUpdateFunc: func(phone string, ur *dto.UserUpdate) error {
					assert.Equal(t, correctPhone, phone)
					assert.Equal(t, &birth, ur.BirthDate)
					return nil
				},
			},
			inputPhone: correctPhone,
			inputData: &dto.UserUpdate{
				BirthDate: &birth,
			},
		},
		{
			name: "user not updated",
			storage: &storageMock{
				UserUpdateFunc: func(phone string, ur *dto.UserUpdate) error {
					return errors.New(errUpdateDB)
				},
			},
			inputPhone: wrongPhone,
			inputData:  nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			s := &Service{
				storage: tt.storage,
			}

			// вызов тестируемой функции
			err := s.UserUpdate(tt.inputPhone, tt.inputData)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errUpdateDB)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}

func TestService_UserGenerateToken(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name          string
		storage       *storageMock
		input         uuid.UUID
		wantErr       bool
		wantErrStatus int
	}{
		{
			name: "token generated",
			storage: &storageMock{
				RefreshTokenSaveFunc: func(rt *models.RefreshToken) error {
					assert.NotEmpty(t, rt.ID)
					assert.NotEmpty(t, rt.Token)
					assert.Equal(t, correctID, rt.UserID)
					assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), rt.ExpiresAt, 2*time.Second)
					assert.False(t, rt.Revoked)
					return nil
				},
			},
			input: correctID,
		},
		{
			name: "token not saved",
			storage: &storageMock{
				RefreshTokenSaveFunc: func(rt *models.RefreshToken) error {
					return errors.New(errSaveDB)
				},
			},
			input:         correctID,
			wantErr:       true,
			wantErrStatus: 500,
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
			accessToken, err := s.UserGenerateToken(tt.input)

			// проверка результата
			if tt.wantErr {
				assert.ErrorContains(t, err, errSaveDB)
				assert.Equal(t, tt.wantErrStatus, err.Status)
				assert.Empty(t, accessToken)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, accessToken)
			}
		})
	}
}
