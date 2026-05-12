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

	service := application.NewInventoryService(movementRepo, balanceRepo)

	inventoryhttp.NewInventoryHandler(service).RegisterRoutes(r)
}
