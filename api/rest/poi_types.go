package rest_api

type CreatePoiRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=200"`
	Address      string `json:"address"`
	Website      string `json:"website"`
	CityId       string `json:"city_id" validate:"required,min=2,max=255"`
	SportType    string `json:"sport_type" validate:"required,min=2,max=50"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Description  string `json:"description"`
	Note         string `json:"note"`
}

type PoiResponse struct {
	ID          string `json:"id"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	SportType   string `json:"sport_type"`
	Description string `json:"description"`
	Note        string `json:"note"`
}
