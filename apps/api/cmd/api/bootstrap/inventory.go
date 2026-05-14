package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/infrastructure"
	inventoryhttp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/interfaces/http"
)

func RegisterInventory(r chi.Router, pool *pgxpool.Pool) {
	movementRepo := infrastructure.NewPostgresStockMovementRepository(pool)
	balanceRepo := infrastructure.NewPostgresStockBalanceRepository(pool)
	transferRepo := infrastructure.NewPostgresStockTransferRepository(pool)

	inventoryService := application.NewInventoryService(movementRepo, balanceRepo)
	transferService := application.NewTransferService(transferRepo, balanceRepo, movementRepo)

	inventoryhttp.NewInventoryHandler(inventoryService).RegisterRoutes(r)
	inventoryhttp.NewTransferHandler(transferService).RegisterRoutes(r)
}
