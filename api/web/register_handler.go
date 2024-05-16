package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/service/user"
)

const (
	passwordMinLength = 6
)

type RegisterHandler struct {
	userService *user.UserService
	logger      *slog.Logger
}

func NewRegisterHandler(userService *user.UserService, logger *slog.Logger) *RegisterHandler {
	return &RegisterHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *RegisterHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.serveRegisterPageHTML).Methods(http.MethodGet)
	router.HandleFunc("/register", h.submitRegisterHTML).Methods(http.MethodPost)
}

func (h *RegisterHandler) serveRegisterPageHTML(w http.ResponseWriter, r *http.Request) {
	content := templates.RegisterPage()
	err := templates.Layout(content, "Sportspazz").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func (h *RegisterHandler) submitRegisterHTML(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	h.logger.Info(fmt.Sprintf("govalidator.IsEmail(email): %v", govalidator.IsEmail(email)))
	if !govalidator.IsEmail(email) {
		templates.RegisterError("Invalid email address").Render(r.Context(), w)
		return
	}

	if len(password) < passwordMinLength {
		templates.RegisterError(fmt.Sprintf("Password has to be at least %d characters", passwordMinLength)).Render(r.Context(), w)
		return
	}

	if _, err := h.userService.RegisterUser(email, password); err != nil {
		templates.RegisterError(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}
