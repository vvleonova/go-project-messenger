package database

import (
	"go-messenger/internal/api/dto"
	"go-messenger/internal/models"

	"github.com/google/uuid"
)

// сохранение пользователя
func (pq *PgStorage) UserInsert(u *models.User) error {
	query := `
	INSERT INTO users (id, login, phone, birth_date, password, created_on)
	VALUES (:id, :login, :phone, :birth_date, :password, :created_on)
	`
	_, err := pq.NamedExec(query, u)

	return err
}

// получение данных о пользователе по номеру телефона
func (pq *PgStorage) UserGetPhone(phone string) (*models.User, error) {
	query := "SELECT * FROM users WHERE phone = $1"

	u := &models.User{}
	if err := pq.Get(u, query, phone); err != nil {
		return nil, err
	}

	return u, nil
}

// получение данных о пользователе по номеру телефона
func (pq *PgStorage) UserGetID(id uuid.UUID) (*models.User, error) {
	query := "SELECT * FROM users WHERE id = $1"

	u := &models.User{}
	if err := pq.Get(u, query, id); err != nil {
		return nil, err
	}

	return u, nil
}

// удаление данных о пользователе по номеру телефона
func (pq *PgStorage) UserDelete(phone string) error {
	query := "DELETE FROM users WHERE phone = $1"

	_, err := pq.Exec(query, phone)

	return err
}

// обновление данных о пользователе по номеру телефона
func (pq *PgStorage) UserUpdate(phone string, u *dto.UserUpdate) error {
	query := `
	UPDATE users SET 
		login = COALESCE(NULLIF(:login, ''), login),
		birth_date = COALESCE(:birth_date, birth_date)
	WHERE phone = :phone
	`

	params := map[string]any{
		"login":      u.Login,
		"birth_date": u.BirthDate.Time,
		"phone":      phone,
	}

	_, err := pq.NamedExec(query, params)

	return err
}
