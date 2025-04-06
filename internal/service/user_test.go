package service

import (
	"errors"
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	login        = "testuser"
	correctPhone = "79999999999"
	wrongPhone   = "00000000000"
)

func TestService_UserCreate(t *testing.T) {
	// входные данные для теста
	birth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	signUp := &dto.UserSignUp{
		Login:     login,
		Phone:     correctPhone,
		BirthDate: dto.Date{Time: birth},
		Password:  "securepassword",
	}

	// возможные кейсы
	tests := []struct {
		name       string
		input      *dto.UserSignUp
		wantErr    bool
		mockInsert func(user *models.User) error
	}{
		{
			name:    "user creation successful",
			input:   signUp,
			wantErr: false,
			mockInsert: func(user *models.User) error {
				// проверка, что в UserInsert доходят правильные данные
				assert.Equal(t, login, user.Login)
				assert.Equal(t, correctPhone, user.Phone)
				assert.WithinDuration(t, birth, user.BirthDate, time.Second)
				assert.NotEmpty(t, user.PasswordHash)
				return nil
			},
		},
		{
			name:    "user creation fails",
			input:   signUp,
			wantErr: true,
			mockInsert: func(user *models.User) error {
				return errors.New("unable to save user to database")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание мока с нужным поведением
			mockStorage := &storageMock{
				UserInsertFunc: tt.mockInsert,
			}

			// создание сервиса с моком
			s := &Service{
				storage: mockStorage,
			}

			// вызов тестируемой функции
			got, err := s.UserCreate(tt.input)

			// проверка на ошибку
			if tt.wantErr {
				assert.Error(t, err, "expected an error, got nil") // функция должна вернуть ошибку
				assert.Nil(t, got)                                 // функция не должна вернуть объект
			} else {
				assert.NoError(t, err) // функция должна пройти без ошибки
				assert.NotNil(t, got)  // функция должна вернуть не nil объект
				// проверки, что вернулся корректный пользователь
				assert.Equal(t, login, got.Login)
				assert.Equal(t, correctPhone, got.Phone)
				assert.WithinDuration(t, birth, got.BirthDate, time.Second)
				assert.NotEmpty(t, got.PasswordHash)
				assert.NotEqual(t, uuid.Nil, got.ID)
			}
		})
	}
}

func TestService_UserGetPhone(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name       string
		phone      string
		mockReturn *models.User
		mockError  error
		wantErr    bool
		wantLogin  string
		wantPhone  string
	}{
		{
			name:  "user returned",
			phone: correctPhone,
			mockReturn: &models.User{
				Login: login,
				Phone: correctPhone,
			},
			mockError: nil,
			wantErr:   false,
			wantLogin: login,
			wantPhone: correctPhone,
		},
		{
			name:       "user not returned",
			phone:      wrongPhone,
			mockReturn: nil,
			mockError:  errors.New("not found"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание мока с нужным поведением
			mockStorage := &storageMock{
				UserGetPhoneFunc: func(phone string) (*models.User, error) {
					assert.Equal(t, tt.phone, phone)
					return tt.mockReturn, tt.mockError
				},
			}

			// создание сервиса с моком
			s := &Service{
				storage: mockStorage,
			}

			// вызов тестируемой функции
			got, err := s.UserGetPhone(tt.phone)

			// проверка ошибок
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantLogin, got.Login)
				assert.Equal(t, tt.wantPhone, got.Phone)
			}
		})
	}
}

func TestService_UserDelete(t *testing.T) {
	// возможные кейсы
	tests := []struct {
		name      string
		phone     string
		mockError error
		wantErr   bool
	}{
		{
			name:      "user deleted",
			phone:     correctPhone,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "user not deleted",
			phone:     wrongPhone,
			mockError: errors.New("unable to delete user from database"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание мока с нужным поведением
			mockStorage := &storageMock{
				UserDeleteFunc: func(phone string) error {
					assert.Equal(t, tt.phone, phone)
					return tt.mockError
				},
			}

			// создание сервиса с моком
			s := &Service{
				storage: mockStorage,
			}

			// вызов тестируемой функции
			err := s.UserDelete(tt.phone)

			// проверка ошибок
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
