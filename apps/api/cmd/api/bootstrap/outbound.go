package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/infrastructure"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/application"
	outboundinfra "github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/infrastructure"
	outboundhttp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/interfaces/http"
)

func RegisterOutbound(r chi.Router, pool *pgxpool.Pool) {
	orderRepo := outboundinfra.NewPostgresOrderRepository(pool)
	pickingRepo := outboundinfra.NewPostgresPickingRepository(pool)
	packingRepo := outboundinfra.NewPostgresPackingRepository(pool)
	shippingRepo := outboundinfra.NewPostgresShippingRepository(pool)

	balanceRepo := infrastructure.NewPostgresStockBalanceRepository(pool)
	movementRepo := infrastructure.NewPostgresStockMovementRepository(pool)

	orderService := application.NewOrderService(orderRepo)
	pickingService := application.NewPickingService(pickingRepo, orderRepo, balanceRepo, movementRepo)
	packingService := application.NewPackingService(packingRepo, pickingRepo, orderRepo)
	shippingService := application.NewShippingService(shippingRepo, packingRepo, orderRepo, balanceRepo, movementRepo)

	outboundhttp.NewOrderHandler(orderService).RegisterRoutes(r)
	outboundhttp.NewPickingHandler(pickingService).RegisterRoutes(r)
	outboundhttp.NewPackingHandler(packingService).RegisterRoutes(r)
	outboundhttp.NewShippingHandler(shippingService).RegisterRoutes(r)
}
