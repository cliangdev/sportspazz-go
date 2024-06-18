package rest_api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/middleware"
	"github.com/sportspazz/service/poi"
	"github.com/sportspazz/utils"
)

type PoiHandler struct {
	poiService     *poi.PoiService
	firebaseClient *client.FirebaseClient
}

func NewPoiHandler(poiService *poi.PoiService, firebaseClient *client.FirebaseClient) *PoiHandler {
	return &PoiHandler{
		poiService:     poiService,
		firebaseClient: firebaseClient,
	}
}

func (h *PoiHandler) RegisterRoutes(router *mux.Router) {
	router.Handle("/pois", middleware.RestAuthMiddleware(http.HandlerFunc(h.createPoi), h.firebaseClient)).Methods(http.MethodPost)
}

func (p *PoiHandler) createPoi(w http.ResponseWriter, r *http.Request) {
	var poiRequest CreatePoiRequest
	if err := json.NewDecoder(r.Body).Decode(&poiRequest); err != nil {
		InvalidJsonResponse(w)
		return
	}

	if err := validCreatePoiRequest(poiRequest); err != nil {
		ErrorJsonResponse(w, err.Error())
		return
	}

	createdBy := r.Context().Value(utils.UserIdKey).(string)
	newPoi, err := p.poiService.CreatePoi(
		createdBy,
		poiRequest.Name,
		poiRequest.Description,
		poiRequest.Address,
		poiRequest.CityId,
		poiRequest.Website,
		poiRequest.SportType,
		poiRequest.ThumbnailUrl,
		poiRequest.Note)

	if err != nil {
		ErrorJsonResponse(w, err.Error())
		return
	}

	dateFmt := `2006-01-02T15:04:05.000Z`
	poiResponse := PoiResponse{
		ID:          newPoi.ID,
		CreatedOn:   newPoi.CreatedOn.Format(dateFmt),
		UpdatedOn:   newPoi.UpdatedOn.Format(dateFmt),
		CreatedBy:   newPoi.CreatedBy,
		UpdatedBy:   newPoi.UpdatedBy,
		Name:        newPoi.Name,
		Address:     newPoi.Address,
		SportType:   newPoi.SportType,
		Description: newPoi.Description,
		Note:        newPoi.Note,
	}

	JsonResponse(poiResponse, w)
}

func validCreatePoiRequest(poiRequest CreatePoiRequest) error {
	validate := validator.New()
	if err := validate.Struct(poiRequest); err != nil {
		return err
	}

	return nil
}
