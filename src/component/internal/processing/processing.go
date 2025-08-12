package processing

import (
	"compmodule/internal/config"
	"compmodule/internal/models"
	"fmt"
	"sharedmodule/logging"
	"time"
)

type ProcConfig struct {
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
	config        ProcConfig
	logger        logging.Logger
	inputHandler  *InputHandler
	processor     *Processor
	outputHandler *OutputHandler
	inputCh       <-chan *models.ChannelMessage
	outputCh      chan<- *models.ChannelMessage
}

func NewPipeline(config ProcConfig, logger logging.Logger) *Pipeline {
	inputHandler := NewInputHandler(config.Input, logger.WithField("component", "input"))
	outputHandler := NewOutputHandler(config.Output, logger.WithField("component", "output"))
	processor := NewProcessor(config.Processor, logger.WithField("component", "processor"), inputHandler.GetInputChannel(), outputHandler.GetOutputChannel())

	return &Pipeline{
		config:        config,
		logger:        logger,
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

func DefaultConfig(cfg *config.Config) ProcConfig {
	if cfg == nil {
		// Return hardcoded defaults if no config is provided
		return ProcConfig{
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

	// Convert config processing configuration to ProcConfig
	return ProcConfig{
		Input: InputConfig{
			Topics:            cfg.Processing.Input.Topics,
			PollTimeout:       cfg.Processing.Input.PollTimeout,
			ChannelBufferSize: cfg.Processing.Input.ChannelBufferSize,
		},
		Processor: ProcessorConfig{
			ProcessingDelay: cfg.Processing.Processor.ProcessingDelay,
			BatchSize:       cfg.Processing.Processor.BatchSize,
		},
		Output: OutputConfig{
			OutputTopic:       cfg.Processing.Output.OutputTopic,
			BatchSize:         cfg.Processing.Output.BatchSize,
			FlushTimeout:      cfg.Processing.Output.FlushTimeout,
			ChannelBufferSize: cfg.Processing.Output.ChannelBufferSize,
		},
		Channels: ChannelConfig{
			InputBufferSize:  cfg.Processing.Channels.InputBufferSize,
			OutputBufferSize: cfg.Processing.Channels.OutputBufferSize,
		},
	}
}

func ValidateConfig(config ProcConfig) error {
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
