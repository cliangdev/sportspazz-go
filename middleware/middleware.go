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

const (
	loggerKey       = "logger"
	idTokenKey      = "idToken"
	refreshTokenKey = "refreshToken"
)

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
			ctx := r.Context()
			if idTokenCookie, _ := r.Cookie(idTokenKey); idTokenCookie != nil {
				idToken := idTokenCookie.Value
				token, err := firebaseClient.VerifyIDToken(context.Background(), idToken)

				if err != nil {
					refreshToken(firebaseClient, w, r, logger)
				} else {
					email := token.Claims["email"].(string)
					ctx = context.WithValue(ctx, "email", email)
					ctx = context.WithValue(ctx, "name", email)
					ctx = context.WithValue(ctx, "logined", true)
				}
			} else {
				ctx = context.WithValue(ctx, "name", "")
				ctx = context.WithValue(ctx, "logined", false)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func refreshToken(firebaseClient *auth.Client, w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	if refreshTokenCookie, err := r.Cookie("refreshToken"); err == nil {
		refreshToken := refreshTokenCookie.Value
		token, err := firebaseClient.VerifyIDTokenAndCheckRevoked(context.Background(), refreshToken)

		logger.Info("refreshing token", slog.Any("err", err))

		if err == nil {
			logger.Error("refresh error", slog.Any("err", err))
			tokenJSON, _ := json.Marshal(token)
			http.SetCookie(w, &http.Cookie{
				Name:     idTokenKey,
				Value:    string(tokenJSON),
				Path:     "/",
				Expires:  time.Unix(token.Claims["exp"].(int64), 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     idTokenKey,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     refreshTokenKey,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
		}
	}
}
