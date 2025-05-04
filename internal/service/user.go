package service

import (
	"fmt"
	"time"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"go-messenger/internal/jwt"
	"go-messenger/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// обновление и сохранение структуры пользователя
func (s *Service) UserCreate(ur *dto.UserSignUp) (*models.User, error) {
	// структура пользователя
	u := &models.User{}

	// создание ID
	var err error
	u.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to create user uuid: %w", err)
	}

	// перенос полей из запроса
	u.Login = ur.Login
	u.Phone = ur.Phone
	u.BirthDate, err = time.Parse("2006-01-02", ur.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("unable to parse birth date: %w", err)
	}

	// определение времени создания пользователя
	u.CreatedOn = time.Now()

	// хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(ur.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash user password: %w", err)
	}
	u.PasswordHash = string(hash)

	// сохранение в базу данных
	if err := s.storage.UserInsert(u); err != nil {
		return nil, fmt.Errorf("unable to save user to database: %w", err)
	}

	return u, nil
}

// получение данных о пользователе из базы данных по номеру телефона
func (s *Service) UserGetPhone(phone string) (*models.User, error) {
	u, err := s.storage.UserGetPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("unable to get user from database by phone %s: %w", phone, err)
	}

	return u, nil
}

// получение данных о пользователе из базы данных по id
func (s *Service) UserGetID(id uuid.UUID) (*models.User, error) {
	u, err := s.storage.UserGetID(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get user from database by id %s: %w", id, err)
	}

	return u, nil
}

// проверка пользователя
func (s *Service) UserCheck(body *dto.UserSignIn) (*models.User, *errs.HTTPError) {
	// получение данных о пользователе из БД
	user, err := s.UserGetPhone(body.Phone)
	if err != nil {
		return nil, errs.InternalServerError(err.Error())
	}

	// проверка пароля
	err = s.hashComparer([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		return nil, errs.Unauthorized("invalid credentials")
	}

	return user, nil
}

// удаление данных о пользователе из базы данных по номеру телефона
func (s *Service) UserDelete(phone string) error {
	err := s.storage.UserDelete(phone)
	if err != nil {
		return fmt.Errorf("unable to delete user from database by phone %s: %w", phone, err)
	}

	return nil
}

// обновление структуры пользователя
func (s *Service) UserUpdate(phone string, ur *dto.UserUpdate) error {
	err := s.storage.UserUpdate(phone, ur)
	if err != nil {
		return fmt.Errorf("unable to update user with phone %s: %w", phone, err)
	}

	return nil
}

// генерация JWT-токена
func (s *Service) UserGenerateToken(id uuid.UUID) (string, *errs.HTTPError) {
	// генерация access JWT-token
	accessToken, _, err := jwt.TokenGenerate(id, s.config.SecretKey, 1*time.Minute)
	if err != nil {
		return "", errs.InternalServerError(fmt.Errorf("unable to create access token: %w", err).Error())
	}

	// генерация refresh JWT-token
	refreshToken, expiredAt, err := jwt.TokenGenerate(id, s.config.SecretKey, 30*24*time.Hour)
	if err != nil {
		return "", errs.InternalServerError(fmt.Errorf("unable to create refresh token: %w", err).Error())
	}

	// сохранение refresh JWT-token в БД
	rt, err := s.refreshTokenPrepare(refreshToken, id, expiredAt)
	if err != nil {
		return "", errs.InternalServerError(fmt.Errorf("unable to save refresh token to database: %w", err).Error())
	}

	err = s.storage.RefreshTokenSave(rt)
	if err != nil {
		return "", errs.InternalServerError(fmt.Errorf("unable to save refresh token to database: %w", err).Error())
	}

	// возвращение access token клиенту
	return accessToken, nil
}

// подготовка refresh JWT-token
func (s *Service) refreshTokenPrepare(refreshToken string, userID uuid.UUID, expiredAt time.Time) (*models.RefreshToken, error) {
	// структура refresh JWT-token
	rt := &models.RefreshToken{}

	// создание ID
	var err error
	rt.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to create refresh token uuid: %w", err)
	}

	// сохранение токена
	rt.Token = refreshToken

	// сохранение id пользователя
	rt.UserID = userID

	// определение времени создания и expiration
	rt.CreatedOn = time.Now()
	rt.ExpiresAt = expiredAt

	// флаг отзыва
	rt.Revoked = false

	return rt, nil
}
