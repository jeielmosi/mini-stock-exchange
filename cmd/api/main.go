package main

import (
	"context"
	"log"
	"net/http"

	"mini-stock-exchange/internal/controller"
	"mini-stock-exchange/internal/observability"
	"mini-stock-exchange/internal/repository"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	observability.InitLogger()
	ctx := context.Background()
	tp, err := observability.InitTracer(ctx, "stock-exchange-api")
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	db, err := repository.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	controller.RegisterRoutes(r, db)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
