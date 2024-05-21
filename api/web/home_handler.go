package web

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sportspazz/api/web/templates"
)

type HomeHandler struct {
	logger *slog.Logger
}

func NewHomeHandler(logger *slog.Logger) *HomeHandler {
	return &HomeHandler{
		logger: logger,
	}
}

func (h *HomeHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.serveHTTP).Methods(http.MethodGet)
}

func (h *HomeHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Index()
	if err := templates.Layout(c).Render(r.Context(), w); err != nil {
		http.Error(w, "Cannot render home page", http.StatusInternalServerError)
	}
}
