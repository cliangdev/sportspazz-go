package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/utils"
)

type ContextKey string

const (
	loggerKey ContextKey = "logger"
)

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

func AuthenticateMiddleWare(firebaseClient *client.FirebaseClient, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if idTokenCookie, _ := r.Cookie(string(utils.IdTokenKey)); idTokenCookie != nil {
				token, err := firebaseClient.VerifyIDToken(idTokenCookie.Value)

				if err != nil {
					refreshToken(firebaseClient, w, r, logger)
				} else if token.Valid {
					claims := token.Claims.(jwt.MapClaims)

					email := claims[string(utils.Email)].(string)
					ctx = context.WithValue(ctx, utils.Email, email)
					ctx = context.WithValue(ctx, utils.Name, email)
					ctx = context.WithValue(ctx, utils.Logined, true)
				}
			} else {
				ctx = context.WithValue(ctx, utils.Name, "")
				ctx = context.WithValue(ctx, utils.Logined, false)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func refreshToken(firebaseClient *client.FirebaseClient, w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	// TODO: update refresh token
	if refreshTokenCookie, err := r.Cookie(string(utils.RefreshTokenKey)); err == nil {
		refreshToken := refreshTokenCookie.Value
		token, err := firebaseClient.VerifyIDToken(refreshToken)
		if err != nil {
			logger.Error("not able to refresh token", slog.Any("err", err))

			utils.ClearTokenCookies(w)

			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		tokenJSON, _ := json.Marshal(token)
		http.SetCookie(w, &http.Cookie{
			Name:     string(utils.IdTokenKey),
			Value:    string(tokenJSON),
			Path:     "/",
			Expires:  time.Unix(claims["exp"].(int64), 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}
