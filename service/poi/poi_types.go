package poi

import "time"

type Poi struct {
	internalId    uint `gorm:"primaryKey"`
	ID            string
	CreatedOn     time.Time `gorm:"type:timestamp(3) without time zone" json:"created_on"`
	UpdatedOn     time.Time `gorm:"type:timestamp(3) without time zone" json:"updated_on"`
	CreatedBy     string
	UpdatedBy     string
	Name          string
	Address       string
	Website       string
	CityId        string
	GooglePlaceId *string
	SportType     string
	ThumbnailUrl  string
	Description   string
	Note          string
}

type Pois struct {
	Results []Poi
	Cursor  string
}
