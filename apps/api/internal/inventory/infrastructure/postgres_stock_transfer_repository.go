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

type PostgresStockTransferRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresStockTransferRepository(pool *pgxpool.Pool) *PostgresStockTransferRepository {
	return &PostgresStockTransferRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresStockTransferRepository) Create(ctx context.Context, t *domain.StockTransfer) error {
	row, err := r.queries.CreateStockTransfer(ctx, db.CreateStockTransferParams{
		ID:             t.ID,
		ProductID:      t.ProductID,
		FromLocationID: t.FromLocationID,
		ToLocationID:   t.ToLocationID,
		Quantity:       pgNumeric(t.Quantity),
		Status:         db.TransferStatus(t.Status),
		Note:           pgNullText(t.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating stock transfer", err)
	}

	t.CreatedAt = row.CreatedAt.Time
	t.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresStockTransferRepository) FindByID(ctx context.Context, id string) (*domain.StockTransfer, error) {
	row, err := r.queries.GetStockTransferByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding stock transfer by id", err)
	}

	return stockTransferToDomain(row), nil
}

func (r *PostgresStockTransferRepository) FindAllByProduct(ctx context.Context, productID string, limit, offset int32) ([]*domain.StockTransfer, error) {
	rows, err := r.queries.ListStockTransfersByProduct(ctx, db.ListStockTransfersByProductParams{
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock transfers by product", err)
	}

	return stockTransfersToDomain(rows), nil
}

func (r *PostgresStockTransferRepository) FindAllByFromLocation(ctx context.Context, fromLocationID string, limit, offset int32) ([]*domain.StockTransfer, error) {
	rows, err := r.queries.ListStockTransfersByFromLocation(ctx, db.ListStockTransfersByFromLocationParams{
		FromLocationID: fromLocationID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock transfers by from location", err)
	}

	return stockTransfersToDomain(rows), nil
}

func (r *PostgresStockTransferRepository) FindAllByToLocation(ctx context.Context, toLocationID string, limit, offset int32) ([]*domain.StockTransfer, error) {
	rows, err := r.queries.ListStockTransfersByToLocation(ctx, db.ListStockTransfersByToLocationParams{
		ToLocationID: toLocationID,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock transfers by to location", err)
	}

	return stockTransfersToDomain(rows), nil
}

func (r *PostgresStockTransferRepository) UpdateStatus(ctx context.Context, t *domain.StockTransfer) error {
	row, err := r.queries.UpdateStockTransferStatus(ctx, db.UpdateStockTransferStatusParams{
		ID:     t.ID,
		Status: db.TransferStatus(t.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating stock transfer status", err)
	}

	t.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func stockTransferToDomain(row db.StockTransfer) *domain.StockTransfer {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	qty, _ := row.Quantity.Float64Value()

	return domain.RestoreStockTransfer(
		row.ID,
		row.ProductID,
		row.FromLocationID,
		row.ToLocationID,
		qty.Float64,
		domain.TransferStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func stockTransfersToDomain(rows []db.StockTransfer) []*domain.StockTransfer {
	transfers := make([]*domain.StockTransfer, len(rows))
	for i, row := range rows {
		transfers[i] = stockTransferToDomain(row)
	}

	return transfers
}
