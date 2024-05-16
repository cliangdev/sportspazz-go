package web

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/service/user"
)

type LoginHandler struct {
	userService *user.UserService
	logger      *slog.Logger
}

func NewLoginHandler(userService *user.UserService, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *LoginHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.serveLoginPageHTML).Methods(http.MethodGet)
}

func (h *LoginHandler) serveLoginPageHTML(w http.ResponseWriter, r *http.Request) {
	content := templates.LoginPage()
	err := templates.Layout(content, "Sportspazz").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
