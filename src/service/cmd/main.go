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
	"katharos/service/internal/app"
	"katharos/service/internal/config"
	"sharedmodule/logging"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	loggerConfig := &logging.LoggerConfig{
		Level:          logging.InfoLevel,
		FileName:       "/tmp/katharos-service.log",
		LoggerName:     "katharos-service",
		ComponentName:  "main",
		ServiceName:    "katharos-service",
		MaxAge:         30,
		MaxBackups:     10,
		MaxSize:        100,
		IsLogRotatable: false,
	}

	logger, err := logging.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Create application instance
	application := app.NewApplication(cfg, logger)

	// Initialize handlers (without user service for now)
	handler := api.NewHandler()

	// Setup HTTP mux
	mux := setupRouter(handler)

	// Start server
	startServer(mux, cfg, application)
}

func setupRouter(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup routes
	handler.SetupRoutes(mux)

	return mux
}

func startServer(mux *http.ServeMux, cfg *config.Config, application *app.Application) {
	logger := application.Logger()

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests a 30-second deadline to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	// Shutdown the application
	if err := application.Shutdown(); err != nil {
		logger.Errorf("Application shutdown error: %v", err)
	}

	logger.Info("Server exited")
}
