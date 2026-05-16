package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresStockReservationRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresStockReservationRepository(pool *pgxpool.Pool) *PostgresStockReservationRepository {
	return &PostgresStockReservationRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresStockReservationRepository) Create(ctx context.Context, res *domain.StockReservation) error {
	row, err := r.queries.CreateStockReservation(ctx, db.CreateStockReservationParams{
		ID:         res.ID,
		ProductID:  res.ProductID,
		LocationID: res.LocationID,
		Quantity:   pgNumeric(res.Quantity),
		Status:     db.ReservationStatus(res.Status),
		Reference:  pgNullText(res.Reference),
		Note:       pgNullText(res.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating stock reservation", err)
	}

	res.CreatedAt = row.CreatedAt.Time
	res.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresStockReservationRepository) FindByID(ctx context.Context, id string) (*domain.StockReservation, error) {
	row, err := r.queries.GetStockReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding stock reservation by id", err)
	}

	return stockReservationToDomain(row), nil
}

func (r *PostgresStockReservationRepository) FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*domain.StockReservation, error) {
	rows, err := r.queries.ListStockReservationsByProduct(ctx, db.ListStockReservationsByProductParams{
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock reservations by product", err)
	}

	return stockReservationsToDomain(rows), nil
}

func (r *PostgresStockReservationRepository) FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*domain.StockReservation, error) {
	rows, err := r.queries.ListStockReservationsByLocation(ctx, db.ListStockReservationsByLocationParams{
		LocationID: locationID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock reservations by location", err)
	}

	return stockReservationsToDomain(rows), nil
}

func (r *PostgresStockReservationRepository) FindAllByReference(ctx context.Context, reference string) ([]*domain.StockReservation, error) {
	rows, err := r.queries.ListStockReservationsByReference(ctx, pgNullText(reference))
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock reservations by reference", err)
	}

	return stockReservationsToDomain(rows), nil
}

func (r *PostgresStockReservationRepository) UpdateStatus(ctx context.Context, res *domain.StockReservation) error {
	row, err := r.queries.UpdateStockReservationStatus(ctx, db.UpdateStockReservationStatusParams{
		ID:     res.ID,
		Status: db.ReservationStatus(res.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating stock reservation status", err)
	}

	res.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func stockReservationToDomain(row db.StockReservation) *domain.StockReservation {
	reference := ""
	if row.Reference.Valid {
		reference = row.Reference.String
	}

	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	qty, _ := row.Quantity.Float64Value()

	return domain.RestoreStockReservation(
		row.ID,
		row.ProductID,
		row.LocationID,
		qty.Float64,
		domain.ReservationStatus(row.Status),
		reference,
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func stockReservationsToDomain(rows []db.StockReservation) []*domain.StockReservation {
	reservations := make([]*domain.StockReservation, len(rows))
	for i, row := range rows {
		reservations[i] = stockReservationToDomain(row)
	}

	return reservations
}
