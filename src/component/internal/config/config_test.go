package config

import (
	"path/filepath"
	"os"
	"testing"
)

const errCreateTestConfigFile = "Failed to create test config file: %v"

func TestLoadConfigDefaults(t *testing.T) {
	config := LoadConfig()
	
	if config.Server.Host != "localhost" {
		t.Errorf("Expected localhost, got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("Expected 8080, got %d", config.Server.Port)
	}
	if config.Database.Port != 5432 {
		t.Errorf("Expected 5432, got %d", config.Database.Port)
	}
	if config.Logging.Level != "info" {
		t.Errorf("Expected info, got %s", config.Logging.Level)
	}
}

func TestLoadConfigWithEnvVars(t *testing.T) {
	// Save original values
	originalHost := os.Getenv("SERVER_HOST")
	originalPort := os.Getenv("SERVER_PORT")
	
	// Set test values
	os.Setenv("SERVER_HOST", "testhost")
	os.Setenv("SERVER_PORT", "9999")
	
	config := LoadConfig()
	
	if config.Server.Host != "testhost" {
		t.Errorf("Expected testhost, got %s", config.Server.Host)
	}
	if config.Server.Port != 9999 {
		t.Errorf("Expected 9999, got %d", config.Server.Port)
	}
	
	// Restore original values
	if originalHost != "" {
		os.Setenv("SERVER_HOST", originalHost)
	} else {
		os.Unsetenv("SERVER_HOST")
	}
	if originalPort != "" {
		os.Setenv("SERVER_PORT", originalPort)
	} else {
		os.Unsetenv("SERVER_PORT")
	}
}

func TestLoadConfigFromFileSuccess(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()
	
	// Create a temporary config file
	configContent := `
server:
  host: "test.example.com"
  port: 9000
  read_timeout: 30
  write_timeout: 40
database:
  host: "db.test.com"
  port: 3306
  username: "testuser"
  password: "testpass"
  database: "testdb"
logging:
  level: "debug"
  format: "text"
`
	err := os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	config, err := LoadConfigFromFile(filepath.Join(tempDir, "config.yaml"))
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}
	
	// Verify the values were loaded correctly
	if config.Server.Host != "test.example.com" {
		t.Errorf("Expected host 'test.example.com', got %s", config.Server.Host)
	}
	if config.Server.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", config.Server.Port)
	}
	if config.Database.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", config.Database.Username)
	}
	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.Logging.Level)
	}
}

func TestLoadConfigFromFileError(t *testing.T) {
	// Test with non-existent file
	_, err := LoadConfigFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	// Test with invalid YAML
	tempDir := t.TempDir()
	invalidYAML := `
server:
  host: "test.com
  port: invalid
database:
  - this is not valid YAML structure
`
	err = os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	_, err = LoadConfigFromFile(filepath.Join(tempDir, "config.yaml"))
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadConfigWithDefaultsFileExists(t *testing.T) {
	tempDir := t.TempDir()
	
	configContent := `
server:
  host: "file.example.com"
  port: 7777
logging:
  level: "warn"
`
	
	configFile := tempDir + "/existing_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
	}
	
	config := LoadConfigWithDefaults(configFile)
	
	// Should load from file
	if config.Server.Host != "file.example.com" {
		t.Errorf("Expected host from file 'file.example.com', got %s", config.Server.Host)
	}
	if config.Server.Port != 7777 {
		t.Errorf("Expected port from file 7777, got %d", config.Server.Port)
	}
}

func TestLoadConfigWithDefaultsFallback(t *testing.T) {
	// Test with non-existent file - should fall back to defaults
	config := LoadConfigWithDefaults("/nonexistent/config.yaml")
	
	// Should use default values
	if config.Server.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Server.Port)
	}
}

func TestInvalidEnvVars(t *testing.T) {
	// Save original values
	originalPort := os.Getenv("SERVER_PORT")
	originalDbPort := os.Getenv("DB_PORT")
	
	// Set invalid values
	os.Setenv("SERVER_PORT", "not-a-number")
	os.Setenv("DB_PORT", "invalid-port")
	
	config := LoadConfig()
	
	// Should use default values when env vars are invalid
	if config.Server.Port != 8080 {
		t.Errorf("Expected default port 8080 for invalid env var, got %d", config.Server.Port)
	}
	if config.Database.Port != 5432 {
		t.Errorf("Expected default database port 5432 for invalid env var, got %d", config.Database.Port)
	}
	
	// Restore original values
	if originalPort != "" {
		os.Setenv("SERVER_PORT", originalPort)
	} else {
		os.Unsetenv("SERVER_PORT")
	}
	if originalDbPort != "" {
		os.Setenv("DB_PORT", originalDbPort)
	} else {
		os.Unsetenv("DB_PORT")
	}
}

