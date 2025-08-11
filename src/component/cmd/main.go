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

	"compmodule/internal/api"
	"compmodule/internal/app"
	"compmodule/internal/config"
	"compmodule/internal/processing"
	"sharedmodule/logging"
)

func main() {
	// Load configuration with smart fallback (file first, then env+defaults)
	cfg := config.LoadConfigWithDefaults("config.yaml")

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
	defer logger.Close()

	// Log the configuration being used (without sensitive data)
	logger.Infof("Starting katharos-component on %s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Logging to: %s", logFilePath)

	// Create application instance
	application := app.NewApplication(cfg, logger)

	// Initialize handlers (without user service for now)
	handler := api.NewHandler()

	// Setup HTTP mux
	mux := setupRouter(handler)

	// Start processing pipeline for message bus communication
	processingConfig := processing.Config{
		Input: processing.InputConfig{
			Topics:            []string{"test_input"},
			PollTimeout:       1 * time.Second,
			ChannelBufferSize: 100,
		},
		Processor: processing.ProcessorConfig{
			ProcessingDelay: 10 * time.Millisecond,
			BatchSize:       10,
		},
		Output: processing.OutputConfig{
			OutputTopic:       "test_output",
			BatchSize:         10,
			FlushTimeout:      1 * time.Second,
			ChannelBufferSize: 100,
		},
	}

	pipeline := processing.NewPipeline(processingConfig, logger)
	if err := pipeline.Start(); err != nil {
		logger.Fatalf("Failed to start processing pipeline: %v", err)
	}

	// Start server
	startServer(mux, cfg, application, pipeline)
}

func setupRouter(handler *api.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup routes
	handler.SetupRoutes(mux)

	return mux
}

func startServer(mux *http.ServeMux, cfg *config.Config, application *app.Application, pipeline *processing.Pipeline) {
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

	// Shutdown the processing pipeline
	if err := pipeline.Stop(); err != nil {
		logger.Errorf("Pipeline shutdown error: %v", err)
	}

	// Shutdown the application
	if err := application.Shutdown(); err != nil {
		logger.Errorf("Application shutdown error: %v", err)
	}

	logger.Info("Server exited")
}
