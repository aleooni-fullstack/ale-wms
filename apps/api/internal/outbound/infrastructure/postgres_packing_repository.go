package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresPackingRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresPackingRepository(pool *pgxpool.Pool) *PostgresPackingRepository {
	return &PostgresPackingRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresPackingRepository) Create(ctx context.Context, p *domain.Packing) error {
	row, err := r.queries.CreatePacking(ctx, db.CreatePackingParams{
		ID:        p.ID,
		OrderID:   p.OrderID,
		PickingID: p.PickingID,
		Status:    db.PackingStatus(p.Status),
		Note:      pgNullText(p.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating packing", err)
	}

	p.CreatedAt = row.CreatedAt.Time
	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPackingRepository) FindByID(ctx context.Context, id string) (*domain.Packing, error) {
	row, err := r.queries.GetPackingByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding packing by id", err)
	}

	return packingToDomain(row), nil
}

func (r *PostgresPackingRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Packing, error) {
	row, err := r.queries.GetPackingByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding packing by order id", err)
	}

	return packingToDomain(row), nil
}

func (r *PostgresPackingRepository) UpdateStatus(ctx context.Context, p *domain.Packing) error {
	row, err := r.queries.UpdatePackingStatus(ctx, db.UpdatePackingStatusParams{
		ID:     p.ID,
		Status: db.PackingStatus(p.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating packing status", err)
	}

	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func packingToDomain(row db.Packing) *domain.Packing {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestorePacking(
		row.ID,
		row.OrderID,
		row.PickingID,
		domain.PackingStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
