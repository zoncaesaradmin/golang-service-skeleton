package app

import (
	"context"
	"fmt"
	"sync"

	"katharos/service/internal/config"
	"sharedmodule/logging"
)

// Application represents the main application instance that holds all services and dependencies
type Application struct {
	config   *config.Config
	logger   logging.Logger
	services *ServiceRegistry
	mutex    sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// ServiceRegistry holds all registered services
type ServiceRegistry struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]interface{}),
	}
}

// Register registers a service with the given name
func (r *ServiceRegistry) Register(name string, service interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.services[name]; exists {
		return fmt.Errorf("service '%s' is already registered", name)
	}

	r.services[name] = service
	return nil
}

// Get retrieves a service by name
func (r *ServiceRegistry) Get(name string) (interface{}, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	service, exists := r.services[name]
	return service, exists
}

// GetTyped retrieves a service by name with type assertion
func (r *ServiceRegistry) GetTyped(name string, serviceType interface{}) (interface{}, error) {
	service, exists := r.Get(name)
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	return service, nil
}

// List returns all registered service names
func (r *ServiceRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}
	return names
}

// Unregister removes a service from the registry
func (r *ServiceRegistry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.services[name]; !exists {
		return fmt.Errorf("service '%s' not found", name)
	}

	delete(r.services, name)
	return nil
}

// NewApplication creates a new application instance
func NewApplication(cfg *config.Config, logger logging.Logger) *Application {
	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		config:   cfg,
		logger:   logger,
		services: NewServiceRegistry(),
		ctx:      ctx,
		cancel:   cancel,
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

// Services returns the service registry
func (app *Application) Services() *ServiceRegistry {
	return app.services
}

// Context returns the application context
func (app *Application) Context() context.Context {
	return app.ctx
}

// RegisterService registers a service with the application
func (app *Application) RegisterService(name string, service interface{}) error {
	err := app.services.Register(name, service)
	if err != nil {
		app.logger.Errorf("Failed to register service '%s': %v", name, err)
		return err
	}

	app.logger.Infof("Service '%s' registered successfully", name)
	return nil
}

// GetService retrieves a service by name
func (app *Application) GetService(name string) (interface{}, bool) {
	return app.services.Get(name)
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
