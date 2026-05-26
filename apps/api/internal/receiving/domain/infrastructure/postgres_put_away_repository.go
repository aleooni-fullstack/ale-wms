package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresPutAwayRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresPutAwayRepository(pool *pgxpool.Pool) *PostgresPutAwayRepository {
	return &PostgresPutAwayRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresPutAwayRepository) Create(ctx context.Context, p *domain.PutAway) error {
	row, err := r.queries.CreatePutAway(ctx, db.CreatePutAwayParams{
		ID:        p.ID,
		ReceiptID: p.ReceiptID,
		Status:    db.PutAwayStatus(p.Status),
		Note:      pgNullText(p.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating put away", err)
	}

	p.CreatedAt = row.CreatedAt.Time
	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPutAwayRepository) FindByID(ctx context.Context, id string) (*domain.PutAway, error) {
	row, err := r.queries.GetPutAwayByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding put away by id", err)
	}

	return putAwayToDomain(row), nil
}

func (r *PostgresPutAwayRepository) FindByReceiptID(ctx context.Context, receiptID string) (*domain.PutAway, error) {
	row, err := r.queries.GetPutAwayByReceiptID(ctx, receiptID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding put away by receipt id", err)
	}

	return putAwayToDomain(row), nil
}

func (r *PostgresPutAwayRepository) UpdateStatus(ctx context.Context, p *domain.PutAway) error {
	row, err := r.queries.UpdatePutAwayStatus(ctx, db.UpdatePutAwayStatusParams{
		ID:     p.ID,
		Status: db.PutAwayStatus(p.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating put away status", err)
	}

	p.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPutAwayRepository) AddItem(ctx context.Context, item *domain.PutAwayItem) error {
	row, err := r.queries.CreatePutAwayItem(ctx, db.CreatePutAwayItemParams{
		ID:         item.ID,
		PutAwayID:  item.PutAwayID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   pgNumeric(item.Quantity),
		PutAway:    item.PutAway,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding put away item", err)
	}

	item.CreatedAt = row.CreatedAt.Time
	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPutAwayRepository) FindAllItems(ctx context.Context, putAwayID string) ([]*domain.PutAwayItem, error) {
	rows, err := r.queries.ListPutAwayItems(ctx, putAwayID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing put away items", err)
	}

	items := make([]*domain.PutAwayItem, len(rows))
	for i, row := range rows {
		items[i] = putAwayItemToDomain(row)
	}

	return items, nil
}

func (r *PostgresPutAwayRepository) UpdateItemPutAway(ctx context.Context, item *domain.PutAwayItem) error {
	row, err := r.queries.UpdatePutAwayItemPutAway(ctx, db.UpdatePutAwayItemPutAwayParams{
		PutAwayID: item.PutAwayID,
		ProductID: item.ProductID,
		PutAway:   item.PutAway,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating put away item", err)
	}

	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func putAwayToDomain(row db.PutAway) *domain.PutAway {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestorePutAway(
		row.ID,
		row.ReceiptID,
		domain.PutAwayStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func putAwayItemToDomain(row db.PutAwayItem) *domain.PutAwayItem {
	qty, _ := row.Quantity.Float64Value()

	return domain.RestorePutAwayItem(
		row.ID,
		row.PutAwayID,
		row.ProductID,
		row.LocationID,
		qty.Float64,
		row.PutAway,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
