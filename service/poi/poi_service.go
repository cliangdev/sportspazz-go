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

func (p *PoiService) CreatePoi(createdBy, name, description, address, cityId, website, sportType, thumbnailUrl, note string) (Poi, error) {
	return p.store.CreatePoi(createdBy, name, description, address, cityId, website, sportType, thumbnailUrl, note), nil
}

func (p *PoiService) SearchPois(cityId, sport, cursor string, pageSize int) Pois {
	internalCursor := p.getInternalCursor(cursor)
	
	pois := p.store.GetPois(cityId, sport, internalCursor, pageSize+1)
	nextCursor := ""
	if len(pois) > pageSize {
		nextCursor = pois[len(pois)-1].ID
		pois = pois[:len(pois)-1]
	}

	return Pois{
		Results: pois,
		Cursor:  nextCursor,
	}
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
