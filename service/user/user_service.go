package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"firebase.google.com/go/v4/auth"
)

type UserService struct {
	store          *UserStore
	firebaseClient *auth.Client
	logger         *slog.Logger
}

func NewUserService(store *UserStore, firebaseClient *auth.Client, logger *slog.Logger) *UserService {
	return &UserService{
		store:          store,
		firebaseClient: firebaseClient,
		logger:         logger,
	}
}

func (u *UserService) RegisterUser(email, password string) (*User, error) {
	if user := u.store.GetUserByEmail(email); user != nil {
		return nil, fmt.Errorf("user %s is already registered", email)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(password).
		Disabled(false)
	newUser, err := u.firebaseClient.CreateUser(context.Background(), params)

	if err != nil {
		u.logger.Error("Unable to create user in Firebase", slog.Any("err", err))
		return nil, errors.New("unable to register due to internal error")
	}
	u.logger.Info("New user created in Firebase", slog.Any("user", newUser))

	return u.store.CreateUser(email), nil
}
