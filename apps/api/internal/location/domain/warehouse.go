package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type Warehouse struct {
	domain.AggregateRoot[string]
	Code      string
	Name      string
	Address   string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWarehouse(code, name, address string) (*Warehouse, error) {
	if code == "" {
		return nil, dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return nil, dderr.New("INVALID_NAME", "name is required", nil)
	}

	return &Warehouse{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		Code:          code,
		Name:          name,
		Address:       address,
		Active:        true,
	}, nil
}

func RestoreWarehouse(id, code, name, address string, active bool, createdAt, updatedAt time.Time) *Warehouse {
	return &Warehouse{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		Code:          code,
		Name:          name,
		Address:       address,
		Active:        active,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (w *Warehouse) Update(code, name, address string) error {
	if code == "" {
		return dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return dderr.New("INVALID_NAME", "name is required", nil)
	}

	w.Code = code
	w.Name = name
	w.Address = address

	return nil
}

func (w *Warehouse) Deactivate() {
	w.Active = false
}
