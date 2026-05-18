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

type PostgresPickingRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresPickingRepository(pool *pgxpool.Pool) *PostgresPickingRepository {
	return &PostgresPickingRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresPickingRepository) Create(ctx context.Context, p *domain.Picking) error {
	row, err := r.queries.CreatePicking(ctx, db.CreatePickingParams{
		ID:      p.ID,
		OrderID: p.OrderID,
		Status:  db.PickingStatus(p.Status),
		Note:    pgNullText(p.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating picking", err)
	}

	p.CreatedAt = row.CreatedAt.Time
	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPickingRepository) FindByID(ctx context.Context, id string) (*domain.Picking, error) {
	row, err := r.queries.GetPickingByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding picking by id", err)
	}

	return pickingToDomain(row), nil
}

func (r *PostgresPickingRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Picking, error) {
	row, err := r.queries.GetPickingByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding picking by order id", err)
	}

	return pickingToDomain(row), nil
}

func (r *PostgresPickingRepository) UpdateStatus(ctx context.Context, p *domain.Picking) error {
	row, err := r.queries.UpdatePickingStatus(ctx, db.UpdatePickingStatusParams{
		ID:     p.ID,
		Status: db.PickingStatus(p.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating picking status", err)
	}

	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPickingRepository) AddItem(ctx context.Context, item *domain.PickingItem) error {
	row, err := r.queries.CreatePickingItem(ctx, db.CreatePickingItemParams{
		ID:         item.ID,
		PickingID:  item.PickingID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   pgNumeric(item.Quantity),
		Picked:     item.Picked,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding picking item", err)
	}

	item.CreatedAt = row.CreatedAt.Time
	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPickingRepository) FindAllItems(ctx context.Context, pickingID string) ([]*domain.PickingItem, error) {
	rows, err := r.queries.ListPickingItems(ctx, pickingID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing picking items", err)
	}

	items := make([]*domain.PickingItem, len(rows))
	for i, row := range rows {
		items[i] = pickingItemToDomain(row)
	}

	return items, nil
}

func (r *PostgresPickingRepository) UpdateItemPicked(ctx context.Context, item *domain.PickingItem) error {
	row, err := r.queries.UpdatePickingItemPicked(ctx, db.UpdatePickingItemPickedParams{
		PickingID: item.PickingID,
		ProductID: item.ProductID,
		Picked:    item.Picked,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating picking item picked", err)
	}

	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func pickingToDomain(row db.Picking) *domain.Picking {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestorePicking(
		row.ID,
		row.OrderID,
		domain.PickingStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func pickingItemToDomain(row db.PickingItem) *domain.PickingItem {
	qty, _ := row.Quantity.Float64Value()

	return domain.RestorePickingItem(
		row.ID,
		row.PickingID,
		row.ProductID,
		row.LocationID,
		qty.Float64,
		row.Picked,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
