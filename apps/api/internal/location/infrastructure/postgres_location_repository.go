package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresLocationRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresLocationRepository(pool *pgxpool.Pool) *PostgresLocationRepository {
	return &PostgresLocationRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresLocationRepository) FindByID(ctx context.Context, id string) (*domain.Location, error) {
	row, err := r.queries.GetLocationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding location by id", err)
	}

	return locationToDomain(row), nil
}

func (r *PostgresLocationRepository) FindByCode(ctx context.Context, zoneID, code string) (*domain.Location, error) {
	row, err := r.queries.GetLocationByCode(ctx, db.GetLocationByCodeParams{
		ZoneID: zoneID,
		Code:   code,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding location by code", err)
	}

	return locationToDomain(row), nil
}

func (r *PostgresLocationRepository) FindAllByZone(ctx context.Context, zoneID string, limit, offset int32) ([]*domain.Location, error) {
	rows, err := r.queries.ListLocationsByZone(ctx, db.ListLocationsByZoneParams{
		ZoneID: zoneID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing locations", err)
	}

	locations := make([]*domain.Location, len(rows))
	for i, row := range rows {
		locations[i] = locationToDomain(row)
	}

	return locations, nil
}

func (r *PostgresLocationRepository) Create(ctx context.Context, l *domain.Location) error {
	_, err := r.queries.CreateLocation(ctx, db.CreateLocationParams{
		ID:     l.ID,
		ZoneID: l.ZoneID,
		Code:   l.Code,
		Name:   l.Name,
		Active: l.Active,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating location", err)
	}

	return nil
}

func (r *PostgresLocationRepository) Update(ctx context.Context, l *domain.Location) error {
	_, err := r.queries.UpdateLocation(ctx, db.UpdateLocationParams{
		ID:   l.ID,
		Code: l.Code,
		Name: l.Name,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating location", err)
	}

	return nil
}

func (r *PostgresLocationRepository) Delete(ctx context.Context, l *domain.Location) error {
	err := r.queries.DeactivateLocation(ctx, l.ID)
	if err != nil {
		return dderr.New("DB_ERROR", "error deactivating location", err)
	}

	return nil
}

func locationToDomain(row db.Location) *domain.Location {
	return domain.RestoreLocation(
		row.ID,
		row.ZoneID,
		row.Code,
		row.Name,
		row.Active,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
