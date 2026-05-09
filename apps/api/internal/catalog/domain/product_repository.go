package domain

import (
	"context"

	"github.com/aleodoni/go-ddd/repository"
)

type ProductRepository interface {
	repository.Repository[string, *Product]
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*Product, error)
}
