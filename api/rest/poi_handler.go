package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sportspazz/api/client"
	"github.com/sportspazz/middleware"
	"github.com/sportspazz/service/poi"
	"github.com/sportspazz/utils"
)

type PoiHandler struct {
	poiService     *poi.PoiService
	firebaseClient *client.FirebaseClient
	cloudStorage   *storage.Client
	bucket         string
}

func NewPoiHandler(poiService *poi.PoiService, firebaseClient *client.FirebaseClient, cloudStorage *storage.Client, bucket string) *PoiHandler {
	return &PoiHandler{
		poiService:     poiService,
		firebaseClient: firebaseClient,
		cloudStorage:   cloudStorage,
		bucket:         bucket,
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

	if poiRequest.GooglePlaceId != "" && p.poiService.GetPoiByGooglePlaceId(poiRequest.GooglePlaceId) != nil {
		ErrorJsonResponseWithCode(w, http.StatusConflict,
			fmt.Sprintf("Place with google place id %s already exists", poiRequest.GooglePlaceId))
		return
	}

	if poiRequest.ThumbnailUrl != "" {
		resp, err := http.Get(poiRequest.ThumbnailUrl)
		if err != nil || resp.StatusCode != http.StatusOK {
			ErrorJsonResponse(w, "cannot download image "+poiRequest.ThumbnailUrl)
			return
		}
		defer resp.Body.Close()

		objectName := "poi/thumbnails/" + uuid.New().String() + "/" + poiRequest.Name
		wc := p.cloudStorage.Bucket(p.bucket).
			Object(objectName).
			NewWriter(context.Background())
		defer wc.Close()

		if _, err := io.Copy(wc, resp.Body); err != nil {
			ErrorJsonResponse(w, "cannot download image "+poiRequest.ThumbnailUrl)
			return
		}
		poiRequest.ThumbnailUrl = fmt.Sprintf("https://storage.googleapis.com/%s/%s", p.bucket, objectName)
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
