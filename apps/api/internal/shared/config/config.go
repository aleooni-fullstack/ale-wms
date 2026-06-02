package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL   string
	Port          string
	KeycloakURL   string
	KeycloakRealm string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	keycloakURL := os.Getenv("KEYCLOAK_URL")
	if keycloakURL == "" {
		log.Fatal("KEYCLOAK_URL environment variable is required")
	}

	keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	if keycloakRealm == "" {
		log.Fatal("KEYCLOAK_REALM environment variable is required")
	}

	return &Config{
		DatabaseURL:   dbURL,
		Port:          port,
		KeycloakURL:   keycloakURL,
		KeycloakRealm: keycloakRealm,
	}
}
