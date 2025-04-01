package service

import (
	"fmt"
	"time"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// обновление и сохранение структуры пользователя
func (s *Service) UserCreate(ur *dto.UserSignUp) (*models.User, error) {
	var err error

	// структура пользователя
	u := &models.User{}

	// создание ID
	u.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to create user uuid: %s", err)
	}

	// перенос полей из запроса
	u.Login = ur.Login
	u.Phone = ur.Phone
	u.BirthDate = ur.BirthDate.Time

	// определение времени создания пользователя
	u.CreatedOn = time.Now()

	// хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(ur.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash user password: %s", err)
	}
	u.PasswordHash = string(hash)

	// сохранение в базу данных
	if err := s.storage.UserInsert(u); err != nil {
		return nil, fmt.Errorf("unable to save user to database: %s", err)
	}

	return u, nil
}

// получение данных о пользователе из базы данных по номеру телефона
func (s *Service) UserGet(phone string) (*models.User, error) {
	u, err := s.storage.UserGet(phone)
	if err != nil {
		return nil, fmt.Errorf("unable to get user from database by phone %s: %s", phone, err)
	}

	return u, nil
}

// удаление данных о пользователе из базы данных по номеру телефона
func (s *Service) UserDelete(phone string) error {
	err := s.storage.UserDelete(phone)
	if err != nil {
		return fmt.Errorf("unable to delete user from database by phone %s: %s", phone, err)
	}

	return nil
}

// обновление структуры пользователя
func (s *Service) UserUpdate(phone string, ur *dto.UserUpdate) error {
	err := s.storage.UserUpdate(phone, ur)
	if err != nil {
		return fmt.Errorf("unable to update user with phone %s: %s", phone, err)
	}

	return nil
}
