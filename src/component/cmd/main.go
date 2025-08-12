package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"compmodule/internal/api"
	"compmodule/internal/app"
	"compmodule/internal/config"
	"sharedmodule/logging"
)

func main() {
	cfg := loadConfig()

	logger := initLogger(cfg)
	defer logger.Close()

	// Create application instance
	application := app.NewApplication(cfg, logger)

	// Initialize handlers and setup HTTP mux
	handler := api.NewHandler()
	mux := setupRouter(handler)

	// Start server
	startServer(logger, mux, cfg, application)
}

func setupRouter(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup routes
	handler.SetupRoutes(mux)

	return mux
}

func startServer(logger logging.Logger, mux *http.ServeMux, cfg *config.Config, application *app.Application) {

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting http server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down application ...")

	// Give outstanding requests a 10-second deadline to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func loadConfig() *config.Config {
	// Load configuration using absolute paths based on HOME_DIR environment variable
	homeDir := os.Getenv("HOME_DIR")
	if homeDir == "" {
		log.Fatal("HOME_DIR environment variable is required and must point to the repository root")
	}

	// Load configuration from the centralized config file
	configPath := filepath.Join(homeDir, "conf", "config.yaml")

	cfg, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configPath, err)
	}
	log.Printf("Loaded configuration from: %s", configPath)

	// If no config file was found, use defaults
	if cfg == nil {
		cfg = config.LoadConfig()
		log.Printf("No configuration file found, using environment variables and defaults")
	}

	return cfg
}

func initLogger(cfg *config.Config) logging.Logger {
	// Determine log file path from configuration
	logFilePath := cfg.Logging.FilePath
	if logFilePath == "" {
		logFilePath = "/tmp/katharos-component.log"
	}

	// Initialize logger with configurable path
	loggerConfig := &logging.LoggerConfig{
		Level:         logging.InfoLevel,
		FileName:      logFilePath,
		LoggerName:    "katharos-component",
		ComponentName: "main",
		ServiceName:   "katharos-component",
	}

	logger, err := logging.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}
