package logging

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestNewZerologLoggerWithConfig(t *testing.T) {
	logFile := "/tmp/test_new_logger.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:          DebugLevel,
		FileName:       logFile,
		LoggerName:     testLoggerName,
		ComponentName:  testComponentName,
		ServiceName:    testServiceName,
		MaxAge:         7,
		MaxBackups:     5,
		MaxSize:        100,
		IsLogRotatable: false,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test basic functionality
	if logger.GetLevel() != DebugLevel {
		t.Errorf("GetLevel() = %v, want %v", logger.GetLevel(), DebugLevel)
	}

	// Test logging
	logger.Info("Test message")
	logger.Debug("Debug message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Check if log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf(logFileNotCreatedFmt, err)
	}

	os.Remove(logFile)
}

func TestNewZerologLoggerWithConfigFileError(t *testing.T) {
	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      "/invalid/path/test.log",
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err == nil {
		if logger != nil {
			logger.Close()
		}
		t.Error("NewZerologLoggerWithConfig() expected error for invalid file path but got none")
	}
}

func TestZerologLoggerSetLevel(t *testing.T) {
	logFile := "/tmp/test_set_level.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test initial level
	if logger.GetLevel() != InfoLevel {
		t.Errorf("Initial level = %v, want %v", logger.GetLevel(), InfoLevel)
	}

	// Test setting level
	logger.SetLevel(ErrorLevel)
	if logger.GetLevel() != ErrorLevel {
		t.Errorf("After SetLevel(ErrorLevel) = %v, want %v", logger.GetLevel(), ErrorLevel)
	}

	// Test level filtering after change
	if logger.IsLevelEnabled(InfoLevel) {
		t.Error("Info level should not be enabled after setting level to Error")
	}

	if !logger.IsLevelEnabled(ErrorLevel) {
		t.Error("Error level should be enabled after setting level to Error")
	}

	os.Remove(logFile)
}

