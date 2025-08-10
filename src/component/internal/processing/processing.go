package processing

import (
	"fmt"
	"sharedmodule/logging"
	"sharedmodule/messagebus"
	"compmodule/internal/models"
	"time"
)

type Config struct {
	Input     InputConfig
	Processor ProcessorConfig
	Output    OutputConfig
	Channels  ChannelConfig
}

type ChannelConfig struct {
	InputBufferSize  int
	OutputBufferSize int
}

type Pipeline struct {
	config        Config
	logger        logging.Logger
	producer      messagebus.Producer
	inputHandler  *InputHandler
	processor     *Processor
	outputHandler *OutputHandler
	inputCh       <-chan *models.ChannelMessage
	outputCh      chan<- *models.ChannelMessage
}

func NewPipeline(config Config, logger logging.Logger) *Pipeline {
	producer := messagebus.NewProducer()
	
	inputHandler := NewInputHandler(producer, config.Input, logger.WithField("component", "input"))
	outputHandler := NewOutputHandler(producer, config.Output, logger.WithField("component", "output"))
	processor := NewProcessor(config.Processor, logger.WithField("component", "processor"), inputHandler.GetInputChannel(), outputHandler.GetOutputChannel())

	return &Pipeline{
		config:        config,
		logger:        logger,
		producer:      producer,
		inputHandler:  inputHandler,
		processor:     processor,
		outputHandler: outputHandler,
		inputCh:       inputHandler.GetInputChannel(),
		outputCh:      outputHandler.GetOutputChannel(),
	}
}

func (p *Pipeline) Start() error {
	p.logger.Info("Starting processing pipeline")

	if err := p.outputHandler.Start(); err != nil {
		return fmt.Errorf("failed to start output handler: %w", err)
	}

	if err := p.processor.Start(); err != nil {
		p.outputHandler.Stop()
		return fmt.Errorf("failed to start processor: %w", err)
	}

	if err := p.inputHandler.Start(); err != nil {
		p.processor.Stop()
		p.outputHandler.Stop()
		return fmt.Errorf("failed to start input handler: %w", err)
	}

	p.logger.Info("Processing pipeline started successfully")
	return nil
}

func (p *Pipeline) Stop() error {
	p.logger.Info("Stopping processing pipeline")

	var errs []error

	if err := p.inputHandler.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("error stopping input handler: %w", err))
	}

	if err := p.processor.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("error stopping processor: %w", err))
	}

	if err := p.outputHandler.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("error stopping output handler: %w", err))
	}

	if len(errs) > 0 {
		p.logger.Errorw("Errors occurred during pipeline shutdown", "error_count", len(errs))
		return fmt.Errorf("pipeline shutdown errors: %v", errs)
	}

	p.logger.Info("Processing pipeline stopped successfully")
	return nil
}

func (p *Pipeline) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"pipeline_status": "running",
		"input_stats":     p.inputHandler.GetStats(),
		"processor_stats": p.processor.GetStats(),
		"output_stats":    p.outputHandler.GetStats(),
	}
}

func DefaultConfig() Config {
	return Config{
		Input: InputConfig{
			Topics:            []string{"input-topic"},
			PollTimeout:       1 * time.Second,
			ChannelBufferSize: 1000,
		},
		Processor: ProcessorConfig{
			ProcessingDelay: 10 * time.Millisecond,
			BatchSize:       100,
		},
		Output: OutputConfig{
			OutputTopic:       "output-topic",
			BatchSize:         50,
			FlushTimeout:      5 * time.Second,
			ChannelBufferSize: 1000,
		},
		Channels: ChannelConfig{
			InputBufferSize:  1000,
			OutputBufferSize: 1000,
		},
	}
}

func ValidateConfig(config Config) error {
	if len(config.Input.Topics) == 0 {
		return fmt.Errorf("input topics cannot be empty")
	}
	if config.Input.PollTimeout <= 0 {
		return fmt.Errorf("poll timeout must be positive")
	}
	if config.Input.ChannelBufferSize <= 0 {
		return fmt.Errorf("input channel buffer size must be positive")
	}

	if config.Processor.BatchSize <= 0 {
		return fmt.Errorf("processor batch size must be positive")
	}

	if config.Output.OutputTopic == "" {
		return fmt.Errorf("output topic cannot be empty")
	}
	if config.Output.BatchSize <= 0 {
		return fmt.Errorf("output batch size must be positive")
	}
	if config.Output.FlushTimeout <= 0 {
		return fmt.Errorf("flush timeout must be positive")
	}
	if config.Output.ChannelBufferSize <= 0 {
		return fmt.Errorf("output channel buffer size must be positive")
	}

	if config.Channels.InputBufferSize <= 0 {
		return fmt.Errorf("input buffer size must be positive")
	}
	if config.Channels.OutputBufferSize <= 0 {
		return fmt.Errorf("output buffer size must be positive")
	}

	return nil
}
