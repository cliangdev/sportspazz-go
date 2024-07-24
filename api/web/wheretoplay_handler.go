package web

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"fmt"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/web/templates"
	"github.com/sportspazz/api/web/types"
	"github.com/sportspazz/service/poi"
	"github.com/sportspazz/utils"
)

const maxThumbnailSize = 100 * 1024 // 100 KB

const cityPlaceIdParam = "cityPlaceId"
const sportParam = "sport"
const pageSizeParam = "pageSize"
const cursorParam = "cursor"

type WhereToPlayHandler struct {
	logger          *slog.Logger
	poiService      *poi.PoiService
	cloudStorage    *storage.Client
	bucket          string
	googleMapApiKey string
}

func NewWhereToPlayHandler(logger *slog.Logger, poiService *poi.PoiService, cloudStorage *storage.Client, bucket, googleMapApiKey string) *WhereToPlayHandler {
	return &WhereToPlayHandler{
		logger:          logger,
		poiService:      poiService,
		cloudStorage:    cloudStorage,
		bucket:          bucket,
		googleMapApiKey: googleMapApiKey,
	}
}

func (h *WhereToPlayHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/wheretoplay", h.serveWhereToPlayPageHTML).Methods(http.MethodGet)
	router.HandleFunc("/wheretoplay/search", h.searchWhereToPlay).Methods(http.MethodGet)
	router.HandleFunc("/wheretoplay/new", h.serveCreateNewPlacePageHTML).Methods(http.MethodGet)
	router.HandleFunc("/wheretoplay/new", h.createNewPlace).Methods(http.MethodPost)
	router.HandleFunc("/wheretoplay/{sport}/{placeId}", h.placeDetails).Methods(http.MethodGet)
}

func (h *WhereToPlayHandler) placeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placeId := vars["placeId"]

	poi := h.poiService.GetPoiById(placeId)
	if poi == nil {
		content := templates.NotFoundMessage()
		err := templates.Layout(content).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
		return
	}
	if poi.GooglePlaceId != nil {
		details, err := getGooglePlaceDetails(*poi.GooglePlaceId, h.googleMapApiKey)
		if err == nil && details.Status == "OK" {
			w.WriteHeader(http.StatusOK)

			content := templates.PlaceDetais(details.Result)
			err := templates.MapLayout(content).Render(r.Context(), w)

			if err != nil {
				http.Error(w, "Error rendering page", http.StatusInternalServerError)
			}
			return
		}
		if err != nil {
			h.logger.Error("Cannot get place details from google api", slog.Any("err", err), slog.Any("details", details))
		}
	}
	templates.Layout(templates.NotFoundMessage()).Render(r.Context(), w)
}

func getGooglePlaceDetails(googlePlaceId, apiKey string) (*types.GooglePlaceResponse, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&key=%s&fields=current_opening_hours,formatted_address,formatted_phone_number,name,opening_hours,photos,rating,website",
		googlePlaceId, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response types.GooglePlaceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
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
	cityPlaceId := r.FormValue(cityPlaceIdParam)
	sport := r.FormValue(sportParam)

	pageSize := 15
	cursor := r.URL.Query().Get(cursorParam)
	pageSizeParam := r.URL.Query().Get(pageSizeParam)
	if pageSizeParam != "" {
		pageSize, _ = strconv.Atoi(pageSizeParam)
	}

	if cityPlaceId == "" || sport == "" {
		templates.SearchError("Pick a sport and city!").Render(r.Context(), w)
		return
	}

	pois := h.poiService.SearchPois(cityPlaceId, sport, cursor, pageSize)

	w.WriteHeader(http.StatusOK)
	templates.SearchResult(pois, sport, cityPlaceId, pageSize).Render(r.Context(), w)
}

func (h *WhereToPlayHandler) serveCreateNewPlacePageHTML(w http.ResponseWriter, r *http.Request) {
	if utils.Logined(r.Context()) {
		content := templates.CreateNewPlace()
		err := templates.MapLayout(content).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *WhereToPlayHandler) createNewPlace(w http.ResponseWriter, r *http.Request) {
	input, err := h.parseCreateNewPlaceFormInputAndValidate(r)
	if err != nil {
		templates.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}
	defer input.Thumbnail.Close()

	objectName := "poi/thumbnails/" + uuid.New().String() + "/" + input.ThumbnailFilename
	wc := h.cloudStorage.Bucket(h.bucket).
		Object(objectName).
		NewWriter(r.Context())
	defer wc.Close()

	if _, err := io.Copy(wc, input.Thumbnail); err != nil {
		templates.ErrorMessage("cannot upload thumbnail").Render(r.Context(), w)
		return
	}

	thumbnailUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s", h.bucket, objectName)

	createdBy := r.Context().Value(utils.UserIdKey).(string)
	h.poiService.CreatePoi(
		createdBy,
		input.Name,
		input.Description,
		input.Address,
		input.CityId,
		nil,
		input.Website,
		input.Sport,
		thumbnailUrl,
		"",
	)
	w.Header().Set("HX-Redirect", "/wheretoplay")
	w.WriteHeader(http.StatusSeeOther)
}

func (h *WhereToPlayHandler) parseCreateNewPlaceFormInputAndValidate(r *http.Request) (*templates.CreateNewPlaceFormInput, error) {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return nil, err
	}
	thumbnail, thumbnailHeader, err := r.FormFile("thumbnail")
	if err != nil {
		return nil, err
	}
	if thumbnailHeader.Size > maxThumbnailSize {
		return nil, fmt.Errorf("thumbnail must be less than 100 KB")
	}

	input := templates.CreateNewPlaceFormInput{
		Name:              r.FormValue("name"),
		Description:       r.FormValue("description"),
		Address:           r.FormValue("address"),
		CityId:            r.FormValue("cityPlaceId"),
		Website:           r.FormValue("website"),
		Sport:             r.FormValue("sport"),
		Thumbnail:         thumbnail,
		ThumbnailFilename: thumbnailHeader.Filename,
	}

	if len(input.Name) < 3 || len(input.Name) > 100 {
		return nil, fmt.Errorf("name must be 3 to 100 characters")
	}
	if len(input.Description) < 50 || len(input.Description) > 8000 {
		return nil, fmt.Errorf("description must be 550 to 8000 characters")
	}

	return &input, nil
}
