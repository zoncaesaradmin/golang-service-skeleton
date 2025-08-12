package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"sharedmodule/logging"
	"sharedmodule/utils"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Logging    LoggingConfig    `yaml:"logging"`
	Processing ProcessingConfig `yaml:"processing"`
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
	Level         string `yaml:"level"`          // Log level: debug, info, warn, error, fatal, panic
	Format        string `yaml:"format"`         // Log format: json, text (not used in current implementation)
	FilePath      string `yaml:"file_path"`      // Path to the log file
	LoggerName    string `yaml:"logger_name"`    // Name identifier for the logger
	ComponentName string `yaml:"component_name"` // Component/module name for structured logging
	ServiceName   string `yaml:"service_name"`   // Service name for structured logging
}

// ProcessingConfig holds processing pipeline configuration
type ProcessingConfig struct {
	Input         InputConfig     `yaml:"input"`
	Processor     ProcessorConfig `yaml:"processor"`
	Output        OutputConfig    `yaml:"output"`
	Channels      ChannelConfig   `yaml:"channels"`
	PloggerConfig LoggingConfig   `yaml:"logging"`
}

// InputConfig holds input handler configuration
type InputConfig struct {
	Topics            []string      `yaml:"topics"`
	PollTimeout       time.Duration `yaml:"poll_timeout"`
	ChannelBufferSize int           `yaml:"channel_buffer_size"`
}

// ProcessorConfig holds processor configuration
type ProcessorConfig struct {
	ProcessingDelay time.Duration `yaml:"processing_delay"`
	BatchSize       int           `yaml:"batch_size"`
}

// OutputConfig holds output handler configuration
type OutputConfig struct {
	OutputTopic       string        `yaml:"output_topic"`
	BatchSize         int           `yaml:"batch_size"`
	FlushTimeout      time.Duration `yaml:"flush_timeout"`
	ChannelBufferSize int           `yaml:"channel_buffer_size"`
}

// ChannelConfig holds channel buffer configuration
type ChannelConfig struct {
	InputBufferSize  int `yaml:"input_buffer_size"`
	OutputBufferSize int `yaml:"output_buffer_size"`
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
			Level:         utils.GetEnv("LOG_LEVEL", "info"),
			Format:        utils.GetEnv("LOG_FORMAT", "json"),
			FilePath:      utils.GetEnv("LOG_FILE_PATH", "/tmp/katharos-component.log"),
			LoggerName:    utils.GetEnv("LOG_LOGGER_NAME", "katharos-component"),
			ComponentName: utils.GetEnv("LOG_COMPONENT_NAME", "main"),
			ServiceName:   utils.GetEnv("LOG_SERVICE_NAME", "katharos-component"),
		},
		Processing: ProcessingConfig{
			Input: InputConfig{
				Topics:            parseTopics(utils.GetEnv("PROCESSING_INPUT_TOPICS", "input-topic")),
				PollTimeout:       time.Duration(utils.GetEnvInt("PROCESSING_INPUT_POLL_TIMEOUT_MS", 1000)) * time.Millisecond,
				ChannelBufferSize: utils.GetEnvInt("PROCESSING_INPUT_BUFFER_SIZE", 1000),
			},
			Processor: ProcessorConfig{
				ProcessingDelay: time.Duration(utils.GetEnvInt("PROCESSING_DELAY_MS", 10)) * time.Millisecond,
				BatchSize:       utils.GetEnvInt("PROCESSING_BATCH_SIZE", 100),
			},
			Output: OutputConfig{
				OutputTopic:       utils.GetEnv("PROCESSING_OUTPUT_TOPIC", "output-topic"),
				BatchSize:         utils.GetEnvInt("PROCESSING_OUTPUT_BATCH_SIZE", 50),
				FlushTimeout:      time.Duration(utils.GetEnvInt("PROCESSING_OUTPUT_FLUSH_TIMEOUT_MS", 5000)) * time.Millisecond,
				ChannelBufferSize: utils.GetEnvInt("PROCESSING_OUTPUT_BUFFER_SIZE", 1000),
			},
			Channels: ChannelConfig{
				InputBufferSize:  utils.GetEnvInt("PROCESSING_CHANNELS_INPUT_BUFFER_SIZE", 1000),
				OutputBufferSize: utils.GetEnvInt("PROCESSING_CHANNELS_OUTPUT_BUFFER_SIZE", 1000),
			},
			PloggerConfig: LoggingConfig{
				Level:         utils.GetEnv("PROCESSING_PLOGGER_LEVEL", "info"),
				Format:        "json", // Pipeline logger uses same format as main logger
				FilePath:      utils.GetEnv("PROCESSING_PLOGGER_FILE_NAME", "/tmp/katharos-pipeline.log"),
				LoggerName:    utils.GetEnv("PROCESSING_PLOGGER_LOGGER_NAME", "pipeline"),
				ComponentName: utils.GetEnv("PROCESSING_PLOGGER_COMPONENT_NAME", "processing"),
				ServiceName:   utils.GetEnv("PROCESSING_PLOGGER_SERVICE_NAME", "katharos"),
			},
		},
	}

	return config
}

