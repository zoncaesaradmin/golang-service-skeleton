package logging

import (
	"os"
	"strings"
	"testing"
)

const (
	testLogFile          = "/tmp/test.log"
	testLoggerName       = "test-logger"
	testServiceName      = "test-service"
	testComponentName    = "test"
	newLoggerErrorFmt    = "NewZerologLoggerWithConfig() error = %v"
	logFileNotCreatedFmt = "Log file was not created: %v"
	logFileEmptyMsg      = "Log file is empty"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{PanicLevel, "PANIC"},
		{Level(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoggerConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  LoggerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: LoggerConfig{
				Level:         InfoLevel,
				FileName:      testLogFile,
				LoggerName:    testLoggerName,
				ComponentName: testComponentName,
				ServiceName:   testServiceName,
			},
			wantErr: false,
		},
		{
			name: "missing filename",
			config: LoggerConfig{
				Level:         InfoLevel,
				LoggerName:    testLoggerName,
				ComponentName: testComponentName,
				ServiceName:   testServiceName,
			},
			wantErr: true,
			errMsg:  "filename is required",
		},
		{
			name: "missing logger name",
			config: LoggerConfig{
				Level:         InfoLevel,
				FileName:      testLogFile,
				ComponentName: testComponentName,
				ServiceName:   testServiceName,
			},
			wantErr: true,
			errMsg:  "logger name is required",
		},
		{
			name: "missing component name",
			config: LoggerConfig{
				Level:       InfoLevel,
				FileName:    testLogFile,
				LoggerName:  testLoggerName,
				ServiceName: testServiceName,
			},
			wantErr: true,
			errMsg:  "component name is required",
		},
		{
			name: "missing service name",
			config: LoggerConfig{
				Level:         InfoLevel,
				FileName:      testLogFile,
				LoggerName:    testLoggerName,
				ComponentName: testComponentName,
			},
			wantErr: true,
			errMsg:  "service name is required",
		},
	}

	checkValidation := func(t *testing.T, tt struct {
		name    string
		config  LoggerConfig
		wantErr bool
		errMsg  string
	}) {
		err := tt.config.Validate()
		if tt.wantErr {
			if err == nil {
				t.Errorf("LoggerConfig.Validate() expected error but got none")
				return
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("LoggerConfig.Validate() error = %v, want error containing %v", err, tt.errMsg)
			}
		} else if err != nil {
			t.Errorf("LoggerConfig.Validate() unexpected error = %v", err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkValidation(t, tt)
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != InfoLevel {
		t.Errorf("DefaultConfig() Level = %v, want %v", config.Level, InfoLevel)
	}

	if config.LoggerName != "default" {
		t.Errorf("DefaultConfig() LoggerName = %v, want %v", config.LoggerName, "default")
	}

	if config.ComponentName != "application" {
		t.Errorf("DefaultConfig() ComponentName = %v, want %v", config.ComponentName, "application")
	}

	if config.ServiceName != "service" {
		t.Errorf("DefaultConfig() ServiceName = %v, want %v", config.ServiceName, "service")
	}
}

func TestKeysAndValuesToFields(t *testing.T) {
	tests := []struct {
		name           string
		keysAndValues  []interface{}
		expectedFields Fields
	}{
		{
			name:           "empty",
			keysAndValues:  []interface{}{},
			expectedFields: Fields{},
		},
		{
			name:           "single pair",
			keysAndValues:  []interface{}{"key", "value"},
			expectedFields: Fields{"key": "value"},
		},
		{
			name:           "multiple pairs",
			keysAndValues:  []interface{}{"key1", "value1", "key2", 42, "key3", true},
			expectedFields: Fields{"key1": "value1", "key2": 42, "key3": true},
		},
		{
			name:           "odd number of arguments",
			keysAndValues:  []interface{}{"key1", "value1", "key2"},
			expectedFields: Fields{"key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keysAndValuesToFields(tt.keysAndValues...)

			if len(result) != len(tt.expectedFields) {
				t.Errorf("Expected %d fields, got %d", len(tt.expectedFields), len(result))
			}

			for key, expectedValue := range tt.expectedFields {
				if result[key] != expectedValue {
					t.Errorf("Expected field %s=%v, got %v", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	logFile := "/tmp/test_new_logger_wrapper.log"
	defer func() {
		// Clean up
		_ = os.Remove(logFile)
	}()

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    "wrapper-test-logger",
		ComponentName: "wrapper-test",
		ServiceName:   "wrapper-test-service",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// Test basic functionality
	if logger.GetLevel() != InfoLevel {
		t.Errorf("GetLevel() = %v, want %v", logger.GetLevel(), InfoLevel)
	}

	// Test logging
	logger.Info("Test message from wrapper")
	logger.Warn("Warning message from wrapper")
	logger.Error("Error message from wrapper")

	// Check if log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf(logFileNotCreatedFmt, err)
	}
}

func TestNewLoggerWithNilConfig(t *testing.T) {
	logger, err := NewLogger(nil)
	if err != nil {
		t.Fatalf("NewLogger(nil) error = %v", err)
	}
	defer logger.Close()

	// Should use default config
	if logger.GetLevel() != InfoLevel {
		t.Errorf("GetLevel() = %v, want %v", logger.GetLevel(), InfoLevel)
	}
}

func TestNewLoggerWithInvalidConfig(t *testing.T) {
	config := &LoggerConfig{
		// Missing required fields
	}

	logger, err := NewLogger(config)
	if err == nil {
		if logger != nil {
			logger.Close()
		}
		t.Error("NewLogger() expected error for invalid config but got none")
	}
}
