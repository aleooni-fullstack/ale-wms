package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/config"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/shared/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	pool := database.New(cfg.DatabaseURL)
	defer pool.Close()

	server := newServer(cfg, pool)

	log.Printf("API running on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
