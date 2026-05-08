package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() *pgxpool.Pool {
	db, err := pgxpool.New(
		context.Background(),
		"postgres://wms:wms@localhost:5432/wms",
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
