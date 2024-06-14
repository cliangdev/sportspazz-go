package cli

type PlacesResponse struct {
	NextPageToken string `json:"next_page_token"`
	Results       []Place `json:"results"`
}

// Place represents a place from the Google Places API response.
type Place struct {
	Name              string       `json:"name"`
	PlaceID           string       `json:"place_id"`
	Address           string       `json:"formatted_address"`
	Types             []string     `json:"types"`
	Geometry          Geometry     `json:"geometry"`
	OpeningHours      OpeningHours `json:"opening_hours,omitempty"`
	Photos            []Photo      `json:"photos,omitempty"`
	Rating            float64      `json:"rating,omitempty"`
	UserRatingsTotal  int          `json:"user_ratings_total,omitempty"`
	PriceLevel        int          `json:"price_level,omitempty"`
	Vicinity          string       `json:"vicinity,omitempty"`
	PermanentlyClosed bool         `json:"permanently_closed,omitempty"`
	PlusCode          PlusCode     `json:"plus_code,omitempty"`
	Website           string       `json:"website,omitempty"`
	PhoneNumber       string       `json:"formatted_phone_number,omitempty"`
	BusinessStatus    string       `json:"business_status,omitempty"`
}

// Geometry represents the geometry field in the API response.
type Geometry struct {
	Location Location `json:"location"`
	Viewport Viewport `json:"viewport"`
}

// Location represents a geographical location.
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Viewport represents a viewport bounding box.
type Viewport struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

// OpeningHours represents the opening hours information.
type OpeningHours struct {
	OpenNow     bool     `json:"open_now"`
	Periods     []Period `json:"periods,omitempty"`
	WeekdayText []string `json:"weekday_text,omitempty"`
}

// Period represents a period in the opening hours.
type Period struct {
	Open  TimeOfDay `json:"open"`
	Close TimeOfDay `json:"close,omitempty"`
}

// TimeOfDay represents a time of day.
type TimeOfDay struct {
	Day  int    `json:"day"`
	Time string `json:"time"`
}

// Photo represents a photo from the API response.
type Photo struct {
	PhotoReference   string   `json:"photo_reference"`
	Height           int      `json:"height"`
	Width            int      `json:"width"`
	HTMLAttributions []string `json:"html_attributions"`
}

// PlusCode represents a plus code for the location.
type PlusCode struct {
	GlobalCode   string `json:"global_code"`
	CompoundCode string `json:"compound_code"`
}

type PlaceDetails struct {
	Name             string   `json:"name"`
	PlaceID          string   `json:"place_id"`
	Address          string   `json:"formatted_address"`
	Types            []string `json:"types"`
	Rating           float64  `json:"rating"`
	UserRatingsTotal int      `json:"user_ratings_total"`
	Geometry         Geometry `json:"geometry"`
	Photos           []Photo  `json:"photos"`
	Website          string   `json:"website"`
	PhoneNumber      string   `json:"formatted_phone_number"`
	BusinessStatus   string   `json:"business_status"`
	Reviews          []Review `json:"reviews"`
}

type Review struct {
	AuthorName string `json:"author_name"`
	Rating     int    `json:"rating"`
	Text       string `json:"text"`
	Time       int64  `json:"time"`
}

type POI struct {
    Name          string `json:"name"`
    Address       string `json:"address"`
    CityID        string `json:"city_id"`
    SportType     string `json:"sport_type"`
    ThumbnailURL  string `json:"thumbnail_url"`
    Description   string `json:"description"`
}