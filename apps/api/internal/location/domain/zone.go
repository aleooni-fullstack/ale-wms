package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type Zone struct {
	domain.AggregateRoot[string]
	WarehouseID string
	Code        string
	Name        string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewZone(warehouseID, code, name string) (*Zone, error) {
	if warehouseID == "" {
		return nil, dderr.New("INVALID_WAREHOUSE_ID", "warehouse_id is required", nil)
	}
	if code == "" {
		return nil, dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return nil, dderr.New("INVALID_NAME", "name is required", nil)
	}

	return &Zone{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		WarehouseID:   warehouseID,
		Code:          code,
		Name:          name,
		Active:        true,
	}, nil
}

func RestoreZone(id, warehouseID, code, name string, active bool, createdAt, updatedAt time.Time) *Zone {
	return &Zone{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		WarehouseID:   warehouseID,
		Code:          code,
		Name:          name,
		Active:        active,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (z *Zone) Update(code, name string) error {
	if code == "" {
		return dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return dderr.New("INVALID_NAME", "name is required", nil)
	}

	z.Code = code
	z.Name = name

	return nil
}

func (z *Zone) Deactivate() {
	z.Active = false
}
