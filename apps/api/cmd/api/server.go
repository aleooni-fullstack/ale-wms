package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleooni-fullstack/ale-wms/apps/api/cmd/api/bootstrap"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/config"
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

	bootstrap.RegisterCatalog(r, pool)
	bootstrap.RegisterLocation(r, pool)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
}
