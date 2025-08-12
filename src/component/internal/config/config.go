package config

import (
	"fmt"
	"os"

	"sharedmodule/utils"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	FilePath string `yaml:"file_path"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Host:         utils.GetEnv("SERVER_HOST", "localhost"),
			Port:         utils.GetEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  utils.GetEnvInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: utils.GetEnvInt("SERVER_WRITE_TIMEOUT", 10),
		},
		Logging: LoggingConfig{
			Level:    utils.GetEnv("LOG_LEVEL", "info"),
			Format:   utils.GetEnv("LOG_FORMAT", "json"),
			FilePath: utils.GetEnv("LOG_FILE_PATH", "/tmp/katharos-component.log"),
		},
	}

	return config
}

// LoadConfigFromFile loads configuration from a YAML file with optional environment variable overrides
func LoadConfigFromFile(configPath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	// Parse YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing YAML config file %s: %w", configPath, err)
	}

	// Override with environment variables if they exist
	overrideWithEnvVars(config)

	return config, nil
}

// LoadConfigWithDefaults loads configuration from file if it exists, falling back to environment variables and defaults
func LoadConfigWithDefaults(configPath string) *Config {
	// Try to load from file first
	if config, err := LoadConfigFromFile(configPath); err == nil {
		return config
	}

	// Fallback to environment variables and defaults
	return LoadConfig()
}

// overrideWithEnvVars overrides config values with environment variables if they are set
func overrideWithEnvVars(config *Config) {
	// Server configuration overrides
	if host := utils.GetEnv("SERVER_HOST", ""); host != "" {
		config.Server.Host = host
	}
	if port := utils.GetEnvInt("SERVER_PORT", -1); port != -1 {
		config.Server.Port = port
	}
	if readTimeout := utils.GetEnvInt("SERVER_READ_TIMEOUT", -1); readTimeout != -1 {
		config.Server.ReadTimeout = readTimeout
	}
	if writeTimeout := utils.GetEnvInt("SERVER_WRITE_TIMEOUT", -1); writeTimeout != -1 {
		config.Server.WriteTimeout = writeTimeout
	}

	// Logging configuration overrides
	if level := utils.GetEnv("LOG_LEVEL", ""); level != "" {
		config.Logging.Level = level
	}
	if format := utils.GetEnv("LOG_FORMAT", ""); format != "" {
		config.Logging.Format = format
	}
	if filePath := utils.GetEnv("LOG_FILE_PATH", ""); filePath != "" {
		config.Logging.FilePath = filePath
	}
}
