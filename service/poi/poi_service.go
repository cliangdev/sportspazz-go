package poi

import "log/slog"

type PoiService struct {
	store  *PoiStore
	logger *slog.Logger
}

func NewPoiService(store *PoiStore, logger *slog.Logger) *PoiService {
	return &PoiService{
		store:  store,
		logger: logger,
	}
}

func (p *PoiService) CreatePoi(createdBy, name, description, address, cityId, sportType, thumbnailUrl, note string) (Poi, error) {
	return p.store.CreatePoi(createdBy, name, description, address, cityId, sportType, thumbnailUrl, note), nil
}

func (p *PoiService) SearchPois(cityId, sport, cursor string, pageSize int) []Poi {
	internalCursor := p.getInternalCursor(cursor)

	return p.store.GetPois(cityId, sport, internalCursor, pageSize)
}

func (p *PoiService) getInternalCursor(cursor string) uint {
	if cursor == "" {
		return p.store.GetLatestPoiInternalId()
	}
	internalCursor, err := p.store.GetInternalCursor(cursor)
	if err != nil {
		return p.store.GetLatestPoiInternalId()
	}
	return internalCursor
}
