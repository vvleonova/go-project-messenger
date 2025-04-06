package dto

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// обработка ошибок
var errShortLogin = errors.New("user login is not long enough")
var errWrongPhone = errors.New("wrong phone number format")
var errYoungAge = errors.New("user must be older than 18")
var errFutureBirth = errors.New("birth date can't be in future")
var errShortPassword = errors.New("password must be at least 8 symbols")
var errWrongPassword = errors.New("password can contain only latin letters, numbers, ! and _")
var errLowPassword = errors.New("password must contain at least one capital letter")

// структура для даты
type Date struct {
	time.Time
}

// структура пользователя для регистрации
type UserSignUp struct {
	Login     string `json:"login"`
	Phone     string `json:"phone"`
	BirthDate Date   `json:"birth_date"`
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
	BirthDate *Date   `json:"birth_date,omitempty" db:"birth_date"`
	// Photo
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
	if !u.BirthDate.Before(time.Now().AddDate(-18, 0, 0)) {
		return errYoungAge
	}

	// проверка даты рождения
	if u.BirthDate.After(time.Now()) {
		return errFutureBirth
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
	if u.BirthDate != nil && !u.BirthDate.Before(time.Now().AddDate(-18, 0, 0)) {
		return errYoungAge
	}

	// проверка даты рождения
	if u.BirthDate != nil && u.BirthDate.After(time.Now()) {
		return errFutureBirth
	}

	return nil
}

// кастомный парсер JSON для Date
func (ct *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime

	return nil
}
