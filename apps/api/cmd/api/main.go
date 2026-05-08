package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("API running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
