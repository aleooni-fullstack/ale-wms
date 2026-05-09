package domain

import (
	"context"

	"github.com/aleodoni/go-ddd/repository"
)

type LocationRepository interface {
	repository.Repository[string, *Location]
	FindByCode(ctx context.Context, zoneID, code string) (*Location, error)
	FindAllByZone(ctx context.Context, zoneID string, limit, offset int32) ([]*Location, error)
}
