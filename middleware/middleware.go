package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
)

type contextKey string

const loggerKey contextKey = "logger"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func LoggerMiddleWare(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request", slog.String("method", r.Method), slog.String("path", r.URL.Path))

			ctx := context.WithValue(r.Context(), loggerKey, logger)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func ContentTypeHeaderMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		var errMsg string
		if strings.HasPrefix(r.URL.Path, "/api/v1") && contentType != "application/json" {
			errMsg = "Unsupported Content-Type, please use application/json"
		}

		if errMsg != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader((http.StatusUnsupportedMediaType))
			json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthenticateMiddleWare(firebaseClient *auth.Client, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if idTokenCookie, _ := r.Cookie("idToken"); cookie != nil {
				idToken := idTokenCookie.Value
				token, err := firebaseClient.VerifyIDToken(context.Background(), idToken)

				if err != nil {
					refreshToken(firebaseClient, w, r, logger)
				} else {
					email := token.Claims["email"].(string)
					ctx := context.WithValue(r.Context(), "email", email)

					logger.Debug("authenticated user", slog.Any("email", email))
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func refreshToken(firebaseClient *auth.Client, w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	if refreshTokenCookie, err := r.Cookie("refreshToken"); err == nil {
		refreshToken := refreshTokenCookie.Value
		token, err := firebaseClient.VerifyIDTokenAndCheckRevoked(context.Background(), refreshToken)

		if err == nil {
			logger.Error("refresh error", slog.Any("err", err))
			tokenJSON, _ := json.Marshal(token)
			http.SetCookie(w, &http.Cookie{
				Name:     "idToken",
				Value:    string(tokenJSON),
				Path:     "/",
				Expires:  time.Unix(token.Claims["exp"].(int64), 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}
	}
}
