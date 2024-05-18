package poi

import "time"

type Poi struct {
	internalId  uint `gorm:"primaryKey"`
	ID          string
	CreatedOn   time.Time `gorm:"type:datetime(3)"`
	UpdatedOn   time.Time `gorm:"type:datetime(3)"`
	CreatedBy   string
	UpdatedBy   string
	Name        string
	Address     string
	City        string
	SportType   string
	Description string
	Note        string
}
