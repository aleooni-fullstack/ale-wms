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

type PostgresShippingRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresShippingRepository(pool *pgxpool.Pool) *PostgresShippingRepository {
	return &PostgresShippingRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresShippingRepository) Create(ctx context.Context, s *domain.Shipping) error {
	row, err := r.queries.CreateShipping(ctx, db.CreateShippingParams{
		ID:           s.ID,
		OrderID:      s.OrderID,
		PackingID:    s.PackingID,
		Status:       db.ShippingStatus(s.Status),
		TrackingCode: pgNullText(s.TrackingCode),
		Note:         pgNullText(s.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating shipping", err)
	}

	s.CreatedAt = row.CreatedAt.Time
	s.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresShippingRepository) FindByID(ctx context.Context, id string) (*domain.Shipping, error) {
	row, err := r.queries.GetShippingByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding shipping by id", err)
	}

	return shippingToDomain(row), nil
}

func (r *PostgresShippingRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Shipping, error) {
	row, err := r.queries.GetShippingByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding shipping by order id", err)
	}

	return shippingToDomain(row), nil
}

func (r *PostgresShippingRepository) UpdateStatus(ctx context.Context, s *domain.Shipping) error {
	row, err := r.queries.UpdateShippingStatus(ctx, db.UpdateShippingStatusParams{
		ID:           s.ID,
		Status:       db.ShippingStatus(s.Status),
		TrackingCode: pgNullText(s.TrackingCode),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating shipping status", err)
	}

	s.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func shippingToDomain(row db.Shipping) *domain.Shipping {
	trackingCode := ""
	if row.TrackingCode.Valid {
		trackingCode = row.TrackingCode.String
	}

	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestoreShipping(
		row.ID,
		row.OrderID,
		row.PackingID,
		domain.ShippingStatus(row.Status),
		trackingCode,
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
