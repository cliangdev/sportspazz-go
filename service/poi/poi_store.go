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

func (s *PoiStore) CreatePoi(createdBy, name, description, address, cityId, sportType, thumbnailUrl, note string) Poi {
	now := time.Now().UTC()
	poi := Poi{
		ID:           uuid.New().String(),
		CreatedOn:    now,
		UpdatedOn:    now,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
		Name:         name,
		Address:      address,
		CityId:       cityId,
		SportType:    sportType,
		Description:  description,
		ThumbnailUrl: thumbnailUrl,
		Note:         note,
	}

	if err := s.db.Create(poi).Error; err != nil {
		s.logger.Error("not able to create a new poi", slog.Any("err", err))
	}

	return poi
}

func (s *PoiStore) GetInternalCursor(cursor string) (uint, error) {
	var id uint
	if err := s.db.Model(&Poi{}).
		Select("internal_id").
		Where("id = ?", cursor).
		First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PoiStore) GetLatestPoiInternalId() uint {
	var internalId uint
	if err := s.db.Model(&Poi{}).
		Select("internal_id").
		Order("internal_id DESC").
		Limit(1).Scan(&internalId).Error; err != nil {
		return 0
	}

	return internalId
}

func (s *PoiStore) GetPois(cityId, sport string, cursor uint, pageSize int) []Poi {
	var pois []Poi
	s.db.Where("city_id = ? AND sport_type = ? AND internal_id <= ?",
		cityId, sport, cursor).
		Order("internal_id DESC").
		Limit(pageSize).
		Find(&pois)

	return pois
}
