package infrastructure

import (
	"context"
	"errors"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/domain"
	db "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database/generated"
)

type PostgresProductRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	row, err := r.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding product by id", err)
	}

	return toDomain(row), nil
}

func (r *PostgresProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	row, err := r.queries.GetProductBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dderr.ErrNotFound
		}
		return nil, dderr.New("DB_ERROR", "error finding product by sku", err)
	}

	return toDomain(row), nil
}

func (r *PostgresProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.queries.ListProducts(ctx)
	if err != nil {
		return nil, dderr.New("DB_ERROR", "error listing products", err)
	}

	products := make([]*domain.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomain(row)
	}

	return products, nil
}

func (r *PostgresProductRepository) Create(ctx context.Context, p *domain.Product) error {
	_, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		ID:          p.ID,
		Sku:         p.SKU,
		Name:        p.Name,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
		Unit:        p.Unit,
		Active:      p.Active,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error creating product", err)
	}

	return nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, p *domain.Product) error {
	_, err := r.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          p.ID,
		Sku:         p.SKU,
		Name:        p.Name,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
		Unit:        p.Unit,
	})
	if err != nil {
		return dderr.New("DB_ERROR", "error updating product", err)
	}

	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, p *domain.Product) error {
	err := r.queries.DeactivateProduct(ctx, p.ID)
	if err != nil {
		return dderr.New("DB_ERROR", "error deactivating product", err)
	}

	return nil
}

func toDomain(row db.Product) *domain.Product {
	desc := ""
	if row.Description.Valid {
		desc = row.Description.String
	}

	return domain.RestoreProduct(
		row.ID,
		row.Sku,
		row.Name,
		desc,
		row.Unit,
		row.Active,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}
