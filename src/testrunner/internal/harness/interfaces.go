package harness

import (
	"fmt"
	"time"

	"katharos/testrunner/internal/config"
)

type TestHarness interface {
	Initialize() error
	SendMessage(data map[string]interface{}) error
	ReceiveMessage(timeout time.Duration) (map[string]interface{}, error)
	Cleanup() error
}

type LocalHarness struct {
	messageQueue chan map[string]interface{}
	config       config.MessageBusConfig
}

func NewTestHarness(cfg config.MessageBusConfig) (TestHarness, error) {
	return NewLocalHarness(cfg), nil
}

func NewLocalHarness(cfg config.MessageBusConfig) *LocalHarness {
	return &LocalHarness{
		messageQueue: make(chan map[string]interface{}, 10),
		config:       cfg,
	}
}

func (h *LocalHarness) Initialize() error {
	return nil
}

func (h *LocalHarness) SendMessage(data map[string]interface{}) error {
	select {
	case h.messageQueue <- data:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message")
	}
}

func (h *LocalHarness) ReceiveMessage(timeout time.Duration) (map[string]interface{}, error) {
	select {
	case msg := <-h.messageQueue:
		return msg, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout receiving message")
	}
}

func (h *LocalHarness) Cleanup() error {
	close(h.messageQueue)
	return nil
}
