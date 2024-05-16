package web

import (
	"log/slog"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/service/user"
)

type LoginHandler struct {
	userService         *user.UserService
	firebaseAdminClient *auth.Client
	firebaseRestClient  *client.FirebaseClient
	logger              *slog.Logger
}

func NewLoginHandler(userService *user.UserService, firebaseAdminClient *auth.Client, firebaseRestClient *client.FirebaseClient, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		userService:         userService,
		firebaseAdminClient: firebaseAdminClient,
		firebaseRestClient:  firebaseRestClient,
		logger:              logger,
	}
}

func (h *LoginHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.serveLoginPageHTML).Methods(http.MethodGet)
	router.HandleFunc("/login", h.loginHTML).Methods(http.MethodPost)
}

func (h *LoginHandler) serveLoginPageHTML(w http.ResponseWriter, r *http.Request) {
	content := templates.LoginPage()
	err := templates.Layout(content, "Sportspazz").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) loginHTML(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	h.logger.Info("", slog.String("email", email), slog.String("password", password))

	if email == "" || password == "" {
		templates.LoginError("Email and Passowrd are required!").Render(r.Context(), w)
		return
	}

	resp, err := h.firebaseRestClient.SignInWithEmailAndPassword(client.NewSignInWithPasswordRequest(email, password))
	if err != nil {
		templates.LoginError("Invalid email or passowrd").Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "idToken",
		Value:    resp.IdToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("HX-Redirect", "/")
}
