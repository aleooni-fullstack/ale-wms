package infrastructure

import (
	"context"
	"errors"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresStockBalanceRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresStockBalanceRepository(pool *pgxpool.Pool) *PostgresStockBalanceRepository {
	return &PostgresStockBalanceRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresStockBalanceRepository) FindByProductAndLocation(ctx context.Context, productID, locationID string) (*domain.StockBalance, error) {
	row, err := r.queries.GetStockBalance(ctx, db.GetStockBalanceParams{
		ProductID:  productID,
		LocationID: locationID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding stock balance", err)
	}

	return stockBalanceToDomain(row), nil
}

func (r *PostgresStockBalanceRepository) FindAllByProduct(ctx context.Context, productID string) ([]*domain.StockBalance, error) {
	rows, err := r.queries.ListStockBalancesByProduct(ctx, productID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock balances by product", err)
	}

	balances := make([]*domain.StockBalance, len(rows))
	for i, row := range rows {
		balances[i] = stockBalanceToDomain(row)
	}

	return balances, nil
}

func (r *PostgresStockBalanceRepository) FindAllByLocation(ctx context.Context, locationID string) ([]*domain.StockBalance, error) {
	rows, err := r.queries.ListStockBalancesByLocation(ctx, locationID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing stock balances by location", err)
	}

	balances := make([]*domain.StockBalance, len(rows))
	for i, row := range rows {
		balances[i] = stockBalanceToDomain(row)
	}

	return balances, nil
}

func (r *PostgresStockBalanceRepository) Upsert(ctx context.Context, b *domain.StockBalance) error {
	_, err := r.queries.UpsertStockBalance(ctx, db.UpsertStockBalanceParams{
		ID:               b.ID,
		ProductID:        b.ProductID,
		LocationID:       b.LocationID,
		Quantity:         pgNumeric(b.Quantity),
		ReservedQuantity: pgNumeric(b.ReservedQuantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error upserting stock balance", err)
	}

	return nil
}

func stockBalanceToDomain(row db.StockBalance) *domain.StockBalance {
	qty, _ := row.Quantity.Float64Value()
	reservedQty, _ := row.ReservedQuantity.Float64Value()

	return domain.RestoreStockBalance(
		row.ID,
		row.ProductID,
		row.LocationID,
		qty.Float64,
		reservedQty.Float64,
		row.UpdatedAt.Time,
	)
}

func pgNullText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func pgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(strconv.FormatFloat(f, 'f', -1, 64)); err != nil {
		return pgtype.Numeric{}
	}
	return n
}
