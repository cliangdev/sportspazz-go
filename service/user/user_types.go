package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	internalId uint `gorm:"primaryKey"`
	ID         string
	CreatedOn  time.Time `gorm:"type:datetime(3)"`
	UpdatedOn  time.Time `gorm:"type:datetime(3)"`
	Email      string
}

func NewUser(email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		CreatedOn: time.Now().UTC(),
		UpdatedOn: time.Now().UTC(),
	}
}
