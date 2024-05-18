package utils

import (
	"net/http"
	"time"
)

type ContextKey string

const (
	IdTokenKey      ContextKey = "idToken"
	RefreshTokenKey ContextKey = "refreshToken"
	LoggerKey       ContextKey = "logger"
	// For UI
	UserIdKey  ContextKey = "userId"
	EmailKey   ContextKey = "email"
	NameKey    ContextKey = "name"
	LoginedKey ContextKey = "logined"
)

func ClearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(IdTokenKey),
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     string(RefreshTokenKey),
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
