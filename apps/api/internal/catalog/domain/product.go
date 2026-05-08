package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type Product struct {
	domain.AggregateRoot[string]
	SKU         string
	Name        string
	Description string
	Unit        string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(sku, name, description, unit string) (*Product, error) {
	if sku == "" {
		return nil, dderr.New("INVALID_SKU", "sku is required", nil)
	}
	if name == "" {
		return nil, dderr.New("INVALID_NAME", "name is required", nil)
	}
	if unit == "" {
		return nil, dderr.New("INVALID_UNIT", "unit is required", nil)
	}

	p := &Product{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		SKU:           sku,
		Name:          name,
		Description:   description,
		Unit:          unit,
		Active:        true,
	}

	return p, nil
}

func RestoreProduct(id, sku, name, description, unit string, active bool, createdAt, updatedAt time.Time) *Product {
	return &Product{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		SKU:           sku,
		Name:          name,
		Description:   description,
		Unit:          unit,
		Active:        active,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (p *Product) Update(sku, name, description, unit string) error {
	if sku == "" {
		return dderr.New("INVALID_SKU", "sku is required", nil)
	}
	if name == "" {
		return dderr.New("INVALID_NAME", "name is required", nil)
	}
	if unit == "" {
		return dderr.New("INVALID_UNIT", "unit is required", nil)
	}

	p.SKU = sku
	p.Name = name
	p.Description = description
	p.Unit = unit

	return nil
}

func (p *Product) Deactivate() {
	p.Active = false
}
