package infrastructure

import (
	"context"
	"errors"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresOrderRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	row, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		ID:        o.ID,
		Reference: o.Reference,
		Status:    db.OrderStatus(o.Status),
		Note:      pgNullText(o.Note),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating order", err)
	}

	o.CreatedAt = row.CreatedAt.Time
	o.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding order by id", err)
	}

	return orderToDomain(row), nil
}

func (r *PostgresOrderRepository) FindByReference(ctx context.Context, reference string) (*domain.Order, error) {
	row, err := r.queries.GetOrderByReference(ctx, reference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding order by reference", err)
	}

	return orderToDomain(row), nil
}

func (r *PostgresOrderRepository) FindAll(ctx context.Context, limit, offset int32) ([]*domain.Order, error) {
	rows, err := r.queries.ListOrders(ctx, db.ListOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing orders", err)
	}

	orders := make([]*domain.Order, len(rows))
	for i, row := range rows {
		orders[i] = orderToDomain(row)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, o *domain.Order) error {
	row, err := r.queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     o.ID,
		Status: db.OrderStatus(o.Status),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating order status", err)
	}

	o.UpdatedAt = row.UpdatedAt.Time

	return nil
}

func (r *PostgresOrderRepository) AddItem(ctx context.Context, item *domain.OrderItem) error {
	row, err := r.queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
		ID:         item.ID,
		OrderID:    item.OrderID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   pgNumeric(item.Quantity),
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error adding order item", err)
	}

	item.CreatedAt = row.CreatedAt.Time

	return nil
}

func (r *PostgresOrderRepository) FindAllItems(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	rows, err := r.queries.ListOrderItems(ctx, orderID)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing order items", err)
	}

	items := make([]*domain.OrderItem, len(rows))
	for i, row := range rows {
		items[i] = orderItemToDomain(row)
	}

	return items, nil
}

func orderToDomain(row db.Order) *domain.Order {
	note := ""
	if row.Note.Valid {
		note = row.Note.String
	}

	return domain.RestoreOrder(
		row.ID,
		row.Reference,
		domain.OrderStatus(row.Status),
		note,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func orderItemToDomain(row db.OrderItem) *domain.OrderItem {
	qty, _ := row.Quantity.Float64Value()

	return domain.RestoreOrderItem(
		row.ID,
		row.OrderID,
		row.ProductID,
		row.LocationID,
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