func TestZerologLoggerFormattedLogging(t *testing.T) {
	logFile := "/tmp/test_formatted_logging.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         DebugLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test formatted logging methods
	logger.Debugf("Debug message with %s and %d", "string", 42)
	logger.Infof("Info message with %s", "parameter")
	logger.Warnf("Warning message with %d warnings", 3)
	logger.Errorf("Error message with error code %d", 500)

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerVariadicLogging(t *testing.T) {
	logFile := "/tmp/test_variadic_logging.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         DebugLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test variadic logging methods
	logger.Debugw("Debug message", "key1", "value1", "key2", 42)
	logger.Infow("Info message", "user", "john", "action", "login")
	logger.Warnw("Warning message", "count", 3, "threshold", 5)
	logger.Errorw("Error message", "error_code", 500, "message", "internal error")

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerWithFields(t *testing.T) {
	logFile := "/tmp/test_with_fields.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test WithFields
	fieldLogger := logger.WithFields(Fields{
		"user_id":    123,
		"session_id": "abc-def-ghi",
		"module":     "auth",
	})

	fieldLogger.Info("User logged in")
	fieldLogger.Warn("Password attempt failed")

	// Test WithField
	singleFieldLogger := logger.WithField("request_id", "req-12345")
	singleFieldLogger.Info("Processing request")

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerWithError(t *testing.T) {
	logFile := "/tmp/test_with_error.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test WithError with actual error
	testErr := errors.New("test error message")
	errorLogger := logger.WithError(testErr)
	errorLogger.Error("An error occurred")

	// Test WithError with nil error (should not add error field)
	nilErrorLogger := logger.WithError(nil)
	nilErrorLogger.Info("No error here")

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerWithContext(t *testing.T) {
	logFile := "/tmp/test_with_context.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test WithContext
	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	contextLogger := logger.WithContext(ctx)
	contextLogger.Info("Request processed")

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerGenericLogging(t *testing.T) {
	logFile := "/tmp/test_generic_logging.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         DebugLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test Log, Logf, and Logw methods
	logger.Log(InfoLevel, "Generic info message")
	logger.Logf(WarnLevel, "Generic warning with %d items", 5)
	logger.Logw(ErrorLevel, "Generic error", "code", 404, "path", "/api/users")

	// Test with disabled level
	logger.Log(Level(999), "This should not log")

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerClone(t *testing.T) {
	logFile := "/tmp/test_clone.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Add some fields to original logger
	originalWithFields := logger.WithFields(Fields{"original": "value", "shared": "original"})

	// Clone the logger
	clonedLogger := originalWithFields.Clone()

	// Add different fields to cloned logger
	clonedWithFields := clonedLogger.WithFields(Fields{"cloned": "value", "shared": "cloned"})

	// Log with both loggers
	originalWithFields.Info("Message from original logger")
	clonedWithFields.Info("Message from cloned logger")

	// Verify both can log independently
	if originalWithFields.GetLevel() != clonedLogger.GetLevel() {
		t.Error("Cloned logger should have same level as original")
	}

	// Check if log file was created and has content
	if stat, err := os.Stat(logFile); err != nil {
		t.Errorf(logFileNotCreatedFmt, err)
	} else if stat.Size() == 0 {
		t.Error(logFileEmptyMsg)
	}

	os.Remove(logFile)
}

func TestZerologLoggerLevelConversion(t *testing.T) {
	tests := []struct {
		ourLevel     Level
		zerologLevel string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{PanicLevel, "panic"},
		{Level(999), "info"}, // Unknown level defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.ourLevel.String(), func(t *testing.T) {
			logFile := "/tmp/test_level_conversion.log"
			os.Remove(logFile)

			config := &LoggerConfig{
				Level:         tt.ourLevel,
				FileName:      logFile,
				LoggerName:    testLoggerName,
				ComponentName: testComponentName,
				ServiceName:   testServiceName,
			}

			logger, err := NewZerologLoggerWithConfig(config)
			if err != nil {
				t.Fatalf(newLoggerErrorFmt, err)
			}
			defer logger.Close()

			// Just verify the logger was created successfully with the level
			if logger.GetLevel() != tt.ourLevel {
				t.Errorf("Logger level = %v, want %v", logger.GetLevel(), tt.ourLevel)
			}

			os.Remove(logFile)
		})
	}
}

func TestZerologLoggerClose(t *testing.T) {
	logFile := "/tmp/test_close.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         InfoLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}

	// Log something first
	logger.Info("Test message before close")

	// Close the logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Try to close again (should not error)
	err = logger.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}

	os.Remove(logFile)
}

func TestZerologLoggerLevelFiltering(t *testing.T) {
	logFile := "/tmp/test_level_filtering_detailed.log"
	os.Remove(logFile)

	config := &LoggerConfig{
		Level:         WarnLevel,
		FileName:      logFile,
		LoggerName:    testLoggerName,
		ComponentName: testComponentName,
		ServiceName:   testServiceName,
	}

	logger, err := NewZerologLoggerWithConfig(config)
	if err != nil {
		t.Fatalf(newLoggerErrorFmt, err)
	}
	defer logger.Close()

	// Test that debug and info are filtered out
	logger.Debug("This debug message should not appear")
	logger.Info("This info message should not appear")
	logger.Debugf("This debug formatted message should not appear: %s", "test")
	logger.Infof("This info formatted message should not appear: %s", "test")

	// Test that warn and error go through
	logger.Warn("This warning should appear")
	logger.Error("This error should appear")
	logger.Warnf("This warning formatted message should appear: %s", "test")
	logger.Errorf("This error formatted message should appear: %s", "test")

	// Read the log file and verify only warn/error messages are present
	if content, err := os.ReadFile(logFile); err != nil {
		t.Errorf("Failed to read log file: %v", err)
	} else {
		contentStr := string(content)
		if strings.Contains(contentStr, "debug") || strings.Contains(contentStr, "info") {
			// This might contain the level field, so check for the actual message content
			if strings.Contains(contentStr, "This debug message") || strings.Contains(contentStr, "This info message") {
				t.Error("Debug/Info messages should be filtered out but were found in log file")
			}
		}
		if !strings.Contains(contentStr, "warning") || !strings.Contains(contentStr, "error") {
			t.Error("Warning/Error messages should be present in log file")
		}
	}

	os.Remove(logFile)
}
