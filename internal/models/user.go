package models

import (
	"fmt"

	"time"

	"github.com/google/uuid"
)

// структура пользователя
type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	Phone        string    `db:"phone"`
	BirthDate    time.Time `db:"birth_date"`
	PasswordHash string    `db:"password" json:"-"`
	CreatedOn    time.Time `db:"created_on"`
	// LastOnline time.Time `db:"last_online"`
	// Photo
}

// логирование
func (u *User) String() string {
	return fmt.Sprintf("user %s with phone number %s", u.Login, u.Phone)
}
