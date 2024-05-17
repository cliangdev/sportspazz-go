package user

import (
	"time"
)

type User struct {
	internalId uint `gorm:"primaryKey"`
	ID         string
	CreatedOn  time.Time `gorm:"type:datetime(3)"`
	UpdatedOn  time.Time `gorm:"type:datetime(3)"`
	Email      string
}

func NewUser(id, email string) *User {
	return &User{
		ID:        id,
		Email:     email,
		CreatedOn: time.Now().UTC(),
		UpdatedOn: time.Now().UTC(),
	}
}
