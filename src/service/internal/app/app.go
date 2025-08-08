package app

import (
	"context"
	"sync"

	"katharos/service/internal/config"
	"sharedmodule/logging"
)

// Application represents the main application instance that holds configuration and dependencies
type Application struct {
	config *config.Config
	logger logging.Logger
	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApplication creates a new application instance
func NewApplication(cfg *config.Config, logger logging.Logger) *Application {
	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		config: cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Config returns the application configuration
func (app *Application) Config() *config.Config {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.config
}

// Logger returns the application logger
func (app *Application) Logger() logging.Logger {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.logger
}

// Context returns the application context
func (app *Application) Context() context.Context {
	return app.ctx
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	app.logger.Info("Shutting down application...")

	// Cancel the application context
	app.cancel()

	// Here you can add cleanup logic for services
	// For example, calling Close() on services that implement io.Closer

	app.logger.Info("Application shutdown completed")
	return nil
}

// IsShuttingDown returns true if the application is shutting down
func (app *Application) IsShuttingDown() bool {
	select {
	case <-app.ctx.Done():
		return true
	default:
		return false
	}
}
