package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Clear any existing env vars that might interfere
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_DATABASE",
		"LOG_LEVEL", "LOG_FORMAT",
	}

	originalValues := make(map[string]string)
	for _, key := range envVars {
		originalValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalValues {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	config := LoadConfig()

	// Test server defaults
	if config.Server.Host != "localhost" {
		t.Errorf("expected Server.Host to be 'localhost', got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("expected Server.Port to be 8080, got %d", config.Server.Port)
	}
	if config.Server.ReadTimeout != 10 {
		t.Errorf("expected Server.ReadTimeout to be 10, got %d", config.Server.ReadTimeout)
	}
	if config.Server.WriteTimeout != 10 {
		t.Errorf("expected Server.WriteTimeout to be 10, got %d", config.Server.WriteTimeout)
	}

	// Test database defaults
	if config.Database.Host != "localhost" {
		t.Errorf("expected Database.Host to be 'localhost', got %s", config.Database.Host)
	}
	if config.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432, got %d", config.Database.Port)
	}
	if config.Database.Username != "user" {
		t.Errorf("expected Database.Username to be 'user', got %s", config.Database.Username)
	}
	if config.Database.Password != "password" {
		t.Errorf("expected Database.Password to be 'password', got %s", config.Database.Password)
	}
	if config.Database.Database != "mydb" {
		t.Errorf("expected Database.Database to be 'mydb', got %s", config.Database.Database)
	}

	// Test logging defaults
	if config.Logging.Level != "info" {
		t.Errorf("expected Logging.Level to be 'info', got %s", config.Logging.Level)
	}
	if config.Logging.Format != "json" {
		t.Errorf("expected Logging.Format to be 'json', got %s", config.Logging.Format)
	}
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	testEnvVars := map[string]string{
		"SERVER_HOST":          "example.com",
		"SERVER_PORT":          "9090",
		"SERVER_READ_TIMEOUT":  "15",
		"SERVER_WRITE_TIMEOUT": "20",
		"DB_HOST":              "db.example.com",
		"DB_PORT":              "3306",
		"DB_USERNAME":          "testuser",
		"DB_PASSWORD":          "testpass",
		"DB_DATABASE":          "testdb",
		"LOG_LEVEL":            "debug",
		"LOG_FORMAT":           "text",
	}

	// Set env vars
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	config := LoadConfig()

	// Test server config from env
	if config.Server.Host != "example.com" {
		t.Errorf("expected Server.Host to be 'example.com', got %s", config.Server.Host)
	}
	if config.Server.Port != 9090 {
		t.Errorf("expected Server.Port to be 9090, got %d", config.Server.Port)
	}
	if config.Server.ReadTimeout != 15 {
		t.Errorf("expected Server.ReadTimeout to be 15, got %d", config.Server.ReadTimeout)
	}
	if config.Server.WriteTimeout != 20 {
		t.Errorf("expected Server.WriteTimeout to be 20, got %d", config.Server.WriteTimeout)
	}

	// Test database config from env
	if config.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com', got %s", config.Database.Host)
	}
	if config.Database.Port != 3306 {
		t.Errorf("expected Database.Port to be 3306, got %d", config.Database.Port)
	}
	if config.Database.Username != "testuser" {
		t.Errorf("expected Database.Username to be 'testuser', got %s", config.Database.Username)
	}
	if config.Database.Password != "testpass" {
		t.Errorf("expected Database.Password to be 'testpass', got %s", config.Database.Password)
	}
	if config.Database.Database != "testdb" {
		t.Errorf("expected Database.Database to be 'testdb', got %s", config.Database.Database)
	}

	// Test logging config from env
	if config.Logging.Level != "debug" {
		t.Errorf("expected Logging.Level to be 'debug', got %s", config.Logging.Level)
	}
	if config.Logging.Format != "text" {
		t.Errorf("expected Logging.Format to be 'text', got %s", config.Logging.Format)
	}
}

func TestGetEnv(t *testing.T) {
	const testKey = "TEST_GET_ENV_KEY"
	const testValue = "test_value"
	const defaultValue = "default_value"

	// Test with env var set
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := getEnv(testKey, defaultValue)
	if result != testValue {
		t.Errorf("expected '%s', got '%s'", testValue, result)
	}

	// Test with env var unset
	os.Unsetenv(testKey)
	result = getEnv(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected '%s', got '%s'", defaultValue, result)
	}

	// Test with empty env var
	os.Setenv(testKey, "")
	result = getEnv(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected '%s' for empty env var, got '%s'", defaultValue, result)
	}
}

func TestGetEnvInt(t *testing.T) {
	const testKey = "TEST_GET_ENV_INT_KEY"
	const defaultValue = 42

	// Test with valid integer env var
	os.Setenv(testKey, "123")
	defer os.Unsetenv(testKey)

	result := getEnvInt(testKey, defaultValue)
	if result != 123 {
		t.Errorf("expected 123, got %d", result)
	}

	// Test with invalid integer env var
	os.Setenv(testKey, "invalid")
	result = getEnvInt(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d for invalid env var, got %d", defaultValue, result)
	}

	// Test with env var unset
	os.Unsetenv(testKey)
	result = getEnvInt(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d for unset env var, got %d", defaultValue, result)
	}

	// Test with empty env var
	os.Setenv(testKey, "")
	result = getEnvInt(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d for empty env var, got %d", defaultValue, result)
	}
}

func TestConfigStructFields(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host:         "test-host",
			Port:         8000,
			ReadTimeout:  5,
			WriteTimeout: 10,
		},
		Database: DatabaseConfig{
			Host:     "db-host",
			Port:     5432,
			Username: "dbuser",
			Password: "dbpass",
			Database: "dbname",
		},
		Logging: LoggingConfig{
			Level:  "warn",
			Format: "plain",
		},
	}

	// Test that all fields are accessible and correct
	if config.Server.Host != "test-host" {
		t.Errorf("expected Server.Host to be 'test-host', got %s", config.Server.Host)
	}
	if config.Database.Username != "dbuser" {
		t.Errorf("expected Database.Username to be 'dbuser', got %s", config.Database.Username)
	}
	if config.Logging.Level != "warn" {
		t.Errorf("expected Logging.Level to be 'warn', got %s", config.Logging.Level)
	}
}
