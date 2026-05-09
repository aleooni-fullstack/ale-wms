package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	locationapp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/application"
	locationinfra "github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/infrastructure"
	locationhttp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/interfaces/http"
)

func RegisterLocation(r chi.Router, pool *pgxpool.Pool) {
	warehouseRepo := locationinfra.NewPostgresWarehouseRepository(pool)
	zoneRepo := locationinfra.NewPostgresZoneRepository(pool)
	locationRepo := locationinfra.NewPostgresLocationRepository(pool)

	warehouseService := locationapp.NewWarehouseService(warehouseRepo)
	zoneService := locationapp.NewZoneService(zoneRepo, warehouseRepo)
	locationService := locationapp.NewLocationService(locationRepo, zoneRepo)

	locationhttp.NewWarehouseHandler(warehouseService).RegisterRoutes(r)
	locationhttp.NewZoneHandler(zoneService).RegisterRoutes(r)
	locationhttp.NewLocationHandler(locationService).RegisterRoutes(r)
}
