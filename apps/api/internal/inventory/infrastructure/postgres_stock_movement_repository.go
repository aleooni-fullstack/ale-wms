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

type PostgresStockMovementRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresStockMovementRepository(pool *pgxpool.Pool) *PostgresStockMovementRepository {
	return &PostgresStockMovementRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresStockMovementRepository) Create(ctx context.Context, m *domain.StockMovement) error {
	row, err := r.queries.CreateStockMovement(ctx, db.CreateStockMovementParams{
		ID:         m.ID,
		ProductID:  m.ProductID,
		LocationID: m.LocationID,
		Type:       db.MovementType(m.Type),
		Quantity:   pgNumeric(m.Quantity),
		Note:       pgNullText(m.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating stock movement", err)
	}

	m.CreatedAt = row.CreatedAt.Time

	return nil
}

func (r *PostgresStockMovementRepository) FindByID(ctx context.Context, id string) (*domain.StockMovement, error) {
	row, err := r.queries.GetStockMovementByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding stock movement by id", err)
	}

	return stockMovementToDomain(row), nil
}

func (r *PostgresStockMovementRepository) FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*domain.StockMovement, error) {
	rows, err := r.queries.ListStockMovementsByProduct(ctx, db.ListStockMovementsByProductParams{
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock movements by product", err)
	}

	movements := make([]*domain.StockMovement, len(rows))
	for i, row := range rows {
		movements[i] = stockMovementToDomain(row)
	}

	return movements, nil
}

func (r *PostgresStockMovementRepository) FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*domain.StockMovement, error) {
	rows, err := r.queries.ListStockMovementsByLocation(ctx, db.ListStockMovementsByLocationParams{
		LocationID: locationID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock movements by location", err)
	}

	movements := make([]*domain.StockMovement, len(rows))
	for i, row := range rows {
		movements[i] = stockMovementToDomain(row)
	}

	return movements, nil
}

func stockMovementToDomain(row db.StockMovement) *domain.StockMovement {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	qty, _ := row.Quantity.Float64Value()

	return domain.RestoreStockMovement(
		row.ID,
		row.ProductID,
		row.LocationID,
		domain.MovementType(row.Type),
		qty.Float64,
		note,
		row.CreatedAt.Time,
	)
}
