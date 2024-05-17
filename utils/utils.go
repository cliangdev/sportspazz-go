package utils

import (
	"net/http"
	"time"
)

type ContextKey string

const (
	IdTokenKey      ContextKey = "id-token"
	RefreshTokenKey ContextKey = "refresh-token"
	// For UI
	Email ContextKey = "email"
	Name ContextKey = "name"
	Logined ContextKey = "logined"
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
