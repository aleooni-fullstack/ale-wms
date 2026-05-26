package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	inventoryinfra "github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/infrastructure"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/application"
	receivinginfra "github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/infrastructure"
	receivinghttp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/interfaces/http"
)

func RegisterReceiving(r chi.Router, pool *pgxpool.Pool) {
	poRepo := receivinginfra.NewPostgresPurchaseOrderRepository(pool)
	receiptRepo := receivinginfra.NewPostgresReceiptRepository(pool)
	putAwayRepo := receivinginfra.NewPostgresPutAwayRepository(pool)

	balanceRepo := inventoryinfra.NewPostgresStockBalanceRepository(pool)
	movementRepo := inventoryinfra.NewPostgresStockMovementRepository(pool)

	poService := application.NewPurchaseOrderService(poRepo)
	receiptService := application.NewReceiptService(receiptRepo, poRepo)
	putAwayService := application.NewPutAwayService(putAwayRepo, receiptRepo, poRepo, balanceRepo, movementRepo)

	receivinghttp.NewPurchaseOrderHandler(poService).RegisterRoutes(r)
	receivinghttp.NewReceiptHandler(receiptService).RegisterRoutes(r)
	receivinghttp.NewPutAwayHandler(putAwayService).RegisterRoutes(r)
}
