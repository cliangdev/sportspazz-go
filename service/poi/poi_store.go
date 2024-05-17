package poi

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PoiStore struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPoiStore(db *gorm.DB, logger *slog.Logger) *PoiStore {
	return &PoiStore{
		db:     db,
		logger: logger,
	}
}

func (p PoiStore) CreatePoi(createdBy, name, description, address, city, sportType, note string) Poi {
	now := time.Now().UTC()
	poi := Poi{
		ID:          uuid.New().String(),
		CreatedOn:   now,
		UpdatedOn:   now,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
		Name:        name,
		Address:     address,
		City:        city,
		SportType:   sportType,
		Description: description,
		Note:        note,
	}

	if err := p.db.Create(poi).Error; err != nil {
		p.logger.Error("not able to create a new poi", slog.Any("err", err))
	}

	return poi
}
