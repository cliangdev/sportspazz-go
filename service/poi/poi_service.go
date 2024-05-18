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

func (p *PoiService) CreatePOI(createdBy, name, description, address, city, sportType, note string) (Poi, error) {
	return p.store.CreatePoi(createdBy, name, description, address, city, sportType, note), nil
}
