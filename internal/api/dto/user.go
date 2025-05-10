package dto

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// обработка ошибок
var (
	errShortLogin    = errors.New("user login is not long enough")
	errWrongPhone    = errors.New("wrong phone number format")
	errWrongDate     = errors.New("birth date must be in format 'YYYY-MM-DD'")
	errYoungAge      = errors.New("user must be older than 18")
	errShortPassword = errors.New("password must be at least 8 symbols")
	errWrongPassword = errors.New("password can contain only latin letters, numbers, ! and _")
	errLowPassword   = errors.New("password must contain at least one capital letter")
)

// структура пользователя для регистрации
type UserSignUp struct {
	Login     string `json:"login"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	Password  string `json:"password"`
	// Photo
}

// структура пользователя для авторизации
type UserSignIn struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// структура пользователя для обновления
type UserUpdate struct {
	Login     *string `json:"login,omitempty" db:"login"`
	BirthDate *string `json:"birth_date,omitempty" db:"birth_date"`
	// Photo
}

// структура пользователя для вывода результата
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	Phone     string    `json:"phone"`
	BirthDate string    `json:"birth_date"`
	CreatedOn string    `json:"created_on"`
}

// валидация структуры пользователя для регистрации
func (u UserSignUp) ValidateUserSignUp() error {
	// проверка длины логина
	if len(u.Login) < 5 {
		return errShortLogin
	}

	// проверка номера телефона
	if !regexp.MustCompile(`^\d{11}$`).MatchString(u.Phone) {
		return errWrongPhone
	}

	// проверка возраста
	birthDate, err := time.Parse("2006-01-02", u.BirthDate)
	if err != nil {
		return errWrongDate
	}

	if !birthDate.Before(time.Now().AddDate(-18, 0, 0)) {
		return errYoungAge
	}

	// проверка длины пароля
	if len(u.Password) < 8 {
		return errShortPassword
	}

	// проверка символов пароля
	if !regexp.MustCompile(`^[a-zA-Z0-9!_]+$`).MatchString(u.Password) {
		return errWrongPassword
	}

	// проверка наличия заглавных букв в пароле
	hasUpper := false
	for _, r := range u.Password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errLowPassword
	}

	return nil
}

// валидация структуры пользователя для авторизации
func (u UserSignIn) ValidateUserSignIn() error {
	// проверка номера телефона
	if !regexp.MustCompile(`^\d{11}$`).MatchString(u.Phone) {
		return errWrongPhone
	}

	return nil
}

// валидация структуры пользователя для обновления
func (u UserUpdate) ValidateUserUpdate() error {
	// проверка длины логина
	if u.Login != nil && len(*u.Login) < 5 {
		return errShortLogin
	}

	// проверка возраста
	birthDate, err := time.Parse("2006-01-02", *u.BirthDate)
	if err != nil {
		return errWrongDate
	}

	if u.BirthDate != nil && !birthDate.Before(time.Now().AddDate(-18, 0, 0)) {
		return errYoungAge
	}

	return nil
}
