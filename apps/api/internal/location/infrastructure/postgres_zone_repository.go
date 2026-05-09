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

type PostgresZoneRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresZoneRepository(pool *pgxpool.Pool) *PostgresZoneRepository {
	return &PostgresZoneRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresZoneRepository) FindByID(ctx context.Context, id string) (*domain.Zone, error) {
	row, err := r.queries.GetZoneByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding zone by id", err)
	}

	return zoneToDomain(row), nil
}

func (r *PostgresZoneRepository) FindByCode(ctx context.Context, warehouseID, code string) (*domain.Zone, error) {
	row, err := r.queries.GetZoneByCode(ctx, db.GetZoneByCodeParams{
		WarehouseID: warehouseID,
		Code:        code,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding zone by code", err)
	}

	return zoneToDomain(row), nil
}

func (r *PostgresZoneRepository) FindAllByWarehouse(ctx context.Context, warehouseID string, limit, offset int32) ([]*domain.Zone, error) {
	rows, err := r.queries.ListZonesByWarehouse(ctx, db.ListZonesByWarehouseParams{
		WarehouseID: warehouseID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing zones", err)
	}

	zones := make([]*domain.Zone, len(rows))
	for i, row := range rows {
		zones[i] = zoneToDomain(row)
	}

	return zones, nil
}

func (r *PostgresZoneRepository) Create(ctx context.Context, z *domain.Zone) error {
	_, err := r.queries.CreateZone(ctx, db.CreateZoneParams{
		ID:          z.ID,
		WarehouseID: z.WarehouseID,
		Code:        z.Code,
		Name:        z.Name,
		Active:      z.Active,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating zone", err)
	}

	return nil
}

func (r *PostgresZoneRepository) Update(ctx context.Context, z *domain.Zone) error {
	_, err := r.queries.UpdateZone(ctx, db.UpdateZoneParams{
		ID:   z.ID,
		Code: z.Code,
		Name: z.Name,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating zone", err)
	}

	return nil
}

func (r *PostgresZoneRepository) Delete(ctx context.Context, z *domain.Zone) error {
	err := r.queries.DeactivateZone(ctx, z.ID)
	if err != nil {
		return dderr.New("DB_ERROR", "error deactivating zone", err)
	}

	return nil
}

func zoneToDomain(row db.Zone) *domain.Zone {
	return domain.RestoreZone(
		row.ID,
		row.WarehouseID,
		row.Code,
		row.Name,
		row.Active,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