// parseTopics parses comma-separated topics from a string
func parseTopics(topicsStr string) []string {
	if topicsStr == "" {
		return []string{}
	}
	topics := strings.Split(topicsStr, ",")
	for i, topic := range topics {
		topics[i] = strings.TrimSpace(topic)
	}
	return topics
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
	if loggerName := utils.GetEnv("LOG_LOGGER_NAME", ""); loggerName != "" {
		config.Logging.LoggerName = loggerName
	}
	if componentName := utils.GetEnv("LOG_COMPONENT_NAME", ""); componentName != "" {
		config.Logging.ComponentName = componentName
	}
	if serviceName := utils.GetEnv("LOG_SERVICE_NAME", ""); serviceName != "" {
		config.Logging.ServiceName = serviceName
	}

	// Processing configuration overrides
	if topics := utils.GetEnv("PROCESSING_INPUT_TOPICS", ""); topics != "" {
		config.Processing.Input.Topics = parseTopics(topics)
	}
	if pollTimeout := utils.GetEnvInt("PROCESSING_INPUT_POLL_TIMEOUT_MS", -1); pollTimeout != -1 {
		config.Processing.Input.PollTimeout = time.Duration(pollTimeout) * time.Millisecond
	}
	if bufferSize := utils.GetEnvInt("PROCESSING_INPUT_BUFFER_SIZE", -1); bufferSize != -1 {
		config.Processing.Input.ChannelBufferSize = bufferSize
	}
	if delay := utils.GetEnvInt("PROCESSING_DELAY_MS", -1); delay != -1 {
		config.Processing.Processor.ProcessingDelay = time.Duration(delay) * time.Millisecond
	}
	if batchSize := utils.GetEnvInt("PROCESSING_BATCH_SIZE", -1); batchSize != -1 {
		config.Processing.Processor.BatchSize = batchSize
	}
	if outputTopic := utils.GetEnv("PROCESSING_OUTPUT_TOPIC", ""); outputTopic != "" {
		config.Processing.Output.OutputTopic = outputTopic
	}
	if outputBatchSize := utils.GetEnvInt("PROCESSING_OUTPUT_BATCH_SIZE", -1); outputBatchSize != -1 {
		config.Processing.Output.BatchSize = outputBatchSize
	}
	if flushTimeout := utils.GetEnvInt("PROCESSING_OUTPUT_FLUSH_TIMEOUT_MS", -1); flushTimeout != -1 {
		config.Processing.Output.FlushTimeout = time.Duration(flushTimeout) * time.Millisecond
	}
	if outputBufferSize := utils.GetEnvInt("PROCESSING_OUTPUT_BUFFER_SIZE", -1); outputBufferSize != -1 {
		config.Processing.Output.ChannelBufferSize = outputBufferSize
	}
	if inputBufferSize := utils.GetEnvInt("PROCESSING_CHANNELS_INPUT_BUFFER_SIZE", -1); inputBufferSize != -1 {
		config.Processing.Channels.InputBufferSize = inputBufferSize
	}
	if outputChannelBufferSize := utils.GetEnvInt("PROCESSING_CHANNELS_OUTPUT_BUFFER_SIZE", -1); outputChannelBufferSize != -1 {
		config.Processing.Channels.OutputBufferSize = outputChannelBufferSize
	}

	// Pipeline logger configuration overrides
	if ploggerLevel := utils.GetEnv("PROCESSING_PLOGGER_LEVEL", ""); ploggerLevel != "" {
		config.Processing.PloggerConfig.Level = ploggerLevel
	}
	if ploggerFileName := utils.GetEnv("PROCESSING_PLOGGER_FILE_NAME", ""); ploggerFileName != "" {
		config.Processing.PloggerConfig.FilePath = ploggerFileName
	}
	if ploggerLoggerName := utils.GetEnv("PROCESSING_PLOGGER_LOGGER_NAME", ""); ploggerLoggerName != "" {
		config.Processing.PloggerConfig.LoggerName = ploggerLoggerName
	}
	if ploggerComponentName := utils.GetEnv("PROCESSING_PLOGGER_COMPONENT_NAME", ""); ploggerComponentName != "" {
		config.Processing.PloggerConfig.ComponentName = ploggerComponentName
	}
	if ploggerServiceName := utils.GetEnv("PROCESSING_PLOGGER_SERVICE_NAME", ""); ploggerServiceName != "" {
		config.Processing.PloggerConfig.ServiceName = ploggerServiceName
	}
}

// convertLogLevel converts a string log level to logging.Level
func convertLogLevel(levelStr string) logging.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return logging.DebugLevel
	case "info":
		return logging.InfoLevel
	case "warn":
		return logging.WarnLevel
	case "error":
		return logging.ErrorLevel
	case "fatal":
		return logging.FatalLevel
	case "panic":
		return logging.PanicLevel
	default:
		return logging.InfoLevel // Default to info if invalid level
	}
}

// ConvertLoggingConfig converts LoggingConfig to logging.LoggerConfig
func (cfg LoggingConfig) ConvertToLoggerConfig() logging.LoggerConfig {
	return logging.LoggerConfig{
		Level:         convertLogLevel(cfg.Level),
		FileName:      cfg.FilePath,
		LoggerName:    cfg.LoggerName,
		ComponentName: cfg.ComponentName,
		ServiceName:   cfg.ServiceName,
	}
}
