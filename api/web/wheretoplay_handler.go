package web

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/service/poi"
)

type WhereToPlayHandler struct {
	logger     *slog.Logger
	poiService *poi.PoiService
}

func NewWhereToPlayHandler(logger *slog.Logger, poiService *poi.PoiService) *WhereToPlayHandler {
	return &WhereToPlayHandler{
		logger:     logger,
		poiService: poiService,
	}
}

func (h *WhereToPlayHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/wheretoplay", h.serveWhereToPlayPageHTML).Methods(http.MethodGet)
	router.HandleFunc("/wheretoplay", h.searchWhereToPlay).Methods(http.MethodPost)
}

func (h *WhereToPlayHandler) serveWhereToPlayPageHTML(w http.ResponseWriter, r *http.Request) {
	content := templates.WhereToPlayPage()
	err := templates.MapLayout(content).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func (h *WhereToPlayHandler) searchWhereToPlay(w http.ResponseWriter, r *http.Request) {
	cityPlaceId := r.FormValue("cityPlaceId")
	sport := r.FormValue("sport")

	if cityPlaceId == "" || sport == "" {
		templates.SearchError("Pick a sport and city!").Render(r.Context(), w)
		return
	}

	h.logger.Info("", slog.Any("cityPlaceId", cityPlaceId))

	w.WriteHeader(http.StatusOK)
	templates.SearchResult(h.poiService.SearchPois(cityPlaceId, sport, "", 50)).Render(r.Context(), w)
}
