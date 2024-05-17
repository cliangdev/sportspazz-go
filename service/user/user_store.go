package user

import (
	"log/slog"

	"gorm.io/gorm"
)

type UserStore struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserStore(db *gorm.DB, logger *slog.Logger) *UserStore {
	return &UserStore{
		db:     db,
		logger: logger,
	}
}

func (s *UserStore) GetUserByEmail(email string) *User {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func (s *UserStore) CreateUser(id, email string) *User {
	user := NewUser(id, email)

	if err := s.db.Create(user).Error; err != nil {
		s.logger.Error("not able to create a new user", slog.Any("err", err))
	}

	return user
}
