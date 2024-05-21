package rest_api

type CreatePoiRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=200"`
	Address      string `json:"address"`
	CityId       string `json:"city_id" validate:"required,min=2,max=255"`
	SportType    string `json:"sport_type" validate:"required,min=2,max=50"`
	ThumbnailUrl string `json:"thumbnail_url" validate:"required,min=10"`
	Description  string `json:"description" validate:"required,min=50,max=8000"`
	Note         string `json:"note" validate:"min=50,max=8000"`
}

type PoiResponse struct {
	ID          string `json:"id"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	SportType   string `json:"sport_type"`
	Description string `json:"description"`
	Note        string `json:"note"`
}
