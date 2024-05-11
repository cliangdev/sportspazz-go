package rest_api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
	"github.com/sportspazz/service/user"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func getLogger(r *http.Request) *slog.Logger {
	return r.Context().Value("logger").(*slog.Logger)
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.registerUser).Methods(http.MethodPost)
}

func (h *UserHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var userRequest RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		getLogger(r).Error(err.Error())
		InvalidJsonResponse(w)
		return
	}

	if err := validRegisterUserRequest(userRequest); err != nil {
		ErrorJsonResponse(w, err.Error())
		return
	}

	defer r.Body.Close()

	newUser, err := h.userService.RegisterUser(userRequest.Email, userRequest.Password)
	if err != nil {
		ErrorJsonResponse(w, err.Error())
		return
	}

	userResponse := UserResponse{
		ID:        newUser.ID,
		CreatedOn: newUser.CreatedOn,
		UpdatedOn: newUser.UpdatedOn,
		Email:     newUser.Email,
	}

	JsonResponse(userResponse, w)
}

func validRegisterUserRequest(payload RegisterUserRequest) error {
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		return errors.New("invalid email address")
	}
	return nil
}
