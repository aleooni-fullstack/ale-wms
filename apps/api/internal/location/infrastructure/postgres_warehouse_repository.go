package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresWarehouseRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresWarehouseRepository(pool *pgxpool.Pool) *PostgresWarehouseRepository {
	return &PostgresWarehouseRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresWarehouseRepository) FindByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	row, err := r.queries.GetWarehouseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding warehouse by id", err)
	}

	return warehouseToDomain(row), nil
}

func (r *PostgresWarehouseRepository) FindByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	row, err := r.queries.GetWarehouseByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding warehouse by code", err)
	}

	return warehouseToDomain(row), nil
}

func (r *PostgresWarehouseRepository) FindAll(ctx context.Context, limit, offset int32) ([]*domain.Warehouse, error) {
	rows, err := r.queries.ListWarehouses(ctx, db.ListWarehousesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing warehouses", err)
	}

	warehouses := make([]*domain.Warehouse, len(rows))
	for i, row := range rows {
		warehouses[i] = warehouseToDomain(row)
	}

	return warehouses, nil
}

func (r *PostgresWarehouseRepository) Create(ctx context.Context, w *domain.Warehouse) error {
	_, err := r.queries.CreateWarehouse(ctx, db.CreateWarehouseParams{
		ID:      w.ID,
		Code:    w.Code,
		Name:    w.Name,
		Address: pgtype.Text{String: w.Address, Valid: w.Address != ""},
		Active:  w.Active,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating warehouse", err)
	}

	return nil
}

func (r *PostgresWarehouseRepository) Update(ctx context.Context, w *domain.Warehouse) error {
	_, err := r.queries.UpdateWarehouse(ctx, db.UpdateWarehouseParams{
		ID:      w.ID,
		Code:    w.Code,
		Name:    w.Name,
		Address: pgtype.Text{String: w.Address, Valid: w.Address != ""},
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating warehouse", err)
	}

	return nil
}

func (r *PostgresWarehouseRepository) Delete(ctx context.Context, w *domain.Warehouse) error {
	err := r.queries.DeactivateWarehouse(ctx, w.ID)
	if err != nil {
		return dderr.New("DB_ERROR", "error deactivating warehouse", err)
	}

	return nil
}

func warehouseToDomain(row db.Warehouse) *domain.Warehouse {
	address := ""
	if row.Address.Valid {
		address = row.Address.String
	}

	return domain.RestoreWarehouse(
		row.ID,
		row.Code,
		row.Name,
		address,
		row.Active,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
