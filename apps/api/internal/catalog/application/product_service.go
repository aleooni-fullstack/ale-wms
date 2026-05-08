package application

import (
	"context"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/domain"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductInput struct {
	SKU         string
	Name        string
	Description string
	Unit        string
}

type UpdateProductInput struct {
	SKU         string
	Name        string
	Description string
	Unit        string
}

func (s *ProductService) Create(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	_, err := s.repo.FindBySKU(ctx, input.SKU)
	if err == nil {
		return nil, dderr.New("SKU_ALREADY_EXISTS", "a product with this SKU already exists", nil)
	}
	if err != dderr.ErrNotFound {
		return nil, err
	}

	product, err := domain.NewProduct(input.SKU, input.Name, input.Description, input.Unit)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProductService) Update(ctx context.Context, id string, input UpdateProductInput) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product.SKU != input.SKU {
		_, err := s.repo.FindBySKU(ctx, input.SKU)
		if err == nil {
			return nil, dderr.New("SKU_ALREADY_EXISTS", "a product with this SKU already exists", nil)
		}
		if err != dderr.ErrNotFound {
			return nil, err
		}
	}

	if err := product.Update(input.SKU, input.Name, input.Description, input.Unit); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Deactivate(ctx context.Context, id string) error {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	product.Deactivate()

	return s.repo.Delete(ctx, product)
}
