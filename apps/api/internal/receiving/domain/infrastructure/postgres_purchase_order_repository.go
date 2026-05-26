package infrastructure

import (
	"context"
	"errors"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresPurchaseOrderRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresPurchaseOrderRepository(pool *pgxpool.Pool) *PostgresPurchaseOrderRepository {
	return &PostgresPurchaseOrderRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresPurchaseOrderRepository) Create(ctx context.Context, po *domain.PurchaseOrder) error {
	row, err := r.queries.CreatePurchaseOrder(ctx, db.CreatePurchaseOrderParams{
		ID:        po.ID,
		Reference: po.Reference,
		Supplier:  pgNullText(po.Supplier),
		Status:    db.PurchaseOrderStatus(po.Status),
		Note:      pgNullText(po.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating purchase order", err)
	}

	po.CreatedAt = row.CreatedAt.Time
	po.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPurchaseOrderRepository) FindByID(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	row, err := r.queries.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding purchase order by id", err)
	}

	return purchaseOrderToDomain(row), nil
}

func (r *PostgresPurchaseOrderRepository) FindByReference(ctx context.Context, reference string) (*domain.PurchaseOrder, error) {
	row, err := r.queries.GetPurchaseOrderByReference(ctx, reference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding purchase order by reference", err)
	}

	return purchaseOrderToDomain(row), nil
}

func (r *PostgresPurchaseOrderRepository) FindAll(ctx context.Context, limit, offset int32) ([]*domain.PurchaseOrder, error) {
	rows, err := r.queries.ListPurchaseOrders(ctx, db.ListPurchaseOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing purchase orders", err)
	}

	pos := make([]*domain.PurchaseOrder, len(rows))
	for i, row := range rows {
		pos[i] = purchaseOrderToDomain(row)
	}

	return pos, nil
}

func (r *PostgresPurchaseOrderRepository) UpdateStatus(ctx context.Context, po *domain.PurchaseOrder) error {
	row, err := r.queries.UpdatePurchaseOrderStatus(ctx, db.UpdatePurchaseOrderStatusParams{
		ID:     po.ID,
		Status: db.PurchaseOrderStatus(po.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating purchase order status", err)
	}

	po.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresPurchaseOrderRepository) AddItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	row, err := r.queries.CreatePurchaseOrderItem(ctx, db.CreatePurchaseOrderItemParams{
		ID:              item.ID,
		PurchaseOrderID: item.PurchaseOrderID,
		ProductID:       item.ProductID,
		Quantity:        pgNumeric(item.Quantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding purchase order item", err)
	}

	item.CreatedAt = row.CreatedAt.Time

	return nil
}

func (r *PostgresPurchaseOrderRepository) FindAllItems(ctx context.Context, purchaseOrderID string) ([]*domain.PurchaseOrderItem, error) {
	rows, err := r.queries.ListPurchaseOrderItems(ctx, purchaseOrderID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing purchase order items", err)
	}

	items := make([]*domain.PurchaseOrderItem, len(rows))
	for i, row := range rows {
		items[i] = purchaseOrderItemToDomain(row)
	}

	return items, nil
}

func purchaseOrderToDomain(row db.PurchaseOrder) *domain.PurchaseOrder {
	supplier := ""
	if row.Supplier.Valid {
		supplier = row.Supplier.String
	}

	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestorePurchaseOrder(
		row.ID,
		row.Reference,
		supplier,
		domain.PurchaseOrderStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func purchaseOrderItemToDomain(row db.PurchaseOrderItem) *domain.PurchaseOrderItem {
	qty, _ := row.Quantity.Float64Value()

	return domain.RestorePurchaseOrderItem(
		row.ID,
		row.PurchaseOrderID,
		row.ProductID,
		qty.Float64,
		row.CreatedAt.Time,
	)
}

func pgNullText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func pgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(strconv.FormatFloat(f, 'f', -1, 64))
	return n
}
