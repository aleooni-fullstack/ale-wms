package domain

import (
	"time"

	"github.com/aleodoni/go-ddd/domain"
	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/nrednav/cuid2"
)

type Location struct {
	domain.AggregateRoot[string]
	ZoneID    string
	Code      string
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLocation(zoneID, code, name string) (*Location, error) {
	if zoneID == "" {
		return nil, dderr.New("INVALID_ZONE_ID", "zone_id is required", nil)
	}
	if code == "" {
		return nil, dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return nil, dderr.New("INVALID_NAME", "name is required", nil)
	}

	return &Location{
		AggregateRoot: domain.NewAggregateRoot[string](cuid2.Generate()),
		ZoneID:        zoneID,
		Code:          code,
		Name:          name,
		Active:        true,
	}, nil
}

func RestoreLocation(id, zoneID, code, name string, active bool, createdAt, updatedAt time.Time) *Location {
	return &Location{
		AggregateRoot: domain.NewAggregateRoot[string](id),
		ZoneID:        zoneID,
		Code:          code,
		Name:          name,
		Active:        active,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (l *Location) Update(code, name string) error {
	if code == "" {
		return dderr.New("INVALID_CODE", "code is required", nil)
	}
	if name == "" {
		return dderr.New("INVALID_NAME", "name is required", nil)
	}

	l.Code = code
	l.Name = name

	return nil
}

func (l *Location) Deactivate() {
	l.Active = false
}
