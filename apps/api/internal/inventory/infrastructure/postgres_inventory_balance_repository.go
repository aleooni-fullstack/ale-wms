package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresInventoryBalanceRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresInventoryBalanceRepository(pool *pgxpool.Pool) *PostgresInventoryBalanceRepository {
	return &PostgresInventoryBalanceRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresInventoryBalanceRepository) Create(ctx context.Context, b *domain.InventoryBalance) error {
	row, err := r.queries.CreateInventoryBalance(ctx, db.CreateInventoryBalanceParams{
		ID:         b.ID,
		LocationID: b.LocationID,
		Status:     db.InventoryBalanceStatus(b.Status),
		Note:       pgNullText(b.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating inventory balance", err)
	}

	b.CreatedAt = row.CreatedAt.Time
	b.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresInventoryBalanceRepository) FindByID(ctx context.Context, id string) (*domain.InventoryBalance, error) {
	row, err := r.queries.GetInventoryBalanceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding inventory balance by id", err)
	}

	return inventoryBalanceToDomain(row), nil
}

func (r *PostgresInventoryBalanceRepository) FindAllByLocation(ctx context.Context, locationID string, limit, offset int32) ([]*domain.InventoryBalance, error) {
	rows, err := r.queries.ListInventoryBalancesByLocation(ctx, db.ListInventoryBalancesByLocationParams{
		LocationID: locationID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing inventory balances by location", err)
	}

	balances := make([]*domain.InventoryBalance, len(rows))
	for i, row := range rows {
		balances[i] = inventoryBalanceToDomain(row)
	}

	return balances, nil
}

func (r *PostgresInventoryBalanceRepository) UpdateStatus(ctx context.Context, b *domain.InventoryBalance) error {
	row, err := r.queries.UpdateInventoryBalanceStatus(ctx, db.UpdateInventoryBalanceStatusParams{
		ID:     b.ID,
		Status: db.InventoryBalanceStatus(b.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating inventory balance status", err)
	}

	b.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresInventoryBalanceRepository) AddItem(ctx context.Context, item *domain.InventoryBalanceItem) error {
	row, err := r.queries.CreateInventoryBalanceItem(ctx, db.CreateInventoryBalanceItemParams{
		ID:                 item.ID,
		InventoryBalanceID: item.InventoryBalanceID,
		ProductID:          item.ProductID,
		SystemQuantity:     pgNumeric(item.SystemQuantity),
		CountedQuantity:    pgNumericPtr(item.CountedQuantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding inventory balance item", err)
	}

	item.CreatedAt = row.CreatedAt.Time
	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresInventoryBalanceRepository) FindItem(ctx context.Context, inventoryBalanceID, productID string) (*domain.InventoryBalanceItem, error) {
	row, err := r.queries.GetInventoryBalanceItem(ctx, db.GetInventoryBalanceItemParams{
		InventoryBalanceID: inventoryBalanceID,
		ProductID:          productID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding inventory balance item", err)
	}

	return inventoryBalanceItemToDomain(row), nil
}

func (r *PostgresInventoryBalanceRepository) FindAllItems(ctx context.Context, inventoryBalanceID string) ([]*domain.InventoryBalanceItem, error) {
	rows, err := r.queries.ListInventoryBalanceItems(ctx, inventoryBalanceID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing inventory balance items", err)
	}

	items := make([]*domain.InventoryBalanceItem, len(rows))
	for i, row := range rows {
		items[i] = inventoryBalanceItemToDomain(row)
	}

	return items, nil
}

func (r *PostgresInventoryBalanceRepository) UpdateItemCountedQuantity(ctx context.Context, item *domain.InventoryBalanceItem) error {
	row, err := r.queries.UpdateInventoryBalanceItemCountedQuantity(ctx, db.UpdateInventoryBalanceItemCountedQuantityParams{
		InventoryBalanceID: item.InventoryBalanceID,
		ProductID:          item.ProductID,
		CountedQuantity:    pgNumericPtr(item.CountedQuantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating inventory balance item counted quantity", err)
	}

	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func inventoryBalanceToDomain(row db.InventoryBalance) *domain.InventoryBalance {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestoreInventoryBalance(
		row.ID,
		row.LocationID,
		domain.InventoryBalanceStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func inventoryBalanceItemToDomain(row db.InventoryBalanceItem) *domain.InventoryBalanceItem {
	sysQty, _ := row.SystemQuantity.Float64Value()

	var countedQty *float64
	if row.CountedQuantity.Valid {
		v, _ := row.CountedQuantity.Float64Value()
		countedQty = &v.Float64
	}

	return domain.RestoreInventoryBalanceItem(
		row.ID,
		row.InventoryBalanceID,
		row.ProductID,
		sysQty.Float64,
		countedQty,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func pgNumericPtr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{}
	}
	return pgNumeric(*f)
}
