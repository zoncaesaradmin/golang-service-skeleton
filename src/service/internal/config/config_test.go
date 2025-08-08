package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"sharedmodule/utils"
)

const (
	testConfigFileName    = "test_config.yaml"
	testConfigCreateError = "Failed to create test config file: %v"
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

	assertServerDefaults(t, config.Server)
	assertDatabaseDefaults(t, config.Database)
	assertLoggingDefaults(t, config.Logging)
}

func assertServerDefaults(t *testing.T, server ServerConfig) {
	if server.Host != "localhost" {
		t.Errorf("expected Server.Host to be 'localhost', got %s", server.Host)
	}
	if server.Port != 8080 {
		t.Errorf("expected Server.Port to be 8080, got %d", server.Port)
	}
	if server.ReadTimeout != 10 {
		t.Errorf("expected Server.ReadTimeout to be 10, got %d", server.ReadTimeout)
	}
	if server.WriteTimeout != 10 {
		t.Errorf("expected Server.WriteTimeout to be 10, got %d", server.WriteTimeout)
	}
}

func assertDatabaseDefaults(t *testing.T, db DatabaseConfig) {
	if db.Host != "localhost" {
		t.Errorf("expected Database.Host to be 'localhost', got %s", db.Host)
	}
	if db.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432, got %d", db.Port)
	}
	if db.Username != "user" {
		t.Errorf("expected Database.Username to be 'user', got %s", db.Username)
	}
	if db.Password != "password" {
		t.Errorf("expected Database.Password to be 'password', got %s", db.Password)
	}
	if db.Database != "mydb" {
		t.Errorf("expected Database.Database to be 'mydb', got %s", db.Database)
	}
}

func assertLoggingDefaults(t *testing.T, logging LoggingConfig) {
	if logging.Level != "info" {
		t.Errorf("expected Logging.Level to be 'info', got %s", logging.Level)
	}
	if logging.Format != "json" {
		t.Errorf("expected Logging.Format to be 'json', got %s", logging.Format)
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

	result := utils.GetEnv(testKey, defaultValue)
	if result != testValue {
		t.Errorf("expected '%s', got '%s'", testValue, result)
	}

	// Test with env var unset
	os.Unsetenv(testKey)
	result = utils.GetEnv(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected '%s', got '%s'", defaultValue, result)
	}

	// Test with empty env var
	os.Setenv(testKey, "")
	result = utils.GetEnv(testKey, defaultValue)
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

	result := utils.GetEnvInt(testKey, defaultValue)
	if result != 123 {
		t.Errorf("expected 123, got %d", result)
	}

	// Test with invalid integer env var
	os.Setenv(testKey, "invalid")
	result = utils.GetEnvInt(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d for invalid env var, got %d", defaultValue, result)
	}

	// Test with env var unset
	os.Unsetenv(testKey)
	result = utils.GetEnvInt(testKey, defaultValue)
	if result != defaultValue {
		t.Errorf("expected %d for unset env var, got %d", defaultValue, result)
	}

	// Test with empty env var
	os.Setenv(testKey, "")
	result = utils.GetEnvInt(testKey, defaultValue)
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

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	configData := `
server:
  host: "filehost"
  port: 9000
  read_timeout: 15
  write_timeout: 20

database:
  host: "filedbhost"
  port: 3306
  username: "fileuser"
  password: "filepass"
  database: "filedb"

logging:
  level: "debug"
  format: "text"
`

	// Create temporary file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, testConfigFileName)
	err := ioutil.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf(testConfigCreateError, err)
	}

	// Load config from file
	config, err := LoadConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	// Test server config from file
	if config.Server.Host != "filehost" {
		t.Errorf("expected Server.Host to be 'filehost', got %s", config.Server.Host)
	}
	if config.Server.Port != 9000 {
		t.Errorf("expected Server.Port to be 9000, got %d", config.Server.Port)
	}

	// Test database config from file
	if config.Database.Host != "filedbhost" {
		t.Errorf("expected Database.Host to be 'filedbhost', got %s", config.Database.Host)
	}
	if config.Database.Username != "fileuser" {
		t.Errorf("expected Database.Username to be 'fileuser', got %s", config.Database.Username)
	}

	// Test logging config from file
	if config.Logging.Level != "debug" {
		t.Errorf("expected Logging.Level to be 'debug', got %s", config.Logging.Level)
	}
}

func TestLoadConfigFromFileWithEnvOverrides(t *testing.T) {
	// Create a temporary config file
	configData := `
server:
  host: "filehost"
  port: 9000
  read_timeout: 15
  write_timeout: 20

database:
  host: "filedbhost"
  port: 3306
  username: "fileuser"
  password: "filepass"
  database: "filedb"

logging:
  level: "debug"
  format: "text"
`

	// Create temporary file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, testConfigFileName)
	err := ioutil.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf(testConfigCreateError, err)
	}

	// Set some environment variables to override file values
	testEnvVars := map[string]string{
		"SERVER_HOST": "envhost",
		"DB_PORT":     "5432",
		"LOG_LEVEL":   "warn",
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

	// Load config from file
	config, err := LoadConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	// Test that env vars override file values
	if config.Server.Host != "envhost" {
		t.Errorf("expected Server.Host to be 'envhost' (from env), got %s", config.Server.Host)
	}
	if config.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432 (from env), got %d", config.Database.Port)
	}
	if config.Logging.Level != "warn" {
		t.Errorf("expected Logging.Level to be 'warn' (from env), got %s", config.Logging.Level)
	}

	// Test that non-overridden values still come from file
	if config.Server.Port != 9000 {
		t.Errorf("expected Server.Port to be 9000 (from file), got %d", config.Server.Port)
	}
	if config.Database.Username != "fileuser" {
		t.Errorf("expected Database.Username to be 'fileuser' (from file), got %s", config.Database.Username)
	}
}

func TestLoadConfigFromFileNotFound(t *testing.T) {
	_, err := LoadConfigFromFile("nonexistent.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestLoadConfigFromFileInvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid.yaml")
	err := ioutil.WriteFile(configFile, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf(testConfigCreateError, err)
	}

	_, err = LoadConfigFromFile(configFile)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Test with existing config file
	configData := `
server:
  host: "filehost"
  port: 9000

database:
  host: "filedbhost"
  port: 3306

logging:
  level: "debug"
  format: "text"
`

	// Create temporary file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, testConfigFileName)
	err := ioutil.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf(testConfigCreateError, err)
	}

	config := LoadConfigWithDefaults(configFile)

	// Should load from file
	if config.Server.Host != "filehost" {
		t.Errorf("expected Server.Host to be 'filehost' (from file), got %s", config.Server.Host)
	}

	// Test with nonexistent file - should fallback to env vars and defaults
	config2 := LoadConfigWithDefaults("nonexistent.yaml")

	// Should use defaults
	if config2.Server.Host != "localhost" {
		t.Errorf("expected Server.Host to be 'localhost' (default), got %s", config2.Server.Host)
	}
	if config2.Server.Port != 8080 {
		t.Errorf("expected Server.Port to be 8080 (default), got %d", config2.Server.Port)
	}
}
