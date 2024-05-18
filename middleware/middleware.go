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

func LoggerMiddleWare(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request", slog.String("method", r.Method), slog.String("path", r.URL.Path))

			ctx := context.WithValue(r.Context(), utils.LoggerKey, logger)
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
					ctx = updateContext(token, ctx)
				}
			} else {
				ctx = context.WithValue(ctx, utils.LoginedKey, false)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// authentication middleware for secured REST endpoints
func RestAuthMiddleware(next http.Handler, firebaseClient *client.FirebaseClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		idToken := parts[1]
		if firebaseClient == nil {
			return
		}

		token, err := firebaseClient.VerifyIDToken(idToken)
		if err != nil || !token.Valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		ctx := updateContext(token, r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func updateContext(token *jwt.Token, ctx context.Context) context.Context {
	claims := token.Claims.(jwt.MapClaims)

	userId := claims["user_id"]
	email := claims[string(utils.EmailKey)].(string)

	ctx = context.WithValue(ctx, utils.UserIdKey, userId)
	ctx = context.WithValue(ctx, utils.EmailKey, email)
	ctx = context.WithValue(ctx, utils.NameKey, email)
	ctx = context.WithValue(ctx, utils.LoginedKey, true)

	return ctx
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
