package database

import (
	"fmt"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/config"
	"go-messenger/internal/models"

	"github.com/jmoiron/sqlx"
)

// структура хранилища
type PgStorage struct {
	*sqlx.DB
}

// создание структуры хранилища
func New(config *config.Config) (*PgStorage, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)

	return &PgStorage{db}, err
}

// сохранение пользователя
func (pq *PgStorage) UserInsert(u *models.User) error {
	query := `
	INSERT INTO users (id, login, phone, birth_date, password, created_on)
	VALUES (:id, :login, :phone, :birth_date, :password, :created_on)
	`
	_, err := pq.NamedExec(query, u)

	return err
}

// получение данных о пользователе из базы данных по номеру телефона
func (pq *PgStorage) UserGet(phone string) (*models.User, error) {
	query := "SELECT * FROM users WHERE phone = $1"

	u := &models.User{}
	if err := pq.Get(u, query, phone); err != nil {
		return nil, err
	}

	return u, nil
}

// удаление данных о пользователе из базы данных по номеру телефона
func (pq *PgStorage) UserDelete(phone string) error {
	query := "DELETE FROM users WHERE phone = $1"

	_, err := pq.Exec(query, phone)

	return err
}

// обновление данных о пользователе в базе данных по номеру телефона
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
