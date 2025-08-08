package app

import (
	"context"
	"testing"

	"katharos/service/internal/config"
	"sharedmodule/logging"
)

// Mock logger for testing
type mockLogger struct{}

func (m *mockLogger) SetLevel(level logging.Level)                           { /* no-op for testing */ }
func (m *mockLogger) GetLevel() logging.Level                                { return logging.InfoLevel }
func (m *mockLogger) IsLevelEnabled(level logging.Level) bool                { return true }
func (m *mockLogger) Debug(msg string)                                       { /* no-op for testing */ }
func (m *mockLogger) Info(msg string)                                        { /* no-op for testing */ }
func (m *mockLogger) Warn(msg string)                                        { /* no-op for testing */ }
func (m *mockLogger) Error(msg string)                                       { /* no-op for testing */ }
func (m *mockLogger) Fatal(msg string)                                       { /* no-op for testing */ }
func (m *mockLogger) Panic(msg string)                                       { /* no-op for testing */ }
func (m *mockLogger) Debugf(format string, args ...interface{})              { /* no-op for testing */ }
func (m *mockLogger) Infof(format string, args ...interface{})               { /* no-op for testing */ }
func (m *mockLogger) Warnf(format string, args ...interface{})               { /* no-op for testing */ }
func (m *mockLogger) Errorf(format string, args ...interface{})              { /* no-op for testing */ }
func (m *mockLogger) Fatalf(format string, args ...interface{})              { /* no-op for testing */ }
func (m *mockLogger) Panicf(format string, args ...interface{})              { /* no-op for testing */ }
func (m *mockLogger) Debugw(msg string, keysAndValues ...interface{})        { /* no-op for testing */ }
func (m *mockLogger) Infow(msg string, keysAndValues ...interface{})         { /* no-op for testing */ }
func (m *mockLogger) Warnw(msg string, keysAndValues ...interface{})         { /* no-op for testing */ }
func (m *mockLogger) Errorw(msg string, keysAndValues ...interface{})        { /* no-op for testing */ }
func (m *mockLogger) Fatalw(msg string, keysAndValues ...interface{})        { /* no-op for testing */ }
func (m *mockLogger) Panicw(msg string, keysAndValues ...interface{})        { /* no-op for testing */ }
func (m *mockLogger) WithFields(fields logging.Fields) logging.Logger        { return m }
func (m *mockLogger) WithField(key string, value interface{}) logging.Logger { return m }
func (m *mockLogger) WithError(err error) logging.Logger                     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logging.Logger         { return m }
func (m *mockLogger) Log(level logging.Level, msg string)                    { /* no-op for testing */ }
func (m *mockLogger) Logf(level logging.Level, format string, args ...interface{}) { /* no-op for testing */
}
func (m *mockLogger) Logw(level logging.Level, msg string, keysAndValues ...interface{}) { /* no-op for testing */
}
func (m *mockLogger) Clone() logging.Logger { return &mockLogger{} }
func (m *mockLogger) Close() error          { return nil }

func TestNewServiceRegistry(t *testing.T) {
	registry := NewServiceRegistry()

	if registry == nil {
		t.Fatal("expected registry to not be nil")
	}

	if registry.services == nil {
		t.Error("expected services map to be initialized")
	}
}

func TestServiceRegistryRegister(t *testing.T) {
	registry := NewServiceRegistry()
	testService := "test-service-instance"

	err := registry.Register("test-service", testService)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test duplicate registration
	err = registry.Register("test-service", testService)
	if err == nil {
		t.Error("expected error for duplicate service registration")
	}
	expectedError := "service 'test-service' is already registered"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestServiceRegistryGet(t *testing.T) {
	registry := NewServiceRegistry()
	testService := "test-service-instance"

	// Test getting non-existent service
	service, exists := registry.Get("non-existent")
	if exists {
		t.Error("expected service to not exist")
	}
	if service != nil {
		t.Error("expected service to be nil")
	}

	// Register and get service
	registry.Register("test-service", testService)
	service, exists = registry.Get("test-service")
	if !exists {
		t.Error("expected service to exist")
	}
	if service != testService {
		t.Errorf("expected service to be '%s', got %v", testService, service)
	}
}

func TestServiceRegistryGetTyped(t *testing.T) {
	registry := NewServiceRegistry()
	testService := "test-service-instance"

	// Test getting non-existent service
	service, err := registry.GetTyped("non-existent", testService)
	if err == nil {
		t.Error("expected error for non-existent service")
	}
	expectedError := "service 'non-existent' not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}

	// Register and get typed service
	registry.Register("test-service", testService)
	service, err = registry.GetTyped("test-service", testService)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if service != testService {
		t.Errorf("expected service to be '%s', got %v", testService, service)
	}
}

func TestServiceRegistryList(t *testing.T) {
	registry := NewServiceRegistry()

	// Test empty registry
	names := registry.List()
	if len(names) != 0 {
		t.Errorf("expected 0 services, got %d", len(names))
	}

	// Add services and test list
	registry.Register("service1", "instance1")
	registry.Register("service2", "instance2")

	names = registry.List()
	if len(names) != 2 {
		t.Errorf("expected 2 services, got %d", len(names))
	}

	// Check that all service names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	if !nameMap["service1"] {
		t.Error("expected 'service1' to be in the list")
	}
	if !nameMap["service2"] {
		t.Error("expected 'service2' to be in the list")
	}
}

