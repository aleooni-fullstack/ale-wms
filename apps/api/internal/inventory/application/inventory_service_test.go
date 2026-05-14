package application

import (
	"context"
	"errors"
	"testing"

	dderr "github.com/aleodoni/go-ddd/errors"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type mockMovementRepo struct {
	createFn            func(ctx context.Context, movement *domain.StockMovement) error
	findByIDFn          func(ctx context.Context, id string) (*domain.StockMovement, error)
	findAllByProductFn  func(ctx context.Context, productID string, limit, offset int32) ([]*domain.StockMovement, error)
	findAllByLocationFn func(ctx context.Context, locationID string, limit, offset int32) ([]*domain.StockMovement, error)
}

func (m *mockMovementRepo) Create(
	ctx context.Context,
	movement *domain.StockMovement,
) error {
	if m.createFn != nil {
		return m.createFn(ctx, movement)
	}

	return nil
}

func (m *mockMovementRepo) FindByID(
	ctx context.Context,
	id string,
) (*domain.StockMovement, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}

	return nil, dderr.ErrNotFound
}

func (m *mockMovementRepo) FindAllByProduct(
	ctx context.Context,
	productID string,
	limit,
	offset int32,
) ([]*domain.StockMovement, error) {
	if m.findAllByProductFn != nil {
		return m.findAllByProductFn(ctx, productID, limit, offset)
	}

	return nil, nil
}

func (m *mockMovementRepo) FindAllByLocation(
	ctx context.Context,
	locationID string,
	limit,
	offset int32,
) ([]*domain.StockMovement, error) {
	if m.findAllByLocationFn != nil {
		return m.findAllByLocationFn(ctx, locationID, limit, offset)
	}

	return nil, nil
}

type mockBalanceRepo struct {
	findByProductAndLocationFn func(ctx context.Context, productID, locationID string) (*domain.StockBalance, error)
	findAllByProductFn         func(ctx context.Context, productID string) ([]*domain.StockBalance, error)
	findAllByLocationFn        func(ctx context.Context, locationID string) ([]*domain.StockBalance, error)
	upsertFn                   func(ctx context.Context, balance *domain.StockBalance) error
}

func (m *mockBalanceRepo) FindByProductAndLocation(
	ctx context.Context,
	productID,
	locationID string,
) (*domain.StockBalance, error) {
	if m.findByProductAndLocationFn != nil {
		return m.findByProductAndLocationFn(ctx, productID, locationID)
	}

	return nil, dderr.ErrNotFound
}

func (m *mockBalanceRepo) FindAllByProduct(
	ctx context.Context,
	productID string,
) ([]*domain.StockBalance, error) {
	if m.findAllByProductFn != nil {
		return m.findAllByProductFn(ctx, productID)
	}

	return nil, nil
}

func (m *mockBalanceRepo) FindAllByLocation(
	ctx context.Context,
	locationID string,
) ([]*domain.StockBalance, error) {
	if m.findAllByLocationFn != nil {
		return m.findAllByLocationFn(ctx, locationID)
	}

	return nil, nil
}

func (m *mockBalanceRepo) Upsert(
	ctx context.Context,
	balance *domain.StockBalance,
) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, balance)
	}

	return nil
}

