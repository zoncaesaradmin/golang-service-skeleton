package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"katharos/service/internal/api"
	"katharos/service/internal/config"
	"katharos/service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize services
	userService := service.NewUserService()

	// Initialize handlers
	handler := api.NewHandler(userService)

	// Setup HTTP mux
	mux := setupRouter(handler)

	// Start server
	startServer(mux, cfg)
}

func setupRouter(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup routes
	handler.SetupRoutes(mux)

	return mux
}

func startServer(mux *http.ServeMux, cfg *config.Config) {
	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a 30-second deadline to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