func TestServiceRegistryUnregister(t *testing.T) {
	registry := NewServiceRegistry()
	testService := "test-service-instance"

	// Test unregistering non-existent service
	err := registry.Unregister("non-existent")
	if err == nil {
		t.Error("expected error for non-existent service")
	}
	expectedError := "service 'non-existent' not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}

	// Register, then unregister service
	registry.Register("test-service", testService)
	err = registry.Unregister("test-service")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify service is gone
	_, exists := registry.Get("test-service")
	if exists {
		t.Error("expected service to be unregistered")
	}
}

func TestNewApplication(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)

	if app == nil {
		t.Fatal("expected application to not be nil")
	}

	if app.config != cfg {
		t.Error("expected config to be set correctly")
	}

	if app.logger != logger {
		t.Error("expected logger to be set correctly")
	}

	if app.services == nil {
		t.Error("expected services registry to be initialized")
	}

	if app.ctx == nil {
		t.Error("expected context to be initialized")
	}

	if app.cancel == nil {
		t.Error("expected cancel function to be initialized")
	}
}

func TestApplicationConfig(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "test-host"},
	}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	retrievedConfig := app.Config()

	if retrievedConfig != cfg {
		t.Error("expected retrieved config to match original")
	}
	if retrievedConfig.Server.Host != "test-host" {
		t.Errorf("expected host to be 'test-host', got %s", retrievedConfig.Server.Host)
	}
}

func TestApplicationLogger(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	retrievedLogger := app.Logger()

	if retrievedLogger != logger {
		t.Error("expected retrieved logger to match original")
	}
}

func TestApplicationServices(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	services := app.Services()

	if services == nil {
		t.Error("expected services to not be nil")
	}
	if services != app.services {
		t.Error("expected returned services to match internal registry")
	}
}

func TestApplicationContext(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	ctx := app.Context()

	if ctx == nil {
		t.Error("expected context to not be nil")
	}

	// Test that context is not cancelled initially
	select {
	case <-ctx.Done():
		t.Error("expected context to not be cancelled initially")
	default:
		// Context is not cancelled, which is expected
	}
}

func TestApplicationRegisterService(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	testService := "test-service-instance"

	err := app.RegisterService("test-service", testService)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify service is registered
	service, exists := app.GetService("test-service")
	if !exists {
		t.Error("expected service to be registered")
	}
	if service != testService {
		t.Errorf("expected service to be '%s', got %v", testService, service)
	}
}

func TestApplicationGetService(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	testService := "test-service-instance"

	// Test getting non-existent service
	service, exists := app.GetService("non-existent")
	if exists {
		t.Error("expected service to not exist")
	}
	if service != nil {
		t.Error("expected service to be nil")
	}

	// Register and get service
	app.RegisterService("test-service", testService)
	service, exists = app.GetService("test-service")
	if !exists {
		t.Error("expected service to exist")
	}
	if service != testService {
		t.Errorf("expected service to be '%s', got %v", testService, service)
	}
}

func TestApplicationShutdown(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	ctx := app.Context()

	// Verify context is not cancelled before shutdown
	select {
	case <-ctx.Done():
		t.Error("expected context to not be cancelled before shutdown")
	default:
		// Expected
	}

	err := app.Shutdown()
	if err != nil {
		t.Errorf("expected no error during shutdown, got %v", err)
	}

	// Verify context is cancelled after shutdown
	select {
	case <-ctx.Done():
		// Expected - context should be cancelled
	default:
		t.Error("expected context to be cancelled after shutdown")
	}
}

func TestApplicationIsShuttingDown(t *testing.T) {
	cfg := &config.Config{}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)

	// Initially should not be shutting down
	if app.IsShuttingDown() {
		t.Error("expected application to not be shutting down initially")
	}

	// After shutdown, should be shutting down
	app.Shutdown()
	if !app.IsShuttingDown() {
		t.Error("expected application to be shutting down after shutdown")
	}
}
