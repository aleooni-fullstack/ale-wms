package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/infrastructure"
	cataloghttp "github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/interfaces/http"
)

func RegisterCatalog(r chi.Router, pool *pgxpool.Pool) {
	repo := infrastructure.NewPostgresProductRepository(pool)
	service := application.NewProductService(repo)
	handler := cataloghttp.NewProductHandler(service)
	handler.RegisterRoutes(r)
}
