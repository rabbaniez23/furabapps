// Package main is the entry point for promo-service.
package main

import (
	"log"
	"net/http"
	"time"

	"furab-backend/services/promo-service/internal/client"
	"furab-backend/services/promo-service/internal/handler"
	"furab-backend/services/promo-service/internal/repository"
	"furab-backend/services/promo-service/internal/service"
	"furab-backend/shared/config"
	sharedlogger "furab-backend/shared/logger"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load("promo-service")
	logger := sharedlogger.New(cfg.ServiceName, cfg.Environment)

	logger.Info("starting promo-service", "port", cfg.ServerPort)

	repo := repository.NewInMemoryPromoRepository()
	orderClient := client.NewDummyOrderClient()
	userClient := client.NewDummyUserClient()
	promoService := service.NewPromoService(repo, orderClient, userClient)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// Register routes
	h := handler.NewPromoHandler(promoService)
	h.RegisterRoutes(r)

	// Start server
	logger.Info("server listening", "address", cfg.ServerAddr())
	if err := http.ListenAndServe(cfg.ServerAddr(), r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
