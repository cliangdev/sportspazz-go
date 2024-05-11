package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const traceIDKey = "trace.id"

func LoggerMiddleWare(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			func() {
				logger = logger.With(traceIDKey, uuid.New().String())

				ctx := context.WithValue(r.Context(), "logger", logger)
				r = r.WithContext(ctx)
			}()
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