func TestInventoryService_RegisterMovement(t *testing.T) {
	t.Run("should create IN movement and update balance", func(t *testing.T) {
		movementRepo := &mockMovementRepo{}

		var updatedBalance *domain.StockBalance

		balanceRepo := &mockBalanceRepo{
			findByProductAndLocationFn: func(
				ctx context.Context,
				productID,
				locationID string,
			) (*domain.StockBalance, error) {
				return domain.NewStockBalance(
					productID,
					locationID,
					10,
				), nil
			},
			upsertFn: func(
				ctx context.Context,
				balance *domain.StockBalance,
			) error {
				updatedBalance = balance
				return nil
			},
		}

		service := NewInventoryService(
			movementRepo,
			balanceRepo,
		)

		movement, err := service.RegisterMovement(
			context.Background(),
			RegisterMovementInput{
				ProductID:  "product-1",
				LocationID: "location-1",
				Type:       string(domain.MovementTypeIn),
				Quantity:   5,
				Note:       "restock",
			},
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if movement == nil {
			t.Fatal("expected movement to be created")
		}

		if updatedBalance == nil {
			t.Fatal("expected balance to be updated")
		}

		if updatedBalance.Quantity != 15 {
			t.Fatalf(
				"expected balance quantity to be 15, got %v",
				updatedBalance.Quantity,
			)
		}
	})

	t.Run("should reject OUT movement with insufficient stock", func(t *testing.T) {
		movementRepo := &mockMovementRepo{}

		balanceRepo := &mockBalanceRepo{
			findByProductAndLocationFn: func(
				ctx context.Context,
				productID,
				locationID string,
			) (*domain.StockBalance, error) {
				return domain.NewStockBalance(
					productID,
					locationID,
					3,
				), nil
			},
		}

		service := NewInventoryService(
			movementRepo,
			balanceRepo,
		)

		movement, err := service.RegisterMovement(
			context.Background(),
			RegisterMovementInput{
				ProductID:  "product-1",
				LocationID: "location-1",
				Type:       string(domain.MovementTypeOut),
				Quantity:   10,
			},
		)

		if err == nil {
			t.Fatal("expected error")
		}

		if movement != nil {
			t.Fatal("expected movement to be nil")
		}

		var domainErr *dderr.DomainError

		if !errors.As(err, &domainErr) {
			t.Fatalf("expected domain error, got %v", err)
		}

		if domainErr.Code != "INSUFFICIENT_STOCK" {
			t.Fatalf(
				"expected INSUFFICIENT_STOCK, got %s",
				domainErr.Code,
			)
		}
	})

	t.Run("should create balance when balance does not exist", func(t *testing.T) {
		movementRepo := &mockMovementRepo{}

		var updatedBalance *domain.StockBalance

		balanceRepo := &mockBalanceRepo{
			findByProductAndLocationFn: func(
				ctx context.Context,
				productID,
				locationID string,
			) (*domain.StockBalance, error) {
				return nil, dderr.ErrNotFound
			},
			upsertFn: func(
				ctx context.Context,
				balance *domain.StockBalance,
			) error {
				updatedBalance = balance
				return nil
			},
		}

		service := NewInventoryService(
			movementRepo,
			balanceRepo,
		)

		_, err := service.RegisterMovement(
			context.Background(),
			RegisterMovementInput{
				ProductID:  "product-1",
				LocationID: "location-1",
				Type:       string(domain.MovementTypeIn),
				Quantity:   7,
			},
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if updatedBalance == nil {
			t.Fatal("expected balance to be created")
		}

		if updatedBalance.Quantity != 7 {
			t.Fatalf(
				"expected balance quantity to be 7, got %v",
				updatedBalance.Quantity,
			)
		}
	})

	t.Run("should return repository error when create movement fails", func(t *testing.T) {
		expectedErr := errors.New("database error")

		movementRepo := &mockMovementRepo{
			createFn: func(
				ctx context.Context,
				movement *domain.StockMovement,
			) error {
				return expectedErr
			},
		}

		service := NewInventoryService(
			movementRepo,
			&mockBalanceRepo{},
		)

		_, err := service.RegisterMovement(
			context.Background(),
			RegisterMovementInput{
				ProductID:  "product-1",
				LocationID: "location-1",
				Type:       string(domain.MovementTypeIn),
				Quantity:   5,
			},
		)

		if !errors.Is(err, expectedErr) {
			t.Fatalf(
				"expected %v, got %v",
				expectedErr,
				err,
			)
		}
	})
}

func TestInventoryService_GetBalance(t *testing.T) {
	expectedBalance := domain.NewStockBalance(
		"product-1",
		"location-1",
		25,
	)

	service := NewInventoryService(
		&mockMovementRepo{},
		&mockBalanceRepo{
			findByProductAndLocationFn: func(
				ctx context.Context,
				productID,
				locationID string,
			) (*domain.StockBalance, error) {
				return expectedBalance, nil
			},
		},
	)

	balance, err := service.GetBalance(
		context.Background(),
		"product-1",
		"location-1",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if balance != expectedBalance {
		t.Fatal("expected returned balance to match")
	}
}

func TestInventoryService_ListMovements(t *testing.T) {
	t.Run("should list movements by product", func(t *testing.T) {
		expectedMovements := []*domain.StockMovement{}

		movementRepo := &mockMovementRepo{
			findAllByProductFn: func(
				ctx context.Context,
				productID string,
				limit,
				offset int32,
			) ([]*domain.StockMovement, error) {
				if productID != "product-1" {
					t.Fatalf(
						"expected product-1, got %s",
						productID,
					)
				}

				if limit != 20 {
					t.Fatalf(
						"expected limit 20, got %d",
						limit,
					)
				}

				if offset != 0 {
					t.Fatalf(
						"expected offset 0, got %d",
						offset,
					)
				}

				return expectedMovements, nil
			},
		}

		service := NewInventoryService(
			movementRepo,
			&mockBalanceRepo{},
		)

		result, err := service.ListMovements(
			context.Background(),
			ListMovementsInput{
				ProductID: "product-1",
			},
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("expected result")
		}

		if result.Page != 1 {
			t.Fatalf(
				"expected page 1, got %d",
				result.Page,
			)
		}

		if result.PerPage != 20 {
			t.Fatalf(
				"expected per_page 20, got %d",
				result.PerPage,
			)
		}
	})

	t.Run("should return error when no filter is provided", func(t *testing.T) {
		service := NewInventoryService(
			&mockMovementRepo{},
			&mockBalanceRepo{},
		)

		_, err := service.ListMovements(
			context.Background(),
			ListMovementsInput{},
		)

		if err == nil {
			t.Fatal("expected error")
		}

		var domainErr *dderr.DomainError

		if !errors.As(err, &domainErr) {
			t.Fatalf(
				"expected domain error, got %v",
				err,
			)
		}

		if domainErr.Code != "INVALID_FILTER" {
			t.Fatalf(
				"expected INVALID_FILTER, got %s",
				domainErr.Code,
			)
		}
	})
}
