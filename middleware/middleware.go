package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
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