func TestOverrideWithEnvVars(t *testing.T) {
	// Save original values
	originalHost := os.Getenv("SERVER_HOST")
	originalDbPort := os.Getenv("DB_PORT")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	
	// Test override functionality
	config := &Config{
		Server: ServerConfig{
			Host: "original.com",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host: "orig-db.com", 
			Port: 5432,
		},
		Logging: LoggingConfig{
			Level: "info",
			Format: "json",
		},
	}
	
	// Set environment variables to override some values
	os.Setenv("SERVER_HOST", "override.com")
	os.Setenv("DB_PORT", "9999")
	os.Setenv("LOG_LEVEL", "debug")
	
	// Call the function
	overrideWithEnvVars(config)
	
	// Check that specified values were overridden
	if config.Server.Host != "override.com" {
		t.Errorf("Expected server host 'override.com', got %s", config.Server.Host)
	}
	if config.Database.Port != 9999 {
		t.Errorf("Expected database port 9999, got %d", config.Database.Port)
	}
	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.Logging.Level)
	}
	
	// Check that non-overridden values remained the same
	if config.Server.Port != 8080 {
		t.Errorf("Expected server port 8080 (not overridden), got %d", config.Server.Port)
	}
	if config.Database.Host != "orig-db.com" {
		t.Errorf("Expected database host 'orig-db.com' (not overridden), got %s", config.Database.Host)
	}
	if config.Logging.Format != "json" {
		t.Errorf("Expected log format 'json' (not overridden), got %s", config.Logging.Format)
	}
	
	// Restore original values
	if originalHost != "" {
		os.Setenv("SERVER_HOST", originalHost)
	} else {
		os.Unsetenv("SERVER_HOST")
	}
	if originalDbPort != "" {
		os.Setenv("DB_PORT", originalDbPort)
	} else {
		os.Unsetenv("DB_PORT")
	}
	if originalLogLevel != "" {
		os.Setenv("LOG_LEVEL", originalLogLevel)
	} else {
		os.Unsetenv("LOG_LEVEL")
	}
}

func TestOverrideWithEnvVarsAllFields(t *testing.T) {
	// Save original values
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_DATABASE",
		"LOG_LEVEL", "LOG_FORMAT",
	}
	originalValues := make(map[string]string)
	for _, env := range envVars {
		originalValues[env] = os.Getenv(env)
	}
	
	// Test override functionality for all fields
	config := &Config{
		Server: ServerConfig{
			Host: "original.com",
			Port: 8080,
			ReadTimeout: 10,
			WriteTimeout: 10,
		},
		Database: DatabaseConfig{
			Host: "orig-db.com", 
			Port: 5432,
			Username: "original_user",
			Password: "original_pass",
			Database: "original_db",
		},
		Logging: LoggingConfig{
			Level: "info",
			Format: "json",
		},
	}
	
	// Set all environment variables
	os.Setenv("SERVER_HOST", "new.com")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_READ_TIMEOUT", "20")
	os.Setenv("SERVER_WRITE_TIMEOUT", "30")
	os.Setenv("DB_HOST", "new-db.com")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USERNAME", "new_user")
	os.Setenv("DB_PASSWORD", "new_pass")
	os.Setenv("DB_DATABASE", "new_db")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "text")
	
	// Call the function
	overrideWithEnvVars(config)
	
	// Check all values were overridden
	if config.Server.Host != "new.com" {
		t.Errorf("Expected server host 'new.com', got %s", config.Server.Host)
	}
	if config.Server.Port != 9090 {
		t.Errorf("Expected server port 9090, got %d", config.Server.Port)
	}
	if config.Server.ReadTimeout != 20 {
		t.Errorf("Expected read timeout 20, got %d", config.Server.ReadTimeout)
	}
	if config.Server.WriteTimeout != 30 {
		t.Errorf("Expected write timeout 30, got %d", config.Server.WriteTimeout)
	}
	if config.Database.Host != "new-db.com" {
		t.Errorf("Expected database host 'new-db.com', got %s", config.Database.Host)
	}
	if config.Database.Port != 3306 {
		t.Errorf("Expected database port 3306, got %d", config.Database.Port)
	}
	if config.Database.Username != "new_user" {
		t.Errorf("Expected username 'new_user', got %s", config.Database.Username)
	}
	if config.Database.Password != "new_pass" {
		t.Errorf("Expected password 'new_pass', got %s", config.Database.Password)
	}
	if config.Database.Database != "new_db" {
		t.Errorf("Expected database 'new_db', got %s", config.Database.Database)
	}
	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.Logging.Level)
	}
	if config.Logging.Format != "text" {
		t.Errorf("Expected log format 'text', got %s", config.Logging.Format)
	}
	
	// Restore original values
	for _, env := range envVars {
		if originalValues[env] != "" {
			os.Setenv(env, originalValues[env])
		} else {
			os.Unsetenv(env)
		}
	}
}
