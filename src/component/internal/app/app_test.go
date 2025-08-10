package app

import (
	"context"
	"testing"

	"compmodule/internal/config"
	"sharedmodule/logging"
)

const (
	// Test constants
	testHost           = "test-host"
	expectedNoErrorMsg = "expected no error, got %v"
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

	if app.ctx == nil {
		t.Error("expected context to be initialized")
	}

	if app.cancel == nil {
		t.Error("expected cancel function to be initialized")
	}
}

func TestApplicationConfig(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: testHost},
	}
	logger := &mockLogger{}

	app := NewApplication(cfg, logger)
	retrievedConfig := app.Config()

	if retrievedConfig != cfg {
		t.Error("expected retrieved config to match original")
	}
	if retrievedConfig.Server.Host != testHost {
		t.Errorf("expected host to be '%s', got %s", testHost, retrievedConfig.Server.Host)
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
		t.Errorf(expectedNoErrorMsg, err)
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
