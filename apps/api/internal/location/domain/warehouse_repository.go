package domain

import (
	"context"

	"github.com/aleodoni/go-ddd/repository"
)

type WarehouseRepository interface {
	repository.Repository[string, *Warehouse]
	FindByCode(ctx context.Context, code string) (*Warehouse, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*Warehouse, error)
}
