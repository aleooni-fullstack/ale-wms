package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/cmd/api/bootstrap"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/config"
	wmsmiddleware "github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/middleware"
)

func newServer(cfg *config.Config, pool *pgxpool.Pool) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	auth := wmsmiddleware.NewAuthMiddleware(cfg.KeycloakURL, cfg.KeycloakRealm)

	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)

		// qualquer usuário autenticado
		bootstrap.RegisterCatalog(r, pool)
		bootstrap.RegisterLocation(r, pool)

		// operator ou admin
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireRole("ADMIN", "OPERATOR"))
			bootstrap.RegisterInventory(r, pool)
			bootstrap.RegisterOutbound(r, pool)
			bootstrap.RegisterReceiving(r, pool)
		})
	})

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
}
