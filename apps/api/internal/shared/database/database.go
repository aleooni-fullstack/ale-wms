package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(databaseURL string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
