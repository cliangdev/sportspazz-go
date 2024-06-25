package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/service/user"
	"github.com/sportspazz/utils"
)

type LoginHandler struct {
	userService        *user.UserService
	firebaseRestClient *client.FirebaseClient
	logger             *slog.Logger
}

func NewLoginHandler(userService *user.UserService, firebaseRestClient *client.FirebaseClient, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		userService:        userService,
		firebaseRestClient: firebaseRestClient,
		logger:             logger,
	}
}

func (h *LoginHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.serveLoginPageHTML).Methods(http.MethodGet)
	router.HandleFunc("/login", h.loginHTML).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.logoutHTML).Methods(http.MethodPost)
}

func (h *LoginHandler) serveLoginPageHTML(w http.ResponseWriter, r *http.Request) {
	content := templates.LoginPage()
	err := templates.Layout(content).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) loginHTML(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		templates.LoginError("Email and Passowrd are required!").Render(r.Context(), w)
		return
	}

	resp, err := h.firebaseRestClient.SignInWithEmailAndPassword(client.NewSignInWithPasswordRequest(email, password))
	if err != nil {
		templates.LoginError("Invalid email or passowrd").Render(r.Context(), w)
		return
	}

	h.logger.Info("User is authenticated", slog.Any("SignInWithPasswordResponse", resp))

	http.SetCookie(w, &http.Cookie{
		Name:     string(utils.IdTokenKey),
		Value:    resp.IdToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(2 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     string(utils.RefreshTokenKey),
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	w.Header().Set("HX-Redirect", "/")
}

func (h *LoginHandler) logoutHTML(w http.ResponseWriter, r *http.Request) {
	utils.ClearTokenCookies(w)

	w.Header().Set("HX-Redirect", "/")
}
