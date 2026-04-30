// Package main is the entry point for pricing-service.
package main

import (
	"log"
	"net/http"
	"time"

	"furab-backend/services/pricing-service/internal/client"
	"furab-backend/services/pricing-service/internal/handler"
	"furab-backend/services/pricing-service/internal/repository"
	"furab-backend/services/pricing-service/internal/service"
	"furab-backend/shared/config"
	sharedlogger "furab-backend/shared/logger"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load("pricing-service")
	logger := sharedlogger.New(cfg.ServiceName, cfg.Environment)

	logger.Info("starting pricing-service", "port", cfg.ServerPort)

	// Compose dependencies
	repo := repository.NewInMemoryPriceRepository()
	orderClient := client.NewDummyOrderClient()
	locationClient := client.NewDummyLocationClient()
	priceService := service.NewPriceService(repo, orderClient, locationClient)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// Register routes
	h := handler.NewPriceHandler(priceService)
	h.RegisterRoutes(r)

	// Start server
	logger.Info("server listening", "address", cfg.ServerAddr())
	if err := http.ListenAndServe(cfg.ServerAddr(), r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
