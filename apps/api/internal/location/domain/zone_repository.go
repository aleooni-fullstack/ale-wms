package domain

import (
	"context"

	"github.com/aleodoni/go-ddd/repository"
)

type ZoneRepository interface {
	repository.Repository[string, *Zone]
	FindByCode(ctx context.Context, warehouseID, code string) (*Zone, error)
	FindAllByWarehouse(ctx context.Context, warehouseID string, limit, offset int32) ([]*Zone, error)
}
