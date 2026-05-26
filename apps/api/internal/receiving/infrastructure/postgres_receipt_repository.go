package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresReceiptRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresReceiptRepository(pool *pgxpool.Pool) *PostgresReceiptRepository {
	return &PostgresReceiptRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresReceiptRepository) Create(ctx context.Context, receipt *domain.Receipt) error {
	row, err := r.queries.CreateReceipt(ctx, db.CreateReceiptParams{
		ID:              receipt.ID,
		PurchaseOrderID: receipt.PurchaseOrderID,
		Status:          db.ReceiptStatus(receipt.Status),
		Note:            pgNullText(receipt.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating receipt", err)
	}

	receipt.CreatedAt = row.CreatedAt.Time
	receipt.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresReceiptRepository) FindByID(ctx context.Context, id string) (*domain.Receipt, error) {
	row, err := r.queries.GetReceiptByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding receipt by id", err)
	}

	return receiptToDomain(row), nil
}

func (r *PostgresReceiptRepository) FindByPurchaseOrderID(ctx context.Context, purchaseOrderID string) (*domain.Receipt, error) {
	row, err := r.queries.GetReceiptByPurchaseOrderID(ctx, purchaseOrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding receipt by purchase order id", err)
	}

	return receiptToDomain(row), nil
}

func (r *PostgresReceiptRepository) UpdateStatus(ctx context.Context, receipt *domain.Receipt) error {
	row, err := r.queries.UpdateReceiptStatus(ctx, db.UpdateReceiptStatusParams{
		ID:     receipt.ID,
		Status: db.ReceiptStatus(receipt.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating receipt status", err)
	}

	receipt.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresReceiptRepository) AddItem(ctx context.Context, item *domain.ReceiptItem) error {
	row, err := r.queries.CreateReceiptItem(ctx, db.CreateReceiptItemParams{
		ID:               item.ID,
		ReceiptID:        item.ReceiptID,
		ProductID:        item.ProductID,
		ExpectedQuantity: pgNumeric(item.ExpectedQuantity),
		ReceivedQuantity: pgNumericPtr(item.ReceivedQuantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding receipt item", err)
	}

	item.CreatedAt = row.CreatedAt.Time
	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresReceiptRepository) FindAllItems(ctx context.Context, receiptID string) ([]*domain.ReceiptItem, error) {
	rows, err := r.queries.ListReceiptItems(ctx, receiptID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing receipt items", err)
	}

	items := make([]*domain.ReceiptItem, len(rows))
	for i, row := range rows {
		items[i] = receiptItemToDomain(row)
	}

	return items, nil
}

func (r *PostgresReceiptRepository) UpdateItemReceivedQuantity(ctx context.Context, item *domain.ReceiptItem) error {
	row, err := r.queries.UpdateReceiptItemReceivedQuantity(ctx, db.UpdateReceiptItemReceivedQuantityParams{
		ReceiptID:        item.ReceiptID,
		ProductID:        item.ProductID,
		ReceivedQuantity: pgNumericPtr(item.ReceivedQuantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating receipt item received quantity", err)
	}

	item.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func receiptToDomain(row db.Receipt) *domain.Receipt {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestoreReceipt(
		row.ID,
		row.PurchaseOrderID,
		domain.ReceiptStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func receiptItemToDomain(row db.ReceiptItem) *domain.ReceiptItem {
	expectedQty, _ := row.ExpectedQuantity.Float64Value()

	var receivedQty *float64
	if row.ReceivedQuantity.Valid {
		v, _ := row.ReceivedQuantity.Float64Value()
		receivedQty = &v.Float64
	}

	return domain.RestoreReceiptItem(
		row.ID,
		row.ReceiptID,
		row.ProductID,
		expectedQty.Float64,
		receivedQty,
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
